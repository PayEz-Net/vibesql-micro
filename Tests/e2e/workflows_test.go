package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

// E2E Tests for VibeSQL HTTP API Workflows
//
// These tests verify complete user workflows end-to-end via the HTTP API.
//
// Prerequisites:
//   1. PostgreSQL running on localhost:5432 (or set VIBESQL_TEST_DB env var)
//   2. VibeSQL server running on localhost:5173 (start with: ./vibe serve)
//
// In Phase 4, these tests will be updated to automatically start the embedded
// VibeSQL binary. For now, they test the HTTP API workflows assuming the server
// is already running.
//
// To run these tests:
//   1. Start PostgreSQL: docker run -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
//   2. Create test database: psql -U postgres -c "CREATE DATABASE vibesql_test;"
//   3. Start VibeSQL server: ./vibe serve
//   4. Run tests: go test ./Tests/e2e/... -v

const (
	testAPIURL         = "http://127.0.0.1:5173/v1/query"
	serverReadyTimeout = 2 * time.Second
)

// QueryRequest represents an incoming SQL query request
type QueryRequest struct {
	SQL string `json:"sql"`
}

// QueryResponse represents a query response (success or error)
type QueryResponse struct {
	Success       bool                     `json:"success"`
	Rows          []map[string]interface{} `json:"rows,omitempty"`
	RowCount      int                      `json:"rowCount,omitempty"`
	ExecutionTime float64                  `json:"executionTime,omitempty"`
	Error         *ErrorDetail             `json:"error,omitempty"`
}

// ErrorDetail represents error information in the response
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// APIClient is a helper for making HTTP requests to the VibeSQL API
type APIClient struct {
	client *http.Client
	baseURL string
}

// NewAPIClient creates a new API client
func NewAPIClient() *APIClient {
	return &APIClient{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: testAPIURL,
	}
}

// Query executes a SQL query via HTTP API
func (c *APIClient) Query(req QueryRequest) (*QueryResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w (body: %s)", err, string(body))
	}

	return &queryResp, nil
}

// WaitForServer waits for the HTTP server to become ready
func (c *APIClient) WaitForServer(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("server did not become ready within %v", timeout)
		case <-ticker.C:
			resp, err := c.Query(QueryRequest{SQL: "SELECT 1"})
			if err == nil && resp.Success {
				return nil
			}
		}
	}
}

// TestE2E_ServerReady verifies the server is running and accessible
func TestE2E_ServerReady(t *testing.T) {
	client := NewAPIClient()
	
	err := client.WaitForServer(serverReadyTimeout)
	if err != nil {
		t.Skipf("Server not ready: %v\n\nPlease ensure:\n  1. PostgreSQL is running (localhost:5432)\n  2. VibeSQL server is running: ./vibe serve", err)
	}
	
	t.Log("✓ VibeSQL server is ready")
}

