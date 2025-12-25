# Kubernetes Deployment Manifests

This directory contains Kubernetes manifests for deploying all Woragis backend services.

## Structure

```
k8s/
├── server/              # Main API server
├── email-worker/        # Email processing worker
├── translation-worker/   # Translation processing worker
├── job-application-worker/  # Job application automation worker
├── resume-worker/       # Resume generation worker
├── whatsapp-worker/     # WhatsApp messaging worker (already exists)
├── ai-service/          # AI/LLM service
├── creative-service/    # Creative content generation service
└── docs-service/        # Documentation service
```

## Prerequisites

1. **Kubernetes cluster** (v1.24+)
2. **kubectl** configured to access your cluster
3. **Docker images** built and pushed to your registry
4. **Secrets** created (see below)

## Quick Start

### 1. Create Namespace (Optional)

```bash
kubectl create namespace woragis
```

### 2. Create Secrets

You need to create secrets for each service. Example for server:

```bash
kubectl create secret generic server-secrets \
  --from-literal=database_url="postgres://user:pass@host:5432/db" \
  --from-literal=redis_url="redis://host:6379/0" \
  --from-literal=jwt_secret="your-jwt-secret" \
  --from-literal=rabbitmq_url="amqp://user:pass@host:5672/vhost"
```

Repeat for all services that need secrets:
- `server-secrets`
- `email-worker-secrets`
- `translation-worker-secrets`
- `job-application-worker-secrets`
- `resume-worker-secrets`
- `ai-service-secrets`
- `creative-service-secrets`

### 3. Deploy Services

Deploy in order (dependencies first):

```bash
# Deploy services (AI, Creative, Docs)
kubectl apply -f ai-service/
kubectl apply -f creative-service/
kubectl apply -f docs-service/

# Deploy workers
kubectl apply -f email-worker/
kubectl apply -f translation-worker/
kubectl apply -f job-application-worker/
kubectl apply -f resume-worker/
kubectl apply -f whatsapp-worker/

# Deploy server (depends on services)
kubectl apply -f server/
```

Or deploy all at once:

```bash
kubectl apply -f .
```

### 4. Verify Deployment

```bash
# Check pods
kubectl get pods

# Check services
kubectl get services

# Check deployments
kubectl get deployments
```

## Configuration

### ConfigMaps

ConfigMaps contain non-sensitive configuration:
- `server-config` - Server configuration
- `email-worker-config` - Email worker configuration
- `translation-worker-config` - Translation worker configuration
- `job-application-worker-config` - Job application worker configuration
- `resume-worker-config` - Resume worker configuration

### Secrets

Secrets contain sensitive data (API keys, passwords, connection strings):
- `server-secrets` - Database, Redis, RabbitMQ, JWT secret
- `email-worker-secrets` - RabbitMQ, SMTP credentials
- `translation-worker-secrets` - RabbitMQ, Database, API keys
- `job-application-worker-secrets` - RabbitMQ, Database
- `resume-worker-secrets` - RabbitMQ, Database
- `ai-service-secrets` - AI provider API keys
- `creative-service-secrets` - Creative service API keys

## Resource Limits

### Server
- Requests: 500m CPU, 512Mi memory
- Limits: 1000m CPU, 1Gi memory
- Replicas: 2-5 (HPA configured)

### Go Workers (Email, Translation, WhatsApp)
- Requests: 200m CPU, 256Mi memory
- Limits: 500m CPU, 512Mi memory
- Replicas: 1-2

### Node.js/Python Workers (Job Application, Resume)
- Requests: 300m CPU, 512Mi memory
- Limits: 1000m CPU, 1Gi memory
- Replicas: 1-2

### Services (AI, Creative, Docs)
- Requests: 300m CPU, 512Mi memory
- Limits: 1000m CPU, 1Gi memory
- Replicas: 1-2

## Health Checks

All services have:
- **Liveness Probe**: Restarts container if unhealthy
- **Readiness Probe**: Removes from service endpoints if not ready
- **Startup Probe**: Allows time for initialization

Health check endpoints:
- Server: `/healthz`
- Workers: `/healthz` or `/health`
- Services: `/health`

## Scaling

### Horizontal Pod Autoscaler (HPA)

Server has HPA configured:
- Min replicas: 2
- Max replicas: 5
- CPU target: 70%
- Memory target: 80%

Workers can be scaled manually based on queue depth:

```bash
kubectl scale deployment email-worker --replicas=3
```

## Ingress

Server has Ingress configured for external access:
- Host: `api.woragis.com`
- TLS: Managed by cert-manager
- Ingress class: `nginx`

## Monitoring

All services expose metrics on port `8081`:
- Prometheus can scrape from `/metrics` endpoint
- Service discovery via Kubernetes annotations

## Troubleshooting

### Check Pod Logs

```bash
kubectl logs -f deployment/server
kubectl logs -f deployment/email-worker
```

### Check Pod Status

```bash
kubectl describe pod <pod-name>
```

### Check Service Endpoints

```bash
kubectl get endpoints
```

### Port Forward for Local Testing

```bash
# Server
kubectl port-forward service/server 8080:80

# AI Service
kubectl port-forward service/ai-service 8000:8000
```

## Updating Deployments

### Update Image

```bash
kubectl set image deployment/server server=woragis/server:v0.0.2
```

### Rolling Update

Deployments use `RollingUpdate` strategy by default:
- `maxSurge: 1` - Can create 1 extra pod during update
- `maxUnavailable: 0` - No pods unavailable during update

### Rollback

```bash
kubectl rollout undo deployment/server
```

## Production Considerations

1. **Use versioned image tags** (not `latest`)
2. **Set up proper secrets management** (external-secrets-operator, sealed-secrets)
3. **Configure resource quotas** per namespace
4. **Set up network policies** for security
5. **Enable pod security standards**
6. **Configure backup for persistent volumes** (WhatsApp worker)
7. **Set up monitoring and alerting**
8. **Use separate namespaces** for dev/staging/prod

## Local Development

For local testing, use `minikube` or `kind`:

```bash
# Start minikube
minikube start

# Enable ingress
minikube addons enable ingress

# Deploy
kubectl apply -f .
```

## Next Steps

1. Create secrets for all services
2. Update image tags to specific versions
3. Configure Ingress for your domain
4. Set up TLS certificates (cert-manager)
5. Configure monitoring (Prometheus, Grafana)
6. Set up CI/CD for automated deployments

