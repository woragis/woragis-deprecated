# Monitoring: Health Checks + Metrics

## Overview
How health checks and metrics work together for comprehensive system monitoring.

## Key Points

### Health Checks
- Service availability
- Dependency status
- Binary status (healthy/unhealthy)
- Fast checks (cached)

### Metrics
- Request rate
- Latency
- Error rate
- Resource usage

### Combined Monitoring
- Health checks for availability
- Metrics for performance
- Both needed for full picture

## Implementation

### Health Checks
- `/healthz` endpoint
- Dependency checks
- Caching (5s TTL)
- Kubernetes probes

### Metrics
- Prometheus `/metrics` endpoint
- Service-specific metrics
- Business metrics
- Infrastructure metrics

## Benefits
- Availability monitoring
- Performance monitoring
- Issue detection
- Capacity planning

## Challenges
- Two systems to maintain
- Different update frequencies
- Alerting complexity
- Dashboard integration

## Lessons Learned
- Health checks + metrics essential
- Different purposes
- Both needed
- Integration important

## Future Improvements
- Unified dashboard
- Combined alerting
- SLO tracking
- Capacity planning
