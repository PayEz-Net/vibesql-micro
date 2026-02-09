#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

VERSION="${VERSION:-1.0.0}"
GIT_COMMIT="${GIT_COMMIT:-$(git -C "$PROJECT_ROOT" rev-parse --short HEAD 2>/dev/null || echo dev)}"
BUILD_DATE="${BUILD_DATE:-$(date -u '+%Y-%m-%d_%H:%M:%S')}"

LDFLAGS="-s -w"
LDFLAGS="$LDFLAGS -X github.com/vibesql/vibe/internal/version.Version=$VERSION"
LDFLAGS="$LDFLAGS -X github.com/vibesql/vibe/internal/version.GitCommit=$GIT_COMMIT"
LDFLAGS="$LDFLAGS -X github.com/vibesql/vibe/internal/version.BuildDate=$BUILD_DATE"

OUTPUT="${PROJECT_ROOT}/vibe-linux-amd64"

echo "Building VibeSQL for Linux x64..."
echo "  Version:    $VERSION"
echo "  Commit:     $GIT_COMMIT"
echo "  Build Date: $BUILD_DATE"

REQUIRED_EMBEDS=(
    "internal/postgres/embed/postgres_micro_linux_amd64"
    "internal/postgres/embed/initdb_linux_amd64"
    "internal/postgres/embed/pg_ctl_linux_amd64"
    "internal/postgres/embed/libpq.so.5"
    "internal/postgres/embed/share.tar.gz"
)

for f in "${REQUIRED_EMBEDS[@]}"; do
    if [ ! -f "$PROJECT_ROOT/$f" ]; then
        echo "ERROR: Missing embedded file: $f"
        echo "Run the PostgreSQL build first: bash build/build_postgres.sh"
        exit 1
    fi
done

cd "$PROJECT_ROOT"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build \
    -ldflags="$LDFLAGS" \
    -o "$OUTPUT" \
    ./cmd/vibe

SIZE=$(stat -c%s "$OUTPUT")
SIZE_MB=$((SIZE / 1024 / 1024))

echo ""
echo "Build complete: $OUTPUT"
echo "  Size: ${SIZE_MB}MB (${SIZE} bytes)"

if [ "$SIZE" -gt 26214400 ]; then
    echo "ERROR: Binary exceeds 25MB hard limit!"
    exit 1
elif [ "$SIZE" -gt 20971520 ]; then
    echo "WARNING: Binary exceeds 20MB preferred target"
else
    echo "  Size OK (under 20MB target)"
fi

file "$OUTPUT"
echo ""
echo "To test: $OUTPUT serve"