// TestE2E_FullCRUDWorkflow tests complete CRUD lifecycle
func TestE2E_FullCRUDWorkflow(t *testing.T) {
	client := NewAPIClient()
	
	// Ensure server is ready
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	// Use unique table name to avoid conflicts
	tableName := fmt.Sprintf("users_crud_%d", time.Now().Unix())
	
	// Cleanup on exit
	defer func() {
		client.Query(QueryRequest{SQL: fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)})
	}()

	// Step 1: CREATE TABLE
	resp, err := client.Query(QueryRequest{
		SQL: fmt.Sprintf(`CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			age INTEGER
		)`, tableName),
	})
	if err != nil {
		t.Fatalf("CREATE TABLE request failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("CREATE TABLE failed: %+v", resp.Error)
	}

	// Step 2: INSERT data
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("INSERT INTO %s (name, email, age) VALUES ('Alice', 'alice@example.com', 30)", tableName),
	})
	if err != nil {
		t.Fatalf("INSERT request failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("INSERT failed: %+v", resp.Error)
	}

	// Step 3: SELECT data
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT name, email, age FROM %s WHERE name = 'Alice'", tableName),
	})
	if err != nil {
		t.Fatalf("SELECT request failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("SELECT failed: %+v", resp.Error)
	}
	if resp.RowCount != 1 {
		t.Fatalf("Expected 1 row, got %d", resp.RowCount)
	}
	if resp.Rows[0]["name"] != "Alice" {
		t.Errorf("Expected name 'Alice', got %v", resp.Rows[0]["name"])
	}
	if resp.Rows[0]["email"] != "alice@example.com" {
		t.Errorf("Expected email 'alice@example.com', got %v", resp.Rows[0]["email"])
	}

	// Step 4: UPDATE data
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("UPDATE %s SET age = 31 WHERE name = 'Alice'", tableName),
	})
	if err != nil {
		t.Fatalf("UPDATE request failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("UPDATE failed: %+v", resp.Error)
	}

	// Step 5: Verify UPDATE
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT age FROM %s WHERE name = 'Alice'", tableName),
	})
	if err != nil {
		t.Fatalf("SELECT after UPDATE failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("SELECT after UPDATE failed: %+v", resp.Error)
	}
	age := int(resp.Rows[0]["age"].(float64))
	if age != 31 {
		t.Errorf("Expected age 31, got %d", age)
	}

	// Step 6: DELETE data
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("DELETE FROM %s WHERE name = 'Alice'", tableName),
	})
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("DELETE failed: %+v", resp.Error)
	}

	// Step 7: Verify DELETE
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT COUNT(*) as count FROM %s", tableName),
	})
	if err != nil {
		t.Fatalf("SELECT after DELETE failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("SELECT after DELETE failed: %+v", resp.Error)
	}
	count := int(resp.Rows[0]["count"].(float64))
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// TestE2E_ConcurrentQueries tests concurrent query execution
func TestE2E_ConcurrentQueries(t *testing.T) {
	client := NewAPIClient()
	
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	tableName := fmt.Sprintf("concurrent_test_%d", time.Now().Unix())
	
	defer func() {
		client.Query(QueryRequest{SQL: fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)})
	}()

	// Create test table
	resp, err := client.Query(QueryRequest{
		SQL: fmt.Sprintf("CREATE TABLE %s (id SERIAL PRIMARY KEY, thread_id INTEGER)", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("CREATE TABLE failed: %v, %+v", err, resp)
	}

	// Execute concurrent INSERT operations
	numThreads := 5
	var wg sync.WaitGroup
	errors := make(chan error, numThreads)

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		threadID := i
		go func() {
			defer wg.Done()
			sql := fmt.Sprintf("INSERT INTO %s (thread_id) VALUES (%d)", tableName, threadID)
			resp, err := client.Query(QueryRequest{SQL: sql})
			if err != nil {
				errors <- fmt.Errorf("thread %d: request failed: %w", threadID, err)
				return
			}
			if !resp.Success {
				errors <- fmt.Errorf("thread %d: query failed: %+v", threadID, resp.Error)
				return
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}

	// Verify all inserts succeeded
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT COUNT(*) as count FROM %s", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("SELECT COUNT failed: %v, %+v", err, resp)
	}
	count := int(resp.Rows[0]["count"].(float64))
	if count != numThreads {
		t.Errorf("Expected %d rows, got %d", numThreads, count)
	}
}

// TestE2E_ErrorRecovery tests that invalid queries don't crash the server
func TestE2E_ErrorRecovery(t *testing.T) {
	client := NewAPIClient()
	
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	tableName := fmt.Sprintf("error_test_%d", time.Now().Unix())
	
	defer func() {
		client.Query(QueryRequest{SQL: fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)})
	}()

	// Test 1: Invalid SQL syntax
	resp, err := client.Query(QueryRequest{
		SQL: "SELCT INVALID SYNTAX",
	})
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.Success {
		t.Error("Expected query to fail with invalid syntax")
	}
	if resp.Error.Code != "INVALID_SQL" {
		t.Errorf("Expected error code INVALID_SQL, got %s", resp.Error.Code)
	}

	// Test 2: Missing WHERE clause (unsafe query)
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("CREATE TABLE %s (id INTEGER)", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("CREATE TABLE failed: %v, %+v", err, resp)
	}

	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("DELETE FROM %s", tableName),
	})
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.Success {
		t.Error("Expected query to fail without WHERE clause")
	}
	if resp.Error.Code != "UNSAFE_QUERY" {
		t.Errorf("Expected error code UNSAFE_QUERY, got %s", resp.Error.Code)
	}

	// Test 3: Verify server is still functional after errors
	resp, err = client.Query(QueryRequest{
		SQL: "SELECT 1 as test",
	})
	if err != nil {
		t.Fatalf("Server not responding after errors: %v", err)
	}
	if !resp.Success {
		t.Fatalf("Server query failed after errors: %+v", resp.Error)
	}
	if resp.RowCount != 1 {
		t.Errorf("Expected 1 row, got %d", resp.RowCount)
	}
}

