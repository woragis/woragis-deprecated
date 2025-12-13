# Logging Strategy for Woragis Backend

## Overview

This document outlines the comprehensive logging strategy for all Woragis backend services and workers, focusing on structured logging, log rotation, distributed tracing, and scalability.

## Architecture Overview

The Woragis backend consists of:

**Main Services:**
- **Main Server** (Go/Fiber) - Main API server
- **AI Service** (Python/FastAPI) - AI/LLM service
- **Creative Service** (Python/FastAPI) - Image/diagram/video generation

**Workers:**
- **Resume Worker** (Python) - Resume generation from queue
- **Translation Worker** (Go) - Translation processing
- **Job Application Worker** (Go + Node.js) - Job application scraping and processing
- **WhatsApp Worker** (Go) - WhatsApp notifications with leader election

## Principles

1. **Structured Logging**: Use JSON format for easy parsing and analysis
2. **Log Levels**: Appropriate use of ERROR, WARN, INFO, and DEBUG
3. **Log Rotation**: Automatic rotation to prevent disk space issues
4. **Separation of Concerns**: Separate log files for different services and log levels
5. **Performance**: Async/buffered logging to avoid blocking operations
6. **Security**: Never log sensitive data (API keys, passwords, full user data)

## Directory Structure

```
backend/
  server/
    logs/
      app.log                          # Main server log (all levels)
      app.error.log                    # Main server errors only
      resume-worker.log                # Resume worker log
      resume-worker.error.log          # Resume worker errors
      translation-worker.log           # Translation worker log
      translation-worker.error.log     # Translation worker errors
      job-application-worker.log       # Job application worker log
      job-application-worker.error.log # Job application worker errors
      whatsapp-worker.log              # WhatsApp worker log
      whatsapp-worker.error.log       # WhatsApp worker errors
      archive/                         # Rotated logs
        2024/
          01/
            app.2024-01-01.log.gz
            resume-worker.2024-01-01.log.gz
  creative-service/
    logs/
      creative-service.log
      creative-service.error.log
  ai-service/
    logs/
      ai-service.log
      ai-service.error.log
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

## Service-Specific Logging Strategies

### 1. Main Server (Go/Fiber)

**Current Implementation:** Uses `slog` with JSON handler in production, text handler in development.

**Log Files:**
- `logs/app.log` - All levels
- `logs/app.error.log` - Errors only

**Key Events to Log:**
- Server startup/shutdown with configuration
- HTTP request/response (with request_id)
- Database connection status
- Redis connection status
- Authentication/authorization events
- Domain-specific operations (posts, job applications, etc.)

**Best Practices:**
```go
// Use structured logging with context
logger.Info("request completed",
    slog.String("method", c.Method()),
    slog.String("path", c.Path()),
    slog.String("request_id", requestID),
    slog.Int("status", statusCode),
    slog.Int("duration_ms", duration),
)

// Include correlation IDs
logger.Error("database operation failed",
    slog.String("request_id", requestID),
    slog.String("user_id", userID),
    slog.String("operation", "create_post"),
    slog.Any("error", err),
)
```

**Improvements Needed:**
- Add file rotation with lumberjack
- Add request ID middleware for correlation
- Separate domain-specific log files
- Add performance metrics logging

### 2. Resume Worker (Python)

**Current Implementation:** ✅ Well-implemented with JSON formatter, rotation, and structured logging.

**Log Files:**
- `logs/resume-worker.log` - All levels
- `logs/resume-worker.error.log` - Errors only

**Key Events:**
- Worker startup/shutdown
- Job dequeued from queue
- Job processing started
- AI service calls (with duration)
- Resume generation steps
- Job completed (with metrics)
- Job failed (with error classification)
- Retry attempts with backoff
- Dead letter queue entries

**Current Best Practices (Already Implemented):**
```python
# Structured logging with context
logger.info(
    "Resume generation completed",
    extra={
        "job_id": job_id,
        "user_id": user_id,
        "duration_ms": duration_ms,
        "file_size": file_size,
        "output_path": output_path,
    },
)

