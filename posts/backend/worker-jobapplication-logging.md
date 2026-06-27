# Job Application Worker Logging Strategy

## Overview
Logging approach for the Node.js job application worker.

## Key Points

### Logging Framework
- Custom logger utility (likely Winston or Pino)
- Structured logging with context fields
- Console output for Docker logs

### Log Events

#### Worker Lifecycle
- Worker start/stop
- Connection initialization (queue, database, scraper)
- Resource cleanup on shutdown

#### Job Processing
- Job dequeued (with job ID, company, website)
- Website rate limit checks
- Job processing start
- Job processing completion (status: applied/failed)
- Job re-enqueuing (when rate limit reached)

#### Application Flow
- User profile data fetched
- Cover letter generation
- Scraping operations (Playwright)
- Form filling and submission
- Status updates (processing → applied/failed)

### Log Context Fields
- `jobId`: Job application identifier
- `company`: Company name
- `website`: Job website (linkedin, indeed, etc.)
- `userId`: User identifier
- `status`: Application status
- `error`: Error message when failures occur

### Error Logging
- Errors during scraping
- Errors during cover letter generation
- Errors during form submission
- Network errors
- Playwright errors
- Full error stack traces

## Implementation Details

### Structured Logging
- Uses structured logger with context objects
- All logs include relevant context
- Error logging includes full error details and stack traces

### Log Levels
- **INFO**: Normal operations, job lifecycle
- **ERROR**: Failed operations, scraping errors
- **WARN**: Recoverable issues (rate limits, website issues)

## Potential Improvements
- Add request ID/trace ID for distributed tracing
- Add performance metrics (job processing duration, scraping time)
- Implement log rotation (file-based logging)
- Add separate error log file
- Add structured error context (error type, retry count, website)
- Add job statistics logging (success rate, average duration per website)
- Add queue metrics (queue length, wait time)
- Implement log sampling for high-volume scenarios
- Add correlation IDs between related logs
- Add scraping metrics (selectors found, form fields filled, screenshots taken)
- Implement log aggregation integration
- Add AI service interaction logs (cover letter generation requests)
- Add rate limiting logs (current count, daily limit, reset time)

