# vibesql-micro

SQLite + JSON + HTTP API in one command.

---

## What is this?

`vibesql-micro` is a lightweight database server for local development.

- **SQLite-based** — Fast, embedded, zero config
- **JSON support** — Flexible schema, JSONB-style queries
- **HTTP API** — Query via curl, Postman, or any HTTP client
- **Single command** — `npx vibesql-micro` and you're running

Perfect for prototyping, testing, and local development.

---

## Quick Start

```bash
npx vibesql-micro
# → Running at http://localhost:5173
```

Query your database:

```bash
curl -X POST http://localhost:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT * FROM users LIMIT 10"}'
```

---

## Installation

### npm (recommended)

```bash
npx vibesql-micro
```

### Global install

```bash
npm install -g vibesql-micro
vibesql-micro
```

### Configuration

Set environment variables:

```bash
VIBESQL_PORT=5173          # HTTP port (default: 5173)
VIBESQL_DATA=./vibe-data   # Data directory (default: ./vibe-data)
```

---

## API Reference

### POST /v1/query

Execute SQL queries.

**Request:**

```json
{
  "sql": "SELECT * FROM users WHERE active = 1"
}
```

**Response (success):**

```json
{
  "success": true,
  "data": [
    { "id": 1, "name": "Alice", "email": "alice@example.com" }
  ],
  "meta": {
    "row_count": 1,
    "execution_time_ms": 5
  }
}
```

**Response (error):**

```json
{
  "success": false,
  "error": {
    "code": "SYNTAX_ERROR",
    "message": "near \"FROM\": syntax error"
  }
}
```

---

### GET /v1/health

Check server status.

**Response:**

```json
{
  "status": "healthy",
  "version": "0.1.0",
  "database": "sqlite",
  "uptime_seconds": 3600
}
```

---

## JSON Queries

VibeSQL supports JSON1 extension for flexible schema:

```sql
-- Store JSON
INSERT INTO users (data) VALUES ('{"name": "Alice", "tags": ["developer", "golang"]}');

-- Query JSON fields
SELECT data->>'name' as name FROM users;

-- Filter by JSON
SELECT * FROM users WHERE data->>'active' = 'true';

-- Array operations
SELECT * FROM users WHERE json_array_length(data->'tags') > 1;
```

---

## Use Cases

### Prototyping

```bash
# Start database
npx vibesql-micro

# Create table and insert data
curl -X POST http://localhost:5173/v1/query \
  -d '{"sql": "CREATE TABLE todos (id INTEGER PRIMARY KEY, title TEXT, done INTEGER)"}'

curl -X POST http://localhost:5173/v1/query \
  -d '{"sql": "INSERT INTO todos (title, done) VALUES ('\''Buy milk'\'', 0)"}'
```

### Testing

```javascript
// test-helper.js
import { spawn } from 'child_process';

export async function startTestDB() {
  const proc = spawn('npx', ['vibesql-micro'], {
    env: { ...process.env, VIBESQL_PORT: 5174 }
  });

  // Wait for startup
  await new Promise(resolve => setTimeout(resolve, 1000));

  return proc;
}

export async function stopTestDB(proc) {
  proc.kill();
}
```

### Local Development

```bash
# Terminal 1: Run database
npx vibesql-micro

# Terminal 2: Run your app
npm run dev

# Your app connects to http://localhost:5173
```

---

## Admin UI

Want a visual interface? Use [vibesql-admin](https://github.com/vibesql/vibesql-admin):

```bash
# Terminal 1: Database
npx vibesql-micro

# Terminal 2: Admin UI
npx vibesql-admin
# → Opens browser at http://localhost:5174
```

---

## Comparison

| Feature | VibeSQL Micro | SQLite CLI | PostgreSQL |
|---------|---------------|------------|------------|
| Installation | `npx` command | Download binary | Homebrew/apt/installer |
| HTTP API | ✅ Built-in | ❌ No | ❌ No (need PostgREST) |
| JSON support | ✅ JSON1 | ✅ JSON1 | ✅ JSONB (more powerful) |
| Setup time | < 10 seconds | ~1 minute | ~5 minutes |
| Use case | Local dev, prototyping | Embedded apps | Production |

---

## Limitations

- **Not for production** — Use PostgreSQL, MySQL, or managed services for production
- **SQLite-based** — Subject to SQLite limitations (no parallel writes)
- **JSON1 extension** — Less powerful than PostgreSQL's JSONB
- **Single file database** — Not suitable for high-concurrency workloads

For production features (replication, backups, auth, scaling), see [VibeSQL Cloud](https://vibesql.online/cloud).

---

## Development

Clone the repo:

```bash
git clone https://github.com/vibesql/vibesql-micro.git
cd vibesql-micro
```

Build:

```bash
go build -o vibesql-micro ./cmd/server
```

Run:

```bash
./vibesql-micro
```

Test:

```bash
go test ./...
```

---

## Tech Stack

- **Language:** Go
- **Database:** SQLite + JSON1 extension
- **HTTP:** Standard library (`net/http`)
- **Packaging:** npm (via npx)

---

## Contributing

Contributions welcome. Open an issue or pull request.

---

## License

MIT License. See [LICENSE](LICENSE).

---

## Links

- **Website:** [vibesql.online](https://vibesql.online)
- **Admin UI:** [github.com/vibesql/vibesql-admin](https://github.com/vibesql/vibesql-admin)
- **Docs:** [vibesql.online/docs](https://vibesql.online/docs)
- **Discord:** [discord.gg/vibesql](https://discord.gg/vibesql)

---

Built for developers. Zero config. Just works.
