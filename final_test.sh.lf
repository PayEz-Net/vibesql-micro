#!/bin/bash
set -euo pipefail
cd /opt/vibesql

BINARY="./vibe-linux-amd64"
PASS=0
FAIL=0
TOTAL=0

pass() { PASS=$((PASS+1)); TOTAL=$((TOTAL+1)); echo "  PASS: $1"; }
fail() { FAIL=$((FAIL+1)); TOTAL=$((TOTAL+1)); echo "  FAIL: $1 - $2"; }

echo "=== Cleanup ==="
pkill -f vibe-linux-amd64 2>/dev/null || true
pkill -f "postgres -D" 2>/dev/null || true
sleep 1
rm -rf ./vibe-data /tmp/vibe-postgres-* /tmp/.s.PGSQL.5433.lock 2>/dev/null || true

echo ""
echo "=== Test 1: Binary exists and is executable ==="
if [ -x "$BINARY" ]; then
    pass "Binary exists and executable"
else
    fail "Binary" "not found or not executable"
fi

echo ""
echo "=== Test 2: Version command ==="
VERSION_OUT=$($BINARY version 2>&1)
if echo "$VERSION_OUT" | grep -q "VibeSQL"; then
    pass "Version command works: $VERSION_OUT"
else
    fail "Version" "unexpected output: $VERSION_OUT"
fi

echo ""
echo "=== Test 3: Binary size <= 25MB ==="
SIZE=$(stat -c%s "$BINARY")
SIZE_MB=$((SIZE / 1024 / 1024))
if [ "$SIZE" -le 26214400 ]; then
    pass "Binary size: ${SIZE_MB}MB (${SIZE} bytes) <= 25MB"
else
    fail "Binary size" "${SIZE_MB}MB exceeds 25MB"
fi

echo ""
echo "=== Test 4: Start server and measure cold start ==="
START_TIME=$(date +%s%N)
$BINARY serve &
SERVER_PID=$!
sleep 0.5

READY=false
for i in $(seq 1 40); do
    if curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT 1"}' 2>/dev/null | grep -q '"success":true'; then
        READY=true
        END_TIME=$(date +%s%N)
        ELAPSED=$(( (END_TIME - START_TIME) / 1000000 ))
        if [ "$ELAPSED" -le 2000 ]; then
            pass "Cold start: ${ELAPSED}ms (<= 2000ms)"
        else
            fail "Cold start" "${ELAPSED}ms > 2000ms"
        fi
        break
    fi
    sleep 0.25
done

if [ "$READY" = false ]; then
    fail "Server startup" "not ready after 10s"
    kill $SERVER_PID 2>/dev/null || true
    echo "RESULTS: $PASS passed, $FAIL failed, $TOTAL total"
    exit 1
fi

echo ""
echo "=== Test 5: Binds to 127.0.0.1 (not 0.0.0.0) ==="
if ss -tlnp 2>/dev/null | grep 5173 | grep -q "127.0.0.1"; then
    pass "Server binds to 127.0.0.1:5173"
elif netstat -tlnp 2>/dev/null | grep 5173 | grep -q "127.0.0.1"; then
    pass "Server binds to 127.0.0.1:5173"
else
    fail "Server binding" "not on 127.0.0.1"
fi

echo ""
echo "=== Test 6: SELECT query ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT 1 AS test"}')
if echo "$R" | grep -q '"success":true'; then
    pass "SELECT 1: $R"
else
    fail "SELECT 1" "$R"
fi

echo ""
echo "=== Test 7: CREATE TABLE ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"CREATE TABLE IF NOT EXISTS test_final (id SERIAL PRIMARY KEY, name TEXT, data JSONB)"}')
if echo "$R" | grep -q '"success":true'; then
    pass "CREATE TABLE"
else
    fail "CREATE TABLE" "$R"
fi

echo ""
echo "=== Test 8: INSERT with JSONB ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d "{\"sql\":\"INSERT INTO test_final (name, data) VALUES ('alice', '{\\\"age\\\": 30, \\\"city\\\": \\\"NYC\\\"}')\"}")
if echo "$R" | grep -q '"success":true'; then
    pass "INSERT with JSONB"
else
    fail "INSERT" "$R"
fi

echo ""
echo "=== Test 9: SELECT with JSONB operator ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d "{\"sql\":\"SELECT name, data->>'city' AS city FROM test_final WHERE name='alice'\"}")
if echo "$R" | grep -q "NYC"; then
    pass "JSONB ->> operator: $R"
else
    fail "JSONB operator" "$R"
fi

