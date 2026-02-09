# VibeSQL Local v1 - Product Requirements Document

**Version:** 1.0.0  
**Status:** Approved for Implementation  
**Author:** Requirements Analysis based on QAPert-approved specifications  
**Date:** 2026-02-05  
**Approval:** QAPert 9.5/10 (Specification Approval)

---

## Executive Summary

VibeSQL Local v1 is a minimal, self-contained PostgreSQL-based database packaged as a single 20MB binary with HTTP API support for executing SQL queries with full JSONB capabilities. This product enables developers to use PostgreSQL's powerful JSONB features in a lightweight, zero-configuration local environment with a seamless migration path to cloud deployment.

### Core Value Proposition

- **Full PostgreSQL JSONB Compatibility**: 100% JSONB feature support, not a subset
- **Zero Migration Friction**: Local PostgreSQL → Cloud PostgreSQL (same database, same queries)
- **Minimal Footprint**: 20MB binary (50% under the 40MB target)
- **Developer Trust**: "It's just Postgres" - battle-tested, predictable behavior
- **Zero Configuration**: No setup, no auth, no config files - just works

---

## Product Vision

### Long-term Strategy

**Give away VibeSQL Local as a tiny dev tool. Get the standard adopted. Sell the cloud version.**

1. **Phase 1 (v1)**: Ship free, open-source local database for developers
2. **Phase 2**: Build community adoption and developer trust
3. **Phase 3**: Launch VibeSQL Cloud with enterprise features (auth, backups, teams)
4. **Phase 4**: Monetize via cloud hosting with seamless `vibe deploy` workflow

### Target Users

**Primary**: Individual developers and small teams who need:
- Quick local database setup for prototyping
- Flexible schema with JSONB support
- SQL power without database administration overhead
- Future-proof path to production (cloud deployment)

**Secondary**: Agents/automation systems that need:
- Local data persistence (e.g., agent memory, relevance tracking)
- Offline operation
- No network dependencies

---

## Product Requirements

### 1. Core Functionality

#### 1.1 Database Engine

**REQ-1.1.1**: PostgreSQL Micro Build
- **Requirement**: Strip PostgreSQL to 20MB binary (hard limit: 25MB)
- **Includes**: Core engine, JSONB support, basic SQL, plpgsql
- **Excludes**: SSL/TLS, replication, LDAP/PAM/GSSAPI, extra CLIs (psql, pg_dump)
- **Success Criteria**: 31/31 JSONB + vsql-specific tests passing

**REQ-1.1.2**: JSONB Support
- **Requirement**: 100% PostgreSQL JSONB operator compatibility
- **Operators**: `->`, `->>`, `#>`, `#>>`, `@>`, `<@`, `?`, `?|`, `?&`
- **Functions**: `jsonb_array_length()`, `jsonb_typeof()`, `jsonb_set()` (basic set)
- **Success Criteria**: All PostgreSQL JSONB tests pass natively

**REQ-1.1.3**: Data Storage
- **Requirement**: Local filesystem storage in `./vibe-data/` directory
- **Format**: Standard PostgreSQL data directory structure
- **Portability**: Data files compatible with standard PostgreSQL (migration ready)

#### 1.2 HTTP API

**REQ-1.2.1**: Single Query Endpoint
- **Endpoint**: `POST /v1/query`
- **Request Format**: JSON with `sql` field (raw SQL string, max 10KB)
- **Response Format**: JSON with `success`, `rows`, `rowCount`, `executionTime`
- **Port**: 5173 (default, configurable)

**REQ-1.2.2**: Supported SQL Subset
- **SELECT**: Single table, WHERE, ORDER BY, LIMIT, OFFSET
- **INSERT**: Single row with optional RETURNING
- **UPDATE**: With mandatory WHERE (safety), optional RETURNING
- **DELETE**: With mandatory WHERE (safety), optional RETURNING
- **CREATE TABLE**: Basic schema with SERIAL, JSONB, PRIMARY KEY
- **DROP TABLE**: Including `IF EXISTS` support

**REQ-1.2.3**: SQL Restrictions
- **Not Supported**: JOINs, subqueries, transactions, GROUP BY/HAVING, window functions, CTEs
- **Rationale**: Minimal v1 scope; features deferred to v2

