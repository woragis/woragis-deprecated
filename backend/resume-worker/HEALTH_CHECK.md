# Health Check - Resume Worker

## Overview

The Resume Worker provides a health check endpoint at `GET /healthz` on port 8080 that reports the status of RabbitMQ connection.

## Endpoint

**URL:** `GET http://localhost:8080/healthz`  
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
    },
    {
      "name": "rabbitmq",
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
      "status": "ok"
    },
    {
      "name": "rabbitmq",
      "status": "error",
      "message": "connection closed"
    }
  ]
}
```

## Health Checks

### Service Check
- **Name:** `service`
- **Critical:** Yes
- **Check:** Service availability
- **Failure Impact:** Worker cannot function

### RabbitMQ Check
- **Name:** `rabbitmq`
- **Critical:** Yes
- **Check:** RabbitMQ connection status
- **Failure Impact:** Worker cannot consume resume generation jobs

## Performance

- **Caching:** Results are cached for 5 seconds
- **Lightweight:** Fast, non-blocking checks

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

## Usage

```bash
curl http://localhost:8080/healthz
```

## Implementation

- **HTTP Server:** Python's HTTPServer on port 8080 (daemon thread)
- **Thread-safe:** Uses context variables for caching
- **Non-blocking:** Health checks don't interfere with message processing
