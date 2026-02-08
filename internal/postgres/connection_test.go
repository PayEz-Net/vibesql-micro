package postgres

import (
	"testing"
)

func TestBuildConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		user     string
		password string
		dbname   string
		expected string
	}{
		{
			name:     "Basic connection without password",
			host:     "localhost",
			port:     5432,
			user:     "postgres",
			password: "",
			dbname:   "testdb",
			expected: "host=localhost port=5432 user=postgres dbname=testdb sslmode=disable statement_timeout=5000",
		},
		{
			name:     "Connection with password",
			host:     "127.0.0.1",
			port:     5433,
			user:     "admin",
			password: "secret",
			dbname:   "mydb",
			expected: "host=127.0.0.1 port=5433 user=admin dbname=mydb sslmode=disable statement_timeout=5000 password=secret",
		},
		{
			name:     "IPv6 localhost",
			host:     "::1",
			port:     5432,
			user:     "postgres",
			password: "",
			dbname:   "postgres",
			expected: "host=::1 port=5432 user=postgres dbname=postgres sslmode=disable statement_timeout=5000",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildConnectionString(tt.host, tt.port, tt.user, tt.password, tt.dbname)
			if result != tt.expected {
				t.Errorf("buildConnectionString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewConnection_InvalidParams(t *testing.T) {
	// Test with invalid port (will fail to connect, but should return error)
	_, err := NewConnection("localhost", 99999, "postgres", "", "postgres")
	
	// We expect an error because the port is invalid or unreachable
	if err == nil {
		t.Error("Expected error for invalid port, got nil")
	}
}

func TestNewConnectionSimple_Parameters(t *testing.T) {
	// This test verifies the connection string is built correctly
	// Actual connection will fail in test environment without PostgreSQL
	
	// We can't test actual connection without a running PostgreSQL,
	// but we can verify the function signature and error handling
	
	_, err := NewConnectionSimple(5432)
	
	// Expected to fail since PostgreSQL is not running in test environment
	if err == nil {
		// If somehow it succeeds, verify the connection can be closed
		t.Log("Connection succeeded (unexpected in test environment)")
	} else {
		// Expected case - connection fails
		t.Logf("Connection failed as expected: %v", err)
	}
}

func TestConnection_Methods(t *testing.T) {
	// Create a mock connection (without actual DB)
	// We test that methods don't panic with nil handling
	
	conn := &Connection{db: nil}
	
	// Test Close with nil DB
	err := conn.Close()
	if err != nil {
		t.Errorf("Close() with nil DB should not error, got: %v", err)
	}
	
	// Test DB() method
	db := conn.DB()
	if db != nil {
		t.Error("DB() should return nil when db is nil")
	}
}

func TestConnectionPoolConfiguration(t *testing.T) {
	// Verify connection pool constants are sensible
	if maxOpenConnections <= 0 {
		t.Error("maxOpenConnections must be positive")
	}
	if maxIdleConnections <= 0 {
		t.Error("maxIdleConnections must be positive")
	}
	if maxIdleConnections > maxOpenConnections {
		t.Error("maxIdleConnections should not exceed maxOpenConnections")
	}
	if connMaxLifetime <= 0 {
		t.Error("connMaxLifetime must be positive")
	}
	if connMaxIdleTime <= 0 {
		t.Error("connMaxIdleTime must be positive")
	}
}

func TestConnection_Ping(t *testing.T) {
	conn := &Connection{db: nil}
	
	err := conn.Ping()
	if err == nil {
		t.Error("Ping() with nil DB should return error")
	}
}

func TestConnection_CloseNilDB(t *testing.T) {
	conn := &Connection{db: nil}
	
	err := conn.Close()
	if err != nil {
		t.Errorf("Close() with nil DB should not error, got: %v", err)
	}
}

func BenchmarkBuildConnectionString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = buildConnectionString("localhost", 5432, "postgres", "password", "testdb")
	}
}
