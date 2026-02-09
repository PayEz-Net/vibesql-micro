package postgres

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

//go:embed embed/*
var embeddedPostgres embed.FS

const (
	defaultDataDir  = "./vibe-data"
	defaultPort     = 5432
	shutdownTimeout = 10 * time.Second
)

var (
	startupTimeout = 30 * time.Second
)

type Manager struct {
	dataDir     string
	port        int
	process     *exec.Cmd
	processLock sync.Mutex
	running     bool

	postgresBinPath string
	initdbBinPath   string
	pgCtlBinPath    string
	libDir          string
	shareDir        string
	tmpDir          string

	// Windows workaround: EDB binaries have hardcoded /share and $libdir paths
	// which Windows interprets as <drive>:\share and <drive>:\lib
	winShareDir string
	winLibDir   string

	ctx    context.Context
	cancel context.CancelFunc
	errCh  chan error
}

func NewManager(dataDir string, port int) *Manager {
	if dataDir == "" {
		dataDir = defaultDataDir
	}
	if port == 0 {
		port = defaultPort
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		dataDir: dataDir,
		port:    port,
		ctx:     ctx,
		cancel:  cancel,
		errCh:   make(chan error, 1),
	}
}

func (m *Manager) Start() error {
	m.processLock.Lock()
	defer m.processLock.Unlock()

	if m.running {
		return fmt.Errorf("postgres manager already running")
	}

	if err := m.extractBinaries(); err != nil {
		return fmt.Errorf("failed to extract postgres binaries: %w", err)
	}

	if err := m.initializeDataDir(); err != nil {
		return fmt.Errorf("failed to initialize data directory: %w", err)
	}

	if err := m.startPostgres(); err != nil {
		return fmt.Errorf("failed to start postgres: %w", err)
	}

	if err := m.waitForReady(); err != nil {
		_ = m.stopPostgres()
		return fmt.Errorf("postgres failed to become ready: %w", err)
	}

	m.running = true

	go m.monitorProcess()

	return nil
}

func (m *Manager) Stop() error {
	m.processLock.Lock()
	defer m.processLock.Unlock()

	if !m.running {
		return nil
	}

	m.cancel()
	m.running = false

	err := m.stopPostgres()

	if m.tmpDir != "" {
		_ = os.RemoveAll(m.tmpDir)
		m.tmpDir = ""
	}

	// Clean up Windows workaround directories
	if m.winShareDir != "" {
		_ = os.RemoveAll(m.winShareDir)
		m.winShareDir = ""
	}
	if m.winLibDir != "" {
		_ = os.RemoveAll(m.winLibDir)
		m.winLibDir = ""
	}

	return err
}

func platformBinExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func libpqName() string {
	switch runtime.GOOS {
	case "darwin":
		return "libpq.5.dylib"
	case "windows":
		return "libpq-5.dll"
	default:
		return "libpq.so.5"
	}
}

func libPathEnvVar() string {
	switch runtime.GOOS {
	case "darwin":
		return "DYLD_LIBRARY_PATH"
	case "windows":
		return "PATH"
	default:
		return "LD_LIBRARY_PATH"
	}
}

func supportedPlatform() bool {
	switch runtime.GOOS {
	case "linux", "darwin":
		switch runtime.GOARCH {
		case "amd64", "arm64":
			return true
		}
	case "windows":
		if runtime.GOARCH == "amd64" {
			return true
		}
	}
	return false
}

