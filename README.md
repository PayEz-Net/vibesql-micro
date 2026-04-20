# vibesql-micro

PostgreSQL that feels like SQLite. Sunday morning easy.

```bash
$ go install github.com/vibesql/vibesql-micro/cmd/vsql-micro@latest

$ vsql-micro ./app.vsql "SELECT 1"
[{"?column?": 1}]
```

## What is it?

**vibesql-micro** embeds a full PostgreSQL 16 engine inside a single binary. No external PostgreSQL installation, no configuration files, no service management. Just open a `.vsql` file and start querying.

We looked at how projects like **PGlite** handle single-binary distribution, then adapted their best ideas to our own wrapper:
- **Streamlined extraction** — binaries cache once and reuse, with progress feedback on first run.
- **Simplified lifecycle** — dynamic port allocation, lock-file guarding, and graceful shutdown remove the need for manual config.
- **Truly zero-config** — `vsql.Open("app.vsql")` is all you need.

The result is a **~22 MB Linux binary** and **~67 MB Windows binary** that passes an automated test suite across both platforms, with startup times under a second on warm open.

## Quick Start

### Install

```bash
# npm / npx — thin wrapper that downloads the native binary from GitHub Releases
npx vibesql-micro version
# or install as a project dep
npm install vibesql-micro

# or install directly with Go
go install github.com/vibesql/vibesql-micro/cmd/vsql-micro@latest
```

Or download pre-built binaries from [GitHub Releases](https://github.com/PayEz-Net/vibesql-micro/releases).

### CLI

```bash
# Interactive shell
$ vsql-micro ./app.vsql
app.vsql> CREATE TABLE users (id SERIAL, data JSONB);
app.vsql> INSERT INTO users (data) VALUES ('{"name": "Alice"}');
app.vsql> SELECT * FROM users;
[{"id": 1, "data": {"name": "Alice"}}]
app.vsql> \q

# Single query
$ vsql-micro ./app.vsql "SELECT * FROM users"
[{"id": 1, "data": {"name": "Alice"}}]
```

### Server mode (long-lived, pinned port)

For service-style deployments — e.g. hosting the embedded PostgreSQL as a backend for another process on the same box — run `vsql-micro serve`:

```bash
$ vsql-micro serve --listen 127.0.0.1:5433 --data ./vault.vsql
Creating database... done
vsql-micro serving /abs/path/vault.vsql on 127.0.0.1:5433 (user=postgres, trust auth)
Ctrl+C or SIGTERM to shut down.
```

- Binds a fixed TCP port on 127.0.0.1 (pod-internal — file-system ACLs are the boundary)
- `user=postgres` with trust auth
- Graceful shutdown on SIGINT / SIGTERM — signals the embedded PostgreSQL cleanly before exiting

This is the mode used by [vsql-vault](https://github.com/PayEz-Net/vibesql-vault) as its storage backend.

### Go API

```go
package main

import (
    "fmt"
    "log"
    "github.com/vibesql/vibesql-micro/pkg/vsql"
)

func main() {
    // Open database (creates if doesn't exist)
    db, err := vsql.Open("./app.vsql")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Execute SQL
    db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL, data JSONB)")

    // Query with JSONB
    rows, err := db.Query("SELECT * FROM users WHERE data @> $1", `{"active": true}`)
    if err != nil {
        log.Fatal(err)
    }

    for _, row := range rows {
        fmt.Printf("User %d: %v\n", row["id"], row["data"])
    }
}
```

## Features

- **Zero configuration** — Works out of the box
- **Single file database** — Your database is `app.vsql` (marker file + hidden `.data` directory)
- **JSON-first** — JSON output by default, perfect for JSONB workflows
- **Full PostgreSQL 16** — Real PostgreSQL, not a reimplementation
- **Embedded** — Single binary, no external PostgreSQL installation
- **Cross-platform** — Tested on Windows and Linux
- **Concurrent safe** — Lock-file protection prevents multiple processes from opening the same database

## Testing

The project includes a comprehensive test suite covering unit tests, integration tests, and CLI validation:

```bash
# Unit tests
go test ./pkg/vsql/ -v

# Integration tests (Go)
go run test_integration.go
go run test_comprehensive.go

# CLI tests (Linux)
bash test_linux_cli.sh
```

### Test Coverage

- **23 unit tests** — Open lifecycle, query execution, JSONB, unicode, error handling
- **26 integration tests** — Full workflow, concurrency, persistence, performance
- **18 CLI tests** — Version, CRUD, complex queries, lock detection, large result sets
- **Verified platforms:** Linux (x64) and Windows (x64)

See [TEST_STRATEGY.md](TEST_STRATEGY.md) for the full test pyramid.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
