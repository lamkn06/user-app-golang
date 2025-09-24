#!/bin/bash

# Quick start script - minimal version
echo "🚀 Starting database..."

# Start database
docker-compose up -d db

# Wait for database
echo "⏳ Waiting for database..."
until nc -z localhost 5432; do
    sleep 1
done

echo "✅ Database ready!"

# Run migrations
echo "🔄 Running migrations..."
migrate -path migrations -database "postgres://local:local@localhost:5432/db_name?sslmode=disable" up

if [ $? -eq 0 ]; then
    echo "✅ Migrations completed!"
else
    echo "❌ Migration failed!"
    exit 1
fi