func (m *Manager) extractBinaries() error {
	// Check for system PostgreSQL via environment variable
	if postgresBin := os.Getenv("POSTGRES_BIN"); postgresBin != "" {
		log.Printf("[INFO] Using system PostgreSQL from POSTGRES_BIN: %s", postgresBin)
		m.postgresBinPath = postgresBin
		m.initdbBinPath = filepath.Join(filepath.Dir(postgresBin), "initdb"+platformBinExt())
		m.pgCtlBinPath = filepath.Join(filepath.Dir(postgresBin), "pg_ctl"+platformBinExt())

		// Check if required binaries exist
		if _, err := os.Stat(m.postgresBinPath); err != nil {
			return fmt.Errorf("POSTGRES_BIN specified but postgres not found at %s: %w", m.postgresBinPath, err)
		}
		if _, err := os.Stat(m.initdbBinPath); err != nil {
			return fmt.Errorf("POSTGRES_BIN specified but initdb not found at %s: %w", m.initdbBinPath, err)
		}

		// For system PostgreSQL, use system share directory
		if shareDir := os.Getenv("PGSHAREDIR"); shareDir != "" {
			m.shareDir = shareDir
		}

		return nil
	}

	if !supportedPlatform() {
		return fmt.Errorf(
			"unsupported platform: %s/%s\n\n"+
				"VibeSQL supports: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64\n"+
				"Build PostgreSQL manually and set POSTGRES_BIN environment variable",
			runtime.GOOS, runtime.GOARCH)
	}

	platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	ext := platformBinExt()

	tmpDir, err := os.MkdirTemp("", "vibe-postgres-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	m.tmpDir = tmpDir

	postgresEmbedPath := fmt.Sprintf("embed/postgres_micro_%s%s", platform, ext)
	postgresData, err := embeddedPostgres.ReadFile(postgresEmbedPath)
	if err != nil {
		return fmt.Errorf("embedded postgres binary not found for platform %s: %w", platform, err)
	}
	m.postgresBinPath = filepath.Join(tmpDir, "postgres"+ext)
	if err := os.WriteFile(m.postgresBinPath, postgresData, 0755); err != nil {
		return fmt.Errorf("failed to write postgres binary: %w", err)
	}

	initdbEmbedPath := fmt.Sprintf("embed/initdb_%s%s", platform, ext)
	initdbData, err := embeddedPostgres.ReadFile(initdbEmbedPath)
	if err != nil {
		return fmt.Errorf("embedded initdb binary not found for platform %s: %w", platform, err)
	}
	m.initdbBinPath = filepath.Join(tmpDir, "initdb"+ext)
	if err := os.WriteFile(m.initdbBinPath, initdbData, 0755); err != nil {
		return fmt.Errorf("failed to write initdb binary: %w", err)
	}

	pgCtlEmbedPath := fmt.Sprintf("embed/pg_ctl_%s%s", platform, ext)
	pgCtlData, err := embeddedPostgres.ReadFile(pgCtlEmbedPath)
	if err == nil {
		m.pgCtlBinPath = filepath.Join(tmpDir, "pg_ctl"+ext)
		if writeErr := os.WriteFile(m.pgCtlBinPath, pgCtlData, 0755); writeErr != nil {
			m.pgCtlBinPath = ""
		}
	}

	libDir := filepath.Join(tmpDir, "lib")
	if err := os.MkdirAll(libDir, 0755); err != nil {
		return fmt.Errorf("failed to create lib directory: %w", err)
	}

	libName := libpqName()
	libpqData, err := embeddedPostgres.ReadFile("embed/" + libName)
	if err != nil {
		return fmt.Errorf("embedded %s not found: %w", libName, err)
	}
	libpqPath := filepath.Join(libDir, libName)
	if err := os.WriteFile(libpqPath, libpqData, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", libName, err)
	}

	if runtime.GOOS == "windows" {
		// Copy libpq to tmpDir for Windows (both names needed)
		_ = os.WriteFile(filepath.Join(tmpDir, libName), libpqData, 0644)
		_ = os.WriteFile(filepath.Join(tmpDir, "LIBPQ.dll"), libpqData, 0644)

		// Extract all Windows DLLs needed by PostgreSQL binaries
		windowsDLLs := []string{
			"libcrypto-3-x64.dll",
			"libssl-3-x64.dll",
			"libiconv-2.dll",
			"libintl-9.dll",
			"zlib1.dll",
			"icudt67.dll",
			"icuin67.dll",
			"icuio67.dll",
			"icutu67.dll",
			"icuuc67.dll",
			"libwinpthread-1.dll",
			"libzstd.dll",
			"liblz4.dll",
			"libxml2.dll",
		}
		for _, dllName := range windowsDLLs {
			dllData, dllErr := embeddedPostgres.ReadFile("embed/" + dllName)
			if dllErr == nil {
				_ = os.WriteFile(filepath.Join(tmpDir, dllName), dllData, 0644)
			}
		}

		// Extract PostgreSQL extension DLLs to lib directory (for $libdir)
		libExtDLLs := []string{
			"plpgsql.dll",
			"dict_snowball.dll",
		}
		for _, dllName := range libExtDLLs {
			dllData, dllErr := embeddedPostgres.ReadFile("embed/" + dllName)
			if dllErr == nil {
				_ = os.WriteFile(filepath.Join(libDir, dllName), dllData, 0644)
			}
		}
	}

	m.libDir = libDir

	shareTarData, err := embeddedPostgres.ReadFile("embed/share.tar.gz")
	if err != nil {
		return fmt.Errorf("embedded share.tar.gz not found: %w", err)
	}

	if err := extractShareTarGz(shareTarData, tmpDir); err != nil {
		return fmt.Errorf("failed to extract share directory: %w", err)
	}

	m.shareDir = filepath.Join(tmpDir, "share")

	// Windows workaround: EDB binaries have hardcoded /share and $libdir paths
	// which Windows interprets as <drive>:\share and <drive>:\lib
	// We create these directories at the drive root and clean them up on Stop()
	// IMPORTANT: Use the CURRENT WORKING DIRECTORY's drive, not tmpDir's drive,
	// because that's what postgres.exe will use when resolving /share
	if runtime.GOOS == "windows" {
		cwd, _ := os.Getwd()
		driveLetter := filepath.VolumeName(cwd)
		if driveLetter == "" {
			driveLetter = filepath.VolumeName(tmpDir)
		}
		if driveLetter != "" {
			// Create <drive>:\share by copying our extracted share
			m.winShareDir = filepath.Join(driveLetter, "\\share")
			if err := copyDir(m.shareDir, m.winShareDir); err != nil {
				return fmt.Errorf("failed to create Windows share directory: %w", err)
			}

			// Create <drive>:\lib with extension DLLs
			m.winLibDir = filepath.Join(driveLetter, "\\lib")
			if err := os.MkdirAll(m.winLibDir, 0755); err != nil {
				return fmt.Errorf("failed to create Windows lib directory: %w", err)
			}
			// Copy extension DLLs to drive root lib
			libExtDLLs := []string{"plpgsql.dll", "dict_snowball.dll"}
			for _, dllName := range libExtDLLs {
				srcPath := filepath.Join(libDir, dllName)
				if _, err := os.Stat(srcPath); err == nil {
					dstPath := filepath.Join(m.winLibDir, dllName)
					data, _ := os.ReadFile(srcPath)
					_ = os.WriteFile(dstPath, data, 0644)
				}
			}
		}
	}

	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}

func extractShareTarGz(data []byte, targetDir string) error {
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar read error: %w", err)
		}

		targetPath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			outFile.Close()
		}
	}

	return nil
}

