#!/bin/bash

set -e

# Set up the test database
./setup_test_db.sh

# Run Go tests
echo "Running Go tests..."
go test ./... -v

if [ $? -eq 0 ]; then
  echo "All tests passed!"
else
  echo "Some tests failed."
  exit 1
fi 