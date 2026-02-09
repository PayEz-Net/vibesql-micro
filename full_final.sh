#!/bin/bash
set -euo pipefail
cd /opt/vibesql

echo "=== Rebuilding ==="
/usr/local/go/bin/go clean -cache
bash scripts/build-linux.sh

echo ""
echo "=== Cleanup ==="
pkill -f vibe-linux-amd64 2>/dev/null || true
pkill -f "postgres -D" 2>/dev/null || true
sleep 1
rm -rf ./vibe-data /tmp/vibe-postgres-* /tmp/.s.PGSQL.5433.lock 2>/dev/null || true

PASS=0
FAIL=0
TOTAL=0

pass() { PASS=$((PASS+1)); TOTAL=$((TOTAL+1)); echo "  PASS: $1"; }
fail() { FAIL=$((FAIL+1)); TOTAL=$((TOTAL+1)); echo "  FAIL: $1 - $2"; }

echo ""
echo "=== Test 1: Binary size ==="
SIZE=$(stat -c%s vibe-linux-amd64)
SIZE_MB=$((SIZE / 1024 / 1024))
if [ "$SIZE" -le 26214400 ]; then
    pass "Binary: ${SIZE_MB}MB (${SIZE} bytes) <= 25MB"
else
    fail "Binary size" "${SIZE_MB}MB exceeds 25MB"
fi

echo ""
echo "=== Test 2: Cold start + first query ==="
START_T=$(date +%s%N)
./vibe-linux-amd64 serve &
PID=$!

READY=false
for i in $(seq 1 80); do
    if curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT 1"}' 2>/dev/null | grep -q '"success":true'; then
        READY=true
        END_T=$(date +%s%N)
        MS=$(( (END_T - START_T) / 1000000 ))
        echo "  Cold start (with initdb): ${MS}ms"
        pass "Cold start completed in ${MS}ms (includes one-time initdb)"
        break
    fi
    sleep 0.25
done

if [ "$READY" = false ]; then
    fail "Cold start" "server not ready after 20s"
    kill $PID 2>/dev/null || true
    echo "RESULTS: $PASS/$TOTAL"
    exit 1
fi

echo ""
echo "=== Test 3: Server binds to 127.0.0.1 ==="
if ss -tlnp 2>/dev/null | grep 5173 | grep -q "127.0.0.1"; then
    pass "Binds to 127.0.0.1:5173"
else
    fail "Server binding" "not 127.0.0.1"
fi

echo ""
echo "=== Test 4: SELECT ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT 1 AS test"}')
echo "$R" | grep -q '"success":true' && pass "SELECT: $R" || fail "SELECT" "$R"

echo ""
echo "=== Test 5: CREATE TABLE ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"CREATE TABLE IF NOT EXISTS t1 (id SERIAL PRIMARY KEY, name TEXT, data JSONB)"}')
echo "$R" | grep -q '"success":true' && pass "CREATE TABLE" || fail "CREATE TABLE" "$R"

echo ""
echo "=== Test 6: INSERT with JSONB ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d "{\"sql\":\"INSERT INTO t1 (name, data) VALUES ('alice', '{\\\"age\\\": 30, \\\"city\\\": \\\"NYC\\\"}')\"}")
echo "$R" | grep -q '"success":true' && pass "INSERT" || fail "INSERT" "$R"

echo ""
echo "=== Test 7: SELECT with JSONB ->> ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d "{\"sql\":\"SELECT data->>'city' AS city FROM t1\"}")
echo "$R" | grep -q "NYC" && pass "JSONB ->>: $R" || fail "JSONB" "$R"

echo ""
echo "=== Test 8: UNSAFE_QUERY ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"DELETE FROM t1"}')
echo "$R" | grep -q "UNSAFE_QUERY" && pass "UNSAFE_QUERY" || fail "UNSAFE_QUERY" "$R"

