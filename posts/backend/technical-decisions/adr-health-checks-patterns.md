# ADR-005: Health Checks: Consistent Patterns Across Services

## Context
We need health checks for all services and workers. Should be consistent but allow for service-specific needs.

## Decision
All components implement `GET /healthz` endpoint with JSON response. Consistent structure with service-specific checks.

## Rationale

### Consistent Endpoint
- **Standardization**: Same endpoint everywhere (`/healthz`)
- **Tooling**: Monitoring tools can use same endpoint
- **Kubernetes**: Works with liveness/readiness probes

### Consistent Structure
- **Status**: `healthy`, `degraded`, `unhealthy`
- **Checks**: Array of dependency checks
- **Caching**: 5-second cache to reduce load
- **HTTP Codes**: 200 for healthy/degraded, 503 for unhealthy

### Service-Specific Checks
- **Server**: Database, Redis, RabbitMQ (optional)
- **Workers**: RabbitMQ (critical)
- **Services**: Service availability

## Consequences

### Pros
- Consistent monitoring
- Easy to implement
- Kubernetes compatible
- Service-specific flexibility

### Cons
- Need to implement in each service
- Caching adds complexity
- Status granularity varies

## Status
Accepted - 2024-01-15

## Alternatives Considered
- Different endpoints per service: Inconsistent
- No caching: Too much load
- Binary health (healthy/unhealthy only): Less informative
