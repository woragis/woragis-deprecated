# Installation Guide

## Prerequisites

1. Kubernetes cluster (1.20+)
2. kubectl configured
3. Docker registry access
4. Redis service accessible from cluster

## Step 1: Add Kubernetes Client Dependency

The worker needs `k8s.io/client-go` for leader election. Add it to `go.mod`:

```bash
cd backend/server
go get k8s.io/client-go@latest
go mod tidy
```

Then uncomment the Kubernetes imports in `app/cmd/whatsapp-worker/main.go`.

## Step 2: Build Docker Image

```bash
cd backend/server
docker build -f Dockerfile.whatsapp-worker -t your-registry/woragis-whatsapp-worker:latest .
docker push your-registry/woragis-whatsapp-worker:latest
```

## Step 3: Update Image in Manifests

Edit the Kubernetes manifests and replace:
```yaml
image: your-registry/woragis-whatsapp-worker:latest
```

## Step 4: Create ConfigMap (if needed)

```bash
kubectl create configmap whatsapp-worker-config \
  --from-literal=redis_url=redis://redis-service:6379/0 \
  --from-literal=log_level=info
```

## Step 5: Deploy

### Option A: StatefulSet (Simpler, Single Pod)

```bash
kubectl apply -f k8s/whatsapp-worker/statefulset.yaml
```

### Option B: Deployment with Leader Election (HA, Multiple Pods)

```bash
kubectl apply -f k8s/whatsapp-worker/deployment-leader-election.yaml
```

## Step 6: Verify Deployment

```bash
# Check pods
kubectl get pods -l app=whatsapp-worker

# Check logs
kubectl logs -f statefulset/whatsapp-worker
# or
kubectl logs -f deployment/whatsapp-worker -l app=whatsapp-worker

# Check leader (if using leader election)
kubectl get leases -n default
```

## Step 7: Connect WhatsApp

1. Get the pod name:
   ```bash
   kubectl get pods -l app=whatsapp-worker
   ```

2. Check logs for QR code:
   ```bash
   kubectl logs <pod-name> | grep "QR Code"
   ```

3. Or access QR code via API (if exposed):
   - The main API server should have `/api/whatsapp/qr` endpoint
   - Only the leader pod (or single pod in StatefulSet) will have the QR code

## Troubleshooting

### Pod not starting
```bash
kubectl describe pod <pod-name>
kubectl logs <pod-name>
```

### Session conflicts
- Ensure PVC is mounted correctly
- Check file permissions: `kubectl exec <pod-name> -- ls -la /data/whatsapp-session`

### Leader election issues
```bash
# Check RBAC
kubectl auth can-i create leases --as=system:serviceaccount:default:whatsapp-worker

# Check service account
kubectl get serviceaccount whatsapp-worker
```

### Redis connection issues
- Verify Redis service is accessible: `kubectl get svc redis-service`
- Test from pod: `kubectl exec <pod-name> -- nc -zv redis-service 6379`

