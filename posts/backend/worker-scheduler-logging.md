# Scheduler Worker Logging Strategy

## Overview
Logging approach for the scheduler worker.

## Key Points

### Logging Framework
- Go `slog` (structured logging)
- Context-aware logging
- Error-focused logging

### Log Events

#### Worker Lifecycle
- Worker start (implicit via ticker start)
- Worker shutdown (context cancellation)

#### Schedule Processing
- ListDue errors (with error details)
- Execute failures (with schedule_id and error)

### Log Levels
- **ERROR**: ListDue failures, Execute failures
- No INFO/WARN logging currently

### Log Context Fields
- `schedule_id`: Schedule identifier (for execute failures)
- `error`: Error message/details

## Implementation Details

### Structured Logging
- Uses slog's structured fields
- Error logging includes schedule context
- Minimal logging (only errors)

## Potential Improvements
- Add INFO logs for successful executions
- Add schedule processing metrics (schedules processed, execution time)
- Add request ID/trace ID for tracking
- Add schedule details logging (type, due time, execution duration)
- Implement log rotation
- Add separate error log file
- Add performance metrics (list duration, execution duration)
- Add schedule execution statistics
- Implement log aggregation integration
- Add correlation IDs between related logs
- Add schedule execution history logging
- Support log sampling for high-volume scenarios
- Add schedule health monitoring logs

