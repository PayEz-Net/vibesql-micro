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
rm -f /tmp/.s.PGSQL.5433.lock 2>/dev/null || true

echo ""
echo "=== Starting server (warm start - data already exists) ==="
START_T=$(date +%s%N)
./vibe-linux-amd64 serve &
PID=$!

READY=false
for i in $(seq 1 40); do
    if curl -s http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT 1"}' 2>/dev/null | grep -q '"success":true'; then
        READY=true
        END_T=$(date +%s%N)
        MS=$(( (END_T - START_T) / 1000000 ))
        echo "WARM START: ${MS}ms"
        break
    fi
    sleep 0.25
done

if [ "$READY" = false ]; then
    echo "FAIL: Server not ready"
    kill $PID 2>/dev/null || true
    exit 1
fi

echo ""
echo "=== Testing query timeout ==="
START_T=$(date +%s%N)
R=$(curl -s --max-time 10 http://127.0.0.1:5173/v1/query -X POST -H "Content-Type: application/json" -d '{"sql":"SELECT pg_sleep(10)"}')
END_T=$(date +%s%N)
MS=$(( (END_T - START_T) / 1000000 ))
echo "TIMEOUT RESPONSE: $R"
echo "ELAPSED: ${MS}ms"

if echo "$R" | grep -q "QUERY_TIMEOUT"; then
    if [ "$MS" -ge 4900 ] && [ "$MS" -le 5500 ]; then
        echo "PASS: Timeout at ${MS}ms"
    else
        echo "MARGINAL: Timeout at ${MS}ms (expected 4900-5500ms)"
    fi
else
    echo "FAIL: No QUERY_TIMEOUT error"
fi

echo ""
echo "=== Cleanup ==="
kill -TERM $PID 2>/dev/null
wait $PID 2>/dev/null
echo "Done"
