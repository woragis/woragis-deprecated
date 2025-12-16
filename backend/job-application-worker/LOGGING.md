# Structured Logging - Job Application Worker

## Overview

The Job Application Worker uses structured JSON logging for consistent, searchable logs across all environments.

## Features

- **JSON format** - Always JSON for easy parsing and aggregation
- **Automatic service name** - All logs include `service: "job-application-worker"`
- **Trace ID support** - Distributed tracing via context
- **Structured fields** - All log entries use key-value pairs
- **ISO 8601 timestamps** - Standard timestamp format

## Configuration

### Environment Variables

```bash
# Environment (affects log level)
ENV=development  # or "production"
```

### Default Behavior

- **Development:** JSON format, all levels (including debug)
- **Production:** JSON format, info/warn/error only

## Log Format

### All Environments (JSON)
```json
{
  "timestamp": "2024-01-15T10:30:45.123456789Z",
  "level": "info",
  "service": "job-application-worker",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Job application completed",
  "jobId": "123e4567-e89b-12d3-a456-426614174000",
  "user_id": "6ad0d828-1234-5678-90ab-cdef12345678",
  "website": "linkedin.com",
  "status": "success"
}
```

## Usage

### Basic Logging

```javascript
import { logger, setTraceId, getTraceId } from './utils/logger.js';

// Info log
logger.info('Job application completed', {
  jobId: job.id,
  user_id: job.userId,
  website: job.website,
  status: 'success'
});

// Error log with context
logger.error('Error processing job', {
  jobId: job.id,
  user_id: job.userId,
  error: err.message,
  stack: err.stack
});
```

### Trace ID Support

```javascript
import { logger, setTraceId } from './utils/logger.js';

// Set trace_id from message headers
setTraceId(message.headers?.['x-trace-id'] || '');

// All subsequent logs will include trace_id
logger.info('Processing job application', { jobId: job.id });
```

### Debug Logging

```javascript
// Debug logs only appear in development
logger.debug('Selector cache hit', {
  selector: selector,
  cacheKey: cacheKey
});
```

## Log Storage

- **Development:** stdout (JSON format)
- **Production:** stdout (JSON format, collected by Kubernetes/log aggregator)

## Best Practices

1. **Use structured fields** - Always use objects, not string interpolation
   ```javascript
   // Good
   logger.info('Job processed', { jobId, status, duration_ms });
   
   // Bad
   logger.info(`Job ${jobId} processed with status ${status}`);
   ```

2. **Include context** - Add relevant fields to help debugging
   ```javascript
   logger.error('Scraping failed', {
     jobId: job.id,
     website: job.website,
     attempt: retryCount,
     error: err.message,
     selector: failedSelector
   });
   ```

3. **Never log sensitive data** - No passwords, tokens, or PII

4. **Use appropriate levels**
   - `debug`: Detailed diagnostic info (dev only)
   - `info`: General operational messages
   - `warn`: Warnings, retries, degraded functionality
   - `error`: Errors, failures, exceptions

## Integration with Log Aggregation

In production, logs go to stdout and are collected by:
- Kubernetes → Log aggregator (Fluentd/Fluent Bit) → Database/ELK/Loki

See main backend TODO.md for production log aggregation setup.