# Error classification
logger.error(
    "Resume generation failed",
    extra={
        "job_id": job_id,
        "error": error_msg,
        "error_type": error_type,  # transient/permanent
        "retry_count": retry_count,
    },
    exc_info=True,
)
```

**Additional Recommendations:**
- Add correlation IDs for tracing across services
- Log AI service response times separately
- Add metrics for queue depth and processing rate

### 3. Translation Worker (Go)

**Current Implementation:** Uses `slog` logger from `app/pkg/logger`.

**Log Files:**
- `logs/translation-worker.log` - All levels
- `logs/translation-worker.error.log` - Errors only

**Key Events to Log:**
- Worker startup/shutdown
- Translation job dequeued
- Translation processing started
- AI service calls for translation
- Translation completed
- Translation failed
- Retry attempts
- Queue depth monitoring

**Recommended Implementation:**
```go
// Add structured logging with worker context
logger.Info("translation job started",
    slog.String("job_id", jobID),
    slog.String("source_lang", sourceLang),
    slog.String("target_lang", targetLang),
    slog.String("entity_type", entityType),
)

logger.Info("translation completed",
    slog.String("job_id", jobID),
    slog.Int("duration_ms", duration),
    slog.Int("char_count", charCount),
)

logger.Error("translation failed",
    slog.String("job_id", jobID),
    slog.String("error", err.Error()),
    slog.String("error_type", classifyError(err)),
    slog.Int("retry_count", retryCount),
)
```

**Improvements Needed:**
- Add file rotation
- Add correlation IDs
- Add performance metrics
- Add queue monitoring

### 4. Job Application Worker (Go + Node.js)

**Current Implementation:**
- **Go Worker:** Uses `slog` logger
- **Node.js Worker:** Has basic JSON logger utility

**Log Files:**
- `logs/job-application-worker.log` - All levels
- `logs/job-application-worker.error.log` - Errors only

**Key Events to Log:**
- Worker startup/shutdown
- Job dequeued
- Website scraping started
- Selector cache hits/misses
- Scraping attempts and results
- Cover letter generation
- Application record created/updated
- Rate limiting events
- Website processing limits

**Go Worker Best Practices:**
```go
logger.Info("job application processing started",
    slog.String("job_id", jobID),
    slog.String("user_id", userID),
    slog.String("website", website),
    slog.String("company", companyName),
)

logger.Info("scraping completed",
    slog.String("job_id", jobID),
    slog.String("website", website),
    slog.Int("duration_ms", duration),
    slog.Bool("cache_hit", cacheHit),
)

logger.Warn("rate limit reached",
    slog.String("website", website),
    slog.Int("requests_today", count),
    slog.String("action", "re_enqueue"),
)
```

**Node.js Worker Best Practices:**
```javascript
// Enhance existing logger with more context
logger.info('Processing job application', {
  jobId: job.id,
  userId: job.userId,
  company: job.companyName,
  website: job.website,
  requestId: generateRequestId(),
});

logger.info('Scraping completed', {
  jobId: job.id,
  website: job.website,
  durationMs: duration,
  selectorCacheHit: cacheHit,
  fieldsExtracted: Object.keys(extracted).length,
});

logger.error('Scraping failed', {
  jobId: job.id,
  website: job.website,
  error: error.message,
  errorType: classifyError(error),
  retryCount: retryCount,
  stack: error.stack,
});
```

**Improvements Needed:**
- Add file rotation for Node.js worker
- Standardize logging format between Go and Node.js
- Add correlation IDs
- Add scraping performance metrics

### 5. WhatsApp Worker (Go)

**Current Implementation:** Uses `slog` with JSON handler.

**Log Files:**
- `logs/whatsapp-worker.log` - All levels
- `logs/whatsapp-worker.error.log` - Errors only

**Key Events to Log:**
- Worker startup/shutdown
- Leader election events (if enabled)
- WhatsApp connection status
- Notification queued
- Notification sent
- Notification failed
- Rate limiting events
- Queue depth

**Best Practices:**
```go
logger.Info("whatsapp worker started",
    slog.String("mode", mode), // standalone or leader-election
    slog.String("pod_name", podName),
)

