# Resume Worker Logging Strategy

## Overview
Comprehensive logging strategy for the Python-based resume generation worker.

## Key Points

### Logging Framework
- Python `logging` module with structured JSON format
- Custom JSONFormatter for parseable logs
- Multiple handlers: main log, error log, console
- Log rotation: daily rotation with 30-day retention

### Log File Structure
- Main log: `resume-worker.log` (all levels)
- Error log: `resume-worker.error.log` (ERROR level only)
- Console output: stdout for Docker logs
- Log directory: `/app/logs` (configurable via LOG_DIR env var)

### JSON Log Format
- Structured JSON logs for easy parsing
- Standard fields: timestamp, level, service, message
- Context field: additional structured data
- Exception field: formatted exception info when present
- Timestamp: ISO 8601 format with milliseconds

### Log Events

#### Worker Lifecycle
- Worker start/stop
- Connection initialization (queue, database, AI service, translation helper)
- Resource cleanup on shutdown

#### Job Processing
- Job dequeued (with job_id, user_id, job_application_id)
- Job status updates (pending → processing → completed/failed)
- Job processing start
- Job processing completion (with duration, file_size, output_path)
- Job processing failure (with error type, error message, retry count)

#### Error Classification
- Transient errors (network issues, timeouts)
- Permanent errors (authentication, not found, permission denied, disk full)
- Error classification logged for retry decisions

#### Performance Metrics
- Job processing duration (milliseconds)
- File generation size
- Retry attempts with backoff timing

## Implementation Details

### Logging Configuration
- TimedRotatingFileHandler for daily rotation
- Backup count: 30 days
- Encoding: UTF-8
- Log levels: DEBUG (file), INFO (console), ERROR (error file)

### Context Extraction
- Extracts extra fields from log record
- Filters out internal Python logging fields
- Includes JSON-serializable values only
- Converts non-serializable values to strings

### Error Handling
- Comprehensive exception logging with stack traces
- Error type classification (transient vs permanent)
- Retry count tracking in logs
- Backoff calculation logging

## Potential Improvements
- Add request ID/trace ID for distributed tracing
- Add correlation IDs between related log entries
- Implement log sampling for high-volume scenarios
- Add structured metrics logging (success rate, average duration)
- Add queue metrics (queue length, wait time)
- Implement log aggregation integration (ELK, Loki, CloudWatch)
- Add performance profiling logs (memory usage, CPU time)
- Add AI service interaction logs (request/response times)
- Implement log alerting based on error patterns
- Add log search/filtering capabilities
- Add log compression for archived logs
- Implement log streaming for real-time monitoring

