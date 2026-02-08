package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vibesql/vibe/internal/postgres"
)

func TestNewSuccessResponse(t *testing.T) {
	tests := []struct {
		name          string
		rows          []map[string]interface{}
		executionTime float64
		wantRowCount  int
		wantSuccess   bool
	}{
		{
			name: "with data rows",
			rows: []map[string]interface{}{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
			executionTime: 5.2,
			wantRowCount:  2,
			wantSuccess:   true,
		},
		{
			name:          "empty rows",
			rows:          []map[string]interface{}{},
			executionTime: 1.1,
			wantRowCount:  0,
			wantSuccess:   true,
		},
		{
			name:          "nil rows",
			rows:          nil,
			executionTime: 0.5,
			wantRowCount:  0,
			wantSuccess:   true,
		},
		{
			name: "single row",
			rows: []map[string]interface{}{
				{"count": 42},
			},
			executionTime: 2.3,
			wantRowCount:  1,
			wantSuccess:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewSuccessResponse(tt.rows, tt.executionTime)

			if resp.Success != tt.wantSuccess {
				t.Errorf("NewSuccessResponse().Success = %v, want %v", resp.Success, tt.wantSuccess)
			}

			if resp.RowCount != tt.wantRowCount {
				t.Errorf("NewSuccessResponse().RowCount = %v, want %v", resp.RowCount, tt.wantRowCount)
			}

			if resp.ExecutionTime != tt.executionTime {
				t.Errorf("NewSuccessResponse().ExecutionTime = %v, want %v", resp.ExecutionTime, tt.executionTime)
			}

			if resp.Error != nil {
				t.Errorf("NewSuccessResponse().Error = %v, want nil", resp.Error)
			}
		})
	}
}

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name        string
		err         *postgres.VibeError
		wantSuccess bool
		wantCode    string
		wantMessage string
		wantDetail  string
	}{
		{
			name: "invalid SQL error",
			err: postgres.NewVibeError(
				postgres.ErrorCodeInvalidSQL,
				"Invalid SQL syntax",
				"PostgreSQL error: syntax error at or near \"SELCT\"",
			),
			wantSuccess: false,
			wantCode:    postgres.ErrorCodeInvalidSQL,
			wantMessage: "Invalid SQL syntax",
			wantDetail:  "PostgreSQL error: syntax error at or near \"SELCT\"",
		},
		{
			name: "query timeout error",
			err: postgres.NewVibeError(
				postgres.ErrorCodeQueryTimeout,
				"Query execution timeout",
				"Query exceeded the maximum execution time of 5 seconds",
			),
			wantSuccess: false,
			wantCode:    postgres.ErrorCodeQueryTimeout,
			wantMessage: "Query execution timeout",
			wantDetail:  "Query exceeded the maximum execution time of 5 seconds",
		},
		{
			name: "error without detail",
			err: postgres.NewVibeError(
				postgres.ErrorCodeMissingRequiredField,
				"Missing required field",
				"",
			),
			wantSuccess: false,
			wantCode:    postgres.ErrorCodeMissingRequiredField,
			wantMessage: "Missing required field",
			wantDetail:  "",
		},
		{
			name:        "nil error",
			err:         nil,
			wantSuccess: false,
			wantCode:    postgres.ErrorCodeInternalError,
			wantMessage: "Unknown error occurred",
			wantDetail:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewErrorResponse(tt.err)

			if resp.Success != tt.wantSuccess {
				t.Errorf("NewErrorResponse().Success = %v, want %v", resp.Success, tt.wantSuccess)
			}

			if resp.Error == nil {
				t.Fatal("NewErrorResponse().Error is nil, want error detail")
			}

			if resp.Error.Code != tt.wantCode {
				t.Errorf("NewErrorResponse().Error.Code = %v, want %v", resp.Error.Code, tt.wantCode)
			}

			if resp.Error.Message != tt.wantMessage {
				t.Errorf("NewErrorResponse().Error.Message = %v, want %v", resp.Error.Message, tt.wantMessage)
			}

			if resp.Error.Detail != tt.wantDetail {
				t.Errorf("NewErrorResponse().Error.Detail = %v, want %v", resp.Error.Detail, tt.wantDetail)
			}

			if resp.Rows != nil {
				t.Errorf("NewErrorResponse().Rows = %v, want nil", resp.Rows)
			}

			if resp.RowCount != 0 {
				t.Errorf("NewErrorResponse().RowCount = %v, want 0", resp.RowCount)
			}

			if resp.ExecutionTime != 0 {
				t.Errorf("NewErrorResponse().ExecutionTime = %v, want 0", resp.ExecutionTime)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name           string
		response       *QueryResponse
		statusCode     int
		wantStatusCode int
		wantSuccess    bool
	}{
		{
			name: "success response",
			response: &QueryResponse{
				Success: true,
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
				},
				RowCount:      1,
				ExecutionTime: 2.5,
			},
			statusCode:     http.StatusOK,
			wantStatusCode: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name: "error response",
			response: &QueryResponse{
				Success: false,
				Error: &ErrorDetail{
					Code:    postgres.ErrorCodeInvalidSQL,
					Message: "Invalid SQL syntax",
					Detail:  "PostgreSQL error: syntax error",
				},
			},
			statusCode:     http.StatusBadRequest,
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := WriteJSON(w, tt.statusCode, tt.response)
			if err != nil {
				t.Fatalf("WriteJSON() error = %v", err)
			}

			if w.Code != tt.wantStatusCode {
				t.Errorf("WriteJSON() status code = %v, want %v", w.Code, tt.wantStatusCode)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("WriteJSON() Content-Type = %v, want application/json", contentType)
			}

			var decoded QueryResponse
			if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
				t.Fatalf("Failed to decode response JSON: %v", err)
			}

			if decoded.Success != tt.wantSuccess {
				t.Errorf("Decoded response.Success = %v, want %v", decoded.Success, tt.wantSuccess)
			}
		})
	}
}

