# Structured Logging - Translation Worker

## Overview

The Translation Worker uses structured logging with Go's `log/slog` for consistent, searchable logs across all environments.

## Features

- **JSON format in production** - Easy to parse and aggregate
- **Text format in development** - Human-readable for debugging
- **Automatic service name** - All logs include `service: "translation-worker"`
- **Trace ID support** - Distributed tracing via context
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
  "service": "translation-worker",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Translation job completed successfully",
  "job_id": "abc123",
  "entity_type": "project",
  "entity_id": "6ad0d828-1234-5678-90ab-cdef12345678",
  "language": "pt-BR"
}
```

### Development (Text)
```
2024-01-15T10:30:45.123456789Z INFO service=translation-worker trace_id=550e8400... message="Translation job completed successfully" job_id=abc123 entity_type=project entity_id=6ad0d828... language=pt-BR
```

## Usage

### Basic Logging

```go
import "github.com/woragis/backend/translation-worker/pkg/logger"

logger := logger.New(os.Getenv("ENV"))

// Info log
logger.Info("Translation job completed successfully",
    slog.String("job_id", jobID),
    slog.String("entity_type", entityType),
    slog.String("entity_id", entityID),
    slog.String("language", language))

// Error log with context
logger.Error("Failed to translate field",
    slog.String("job_id", jobID),
    slog.String("field", field),
    slog.Any("error", err))
```

### Trace ID Support

```go
import "github.com/woragis/backend/translation-worker/pkg/logger"

// Set trace_id in context (from message headers or generate new)
ctx := logger.WithTraceID(ctx, traceID)

// All logs in this context will include trace_id
logger.InfoContext(ctx, "Processing translation job", ...)
```

## Log Storage

- **Development:** stdout (default) or files (if `LOG_TO_FILE=true`)
- **Production:** stdout (collected by Kubernetes/log aggregator)

## Best Practices

1. **Use structured fields** - Always use key-value pairs, not string interpolation
   ```go
   // Good
   logger.Info("Translation completed", 
       slog.String("job_id", jobID),
       slog.Duration("duration", duration))
   
   // Bad
   logger.Info(fmt.Sprintf("Translation completed for job %s in %v", jobID, duration))
   ```

2. **Include context** - Add relevant fields to help debugging
   ```go
   logger.Error("Translation failed",
       slog.String("job_id", job.ID),
       slog.String("entity_type", job.EntityType),
       slog.String("entity_id", job.EntityID),
       slog.String("language", job.Language),
       slog.Any("error", err))
   ```

3. **Never log sensitive data** - No API keys, tokens, or PII

4. **Use appropriate levels**
   - `debug`: Detailed diagnostic info (dev only)
   - `info`: General operational messages
   - `warn`: Warnings, retries, API rate limits
   - `error`: Errors, failures, translation API failures

## Integration with Log Aggregation

In production, logs go to stdout and are collected by:
- Kubernetes → Log aggregator (Fluentd/Fluent Bit) → Database/ELK/Loki

See main backend TODO.md for production log aggregation setup.
