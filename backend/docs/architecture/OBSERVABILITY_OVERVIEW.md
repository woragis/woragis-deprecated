# Observability Overview - Backend Architecture

## General Architecture

Observability is the ability to understand what's happening inside a system by examining its outputs (logs, metrics, traces). Our backend implements a **three-pillar approach** to observability: **Logging**, **Metrics**, and **Distributed Tracing**.

### Current State

- ✅ **Structured Logging**: Fully implemented across all components
- ✅ **Metrics Collection**: Fully implemented (Prometheus client libraries, `/metrics` endpoints)
- ⚠️ **Metrics Infrastructure**: Not deployed (Prometheus server, Grafana dashboards)
- ⚠️ **Distributed Tracing**: Partial (trace IDs exist, but no OpenTelemetry/Jaeger)
- ⚠️ **Dashboards**: Not implemented (Grafana)
- ⚠️ **Testing**: Metrics collection needs testing (cannot test right now)

---

## The Three Pillars of Observability

### 1. **Logging** (✅ Implemented)

**Status**: Fully implemented across all components

**What we have:**
- Structured logging (JSON format in production)
- Trace IDs for request correlation
- Consistent log levels (debug, info, warn, error, fatal)
- Service identification in logs
- Contextual fields (request ID, user ID, job ID, etc.)

**See**: `STRUCTURED_LOGGING_OVERVIEW.md` for detailed implementation

**Limitations:**
- Logs are reactive (you need to know what to look for)
- No aggregation/search interface (yet)
- No alerting based on log patterns

---

### 2. **Metrics** (✅ Collection Implemented, ⚠️ Infrastructure Pending)

**Status**: Metrics collection implemented - **Infrastructure (Prometheus/Grafana) pending**

**What is implemented:**
- ✅ Prometheus client libraries added to all components
- ✅ `/metrics` endpoints exposed on all services/workers
- ✅ HTTP request metrics (server, services)
- ✅ Worker job processing metrics (all workers)
- ✅ Queue metrics (workers)
- ✅ Metrics recording in job processing functions
- ✅ Automatic HTTP metrics for FastAPI services (via prometheus-fastapi-instrumentator)

**What is NOT implemented:**
- ❌ Prometheus server deployment
- ❌ Prometheus scraping configuration
- ❌ Grafana dashboards
- ❌ Metrics testing/verification

**What we need:**
- **Prometheus** for metrics collection
- **Grafana** for visualization
- Key metrics:
  - Request rate (requests per second)
  - Latency (p50, p95, p99)
  - Error rate
  - Queue depth and throughput
  - Database connection pool metrics
  - Worker processing metrics

**Why it matters:**
- Proactive monitoring (detect issues before users notice)
- Performance optimization (identify bottlenecks)
- Capacity planning (know when to scale)
- SLA monitoring (ensure service quality)

**Implementation Plan:**

#### Prometheus Metrics

**Metrics to Expose:**

1. **HTTP Request Metrics:**
   - `http_requests_total` (counter) - Total HTTP requests
   - `http_request_duration_seconds` (histogram) - Request duration
   - `http_requests_in_flight` (gauge) - Current requests being processed

2. **Queue Metrics:**
   - `queue_depth` (gauge) - Current queue size
   - `queue_messages_processed_total` (counter) - Total messages processed
   - `queue_processing_duration_seconds` (histogram) - Message processing time
   - `queue_dlq_size` (gauge) - Dead letter queue size

3. **Database Metrics:**
   - `db_connections_active` (gauge) - Active database connections
   - `db_query_duration_seconds` (histogram) - Query duration
   - `db_queries_total` (counter) - Total queries executed

4. **Worker Metrics:**
   - `worker_jobs_processed_total` (counter) - Total jobs processed
   - `worker_job_duration_seconds` (histogram) - Job processing time
   - `worker_jobs_failed_total` (counter) - Failed jobs
   - `worker_jobs_retried_total` (counter) - Retried jobs

