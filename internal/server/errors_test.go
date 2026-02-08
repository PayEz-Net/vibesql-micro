package server

import (
	"net/http"
	"testing"

	"github.com/vibesql/vibe/internal/postgres"
)

// TestGetHTTPStatusCode tests the HTTP status code mapping for all error codes
func TestGetHTTPStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		errorCode      string
		expectedStatus int
	}{
		// 400 Bad Request errors
		{
			name:           "INVALID_SQL returns 400",
			errorCode:      ErrorCodeInvalidSQL,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "MISSING_REQUIRED_FIELD returns 400",
			errorCode:      ErrorCodeMissingRequiredField,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "UNSAFE_QUERY returns 400",
			errorCode:      ErrorCodeUnsafeQuery,
			expectedStatus: http.StatusBadRequest,
		},
		// 408 Request Timeout errors
		{
			name:           "QUERY_TIMEOUT returns 408",
			errorCode:      ErrorCodeQueryTimeout,
			expectedStatus: http.StatusRequestTimeout,
		},
		// 413 Payload Too Large errors
		{
			name:           "QUERY_TOO_LARGE returns 413",
			errorCode:      ErrorCodeQueryTooLarge,
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "RESULT_TOO_LARGE returns 413",
			errorCode:      ErrorCodeResultTooLarge,
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "DOCUMENT_TOO_LARGE returns 413",
			errorCode:      ErrorCodeDocumentTooLarge,
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		// 500 Internal Server Error
		{
			name:           "INTERNAL_ERROR returns 500",
			errorCode:      ErrorCodeInternalError,
			expectedStatus: http.StatusInternalServerError,
		},
		// 503 Service Unavailable errors
		{
			name:           "SERVICE_UNAVAILABLE returns 503",
			errorCode:      ErrorCodeServiceUnavailable,
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "DATABASE_UNAVAILABLE returns 503",
			errorCode:      ErrorCodeDatabaseUnavailable,
			expectedStatus: http.StatusServiceUnavailable,
		},
		// Unknown error code
		{
			name:           "Unknown error code returns 500",
			errorCode:      "UNKNOWN_ERROR",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := GetHTTPStatusCode(tt.errorCode)
			if statusCode != tt.expectedStatus {
				t.Errorf("GetHTTPStatusCode(%s) = %d, want %d", tt.errorCode, statusCode, tt.expectedStatus)
			}
		})
	}
}

// TestNewMissingFieldError tests the NewMissingFieldError helper
func TestNewMissingFieldError(t *testing.T) {
	err := NewMissingFieldError("sql")

	if err.Code != ErrorCodeMissingRequiredField {
		t.Errorf("Expected error code %s, got %s", ErrorCodeMissingRequiredField, err.Code)
	}

	if err.Message != "Missing required field: sql" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	if err.Detail != "The request must include a 'sql' field" {
		t.Errorf("Unexpected detail: %s", err.Detail)
	}

	// Test HTTP status mapping
	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400, got %d", statusCode)
	}
}

// TestNewInvalidSQLError tests the NewInvalidSQLError helper
func TestNewInvalidSQLError(t *testing.T) {
	detail := "syntax error at position 5"
	err := NewInvalidSQLError(detail)

	if err.Code != ErrorCodeInvalidSQL {
		t.Errorf("Expected error code %s, got %s", ErrorCodeInvalidSQL, err.Code)
	}

	if err.Message != "Invalid SQL syntax" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	if err.Detail != detail {
		t.Errorf("Expected detail '%s', got '%s'", detail, err.Detail)
	}

	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP status 400, got %d", statusCode)
	}
}

