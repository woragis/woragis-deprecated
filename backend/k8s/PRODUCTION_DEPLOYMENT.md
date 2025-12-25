# Kubernetes Production Deployment Guide

**Date:** 2025-12-23  
**Status:** Ready for Production Deployment

---

## Overview

This guide covers deploying all Woragis backend services to a Kubernetes cluster with production-ready configurations including TLS, monitoring, and auto-scaling.

---

## Prerequisites

### 1. Kubernetes Cluster
- **Version:** v1.24 or higher
- **Options:**
  - GKE (Google Kubernetes Engine)
  - EKS (Amazon Elastic Kubernetes Service)
  - AKS (Azure Kubernetes Service)
  - Self-hosted cluster
  - Local: minikube or kind (for testing)

### 2. Required Tools
- `kubectl` configured to access your cluster
- `helm` (for cert-manager and monitoring)
- Docker images built and pushed to registry

### 3. External Services
- PostgreSQL database (managed or self-hosted)
- Redis (managed or self-hosted)
- RabbitMQ (managed or self-hosted)
- Domain name for Ingress

---

## Step 1: Prepare Cluster

### 1.1 Create Namespace

```bash
kubectl create namespace woragis-production
kubectl config set-context --current --namespace=woragis-production
```

### 1.2 Set Up Secrets Management (Optional but Recommended)

**Option A: External Secrets Operator**
```bash
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets external-secrets/external-secrets -n external-secrets-system --create-namespace
```

**Option B: Sealed Secrets**
```bash
helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets
helm install sealed-secrets sealed-secrets/sealed-secrets
```

**Option C: Manual Secrets (for testing)**
See Step 2 below.

---

## Step 2: Create Secrets

### 2.1 Server Secrets

```bash
kubectl create secret generic server-secrets \
  --from-literal=database_url="postgres://user:pass@host:5432/woragis" \
  --from-literal=redis_url="redis://host:6379/0" \
  --from-literal=jwt_secret="$(openssl rand -base64 32)" \
  --from-literal=rabbitmq_url="amqp://user:pass@host:5672/vhost" \
  -n woragis-production
```

### 2.2 Email Worker Secrets

```bash
kubectl create secret generic email-worker-secrets \
  --from-literal=rabbitmq_url="amqp://user:pass@host:5672/vhost" \
  --from-literal=smtp_host="smtp.example.com" \
  --from-literal=smtp_port="587" \
  --from-literal=smtp_username="user" \
  --from-literal=smtp_password="pass" \
  -n woragis-production
```

### 2.3 Translation Worker Secrets

```bash
kubectl create secret generic translation-worker-secrets \
  --from-literal=rabbitmq_url="amqp://user:pass@host:5672/vhost" \
  --from-literal=database_url="postgres://user:pass@host:5432/woragis" \
  --from-literal=google_translate_api_key="your-key" \
  --from-literal=deepl_api_key="your-key" \
  -n woragis-production
```

### 2.4 Job Application Worker Secrets

```bash
kubectl create secret generic job-application-worker-secrets \
  --from-literal=rabbitmq_url="amqp://user:pass@host:5672/vhost" \
  --from-literal=database_url="postgres://user:pass@host:5432/woragis" \
  -n woragis-production
```

### 2.5 Resume Worker Secrets

```bash
kubectl create secret generic resume-worker-secrets \
  --from-literal=rabbitmq_url="amqp://user:pass@host:5672/vhost" \
  --from-literal=database_url="postgres://user:pass@host:5432/woragis" \
  -n woragis-production
```

### 2.6 AI Service Secrets

```bash
kubectl create secret generic ai-service-secrets \
  --from-literal=openai_api_key="sk-..." \
  --from-literal=anthropic_api_key="sk-ant-..." \
  --from-literal=grok_api_key="..." \
  --from-literal=manus_api_key="..." \
  --from-literal=cipher_api_key="..." \
  -n woragis-production
```

### 2.7 Creative Service Secrets

```bash
kubectl create secret generic creative-service-secrets \
  --from-literal=openai_api_key="sk-..." \
  --from-literal=stability_ai_api_key="..." \
  -n woragis-production
```

### 2.8 Verify Secrets

```bash
kubectl get secrets -n woragis-production
```

---

## Step 3: Update Image Tags

Before deploying, update image tags in deployment manifests to use specific versions (not `latest`):

