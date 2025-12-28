# Deploying Services and Workers

## Overview
This runbook covers the procedures for deploying backend services and workers to production, staging, and development environments.

## Prerequisites
- Access to deployment platform (Railway, Kubernetes, etc.)
- Docker images built and pushed to registry
- Environment variables configured
- Database migrations run (if needed)
- Health checks passing

## Pre-Deployment Checklist

### 1. Code Review
- [ ] Code reviewed and approved
- [ ] Tests passing (unit, integration)
- [ ] No breaking changes (or migration plan in place)
- [ ] Documentation updated

### 2. Build Verification
- [ ] Docker images build successfully
- [ ] Docker images pushed to registry
- [ ] Image tags correct (version, commit SHA)
- [ ] Image sizes reasonable

### 3. Configuration
- [ ] Environment variables configured
- [ ] Secrets configured (API keys, passwords)
- [ ] Database migrations ready (if needed)
- [ ] Feature flags configured (if applicable)

### 4. Infrastructure
- [ ] Database available and healthy
- [ ] RabbitMQ available and healthy
- [ ] Redis available and healthy
- [ ] Network connectivity verified

### 5. Monitoring
- [ ] Health checks configured
- [ ] Metrics endpoints accessible
- [ ] Logging configured
- [ ] Alerts configured

## Deployment Procedures

### Server Deployment

#### 1. Build Docker Image
```bash
cd backend/server
docker build -t woragis/server:latest -t woragis/server:$(git rev-parse --short HEAD) .
```

#### 2. Push to Registry
```bash
docker push woragis/server:latest
docker push woragis/server:$(git rev-parse --short HEAD)
```

#### 3. Deploy to Railway
```bash
# Via Railway CLI
railway up --service server

# Or via Railway Dashboard
# 1. Go to Railway dashboard
# 2. Select server service
# 3. Click "Deploy"
# 4. Select image tag
```

#### 4. Verify Deployment
```bash
# Check health
curl https://api.woragis.com/healthz

# Check logs
railway logs --service server

# Check metrics
curl https://api.woragis.com/metrics
```

### Worker Deployment

#### 1. Build Docker Image
```bash
cd backend/{worker-name}
docker build -t woragis/{worker-name}:latest -t woragis/{worker-name}:$(git rev-parse --short HEAD) .
```

#### 2. Push to Registry
```bash
docker push woragis/{worker-name}:latest
docker push woragis/{worker-name}:$(git rev-parse --short HEAD)
```

#### 3. Deploy to Railway/Kubernetes
```bash
# Railway
railway up --service {worker-name}

# Kubernetes
kubectl set image deployment/{worker-name} {worker-name}=woragis/{worker-name}:$(git rev-parse --short HEAD)
kubectl rollout status deployment/{worker-name}
```

#### 4. Verify Deployment
```bash
# Check health
curl http://{worker-name}:8080/healthz

# Check logs
railway logs --service {worker-name}
# or
kubectl logs -f deployment/{worker-name}

# Check metrics
curl http://{worker-name}:8080/metrics
```

### Service Deployment (AI, Creative)

#### 1. Build Docker Image
```bash
cd backend/{service-name}
docker build -t woragis/{service-name}:latest -t woragis/{service-name}:$(git rev-parse --short HEAD) .
```

#### 2. Push to Registry
```bash
docker push woragis/{service-name}:latest
docker push woragis/{service-name}:$(git rev-parse --short HEAD)
```

#### 3. Deploy
```bash
# Railway
railway up --service {service-name}

# Kubernetes
kubectl set image deployment/{service-name} {service-name}=woragis/{service-name}:$(git rev-parse --short HEAD)
kubectl rollout status deployment/{service-name}
```

#### 4. Verify Deployment
```bash
# Check health
curl http://{service-name}:8000/healthz

# Check logs
railway logs --service {service-name}
# or
kubectl logs -f deployment/{service-name}

# Check metrics
curl http://{service-name}:8000/metrics
```

## Deployment Strategies

### Rolling Deployment (Default)
- Deploy new version gradually
- Old version remains running during deployment
- Zero downtime (if configured correctly)
- Automatic rollback on failure

### Blue-Green Deployment
- Deploy new version alongside old version
- Switch traffic to new version
- Keep old version for quick rollback
- Requires more resources

### Canary Deployment
- Deploy new version to small subset
- Monitor metrics and errors
- Gradually increase traffic
- Rollback if issues detected

## Environment-Specific Procedures

### Development Environment

