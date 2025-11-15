# WhatsApp Worker - Kubernetes Deployment Guide

## Architecture Decision: Separate Service

**Yes, the WhatsApp worker should be a separate Docker image and Kubernetes deployment** for the following reasons:

### Benefits:
1. **Independent Scaling**: Scale workers independently from the main API server
2. **Resource Isolation**: WhatsApp connection consumes resources separately
3. **Fault Isolation**: Worker failures don't affect the main API
4. **Deployment Flexibility**: Update workers without restarting the API
5. **Better Observability**: Separate metrics and logs

### Challenges with Shared Session:

WhatsApp Web API (whatsmeow) has a critical constraint:
- **Only ONE active connection per account** is allowed
- Multiple pods trying to connect simultaneously will conflict
- WhatsApp will disconnect one session when another connects

## Solution: Leader Election Pattern

We need to ensure only **one worker pod is active** at a time using:

1. **Kubernetes Leader Election**: Using Kubernetes leases API
2. **Shared Persistent Volume**: For session database (SQLite)
3. **File Locking**: Additional protection at application level

## Deployment Strategy

### Option 1: StatefulSet with Single Replica (Recommended for Production)
- Guarantees only one pod
- Automatic restart on failure
- Persistent storage
- **File**: `statefulset.yaml`

### Option 2: Deployment with Leader Election (Recommended for High Availability)
- Multiple pods, but only one active
- Automatic failover if leader dies
- More complex but more resilient
- **File**: `deployment-leader-election.yaml`

## Implementation

See the Kubernetes manifests in this directory:
- `deployment-leader-election.yaml` - Deployment with leader election
- `statefulset.yaml` - StatefulSet (single replica, simpler)
- `configmap.yaml` - Configuration (create this)
- `secrets.yaml.example` - Example secrets

## Session Storage

The WhatsApp session is stored in SQLite at `{WHATSAPP_SESSION_PATH}/whatsapp.db`.

For Kubernetes:
- Use a **PersistentVolumeClaim (PVC)** mounted at `/data/whatsapp-session`
- For StatefulSet: `ReadWriteOnce` is sufficient (single pod)
- For Deployment with leader election: `ReadWriteMany` is required (multiple pods, shared storage)
- Use file locking to prevent concurrent access

## Leader Election Implementation

The worker uses Kubernetes leader election:
- Creates a Lease resource in the specified namespace
- Only the pod holding the lease connects to WhatsApp
- Other pods wait in standby mode
- Automatic failover when leader pod dies

## Environment Variables

```bash
# Required
REDIS_URL=redis://redis-service:6379/0
WHATSAPP_SESSION_PATH=/data/whatsapp-session
WHATSAPP_ENABLED=true

# Leader Election (if enabled)
LEADER_ELECTION_ENABLED=true
LEADER_ELECTION_LEASE_NAME=whatsapp-worker-leader
LEADER_ELECTION_NAMESPACE=default
POD_NAME=<auto-set-by-k8s>

# Optional
LOG_LEVEL=info
```

## Building and Deploying

### 1. Build the Docker image:

```bash
docker build -f Dockerfile.whatsapp-worker -t your-registry/woragis-whatsapp-worker:latest .
docker push your-registry/woragis-whatsapp-worker:latest
```

### 2. Update image in Kubernetes manifests:

Edit `statefulset.yaml` or `deployment-leader-election.yaml` and replace:
```yaml
image: your-registry/woragis-whatsapp-worker:latest
```

### 3. Deploy:

**For StatefulSet (simpler):**
```bash
kubectl apply -f statefulset.yaml
```

**For Deployment with Leader Election (HA):**
```bash
kubectl apply -f deployment-leader-election.yaml
```

## Monitoring

- Health check endpoint: `/health` (you'll need to add this)
- Metrics endpoint: `/metrics` (if Prometheus enabled)
- Watch pod logs: `kubectl logs -f statefulset/whatsapp-worker` or `kubectl logs -f deployment/whatsapp-worker -l app=whatsapp-worker`

## Scaling Considerations

- **StatefulSet**: Set `replicas: 1` (only one active connection allowed)
- **Deployment with Leader Election**: Can have multiple replicas (e.g., 3), but only leader is active
- **Horizontal Pod Autoscaler**: Not recommended (WhatsApp limitation)

## Troubleshooting

1. **Multiple connections detected**: 
   - Check leader election is working: `kubectl get leases -n default`
   - Verify only one pod is leader

2. **Session conflicts**: 
   - Ensure PVC is correct access mode (ReadWriteOnce for StatefulSet, ReadWriteMany for Deployment)
   - Check file permissions on session directory

3. **QR code not accessible**: 
   - Only leader pod generates QR codes
   - Check which pod is leader: `kubectl get pods -l app=whatsapp-worker -o wide`

4. **Messages not sending**: 
   - Check leader pod logs: `kubectl logs -f <leader-pod-name>`
   - Verify WhatsApp connection status

5. **Leader election not working**:
   - Check RBAC permissions: `kubectl auth can-i create leases --as=system:serviceaccount:default:whatsapp-worker`
   - Verify service account exists

## Notes

- The worker needs access to Redis (same Redis as main API)
- Session database is SQLite, so concurrent access from multiple pods requires careful coordination
- Leader election ensures only one pod accesses the session at a time
- If leader pod dies, Kubernetes automatically elects a new leader
