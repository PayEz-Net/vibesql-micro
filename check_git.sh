#!/bin/bash
cd /home/zenflow/vibesql
echo "=== REMOTES ==="
git remote -v
echo "=== BRANCH ==="
git branch
echo "=== STATUS ==="
git status --short | head -20
echo "=== LOG ==="
git log --oneline -3
