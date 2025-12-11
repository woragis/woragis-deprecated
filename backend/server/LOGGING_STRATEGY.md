# Logging Strategy for Woragis Backend

## Overview

This document outlines the logging strategy for the Woragis backend services, focusing on structured logging, log rotation, and scalability.

## Principles

1. **Structured Logging**: Use JSON format for easy parsing and analysis
2. **Log Levels**: Appropriate use of ERROR, WARN, INFO, and DEBUG
3. **Log Rotation**: Automatic rotation to prevent disk space issues
4. **Separation of Concerns**: Separate log files for different services and log levels
5. **Performance**: Async/buffered logging to avoid blocking operations
6. **Security**: Never log sensitive data (API keys, passwords, full user data)

## Directory Structure

```
backend/server/
  logs/
    app.log                    # Main application log (all levels)
    app.error.log              # Errors only
    resumes.log                # Resume domain logs
    resumes.error.log          # Resume errors only
    job-applications.log       # Job application domain logs
    job-applications.error.log # Job application errors only
    archive/                   # Rotated logs
      2024/
        01/
          app.2024-01-01.log.gz
          app.error.2024-01-01.log.gz
```

## Log Levels

### ERROR
- System failures
- Critical errors that require immediate attention
- Failed operations that cannot be recovered
- Authentication/authorization failures
- Database connection failures

**Example:**
```json
{
  "timestamp": "2024-01-01T12:00:00.123Z",
  "level": "ERROR",
  "service": "resumes",
  "message": "Failed to generate resume",
  "error": "AI service timeout",
  "job_id": "550e8400-...",
  "user_id": "6ad0d828-...",
  "retry_count": 3
}
```

### WARN
- Recoverable errors
- Retry attempts
- Rate limiting
- Degraded performance
- Non-critical failures

**Example:**
```json
{
  "timestamp": "2024-01-01T12:00:00.123Z",
  "level": "WARN",
  "service": "resumes",
  "message": "Retrying job after transient error",
  "job_id": "550e8400-...",
  "retry_count": 2,
  "backoff_seconds": 4
}
```

### INFO
- Important business events
- Job lifecycle events (enqueued, started, completed)
- Service startup/shutdown
- Successful operations
- Status changes

**Example:**
```json
{
  "timestamp": "2024-01-01T12:00:00.123Z",
  "level": "INFO",
  "service": "resumes",
  "message": "Resume generation completed",
  "job_id": "550e8400-...",
  "user_id": "6ad0d828-...",
  "duration_ms": 4523,
  "file_size": 123456
}
```

### DEBUG
- Detailed processing steps
- Request/response details (sanitized)
- Internal state changes
- Performance metrics
- Only enabled in development

**Example:**
```json
{
  "timestamp": "2024-01-01T12:00:00.123Z",
  "level": "DEBUG",
  "service": "resumes",
  "message": "AI service request",
  "endpoint": "/generate",
  "duration_ms": 1234
}
```

## Structured Log Format

All logs should follow this JSON structure:

```json
{
  "timestamp": "ISO8601 timestamp",
  "level": "ERROR|WARN|INFO|DEBUG",
  "service": "service-name",
  "message": "Human-readable message",
  "context": {
    // Additional context fields
    "job_id": "uuid",
    "user_id": "uuid",
    "duration_ms": 1234,
    "error": "error message",
    "error_type": "transient|permanent"
  }
}
```

### Standard Context Fields

- `job_id`: UUID of the job/operation
- `user_id`: UUID of the user
- `request_id`: UUID for request tracing
- `duration_ms`: Operation duration in milliseconds
- `error`: Error message
- `error_type`: Type of error (transient/permanent)
- `retry_count`: Number of retry attempts
- `status`: Current status

## Log Rotation

### Strategy
- **Rotation**: Daily at midnight UTC
- **Retention**: 30 days
- **Compression**: Gzip old logs
- **Max File Size**: 100MB per file (before rotation)

