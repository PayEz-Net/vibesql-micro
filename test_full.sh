#!/bin/bash
# Full integration test suite for vibesql-micro
# Run this on Linux via ZeroClaw

set -e

VSQL="${VSQL:-./vsql-micro-linux}"
TEST_DIR="/tmp/vsql-full-test-$$"
FAILED=0
PASSED=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASSED++))
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAILED++))
}

# Setup
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

echo "=========================================="
echo "vibesql-micro Full Test Suite"
echo "Binary: $VSQL"
echo "Test Dir: $TEST_DIR"
echo "=========================================="

# Test 1: Version works
log_test "Version command"
if $VSQL version 2>&1 | grep -q "vsql-micro"; then
    log_pass "Version shows correctly"
else
    log_fail "Version command failed"
fi

# Test 2: Help works
log_test "Help command"
if $VSQL --help 2>&1 | grep -q "Usage:"; then
    log_pass "Help shows correctly"
else
    log_fail "Help command failed"
fi

# Test 3: Fresh database creation
log_test "Fresh database creation"
rm -rf test1.vsql*
OUTPUT=$($VSQL test1.vsql "SELECT 'fresh' as status" 2>&1)
if echo "$OUTPUT" | grep -q "fresh"; then
    log_pass "Fresh database created and query works"
else
    log_fail "Fresh database failed: $OUTPUT"
fi

# Test 4: Progress indicator on first run
log_test "Progress indicator shown"
rm -rf test2.vsql*
OUTPUT=$($VSQL test2.vsql "SELECT 1" 2>&1)
if echo "$OUTPUT" | grep -qi "setting up\|creating\|done"; then
    log_pass "Progress indicator shown"
else
    log_fail "No progress indicator"
fi

# Test 5: Extraction created cache
log_test "Binary extraction cache"
if [ -d "$HOME/.cache/vibesql-micro" ] || [ -d "$HOME/vibesql-micro" ]; then
    log_pass "Cache directory exists"
else
    log_fail "Cache directory not found"
fi

# Test 6: CREATE TABLE
log_test "CREATE TABLE"
rm -rf test3.vsql*
$VSQL test3.vsql "CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    log_pass "CREATE TABLE succeeded"
else
    log_fail "CREATE TABLE failed"
fi

# Test 7: INSERT
log_test "INSERT data"
$VSQL test3.vsql "INSERT INTO users (name) VALUES ('Alice'), ('Bob'), ('Charlie')" > /dev/null 2>&1
COUNT=$($VSQL test3.vsql "SELECT COUNT(*) as c FROM users" 2>&1 | grep -o '"c":[0-9]*' | cut -d: -f2)
if [ "$COUNT" = "3" ]; then
    log_pass "INSERT 3 rows succeeded"
else
    log_fail "INSERT failed or wrong count: $COUNT"
fi

# Test 8: SELECT returns JSON
log_test "SELECT returns valid JSON"
OUTPUT=$($VSQL test3.vsql "SELECT * FROM users WHERE name='Alice'" 2>&1)
if echo "$OUTPUT" | grep -q '"id":1' && echo "$OUTPUT" | grep -q '"name":"Alice"'; then
    log_pass "SELECT returns correct JSON"
else
    log_fail "SELECT JSON incorrect: $OUTPUT"
fi

# Test 9: Complex query
log_test "Complex query (JOIN, aggregate)"
$VSQL test3.vsql "CREATE TABLE orders (id SERIAL, user_id INT, amount DECIMAL)" > /dev/null 2>&1
$VSQL test3.vsql "INSERT INTO orders (user_id, amount) VALUES (1, 100.00), (1, 200.00), (2, 50.00)" > /dev/null 2>&1
OUTPUT=$($VSQL test3.vsql "SELECT u.name, SUM(o.amount) as total FROM users u JOIN orders o ON u.id = o.user_id GROUP BY u.name ORDER BY total DESC" 2>&1)
if echo "$OUTPUT" | grep -q '"total":300'; then
    log_pass "Complex query works"
else
    log_fail "Complex query failed: $OUTPUT"
fi

# Test 10: Data persistence (close and reopen)
log_test "Data persistence across reopen"
rm -rf test4.vsql*
$VSQL test4.vsql "CREATE TABLE persistent (id SERIAL, data TEXT)" > /dev/null 2>&1
$VSQL test4.vsql "INSERT INTO persistent (data) VALUES ('survive')" > /dev/null 2>&1
# Database closes here, reopen
OUTPUT=$($VSQL test4.vsql "SELECT data FROM persistent" 2>&1)
if echo "$OUTPUT" | grep -q "survive"; then
    log_pass "Data persisted after close/reopen"
