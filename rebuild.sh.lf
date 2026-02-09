#!/bin/bash
set -euo pipefail
cd /opt/vibesql

echo "=== Cleaning Go build cache ==="
/usr/local/go/bin/go clean -cache

echo "=== Building VibeSQL ==="
bash scripts/build-linux.sh

echo ""
echo "=== Verifying DefaultHost ==="
strings vibe-linux-amd64 | grep -c "0.0.0.0" || echo "No 0.0.0.0 found (good)"
strings vibe-linux-amd64 | grep "127.0.0.1" | head -3

echo ""
echo "=== Binary info ==="
ls -lh vibe-linux-amd64
file vibe-linux-amd64
