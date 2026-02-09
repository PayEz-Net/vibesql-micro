package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

// Performance Tests for VibeSQL
//
// These tests verify performance requirements and system stability under load.
//
// Prerequisites:
//   1. PostgreSQL running on localhost:5432
//   2. VibeSQL server running on localhost:5173
//
// Run benchmarks:
//   go test ./Tests/performance/... -bench=. -benchmem
//
// Run load tests:
//   go test ./Tests/performance/... -v -run TestLoad

const (
	testAPIURL = "http://127.0.0.1:5173/v1/query"
	testTimeout = 30 * time.Second
)

// QueryRequest represents an SQL query request
type QueryRequest struct {
	SQL string `json:"sql"`
}

// QueryResponse represents a query response
type QueryResponse struct {
	Success       bool                     `json:"success"`
	Rows          []map[string]interface{} `json:"rows,omitempty"`
	RowCount      int                      `json:"rowCount,omitempty"`
	ExecutionTime float64                  `json:"executionTime,omitempty"`
	Error         *ErrorDetail             `json:"error,omitempty"`
}

// ErrorDetail represents error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Helper function to execute a query
func executeQuery(sql string) (*QueryResponse, error) {
	reqBody, err := json.Marshal(QueryRequest{SQL: sql})
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: testTimeout}
	resp, err := client.Post(testAPIURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return nil, err
	}

	return &queryResp, nil
}

// checkServerReady verifies the server is accessible
func checkServerReady(t testing.TB) {
	resp, err := executeQuery("SELECT 1")
	if err != nil {
		t.Skipf("Server not ready: %v\nPlease start the server: ./vibe serve", err)
	}
	if !resp.Success {
		t.Skipf("Server not healthy: %+v", resp.Error)
	}
}

// BenchmarkSimpleSelect benchmarks a simple SELECT 1 query (target: <10ms)
func BenchmarkSimpleSelect(b *testing.B) {
	checkServerReady(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := executeQuery("SELECT 1 as test")
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
		if !resp.Success {
			b.Fatalf("Query returned error: %+v", resp.Error)
		}
	}
}

// BenchmarkSelectWithWhere benchmarks SELECT with WHERE clause
func BenchmarkSelectWithWhere(b *testing.B) {
	checkServerReady(b)

	// Setup: Create test table
	tableName := fmt.Sprintf("perf_test_%d", time.Now().Unix())
	_, err := executeQuery(fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			value INTEGER,
			name TEXT
		)
	`, tableName))
	if err != nil {
		b.Fatalf("Setup failed: %v", err)
	}
	defer executeQuery(fmt.Sprintf("DROP TABLE %s", tableName))

	// Insert test data
	for i := 1; i <= 100; i++ {
		_, err := executeQuery(fmt.Sprintf("INSERT INTO %s (value, name) VALUES (%d, 'test%d')", tableName, i, i))
		if err != nil {
			b.Fatalf("Data insert failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sql := fmt.Sprintf("SELECT * FROM %s WHERE value = %d", tableName, (i%100)+1)
		resp, err := executeQuery(sql)
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
		if !resp.Success {
			b.Fatalf("Query returned error: %+v", resp.Error)
		}
	}
}

// BenchmarkJSONBFieldAccess benchmarks JSONB field access with -> operator
func BenchmarkJSONBFieldAccess(b *testing.B) {
	checkServerReady(b)

	tableName := fmt.Sprintf("perf_jsonb_%d", time.Now().Unix())
	_, err := executeQuery(fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			data JSONB
		)
	`, tableName))
	if err != nil {
		b.Fatalf("Setup failed: %v", err)
	}
	defer executeQuery(fmt.Sprintf("DROP TABLE %s", tableName))

	// Insert JSONB test data
	_, err = executeQuery(fmt.Sprintf(`
		INSERT INTO %s (data) VALUES 
		('{"name": "test1", "value": 100, "nested": {"key": "value1"}}'),
		('{"name": "test2", "value": 200, "nested": {"key": "value2"}}'),
		('{"name": "test3", "value": 300, "nested": {"key": "value3"}}')
	`, tableName))
	if err != nil {
		b.Fatalf("Data insert failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := executeQuery(fmt.Sprintf("SELECT data->'name' as name FROM %s", tableName))
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
		if !resp.Success {
			b.Fatalf("Query returned error: %+v", resp.Error)
		}
	}
}