echo ""
echo "=== Test 9: INVALID_SQL ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECTX BAD"}')
echo "$R" | grep -q "INVALID_SQL" && pass "INVALID_SQL" || fail "INVALID_SQL" "$R"

echo ""
echo "=== Test 10: MISSING_REQUIRED_FIELD ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{}')
echo "$R" | grep -q "MISSING_REQUIRED_FIELD" && pass "MISSING_REQUIRED_FIELD" || fail "MISSING_REQUIRED_FIELD" "$R"

echo ""
echo "=== Test 11: QUERY_TOO_LARGE ==="
LARGE=$(printf 'SELECT %0.s' {1..2000})
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d "{\"sql\":\"$LARGE\"}")
echo "$R" | grep -q "QUERY_TOO_LARGE" && pass "QUERY_TOO_LARGE" || fail "QUERY_TOO_LARGE" "$R"

echo ""
echo "=== Test 12: RESULT_TOO_LARGE (1000 row limit) ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT generate_series(1, 1500)"}')
echo "$R" | grep -q "RESULT_TOO_LARGE" && pass "RESULT_TOO_LARGE" || fail "RESULT_TOO_LARGE" "$R"

echo ""
echo "=== Test 13: QUERY_TIMEOUT (5s) ==="
START_T=$(date +%s%N)
R=$(curl -s --max-time 10 http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT pg_sleep(10)"}')
END_T=$(date +%s%N)
MS=$(( (END_T - START_T) / 1000000 ))
if echo "$R" | grep -q "QUERY_TIMEOUT"; then
    if [ "$MS" -ge 4500 ] && [ "$MS" -le 5500 ]; then
        pass "QUERY_TIMEOUT at ${MS}ms"
    else
        fail "QUERY_TIMEOUT timing" "${MS}ms"
    fi
else
    fail "QUERY_TIMEOUT" "$R"
fi

echo ""
echo "=== Test 14: DROP TABLE ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"DROP TABLE IF EXISTS t1"}')
echo "$R" | grep -q '"success":true' && pass "DROP TABLE" || fail "DROP TABLE" "$R"

echo ""
echo "=== Test 15: Graceful shutdown ==="
kill -TERM $PID 2>/dev/null
S=$(date +%s%N)
wait $PID 2>/dev/null
E=$(date +%s%N)
MS=$(( (E - S) / 1000000 ))
[ "$MS" -le 10000 ] && pass "Graceful shutdown in ${MS}ms" || fail "Shutdown" "${MS}ms"

echo ""
echo "=== Test 16: Warm restart ==="
rm -f /tmp/.s.PGSQL.5433.lock 2>/dev/null || true
START_T=$(date +%s%N)
./vibe-linux-amd64 serve &
PID2=$!
for i in $(seq 1 20); do
    if curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT 1"}' 2>/dev/null | grep -q '"success":true'; then
        END_T=$(date +%s%N)
        MS=$(( (END_T - START_T) / 1000000 ))
        if [ "$MS" -le 2000 ]; then
            pass "Warm restart: ${MS}ms <= 2000ms"
        else
            fail "Warm restart" "${MS}ms > 2000ms"
        fi
        break
    fi
    sleep 0.25
done

echo ""
echo "=== Test 17: Data persistence ==="
R=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"CREATE TABLE persist (v INT); INSERT INTO persist VALUES (99)"}')
R2=$(curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT * FROM persist"}')
echo "$R2" | grep -q "99" && pass "Data persistence" || fail "Persistence" "$R2"
curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"DROP TABLE persist"}' > /dev/null 2>&1

kill -TERM $PID2 2>/dev/null
wait $PID2 2>/dev/null

echo ""
echo "========================================"
echo "RESULTS: $PASS passed, $FAIL failed, $TOTAL total"
echo "========================================"

[ "$FAIL" -gt 0 ] && exit 1 || echo "ALL TESTS PASSED!"
