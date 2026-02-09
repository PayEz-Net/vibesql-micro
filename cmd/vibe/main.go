package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vibesql/vibe/internal/postgres"
	"github.com/vibesql/vibe/internal/query"
	"github.com/vibesql/vibe/internal/server"
	"github.com/vibesql/vibe/internal/version"
)

const (
	usageText = `VibeSQL - Lightweight PostgreSQL HTTP API

Usage:
  vibe <command> [options]

Commands:
  serve      Start the HTTP server and embedded PostgreSQL
  version    Print version information
  help       Display this help message

Examples:
  vibe serve           Start server on 127.0.0.1:5173
  vibe version         Show version and build info
  vibe help            Show this help

For more information, visit: https://vibesql.dev
`
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "serve":
		if err := runServe(); err != nil {
			log.Printf("[FATAL] %v", err)
			os.Exit(1)
		}
	case "version":
		printVersion()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func runServe() error {
	log.Printf("[INFO] Starting VibeSQL %s", version.Get().Short())

	pgManager := postgres.NewManager("", 5433) // Use 5433 to avoid conflict with system PostgreSQL
	
	log.Printf("[INFO] Starting PostgreSQL...")
	startTime := time.Now()
	if err := pgManager.Start(); err != nil {
		return fmt.Errorf("failed to start PostgreSQL: %w", err)
	}
	defer func() {
		log.Printf("[INFO] Stopping PostgreSQL...")
		if err := pgManager.Stop(); err != nil {
			log.Printf("[ERROR] Failed to stop PostgreSQL: %v", err)
		}
	}()
	
	pgStartupTime := time.Since(startTime)
	log.Printf("[INFO] PostgreSQL started in %v", pgStartupTime)

	conn, err := pgManager.CreateConnection()
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("[ERROR] Failed to close database connection: %v", err)
		}
	}()

	executor := query.NewExecutor(conn.DB())

	httpServer := server.NewServer(executor)
	
	log.Printf("[INFO] Starting HTTP server...")
	if err := httpServer.Start(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	defer func() {
		if err := httpServer.Stop(); err != nil {
			log.Printf("[ERROR] Failed to stop HTTP server: %v", err)
		}
	}()

	totalStartupTime := time.Since(startTime)
	log.Printf("[INFO] VibeSQL ready in %v", totalStartupTime)
	log.Printf("[INFO] HTTP API: http://%s", httpServer.Addr())
	log.Printf("[INFO] Press Ctrl+C to stop")

	httpServer.WaitForShutdown()
	
	log.Printf("[INFO] Shutdown complete")
	return nil
}

func printVersion() {
	info := version.Get()
	fmt.Println(info.Full())
}

func printUsage() {
	fmt.Print(usageText)
}
