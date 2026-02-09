# VibeSQL Local v1 - Technical Specification

**Version:** 1.0.0  
**Status:** Draft  
**Based on:** requirements.md v1.0.0  
**Author:** Technical Specification Phase  
**Date:** 2026-02-05

---

## 1. Technical Context

### 1.1 Implementation Language

**Selected: Go (Golang) 1.21+**

**Rationale:**
- **Binary size**: Go produces small, statically-linked binaries (~5-10MB base, achievable 20MB target with embedded Postgres)
- **PostgreSQL integration**: Strong C interop via cgo for embedding libpq/PostgreSQL
- **HTTP server**: Standard library `net/http` (zero dependencies, production-ready)
- **Cross-platform**: Native cross-compilation for Linux/macOS/Windows
- **Performance**: Fast startup, low memory overhead, excellent concurrency
- **Tooling**: Built-in testing framework, race detector, benchmarking

**Alternatives considered:**
- **Rust**: Superior binary size/performance but steeper learning curve, longer compile times
- **C**: Optimal size but significantly more development effort, memory safety risks
- **Zig**: Promising but ecosystem too immature for PostgreSQL embedding

### 1.2 Core Dependencies

#### Embedded PostgreSQL
- **Version**: PostgreSQL 16.x (latest stable)
- **Source**: Official PostgreSQL source tarball
- **Build**: Custom minimal build configuration
- **Integration**: Embedded as shared library or statically linked binary

#### Go Standard Library
- `net/http` - HTTP server (no external HTTP framework)
- `encoding/json` - JSON parsing/serialization
- `database/sql` - Database abstraction
- `context` - Timeout management
- `os/exec` - Process management for embedded Postgres

#### External Go Packages (Minimal Set)
- `github.com/lib/pq` - Pure Go PostgreSQL driver (BSD license, 150KB)
- **Optional**: `github.com/mattn/go-sqlite3` - Fallback if Postgres embedding proves difficult (can be removed post-Phase 1)

#### Build Dependencies
- Docker - Reproducible builds
- PostgreSQL build tools - autoconf, make, gcc
- Go toolchain - go 1.21+

### 1.3 Target Platforms

**Priority 1 (v1.0):**
- Linux x86_64 (GOOS=linux GOARCH=amd64)

**Priority 2 (v1.1):**
- Linux ARM64 (GOOS=linux GOARCH=arm64)
- macOS Apple Silicon (GOOS=darwin GOARCH=arm64)
- macOS Intel (GOOS=darwin GOARCH=amd64)

**Priority 3 (v1.2):**
- Windows x86_64 (GOOS=windows GOARCH=amd64)

---

## 2. Implementation Approach

### 2.1 Architecture Overview

```
┌─────────────────────────────────────────────────┐
│              vibe CLI Binary                    │
│                                                 │
│  ┌───────────────────────────────────────────┐ │
│  │   Command Layer (main.go)                 │ │
│  │   - serve, version, help                  │ │
│  └─────────────────┬─────────────────────────┘ │
│                    │                            │
│  ┌─────────────────▼─────────────────────────┐ │
│  │   HTTP Server Layer (server/)             │ │
│  │   - Request parsing/validation            │ │
│  │   - Response formatting                   │ │
│  │   - Error mapping                         │ │
│  └─────────────────┬─────────────────────────┘ │
│                    │                            │
│  ┌─────────────────▼─────────────────────────┐ │
│  │   Query Processor (query/)                │ │
│  │   - SQL validation                        │ │
│  │   - Safety checks (WHERE enforcement)     │ │
│  │   - Timeout management                    │ │
│  │   - Result limiting                       │ │
│  └─────────────────┬─────────────────────────┘ │
│                    │                            │
│  ┌─────────────────▼─────────────────────────┐ │
│  │   PostgreSQL Adapter (postgres/)          │ │
│  │   - Connection pool                       │ │
│  │   - Query execution                       │ │
│  │   - SQLSTATE error mapping                │ │
│  └─────────────────┬─────────────────────────┘ │
│                    │                            │
│  ┌─────────────────▼─────────────────────────┐ │
│  │   Embedded PostgreSQL Engine              │ │
│  │   - libpq (minimal build)                 │ │
│  │   - Data directory: ./vibe-data/          │ │
│  └───────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘
```