else
    log_fail "Data not persisted: $OUTPUT"
fi

# Test 11: Lock detection
log_test "Lock detection (concurrent access)"
rm -rf test5.vsql*
# Create database first
$VSQL test5.vsql "SELECT 1" > /dev/null 2>&1
# Start a long-running query in background
$VSQL test5.vsql "SELECT pg_sleep(10)" &
BG_PID=$!
sleep 2
# Try to open same database - should fail with lock error
OUTPUT=$($VSQL test5.vsql "SELECT 1" 2>&1)
kill $BG_PID 2>/dev/null || true
wait $BG_PID 2>/dev/null || true
if echo "$OUTPUT" | grep -qi "busy\|lock\|in use"; then
    log_pass "Lock detection works"
else
    log_fail "No lock detected or wrong error: $OUTPUT"
fi

# Test 12: JSONB support
log_test "JSONB data type"
rm -rf test6.vsql*
$VSQL test6.vsql "CREATE TABLE json_test (id SERIAL, data JSONB)" > /dev/null 2>&1
$VSQL test6.vsql "INSERT INTO json_test (data) VALUES ('{\"key\": \"value\", \"num\": 42}')" > /dev/null 2>&1
OUTPUT=$($VSQL test6.vsql "SELECT data->>'key' as val FROM json_test" 2>&1)
if echo "$OUTPUT" | grep -q "value"; then
    log_pass "JSONB works"
else
    log_fail "JSONB failed: $OUTPUT"
fi

# Test 13: Error handling (invalid SQL)
log_test "Error handling - invalid SQL"
rm -rf test7.vsql*
OUTPUT=$($VSQL test7.vsql "INVALID SYNTAX HERE" 2>&1)
if echo "$OUTPUT" | grep -qi "error\|syntax"; then
    log_pass "Error message shown for invalid SQL"
else
    log_fail "No error for invalid SQL: $OUTPUT"
fi

# Test 14: Special characters in data
log_test "Unicode and special characters"
rm -rf test8.vsql*
$VSQL test8.vsql "CREATE TABLE special (id SERIAL, data TEXT)" > /dev/null 2>&1
$VSQL test8.vsql "INSERT INTO special (data) VALUES ('日本語'), ('🎉'), ('''quotes'''), ('new\nline')" > /dev/null 2>&1
OUTPUT=$($VSQL test8.vsql "SELECT * FROM special ORDER BY id" 2>&1)
if echo "$OUTPUT" | grep -q "日本語" && echo "$OUTPUT" | grep -q "🎉"; then
    log_pass "Unicode and special chars work"
else
    log_fail "Special chars failed: $OUTPUT"
fi

# Test 15: Large result set
log_test "Large result set (1000 rows)"
rm -rf test9.vsql*
$VSQL test9.vsql "CREATE TABLE large (id SERIAL)" > /dev/null 2>&1
$VSQL test9.vsql "INSERT INTO large SELECT generate_series(1, 1000)" > /dev/null 2>&1
COUNT=$($VSQL test9.vsql "SELECT COUNT(*) as c FROM large" 2>&1 | grep -o '"c":[0-9]*' | cut -d: -f2)
if [ "$COUNT" = "1000" ]; then
    log_pass "Large result set works"
else
    log_fail "Large result set failed: $COUNT"
fi

# Test 16: Binary files in cache
log_test "Shared libraries extracted"
CACHE_DIR=$(find "$HOME/.cache" -name "vibesql-micro" -type d 2>/dev/null | head -1)
if [ -f "$CACHE_DIR"/*/postgres ] && [ -f "$CACHE_DIR"*/libpq.so.5 ]; then
    log_pass "Binaries and libraries extracted"
else
    log_fail "Missing binaries in cache: $CACHE_DIR"
fi

# Cleanup
echo ""
echo "=========================================="
echo "Test Results"
echo "=========================================="
echo -e "${GREEN}PASSED: $PASSED${NC}"
echo -e "${RED}FAILED: $FAILED${NC}"

# Cleanup test files
cd /
rm -rf "$TEST_DIR"

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}SOME TESTS FAILED${NC}"
    exit 1
fi
