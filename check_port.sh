#!/bin/bash
echo "=== What's on port 5433? ==="
ss -tlnp | grep 5433
echo ""
echo "=== What's on port 5432? ==="
ss -tlnp | grep 5432
echo ""
echo "=== PostgreSQL processes ==="
ps aux | grep postgres | grep -v grep
echo ""
echo "=== Try port 5434 ==="
ss -tlnp | grep 5434 || echo "Port 5434 is free"
echo ""
echo "=== Try port 5435 ==="
ss -tlnp | grep 5435 || echo "Port 5435 is free"