### 2.2 Component Design

#### 2.2.1 Command Layer (`cmd/vibe/`)

**Responsibility**: CLI command routing and orchestration

**Commands:**
- `vibe serve` → Start HTTP server and embedded Postgres
- `vibe version` → Print version info and exit
- `vibe help` → Display usage information

**Implementation:**
```go
// cmd/vibe/main.go
package main

import (
    "fmt"
    "os"
    "github.com/vibesql/vibe/internal/server"
    "github.com/vibesql/vibe/internal/postgres"
)

func main() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    switch os.Args[1] {
    case "serve":
        serve()
    case "version":
        printVersion()
    case "help":
        printUsage()
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}
```

#### 2.2.2 HTTP Server Layer (`internal/server/`)

**Responsibility**: HTTP request/response handling, routing

**Key components:**
- `server.go` - HTTP server lifecycle (start/stop/graceful shutdown)
- `handler.go` - Query endpoint handler
- `response.go` - JSON response formatting
- `errors.go` - HTTP error mapping

**API Contract:**

```
POST /v1/query
Content-Type: application/json

Request:
{
  "sql": "SELECT * FROM users WHERE id = 1"
}

Response (success):
{
  "success": true,
  "rows": [{"id": 1, "data": {"name": "Alice"}}],
  "rowCount": 1,
  "executionTime": 5.2
}

Response (error):
{
  "success": false,
  "error": {
    "code": "INVALID_SQL",
    "message": "Syntax error near 'SELCT'",
    "detail": "PostgreSQL error: syntax error at or near \"SELCT\""
  }
}
```

**Implementation sketch:**
```go
// internal/server/handler.go
type QueryRequest struct {
    SQL string `json:"sql"`
}

type QueryResponse struct {
    Success       bool                     `json:"success"`
    Rows          []map[string]interface{} `json:"rows,omitempty"`
    RowCount      int                      `json:"rowCount,omitempty"`
    ExecutionTime float64                  `json:"executionTime,omitempty"`
    Error         *ErrorDetail             `json:"error,omitempty"`
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    // 2. Validate (length, format)
    // 3. Execute via query processor
    // 4. Format response
    // 5. Write JSON
}
```

#### 2.2.3 Query Processor (`internal/query/`)

**Responsibility**: SQL validation, safety enforcement, execution orchestration

**Key components:**
- `validator.go` - SQL syntax pre-checks, length limits
- `safety.go` - UPDATE/DELETE WHERE clause enforcement
- `executor.go` - Query execution with timeout/limits
- `limiter.go` - Result row limiting

**Safety enforcement:**
```go
// internal/query/safety.go
func enforceWhereClause(sql string) error {
    // Parse SQL (simple regex or basic parser)
    if isUpdateOrDelete(sql) && !hasWhereClause(sql) {
        return ErrMissingWhereClause
    }
    return nil
}
```

**Timeout management:**
```go
// internal/query/executor.go
func (e *Executor) Execute(sql string) (*Result, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    rows, err := e.db.QueryContext(ctx, sql)
    // Handle context.DeadlineExceeded → QUERY_TIMEOUT error
}
```

#### 2.2.4 PostgreSQL Adapter (`internal/postgres/`)

**Responsibility**: Database connection, query execution, error translation

**Key components:**
- `manager.go` - Postgres lifecycle (start/stop embedded process)
- `connection.go` - Connection pool management
- `errors.go` - SQLSTATE → VibeSQL error code mapping

**Error mapping:**
```go
// internal/postgres/errors.go
var sqlStateToVibeCode = map[string]string{
    "42601": "INVALID_SQL",      // syntax_error
    "42703": "INVALID_SQL",      // undefined_column
    "42P01": "INVALID_SQL",      // undefined_table
    "57014": "QUERY_TIMEOUT",    // query_canceled
    "53400": "DATABASE_UNAVAILABLE", // configuration_limit_exceeded
}
```

**Embedded Postgres management:**
```go
// internal/postgres/manager.go
type Manager struct {
    process  *exec.Cmd
    dataDir  string
    port     int
}

func (m *Manager) Start() error {
    // 1. Initialize data directory (if first run)
    // 2. Start postgres process
    // 3. Wait for ready signal
    // 4. Create connection pool
}
```

