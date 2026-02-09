#!/bin/bash
echo "=== What process is on 5173? ==="
sudo ss -tlnp | grep 5173
echo ""
echo "=== Process details ==="
sudo lsof -i :5173 2>/dev/null || sudo fuser -v 5173/tcp 2>/dev/null || echo "Could not identify process"
echo ""
echo "=== Force kill ==="
sudo fuser -k -9 5173/tcp 2>/dev/null || true
sleep 2
echo ""
echo "=== Check again ==="
ss -tlnp | grep 5173 || echo "Port 5173 is now free"
