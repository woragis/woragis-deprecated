# Health Check - Server

## Overview

The server provides a health check endpoint at `GET /healthz` that reports the status of critical dependencies.

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
      "name": "database",
      "status": "ok"
    },
    {
      "name": "redis",
      "status": "ok"
    },
    {
      "name": "rabbitmq",
      "status": "ok"
    }
  ]
}
```

### Degraded Response (200 OK)
```json
{
  "status": "degraded",
  "checks": [
    {
      "name": "database",
      "status": "ok"
    },
    {
      "name": "redis",
      "status": "ok"
    },
    {
      "name": "rabbitmq",
      "status": "error",
      "message": "not connected"
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
      "name": "database",
      "status": "error",
      "message": "connection timeout"
    },
    {
      "name": "redis",
      "status": "ok"
    },
    {
      "name": "rabbitmq",
      "status": "ok"
    }
  ]
}
```

## Status Values

- **healthy:** All critical dependencies are operational
- **degraded:** Some non-critical dependencies failed (e.g., RabbitMQ optional)
- **unhealthy:** Critical dependencies failed (e.g., database, Redis)

## Health Checks

### Database Check
- **Name:** `database`
- **Critical:** Yes
- **Check:** PostgreSQL connection ping with 2s timeout
- **Failure Impact:** Server cannot function without database

### Redis Check
- **Name:** `redis`
- **Critical:** Yes
- **Check:** Redis connection ping with 2s timeout
- **Failure Impact:** Server cannot function without Redis

### RabbitMQ Check
- **Name:** `rabbitmq`
- **Critical:** No (optional)
- **Check:** Connection status check
- **Failure Impact:** Notifications fall back to Redis-only mode

## Performance

- **Caching:** Results are cached for 5 seconds to avoid excessive load
- **Timeout:** Each check has a 2-second timeout
- **Lightweight:** Checks are non-blocking and fast

## Kubernetes Integration

### Liveness Probe
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
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
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

### Startup Probe
```yaml
startupProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 30
```

## Usage

### Manual Check
```bash
curl http://localhost:8080/healthz
```

### Monitoring Integration
- Use HTTP status code for alerting (200 = healthy, 503 = unhealthy)
- Parse JSON response for detailed component status
- Monitor `status` field for overall health state

## Implementation Details

- **Package:** `app/pkg/health`
- **Caching:** 5-second TTL to reduce database/Redis load
- **Thread-safe:** Uses mutex for concurrent access
- **Non-blocking:** All checks use timeouts to prevent hanging
