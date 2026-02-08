package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	maxOpenConnections = 5
	maxIdleConnections = 2
	connMaxLifetime    = 1 * time.Hour
	connMaxIdleTime    = 10 * time.Minute
)

// Connection represents a PostgreSQL database connection pool
type Connection struct {
	db *sql.DB
}

// NewConnection creates a new connection pool to the PostgreSQL database
func NewConnection(host string, port int, user string, password string, dbname string) (*Connection, error) {
	connStr := buildConnectionString(host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	
	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &Connection{db: db}, nil
}

// NewConnectionSimple creates a connection with simplified parameters for localhost
func NewConnectionSimple(port int) (*Connection, error) {
	return NewConnection("127.0.0.1", port, "postgres", "", "postgres")
}

// buildConnectionString constructs a PostgreSQL connection string
func buildConnectionString(host string, port int, user string, password string, dbname string) string {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable statement_timeout=5000",
		host, port, user, dbname)
	
	if password != "" {
		connStr += fmt.Sprintf(" password=%s", password)
	}
	
	return connStr
}

// DB returns the underlying database connection pool
func (c *Connection) DB() *sql.DB {
	return c.db
}

// Close closes the database connection pool
func (c *Connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Ping verifies the connection is still alive
func (c *Connection) Ping() error {
	if c.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	return c.db.Ping()
}
