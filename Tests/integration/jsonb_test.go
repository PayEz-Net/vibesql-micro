package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// Test configuration
const (
	defaultTestDB = "host=localhost port=5432 user=postgres password=postgres dbname=vibesql_test sslmode=disable"
)

// getTestDB returns a database connection for testing
// Reads from VIBESQL_TEST_DB environment variable, or uses default localhost connection
func getTestDB(t *testing.T) *sql.DB {
	connStr := os.Getenv("VIBESQL_TEST_DB")
	if connStr == "" {
		connStr = defaultTestDB
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v\nConnection string: %s\nIs PostgreSQL running?", err, connStr)
	}

	return db
}

// setupTestTable creates a test table for JSONB testing
func setupTestTable(t *testing.T, db *sql.DB) {
	// Drop table if exists
	_, err := db.Exec("DROP TABLE IF EXISTS jsonb_test")
	if err != nil {
		t.Fatalf("Failed to drop test table: %v", err)
	}

	// Create test table with JSONB column
	_, err = db.Exec(`
		CREATE TABLE jsonb_test (
			id SERIAL PRIMARY KEY,
			data JSONB NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
}

// insertTestData inserts sample JSONB data for testing
func insertTestData(t *testing.T, db *sql.DB) {
	testData := []string{
		`{"name": "Alice", "age": 30, "city": "NYC"}`,
		`{"name": "Bob", "age": 25, "city": "LA", "hobbies": ["coding", "gaming"]}`,
		`{"name": "Charlie", "age": 35, "address": {"street": "123 Main St", "zip": "10001"}}`,
		`{"name": "Diana", "tags": ["developer", "designer"], "active": true}`,
		`{"name": "Eve", "metadata": {"level": 5, "score": 9500}}`,
		`{"products": [{"id": 1, "name": "Widget"}, {"id": 2, "name": "Gadget"}]}`,
		`{"config": {"enabled": true, "options": {"timeout": 30, "retry": 3}}}`,
	}

	for _, data := range testData {
		_, err := db.Exec("INSERT INTO jsonb_test (data) VALUES ($1)", data)
		if err != nil {
			t.Fatalf("Failed to insert test data: %v", err)
		}
	}
}

// TestJSONB_BasicFieldAccess tests basic field access with -> operator
func TestJSONB_BasicFieldAccess(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 1: -> operator returns JSONB type
	var result sql.NullString
	err := db.QueryRow("SELECT data->'name' FROM jsonb_test WHERE data->>'name' = 'Alice'").Scan(&result)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if !result.Valid || result.String != `"Alice"` {
		t.Errorf("Expected '\"Alice\"', got '%s'", result.String)
	}
}

// TestJSONB_TextExtraction tests text extraction with ->> operator
func TestJSONB_TextExtraction(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 2: ->> operator returns text type
	var name string
	err := db.QueryRow("SELECT data->>'name' FROM jsonb_test WHERE data->>'name' = 'Bob'").Scan(&name)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if name != "Bob" {
		t.Errorf("Expected 'Bob', got '%s'", name)
	}

	// Test 3: Extract numeric field as text
	var age string
	err = db.QueryRow("SELECT data->>'age' FROM jsonb_test WHERE data->>'name' = 'Alice'").Scan(&age)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if age != "30" {
		t.Errorf("Expected '30', got '%s'", age)
	}
}

// TestJSONB_NestedPathAccess tests nested path access with #> and #>> operators
func TestJSONB_NestedPathAccess(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 4: #> operator for nested JSONB extraction
	var result sql.NullString
	err := db.QueryRow("SELECT data#>'{address,street}' FROM jsonb_test WHERE data->>'name' = 'Charlie'").Scan(&result)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if !result.Valid || result.String != `"123 Main St"` {
		t.Errorf("Expected '\"123 Main St\"', got '%s'", result.String)
	}

	// Test 5: #>> operator for nested text extraction
	var street string
	err = db.QueryRow("SELECT data#>>'{address,street}' FROM jsonb_test WHERE data->>'name' = 'Charlie'").Scan(&street)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if street != "123 Main St" {
		t.Errorf("Expected '123 Main St', got '%s'", street)
	}

	// Test 6: Deep nested path access
	var timeout string
	err = db.QueryRow("SELECT data#>>'{config,options,timeout}' FROM jsonb_test WHERE data->'config' IS NOT NULL").Scan(&timeout)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if timeout != "30" {
		t.Errorf("Expected '30', got '%s'", timeout)
	}
}

// TestJSONB_ContainmentOperators tests @> and <@ containment operators
func TestJSONB_ContainmentOperators(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 7: @> operator - top-level object contains smaller object
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data @> '{\"name\": \"Alice\"}'").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row with name Alice, got %d", count)
	}

	// Test 8: @> operator with nested containment
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data @> '{\"address\": {\"zip\": \"10001\"}}'").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row with zip 10001, got %d", count)
	}

	// Test 9: <@ operator - smaller object is contained by larger
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE '{\"name\": \"Bob\"}' <@ data").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row containing name Bob, got %d", count)
	}

	// Test 10: Array containment
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data @> '{\"tags\": [\"developer\"]}'").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row with developer tag, got %d", count)
	}
}

// TestJSONB_KeyExistence tests ?, ?|, and ?& key existence operators
func TestJSONB_KeyExistence(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 11: ? operator - key exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data ? 'age'").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 rows with 'age' key, got %d", count)
	}

	// Test 12: ? operator - key does not exist
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data ? 'nonexistent'").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 rows with 'nonexistent' key, got %d", count)
	}

	// Test 13: ?| operator - any of the keys exist
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data ?| array['age', 'tags']").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count >= 3 {
		t.Logf("Found %d rows with 'age' or 'tags' key", count)
	} else {
		t.Errorf("Expected at least 3 rows with 'age' or 'tags', got %d", count)
	}

	// Test 14: ?& operator - all keys exist
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data ?& array['name', 'age']").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 rows with both 'name' and 'age' keys, got %d", count)
	}

	// Test 15: ?& operator - not all keys exist
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data ?& array['name', 'nonexistent']").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 rows with both 'name' and 'nonexistent' keys, got %d", count)
	}
}

// TestJSONB_ArrayOperations tests JSONB array operations
func TestJSONB_ArrayOperations(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 16: Access array element by index
	var hobby string
	err := db.QueryRow("SELECT data->'hobbies'->0 FROM jsonb_test WHERE data->>'name' = 'Bob'").Scan(&hobby)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if hobby != `"coding"` {
		t.Errorf("Expected '\"coding\"', got '%s'", hobby)
	}

	// Test 17: Access array element as text
	err = db.QueryRow("SELECT data->'hobbies'->>1 FROM jsonb_test WHERE data->>'name' = 'Bob'").Scan(&hobby)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if hobby != "gaming" {
		t.Errorf("Expected 'gaming', got '%s'", hobby)
	}

	// Test 18: Access nested array element
	var productName string
	err = db.QueryRow("SELECT data->'products'->0->>'name' FROM jsonb_test WHERE data ? 'products'").Scan(&productName)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if productName != "Widget" {
		t.Errorf("Expected 'Widget', got '%s'", productName)
	}
}

// TestJSONB_ArrayLength tests jsonb_array_length() function
func TestJSONB_ArrayLength(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 19: Get array length
	var length int
	err := db.QueryRow("SELECT jsonb_array_length(data->'hobbies') FROM jsonb_test WHERE data->>'name' = 'Bob'").Scan(&length)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if length != 2 {
		t.Errorf("Expected array length 2, got %d", length)
	}

	// Test 20: Get length of products array
	err = db.QueryRow("SELECT jsonb_array_length(data->'products') FROM jsonb_test WHERE data ? 'products'").Scan(&length)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if length != 2 {
		t.Errorf("Expected array length 2, got %d", length)
	}
}

// TestJSONB_TypeOf tests jsonb_typeof() function
func TestJSONB_TypeOf(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 21: Type of string value
	var typeStr string
	err := db.QueryRow("SELECT jsonb_typeof(data->'name') FROM jsonb_test WHERE data->>'name' = 'Alice'").Scan(&typeStr)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if typeStr != "string" {
		t.Errorf("Expected type 'string', got '%s'", typeStr)
	}

	// Test 22: Type of number value
	err = db.QueryRow("SELECT jsonb_typeof(data->'age') FROM jsonb_test WHERE data->>'name' = 'Alice'").Scan(&typeStr)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if typeStr != "number" {
		t.Errorf("Expected type 'number', got '%s'", typeStr)
	}

	// Test 23: Type of array value
	err = db.QueryRow("SELECT jsonb_typeof(data->'hobbies') FROM jsonb_test WHERE data->>'name' = 'Bob'").Scan(&typeStr)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if typeStr != "array" {
		t.Errorf("Expected type 'array', got '%s'", typeStr)
	}

	// Test 24: Type of object value
	err = db.QueryRow("SELECT jsonb_typeof(data->'address') FROM jsonb_test WHERE data->>'name' = 'Charlie'").Scan(&typeStr)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if typeStr != "object" {
		t.Errorf("Expected type 'object', got '%s'", typeStr)
	}

	// Test 25: Type of boolean value
	err = db.QueryRow("SELECT jsonb_typeof(data->'active') FROM jsonb_test WHERE data->>'name' = 'Diana'").Scan(&typeStr)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if typeStr != "boolean" {
		t.Errorf("Expected type 'boolean', got '%s'", typeStr)
	}
}

// TestJSONB_Set tests jsonb_set() function
func TestJSONB_Set(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 26: Update existing top-level field
	var updated string
	err := db.QueryRow(`
		SELECT jsonb_set(data, '{age}', '31') 
		FROM jsonb_test 
		WHERE data->>'name' = 'Alice'
	`).Scan(&updated)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	var updatedData map[string]interface{}
	if err := json.Unmarshal([]byte(updated), &updatedData); err != nil {
		t.Fatalf("Failed to parse updated JSON: %v", err)
	}
	if age, ok := updatedData["age"].(float64); !ok || age != 31 {
		t.Errorf("Expected age to be 31, got %v", updatedData["age"])
	}

	// Test 27: Update nested field
	err = db.QueryRow(`
		SELECT jsonb_set(data, '{address,zip}', '"10002"') 
		FROM jsonb_test 
		WHERE data->>'name' = 'Charlie'
	`).Scan(&updated)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if err := json.Unmarshal([]byte(updated), &updatedData); err != nil {
		t.Fatalf("Failed to parse updated JSON: %v", err)
	}
	if address, ok := updatedData["address"].(map[string]interface{}); !ok {
		t.Error("Expected 'address' to be an object")
	} else if zip, ok := address["zip"].(string); !ok || zip != "10002" {
		t.Errorf("Expected zip to be '10002', got '%v'", address["zip"])
	}

	// Test 28: Add new field with create_missing=true
	err = db.QueryRow(`
		SELECT jsonb_set(data, '{email}', '"alice@example.com"', true) 
		FROM jsonb_test 
		WHERE data->>'name' = 'Alice'
	`).Scan(&updated)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if err := json.Unmarshal([]byte(updated), &updatedData); err != nil {
		t.Fatalf("Failed to parse updated JSON: %v", err)
	}
	if email, ok := updatedData["email"].(string); !ok || email != "alice@example.com" {
		t.Errorf("Expected email to be 'alice@example.com', got '%v'", updatedData["email"])
	}
}

// TestJSONB_WhereClauseFiltering tests JSONB in WHERE clauses
func TestJSONB_WhereClauseFiltering(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)
	insertTestData(t, db)

	// Test 29: Filter by exact JSONB match
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data->>'city' = 'NYC'").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row with city NYC, got %d", count)
	}

	// Test 30: Filter by numeric comparison
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE (data->>'age')::int > 28").Scan(&count)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 rows with age > 28, got %d", count)
	}
}

// TestJSONB_GINIndex tests GIN index support for JSONB
func TestJSONB_GINIndex(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupTestTable(t, db)

	// Test 31: Create GIN index on JSONB column
	_, err := db.Exec("CREATE INDEX idx_jsonb_test_data ON jsonb_test USING GIN (data)")
	if err != nil {
		t.Fatalf("Failed to create GIN index: %v", err)
	}

	// Insert data and verify index is used for containment queries
	insertTestData(t, db)

	// Query using index (containment operator)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM jsonb_test WHERE data @> '{\"name\": \"Alice\"}'").Scan(&count)
	if err != nil {
		t.Fatalf("Query with GIN index failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}

	// Verify index exists
	var indexExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes 
			WHERE tablename = 'jsonb_test' 
			AND indexname = 'idx_jsonb_test_data'
		)
	`).Scan(&indexExists)
	if err != nil {
		t.Fatalf("Failed to check index existence: %v", err)
	}
	if !indexExists {
		t.Error("GIN index was not created")
	}

	t.Logf("✓ GIN index created and functional")
}

