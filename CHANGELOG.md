# Changelog

All notable changes to VibeSQL Micro are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.6] - 2026-02-24

Install script version check — npm upgrades now correctly replace outdated binaries.

### Fixed
- `install.js` now checks the installed binary version before skipping download — previously any existing binary would skip the download, even if it was an older version
- Upgrading via `npm install -g vibesql-micro@<new>` now correctly replaces the old binary

## [1.0.5] - 2026-02-24

Linux initdb fix — embedded `postgres.bki` had CRLF line endings causing startup failure on Linux.

### Fixed
- Linux: `initdb` failed with `input file "...postgres.bki" does not belong to PostgreSQL 16.1` on every startup
- Root cause: `postgres.bki` had Windows CRLF line endings (`\r\n`), causing the version header to read as `# PostgreSQL 16\r` instead of `# PostgreSQL 16` — `strcmp` failed on Linux
- Windows was unaffected (C runtime strips `\r` in text mode)
- Converted `postgres.bki` to Unix LF line endings in `share.tar.gz`

## [1.0.4] - 2026-02-23

Graceful shutdown — zero orphan PostgreSQL processes in any scenario.

### Added
- `vibesql-micro stop` command — cleanly stops a running instance via HTTP shutdown endpoint
- `POST /v1/shutdown` endpoint — localhost-only, triggers graceful shutdown sequence
- Windows Job Object with `KILL_ON_JOB_CLOSE` — guarantees all PostgreSQL children are killed even if Go process crashes
- Unix process groups — kills entire process tree on shutdown
- PID file written on start, cleaned up on stop

### Fixed
- `pg_ctl` race condition — graceful shutdown now completes before process termination
- Process tree kill fallback if `pg_ctl` fails

## [1.0.3] - 2026-02-23

Parameterized query support — pass `params` alongside `sql` for safe, typed queries.

### Added
- Optional `params` array in request body: `{"sql": "SELECT $1::text", "params": ["hello"]}`
- PostgreSQL native `$1`, `$2`, ... placeholders with automatic type inference
- Full CRUD support: SELECT, INSERT, UPDATE, DELETE all work with params
- Backward compatible — existing raw SQL requests without params work unchanged

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

[1.0.6]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.6
[1.0.5]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.5
[1.0.4]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.4
[1.0.3]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.3
[1.0.2]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.2
[1.0.1]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.1
[1.0.0]: https://github.com/PayEz-Net/vibesql-micro/releases/tag/v1.0.0
