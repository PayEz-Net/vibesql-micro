package postgres

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name        string
		dataDir     string
		port        int
		wantDataDir string
		wantPort    int
	}{
		{
			name:        "default values",
			dataDir:     "",
			port:        0,
			wantDataDir: defaultDataDir,
			wantPort:    defaultPort,
		},
		{
			name:        "custom values",
			dataDir:     "/tmp/test-data",
			port:        5433,
			wantDataDir: "/tmp/test-data",
			wantPort:    5433,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.dataDir, tt.port)
			if m == nil {
				t.Fatal("NewManager returned nil")
			}
			if m.dataDir != tt.wantDataDir {
				t.Errorf("dataDir = %s, want %s", m.dataDir, tt.wantDataDir)
			}
			if m.port != tt.wantPort {
				t.Errorf("port = %d, want %d", m.port, tt.wantPort)
			}
			if m.running {
				t.Error("manager should not be running initially")
			}
			if m.ctx == nil {
				t.Error("context should be initialized")
			}
			if m.cancel == nil {
				t.Error("cancel function should be initialized")
			}
		})
	}
}

func TestManager_ExtractBinaries(t *testing.T) {
	m := NewManager("", 0)
	
	err := m.extractBinaries()
	if err != nil {
		t.Skipf("Skipping binary extraction test: %v (expected on platforms without embedded binary)", err)
		return
	}
	
	// Verify binary was extracted
	if m.postgresBinPath == "" {
		t.Error("postgresBinPath should be set")
	}
	
	// Verify file exists and is executable
	info, err := os.Stat(m.postgresBinPath)
	if err != nil {
		t.Errorf("postgres binary not found: %v", err)
	}
	
	if info.IsDir() {
		t.Error("postgres binary path points to directory")
	}
	
	// Check permissions (should be executable)
	mode := info.Mode()
	if mode&0100 == 0 {
		t.Error("postgres binary is not executable")
	}
}

func TestManager_InitializeDataDir(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "vibe-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	testDataDir := filepath.Join(tmpDir, "pgdata")
	m := NewManager(testDataDir, 0)
	
	// Extract binaries first
	err = m.extractBinaries()
	if err != nil {
		t.Skipf("Skipping data dir initialization test: %v", err)
		return
	}
	
	// Initialize data directory
	err = m.initializeDataDir()
	if err != nil {
		t.Fatalf("initializeDataDir failed: %v", err)
	}
	
	// Verify data directory structure
	expectedDirs := []string{
		"base",
		"global",
		"pg_wal",
		"pg_stat",
	}
	
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(testDataDir, dir)
		info, err := os.Stat(dirPath)
		if err != nil {
			t.Errorf("expected directory %s not found: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", dir)
		}
	}
	
	// Verify configuration files
	expectedFiles := []string{
		"postgresql.conf",
		"pg_hba.conf",
	}
	
	for _, file := range expectedFiles {
		filePath := filepath.Join(testDataDir, file)
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("expected file %s not found: %v", file, err)
		}
	}
	
	// Test idempotency - calling again should not error
	err = m.initializeDataDir()
	if err != nil {
		t.Errorf("second initializeDataDir call failed: %v", err)
	}
}

func TestManager_InitializeDataDir_AlreadyInitialized(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "vibe-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	testDataDir := filepath.Join(tmpDir, "pgdata")
	
	// Create data directory and PG_VERSION file
	err = os.MkdirAll(testDataDir, 0700)
	if err != nil {
		t.Fatalf("failed to create data dir: %v", err)
	}
	
	pgVersionPath := filepath.Join(testDataDir, "PG_VERSION")
	err = os.WriteFile(pgVersionPath, []byte("16\n"), 0600)
	if err != nil {
		t.Fatalf("failed to write PG_VERSION: %v", err)
	}
	
	m := NewManager(testDataDir, 0)
	
	// Extract binaries
	err = m.extractBinaries()
	if err != nil {
		t.Skipf("Skipping test: %v", err)
		return
	}
	
	// Should skip initialization
	err = m.initializeDataDir()
	if err != nil {
		t.Errorf("initializeDataDir failed on already initialized directory: %v", err)
	}
}

func TestManager_CreateConfigFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibe-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	m := NewManager(tmpDir, 5433)
	
	err = m.createConfigFiles()
	if err != nil {
		t.Fatalf("createConfigFiles failed: %v", err)
	}
	
	// Check postgresql.conf
	confPath := filepath.Join(tmpDir, "postgresql.conf")
	confData, err := os.ReadFile(confPath)
	if err != nil {
		t.Fatalf("failed to read postgresql.conf: %v", err)
	}
	
	confStr := string(confData)
	expectedSettings := []string{
		"listen_addresses = '127.0.0.1'",
		"port = 5433",
		"max_connections = 10",
	}
	
	for _, setting := range expectedSettings {
		if !strings.Contains(confStr, setting) {
			t.Errorf("postgresql.conf missing expected setting: %s", setting)
		}
	}
	
	// Check pg_hba.conf
	hbaPath := filepath.Join(tmpDir, "pg_hba.conf")
	hbaData, err := os.ReadFile(hbaPath)
	if err != nil {
		t.Fatalf("failed to read pg_hba.conf: %v", err)
	}
	
	hbaStr := string(hbaData)
	if !strings.Contains(hbaStr, "127.0.0.1/32") {
		t.Error("pg_hba.conf missing localhost entry")
	}
}

