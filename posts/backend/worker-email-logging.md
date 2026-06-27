# Email Worker Logging Strategy

## Overview
Logging approach for the email worker.

## Key Points

### Logging Framework
- Go `slog` (structured logging)
- Context-aware logging
- Error-focused logging

### Log Events

#### Worker Lifecycle
- Worker start (implicit via subscription)
- Worker shutdown (context cancellation)

#### Message Processing
- Invalid payload errors (with error details)
- Send failures (with user_id and error)

### Log Levels
- **ERROR**: Invalid payloads, send failures
- No INFO/WARN logging currently

### Log Context Fields
- `user_id`: User identifier (for send failures)
- `error`: Error message/details

## Implementation Details

### Structured Logging
- Uses slog's structured fields
- Error logging includes user context
- Minimal logging (only errors)

## Potential Improvements
- Add INFO logs for successful sends
- Add message processing metrics (throughput)
- Add request ID/trace ID for tracking
- Add email details logging (subject, recipient hash)
- Implement log rotation
- Add separate error log file
- Add performance metrics (send duration)
- Add email service health logging
- Implement log aggregation integration
- Add correlation IDs between related logs
- Add email delivery status logging
- Support log sampling for high-volume scenarios

