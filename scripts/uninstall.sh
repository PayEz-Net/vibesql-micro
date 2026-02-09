#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${VIBESQL_INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="vibe"
DATA_DIR="./vibe-data"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}VibeSQL Uninstaller${NC}"
echo "================================"

if [ ! -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo "VibeSQL is not installed at $INSTALL_DIR/$BINARY_NAME"
    exit 0
fi

echo "This will remove:"
echo "  - $INSTALL_DIR/$BINARY_NAME"

if [ -d "$DATA_DIR" ]; then
    echo -e "  - ${RED}$DATA_DIR (DATABASE DATA)${NC}"
    echo ""
    read -p "Remove database data too? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        REMOVE_DATA=true
    else
        REMOVE_DATA=false
    fi
else
    REMOVE_DATA=false
fi

echo ""
read -p "Proceed with uninstall? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

if [ -w "$INSTALL_DIR" ]; then
    rm -f "$INSTALL_DIR/$BINARY_NAME"
else
    sudo rm -f "$INSTALL_DIR/$BINARY_NAME"
fi

if [ "$REMOVE_DATA" = true ] && [ -d "$DATA_DIR" ]; then
    rm -rf "$DATA_DIR"
    echo -e "${GREEN}Database data removed${NC}"
fi

echo -e "${GREEN}VibeSQL uninstalled successfully${NC}"
