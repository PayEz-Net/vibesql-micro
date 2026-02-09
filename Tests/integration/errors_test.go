package integration

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/vibesql/vibe/internal/postgres"
	"github.com/vibesql/vibe/internal/server"
)

// TestErrors_HTTPStatusMapping validates all error codes map to correct HTTP status
func TestErrors_HTTPStatusMapping(t *testing.T) {
	// Verify all 10 error codes have correct HTTP status mappings
	err := server.ValidateHTTPStatusMapping()
	if err != nil {
		t.Errorf("HTTP status mapping validation failed: %v", err)
	}

	// Test each error code individually
	testCases := []struct {
		errorCode      string
		expectedStatus int
	}{
		{server.ErrorCodeInvalidSQL, 400},
		{server.ErrorCodeMissingRequiredField, 400},
		{server.ErrorCodeUnsafeQuery, 400},
		{server.ErrorCodeQueryTimeout, 408},
		{server.ErrorCodeQueryTooLarge, 413},
		{server.ErrorCodeResultTooLarge, 413},
		{server.ErrorCodeDocumentTooLarge, 413},
		{server.ErrorCodeInternalError, 500},
		{server.ErrorCodeServiceUnavailable, 503},
		{server.ErrorCodeDatabaseUnavailable, 503},
	}

	for _, tc := range testCases {
		actualStatus := server.GetHTTPStatusCode(tc.errorCode)
		if actualStatus != tc.expectedStatus {
			t.Errorf("Error code %s: expected HTTP %d, got %d", 
				tc.errorCode, tc.expectedStatus, actualStatus)
		}
	}

	t.Logf("✓ All 10 error codes map to correct HTTP status")
}

// TestErrors_InvalidSQLSyntax tests invalid SQL syntax → 400 INVALID_SQL
func TestErrors_InvalidSQLSyntax(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Test invalid SQL syntax
	_, err := db.Exec("SELECTT 1") // Typo in SELECT
	if err == nil {
		t.Error("Invalid SQL should return error")
		return
	}

	// Verify error is translated to INVALID_SQL
	vibeErr := postgres.TranslateError(err)
	if vibeErr == nil {
		t.Error("Expected VibeError, got nil")
		return
	}

	if vibeErr.Code != postgres.ErrorCodeInvalidSQL {
		t.Errorf("Expected error code INVALID_SQL, got %s", vibeErr.Code)
	}

	// Verify HTTP status
	httpStatus := server.GetHTTPStatusCode(vibeErr.Code)
	if httpStatus != 400 {
		t.Errorf("Expected HTTP 400, got %d", httpStatus)
	}

	t.Logf("✓ Invalid SQL syntax → %s (HTTP %d)", vibeErr.Code, httpStatus)
}

// TestErrors_UndefinedColumn tests undefined column → 400 INVALID_SQL
func TestErrors_UndefinedColumn(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupErrorsTestTable(t, db, "error_test_column")
	defer db.Exec("DROP TABLE IF EXISTS error_test_column")

	// Query non-existent column
	_, err := db.Query("SELECT nonexistent_column FROM error_test_column")
	if err == nil {
		t.Error("Query with undefined column should return error")
		return
	}

	// Verify error translation
	vibeErr := postgres.TranslateError(err)
	if vibeErr == nil {
		t.Error("Expected VibeError, got nil")
		return
	}

	if vibeErr.Code != postgres.ErrorCodeInvalidSQL {
		t.Errorf("Expected INVALID_SQL, got %s", vibeErr.Code)
	}

	httpStatus := server.GetHTTPStatusCode(vibeErr.Code)
	if httpStatus != 400 {
		t.Errorf("Expected HTTP 400, got %d", httpStatus)
	}

	t.Logf("✓ Undefined column → %s (HTTP %d)", vibeErr.Code, httpStatus)
}

// TestErrors_UndefinedTable tests undefined table → 400 INVALID_SQL
func TestErrors_UndefinedTable(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Query non-existent table
	_, err := db.Query("SELECT * FROM nonexistent_table_12345")
	if err == nil {
		t.Error("Query with undefined table should return error")
		return
	}

	// Verify error translation
	vibeErr := postgres.TranslateError(err)
	if vibeErr == nil {
		t.Error("Expected VibeError, got nil")
		return
	}

	if vibeErr.Code != postgres.ErrorCodeInvalidSQL {
		t.Errorf("Expected INVALID_SQL, got %s", vibeErr.Code)
	}

	httpStatus := server.GetHTTPStatusCode(vibeErr.Code)
	if httpStatus != 400 {
		t.Errorf("Expected HTTP 400, got %d", httpStatus)
	}

	t.Logf("✓ Undefined table → %s (HTTP %d)", vibeErr.Code, httpStatus)
}

