package vsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Row represents a single row as a map of column names to values
type Row map[string]interface{}

// Result represents the result of an Exec operation
type Result struct {
	RowsAffected int64
	LastInsertID int64
}

// Query executes a SQL query and returns rows as JSON-friendly maps.
func (db *DB) Query(sqlStr string, args ...interface{}) ([]Row, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if db.closed {
		return nil, fmt.Errorf("database is closed")
	}
	
	if db.postgres == nil || db.postgres.DB() == nil {
		return nil, fmt.Errorf("database not connected")
	}
	
	// Execute query
	rows, err := db.postgres.DB().Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("get columns: %w", err)
	}
	
	// Scan results
	var results []Row
	
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		
		// Convert to map
		row := make(Row)
		for i, col := range columns {
			val := values[i]
			
			// Handle []byte (from database) as string
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		
		results = append(results, row)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	
	return results, nil
}

// Exec executes a SQL statement (INSERT, UPDATE, DELETE, etc.).
func (db *DB) Exec(sqlStr string, args ...interface{}) (*Result, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if db.closed {
		return nil, fmt.Errorf("database is closed")
	}
	
	if db.postgres == nil || db.postgres.DB() == nil {
		return nil, fmt.Errorf("database not connected")
	}
	
	result, err := db.postgres.DB().Exec(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	
	rowsAffected, _ := result.RowsAffected()
	lastID, _ := result.LastInsertId()
	
	return &Result{
		RowsAffected: rowsAffected,
		LastInsertID: lastID,
	}, nil
}

// QueryJSON executes a query and returns the result as a JSON string.
// This is the raw pass-through format for CLI output.
func (db *DB) QueryJSON(sqlStr string, args ...interface{}) (string, error) {
	rows, err := db.Query(sqlStr, args...)
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
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if db.closed {
		return nil, fmt.Errorf("database is closed")
	}
	
	if db.postgres == nil {
		return nil, fmt.Errorf("database not connected")
	}
	
	return db.postgres.DB(), nil
}