// BenchmarkJSONBTextExtraction benchmarks JSONB text extraction with ->> operator
func BenchmarkJSONBTextExtraction(b *testing.B) {
	checkServerReady(b)

	tableName := fmt.Sprintf("perf_jsonb_text_%d", time.Now().Unix())
	_, err := executeQuery(fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			data JSONB
		)
	`, tableName))
	if err != nil {
		b.Fatalf("Setup failed: %v", err)
	}
	defer executeQuery(fmt.Sprintf("DROP TABLE %s", tableName))

	_, err = executeQuery(fmt.Sprintf(`
		INSERT INTO %s (data) VALUES 
		('{"name": "Alice", "city": "NYC"}'),
		('{"name": "Bob", "city": "LA"}'),
		('{"name": "Charlie", "city": "SF"}')
	`, tableName))
	if err != nil {
		b.Fatalf("Data insert failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := executeQuery(fmt.Sprintf("SELECT data->>'name' as name FROM %s WHERE data->>'city' = 'NYC'", tableName))
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
		if !resp.Success {
			b.Fatalf("Query returned error: %+v", resp.Error)
		}
	}
}

// BenchmarkInsert benchmarks INSERT operations
func BenchmarkInsert(b *testing.B) {
	checkServerReady(b)

	tableName := fmt.Sprintf("perf_insert_%d", time.Now().Unix())
	_, err := executeQuery(fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			value INTEGER
		)
	`, tableName))
	if err != nil {
		b.Fatalf("Setup failed: %v", err)
	}
	defer executeQuery(fmt.Sprintf("DROP TABLE %s", tableName))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := executeQuery(fmt.Sprintf("INSERT INTO %s (value) VALUES (%d)", tableName, i))
		if err != nil {
			b.Fatalf("Insert failed: %v", err)
		}
		if !resp.Success {
			b.Fatalf("Insert returned error: %+v", resp.Error)
		}
	}
}

// BenchmarkUpdate benchmarks UPDATE operations
func BenchmarkUpdate(b *testing.B) {
	checkServerReady(b)

	tableName := fmt.Sprintf("perf_update_%d", time.Now().Unix())
	_, err := executeQuery(fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			value INTEGER
		)
	`, tableName))
	if err != nil {
		b.Fatalf("Setup failed: %v", err)
	}
	defer executeQuery(fmt.Sprintf("DROP TABLE %s", tableName))

	// Insert initial data
	for i := 1; i <= 100; i++ {
		_, err := executeQuery(fmt.Sprintf("INSERT INTO %s (value) VALUES (%d)", tableName, i))
		if err != nil {
			b.Fatalf("Data insert failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := (i % 100) + 1
		resp, err := executeQuery(fmt.Sprintf("UPDATE %s SET value = %d WHERE id = %d", tableName, i*10, id))
		if err != nil {
			b.Fatalf("Update failed: %v", err)
		}
		if !resp.Success {
			b.Fatalf("Update returned error: %+v", resp.Error)
		}
	}
}

// TestLoadSequential tests 100 sequential queries (load test)
func TestLoadSequential(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	checkServerReady(t)

	const queryCount = 100
	start := time.Now()

	for i := 0; i < queryCount; i++ {
		resp, err := executeQuery("SELECT 1 as test")
		if err != nil {
			t.Fatalf("Query %d failed: %v", i+1, err)
		}
		if !resp.Success {
			t.Fatalf("Query %d returned error: %+v", i+1, resp.Error)
		}
	}

	duration := time.Since(start)
	avgTime := duration / queryCount

	t.Logf("Sequential load test completed:")
	t.Logf("  Total queries: %d", queryCount)
	t.Logf("  Total time: %v", duration)
	t.Logf("  Average time per query: %v", avgTime)
	t.Logf("  Queries per second: %.2f", float64(queryCount)/duration.Seconds())

	// Verify performance target (should be well under 10ms per query on average)
	if avgTime > 50*time.Millisecond {
		t.Errorf("Average query time %v exceeds 50ms threshold", avgTime)
	}
}

// TestLoadConcurrent tests concurrent queries (respecting 2 connection limit)
func TestLoadConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	checkServerReady(t)

	const concurrency = 2 // Match server's max connections
	const queriesPerWorker = 50
	const totalQueries = concurrency * queriesPerWorker

	var wg sync.WaitGroup
	errors := make(chan error, totalQueries)
	start := time.Now()

	for worker := 0; worker < concurrency; worker++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < queriesPerWorker; i++ {
				resp, err := executeQuery("SELECT 1 as test")
				if err != nil {
					errors <- fmt.Errorf("worker %d, query %d: %w", workerID, i+1, err)
					return
				}
				if !resp.Success {
					errors <- fmt.Errorf("worker %d, query %d: %+v", workerID, i+1, resp.Error)
					return
				}
			}
		}(worker)
	}

	wg.Wait()
	close(errors)
	duration := time.Since(start)

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Error(err)
		errorCount++
	}

	if errorCount > 0 {
		t.Fatalf("Concurrent load test failed with %d errors", errorCount)
	}

	avgTime := duration / totalQueries

	t.Logf("Concurrent load test completed:")
	t.Logf("  Concurrency: %d workers", concurrency)
	t.Logf("  Total queries: %d", totalQueries)
	t.Logf("  Total time: %v", duration)
	t.Logf("  Average time per query: %v", avgTime)
	t.Logf("  Throughput: %.2f queries/sec", float64(totalQueries)/duration.Seconds())
}

// TestMemoryUsage monitors memory usage during query execution
func TestMemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory test in short mode")
	}

	checkServerReady(t)

	tableName := fmt.Sprintf("perf_memory_%d", time.Now().Unix())
	_, err := executeQuery(fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			data TEXT
		)
	`, tableName))
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	defer executeQuery(fmt.Sprintf("DROP TABLE %s", tableName))

	// Insert test data
	for i := 1; i <= 100; i++ {
		_, err := executeQuery(fmt.Sprintf("INSERT INTO %s (data) VALUES ('test data %d')", tableName, i))
		if err != nil {
			t.Fatalf("Data insert failed: %v", err)
		}
	}

	// Record initial memory
	runtime.GC() // Force GC to get accurate baseline
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Execute many queries
	const queryCount = 1000
	for i := 0; i < queryCount; i++ {
		resp, err := executeQuery(fmt.Sprintf("SELECT * FROM %s WHERE id = %d", tableName, (i%100)+1))
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if !resp.Success {
			t.Fatalf("Query returned error: %+v", resp.Error)
		}
	}

	// Record final memory
	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	allocDiff := int64(m2.Alloc) - int64(m1.Alloc)
	totalAllocDiff := int64(m2.TotalAlloc) - int64(m1.TotalAlloc)

	t.Logf("Memory usage after %d queries:", queryCount)
	t.Logf("  Current alloc delta: %d bytes (%.2f MB)", allocDiff, float64(allocDiff)/(1024*1024))
	t.Logf("  Total alloc delta: %d bytes (%.2f MB)", totalAllocDiff, float64(totalAllocDiff)/(1024*1024))
	t.Logf("  Allocations per query: %.2f bytes", float64(totalAllocDiff)/float64(queryCount))

	// Check for memory leaks (current alloc should not grow excessively)
	// Allow up to 10MB growth for client-side state
	if allocDiff > 10*1024*1024 {
		t.Errorf("Potential memory leak detected: current alloc grew by %.2f MB", float64(allocDiff)/(1024*1024))
	}
}