func TestManager_GetConnectionString(t *testing.T) {
	tests := []struct {
		name string
		port int
		want string
	}{
		{
			name: "default port",
			port: 5432,
			want: "host=127.0.0.1 port=5432 dbname=postgres user=postgres sslmode=disable",
		},
		{
			name: "custom port",
			port: 5433,
			want: "host=127.0.0.1 port=5433 dbname=postgres user=postgres sslmode=disable",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager("", tt.port)
			got := m.GetConnectionString()
			if got != tt.want {
				t.Errorf("GetConnectionString() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestManager_GetDataDir(t *testing.T) {
	dataDir := "/tmp/test-data"
	m := NewManager(dataDir, 0)
	
	got := m.GetDataDir()
	if got != dataDir {
		t.Errorf("GetDataDir() = %s, want %s", got, dataDir)
	}
}

func TestManager_IsRunning(t *testing.T) {
	m := NewManager("", 0)
	
	if m.IsRunning() {
		t.Error("IsRunning() should return false initially")
	}
	
	// Simulate running state
	m.processLock.Lock()
	m.running = true
	m.processLock.Unlock()
	
	if !m.IsRunning() {
		t.Error("IsRunning() should return true when running")
	}
	
	// Reset state
	m.processLock.Lock()
	m.running = false
	m.processLock.Unlock()
	
	if m.IsRunning() {
		t.Error("IsRunning() should return false after reset")
	}
}

func TestManager_Start_AlreadyRunning(t *testing.T) {
	m := NewManager("", 0)
	
	// Set running state
	m.processLock.Lock()
	m.running = true
	m.processLock.Unlock()
	
	err := m.Start()
	if err == nil {
		t.Error("Start() should return error when already running")
	}
	
	if err.Error() != "postgres manager already running" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestManager_Stop_NotRunning(t *testing.T) {
	m := NewManager("", 0)
	
	// Stop when not running should not error
	err := m.Stop()
	if err != nil {
		t.Errorf("Stop() returned unexpected error: %v", err)
	}
}

func TestManager_Stop_CancelsContext(t *testing.T) {
	m := NewManager("", 0)
	
	// Simulate running state
	m.processLock.Lock()
	m.running = true
	m.processLock.Unlock()
	
	// Stop should cancel context
	err := m.Stop()
	if err != nil {
		t.Errorf("Stop() returned error: %v", err)
	}
	
	// Check context is cancelled
	select {
	case <-m.ctx.Done():
		// Context cancelled as expected
	case <-time.After(100 * time.Millisecond):
		t.Error("context was not cancelled")
	}
	
	// Check running state
	if m.IsRunning() {
		t.Error("manager should not be running after Stop()")
	}
}

// Integration test - only runs if embedded binary is available
func TestManager_StartStop_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	
	tmpDir, err := os.MkdirTemp("", "vibe-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	testDataDir := filepath.Join(tmpDir, "pgdata")
	m := NewManager(testDataDir, 5433)
	
	// Start PostgreSQL
	err = m.Start()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}
	
	// Verify running
	if !m.IsRunning() {
		t.Error("manager should be running after Start()")
	}
	
	// Wait a bit to ensure it's stable
	time.Sleep(500 * time.Millisecond)
	
	// Stop PostgreSQL
	err = m.Stop()
	if err != nil {
		t.Errorf("Stop() failed: %v", err)
	}
	
	// Verify stopped
	if m.IsRunning() {
		t.Error("manager should not be running after Stop()")
	}
}

func TestManager_GetPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"default port", 5432},
		{"custom port", 5433},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager("", tt.port)
			got := m.GetPort()
			if got != tt.port {
				t.Errorf("GetPort() = %d, want %d", got, tt.port)
			}
		})
	}
}

func TestManager_CreateConnection_NotRunning(t *testing.T) {
	m := NewManager("", 0)

	_, err := m.CreateConnection()
	if err == nil {
		t.Error("CreateConnection() should return error when not running")
	}

	if err.Error() != "postgres manager is not running" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func BenchmarkNewManager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewManager("", 0)
	}
}

func BenchmarkManager_IsRunning(b *testing.B) {
	m := NewManager("", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.IsRunning()
	}
}


