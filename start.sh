#!/bin/sh

set -e

ls -la 

echo "run db migration"
./goose postgres "$DB_SOURCE" up

echo "run app"
exec "$@"