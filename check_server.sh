#!/bin/bash
echo "=== server.go DefaultHost ==="
grep -n "DefaultHost" /opt/vibesql/internal/server/server.go
echo ""
echo "=== What's on port 5173? ==="
ss -tlnp | grep 5173 || echo "Nothing on 5173"
echo ""
echo "=== server.go Start function ==="
grep -A5 "func.*Start" /opt/vibesql/internal/server/server.go | head -15
