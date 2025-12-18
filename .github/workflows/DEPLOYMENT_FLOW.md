# Deployment Flow

This document explains how the CI/CD pipeline deploys Docker images to Railway.

## Complete Flow

### 1. **Build & Push to Docker Hub** (`build-all.yml`)

When you push a version tag (e.g., `v1.0.0`):

1. **Extract Tag**: Gets the version tag from the git ref
2. **Run Tests**: Tests all services in parallel
3. **Build Images**: Builds Docker images for all services in parallel
4. **Push to Docker Hub**: Pushes images with two tags:
   - `username/woragis-server:v1.0.0` (version tag)
   - `username/woragis-server:latest` (latest tag)

**Result**: Images are available on Docker Hub at:
- `docker.io/username/woragis-server:v1.0.0`
- `docker.io/username/woragis-server:latest`

### 2. **Deploy to Railway** (`deploy-all.yml`)

After `build-all.yml` completes successfully:

1. **Extract Tag**: Gets the same version tag from the build workflow
2. **Deploy All Services**: Deploys each service in parallel using the reusable workflow
3. **Update Railway**: For each service:
   - Sets `DOCKER_IMAGE` environment variable to `docker.io/username/woragis-server:v1.0.0`
   - Railway watches this variable and automatically redeploys
   - Sets additional environment variables (DATABASE_URL, REDIS_URL) if needed

**Result**: Railway services are updated with the new Docker images and automatically redeployed.

## How Railway Deployment Works

The deployment workflow uses Railway's environment variable mechanism:

1. **Set DOCKER_IMAGE Variable**: 
   ```bash
   railway variables set "DOCKER_IMAGE=docker.io/username/woragis-server:v1.0.0" --service woragis-server
   ```

2. **Railway Auto-Redeploy**: 
   - Railway monitors the `DOCKER_IMAGE` variable
   - When it changes, Railway automatically:
     - Pulls the new image from Docker Hub
     - Stops the old container
     - Starts a new container with the new image
     - Performs a rolling update (zero-downtime if configured)

3. **Fallback Methods**: 
   - If setting the variable fails, the workflow tries:
     - Railway service update commands
     - Manual redeploy triggers
   - Provides instructions for manual deployment if all methods fail

## Prerequisites

For this to work, your Railway services must be configured to:

1. **Use Docker Hub as Image Source**: 
   - In Railway dashboard, each service should be configured to pull from Docker Hub
   - The service should watch the `DOCKER_IMAGE` environment variable

2. **Have Required Secrets**:
   - `DOCKER_HUB_USERNAME` - Your Docker Hub username
   - `DOCKER_HUB_TOKEN` - Docker Hub access token
   - `RAILWAY_TOKEN` - Railway API token
   - `RAILWAY_PROJECT_ID` - Your Railway project ID
   - `RAILWAY_DATABASE_URL` - Database connection (if services need it)
   - `RAILWAY_REDIS_URL` - Redis connection (if services need it)

## Verification

After deployment, verify in Railway dashboard:

1. Go to your Railway project
2. Check each service:
   - The `DOCKER_IMAGE` variable should be set to the new image
   - The service should show "Deploying" or "Active" status
   - Check logs to confirm the new image is running

## Troubleshooting

If deployment fails:

1. **Check Railway CLI**: Verify Railway CLI is installed and authenticated
2. **Check Service Names**: Ensure Railway service names match exactly (case-sensitive)
3. **Check Image Exists**: Verify the image exists on Docker Hub with the correct tag
4. **Manual Deployment**: Use the instructions provided in the workflow output to manually update in Railway dashboard

## Example Flow

```
1. Developer pushes tag: git push origin v1.2.0
   ↓
2. build-all.yml triggers
   ↓
3. Tests run (all services in parallel)
   ↓
4. Docker images built (all services in parallel)
   ↓
5. Images pushed to Docker Hub:
   - docker.io/username/woragis-server:v1.2.0
   - docker.io/username/woragis-server:latest
   ↓
6. deploy-all.yml triggers automatically
   ↓
7. For each service:
   - Sets DOCKER_IMAGE=docker.io/username/woragis-server:v1.2.0
   - Railway detects change and redeploys
   ↓
8. All services running new version on Railway
```
