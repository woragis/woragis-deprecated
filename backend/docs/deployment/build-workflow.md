# Build Workflow Documentation

**Last Updated:** 2025-12-23  
**Purpose:** Documentation for the CI/CD build workflow

---

## Overview

The build workflow automatically builds Docker images for all services when a version tag is pushed to the repository. This document explains how the workflow works and how to use it.

---

## Workflow File

**Location:** `.github/workflows/build-all.yml`

**Trigger:** Pushes to tags matching `v*` (e.g., `v1.0.0`, `v8.5`, `v2.1.3`)

---

## Services Built

The workflow builds Docker images for all 9 backend services:

### Go Services
1. **Server** (`backend/server/app`)
2. **Email Worker** (`backend/email-worker`)
3. **Translation Worker** (`backend/translation-worker`)
4. **WhatsApp Worker** (`backend/whatsapp-worker`)

### Python Services
5. **AI Service** (`backend/ai-service`)
6. **Creative Service** (`backend/creative-service`)
7. **Docs Service** (`backend/docs-service`)
8. **Resume Worker** (`backend/resume-worker`)

### Node.js Services
9. **Job Application Worker** (`backend/job-application-worker`)

---

## Workflow Steps

### 1. Extract Tag
- Extracts the tag name from the Git reference
- Outputs tag for use in subsequent jobs

### 2. Test Services
- Runs unit tests for all services
- Uses matrix strategy to test each service in parallel
- Uploads test coverage reports
- **Languages:** Go, Python, Node.js

### 3. Build Docker Images
- Builds Docker images for all services
- Tags images with:
  - `latest` (for latest tag)
  - `{tag-name}` (e.g., `v1.0.0`)
  - `{tag-name}-{service-name}` (e.g., `v1.0.0-server`)

### 4. Push to Docker Hub
- Pushes all images to Docker Hub
- Uses credentials from GitHub secrets:
  - `DOCKER_HUB_USERNAME`
  - `DOCKER_HUB_TOKEN`

### 5. Deploy to Railway (Optional)
- Deploys services to Railway if configured
- Uses Railway CLI to update service images
- Requires:
  - `RAILWAY_TOKEN`
  - `RAILWAY_PROJECT_ID`

---

## Usage

### Creating a Release

1. **Create and push a tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Monitor the workflow:**
   - Go to GitHub Actions tab
   - Watch the "Build All Services" workflow
   - Check each job for status

3. **Verify images:**
   ```bash
   # Check Docker Hub
   docker pull yourusername/woragis-server:v1.0.0
   ```

### Manual Trigger

You can also trigger the workflow manually:

1. Go to GitHub Actions
2. Select "Build All Services"
3. Click "Run workflow"
4. Enter a tag name (e.g., `v1.0.0`)
5. Click "Run"

---

## Docker Images

### Image Naming Convention

Images are named: `{docker-hub-username}/woragis-{service-name}:{tag}`

**Examples:**
- `yourusername/woragis-server:v1.0.0`
- `yourusername/woragis-ai-service:v1.0.0`
- `yourusername/woragis-email-worker:v1.0.0`

### Image Tags

Each image is tagged with:
- `latest` - Latest version (only for most recent tag)
- `{tag-name}` - Specific version tag
- `{tag-name}-{service-name}` - Service-specific tag

---

## Required Secrets

### Docker Hub
- `DOCKER_HUB_USERNAME` - Your Docker Hub username
- `DOCKER_HUB_TOKEN` - Docker Hub access token

### Railway (Optional)
- `RAILWAY_TOKEN` - Railway project token
- `RAILWAY_PROJECT_ID` - Railway project ID

---

## Build Configuration

### Go Services

**Build Command:**
```bash
docker build -t {image-name}:{tag} -f Dockerfile .
```

**Dockerfile Location:**
- Server: `backend/server/app/Dockerfile`
- Workers: `backend/{worker-name}/Dockerfile`

### Python Services

**Build Command:**
```bash
docker build -t {image-name}:{tag} -f Dockerfile .
```

**Dockerfile Location:**
- Services: `backend/{service-name}/Dockerfile`
- Workers: `backend/{worker-name}/Dockerfile`

### Node.js Services

**Build Command:**
```bash
docker build -t {image-name}:{tag} -f Dockerfile .
```

**Dockerfile Location:**
- `backend/job-application-worker/Dockerfile`

---

## Test Execution

Before building, the workflow runs tests for each service:

### Go Services
```bash
go test ./... -v -coverprofile=coverage.out
```

### Python Services
```bash
pytest tests/ -v --cov=. --cov-report=xml
```

### Node.js Services
```bash
npm test -- --coverage
```

---

## Build Matrix

The workflow uses a matrix strategy to build services in parallel:

```yaml
strategy:
  fail-fast: false
  matrix:
    include:
      - service_name: server
        language: go
        working_directory: ./backend/server/app
      - service_name: email-worker
        language: go
        working_directory: ./backend/email-worker
      # ... more services
```

**Benefits:**
- Parallel execution (faster builds)
- Independent failure (one service failure doesn't stop others)
- Better visibility (see status per service)

---

## Troubleshooting

### Build Fails

**Check:**
1. Dockerfile syntax
2. Dependencies in requirements.txt/package.json/go.mod
3. Build context paths
4. GitHub Actions logs

**Common Issues:**
- Missing dependencies
- Incorrect Dockerfile paths
- Build context issues
- Resource limits

### Image Push Fails

**Check:**
1. Docker Hub credentials
2. Image name format
3. Docker Hub rate limits
4. Network connectivity

**Common Issues:**
- Invalid credentials
- Rate limit exceeded
- Image name too long
- Network timeout

### Tests Fail

**Check:**
1. Test code
2. Test dependencies
3. Test environment variables
4. Test database/Redis/RabbitMQ connections

**Common Issues:**
- Missing test dependencies
- Flaky tests
- Environment configuration
- Service availability

---

## Best Practices

### 1. Use Semantic Versioning
- Use tags like `v1.0.0`, `v1.1.0`, `v2.0.0`
- Follow semantic versioning (major.minor.patch)

### 2. Test Before Tagging
- Run tests locally before creating a tag
- Fix any failing tests
- Ensure all services build successfully

### 3. Monitor Build Times
- Keep build times reasonable (< 30 minutes)
- Optimize Dockerfiles for faster builds
- Use build cache effectively

### 4. Tag Management
- Don't reuse tags
- Use descriptive tag names
- Document what changed in each version

### 5. Security
- Don't commit secrets
- Use GitHub Secrets for sensitive data
- Scan images for vulnerabilities

---

## Related Documentation

- **Railway Setup:** `.github/workflows/RAILWAY_SETUP.md`
- **Deployment Guide:** `docs/deployment/deployment-procedures.md`
- **Docker Setup:** `docs/deployment/docker-setup.md`
- **CI/CD Overview:** `.github/workflows/README.md`

---

## Workflow Status

You can check workflow status:
- **GitHub Actions:** Repository → Actions tab
- **Docker Hub:** Your Docker Hub repository
- **Railway:** Railway dashboard (if deployed)

---

**Last Updated:** 2025-12-23

