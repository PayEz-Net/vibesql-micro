#!/bin/bash
VSQL=/home/dotnetpert/test/vsql-micro-linux
cd /tmp
rm -rf dbg.vsql*
echo "=== CREATE TABLE ==="
"$VSQL" dbg.vsql "CREATE TABLE t (id INT)"
echo "=== INSERT ==="
"$VSQL" dbg.vsql "INSERT INTO t VALUES (1), (2), (3)"
echo "=== SELECT COUNT ==="
"$VSQL" dbg.vsql "SELECT COUNT(*) as c FROM t" 2>&1
echo "=== Check pattern ==="
"$VSQL" dbg.vsql "SELECT COUNT(*) as c FROM t" 2>&1 | grep -o '"c":[0-9]*'
echo "=== Cache check ==="
ls ~/.cache/vibesql-micro/bin-*/postgres
ls ~/.cache/vibesql-micro/bin-*/libpq.so.5