logger.Info("became leader",
    slog.String("pod", podName),
    slog.String("lease_name", leaseName),
)

logger.Info("notification sent",
    slog.String("notification_id", notificationID),
    slog.String("recipient", recipient),
    slog.Int("duration_ms", duration),
)

logger.Error("notification failed",
    slog.String("notification_id", notificationID),
    slog.String("error", err.Error()),
    slog.String("error_type", classifyError(err)),
)
```

**Improvements Needed:**
- Add file rotation
- Add correlation IDs
- Add notification metrics
- Log leader election transitions

### 6. AI Service (Python/FastAPI)

**Current Implementation:** Basic logging with standard Python logging.

**Log Files:**
- `logs/ai-service.log` - All levels
- `logs/ai-service.error.log` - Errors only

**Key Events to Log:**
- Service startup/shutdown
- Chat requests (with agent, provider, model)
- Chat completions (with duration, token counts)
- Image generation requests
- Provider selection and fallbacks
- Rate limiting events
- API errors from providers

**Recommended Implementation:**
```python
# Enhanced structured logging
logger.info(
    "chat request",
    extra={
        "request_id": request_id,
        "agent": req.agent,
        "provider": req.provider,
        "model": req.model,
        "input_length": len(req.input),
        "has_system": bool(req.system),
    },
)

logger.info(
    "chat completed",
    extra={
        "request_id": request_id,
        "agent": agent_name,
        "provider": provider,
        "duration_ms": duration,
        "output_length": len(output_text),
        "tokens_used": tokens_used,  # if available
    },
)

logger.error(
    "provider error",
    extra={
        "request_id": request_id,
        "provider": provider,
        "error": str(e),
        "error_type": type(e).__name__,
    },
    exc_info=True,
)
```

**Improvements Needed:**
- Add structured JSON logging
- Add file rotation
- Add request ID middleware
- Add provider performance metrics
- Add token usage tracking

### 7. Creative Service (Python/FastAPI)

**Current Implementation:** Basic logging with standard Python logging.

**Log Files:**
- `logs/creative-service.log` - All levels
- `logs/creative-service.error.log` - Errors only

**Key Events to Log:**
- Service startup/shutdown
- Image generation requests (provider, style, size)
- Diagram generation requests (type, format)
- Video generation requests
- Provider selection
- Generation duration
- File sizes generated
- Provider errors and fallbacks

**Recommended Implementation:**
```python
logger.info(
    "image generation request",
    extra={
        "request_id": request_id,
        "provider": req.provider,
        "style": req.style,
        "size": req.size,
        "prompt_length": len(req.prompt),
    },
)

logger.info(
    "image generation completed",
    extra={
        "request_id": request_id,
        "provider": req.provider,
        "duration_ms": duration,
        "images_generated": len(results),
        "total_size_bytes": total_size,
    },
)

logger.error(
    "generation failed",
    extra={
        "request_id": request_id,
        "provider": req.provider,
        "error": str(e),
        "error_type": type(e).__name__,
    },
    exc_info=True,
)
```

**Improvements Needed:**
- Add structured JSON logging
- Add file rotation
- Add request ID middleware
- Add provider performance metrics
- Add cost tracking (if applicable)

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

## Distributed Tracing and Correlation IDs

### Request ID Propagation

All services should propagate and log correlation IDs for request tracing:

**Implementation Pattern:**

1. **Main Server (Go):**
```go
// Middleware to generate/read request ID
func RequestIDMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        requestID := c.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("X-Request-ID", requestID)
        c.Locals("request_id", requestID)
        return c.Next()
    }
}

// Use in handlers
requestID := c.Locals("request_id").(string)
logger.Info("processing request", slog.String("request_id", requestID))
```

2. **Python Services (FastAPI):**
```python
from fastapi import Request
import uuid

