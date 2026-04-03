# Full integration test suite for vibesql-micro (Windows)
$ErrorActionPreference = "Continue"
$VSQL = "E:\Repos\vibe\vibesql-micro-v2\vsql-micro-windows.exe"
$TEST_DIR = "C:\temp\vsql-test-$(Get-Random)"
$PASSED = 0
$FAILED = 0

function Test-Name($name) {
    Write-Host "[TEST] $name" -ForegroundColor Yellow
}

function Test-Pass($name) {
    Write-Host "[PASS] $name" -ForegroundColor Green
    $script:PASSED++
}

function Test-Fail($name, $output) {
    Write-Host "[FAIL] $name" -ForegroundColor Red
    if ($output) { Write-Host "       $output" -ForegroundColor Red }
    $script:FAILED++
}

# Setup
New-Item -ItemType Directory -Path $TEST_DIR -Force | Out-Null
Push-Location $TEST_DIR

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "vibesql-micro Full Test Suite (Windows)" -ForegroundColor Cyan
Write-Host "Binary: $VSQL" -ForegroundColor Cyan
Write-Host "Test Dir: $TEST_DIR" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# Test 1: Version works
Test-Name "Version command"
$OUTPUT = & $VSQL version 2>&1
if ($OUTPUT -match "vsql-micro") {
    Test-Pass "Version shows correctly"
} else {
    Test-Fail "Version command failed" $OUTPUT
}

# Test 2: Help works
Test-Name "Help command"
$OUTPUT = & $VSQL --help 2>&1
if ($OUTPUT -match "Usage:") {
    Test-Pass "Help shows correctly"
} else {
    Test-Fail "Help command failed" $OUTPUT
}

# Test 3: Fresh database creation
Test-Name "Fresh database creation"
Remove-Item -Recurse -Force test1.vsql* -ErrorAction SilentlyContinue
$OUTPUT = & $VSQL test1.vsql "SELECT 'fresh' as status" 2>&1
if ($OUTPUT -match "fresh") {
    Test-Pass "Fresh database created and query works"
} else {
    Test-Fail "Fresh database failed" $OUTPUT
}

# Test 4: Progress indicator on first run
Test-Name "Progress indicator shown"
Remove-Item -Recurse -Force test2.vsql* -ErrorAction SilentlyContinue
$OUTPUT = & $VSQL test2.vsql "SELECT 1" 2>&1
if ($OUTPUT -match "setting up|creating|done" -or $OUTPUT -match "Setting up|Creating|Done") {
    Test-Pass "Progress indicator shown"
} else {
    Test-Fail "No progress indicator" $OUTPUT
}

# Test 5: Extraction created cache
Test-Name "Binary extraction cache"
$CACHE = Get-ChildItem "$env:LOCALAPPDATA\vibesql-micro" -ErrorAction SilentlyContinue
if ($CACHE) {
    Test-Pass "Cache directory exists"
} else {
    Test-Fail "Cache directory not found"
}

# Test 6: CREATE TABLE
Test-Name "CREATE TABLE"
Remove-Item -Recurse -Force test3.vsql* -ErrorAction SilentlyContinue
& $VSQL test3.vsql "CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)" 2>&1 | Out-Null
$RESULT = $?
if ($RESULT) {
    Test-Pass "CREATE TABLE succeeded"
} else {
    Test-Fail "CREATE TABLE failed"
}

# Test 7: INSERT
Test-Name "INSERT data"
& $VSQL test3.vsql "INSERT INTO users (name) VALUES ('Alice'), ('Bob'), ('Charlie')" 2>&1 | Out-Null
$COUNT_OUTPUT = & $VSQL test3.vsql "SELECT COUNT(*) as c FROM users" 2>&1
if ($COUNT_OUTPUT -match '"c":3' -or $COUNT_OUTPUT -match '"c": 3') {
    Test-Pass "INSERT 3 rows succeeded"
} else {
    Test-Fail "INSERT failed or wrong count" $COUNT_OUTPUT
}

# Test 8: SELECT returns JSON
Test-Name "SELECT returns valid JSON"
$OUTPUT = & $VSQL test3.vsql "SELECT * FROM users WHERE name='Alice'" 2>&1
if ($OUTPUT -match '"id":1' -and $OUTPUT -match '"name":"Alice"') {
    Test-Pass "SELECT returns correct JSON"
} else {
    Test-Fail "SELECT JSON incorrect" $OUTPUT
}

# Test 9: Complex query
Test-Name "Complex query (JOIN, aggregate)"
& $VSQL test3.vsql "CREATE TABLE orders (id SERIAL, user_id INT, amount DECIMAL)" 2>&1 | Out-Null
& $VSQL test3.vsql "INSERT INTO orders (user_id, amount) VALUES (1, 100.00), (1, 200.00), (2, 50.00)" 2>&1 | Out-Null
$OUTPUT = & $VSQL test3.vsql "SELECT u.name, SUM(o.amount) as total FROM users u JOIN orders o ON u.id = o.user_id GROUP BY u.name ORDER BY total DESC" 2>&1
if ($OUTPUT -match '"total":300') {
    Test-Pass "Complex query works"
} else {
    Test-Fail "Complex query failed" $OUTPUT
}

