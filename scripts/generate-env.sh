#!/bin/bash

# Script to generate .env files from GitHub Variables and Secrets
# This script reads from GitHub Actions environment variables and creates proper .env files

set -e  # Exit on any error

echo "Generating .env files from GitHub Variables and Secrets..."

# Create backend app .env file
mkdir -p backend/backend-app
cat > backend/backend-app/.env << EOF
# Настройки сервера
PORT=$BACKEND_PORT
UPLOAD_PATH=$BACKEND_UPLOAD_PATH
MAX_FILE_SIZE=$BACKEND_MAX_FILE_SIZE  # 50MB в байтах
JWT_SECRET=$JWT_SECRET
MAX_ATTEMPTS_HANDLED=$BACKEND_MAX_ATTEMPTS_HANDLED  # Количество попыток обработки файла
HANDLER_TIMEOUT=$BACKEND_HANDLER_TIMEOUT # Время ожидания после бесплатного лимита обработок, в часах

# Настройки базы данных PostgreSQL
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME

ML_SERVICE_URL=$ML_SERVICE_URL
ML_SERVICE_ENABLED=$ML_SERVICE_ENABLED
EOF

echo "Created backend/backend-app/.env file"
cat backend/backend-app/.env
echo ""

# Create ML service .env file
mkdir -p backend/ml
cat > backend/ml/.env << EOF
ML_WORKERS=$ML_WORKERS
EOF

echo "Created backend/ml/.env file"
cat backend/ml/.env
echo ""

# Create frontend .env file
mkdir -p frontend
cat > frontend/.env << EOF
NEXT_PUBLIC_API_LINK=$NEXT_PUBLIC_API_LINK
EOF

echo "Created frontend/.env file"
cat frontend/.env
echo ""

echo ".env files generation completed successfully!"