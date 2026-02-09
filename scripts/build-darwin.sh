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

ARCH="${GOARCH:-$(uname -m)}"
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo "ERROR: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

OUTPUT="${PROJECT_ROOT}/vibe-darwin-${ARCH}"

echo "Building VibeSQL for macOS ${ARCH}..."
echo "  Version:    $VERSION"
echo "  Commit:     $GIT_COMMIT"
echo "  Build Date: $BUILD_DATE"

REQUIRED_EMBEDS=(
    "internal/postgres/embed/postgres_micro_darwin_${ARCH}"
    "internal/postgres/embed/initdb_darwin_${ARCH}"
    "internal/postgres/embed/pg_ctl_darwin_${ARCH}"
    "internal/postgres/embed/libpq.5.dylib"
    "internal/postgres/embed/share.tar.gz"
)

for f in "${REQUIRED_EMBEDS[@]}"; do
    if [ ! -f "$PROJECT_ROOT/$f" ]; then
        echo "ERROR: Missing embedded file: $f"
        echo "Run the PostgreSQL build first: bash build/build_postgres_darwin.sh"
        exit 1
    fi
done

cd "$PROJECT_ROOT"
CGO_ENABLED=0 GOOS=darwin GOARCH="$ARCH" go build \
    -ldflags="$LDFLAGS" \
    -o "$OUTPUT" \
    ./cmd/vibe

SIZE=$(stat -f%z "$OUTPUT" 2>/dev/null || stat -c%s "$OUTPUT")
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
