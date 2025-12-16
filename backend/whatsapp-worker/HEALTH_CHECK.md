# Health Check - WhatsApp Worker

## Overview

The WhatsApp Worker provides a health check endpoint at `GET /healthz` on port 8080 that reports the status of RabbitMQ connection.

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
      "name": "rabbitmq",
      "status": "error",
      "message": "connection closed"
    }
  ]
}
```

## Health Checks

### RabbitMQ Check
- **Name:** `rabbitmq`
- **Critical:** Yes
- **Check:** RabbitMQ connection status (2s timeout)
- **Failure Impact:** Worker cannot process WhatsApp messages

## Performance

- **Caching:** Results are cached for 5 seconds
- **Timeout:** 2 seconds per check
- **Lightweight:** Non-blocking connection status check

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

- **HTTP Server:** Runs on port 8080 alongside the worker
- **Graceful Shutdown:** Health server shuts down with the worker
- **Thread-safe:** Uses mutex for concurrent access
