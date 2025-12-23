# Log Format Specification

**Version:** 1.0  
**Last Updated:** 2025-12-22  
**Status:** Active

## Overview

This document defines the standard log format for all Woragis services. All services must output structured logs that conform to this specification to ensure proper aggregation, searchability, and analysis.

## Standard Fields

All log entries must include the following standard fields:

### Required Fields

| Field | Type | Format | Description | Example |
|-------|------|--------|-------------|---------|
| `timestamp` | string | ISO 8601 | UTC timestamp with nanoseconds | `2025-12-22T10:30:45.123456789Z` |
| `level` | string | lowercase | Log level | `info`, `warn`, `error`, `debug` |
| `service` | string | lowercase with hyphens | Service identifier | `app`, `ai-service`, `resume-worker` |
| `message` | string | plain text | Human-readable log message | `Processing resume generation request` |

### Optional Fields

| Field | Type | Format | Description | Example |
|-------|------|--------|-------------|---------|
| `trace_id` | string | UUID or hex | Distributed tracing ID | `550e8400-e29b-41d4-a716-446655440000` |
| `request_id` | string | UUID or hex | HTTP request ID | `550e8400-e29b-41d4-a716-446655440000` |
| `user_id` | integer/string | number or UUID | User identifier | `123` or `550e8400-e29b-41d4-a716-446655440000` |
| `duration_ms` | number | float | Request/operation duration in milliseconds | `125.5` |
| `status_code` | integer | number | HTTP status code | `200`, `404`, `500` |
| `method` | string | uppercase | HTTP method | `GET`, `POST`, `PUT`, `DELETE` |
| `path` | string | URL path | HTTP request path | `/api/v1/resumes` |
| `ip_address` | string | IP address | Client IP address | `192.168.1.1` |
| `error` | string | plain text | Error message | `Database connection failed` |
| `stack_trace` | string | multi-line | Stack trace (errors only) | Full stack trace |

## Log Levels

Use the following log levels consistently:

- **DEBUG**: Detailed information for debugging (development only)
- **INFO**: General informational messages about normal operation
- **WARN**: Warning messages for potentially problematic situations
- **ERROR**: Error messages for failures that don't stop the service
- **FATAL**: Critical errors that cause the service to stop (use sparingly)

## Service Identifiers

Use the following service identifiers:

- `app` - Main server application
- `ai-service` - AI service
- `creative-service` - Creative service
- `docs-service` - Documentation service
- `resume-worker` - Resume generation worker
- `translation-worker` - Translation worker
- `email-worker` - Email worker
- `whatsapp-worker` - WhatsApp worker
- `job-application-worker` - Job application worker

## Output Format

### Production Environment

In production, all services must output **JSON logs** to stdout/stderr:

```json
{
  "timestamp": "2025-12-22T10:30:45.123456789Z",
  "level": "info",
  "service": "resume-worker",
  "message": "Processing resume generation request",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": 123,
  "job_id": 456
}
```

### Development Environment

In development, services may output human-readable logs for easier debugging, but should still include all standard fields.

## Implementation by Language

### Python Services (structlog)

```python
from logger import get_logger, set_trace_id

logger = get_logger()
set_trace_id("550e8400-e29b-41d4-a716-446655440000")

logger.info(
    "Processing resume generation request",
    user_id=123,
    job_id=456,
    duration_ms=125.5
)
```

**Output (production):**
```json
{
  "timestamp": "2025-12-22T10:30:45.123456789Z",
  "level": "info",
  "service": "resume-worker",
  "message": "Processing resume generation request",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": 123,
  "job_id": 456,
  "duration_ms": 125.5
}
```

### Go Services (slog)

```go
import (
    "context"
    "log/slog"
    "github.com/woragis/worker/pkg/logger"
)

ctx := logger.WithTraceID(context.Background(), "550e8400-e29b-41d4-a716-446655440000")
log := logger.New("production")

log.InfoContext(ctx, "Processing email send request",
    "user_id", 123,
    "email_id", 456,
    "duration_ms", 125.5,
)
```

**Output (production):**
```json
{
  "timestamp": "2025-12-22T10:30:45.123456789Z",
  "level": "info",
  "service": "email-worker",
  "message": "Processing email send request",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": 123,
  "email_id": 456,
  "duration_ms": 125.5
}
```

### Node.js Services

```javascript
import { logger, setTraceId } from './utils/logger.js';

setTraceId("550e8400-e29b-41d4-a716-446655440000");

logger.info("Processing job application", {
  user_id: 123,
  job_id: 456,
  duration_ms: 125.5
});
```

**Output (production):**
```json
{
  "timestamp": "2025-12-22T10:30:45.123456789Z",
  "level": "info",
  "service": "job-application-worker",
  "message": "Processing job application",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": 123,
  "job_id": 456,
  "duration_ms": 125.5
}
```

## Error Logging

When logging errors, include:

1. Error message in `message` field
2. Error details in `error` field
3. Stack trace in `stack_trace` field (if available)
4. Context information (user_id, request_id, etc.)

**Example:**
```json
{
  "timestamp": "2025-12-22T10:30:45.123456789Z",
  "level": "error",
  "service": "resume-worker",
  "message": "Failed to generate resume",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": 123,
  "job_id": 456,
  "error": "Database connection timeout",
  "stack_trace": "Traceback (most recent call last):\n  File..."
}
```

## HTTP Request Logging

For HTTP requests, include:

- `method`: HTTP method
- `path`: Request path
- `status_code`: Response status code
- `duration_ms`: Request duration
- `ip_address`: Client IP (if available)
- `request_id`: Unique request identifier

**Example:**
```json
{
  "timestamp": "2025-12-22T10:30:45.123456789Z",
  "level": "info",
  "service": "app",
  "message": "HTTP request completed",
  "method": "POST",
  "path": "/api/v1/resumes",
  "status_code": 200,
  "duration_ms": 125.5,
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": 123
}
```

## Best Practices

1. **Always include service name**: Every log must have the `service` field
2. **Use trace_id for distributed tracing**: Include `trace_id` when available
3. **Include context**: Add relevant context fields (user_id, request_id, etc.)
4. **Structured data**: Use key-value pairs, not string concatenation
5. **No sensitive data**: Never log passwords, tokens, or PII
6. **Consistent levels**: Use appropriate log levels consistently
7. **ISO 8601 timestamps**: Always use UTC with nanoseconds
8. **JSON in production**: Always output JSON in production environments

## Validation

All services should validate their log output format. Use the following tools:

- **Python**: `structlog` with JSONRenderer
- **Go**: `slog` with JSONHandler
- **Node.js**: Custom JSON logger

## Migration Checklist

- [ ] All services output JSON logs in production
- [ ] All services include `timestamp`, `level`, `service`, `message`
- [ ] All services support `trace_id` for distributed tracing
- [ ] Error logs include `error` and `stack_trace` fields
- [ ] HTTP logs include `method`, `path`, `status_code`, `duration_ms`
- [ ] No sensitive data in logs
- [ ] Logs validated and tested

## References

- [Structured Logging Best Practices](https://www.structlog.org/en/stable/)
- [Go slog Package](https://pkg.go.dev/log/slog)
- [ISO 8601 Date Format](https://en.wikipedia.org/wiki/ISO_8601)
