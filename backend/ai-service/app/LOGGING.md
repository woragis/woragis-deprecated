# Structured Logging - AI Service

## Overview

The AI Service uses structured logging with `structlog` for consistent, searchable logs across all environments.

## Features

- **JSON format in production** - Easy to parse and aggregate
- **Text format in development** - Human-readable for debugging
- **Automatic service name** - All logs include `service: "ai-service"`
- **Trace ID support** - Distributed tracing via `X-Trace-ID` headers
- **HTTP request logging** - Automatic logging of all requests with context
- **Structured fields** - All log entries use key-value pairs

## Configuration

### Environment Variables

```bash
# Environment (affects log format and level)
ENV=development  # or "production"

# Optional: Enable file logging in development
LOG_TO_FILE=true
LOG_DIR=logs
```

### Default Behavior

- **Development:** Text format, DEBUG level, stdout
- **Production:** JSON format, INFO level, stdout

## Log Format

### Production (JSON)
```json
{
  "timestamp": "2024-01-15T10:30:45.123456789Z",
  "level": "info",
  "service": "ai-service",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "chat request",
  "agent": "economist",
  "provider": "openai",
  "model": "gpt-4o-mini"
}
```

### Development (Text)
```
2024-01-15T10:30:45.123456789Z [info     ] chat request                    agent=economist provider=openai model=gpt-4o-mini service=ai-service trace_id=550e8400...
```

## Usage

### Basic Logging

```python
from app.logger import get_logger

logger = get_logger()

# Info log
logger.info("Service started", port=8000)

# Error log with context
logger.error("Failed to process request", 
             agent="economist", 
             error=str(e),
             exc_info=True)
```

### HTTP Request Logging

Automatic via middleware - all requests are logged with:
- Method, path, IP, status code
- Duration in milliseconds
- Trace ID (if available)
- Query parameters (if present)

### Trace ID Propagation

Trace IDs are automatically:
- Generated for each request (UUID)
- Read from `X-Trace-ID` header (for distributed tracing)
- Added to response headers
- Included in all log entries

## Log Storage

- **Development:** stdout (default) or files (if `LOG_TO_FILE=true`)
- **Production:** stdout (collected by Kubernetes/log aggregator)

## Best Practices

1. **Use structured fields** - Always use key-value pairs, not string interpolation
   ```python
   # Good
   logger.info("Request processed", user_id=user_id, duration_ms=duration)
   
   # Bad
   logger.info(f"Request processed for user {user_id} in {duration}ms")
   ```

2. **Include context** - Add relevant fields to help debugging
   ```python
   logger.error("Chat failed", 
                agent=agent_name,
                provider=provider,
                input_length=len(input_text),
                error=str(e))
   ```

3. **Never log sensitive data** - No passwords, tokens, or PII

4. **Use appropriate levels**
   - `debug`: Detailed diagnostic info (dev only)
   - `info`: General operational messages
   - `warn`: Warnings, retries, degraded functionality
   - `error`: Errors, failures, exceptions

## Integration with Log Aggregation

In production, logs go to stdout and are collected by:
- Kubernetes → Log aggregator (Fluentd/Fluent Bit) → Database/ELK/Loki

See main backend TODO.md for production log aggregation setup.
