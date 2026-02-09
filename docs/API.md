# VibeSQL API Reference

## Endpoint

```
POST /v1/query
```

All SQL queries are sent as JSON to this single endpoint.

## Request Format

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "sql": "SELECT * FROM users WHERE id = 1"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `sql` | string | Yes | SQL query to execute (max 10KB) |

## Response Format

### Success (HTTP 200)

```json
{
  "success": true,
  "rows": [
    {"id": 1, "name": "Alice", "email": "alice@example.com"}
  ],
  "rowCount": 1,
  "executionTime": 0.42
}
```

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Always `true` for successful queries |
| `rows` | array | Array of row objects (column name → value) |
| `rowCount` | integer | Number of rows returned |
| `executionTime` | float | Execution time in milliseconds |

### Error (HTTP 4xx/5xx)

```json
{
  "success": false,
  "error": {
    "code": "INVALID_SQL",
    "message": "Invalid SQL syntax",
    "detail": "PostgreSQL error: relation \"nonexistent\" does not exist"
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Always `false` for errors |
| `error.code` | string | Machine-readable error code (see [ERRORS.md](ERRORS.md)) |
| `error.message` | string | Human-readable error message |
| `error.detail` | string | Additional context (optional) |

## Query Examples

### SELECT

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT * FROM users ORDER BY id LIMIT 10"}'
```

```json
{
  "success": true,
  "rows": [
    {"id": 1, "name": "Alice", "email": "alice@example.com"},
    {"id": 2, "name": "Bob", "email": "bob@example.com"}
  ],
  "rowCount": 2,
  "executionTime": 0.85
}
```

### SELECT with WHERE

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT name, email FROM users WHERE name = '\''Alice'\''"}'
```

### SELECT with ORDER BY, LIMIT, OFFSET

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT * FROM users ORDER BY name ASC LIMIT 10 OFFSET 20"}'
```

### INSERT

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "INSERT INTO users (name, email) VALUES ('\''Charlie'\'', '\''charlie@example.com'\'')"}'
```

```json
{
  "success": true,
  "rows": null,
  "rowCount": 0,
  "executionTime": 1.12
}
```

### INSERT with RETURNING

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "INSERT INTO users (name, email) VALUES ('\''Dave'\'', '\''dave@example.com'\'') RETURNING id, name"}'
```

```json
{
  "success": true,
  "rows": [{"id": 3, "name": "Dave"}],
  "rowCount": 1,
  "executionTime": 1.34
}
```

### UPDATE (WHERE required)

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "UPDATE users SET email = '\''newemail@example.com'\'' WHERE id = 1"}'
```

UPDATE without WHERE returns `UNSAFE_QUERY` (HTTP 400):

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "UPDATE users SET email = '\''all@example.com'\''"}'
```

```json
{
  "success": false,
  "error": {
    "code": "UNSAFE_QUERY",
    "message": "Unsafe query: UPDATE without WHERE clause",
    "detail": "UPDATE queries must include a WHERE clause. Use 'WHERE 1=1' to update all rows explicitly"
  }
}
```

### DELETE (WHERE required)

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "DELETE FROM users WHERE id = 1"}'
```

### CREATE TABLE

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "CREATE TABLE products (id SERIAL PRIMARY KEY, name TEXT NOT NULL, price NUMERIC(10,2), metadata JSONB)"}'
```

### DROP TABLE

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "DROP TABLE IF EXISTS products"}'
```

### JSONB Queries

See [JSONB.md](JSONB.md) for comprehensive JSONB examples.

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT data->>'\''name'\'' AS name FROM documents WHERE data @> '\''{ \"type\": \"invoice\" }'\'' ORDER BY data->>'\''date'\'' DESC LIMIT 10"}'
```

## Limits

| Limit | Value | Error Code |
|-------|-------|------------|
| Max query size | 10KB (10,240 bytes) | `QUERY_TOO_LARGE` (413) |
| Max result rows | 1,000 | `RESULT_TOO_LARGE` (413) |
| Query timeout | 5 seconds | `QUERY_TIMEOUT` (408) |
| Max concurrent connections | 2 | — |
| HTTP read timeout | 10 seconds | — |
| HTTP write timeout | 10 seconds | — |

## HTTP Status Codes

| Status | Meaning |
|--------|---------|
| 200 | Query executed successfully |
| 400 | Invalid SQL, missing field, or unsafe query |
| 408 | Query timed out (exceeded 5 seconds) |
| 413 | Query or result too large |
| 500 | Internal server error |
| 503 | Database unavailable |

## Notes

- All queries run against the embedded PostgreSQL instance on port 5433
- The HTTP server binds to `127.0.0.1:5173` (localhost only)
- Data is stored in `./vibe-data/` relative to where `vibe serve` is run
- Results are JSON objects with column names as keys
- NULL values are returned as JSON `null`
- JSONB columns are returned as JSON objects/arrays
- `executionTime` is in milliseconds