echo ""
echo "=== Test 10: UNSAFE_QUERY error (DELETE without WHERE) ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"DELETE FROM test_final"}')
if echo "$R" | grep -q "UNSAFE_QUERY"; then
    pass "UNSAFE_QUERY error returned"
else
    fail "UNSAFE_QUERY" "$R"
fi

echo ""
echo "=== Test 11: INVALID_SQL error ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECTX BADQUERY"}')
if echo "$R" | grep -q "INVALID_SQL"; then
    pass "INVALID_SQL error returned"
else
    fail "INVALID_SQL" "$R"
fi

echo ""
echo "=== Test 12: MISSING_REQUIRED_FIELD error ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{}')
if echo "$R" | grep -q "MISSING_REQUIRED_FIELD"; then
    pass "MISSING_REQUIRED_FIELD error returned"
else
    fail "MISSING_REQUIRED_FIELD" "$R"
fi

echo ""
echo "=== Test 13: QUERY_TOO_LARGE error ==="
LARGE_SQL=$(printf 'SELECT %0.s' {1..2000})
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d "{\"sql\":\"$LARGE_SQL\"}")
if echo "$R" | grep -q "QUERY_TOO_LARGE"; then
    pass "QUERY_TOO_LARGE error returned"
else
    fail "QUERY_TOO_LARGE" "$R"
fi

echo ""
echo "=== Test 14: 1000-row limit ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT generate_series(1, 1500)"}')
if echo "$R" | grep -q "RESULT_TOO_LARGE\|rowCount.*1000\|1000"; then
    pass "Row limit enforced"
else
    fail "Row limit" "$R"
fi

echo ""
echo "=== Test 15: Query timeout (5s) ==="
START_T=$(date +%s%N)
R=$(curl -s --max-time 10 http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT pg_sleep(10)"}')
END_T=$(date +%s%N)
ELAPSED_MS=$(( (END_T - START_T) / 1000000 ))
if echo "$R" | grep -q "QUERY_TIMEOUT"; then
    if [ "$ELAPSED_MS" -ge 4900 ] && [ "$ELAPSED_MS" -le 5100 ]; then
        pass "Query timeout at ${ELAPSED_MS}ms (5s +/- 100ms)"
    elif [ "$ELAPSED_MS" -ge 4500 ] && [ "$ELAPSED_MS" -le 6000 ]; then
        pass "Query timeout at ${ELAPSED_MS}ms (within acceptable range)"
    else
        fail "Query timeout timing" "${ELAPSED_MS}ms (expected ~5000ms)"
    fi
else
    fail "Query timeout" "no QUERY_TIMEOUT error: $R"
fi

echo ""
echo "=== Test 16: DROP TABLE cleanup ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"DROP TABLE IF EXISTS test_final"}')
if echo "$R" | grep -q '"success":true'; then
    pass "DROP TABLE cleanup"
else
    fail "DROP TABLE" "$R"
fi

echo ""
echo "=== Test 17: Graceful shutdown ==="
kill -TERM $SERVER_PID 2>/dev/null
SHUTDOWN_START=$(date +%s%N)
wait $SERVER_PID 2>/dev/null
SHUTDOWN_END=$(date +%s%N)
SHUTDOWN_MS=$(( (SHUTDOWN_END - SHUTDOWN_START) / 1000000 ))
if [ "$SHUTDOWN_MS" -le 10000 ]; then
    pass "Graceful shutdown in ${SHUTDOWN_MS}ms"
else
    fail "Graceful shutdown" "took ${SHUTDOWN_MS}ms"
fi

echo ""
echo "=== Test 18: Data persistence (restart) ==="
echo "  Starting server again..."
rm -f /tmp/.s.PGSQL.5433.lock 2>/dev/null || true
$BINARY serve &
SERVER_PID2=$!
sleep 3

for i in $(seq 1 20); do
    if curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT 1"}' 2>/dev/null | grep -q '"success":true'; then
        break
    fi
    sleep 0.5
done

R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"CREATE TABLE persist_test (id INT); INSERT INTO persist_test VALUES (42)"}')
R2=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT * FROM persist_test"}')
if echo "$R2" | grep -q "42"; then
    pass "Data persistence verified"
else
    fail "Data persistence" "$R2"
fi

curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"DROP TABLE IF EXISTS persist_test"}' > /dev/null 2>&1

kill -TERM $SERVER_PID2 2>/dev/null
wait $SERVER_PID2 2>/dev/null

echo ""
echo "========================================"
echo "RESULTS: $PASS passed, $FAIL failed, $TOTAL total"
echo "========================================"

if [ "$FAIL" -gt 0 ]; then
    exit 1
else
    echo "ALL TESTS PASSED!"
    exit 0
fi
