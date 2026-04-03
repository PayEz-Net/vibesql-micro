package postgres

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
	
	_ "github.com/lib/pq"
)

// Process manages a postgres subprocess
type Process struct {
	cmd      *exec.Cmd
	dataDir  string
	port     int
	db       *sql.DB
	shutdown chan struct{}
}

// findFreePort finds an available TCP port
func findFreePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// Start initializes and starts postgres
func Start(postgresBin, initdbBin, shareDir, dataDir string) (*Process, error) {
	// Find a free port
	port, err := findFreePort()
	if err != nil {
		return nil, fmt.Errorf("find free port: %w", err)
	}
	
	// Check if already initialized
	pgVersionFile := filepath.Join(dataDir, "PG_VERSION")
	if _, err := os.Stat(pgVersionFile); os.IsNotExist(err) {
		// Need to run initdb
		if err := runInitdb(initdbBin, shareDir, dataDir); err != nil {
			return nil, fmt.Errorf("initdb failed: %w", err)
		}
	}
	
	// Start postgres
	p := &Process{
		dataDir:  dataDir,
		port:     port,
		shutdown: make(chan struct{}),
	}
	
	if err := p.start(postgresBin, shareDir); err != nil {
		return nil, err
	}
	
	// Connect to the database
	connStr := fmt.Sprintf("host=localhost port=%d user=postgres dbname=postgres sslmode=disable", port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		p.Stop()
		return nil, fmt.Errorf("open database: %w", err)
	}
	
	// Wait for connection to be ready
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	
	if err := db.Ping(); err != nil {
		p.Stop()
		return nil, fmt.Errorf("database not ready: %w", err)
	}
	
	p.db = db
	return p, nil
}

func runInitdb(initdbBin, shareDir, dataDir string) error {
	args := []string{
		"-D", dataDir,
		"--no-locale",
		"--encoding=UTF8",
		"--auth=trust",
		"--username=postgres",
	}
	
	binDir := filepath.Dir(initdbBin)
	cmd := exec.Command(initdbBin, args...)
	cmd.Dir = binDir
	
	// Set up environment
	env := os.Environ()
	env = append(env, "PGSHAREDIR="+shareDir)
	env = append(env, "LD_LIBRARY_PATH="+binDir)
	
	// Prepend binary directory to PATH for Windows DLL loading
	for i, e := range env {
		if len(e) > 5 && e[:5] == "PATH=" {
			if runtime.GOOS == "windows" {
				env[i] = "PATH=" + binDir + ";" + e[5:]
			} else {
				env[i] = "PATH=" + binDir + ":" + e[5:]
			}
			break
		}
	}
	cmd.Env = env
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("initdb: %w\nOutput: %s", err, string(output))
	}
	
	return nil
}

func (p *Process) start(postgresBin, shareDir string) error {
	args := []string{
		"-D", p.dataDir,
		"-p", fmt.Sprintf("%d", p.port),
		"-h", "127.0.0.1",
	}
	
	// Disable Unix sockets on Windows
	if runtime.GOOS == "windows" {
		args = append(args, "-k", "")
	}
	
	binDir := filepath.Dir(postgresBin)
	p.cmd = exec.Command(postgresBin, args...)
	p.cmd.Dir = binDir
	
	// Set up environment
	env := os.Environ()
	env = append(env, "PGSHAREDIR="+shareDir)
	env = append(env, "LD_LIBRARY_PATH="+binDir)
	
	// Prepend binary directory to PATH so Windows can find DLLs
	for i, e := range env {
		if len(e) > 5 && e[:5] == "PATH=" {
			env[i] = "PATH=" + binDir + ";" + e[5:]
			break
		}
	}
	p.cmd.Env = env
	
	// Capture output for debugging
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}
	
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("start postgres: %w", err)
	}
	
	// Log output in background
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			// Discard or log
		}
	}()
	
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			// Discard or log
		}
	}()
	
	return nil
}

// Stop shuts down postgres gracefully
func (p *Process) Stop() error {
	if p.db != nil {
		p.db.Close()
	}
	
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	
	// Try graceful shutdown first
	p.cmd.Process.Signal(os.Interrupt)
	
	// Wait a bit
	done := make(chan error, 1)
	go func() {
		done <- p.cmd.Wait()
	}()
	
	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
		// Force kill
		p.cmd.Process.Kill()
		return nil
	}
}

// Port returns the postgres port
func (p *Process) Port() int {
	return p.port
}

// DB returns the database connection
func (p *Process) DB() *sql.DB {
	return p.db
}
