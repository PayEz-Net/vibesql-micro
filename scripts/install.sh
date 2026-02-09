#!/usr/bin/env bash
set -euo pipefail

VERSION="${VIBESQL_VERSION:-latest}"
BINARY_NAME="vibe"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}VibeSQL Installer${NC}"
echo "================================"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        echo "VibeSQL supports: x86_64 (amd64), arm64"
        exit 1
        ;;
esac

case "$OS" in
    linux)
        INSTALL_DIR="${VIBESQL_INSTALL_DIR:-/usr/local/bin}"
        BINARY_FILE="vibe-linux-${ARCH}"
        ;;
    darwin)
        INSTALL_DIR="${VIBESQL_INSTALL_DIR:-/usr/local/bin}"
        BINARY_FILE="vibe-darwin-${ARCH}"
        ;;
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        echo "VibeSQL supports: Linux, macOS"
        echo ""
        echo "For Windows, use install.ps1 instead:"
        echo "  powershell -ExecutionPolicy Bypass -File install.ps1"
        exit 1
        ;;
esac

if [ "$OS" = "linux" ] && [ "$ARCH" != "amd64" ]; then
    echo -e "${YELLOW}Warning: Only amd64 binaries are currently available for Linux. ARM64 support coming soon.${NC}"
    exit 1
fi

echo "Platform: $OS/$ARCH"
echo "Install to: $INSTALL_DIR/$BINARY_NAME"
echo ""

if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    EXISTING_VERSION=$("$INSTALL_DIR/$BINARY_NAME" version 2>/dev/null || echo "unknown")
    echo -e "${YELLOW}Existing installation found: $EXISTING_VERSION${NC}"
    echo "Upgrading..."
fi

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

if [ -f "./$BINARY_FILE" ]; then
    echo "Installing from local file..."
    cp "./$BINARY_FILE" "$TMPDIR/$BINARY_NAME"
else
    echo -e "${RED}Binary not found: ./$BINARY_FILE${NC}"
    echo ""
    echo "To install VibeSQL:"
    if [ "$OS" = "darwin" ]; then
        echo "  1. Build from source:  bash scripts/build-darwin.sh"
    else
        echo "  1. Build from source:  make build-linux"
    fi
    echo "  2. Place the binary in the current directory"
    echo "  3. Run this script again"
    exit 1
fi

chmod +x "$TMPDIR/$BINARY_NAME"

"$TMPDIR/$BINARY_NAME" version >/dev/null 2>&1 || {
    echo -e "${RED}Binary verification failed${NC}"
    exit 1
}

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMPDIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Installing to $INSTALL_DIR requires root privileges..."
    sudo mv "$TMPDIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

echo ""
echo -e "${GREEN}VibeSQL installed successfully!${NC}"
echo ""
$INSTALL_DIR/$BINARY_NAME version
echo ""
echo "Quick start:"
echo "  vibe serve        Start the server"
echo "  vibe version      Show version info"
echo "  vibe help         Show help"
echo ""
echo "API endpoint: http://127.0.0.1:5173/v1/query"