// TestQueryTimeout verifies query timeout is enforced correctly (5s ± 100ms)
func TestQueryTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	checkServerReady(t)

	// Test query that should timeout (6 seconds > 5 second limit)
	start := time.Now()
	resp, err := executeQuery("SELECT pg_sleep(6)")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.Success {
		t.Error("Expected query to timeout, but it succeeded")
	}

	if resp.Error.Code != "QUERY_TIMEOUT" {
		t.Errorf("Expected error code QUERY_TIMEOUT, got %s", resp.Error.Code)
	}

	// Verify timeout occurred around 5 seconds (±100ms tolerance)
	expectedTimeout := 5 * time.Second
	if duration < expectedTimeout-100*time.Millisecond || duration > expectedTimeout+100*time.Millisecond {
		t.Errorf("Timeout duration %v outside acceptable range (5s ± 100ms)", duration)
	}

	t.Logf("Query timeout enforced correctly:")
	t.Logf("  Duration: %v", duration)
	t.Logf("  Expected: 5s ± 100ms")
	t.Logf("  Error code: %s", resp.Error.Code)
}

// TestStartupTime measures cold start time (target: <2s)
// Note: This test cannot programmatically start the server without embedded binary support (Phase 4)
// For now, it documents the manual test procedure
func TestStartupTime(t *testing.T) {
	t.Skip("Startup time test requires manual measurement until Phase 4")

	// Manual test procedure:
	// 1. Stop any running vibe server
	// 2. Delete ./vibe-data directory (cold start)
	// 3. Run: time ./vibe serve
	// 4. Measure time until "VibeSQL ready" log message
	// 5. Verify: < 2 seconds total startup time

	t.Log("Manual startup time test:")
	t.Log("  1. Stop server: pkill vibe")
	t.Log("  2. Remove data: rm -rf ./vibe-data")
	t.Log("  3. Time startup: time ./vibe serve")
	t.Log("  4. Target: <2 seconds to 'VibeSQL ready'")
}

// TestMain provides setup and summary reporting
func TestMain(m *testing.M) {
	// Check if server is accessible
	resp, err := executeQuery("SELECT 1")
	if err != nil || !resp.Success {
		fmt.Fprintf(os.Stderr, "\n⚠️  VibeSQL server not accessible at %s\n\n", testAPIURL)
		fmt.Fprintf(os.Stderr, "Performance tests require a running VibeSQL server.\n\n")
		fmt.Fprintf(os.Stderr, "Setup:\n")
		fmt.Fprintf(os.Stderr, "  1. Ensure PostgreSQL is running\n")
		fmt.Fprintf(os.Stderr, "  2. Start VibeSQL: ./vibe serve\n")
		fmt.Fprintf(os.Stderr, "  3. Run benchmarks: go test ./Tests/performance/... -bench=.\n\n")
	}

	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}
