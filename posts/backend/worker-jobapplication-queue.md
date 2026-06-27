# Job Application Worker Queue Strategy

## Overview
How the job application worker manages job queuing using Redis for automated job applications.

## Key Points

### Queue Architecture
- Redis-based queue using LPush/BRPop pattern
- Similar to translation worker pattern
- Queue key: `job-applications:queue`
- Job prefix: `job-applications:job:{job_id}`

### Enqueue Process
- Job ID auto-generation if not provided
- Job data stored as JSON with 24-hour TTL
- Uses LPush for FIFO queue
- Error handling for Redis unavailability

### Dequeue Process
- Blocking pop with configurable timeout (5 seconds default)
- BRPop for efficient blocking wait
- Job retrieval by ID after dequeue
- Handles nil/empty queue gracefully

### Job Lifecycle
- Enqueue → Processing → Complete/Failed
- Job completion: removes job data from Redis
- Job failure: marks as failed
- Rate limiting handled at orchestration level (not queue level)

### Job Data Structure
- ID, UserID
- CompanyName, JobTitle, JobUrl
- Website (e.g., "linkedin", "indeed")
- Location
- Additional metadata

## Implementation Details

### Redis Operations
- `LPush`: Add job ID to queue
- `BRPop`: Blocking pop with timeout
- `Set`: Store job data with TTL
- `Get`: Retrieve job data
- `Del`: Remove job on completion

### Rate Limiting Integration
- Rate limiting handled by Orchestrator (not queue)
- Website-based daily limits
- Jobs re-enqueued if rate limit reached
- Separate tracking for each job website

## Potential Improvements
- Add priority queue support (urgent applications first)
- Implement dead letter queue for permanently failed jobs
- Add job status tracking (pending, processing, applied, failed, retrying)
- Implement job timeout mechanism
- Add job cancellation support
- Support job scheduling (delayed applications)
- Add job deduplication (same URL + user)
- Implement job batching
- Add metrics/observability (queue length, processing time, success rate)
- Support job retry with exponential backoff
- Add job result caching (successful application patterns)
- Implement job dependency management (wait for resume generation)
- Add job website-specific queues (separate queues per website)