### 2.3 PostgreSQL Minimal Build Strategy

**Phase 1 Approach: Custom PostgreSQL Build**

#### Build Configuration
```bash
# Dockerfile for Postgres micro build
FROM postgres:16 as builder

# Configure minimal build
./configure \
    --prefix=/usr/local/pgsql \
    --without-readline \
    --without-zlib \
    --without-openssl \
    --without-ldap \
    --without-pam \
    --without-gssapi \
    --without-icu \
    --disable-rpath \
    --disable-debug \
    --disable-profiling \
    --with-system-tzdata=/usr/share/zoneinfo

# Build only required components
make -j$(nproc) \
    src/backend/postgres \
    src/backend/utils/fmgr/fmgrtab.o \
    src/backend/catalog/postgres.bki

# Strip symbols
strip --strip-all postgres
```

#### Binary Size Targets
- Base PostgreSQL: ~15MB (after stripping)
- Go HTTP server: ~5MB
- **Total target**: 20MB (hard limit: 25MB)

#### Excluded Components (Size Savings)
- SSL/TLS libraries (~3MB)
- Replication support (~2MB)
- Authentication modules (LDAP/PAM/GSSAPI) (~1MB)
- Extra command-line tools (psql, pg_dump) (~5MB)
- Extra contrib modules (~2MB)

#### Required Components (Must Keep)
- Core SQL engine
- JSONB data type and operators
- plpgsql procedural language
- Basic aggregate functions
- Index support (B-tree, GIN for JSONB)

### 2.4 Data Storage

**Data directory structure:**
```
./vibe-data/
├── base/          # Database files
├── global/        # Cluster-wide tables
├── pg_wal/        # Write-ahead log
├── pg_stat/       # Statistics
└── PG_VERSION     # Version marker
```

**Initialization:**
- Auto-create `./vibe-data/` on first run
- Run `initdb` equivalent programmatically
- No user-facing configuration required

---

## 3. Source Code Structure

### 3.1 Repository Layout

```
vibesql-local/
├── cmd/
│   └── vibe/
│       └── main.go                  # CLI entry point
│
├── internal/
│   ├── server/                      # HTTP server
│   │   ├── server.go               # Server lifecycle
│   │   ├── handler.go              # /v1/query handler
│   │   ├── response.go             # JSON response formatting
│   │   └── errors.go               # HTTP error codes
│   │
│   ├── query/                       # Query processing
│   │   ├── validator.go            # SQL validation
│   │   ├── safety.go               # WHERE clause enforcement
│   │   ├── executor.go             # Query execution
│   │   └── limiter.go              # Result limiting
│   │
│   ├── postgres/                    # PostgreSQL adapter
│   │   ├── manager.go              # Embedded Postgres lifecycle
│   │   ├── connection.go           # Connection pool
│   │   ├── errors.go               # SQLSTATE mapping
│   │   └── embed/                  # Embedded Postgres binary
│   │       └── postgres            # (Linux x64)
│   │
│   └── version/
│       └── version.go              # Version constants
│
├── tests/
│   ├── unit/                        # Unit tests (~50)
│   │   ├── server_test.go
│   │   ├── validator_test.go
│   │   └── safety_test.go
│   │
│   ├── integration/                 # Integration tests (~100)
│   │   ├── jsonb_test.go           # JSONB operators
│   │   ├── sql_subset_test.go      # Supported SQL
│   │   └── limits_test.go          # Limit enforcement
│   │
│   └── e2e/                         # End-to-end tests (~10)
│       └── workflows_test.go       # Full CRUD workflows
│
├── build/
│   ├── Dockerfile.postgres         # Postgres minimal build
│   ├── Dockerfile.vibe             # Final binary build
│   └── build.sh                    # Build orchestration script
│
├── docs/
│   ├── README.md                   # Quick start guide
│   ├── API.md                      # API reference
│   ├── ERRORS.md                   # Error code reference
│   └── JSONB.md                    # JSONB examples
│
├── scripts/
│   └── install.sh                  # Installer script
│
├── .github/
│   └── workflows/
│       ├── build.yml               # CI build pipeline
│       └── test.yml                # Test pipeline
│
├── go.mod                          # Go dependencies
├── go.sum                          # Dependency checksums
├── Makefile                        # Build automation
└── README.md                       # Project overview
```

