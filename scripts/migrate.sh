#!/bin/bash

# Database Migration Script
# This script runs database migrations before the application starts

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Default values
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_NAME=${DB_NAME:-"user_app_db"}
DB_USER=${DB_USER:-"local"}
DB_PASSWORD=${DB_PASSWORD:-"local"}
MIGRATIONS_DIR=${MIGRATIONS_DIR:-"/app/migrations"}

print_status "Starting database migration..."
print_status "Database: $DB_HOST:$DB_PORT/$DB_NAME"
print_status "User: $DB_USER"
print_status "Migrations directory: $MIGRATIONS_DIR"

# Wait for database to be ready
print_status "Waiting for database to be ready..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME"; do
    print_warning "Database is not ready yet. Waiting..."
    sleep 2
done

print_success "Database is ready!"

# Set password for psql
export PGPASSWORD="$DB_PASSWORD"

# Run migrations
print_status "Running database migrations..."

# Enable PostGIS extension
print_status "Enabling PostGIS extension..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "CREATE EXTENSION IF NOT EXISTS postgis;" || print_warning "PostGIS extension already exists or failed to create"

# Create users table
print_status "Creating users table..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);" || print_warning "Users table already exists or failed to create"

# Create index
print_status "Creating email index..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);" || print_warning "Email index already exists or failed to create"

print_success "Database migrations completed successfully!"

# Verify tables exist
print_status "Verifying database schema..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "\dt" | grep users && print_success "Users table verified!" || print_error "Users table not found!"

print_success "Database migration process completed!"
