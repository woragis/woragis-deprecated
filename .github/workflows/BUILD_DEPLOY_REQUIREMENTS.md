# Build and Deploy Workflows - Requirements & Status

## Overview

This document outlines the requirements for the **Build All Services** (`build-all.yml`) and **Deploy All Services** (`deploy-all.yml`) workflows, including required secrets, variables, inputs, and their current status.

---

## 🔨 Build All Services Workflow (`build-all.yml`)

### Trigger
- **Event**: Push to tags matching `v*` (e.g., `v1.0.0`, `v2.5.3`)
- **Workflow**: Runs tests first, then builds and pushes Docker images to Docker Hub

### Required Secrets

| Secret | Required | Purpose | Status |
|--------|----------|---------|--------|
| `DOCKER_HUB_USERNAME` | ✅ **YES** | Docker Hub username for pushing images | ❓ **NEEDS SETUP** |
| `DOCKER_HUB_TOKEN` | ✅ **YES** | Docker Hub access token/password | ❓ **NEEDS SETUP** |
| `DATABASE_URL` | ⚠️ Optional | For test services (defaults to local test DB) | ✅ Optional |
| `REDIS_URL` | ⚠️ Optional | For test services (defaults to local test Redis) | ✅ Optional |

### Required Inputs/Variables
- **None** - All inputs come from the git tag trigger

### Services Built
1. ✅ `server` → `woragis-server`
2. ✅ `email-worker` → `woragis-email-worker`
3. ✅ `translation-worker` → `woragis-translation-worker`
4. ✅ `whatsapp-worker` → `woragis-whatsapp-worker`
5. ✅ `job-application-worker` → `woragis-job-application-worker`
6. ✅ `resume-worker` → `woragis-resume-worker`
7. ✅ `ai-service` → `woragis-ai-service`
8. ✅ `creative-service` → `woragis-creative-service`

### Dockerfile Requirements

| Service | Dockerfile Path | Status |
|---------|---------------|--------|
| server | `./backend/server/Dockerfile` | ✅ Exists |
| email-worker | `./backend/Dockerfile.email-worker` | ✅ Exists |
| translation-worker | `./backend/Dockerfile.translation-worker` | ✅ Exists |
| whatsapp-worker | `./backend/Dockerfile.whatsapp-worker` | ✅ Exists |
| job-application-worker | `./backend/Dockerfile.job-application-worker` | ✅ Exists |
| resume-worker | `./backend/Dockerfile.resume-worker` | ✅ Exists |
| ai-service | `./backend/ai-service/Dockerfile` | ✅ Exists |
| creative-service | `./backend/creative-service/Dockerfile` | ✅ Exists |

### Workflow Status: ⚠️ **WILL FAIL WITHOUT SECRETS**

**Current State:**
- ✅ All Dockerfiles exist
- ✅ All service paths are correct
- ✅ Test workflow is properly configured
- ❌ **Missing required secrets**: `DOCKER_HUB_USERNAME`, `DOCKER_HUB_TOKEN`

