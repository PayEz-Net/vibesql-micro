#!/bin/bash
set -euo pipefail

echo "=== Kill old vibe process on 5173 ==="
sudo fuser -k 5173/tcp 2>/dev/null || true
sleep 1
ss -tlnp | grep 5173 || echo "Port 5173 is free"

echo ""
echo "=== Kill old vibe binary ==="
sudo pkill -f "/opt/vibesql/vibe" 2>/dev/null || true
sudo pkill -f "vibe serve" 2>/dev/null || true
sleep 1

echo ""
echo "=== Sync will happen via scp ==="
echo "Ready for code sync"