func (m *Manager) initializeDataDir() error {
	pgVersionPath := filepath.Join(m.dataDir, "PG_VERSION")
	if _, err := os.Stat(pgVersionPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(m.dataDir, 0700); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	initdbArgs := []string{
		"-D", m.dataDir,
		"--no-locale",
		"--encoding=UTF8",
		"--auth=trust",
		"--username=postgres",
		"--nosync",
	}

	if m.shareDir != "" {
		initdbArgs = append(initdbArgs, "-L", m.shareDir)
	}

	initdbArgs = append(initdbArgs, "-c", "timezone=+00", "-c", "log_timezone=+00")

	cmd := exec.Command(m.initdbBinPath, initdbArgs...)
	cmd.Env = m.buildEnv()

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if initdb partially succeeded (PG_VERSION exists) - can happen on macOS
		// when plpgsql extension fails to load due to .dylib symbol issues
		if _, statErr := os.Stat(pgVersionPath); statErr == nil {
			log.Printf("[WARN] initdb reported error but data directory was created, continuing...")
			log.Printf("[WARN] initdb output: %s", string(output))
		} else {
			return fmt.Errorf("initdb failed: %w\nOutput: %s", err, string(output))
		}
	}

	if err := m.createConfigFiles(); err != nil {
		return fmt.Errorf("failed to create config files: %w", err)
	}

	return nil
}

func (m *Manager) buildEnv() []string {
	env := os.Environ()
	if m.shareDir != "" {
		env = append(env, "PGSHAREDIR="+m.shareDir)
	}
	if m.libDir != "" {
		env = append(env, "PKGLIBDIR="+m.libDir)
	}
	if m.libDir == "" {
		return env
	}

	if runtime.GOOS == "windows" {
		existing := os.Getenv("PATH")
		env = append(env, "PATH="+m.libDir+";"+m.tmpDir+";"+existing)
	} else if runtime.GOOS == "darwin" {
		// On macOS, extensions need to find symbols from the postgres binary
		// DYLD_LIBRARY_PATH needs to include both lib dir and the dir with postgres binary
		env = append(env, libPathEnvVar()+"="+m.libDir+":"+m.tmpDir)
	} else {
		env = append(env, libPathEnvVar()+"="+m.libDir)
	}
	return env
}

