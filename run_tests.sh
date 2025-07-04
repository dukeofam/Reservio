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
GO_TEST_OUTPUT=$(go test ./... -v)
echo "$GO_TEST_OUTPUT"

EXIT_CODE=$?

# ANSI colours
GREEN="$(tput setaf 2)"
RED="$(tput setaf 1)"
RESET="$(tput sgr0)"

if [ $EXIT_CODE -eq 0 ]; then
  echo "${GREEN}OK${RESET}  All tests passed!"
else
  echo "${RED}FAIL${RESET}  Some tests failed."
  exit $EXIT_CODE
fi 