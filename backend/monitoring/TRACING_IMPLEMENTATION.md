# Distributed Tracing Implementation Guide

**Status:** In Progress  
**Last Updated:** 2025-12-22

## Overview

This guide documents the implementation of distributed tracing using Jaeger and OpenTelemetry across all Woragis services.

## Architecture

- **Tracing Backend:** Jaeger (all-in-one container)
- **Protocol:** OTLP (OpenTelemetry Protocol) over HTTP
- **Instrumentation:** OpenTelemetry SDK per language
- **Sampling:** 100% in development, 10% in production

## Services Status

### ✅ Infrastructure
- [x] Jaeger service added to docker-compose.yml
- [x] OTLP endpoints configured (HTTP: 4318, gRPC: 4317)

### 🔄 Go Services (In Progress)
- [x] Main server (app) - tracing package created, middleware added
- [ ] email-worker - pending
- [ ] translation-worker - pending
- [ ] whatsapp-worker - pending

### ⏳ Python Services (Pending)
- [ ] ai-service
- [ ] creative-service
- [ ] docs-service
- [ ] resume-worker

### ⏳ Node.js Service (Pending)
- [ ] job-application-worker

## Implementation Details

### Go Services

**Dependencies Required:**
```go
go.opentelemetry.io/otel
go.opentelemetry.io/otel/trace
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
go.opentelemetry.io/otel/propagation
go.opentelemetry.io/otel/sdk/resource
go.opentelemetry.io/otel/sdk/trace
go.opentelemetry.io/otel/semconv/v1.24.0
```

**Implementation Pattern:**
1. Initialize tracing in main()
2. Add tracing middleware (for HTTP servers)
3. Instrument queue consumers
4. Instrument external calls
5. Connect to existing trace_id

### Python Services

**Dependencies Required:**
```python
opentelemetry-api
opentelemetry-sdk
opentelemetry-exporter-otlp-proto-http
opentelemetry-instrumentation-fastapi
opentelemetry-instrumentation-requests
opentelemetry-instrumentation-sqlalchemy
```

**Implementation Pattern:**
1. Initialize OpenTelemetry SDK
2. Configure OTLP HTTP exporter
3. Auto-instrument FastAPI
4. Instrument database operations
5. Connect to existing trace_id context

### Node.js Service

**Dependencies Required:**
```json
"@opentelemetry/api"
"@opentelemetry/sdk-node"
"@opentelemetry/exporter-otlp-http"
"@opentelemetry/instrumentation-http"
"@opentelemetry/instrumentation-amqplib"
```

**Implementation Pattern:**
1. Initialize OpenTelemetry SDK
2. Configure OTLP HTTP exporter
3. Auto-instrument HTTP
4. Instrument queue operations
5. Connect to existing trace_id

## Trace Context Propagation

### HTTP Requests
- **W3C Trace Context:** `traceparent`, `tracestate` headers
- **Custom Header:** `X-Trace-ID` (for compatibility)
- **Propagation:** Automatic via OpenTelemetry propagator

### Queue Messages
- Extract trace context from message headers
- Continue trace across queue boundaries
- Add trace_id to message metadata

### Database Queries
- Add trace context as query metadata
- Log trace_id with queries
- Correlate slow queries with traces

## Configuration

### Environment Variables

```bash
# Jaeger OTLP endpoint (default: http://jaeger:4318)
OTLP_ENDPOINT=http://jaeger:4318

# Sampling rate (0.0 to 1.0, default: auto based on ENV)
TRACING_SAMPLING_RATE=1.0  # 100% in dev, 0.1 (10%) in prod
```

### Service-Specific Configuration

Each service should initialize tracing with:
- Service name
- Service version
- Environment
- OTLP endpoint
- Sampling rate

## Next Steps

1. Complete Go services instrumentation
2. Implement Python services instrumentation
3. Implement Node.js service instrumentation
4. Configure Grafana Jaeger data source
5. Create trace visualization dashboard
6. Link traces to logs and metrics
7. Create documentation

## References

- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [OpenTelemetry Python](https://opentelemetry.io/docs/instrumentation/python/)
- [OpenTelemetry Node.js](https://opentelemetry.io/docs/instrumentation/js/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
