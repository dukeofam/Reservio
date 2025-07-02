#!/bin/bash

# Exit on error
set -e

DB_USER="reservio"
DB_PASS="reservio"
DB_NAME="reservio_test"

# Create user if not exists
psql postgres -tc "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'" | grep -q 1 || \
  psql postgres -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASS';"

# Create test database if not exists
psql postgres -tc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" | grep -q 1 || \
  psql postgres -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;"

# Grant privileges
psql postgres -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"

echo "Test database setup complete." 