```bash
# Example: Update server deployment
sed -i 's|image:.*server.*|image: your-registry/woragis-server:v1.0.0|g' server/deployment.yaml
```

Or use environment variables in your CI/CD pipeline.

---

## Step 4: Deploy Services

### 4.1 Deploy Using Script

```bash
cd backend/k8s
NAMESPACE=woragis-production ./deploy.sh
```

### 4.2 Deploy Manually

Deploy in dependency order:

```bash
# Services first
kubectl apply -f ai-service/ -n woragis-production
kubectl apply -f creative-service/ -n woragis-production
kubectl apply -f docs-service/ -n woragis-production

# Workers
kubectl apply -f email-worker/ -n woragis-production
kubectl apply -f translation-worker/ -n woragis-production
kubectl apply -f job-application-worker/ -n woragis-production
kubectl apply -f resume-worker/ -n woragis-production
kubectl apply -f whatsapp-worker/ -n woragis-production

# Server last
kubectl apply -f server/ -n woragis-production
```

### 4.3 Verify Deployment

```bash
# Check pods
kubectl get pods -n woragis-production

# Check services
kubectl get services -n woragis-production

# Check deployments
kubectl get deployments -n woragis-production

# Watch pod status
kubectl get pods -n woragis-production -w
```

---

## Step 5: Configure Ingress with TLS

### 5.1 Install Ingress Controller

**NGINX Ingress:**
```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace
```

### 5.2 Install cert-manager

```bash
helm repo add jetstack https://charts.jetstack.io
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true
```

### 5.3 Create ClusterIssuer (Let's Encrypt)

```yaml
# cert-manager/cluster-issuer.yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
```

```bash
kubectl apply -f cert-manager/cluster-issuer.yaml
```

### 5.4 Update Ingress with TLS

Edit `server/ingress.yaml`:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: server-ingress
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - api.woragis.com
    secretName: server-tls
  rules:
  - host: api.woragis.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: server
            port:
              number: 80
```

Apply:
```bash
kubectl apply -f server/ingress.yaml -n woragis-production
```

### 5.5 Verify TLS Certificate

```bash
# Check certificate status
kubectl get certificate -n woragis-production

# Check certificate details
kubectl describe certificate server-tls -n woragis-production

# Test HTTPS
curl -I https://api.woragis.com/healthz
```

---

## Step 6: Set Up Monitoring

### 6.1 Install Prometheus Operator

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace
```

### 6.2 Configure ServiceMonitor

Create `server/servicemonitor.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: server-metrics
  labels:
    app: server
spec:
  selector:
    matchLabels:
      app: server
  endpoints:
  - port: metrics
    path: /metrics
    interval: 30s
```

Apply:
```bash
kubectl apply -f server/servicemonitor.yaml -n woragis-production
```

### 6.3 Access Grafana

```bash
# Get Grafana admin password
kubectl get secret --namespace monitoring kube-prometheus-stack-grafana -o jsonpath="{.data.admin-password}" | base64 --decode

# Port forward
kubectl port-forward -n monitoring service/kube-prometheus-stack-grafana 3000:80
```

Access: http://localhost:3000 (admin / password from above)

### 6.4 Create Dashboards

Import dashboards for:
- Kubernetes cluster metrics
- Application metrics (Prometheus)
- Service health and performance

---

## Step 7: Test Auto-Scaling

### 7.1 Verify HPA

```bash
# Check HPA status
kubectl get hpa -n woragis-production

# Describe HPA
kubectl describe hpa server-hpa -n woragis-production
```

### 7.2 Generate Load

```bash
# Install hey (load testing tool)
go install github.com/rakyll/hey@latest

# Generate load
hey -n 10000 -c 50 https://api.woragis.com/healthz
```

### 7.3 Watch Scaling

```bash
# Watch pods scaling
kubectl get pods -n woragis-production -w

# Watch HPA
kubectl get hpa -n woragis-production -w
```

### 7.4 Verify Scaling

HPA should scale server pods from 2 to 5 based on CPU/memory usage.

---

## Step 8: Test Rolling Updates

### 8.1 Update Deployment

```bash
# Update image
kubectl set image deployment/server server=your-registry/woragis-server:v1.0.1 -n woragis-production

# Watch rollout
kubectl rollout status deployment/server -n woragis-production
```

