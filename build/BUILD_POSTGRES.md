# Building Minimal PostgreSQL for VibeSQL

## Overview

This guide explains how to build a minimal PostgreSQL binary (<20MB) for embedding in VibeSQL.

**Target**: Linux x64 (amd64)  
**Size Goal**: ≤20MB (target: ~15MB after stripping)

---

## Prerequisites

### Option 1: Docker (Recommended for Windows/macOS)

**Easiest approach** - Works on all platforms:

```bash
# Windows
build_postgres.cmd

# Linux/macOS
./build_postgres_docker.sh
```

**Requirements**:
- Docker Desktop installed and running
- ~2GB disk space for build

### Option 2: Native Linux Build

You need a Linux x64 environment with development tools:

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y build-essential wget tar

# RHEL/CentOS/Fedora
sudo yum groupinstall "Development Tools"
sudo yum install wget tar
```

**Required Tools**:
- **GCC**: C compiler
- **Make**: Build automation
- **wget/curl**: Download source
- **tar**: Extract archives
- **strip**: Binary optimization (usually included in binutils)

---

## Step 1: Download PostgreSQL Source

```bash
cd build

# Download PostgreSQL 16.1 (latest stable as of Dec 2023)
POSTGRES_VERSION=16.1
wget https://ftp.postgresql.org/pub/source/v${POSTGRES_VERSION}/postgresql-${POSTGRES_VERSION}.tar.gz

# Extract
tar -xzf postgresql-${POSTGRES_VERSION}.tar.gz
cd postgresql-${POSTGRES_VERSION}
```

**Size**: ~30MB compressed, ~180MB extracted

---

## Step 2: Configure Minimal Build

Configure PostgreSQL with minimal features to reduce binary size:

```bash
./configure \
  --prefix=/tmp/postgres_micro \
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
  --disable-spinlocks \
  --disable-thread-safety
```

### Configuration Flags Explained

| Flag | Purpose | Size Impact |
|------|---------|-------------|
| `--without-readline` | No interactive CLI editing | -500KB |
| `--without-zlib` | No compression support | -100KB |
| `--without-openssl` | No SSL/TLS (localhost only) | -2MB |
| `--without-gssapi` | No Kerberos auth | -200KB |
| `--without-ldap` | No LDAP auth | -100KB |
| `--without-pam` | No PAM auth | -50KB |
| `--without-systemd` | No systemd integration | -50KB |
| `--without-libxml` | No XML functions | -1MB |
| `--without-libxslt` | No XSLT functions | -500KB |
| `--without-perl` | No PL/Perl | -1MB |
| `--without-python` | No PL/Python | -2MB |
| `--without-tcl` | No PL/Tcl | -500KB |
| `--without-icu` | No international support | -10MB |
| `--disable-rpath` | No runtime library paths | Smaller |
| `--disable-spinlocks` | Use semaphores (slower, smaller) | -100KB |
| `--disable-thread-safety` | No libpq threading | -50KB |

**Estimated Savings**: ~18MB+ from configuration alone

---

## Step 3: Build Only Required Components

Build only the core PostgreSQL server (no client tools):

```bash
# Build postgres server only (not psql, pg_dump, etc.)
cd src/backend
make -j$(nproc)

# Build essential utilities
cd ../bin/initdb
make -j$(nproc)

cd ../pg_ctl
make -j$(nproc)
```

**What We're Building**:
- `postgres` - Database server (core engine)
- `initdb` - Database initialization utility
- `pg_ctl` - Server control utility

**What We're Skipping**:
- `psql` - Interactive terminal
- `pg_dump` / `pg_restore` - Backup tools
- `createdb` / `dropdb` - Database management
- All other client utilities

**Size Savings**: ~5MB by skipping unused tools

---

## Step 4: Strip Symbols and Debug Info

Strip debugging symbols to minimize binary size:

```bash
# Go to build output directory
cd /tmp/postgres_micro/bin

# Strip all binaries
strip --strip-all postgres
strip --strip-all initdb  
strip --strip-all pg_ctl

# Check sizes
ls -lh postgres initdb pg_ctl
```

**Expected Sizes** (after stripping):
- `postgres`: ~12-15MB
- `initdb`: ~1-2MB
- `pg_ctl`: ~500KB

**Total**: ~14-18MB (within 20MB target)

---

## Step 5: Create Single Binary Package

For VibeSQL, we only need the `postgres` binary:

```bash
# Copy just the postgres binary
cp /tmp/postgres_micro/bin/postgres ./postgres_micro_linux_amd64

# Verify size
ls -lh postgres_micro_linux_amd64

# Verify it's executable
chmod +x postgres_micro_linux_amd64
```

---

## Step 6: Verify Functionality

Test that the minimal PostgreSQL binary works:

```bash
# Create test data directory
mkdir -p /tmp/pg_test_data

