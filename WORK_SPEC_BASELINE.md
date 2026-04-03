# Work Spec: Warm Brain MVP - vibesql-micro

## Project
**Name:** vibesql-micro Warm Brain MVP  
**Goal:** Sunday-morning-easy PostgreSQL for Go  
**Duration:** 1 week  
**Status:** Ready for Review

---

## Overview

Transform vibesql-micro from "works but requires setup" to "just works on first try." The target is a developer in pajamas at 9am Sunday who wants to add a database to their side project without friction.

---

## Deliverables

### 1. Zero-Friction Installation
**Acceptance Criteria:**
- [ ] `go install github.com/vibesql/vsql-micro/cmd/vsql-micro@latest` completes in < 30 seconds
- [ ] No additional dependencies required
- [ ] Works on fresh Windows, macOS, and Linux machines
- [ ] First run shows gentle progress indicator (not silent, not verbose)

**Technical Approach:**
- Embed all platform postgres binaries in Go binary
- Extract to OS temp directory on first run
- Cache and reuse extracted binary
- Auto-clean old extractions

---

### 2. Gentle CLI
**Acceptance Criteria:**
- [ ] `vsql-micro ./app.vsql` opens interactive shell
- [ ] `vsql-micro ./app.vsql "SELECT 1"` runs single query
- [ ] `vsql-micro version` shows version info
- [ ] `vsql-micro --help` is actually helpful (examples, not flags)
- [ ] Exit with `\q` or Ctrl+D (both work)
- [ ] Errors are human-readable with suggested fixes

**Technical Approach:**
- Use existing Go REPL libraries or build simple one
- Table output with `github.com/olekukonko/tablewriter` or similar
- Error wrapping with context and fixes
- Signal handling for clean shutdown

---

### 3. Single-File Database UX
**Acceptance Criteria:**
- [ ] `vsql.Open("./app.vsql")` creates marker file + hidden data dir
- [ ] User sees only `app.vsql` (4KB marker), not implementation details
- [ ] Data stored in `.app.vsql-data/` (hidden directory)
- [ ] Close() cleans up temp files, leaves data intact
- [ ] Re-opening same path reconnects to existing data

**Technical Approach:**
- Marker file contains metadata (version, pid, etc.)
- Hidden directory name: `.` + marker name + `-data`
- Lock file in hidden directory for concurrent access detection
- Auto-create directories as needed

---

### 4. Go API Simplicity
**Acceptance Criteria:**
- [ ] `db, err := vsql.Open("./app.vsql")` is the only required call
- [ ] `db.Query()` returns rows as `[]map[string]any`
- [ ] `db.Exec()` returns simple result struct
- [ ] `db.Close()` shuts down cleanly
- [ ] All methods have sensible timeouts (5s default)
- [ ] Errors include what went wrong and how to fix it

**Technical Approach:**
- Keep existing subprocess architecture (proven)
- Simplify connection handling
- Better error wrapping
- Remove configuration complexity from public API

---

### 5. Warm Brain Polish
**Acceptance Criteria:**
- [ ] First run shows: "Setting up vibesql... done" (not silence)
- [ ] Creating database shows: "Creating database... done"
- [ ] Opening existing database is silent (just works)
- [ ] Table output is aligned and readable
- [ ] Error: "That database is busy (another program is using it)"
- [ ] Not error: "FATAL: could not create lock file..."

**Technical Approach:**
- Progress callbacks for long operations
- Error types with What/Why/Fix fields
- Terminal detection for colors
- Sensible defaults everywhere

---

## Non-Goals (Out of Scope)

These are intentionally NOT in this spec:

- ❌ JSONB sugar layer (InsertDoc, FindDoc helpers) - v2
- ❌ Windows service support - later
- ❌ Database migrations - later  
- ❌ Backup/restore commands - later
- ❌ Performance optimizations - later
- ❌ WASM runtime - abandoned

---

## File Structure

