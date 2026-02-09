package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/lib/pq"
)

func TestTranslateError_NilError(t *testing.T) {
	result := TranslateError(nil)
	if result != nil {
		t.Errorf("Expected nil for nil error, got %v", result)
	}
}

func TestTranslateError_VibeError(t *testing.T) {
	original := NewVibeError(ErrorCodeInvalidSQL, "Test message", "Test detail")
	result := TranslateError(original)
	
	if result.Code != ErrorCodeInvalidSQL {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidSQL, result.Code)
	}
	if result.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %s", result.Message)
	}
}

func TestTranslateError_SyntaxError(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "42601",
		Message: "syntax error at or near \"SELCT\"",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeInvalidSQL {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidSQL, result.Code)
	}
	if result.Message != "Invalid SQL syntax" {
		t.Errorf("Expected message 'Invalid SQL syntax', got %s", result.Message)
	}
}

func TestTranslateError_UndefinedColumn(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "42703",
		Message: "column \"foo\" does not exist",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeInvalidSQL {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidSQL, result.Code)
	}
}

func TestTranslateError_UndefinedTable(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "42P01",
		Message: "relation \"users\" does not exist",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeInvalidSQL {
		t.Errorf("Expected code %s, got %s", ErrorCodeInvalidSQL, result.Code)
	}
}

func TestTranslateError_QueryCanceled(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "57014",
		Message: "canceling statement due to user request",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeQueryTimeout {
		t.Errorf("Expected code %s, got %s", ErrorCodeQueryTimeout, result.Code)
	}
	if result.Message != "Query execution timeout" {
		t.Errorf("Expected message 'Query execution timeout', got %s", result.Message)
	}
}

func TestTranslateError_ConfigurationLimitExceeded(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "53400",
		Message: "configuration limit exceeded",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeDatabaseUnavailable {
		t.Errorf("Expected code %s, got %s", ErrorCodeDatabaseUnavailable, result.Code)
	}
}

func TestTranslateError_TooManyConnections(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "53300",
		Message: "too many connections",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeDatabaseUnavailable {
		t.Errorf("Expected code %s, got %s", ErrorCodeDatabaseUnavailable, result.Code)
	}
}

func TestTranslateError_ConnectionFailure(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "08006",
		Message: "connection failure",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeDatabaseUnavailable {
		t.Errorf("Expected code %s, got %s", ErrorCodeDatabaseUnavailable, result.Code)
	}
}

func TestTranslateError_ProgramLimitExceeded(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "54000",
		Message: "program limit exceeded",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeDocumentTooLarge {
		t.Errorf("Expected code %s, got %s", ErrorCodeDocumentTooLarge, result.Code)
	}
}

func TestTranslateError_UnknownSQLSTATE(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "99999",
		Message: "unknown error",
	}
	
	result := TranslateError(pqErr)
	
	if result.Code != ErrorCodeInternalError {
		t.Errorf("Expected code %s for unknown SQLSTATE, got %s", ErrorCodeInternalError, result.Code)
	}
}

func TestTranslateError_GenericError(t *testing.T) {
	genericErr := errors.New("generic error")
	
	result := TranslateError(genericErr)
	
	if result.Code != ErrorCodeInternalError {
		t.Errorf("Expected code %s for generic error, got %s", ErrorCodeInternalError, result.Code)
	}
	if result.Detail != "generic error" {
		t.Errorf("Expected detail 'generic error', got %s", result.Detail)
	}
}

func TestTranslateError_ContextDeadlineExceeded(t *testing.T) {
	err := context.DeadlineExceeded
	
	result := TranslateError(err)
	
	if result.Code != ErrorCodeQueryTimeout {
		t.Errorf("Expected code %s for context.DeadlineExceeded, got %s", ErrorCodeQueryTimeout, result.Code)
	}
	if result.Message != "Query execution timeout" {
		t.Errorf("Expected message 'Query execution timeout', got %s", result.Message)
	}
}

func TestTranslateError_ContextCanceled(t *testing.T) {
	err := context.Canceled
	
	result := TranslateError(err)
	
	if result.Code != ErrorCodeQueryTimeout {
		t.Errorf("Expected code %s for context.Canceled, got %s", ErrorCodeQueryTimeout, result.Code)
	}
	if result.Message != "Query execution canceled" {
		t.Errorf("Expected message 'Query execution canceled', got %s", result.Message)
	}
}