# Initialize database (VibeSQL will do this programmatically)
./postgres_micro_linux_amd64 --single -D /tmp/pg_test_data template1 <<EOF
CREATE DATABASE test;
EOF

# Start postgres in background
./postgres_micro_linux_amd64 -D /tmp/pg_test_data -p 5555 &

# Wait for startup
sleep 2

# Test connection (requires psql or libpq client)
# For VibeSQL, we'll connect via Go's lib/pq driver

# Stop postgres
pkill postgres_micro
```

---

## Automated Build Script

See `build/build_postgres.sh` for automated build process:

```bash
cd build
chmod +x build_postgres.sh
./build_postgres.sh
```

This script automates all steps above and produces `postgres_micro_linux_amd64`.

---

## Troubleshooting

### Build Fails with Missing Dependencies

**Error**: `configure: error: readline library not found`

**Solution**: You likely need to install development headers:
```bash
# Ubuntu/Debian
sudo apt-get install libreadline-dev

# Or disable readline (we already do this)
./configure --without-readline ...
```

### Binary Size Too Large (>20MB)

**Solutions**:
1. Ensure you stripped symbols: `strip --strip-all postgres`
2. Check if debug build: Add `CFLAGS="-O2"` to configure
3. Verify minimal configure flags are used

### Binary Won't Execute

**Error**: `./postgres: cannot execute binary file`

**Cause**: Wrong architecture or corrupted binary

**Solution**:
```bash
# Check binary type
file postgres_micro_linux_amd64
# Should show: ELF 64-bit LSB executable, x86-64

# Check if stripped
ls -lh postgres_micro_linux_amd64
# Should be 12-15MB
```

### Postgres Crashes on Startup

**Error**: `FATAL: could not create shared memory segment`

**Cause**: Insufficient shared memory limits

**Solution**: VibeSQL uses minimal shared memory settings, but if testing manually:
```bash
# Increase system limits
sudo sysctl -w kernel.shmmax=17179869184
sudo sysctl -w kernel.shmall=4194304
```

---

## Next Steps

After building `postgres_micro_linux_amd64`:

1. ✅ **Copy to embed directory**: `internal/postgres/embed/postgres_micro_linux_amd64`
2. ✅ **Test embedding**: Update `manager.go` to use embedded binary
3. ✅ **Test extraction**: Verify VibeSQL can extract and run it
4. ✅ **Integration test**: Run full test suite with embedded postgres

---

## Size Optimization Tips

### Further Reductions (if needed)

If binary is still too large, try:

1. **Use `strip -s` instead of `--strip-all`**: More aggressive stripping
2. **Compile with `-Os`**: Optimize for size instead of speed
3. **Use UPX compression**: Compress the binary (adds runtime overhead)
4. **Remove more features**: Disable JSONB, triggers, foreign keys

### Size vs Functionality Tradeoff

| Feature | Size Impact | VibeSQL Needs? |
|---------|-------------|----------------|
| JSONB support | ~2MB | ✅ YES (required) |
| Triggers | ~500KB | ❌ NO |
| Foreign keys | ~200KB | ❌ NO |
| Full text search | ~1MB | ❌ NO |
| PostGIS/extensions | ~10MB | ❌ NO |

**Current Config**: Keeps JSONB, removes everything else

---

## Security Considerations

### Disabled Features Impact

1. **No SSL/TLS** (`--without-openssl`)
   - Impact: Unencrypted connections
   - Mitigation: VibeSQL binds to localhost only (127.0.0.1)
   - Risk: Low (no network exposure)

2. **No Authentication Libs** (`--without-pam`, `--without-ldap`)
   - Impact: No external auth
   - Mitigation: VibeSQL uses local trust authentication
   - Risk: Low (embedded, localhost only)

3. **No Compression** (`--without-zlib`)
   - Impact: No TOAST compression
   - Mitigation: VibeSQL handles small documents
   - Risk: Low (10KB query limit)

### Recommendations

- ✅ Use VibeSQL on localhost only (already enforced)
- ✅ Don't expose VibeSQL port to network
- ✅ Run in container for isolation (optional)

---

## References

- PostgreSQL Build Docs: https://www.postgresql.org/docs/16/installation.html
- Configure Options: https://www.postgresql.org/docs/16/install-procedure.html
- Size Optimization: https://wiki.postgresql.org/wiki/Minimizing_PostgreSQL_footprint

---

## Summary

**Target**: Linux x64 PostgreSQL binary ≤20MB  
**Expected Size**: ~14-16MB  
**Build Time**: ~10-15 minutes  
**Disk Space**: ~500MB for source + build artifacts

**Next**: Run `build/build_postgres.sh` to create `postgres_micro_linux_amd64`
