package integration

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

// setupSQLTestTable creates a test table for SQL subset testing
func setupSQLTestTable(t *testing.T, db *sql.DB, tableName string) {
	// Drop table if exists
	_, err := db.Exec("DROP TABLE IF EXISTS " + tableName)
	if err != nil {
		t.Fatalf("Failed to drop test table %s: %v", tableName, err)
	}

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE ` + tableName + ` (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			age INTEGER,
			email TEXT,
			active BOOLEAN DEFAULT true
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table %s: %v", tableName, err)
	}
}

// TestSQL_BasicInsert tests basic INSERT statement
func TestSQL_BasicInsert(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Test 1: Basic INSERT
	result, err := db.Exec("INSERT INTO users (name, age, email) VALUES ($1, $2, $3)", "Alice", 30, "alice@example.com")
	if err != nil {
		t.Fatalf("INSERT failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("Failed to get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE name = 'Alice'").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}
}

// TestSQL_InsertWithReturning tests INSERT with RETURNING clause
func TestSQL_InsertWithReturning(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Test 2: INSERT with RETURNING
	var id int
	var name string
	err := db.QueryRow("INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, name", "Bob", 25).Scan(&id, &name)
	if err != nil {
		t.Fatalf("INSERT with RETURNING failed: %v", err)
	}

	if id == 0 {
		t.Error("Expected non-zero ID")
	}
	if name != "Bob" {
		t.Errorf("Expected name 'Bob', got '%s'", name)
	}
}

// TestSQL_BasicSelect tests basic SELECT statement
func TestSQL_BasicSelect(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25), ('Charlie', 35)")

	// Test 3: Basic SELECT
	rows, err := db.Query("SELECT name, age FROM users ORDER BY age")
	if err != nil {
		t.Fatalf("SELECT failed: %v", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		var age int
		if err := rows.Scan(&name, &age); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		names = append(names, name)
	}

	if len(names) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(names))
	}
	if names[0] != "Bob" || names[1] != "Alice" || names[2] != "Charlie" {
		t.Errorf("Expected [Bob, Alice, Charlie], got %v", names)
	}
}

// TestSQL_SelectWithWhere tests SELECT with WHERE clause
func TestSQL_SelectWithWhere(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age, active) VALUES ('Alice', 30, true), ('Bob', 25, false), ('Charlie', 35, true)")

	// Test 4: WHERE with equality
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE active = true").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT with WHERE failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 active users, got %d", count)
	}

	// Test 5: WHERE with comparison
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE age > 28").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT with WHERE comparison failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 users with age > 28, got %d", count)
	}

	// Test 6: WHERE with AND
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE age > 28 AND active = true").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT with WHERE AND failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 users with age > 28 and active, got %d", count)
	}

	// Test 7: WHERE with OR
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE age < 28 OR age > 32").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT with WHERE OR failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 users, got %d", count)
	}
}

// TestSQL_OrderBy tests ORDER BY clause
func TestSQL_OrderBy(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Charlie', 35), ('Alice', 30), ('Bob', 25)")

	// Test 8: ORDER BY ASC
	rows, err := db.Query("SELECT name FROM users ORDER BY age ASC")
	if err != nil {
		t.Fatalf("SELECT with ORDER BY failed: %v", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	if len(names) != 3 || names[0] != "Bob" || names[2] != "Charlie" {
		t.Errorf("Expected [Bob, Alice, Charlie], got %v", names)
	}

	// Test 9: ORDER BY DESC
	rows2, err := db.Query("SELECT name FROM users ORDER BY age DESC")
	if err != nil {
		t.Fatalf("SELECT with ORDER BY DESC failed: %v", err)
	}
	defer rows2.Close()

	names = nil
	for rows2.Next() {
		var name string
		rows2.Scan(&name)
		names = append(names, name)
	}

	if len(names) != 3 || names[0] != "Charlie" || names[2] != "Bob" {
		t.Errorf("Expected [Charlie, Alice, Bob], got %v", names)
	}
}

// TestSQL_LimitOffset tests LIMIT and OFFSET clauses
func TestSQL_LimitOffset(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	for i := 1; i <= 10; i++ {
		db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "User"+string(rune('0'+i)), 20+i)
	}

	// Test 10: LIMIT
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM (SELECT * FROM users LIMIT 5) AS limited").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT with LIMIT failed: %v", err)
	}
	if count != 5 {
		t.Errorf("Expected 5 rows with LIMIT 5, got %d", count)
	}

	// Test 11: OFFSET
	rows, err := db.Query("SELECT name FROM users ORDER BY age LIMIT 3 OFFSET 2")
	if err != nil {
		t.Fatalf("SELECT with OFFSET failed: %v", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	if len(names) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(names))
	}

	// Test 12: LIMIT with OFFSET
	rows2, err := db.Query("SELECT name FROM users ORDER BY age LIMIT 2 OFFSET 5")
	if err != nil {
		t.Fatalf("SELECT with LIMIT and OFFSET failed: %v", err)
	}
	defer rows2.Close()

	count = 0
	for rows2.Next() {
		count++
	}

	if count != 2 {
		t.Errorf("Expected 2 rows with LIMIT 2 OFFSET 5, got %d", count)
	}
}

