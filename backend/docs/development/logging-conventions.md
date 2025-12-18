# Logging Conventions

## Overview

This guide covers logging conventions and best practices across all backend components. All components use structured logging with consistent patterns.

## Logging Principles

1. **Structured Logging**: Use key-value pairs, not string interpolation
2. **Service Identification**: Include service name in all logs
3. **Trace IDs**: Include trace IDs for request correlation
4. **Appropriate Levels**: Use correct log levels
5. **Context**: Include relevant context (user ID, job ID, etc.)

## Log Levels

### DEBUG
- Detailed information for debugging
- Only in development
- Verbose details

**When to use**:
- Function entry/exit
- Variable values
- Detailed execution flow

**Example**:
```go
logger.Debug("Processing job", "job_id", jobID, "step", "validation")
```

### INFO
- General informational messages
- Normal operations
- Important events

**When to use**:
- Request received
- Job started/completed
- Important state changes

**Example**:
```go
logger.Info("Job processed", "job_id", jobID, "status", "success", "duration_ms", 150)
```

### WARN
- Warning messages
- Non-critical issues
- Recoverable problems

**When to use**:
- Fallback activated
- Retry attempts
- Degraded functionality

**Example**:
```go
logger.Warn("RabbitMQ unavailable, using Redis fallback", "queue", queueName)
```

### ERROR
- Error messages
- Handled errors
- Failed operations

**When to use**:
- Operation failed
- Error caught and handled
- External service errors

**Example**:
```go
logger.Error("Failed to process job", "job_id", jobID, "error", err)
```

### FATAL
- Fatal errors
- Service should exit
- Unrecoverable errors

**When to use**:
- Critical initialization failure
- Unrecoverable state
- Service must exit

**Example**:
```go
logger.Fatal("Failed to connect to database", "error", err)
```

## Structured Logging

### Go (log/slog)

**Good**:
```go
logger.Info("Request processed",
    slog.String("method", "GET"),
    slog.String("path", "/api/projects"),
    slog.Int("status", 200),
    slog.Int("duration_ms", 45))
```

**Bad**:
```go
logger.Info(fmt.Sprintf("Request processed: %s %s - %d (%dms)", method, path, status, duration))
```

### Python (structlog)

**Good**:
```python
logger.info("Request processed",
    method="GET",
    path="/api/projects",
    status=200,
    duration_ms=45)
```

**Bad**:
```python
logger.info(f"Request processed: {method} {path} - {status} ({duration}ms)")
```

### Node.js

**Good**:
```javascript
logger.info('Request processed', {
    method: 'GET',
    path: '/api/projects',
    status: 200,
    duration_ms: 45
});
```

**Bad**:
```javascript
logger.info(`Request processed: ${method} ${path} - ${status} (${duration}ms)`);
```

## Service Identification

All logs should include service name:

**Go**:
```go
logger := logger.New(os.Getenv("ENV"))
logger = logger.With("service", "server")
```

**Python**:
```python
logger = structlog.get_logger().bind(service="ai-service")
```

**Node.js**:
```javascript
const logger = createLogger({ service: 'job-application-worker' });
```

## Trace IDs

Include trace IDs for request correlation:

**Go**:
```go
ctx := logger.WithTraceID(ctx, traceID)
logger.InfoContext(ctx, "Processing request", ...)
```

**Python**:
```python
logger = logger.bind(trace_id=trace_id)
logger.info("Processing request", ...)
```

**Node.js**:
```javascript
logger.info('Processing request', { trace_id, ... });
```

## Context Fields

Include relevant context in logs:

### Common Fields

- `user_id` - User identifier
- `job_id` - Job identifier
- `request_id` - Request identifier
- `trace_id` - Trace identifier
- `duration_ms` - Operation duration
- `status` - Operation status

### Component-Specific Fields

**Server**:
- `method` - HTTP method
- `path` - Request path
- `status_code` - HTTP status code

**Workers**:
- `queue_name` - Queue name
- `job_type` - Job type
- `retry_count` - Retry attempt number

**Services**:
- `provider` - External provider
- `model` - Model used
- `input_length` - Input text length

## Logging Patterns

### Request Logging

**Server**:
```go
logger.Info("Request received",
    slog.String("method", c.Method()),
    slog.String("path", c.Path()),
    slog.String("user_id", userID.String()))
```

**Services**:
```python
logger.info("Request received",
    method=request.method,
    path=request.path,
    user_id=user_id)
```

### Job Processing

**Workers**:
```go
logger.Info("Job processing started",
    slog.String("job_id", jobID.String()),
    slog.String("job_type", jobType))

// ... processing ...

logger.Info("Job processing completed",
    slog.String("job_id", jobID.String()),
    slog.String("status", "success"),
    slog.Int("duration_ms", duration.Milliseconds()))
```

### Error Logging

**Always include error and context**:
```go
logger.Error("Operation failed",
    slog.String("operation", "create_entity"),
    slog.String("entity_id", entityID.String()),
    slog.String("error", err.Error()))
```

**With stack trace (if available)**:
```python
logger.exception("Operation failed",
    operation="create_entity",
    entity_id=entity_id,
    exc_info=True)
```

## What to Log

### Do Log

- ✅ Request/response (with sanitized data)
- ✅ Job start/end
- ✅ Important state changes
- ✅ Errors (with context)
- ✅ External API calls
- ✅ Database operations (if slow or failed)
- ✅ Configuration changes
- ✅ Security events

### Don't Log

- ❌ Passwords or secrets
- ❌ Full request/response bodies (may contain sensitive data)
- ❌ Credit card numbers, SSNs, etc.
- ❌ Excessive detail (use DEBUG level)
- ❌ Every database query (only slow/failed ones)
- ❌ Sensitive user data

## Log Format

### Production (JSON)

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

### Development (Text)

```
2024-01-15T10:30:45.123456789Z INFO service=server trace_id=550e8400... message="Request processed" method=GET path=/api/projects status=200 duration_ms=45
```

## Performance Considerations

1. **Avoid Expensive Operations**:
   - Don't serialize large objects in logs
   - Don't call external services for logging
   - Use lazy evaluation when possible

2. **Log Sampling**:
   - Consider sampling for high-volume logs
   - Always log errors (no sampling)

3. **Async Logging** (if needed):
   - Use async handlers for high-throughput scenarios
   - Ensure logs are flushed on shutdown

## Best Practices

1. **Consistency**:
   - Use consistent field names across components
   - Use consistent log levels
   - Use consistent message formats

2. **Context**:
   - Include enough context to understand the log
   - Include identifiers (user ID, job ID, etc.)
   - Include operation details

3. **Levels**:
   - Use appropriate log levels
   - Don't log everything at ERROR
   - Use DEBUG sparingly (development only)

4. **Structured Data**:
   - Use structured fields, not string interpolation
   - Use appropriate types (int, string, bool)
   - Keep field names consistent

5. **Security**:
   - Never log secrets or passwords
   - Sanitize sensitive data
   - Be careful with user input in logs

## Related Documentation

- [Structured Logging Overview](../architecture/STRUCTURED_LOGGING_OVERVIEW.md) - Comprehensive logging guide
- [ADR-003: Structured Logging Implementation](../adr/003-structured-logging.md) - Architecture decision
- [Adding a Domain](./adding-domain.md) - Domain logging examples
- [Adding a Worker](./adding-worker.md) - Worker logging examples
- [Adding a Service](./adding-service.md) - Service logging examples
