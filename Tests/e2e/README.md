# E2E Tests for VibeSQL

## Overview

These end-to-end tests verify complete user workflows via the HTTP API. They test real-world scenarios from starting the server to making HTTP requests and verifying responses.

## Prerequisites

### 1. PostgreSQL

You need a running PostgreSQL instance. The easiest way is via Docker:

```bash
docker run -d \
  --name vibesql-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:latest
```

Create the test database:

```bash
docker exec -it vibesql-postgres psql -U postgres -c "CREATE DATABASE vibesql_test;"
```

### 2. VibeSQL Server

Start the VibeSQL server:

```bash
# Build the binary
go build -o vibe cmd/vibe/main.go

# Start the server
./vibe serve
```

The server should start on `http://localhost:5173`.

## Running Tests

### Run all E2E tests

```bash
go test ./Tests/e2e/... -v
```

### Run specific test

```bash
go test ./Tests/e2e/... -v -run TestE2E_FullCRUDWorkflow
```

### Skip timeout tests (faster)

```bash
go test ./Tests/e2e/... -v -short
```

### Skip E2E tests entirely

If you don't have the prerequisites set up:

```bash
export VIBESQL_E2E_SKIP=1
go test ./Tests/e2e/... -v
```

## Test Scenarios

### 1. Server Ready Check
- **Test**: `TestE2E_ServerReady`
- **Purpose**: Verifies the server is running and accessible
- **Duration**: ~2 seconds

### 2. Full CRUD Workflow
- **Test**: `TestE2E_FullCRUDWorkflow`
- **Purpose**: Tests CREATE → INSERT → SELECT → UPDATE → DELETE lifecycle
- **Coverage**: Complete database operations workflow
- **Duration**: ~2 seconds

### 3. Concurrent Queries
- **Test**: `TestE2E_ConcurrentQueries`
- **Purpose**: Verifies multiple simultaneous requests work correctly
- **Coverage**: Concurrent INSERT operations, race condition testing
- **Duration**: ~1 second

### 4. Error Recovery
- **Test**: `TestE2E_ErrorRecovery`
- **Purpose**: Ensures invalid queries don't crash the server
- **Coverage**: Invalid SQL, unsafe queries, server stability
- **Duration**: ~1 second

### 5. Timeout Handling
- **Test**: `TestE2E_TimeoutHandling`
- **Purpose**: Verifies query timeout enforcement (5 seconds)
- **Coverage**: Long-running queries, timeout errors, server recovery
- **Duration**: ~7 seconds
- **Note**: Skipped in `-short` mode

### 6. JSONB Workflow
- **Test**: `TestE2E_JSONBWorkflow`
- **Purpose**: Tests JSONB data type operations
- **Coverage**: JSONB insert, query with `->` and `->>` operators
- **Duration**: ~1 second

### 7. Limit Enforcement
- **Test**: `TestE2E_LimitEnforcement`
- **Purpose**: Verifies 1000-row result limit
- **Coverage**: Row limit enforcement, RESULT_TOO_LARGE error
- **Duration**: ~3 seconds

### 8. Data Persistence (Server Restart)
- **Test**: `TestE2E_DataPersistence`
- **Purpose**: Verifies data persists across server restarts
- **Coverage**: Two-phase test (create data, restart, verify data)
- **Duration**: ~2 seconds (per phase)
- **Note**: Requires manual server restart between test runs (automated in Phase 4)

### 9. Graceful Shutdown
- **Test**: `TestE2E_GracefulShutdown`
- **Purpose**: Verifies server handles shutdown gracefully
- **Coverage**: Concurrent queries during potential shutdown window
- **Duration**: ~3 seconds
- **Note**: Full SIGTERM testing will be added in Phase 4

### 10. Multiple Tables Workflow
- **Test**: `TestE2E_MultipleTablesWorkflow`
- **Purpose**: Tests working with multiple tables simultaneously
- **Coverage**: Multiple CREATE/INSERT/SELECT operations
- **Duration**: ~2 seconds

## Environment Variables

- **`VIBESQL_TEST_DB`**: Custom PostgreSQL connection string (optional)
  - Default: `host=localhost port=5432 user=postgres password=postgres dbname=vibesql_test sslmode=disable`
- **`VIBESQL_E2E_SKIP`**: Set to `1` to skip E2E tests if prerequisites aren't met

## Test Architecture

### API Client
The tests use a lightweight HTTP client (`APIClient`) to make requests to the VibeSQL server. This mimics real-world usage and tests the full HTTP stack.

### Test Isolation
Each test:
- Uses unique table names (timestamped) to avoid conflicts
- Cleans up tables on exit using `defer`
- Can run in parallel safely

### Error Handling
Tests verify both success and error cases:
- Success responses include `success: true`, `rows`, `rowCount`, `executionTime`
- Error responses include `success: false`, `error.code`, `error.message`, `error.detail`

## Troubleshooting

### "Server not ready" error

```
Server not ready: server did not become ready within 2s
```

**Solution**: Ensure the VibeSQL server is running on `localhost:5173`:
```bash
./vibe serve
```

### "Failed to ping test database" error

```
Failed to ping test database: connection refused
```

**Solution**: Ensure PostgreSQL is running:
```bash
docker ps | grep postgres
```

If not running:
```bash
docker run -d -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
```

### Tests timeout or hang

**Solution**: Check server logs for errors. Common issues:
- Database connection pool exhausted
- PostgreSQL not accessible
- Port 5173 already in use

## Phase 4 Updates

In **Phase 4** (Binary Embedding), these tests will be updated to:
1. Automatically start the embedded VibeSQL binary
2. Manage server lifecycle programmatically
3. Test data persistence across server restarts

For now, they require manual server setup to test the HTTP API workflows.

## Contributing

When adding new E2E tests:
1. Use unique table names (timestamp suffix)
2. Add cleanup with `defer`
3. Check server readiness with `WaitForServer()`
4. Skip tests gracefully if server isn't available
5. Verify both success and error cases
6. Update this README with test description