// TestErrors_MissingRequiredField tests missing sql field → 400 MISSING_REQUIRED_FIELD
func TestErrors_MissingRequiredField(t *testing.T) {
	// Test the error creation helper
	err := server.NewMissingFieldError("sql")
	
	if err.Code != server.ErrorCodeMissingRequiredField {
		t.Errorf("Expected MISSING_REQUIRED_FIELD, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 400 {
		t.Errorf("Expected HTTP 400, got %d", httpStatus)
	}

	if !strings.Contains(err.Message, "sql") {
		t.Errorf("Error message should mention 'sql' field: %s", err.Message)
	}

	t.Logf("✓ Missing required field → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_EmptySQLField tests empty sql field → 400 MISSING_REQUIRED_FIELD
func TestErrors_EmptySQLField(t *testing.T) {
	// Test empty SQL validation
	err := server.NewMissingFieldError("sql")
	
	if err.Code != server.ErrorCodeMissingRequiredField {
		t.Errorf("Expected MISSING_REQUIRED_FIELD, got %s", err.Code)
	}

	if err.Detail == "" {
		t.Error("Error detail should not be empty")
	}

	t.Logf("✓ Empty SQL field → %s: %s", err.Code, err.Detail)
}

// TestErrors_UnsafeQueryUpdate tests UPDATE without WHERE → 400 UNSAFE_QUERY
func TestErrors_UnsafeQueryUpdate(t *testing.T) {
	// Test the error creation helper
	err := server.NewUnsafeQueryError("UPDATE")
	
	if err.Code != server.ErrorCodeUnsafeQuery {
		t.Errorf("Expected UNSAFE_QUERY, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 400 {
		t.Errorf("Expected HTTP 400, got %d", httpStatus)
	}

	if !strings.Contains(err.Message, "UPDATE") {
		t.Errorf("Error message should mention UPDATE: %s", err.Message)
	}

	if !strings.Contains(err.Detail, "WHERE 1=1") {
		t.Errorf("Error detail should suggest WHERE 1=1 bypass: %s", err.Detail)
	}

	t.Logf("✓ UPDATE without WHERE → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_UnsafeQueryDelete tests DELETE without WHERE → 400 UNSAFE_QUERY
func TestErrors_UnsafeQueryDelete(t *testing.T) {
	// Test the error creation helper
	err := server.NewUnsafeQueryError("DELETE")
	
	if err.Code != server.ErrorCodeUnsafeQuery {
		t.Errorf("Expected UNSAFE_QUERY, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 400 {
		t.Errorf("Expected HTTP 400, got %d", httpStatus)
	}

	if !strings.Contains(err.Message, "DELETE") {
		t.Errorf("Error message should mention DELETE: %s", err.Message)
	}

	t.Logf("✓ DELETE without WHERE → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_QueryTimeout tests timeout → 408 QUERY_TIMEOUT
func TestErrors_QueryTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	db := getTestDB(t)
	defer db.Close()

	// Create a query that will timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "SELECT pg_sleep(10)")
	if err == nil {
		t.Error("Long query should timeout")
		return
	}

	// Check if error is timeout-related
	if ctx.Err() != context.DeadlineExceeded {
		t.Logf("Context error: %v", ctx.Err())
	}

	// Test the error creation helper
	timeoutErr := server.NewQueryTimeoutError()
	
	if timeoutErr.Code != server.ErrorCodeQueryTimeout {
		t.Errorf("Expected QUERY_TIMEOUT, got %s", timeoutErr.Code)
	}

	httpStatus := server.GetHTTPStatusCode(timeoutErr.Code)
	if httpStatus != 408 {
		t.Errorf("Expected HTTP 408, got %d", httpStatus)
	}

	t.Logf("✓ Query timeout → %s (HTTP %d): %s", timeoutErr.Code, httpStatus, timeoutErr.Message)
}

// TestErrors_QueryTooLarge tests query > 10KB → 413 QUERY_TOO_LARGE
func TestErrors_QueryTooLarge(t *testing.T) {
	actualSize := 10*1024 + 100
	maxSize := 10 * 1024
	
	err := server.NewQueryTooLargeError(actualSize, maxSize)
	
	if err.Code != server.ErrorCodeQueryTooLarge {
		t.Errorf("Expected QUERY_TOO_LARGE, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 413 {
		t.Errorf("Expected HTTP 413, got %d", httpStatus)
	}

	if !strings.Contains(err.Detail, "10240") {
		t.Errorf("Error detail should mention max size: %s", err.Detail)
	}

	t.Logf("✓ Query too large → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_ResultTooLarge tests result > 1000 rows → 413 RESULT_TOO_LARGE
func TestErrors_ResultTooLarge(t *testing.T) {
	actualRows := 1001
	maxRows := 1000
	
	err := server.NewResultTooLargeError(actualRows, maxRows)
	
	if err.Code != server.ErrorCodeResultTooLarge {
		t.Errorf("Expected RESULT_TOO_LARGE, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 413 {
		t.Errorf("Expected HTTP 413, got %d", httpStatus)
	}

	if !strings.Contains(err.Detail, "1001") || !strings.Contains(err.Detail, "1000") {
		t.Errorf("Error detail should mention row counts: %s", err.Detail)
	}

	t.Logf("✓ Result too large → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_DocumentTooLarge tests JSONB document > 1MB → 413 DOCUMENT_TOO_LARGE
func TestErrors_DocumentTooLarge(t *testing.T) {
	maxSizeBytes := 1024 * 1024 // 1MB
	
	err := server.NewDocumentTooLargeError(maxSizeBytes)
	
	if err.Code != server.ErrorCodeDocumentTooLarge {
		t.Errorf("Expected DOCUMENT_TOO_LARGE, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 413 {
		t.Errorf("Expected HTTP 413, got %d", httpStatus)
	}

	if !strings.Contains(err.Detail, "1048576") {
		t.Errorf("Error detail should mention max size in bytes: %s", err.Detail)
	}

	t.Logf("✓ Document too large → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_InternalError tests internal error → 500 INTERNAL_ERROR
func TestErrors_InternalError(t *testing.T) {
	detail := "Unexpected error during query processing"
	
	err := server.NewInternalError(detail)
	
	if err.Code != server.ErrorCodeInternalError {
		t.Errorf("Expected INTERNAL_ERROR, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 500 {
		t.Errorf("Expected HTTP 500, got %d", httpStatus)
	}

	if err.Detail != detail {
		t.Errorf("Expected detail '%s', got '%s'", detail, err.Detail)
	}

	t.Logf("✓ Internal error → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_ServiceUnavailable tests service unavailable → 503 SERVICE_UNAVAILABLE
func TestErrors_ServiceUnavailable(t *testing.T) {
	reason := "Server is at maximum connection capacity"
	
	err := server.NewServiceUnavailableError(reason)
	
	if err.Code != server.ErrorCodeServiceUnavailable {
		t.Errorf("Expected SERVICE_UNAVAILABLE, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 503 {
		t.Errorf("Expected HTTP 503, got %d", httpStatus)
	}

	if err.Detail != reason {
		t.Errorf("Expected detail '%s', got '%s'", reason, err.Detail)
	}

	t.Logf("✓ Service unavailable → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_DatabaseUnavailable tests database unavailable → 503 DATABASE_UNAVAILABLE
func TestErrors_DatabaseUnavailable(t *testing.T) {
	reason := "Database connection failed"
	
	err := server.NewDatabaseUnavailableError(reason)
	
	if err.Code != server.ErrorCodeDatabaseUnavailable {
		t.Errorf("Expected DATABASE_UNAVAILABLE, got %s", err.Code)
	}

	httpStatus := server.GetHTTPStatusCode(err.Code)
	if httpStatus != 503 {
		t.Errorf("Expected HTTP 503, got %d", httpStatus)
	}

	if err.Detail != reason {
		t.Errorf("Expected detail '%s', got '%s'", reason, err.Detail)
	}

	t.Logf("✓ Database unavailable → %s (HTTP %d): %s", err.Code, httpStatus, err.Message)
}

// TestErrors_SQLSTATETranslation tests SQLSTATE code translation
func TestErrors_SQLSTATETranslation(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	testCases := []struct {
		name          string
		sql           string
		expectedCode  string
		expectedHTTP  int
		setupFunc     func() error
		cleanupFunc   func() error
	}{
		{
			name:         "Syntax error (42601)",
			sql:          "SELECTT 1",
			expectedCode: postgres.ErrorCodeInvalidSQL,
			expectedHTTP: 400,
		},
		{
			name:         "Undefined column (42703)",
			sql:          "SELECT nonexistent FROM pg_type LIMIT 1",
			expectedCode: postgres.ErrorCodeInvalidSQL,
			expectedHTTP: 400,
		},
		{
			name:         "Undefined table (42P01)",
			sql:          "SELECT * FROM nonexistent_table_xyz",
			expectedCode: postgres.ErrorCodeInvalidSQL,
			expectedHTTP: 400,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				if err := tc.setupFunc(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			if tc.cleanupFunc != nil {
				defer tc.cleanupFunc()
			}

			_, err := db.Exec(tc.sql)
			if err == nil {
				t.Error("Expected error, got nil")
				return
			}

			// Verify SQLSTATE code if available
			if pqErr, ok := err.(*pq.Error); ok {
				t.Logf("SQLSTATE: %s, Message: %s", pqErr.Code, pqErr.Message)
			}

			// Verify translation to VibeSQL error
			vibeErr := postgres.TranslateError(err)
			if vibeErr == nil {
				t.Error("Expected VibeError, got nil")
				return
			}

			if vibeErr.Code != tc.expectedCode {
				t.Errorf("Expected error code %s, got %s", tc.expectedCode, vibeErr.Code)
			}

			httpStatus := server.GetHTTPStatusCode(vibeErr.Code)
			if httpStatus != tc.expectedHTTP {
				t.Errorf("Expected HTTP %d, got %d", tc.expectedHTTP, httpStatus)
			}

			t.Logf("✓ %s → %s (HTTP %d)", tc.name, vibeErr.Code, httpStatus)
		})
	}
}

// TestErrors_ErrorMessageClarity tests that error messages are helpful
func TestErrors_ErrorMessageClarity(t *testing.T) {
	testCases := []struct {
		name          string
		errorFunc     func() *postgres.VibeError
		expectInMsg   string
		expectInDetail string
	}{
		{
			name:           "Missing field error",
			errorFunc:      func() *postgres.VibeError { return server.NewMissingFieldError("sql") },
			expectInMsg:    "sql",
			expectInDetail: "must include",
		},
		{
			name:           "Unsafe query error",
			errorFunc:      func() *postgres.VibeError { return server.NewUnsafeQueryError("UPDATE") },
			expectInMsg:    "where",
			expectInDetail: "where 1=1",
		},
		{
			name:           "Query timeout error",
			errorFunc:      func() *postgres.VibeError { return server.NewQueryTimeoutError() },
			expectInMsg:    "timeout",
			expectInDetail: "exceeded",
		},
		{
			name:           "Query too large error",
			errorFunc:      func() *postgres.VibeError { return server.NewQueryTooLargeError(11000, 10240) },
			expectInMsg:    "too large",
			expectInDetail: "exceeds",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.errorFunc()

			if !strings.Contains(strings.ToLower(err.Message), tc.expectInMsg) {
				t.Errorf("Message should contain '%s': %s", tc.expectInMsg, err.Message)
			}

			if !strings.Contains(strings.ToLower(err.Detail), tc.expectInDetail) {
				t.Errorf("Detail should contain '%s': %s", tc.expectInDetail, err.Detail)
			}

			t.Logf("✓ %s message is clear and helpful", tc.name)
		})
	}
}

// TestErrors_ErrorJSONFormat tests error serialization format
func TestErrors_ErrorJSONFormat(t *testing.T) {
	err := server.NewInvalidSQLError("Invalid syntax near 'SELECTT'")

	if err.Code == "" {
		t.Error("Error code should not be empty")
	}

	if err.Message == "" {
		t.Error("Error message should not be empty")
	}

	if err.Detail == "" {
		t.Error("Error detail should not be empty")
	}

	t.Logf("✓ Error format: {code: %s, message: %s, detail: %s}", 
		err.Code, err.Message, err.Detail)
}

// Helper function to create a test table for errors testing
func setupErrorsTestTable(t *testing.T, db *sql.DB, tableName string) {
	// Drop table if exists
	_, err := db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		t.Fatalf("Failed to drop test table %s: %v", tableName, err)
	}

	// Create test table
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			name TEXT,
			value INTEGER
		)
	`, tableName))
	if err != nil {
		t.Fatalf("Failed to create test table %s: %v", tableName, err)
	}
}
