#!/bin/sh

set -e 

echo "Waiting for PostgreSQL to be ready..."
while ! pg_isready -h db -p 5432 -U root; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done

echo "run db migration"
/app/migrate -path /app/migration -database postgres://root:secret@db:5432/simplebank?sslmode=disable -verbose up

echo "starting the application"
exec "$@"