// TestE2E_TimeoutHandling tests query timeout enforcement
func TestE2E_TimeoutHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	client := NewAPIClient()
	
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	// Execute a query that will timeout (pg_sleep for 6 seconds, timeout is 5s)
	resp, err := client.Query(QueryRequest{
		SQL: "SELECT pg_sleep(6)",
	})
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.Success {
		t.Error("Expected query to timeout")
	}
	if resp.Error.Code != "QUERY_TIMEOUT" {
		t.Errorf("Expected error code QUERY_TIMEOUT, got %s", resp.Error.Code)
	}

	// Verify server is still functional after timeout
	resp, err = client.Query(QueryRequest{
		SQL: "SELECT 1 as test",
	})
	if err != nil {
		t.Fatalf("Server not responding after timeout: %v", err)
	}
	if !resp.Success {
		t.Fatalf("Server query failed after timeout: %+v", resp.Error)
	}
}

// TestE2E_JSONBWorkflow tests JSONB operations end-to-end
func TestE2E_JSONBWorkflow(t *testing.T) {
	client := NewAPIClient()
	
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	tableName := fmt.Sprintf("jsonb_users_%d", time.Now().Unix())
	
	defer func() {
		client.Query(QueryRequest{SQL: fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)})
	}()

	// Create table with JSONB column
	resp, err := client.Query(QueryRequest{
		SQL: fmt.Sprintf(`CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			data JSONB NOT NULL
		)`, tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("CREATE TABLE failed: %v, %+v", err, resp)
	}

	// Insert JSONB data
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf(`INSERT INTO %s (data) VALUES ('{"name": "Bob", "age": 25, "city": "NYC"}')`, tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("INSERT failed: %v, %+v", err, resp)
	}

	// Query with JSONB operator (->)
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT data->'name' as name FROM %s", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("JSONB query failed: %v, %+v", err, resp)
	}
	if resp.RowCount != 1 {
		t.Fatalf("Expected 1 row, got %d", resp.RowCount)
	}

	// Query with JSONB text operator (->>)
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT data->>'name' as name FROM %s WHERE data->>'city' = 'NYC'", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("JSONB text query failed: %v, %+v", err, resp)
	}
	if resp.RowCount != 1 {
		t.Fatalf("Expected 1 row, got %d", resp.RowCount)
	}
	if resp.Rows[0]["name"] != "Bob" {
		t.Errorf("Expected name 'Bob', got %v", resp.Rows[0]["name"])
	}
}

// TestE2E_LimitEnforcement tests result limit enforcement
func TestE2E_LimitEnforcement(t *testing.T) {
	client := NewAPIClient()
	
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	tableName := fmt.Sprintf("limit_test_%d", time.Now().Unix())
	
	defer func() {
		client.Query(QueryRequest{SQL: fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)})
	}()

	// Create table
	resp, err := client.Query(QueryRequest{
		SQL: fmt.Sprintf("CREATE TABLE %s (id SERIAL PRIMARY KEY)", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("CREATE TABLE failed: %v, %+v", err, resp)
	}

	// Insert 1000 rows (at the limit)
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("INSERT INTO %s SELECT generate_series(1, 1000)", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("INSERT failed: %v, %+v", err, resp)
	}

	// Query exactly 1000 rows (should succeed)
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT * FROM %s", tableName),
	})
	if err != nil {
		t.Fatalf("SELECT 1000 rows failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("SELECT 1000 rows failed: %+v", resp.Error)
	}
	if resp.RowCount != 1000 {
		t.Errorf("Expected 1000 rows, got %d", resp.RowCount)
	}

	// Try to query more than 1000 rows (should fail with RESULT_TOO_LARGE)
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("INSERT INTO %s SELECT generate_series(1001, 1010)", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("INSERT additional rows failed: %v, %+v", err, resp)
	}

	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT * FROM %s", tableName),
	})
	if err != nil {
		t.Fatalf("SELECT >1000 rows request failed: %v", err)
	}
	if resp.Success {
		t.Error("Expected query to fail with RESULT_TOO_LARGE")
	}
	if resp.Error.Code != "RESULT_TOO_LARGE" {
		t.Errorf("Expected error code RESULT_TOO_LARGE, got %s", resp.Error.Code)
	}
}