@app.middleware("http")
async def add_request_id(request: Request, call_next):
    request_id = request.headers.get("X-Request-ID") or str(uuid.uuid4())
    request.state.request_id = request_id
    response = await call_next(request)
    response.headers["X-Request-ID"] = request_id
    return response

# Use in endpoints
request_id = request.state.request_id
logger.info("processing request", extra={"request_id": request_id})
```

3. **Workers:**
```python
# Include job_id as correlation ID
logger.info("processing job", extra={
    "job_id": job_id,
    "correlation_id": job.get("correlation_id", job_id),
})
```

### Trace Context

For cross-service tracing, include:
- `request_id` or `correlation_id` - Unique identifier for the request/job
- `parent_request_id` - If this is a child request
- `service_name` - Name of the service logging
- `operation` - The operation being performed

## Performance

### Async Logging
- Use buffered writes
- Don't block main operations
- Batch log writes when possible

**Python Implementation:**
```python
from logging.handlers import QueueHandler, QueueListener
import queue

log_queue = queue.Queue(-1)
queue_handler = QueueHandler(log_queue)
queue_listener = QueueListener(log_queue, main_handler, error_handler)
queue_listener.start()
```

**Go Implementation:**
```go
// slog handlers are already non-blocking, but for high-volume scenarios:
// Consider using a buffered channel with a separate goroutine
```

### Log Volume Management
- Use appropriate log levels
- Don't log in tight loops
- Aggregate similar events
- Sample DEBUG logs in production (e.g., 1 in 100)
- Use rate limiting for repeated errors

**Sampling Example (Python):**
```python
import random

def should_log_debug():
    return random.random() < 0.01  # 1% sampling

if should_log_debug():
    logger.debug("detailed debug info", extra={...})
```

### Performance Metrics Logging

Log key performance indicators:

**For Workers:**
- Queue depth
- Processing rate (jobs/second)
- Average processing time
- Error rate
- Retry rate

**For Services:**
- Request rate
- Response time percentiles (p50, p95, p99)
- Error rate
- Database query time
- External API call time

**Example:**
```python
# Periodic metrics logging
logger.info("worker metrics",
    extra={
        "queue_depth": queue_depth,
        "jobs_processed_last_minute": count,
        "avg_processing_time_ms": avg_time,
        "error_rate": error_rate,
    },
)
```

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

### Infrastructure
- [x] Resume Worker: Structured JSON logging with rotation ✅
- [ ] Main Server: Add file rotation with lumberjack
- [ ] Translation Worker: Add file rotation
- [ ] Job Application Worker (Go): Add file rotation
- [ ] Job Application Worker (Node.js): Add file rotation
- [ ] WhatsApp Worker: Add file rotation
- [ ] AI Service: Add structured JSON logging with rotation
- [ ] Creative Service: Add structured JSON logging with rotation

### Cross-Cutting Concerns
- [ ] Add request ID middleware to all services
- [ ] Implement correlation ID propagation
- [ ] Add log sanitization utilities
- [ ] Configure log level via environment variables
- [ ] Set up Docker volume mounts for all services
- [ ] Add health check endpoints that include log status

### Service-Specific
- [ ] Main Server: Add domain-specific log files
- [ ] Translation Worker: Add structured logging with context
- [ ] Job Application Worker: Standardize logging between Go and Node.js
- [ ] WhatsApp Worker: Add leader election event logging
- [ ] AI Service: Add provider performance metrics
- [ ] Creative Service: Add generation metrics

### Monitoring and Observability
- [ ] Set up log aggregation (ELK/Loki/CloudWatch)
- [ ] Create log analysis dashboards
- [ ] Set up alerts for error rates
- [ ] Add performance metrics extraction
- [ ] Implement distributed tracing (optional)

### Documentation
- [ ] Document log file locations
- [ ] Document log query patterns
- [ ] Create runbooks for common log analysis tasks
- [ ] Document alert thresholds

## Future Enhancements

1. **Centralized Logging**: Send logs to ELK stack, Loki, or CloudWatch
2. **Log Analytics**: Dashboard for log analysis with Grafana/Kibana
3. **Real-time Alerts**: Alert on error patterns, high error rates, slow requests
4. **Distributed Tracing**: Add OpenTelemetry or similar for full request tracing
5. **Log Sampling**: Sample DEBUG logs in production (1-10% sampling rate)
6. **Cost Tracking**: Log API costs for AI services
7. **Performance Baselines**: Establish and alert on performance degradation
8. **Anomaly Detection**: Detect unusual patterns in logs
9. **Log Retention Policies**: Different retention for different log levels
10. **Compliance Logging**: Audit logs for sensitive operations

## Quick Reference

### Log Levels Decision Tree

```
Is it a system failure or unrecoverable error?
├─ YES → ERROR
└─ NO → Is it recoverable but concerning?
    ├─ YES → WARN
    └─ NO → Is it a normal business event?
        ├─ YES → INFO
        └─ NO → Is it detailed debugging info?
            ├─ YES → DEBUG (only in dev)
            └─ NO → Don't log