```
vibesql-micro/
├── cmd/vsql-micro/
│   ├── main.go              # CLI entry point
│   ├── shell.go             # Interactive REPL
│   ├── query.go             # Single query execution
│   └── version.go           # Version info
├── pkg/vsql/
│   ├── open.go              # Open() function
│   ├── query.go             # Query() and Exec()
│   ├── close.go             # Close() cleanup
│   ├── errors.go            # WarmError types
│   ├── progress.go          # Progress callbacks
│   └── embed/               # Embedded binaries
│       ├── postgres_linux_amd64
│       ├── postgres_darwin_amd64
│       └── postgres_windows_amd64.exe
└── internal/
    ├── binary/              # Binary extraction/management
    ├── postgres/            # Subprocess lifecycle
    └── protocol/            # JSON protocol
```

---

## Testing Plan

### Unit Tests
- [ ] Open creates marker file and data directory
- [ ] Open returns error for locked database
- [ ] Query returns expected rows
- [ ] Close cleans up but preserves data
- [ ] Error messages include fix suggestions

### Integration Tests
- [ ] Full workflow: Open → Query → Close → Reopen
- [ ] Concurrent access detection
- [ ] Binary extraction and caching
- [ ] Cross-platform (Windows, macOS, Linux)

### Warm Brain Test
Fresh machine, fresh user:
```bash
# 1. Install
go install github.com/vibesql/vsql-micro/cmd/vsql-micro@latest
# Expected: Success in < 30s

# 2. First use
vsql-micro ./test.vsql "SELECT 1"
# Expected: "Setting up... done", then "1"

# 3. Shell
vsql-micro ./test.vsql
# Expected: Interactive prompt, \q exits cleanly

# 4. From Go
# Copy 5-line example from README
# Expected: Works immediately
```

---

## Timeline

| Day | Focus | Deliverable |
|-----|-------|-------------|
| 1 | Installation | `go install` works with progress |
| 2 | CLI basics | Shell and single query work |
| 3 | Single-file UX | Marker file + hidden data dir |
| 4 | Go API | Open/Query/Exec/Close simplified |
| 5 | Error handling | Warm errors with fixes |
| 6 | Polish | Colors, alignment, testing |
| 7 | Release | GitHub release with binaries |

---

## Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Installation time | < 30s | `time go install` |
| First query time | < 10s | From install to SELECT 1 |
| Binary size | < 25MB | `du -h vsql-micro` |
| Test pass rate | 100% | `go test ./...` |
| Warm brain score | 4/5 | Beta tester survey |

---

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Binary extraction slow | Medium | Annoyance | Show progress, cache aggressively |
| Cross-platform issues | Medium | Blocker | Test on all platforms early |
| Subprocess flaky | Low | Blocker | Keep existing proven code |
| Scope creep | High | Delay | Strict non-goals list |

---

## Dependencies

### External
- `github.com/lib/pq` - PostgreSQL driver (already using)
- `github.com/olekukonko/tablewriter` - Table output (new)

### Internal
- Existing vibesql-micro code
- Existing postgres-micro binaries

---

## Review Checklist

Before starting work, confirm:

- [ ] Scope is clear and bounded
- [ ] Non-goals are explicitly listed
- [ ] Timeline is realistic
- [ ] Success metrics are measurable
- [ ] Risks are identified

---

## Review Decisions

| Question | Decision |
|----------|----------|
| **1. Platform priority** | Windows first (we're on Windows, easiest commands), Linux/macOS later |
| **2. Go version** | 1.22+ required |
| **3. Shell features** | Just SQL (no meta-commands like \d, \dt) |
| **4. Output format** | JSON by default (non-transforming, pass-through) |
| **5. Timeline** | 1 week (doesn't matter, but 1 week is fine) |

### JSON-First Rationale

Since our use case is JSONB-heavy, default to JSON output:
```bash
$ vsql-micro ./app.vsql "SELECT * FROM users"
[{"id": 1, "name": "Alice", "data": {"active": true}}]
```

Non-transforming = PostgreSQL returns JSON, we don't parse/reformat it. Just pass through. This is:
- Faster (no table formatting)
- More useful for JSONB workflows
- Pipe-friendly (`| jq`)
- Programmable from other tools

Optional table view later: `--table` flag.

---

## Sign-Off

| Role | Name | Approved | Notes |
|------|------|----------|-------|
| Product | | | |
| Technical | | | |
| QA | | | |

---

*This is the baseline. Adjust as needed, then execute.*
