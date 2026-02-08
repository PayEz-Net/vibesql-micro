.PHONY: all build build-linux build-darwin build-windows remote-build clean test lint docker-build build-postgres build-postgres-darwin build-postgres-windows help

BINARY_NAME := vibe
BUILD_DIR := build
OUTPUT_DIR := $(BUILD_DIR)/output
PG_VERSION := 16.1

GO := go
GOOS := $(shell go env GOOS 2>/dev/null || echo linux)
GOARCH := $(shell go env GOARCH 2>/dev/null || echo amd64)

VERSION ?= 1.0.0
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S' 2>/dev/null || echo "unknown")

LDFLAGS := -s -w \
	-X github.com/vibesql/vibe/internal/version.Version=$(VERSION) \
	-X github.com/vibesql/vibe/internal/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/vibesql/vibe/internal/version.BuildDate=$(BUILD_DATE)
GOFLAGS := -ldflags="$(LDFLAGS)"

GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m

all: build ## Build everything (default target)

help: ## Show this help message
	@echo "VibeSQL Local - Makefile Commands"
	@echo "=================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-24s$(NC) %s\n", $$1, $$2}'

build-postgres: ## Build minimal PostgreSQL binary (Linux)
	@echo "$(YELLOW)Building PostgreSQL micro binary (Linux)...$(NC)"
	@bash $(BUILD_DIR)/build.sh

build-postgres-darwin: ## Build minimal PostgreSQL binary (macOS)
	@echo "$(YELLOW)Building PostgreSQL micro binary (macOS)...$(NC)"
	@bash $(BUILD_DIR)/build_postgres_darwin.sh

build-postgres-windows: ## Build/download PostgreSQL binaries (Windows) -- run from PowerShell
	@echo "$(YELLOW)Use PowerShell: .\\build\\build_postgres_windows.ps1$(NC)"

build: build-postgres ## Build complete VibeSQL binary (current platform)
	@echo "$(YELLOW)Building VibeSQL Go binary...$(NC)"
	@$(GO) build $(GOFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/vibe
	@echo "$(GREEN)✓ Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)$(NC)"

docker-build: ## Build using Docker (reproducible builds)
	@echo "$(YELLOW)Building PostgreSQL with Docker...$(NC)"
	@docker build -f $(BUILD_DIR)/Dockerfile.postgres -t vibesql-postgres:micro .
	@echo "$(GREEN)✓ Docker build complete$(NC)"

clean: ## Clean build artifacts (Unix only; on Windows use: del /S /Q build\output)
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@rm -f internal/postgres/embed/postgres_micro_*
	@rm -f internal/postgres/embed/initdb_*
	@rm -f internal/postgres/embed/pg_ctl_*
	@rm -f internal/postgres/embed/libpq*
	@$(GO) clean
	@echo "$(GREEN)✓ Clean complete$(NC)"

test: ## Run all tests
	@echo "Running tests..."
	@$(GO) test ./... -v -race -cover

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@$(GO) test ./tests/unit/... -v -race

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@$(GO) test ./tests/integration/... -v -timeout=60s

test-e2e: ## Run end-to-end tests
	@echo "Running E2E tests..."
	@$(GO) test ./tests/e2e/... -v

test-coverage: ## Generate test coverage report
	@echo "Generating coverage report..."
	@$(GO) test ./... -coverprofile=coverage.out
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

lint: ## Run linters
	@echo "Running linters..."
	@$(GO) vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not installed, skipping$(NC)"; \
	fi
	@echo "$(GREEN)✓ Lint complete$(NC)"

fmt: ## Format Go code
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@echo "$(GREEN)✓ Format complete$(NC)"

run: build ## Build and run VibeSQL server
	@echo "Starting VibeSQL server..."
	@$(OUTPUT_DIR)/$(BINARY_NAME) serve

install: build ## Install binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(OUTPUT_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)✓ Installed$(NC)"

build-linux: ## Build for Linux x64
	@echo "$(YELLOW)Building VibeSQL for Linux x64...$(NC)"
	@bash scripts/build-linux.sh
	@echo "$(GREEN)✓ Linux x64 build complete$(NC)"

build-darwin: ## Build for macOS (auto-detect amd64/arm64)
	@echo "$(YELLOW)Building VibeSQL for macOS...$(NC)"
	@bash scripts/build-darwin.sh
	@echo "$(GREEN)✓ macOS build complete$(NC)"

build-darwin-arm64: ## Build for macOS Apple Silicon (arm64)
	@echo "$(YELLOW)Building VibeSQL for macOS arm64...$(NC)"
	@GOARCH=arm64 bash scripts/build-darwin.sh
	@echo "$(GREEN)✓ macOS arm64 build complete$(NC)"

build-darwin-amd64: ## Build for macOS Intel (amd64)
	@echo "$(YELLOW)Building VibeSQL for macOS amd64...$(NC)"
	@GOARCH=amd64 bash scripts/build-darwin.sh
	@echo "$(GREEN)✓ macOS amd64 build complete$(NC)"

build-windows: ## Build for Windows x64 -- run from PowerShell
	@echo "$(YELLOW)Use PowerShell: .\\scripts\\build-windows.ps1$(NC)"

remote-build: ## Build on remote Linux dev box (10.0.0.93)
	@echo "$(YELLOW)Building on remote Linux box...$(NC)"
	@bash scripts/remote-build.sh
	@echo "$(GREEN)✓ Remote build complete$(NC)"

size-check: ## Check binary sizes
	@echo "Binary sizes:"
	@ls -lh vibe-* 2>/dev/null || ls -lh $(OUTPUT_DIR)/$(BINARY_NAME)* 2>/dev/null || echo "No binaries found. Run 'make build-linux' first."

verify: test lint ## Run all verification checks

release: clean build-linux ## Build release binary for Linux x64
	@sha256sum vibe-linux-amd64 > checksums.txt
	@echo "$(GREEN)✓ Release ready$(NC)"
	@cat checksums.txt

release-all: clean build-linux build-darwin ## Build release binaries for all Unix platforms
	@sha256sum vibe-linux-amd64 vibe-darwin-* > checksums.txt 2>/dev/null || true
	@echo "$(GREEN)✓ Release ready (Linux + macOS)$(NC)"
	@cat checksums.txt

.DEFAULT_GOAL := help
