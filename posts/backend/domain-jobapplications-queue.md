# Job Applications Domain - Queue Strategy

## Overview
How job application processing jobs are managed via Redis queue.

## Key Points

### Queue Architecture
- Redis-based queue (similar to translation worker)
- Queue key: `job-applications:queue`
- Job prefix: `job-applications:job:{job_id}`
- 24-hour TTL for job data

### Job Lifecycle
- Enqueue: Job created, stored in Redis, added to queue
- Dequeue: Worker picks up job (blocking pop with timeout)
- Processing: Worker processes application
- Complete: Job removed from storage
- Failed: Job marked as failed (could implement dead letter queue)

### Job Data Structure
- ID, UserID
- CompanyName, JobTitle, JobUrl
- Website (linkedin, indeed, etc.)
- Location
- Additional metadata

### Integration with Worker
- Job Application Worker (Node.js) processes jobs
- Rate limiting handled at orchestration level
- Jobs re-enqueued if rate limit reached

### Error Handling
- MarkJobFailed stores error message
- Failed jobs could go to dead letter queue
- Retry logic at worker level

## Potential Improvements
- Add priority queue support (urgent applications first)
- Implement dead letter queue for permanently failed jobs
- Add job status tracking (pending, processing, applied, failed, retrying)
- Implement job timeout mechanism
- Add job cancellation support
- Support job scheduling (delayed applications)
- Add job deduplication (same URL + user)
- Implement job batching
- Add metrics/observability
- Support job retry with exponential backoff
- Add job result caching
- Implement job dependency management (wait for resume generation)
- Add job website-specific queues

