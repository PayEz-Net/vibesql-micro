package postgres

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestManager_InitializeDataDir_InvalidBinary(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManager(tempDir, 5433)
	
	m.initdbBinPath = "/nonexistent/postgres"
	
	err := m.initializeDataDir()
	if err == nil {
		t.Error("expected error when initdb binary doesn't exist")
	}
}

func TestManager_SupportedPlatform(t *testing.T) {
	supported := supportedPlatform()
	if !supported {
		t.Skipf("current platform not supported, skipping assertion")
	}
}

func TestManager_PlatformBinExt(t *testing.T) {
	ext := platformBinExt()
	if ext != "" && ext != ".exe" {
		t.Errorf("unexpected binary extension: %s", ext)
	}
}

func TestManager_LibpqName(t *testing.T) {
	name := libpqName()
	if name == "" {
		t.Error("libpqName returned empty string")
	}
	validNames := []string{"libpq.so.5", "libpq.5.dylib", "libpq-5.dll"}
	found := false
	for _, v := range validNames {
		if name == v {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("unexpected libpq name: %s", name)
	}
}

func TestManager_LibPathEnvVar(t *testing.T) {
	envVar := libPathEnvVar()
	validVars := []string{"LD_LIBRARY_PATH", "DYLD_LIBRARY_PATH", "PATH"}
	found := false
	for _, v := range validVars {
		if envVar == v {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("unexpected lib path env var: %s", envVar)
	}
}

func TestManager_StartPostgres_InvalidBinary(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManager(tempDir, 5433)
	m.postgresBinPath = "/nonexistent/postgres"
	
	err := m.startPostgres()
	if err == nil {
		t.Fatal("expected error for invalid binary path, got nil")
	}
}

func TestManager_StopPostgres_NoProcess(t *testing.T) {
	m := NewManager("", 0)
	
	err := m.stopPostgres()
	if err != nil {
		t.Errorf("stopPostgres with no process should not error, got: %v", err)
	}
}

func TestManager_StopPostgres_WithProcess(t *testing.T) {
	m := NewManager("", 0)
	
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Skipf("cannot start sleep command: %v", err)
	}
	
	m.process = cmd
	
	err := m.stopPostgres()
	if err != nil {
		t.Errorf("stopPostgres failed: %v", err)
	}
	
	if m.process != nil {
		t.Error("process should be nil after stop")
	}
}

func TestManager_IsReady_NotInitialized(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManager(tempDir, 5433)
	
	ready := m.isReady()
	if ready {
		t.Error("expected isReady to return false for uninitialized manager")
	}
}

func TestManager_IsReady_NoProcess(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManager(tempDir, 5433)
	
	pidPath := filepath.Join(tempDir, "postmaster.pid")
	if err := os.WriteFile(pidPath, []byte("12345\n"), 0600); err != nil {
		t.Fatalf("failed to create postmaster.pid: %v", err)
	}
	
	ready := m.isReady()
	if ready {
		t.Error("expected isReady to return false when process is nil")
	}
}

func TestManager_WaitForReady_Timeout(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManager(tempDir, 5433)
	
	originalTimeout := startupTimeout
	defer func() { startupTimeout = originalTimeout }()
	
	startupTimeout = 100 * time.Millisecond
	
	err := m.waitForReady()
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	
	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error message, got: %v", err)
	}
}

func TestManager_WaitForReady_ErrorChannel(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManager(tempDir, 5433)
	
	expectedErr := fmt.Errorf("mock startup error")
	go func() {
		time.Sleep(10 * time.Millisecond)
		m.errCh <- expectedErr
	}()
	
	err := m.waitForReady()
	if err == nil {
		t.Fatal("expected error from error channel, got nil")
	}
	
	if !strings.Contains(err.Error(), "mock startup error") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestManager_LogOutput(t *testing.T) {
	m := NewManager("", 0)
	
	r, w := io.Pipe()
	
	done := make(chan bool)
	go func() {
		m.logOutput(r, "test")
		done <- true
	}()
	
	w.Write([]byte("test log line\n"))
	w.Close()
	
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("logOutput did not complete in time")
	}
}

func TestManager_MonitorProcess_ProcessExit(t *testing.T) {
	tempDir := t.TempDir()
	m := NewManager(tempDir, 5433)
	
	cmd := exec.Command("echo", "test")
	if err := cmd.Start(); err != nil {
		t.Skipf("cannot start echo command: %v", err)
	}
	
	m.process = cmd
	m.running = true
	
	done := make(chan bool)
	go func() {
		m.monitorProcess()
		done <- true
	}()
	
	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}
}