### 8.2 Verify Rolling Update

```bash
# Check rollout history
kubectl rollout history deployment/server -n woragis-production

# Check current revision
kubectl get deployment server -n woragis-production -o jsonpath='{.metadata.annotations.deployment\.kubernetes\.io/revision}'
```

### 8.3 Rollback if Needed

```bash
# Rollback to previous version
kubectl rollout undo deployment/server -n woragis-production

# Rollback to specific revision
kubectl rollout undo deployment/server --to-revision=2 -n woragis-production
```

---

## Step 9: Production Hardening

### 9.1 Network Policies

Create network policies to restrict pod-to-pod communication:

```yaml
# network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: server-network-policy
spec:
  podSelector:
    matchLabels:
      app: server
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 80
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: database
    ports:
    - protocol: TCP
      port: 5432
```

### 9.2 Pod Security Standards

Enable Pod Security Standards:

```yaml
# namespace-pod-security.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: woragis-production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### 9.3 Resource Quotas

Create resource quotas:

```yaml
# resource-quota.yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: woragis-quota
spec:
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
    persistentvolumeclaims: "10"
    pods: "50"
```

### 9.4 Limit Ranges

Create limit ranges:

```yaml
# limit-range.yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: woragis-limits
spec:
  limits:
  - default:
      cpu: "500m"
      memory: "512Mi"
    defaultRequest:
      cpu: "200m"
      memory: "256Mi"
    type: Container
```

---

## Step 10: Backup and Disaster Recovery

### 10.1 Database Backups

Set up automated database backups:

```bash
# Create backup job
kubectl create job --from=cronjob/database-backup database-backup-manual -n woragis-production
```

### 10.2 Persistent Volume Backups

For WhatsApp worker (StatefulSet with persistent volumes):

```bash
# Install Velero for backup
helm repo add vmware-tanzu https://vmware-tanzu.github.io/helm-charts
helm install velero vmware-tanzu/velero \
  --namespace velero \
  --create-namespace \
  --set configuration.provider=aws \
  --set configuration.backupStorageLocation.bucket=your-backup-bucket
```

---

## Troubleshooting

### Pods Not Starting

```bash
# Check pod events
kubectl describe pod <pod-name> -n woragis-production

# Check logs
kubectl logs <pod-name> -n woragis-production

# Check previous container logs
kubectl logs <pod-name> -n woragis-production --previous
```

### Services Not Accessible

```bash
# Check service endpoints
kubectl get endpoints -n woragis-production

# Check service details
kubectl describe service <service-name> -n woragis-production

# Test service from within cluster
kubectl run -it --rm debug --image=busybox --restart=Never -- sh
# Then: wget -O- http://server:80/healthz
```

### Ingress Not Working

```bash
# Check ingress status
kubectl describe ingress server-ingress -n woragis-production

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

### Certificate Issues

```bash
# Check certificate status
kubectl get certificate -n woragis-production

# Check certificate request
kubectl get certificaterequest -n woragis-production

# Check cert-manager logs
kubectl logs -n cert-manager deployment/cert-manager
```

---

## Monitoring and Alerting

### Key Metrics to Monitor

1. **Pod Health**
   - Pod restarts
   - Pod status (Running/Error/CrashLoopBackOff)

2. **Resource Usage**
   - CPU usage per pod
   - Memory usage per pod
   - Disk usage

3. **Application Metrics**
   - Request rate
   - Error rate
   - Response time (p50, p95, p99)

4. **Queue Depth**
   - RabbitMQ queue depth
   - Worker processing rate

### Alerting Rules

Create Prometheus alerting rules for:
- High error rate (> 5%)
- High response time (> 1s p95)
- Pod crashes
- Resource exhaustion
- Queue backlog

---

## Next Steps

1. ✅ **Kubernetes Deployment** - Complete
2. **CI/CD Integration** - Set up automated deployments
3. **Performance Testing** - Load test the cluster
4. **Cost Optimization** - Right-size resources
5. **Multi-Region Deployment** - Deploy to multiple regions

---

## Related Documents

- **Kubernetes Manifests:** `backend/k8s/`
- **Deployment Script:** `backend/k8s/deploy.sh`
- **Secrets Template:** `backend/k8s/secrets-template.yaml`
- **Main README:** `backend/k8s/README.md`

---

**Last Updated:** 2025-12-23