// TestNewUnsafeQueryError tests the NewUnsafeQueryError helper
func TestNewUnsafeQueryError(t *testing.T) {
	tests := []struct {
		name      string
		queryType string
	}{
		{"UPDATE without WHERE", "UPDATE"},
		{"DELETE without WHERE", "DELETE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUnsafeQueryError(tt.queryType)

			if err.Code != ErrorCodeUnsafeQuery {
				t.Errorf("Expected error code %s, got %s", ErrorCodeUnsafeQuery, err.Code)
			}

			expectedMessage := tt.queryType + " without WHERE clause is not allowed"
			if err.Message != expectedMessage {
				t.Errorf("Expected message '%s', got '%s'", expectedMessage, err.Message)
			}

			statusCode := GetHTTPStatusCode(err.Code)
			if statusCode != http.StatusBadRequest {
				t.Errorf("Expected HTTP status 400, got %d", statusCode)
			}
		})
	}
}

// TestNewQueryTimeoutError tests the NewQueryTimeoutError helper
func TestNewQueryTimeoutError(t *testing.T) {
	err := NewQueryTimeoutError()

	if err.Code != ErrorCodeQueryTimeout {
		t.Errorf("Expected error code %s, got %s", ErrorCodeQueryTimeout, err.Code)
	}

	if err.Message != "Query execution timeout" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	if err.Detail != "Query exceeded the maximum execution time" {
		t.Errorf("Unexpected detail: %s", err.Detail)
	}

	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusRequestTimeout {
		t.Errorf("Expected HTTP status 408, got %d", statusCode)
	}
}

// TestNewQueryTooLargeError tests the NewQueryTooLargeError helper
func TestNewQueryTooLargeError(t *testing.T) {
	actualSize := 15000
	maxSize := 10240
	err := NewQueryTooLargeError(actualSize, maxSize)

	if err.Code != ErrorCodeQueryTooLarge {
		t.Errorf("Expected error code %s, got %s", ErrorCodeQueryTooLarge, err.Code)
	}

	if err.Message != "Query too large" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	expectedDetail := "Query size (15000 bytes) exceeds maximum allowed size (10240 bytes)"
	if err.Detail != expectedDetail {
		t.Errorf("Expected detail '%s', got '%s'", expectedDetail, err.Detail)
	}

	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected HTTP status 413, got %d", statusCode)
	}
}

// TestNewResultTooLargeError tests the NewResultTooLargeError helper
func TestNewResultTooLargeError(t *testing.T) {
	actualRows := 1500
	maxRows := 1000
	err := NewResultTooLargeError(actualRows, maxRows)

	if err.Code != ErrorCodeResultTooLarge {
		t.Errorf("Expected error code %s, got %s", ErrorCodeResultTooLarge, err.Code)
	}

	if err.Message != "Result set too large" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	expectedDetail := "Query returned 1500 rows, exceeding the maximum limit of 1000 rows"
	if err.Detail != expectedDetail {
		t.Errorf("Expected detail '%s', got '%s'", expectedDetail, err.Detail)
	}

	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected HTTP status 413, got %d", statusCode)
	}
}

// TestNewDocumentTooLargeError tests the NewDocumentTooLargeError helper
func TestNewDocumentTooLargeError(t *testing.T) {
	maxSizeBytes := 1048576 // 1MB
	err := NewDocumentTooLargeError(maxSizeBytes)

	if err.Code != ErrorCodeDocumentTooLarge {
		t.Errorf("Expected error code %s, got %s", ErrorCodeDocumentTooLarge, err.Code)
	}

	if err.Message != "Document too large" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	expectedDetail := "JSONB document exceeds maximum size of 1048576 bytes"
	if err.Detail != expectedDetail {
		t.Errorf("Expected detail '%s', got '%s'", expectedDetail, err.Detail)
	}

	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected HTTP status 413, got %d", statusCode)
	}
}

// TestNewInternalError tests the NewInternalError helper
func TestNewInternalError(t *testing.T) {
	detail := "unexpected database connection failure"
	err := NewInternalError(detail)

	if err.Code != ErrorCodeInternalError {
		t.Errorf("Expected error code %s, got %s", ErrorCodeInternalError, err.Code)
	}

	if err.Message != "An internal error occurred" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	if err.Detail != detail {
		t.Errorf("Expected detail '%s', got '%s'", detail, err.Detail)
	}

	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status 500, got %d", statusCode)
	}
}

