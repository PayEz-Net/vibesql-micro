# VibeSQL Error Codes

All errors return a JSON response with `"success": false` and an `error` object containing `code`, `message`, and optional `detail`.

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "detail": "Additional context"
  }
}
```

## Error Code Reference

### INVALID_SQL (HTTP 400)

Returned when the SQL query has syntax errors, references undefined tables or columns, or uses unsupported functions.

**Triggers:**
- SQL syntax errors
- Undefined table (`42P01`)
- Undefined column (`42703`)
- Undefined function (`42883`)
- Data type mismatch (`42804`)
- Query doesn't start with a valid SQL keyword

**Example:**
```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELEC * FROM users"}'
```

```json
{
  "success": false,
  "error": {
    "code": "INVALID_SQL",
    "message": "Invalid SQL syntax",
    "detail": "Query must start with a valid SQL keyword (SELECT, INSERT, UPDATE, DELETE, CREATE, DROP)"
  }
}
```

**Resolution:**
- Check SQL syntax
- Verify table and column names exist
- Ensure you're using a supported SQL keyword (SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, TRUNCATE)

---

### MISSING_REQUIRED_FIELD (HTTP 400)

Returned when the `sql` field is missing or empty in the request body.

**Example:**
```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"query": "SELECT 1"}'
```

```json
{
  "success": false,
  "error": {
    "code": "MISSING_REQUIRED_FIELD",
    "message": "Missing required field: sql",
    "detail": "The 'sql' field is required"
  }
}
```

**Resolution:**
- Use `{"sql": "..."}` as the request body
- Ensure the `sql` field is not empty

---

### UNSAFE_QUERY (HTTP 400)

Returned when an UPDATE or DELETE statement lacks a WHERE clause.

**Example:**
```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "DELETE FROM users"}'
```

```json
{
  "success": false,
  "error": {
    "code": "UNSAFE_QUERY",
    "message": "Unsafe query: DELETE without WHERE clause",
    "detail": "DELETE queries must include a WHERE clause. Use 'WHERE 1=1' to delete all rows explicitly"
  }
}
```

**Resolution:**
- Add a WHERE clause to UPDATE and DELETE queries
- To affect all rows intentionally, use `WHERE 1=1`

---

### QUERY_TIMEOUT (HTTP 408)

Returned when a query exceeds the 5-second execution limit.

**Triggers:**
- Query runs longer than 5 seconds
- Query is canceled due to context deadline
- PostgreSQL SQLSTATE `57014` (query_canceled)

**Example:**
```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT pg_sleep(10)"}'
```

```json
{
  "success": false,
  "error": {
    "code": "QUERY_TIMEOUT",
    "message": "Query execution timeout",
    "detail": "Query exceeded the maximum execution time of 5 seconds"
  }
}
```

**Resolution:**
- Optimize the query (add indexes, reduce data scanned)
- Add LIMIT to constrain result size
- Break complex queries into smaller operations

---

### QUERY_TOO_LARGE (HTTP 413)

Returned when the SQL query exceeds the 10KB size limit.

**Example response:**
```json
{
  "success": false,
  "error": {
    "code": "QUERY_TOO_LARGE",
    "message": "Query too large",
    "detail": "SQL query exceeds the maximum allowed size of 10KB"
  }
}
```

**Resolution:**
- Reduce query size below 10,240 bytes
- Break large INSERT statements into multiple smaller ones
- Use parameterized values instead of inline data

---

### RESULT_TOO_LARGE (HTTP 413)

Returned when a query result exceeds 1,000 rows.

**Example:**
```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT * FROM large_table"}'
```

```json
{
  "success": false,
  "error": {
    "code": "RESULT_TOO_LARGE",
    "message": "Result set too large",
    "detail": "Query returned more than the maximum allowed 1000 rows"
  }
}
```

**Resolution:**
- Add `LIMIT 1000` (or smaller) to your query
- Use `OFFSET` for pagination
- Add WHERE clauses to filter results

---

### DOCUMENT_TOO_LARGE (HTTP 413)

Returned when a JSONB document or statement exceeds PostgreSQL's internal limits.

**Triggers:**
- PostgreSQL SQLSTATE `54000` (program_limit_exceeded)
- PostgreSQL SQLSTATE `54001` (statement_too_complex)

**Resolution:**
- Reduce JSONB document size
- Simplify deeply nested JSON structures
- Break large documents into smaller related records

---

### INTERNAL_ERROR (HTTP 500)

Returned for unexpected server errors not covered by other error codes.

**Example response:**
```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An internal error occurred",
    "detail": "..."
  }
}
```

**Resolution:**
- Check server logs for details
- Retry the request
- If persistent, restart the server

---

### SERVICE_UNAVAILABLE (HTTP 503)

Returned when the server is not ready to handle requests.

**Resolution:**
- Wait for the server to finish starting up
- Check that `vibe serve` is running
- Check server logs for startup errors

---

### DATABASE_UNAVAILABLE (HTTP 503)

Returned when the embedded PostgreSQL is unreachable.

**Triggers:**
- PostgreSQL process crashed or is not running
- Connection failure (SQLSTATE `08xxx`)
- Insufficient resources (SQLSTATE `53xxx`)
- Too many connections (SQLSTATE `53300`)

**Resolution:**
- Restart the server: `vibe serve`
- Check disk space (PostgreSQL needs space for WAL files)
- Check if another process is using port 5433

## PostgreSQL SQLSTATE Mapping

| SQLSTATE | VibeSQL Code | Description |
|----------|-------------|-------------|
| `42601` | `INVALID_SQL` | syntax_error |
| `42703` | `INVALID_SQL` | undefined_column |
| `42P01` | `INVALID_SQL` | undefined_table |
| `42P02` | `INVALID_SQL` | undefined_parameter |
| `42883` | `INVALID_SQL` | undefined_function |
| `42804` | `INVALID_SQL` | datatype_mismatch |
| `57014` | `QUERY_TIMEOUT` | query_canceled |
| `53000` | `DATABASE_UNAVAILABLE` | insufficient_resources |
| `53100` | `DATABASE_UNAVAILABLE` | disk_full |
| `53200` | `DATABASE_UNAVAILABLE` | out_of_memory |
| `53300` | `DATABASE_UNAVAILABLE` | too_many_connections |
| `53400` | `DATABASE_UNAVAILABLE` | configuration_limit_exceeded |
| `08000` | `DATABASE_UNAVAILABLE` | connection_exception |
| `08003` | `DATABASE_UNAVAILABLE` | connection_does_not_exist |
| `08006` | `DATABASE_UNAVAILABLE` | connection_failure |
| `08001` | `DATABASE_UNAVAILABLE` | sqlclient_unable_to_establish |
| `08004` | `DATABASE_UNAVAILABLE` | sqlserver_rejected_establishment |
| `54000` | `DOCUMENT_TOO_LARGE` | program_limit_exceeded |
| `54001` | `DOCUMENT_TOO_LARGE` | statement_too_complex |
