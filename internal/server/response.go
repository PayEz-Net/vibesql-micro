package server

import (
	"encoding/json"
	"net/http"

	"github.com/vibesql/vibe/internal/postgres"
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

// NewSuccessResponse creates a successful query response
func NewSuccessResponse(rows []map[string]interface{}, executionTime float64) *QueryResponse {
	rowCount := 0
	if rows != nil {
		rowCount = len(rows)
	}

	return &QueryResponse{
		Success:       true,
		Rows:          rows,
		RowCount:      rowCount,
		ExecutionTime: executionTime,
	}
}

// NewErrorResponse creates an error response from a VibeError
func NewErrorResponse(err *postgres.VibeError) *QueryResponse {
	if err == nil {
		return &QueryResponse{
			Success: false,
			Error: &ErrorDetail{
				Code:    postgres.ErrorCodeInternalError,
				Message: "Unknown error occurred",
			},
		}
	}

	return &QueryResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:    err.Code,
			Message: err.Message,
			Detail:  err.Detail,
		},
	}
}

// WriteJSON writes a QueryResponse as JSON to the HTTP response writer
func WriteJSON(w http.ResponseWriter, statusCode int, response *QueryResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	return encoder.Encode(response)
}

// WriteSuccess writes a successful query response with 200 OK status
func WriteSuccess(w http.ResponseWriter, rows []map[string]interface{}, executionTime float64) error {
	response := NewSuccessResponse(rows, executionTime)
	return WriteJSON(w, http.StatusOK, response)
}

// WriteError writes an error response with appropriate HTTP status code
func WriteError(w http.ResponseWriter, err *postgres.VibeError) error {
	response := NewErrorResponse(err)
	// Use response.Error.Code instead of err.Code to safely handle nil errors
	statusCode := postgres.GetHTTPStatusCode(response.Error.Code)
	return WriteJSON(w, statusCode, response)
}