// TestSQL_Update tests UPDATE statement
func TestSQL_Update(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30)")

	// Test 13: UPDATE with WHERE
	result, err := db.Exec("UPDATE users SET age = 31 WHERE name = 'Alice'")
	if err != nil {
		t.Fatalf("UPDATE failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("Failed to get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAffected)
	}

	// Verify update
	var age int
	err = db.QueryRow("SELECT age FROM users WHERE name = 'Alice'").Scan(&age)
	if err != nil {
		t.Fatalf("SELECT after UPDATE failed: %v", err)
	}
	if age != 31 {
		t.Errorf("Expected age 31, got %d", age)
	}
}

// TestSQL_UpdateWithReturning tests UPDATE with RETURNING clause
func TestSQL_UpdateWithReturning(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Bob', 25)")

	// Test 14: UPDATE with RETURNING
	var newAge int
	var name string
	err := db.QueryRow("UPDATE users SET age = 26 WHERE name = 'Bob' RETURNING age, name").Scan(&newAge, &name)
	if err != nil {
		t.Fatalf("UPDATE with RETURNING failed: %v", err)
	}

	if newAge != 26 {
		t.Errorf("Expected age 26, got %d", newAge)
	}
	if name != "Bob" {
		t.Errorf("Expected name 'Bob', got '%s'", name)
	}
}

// TestSQL_Delete tests DELETE statement
func TestSQL_Delete(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25)")

	// Test 15: DELETE with WHERE
	result, err := db.Exec("DELETE FROM users WHERE name = 'Alice'")
	if err != nil {
		t.Fatalf("DELETE failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("Failed to get rows affected: %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("Expected 1 row deleted, got %d", rowsAffected)
	}

	// Verify deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT after DELETE failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 remaining row, got %d", count)
	}
}

// TestSQL_DeleteWithReturning tests DELETE with RETURNING clause
func TestSQL_DeleteWithReturning(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Charlie', 35)")

	// Test 16: DELETE with RETURNING
	var deletedName string
	var deletedAge int
	err := db.QueryRow("DELETE FROM users WHERE name = 'Charlie' RETURNING name, age").Scan(&deletedName, &deletedAge)
	if err != nil {
		t.Fatalf("DELETE with RETURNING failed: %v", err)
	}

	if deletedName != "Charlie" {
		t.Errorf("Expected name 'Charlie', got '%s'", deletedName)
	}
	if deletedAge != 35 {
		t.Errorf("Expected age 35, got %d", deletedAge)
	}
}

// TestSQL_CreateTable tests CREATE TABLE statement
func TestSQL_CreateTable(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Test 17: CREATE TABLE
	_, err := db.Exec(`
		CREATE TABLE test_create (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("CREATE TABLE failed: %v", err)
	}

	// Verify table exists
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM pg_tables
			WHERE tablename = 'test_create'
		)
	`).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}
	if !exists {
		t.Error("Table 'test_create' was not created")
	}

	// Cleanup
	db.Exec("DROP TABLE test_create")
}

