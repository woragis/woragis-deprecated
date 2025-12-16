# Structured Logging - Resume Worker

## Overview

The Resume Worker uses structured logging with `structlog` for consistent, searchable logs across all environments.

## Features

- **JSON format in production** - Easy to parse and aggregate
- **Text format in development** - Human-readable for debugging
- **Automatic service name** - All logs include `service: "resume-worker"`
- **Trace ID support** - Distributed tracing via context variables
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
  "service": "resume-worker",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Resume successfully generated",
  "user_id": "6ad0d828-1234-5678-90ab-cdef12345678",
  "output_path": "/app/output/resume.pdf",
  "file_size": 123456,
  "projects_count": 5,
  "certifications_count": 3
}
```

### Development (Text)
```
2024-01-15T10:30:45.123456789Z [info     ] Resume successfully generated     user_id=6ad0d828... output_path=/app/output/resume.pdf file_size=123456 service=resume-worker trace_id=550e8400...
```

## Usage

### Basic Logging

```python
from logger import get_logger

logger = get_logger()

# Info log
logger.info("Resume successfully generated",
            user_id=user_id,
            output_path=output_path,
            file_size=file_size,
            projects_count=projects_count)

# Error log with context
logger.error("Error generating resume",
             user_id=user_id,
             language=language,
             error=str(e),
             exc_info=True)
```

### Trace ID Support

```python
from logger import set_trace_id, get_trace_id

# Set trace_id from message headers
set_trace_id(message.get('trace_id', ''))

# All subsequent logs will include trace_id
logger.info("Processing resume job", user_id=user_id)
```

## Log Storage

- **Development:** stdout (default) or files (if `LOG_TO_FILE=true`)
- **Production:** stdout (collected by Kubernetes/log aggregator)

## Best Practices

1. **Use structured fields** - Always use key-value pairs, not f-strings
   ```python
   # Good
   logger.info("Resume generated", user_id=user_id, file_size=file_size)
   
   # Bad
   logger.info(f"Resume generated for user {user_id}, size {file_size}")
   ```

2. **Include context** - Add relevant fields to help debugging
3. **Never log sensitive data** - No passwords, tokens, or PII
4. **Use appropriate levels** - debug/info/warn/error

## Integration with Log Aggregation

In production, logs go to stdout and are collected by:
- Kubernetes → Log aggregator (Fluentd/Fluent Bit) → Database/ELK/Loki

See main backend TODO.md for production log aggregation setup.
