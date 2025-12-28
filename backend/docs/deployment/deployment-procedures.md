# Production Deployment Procedures

**Last Updated:** 2025-12-23  
**Purpose:** Step-by-step guide for deploying to production

---

## Overview

This guide provides detailed step-by-step procedures for deploying all Woragis backend services to production environments (Railway, Kubernetes, or other platforms).

---

## Pre-Deployment Checklist

### Code Quality
- [ ] All tests passing (unit, integration, security)
- [ ] Code reviewed and approved
- [ ] No breaking changes (or migration plan documented)
- [ ] Documentation updated

### Build Verification
- [ ] Docker images build successfully
- [ ] Images pushed to registry
- [ ] Image tags correct (version, not `latest`)
- [ ] Image sizes reasonable

### Configuration
- [ ] Environment variables documented
- [ ] Secrets configured
- [ ] Database migrations ready (if needed)
- [ ] Feature flags configured

### Infrastructure
- [ ] Database available and healthy
- [ ] RabbitMQ available and healthy
- [ ] Redis available and healthy
- [ ] Network connectivity verified

### Monitoring
- [ ] Health checks configured
- [ ] Metrics endpoints accessible
- [ ] Logging configured
- [ ] Alerts configured

---

## Deployment Methods

### Method 1: Railway Deployment

#### Step 1: Prepare Release

```bash
# 1. Ensure all changes are committed
git status

# 2. Run tests locally
cd backend/server/app
go test ./...

# 3. Create release tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

#### Step 2: Monitor Build

1. Go to GitHub Actions
2. Watch "Build All Services" workflow
3. Verify all services build successfully
4. Check Docker Hub for images

#### Step 3: Deploy to Railway

**Automatic (via GitHub Actions):**
- Deployment happens automatically after build
- Check Railway dashboard for status

**Manual:**
```bash
# Install Railway CLI
curl -fsSL https://railway.app/install.sh | sh

# Login
railway login

# Link to project
railway link <project-id>

# Deploy each service
railway up --service woragis-server
railway up --service woragis-ai-service
# ... repeat for all services
```

#### Step 4: Verify Deployment

```bash
# Check health
curl https://api.woragis.com/healthz

# Check logs
railway logs --service woragis-server

# Check metrics
curl https://api.woragis.com/metrics
```

---

### Method 2: Kubernetes Deployment

#### Step 1: Prepare Kubernetes Cluster

```bash
# Verify cluster access
kubectl cluster-info

# Create namespace
kubectl create namespace woragis-production

# Set context
kubectl config set-context --current --namespace=woragis-production
```

#### Step 2: Create Secrets

```bash
# Server secrets
kubectl create secret generic server-secrets \
  --from-literal=database_url="postgres://..." \
  --from-literal=redis_url="redis://..." \
  --from-literal=jwt_secret="..." \
  --from-literal=rabbitmq_url="amqp://..." \
  -n woragis-production

# Repeat for all services
```

#### Step 3: Update Image Tags

```bash
# Update deployment manifests with new image tag
cd backend/k8s
sed -i 's|image:.*server.*|image: yourregistry/woragis-server:v1.0.0|g' server/deployment.yaml
```

#### Step 4: Deploy Services

```bash
# Deploy in dependency order
kubectl apply -f ai-service/ -n woragis-production
kubectl apply -f creative-service/ -n woragis-production
kubectl apply -f docs-service/ -n woragis-production
kubectl apply -f email-worker/ -n woragis-production
kubectl apply -f translation-worker/ -n woragis-production
kubectl apply -f job-application-worker/ -n woragis-production
kubectl apply -f resume-worker/ -n woragis-production
kubectl apply -f whatsapp-worker/ -n woragis-production
kubectl apply -f server/ -n woragis-production
```

#### Step 5: Verify Deployment

```bash
# Check pods
kubectl get pods -n woragis-production

# Check services
kubectl get services -n woragis-production

# Check deployments
kubectl get deployments -n woragis-production

# Watch rollout
kubectl rollout status deployment/server -n woragis-production
```

---

## Deployment Order

Deploy services in this order to ensure dependencies are available:

1. **Database Services** (PostgreSQL, Redis, RabbitMQ)
2. **Supporting Services** (AI, Creative, Docs)
3. **Workers** (Email, Translation, Job Application, Resume, WhatsApp)
4. **Server** (depends on all above)

---

## Post-Deployment Verification

### 1. Health Checks

```bash
# Server
curl https://api.woragis.com/healthz

# AI Service
curl https://ai-service-url/health

