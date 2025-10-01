# GitHub Secrets and Variables Setup

This document describes how to set up GitHub Secrets and Variables for the Obscura project CI/CD pipeline.

## GitHub Variables (Non-Sensitive Values)

The following variables should be added to your GitHub repository's Variables:

### Backend Variables
- `BACKEND_PORT`: Port for the backend server (default: 8080)
- `BACKEND_UPLOAD_PATH`: Path for file uploads (default: ./app/uploads)
- `BACKEND_MAX_FILE_SIZE`: Maximum file size in bytes (default: 52428800 for 50MB)
- `BACKEND_MAX_ATTEMPTS_HANDLED`: Number of handling attempts (default: 3)
- `BACKEND_HANDLER_TIMEOUT`: Timeout in hours (default: 24)
- `DB_HOST`: Database host (default: db)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_NAME`: Database name (default: obscura)
- `ML_SERVICE_URL`: URL for ML service (default: http://ml-service:8000)
- `ML_SERVICE_ENABLED`: Whether ML service is enabled (default: true)
- `ML_WORKERS`: Number of ML workers (default: 3)
- `NEXT_PUBLIC_API_LINK`: Frontend API link (default: http://localhost:8080)

## GitHub Secrets (Sensitive Values)

The following secrets should be added to your GitHub repository's Secrets:

### Backend Secrets
- `JWT_SECRET`: Secret key for JWT authentication (required)
- `DB_PASSWORD`: Database password (default: postgres for testing)

## Docker Compose Integration

The generated .env files will be used by the docker-compose.yml file to configure the services during CI/CD deployment.

## CI/CD Workflow

The CI/CD pipeline in `.github/workflows/ci.yml` will:
1. Fetch the variables and secrets from GitHub
2. Generate appropriate .env files in each service directory
3. Use these files in the Docker Compose setup

## Local Development

For local development, you can continue to use the .env files directly. The CI/CD approach only affects the automated deployment process.