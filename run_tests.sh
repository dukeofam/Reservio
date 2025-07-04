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

# ANSI colours (optional – some CI environments lack a valid TERM entry)
if command -v tput >/dev/null 2>&1; then
  # Disable 'exit on error' temporarily in case tput fails
  set +e
  GREEN="$(tput setaf 2 2>/dev/null)"
  RED="$(tput setaf 1 2>/dev/null)"
  RESET="$(tput sgr0 2>/dev/null)"
  # Restore strict mode
  set -e
fi

# Fallback to empty strings if colours couldn't be determined
GREEN=${GREEN:-""}
RED=${RED:-""}
RESET=${RESET:-""}

# Run Go tests
echo "Running Go tests..."
GO_TEST_OUTPUT=$(go test ./... -v)
echo "$GO_TEST_OUTPUT"

EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
  echo "${GREEN}OK${RESET}  All tests passed!"
else
  echo "${RED}FAIL${RESET}  Some tests failed."
  exit $EXIT_CODE
fi 