// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vibesql/vibesql-micro/pkg/vsql"
)

var passed, failed int

func main() {
	tmpDir, err := os.MkdirTemp("", "vsql-comprehensive-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Println("========================================")
	fmt.Println("vsql-micro-v2 Comprehensive Linux Test Suite")
	fmt.Println("Test Dir:", tmpDir)
	fmt.Println("========================================")

	// Group 1: Basic Lifecycle
	test("Open fresh database", func() error {
		dbPath := filepath.Join(tmpDir, "t1.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		if _, err := os.Stat(dbPath); err != nil {
			return fmt.Errorf("marker file not created")
		}
		if _, err := os.Stat(dbPath + ".data"); err != nil {
			return fmt.Errorf("data dir not created")
		}
		return nil
	})

	test("Open with empty path defaults to default.vsql", func() error {
		wd := filepath.Join(tmpDir, "empty_test")
		os.MkdirAll(wd, 0755)
		os.Chdir(wd)
		defer os.Chdir(tmpDir)
		db, err := vsql.Open("")
		if err != nil {
			return err
		}
		defer db.Close()
		if _, err := os.Stat("default.vsql"); err != nil {
			return fmt.Errorf("default.vsql not created")
		}
		return nil
	})

	test("Open adds .vsql extension", func() error {
		dbPath := filepath.Join(tmpDir, "noext")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		if _, err := os.Stat(dbPath + ".vsql"); err != nil {
			return fmt.Errorf(".vsql extension not added")
		}
		return nil
	})

	test("Reopen existing database", func() error {
		dbPath := filepath.Join(tmpDir, "reopen.vsql")
		db1, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		db1.Exec("CREATE TABLE reopen_test (id INT)")
		db1.Close()

		db2, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db2.Close()
		rows, err := db2.Query("SELECT * FROM reopen_test")
		if err != nil {
			return err
		}
		if len(rows) != 0 {
			return fmt.Errorf("expected 0 rows, got %d", len(rows))
		}
		return nil
	})

	test("Lock prevents concurrent open", func() error {
		dbPath := filepath.Join(tmpDir, "locked.vsql")
		db1, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db1.Close()

		_, err = vsql.Open(dbPath)
		if err == nil {
			return fmt.Errorf("expected error for locked database")
		}
		if !strings.Contains(err.Error(), "busy") && !strings.Contains(err.Error(), "use") && !strings.Contains(err.Error(), "lock") {
			return fmt.Errorf("expected lock-related error, got: %v", err)
		}
		return nil
	})

	// Group 2: Query/Exec
	test("CREATE TABLE", func() error {
		dbPath := filepath.Join(tmpDir, "query.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)")
		return err
	})

	test("INSERT with parameters", func() error {
		dbPath := filepath.Join(tmpDir, "insert.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)")
		res, err := db.Exec("INSERT INTO users (name) VALUES ($1), ($2)", "Alice", "Bob")
		if err != nil {
			return err
		}
		if res.RowsAffected != 2 {
			return fmt.Errorf("expected 2 rows affected, got %d", res.RowsAffected)
		}
		return nil
	})

	test("SELECT returns correct data", func() error {
		dbPath := filepath.Join(tmpDir, "select.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)")
		db.Exec("INSERT INTO users (name) VALUES ($1)", "Charlie")
		rows, err := db.Query("SELECT * FROM users WHERE name = $1", "Charlie")
		if err != nil {
			return err
		}
		if len(rows) != 1 {
			return fmt.Errorf("expected 1 row, got %d", len(rows))
		}
		if rows[0]["name"] != "Charlie" {
			return fmt.Errorf("expected Charlie, got %v", rows[0]["name"])
		}
		return nil
	})

	test("QueryJSON returns valid JSON", func() error {
		dbPath := filepath.Join(tmpDir, "json.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)")
		db.Exec("INSERT INTO users (name) VALUES ($1)", "Dave")
		jsonStr, err := db.QueryJSON("SELECT * FROM users WHERE name = $1", "Dave")
		if err != nil {
			return err
		}
		var result []map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			return fmt.Errorf("invalid JSON: %v", err)
		}
		if len(result) != 1 {
			return fmt.Errorf("expected 1 result, got %d", len(result))
		}
		return nil
	})

	test("JOIN and aggregate query", func() error {
		dbPath := filepath.Join(tmpDir, "join.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)")
		db.Exec("CREATE TABLE orders (id SERIAL, user_id INT, amount DECIMAL)")
		db.Exec("INSERT INTO users (name) VALUES ($1), ($2)", "Alice", "Bob")
		db.Exec("INSERT INTO orders (user_id, amount) VALUES (1, 100), (1, 200), (2, 50)")
		rows, err := db.Query("SELECT u.name, SUM(o.amount) as total FROM users u JOIN orders o ON u.id = o.user_id GROUP BY u.name ORDER BY total DESC")
		if err != nil {
			return err
		}
		if len(rows) != 2 {
			return fmt.Errorf("expected 2 rows, got %d", len(rows))
		}
		return nil
	})

	test("Query after close returns error", func() error {
		dbPath := filepath.Join(tmpDir, "closed.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		db.Close()
		_, err = db.Query("SELECT 1")
		if err == nil {
			return fmt.Errorf("expected error after close")
		}
		return nil
	})

	// Group 3: Data Types
	test("JSONB operations", func() error {
		dbPath := filepath.Join(tmpDir, "jsonb.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE docs (id SERIAL, data JSONB)")
		db.Exec("INSERT INTO docs (data) VALUES ($1)", `{"name": "test", "active": true}`)
		rows, err := db.Query("SELECT data->>$1 as name FROM docs WHERE data @> $2", "name", `{"active": true}`)
		if err != nil {
			return err
		}
		if len(rows) != 1 {
			return fmt.Errorf("expected 1 row, got %d", len(rows))
		}
		if rows[0]["name"] != "test" {
			return fmt.Errorf("expected test, got %v", rows[0]["name"])
		}
		return nil
	})

	test("Unicode and special characters", func() error {
		dbPath := filepath.Join(tmpDir, "unicode.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE special (id SERIAL, data TEXT)")
		db.Exec("INSERT INTO special (data) VALUES ($1), ($2), ($3)", "日本語", "🎉", "'quotes'")
		rows, err := db.Query("SELECT * FROM special ORDER BY id")
		if err != nil {
			return err
		}
		if len(rows) != 3 {
			return fmt.Errorf("expected 3 rows, got %d", len(rows))
		}
		found := false
		for _, row := range rows {
			if row["data"] == "日本語" {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("unicode data not preserved correctly")
		}
		return nil
	})

	test("Large result set (1000 rows)", func() error {
		dbPath := filepath.Join(tmpDir, "large.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE large (id SERIAL)")
		db.Exec("INSERT INTO large SELECT generate_series(1, 1000)")
		rows, err := db.Query("SELECT COUNT(*) as c FROM large")
		if err != nil {
			return err
		}
		count := fmt.Sprintf("%v", rows[0]["c"])
		if count != "1000" {
			return fmt.Errorf("expected 1000, got %s", count)
		}
		return nil
	})

	// Group 4: Concurrency
	test("Concurrent queries", func() error {
		dbPath := filepath.Join(tmpDir, "concurrent.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE concurrent (id SERIAL)")
		db.Exec("INSERT INTO concurrent SELECT generate_series(1, 100)")

		var wg sync.WaitGroup
		errors := make(chan error, 10)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := db.Query("SELECT COUNT(*) FROM concurrent")
				if err != nil {
					errors <- err
				}
			}()
		}
		wg.Wait()
		close(errors)
		for err := range errors {
			return err
		}
		return nil
	})

	// Group 5: Error Handling
	test("Syntax error returns error", func() error {
		dbPath := filepath.Join(tmpDir, "error.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.Query("INVALID SYNTAX HERE")
		if err == nil {
			return fmt.Errorf("expected syntax error")
		}
		return nil
	})

	test("Wrong parameter count returns error", func() error {
		dbPath := filepath.Join(tmpDir, "param_err.vsql")
		db, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		db.Exec("CREATE TABLE params (id SERIAL, name TEXT)")
		_, err = db.Query("SELECT * FROM params WHERE name = $1 AND id = $2", "only_one")
		if err == nil {
			return fmt.Errorf("expected parameter error")
		}
		return nil
	})

	// Group 6: Persistence
	test("Data persists across close and reopen", func() error {
		dbPath := filepath.Join(tmpDir, "persist.vsql")
		db1, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		db1.Exec("CREATE TABLE persist (id SERIAL, data TEXT)")
		db1.Exec("INSERT INTO persist (data) VALUES ($1)", "survive")
		db1.Close()

		db2, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		defer db2.Close()
		rows, err := db2.Query("SELECT data FROM persist")
		if err != nil {
			return err
		}
		if len(rows) != 1 || rows[0]["data"] != "survive" {
			return fmt.Errorf("data did not persist")
		}
		return nil
	})

	// Group 7: Progress callback
	test("Progress callback fires", func() error {
		dbPath := filepath.Join(tmpDir, "progress.vsql")
		progressMsgs := []string{}
		db, err := vsql.OpenWithProgress(dbPath, func(msg string) {
			progressMsgs = append(progressMsgs, msg)
		})
		if err != nil {
			return err
		}
		defer db.Close()
		if len(progressMsgs) == 0 {
			return fmt.Errorf("expected progress messages")
		}
		return nil
	})

	// Group 8: Performance
	test("Open warm database is fast", func() error {
		dbPath := filepath.Join(tmpDir, "perf.vsql")
		db1, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		db1.Close()

		start := time.Now()
		db2, err := vsql.Open(dbPath)
		if err != nil {
			return err
		}
		db2.Close()
		elapsed := time.Since(start)
		if elapsed > 2*time.Second {
			return fmt.Errorf("warm open took %v, expected < 2s", elapsed)
		}
		return nil
	})

	// Summary
	fmt.Println("========================================")
	fmt.Printf("PASSED: %d\n", passed)
	fmt.Printf("FAILED: %d\n", failed)
	fmt.Println("========================================")
	if failed > 0 {
		os.Exit(1)
	}
}

func test(name string, fn func() error) {
	fmt.Printf("[TEST] %s ... ", name)
	if err := fn(); err != nil {
		fmt.Printf("FAIL\n       %v\n", err)
		failed++
	} else {
		fmt.Println("PASS")
		passed++
	}
}
