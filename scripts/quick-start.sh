#!/bin/bash

# Quick start script - minimal version
echo "🚀 Starting database..."

# Start database
docker-compose up -d db

# Wait for database
echo "⏳ Waiting for database..."
until pg_isready -h localhost -p 5432 -U local -d db_name; do
    sleep 1
done

echo "✅ Database ready!"
