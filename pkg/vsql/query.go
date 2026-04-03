package vsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	
	_ "github.com/lib/pq"
)

// Row represents a single row as a map of column names to values
type Row map[string]interface{}

// Result represents the result of an Exec operation
type Result struct {
	RowsAffected int64
	LastInsertID int64
}

// Query executes a SQL query and returns rows as JSON-friendly maps.
func (db *DB) Query(sql string, args ...interface{}) ([]Row, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if db.closed {
		return nil, fmt.Errorf("database is closed")
	}
	
	// TODO: Connect to postgres and execute query
	// For now, return mock data
	return []Row{
		{"result": "not implemented - work in progress"},
	}, nil
}

// Exec executes a SQL statement (INSERT, UPDATE, DELETE, etc.).
func (db *DB) Exec(sql string, args ...interface{}) (*Result, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if db.closed {
		return nil, fmt.Errorf("database is closed")
	}
	
	// TODO: Connect to postgres and execute statement
	return &Result{}, fmt.Errorf("not implemented - work in progress")
}

// QueryJSON executes a query and returns the result as a JSON string.
// This is the raw pass-through format for CLI output.
func (db *DB) QueryJSON(sql string, args ...interface{}) (string, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return "", err
	}
	
	data, err := json.Marshal(rows)
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}
	
	return string(data), nil
}

// OpenDB opens a direct database/sql connection for advanced use.
// Most users should use Query/Exec instead.
func (db *DB) OpenDB() (*sql.DB, error) {
	// TODO: Return connection to postgres
	return nil, fmt.Errorf("not implemented")
}
