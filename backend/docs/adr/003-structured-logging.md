# ADR-003: Structured Logging Implementation

## Status
Accepted

## Context
The backend consists of multiple services and workers written in different languages (Go, Python, Node.js). We need a consistent logging strategy that:
- Provides observability across all components
- Supports log aggregation and analysis
- Enables distributed tracing
- Works well with containerized deployments
- Is easy to use for developers

**Requirements:**
- Consistent log format across all components
- Structured data (key-value pairs, not string interpolation)
- Trace ID support for distributed tracing
- Service identification in logs
- Production-ready (JSON format)
- Development-friendly (human-readable format)

**Constraints:**
- Different languages have different logging libraries
- Need to balance consistency with language-specific best practices
- Logs should be suitable for stdout (Kubernetes/log aggregation)

## Decision
We will implement **structured logging** with **language-specific implementations** that follow consistent patterns.

**Patterns:**
- **Format**: JSON in production, human-readable text in development
- **Service Name**: Automatically included in all logs
- **Trace ID**: Support for distributed tracing via context
- **ISO 8601 Timestamps**: Standard timestamp format
- **Structured Fields**: Key-value pairs (no string interpolation)
- **Output**: stdout (for Kubernetes/log aggregation)

**Language-Specific Implementations:**

### Go Components
- **Library**: `log/slog` (standard library)
- **Format**: JSON handler in production, text handler in development
- **Service Name**: Injected via custom handler wrapper
- **Trace ID**: From context (`trace_id` key)

### Python Components
- **Library**: `structlog`
- **Format**: JSON processor in production, console processor in development
- **Service Name**: Configured per service
- **Trace ID**: From context/bound logger

### Node.js Components
- **Library**: Custom logger (can use `pino` or `winston` in future)
- **Format**: JSON in production, text in development
- **Service Name**: Configured per service
- **Trace ID**: From context

## Consequences

### Positive
- ✅ **Consistency**: Similar log structure across all components
- ✅ **Observability**: Easy to search and analyze logs
- ✅ **Distributed Tracing**: Trace IDs enable request correlation
- ✅ **Production Ready**: JSON format works with log aggregation tools
- ✅ **Developer Friendly**: Human-readable format in development
- ✅ **Container Friendly**: stdout output works with Kubernetes

### Negative
- ⚠️ **Code Duplication**: Similar logging code duplicated across components
- ⚠️ **Language Differences**: Slight differences in implementation across languages
- ⚠️ **Learning Curve**: Developers need to learn structured logging patterns

### Neutral
- Each language uses its standard/recommended logging library
- Log format is consistent but implementation details differ

## Implementation Details

### Go Implementation

#### Logger Package
```go
// pkg/logger/logger.go
func New(env string) *slog.Logger {
    var handler slog.Handler
    if env == "production" {
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        })
    } else {
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        })
    }
    
    return slog.New(&serviceHandler{
        handler: handler,
        service: "service-name",
    })
}
```

#### Usage
```go
logger := logger.New(os.Getenv("ENV"))
logger.Info("Request processed",
    slog.String("method", "GET"),
    slog.String("path", "/api/projects"),
    slog.Int("status", 200))
```

### Python Implementation

#### Logger Setup
```python
# Using structlog
import structlog

logger = structlog.get_logger()
logger = logger.bind(service="service-name")

# In production: JSON processor
# In development: Console processor
```

#### Usage
```python
logger.info("request_processed",
    method="GET",
    path="/api/projects",
    status=200)
```

### Node.js Implementation

#### Logger Setup
```javascript
// Custom logger
const logger = createLogger({
    service: 'service-name',
    format: process.env.ENV === 'production' ? 'json' : 'text'
});
```

#### Usage
```javascript
logger.info('Request processed', {
    method: 'GET',
    path: '/api/projects',
    status: 200
});
```

### Log Format Examples

#### Production (JSON)
```json
{
  "timestamp": "2024-01-15T10:30:45.123456789Z",
  "level": "info",
  "service": "server",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Request processed",
  "method": "GET",
  "path": "/api/projects",
  "status": 200,
  "duration_ms": 45
}
```

#### Development (Text)
```
2024-01-15T10:30:45.123456789Z INFO service=server trace_id=550e8400... message="Request processed" method=GET path=/api/projects status=200 duration_ms=45
```

## Alternatives Considered

### 1. Single Logging Library Across All Languages
- **Pros**: Perfect consistency
- **Cons**: Not idiomatic, forces unnatural patterns
- **Rejected**: Goes against language best practices

### 2. String Interpolation Logging
- **Pros**: Simple, familiar
- **Cons**: Hard to parse, no structured data, poor for log aggregation
- **Rejected**: Doesn't meet observability requirements

### 3. Binary Log Format (Protobuf, etc.)
- **Pros**: Efficient, structured
- **Cons**: Not human-readable, requires special tools
- **Rejected**: Too complex, not developer-friendly

### 4. Centralized Logging Service
- **Pros**: Single point of control
- **Cons**: Additional infrastructure, single point of failure
- **Rejected**: Overkill for current scale, prefer stdout approach

## Log Levels

### Standard Levels
- **DEBUG**: Detailed information for debugging
- **INFO**: General informational messages
- **WARN**: Warning messages (non-critical issues)
- **ERROR**: Error messages (handled errors)
- **FATAL**: Fatal errors (service should exit)

### Usage Guidelines
- **DEBUG**: Only in development, verbose details
- **INFO**: Normal operations, important events
- **WARN**: Recoverable issues, fallbacks activated
- **ERROR**: Errors that are handled but logged
- **FATAL**: Unrecoverable errors, service exit

## Trace ID Propagation

### Context-Based
- Trace ID stored in context (Go) or bound logger (Python)
- Automatically included in all log entries
- Propagated across service boundaries (HTTP headers, RabbitMQ headers)

### Generation
- Generated at request entry point (server)
- Included in HTTP headers (`X-Trace-ID`)
- Included in RabbitMQ message headers
- Workers extract from message headers

## Future Enhancements

### Planned
- Log aggregation setup (Loki, ELK, etc.)
- Log-based alerting
- Integration with distributed tracing (OpenTelemetry)
- Log sampling for high-volume endpoints

### Under Consideration
- Centralized log format validation
- Log retention policies
- Log compression for storage

## Notes
- Logs are written to stdout (not files) for Kubernetes compatibility
- File logging available in development (optional, via `LOG_TO_FILE` env var)
- Service name is injected automatically, not manually added to each log call
- Trace IDs are optional (not all logs have trace IDs)

## Related ADRs
- [ADR-002: Standalone Workers Architecture](./002-standalone-workers.md) - Logging in workers
- [ADR-005: Health Checks Implementation Strategy](./005-health-checks.md) - Health check logging
