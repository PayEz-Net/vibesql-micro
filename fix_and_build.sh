#!/bin/bash
set -euo pipefail

PROJ="/home/zenflow/vibesql"
EMBED="$PROJ/internal/postgres/embed"

echo "=== Step 1: Pull latest code ==="
cd "$PROJ"
git pull origin vibe-sql-micro-server-be11 || echo "Pull failed, continuing with local"

echo ""
echo "=== Step 2: Verify embed files ==="
for f in postgres_micro_linux_amd64 initdb_linux_amd64 pg_ctl_linux_amd64 libpq.so.5 share.tar.gz; do
    if [ -f "$EMBED/$f" ]; then
        SIZE=$(stat -c%s "$EMBED/$f")
        echo "  OK: $f ($SIZE bytes)"
    else
        echo "  MISSING: $f"
        exit 1
    fi
done

echo ""
echo "=== Step 3: Build binary ==="
cd "$PROJ"
bash scripts/build-linux.sh

echo ""
echo "=== Step 4: Verify binary ==="
BINARY="$PROJ/vibe-linux-amd64"
if [ -f "$BINARY" ]; then
    SIZE=$(stat -c%s "$BINARY")
    SIZE_MB=$((SIZE / 1024 / 1024))
    echo "Binary: $BINARY"
    echo "Size: ${SIZE_MB}MB ($SIZE bytes)"
    file "$BINARY"
else
    echo "ERROR: Binary not found!"
    exit 1
fi

echo ""
echo "=== Step 5: Test version command ==="
"$BINARY" version

echo ""
echo "=== Step 6: Start server (background, 10s test) ==="
rm -rf "$PROJ/vibe-data"
"$BINARY" serve &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"

sleep 5

echo ""
echo "=== Step 7: Test HTTP endpoint ==="
RESPONSE=$(curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d '{"sql": "SELECT 1 AS test"}' 2>&1) || true
echo "Response: $RESPONSE"

echo ""
echo "=== Step 8: Test CREATE TABLE ==="
RESPONSE2=$(curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d '{"sql": "CREATE TABLE integration_test (id SERIAL PRIMARY KEY, name TEXT)"}' 2>&1) || true
echo "Response: $RESPONSE2"

echo ""
echo "=== Step 9: Test INSERT ==="
RESPONSE3=$(curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d "{\"sql\": \"INSERT INTO integration_test (name) VALUES ('\''hello'\'')\"}" 2>&1) || true
echo "Response: $RESPONSE3"

echo ""
echo "=== Step 10: Test SELECT ==="
RESPONSE4=$(curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d '{"sql": "SELECT * FROM integration_test"}' 2>&1) || true
echo "Response: $RESPONSE4"

echo ""
echo "=== Step 11: Cleanup ==="
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
rm -rf "$PROJ/vibe-data"

echo ""
echo "=== DONE ==="