### Implementation

**Go (using slog + file rotation):**
```go
// Use lumberjack or file-rotatelogs for rotation
import "gopkg.in/natefinch/lumberjack.v2"

handler := &lumberjack.Logger{
    Filename:   "logs/app.log",
    MaxSize:    100, // megabytes
    MaxBackups: 30,  // days
    MaxAge:     30,  // days
    Compress:   true,
}
```

**Python (using logging.handlers):**
```python
from logging.handlers import TimedRotatingFileHandler

handler = TimedRotatingFileHandler(
    filename='logs/app.log',
    when='midnight',
    interval=1,
    backupCount=30,
    encoding='utf-8'
)
```

## Service-Specific Logging

### Resume Service
- **File**: `logs/resumes.log`
- **Error File**: `logs/resumes.error.log`
- **Key Events**:
  - Job enqueued
  - Job started processing
  - Job completed
  - Job failed
  - Job retrying
  - Dead letter queue entries

### Job Application Service
- **File**: `logs/job-applications.log`
- **Error File**: `logs/job-applications.error.log`
- **Key Events**:
  - Application created
  - Status changes
  - Scraping operations
  - Cover letter generation

### Main Application
- **File**: `logs/app.log`
- **Error File**: `logs/app.error.log`
- **Key Events**:
  - Server startup/shutdown
  - Request handling
  - Database operations
  - Authentication events

## Docker Integration

### Volume Mounts
```yaml
volumes:
  - ./logs:/app/logs
```

### Container Logs
- Log to both file and stdout
- Docker captures stdout/stderr automatically
- File logs provide persistence across container restarts

## Security Considerations

### Never Log
- API keys or secrets
- Passwords or tokens
- Full user personal information
- Complete request/response bodies (log summaries only)
- Credit card numbers or payment info

### Sanitization
- Replace sensitive values with `***REDACTED***`
- Log only necessary identifiers (user_id, not email)
- Truncate long strings (e.g., job descriptions)

## Performance

### Async Logging
- Use buffered writes
- Don't block main operations
- Batch log writes when possible

### Log Volume
- Use appropriate log levels
- Don't log in tight loops
- Aggregate similar events

## Monitoring and Alerting

### Log Aggregation
- Consider centralized logging (ELK, Loki, CloudWatch)
- Parse JSON logs for metrics
- Track error rates, job success rates, durations

### Alerts
- High error rate (> 5% of requests)
- Critical errors (authentication failures, DB connection loss)
- Job failure rate (> 10% of jobs)
- Disk space usage (> 80%)

## Implementation Checklist

- [ ] Set up log directory structure
- [ ] Configure log rotation
- [ ] Implement structured JSON logging
- [ ] Add service-specific log files
- [ ] Set up error log files
- [ ] Configure Docker volume mounts
- [ ] Add log level configuration (env var)
- [ ] Implement log sanitization
- [ ] Add request ID tracing
- [ ] Set up log monitoring (optional)

## Future Enhancements

1. **Centralized Logging**: Send logs to ELK stack or similar
2. **Log Analytics**: Dashboard for log analysis
3. **Real-time Alerts**: Alert on error patterns
4. **Distributed Tracing**: Add trace IDs across services
5. **Log Sampling**: Sample DEBUG logs in production

## Examples

### Go (slog) Example
```go
logger.Info("Resume generation completed",
    slog.String("job_id", jobID),
    slog.String("user_id", userID),
    slog.Int("duration_ms", duration),
    slog.Int64("file_size", fileSize),
)
```

### Python (logging) Example
```python
logger.info("Resume generation completed", extra={
    "job_id": job_id,
    "user_id": user_id,
    "duration_ms": duration,
    "file_size": file_size,
})
```

## References

- [Go slog documentation](https://pkg.go.dev/log/slog)
- [Python logging documentation](https://docs.python.org/3/library/logging.html)
- [12-Factor App: Logs](https://12factor.net/logs)

