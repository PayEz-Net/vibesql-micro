#!/bin/bash
set -euo pipefail

echo "=== Backup embed files ==="
mkdir -p /tmp/embed_backup
cp -f /home/zenflow/vibesql/internal/postgres/embed/* /tmp/embed_backup/ 2>/dev/null || echo "No embed files to backup"
ls -la /tmp/embed_backup/

echo ""
echo "=== Kill any running vibe processes ==="
pkill -f vibe-linux-amd64 2>/dev/null || true
pkill -f "postgres -D" 2>/dev/null || true
sleep 1

echo ""
echo "=== Clean destination ==="
rm -rf /home/zenflow/vibesql-clean

echo ""
echo "=== Ready for rsync ==="
echo "DONE"