### 3.2 Key Files

**Entry Point:**
- `cmd/vibe/main.go` - CLI entry point, command routing

**Server:**
- `internal/server/server.go` - HTTP server setup, graceful shutdown
- `internal/server/handler.go` - Query endpoint logic

**Query Processing:**
- `internal/query/validator.go` - Length checks, basic SQL validation
- `internal/query/safety.go` - WHERE clause enforcement for UPDATE/DELETE

**PostgreSQL Integration:**
- `internal/postgres/manager.go` - Embedded Postgres process management
- `internal/postgres/connection.go` - database/sql connection pool

**Testing:**
- `tests/integration/jsonb_test.go` - All 31 JSONB tests
- `tests/integration/sql_subset_test.go` - Supported SQL statements

---

## 4. Data Model / API Specifications

### 4.1 HTTP API

#### Endpoint: `POST /v1/query`

**Request:**
```json
{
  "sql": "string (max 10KB)"
}
```

**Response (Success):**
```json
{
  "success": true,
  "rows": [
    {"column1": "value1", "column2": "value2"}
  ],
  "rowCount": 1,
  "executionTime": 5.2
}
```

**Response (Error):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_SQL",
    "message": "Syntax error near 'SELCT'",
    "detail": "PostgreSQL error: syntax error at or near \"SELCT\""
  }
}
```

#### Error Codes (HTTP Status → VibeSQL Code)

| HTTP Status | VibeSQL Code | Description |
|-------------|--------------|-------------|
| 400 | INVALID_SQL | Syntax error, unsupported SQL |
| 400 | MISSING_REQUIRED_FIELD | Missing `sql` field in request |
| 400 | UNSAFE_QUERY | UPDATE/DELETE without WHERE |
| 408 | QUERY_TIMEOUT | Query exceeded 5-second timeout |
| 413 | QUERY_TOO_LARGE | SQL string > 10KB |
| 413 | RESULT_TOO_LARGE | Result > 1000 rows |
| 413 | DOCUMENT_TOO_LARGE | JSONB document > 1MB |
| 500 | INTERNAL_ERROR | Server-side error |
| 503 | SERVICE_UNAVAILABLE | HTTP server not ready |
| 503 | DATABASE_UNAVAILABLE | Postgres connection failed |

### 4.2 SQL Support Matrix

#### Supported Statements

| Statement | Support | Notes |
|-----------|---------|-------|
| SELECT | ✅ Full | Single table, WHERE, ORDER BY, LIMIT, OFFSET |
| INSERT | ✅ Full | Single row, RETURNING support |
| UPDATE | ✅ Partial | Requires WHERE (or `WHERE 1=1`), RETURNING support |
| DELETE | ✅ Partial | Requires WHERE (or `WHERE 1=1`), RETURNING support |
| CREATE TABLE | ✅ Partial | Basic schema: SERIAL, JSONB, PRIMARY KEY |
| DROP TABLE | ✅ Full | Including IF EXISTS |

#### Unsupported (v1)

| Feature | Status | Deferred To |
|---------|--------|-------------|
| JOINs | ❌ | v2 |
| Subqueries | ❌ | v2 |
| Transactions | ❌ | v2 |
| GROUP BY / HAVING | ❌ | v2 |
| Window functions | ❌ | v2 |
| CTEs (WITH) | ❌ | v2 |

### 4.3 JSONB Operators (Required)

| Operator | Description | Example |
|----------|-------------|---------|
| `->` | Get JSON object field | `data->'name'` |
| `->>` | Get JSON object field as text | `data->>'name'` |
| `#>` | Get JSON object at path | `data#>'{address,city}'` |
| `#>>` | Get JSON object at path as text | `data#>>'{address,city}'` |
| `@>` | Contains | `data @> '{"name":"Alice"}'` |
| `<@` | Contained by | `'{"name":"Alice"}' <@ data` |
| `?` | Key exists | `data ? 'name'` |
| `?|` | Any key exists | `data ?| array['name','email']` |
| `?&` | All keys exist | `data ?& array['name','email']` |

**Success Criteria**: All 31 PostgreSQL JSONB tests pass (see requirements.md for test list)

---