# All services should return 200 OK
```

### 2. Functionality Tests

```bash
# Test authentication
curl -X POST https://api.woragis.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Test API endpoint
TOKEN="your-token"
curl https://api.woragis.com/api/projects \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Worker Verification

```bash
# Check worker logs
railway logs --service woragis-email-worker
# or
kubectl logs -f deployment/email-worker -n woragis-production

# Verify workers are processing jobs
# Check RabbitMQ queue depth
```

### 4. Performance Check

```bash
# Check response times
time curl https://api.woragis.com/healthz

# Check metrics
curl https://api.woragis.com/metrics | grep http_request_duration
```

---

## Rollback Procedures

### Railway Rollback

**Quick Rollback:**
```bash
railway rollback --service woragis-server
```

**Manual Rollback:**
```bash
# Deploy previous image
railway up --service woragis-server --image yourusername/woragis-server:v0.9.0
```

### Kubernetes Rollback

**Quick Rollback:**
```bash
kubectl rollout undo deployment/server -n woragis-production
kubectl rollout status deployment/server -n woragis-production
```

**Rollback to Specific Revision:**
```bash
# Check rollout history
kubectl rollout history deployment/server -n woragis-production

# Rollback to specific revision
kubectl rollout undo deployment/server --to-revision=2 -n woragis-production
```

---

## Database Migrations

### Before Deployment

```bash
# Check pending migrations
cd backend/server/app
go run cmd/migrate/main.go status

# Run migrations in staging first
go run cmd/migrate/main.go up
```

### During Deployment

**Option 1: Manual Migration**
```bash
# Run migrations before deploying new version
go run cmd/migrate/main.go up
```

**Option 2: Migration Job (Kubernetes)**
```bash
# Create migration job
kubectl create job --from=cronjob/migrate migrate-manual -n woragis-production
```

---

## Monitoring During Deployment

### Key Metrics to Watch

1. **Error Rate**
   - Monitor error rate during deployment
   - Should remain low (< 1%)

2. **Response Time**
   - Monitor p50, p95, p99 latency
   - Should not increase significantly

3. **Resource Usage**
   - CPU and memory usage
   - Should be within limits

4. **Queue Depth**
   - RabbitMQ queue depth
   - Should not increase dramatically

### Alerting

Set up alerts for:
- High error rate (> 5%)
- High response time (> 1s p95)
- Service failures
- Resource exhaustion

---

## Deployment Windows

### Recommended Times

- **Low Traffic:** Deploy during low-traffic periods
- **Maintenance Window:** For major changes, use maintenance window
- **Business Hours:** For critical fixes, deploy during business hours with team available

### Communication

- Notify team before deployment
- Post in team chat
- Update status page (if available)
- Monitor during deployment

---

## Troubleshooting

### Deployment Fails

**Symptoms:**
- Deployment times out
- Health checks failing
- Service not starting

**Steps:**
1. Check deployment logs
2. Check environment variables
3. Check resource limits
4. Check health check configuration
5. Check dependency availability

### Service Unhealthy

**Symptoms:**
- Health checks failing
- Service not responding
- Errors in logs

**Steps:**
1. Check application logs
2. Check dependency health (DB, Redis, RabbitMQ)
3. Check environment variables
4. Check resource usage
5. Rollback if necessary

### Workers Not Processing

**Symptoms:**
- Queue depth increasing
- No job processing logs
- Workers healthy but idle

**Steps:**
1. Check RabbitMQ connection
2. Check queue configuration
3. Check worker logs
4. Check message format
5. Verify workers consuming from correct queue

---

## Best Practices

### 1. Test in Staging First
- Always test in staging before production
- Use staging for integration testing
- Verify all functionality works

### 2. Use Feature Flags
- Use feature flags for gradual rollouts
- Enable features gradually
- Disable features quickly if issues

### 3. Monitor Closely
- Watch metrics during deployment
- Monitor error rates
- Check logs for issues
- Be ready to rollback

### 4. Document Changes
- Document what changed
- Document rollback procedures
- Update runbooks if procedures change

### 5. Version Control
- Use semantic versioning
- Tag all releases
- Document changes in release notes

---

## Related Documentation

- **Build Workflow:** `docs/deployment/build-workflow.md`
- **Railway Setup:** `docs/deployment/railway-setup-and-testing.md`
- **Kubernetes Deployment:** `backend/k8s/PRODUCTION_DEPLOYMENT.md`
- **Deployment Runbook:** `docs/runbooks/deploying-services.md`
- **Troubleshooting:** `docs/runbooks/troubleshooting.md`

---

**Last Updated:** 2025-12-23