**To Make It Work:**
1. Create a Docker Hub account (if you don't have one)
2. Generate a Docker Hub access token:
   - Go to Docker Hub → Account Settings → Security
   - Create a new access token
3. Add secrets to GitHub:
   - Go to Repository → Settings → Secrets and variables → Actions
   - Add `DOCKER_HUB_USERNAME` (your Docker Hub username)
   - Add `DOCKER_HUB_TOKEN` (the access token you created)

---

## 🚀 Deploy All Services Workflow (`deploy-all.yml`)

### Trigger
- **Event**: `workflow_run` (runs after "Build All Services" completes successfully)
- **Condition**: Only runs if build workflow succeeded AND on `main` branch
- **Purpose**: Deploys all built images to Railway

### Required Secrets

| Secret | Required | Purpose | Status |
|--------|----------|---------|--------|
| `RAILWAY_TOKEN` | ✅ **YES** | Railway API token for authentication | ❓ **NEEDS SETUP** |
| `RAILWAY_PROJECT_ID` | ✅ **YES** | Railway project ID to deploy to | ❓ **NEEDS SETUP** |
| `DOCKER_HUB_USERNAME` | ✅ **YES** | Docker Hub username (same as build) | ❓ **NEEDS SETUP** |
| `RAILWAY_DATABASE_URL` | ⚠️ Optional | Database URL for services that need it | ❓ **NEEDS SETUP** |
| `RAILWAY_REDIS_URL` | ⚠️ Optional | Redis URL for services that need it | ❓ **NEEDS SETUP** |

### Required Inputs/Variables
- **None** - All inputs come from the triggering workflow

### Services Deployed
1. ✅ `woragis-server` (needs DB, Redis)
2. ✅ `woragis-ai-service` (needs DB)
3. ✅ `woragis-job-application-worker` (needs DB, Redis)
4. ✅ `woragis-resume-worker` (needs DB, Redis)
5. ✅ `woragis-translation-worker` (needs DB, Redis)
6. ✅ `woragis-whatsapp-worker` (needs DB, Redis)
7. ✅ `woragis-email-worker` (needs DB, Redis)
8. ✅ `woragis-creative-service` (no DB/Redis)

### Railway Environment
- Uses `environment: production` protection
- Requires Railway services to be pre-configured in Railway dashboard

### Workflow Status: ⚠️ **WILL FAIL WITHOUT SECRETS & RAILWAY SETUP**

**Current State:**
- ✅ All service names are correct
- ✅ Deployment logic is properly configured
- ✅ Environment variable handling is correct
- ❌ **Missing required secrets**: `RAILWAY_TOKEN`, `RAILWAY_PROJECT_ID`, `DOCKER_HUB_USERNAME`
- ❌ **Missing optional secrets**: `RAILWAY_DATABASE_URL`, `RAILWAY_REDIS_URL` (needed for most services)

**To Make It Work:**
1. **Set up Railway:**
   - Create a Railway account (if you don't have one)
   - Create a new project in Railway
   - Get your Railway project ID (from Railway dashboard URL or CLI)
   - Generate a Railway API token:
     - Go to Railway → Account Settings → Tokens
     - Create a new token

2. **Set up Railway Services:**
   - Create services in Railway dashboard for each service:
     - `woragis-server`
     - `woragis-ai-service`
     - `woragis-job-application-worker`
     - `woragis-resume-worker`
     - `woragis-translation-worker`
     - `woragis-whatsapp-worker`
     - `woragis-email-worker`
     - `woragis-creative-service`

3. **Set up Railway Database & Redis:**
   - Create a PostgreSQL database in Railway
   - Create a Redis instance in Railway
   - Get connection URLs from Railway dashboard

4. **Add secrets to GitHub:**
   - Go to Repository → Settings → Secrets and variables → Actions
   - Add `RAILWAY_TOKEN` (your Railway API token)
   - Add `RAILWAY_PROJECT_ID` (your Railway project ID)
   - Add `RAILWAY_DATABASE_URL` (PostgreSQL connection string from Railway)
   - Add `RAILWAY_REDIS_URL` (Redis connection string from Railway)
   - Add `DOCKER_HUB_USERNAME` (same as build workflow)

5. **Configure Railway Environment:**
   - Go to Repository → Settings → Environments
   - Create/configure `production` environment
   - Add protection rules if needed (required reviewers, etc.)

---

## 📋 Complete Checklist

### For Build Workflow to Work:
- [ ] Docker Hub account created
- [ ] Docker Hub access token generated
- [ ] `DOCKER_HUB_USERNAME` secret added to GitHub
- [ ] `DOCKER_HUB_TOKEN` secret added to GitHub
- [ ] All Dockerfiles exist (✅ Already done)
- [ ] Test workflow passes (✅ Should work based on integration tests)

### For Deploy Workflow to Work:
- [ ] Railway account created
- [ ] Railway project created
- [ ] Railway API token generated
- [ ] Railway services created for all 8 services
- [ ] Railway PostgreSQL database created
- [ ] Railway Redis instance created
- [ ] `RAILWAY_TOKEN` secret added to GitHub
- [ ] `RAILWAY_PROJECT_ID` secret added to GitHub
- [ ] `RAILWAY_DATABASE_URL` secret added to GitHub
- [ ] `RAILWAY_REDIS_URL` secret added to GitHub
- [ ] `DOCKER_HUB_USERNAME` secret added to GitHub (same as build)
- [ ] `production` environment configured in GitHub (if using protection)

---

## 🔍 Workflow Dependencies

### Build Workflow Flow:
```
Push Tag (v*) 
  → Extract Tag
  → Test All Services (parallel)
  → Build All Services (parallel, if tests pass)
  → Push to Docker Hub
  → Build Summary
```

### Deploy Workflow Flow:
```
Build Workflow Completes Successfully
  → Extract Tag from Build Workflow
  → Deploy All Services to Railway (parallel)
  → Set Environment Variables (DATABASE_URL, REDIS_URL)
  → Deployment Summary
```

---

## ⚠️ Potential Issues

1. **Tag Extraction in Deploy Workflow:**
   - The deploy workflow tries to extract the tag from the build workflow
   - Uses GitHub API to get the latest tag if extraction fails
   - May need adjustment if tag format changes

2. **Railway Deployment Method:**
   - The workflow tries multiple methods to deploy:
     - Setting `DOCKER_IMAGE` environment variable (preferred)
     - Railway service update command (fallback)
     - Manual redeploy trigger (fallback)
   - Railway CLI commands may change, requiring workflow updates

3. **Service Names in Railway:**
   - Service names in Railway must match exactly:
     - `woragis-server`
     - `woragis-ai-service`
     - etc.
   - If service names differ, update the matrix in `deploy-all.yml`

4. **Environment Variables:**
   - Services that need `DATABASE_URL` or `REDIS_URL` will fail if secrets are not set
   - The workflow sets these automatically, but they must exist in GitHub secrets

---

## ✅ Summary

**Build Workflow:** Will work once `DOCKER_HUB_USERNAME` and `DOCKER_HUB_TOKEN` are added.

**Deploy Workflow:** Will work once all Railway secrets are added and Railway services are configured.

**Current Status:** Both workflows are properly configured but **will fail** until secrets are set up.