5. **External API Metrics:**
   - `external_api_requests_total` (counter) - External API calls
   - `external_api_duration_seconds` (histogram) - External API call duration
   - `external_api_errors_total` (counter) - External API errors

**Implementation:**

**Go Components (Server, Workers):**
```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)

// Register metrics
func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

// Expose /metrics endpoint
http.Handle("/metrics", promhttp.Handler())
```

**Python Components (Services, Workers):**
```python
from prometheus_client import Counter, Histogram, Gauge, start_http_server

http_requests_total = Counter(
    'http_requests_total',
    'Total number of HTTP requests',
    ['method', 'endpoint', 'status']
)

http_request_duration = Histogram(
    'http_request_duration_seconds',
    'HTTP request duration in seconds',
    ['method', 'endpoint']
)

# Start metrics server
start_http_server(9090)  # Expose /metrics endpoint
```

**Node.js Components:**
```javascript
const promClient = require('prom-client');

const httpRequestsTotal = new promClient.Counter({
  name: 'http_requests_total',
  help: 'Total number of HTTP requests',
  labelNames: ['method', 'endpoint', 'status'],
});

const httpRequestDuration = new promClient.Histogram({
  name: 'http_request_duration_seconds',
  help: 'HTTP request duration in seconds',
  labelNames: ['method', 'endpoint'],
});

// Expose /metrics endpoint
app.get('/metrics', async (req, res) => {
  res.set('Content-Type', promClient.register.contentType);
  res.end(await promClient.register.metrics());
});
```

#### Grafana Dashboards

**Dashboard Structure:**

1. **System Overview Dashboard:**
   - All services health status
   - Request rates across all services
   - Error rates across all services
   - Overall system latency

2. **Service-Specific Dashboards:**
   - Per-service request rate
   - Per-service latency (p50, p95, p99)
   - Per-service error rate
   - Per-service resource usage (CPU, memory)

3. **Queue Monitoring Dashboard:**
   - Queue depths for all queues
   - Message processing rate
   - Dead letter queue sizes
   - Worker processing times

4. **Database Dashboard:**
   - Active connections
   - Query duration
   - Query rate
   - Connection pool usage

5. **Error Tracking Dashboard:**
   - Error rates by type
   - Error rates by service
   - Error trends over time
   - Top error messages

**Example Grafana Queries:**

```promql
# Request rate
rate(http_requests_total[5m])

# Error rate
rate(http_requests_total{status=~"5.."}[5m])

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Queue depth
queue_depth

# Worker processing rate
rate(worker_jobs_processed_total[5m])
```

---

### 3. **Distributed Tracing** (⚠️ Partial)

**Status**: Partial - Trace IDs exist, but no full tracing infrastructure

**What we have:**
- Trace IDs generated and propagated in logs
- Request ID middleware in server
- Trace ID context in workers

**What we need:**
- OpenTelemetry instrumentation
- Jaeger for trace visualization
- Trace context propagation across service boundaries
- Span creation for operations (HTTP requests, DB queries, queue operations)

**Why it matters:**
- Understand request flow across services
- Identify bottlenecks in distributed systems
- Debug issues across service boundaries
- Performance optimization

**Implementation Plan:**

#### OpenTelemetry Setup

**Go Components:**
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/trace"
)

// Initialize tracer
func initTracer() {
    exporter, _ := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint("http://jaeger:14268/api/traces"),
    ))
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("server"),
        )),
    )
    
    otel.SetTracerProvider(tp)
}

// Create spans
func handleRequest(ctx context.Context) {
    tracer := otel.Tracer("server")
    ctx, span := tracer.Start(ctx, "handleRequest")
    defer span.End()
    
    // Your code here
    span.SetAttributes(
        attribute.String("user.id", userID),
        attribute.String("request.method", method),
    )
}
```

**Python Components:**
```python
from opentelemetry import trace
from opentelemetry.exporter.jaeger import JaegerExporter
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

