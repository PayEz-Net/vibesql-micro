#!/bin/bash
echo "=== Kill old vibesql postgres processes ==="
sudo pkill -u vibesql -f postgres || true
sleep 2
echo ""
echo "=== Check port 5433 ==="
ss -tlnp | grep 5433 || echo "Port 5433 is now free"
echo ""
echo "=== Check for any remaining postgres ==="
ps aux | grep "vibe-postgres" | grep -v grep || echo "No vibe-postgres processes"
