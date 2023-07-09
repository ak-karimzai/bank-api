#!/bin/sh

set -e

echo "Run db migration"

cd migration && /app/goose postgres "$DB_SOURCE" up && cd ..

echo "start the app"
exec "$@"