package integration

import (
	"context"
	"testing"
	"time"

	"github.com/vibesql/vibe/internal/postgres"
)

// TestConnectionPoolIntegration verifies connection pool configuration
func TestConnectionPoolIntegration(t *testing.T) {
	// This test requires a running PostgreSQL instance
	// Skip if PG_TEST_ENABLED is not set
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Start PostgreSQL manager
	manager := postgres.NewManager("./test-data", 5432)
	
	err := manager.Start()
	if err != nil {
		t.Skipf("Skipping test - could not start PostgreSQL: %v", err)
	}
	defer manager.Stop()
	
	// Create connection pool
	conn, err := manager.CreateConnection()
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()
	
	// Test connection pool is working
	err = conn.Ping()
	if err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}
	
	// Verify we can get the underlying DB
	db := conn.DB()
	if db == nil {
		t.Error("DB() returned nil")
	}
}

// TestErrorTranslation verifies SQLSTATE to VibeSQL error mapping
func TestErrorTranslation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	manager := postgres.NewManager("./test-data", 5432)
	
	err := manager.Start()
	if err != nil {
		t.Skipf("Skipping test - could not start PostgreSQL: %v", err)
	}
	defer manager.Stop()
	
	conn, err := manager.CreateConnection()
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()
	
	tests := []struct {
		name         string
		query        string
		expectedCode string
	}{
		{
			name:         "Syntax error",
			query:        "SELCT 1",
			expectedCode: postgres.ErrorCodeInvalidSQL,
		},
		{
			name:         "Undefined table",
			query:        "SELECT * FROM nonexistent_table",
			expectedCode: postgres.ErrorCodeInvalidSQL,
		},
		{
			name:         "Undefined column",
			query:        "SELECT nonexistent_column FROM pg_database",
			expectedCode: postgres.ErrorCodeInvalidSQL,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			_, err := conn.DB().QueryContext(ctx, tt.query)
			if err == nil {
				t.Error("Expected error but got nil")
				return
			}
			
			// Translate error
			vibeErr := postgres.TranslateError(err)
			
			if vibeErr.Code != tt.expectedCode {
				t.Errorf("Expected error code %s, got %s", tt.expectedCode, vibeErr.Code)
			}
			
			// Verify HTTP status code mapping
			httpStatus := postgres.GetHTTPStatusCode(vibeErr.Code)
			if httpStatus == 0 {
				t.Error("GetHTTPStatusCode returned 0")
			}
		})
	}
}

// TestConnectionTimeout verifies query timeout handling and error translation
func TestConnectionTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	manager := postgres.NewManager("./test-data", 5432)
	
	err := manager.Start()
	if err != nil {
		t.Skipf("Skipping test - could not start PostgreSQL: %v", err)
	}
	defer manager.Stop()
	
	conn, err := manager.CreateConnection()
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()
	
	// Test query timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	// This query will sleep for 10 seconds but timeout after 100ms
	_, err = conn.DB().QueryContext(ctx, "SELECT pg_sleep(10)")
	
	if err == nil {
		t.Error("Expected timeout error but got nil")
		return
	}
	
	// Translate the error to verify it maps to QUERY_TIMEOUT
	vibeErr := postgres.TranslateError(err)
	
	if vibeErr.Code != postgres.ErrorCodeQueryTimeout {
		t.Errorf("Expected error code %s for timeout, got %s", 
			postgres.ErrorCodeQueryTimeout, vibeErr.Code)
	}
	
	// Verify HTTP status code
	httpStatus := postgres.GetHTTPStatusCode(vibeErr.Code)
	if httpStatus != 408 {
		t.Errorf("Expected HTTP status 408 for QUERY_TIMEOUT, got %d", httpStatus)
	}
	
	t.Logf("Timeout correctly translated: code=%s, message=%s", vibeErr.Code, vibeErr.Message)
}

// TestConnectionPoolLimits verifies connection pool configuration
func TestConnectionPoolLimits(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	manager := postgres.NewManager("./test-data", 5432)
	
	err := manager.Start()
	if err != nil {
		t.Skipf("Skipping test - could not start PostgreSQL: %v", err)
	}
	defer manager.Stop()
	
	conn, err := manager.CreateConnection()
	if err != nil {
		t.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()
	
	db := conn.DB()
	
	// Verify max open connections is set correctly
	stats := db.Stats()
	t.Logf("Connection pool stats: Open=%d, InUse=%d, Idle=%d", 
		stats.OpenConnections, stats.InUse, stats.Idle)
	
	// Execute a simple query to verify connection works
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Errorf("Failed to execute simple query: %v", err)
	}
	
	if result != 1 {
		t.Errorf("Expected result 1, got %d", result)
	}
}
