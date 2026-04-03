package postgres

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Process manages a postgres subprocess
type Process struct {
	cmd      *exec.Cmd
	dataDir  string
	port     int
	ready    bool
	shutdown chan struct{}
}

// Start initializes and starts postgres
func Start(postgresBin, initdbBin, shareDir, dataDir string, port int) (*Process, error) {
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
	
	cmd := exec.Command(initdbBin, args...)
	cmd.Env = append(os.Environ(),
		"PGSHAREDIR="+shareDir,
	)
	
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
		"-k", "", // Disable Unix sockets on Windows
	}
	
	p.cmd = exec.Command(postgresBin, args...)
	p.cmd.Env = append(os.Environ(),
		"PGSHAREDIR="+shareDir,
	)
	
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
			// Log or discard
		}
	}()
	
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			// Log or discard
		}
	}()
	
	// Wait for ready
	if err := p.waitForReady(); err != nil {
		p.Stop()
		return err
	}
	
	return nil
}

func (p *Process) waitForReady() error {
	// Simple wait - in production would check port
	time.Sleep(1 * time.Second)
	return nil
}

// Stop shuts down postgres gracefully
func (p *Process) Stop() error {
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