**REQ-1.2.4**: Safety Features
- **UPDATE/DELETE Without WHERE**: Rejected with error (use `WHERE 1=1` to bypass)
- **Rationale**: Prevent accidental full-table modifications

#### 1.3 Limits and Guardrails

**REQ-1.3.1**: Query Limits
- Max query length: 10,000 characters (10KB)
- Max JSON document: 1MB per row
- Max result rows: 1,000 rows
- Query timeout: 5 seconds
- Max concurrent HTTP connections: 2 (browser + CLI)

**REQ-1.3.2**: Enforcement
- Length check before query execution
- Timeout via context cancellation
- Result row limit enforced by query validator
- Connection limit enforced by HTTP server

#### 1.4 Error Handling

**REQ-1.4.1**: Error Code Mapping
- PostgreSQL SQLSTATE → VibeSQL error codes
- 10 error codes total:
  - **400**: INVALID_SQL, MISSING_REQUIRED_FIELD, UNSAFE_QUERY
  - **408**: QUERY_TIMEOUT
  - **413**: QUERY_TOO_LARGE, RESULT_TOO_LARGE, DOCUMENT_TOO_LARGE
  - **500**: INTERNAL_ERROR
  - **503**: SERVICE_UNAVAILABLE, DATABASE_UNAVAILABLE

**REQ-1.4.2**: Error Response Format
```json
{
  "success": false,
  "error": {
    "code": "INVALID_SQL",
    "message": "Human-readable error message",
    "detail": "PostgreSQL error details or additional context"
  }
}
```

#### 1.5 CLI Tool

**REQ-1.5.1**: Core Commands (v1 Scope)
- `vibe serve`: Start HTTP API server
- `vibe version`: Display version information
- `vibe help`: Show usage information

**REQ-1.5.2**: Deferred Commands (v2)
- `vibe init <project>`: Initialize database
- `vibe query "<sql>"`: Execute SQL directly
- `vibe collections`: List tables/collections
- `vibe import/export`: Data import/export

**REQ-1.5.3**: Binary Packaging
- Single executable: `vibe` (includes PostgreSQL engine + HTTP server)
- Embedded Postgres binary (no external dependencies)
- Self-initializing on first run

### 2. Platform Support

**REQ-2.1**: Supported Platforms (Priority Order)
1. Linux (x64) - Primary
2. Linux (ARM64) - Secondary
3. macOS (Apple Silicon) - Secondary
4. macOS (Intel) - Secondary
5. Windows (x64) - Tertiary (may defer to v1.1 if complex)

**REQ-2.2**: Installation Methods
- Direct binary download from GitHub releases
- `curl` installer script: `curl -fsSL https://vibesql.dev/install.sh | sh`
- Future: Homebrew (macOS), Scoop (Windows)

### 3. Performance Requirements

**REQ-3.1**: Startup Performance
- Cold start: < 2 seconds
- Database initialization (first run): < 5 seconds
- HTTP server ready: < 1 second after database ready

**REQ-3.2**: Query Performance
- Simple SELECT: < 10ms (small datasets)
- JSONB operations: Comparable to standard PostgreSQL
- Timeout enforcement: Exactly 5 seconds (no grace period)

### 4. Security Requirements

**REQ-4.1**: Local-Only Security Model
- **Authentication**: None (localhost-only, single-user desktop use)
- **Authorization**: None (full access to all data)
- **Encryption**: None (data stored unencrypted locally)
- **Network**: Binds to 127.0.0.1 only (no external access)

**REQ-4.2**: Safety Features
- UPDATE/DELETE without WHERE rejection
- Query length limits
- Timeout enforcement
- Result size limits

**REQ-4.3**: Deferred to v2 (Cloud)
- Multi-user authentication
- Role-based access control
- Encryption at rest
- TLS/SSL support
- Row-level security

### 5. Operational Requirements

**REQ-5.1**: Zero Configuration
- No config files required
- Sensible defaults baked into binary
- Optional config file support deferred to v2

**REQ-5.2**: Data Persistence
- Data directory: `./vibe-data/` (relative to working directory)
- Auto-create on first run
- Standard PostgreSQL directory structure

**REQ-5.3**: Graceful Shutdown
- SIGTERM/SIGINT handling
- Clean database shutdown
- No data corruption on kill (within reason)

