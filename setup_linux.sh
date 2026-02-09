#!/bin/bash
set -euo pipefail

echo "=== Step 1: Backup embed files ==="
mkdir -p /home/zenflow/embed_backup
cp -f /home/zenflow/vibesql/internal/postgres/embed/* /home/zenflow/embed_backup/ 2>/dev/null || true
ls -la /home/zenflow/embed_backup/

echo ""
echo "=== Step 2: Init git repo in existing vibesql dir ==="
cd /home/zenflow/vibesql

# If not a git repo, init one
if [ ! -d .git ]; then
    git init
    git checkout -b vibe-sql-micro-server-be11
fi

echo ""
echo "=== Step 3: Verify embed files present ==="
ls -la internal/postgres/embed/

echo ""
echo "=== Step 4: Build binary ==="
bash scripts/build-linux.sh

echo ""
echo "=== Step 5: Test version ==="
./vibe-linux-amd64 version

echo ""
echo "=== Step 6: Start server (background) ==="
rm -rf ./vibe-data
./vibe-linux-amd64 serve &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"

echo "Waiting 5 seconds for server to start..."
sleep 5

echo ""
echo "=== Step 7: Test SELECT 1 ==="
curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d '{"sql": "SELECT 1 AS test"}' || echo "FAILED"

echo ""
echo "=== Step 8: Test CREATE TABLE ==="
curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d '{"sql": "CREATE TABLE e2e_verify (id SERIAL PRIMARY KEY, name TEXT NOT NULL)"}' || echo "FAILED"

echo ""
echo "=== Step 9: Test INSERT ==="
curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d '{"sql": "INSERT INTO e2e_verify (name) VALUES ('\''hello'\'')"}' || echo "FAILED"

echo ""
echo "=== Step 10: Test SELECT ==="
curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d '{"sql": "SELECT * FROM e2e_verify"}' || echo "FAILED"

echo ""
echo "=== Step 11: Test error handling ==="
curl -s -X POST http://127.0.0.1:5173/v1/query \
    -H "Content-Type: application/json" \
    -d '{"sql": "DELETE FROM e2e_verify"}' || echo "FAILED"

echo ""
echo "=== Step 12: Cleanup ==="
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true
rm -rf ./vibe-data
echo "Server stopped."

echo ""
echo "=== ALL TESTS COMPLETE ==="