## 5. Delivery Phases

### Phase 1: PostgreSQL Micro Build (Week 1)
**Objective**: Produce 20MB stripped PostgreSQL binary with JSONB support

**Tasks:**
1. Create `build/Dockerfile.postgres` with minimal build configuration
2. Automate build process (strip symbols, remove unused modules)
3. Extract and test JSONB functionality
4. Run 31 JSONB tests against micro build
5. Document build process

**Deliverables:**
- `postgres_micro` binary (≤20MB)
- Build documentation
- JSONB test results (31/31 passing)

**Success Criteria:**
- [ ] Binary size ≤ 20MB
- [ ] All JSONB operators functional
- [ ] 31/31 JSONB tests passing
- [ ] Build reproducible via Docker

### Phase 2: HTTP API Wrapper (Week 2)
**Objective**: Build Go HTTP server that wraps embedded Postgres

**Tasks:**
1. Implement HTTP server (`internal/server/`)
   - Request parsing/validation
   - Response formatting
   - Error mapping
2. Implement query processor (`internal/query/`)
   - SQL validation
   - Safety checks (WHERE enforcement)
   - Timeout management
   - Result limiting
3. Implement Postgres adapter (`internal/postgres/`)
   - Embedded Postgres lifecycle
   - Connection pool
   - SQLSTATE error mapping
4. Integrate components
5. Manual testing

**Deliverables:**
- HTTP server responding to `/v1/query`
- All 10 error codes implemented
- Query validation and limits enforced

**Success Criteria:**
- [ ] Server starts and accepts connections
- [ ] All error codes return correct HTTP status
- [ ] Query timeout enforced at 5s
- [ ] Result limiting works (1000 rows max)
- [ ] UPDATE/DELETE without WHERE rejected

### Phase 3: Testing (Week 3)
**Objective**: Achieve 200/200 tests passing

**Tasks:**
1. Write unit tests (~50)
   - HTTP parsing/validation
   - Error code mapping
   - Safety checks
2. Write integration tests (~100)
   - SQL subset coverage
   - JSONB operators
   - Limit enforcement
3. Write E2E tests (~10)
   - Full CRUD workflows
   - Error handling
4. Load testing
   - Concurrent queries
   - Timeout behavior
   - Large results
5. Fix failing tests

**Deliverables:**
- 200 tests passing
- Test coverage report
- Performance benchmarks

**Success Criteria:**
- [ ] 200/200 tests passing
- [ ] All JSONB tests pass
- [ ] All error paths covered
- [ ] No memory leaks under load

### Phase 4: Packaging & CLI (Week 4)
**Objective**: Single binary distribution with CLI

**Tasks:**
1. Embed PostgreSQL binary in Go binary
   - Use `go:embed` for static assets
   - Extract on first run
2. Implement CLI commands
   - `vibe serve`
   - `vibe version`
   - `vibe help`
3. Cross-platform builds
   - Linux x64 (primary)
   - Linux ARM64, macOS (stretch)
4. Create installation script (`install.sh`)
5. Write documentation
   - README.md (quick start)
   - API.md (endpoint reference)
   - ERRORS.md (error code reference)
   - JSONB.md (operator examples)
6. GitHub release

**Deliverables:**
- `vibe` binary (Linux x64)
- Installation script
- Complete documentation
- GitHub release

**Success Criteria:**
- [ ] Single binary works standalone
- [ ] Installation script works
- [ ] Documentation complete
- [ ] Binary size ≤ 25MB (target 20MB)

---

## 6. Verification Approach

### 6.1 Build Verification

**Binary size check (CI):**
```bash
# Must fail if binary exceeds 25MB
make build
SIZE=$(stat -c%s vibe)
if [ $SIZE -gt 26214400 ]; then
  echo "ERROR: Binary size $SIZE exceeds 25MB limit"
  exit 1
fi
```

**Build reproducibility:**
```bash
# Builds on different machines must produce identical binaries
docker build -f build/Dockerfile.vibe -t vibesql:test .
docker run vibesql:test sha256sum /app/vibe
# Compare checksums across builds
```

### 6.2 Test Verification

**Unit tests:**
```bash
go test ./internal/... -v -race -cover
```

**Integration tests:**
```bash
go test ./tests/integration/... -v -timeout=30s
```

