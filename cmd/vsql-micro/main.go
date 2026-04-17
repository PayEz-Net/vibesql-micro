package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/vibesql/vibesql-micro/pkg/vsql"
)

var version = "0.3.0"

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

	case "serve":
		runServe(os.Args[2:])

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

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	listen := fs.String("listen", "127.0.0.1:5433", "address to bind postgres on (host:port or :port)")
	data := fs.String("data", "./vault.vsql", "path to the .vsql database file")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "vsql-micro serve — run a long-lived embedded postgres on a pinned port")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  vsql-micro serve [--listen 127.0.0.1:5433] [--data ./vault.vsql]")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Authentication: trust (pod-internal use only).")
		fmt.Fprintln(os.Stderr, "Graceful shutdown on SIGINT or SIGTERM.")
	}
	_ = fs.Parse(args)

	host, portStr, err := splitHostPort(*listen)
	if err != nil {
		fatal("serve", err)
	}
	if host != "" && host != "127.0.0.1" && host != "localhost" {
		fmt.Fprintf(os.Stderr, "warning: vsql-micro binds 127.0.0.1 only; ignoring host %q from --listen\n", host)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fatal("serve: parse port", err)
	}

	progress := func(msg string) {
		if msg == "done" {
			fmt.Fprintln(os.Stderr, " done")
		} else {
			fmt.Fprintf(os.Stderr, "%s... ", msg)
		}
	}

	db, err := vsql.OpenOnPort(*data, port, progress)
	if err != nil {
		fatal("serve: open", err)
	}

	absData, _ := filepath.Abs(*data)
	fmt.Fprintf(os.Stderr, "\nvsql-micro serving %s on 127.0.0.1:%d (user=postgres, trust auth)\n", absData, port)
	fmt.Fprintln(os.Stderr, "Ctrl+C or SIGTERM to shut down.")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	sig := <-sigs
	fmt.Fprintf(os.Stderr, "\nreceived %s, shutting down...\n", sig)

	if err := db.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "close:", err)
	}
}

func splitHostPort(addr string) (string, string, error) {
	if strings.HasPrefix(addr, ":") {
		return "", addr[1:], nil
	}
	idx := strings.LastIndex(addr, ":")
	if idx < 0 {
		return "", "", fmt.Errorf("invalid listen addr %q (want host:port or :port)", addr)
	}
	return addr[:idx], addr[idx+1:], nil
}

func fatal(ctx string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", ctx, err)
	os.Exit(1)
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
	fmt.Println("  vsql-micro serve [flags]         # Run as a long-lived server on a pinned port")
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
