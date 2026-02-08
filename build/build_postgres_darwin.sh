#!/bin/bash
set -euo pipefail

POSTGRES_VERSION="16.1"
POSTGRES_URL="https://ftp.postgresql.org/pub/source/v${POSTGRES_VERSION}/postgresql-${POSTGRES_VERSION}.tar.gz"
BUILD_DIR="$(pwd)/postgres-build-darwin"
INSTALL_PREFIX="/tmp/postgres_micro_install_darwin"

ARCH="${GOARCH:-$(uname -m)}"
case "$ARCH" in
    x86_64|amd64) ARCH="amd64"; ARCH_NATIVE="x86_64" ;;
    aarch64|arm64) ARCH="arm64"; ARCH_NATIVE="arm64" ;;
    *)
        echo "ERROR: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

check_requirements() {
    log_info "Checking macOS build requirements..."

    local missing_tools=()
    for tool in gcc make wget tar strip; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done

    if [ ${#missing_tools[@]} -gt 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        log_error "Install Xcode Command Line Tools: xcode-select --install"
        log_error "Install wget: brew install wget"
        exit 1
    fi

    log_info "All requirements satisfied"
}

download_postgres() {
    log_info "Downloading PostgreSQL ${POSTGRES_VERSION}..."

    local tarball="postgresql-${POSTGRES_VERSION}.tar.gz"
    if [ -f "$tarball" ]; then
        log_warn "Tarball already exists, skipping download"
    else
        wget "$POSTGRES_URL" -O "$tarball"
    fi

    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    tar -xzf "$tarball" -C "$BUILD_DIR" --strip-components=1
    log_info "Extraction complete"
}

configure_postgres() {
    log_info "Configuring minimal PostgreSQL build for macOS ${ARCH}..."

    cd "$BUILD_DIR"

    ./configure \
        --prefix="$INSTALL_PREFIX" \
        --without-readline \
        --without-zlib \
        --without-openssl \
        --without-gssapi \
        --without-ldap \
        --without-pam \
        --without-systemd \
        --without-libxml \
        --without-libxslt \
        --without-perl \
        --without-python \
        --without-tcl \
        --without-icu \
        --disable-rpath \
        CFLAGS="-O2 -Wno-unguarded-availability-new" \
        > configure.log 2>&1

    log_info "Configuration complete"
}

build_postgres() {
    log_info "Building PostgreSQL (this may take 10-15 minutes)..."

    cd "$BUILD_DIR"

    local num_cores=$(sysctl -n hw.ncpu 2>/dev/null || echo 4)
    log_info "Using $num_cores parallel jobs"

    make -j"$num_cores" > build.log 2>&1
    log_info "Build complete"

    make install > install.log 2>&1
    log_info "Install complete"
    
    # Remove problematic extensions that cause symbol resolution issues
    log_info "Removing problematic extensions..."
    rm -f "$INSTALL_PREFIX/lib/dict_snowball.dylib"
    rm -f "$INSTALL_PREFIX/lib/plpgsql.dylib"
    rm -f "$INSTALL_PREFIX/lib"/*_and_*.dylib
    rm -f "$INSTALL_PREFIX/share/snowball_create.sql"
    
    # Patch postgres.bki to skip plpgsql initialization
    sed -i '' '/CREATE EXTENSION plpgsql/d' "$INSTALL_PREFIX/share/postgres.bki" 2>/dev/null || true
    
    log_info "Extensions cleaned up"
}

strip_binary() {
    log_info "Stripping debug symbols..."

    local postgres_bin="$INSTALL_PREFIX/bin/postgres"
    if [ ! -f "$postgres_bin" ]; then
        log_error "Postgres binary not found at $postgres_bin"
        exit 1
    fi

    local size_before=$(stat -f%z "$postgres_bin")
    local size_mb_before=$(echo "scale=2; $size_before / 1024 / 1024" | bc)
    log_info "Size before stripping: ${size_mb_before}MB"

    strip "$postgres_bin" 2>/dev/null || strip -x "$postgres_bin"

    local size_after=$(stat -f%z "$postgres_bin")
    local size_mb_after=$(echo "scale=2; $size_after / 1024 / 1024" | bc)
    log_info "Size after stripping: ${size_mb_after}MB"
}

package_binaries() {
    log_info "Packaging binaries for embed..."

    local SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
    local EMBED_DIR="$PROJECT_ROOT/internal/postgres/embed"

    mkdir -p "$EMBED_DIR"

    cp "$INSTALL_PREFIX/bin/postgres" "$EMBED_DIR/postgres_micro_darwin_${ARCH}"
    chmod +x "$EMBED_DIR/postgres_micro_darwin_${ARCH}"

    cp "$INSTALL_PREFIX/bin/initdb" "$EMBED_DIR/initdb_darwin_${ARCH}"
    chmod +x "$EMBED_DIR/initdb_darwin_${ARCH}"

    cp "$INSTALL_PREFIX/bin/pg_ctl" "$EMBED_DIR/pg_ctl_darwin_${ARCH}"
    chmod +x "$EMBED_DIR/pg_ctl_darwin_${ARCH}"

    local LIBPQ="$INSTALL_PREFIX/lib/libpq.5.dylib"
    if [ -f "$LIBPQ" ]; then
        cp "$LIBPQ" "$EMBED_DIR/libpq.5.dylib"
    elif [ -f "$INSTALL_PREFIX/lib/libpq.dylib" ]; then
        cp "$INSTALL_PREFIX/lib/libpq.dylib" "$EMBED_DIR/libpq.5.dylib"
    else
        log_warn "libpq dylib not found, initdb may fail at runtime"
    fi

    if [ ! -f "$EMBED_DIR/share.tar.gz" ]; then
        log_info "Creating share.tar.gz..."
        cd "$INSTALL_PREFIX"
        tar -czf "$EMBED_DIR/share.tar.gz" share/
    else
        log_info "share.tar.gz already exists, skipping"
    fi

    echo ""
    echo "========================================="
    echo "macOS PostgreSQL Build Complete ($ARCH)"
    echo "========================================="
    echo "Binaries placed in: $EMBED_DIR"
    ls -lh "$EMBED_DIR/"*darwin* 2>/dev/null
    ls -lh "$EMBED_DIR/libpq.5.dylib" 2>/dev/null
    echo ""
}

cleanup() {
    log_info "Cleaning up build artifacts..."
    rm -rf "$BUILD_DIR" "$INSTALL_PREFIX"
    log_info "Cleanup complete"
}

main() {
    echo "========================================="
    echo "VibeSQL PostgreSQL Build (macOS $ARCH)"
    echo "========================================="
    echo "Version: PostgreSQL ${POSTGRES_VERSION}"
    echo "Target: macOS ${ARCH}"
    echo "Size Goal: ≤20MB"
    echo ""

    check_requirements
    download_postgres
    configure_postgres
    build_postgres
    strip_binary
    package_binaries

    if [ "${CI:-}" = "true" ] || [ "${NONINTERACTIVE:-}" = "true" ]; then
        cleanup
    else
        read -p "Remove build artifacts to save space? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            cleanup
        fi
    fi

    log_info "Build process complete!"
}

main
