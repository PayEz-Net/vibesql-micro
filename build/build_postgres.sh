#!/bin/bash
#
# Build Minimal PostgreSQL for VibeSQL Embedding
# Target: Linux x64 (amd64)
# Size Goal: ≤20MB
#

set -euo pipefail

# Configuration
POSTGRES_VERSION="16.1"
POSTGRES_URL="https://ftp.postgresql.org/pub/source/v${POSTGRES_VERSION}/postgresql-${POSTGRES_VERSION}.tar.gz"
BUILD_DIR="$(pwd)/postgres-build"
INSTALL_PREFIX="/tmp/postgres_micro_install"
OUTPUT_BINARY="postgres_micro_linux_amd64"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_requirements() {
    log_info "Checking build requirements..."
    
    local missing_tools=()
    
    for tool in gcc make wget tar strip; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        log_error "Install with:"
        log_error "  Ubuntu/Debian: sudo apt-get install build-essential wget tar"
        log_error "  RHEL/CentOS:   sudo yum groupinstall 'Development Tools' && sudo yum install wget tar"
        exit 1
    fi
    
    log_info "All requirements satisfied ✓"
}

download_postgres() {
    log_info "Downloading PostgreSQL ${POSTGRES_VERSION}..."
    
    local tarball="postgresql-${POSTGRES_VERSION}.tar.gz"
    
    if [ -f "$tarball" ]; then
        log_warn "Tarball already exists, skipping download"
    else
        wget "$POSTGRES_URL" -O "$tarball"
        log_info "Download complete ✓"
    fi
    
    # Extract
    log_info "Extracting source..."
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    tar -xzf "$tarball" -C "$BUILD_DIR" --strip-components=1
    log_info "Extraction complete ✓"
}

configure_postgres() {
    log_info "Configuring minimal PostgreSQL build..."
    
    cd "$BUILD_DIR"
    
    # Configure with minimal features
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
        CFLAGS="-O2" \
        > configure.log 2>&1
    
    if [ $? -eq 0 ]; then
        log_info "Configuration complete ✓"
    else
        log_error "Configuration failed! Check configure.log"
        tail -n 20 configure.log
        exit 1
    fi
}

build_postgres() {
    log_info "Building PostgreSQL (this may take 10-15 minutes)..."
    
    cd "$BUILD_DIR"
    
    # Build in parallel
    local num_cores=$(nproc 2>/dev/null || echo 4)
    log_info "Using $num_cores parallel jobs"
    
    # Build server backend
    make -j"$num_cores" > build.log 2>&1
    
    if [ $? -eq 0 ]; then
        log_info "Build complete ✓"
    else
        log_error "Build failed! Check build.log"
        tail -n 20 build.log
        exit 1
    fi
    
    # Install to temporary location
    log_info "Installing to temporary location..."
    make install > install.log 2>&1
    
    log_info "Install complete ✓"
}

strip_binary() {
    log_info "Stripping debug symbols..."
    
    local postgres_bin="$INSTALL_PREFIX/bin/postgres"
    
    if [ ! -f "$postgres_bin" ]; then
        log_error "Postgres binary not found at $postgres_bin"
        exit 1
    fi
    
    # Show size before stripping
    local size_before=$(stat -f%z "$postgres_bin" 2>/dev/null || stat -c%s "$postgres_bin")
    local size_mb_before=$(echo "scale=2; $size_before / 1024 / 1024" | bc)
    log_info "Size before stripping: ${size_mb_before}MB"
    
    # Strip symbols
    strip --strip-all "$postgres_bin"
    
    # Show size after stripping
    local size_after=$(stat -f%z "$postgres_bin" 2>/dev/null || stat -c%s "$postgres_bin")
    local size_mb_after=$(echo "scale=2; $size_after / 1024 / 1024" | bc)
    log_info "Size after stripping: ${size_mb_after}MB"
    
    local savings=$(echo "scale=2; $size_mb_before - $size_mb_after" | bc)
    log_info "Saved: ${savings}MB"
    
    # Check if within target
    if (( $(echo "$size_after > 20971520" | bc -l) )); then
        log_warn "Binary size ${size_mb_after}MB exceeds 20MB target!"
        log_warn "Consider additional optimizations"
    else
        log_info "Binary size within 20MB target ✓"
    fi
}