func TestTranslateError_WithDetailAndHint(t *testing.T) {
	pqErr := &pq.Error{
		Code:    "42601",
		Message: "syntax error",
		Detail:  "Unexpected token",
		Hint:    "Check your SQL syntax",
	}
	
	result := TranslateError(pqErr)
	
	if result.Detail == "" {
		t.Error("Expected non-empty detail")
	}
	
	// Verify detail includes PostgreSQL message, detail, and hint
	expectedSubstrings := []string{"syntax error", "Unexpected token", "Check your SQL syntax"}
	for _, substr := range expectedSubstrings {
		if !contains(result.Detail, substr) {
			t.Errorf("Expected detail to contain '%s', got: %s", substr, result.Detail)
		}
	}
}

func TestGetHTTPStatusCode_AllErrorCodes(t *testing.T) {
	tests := []struct {
		errorCode      string
		expectedStatus int
	}{
		{ErrorCodeInvalidSQL, 400},
		{ErrorCodeMissingRequiredField, 400},
		{ErrorCodeUnsafeQuery, 400},
		{ErrorCodeQueryTimeout, 408},
		{ErrorCodeQueryTooLarge, 413},
		{ErrorCodeResultTooLarge, 413},
		{ErrorCodeDocumentTooLarge, 413},
		{ErrorCodeInternalError, 500},
		{ErrorCodeServiceUnavailable, 503},
		{ErrorCodeDatabaseUnavailable, 503},
		{"UNKNOWN_CODE", 500}, // Default to 500
	}
	
	for _, tt := range tests {
		t.Run(tt.errorCode, func(t *testing.T) {
			status := GetHTTPStatusCode(tt.errorCode)
			if status != tt.expectedStatus {
				t.Errorf("Expected status %d for code %s, got %d", 
					tt.expectedStatus, tt.errorCode, status)
			}
		})
	}
}

func TestVibeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		vibeErr  *VibeError
		expected string
	}{
		{
			name: "With detail",
			vibeErr: &VibeError{
				Code:    ErrorCodeInvalidSQL,
				Message: "Invalid SQL",
				Detail:  "Additional info",
			},
			expected: "INVALID_SQL: Invalid SQL (Additional info)",
		},
		{
			name: "Without detail",
			vibeErr: &VibeError{
				Code:    ErrorCodeQueryTimeout,
				Message: "Timeout",
				Detail:  "",
			},
			expected: "QUERY_TIMEOUT: Timeout",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.vibeErr.Error()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestNewVibeError(t *testing.T) {
	err := NewVibeError("TEST_CODE", "Test message", "Test detail")
	
	if err.Code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got %s", err.Code)
	}
	if err.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %s", err.Message)
	}
	if err.Detail != "Test detail" {
		t.Errorf("Expected detail 'Test detail', got %s", err.Detail)
	}
}

// Test all SQLSTATE mappings from the spec
func TestSQLStateMapping_AllCodes(t *testing.T) {
	tests := []struct {
		sqlState     string
		expectedCode string
	}{
		// Syntax errors
		{"42601", ErrorCodeInvalidSQL},
		{"42703", ErrorCodeInvalidSQL},
		{"42P01", ErrorCodeInvalidSQL},
		{"42P02", ErrorCodeInvalidSQL},
		{"42883", ErrorCodeInvalidSQL},
		{"42804", ErrorCodeInvalidSQL},
		
		// Query cancellation
		{"57014", ErrorCodeQueryTimeout},
		
		// Resource limits
		{"53000", ErrorCodeDatabaseUnavailable},
		{"53100", ErrorCodeDatabaseUnavailable},
		{"53200", ErrorCodeDatabaseUnavailable},
		{"53300", ErrorCodeDatabaseUnavailable},
		{"53400", ErrorCodeDatabaseUnavailable},
		
		// Connection errors
		{"08000", ErrorCodeDatabaseUnavailable},
		{"08003", ErrorCodeDatabaseUnavailable},
		{"08006", ErrorCodeDatabaseUnavailable},
		{"08001", ErrorCodeDatabaseUnavailable},
		{"08004", ErrorCodeDatabaseUnavailable},
		
		// Document size errors
		{"54000", ErrorCodeDocumentTooLarge},
		{"54001", ErrorCodeDocumentTooLarge},
	}
	
	for _, tt := range tests {
		t.Run(tt.sqlState, func(t *testing.T) {
			pqErr := &pq.Error{
				Code:    pq.ErrorCode(tt.sqlState),
				Message: "test error",
			}
			
			result := TranslateError(pqErr)
			
			if result.Code != tt.expectedCode {
				t.Errorf("SQLSTATE %s: expected code %s, got %s", 
					tt.sqlState, tt.expectedCode, result.Code)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
