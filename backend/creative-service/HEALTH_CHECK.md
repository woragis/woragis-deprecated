# Health Check - Creative Service

## Overview

The Creative Service provides a health check endpoint at `GET /healthz` that reports service availability.

## Endpoint

**URL:** `GET /healthz`  
**Authentication:** Not required (public endpoint)  
**Response Format:** JSON

## Response Format

### Healthy Response (200 OK)
```json
{
  "status": "healthy",
  "checks": [
    {
      "name": "service",
      "status": "ok"
    }
  ]
}
```

### Unhealthy Response (503 Service Unavailable)
```json
{
  "status": "unhealthy",
  "checks": [
    {
      "name": "service",
      "status": "error",
      "message": "service unavailable"
    }
  ]
}
```

## Health Checks

### Service Check
- **Name:** `service`
- **Critical:** Yes
- **Check:** Service availability (always ok if endpoint is reachable)
- **Failure Impact:** Service cannot function

## Performance

- **Caching:** Results are cached for 5 seconds
- **Lightweight:** Fast, non-blocking checks

## Kubernetes Integration

### Liveness Probe
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8000
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

### Readiness Probe
```yaml
readinessProbe:
  httpGet:
    path: /healthz
    port: 8000
  initialDelaySeconds: 10
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

## Usage

```bash
curl http://localhost:8000/healthz
```