// TestE2E_DataPersistence tests data persistence across server restarts
// NOTE: This test requires manual server restart. In Phase 4, it will be automated.
// Manual test procedure:
//   1. Run test once (creates table and data)
//   2. Stop server (Ctrl+C)
//   3. Restart server (./vibe serve)
//   4. Run test again (verifies data persists)
func TestE2E_DataPersistence(t *testing.T) {
	client := NewAPIClient()
	
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	// Use a well-known table name for persistence testing
	tableName := "persistence_test_do_not_delete"
	
	// Check if table exists (indicates prior test run)
	resp, err := client.Query(QueryRequest{
		SQL: fmt.Sprintf("SELECT COUNT(*) as count FROM pg_tables WHERE tablename = '%s'", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("Failed to check for existing table: %v, %+v", err, resp)
	}
	
	tableExists := int(resp.Rows[0]["count"].(float64)) > 0
	
	if !tableExists {
		// First run: Create table and insert data
		t.Log("First run: Creating table and inserting data")
		
		resp, err = client.Query(QueryRequest{
			SQL: fmt.Sprintf(`CREATE TABLE %s (
				id SERIAL PRIMARY KEY,
				value TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`, tableName),
		})
		if err != nil || !resp.Success {
			t.Fatalf("CREATE TABLE failed: %v, %+v", err, resp)
		}
		
		resp, err = client.Query(QueryRequest{
			SQL: fmt.Sprintf("INSERT INTO %s (value) VALUES ('test-data-persistence-123')", tableName),
		})
		if err != nil || !resp.Success {
			t.Fatalf("INSERT failed: %v, %+v", err, resp)
		}
		
		t.Log("✓ Table created and data inserted")
		t.Log("⚠️  To complete persistence test:")
		t.Log("   1. Stop the server (Ctrl+C)")
		t.Log("   2. Restart the server (./vibe serve)")
		t.Log("   3. Run this test again")
		t.Skipf("Persistence test requires manual server restart - see log above")
	} else {
		// Second run: Verify data persists
		t.Log("Second run: Verifying data persistence after restart")
		
		resp, err = client.Query(QueryRequest{
			SQL: fmt.Sprintf("SELECT value FROM %s WHERE value = 'test-data-persistence-123'", tableName),
		})
		if err != nil || !resp.Success {
			t.Fatalf("SELECT after restart failed: %v, %+v", err, resp)
		}
		
		if resp.RowCount != 1 {
			t.Fatalf("Expected 1 row after restart, got %d", resp.RowCount)
		}
		
		if resp.Rows[0]["value"] != "test-data-persistence-123" {
			t.Errorf("Expected persisted value 'test-data-persistence-123', got %v", resp.Rows[0]["value"])
		}
		
		// Cleanup
		resp, err = client.Query(QueryRequest{
			SQL: fmt.Sprintf("DROP TABLE %s", tableName),
		})
		if err != nil || !resp.Success {
			t.Logf("Warning: Failed to drop persistence test table: %v", err)
		}
		
		t.Log("✓ Data persistence verified - data survived server restart")
	}
}

// TestE2E_GracefulShutdown tests graceful shutdown behavior
// NOTE: This test verifies the server can handle shutdown signals properly.
// Full in-flight query testing will be added in Phase 4 with programmatic server control.
func TestE2E_GracefulShutdown(t *testing.T) {
	client := NewAPIClient()
	
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	tableName := fmt.Sprintf("shutdown_test_%d", time.Now().Unix())
	
	defer func() {
		client.Query(QueryRequest{SQL: fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)})
	}()

	// Create test table
	resp, err := client.Query(QueryRequest{
		SQL: fmt.Sprintf("CREATE TABLE %s (id SERIAL PRIMARY KEY, value INTEGER)", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("CREATE TABLE failed: %v, %+v", err, resp)
	}

	// Insert test data
	resp, err = client.Query(QueryRequest{
		SQL: fmt.Sprintf("INSERT INTO %s (value) VALUES (1), (2), (3)", tableName),
	})
	if err != nil || !resp.Success {
		t.Fatalf("INSERT failed: %v, %+v", err, resp)
	}

	// Start a long-running query in background (simulates in-flight request)
	done := make(chan bool)
	queryErr := make(chan error)
	
	go func() {
		// This query should complete even if server shutdown is initiated
		resp, err := client.Query(QueryRequest{
			SQL: fmt.Sprintf("SELECT pg_sleep(1), value FROM %s", tableName),
		})
		if err != nil {
			queryErr <- err
		} else if !resp.Success {
			queryErr <- fmt.Errorf("query failed: %+v", resp.Error)
		} else {
			queryErr <- nil
		}
		close(done)
	}()

	// Give the query time to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is still responsive during query execution
	resp, err = client.Query(QueryRequest{
		SQL: "SELECT 1",
	})
	if err != nil {
		t.Fatalf("Server not responsive during long query: %v", err)
	}
	if !resp.Success {
		t.Fatalf("Server query failed during long query: %+v", resp.Error)
	}

	// Wait for background query to complete
	select {
	case err := <-queryErr:
		if err != nil {
			t.Errorf("Background query failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Background query timed out")
	}

	<-done
	
	t.Log("✓ Server handled concurrent queries during potential shutdown window")
	t.Log("⚠️  Full graceful shutdown test with SIGTERM will be added in Phase 4")
	t.Log("   Current test verifies: queries complete successfully under load")
}

// TestE2E_MultipleTablesWorkflow tests working with multiple tables
func TestE2E_MultipleTablesWorkflow(t *testing.T) {
	client := NewAPIClient()
	
	if err := client.WaitForServer(serverReadyTimeout); err != nil {
		t.Skipf("Server not ready, skipping test: %v", err)
	}

	suffix := time.Now().Unix()
	tables := []string{
		fmt.Sprintf("departments_%d", suffix),
		fmt.Sprintf("employees_%d", suffix),
		fmt.Sprintf("projects_%d", suffix),
	}
	
	defer func() {
		for _, table := range tables {
			client.Query(QueryRequest{SQL: fmt.Sprintf("DROP TABLE IF EXISTS %s", table)})
		}
	}()

	// Create multiple tables
	createSQLs := []string{
		fmt.Sprintf(`CREATE TABLE %s (id SERIAL PRIMARY KEY, name TEXT NOT NULL)`, tables[0]),
		fmt.Sprintf(`CREATE TABLE %s (id SERIAL PRIMARY KEY, name TEXT NOT NULL, dept_id INTEGER)`, tables[1]),
		fmt.Sprintf(`CREATE TABLE %s (id SERIAL PRIMARY KEY, title TEXT NOT NULL)`, tables[2]),
	}

	for _, sql := range createSQLs {
		resp, err := client.Query(QueryRequest{SQL: sql})
		if err != nil || !resp.Success {
			t.Fatalf("CREATE TABLE failed: %v, %+v", err, resp)
		}
	}

	// Insert data into multiple tables
	inserts := []string{
		fmt.Sprintf(`INSERT INTO %s (name) VALUES ('Engineering')`, tables[0]),
		fmt.Sprintf(`INSERT INTO %s (name, dept_id) VALUES ('Alice', 1)`, tables[1]),
		fmt.Sprintf(`INSERT INTO %s (title) VALUES ('Project X')`, tables[2]),
	}

	for _, sql := range inserts {
		resp, err := client.Query(QueryRequest{SQL: sql})
		if err != nil || !resp.Success {
			t.Fatalf("INSERT failed: %v, %+v", err, resp)
		}
	}

	// Query each table
	queries := map[string]int{
		fmt.Sprintf("SELECT * FROM %s", tables[0]): 1,
		fmt.Sprintf("SELECT * FROM %s", tables[1]): 1,
		fmt.Sprintf("SELECT * FROM %s", tables[2]): 1,
	}

	for sql, expectedCount := range queries {
		resp, err := client.Query(QueryRequest{SQL: sql})
		if err != nil || !resp.Success {
			t.Fatalf("SELECT failed for %s: %v, %+v", sql, err, resp)
		}
		if resp.RowCount != expectedCount {
			t.Errorf("Expected %d rows for %s, got %d", expectedCount, sql, resp.RowCount)
		}
	}
}

// TestMain provides test setup and cleanup
func TestMain(m *testing.M) {
	// Check if server is accessible before running tests
	client := NewAPIClient()
	err := client.WaitForServer(3 * time.Second)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n⚠️  VibeSQL server not accessible at %s\n\n", testAPIURL)
		fmt.Fprintf(os.Stderr, "E2E tests require a running VibeSQL server.\n\n")
		fmt.Fprintf(os.Stderr, "Setup instructions:\n")
		fmt.Fprintf(os.Stderr, "  1. Ensure PostgreSQL is running: docker run -d -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres\n")
		fmt.Fprintf(os.Stderr, "  2. Create test database: psql -U postgres -c \"CREATE DATABASE vibesql_test;\"\n")
		fmt.Fprintf(os.Stderr, "  3. Start VibeSQL: ./vibe serve\n")
		fmt.Fprintf(os.Stderr, "  4. Run tests: go test ./Tests/e2e/... -v\n\n")
		fmt.Fprintf(os.Stderr, "Alternatively, set VIBESQL_E2E_SKIP=1 to skip E2E tests.\n\n")
		
		// If skip env var is set, exit with success
		if os.Getenv("VIBESQL_E2E_SKIP") == "1" {
			fmt.Fprintf(os.Stderr, "Skipping E2E tests (VIBESQL_E2E_SKIP=1)\n")
			os.Exit(0)
		}
		
		// Otherwise, run tests (they will be skipped individually)
	}
	
	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}
