#!/bin/bash

# GPT generated script + some refactor.

set -e

DB_SUPERUSER="labib"
DB_SUPERDB="postgres"

DB_USER="tempuser1"
DB_PASS="secret"
DB_NAME="tempdb1"

# 1. Create role safely
psql -v ON_ERROR_STOP=1 -U $DB_SUPERUSER -d $DB_SUPERDB <<EOF
DO \$\$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_roles WHERE rolname = '$DB_USER'
    ) THEN
        CREATE ROLE $DB_USER WITH LOGIN PASSWORD '$DB_PASS';
    END IF;
END
\$\$;
EOF

# 2. Create database safely (OUTSIDE transaction)
psql -v ON_ERROR_STOP=1 -U $DB_SUPERUSER -d $DB_SUPERDB -tc "
SELECT 1 FROM pg_database WHERE datname='$DB_NAME'
" | grep -q 1 || \

psql -v ON_ERROR_STOP=1 -U $DB_SUPERUSER -d $DB_SUPERDB -c "
CREATE DATABASE $DB_NAME OWNER $DB_USER;
"

# 3. Grant privileges
psql -v ON_ERROR_STOP=1 -U $DB_SUPERUSER -d $DB_SUPERDB -c "
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
"

echo "Database and user created successfully!"