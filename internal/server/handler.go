package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/vibesql/vibe/internal/postgres"
	"github.com/vibesql/vibe/internal/query"
)

type Handler struct {
	executor query.QueryExecutor
}

func NewHandler(executor query.QueryExecutor) *Handler {
	return &Handler{
		executor: executor,
	}
}

func (h *Handler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := NewInvalidSQLError("Only POST method is supported for /v1/query endpoint")
		WriteError(w, err)
		log.Printf("[ERROR] Method not allowed: %s %s", r.Method, r.URL.Path)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		vibeErr := NewInternalError("Failed to read request body: " + err.Error())
		WriteError(w, vibeErr)
		log.Printf("[ERROR] Failed to read request body: %v", err)
		return
	}

	var req QueryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		vibeErr := NewInvalidSQLError("Invalid JSON request body")
		WriteError(w, vibeErr)
		log.Printf("[ERROR] Invalid JSON: %v", err)
		return
	}

	if req.SQL == "" {
		vibeErr := NewMissingFieldError("sql")
		WriteError(w, vibeErr)
		log.Printf("[ERROR] Missing required field: sql")
		return
	}

	log.Printf("[INFO] Executing query: %.100s...", req.SQL)

	if err := query.ValidateQuery(req.SQL); err != nil {
		if vibeErr, ok := err.(*postgres.VibeError); ok {
			WriteError(w, vibeErr)
		} else {
			WriteError(w, NewInternalError(err.Error()))
		}
		log.Printf("[ERROR] Query validation failed: %v", err)
		return
	}

	if err := query.CheckSafety(req.SQL); err != nil {
		if vibeErr, ok := err.(*postgres.VibeError); ok {
			WriteError(w, vibeErr)
		} else {
			WriteError(w, NewInternalError(err.Error()))
		}
		log.Printf("[ERROR] Query safety check failed: %v", err)
		return
	}

	result, err := h.executor.Execute(req.SQL)
	if err != nil {
		if vibeErr, ok := err.(*postgres.VibeError); ok {
			WriteError(w, vibeErr)
		} else {
			WriteError(w, NewInternalError(err.Error()))
		}
		log.Printf("[ERROR] Query execution failed: %v", err)
		return
	}

	executionTimeMs := float64(result.ExecutionTime.Microseconds()) / 1000.0

	if err := WriteSuccess(w, result.Rows, executionTimeMs); err != nil {
		log.Printf("[ERROR] Failed to write response: %v", err)
		return
	}

	log.Printf("[INFO] Query succeeded: %d rows returned in %.2fms", result.RowCount, executionTimeMs)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/query", h.HandleQuery)
}