func (m *Manager) createConfigFiles() error {
	confPath := filepath.Join(m.dataDir, "postgresql.conf")
	shmType := "posix"
	if runtime.GOOS == "windows" {
		shmType = "windows"
	}
	conf := fmt.Sprintf(`
listen_addresses = '127.0.0.1'
port = %d
max_connections = 10
shared_buffers = 12MB
dynamic_shared_memory_type = %s
max_wal_size = 100MB
min_wal_size = 80MB
log_destination = 'stderr'
logging_collector = off
log_statement = 'all'
`, m.port, shmType)

	if err := os.WriteFile(confPath, []byte(conf), 0600); err != nil {
		return err
	}

	hbaPath := filepath.Join(m.dataDir, "pg_hba.conf")
	var hba string
	if runtime.GOOS == "windows" {
		hba = `# TYPE  DATABASE        USER            ADDRESS                 METHOD
host    all             all             127.0.0.1/32            trust
host    all             all             ::1/128                 trust
`
	} else {
		hba = `# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     trust
host    all             all             127.0.0.1/32            trust
host    all             all             ::1/128                 trust
`
	}
	return os.WriteFile(hbaPath, []byte(hba), 0600)
}

func (m *Manager) startPostgres() error {
	args := []string{
		"-D", m.dataDir,
		"-c", fmt.Sprintf("port=%d", m.port),
		"-c", "listen_addresses=127.0.0.1",
		"-c", "max_connections=10",
		"-c", "shared_buffers=12MB",
	}

	m.process = exec.CommandContext(m.ctx, m.postgresBinPath, args...)
	m.process.Env = m.buildEnv()

	stdout, err := m.process.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := m.process.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := m.process.Start(); err != nil {
		return fmt.Errorf("failed to start postgres process: %w", err)
	}

	go m.logOutput(stdout, "stdout")
	go m.logOutput(stderr, "stderr")

	return nil
}

func (m *Manager) stopPostgres() error {
	if m.process == nil {
		return nil
	}

	if m.pgCtlBinPath != "" {
		cmd := exec.Command(m.pgCtlBinPath, "stop", "-D", m.dataDir, "-m", "fast", "-w")
		cmd.Env = m.buildEnv()
		if err := cmd.Run(); err == nil {
			m.process = nil
			return nil
		}
	}

	if m.process.Process != nil {
		if runtime.GOOS == "windows" {
			_ = m.process.Process.Kill()
		} else {
			if err := m.process.Process.Signal(os.Interrupt); err != nil {
				_ = m.process.Process.Kill()
			}
		}

		timer := time.NewTimer(shutdownTimeout)
		defer timer.Stop()
		done := make(chan struct{})
		go func() {
			_ = m.process.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-timer.C:
			if m.process.Process != nil {
				_ = m.process.Process.Kill()
			}
		}
	}

	m.process = nil
	return nil
}

func (m *Manager) waitForReady() error {
	timeout := time.After(startupTimeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("postgres startup timeout after %v", startupTimeout)
		case err := <-m.errCh:
			return fmt.Errorf("postgres startup error: %w", err)
		case <-ticker.C:
			if m.isReady() {
				return nil
			}
		}
	}
}

func (m *Manager) isReady() bool {
	pidPath := filepath.Join(m.dataDir, "postmaster.pid")
	if _, err := os.Stat(pidPath); err != nil {
		return false
	}

	if m.process == nil || m.process.Process == nil {
		return false
	}

	select {
	case err := <-m.errCh:
		m.errCh <- err
		return false
	default:
		conn, err := NewConnectionSimple(m.port)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}
}

func (m *Manager) monitorProcess() {
	if m.process == nil {
		return
	}

	err := m.process.Wait()

	select {
	case <-m.ctx.Done():
		return
	default:
		if err != nil {
			m.errCh <- fmt.Errorf("postgres process exited unexpectedly: %w", err)
		} else {
			m.errCh <- fmt.Errorf("postgres process exited unexpectedly")
		}

		m.processLock.Lock()
		m.running = false
		m.processLock.Unlock()
	}
}

func (m *Manager) logOutput(reader io.Reader, source string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "database system is ready to accept connections") {
			continue
		}

		if strings.Contains(line, "FATAL") || strings.Contains(line, "ERROR") {
			fmt.Fprintf(os.Stderr, "[postgres %s] %s\n", source, line)
		}
	}
}

func (m *Manager) IsRunning() bool {
	m.processLock.Lock()
	defer m.processLock.Unlock()
	return m.running
}

func (m *Manager) GetConnectionString() string {
	return fmt.Sprintf("host=127.0.0.1 port=%d dbname=postgres user=postgres sslmode=disable", m.port)
}

func (m *Manager) GetDataDir() string {
	return m.dataDir
}

func (m *Manager) CreateConnection() (*Connection, error) {
	if !m.running {
		return nil, fmt.Errorf("postgres manager is not running")
	}

	return NewConnectionSimple(m.port)
}

func (m *Manager) GetPort() int {
	return m.port
}
