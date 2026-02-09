#!/bin/bash
set -euo pipefail

PROJ="/home/zenflow/vibesql"
cd "$PROJ"

echo "=== Kill any running vibe/postgres ==="
pkill -f vibe-linux-amd64 2>/dev/null || true
pkill -f "postgres -D" 2>/dev/null || true
sleep 1

echo ""
echo "=== Clean old data dir (may be root-owned) ==="
sudo rm -rf ./vibe-data /tmp/vibe-postgres-*
echo "Cleaned"

echo ""
echo "=== Check port 5433 ==="
ss -tlnp 2>/dev/null | grep 5433 || echo "Port 5433 is free"

echo ""
echo "=== Verify all 5 embed files ==="
for f in postgres_micro_linux_amd64 initdb_linux_amd64 pg_ctl_linux_amd64 libpq.so.5 share.tar.gz; do
    if [ -f "internal/postgres/embed/$f" ]; then
        SIZE=$(stat -c%s "internal/postgres/embed/$f")
        echo "  OK: $f ($SIZE bytes)"
    else
        echo "  MISSING: $f"
        exit 1
    fi
done

echo ""
echo "=== Build ==="
bash scripts/build-linux.sh

echo ""
echo "=== Test version ==="
./vibe-linux-amd64 version

echo ""
echo "=== Start server ==="
./vibe-linux-amd64 serve > /tmp/vibe_test.log 2>&1 &
SERVER_PID=$!
echo "PID: $SERVER_PID"
echo "Waiting 10 seconds for init + start..."
sleep 10

echo ""
echo "=== Server logs ==="
cat /tmp/vibe_test.log

echo ""
echo "=== Test: SELECT 1 ==="
curl -s -m 5 -X POST http://127.0.0.1:5173/v1/query -H "Content-Type: application/json" -d '{"sql":"SELECT 1 AS test"}' && echo "" || echo "CURL FAILED"

echo ""
echo "=== Test: CREATE TABLE ==="
curl -s -m 5 -X POST http://127.0.0.1:5173/v1/query -H "Content-Type: application/json" -d '{"sql":"CREATE TABLE e2e_test (id SERIAL PRIMARY KEY, name TEXT NOT NULL, data JSONB)"}' && echo "" || echo "CURL FAILED"

echo ""
echo "=== Test: INSERT ==="
curl -s -m 5 -X POST http://127.0.0.1:5173/v1/query -H "Content-Type: application/json" -d "{\"sql\":\"INSERT INTO e2e_test (name, data) VALUES ('hello', '{\\\"key\\\": \\\"value\\\"}')\"}" && echo "" || echo "CURL FAILED"

echo ""
echo "=== Test: SELECT * ==="
curl -s -m 5 -X POST http://127.0.0.1:5173/v1/query -H "Content-Type: application/json" -d '{"sql":"SELECT * FROM e2e_test"}' && echo "" || echo "CURL FAILED"

echo ""
echo "=== Test: JSONB operator ==="
curl -s -m 5 -X POST http://127.0.0.1:5173/v1/query -H "Content-Type: application/json" -d "{\"sql\":\"SELECT name, data->>'key' AS key_value FROM e2e_test\"}" && echo "" || echo "CURL FAILED"

echo ""
echo "=== Test: UNSAFE QUERY (should return error) ==="
curl -s -m 5 -X POST http://127.0.0.1:5173/v1/query -H "Content-Type: application/json" -d '{"sql":"DELETE FROM e2e_test"}' && echo "" || echo "CURL FAILED"

echo ""
echo "=== Test: DROP TABLE ==="
curl -s -m 5 -X POST http://127.0.0.1:5173/v1/query -H "Content-Type: application/json" -d '{"sql":"DROP TABLE e2e_test"}' && echo "" || echo "CURL FAILED"

echo ""
echo "=== Cleanup ==="
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
sudo rm -rf ./vibe-data /tmp/vibe-postgres-* /tmp/vibe_test.log
echo "Done."

echo ""
echo "=== ALL TESTS COMPLETE ==="