**E2E tests:**
```bash
# Start server in background
./vibe serve &
sleep 2

# Run E2E tests
go test ./tests/e2e/... -v

# Stop server
pkill vibe
```

**JSONB test suite:**
```bash
go test ./tests/integration/jsonb_test.go -v
# Must report 31/31 passing
```

### 6.3 Performance Verification

**Startup time:**
```bash
time ./vibe serve &
# Must be < 2 seconds to "ready" log
```

**Query timeout:**
```bash
# Query that takes 10s should timeout at 5s
curl -X POST http://localhost:5173/v1/query \
  -d '{"sql": "SELECT pg_sleep(10)"}' \
  -w "Time: %{time_total}s\n"
# Should return QUERY_TIMEOUT in ~5s
```

### 6.4 Lint & Code Quality

**Go linting:**
```bash
go vet ./...
golangci-lint run ./...
```

**Code formatting:**
```bash
gofmt -d .
# Must have no output (all files formatted)
```

**Security scanning:**
```bash
gosec ./...
# Should report no high/medium vulnerabilities
```

### 6.5 Acceptance Testing

**Manual acceptance checklist:**
1. Install via `install.sh` on clean machine
2. Run `vibe serve`
3. Execute all example queries from README
4. Verify error messages are helpful
5. Test graceful shutdown (Ctrl+C)
6. Verify data persistence across restarts

---

## 7. Risk Mitigation

### Technical Risks

**Risk: Binary size exceeds 20MB**
- **Mitigation**: Aggressive stripping, CI size checks, 25MB hard limit
- **Contingency**: Remove non-essential Postgres features, optimize Go build flags

**Risk: Embedding Postgres proves complex**
- **Mitigation**: Phase 1 dedicated to Postgres build, separate testing
- **Contingency**: Ship Postgres as separate binary in v1.0, embed in v1.1

**Risk: Cross-platform builds fail**
- **Mitigation**: Linux x64 only for v1.0, defer others to v1.1
- **Contingency**: Document manual build process for other platforms

**Risk: Performance issues**
- **Mitigation**: Leverage PostgreSQL's proven performance, minimal overhead in Go layer
- **Contingency**: Optimize connection pooling, query parsing

### Process Risks

**Risk: Scope creep**
- **Mitigation**: Strict adherence to requirements.md, no improvisation
- **Contingency**: Defer all non-v1 features to v2

**Risk: Schedule slippage**
- **Mitigation**: Clear phase boundaries, daily progress tracking
- **Contingency**: Ship Linux x64 only, defer other platforms

---

## 8. Open Questions

### 8.1 Resolved

- **Q**: Which language? **A**: Go (binary size, stdlib HTTP, cross-platform)
- **Q**: Embed Postgres or external? **A**: Embed (zero-config requirement)
- **Q**: Support Windows in v1? **A**: Defer to v1.1 if complex

### 8.2 To Be Resolved in Implementation

1. **Postgres embedding strategy**: Static linking vs dynamic library?
2. **Connection pooling**: How many connections for single-user localhost?
3. **Data directory location**: `./vibe-data/` vs `~/.vibe/` vs XDG_DATA_HOME?
4. **Go embed vs external extraction**: Embed Postgres binary in Go binary or extract on first run?

---

## 9. Success Metrics

### Technical Metrics
- Binary size: ≤ 20MB (hard limit 25MB)
- Tests passing: 200/200
- JSONB tests: 31/31
- Cold start: < 2 seconds
- Query timeout: 5s ± 100ms
- Zero external dependencies

### Product Metrics
- Installation time: < 30 seconds (download + install)
- Time to first query: < 5 seconds (install → serve → query)
- Error message clarity: 100% human-readable (no raw Postgres errors)
- Documentation completeness: 100% API coverage

---

## 10. Next Steps

1. **Approval**: Review and approve this specification
2. **Planning**: Break down into detailed implementation tasks
3. **Phase 1**: Start PostgreSQL micro build (Week 1)
4. **Checkpoint**: Review Phase 1 results before proceeding to Phase 2

---

**Document Status**: Ready for Planning phase  
**Next Step**: Create detailed implementation plan in `plan.md`  
**Approval Authority**: Technical Lead / QAPert (for technical decisions)
