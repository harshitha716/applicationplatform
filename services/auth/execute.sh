#!/bin/sh

echo "Running migration"

# Run migrations if needed
# Note that we have a DSN string provided via an environment variable though you may pass a config here as well
kratos migrate sql -e --yes

echo "Starting auth server"

# Run auth server
kratos serve -c /home/ory/kratos.yml