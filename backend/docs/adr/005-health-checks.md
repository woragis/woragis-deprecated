# ADR-005: Health Checks Implementation Strategy

## Status
Accepted

## Context
The backend consists of multiple services and workers that need to report their health status. Health checks are used by:
- Kubernetes (for pod restart decisions)
- Load balancers (for routing decisions)
- Monitoring systems (for alerting)
- Deployment tools (for deployment verification)

**Requirements:**
- Consistent health check pattern across all components
- Dependency checking (database, Redis, RabbitMQ)
- Fast response times (for frequent polling)
- Clear status indication (healthy, degraded, unhealthy)
- No authentication required (internal use only)

**Constraints:**
- Different components have different dependencies
- Health checks should not impact performance
- Need to distinguish between critical and non-critical failures
- Must work with Kubernetes liveness/readiness probes

## Decision
We will implement **consistent health check endpoints** (`GET /healthz`) across all components with:
- **JSON response format** with `status` and `checks` array
- **Caching** (5 seconds) to reduce load on dependencies
- **Timeout** (2 seconds) for all dependency checks
- **HTTP status codes**: 200 (healthy/degraded), 503 (unhealthy)
- **Dependency-specific checks** per component

**Status Values:**
- `healthy`: All critical dependencies OK
- `degraded`: Non-critical dependencies failed (service operational)
- `unhealthy`: Critical dependencies failed (service cannot function)

## Consequences

### Positive
- ✅ **Consistency**: Same pattern across all components
- ✅ **Kubernetes Compatible**: Works with liveness/readiness probes
- ✅ **Performance**: Caching reduces load on dependencies
- ✅ **Clear Status**: Distinguishes between healthy, degraded, unhealthy
- ✅ **Fast**: Timeouts prevent hanging checks

### Negative
- ⚠️ **Caching Trade-off**: 5-second cache means status may be slightly stale
- ⚠️ **Code Duplication**: Similar health check code across components
- ⚠️ **Dependency Coupling**: Health checks depend on dependency availability

### Neutral
- Health checks are lightweight and don't impact normal operations
- Different components check different dependencies (appropriate)

## Implementation Details

### Response Format
```json
{
  "status": "healthy|degraded|unhealthy",
  "checks": [
    {"name": "database", "status": "ok|error"},
    {"name": "redis", "status": "ok|error"},
    {"name": "rabbitmq", "status": "ok|error"}
  ]
}
```

### HTTP Status Codes
- **200 OK**: Service is operational (healthy or degraded)
- **503 Service Unavailable**: Service cannot function (unhealthy)

### Component-Specific Checks

#### Server
- **Database (PostgreSQL)**: CRITICAL - Connection ping
- **Redis**: CRITICAL - Connection ping
- **RabbitMQ**: OPTIONAL - Connection status (degraded if failed)

#### Go Workers (Email, WhatsApp, Translation)
- **RabbitMQ**: CRITICAL - Connection status

#### Python Workers (Resume)
- **RabbitMQ**: CRITICAL - Connection status
- **Database**: OPTIONAL - Connection ping (if needed)

#### Node.js Workers (Job Application)
- **RabbitMQ**: CRITICAL - Connection status
- **Database**: OPTIONAL - Connection ping (if needed)

#### Python Services (AI, Creative)
- **External APIs**: OPTIONAL - Basic connectivity (if needed)
- **No critical dependencies**: Service is healthy if running

### Caching Strategy
- Health check results cached for **5 seconds**
- Cache key: component name + timestamp
- Cache invalidated after 5 seconds
- Reduces load on dependencies (especially database)

### Timeout Strategy
- All dependency checks use **2-second timeout**
- Prevents health checks from hanging
- Fast failure for unresponsive dependencies

## Implementation Examples

### Go Server
```go
type HealthChecker struct {
    db          *gorm.DB
    redis       *redis.Client
    rabbitmq    RabbitMQChecker // Optional
    cache       *sync.Map
    cacheTTL    time.Duration
}

func (h *HealthChecker) Handler() fiber.Handler {
    return func(c *fiber.Ctx) error {
        status, checks := h.checkHealth()
        
        statusCode := 200
        if status == "unhealthy" {
            statusCode = 503
        }
        
        return c.Status(statusCode).JSON(fiber.Map{
            "status": status,
            "checks": checks,
        })
    }
}
```

### Go Worker
```go
type HealthChecker struct {
    conn    *amqp.Connection
    cache   *sync.Map
    cacheTTL time.Duration
}

func (h *HealthChecker) Handler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        status, checks := h.checkHealth()
        
        statusCode := 200
        if status == "unhealthy" {
            statusCode = 503
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status": status,
            "checks": checks,
        })
    }
}
```

### Python Service
```python
from fastapi import APIRouter

router = APIRouter()

@router.get("/healthz")
async def health_check():
    # Check dependencies
    checks = []
    status = "healthy"
    
    # No critical dependencies for services
    # (they're stateless)
    
    return {
        "status": status,
        "checks": checks
    }
```

## Kubernetes Integration

### Liveness Probe
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 2
  failureThreshold: 3
```

### Readiness Probe
```yaml
readinessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 2
  failureThreshold: 3
```

## Alternatives Considered

### 1. Separate Endpoints (liveness vs. readiness)
- **Pros**: More granular control
- **Cons**: More complex, two endpoints to maintain
- **Rejected**: Single endpoint is simpler, sufficient for current needs

### 2. No Caching
- **Pros**: Always up-to-date status
- **Cons**: High load on dependencies, slower responses
- **Rejected**: Caching is necessary for performance

### 3. Longer Timeouts
- **Pros**: More accurate for slow dependencies
- **Cons**: Slower health check responses
- **Rejected**: 2 seconds is sufficient, faster is better

### 4. Health Check Authentication
- **Pros**: Security
- **Cons**: Complexity, Kubernetes needs unauthenticated access
- **Rejected**: Health checks are internal, no auth needed

## Best Practices

### Do's
- ✅ Check all critical dependencies
- ✅ Use timeouts for all checks
- ✅ Cache results to reduce load
- ✅ Return appropriate HTTP status codes
- ✅ Include dependency names in response

### Don'ts
- ❌ Don't check non-critical dependencies as critical
- ❌ Don't make health checks slow (> 2 seconds)
- ❌ Don't require authentication
- ❌ Don't perform expensive operations in health checks
- ❌ Don't fail health checks for transient issues

## Future Enhancements

### Planned
- Health check metrics (duration, status changes)
- Health check aggregation (service mesh)
- Health check dashboards
- Automated health check testing

### Under Consideration
- Separate liveness/readiness endpoints (if needed)
- Health check dependencies (circuit breakers)
- Health check versioning

## Notes
- Health checks are lightweight and don't impact normal operations
- Caching means status may be up to 5 seconds stale (acceptable trade-off)
- Different components have different dependency checks (appropriate)
- Health checks are exposed on the same port as the service (or dedicated port for workers)

## Related ADRs
- [ADR-001: RabbitMQ with Redis Fallback](./001-rabbitmq-redis-fallback.md) - Health checks reflect degraded state
- [ADR-002: Standalone Workers Architecture](./002-standalone-workers.md) - Health checks in workers
- [ADR-003: Structured Logging Implementation](./003-structured-logging.md) - Health check logging
