# Railway Setup and Testing Procedures

**Last Updated:** 2025-12-23  
**Purpose:** Complete guide for setting up and testing Railway deployments

---

## Overview

This guide covers setting up Railway for all Woragis backend services, configuring services, and testing deployments.

---

## Prerequisites

1. **Railway Account**
   - Sign up at [railway.app](https://railway.app)
   - Verify email address

2. **GitHub Repository**
   - Repository with all backend services
   - GitHub Actions configured

3. **Docker Hub Account**
   - Account for storing Docker images
   - Access token for GitHub Actions

---

## Step 1: Create Railway Project

1. **Login to Railway**
   - Go to [railway.app](https://railway.app)
   - Click "Login" and authenticate

2. **Create New Project**
   - Click "New Project"
   - Select "Empty Project"
   - Name it: `woragis-backend`

---

## Step 2: Add Database Services

### PostgreSQL Database

1. **Add PostgreSQL**
   - In Railway project, click "New" → "Database" → "PostgreSQL"
   - Railway automatically creates PostgreSQL service
   - Note the service name (usually `Postgres`)

2. **Get Connection String**
   - Click on PostgreSQL service
   - Go to "Variables" tab
   - Copy `DATABASE_URL` or `POSTGRES_URL`
   - Format: `postgresql://user:password@host:port/database`

### Redis

1. **Add Redis**
   - In Railway project, click "New" → "Database" → "Redis"
   - Railway automatically creates Redis service
   - Note the service name (usually `Redis`)

2. **Get Connection String**
   - Click on Redis service
   - Go to "Variables" tab
   - Copy `REDIS_URL`
   - Format: `redis://user:password@host:port` or `rediss://...` for SSL

### RabbitMQ (Optional - if using managed RabbitMQ)

If using Railway's RabbitMQ (or external managed RabbitMQ):

1. **Add RabbitMQ**
   - In Railway project, click "New" → "Database" → "RabbitMQ" (if available)
   - Or use external RabbitMQ service

2. **Get Connection String**
   - Copy `RABBITMQ_URL`
   - Format: `amqp://user:password@host:port/vhost`

---

## Step 3: Create Application Services

Create one service for each backend component:

### Services to Create

1. **woragis-server** (Main Go server)
2. **woragis-ai-service** (Python AI service)
3. **woragis-creative-service** (Python creative service)
4. **woragis-docs-service** (Python docs service)
5. **woragis-email-worker** (Go email worker)
6. **woragis-translation-worker** (Go translation worker)
7. **woragis-whatsapp-worker** (Go WhatsApp worker)
8. **woragis-job-application-worker** (Node.js worker)
9. **woragis-resume-worker** (Python resume worker)

### Creating Services

For each service:

1. **Click "New" → "Service"**
2. **Select "Empty Service"**
3. **Name the service** (e.g., `woragis-server`)
4. **Configure source:**
   - Go to service settings
   - Under "Source", select "Docker Hub"
   - Set image: `yourusername/woragis-{service-name}:latest`
   - Enable "Watch Docker Hub" (optional, for auto-updates)

---

## Step 4: Configure Environment Variables

For each service, configure environment variables:

### Server (woragis-server)

**Required Variables:**
```
DATABASE_URL={from PostgreSQL service}
REDIS_URL={from Redis service}
RABBITMQ_URL={from RabbitMQ service}
JWT_SECRET={generate secure random string}
PORT=8080
ENVIRONMENT=production
```

**Optional Variables:**
```
LOG_LEVEL=info
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=https://woragis.com,https://www.woragis.com
```

### AI Service (woragis-ai-service)

**Required Variables:**
```
OPENAI_API_KEY={your OpenAI API key}
ANTHROPIC_API_KEY={your Anthropic API key}
PORT=8000
```

### Creative Service (woragis-creative-service)

**Required Variables:**
```
OPENAI_API_KEY={your OpenAI API key}
STABILITY_AI_API_KEY={your Stability AI API key}
PORT=8000
```

### Email Worker (woragis-email-worker)

**Required Variables:**
```
RABBITMQ_URL={from RabbitMQ service}
SMTP_HOST={your SMTP host}
SMTP_PORT=587
SMTP_USERNAME={your SMTP username}
SMTP_PASSWORD={your SMTP password}
```

### Translation Worker (woragis-translation-worker)

**Required Variables:**
```
DATABASE_URL={from PostgreSQL service}
RABBITMQ_URL={from RabbitMQ service}
GOOGLE_TRANSLATE_API_KEY={your Google Translate API key}
DEEPL_API_KEY={your DeepL API key}
```

### Job Application Worker (woragis-job-application-worker)

**Required Variables:**
```
DATABASE_URL={from PostgreSQL service}
RABBITMQ_URL={from RabbitMQ service}
```

### Resume Worker (woragis-resume-worker)

**Required Variables:**
```
DATABASE_URL={from PostgreSQL service}
RABBITMQ_URL={from RabbitMQ service}
```

### WhatsApp Worker (woragis-whatsapp-worker)

**Required Variables:**
```
RABBITMQ_URL={from RabbitMQ service}
WHATSAPP_API_KEY={your WhatsApp API key}
```

### Docs Service (woragis-docs-service)

**Required Variables:**
```
PORT=8000
```

---

## Step 5: Configure GitHub Secrets

Add these secrets to your GitHub repository:

1. **Go to Repository → Settings → Secrets and variables → Actions**
2. **Add the following secrets:**

### Required Secrets

```
RAILWAY_TOKEN={your Railway project token}
RAILWAY_PROJECT_ID={your Railway project ID}
DOCKER_HUB_USERNAME={your Docker Hub username}
DOCKER_HUB_TOKEN={your Docker Hub access token}
```

### How to Get Railway Token

1. Go to Railway project → Settings → Tokens
2. Click "New Token"
3. Name it "GitHub Actions"
4. Copy the token

### How to Get Railway Project ID

1. Go to Railway project → Settings
2. The Project ID is in the URL or settings page
3. Format: Usually a UUID or short identifier

---

## Step 6: Deploy Services

### Automatic Deployment (Recommended)

The GitHub Actions workflow automatically deploys when you push a tag:

```bash
# Create and push a tag
git tag v1.0.0
git push origin v1.0.0
```

The workflow will:
1. Build Docker images
2. Push to Docker Hub
3. Deploy to Railway

### Manual Deployment

If you need to deploy manually:

```bash
# Install Railway CLI
curl -fsSL https://railway.app/install.sh | sh

# Login
railway login

# Link to project
railway link <project-id>

# Deploy a service
railway up --service woragis-server
```

---

## Step 7: Testing Procedures

### 7.1 Health Check Tests

**Test Server Health:**
```bash
# Get server URL from Railway dashboard
curl https://your-server-url.railway.app/healthz

# Expected response:
# {"status":"healthy"}
```

**Test All Services:**
```bash
# Server
curl https://server-url.railway.app/healthz

# AI Service
curl https://ai-service-url.railway.app/health

# Creative Service
curl https://creative-service-url.railway.app/health

# Docs Service
curl https://docs-service-url.railway.app/health
```

### 7.2 Functional Tests

**Test API Endpoints:**
```bash
# Get authentication token
TOKEN=$(curl -X POST https://server-url.railway.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}' \
  | jq -r '.token')

# Test authenticated endpoint
curl https://server-url.railway.app/api/projects \
  -H "Authorization: Bearer $TOKEN"
```

### 7.3 Integration Tests

**Test Service Communication:**
```bash
# Test server → AI service
curl -X POST https://server-url.railway.app/api/ai/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"prompt":"Test prompt"}'
```

### 7.4 Load Tests

**Test Under Load:**
```bash
# Install hey (load testing tool)
go install github.com/rakyll/hey@latest

# Run load test
hey -n 1000 -c 50 https://server-url.railway.app/healthz
```

### 7.5 Worker Tests

**Test Workers:**
```bash
# Publish a test message to RabbitMQ
# (Use RabbitMQ management UI or CLI)

# Check worker logs in Railway dashboard
# Verify message was processed
```

---

## Step 8: Monitoring and Logs

### View Logs

**In Railway Dashboard:**
1. Go to service
2. Click "Logs" tab
3. View real-time logs

**Via Railway CLI:**
```bash
railway logs --service woragis-server
```

### Monitor Metrics

**Check Service Status:**
- Go to Railway dashboard
- View service metrics (CPU, memory, network)
- Check deployment history

**Set Up Alerts:**
- Configure alerts in Railway dashboard
- Set up notifications for:
  - Service failures
  - High resource usage
  - Deployment failures

---

## Step 9: Rollback Procedures

### Quick Rollback

**Via Railway Dashboard:**
1. Go to service
2. Click "Deployments"
3. Select previous deployment
4. Click "Redeploy"

**Via Railway CLI:**
```bash
railway rollback --service woragis-server
```

### Manual Rollback

**Deploy Previous Image:**
```bash
railway up --service woragis-server --image yourusername/woragis-server:v0.9.0
```

---

## Troubleshooting

### Service Not Starting

**Check:**
1. Service logs in Railway dashboard
2. Environment variables are set correctly
3. Database/Redis/RabbitMQ connections
4. Image exists in Docker Hub

**Common Issues:**
- Missing environment variables
- Incorrect connection strings
- Image not found
- Port conflicts

### Service Unhealthy

**Check:**
1. Health check endpoint
2. Service logs
3. Dependencies (database, Redis, RabbitMQ)
4. Resource limits

**Common Issues:**
- Health check failing
- Database connection issues
- Resource exhaustion
- Network connectivity

### Deployment Fails

**Check:**
1. GitHub Actions logs
2. Railway deployment logs
3. Docker Hub image availability
4. Railway token validity

**Common Issues:**
- Invalid Railway token
- Image not in Docker Hub
- Network issues
- Railway service limits

---

## Best Practices

### 1. Use Staging Environment
- Create separate Railway project for staging
- Test in staging before production
- Use different Docker image tags

### 2. Monitor Closely
- Watch logs during deployment
- Monitor metrics after deployment
- Set up alerts for issues

### 3. Version Control
- Use semantic versioning for tags
- Don't reuse tags
- Document changes in each version

### 4. Security
- Use Railway's built-in secrets management
- Don't commit secrets to repository
- Rotate secrets regularly

### 5. Backup
- Regular database backups
- Document restore procedures
- Test backups regularly

---

## Testing Checklist

### Pre-Deployment
- [ ] All tests passing locally
- [ ] Docker images build successfully
- [ ] Environment variables configured
- [ ] Secrets added to GitHub

### Post-Deployment
- [ ] All services healthy
- [ ] Health checks passing
- [ ] API endpoints working
- [ ] Workers processing jobs
- [ ] Monitoring configured
- [ ] Logs accessible

### Ongoing
- [ ] Regular health checks
- [ ] Monitor error rates
- [ ] Check resource usage
- [ ] Review logs weekly
- [ ] Test rollback procedures monthly

---

## Related Documentation

- **Build Workflow:** `docs/deployment/build-workflow.md`
- **Deployment Procedures:** `docs/deployment/deployment-procedures.md`
- **Railway Setup (GitHub Actions):** `.github/workflows/RAILWAY_SETUP.md`
- **Troubleshooting:** `docs/runbooks/troubleshooting.md`

---

**Last Updated:** 2025-12-23

