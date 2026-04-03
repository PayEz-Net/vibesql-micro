# vibesql-micro-v2 Test Strategy

> Goal: "Test the piss out of it" — comprehensive validation before v0.1.0

## Test Pyramid

```
        /\
       /  \     Integration (end-to-end)
      /____\
     /      \   Component (postgres lifecycle)
    /________\  Unit (core logic)
```

---

## 1. Unit Tests (pkg/vsql/)

### 1.1 Open Lifecycle
| Test | What It Proves |
|------|----------------|
| `TestOpenCreatesMarkerAndDataDir` | Open() creates `.vsql` marker and `.data` dir |
| `TestOpenWithEmptyPath` | Empty path defaults to `./default.vsql` |
| `TestOpenWithNonVsqlExtension` | Path gets `.vsql` extension appended |
| `TestOpenAbsolutePath` | Relative paths resolved to absolute |
| `TestOpenExistingDatabase` | Reopening uses existing data dir, no initdb |
| `TestOpenWhileLocked` | Second Open() returns WarmError when locked |
| `TestOpenProgressCallback` | Progress messages fire in correct order |

### 1.2 Query / Exec
| Test | What It Proves |
|------|----------------|
| `TestQueryAfterOpen` | SELECT 1 returns expected result |
| `TestExecCreateTable` | CREATE TABLE succeeds |
| `TestExecInsert` | INSERT returns RowsAffected > 0 |
| `TestQueryWithParams` | `$1` parameter binding works |
| `TestQueryJSON` | QueryJSON returns valid JSON string |
| `TestQueryAfterClose` | Returns "database is closed" error |
| `TestExecAfterClose` | Returns "database is closed" error |

### 1.3 Close / Lock
| Test | What It Proves |
|------|----------------|
| `TestCloseRemovesLock` | `.lock` file deleted after Close() |
| `TestDoubleClose` | No panic or error on second Close() |
| `TestReopenAfterClose` | Can Open() after Close(), data persists |

### 1.4 WarmError
| Test | What It Proves |
|------|----------------|
| `TestWarmErrorString` | Formats What/Why/Fix correctly |
| `TestWarmErrorIsError` | Implements `error` interface |

---

## 2. Component Tests (internal/postgres/)

### 2.1 Process Lifecycle
| Test | What It Proves |
|------|----------------|
| `TestStartFindsFreePort` | Postgres binds to ephemeral port |
| `TestStartCreatesDataDir` | initdb runs when data dir missing |
| `TestStopKillsProcess` | No postgres.exe/pg_ctl orphan after Stop() |
| `TestStopOnUnstartedProcess` | No panic if Stop() called early |

### 2.2 Binary Extraction
| Test | What It Proves |
|------|----------------|
| `TestEnsureBinaryExtractsOnFirstRun` | Embedded files written to temp dir |
| `TestEnsureBinarySkipsOnSubsequentRun` | Marker file prevents re-extraction |
| `TestEnsureBinaryCopiesLIBPQ` | Windows: `libpq-5.dll` → `LIBPQ.dll` |

---

## 3. Integration Tests (test_integration.go)

These tests run the **full workflow** against a real postgres backend.

### 3.1 Smoke Test
```go
// Must pass on both Windows and Linux
Open → Create Table → Insert → Query → Close → Reopen → Query
```

### 3.2 Data Persistence
- Create table, insert row, close, reopen, verify row exists

### 3.3 Concurrent Access
- Open db, run 5 queries concurrently, verify no races

### 3.4 Error Recovery
- Kill postgres process externally, verify Open() recovers (or returns clear error)

### 3.5 JSONB Workflow
- `CREATE TABLE docs (id SERIAL, data JSONB)`
- `INSERT INTO docs (data) VALUES ('{"name":"test"}')`
- `SELECT data->>'name' FROM docs WHERE data @> '{"name":"test"}'`
- Verify JSONB operators work end-to-end

---

## 4. Platform-Specific Tests

### 4.1 Windows
- Run full integration suite
- Verify no `STATUS_DLL_NOT_FOUND` (0xc0000135)
- Verify `.lock` file behavior with backslash paths

### 4.2 Linux
- Run full integration suite on 93
- Verify binary permissions (`chmod +x` on extracted postgres/initdb)
- Verify shared library loading (`LD_LIBRARY_PATH` or `rpath`)
- Verify port cleanup on SIGTERM

---

## 5. Edge Cases & Error Conditions

| Scenario | Expected Behavior |
|----------|-------------------|
| Open on read-only filesystem | WarmError with clear "Why" and "Fix" |
| Open with path containing spaces | Works correctly |
| Open with unicode path | Works correctly |
| Query with syntax error | Returns postgres error message |
| Query with wrong parameter count | Returns postgres error message |
| Database dir deleted while open | Graceful error on next Query |
| Disk full during initdb | WarmError, cleanup partial data dir |

---

## 6. Performance / Stress Tests

| Test | Target |
|------|--------|
| `BenchmarkOpen` | < 500ms on warm open, < 5s on cold open |
| `BenchmarkQueryLatency` | SELECT 1 < 10ms locally |
| `BenchmarkInsert1000` | 1000 INSERTs < 5 seconds |
| `TestBinarySize` | Final binary < 70MB |

---

## 7. Test Execution Plan

### Phase 1: Unit + Component (Local — Windows)
```bash
go test ./pkg/vsql/ -v
go test ./internal/postgres/ -v
```

### Phase 2: Integration (Local — Windows)
```bash
go run test_integration.go
```

### Phase 3: Linux Validation (Remote — 93 via ZeroClaw)
```bash
# On 93
cd ~/repos/vibesql-micro-v2
go test ./... -v
go run test_integration.go
```

### Phase 4: Ship Criteria
All of the following must pass:
- [ ] Windows unit tests: 100% pass
- [ ] Windows integration tests: 100% pass
- [ ] Linux unit tests: 100% pass
- [ ] Linux integration tests: 100% pass
- [ ] JSONB workflow: verified on both platforms
- [ ] Binary size check: < 70MB
- [ ] No orphaned postgres processes after test suite

---

## 8. Known Gaps to Address

1. **Linux embed files** — Need linux postgres binaries in `embed_linux.go`
2. **Integration test path separators** — `test_integration.go` uses `\` (Windows-only)
3. **No concurrency tests yet** — Add `TestConcurrentQueries`
4. **No JSONB-specific unit tests** — Add to `pkg/vsql/query_test.go`
