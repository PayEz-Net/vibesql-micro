# Performance and Load Tests

## Overview

These tests verify VibeSQL meets performance requirements and remains stable under load.

## Prerequisites

1. **PostgreSQL** running on localhost:5432
2. **VibeSQL server** running on localhost:5173

```bash
# Start PostgreSQL
docker run -d -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres

# Start VibeSQL
./vibe serve
```

## Running Tests

### Run All Benchmarks
```bash
go test ./Tests/performance/... -bench=. -benchmem
```

### Run Specific Benchmark
```bash
go test ./Tests/performance/... -bench=BenchmarkSimpleSelect -benchmem
```

### Run Load Tests
```bash
go test ./Tests/performance/... -v -run TestLoad
```

### Run All Tests (Including Load Tests)
```bash
go test ./Tests/performance/... -v
```

### Skip Long-Running Tests
```bash
go test ./Tests/performance/... -v -short
```

## Benchmarks

### 1. BenchmarkSimpleSelect
- **Purpose**: Measures performance of simple `SELECT 1` query
- **Target**: <10ms per operation
- **Verifies**: Basic query execution overhead

### 2. BenchmarkSelectWithWhere
- **Purpose**: Measures SELECT with WHERE clause on indexed data
- **Target**: <20ms per operation
- **Verifies**: Query planning and filtering performance

### 3. BenchmarkJSONBFieldAccess
- **Purpose**: Measures JSONB field access with `->` operator
- **Target**: <30ms per operation
- **Verifies**: JSONB operator performance

### 4. BenchmarkJSONBTextExtraction
- **Purpose**: Measures JSONB text extraction with `->>` and WHERE
- **Target**: <30ms per operation
- **Verifies**: JSONB text operator and filtering performance

### 5. BenchmarkInsert
- **Purpose**: Measures INSERT operation performance
- **Target**: <20ms per operation
- **Verifies**: Write operation overhead

### 6. BenchmarkUpdate
- **Purpose**: Measures UPDATE operation performance
- **Target**: <20ms per operation
- **Verifies**: Update operation overhead

## Load Tests

### TestLoadSequential
- **Purpose**: Executes 100 sequential queries
- **Target**: Average <50ms per query
- **Duration**: ~5-10 seconds
- **Metrics**:
  - Total queries
  - Total time
  - Average time per query
  - Queries per second

### TestLoadConcurrent
- **Purpose**: Executes 100 queries across 2 concurrent workers
- **Target**: No errors, stable throughput
- **Duration**: ~5-10 seconds
- **Metrics**:
  - Concurrency level
  - Total queries
  - Total time
  - Throughput (queries/sec)

### TestMemoryUsage
- **Purpose**: Monitors memory usage during 1000 queries
- **Target**: <10MB memory growth
- **Verifies**: No memory leaks
- **Metrics**:
  - Current alloc delta
  - Total alloc delta
  - Allocations per query

### TestQueryTimeout
- **Purpose**: Verifies query timeout enforcement
- **Target**: 5s ± 100ms
- **Verifies**: Timeout mechanism accuracy
- **Note**: Takes ~5 seconds to run (not skipped by `-short`)

### TestStartupTime
- **Purpose**: Documents manual startup time measurement
- **Target**: <2 seconds cold start
- **Status**: Manual test (automated in Phase 4)
- **Note**: Currently skipped

## Performance Targets

| Metric | Target | Test |
|--------|--------|------|
| **Cold start** | <2 seconds | Manual (Phase 4) |
| **Database init** | <5 seconds | Manual (first run) |
| **HTTP server ready** | <1 second | After DB ready |
| **Simple SELECT** | <10ms | BenchmarkSimpleSelect |
| **Query timeout** | 5s ± 100ms | TestQueryTimeout |
| **Memory leak** | <10MB growth | TestMemoryUsage |
| **Sequential load** | 100 queries, <50ms avg | TestLoadSequential |
| **Concurrent load** | 2 workers, stable | TestLoadConcurrent |

