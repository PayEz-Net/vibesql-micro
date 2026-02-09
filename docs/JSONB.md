# VibeSQL JSONB Guide

VibeSQL supports all 9 PostgreSQL JSONB operators through standard SQL syntax. This guide covers each operator with practical examples.

## Setup

Create a table with a JSONB column:

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "CREATE TABLE documents (id SERIAL PRIMARY KEY, data JSONB NOT NULL)"}'
```

Insert sample data:

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "INSERT INTO documents (data) VALUES ('\''{ \"name\": \"Alice\", \"age\": 30, \"role\": \"admin\", \"tags\": [\"staff\", \"manager\"], \"address\": { \"city\": \"Portland\", \"state\": \"OR\" } }'\'')"}'
```

## Operators

### `->` Get Field as JSON

Returns a JSON element by key (object) or index (array). Result is JSON type.

```sql
SELECT data->'name' FROM documents
-- Returns: "Alice" (as JSON string, with quotes)

SELECT data->'tags'->0 FROM documents
-- Returns: "staff" (first array element)

SELECT data->'address' FROM documents
-- Returns: {"city": "Portland", "state": "OR"}
```

### `->>` Get Field as Text

Returns a JSON element as plain text. Most common operator for extracting values.

```sql
SELECT data->>'name' FROM documents
-- Returns: Alice (as text, no quotes)

SELECT data->>'age' FROM documents
-- Returns: 30 (as text)

SELECT data->'address'->>'city' FROM documents
-- Returns: Portland
```

### `#>` Get Path as JSON

Navigates a JSON path and returns JSON. Path is specified as a text array.

```sql
SELECT data #> '{address,city}' FROM documents
-- Returns: "Portland" (as JSON string)

SELECT data #> '{tags,0}' FROM documents
-- Returns: "staff" (as JSON string)
```

### `#>>` Get Path as Text

Navigates a JSON path and returns text.

```sql
SELECT data #>> '{address,city}' FROM documents
-- Returns: Portland (as text)

SELECT data #>> '{address,state}' FROM documents
-- Returns: OR
```

### `@>` Contains

Tests if the left JSONB value contains the right JSONB value. Useful for filtering.

```sql
SELECT * FROM documents WHERE data @> '{"role": "admin"}'
-- Returns rows where data contains role=admin

SELECT * FROM documents WHERE data @> '{"tags": ["manager"]}'
-- Returns rows where tags array contains "manager"

SELECT * FROM documents WHERE data->'address' @> '{"state": "OR"}'
-- Returns rows where address.state = OR
```

### `<@` Contained By

Tests if the left JSONB value is contained by the right. Reverse of `@>`.

```sql
SELECT * FROM documents
WHERE '{"role": "admin", "name": "Alice"}' <@ data
-- Returns rows where data contains both role=admin AND name=Alice
```

### `?` Key Exists

Tests if a key exists in the top-level of a JSONB object.

```sql
SELECT * FROM documents WHERE data ? 'name'
-- Returns rows that have a "name" key

SELECT * FROM documents WHERE data ? 'phone'
-- Returns rows that have a "phone" key (empty if none)
```

### `?|` Any Key Exists

Tests if any of the given keys exist.

```sql
SELECT * FROM documents WHERE data ?| array['phone', 'email', 'name']
-- Returns rows that have at least one of: phone, email, or name
```

### `?&` All Keys Exist

Tests if all of the given keys exist.

```sql
SELECT * FROM documents WHERE data ?& array['name', 'age', 'role']
-- Returns rows that have ALL of: name, age, and role
```

## Common Patterns

### Filter by Nested Value

```sql
SELECT data->>'name' AS name
FROM documents
WHERE data->'address'->>'city' = 'Portland'
```

### Sort by JSONB Field

```sql
SELECT data->>'name' AS name, (data->>'age')::int AS age
FROM documents
ORDER BY (data->>'age')::int DESC
```

### Count by JSONB Field

```sql
SELECT data->>'role' AS role, COUNT(*) AS count
FROM documents
GROUP BY data->>'role'
```

### Check Array Contains Value

```sql
SELECT * FROM documents
WHERE data->'tags' @> '"manager"'
```

### Extract Multiple Fields

```sql
SELECT
  data->>'name' AS name,
  data->>'role' AS role,
  data #>> '{address,city}' AS city
FROM documents
```

### Filter with Multiple JSONB Conditions

```sql
SELECT data->>'name' AS name
FROM documents
WHERE data @> '{"role": "admin"}'
  AND data ? 'address'
  AND (data->>'age')::int > 25
```

## Data Types

When extracting JSONB values:

| Operator | Returns | Use When |
|----------|---------|----------|
| `->` | JSONB | You need JSON for further operations |
| `->>` | text | You need the value as a string |
| `#>` | JSONB | You need a nested value as JSON |
| `#>>` | text | You need a nested value as text |

Cast text to other types as needed:

```sql
(data->>'age')::int           -- Cast to integer
(data->>'price')::numeric     -- Cast to decimal
(data->>'active')::boolean    -- Cast to boolean
```

## Cleanup

```bash
curl -X POST http://127.0.0.1:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "DROP TABLE IF EXISTS documents"}'
```
