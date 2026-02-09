package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// VibeSQL error codes
const (
	ErrorCodeInvalidSQL          = "INVALID_SQL"
	ErrorCodeMissingRequiredField = "MISSING_REQUIRED_FIELD"
	ErrorCodeUnsafeQuery         = "UNSAFE_QUERY"
	ErrorCodeQueryTimeout        = "QUERY_TIMEOUT"
	ErrorCodeQueryTooLarge       = "QUERY_TOO_LARGE"
	ErrorCodeResultTooLarge      = "RESULT_TOO_LARGE"
	ErrorCodeDocumentTooLarge    = "DOCUMENT_TOO_LARGE"
	ErrorCodeInternalError       = "INTERNAL_ERROR"
	ErrorCodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
	ErrorCodeDatabaseUnavailable = "DATABASE_UNAVAILABLE"
)

// HTTP status codes for VibeSQL errors
const (
	HTTPStatusInvalidSQL          = 400
	HTTPStatusMissingRequiredField = 400
	HTTPStatusUnsafeQuery         = 400
	HTTPStatusQueryTimeout        = 408
	HTTPStatusQueryTooLarge       = 413
	HTTPStatusResultTooLarge      = 413
	HTTPStatusDocumentTooLarge    = 413
	HTTPStatusInternalError       = 500
	HTTPStatusServiceUnavailable  = 503
	HTTPStatusDatabaseUnavailable = 503
)

// VibeError represents a VibeSQL error
type VibeError struct {
	Code    string
	Message string
	Detail  string
}

func (e *VibeError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewVibeError creates a new VibeSQL error
func NewVibeError(code, message, detail string) *VibeError {
	return &VibeError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// SQLSTATE to VibeSQL error code mapping
var sqlStateToVibeCode = map[string]string{
	// Syntax errors → INVALID_SQL
	"42601": ErrorCodeInvalidSQL, // syntax_error
	"42703": ErrorCodeInvalidSQL, // undefined_column
	"42P01": ErrorCodeInvalidSQL, // undefined_table
	"42P02": ErrorCodeInvalidSQL, // undefined_parameter
	"42883": ErrorCodeInvalidSQL, // undefined_function
	"42804": ErrorCodeInvalidSQL, // datatype_mismatch
	
	// Query cancellation → QUERY_TIMEOUT
	"57014": ErrorCodeQueryTimeout, // query_canceled
	
	// Resource limits → DATABASE_UNAVAILABLE
	"53000": ErrorCodeDatabaseUnavailable, // insufficient_resources
	"53100": ErrorCodeDatabaseUnavailable, // disk_full
	"53200": ErrorCodeDatabaseUnavailable, // out_of_memory
	"53300": ErrorCodeDatabaseUnavailable, // too_many_connections
	"53400": ErrorCodeDatabaseUnavailable, // configuration_limit_exceeded
	
	// Connection errors → DATABASE_UNAVAILABLE
	"08000": ErrorCodeDatabaseUnavailable, // connection_exception
	"08003": ErrorCodeDatabaseUnavailable, // connection_does_not_exist
	"08006": ErrorCodeDatabaseUnavailable, // connection_failure
	"08001": ErrorCodeDatabaseUnavailable, // sqlclient_unable_to_establish_sqlconnection
	"08004": ErrorCodeDatabaseUnavailable, // sqlserver_rejected_establishment_of_sqlconnection
	
	// Document size errors → DOCUMENT_TOO_LARGE
	"54000": ErrorCodeDocumentTooLarge, // program_limit_exceeded
	"54001": ErrorCodeDocumentTooLarge, // statement_too_complex
}

// TranslateError translates a PostgreSQL error to a VibeSQL error
func TranslateError(err error) *VibeError {
	if err == nil {
		return nil
	}
	
	// Check if it's already a VibeError
	var vibeErr *VibeError
	if errors.As(err, &vibeErr) {
		return vibeErr
	}
	
	// Check for context timeout/cancellation (critical for query timeouts)
	if errors.Is(err, context.DeadlineExceeded) {
		return NewVibeError(
			ErrorCodeQueryTimeout,
			"Query execution timeout",
			"Query exceeded the maximum execution time of 5 seconds",
		)
	}
	
	if errors.Is(err, context.Canceled) {
		return NewVibeError(
			ErrorCodeQueryTimeout,
			"Query execution canceled",
			"Query was canceled before completion",
		)
	}
	
	// Check if it's a PostgreSQL error
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return translatePQError(pqErr)
	}
	
	// Unknown error type → INTERNAL_ERROR
	return NewVibeError(
		ErrorCodeInternalError,
		"An internal error occurred",
		err.Error(),
	)
}

// translatePQError translates a pq.Error to a VibeError
func translatePQError(pqErr *pq.Error) *VibeError {
	sqlState := string(pqErr.Code)
	
	// Map SQLSTATE to VibeSQL error code
	vibeCode, found := sqlStateToVibeCode[sqlState]
	if !found {
		// Unknown SQLSTATE → INTERNAL_ERROR
		vibeCode = ErrorCodeInternalError
	}
	
	// Build error message
	message := buildErrorMessage(vibeCode, pqErr)
	detail := buildErrorDetail(pqErr)
	
	return NewVibeError(vibeCode, message, detail)
}

// buildErrorMessage creates a user-friendly error message
func buildErrorMessage(vibeCode string, pqErr *pq.Error) string {
	switch vibeCode {
	case ErrorCodeInvalidSQL:
		return "Invalid SQL syntax"
	case ErrorCodeQueryTimeout:
		return "Query execution timeout"
	case ErrorCodeDatabaseUnavailable:
		return "Database is unavailable"
	case ErrorCodeDocumentTooLarge:
		return "Document too large"
	default:
		// Use PostgreSQL's message if available
		if pqErr.Message != "" {
			return pqErr.Message
		}
		return "An error occurred"
	}
}

// buildErrorDetail creates detailed error information
func buildErrorDetail(pqErr *pq.Error) string {
	detail := fmt.Sprintf("PostgreSQL error: %s", pqErr.Message)
	
	if pqErr.Detail != "" {
		detail += fmt.Sprintf(" | Detail: %s", pqErr.Detail)
	}
	
	if pqErr.Hint != "" {
		detail += fmt.Sprintf(" | Hint: %s", pqErr.Hint)
	}
	
	if pqErr.Position != "" {
		detail += fmt.Sprintf(" | Position: %s", pqErr.Position)
	}
	
	return detail
}

// GetHTTPStatusCode returns the HTTP status code for a VibeSQL error code
func GetHTTPStatusCode(errorCode string) int {
	switch errorCode {
	case ErrorCodeInvalidSQL:
		return HTTPStatusInvalidSQL
	case ErrorCodeMissingRequiredField:
		return HTTPStatusMissingRequiredField
	case ErrorCodeUnsafeQuery:
		return HTTPStatusUnsafeQuery
	case ErrorCodeQueryTimeout:
		return HTTPStatusQueryTimeout
	case ErrorCodeQueryTooLarge:
		return HTTPStatusQueryTooLarge
	case ErrorCodeResultTooLarge:
		return HTTPStatusResultTooLarge
	case ErrorCodeDocumentTooLarge:
		return HTTPStatusDocumentTooLarge
	case ErrorCodeInternalError:
		return HTTPStatusInternalError
	case ErrorCodeServiceUnavailable:
		return HTTPStatusServiceUnavailable
	case ErrorCodeDatabaseUnavailable:
		return HTTPStatusDatabaseUnavailable
	default:
		return HTTPStatusInternalError
	}
}
