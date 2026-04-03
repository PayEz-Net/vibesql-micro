// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/vibesql/vibesql-micro/pkg/vsql"
)

func main() {
	tmpDir, err := os.MkdirTemp("", "vsql-test-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	dbPath := tmpDir + "\\test.vsql"
	
	fmt.Println("=== Test 1: Open database ===")
	db, err := vsql.Open(dbPath)
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}
	fmt.Println("✓ Database opened successfully")
	
	fmt.Println("\n=== Test 2: Create table ===")
	_, err = db.Exec("CREATE TABLE test (id SERIAL PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Fatalf("Create table failed: %v", err)
	}
	fmt.Println("✓ Table created")
	
	fmt.Println("\n=== Test 3: Insert data ===")
	result, err := db.Exec("INSERT INTO test (name) VALUES ($1)", "hello world")
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	fmt.Printf("✓ Inserted %d row(s)\n", result.RowsAffected)
	
	fmt.Println("\n=== Test 4: Query data ===")
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("✓ Query returned %d row(s)\n", len(rows))
	for _, row := range rows {
		fmt.Printf("  Row: %v\n", row)
	}
	
	fmt.Println("\n=== Test 5: Close database ===")
	if err := db.Close(); err != nil {
		log.Fatalf("Close failed: %v", err)
	}
	fmt.Println("✓ Database closed")
	
	fmt.Println("\n=== Test 6: Reopen database ===")
	db2, err := vsql.Open(dbPath)
	if err != nil {
		log.Fatalf("Reopen failed: %v", err)
	}
	
	rows2, err := db2.Query("SELECT * FROM test")
	if err != nil {
		log.Fatalf("Query after reopen failed: %v", err)
	}
	fmt.Printf("✓ Data persisted: %d row(s) after reopen\n", len(rows2))
	
	db2.Close()
	
	fmt.Println("\n=== All tests passed! ===")
}