// TestNewServiceUnavailableError tests the NewServiceUnavailableError helper
func TestNewServiceUnavailableError(t *testing.T) {
	reason := "server is shutting down"
	err := NewServiceUnavailableError(reason)

	if err.Code != ErrorCodeServiceUnavailable {
		t.Errorf("Expected error code %s, got %s", ErrorCodeServiceUnavailable, err.Code)
	}

	if err.Message != "Service unavailable" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	if err.Detail != reason {
		t.Errorf("Expected detail '%s', got '%s'", reason, err.Detail)
	}

	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected HTTP status 503, got %d", statusCode)
	}
}

// TestNewDatabaseUnavailableError tests the NewDatabaseUnavailableError helper
func TestNewDatabaseUnavailableError(t *testing.T) {
	reason := "connection pool exhausted"
	err := NewDatabaseUnavailableError(reason)

	if err.Code != ErrorCodeDatabaseUnavailable {
		t.Errorf("Expected error code %s, got %s", ErrorCodeDatabaseUnavailable, err.Code)
	}

	if err.Message != "Database unavailable" {
		t.Errorf("Unexpected message: %s", err.Message)
	}

	if err.Detail != reason {
		t.Errorf("Expected detail '%s', got '%s'", reason, err.Detail)
	}

	statusCode := GetHTTPStatusCode(err.Code)
	if statusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected HTTP status 503, got %d", statusCode)
	}
}

// TestHTTPErrorCodeMapping tests that all error codes have the correct HTTP status mapping
func TestHTTPErrorCodeMapping(t *testing.T) {
	// This test ensures all 10 error codes map to the correct HTTP status
	expectedMappings := map[string]int{
		ErrorCodeInvalidSQL:           400,
		ErrorCodeMissingRequiredField: 400,
		ErrorCodeUnsafeQuery:          400,
		ErrorCodeQueryTimeout:         408,
		ErrorCodeQueryTooLarge:        413,
		ErrorCodeResultTooLarge:       413,
		ErrorCodeDocumentTooLarge:     413,
		ErrorCodeInternalError:        500,
		ErrorCodeServiceUnavailable:   503,
		ErrorCodeDatabaseUnavailable:  503,
	}

	for errorCode, expectedStatus := range expectedMappings {
		t.Run("Mapping for "+errorCode, func(t *testing.T) {
			actualStatus := GetHTTPStatusCode(errorCode)
			if actualStatus != expectedStatus {
				t.Errorf("Error code %s: expected HTTP %d, got %d", errorCode, expectedStatus, actualStatus)
			}
		})
	}

	// Verify we have exactly 10 error codes
	if len(expectedMappings) != 10 {
		t.Errorf("Expected 10 error codes, found %d", len(expectedMappings))
	}
}

// TestValidateHTTPStatusMapping tests the validation function
func TestValidateHTTPStatusMapping(t *testing.T) {
	err := ValidateHTTPStatusMapping()
	if err != nil {
		t.Errorf("HTTP status mapping validation failed: %v", err)
	}
}

