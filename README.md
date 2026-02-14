# vibesql-micro

PostgreSQL + JSONB + HTTP API in one command.

---

## What is this?

`vibesql-micro` is a lightweight database server for local development with **embedded PostgreSQL 16.1**.

- **PostgreSQL-based** — Full PostgreSQL 16.1 embedded in a single binary
- **Native JSONB** — Real PostgreSQL JSONB support, not an extension
- **HTTP API** — Query via curl, Postman, or any HTTP client
- **Single command** — `npx vibesql-micro` and you're running
- **Zero config** — No installation, no setup, no Docker

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

### Windows (Available Now)

Download the latest Windows binary from [Releases](https://github.com/PayEz-Net/vibesql-micro/releases):

```bash
# Download vibesql-micro-windows-x64.exe
# Run it
.\vibesql-micro-windows-x64.exe
# → Running at http://localhost:5173
```

The server will:
1. Auto-create required directories (`<drive>:\share`, `<drive>:\lib`)
2. Start PostgreSQL 16.1 on port 5432
3. Start HTTP API on port 5173
4. Clean up temporary directories on shutdown

### npm (Coming Soon)

```bash
npx vibesql-micro
```

Windows, macOS, and Linux support will be available via npm in v1.0.0 final release.

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
  "version": "1.0.0",
  "database": "postgresql-16.1",
  "uptime_seconds": 3600
}
```

---

## JSONB Queries

VibeSQL has **native PostgreSQL JSONB support** — not an extension, the real thing:

```sql
-- Store JSONB
INSERT INTO users (data) VALUES ('{"name": "Alice", "tags": ["developer", "golang"]}'::jsonb);

-- Query JSONB fields
SELECT data->>'name' as name FROM users;

-- Filter by JSONB
SELECT * FROM users WHERE data->>'active' = 'true';

-- Array operations (native PostgreSQL)
SELECT * FROM users WHERE jsonb_array_length(data->'tags') > 1;

-- JSONB operators (PostgreSQL)
SELECT * FROM users WHERE data @> '{"active": true}'::jsonb;
SELECT * FROM users WHERE data ? 'email';
```

---

## Use Cases

### Prototyping

```bash
# Start database
npx vibesql-micro

# Create table and insert data
curl -X POST http://localhost:5173/v1/query \
  -d '{"sql": "CREATE TABLE todos (id SERIAL PRIMARY KEY, title TEXT, done BOOLEAN DEFAULT false)"}'

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

Want a visual interface? Use [vibesql-admin](https://github.com/PayEz-Net/vibesql-admin):

```bash
# Terminal 1: Database
npx vibesql-micro

# Terminal 2: Admin UI
npx vibesql-admin
# → Opens browser at http://localhost:5174
```

---

## Comparison

| Feature | VibeSQL Micro | Supabase | Railway/Neon | PlanetScale |
|---------|---------------|----------|--------------|-------------|
| Installation | `npx` command | Sign up + API keys | Sign up + deploy | Sign up + configure |
| Setup time | < 10 seconds | ~5 minutes | ~3 minutes | ~5 minutes |
| Local dev | ✅ Localhost only | ❌ Cloud sandbox only | ❌ Cloud only | ❌ Cloud only |
| Cost | ✅ Free (localhost) | Free tier + paid | Free tier + paid | Free tier + paid |
| PostgreSQL | ✅ Native PostgreSQL 16.1 | ✅ PostgreSQL | ✅ PostgreSQL | ❌ MySQL-compatible |
| JSONB | ✅ Full support | ✅ Full support | ✅ Full support | ❌ JSON only |
| Auth built-in | ❌ No* | ✅ Yes | ❌ No | ❌ No |
| Use case | Local dev, prototyping | Production apps | Production apps | Production apps |

**Note:** *VibeSQL Server (production version) includes HMAC authentication and configurable tier limits. See [VibeSQL Server](https://github.com/PayEz-Net/vibesql-server) for production deployments.

---

## Production Use

VibeSQL Micro is **production-ready** and battle-tested. Perfect for:

- **Edge computing** — AI-enhanced devices, IoT sensors, embedded systems
- **Local-first apps** — Offline-first applications with sync
- **Single-tenant deployments** — One database per customer
- **Development tools** — Build tools, CI/CD pipelines, testing frameworks
- **Desktop applications** — Electron, Tauri, native apps

**Included:**
- Comprehensive test suite
- PostgreSQL 16.1 stability and ACID guarantees
- Built-in safety checks and validation
- Production-grade reliability

**Not included (see VibeSQL Cloud):**
- Multi-instance replication
- Managed backups and point-in-time recovery
- Built-in authentication and authorization
- Horizontal scaling and load balancing

---

## Development

Clone the repo:

```bash
git clone https://github.com/PayEz-Net/vibesql-micro.git
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
- **Database:** Embedded PostgreSQL 16.1 (full PostgreSQL, not a fork)
- **HTTP:** Standard library (`net/http`)
- **Packaging:** npm (via npx)
- **Binary size:** ~68MB (includes PostgreSQL binaries)

---

## Contributing

Contributions welcome. Open an issue or pull request.

---

## License

Apache 2.0 License. See [LICENSE](LICENSE).

---

## Links

- **Website:** [vibesql.online](https://vibesql.online)
- **Admin UI:** [github.com/PayEz-Net/vibesql-admin](https://github.com/PayEz-Net/vibesql-admin)
- **Docs:** [vibesql.online/docs](https://vibesql.online/docs)
- **Discord:** [discord.gg/vibesql](https://discord.gg/vibesql)

---

Built for developers. Zero config. Just works.

---

<div align="right">
  <sub>Powered by <a href="https://idealvibe.online">IdealVibe</a></sub>
</div>
