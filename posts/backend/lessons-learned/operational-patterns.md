# Operational Patterns: Health Checks and Monitoring

## Overview
Operational patterns that work well: health checks, monitoring, and observability.

## Key Points

### Health Checks
- `/healthz` endpoint everywhere
- Dependency checking
- Caching (5s TTL)
- Kubernetes compatible

### Monitoring
- Structured logging
- Health check status
- Error tracking
- Performance tracking (future)

### Observability
- Trace ID propagation
- Request ID middleware
- Consistent log format
- Service identification

## What Works Well

### Health Checks
- Simple to implement
- Kubernetes integration
- Dependency visibility
- Caching reduces load

### Structured Logging
- Machine parseable
- Trace correlation
- Consistent format
- Easy aggregation

### Request ID
- Request tracking
- Log correlation
- Debugging help
- Distributed tracing foundation

## What's Missing

### Metrics
- No Prometheus yet
- No performance metrics
- No business metrics
- No capacity metrics

### Dashboards
- No Grafana
- No visualization
- No analytics
- Manual inspection

### Alerting
- No automated alerts
- Manual monitoring
- Reactive only
- No SLO tracking

## Lessons Learned
- Health checks essential
- Structured logging foundation
- Metrics needed for performance
- Dashboards provide visibility

## Future Improvements
- Prometheus metrics
- Grafana dashboards
- OpenTelemetry tracing
- Automated alerting