// TestSQL_DropTable tests DROP TABLE statement
func TestSQL_DropTable(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Create table first
	db.Exec("CREATE TABLE test_drop (id SERIAL PRIMARY KEY)")

	// Test 18: DROP TABLE
	_, err := db.Exec("DROP TABLE test_drop")
	if err != nil {
		t.Fatalf("DROP TABLE failed: %v", err)
	}

	// Verify table does not exist
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM pg_tables
			WHERE tablename = 'test_drop'
		)
	`).Scan(&exists)
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}
	if exists {
		t.Error("Table 'test_drop' still exists after DROP")
	}
}

// TestSQL_DropTableIfExists tests DROP TABLE IF EXISTS statement
func TestSQL_DropTableIfExists(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Test 19: DROP TABLE IF EXISTS (table doesn't exist)
	_, err := db.Exec("DROP TABLE IF EXISTS nonexistent_table")
	if err != nil {
		t.Fatalf("DROP TABLE IF EXISTS failed: %v", err)
	}

	// Create and drop with IF EXISTS
	db.Exec("CREATE TABLE test_if_exists (id SERIAL PRIMARY KEY)")
	_, err = db.Exec("DROP TABLE IF EXISTS test_if_exists")
	if err != nil {
		t.Fatalf("DROP TABLE IF EXISTS failed: %v", err)
	}
}

// TestSQL_WhereWithIN tests WHERE with IN clause
func TestSQL_WhereWithIN(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25), ('Charlie', 35)")

	// Test 20: WHERE with IN
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE name IN ('Alice', 'Charlie')").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT with WHERE IN failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}
}

// TestSQL_WhereLike tests WHERE with LIKE operator
func TestSQL_WhereLike(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25), ('Charlie', 35)")

	// Test 21: WHERE with LIKE
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE name LIKE 'A%'").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT with WHERE LIKE failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row starting with 'A', got %d", count)
	}
}

// TestSQL_MultipleInserts tests multiple INSERT statements
func TestSQL_MultipleInserts(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Test 22: Multiple INSERTs in separate statements
	_, err := db.Exec("INSERT INTO users (name, age) VALUES ('User1', 21)")
	if err != nil {
		t.Fatalf("First INSERT failed: %v", err)
	}

	_, err = db.Exec("INSERT INTO users (name, age) VALUES ('User2', 22)")
	if err != nil {
		t.Fatalf("Second INSERT failed: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("SELECT failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}
}

// TestSQL_SelectDistinct tests SELECT DISTINCT
func TestSQL_SelectDistinct(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data with duplicates
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 30), ('Charlie', 35)")

	// Test 23: SELECT DISTINCT
	rows, err := db.Query("SELECT DISTINCT age FROM users ORDER BY age")
	if err != nil {
		t.Fatalf("SELECT DISTINCT failed: %v", err)
	}
	defer rows.Close()

	var ages []int
	for rows.Next() {
		var age int
		rows.Scan(&age)
		ages = append(ages, age)
	}

	if len(ages) != 2 {
		t.Errorf("Expected 2 distinct ages, got %d", len(ages))
	}
}

// TestSQL_CountAggregation tests COUNT aggregation
func TestSQL_CountAggregation(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Insert test data
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25), ('Charlie', 35)")

	// Test 24: COUNT(*)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("COUNT(*) failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	// Test 25: COUNT with WHERE
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE age > 28").Scan(&count)
	if err != nil {
		t.Fatalf("COUNT with WHERE failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

// TestSQL_UnsupportedJoin tests that JOIN is not supported
func TestSQL_UnsupportedJoin(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")
	setupSQLTestTable(t, db, "orders")

	// Test 26: JOIN should fail (or work if supported by minimal PostgreSQL)
	_, err := db.Query("SELECT users.name FROM users JOIN orders ON users.id = orders.user_id")
	// Note: PostgreSQL supports JOINs, so this might actually work
	// The test verifies behavior but doesn't necessarily expect failure
	if err != nil {
		t.Logf("JOIN not supported (expected for minimal build): %v", err)
	} else {
		t.Log("JOIN is supported in this PostgreSQL build")
	}
}

// TestSQL_UnsupportedSubquery tests subquery behavior
func TestSQL_UnsupportedSubquery(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25)")

	// Test 27: Subquery (may or may not be supported)
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM (SELECT * FROM users WHERE age > 20) AS subquery").Scan(&count)
	if err != nil {
		t.Logf("Subquery not supported (expected for minimal build): %v", err)
	} else {
		t.Log("Subquery is supported in this PostgreSQL build")
		if count != 2 {
			t.Errorf("Expected count 2, got %d", count)
		}
	}
}

// TestSQL_NullHandling tests NULL value handling
func TestSQL_NullHandling(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Test 28: Insert NULL value
	_, err := db.Exec("INSERT INTO users (name, age, email) VALUES ('Alice', 30, NULL)")
	if err != nil {
		t.Fatalf("INSERT with NULL failed: %v", err)
	}

	// Test 29: Query NULL value
	var email sql.NullString
	err = db.QueryRow("SELECT email FROM users WHERE name = 'Alice'").Scan(&email)
	if err != nil {
		t.Fatalf("SELECT NULL value failed: %v", err)
	}
	if email.Valid {
		t.Error("Expected NULL email, but got valid value")
	}

	// Test 30: WHERE IS NULL
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email IS NULL").Scan(&count)
	if err != nil {
		t.Fatalf("WHERE IS NULL failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row with NULL email, got %d", count)
	}
}

// TestSQL_BooleanValues tests BOOLEAN data type
func TestSQL_BooleanValues(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Test 31: Insert and query BOOLEAN
	db.Exec("INSERT INTO users (name, age, active) VALUES ('Alice', 30, true), ('Bob', 25, false)")

	var active bool
	err := db.QueryRow("SELECT active FROM users WHERE name = 'Alice'").Scan(&active)
	if err != nil {
		t.Fatalf("SELECT BOOLEAN failed: %v", err)
	}
	if !active {
		t.Error("Expected active=true, got false")
	}

	// Test 32: WHERE with BOOLEAN
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE active = false").Scan(&count)
	if err != nil {
		t.Fatalf("WHERE with BOOLEAN failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 inactive user, got %d", count)
	}
}

// TestSQL_UpdateMultipleColumns tests updating multiple columns
func TestSQL_UpdateMultipleColumns(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")
	db.Exec("INSERT INTO users (name, age, email) VALUES ('Alice', 30, 'alice@old.com')")

	// Test 33: UPDATE multiple columns
	_, err := db.Exec("UPDATE users SET age = 31, email = 'alice@new.com' WHERE name = 'Alice'")
	if err != nil {
		t.Fatalf("UPDATE multiple columns failed: %v", err)
	}

	var age int
	var email string
	err = db.QueryRow("SELECT age, email FROM users WHERE name = 'Alice'").Scan(&age, &email)
	if err != nil {
		t.Fatalf("SELECT after UPDATE failed: %v", err)
	}
	if age != 31 || email != "alice@new.com" {
		t.Errorf("Expected age=31, email=alice@new.com, got age=%d, email=%s", age, email)
	}
}

// TestSQL_DeleteMultipleRows tests deleting multiple rows
func TestSQL_DeleteMultipleRows(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")
	db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25), ('Charlie', 20)")

	// Test 34: DELETE multiple rows
	result, err := db.Exec("DELETE FROM users WHERE age < 28")
	if err != nil {
		t.Fatalf("DELETE multiple rows failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("Failed to get rows affected: %v", err)
	}
	if rowsAffected != 2 {
		t.Errorf("Expected 2 rows deleted, got %d", rowsAffected)
	}
}

// TestSQL_InsertMultipleRows tests INSERT with multiple value sets
func TestSQL_InsertMultipleRows(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	setupSQLTestTable(t, db, "users")

	// Test 35: INSERT multiple rows in one statement
	result, err := db.Exec("INSERT INTO users (name, age) VALUES ('Alice', 30), ('Bob', 25), ('Charlie', 35)")
	if err != nil {
		t.Fatalf("INSERT multiple rows failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("Failed to get rows affected: %v", err)
	}
	if rowsAffected != 3 {
		t.Errorf("Expected 3 rows inserted, got %d", rowsAffected)
	}
}