```

### Standard Context Fields

Always include these when available:
- `timestamp` - ISO8601 format
- `level` - ERROR, WARN, INFO, DEBUG
- `service` - Service/worker name
- `message` - Human-readable message
- `request_id` / `correlation_id` - For tracing
- `user_id` - When user context exists
- `job_id` - For worker jobs
- `duration_ms` - For performance tracking
- `error` - Error message (for errors)
- `error_type` - Error classification

### Environment Variables

Standardize these across all services:
- `LOG_LEVEL` - DEBUG, INFO, WARN, ERROR (default: INFO)
- `LOG_DIR` - Directory for log files (default: `/app/logs`)
- `LOG_FORMAT` - json or text (default: json in prod, text in dev)
- `LOG_ROTATION_MAX_SIZE` - Max file size in MB (default: 100)
- `LOG_ROTATION_MAX_AGE` - Max age in days (default: 30)
- `LOG_ROTATION_MAX_BACKUPS` - Max backup files (default: 30)

## Error Handling and Classification

### Error Classification

Classify errors to determine retry behavior:

**Transient Errors (Retryable):**
- Network timeouts
- Temporary service unavailability
- Rate limiting (with backoff)
- Database connection errors
- External API timeouts

**Permanent Errors (Non-retryable):**
- Authentication failures
- Invalid input data
- Resource not found
- Permission denied
- Disk full / quota exceeded

**Implementation Example:**
```python
def classify_error(error: Exception) -> str:
    """Classify error as transient or permanent."""
    error_type = type(error).__name__
    error_msg = str(error).lower()
    
    # Permanent errors
    if any(term in error_msg for term in ["authentication", "unauthorized", "not found", "permission denied"]):
        return "permanent"
    
    # Transient errors (default)
    return "transient"
```

### Error Logging Best Practices

1. **Always include context:**
```python
logger.error("operation failed",
    extra={
        "operation": "generate_resume",
        "job_id": job_id,
        "user_id": user_id,
        "error": str(e),
        "error_type": type(e).__name__,
        "error_classification": classify_error(e),
    },
    exc_info=True,  # Include stack trace
)
```

2. **Log at appropriate levels:**
- ERROR: System failures, unrecoverable errors
- WARN: Recoverable errors, retries, degraded performance
- INFO: Normal operation, successful completions
- DEBUG: Detailed debugging information

3. **Avoid logging sensitive data:**
```python
# Bad
logger.info(f"User login: {email} with password {password}")

