package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/vibesql/vibesql-micro/pkg/vsql"
)

var version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}
	
	command := os.Args[1]
	
	switch command {
	case "version", "-v", "--version":
		printVersion()
		
	case "help", "-h", "--help":
		printUsage()
		
	default:
		// Assume first arg is database path
		dbPath := command
		
		if len(os.Args) > 2 {
			// Single query mode
			sql := os.Args[2]
			runQuery(dbPath, sql)
		} else {
			// Interactive shell mode
			runShell(dbPath)
		}
	}
}

func printVersion() {
	fmt.Println("vsql-micro", version)
	fmt.Println("PostgreSQL that feels like SQLite")
}

func printUsage() {
	fmt.Println("vsql-micro - PostgreSQL that feels like SQLite")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vsql-micro <file.vsql>           # Open interactive shell")
	fmt.Println("  vsql-micro <file.vsql> <query>   # Run single query")
	fmt.Println("  vsql-micro version               # Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vsql-micro ./app.vsql")
	fmt.Println("  vsql-micro ./app.vsql 'SELECT 1'")
	fmt.Println("  vsql-micro ./app.vsql 'CREATE TABLE users (id SERIAL, name TEXT)'")
	fmt.Println()
	fmt.Println("From Go:")
	fmt.Println("  import \"github.com/vibesql/vibesql-micro/pkg/vsql\"")
	fmt.Println("  db, _ := vsql.Open(\"./app.vsql\")")
	fmt.Println("  rows, _ := db.Query(\"SELECT * FROM users\")")
}

func runQuery(dbPath string, sql string) {
	// Progress callback for stderr
	progress := func(msg string) {
		if msg == "done" {
			fmt.Fprintln(os.Stderr, " done")
		} else {
			fmt.Fprintf(os.Stderr, "%s... ", msg)
		}
	}
	
	db, err := vsql.OpenWithProgress(dbPath, progress)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	
	// Execute query and output JSON
	jsonResult, err := db.QueryJSON(sql)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	// Pretty print JSON
	var prettyData interface{}
	if err := json.Unmarshal([]byte(jsonResult), &prettyData); err == nil {
		prettyJSON, _ := json.MarshalIndent(prettyData, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println(jsonResult)
	}
}

func runShell(dbPath string) {
	// Progress callback
	progress := func(msg string) {
		if msg == "done" {
			fmt.Fprintln(os.Stderr, " done")
		} else {
			fmt.Fprintf(os.Stderr, "%s... ", msg)
		}
	}
	
	db, err := vsql.OpenWithProgress(dbPath, progress)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	
	// Show prompt
	base := filepath.Base(dbPath)
	fmt.Fprintf(os.Stderr, "\nConnected to %s\n", base)
	fmt.Fprintln(os.Stderr, "Type SQL queries, or \\q to quit\n")
	
	// Simple REPL
	// TODO: Implement proper readline with history
	for {
		fmt.Print(base + "> ")
		
		// Read line (simple version)
		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			break
		}
		
		// Handle commands
		if line == "\\q" || line == "quit" || line == "exit" {
			fmt.Fprintln(os.Stderr, "Bye!")
			break
		}
		
		// Execute query
		jsonResult, err := db.QueryJSON(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		
		// Output result
		var prettyData interface{}
		if err := json.Unmarshal([]byte(jsonResult), &prettyData); err == nil {
			prettyJSON, _ := json.MarshalIndent(prettyData, "", "  ")
			fmt.Println(string(prettyJSON))
		} else {
			fmt.Println(jsonResult)
		}
	}
}
