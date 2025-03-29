#!/bin/sh
set -e

# Check if PG_DATABASE_URL is set
if [ -z "$PG_DATABASE_URL" ]; then
  echo "Error: PG_DATABASE_URL is not set"
  exit 1
fi

echo "Running migrations"
# Run the migration
migrate -database "$PG_DATABASE_URL" -path ./ up