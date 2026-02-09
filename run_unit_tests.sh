#!/bin/bash
set -euo pipefail
cd /opt/vibesql

echo "=== Running Unit Tests ==="
/usr/local/go/bin/go test ./internal/... -v -count=1 -timeout=60s 2>&1

echo ""
echo "=== Unit Tests Complete ==="
