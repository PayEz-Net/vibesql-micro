#!/bin/bash
# Build orchestration script for VibeSQL Local
# Builds minimal PostgreSQL binary and Go wrapper

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
PG_VERSION="${PG_VERSION:-16.1}"
BUILD_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$BUILD_DIR")"
OUTPUT_DIR="$BUILD_DIR/output"
POSTGRES_BINARY="$OUTPUT_DIR/postgres_micro"

echo -e "${GREEN}VibeSQL Build Script${NC}"
echo "================================"
echo "PostgreSQL Version: $PG_VERSION"
echo "Build Directory: $BUILD_DIR"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Step 1: Build PostgreSQL micro binary
echo -e "${YELLOW}Step 1: Building minimal PostgreSQL binary...${NC}"
docker build \
    --build-arg PG_VERSION=$PG_VERSION \
    -f "$BUILD_DIR/Dockerfile.postgres" \
    -t vibesql-postgres:micro \
    "$PROJECT_ROOT"

# Extract binary from Docker image
echo -e "${YELLOW}Step 2: Extracting postgres binary...${NC}"
docker create --name vibesql-temp vibesql-postgres:micro
docker cp vibesql-temp:/postgres "$POSTGRES_BINARY"
docker rm vibesql-temp

# Verify binary size
if [ -f "$POSTGRES_BINARY" ]; then
    SIZE=$(stat -c%s "$POSTGRES_BINARY" 2>/dev/null || stat -f%z "$POSTGRES_BINARY")
    SIZE_MB=$((SIZE / 1024 / 1024))
    
    echo -e "${GREEN}✓ PostgreSQL binary built: $POSTGRES_BINARY${NC}"
    echo "  Size: ${SIZE_MB}MB"
    
    if [ $SIZE_MB -gt 20 ]; then
        echo -e "${RED}⚠ WARNING: Binary size exceeds 20MB target (hard limit: 25MB)${NC}"
    else
        echo -e "${GREEN}✓ Binary size within target (≤20MB)${NC}"
    fi
else
    echo -e "${RED}✗ Failed to extract PostgreSQL binary${NC}"
    exit 1
fi

# Step 3: Copy to embed directory
echo -e "${YELLOW}Step 3: Copying to embed directory...${NC}"
mkdir -p "$PROJECT_ROOT/internal/postgres/embed"
cp "$POSTGRES_BINARY" "$PROJECT_ROOT/internal/postgres/embed/postgres_micro_linux_amd64"
echo -e "${GREEN}✓ Binary ready for embedding${NC}"

# Step 4: Build Go binary (if Go is available)
if command -v go &> /dev/null; then
    echo -e "${YELLOW}Step 4: Building Go binary...${NC}"
    cd "$PROJECT_ROOT"
    go build -o "$OUTPUT_DIR/vibe" ./cmd/vibe
    
    if [ -f "$OUTPUT_DIR/vibe" ]; then
        VIBE_SIZE=$(stat -c%s "$OUTPUT_DIR/vibe" 2>/dev/null || stat -f%z "$OUTPUT_DIR/vibe")
        VIBE_SIZE_MB=$((VIBE_SIZE / 1024 / 1024))
        
        echo -e "${GREEN}✓ VibeSQL binary built: $OUTPUT_DIR/vibe${NC}"
        echo "  Size: ${VIBE_SIZE_MB}MB"
        
        TOTAL_SIZE=$((SIZE_MB + VIBE_SIZE_MB))
        echo ""
        echo "Total estimated binary size: ${TOTAL_SIZE}MB"
        
        if [ $TOTAL_SIZE -gt 25 ]; then
            echo -e "${RED}⚠ CRITICAL: Total size exceeds 25MB hard limit!${NC}"
            exit 1
        elif [ $TOTAL_SIZE -gt 20 ]; then
            echo -e "${YELLOW}⚠ WARNING: Total size exceeds 20MB target${NC}"
        else
            echo -e "${GREEN}✓ Total size within target (≤20MB)${NC}"
        fi
    fi
else
    echo -e "${YELLOW}⚠ Go not found, skipping Go build${NC}"
    echo "  Install Go 1.21+ to build complete binary"
fi

echo ""
echo -e "${GREEN}Build complete!${NC}"
echo "Outputs:"
echo "  - PostgreSQL: $POSTGRES_BINARY"
[ -f "$OUTPUT_DIR/vibe" ] && echo "  - VibeSQL: $OUTPUT_DIR/vibe"
