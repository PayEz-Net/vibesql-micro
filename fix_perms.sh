#!/bin/bash
echo "=== Fix /opt/vibesql ownership ==="
sudo chown -R zenflow:zenflow /opt/vibesql
sudo chmod -R u+rwX /opt/vibesql
echo "Done"
touch /opt/vibesql/test_write && rm /opt/vibesql/test_write && echo "Write test: OK" || echo "Write test: FAILED"
