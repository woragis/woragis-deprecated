# Structured Logging - Email Worker

## Overview

The Email Worker uses structured logging with Go's `log/slog` for consistent, searchable logs across all environments.

## Features

- **JSON format in production** - Easy to parse and aggregate
- **Text format in development** - Human-readable for debugging
- **Automatic service name** - All logs include `service: "email-worker"`
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
  "service": "email-worker",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Email sent successfully",
  "user_id": "6ad0d828-1234-5678-90ab-cdef12345678",
  "destination": "user@example.com",
  "subject": "Welcome to Woragis"
}
```

### Development (Text)
```
2024-01-15T10:30:45.123456789Z INFO service=email-worker trace_id=550e8400... message="Email sent successfully" user_id=6ad0d828... destination=user@example.com subject="Welcome to Woragis"
```

## Usage

### Basic Logging

```go
import "github.com/woragis/backend/email-worker/pkg/logger"

logger := logger.New(os.Getenv("ENV"))

// Info log
logger.Info("Email sent successfully",
    slog.String("user_id", userID),
    slog.String("destination", destination),
    slog.String("subject", subject))

// Error log with context
logger.Error("Failed to send email",
    slog.String("user_id", userID),
    slog.String("destination", destination),
    slog.Any("error", err))
```

### Trace ID Support

```go
import "github.com/woragis/backend/email-worker/pkg/logger"

// Set trace_id in context (from message headers or generate new)
ctx := logger.WithTraceID(ctx, traceID)

// All logs in this context will include trace_id
logger.InfoContext(ctx, "Processing email", ...)
```

## Log Storage

- **Development:** stdout (default) or files (if `LOG_TO_FILE=true`)
- **Production:** stdout (collected by Kubernetes/log aggregator)

## Best Practices

1. **Use structured fields** - Always use key-value pairs, not string interpolation
   ```go
   // Good
   logger.Info("Email sent", 
       slog.String("user_id", userID),
       slog.Duration("duration", duration))
   
   // Bad
   logger.Info(fmt.Sprintf("Email sent for user %s in %v", userID, duration))
   ```

2. **Include context** - Add relevant fields to help debugging
   ```go
   logger.Error("Send failed",
       slog.String("user_id", envelope.UserID),
       slog.String("destination", envelope.Destination),
       slog.String("subject", envelope.Subject),
       slog.Any("error", err))
   ```

3. **Never log sensitive data** - No passwords, tokens, or PII

4. **Use appropriate levels**
   - `debug`: Detailed diagnostic info (dev only)
   - `info`: General operational messages
   - `warn`: Warnings, retries, connection issues
   - `error`: Errors, failures, exceptions

## Integration with Log Aggregation

In production, logs go to stdout and are collected by:
- Kubernetes → Log aggregator (Fluentd/Fluent Bit) → Database/ELK/Loki

See main backend TODO.md for production log aggregation setup.