package_binary() {
    log_info "Packaging binary..."
    
    local postgres_bin="$INSTALL_PREFIX/bin/postgres"
    local output_path="$(dirname $BUILD_DIR)/$OUTPUT_BINARY"
    
    # Copy binary
    cp "$postgres_bin" "$output_path"
    chmod +x "$output_path"
    
    log_info "Binary created: $output_path"
    
    # Show final info
    local size=$(stat -f%z "$output_path" 2>/dev/null || stat -c%s "$output_path")
    local size_mb=$(echo "scale=2; $size / 1024 / 1024" | bc)
    
    echo ""
    echo "========================================="
    echo "PostgreSQL Minimal Binary Build Complete"
    echo "========================================="
    echo "Output: $output_path"
    echo "Size: ${size_mb}MB"
    echo "Architecture: x86_64 (Linux)"
    echo ""
    
    # Verify binary
    file "$output_path"
}

verify_binary() {
    log_info "Verifying binary..."
    
    local output_path="$(dirname $BUILD_DIR)/$OUTPUT_BINARY"
    
    # Check it's a valid ELF binary
    if file "$output_path" | grep -q "ELF 64-bit"; then
        log_info "Binary type verified ✓"
    else
        log_error "Binary is not a valid ELF 64-bit executable"
        file "$output_path"
        exit 1
    fi
    
    # Check if executable
    if [ -x "$output_path" ]; then
        log_info "Binary is executable ✓"
    else
        log_error "Binary is not executable"
        exit 1
    fi
    
    # Test help output
    log_info "Testing binary execution..."
    if "$output_path" --version 2>&1 | grep -q "postgres"; then
        log_info "Binary execution verified ✓"
    else
        log_warn "Could not verify postgres version (may need dependencies)"
    fi
}

cleanup() {
    log_info "Cleaning up build artifacts..."
    
    if [ -d "$BUILD_DIR" ]; then
        log_info "Removing build directory: $BUILD_DIR"
        rm -rf "$BUILD_DIR"
    fi
    
    if [ -d "$INSTALL_PREFIX" ]; then
        log_info "Removing install directory: $INSTALL_PREFIX"
        rm -rf "$INSTALL_PREFIX"
    fi
    
    log_info "Cleanup complete ✓"
}

show_next_steps() {
    echo ""
    echo "========================================="
    echo "Next Steps"
    echo "========================================="
    echo "1. Copy binary to embed directory:"
    echo "   mkdir -p ../internal/postgres/embed"
    echo "   cp $OUTPUT_BINARY ../internal/postgres/embed/"
    echo ""
    echo "2. Update manager.go to use embedded binary"
    echo ""
    echo "3. Test with VibeSQL:"
    echo "   cd .."
    echo "   go build -o vibe cmd/vibe/main.go"
    echo "   ./vibe serve"
    echo ""
}

# Main execution
main() {
    echo "========================================="
    echo "VibeSQL PostgreSQL Minimal Build"
    echo "========================================="
    echo "Version: PostgreSQL ${POSTGRES_VERSION}"
    echo "Target: Linux x64 (amd64)"
    echo "Size Goal: ≤20MB"
    echo ""
    
    check_requirements
    download_postgres
    configure_postgres
    build_postgres
    strip_binary
    package_binary
    verify_binary
    
    # Ask about cleanup
    read -p "Remove build artifacts to save space? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cleanup
    else
        log_info "Build artifacts kept in: $BUILD_DIR"
    fi
    
    show_next_steps
    
    log_info "Build process complete! ✓"
}

# Run main
main
