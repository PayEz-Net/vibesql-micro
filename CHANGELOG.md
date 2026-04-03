# Changelog

All notable changes to vibesql-micro will be documented in this file.

## [0.2.0] - 2026-04-03

### Added
- Comprehensive cross-platform test suite — 61 total tests across unit, integration, and CLI layers.
- `test_comprehensive.go` — 20-test Go integration suite covering JSONB, unicode, concurrency, persistence, JOINs/aggregates, warm-open performance, and error conditions.
- `test_linux_cli.sh` — 18-test bash CLI suite for Linux validating version, help, CRUD, complex queries, lock detection, large result sets, unicode, and path-with-spaces handling.
- `test_full.sh` — Portable bash integration script for manual CI runs.
- `TEST_STRATEGY.md` — Full test pyramid documentation (unit / integration / platform / edge-case / performance).
- Expanded unit tests in `pkg/vsql/open_test.go` and new `pkg/vsql/query_test.go` covering:
  - Open lifecycle (marker creation, empty path, extension appending, absolute path, lock detection)
  - Query/Exec (CREATE TABLE, INSERT, parameterized queries, QueryJSON, JOINs/aggregates)
  - Data types (JSONB, unicode, large result sets)
  - Error handling (syntax errors, wrong parameter counts)
  - Close/reopen persistence

### Verified
- Linux build fully validated on Ubuntu (x64).
  - 23/23 unit tests passing
  - 26/26 integration tests passing
  - 18/18 CLI tests passing
  - Binary size: ~25 MB
- Windows build remains functional (DLL loading fixed, CLI smoke-tested).
- Cross-platform parity confirmed for core PostgreSQL embedding workflow.

### Fixed
- `vsql-cli` (`--profile` flag) now correctly routes requests to non-default profiles (e.g., `--profile rosa`).

## [0.1.0] - 2026-04-02

### Added
- Initial working Windows build with embedded PostgreSQL binaries.
- `Open()`, `OpenWithProgress()`, `Query()`, `Exec()`, `QueryJSON()`, `Close()` API in `pkg/vsql`.
- Dynamic port allocation for postgres subprocess.
- Lock-file based concurrent access protection with `WarmError` user-friendly messages.
- `vsql-micro` CLI tool for interactive and single-shot SQL execution.
- Binary extraction and caching to user cache dir.
