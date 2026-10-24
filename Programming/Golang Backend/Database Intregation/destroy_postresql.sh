#!/bin/bash


# GPT generated script + some refactor.


set -e

DB_SUPERUSER="labib"
DB_SUPERDB="postgres"

DB_USER="tempuser1"
DB_NAME="tempdb1"

# 1. Terminate active connections (required)
psql -v ON_ERROR_STOP=1 -U $DB_SUPERUSER -d $DB_SUPERDB -c "
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = '$DB_NAME' AND pid <> pg_backend_pid();
"

# 2. Drop database if exists
psql -v ON_ERROR_STOP=1 -U $DB_SUPERUSER -d $DB_SUPERDB -c "
DROP DATABASE IF EXISTS $DB_NAME;
"

# 3. Drop role if exists
psql -v ON_ERROR_STOP=1 -U $DB_SUPERUSER -d $DB_SUPERDB -c "
DROP ROLE IF EXISTS $DB_USER;
"

echo "Database and user destroyed successfully!"