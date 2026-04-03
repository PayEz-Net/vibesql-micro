# vibesql-micro

PostgreSQL that feels like SQLite. Sunday morning easy.

```bash
$ go install github.com/vibesql/vibesql-micro/cmd/vsql-micro@latest

$ vsql-micro ./app.vsql "SELECT 1"
[{"?column?": 1}]
```

## Quick Start

### Install

```bash
go install github.com/vibesql/vibesql-micro/cmd/vsql-micro@latest
```

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

- **Zero configuration** - Works out of the box
- **Single file** - Your database is `app.vsql` (marker file + hidden data)
- **JSON-first** - JSON output by default, perfect for JSONB workflows
- **Full PostgreSQL** - Real PostgreSQL 16, not a reimplementation
- **Embedded** - Single binary, no external PostgreSQL installation

## Status

**Work in progress** - Day 1 of development. Basic structure in place.

## License

Apache 2.0