#### Local Docker Compose
```bash
cd backend
docker-compose up -d {service-name}
docker-compose logs -f {service-name}
```

#### Verify
```bash
# Check health
curl http://localhost:8080/healthz

# Check logs
docker-compose logs {service-name}
```

### Staging Environment

#### Deploy
```bash
# Set environment
export RAILWAY_ENVIRONMENT=staging

# Deploy
railway up --service {service-name}
```

#### Verify
```bash
# Check health
curl https://staging-api.woragis.com/healthz

# Run smoke tests
# (if available)
```

### Production Environment

#### Deploy (with caution!)
```bash
# Set environment
export RAILWAY_ENVIRONMENT=production

# Deploy during maintenance window (if needed)
railway up --service {service-name}
```

#### Verify
```bash
# Check health
curl https://api.woragis.com/healthz

# Monitor metrics
# Check error rates
# Check latency
```

## Rollback Procedures

### Quick Rollback (Railway)
```bash
# Rollback to previous deployment
railway rollback --service {service-name}
```

### Rollback (Kubernetes)
```bash
# Rollback to previous revision
kubectl rollout undo deployment/{service-name}
kubectl rollout status deployment/{service-name}
```

### Manual Rollback
```bash
# Deploy previous image tag
railway up --service {service-name} --image woragis/{service-name}:{previous-tag}
# or
kubectl set image deployment/{service-name} {service-name}=woragis/{service-name}:{previous-tag}
```

## Post-Deployment Verification

### 1. Health Checks
- [ ] All health checks passing
- [ ] No degraded status
- [ ] Dependencies healthy

### 2. Functionality
- [ ] Key endpoints working
- [ ] Workers processing jobs
- [ ] Services responding correctly

### 3. Performance
- [ ] Response times normal
- [ ] Error rates normal
- [ ] Resource usage normal

### 4. Monitoring
- [ ] Metrics being collected
- [ ] Logs being generated
- [ ] Alerts configured correctly

## Common Issues

### Deployment Fails

#### Symptoms
- Deployment times out
- Health checks failing
- Service not starting

#### Steps
1. Check deployment logs
2. Check environment variables
3. Check resource limits
4. Check health check configuration
5. Check dependency availability

### Service Unhealthy After Deployment

#### Symptoms
- Health checks failing
- Service not responding
- Errors in logs

#### Steps
1. Check application logs
2. Check dependency health (DB, Redis, RabbitMQ)
3. Check environment variables
4. Check resource usage (CPU, memory)
5. Rollback if necessary

### Workers Not Processing Jobs

#### Symptoms
- Queue depth increasing
- No job processing logs
- Workers healthy but idle

#### Steps
1. Check RabbitMQ connection
2. Check queue configuration
3. Check worker logs
4. Check message format
5. Verify workers are consuming from correct queue

## Best Practices

### 1. Deploy During Low Traffic
- Schedule deployments during low-traffic periods
- Use maintenance windows for major changes
- Notify users of planned downtime (if needed)

### 2. Monitor Closely
- Watch metrics during deployment
- Monitor error rates
- Check logs for issues
- Be ready to rollback

### 3. Test in Staging First
- Always test in staging before production
- Use staging for integration testing
- Verify all functionality works

### 4. Use Feature Flags
- Use feature flags for gradual rollouts
- Enable features gradually
- Disable features quickly if issues

### 5. Document Changes
- Document what changed
- Document rollback procedures
- Update runbooks if procedures change

## Automation

### CI/CD Pipeline
- Automated builds on push
- Automated tests
- Automated deployment to staging
- Manual approval for production

### Deployment Scripts
```bash
#!/bin/bash
# deploy.sh

SERVICE=$1
VERSION=$(git rev-parse --short HEAD)

# Build
docker build -t woragis/$SERVICE:$VERSION .

# Push
docker push woragis/$SERVICE:$VERSION

# Deploy
railway up --service $SERVICE --image woragis/$SERVICE:$VERSION

# Verify
sleep 10
curl -f http://$SERVICE:8080/healthz || exit 1
```

## Related Documentation
- [Architecture Decision Records](../adr/) - Architecture decisions
- [Troubleshooting Guide](./troubleshooting.md) - Troubleshooting procedures
- [Monitoring DLQ](./monitoring-dlq.md) - DLQ monitoring
- [Deployment Procedures](../deployment/deployment-procedures.md) - Detailed deployment guide
- [Railway Setup](../deployment/railway-setup-and-testing.md) - Railway deployment guide
- [Build Workflow](../deployment/build-workflow.md) - CI/CD build process
