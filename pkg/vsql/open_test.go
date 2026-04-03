package vsql

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCreatesMarkerAndDataDir(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Marker file should exist")
	}
	if _, err := os.Stat(dbPath + ".data"); os.IsNotExist(err) {
		t.Error("Data directory should exist")
	}
}

func TestOpenWithEmptyPath(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir("/")

	db, err := Open("")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat("default.vsql"); os.IsNotExist(err) {
		t.Error("default.vsql should be created for empty path")
	}
}

func TestOpenWithNonVsqlExtension(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "noext")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(dbPath + ".vsql"); os.IsNotExist(err) {
		t.Error(".vsql extension should be appended")
	}
}

func TestOpenAbsolutePath(t *testing.T) {
	tmpDir := t.TempDir()
	relPath := filepath.Join(tmpDir, "..", filepath.Base(tmpDir), "abs.vsql")

	db, err := Open(relPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	absPath, _ := filepath.Abs(relPath)
	if db.path != absPath {
		t.Errorf("Expected absolute path %s, got %s", absPath, db.path)
	}
}

func TestOpenExistingDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "existing.vsql")

	db1, err := Open(dbPath)
	if err != nil {
		t.Fatalf("First open failed: %v", err)
	}
	db1.Exec("CREATE TABLE existing_test (id INT)")
	db1.Close()

	db2, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Second open failed: %v", err)
	}
	defer db2.Close()

	rows, err := db2.Query("SELECT * FROM existing_test")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(rows))
	}
}

func TestOpenWhileLocked(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "locked.vsql")

	db1, err := Open(dbPath)
	if err != nil {
		t.Fatalf("First open failed: %v", err)
	}
	defer db1.Close()

	_, err = Open(dbPath)
	if err == nil {
		t.Fatal("Expected error when opening locked database")
	}

	warmErr, ok := err.(WarmError)
	if !ok {
		t.Fatalf("Expected WarmError, got %T", err)
	}
	if warmErr.What == "" {
		t.Error("WarmError.What should not be empty")
	}
}

func TestCloseRemovesLock(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "unlock.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	lockPath := dbPath + ".data.lock"
	db.Close()

	if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
		t.Error("Lock file should be removed after Close")
	}
}

func TestDoubleClose(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "double.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	db.Close()

	// Second close should not panic
	err = db.Close()
	if err != nil {
		t.Logf("Second close returned error (acceptable): %v", err)
	}
}

func TestReopenAfterClose(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "reopen.vsql")

	db1, err := Open(dbPath)
	if err != nil {
		t.Fatalf("First open failed: %v", err)
	}
	db1.Exec("CREATE TABLE reopen_test (id INT)")
	db1.Exec("INSERT INTO reopen_test VALUES (42)")
	db1.Close()

	db2, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Reopen failed: %v", err)
	}
	defer db2.Close()

	rows, err := db2.Query("SELECT * FROM reopen_test")
	if err != nil {
		t.Fatalf("Query after reopen failed: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("Expected 1 row after reopen, got %d", len(rows))
	}
}

func TestOpenProgressCallback(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "progress.vsql")
	messages := []string{}

	db, err := OpenWithProgress(dbPath, func(msg string) {
		messages = append(messages, msg)
	})
	if err != nil {
		t.Fatalf("OpenWithProgress failed: %v", err)
	}
	defer db.Close()

	if len(messages) == 0 {
		t.Error("Progress callback should have fired at least once")
	}
}

func TestWarmErrorString(t *testing.T) {
	err := WarmError{
		What: "database is busy",
		Why:  "another process is using it",
		Fix:  "close the other process",
	}

	msg := err.Error()
	if msg == "" {
		t.Error("WarmError.Error() should not return empty string")
	}
	if msg != "database is busy" {
		t.Logf("Error message: %s", msg)
	}
}
