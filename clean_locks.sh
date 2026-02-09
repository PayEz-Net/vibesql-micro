#!/bin/bash
echo "=== Clean stale PG lock files ==="
sudo rm -f /tmp/.s.PGSQL.5433.lock /tmp/.s.PGSQL.5433
sudo rm -rf /tmp/vibe-postgres-*
echo "Cleaned"
ls -la /tmp/.s.PGSQL* 2>/dev/null || echo "No PG lock files remaining"