**REQ-5.4**: Logging
- Minimal logging to stdout (errors and startup messages)
- Verbose mode deferred to v2

### 6. Developer Experience

**REQ-6.1**: Onboarding Flow
```bash
# Download and install
curl -fsSL https://vibesql.dev/install.sh | sh

# Start server
vibe serve
# → VibeSQL running at http://localhost:5173

# Use from any HTTP client
curl -X POST http://localhost:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT 1"}'
```

**REQ-6.2**: Documentation Requirements
- README with quick start (5-minute setup)
- API reference (endpoint, request/response formats)
- Error code reference
- JSONB operator examples
- Migration guide (local → cloud) - outline only for v1

---

## Out of Scope (v1)

### Explicitly NOT Included

**Features Deferred to v2**:
- ❌ REST endpoints (`/v1/users`, `/v1/posts`, etc.)
- ❌ Auto-schema inference from JSON
- ❌ JOINs, subqueries, transactions, GROUP BY
- ❌ Configuration file support (`vibe.conf`)
- ❌ CLI commands: `init`, `query`, `collections`, `import`, `export`
- ❌ Authentication/authorization
- ❌ Multi-tenancy
- ❌ Migration tooling (`vibe migrate`)
- ❌ Verbose logging/debugging modes

**Never Planned**:
- ❌ Auto-schema inference (requires explicit CREATE TABLE forever)
- ❌ Magic or invented semantics (pure PostgreSQL behavior)

---

## Success Criteria

### Technical Acceptance

- [ ] Binary size ≤ 20MB (hard limit: 25MB)
- [ ] 31/31 JSONB tests passing
- [ ] 200/200 integration tests passing
- [ ] All 10 error codes return correct HTTP status
- [ ] All limits enforced correctly
- [ ] Single binary builds for Linux (x64 minimum)
- [ ] Cold start < 2 seconds
- [ ] Query timeout enforced at 5s ± 100ms

### Product Acceptance

- [ ] Developer can install in one command
- [ ] Works offline (no network dependencies)
- [ ] Error messages are helpful (not cryptic PostgreSQL errors)
- [ ] API matches specification exactly (zero improvisation)
- [ ] Documentation is complete and accurate

### Quality Gates

- [ ] No breaking changes to PostgreSQL behavior
- [ ] No features beyond approved specification
- [ ] All error paths tested
- [ ] No memory leaks under normal operation
- [ ] Build process is reproducible (Docker-based)

---

## Non-Functional Requirements

### NFR-1: Maintainability
- Build process fully automated (Dockerfile)
- CI/CD pipeline for binary builds
- Automated test suite (200+ tests)
- Version-controlled build configuration

### NFR-2: Portability
- Minimal external dependencies
- Statically linked binary (where possible)
- Self-contained (no separate Postgres install needed)

### NFR-3: Upgradeability
- PostgreSQL version upgrades via Dockerfile ARG
- Automated test suite for version validation
- Backward compatibility for data files

### NFR-4: Observability
- Startup/shutdown logs
- Query execution time tracking
- Error logging with context

---

## Assumptions and Constraints

### Assumptions

1. **Target users have basic SQL knowledge**: No SQL tutorial in v1
2. **Desktop use case**: Not optimized for server deployments
3. **Small datasets**: Optimized for < 100MB databases
4. **Offline-first**: No reliance on network services
5. **Trust model**: Single-user, localhost-only, no malicious queries expected

### Constraints

1. **Binary size hard limit**: 25MB (20MB target)
2. **PostgreSQL version**: 16.x (current stable)
3. **No breaking changes**: Must preserve PostgreSQL JSONB semantics
4. **QAPert authorization required**: For any spec changes
5. **Implementation language**: TBD (optimize for binary size - likely Go, Rust, or C)

### Dependencies

1. **PostgreSQL source**: Official PostgreSQL 16.x tarball
2. **Build tools**: Docker (for reproducible builds)
3. **Test framework**: To be determined based on implementation language
4. **HTTP library**: Standard library preferred (minimal dependencies)

---

## Risks and Mitigations

### Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Binary size exceeds 20MB target | High | Low | Aggressive stripping, CI size checks (25MB hard limit) |
| Windows build complexity | Medium | Medium | Ship Linux/macOS first, Windows as v1.1 |
| Performance issues with JSONB | Medium | Low | Leverage PostgreSQL's proven implementation |
| Build reproducibility issues | High | Low | Docker-based builds, version pinning |

