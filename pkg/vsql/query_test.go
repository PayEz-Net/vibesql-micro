package vsql

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestQueryAfterOpen(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "query.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1 as one")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(rows))
	}
	if rows[0]["one"] != int64(1) {
		t.Errorf("Expected 1, got %v (type %T)", rows[0]["one"], rows[0]["one"])
	}
}

func TestExecCreateTable(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "exec.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE test_exec (id SERIAL PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("CREATE TABLE failed: %v", err)
	}
}

func TestExecInsert(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "insert.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE test_insert (id SERIAL PRIMARY KEY, name TEXT)")
	res, err := db.Exec("INSERT INTO test_insert (name) VALUES ($1)", "test")
	if err != nil {
		t.Fatalf("INSERT failed: %v", err)
	}
	if res.RowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", res.RowsAffected)
	}
}

func TestQueryWithParams(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "params.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE test_params (id SERIAL PRIMARY KEY, name TEXT, age INT)")
	db.Exec("INSERT INTO test_params (name, age) VALUES ($1, $2), ($3, $4)", "Alice", 30, "Bob", 25)

	rows, err := db.Query("SELECT * FROM test_params WHERE age > $1", 26)
	if err != nil {
		t.Fatalf("Query with params failed: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}
	if rows[0]["name"] != "Alice" {
		t.Errorf("Expected Alice, got %v", rows[0]["name"])
	}
}

func TestQueryJSON(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "json.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE test_json (id SERIAL PRIMARY KEY, name TEXT)")
	db.Exec("INSERT INTO test_json (name) VALUES ($1)", "json-test")

	jsonStr, err := db.QueryJSON("SELECT * FROM test_json WHERE name = $1", "json-test")
	if err != nil {
		t.Fatalf("QueryJSON failed: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
}

func TestQueryAfterClose(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "closed_query.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	db.Close()

	_, err = db.Query("SELECT 1")
	if err == nil {
		t.Fatal("Expected error when querying after close")
	}
}

func TestExecAfterClose(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "closed_exec.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	db.Close()

	_, err = db.Exec("SELECT 1")
	if err == nil {
		t.Fatal("Expected error when executing after close")
	}
}

func TestJSONBWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "jsonb.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE docs (id SERIAL, data JSONB)")
	if err != nil {
		t.Fatalf("CREATE TABLE failed: %v", err)
	}

	_, err = db.Exec("INSERT INTO docs (data) VALUES ($1)", `{"name": "test", "active": true}`)
	if err != nil {
		t.Fatalf("INSERT failed: %v", err)
	}

	rows, err := db.Query("SELECT data->>$1 as name FROM docs WHERE data @> $2", "name", `{"active": true}`)
	if err != nil {
		t.Fatalf("JSONB query failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(rows))
	}
	if rows[0]["name"] != "test" {
		t.Errorf("Expected 'test', got %v", rows[0]["name"])
	}
}

func TestUnicodeAndSpecialChars(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "unicode.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE special (id SERIAL, data TEXT)")
	if err != nil {
		t.Fatalf("CREATE TABLE failed: %v", err)
	}

	_, err = db.Exec("INSERT INTO special (data) VALUES ($1), ($2), ($3)", "日本語", "🎉", "'quotes'")
	if err != nil {
		t.Fatalf("INSERT failed: %v", err)
	}

	rows, err := db.Query("SELECT * FROM special ORDER BY id")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("Expected 3 rows, got %d", len(rows))
	}

	found := false
	for _, row := range rows {
		if row["data"] == "日本語" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Unicode data not preserved correctly: %v", rows)
	}
}

func TestLargeResultSet(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "large.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE large (id SERIAL)")
	if err != nil {
		t.Fatalf("CREATE TABLE failed: %v", err)
	}

	_, err = db.Exec("INSERT INTO large SELECT generate_series(1, 1000)")
	if err != nil {
		t.Fatalf("INSERT failed: %v", err)
	}

	rows, err := db.Query("SELECT COUNT(*) as c FROM large")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	countStr := ""
	switch v := rows[0]["c"].(type) {
	case int64:
		countStr = string(rune('0' + v))
	case string:
		countStr = v
	default:
		countStr = ""
	}
	// COUNT returns int64
	if rows[0]["c"] != int64(1000) {
		t.Errorf("Expected 1000, got %v (type %T)", rows[0]["c"], rows[0]["c"])
	}
	_ = countStr
}

func TestSyntaxError(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "syntax.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Query("INVALID SYNTAX HERE")
	if err == nil {
		t.Fatal("Expected syntax error")
	}
}

func TestWrongParameterCount(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "paramcount.vsql")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE params (id SERIAL, name TEXT)")
	_, err = db.Query("SELECT * FROM params WHERE name = $1 AND id = $2", "only_one")
	if err == nil {
		t.Fatal("Expected parameter count error")
	}
}
