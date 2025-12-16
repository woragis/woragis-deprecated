# Health Check - Job Application Worker

## Overview

The Job Application Worker provides a health check endpoint at `GET /healthz` on port 8080 that reports the status of RabbitMQ connection.

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
      "message": "connection not initialized"
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
- **Failure Impact:** Worker cannot consume job application jobs

## Performance

- **Caching:** Results are cached for 5 seconds (via module-level variables)
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

- **HTTP Server:** Node.js http.createServer on port 8080
- **Non-blocking:** Health checks don't interfere with message processing
- **Graceful Shutdown:** Server closes with worker shutdown