### Product Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Low developer adoption | High | Medium | Strong documentation, example projects, marketing |
| Competition from PocketBase/Supabase | Medium | High | Emphasize PostgreSQL compatibility and cloud path |
| Feature creep during implementation | Medium | Medium | Strict adherence to spec, QAPert review gates |

---

## Implementation Phases

### Phase 1: PostgreSQL Micro Build (Week 1)
**Objective**: Build stripped Postgres binary at 20MB

**Deliverables**:
- Dockerfile for reproducible builds
- `vsql_postgres_micro` binary (≤ 20MB)
- 31/31 JSONB tests passing
- Build documentation

### Phase 2: HTTP API Wrapper (Week 2)
**Objective**: Build HTTP server that wraps Postgres

**Deliverables**:
- HTTP server (port 5173)
- Query validator (length, timeout, safety checks)
- Error handling (10 error codes)
- PostgreSQL integration (connection pool, SQLSTATE mapping)

### Phase 3: Testing (Week 3)
**Objective**: Pass all 200 tests defined in specification

**Deliverables**:
- Unit tests (~50): HTTP parsing, validation, response formatting
- Integration tests (~100): SQL subset, JSONB operators, WHERE clauses
- Error handling tests (~30): All error codes, limit enforcement
- Load tests (~10): Concurrent queries, large results, timeouts
- E2E tests (~10): Full workflows (CREATE → INSERT → SELECT → UPDATE → DELETE)

### Phase 4: Packaging & CLI (Week 4)
**Objective**: Single binary distribution

**Deliverables**:
- `vibe` binary (includes Postgres engine + HTTP server)
- Platform builds (Linux x64 minimum, ARM64/macOS/Windows if time permits)
- CLI commands: `serve`, `version`, `help`
- Installation script (`install.sh`)
- Documentation (README, API reference, examples)

---

## Appendix A: Recommended Schema Pattern

The canonical VibeSQL Local table structure:

```sql
CREATE TABLE collection_name (
  id SERIAL PRIMARY KEY,
  data JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);
```

- **id**: Auto-incrementing primary key
- **data**: JSONB document (flexible schema)
- **created_at**: Automatic timestamp (optional but recommended)

---

## Appendix B: Example API Calls

### Create Table
```bash
curl -X POST http://localhost:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "CREATE TABLE users (id SERIAL PRIMARY KEY, data JSONB NOT NULL)"}'
```

### Insert with RETURNING
```bash
curl -X POST http://localhost:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "INSERT INTO users (data) VALUES ('\''{"name": "Alice", "email": "alice@example.com"}'\'') RETURNING *"}'
```

### Query with JSONB Operators
```bash
curl -X POST http://localhost:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT * FROM users WHERE data->>'\''name'\'' = '\''Alice'\''"}'
```

### Update with RETURNING
```bash
curl -X POST http://localhost:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "UPDATE users SET data = '\''{"name": "Alice Updated", "email": "alice@example.com"}'\'' WHERE id = 1 RETURNING *"}'
```

### Delete with Safety Bypass
```bash
curl -X POST http://localhost:5173/v1/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "DELETE FROM users WHERE 1=1 RETURNING *"}'
```

---

## Appendix C: Decision Rationale

### Why PostgreSQL Micro Over SQLite?

**Agent team exploration findings (3 minutes, 457K tokens, 4 agents)**:

1. **Binary size differential**: 20MB vs 10-15MB (only 2x, not 4-10x as expected)
2. **JSONB coverage**: 100% vs ~80% (SQLite JSON1 has limitations)
3. **Cloud migration**: Zero friction (same database) vs moderate friction (query rewrites)
4. **Developer trust**: "It's just Postgres" vs "SQLite pretending to be Postgres"

**Conclusion**: 5-10MB premium buys full JSONB power + seamless cloud path + developer confidence = obvious choice.

### Why DuckDB Was Ruled Out

- Different SQL dialect (migration complexity)
- Not as battle-tested for OLTP workloads
- Cloud deployment story unclear
- No significant advantage over PostgreSQL for this use case

---

**Document Status**: Ready for Technical Specification phase  
**Next Step**: Create `spec.md` with implementation details  
**Approval Authority**: QAPert (for any requirements changes)
