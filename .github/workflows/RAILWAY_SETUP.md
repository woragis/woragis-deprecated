# Railway Deployment Setup Guide

This guide explains how to set up Railway deployment for your backend services.

## Prerequisites

1. **Railway Account**: Sign up at [railway.app](https://railway.app)
2. **Railway Project**: Create a new project in Railway
3. **Services Created**: Create one service for each component:
   - `woragis-server` (Main Go server)
   - `woragis-ai-service` (Python AI service)
   - `woragis-job-application-worker` (Node.js worker)
   - `woragis-resume-worker` (Python worker)
   - `woragis-translation-worker` (Go worker)
   - `woragis-whatsapp-worker` (Go worker)

## Railway Setup Steps

### 1. Create Railway Services

In your Railway project dashboard:

1. Click "New" → "Service"
2. For each service, select "Empty Service"
3. Name them exactly:
   - `woragis-server`
   - `woragis-ai-service`
   - `woragis-job-application-worker`
   - `woragis-resume-worker`
   - `woragis-translation-worker`
   - `woragis-whatsapp-worker`

### 2. Add PostgreSQL Database

1. In Railway project, click "New" → "Database" → "PostgreSQL"
2. Railway will automatically create a PostgreSQL service
3. Note the connection string (you'll need it for secrets)

### 3. Add Redis

1. In Railway project, click "New" → "Database" → "Redis"
2. Railway will automatically create a Redis service
3. Note the connection string (you'll need it for secrets)

### 4. Configure Services to Use Docker Images

For each service:

1. Go to the service settings
2. Under "Source", select "Docker Hub"
3. Set the image name format: `yourusername/woragis-{service-name}:latest`
4. Enable "Watch Docker Hub" if you want automatic updates (optional)

**OR** use the Railway CLI (recommended - handled by GitHub Actions):

The GitHub Actions workflow will automatically update the Docker image for each service when you push a tag.

### 5. Get Railway Credentials

1. Go to your Railway project → Settings → Tokens
2. Click "New Token"
3. Name it "GitHub Actions" or similar
4. Copy the token (you'll add it to GitHub secrets)

5. Get your Project ID:
   - Go to project settings
   - The Project ID is in the URL or settings page
   - Format: Usually a UUID or short identifier

### 6. Get Database and Redis Connection Strings

1. **PostgreSQL**:
   - Go to your PostgreSQL service
   - Click "Variables" tab
   - Copy the `DATABASE_URL` or `POSTGRES_URL` value
   - Format: `postgresql://user:password@host:port/database`

2. **Redis**:
   - Go to your Redis service
   - Click "Variables" tab
   - Copy the `REDIS_URL` value
   - Format: `redis://user:password@host:port` or `rediss://...` for SSL

## GitHub Secrets Setup

Add these secrets to your GitHub repository:

1. Go to your repository → Settings → Secrets and variables → Actions
2. Add the following secrets:

### Required Secrets

- `RAILWAY_TOKEN`: Your Railway project token
- `RAILWAY_PROJECT_ID`: Your Railway project ID
- `RAILWAY_DATABASE_URL`: PostgreSQL connection string from Railway
- `RAILWAY_REDIS_URL`: Redis connection string from Railway
- `DOCKER_HUB_USERNAME`: Your Docker Hub username (already set)
- `DOCKER_HUB_TOKEN`: Your Docker Hub token (already set)

### Optional Service-Specific Secrets

You can add service-specific environment variables as secrets and they'll be set automatically:

- `RAILWAY_SERVER_*`: Any environment variable for main server (e.g., `RAILWAY_SERVER_PORT`)
- `RAILWAY_AI_SERVICE_*`: Environment variables for AI service
- `RAILWAY_JOB_WORKER_*`: Environment variables for job application worker
- etc.

## How It Works

1. **You push a tag**: `git tag v8.5 && git push origin v8.5`

2. **GitHub Actions runs**:
   - ✅ Tests all code (Go, Python, Node.js)
   - ✅ Builds Docker images for all services
   - ✅ Pushes images to Docker Hub with tag `v8.5`
   - ✅ Deploys to Railway by updating each service's Docker image

3. **Railway updates**:
   - Each service pulls the new Docker image from Docker Hub
   - Services restart with the new image
   - Environment variables (DATABASE_URL, REDIS_URL) are set automatically

## Deployment Methods

The workflow includes multiple deployment methods:

### Method 1: Railway CLI (Current)
- Uses Railway CLI to update services
- May require manual configuration if CLI doesn't support `--image` flag

### Method 2: Railway API (Alternative)
- Use `railway-deploy-api.yml` workflow
- Uses Railway's GraphQL API to update service images
- More reliable but requires correct API access

### Method 3: Manual Configuration (Fallback)
- Configure each service in Railway dashboard to watch Docker Hub
- Railway will auto-deploy when new tags are pushed
- Set up once, then automatic

**Recommendation**: Start with Method 3 (manual setup), then automate with Method 1 or 2.

## Service Configuration

### Main Server (woragis-server)

**Port**: Railway will auto-detect or set via `PORT` environment variable
**Health Check**: Configure in Railway dashboard if needed

### AI Service (woragis-ai-service)

**Port**: Default 8000, Railway will expose it
**Health Check**: `/health` endpoint (if you have one)

### Workers

Workers don't need exposed ports. They run as background services.

## Environment Variables

The workflow automatically sets:
- `DATABASE_URL` for all services
- `REDIS_URL` for services that need it

You can add more environment variables in Railway dashboard or via the workflow.

## Troubleshooting

### Services not deploying

1. **Check Railway token**: Make sure `RAILWAY_TOKEN` is correct
2. **Check Project ID**: Verify `RAILWAY_PROJECT_ID` matches your project
3. **Check service names**: Service names in Railway must match exactly
4. **Check Docker Hub**: Verify images were pushed successfully

### Services failing to start

1. **Check logs**: Go to Railway dashboard → Service → Logs
2. **Check environment variables**: Verify all required vars are set
3. **Check database connection**: Ensure `DATABASE_URL` is correct
4. **Check image tag**: Verify the image exists in Docker Hub

### Connection issues

1. **Database**: Ensure PostgreSQL service is running in Railway
2. **Redis**: Ensure Redis service is running in Railway
3. **Network**: Services in the same Railway project can communicate via service names

## Manual Deployment

If you need to deploy manually:

```bash
# Install Railway CLI
curl -fsSL https://railway.app/install.sh | sh

# Login
railway login

# Link to project
railway link <project-id>

# Deploy a service
railway service woragis-server --image yourusername/woragis-server:v8.5

# Set environment variables
railway variables set DATABASE_URL="postgresql://..." --service woragis-server
```

## Next Steps

1. Set up Railway project and services
2. Add all secrets to GitHub
3. Push a tag: `git tag v1.0.0 && git push origin v1.0.0`
4. Watch the deployment in GitHub Actions
5. Check Railway dashboard for service status

