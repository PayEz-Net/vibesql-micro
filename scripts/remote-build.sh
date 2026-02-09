#!/bin/bash
set -euo pipefail

SSH_KEY="${SSH_KEY:-$HOME/.ssh/zenflow_93}"
SSH_HOST="${SSH_HOST:-zenflow@10.0.0.93}"
REMOTE_DIR="${REMOTE_DIR:-/home/zenflow/vibesql}"
LOCAL_PROJECT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "=== VibeSQL Remote Build ==="
echo "  Host:       $SSH_HOST"
echo "  Remote Dir: $REMOTE_DIR"
echo "  Local Dir:  $LOCAL_PROJECT"
echo ""

SSH_CMD="ssh -i $SSH_KEY -o StrictHostKeyChecking=no $SSH_HOST"
SCP_CMD="scp -i $SSH_KEY -o StrictHostKeyChecking=no"

echo "Syncing source code to remote..."
$SCP_CMD -r "$LOCAL_PROJECT/cmd" "$SSH_HOST:$REMOTE_DIR/"
$SCP_CMD -r "$LOCAL_PROJECT/internal" "$SSH_HOST:$REMOTE_DIR/"
$SCP_CMD "$LOCAL_PROJECT/go.mod" "$SSH_HOST:$REMOTE_DIR/"
[ -f "$LOCAL_PROJECT/go.sum" ] && $SCP_CMD "$LOCAL_PROJECT/go.sum" "$SSH_HOST:$REMOTE_DIR/"
$SCP_CMD -r "$LOCAL_PROJECT/scripts" "$SSH_HOST:$REMOTE_DIR/"
$SCP_CMD -r "$LOCAL_PROJECT/Tests" "$SSH_HOST:$REMOTE_DIR/"

echo "Building on remote..."
$SSH_CMD "cd $REMOTE_DIR && bash scripts/build-linux.sh"

echo "Copying binary back..."
$SCP_CMD "$SSH_HOST:$REMOTE_DIR/vibe-linux-amd64" "$LOCAL_PROJECT/vibe-linux-amd64"

SIZE=$(stat -c%s "$LOCAL_PROJECT/vibe-linux-amd64" 2>/dev/null || stat -f%z "$LOCAL_PROJECT/vibe-linux-amd64")
echo ""
echo "Binary downloaded: $LOCAL_PROJECT/vibe-linux-amd64"
echo "  Size: $((SIZE / 1024 / 1024))MB (${SIZE} bytes)"
