package vsql

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	
	"github.com/vibesql/vibesql-micro/internal/postgres"
)

// DB represents a vibesql-micro database connection
type DB struct {
	path       string
	dataDir    string
	markerFile string
	postgres   *postgres.Process
	port       int
	mu         sync.Mutex
	closed     bool
}

// Open creates or opens a database at the specified path.
//
// The path should end with .vsql (e.g., "./app.vsql").
// This creates a marker file and a hidden data directory.
func Open(path string) (*DB, error) {
	return OpenWithProgress(path, nil)
}

// OpenWithProgress is like Open but shows progress via the callback.
func OpenWithProgress(path string, progress func(string)) (*DB, error) {
	// Default path
	if path == "" {
		path = "./default.vsql"
	}
	
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	
	// Ensure .vsql extension
	if filepath.Ext(absPath) != ".vsql" {
		absPath += ".vsql"
	}
	
	markerFile := absPath
	dataDir := absPath + ".data"
	
	// Check if locked
	if isLocked(dataDir) {
		return nil, WarmError{
			What: fmt.Sprintf("%s is busy", filepath.Base(absPath)),
			Why:  "Another program is using it",
			Fix:  "Close the other program, or use a different file",
		}
	}
	
	// Ensure binaries are extracted
	postgresBin, initdbBin, shareDir, firstTime, err := ensureBinary(progress)
	if err != nil {
		return nil, fmt.Errorf("setup failed: %w", err)
	}
	
	if firstTime && progress != nil {
		progress("done")
	}
	
	// Create data directory if it doesn't exist
	isNew := false
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		isNew = true
		if progress != nil {
			progress("Creating database...")
		}
		
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("create data directory: %w", err)
		}
		
		// Create marker file (initdb will be run by postgres.Start)
		if err := createMarker(markerFile); err != nil {
			os.RemoveAll(dataDir)
			return nil, fmt.Errorf("create marker: %w", err)
		}
		
		if progress != nil {
			progress("done")
		}
	}
	
	// Create lock
	if err := createLock(dataDir); err != nil {
		return nil, fmt.Errorf("lock database: %w", err)
	}
	
	// Find a free port
	port := 5433 // TODO: Find free port
	
	// Start postgres
	pg, err := postgres.Start(postgresBin, initdbBin, shareDir, dataDir, port)
	if err != nil {
		removeLock(dataDir)
		return nil, fmt.Errorf("start postgres: %w", err)
	}
	
	db := &DB{
		path:       absPath,
		dataDir:    dataDir,
		markerFile: markerFile,
		postgres:   pg,
		port:       port,
	}
	
	// If new and no progress shown, show silent success
	if isNew && progress == nil {
		// Silent - no output needed
	}
	
	return db, nil
}

// Close shuts down the database connection.
func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if db.closed {
		return nil
	}
	db.closed = true
	
	// Stop postgres
	if db.postgres != nil {
		db.postgres.Stop()
	}
	
	// Remove lock
	removeLock(db.dataDir)
	
	return nil
}

// Path returns the database path.
func (db *DB) Path() string {
	return db.path
}

// WarmError is a user-friendly error
type WarmError struct {
	What string
	Why  string
	Fix  string
}

func (e WarmError) Error() string {
	msg := e.What
	if e.Why != "" {
		msg += "\n  Why: " + e.Why
	}
	if e.Fix != "" {
		msg += "\n  Fix: " + e.Fix
	}
	return msg
}

// Helper functions (to be implemented)
func isLocked(dataDir string) bool {
	lockFile := filepath.Join(dataDir, ".lock")
	_, err := os.Stat(lockFile)
	return err == nil
}

func createLock(dataDir string) error {
	lockFile := filepath.Join(dataDir, ".lock")
	return os.WriteFile(lockFile, []byte{}, 0644)
}

func removeLock(dataDir string) {
	lockFile := filepath.Join(dataDir, ".lock")
	os.Remove(lockFile)
}

func createMarker(markerFile string) error {
	return os.WriteFile(markerFile, []byte("vibesql-micro database\n"), 0644)
}