func TestWriteSuccess(t *testing.T) {
	tests := []struct {
		name          string
		rows          []map[string]interface{}
		executionTime float64
		wantRowCount  int
	}{
		{
			name: "with data",
			rows: []map[string]interface{}{
				{"id": 1, "data": map[string]interface{}{"name": "Alice"}},
				{"id": 2, "data": map[string]interface{}{"name": "Bob"}},
			},
			executionTime: 3.7,
			wantRowCount:  2,
		},
		{
			name:          "empty result",
			rows:          []map[string]interface{}{},
			executionTime: 0.8,
			wantRowCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := WriteSuccess(w, tt.rows, tt.executionTime)
			if err != nil {
				t.Fatalf("WriteSuccess() error = %v", err)
			}

			if w.Code != http.StatusOK {
				t.Errorf("WriteSuccess() status code = %v, want %v", w.Code, http.StatusOK)
			}

			var decoded QueryResponse
			if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
				t.Fatalf("Failed to decode response JSON: %v", err)
			}

			if !decoded.Success {
				t.Error("WriteSuccess() response.Success = false, want true")
			}

			if decoded.RowCount != tt.wantRowCount {
				t.Errorf("WriteSuccess() response.RowCount = %v, want %v", decoded.RowCount, tt.wantRowCount)
			}

			if decoded.ExecutionTime != tt.executionTime {
				t.Errorf("WriteSuccess() response.ExecutionTime = %v, want %v", decoded.ExecutionTime, tt.executionTime)
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name               string
		err                *postgres.VibeError
		wantHTTPStatusCode int
		wantErrorCode      string
	}{
		{
			name: "invalid SQL error",
			err: postgres.NewVibeError(
				postgres.ErrorCodeInvalidSQL,
				"Invalid SQL syntax",
				"PostgreSQL error: syntax error",
			),
			wantHTTPStatusCode: http.StatusBadRequest,
			wantErrorCode:      postgres.ErrorCodeInvalidSQL,
		},
		{
			name: "query timeout error",
			err: postgres.NewVibeError(
				postgres.ErrorCodeQueryTimeout,
				"Query execution timeout",
				"Query exceeded 5 seconds",
			),
			wantHTTPStatusCode: http.StatusRequestTimeout,
			wantErrorCode:      postgres.ErrorCodeQueryTimeout,
		},
		{
			name: "query too large error",
			err: postgres.NewVibeError(
				postgres.ErrorCodeQueryTooLarge,
				"Query too large",
				"Query exceeds 10KB limit",
			),
			wantHTTPStatusCode: http.StatusRequestEntityTooLarge,
			wantErrorCode:      postgres.ErrorCodeQueryTooLarge,
		},
		{
			name: "internal error",
			err: postgres.NewVibeError(
				postgres.ErrorCodeInternalError,
				"Internal error",
				"Unexpected error occurred",
			),
			wantHTTPStatusCode: http.StatusInternalServerError,
			wantErrorCode:      postgres.ErrorCodeInternalError,
		},
		{
			name: "database unavailable error",
			err: postgres.NewVibeError(
				postgres.ErrorCodeDatabaseUnavailable,
				"Database unavailable",
				"Cannot connect to database",
			),
			wantHTTPStatusCode: http.StatusServiceUnavailable,
			wantErrorCode:      postgres.ErrorCodeDatabaseUnavailable,
		},
		{
			name:               "nil error",
			err:                nil,
			wantHTTPStatusCode: http.StatusInternalServerError,
			wantErrorCode:      postgres.ErrorCodeInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := WriteError(w, tt.err)
			if err != nil {
				t.Fatalf("WriteError() error = %v", err)
			}

			if w.Code != tt.wantHTTPStatusCode {
				t.Errorf("WriteError() status code = %v, want %v", w.Code, tt.wantHTTPStatusCode)
			}

			var decoded QueryResponse
			if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
				t.Fatalf("Failed to decode response JSON: %v", err)
			}

			if decoded.Success {
				t.Error("WriteError() response.Success = true, want false")
			}

			if decoded.Error == nil {
				t.Fatal("WriteError() response.Error is nil, want error detail")
			}

			if decoded.Error.Code != tt.wantErrorCode {
				t.Errorf("WriteError() response.Error.Code = %v, want %v", decoded.Error.Code, tt.wantErrorCode)
			}
		})
	}
}