// TestErrorCodeConstants tests that error code constants match postgres package
func TestErrorCodeConstants(t *testing.T) {
	tests := []struct {
		name           string
		serverConstant string
		postgresConstant string
	}{
		{"INVALID_SQL", ErrorCodeInvalidSQL, postgres.ErrorCodeInvalidSQL},
		{"MISSING_REQUIRED_FIELD", ErrorCodeMissingRequiredField, postgres.ErrorCodeMissingRequiredField},
		{"UNSAFE_QUERY", ErrorCodeUnsafeQuery, postgres.ErrorCodeUnsafeQuery},
		{"QUERY_TIMEOUT", ErrorCodeQueryTimeout, postgres.ErrorCodeQueryTimeout},
		{"QUERY_TOO_LARGE", ErrorCodeQueryTooLarge, postgres.ErrorCodeQueryTooLarge},
		{"RESULT_TOO_LARGE", ErrorCodeResultTooLarge, postgres.ErrorCodeResultTooLarge},
		{"DOCUMENT_TOO_LARGE", ErrorCodeDocumentTooLarge, postgres.ErrorCodeDocumentTooLarge},
		{"INTERNAL_ERROR", ErrorCodeInternalError, postgres.ErrorCodeInternalError},
		{"SERVICE_UNAVAILABLE", ErrorCodeServiceUnavailable, postgres.ErrorCodeServiceUnavailable},
		{"DATABASE_UNAVAILABLE", ErrorCodeDatabaseUnavailable, postgres.ErrorCodeDatabaseUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.serverConstant != tt.postgresConstant {
				t.Errorf("Constant mismatch for %s: server=%s, postgres=%s",
					tt.name, tt.serverConstant, tt.postgresConstant)
			}
		})
	}
}

// TestErrorHelperReturnTypes tests that all helper functions return VibeError pointers
func TestErrorHelperReturnTypes(t *testing.T) {
	tests := []struct {
		name     string
		errFunc  func() *postgres.VibeError
	}{
		{
			name:    "NewMissingFieldError",
			errFunc: func() *postgres.VibeError { return NewMissingFieldError("test") },
		},
		{
			name:    "NewInvalidSQLError",
			errFunc: func() *postgres.VibeError { return NewInvalidSQLError("test") },
		},
		{
			name:    "NewUnsafeQueryError",
			errFunc: func() *postgres.VibeError { return NewUnsafeQueryError("UPDATE") },
		},
		{
			name:    "NewQueryTimeoutError",
			errFunc: func() *postgres.VibeError { return NewQueryTimeoutError() },
		},
		{
			name:    "NewQueryTooLargeError",
			errFunc: func() *postgres.VibeError { return NewQueryTooLargeError(100, 50) },
		},
		{
			name:    "NewResultTooLargeError",
			errFunc: func() *postgres.VibeError { return NewResultTooLargeError(2000, 1000) },
		},
		{
			name:    "NewDocumentTooLargeError",
			errFunc: func() *postgres.VibeError { return NewDocumentTooLargeError(1048576) },
		},
		{
			name:    "NewInternalError",
			errFunc: func() *postgres.VibeError { return NewInternalError("test") },
		},
		{
			name:    "NewServiceUnavailableError",
			errFunc: func() *postgres.VibeError { return NewServiceUnavailableError("test") },
		},
		{
			name:    "NewDatabaseUnavailableError",
			errFunc: func() *postgres.VibeError { return NewDatabaseUnavailableError("test") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()
			if err == nil {
				t.Errorf("%s returned nil", tt.name)
			}
			if err.Code == "" {
				t.Errorf("%s returned error with empty Code", tt.name)
			}
			if err.Message == "" {
				t.Errorf("%s returned error with empty Message", tt.name)
			}
		})
	}
}

// TestAllHTTPStatusCodesInRange tests that all HTTP status codes are valid
func TestAllHTTPStatusCodesInRange(t *testing.T) {
	validStatuses := map[int]bool{
		400: true, // Bad Request
		408: true, // Request Timeout
		413: true, // Payload Too Large
		500: true, // Internal Server Error
		503: true, // Service Unavailable
	}

	for errorCode := range HTTPErrorCodeMapping {
		statusCode := GetHTTPStatusCode(errorCode)
		if !validStatuses[statusCode] {
			t.Errorf("Error code %s maps to invalid HTTP status %d", errorCode, statusCode)
		}
		if statusCode < 400 || statusCode >= 600 {
			t.Errorf("Error code %s maps to out-of-range HTTP status %d", errorCode, statusCode)
		}
	}
}

func BenchmarkGetHTTPStatusCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetHTTPStatusCode(ErrorCodeInvalidSQL)
	}
}
