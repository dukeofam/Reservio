#!/bin/bash

set -e

# Set up the test database
./setup_test_db.sh

# Connection string for local Postgres test DB
export DATABASE_URL="postgres://reservio:reservio@localhost:5432/reservio_test?sslmode=disable"

# Ensure other required env vars have sane defaults for local test runs
export REDIS_ADDR="${REDIS_ADDR:-localhost:6379}"
export SESSION_SECRET="${SESSION_SECRET:-testsecret}"
export REDIS_PASSWORD="${REDIS_PASSWORD:-}"

# Tell the app we are in test mode (use cookie store instead of redis)
export TEST_MODE=1

# Run Go tests
echo "Running Go tests..."
go test ./... -v

if [ $? -eq 0 ]; then
  echo "All tests passed!"
else
  echo "Some tests failed."
  exit 1
fi 