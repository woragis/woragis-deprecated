# Backend Logging Strategy

## Overview
Comprehensive logging strategy for the entire Go backend (still to be fully implemented per LOGGING_STRATEGY.md).

## Key Points

### Current State
- Logging strategy documented in `LOGGING_STRATEGY.md`
- Structured logging plan defined
- Implementation partially done (workers have logging)

### Planned Architecture
- Structured JSON logging
- Log levels: ERROR, WARN, INFO, DEBUG
- Log rotation: daily, 30-day retention
- Service-specific log files
- Separate error log files

### Log File Structure
```
logs/
  app.log                    # Main application log
  app.error.log              # Errors only
  resumes.log                # Resume domain
  resumes.error.log
  job-applications.log       # Job application domain
  job-applications.error.log
  archive/                   # Rotated logs
    2024/01/
      app.2024-01-01.log.gz
```

### Structured Log Format
- JSON format for easy parsing
- Standard fields: timestamp, level, service, message
- Context field: additional structured data
- Standard context fields: job_id, user_id, request_id, duration_ms, error, error_type

### Log Levels Strategy
- **ERROR**: System failures, critical errors, auth failures, DB connection failures
- **WARN**: Recoverable errors, retries, rate limiting, degraded performance
- **INFO**: Business events, job lifecycle, service startup/shutdown, successful operations
- **DEBUG**: Detailed steps, request/response details, internal state (development only)

### Service-Specific Logging
- Each domain can have dedicated log file
- Worker-specific logging (translation, resume, job-application)
- Main application logging for HTTP requests, middleware

## Implementation Details

### Logging Framework
- Go `slog` (structured logging)
- File rotation: lumberjack or file-rotatelogs
- Docker integration: log to both file and stdout

### Security Considerations
- Never log: API keys, passwords, tokens, full user data, payment info
- Sanitization: Replace sensitive values with `***REDACTED***`
- Log only identifiers (user_id, not email)

### Performance
- Async/buffered logging
- Don't block main operations
- Batch log writes when possible

## Potential Improvements
- **Immediate**: Implement file-based logging per LOGGING_STRATEGY.md
- Add request ID/trace ID middleware for distributed tracing
- Implement log aggregation (ELK, Loki, CloudWatch)
- Add log analytics dashboard
- Implement real-time alerts on error patterns
- Add distributed tracing support
- Implement log sampling for DEBUG logs in production
- Add metrics extraction from logs (error rates, durations)
- Implement log compression for archived logs
- Add log search/filtering capabilities
- Implement log streaming for real-time monitoring
- Add correlation between related logs (request flow tracking)