func TestJSONSerialization(t *testing.T) {
	t.Run("success response serializes correctly", func(t *testing.T) {
		resp := NewSuccessResponse([]map[string]interface{}{
			{"id": 1, "name": "Alice"},
		}, 5.2)

		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		// Verify JSON structure
		var decoded map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		if success, ok := decoded["success"].(bool); !ok || !success {
			t.Error("JSON missing or incorrect 'success' field")
		}

		if rowCount, ok := decoded["rowCount"].(float64); !ok || rowCount != 1 {
			t.Errorf("JSON rowCount = %v, want 1", rowCount)
		}

		if executionTime, ok := decoded["executionTime"].(float64); !ok || executionTime != 5.2 {
			t.Errorf("JSON executionTime = %v, want 5.2", executionTime)
		}
	})

	t.Run("error response serializes correctly", func(t *testing.T) {
		vibeErr := postgres.NewVibeError(
			postgres.ErrorCodeInvalidSQL,
			"Invalid SQL syntax",
			"PostgreSQL error: syntax error",
		)
		resp := NewErrorResponse(vibeErr)

		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		// Verify JSON structure
		var decoded map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		if success, ok := decoded["success"].(bool); !ok || success {
			t.Error("JSON success should be false")
		}

		errorDetail, ok := decoded["error"].(map[string]interface{})
		if !ok {
			t.Fatal("JSON missing 'error' field")
		}

		if code, ok := errorDetail["code"].(string); !ok || code != postgres.ErrorCodeInvalidSQL {
			t.Errorf("JSON error.code = %v, want %v", code, postgres.ErrorCodeInvalidSQL)
		}

		if message, ok := errorDetail["message"].(string); !ok || message != "Invalid SQL syntax" {
			t.Errorf("JSON error.message = %v, want 'Invalid SQL syntax'", message)
		}
	})

	t.Run("success response omits error field", func(t *testing.T) {
		resp := NewSuccessResponse(nil, 1.0)

		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		// Verify 'error' field is not present in JSON
		if bytes.Contains(jsonBytes, []byte("\"error\"")) {
			t.Error("Success response should not contain 'error' field in JSON")
		}
	})

	t.Run("error response omits rows fields", func(t *testing.T) {
		vibeErr := postgres.NewVibeError(
			postgres.ErrorCodeInvalidSQL,
			"Invalid SQL",
			"",
		)
		resp := NewErrorResponse(vibeErr)

		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		// Verify 'rows', 'rowCount', 'executionTime' fields are not present
		if bytes.Contains(jsonBytes, []byte("\"rows\"")) {
			t.Error("Error response should not contain 'rows' field in JSON")
		}

		if bytes.Contains(jsonBytes, []byte("\"rowCount\"")) {
			t.Error("Error response should not contain 'rowCount' field in JSON")
		}

		if bytes.Contains(jsonBytes, []byte("\"executionTime\"")) {
			t.Error("Error response should not contain 'executionTime' field in JSON")
		}
	})
}

func TestQueryResponseEdgeCases(t *testing.T) {
	t.Run("response with JSONB data", func(t *testing.T) {
		rows := []map[string]interface{}{
			{
				"id":   1,
				"data": map[string]interface{}{"name": "Alice", "age": 30},
			},
		}

		resp := NewSuccessResponse(rows, 2.5)

		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal JSONB response: %v", err)
		}

		var decoded QueryResponse
		if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal JSONB response: %v", err)
		}

		if decoded.RowCount != 1 {
			t.Errorf("JSONB response rowCount = %v, want 1", decoded.RowCount)
		}
	})

	t.Run("large row count", func(t *testing.T) {
		rows := make([]map[string]interface{}, 1000)
		for i := 0; i < 1000; i++ {
			rows[i] = map[string]interface{}{"id": i}
		}

		resp := NewSuccessResponse(rows, 12.3)

		if resp.RowCount != 1000 {
			t.Errorf("Large response rowCount = %v, want 1000", resp.RowCount)
		}
	})

	t.Run("error with empty detail", func(t *testing.T) {
		vibeErr := postgres.NewVibeError(
			postgres.ErrorCodeUnsafeQuery,
			"Unsafe query",
			"",
		)

		resp := NewErrorResponse(vibeErr)

		jsonBytes, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal error response: %v", err)
		}

		// Empty detail should be omitted from JSON (omitempty tag)
		var decoded map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		errorDetail := decoded["error"].(map[string]interface{})
		if _, hasDetail := errorDetail["detail"]; hasDetail {
			t.Error("Empty detail should be omitted from JSON")
		}
	})
}

func BenchmarkNewSuccessResponse(b *testing.B) {
	rows := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewSuccessResponse(rows, 5.2)
	}
}

func BenchmarkNewErrorResponse(b *testing.B) {
	err := postgres.NewVibeError(
		postgres.ErrorCodeInvalidSQL,
		"Invalid SQL syntax",
		"Test error detail",
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewErrorResponse(err)
	}
}
