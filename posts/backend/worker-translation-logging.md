# Translation Worker Logging Strategy

## Overview
Logging approach for the translation worker in Go.

## Key Points

### Logging Framework
- Uses Go's `slog` (structured logging)
- Context-aware logging throughout worker lifecycle
- Structured log format for parsing

### Log Events

#### Worker Lifecycle
- Worker start
- Worker stop (context cancelled or stop signal)
- Worker shutdown

#### Job Processing
- Job dequeued (with job ID, entity type, entity ID, language)
- Job processing start
- Job processing completion
- Job processing failure (with error details)

### Log Levels
- **INFO**: Normal operations, job lifecycle
- **ERROR**: Failed operations, job processing errors
- **WARN**: Recoverable issues (e.g., failed to mark job complete)

### Log Context Fields
- `jobId`: Translation job identifier
- `entityType`: Type of entity being translated
- `entityId`: ID of entity being translated
- `language`: Target language
- `error`: Error message when failures occur

### Error Logging
- Logs errors from queue operations
- Logs errors from service.ProcessTranslationJob
- Distinguishes between "job not found" (expected) vs real errors

## Implementation Details

### Structured Logging
- Uses slog's structured fields
- All logs include relevant context
- Error logging includes full error details

### Log Location
- Currently logs to stdout/stderr (Docker captures)
- Could be extended to file-based logging per LOGGING_STRATEGY.md

## Potential Improvements
- Add request ID/trace ID for distributed tracing
- Add performance metrics (job processing duration)
- Add structured error context (error type, retry count)
- Implement log rotation per LOGGING_STRATEGY.md
- Add separate error log file
- Add job statistics logging (success rate, average duration)
- Add queue metrics (queue length, wait time)
- Implement log sampling for high-volume scenarios
- Add correlation IDs between related logs

