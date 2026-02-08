// Package server provides HTTP-specific error handling for VibeSQL.
// It wraps error codes from the postgres package and provides
// convenient helper functions for creating common HTTP errors with
// appropriate status codes and user-friendly messages.
package server

import (
	"fmt"
	"net/http"

	"github.com/vibesql/vibe/internal/postgres"
)

// Error code constants (imported from postgres package for convenience)
const (
	ErrorCodeInvalidSQL           = postgres.ErrorCodeInvalidSQL
	ErrorCodeMissingRequiredField = postgres.ErrorCodeMissingRequiredField
	ErrorCodeUnsafeQuery          = postgres.ErrorCodeUnsafeQuery
	ErrorCodeQueryTimeout         = postgres.ErrorCodeQueryTimeout
	ErrorCodeQueryTooLarge        = postgres.ErrorCodeQueryTooLarge
	ErrorCodeResultTooLarge       = postgres.ErrorCodeResultTooLarge
	ErrorCodeDocumentTooLarge     = postgres.ErrorCodeDocumentTooLarge
	ErrorCodeInternalError        = postgres.ErrorCodeInternalError
	ErrorCodeServiceUnavailable   = postgres.ErrorCodeServiceUnavailable
	ErrorCodeDatabaseUnavailable  = postgres.ErrorCodeDatabaseUnavailable
)

// GetHTTPStatusCode returns the HTTP status code for a given VibeSQL error code
func GetHTTPStatusCode(errorCode string) int {
	return postgres.GetHTTPStatusCode(errorCode)
}

// Helper functions for creating common errors

// NewMissingFieldError creates an error for a missing required field
func NewMissingFieldError(fieldName string) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeMissingRequiredField,
		fmt.Sprintf("Missing required field: %s", fieldName),
		fmt.Sprintf("The request must include a '%s' field", fieldName),
	)
}

// NewInvalidSQLError creates an error for invalid SQL syntax
func NewInvalidSQLError(message string) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeInvalidSQL,
		"Invalid SQL syntax",
		message,
	)
}

// NewUnsafeQueryError creates an error for unsafe queries (UPDATE/DELETE without WHERE)
func NewUnsafeQueryError(queryType string) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeUnsafeQuery,
		fmt.Sprintf("%s without WHERE clause is not allowed", queryType),
		fmt.Sprintf("For safety, %s statements must include a WHERE clause. Use 'WHERE 1=1' to update/delete all rows.", queryType),
	)
}

// NewQueryTimeoutError creates an error for query timeout
func NewQueryTimeoutError() *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeQueryTimeout,
		"Query execution timeout",
		"Query exceeded the maximum execution time",
	)
}

// NewQueryTooLargeError creates an error for queries exceeding size limit
func NewQueryTooLargeError(actualSize, maxSize int) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeQueryTooLarge,
		"Query too large",
		fmt.Sprintf("Query size (%d bytes) exceeds maximum allowed size (%d bytes)", actualSize, maxSize),
	)
}

// NewResultTooLargeError creates an error for result sets exceeding row limit
func NewResultTooLargeError(actualRows, maxRows int) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeResultTooLarge,
		"Result set too large",
		fmt.Sprintf("Query returned %d rows, exceeding the maximum limit of %d rows", actualRows, maxRows),
	)
}

// NewDocumentTooLargeError creates an error for JSONB documents exceeding size limit
// maxSizeBytes is the maximum allowed document size in bytes (e.g., 1048576 for 1MB)
func NewDocumentTooLargeError(maxSizeBytes int) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeDocumentTooLarge,
		"Document too large",
		fmt.Sprintf("JSONB document exceeds maximum size of %d bytes", maxSizeBytes),
	)
}

// NewInternalError creates an error for internal server errors
func NewInternalError(detail string) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeInternalError,
		"An internal error occurred",
		detail,
	)
}

// NewServiceUnavailableError creates an error for service unavailability
func NewServiceUnavailableError(reason string) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeServiceUnavailable,
		"Service unavailable",
		reason,
	)
}

// NewDatabaseUnavailableError creates an error for database unavailability
func NewDatabaseUnavailableError(reason string) *postgres.VibeError {
	return postgres.NewVibeError(
		ErrorCodeDatabaseUnavailable,
		"Database unavailable",
		reason,
	)
}

// HTTPErrorCodeMapping maps VibeSQL error codes to HTTP status codes for reference.
// This map serves as documentation and is used in testing to verify consistency
// with the postgres package implementation.
var HTTPErrorCodeMapping = map[string]int{
	ErrorCodeInvalidSQL:           http.StatusBadRequest,           // 400
	ErrorCodeMissingRequiredField: http.StatusBadRequest,           // 400
	ErrorCodeUnsafeQuery:          http.StatusBadRequest,           // 400
	ErrorCodeQueryTimeout:         http.StatusRequestTimeout,       // 408
	ErrorCodeQueryTooLarge:        http.StatusRequestEntityTooLarge, // 413
	ErrorCodeResultTooLarge:       http.StatusRequestEntityTooLarge, // 413
	ErrorCodeDocumentTooLarge:     http.StatusRequestEntityTooLarge, // 413
	ErrorCodeInternalError:        http.StatusInternalServerError,  // 500
	ErrorCodeServiceUnavailable:   http.StatusServiceUnavailable,   // 503
	ErrorCodeDatabaseUnavailable:  http.StatusServiceUnavailable,   // 503
}

// ValidateHTTPStatusMapping validates that all error codes have correct HTTP status mappings.
// This function is used in testing to ensure consistency between the local reference mapping
// and the actual implementation in the postgres package.
func ValidateHTTPStatusMapping() error {
	for code, expectedStatus := range HTTPErrorCodeMapping {
		actualStatus := GetHTTPStatusCode(code)
		if actualStatus != expectedStatus {
			return fmt.Errorf("HTTP status mismatch for %s: expected %d, got %d", code, expectedStatus, actualStatus)
		}
	}
	return nil
}