# Initialize tracer
trace.set_tracer_provider(TracerProvider())
tracer = trace.get_tracer(__name__)

jaeger_exporter = JaegerExporter(
    agent_host_name="jaeger",
    agent_port=6831,
)

span_processor = BatchSpanProcessor(jaeger_exporter)
trace.get_tracer_provider().add_span_processor(span_processor)

# Create spans
def handle_request():
    with tracer.start_as_current_span("handle_request") as span:
        span.set_attribute("user.id", user_id)
        span.set_attribute("request.method", method)
        # Your code here
```

**Trace Context Propagation:**

- **HTTP Headers**: `traceparent`, `tracestate` (W3C Trace Context)
- **RabbitMQ Headers**: Include trace context in message headers
- **Database**: Include trace context in query metadata (if supported)

**Spans to Create:**

1. **HTTP Request Spans:**
   - Server receives request
   - Handler processing
   - External API calls
   - Database queries
   - Response sent

2. **Worker Spans:**
   - Message consumed
   - Job processing
   - External API calls
   - Database operations
   - Message acknowledgment

3. **Service Spans:**
   - Request received
   - Processing
   - Response sent

---

## Implementation Roadmap

### Phase 1: Metrics (High Priority) - 2-3 weeks

**Week 1:**
- Set up Prometheus
- Implement basic HTTP metrics in server
- Implement queue metrics in workers
- Expose `/metrics` endpoints

**Week 2:**
- Set up Grafana
- Create system overview dashboard
- Create service-specific dashboards
- Create queue monitoring dashboard

**Week 3:**
- Add database metrics
- Add external API metrics
- Create error tracking dashboard
- Set up basic alerts

### Phase 2: Distributed Tracing (Medium Priority) - 2-3 weeks

**Week 1:**
- Set up Jaeger
- Implement OpenTelemetry in server
- Implement trace context propagation

**Week 2:**
- Implement OpenTelemetry in workers
- Implement OpenTelemetry in services
- Add spans for key operations

**Week 3:**
- Trace context propagation across RabbitMQ
- Trace context propagation for external APIs
- Create trace visualization dashboards

### Phase 3: Enhanced Logging (Low Priority) - 1 week

- Set up log aggregation (Loki or ELK)
- Set up log search interface
- Create log-based alerts
- Integrate logs with traces

---

## Best Practices

### Metrics

1. **Use appropriate metric types:**
   - Counters for cumulative values (total requests)
   - Gauges for current values (queue depth)
   - Histograms for distributions (latency)

2. **Label carefully:**
   - Use labels for dimensions you want to query
   - Don't use high-cardinality labels (user IDs)
   - Keep label names consistent across services

3. **Expose metrics on dedicated endpoint:**
   - `/metrics` endpoint (standard Prometheus format)
   - Don't require authentication (internal network only)

4. **Set up alerts:**
   - Alert on error rate spikes
   - Alert on latency increases
   - Alert on queue depth growth
   - Alert on service downtime

### Distributed Tracing

1. **Create spans for operations:**
   - HTTP requests
   - Database queries
   - External API calls
   - Queue operations
   - Business logic operations

2. **Propagate trace context:**
   - Include trace context in HTTP headers
   - Include trace context in RabbitMQ message headers
   - Maintain trace context across async operations

3. **Add meaningful attributes:**
   - User ID, request ID, job ID
   - Operation parameters
   - Error details
   - Performance metrics

4. **Sample appropriately:**
   - Don't trace every request (performance impact)
   - Use sampling (e.g., 10% of requests)
   - Always trace errors
   - Always trace slow requests

### Logging (Already Implemented)

1. **Structured logging:** ✅ Already done
2. **Trace IDs:** ✅ Already done
3. **Consistent format:** ✅ Already done
4. **Appropriate log levels:** ✅ Already done

---

## Tools and Technologies

### Metrics Stack

- **Prometheus**: Metrics collection and storage
- **Grafana**: Metrics visualization and dashboards
- **Alertmanager**: Alert routing and notification

### Tracing Stack

- **OpenTelemetry**: Instrumentation SDK
- **Jaeger**: Distributed tracing backend
- **W3C Trace Context**: Trace context propagation standard

### Logging Stack (Future)

- **Loki**: Log aggregation (Grafana Labs)
- **ELK Stack**: Elasticsearch, Logstash, Kibana
- **Fluentd/Fluent Bit**: Log forwarding

---

## Integration with Existing Systems

### Health Checks

- Expose health check metrics:
  - `health_check_total` (counter)
  - `health_check_duration_seconds` (histogram)
  - `health_check_status` (gauge) - 1 for healthy, 0 for unhealthy

### Structured Logging

- Include trace ID in logs (already done)
- Include metrics in logs (e.g., request duration)
- Correlate logs with traces using trace ID

### CI/CD

- Metrics for build/deployment:
  - `deployment_duration_seconds`
  - `deployment_status` (success/failure)
  - `build_duration_seconds`

---

## Monitoring Strategy

### Key Metrics to Monitor

1. **Availability:**
   - Service uptime
   - Health check success rate
   - Error rate

2. **Performance:**
   - Request latency (p50, p95, p99)
   - Throughput (requests per second)
   - Queue processing rate

3. **Reliability:**
   - Error rate
   - Failed job rate
   - Dead letter queue size

4. **Resource Usage:**
   - CPU usage
   - Memory usage
   - Database connections
   - Queue depth

### Alerting Strategy

**Critical Alerts (Page):**
- Service down
- Error rate > 5%
- Queue depth > 1000
- Database connection pool exhausted

**Warning Alerts (Notify):**
- Error rate > 1%
- Latency p95 > 1 second
- Queue depth > 500
- Dead letter queue growing

**Info Alerts (Log):**
- Deployment completed
- Configuration changed
- Scaling event

---

## Cost Considerations

### Metrics Storage

- **Prometheus**: ~1-2GB per million samples
- **Retention**: 15-30 days (adjustable)
- **Downsampling**: Keep detailed metrics for 7 days, aggregated for 30 days

### Tracing Storage

- **Jaeger**: ~1KB per span
- **Retention**: 7 days (adjustable)
- **Sampling**: 10% of requests (reduces storage by 90%)

### Logging Storage

- **Loki**: Compressed storage (~10x compression)
- **Retention**: 30 days (adjustable)
- **Indexing**: Only index labels, not full log content

---

## Next Steps

1. **Immediate (Week 1):**
   - Set up Prometheus
   - Implement basic HTTP metrics in server
   - Expose `/metrics` endpoints

2. **Short-term (Weeks 2-3):**
   - Set up Grafana
   - Create dashboards
   - Implement queue metrics

3. **Medium-term (Weeks 4-6):**
   - Set up OpenTelemetry/Jaeger
   - Implement distributed tracing
   - Integrate with existing trace IDs

4. **Long-term (Weeks 7+):**
   - Set up log aggregation
   - Create comprehensive dashboards
   - Set up advanced alerting
   - Performance optimization based on metrics

---

## Summary

**Current State:**
- ✅ Structured logging (complete)
- ❌ Metrics (not implemented)
- ⚠️ Distributed tracing (partial - trace IDs only)

**Priority:**
1. **Metrics** (Prometheus + Grafana) - High priority
2. **Distributed Tracing** (OpenTelemetry + Jaeger) - Medium priority
3. **Log Aggregation** (Loki/ELK) - Low priority

**Impact:**
- Metrics: Proactive monitoring, performance optimization
- Tracing: Debug distributed systems, identify bottlenecks
- Logs: Already provide good observability foundation
