#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
VIBE="$PROJECT_ROOT/vibe-linux-amd64"
VIBE_DATA="$PROJECT_ROOT/vibe-test-data"
API="http://127.0.0.1:5173/v1/query"
PASSED=0
FAILED=0
TOTAL=0

cleanup() {
    pkill -f "vibe-linux-amd64" 2>/dev/null || true
    sleep 1
    rm -rf "$VIBE_DATA"
}

assert_contains() {
    local response="$1"
    local expected="$2"
    local test_name="$3"
    TOTAL=$((TOTAL + 1))
    if echo "$response" | grep -q "$expected"; then
        echo "  PASS: $test_name"
        PASSED=$((PASSED + 1))
    else
        echo "  FAIL: $test_name"
        echo "    Expected to contain: $expected"
        echo "    Got: $response"
        FAILED=$((FAILED + 1))
    fi
}

assert_status() {
    local status="$1"
    local expected="$2"
    local test_name="$3"
    TOTAL=$((TOTAL + 1))
    if [ "$status" = "$expected" ]; then
        echo "  PASS: $test_name"
        PASSED=$((PASSED + 1))
    else
        echo "  FAIL: $test_name (status=$status, expected=$expected)"
        FAILED=$((FAILED + 1))
    fi
}

trap cleanup EXIT

echo "=== VibeSQL Integration Test ==="
echo ""

cleanup

echo "1. Starting VibeSQL server..."
cd "$PROJECT_ROOT"
"$VIBE" serve &
VIBE_PID=$!

echo "  Waiting for server to be ready..."
for i in $(seq 1 30); do
    if curl -s -o /dev/null -w "%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "SELECT 1"}' 2>/dev/null | grep -q "200"; then
        echo "  Server ready after ${i}s"
        break
    fi
    if ! kill -0 $VIBE_PID 2>/dev/null; then
        echo "FATAL: Server process died"
        exit 1
    fi
    sleep 1
done

if ! kill -0 $VIBE_PID 2>/dev/null; then
    echo "FATAL: Server failed to start"
    exit 1
fi
echo "  Server started (PID=$VIBE_PID)"

echo ""
echo "2. Testing version endpoint..."
RESP=$(curl -s http://127.0.0.1:5173/health 2>/dev/null || echo "CONN_FAIL")
if [ "$RESP" = "CONN_FAIL" ]; then
    RESP=$(curl -s http://127.0.0.1:5173/ 2>/dev/null || echo "CONN_FAIL")
fi

echo ""
echo "3. Testing CREATE TABLE..."
RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "CREATE TABLE test_users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, email TEXT)"}')
STATUS=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
assert_status "$STATUS" "200" "CREATE TABLE returns 200"
echo "    Response: $BODY"

echo ""
echo "4. Testing INSERT..."
RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "INSERT INTO test_users (name, email) VALUES ('\''Alice'\'', '\''alice@test.com'\'')"}')
STATUS=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
assert_status "$STATUS" "200" "INSERT returns 200"
echo "    Response: $BODY"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "INSERT INTO test_users (name, email) VALUES ('\''Bob'\'', '\''bob@test.com'\'')"}')
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "200" "Second INSERT returns 200"

echo ""
echo "5. Testing SELECT..."
RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "SELECT * FROM test_users ORDER BY id"}')
STATUS=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
assert_status "$STATUS" "200" "SELECT returns 200"
assert_contains "$BODY" "Alice" "SELECT contains Alice"
assert_contains "$BODY" "Bob" "SELECT contains Bob"
echo "    Response: $BODY"

echo ""
echo "6. Testing UPDATE..."
RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "UPDATE test_users SET email = '\''alice_updated@test.com'\'' WHERE name = '\''Alice'\''"}')
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "200" "UPDATE returns 200"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "SELECT email FROM test_users WHERE name = '\''Alice'\''"}')
BODY=$(echo "$RESP" | sed '$d')
assert_contains "$BODY" "alice_updated" "UPDATE persisted correctly"

echo ""
echo "7. Testing DELETE..."
RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "DELETE FROM test_users WHERE name = '\''Bob'\''"}')
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "200" "DELETE returns 200"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "SELECT COUNT(*) as cnt FROM test_users"}')
BODY=$(echo "$RESP" | sed '$d')
assert_contains "$BODY" "1" "DELETE removed one row"

echo ""
echo "8. Testing JSONB..."
RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "CREATE TABLE test_jsonb (id SERIAL PRIMARY KEY, data JSONB)"}')
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "200" "CREATE JSONB TABLE returns 200"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d "{\"sql\": \"INSERT INTO test_jsonb (data) VALUES ('{\\\"name\\\": \\\"test\\\", \\\"tags\\\": [\\\"a\\\", \\\"b\\\"]}')\"}")
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "200" "INSERT JSONB returns 200"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d "{\"sql\": \"SELECT data->>'name' as name FROM test_jsonb\"}")
STATUS=$(echo "$RESP" | tail -1)
BODY=$(echo "$RESP" | sed '$d')
assert_status "$STATUS" "200" "JSONB ->> query returns 200"
assert_contains "$BODY" "test" "JSONB extraction returns correct value"
echo "    Response: $BODY"

echo ""
echo "9. Testing error handling..."
RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "SELECT * FROM nonexistent_table"}')
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "400" "Invalid table returns 400"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "INVALID SQL SYNTAX HERE !!!"}')
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "400" "Invalid SQL returns 400"

echo ""
echo "10. Testing DROP TABLE cleanup..."
RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "DROP TABLE test_users"}')
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "200" "DROP TABLE returns 200"

RESP=$(curl -s -w "\n%{http_code}" -X POST "$API" -H "Content-Type: application/json" -d '{"sql": "DROP TABLE test_jsonb"}')
STATUS=$(echo "$RESP" | tail -1)
assert_status "$STATUS" "200" "DROP JSONB TABLE returns 200"

echo ""
echo "================================"
echo "Results: $PASSED/$TOTAL passed, $FAILED failed"
echo "================================"

if [ "$FAILED" -gt 0 ]; then
    exit 1
fi
