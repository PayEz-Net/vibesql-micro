package query

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/vibesql/vibe/internal/postgres"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
	}

	return db
}

func TestExecutor_Execute_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	result, err := executor.Execute("SELECT 1 as test")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.RowCount != 1 {
		t.Errorf("Expected RowCount = 1, got %d", result.RowCount)
	}

	if len(result.Rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(result.Rows))
	}

	if result.Rows[0]["test"] != int64(1) {
		t.Errorf("Expected test = 1, got %v", result.Rows[0]["test"])
	}

	if result.ExecutionTime <= 0 {
		t.Error("Expected positive execution time")
	}
}

func TestExecutor_Execute_MultipleRows(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	result, err := executor.Execute("SELECT generate_series(1, 10) as num")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.RowCount != 10 {
		t.Errorf("Expected RowCount = 10, got %d", result.RowCount)
	}

	if len(result.Rows) != 10 {
		t.Errorf("Expected 10 rows, got %d", len(result.Rows))
	}
}

func TestExecutor_Execute_Timeout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	startTime := time.Now()
	_, err := executor.Execute("SELECT pg_sleep(10)")
	elapsed := time.Since(startTime)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	vibeErr, ok := err.(*postgres.VibeError)
	if !ok {
		t.Fatalf("Expected VibeError, got %T", err)
	}

	if vibeErr.Code != postgres.ErrorCodeQueryTimeout {
		t.Errorf("Expected QUERY_TIMEOUT error, got %s", vibeErr.Code)
	}

	if elapsed < 4*time.Second || elapsed > 6*time.Second {
		t.Errorf("Expected timeout around 5s, got %v", elapsed)
	}
}

func TestExecutor_Execute_QueryCompletesJustBeforeTimeout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	result, err := executor.Execute("SELECT pg_sleep(3)")
	if err != nil {
		t.Fatalf("Expected no error for query completing in 3s, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestExecutor_Execute_ResultTooLarge(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	_, err := executor.Execute("SELECT generate_series(1, 1001) as num")
	if err == nil {
		t.Fatal("Expected RESULT_TOO_LARGE error, got nil")
	}

	vibeErr, ok := err.(*postgres.VibeError)
	if !ok {
		t.Fatalf("Expected VibeError, got %T", err)
	}

	if vibeErr.Code != postgres.ErrorCodeResultTooLarge {
		t.Errorf("Expected RESULT_TOO_LARGE error, got %s", vibeErr.Code)
	}
}

func TestExecutor_Execute_ExactlyMaxRows(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	result, err := executor.Execute("SELECT generate_series(1, 1000) as num")
	if err != nil {
		t.Fatalf("Expected no error for exactly 1000 rows, got: %v", err)
	}

	if result.RowCount != 1000 {
		t.Errorf("Expected RowCount = 1000, got %d", result.RowCount)
	}
}

func TestExecutor_Execute_InvalidSQL(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	_, err := executor.Execute("SELECT * FROM nonexistent_table")
	if err == nil {
		t.Fatal("Expected error for invalid SQL, got nil")
	}

	vibeErr, ok := err.(*postgres.VibeError)
	if !ok {
		t.Fatalf("Expected VibeError, got %T", err)
	}

	if vibeErr.Code != postgres.ErrorCodeInvalidSQL {
		t.Errorf("Expected INVALID_SQL error, got %s", vibeErr.Code)
	}
}

func TestExecutor_Execute_EmptyResult(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	result, err := executor.Execute("SELECT 1 WHERE false")
	if err != nil {
		t.Fatalf("Expected no error for empty result, got: %v", err)
	}

	if result.RowCount != 0 {
		t.Errorf("Expected RowCount = 0, got %d", result.RowCount)
	}

	if len(result.Rows) != 0 {
		t.Errorf("Expected 0 rows, got %d", len(result.Rows))
	}
}

func TestExecutor_Execute_VariousDataTypes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	sql := `
		SELECT 
			42 as int_col,
			'hello' as text_col,
			true as bool_col,
			3.14::float as float_col,
			'{"key": "value"}'::jsonb as jsonb_col
	`

	result, err := executor.Execute(sql)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.RowCount != 1 {
		t.Fatalf("Expected 1 row, got %d", result.RowCount)
	}

	row := result.Rows[0]

	if row["int_col"] != int64(42) {
		t.Errorf("Expected int_col = 42, got %v (%T)", row["int_col"], row["int_col"])
	}

	if row["text_col"] != "hello" {
		t.Errorf("Expected text_col = 'hello', got %v", row["text_col"])
	}

	if row["bool_col"] != true {
		t.Errorf("Expected bool_col = true, got %v", row["bool_col"])
	}

	if row["jsonb_col"] == nil {
		t.Error("Expected jsonb_col to have value, got nil")
	}
}

func TestExecutor_Execute_TimeoutPrecision(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)

	testCases := []struct {
		name          string
		sleepSeconds  int
		expectTimeout bool
	}{
		{"4 seconds - should succeed", 4, false},
		{"6 seconds - should timeout", 6, true},
		{"7 seconds - should timeout", 7, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sql := fmt.Sprintf("SELECT pg_sleep(%d)", tc.sleepSeconds)
			startTime := time.Now()
			_, err := executor.Execute(sql)
			elapsed := time.Since(startTime)

			if tc.expectTimeout {
				if err == nil {
					t.Errorf("Expected timeout error, got nil (elapsed: %v)", elapsed)
					return
				}
				vibeErr, ok := err.(*postgres.VibeError)
				if !ok {
					t.Errorf("Expected VibeError, got %T", err)
					return
				}
				if vibeErr.Code != postgres.ErrorCodeQueryTimeout {
					t.Errorf("Expected QUERY_TIMEOUT, got %s", vibeErr.Code)
				}
				if elapsed < 4*time.Second || elapsed > 6*time.Second {
					t.Errorf("Expected timeout around 5s, got %v", elapsed)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v (elapsed: %v)", err, elapsed)
				}
			}
		})
	}
}

func TestCheckRowLimit(t *testing.T) {
	testCases := []struct {
		name          string
		currentCount  int
		expectError   bool
		expectedCode  string
	}{
		{"0 rows - should pass", 0, false, ""},
		{"500 rows - should pass", 500, false, ""},
		{"999 rows - should pass", 999, false, ""},
		{"1000 rows - should fail", 1000, true, postgres.ErrorCodeResultTooLarge},
		{"1001 rows - should fail", 1001, true, postgres.ErrorCodeResultTooLarge},
		{"2000 rows - should fail", 2000, true, postgres.ErrorCodeResultTooLarge},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckRowLimit(tc.currentCount)

			if tc.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				vibeErr, ok := err.(*postgres.VibeError)
				if !ok {
					t.Errorf("Expected VibeError, got %T", err)
					return
				}
				if vibeErr.Code != tc.expectedCode {
					t.Errorf("Expected error code %s, got %s", tc.expectedCode, vibeErr.Code)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

func BenchmarkCheckRowLimit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CheckRowLimit(500)
	}
}
