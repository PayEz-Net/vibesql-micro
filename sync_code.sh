#!/bin/bash
set -euo pipefail
echo "=== Current source files on Linux ==="
find /opt/vibesql -name "*.go" -type f 2>/dev/null
echo "---"
ls /opt/vibesql/cmd/vibe/ 2>/dev/null || echo "cmd/vibe/ missing"
echo "---"
ls /opt/vibesql/internal/postgres/manager.go 2>/dev/null || echo "manager.go missing"
echo "---"
ls /opt/vibesql/scripts/ 2>/dev/null || echo "scripts/ missing"
echo "---"
cat /opt/vibesql/go.mod