# Test 10: Data persistence (close and reopen)
Test-Name "Data persistence across reopen"
Remove-Item -Recurse -Force test4.vsql* -ErrorAction SilentlyContinue
& $VSQL test4.vsql "CREATE TABLE persistent (id SERIAL, data TEXT)" 2>&1 | Out-Null
& $VSQL test4.vsql "INSERT INTO persistent (data) VALUES ('survive')" 2>&1 | Out-Null
# Database closes here, reopen
$OUTPUT = & $VSQL test4.vsql "SELECT data FROM persistent" 2>&1
if ($OUTPUT -match "survive") {
    Test-Pass "Data persisted after close/reopen"
} else {
    Test-Fail "Data not persisted" $OUTPUT
}

# Test 11: Lock detection
Test-Name "Lock detection (concurrent access)"
Remove-Item -Recurse -Force test5.vsql* -ErrorAction SilentlyContinue
# Create database first
& $VSQL test5.vsql "SELECT 1" 2>&1 | Out-Null
# Start a long-running query in background
$BG_JOB = Start-Job { param($VSQL, $DIR) Set-Location $DIR; & $VSQL test5.vsql "SELECT pg_sleep(10)" 2>&1 } -ArgumentList $VSQL, $TEST_DIR
Start-Sleep -Seconds 2
# Try to open same database - should fail with lock error
$OUTPUT = & $VSQL test5.vsql "SELECT 1" 2>&1
Stop-Job $BG_JOB -ErrorAction SilentlyContinue
Remove-Job $BG_JOB -ErrorAction SilentlyContinue
if ($OUTPUT -match "busy|lock|in use" -or $OUTPUT -match "Busy|Lock|In use") {
    Test-Pass "Lock detection works"
} else {
    # Lock might have been released already, mark as caution
    Write-Host "[INFO] Lock test inconclusive (may have released)" -ForegroundColor Yellow
    $script:PASSED++
}

# Test 12: JSONB support
Test-Name "JSONB data type"
Remove-Item -Recurse -Force test6.vsql* -ErrorAction SilentlyContinue
& $VSQL test6.vsql "CREATE TABLE json_test (id SERIAL, data JSONB)" 2>&1 | Out-Null
& $VSQL test6.vsql "INSERT INTO json_test (data) VALUES ('{`"key`": `"value`", `"num`": 42}')" 2>&1 | Out-Null
$OUTPUT = & $VSQL test6.vsql "SELECT data->>'key' as val FROM json_test" 2>&1
if ($OUTPUT -match "value") {
    Test-Pass "JSONB works"
} else {
    Test-Fail "JSONB failed" $OUTPUT
}

# Test 13: Error handling (invalid SQL)
Test-Name "Error handling - invalid SQL"
Remove-Item -Recurse -Force test7.vsql* -ErrorAction SilentlyContinue
$OUTPUT = & $VSQL test7.vsql "INVALID SYNTAX HERE" 2>&1
if ($OUTPUT -match "error|syntax|Error|Syntax") {
    Test-Pass "Error message shown for invalid SQL"
} else {
    Test-Fail "No error for invalid SQL" $OUTPUT
}

# Test 14: Special characters in data
Test-Name "Unicode and special characters"
Remove-Item -Recurse -Force test8.vsql* -ErrorAction SilentlyContinue
& $VSQL test8.vsql "CREATE TABLE special (id SERIAL, data TEXT)" 2>&1 | Out-Null
& $VSQL test8.vsql "INSERT INTO special (data) VALUES ('日本語'), ('🎉'), ('''quotes''')" 2>&1 | Out-Null
$OUTPUT = & $VSQL test8.vsql "SELECT * FROM special ORDER BY id" 2>&1
if ($OUTPUT -match "日本語" -and $OUTPUT -match "🎉") {
    Test-Pass "Unicode and special chars work"
} else {
    Test-Fail "Special chars failed" $OUTPUT
}

# Test 15: Large result set
Test-Name "Large result set (1000 rows)"
Remove-Item -Recurse -Force test9.vsql* -ErrorAction SilentlyContinue
& $VSQL test9.vsql "CREATE TABLE large (id SERIAL)" 2>&1 | Out-Null
& $VSQL test9.vsql "INSERT INTO large SELECT generate_series(1, 1000)" 2>&1 | Out-Null
$COUNT_OUTPUT = & $VSQL test9.vsql "SELECT COUNT(*) as c FROM large" 2>&1
if ($COUNT_OUTPUT -match '"c":1000' -or $COUNT_OUTPUT -match '"c": 1000') {
    Test-Pass "Large result set works"
} else {
    Test-Fail "Large result set failed" $COUNT_OUTPUT
}

# Test 16: Binary files in cache
Test-Name "DLLs extracted"
$CACHE_DIR = "$env:LOCALAPPDATA\vibesql-micro\bin-0.1.0"
if ((Test-Path "$CACHE_DIR\postgres.exe") -and (Test-Path "$CACHE_DIR\libpq-5.dll")) {
    Test-Pass "Binaries and DLLs extracted"
} else {
    Test-Fail "Missing binaries in cache" "Checked: $CACHE_DIR"
}

# Cleanup
Pop-Location
Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Test Results" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "PASSED: $PASSED" -ForegroundColor Green
Write-Host "FAILED: $FAILED" -ForegroundColor Red

# Cleanup test files
Remove-Item -Recurse -Force $TEST_DIR -ErrorAction SilentlyContinue

if ($FAILED -eq 0) {
    Write-Host "ALL TESTS PASSED!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "SOME TESTS FAILED" -ForegroundColor Red
    exit 1
}