## Example Output

### Benchmark Results
```
BenchmarkSimpleSelect-8                5000    234567 ns/op    1024 B/op    12 allocs/op
BenchmarkSelectWithWhere-8             3000    345678 ns/op    2048 B/op    18 allocs/op
BenchmarkJSONBFieldAccess-8            2000    567890 ns/op    3072 B/op    24 allocs/op
BenchmarkJSONBTextExtraction-8         2000    678901 ns/op    3584 B/op    28 allocs/op
BenchmarkInsert-8                      4000    289012 ns/op    1536 B/op    15 allocs/op
BenchmarkUpdate-8                      3500    312345 ns/op    1792 B/op    16 allocs/op
```

### Load Test Results
```
=== RUN   TestLoadSequential
    Sequential load test completed:
      Total queries: 100
      Total time: 3.456s
      Average time per query: 34.56ms
      Queries per second: 28.94
--- PASS: TestLoadSequential (3.46s)

=== RUN   TestLoadConcurrent
    Concurrent load test completed:
      Concurrency: 2 workers
      Total queries: 100
      Total time: 2.123s
      Average time per query: 21.23ms
      Throughput: 47.11 queries/sec
--- PASS: TestLoadConcurrent (2.12s)
```

## Interpreting Results

### Good Performance Indicators
- ✅ BenchmarkSimpleSelect: <10ms (100,000 ns/op)
- ✅ Low allocation counts (<20 allocs/op)
- ✅ Sequential load: >20 queries/second
- ✅ Concurrent load: Higher throughput than sequential
- ✅ Memory usage: <10MB growth over 1000 queries
- ✅ Query timeout: 5.0s ± 0.1s

### Warning Signs
- ⚠️ Simple SELECT >20ms: Possible network/server overhead
- ⚠️ High allocations: May indicate unnecessary copying
- ⚠️ Low queries/second: Server may be overloaded
- ⚠️ Memory growth >10MB: Potential memory leak
- ⚠️ Timeout variance >200ms: Inconsistent timeout enforcement

## Troubleshooting

### Benchmarks Very Slow (>100ms per operation)
- **Check**: Is server running locally? Network latency adds overhead
- **Check**: Is PostgreSQL overloaded?
- **Solution**: Restart server, reduce background load

### High Memory Usage
- **Check**: Are tables cleaned up after tests?
- **Check**: Is PostgreSQL caching results?
- **Solution**: Restart PostgreSQL, run tests individually

### Timeout Test Fails
- **Check**: Is server under heavy load?
- **Check**: Are other queries running?
- **Solution**: Run timeout test in isolation

### Load Tests Time Out
- **Check**: Connection limit (max 2 concurrent)
- **Check**: Query backlog
- **Solution**: Reduce concurrency, increase timeout

## CI/CD Integration

These tests are designed to run in CI/CD with the following characteristics:

- **Fast benchmarks**: Run on every commit (use `-short` to skip load tests)
- **Load tests**: Run nightly or pre-release
- **Baseline comparison**: Compare benchmark results to detect regressions

### Example CI Configuration
```yaml
# Fast benchmarks (on every commit)
- name: Run performance benchmarks
  run: go test ./Tests/performance/... -bench=. -benchmem -short

# Full performance suite (nightly)
- name: Run full performance tests
  run: go test ./Tests/performance/... -v -bench=. -benchmem
```

## Future Enhancements (Phase 4)

1. **Automated startup time measurement** with embedded binary
2. **Performance regression detection** (baseline comparison)
3. **Grafana dashboard** integration for metrics
4. **Stress testing** with higher concurrency limits
5. **Long-running stability tests** (hours/days)

## Notes

- Performance results vary by hardware, OS, and system load
- Run benchmarks multiple times for consistency
- Use dedicated test environment for reliable measurements
- Baseline results on reference hardware for comparison