// TestMain runs before all tests - can be used for global setup/teardown
func TestMain(m *testing.M) {
	// Check if test database is accessible
	connStr := os.Getenv("VIBESQL_TEST_DB")
	if connStr == "" {
		connStr = defaultTestDB
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to connect to test database: %v\n", err)
		fmt.Fprintf(os.Stderr, "Set VIBESQL_TEST_DB environment variable or ensure PostgreSQL is running at: %s\n", defaultTestDB)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to ping test database: %v\n", err)
		fmt.Fprintf(os.Stderr, "Connection string: %s\n", connStr)
		fmt.Fprintf(os.Stderr, "\nTo run these tests, you need a PostgreSQL instance running.\n")
		fmt.Fprintf(os.Stderr, "Quick start with Docker:\n")
		fmt.Fprintf(os.Stderr, "  docker run --name vibesql-test -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=vibesql_test -p 5432:5432 -d postgres:16\n")
		fmt.Fprintf(os.Stderr, "\nOr set VIBESQL_TEST_DB to point to your PostgreSQL instance:\n")
		fmt.Fprintf(os.Stderr, "  export VIBESQL_TEST_DB='host=localhost port=5432 user=postgres password=yourpassword dbname=vibesql_test sslmode=disable'\n")
		os.Exit(1)
	}

	fmt.Println("✓ Test database connection successful")
	fmt.Printf("  Connection: %s\n", connStr)
	fmt.Println()

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}
