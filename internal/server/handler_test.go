package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vibesql/vibe/internal/postgres"
	"github.com/vibesql/vibe/internal/query"
)

func setupTestDB(t *testing.T) *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: cannot ping test database: %v", err)
	}

	_, _ = db.Exec("DROP TABLE IF EXISTS test_handler_users")
	_, err = db.Exec(`
		CREATE TABLE test_handler_users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			data JSONB
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return db
}

func teardownTestDB(db *sql.DB) {
	if db != nil {
		_, _ = db.Exec("DROP TABLE IF EXISTS test_handler_users")
		db.Close()
	}
}

func TestHandleQuery_Success(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	reqBody := QueryRequest{SQL: "SELECT 1 as test"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false: %+v", response.Error)
	}

	if response.RowCount != 1 {
		t.Errorf("Expected rowCount=1, got %d", response.RowCount)
	}

	if len(response.Rows) != 1 {
		t.Fatalf("Expected 1 row, got %d", len(response.Rows))
	}

	if response.Rows[0]["test"] == nil {
		t.Errorf("Expected 'test' column in result")
	}

	if response.ExecutionTime <= 0 {
		t.Errorf("Expected executionTime > 0, got %f", response.ExecutionTime)
	}
}

func TestHandleQuery_MethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	req := httptest.NewRequest(http.MethodGet, "/v1/query", nil)
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeInvalidSQL {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeInvalidSQL, response.Error.Code)
	}

	if !strings.Contains(response.Error.Detail, "POST") {
		t.Errorf("Expected error detail to mention POST method, got: %s", response.Error.Detail)
	}
}

func TestHandleQuery_InvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeInvalidSQL {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeInvalidSQL, response.Error.Code)
	}
}

func TestHandleQuery_MissingSQLField(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	reqBody := map[string]string{}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeMissingRequiredField {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeMissingRequiredField, response.Error.Code)
	}

	if !strings.Contains(response.Error.Message, "sql") {
		t.Errorf("Expected error message to mention 'sql' field, got: %s", response.Error.Message)
	}
}

func TestHandleQuery_EmptySQLField(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	reqBody := QueryRequest{SQL: ""}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeMissingRequiredField {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeMissingRequiredField, response.Error.Code)
	}
}

func TestHandleQuery_QueryTooLarge(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	largeSQL := "SELECT " + strings.Repeat("'x',", 10*1024)
	reqBody := QueryRequest{SQL: largeSQL}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status 413, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeQueryTooLarge {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeQueryTooLarge, response.Error.Code)
	}
}

func TestHandleQuery_InvalidSQLSyntax(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	reqBody := QueryRequest{SQL: "not a valid sql keyword"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeInvalidSQL {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeInvalidSQL, response.Error.Code)
	}
}

func TestHandleQuery_UnsafeUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	reqBody := QueryRequest{SQL: "UPDATE test_handler_users SET name = 'test'"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeUnsafeQuery {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeUnsafeQuery, response.Error.Code)
	}
}

func TestHandleQuery_UnsafeDelete(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	reqBody := QueryRequest{SQL: "DELETE FROM test_handler_users"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeUnsafeQuery {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeUnsafeQuery, response.Error.Code)
	}
}

func TestHandleQuery_SafeUpdateWithWhere(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	_, _ = db.Exec("INSERT INTO test_handler_users (name, email) VALUES ('John', 'john@example.com')")

	reqBody := QueryRequest{SQL: "UPDATE test_handler_users SET name = 'Jane' WHERE name = 'John'"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false: %+v", response.Error)
	}
}

func TestHandleQuery_InvalidSQL(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	reqBody := QueryRequest{SQL: "SELECT * FROM nonexistent_table"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}

	if response.Error.Code != postgres.ErrorCodeInvalidSQL {
		t.Errorf("Expected error code %s, got %s", postgres.ErrorCodeInvalidSQL, response.Error.Code)
	}
}

func TestHandleQuery_EmptyResultSet(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	reqBody := QueryRequest{SQL: "SELECT * FROM test_handler_users WHERE id = 999999"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleQuery(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false: %+v", response.Error)
	}

	if response.RowCount != 0 {
		t.Errorf("Expected rowCount=0, got %d", response.RowCount)
	}

	if len(response.Rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(response.Rows))
	}

	if response.ExecutionTime <= 0 {
		t.Errorf("Expected executionTime > 0, got %f", response.ExecutionTime)
	}
}

func TestHandleQuery_FullCRUDWorkflow(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	tests := []struct {
		name          string
		sql           string
		expectSuccess bool
		expectRows    int
	}{
		{
			name:          "Insert row",
			sql:           "INSERT INTO test_handler_users (name, email) VALUES ('Alice', 'alice@example.com') RETURNING id",
			expectSuccess: true,
			expectRows:    1,
		},
		{
			name:          "Select row",
			sql:           "SELECT * FROM test_handler_users WHERE name = 'Alice'",
			expectSuccess: true,
			expectRows:    1,
		},
		{
			name:          "Update row",
			sql:           "UPDATE test_handler_users SET email = 'alice.new@example.com' WHERE name = 'Alice'",
			expectSuccess: true,
			expectRows:    0,
		},
		{
			name:          "Delete row",
			sql:           "DELETE FROM test_handler_users WHERE name = 'Alice'",
			expectSuccess: true,
			expectRows:    0,
		},
		{
			name:          "Verify deletion",
			sql:           "SELECT * FROM test_handler_users WHERE name = 'Alice'",
			expectSuccess: true,
			expectRows:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := QueryRequest{SQL: tt.sql}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleQuery(w, req)

			if tt.expectSuccess && w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			var response QueryResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("Expected success=%v, got %v: %+v", tt.expectSuccess, response.Success, response.Error)
			}

			if response.Success && response.RowCount != tt.expectRows {
				t.Errorf("Expected %d rows, got %d", tt.expectRows, response.RowCount)
			}
		})
	}
}

func TestHandleQuery_JSONBOperations(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	_, _ = db.Exec(`
		INSERT INTO test_handler_users (name, email, data) 
		VALUES ('Bob', 'bob@example.com', '{"age": 30, "city": "NYC", "tags": ["go", "sql"]}')
	`)

	tests := []struct {
		name          string
		sql           string
		expectSuccess bool
	}{
		{
			name:          "JSONB field access",
			sql:           "SELECT data->'age' as age FROM test_handler_users WHERE name = 'Bob'",
			expectSuccess: true,
		},
		{
			name:          "JSONB text extraction",
			sql:           "SELECT data->>'city' as city FROM test_handler_users WHERE name = 'Bob'",
			expectSuccess: true,
		},
		{
			name:          "JSONB containment",
			sql:           "SELECT * FROM test_handler_users WHERE data @> '{\"city\": \"NYC\"}'",
			expectSuccess: true,
		},
		{
			name:          "JSONB key existence",
			sql:           "SELECT * FROM test_handler_users WHERE data ? 'age'",
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := QueryRequest{SQL: tt.sql}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleQuery(w, req)

			if tt.expectSuccess && w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			var response QueryResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("Expected success=%v, got %v: %+v", tt.expectSuccess, response.Success, response.Error)
			}
		})
	}
}

func TestHandleQuery_RegisterRoutes(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	reqBody := QueryRequest{SQL: "SELECT 1 as test"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response QueryResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false")
	}
}

func TestHandleQuery_Concurrency(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	executor := query.NewExecutor(db)
	handler := NewHandler(executor)

	const numRequests = 10

	type result struct {
		statusCode int
		success    bool
		err        error
	}

	results := make(chan result, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			reqBody := QueryRequest{SQL: fmt.Sprintf("SELECT %d as id", id)}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest(http.MethodPost, "/v1/query", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleQuery(w, req)

			var response QueryResponse
			err := json.NewDecoder(w.Body).Decode(&response)

			results <- result{
				statusCode: w.Code,
				success:    response.Success,
				err:        err,
			}
		}(i)
	}

	successCount := 0
	for i := 0; i < numRequests; i++ {
		res := <-results
		if res.err != nil {
			t.Errorf("Request %d failed to decode: %v", i, res.err)
			continue
		}
		if res.statusCode == http.StatusOK && res.success {
			successCount++
		}
	}

	if successCount != numRequests {
		t.Errorf("Expected %d successful requests, got %d", numRequests, successCount)
	}
}