# Good
logger.info("User login", extra={"user_id": user_id, "email_domain": email.split("@")[1]})
```

## Examples

### Go (slog) Examples

**Basic Logging:**
```go
logger.Info("Resume generation completed",
    slog.String("job_id", jobID),
    slog.String("user_id", userID),
    slog.Int("duration_ms", duration),
    slog.Int64("file_size", fileSize),
)
```

**With Request ID:**
```go
logger.Info("request processed",
    slog.String("request_id", requestID),
    slog.String("method", method),
    slog.String("path", path),
    slog.Int("status", statusCode),
    slog.Int("duration_ms", duration),
)
```

**Error Logging:**
```go
logger.Error("operation failed",
    slog.String("operation", "generate_resume"),
    slog.String("job_id", jobID),
    slog.String("error", err.Error()),
    slog.String("error_type", classifyError(err)),
)
```

**With Context:**
```go
ctx := context.WithValue(context.Background(), "request_id", requestID)
logger.Info("processing request",
    slog.String("request_id", requestID),
    slog.String("user_id", userID),
)
```

### Python (logging) Examples

**Basic Structured Logging:**
```python
logger.info("Resume generation completed", extra={
    "job_id": job_id,
    "user_id": user_id,
    "duration_ms": duration,
    "file_size": file_size,
})
```

**With Request ID:**
```python
logger.info("request processed", extra={
    "request_id": request_id,
    "method": method,
    "path": path,
    "status_code": status_code,
    "duration_ms": duration,
})
```

**Error Logging with Classification:**
```python
logger.error("operation failed", extra={
    "operation": "generate_resume",
    "job_id": job_id,
    "user_id": user_id,
    "error": str(e),
    "error_type": type(e).__name__,
    "error_classification": classify_error(e),
    "retry_count": retry_count,
}, exc_info=True)
```

**Worker Job Processing:**
```python
logger.info("Job processing started", extra={
    "job_id": job_id,
    "user_id": user_id,
    "job_type": "resume_generation",
    "queue": "resumes:queue",
    "retry_count": retry_count,
})

logger.info("Job processing completed", extra={
    "job_id": job_id,
    "user_id": user_id,
    "duration_ms": duration,
    "status": "success",
})
```

### Node.js Examples

**Structured Logging:**
```javascript
logger.info('Job processing started', {
  jobId: job.id,
  userId: job.userId,
  website: job.website,
  company: job.companyName,
  requestId: generateRequestId(),
});

logger.error('Job processing failed', {
  jobId: job.id,
  error: error.message,
  errorType: error.constructor.name,
  stack: error.stack,
  retryCount: retryCount,
});
```

## Log Aggregation and Analysis

### Recommended Tools

1. **Local Development:**
   - File-based logging with rotation
   - `tail -f` for real-time monitoring
   - `grep`/`jq` for filtering JSON logs

2. **Production:**
   - **ELK Stack** (Elasticsearch, Logstash, Kibana)
   - **Loki + Grafana** (lightweight alternative)
   - **CloudWatch** (AWS)
   - **Datadog** (SaaS solution)

### Log Query Examples

**Using jq for JSON logs:**
```bash
# Filter errors
cat logs/app.log | jq 'select(.level == "ERROR")'

# Find slow requests
cat logs/app.log | jq 'select(.duration_ms > 1000)'

# Group by error type
cat logs/app.log | jq -r 'select(.level == "ERROR") | .error_type' | sort | uniq -c

# Track job success rate
cat logs/resume-worker.log | jq -r 'select(.message | contains("completed")) | .job_id' | wc -l
```

### Metrics Extraction

Extract metrics from logs for monitoring:

**Key Metrics:**
- Error rate per service
- Average response time
- Queue depth
- Job success/failure rate
- Retry rate
- External API call duration

**Example Prometheus Metrics from Logs:**
```python
# Log metrics that can be scraped
logger.info("metrics", extra={
    "metric_type": "counter",
    "metric_name": "jobs_processed_total",
    "value": 1,
    "labels": {"status": "success", "worker": "resume-worker"},
})
```

## References

- [Go slog documentation](https://pkg.go.dev/log/slog)
- [Python logging documentation](https://docs.python.org/3/library/logging.html)
- [12-Factor App: Logs](https://12factor.net/logs)

