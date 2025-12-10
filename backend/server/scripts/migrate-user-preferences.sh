#!/bin/bash

# Migration script to create user_preferences for existing users
# This script creates default preferences (en, USD) for all users who don't have them yet

echo "Starting user preferences migration..."

# Get database connection details from environment or use defaults
DB_HOST="${DATABASE_HOST:-localhost}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_NAME="${DATABASE_NAME:-woragis}"
DB_USER="${DATABASE_USER:-postgres}"
DB_PASSWORD="${DATABASE_PASSWORD:-postgres}"

export PGPASSWORD="$DB_PASSWORD"

# SQL to create user_preferences for users who don't have them
SQL="
INSERT INTO user_preferences (id, user_id, default_language, default_currency, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    u.id,
    'en',
    'USD',
    NOW(),
    NOW()
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM user_preferences up WHERE up.user_id = u.id
);
"

echo "Executing migration SQL..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL"

if [ $? -eq 0 ]; then
    echo "✅ Migration completed successfully!"
    echo "All existing users now have default preferences (language: en, currency: USD)"
else
    echo "❌ Migration failed!"
    exit 1
fi

unset PGPASSWORD

