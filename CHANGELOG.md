# Changelog

All notable changes to VibeSQL Micro are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.2] - 2026-02-23

Better startup diagnostics — when PostgreSQL fails to start, you now see exactly why.

### Added
- PostgreSQL stderr output included in startup failure messages (last 20 lines)
- Step-by-step progress logging: `Extracting` → `Initializing` → `Starting on port` → `Waiting for connections`
- Real-time PostgreSQL LOG and WARNING lines shown during startup

### Fixed
- Generic "timeout after 30s" message now shows the actual PostgreSQL error (port conflicts, permission issues, missing files)

## [1.0.1] - 2026-02-22

Windows fix for PostgreSQL `/share` directory resolution bug.

### Fixed
- Windows: PostgreSQL resolves `/share` as `<drive>:\share` — now creates share directory on both the temp and CWD drives
- Clear error message when drive-root write fails (suggests running as Administrator)
- Automatic cleanup of `\share` and `\lib` directories on shutdown

## [1.0.0] - 2026-02-22

First public release.

### Added
- Embedded PostgreSQL 16.1 — full ACID, JSONB, arrays, all standard types
- HTTP API at `POST /v1/query` — send SQL, get JSON back
- Single binary for Windows (x64), Linux (x64), and macOS (Intel)
- npm distribution — `npx vibesql-micro` downloads and runs automatically
- Safety: UPDATE/DELETE require WHERE clause, query size limits, row count limits, 5s timeout
- Auto-creates data directory, initializes PostgreSQL, starts HTTP server
- Clean shutdown with Ctrl+C — stops PostgreSQL, removes temp files

---

[1.0.2]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.2
[1.0.1]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.1
[1.0.0]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.0
