# Distributed Tracing Guide

**Version:** 1.0  
**Last Updated:** 2025-12-22

## Overview

This guide covers the distributed tracing system using Jaeger and OpenTelemetry. Distributed tracing enables end-to-end request tracking across all Woragis services.

## Architecture

- **Jaeger**: Tracing backend (all-in-one container)
- **OpenTelemetry**: Instrumentation SDK
- **OTLP**: OpenTelemetry Protocol over HTTP
- **Sampling**: 100% in development, 10% in production

## Quick Start

1. **Start Jaeger:**
   ```bash
   docker-compose up -d jaeger
   ```

2. **Access Jaeger UI:**
   - URL: http://localhost:16686
   - Search for traces by service, operation, or trace ID

3. **Make a request:**
   ```bash
   curl http://localhost:8080/api/health
   ```

4. **View trace:**
   - Go to Jaeger UI
   - Select service: "woragis-server"
   - Click "Find Traces"

## Services Instrumented

### Go Services
- ✅ `app` (main server) - HTTP requests instrumented
- ✅ `email-worker` - Queue operations (pending full implementation)
- ✅ `translation-worker` - Queue operations (pending full implementation)
- ✅ `whatsapp-worker` - Queue operations (pending full implementation)

### Python Services
- ✅ `ai-service` - FastAPI auto-instrumented
- ✅ `creative-service` - FastAPI auto-instrumented
- ✅ `docs-service` - FastAPI auto-instrumented
- ✅ `resume-worker` - HTTP requests instrumented

### Node.js Service
- ✅ `job-application-worker` - HTTP and queue operations instrumented

## Trace Context Propagation

### HTTP Requests
Traces automatically propagate via:
- **W3C Trace Context**: `traceparent`, `tracestate` headers
- **Custom Header**: `X-Trace-ID` (for compatibility)

### Queue Messages
- Trace context extracted from message headers
- Continues trace across queue boundaries
- Trace ID added to message metadata

## Viewing Traces

### Jaeger UI

1. **Search Traces:**
   - Select service from dropdown
   - Choose time range
   - Click "Find Traces"

2. **View Trace Details:**
   - Click on a trace
   - See span timeline
   - View span details and tags
   - See service dependencies

3. **Compare Traces:**
   - Select multiple traces
   - Compare performance
   - Identify differences

### Grafana (when configured)

1. Go to Explore
2. Select Jaeger data source
3. Search for traces
4. View trace details

## Trace-to-Log Correlation

### Using Trace ID

1. **Get trace ID from Jaeger:**
   - Open trace in Jaeger UI
   - Copy trace ID

2. **Search logs in Grafana:**
   ```
   {job="docker", trace_id="YOUR-TRACE-ID"}
   ```

3. **View all logs for a trace:**
   - All services that participated in the trace
   - Chronological order
   - Full request flow

## Trace-to-Metrics Correlation

Traces can be correlated with metrics:
- Request duration in traces matches latency metrics
- Error spans match error rate metrics
- Service dependencies match service health

## Configuration

### Environment Variables

```bash
# OTLP endpoint (default: http://jaeger:4318)
OTLP_ENDPOINT=http://jaeger:4318

# Sampling rate (0.0 to 1.0, auto-based on ENV)
TRACING_SAMPLING_RATE=1.0  # 100% in dev, 0.1 (10%) in prod
```

### Service Configuration

Each service initializes tracing with:
- Service name
- Service version
- Environment
- OTLP endpoint
- Sampling rate

## Sampling Strategy

### Development
- **Rate**: 100% (all traces)
- **Rationale**: Full visibility for debugging

### Production
- **Rate**: 10-20% (configurable)
- **Rationale**: Balance visibility with performance
- **Method**: Head-based sampling
- **Adjustment**: Increase for errors, decrease for normal traffic

## Common Use Cases

### Debugging Slow Requests

1. Find slow trace in Jaeger
2. Identify slowest span
3. Check logs for that span
4. Review metrics for that service

### Tracing Request Flow

1. Start from entry point (main server)
2. Follow trace through services
3. See service dependencies
4. Identify bottlenecks

### Error Investigation

1. Find error trace in Jaeger
2. See which service failed
3. Check error logs
4. Review error metrics

## Best Practices

1. **Use appropriate span names:**
   - Descriptive operation names
   - Include service name

2. **Add relevant attributes:**
   - User ID, request ID
   - Operation parameters
   - Error details

3. **Keep spans focused:**
   - One operation per span
   - Don't create too many spans

4. **Monitor trace volume:**
   - Adjust sampling if needed
   - Monitor Jaeger storage

## Troubleshooting

### No Traces Appearing

1. **Check Jaeger is running:**
   ```bash
   docker-compose ps jaeger
   docker-compose logs jaeger
   ```

2. **Check service is sending traces:**
   - Verify OTLP endpoint is correct
   - Check service logs for errors
   - Test OTLP endpoint connectivity

3. **Check sampling:**
   - Development: 100% sampling
   - Production: 10% sampling (may miss traces)

### Trace ID Mismatch

- Ensure trace_id from OpenTelemetry matches log trace_id
- Check middleware order (tracing before logging)
- Verify trace context propagation

### High Memory Usage

- Reduce sampling rate
- Configure trace retention
- Limit span attributes

## Resources

- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [OpenTelemetry Python](https://opentelemetry.io/docs/instrumentation/python/)
- [OpenTelemetry Node.js](https://opentelemetry.io/docs/instrumentation/js/)
- [W3C Trace Context](https://www.w3.org/TR/trace-context/)
