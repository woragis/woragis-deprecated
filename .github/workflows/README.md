# CI/CD Workflows

This directory contains GitHub Actions workflows for testing, building, and deploying all backend services.

## Structure

### Reusable Workflows (`reusable/`)

Reusable workflows that can be called by other workflows:

- **`test-service.yml`** - Tests a single service (Go, Python, or Node.js)
- **`build-service.yml`** - Builds and pushes a Docker image for a service
- **`deploy-service.yml`** - Deploys a service to Railway

### Main Workflows

- **`test-all.yml`** - Runs tests for all services on push/PR (uses matrix strategy)
- **`build-all.yml`** - Builds all services on version tags (uses matrix strategy)
- **`deploy-all.yml`** - Deploys all services to Railway after successful build (uses matrix strategy)

## How It Works

### Testing

The `test-all.yml` workflow:
- Triggers on push to `main`/`develop` or on pull requests
- Uses a matrix strategy to test all 8 services in parallel
- Each service test runs independently using the reusable `test-service.yml` workflow

### Building

The `build-all.yml` workflow:
- Triggers on version tags (e.g., `v1.0.0`)
- First extracts the tag name
- Runs tests for all services (parallel)
- If tests pass, builds Docker images for all services (parallel)
- Generates a summary of all built images

### Deploying

The `deploy-all.yml` workflow:
- Triggers after `build-all.yml` completes successfully
- Extracts the tag from the triggering workflow
- Deploys all services to Railway in parallel
- Sets environment variables (DATABASE_URL, REDIS_URL) as needed per service
- Generates a deployment summary

## Services

The following services are included in all workflows:

1. **server** (Go) - Main API server
2. **email-worker** (Go) - Email processing worker
3. **translation-worker** (Go) - Translation processing worker
4. **whatsapp-worker** (Go) - WhatsApp messaging worker
5. **job-application-worker** (Node.js) - Job application automation
6. **resume-worker** (Python) - Resume generation worker
7. **ai-service** (Python) - AI service API
8. **creative-service** (Python) - Creative content generation service

## Matrix Strategy Benefits

- **Parallel Execution**: All services are tested/built/deployed in parallel
- **Independent Status**: Each service has its own job status in GitHub Actions
- **Fail-Fast Control**: `fail-fast: false` allows all services to complete even if one fails
- **Easy to Add/Remove**: Just add/remove entries from the matrix

## Adding a New Service

To add a new service:

1. Add an entry to the matrix in `test-all.yml`:
   ```yaml
   - service_name: new-service
     language: go  # or python, node
     working_directory: ./backend/new-service
     cache_dependency_path: backend/new-service/go.sum
   ```

2. Add an entry to the build matrix in `build-all.yml`:
   ```yaml
   - service_name: new-service
     image_name: woragis-new-service
     dockerfile_path: ./backend/new-service/Dockerfile
     context_path: ./backend/new-service
   ```

3. Add an entry to the deploy matrix in `deploy-all.yml`:
   ```yaml
   - service_name: woragis-new-service
     image_name: woragis-new-service
     needs_db: true
     needs_redis: false
   ```

## Secrets Required

- `DOCKER_HUB_USERNAME` - Docker Hub username
- `DOCKER_HUB_TOKEN` - Docker Hub access token
- `RAILWAY_TOKEN` - Railway API token
- `RAILWAY_PROJECT_ID` - Railway project ID
- `RAILWAY_DATABASE_URL` - Database connection string (optional, for services that need it)
- `RAILWAY_REDIS_URL` - Redis connection string (optional, for services that need it)
- `DATABASE_URL` - Test database URL (optional, defaults to local test DB)
- `REDIS_URL` - Test Redis URL (optional, defaults to local test Redis)

## Workflow Triggers

- **Test**: Push to `main`/`develop`, or pull requests to `main`/`develop`
- **Build**: Push of version tags matching `v*` pattern
- **Deploy**: Automatically after successful build workflow completion

## Benefits of This Structure

1. **DRY (Don't Repeat Yourself)**: Common logic is in reusable workflows
2. **Maintainability**: Easy to update test/build/deploy logic in one place
3. **Visibility**: Clear per-service status in GitHub Actions UI
4. **Scalability**: Easy to add new services without duplicating code
5. **Parallelism**: All services run in parallel for faster CI/CD
6. **Flexibility**: Can customize per-service behavior via matrix parameters
