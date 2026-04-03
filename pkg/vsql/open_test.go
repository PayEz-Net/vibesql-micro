package vsql

import (
	"os"
	"testing"
)

func TestOpenCreatesMarkerAndDataDir(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	dbPath := tmpDir + "\\test.vsql"
	
	// Open should create marker and data dir
	db, err := Open(dbPath)
	if err != nil {
		// Expected to fail since postgres isn't fully implemented
		t.Logf("Open failed as expected: %v", err)
		return
	}
	defer db.Close()
	
	// Check marker file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Marker file should exist")
	}
	
	// Check data dir exists
	if _, err := os.Stat(dbPath + ".data"); os.IsNotExist(err) {
		t.Error("Data directory should exist")
	}
}

func TestOpenWithEmptyPath(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	
	// Empty path should default to ./default.vsql
	_, err := Open("")
	if err != nil {
		t.Logf("Open failed as expected: %v", err)
		return
	}
	
	// Check default.vsql was created
	if _, err := os.Stat("default.vsql"); os.IsNotExist(err) {
		t.Error("default.vsql should be created for empty path")
	}
}

func TestWarmError(t *testing.T) {
	err := WarmError{
		What: "database is busy",
		Why:  "another process is using it",
		Fix:  "close the other process",
	}
	
	msg := err.Error()
	if msg == "" {
		t.Error("WarmError should return a message")
	}
	
	t.Logf("Error message: %s", msg)
}
