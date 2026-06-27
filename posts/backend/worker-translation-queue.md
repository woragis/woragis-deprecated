# Translation Worker Queue Strategy

## Overview
How the translation worker manages job queuing using Redis.

## Key Points

### Queue Architecture
- Redis-based queue using LPush/BRPop pattern
- Job storage separated from queue (job data in Redis keys, IDs in list)
- Queue key: `translations:queue`
- Job prefix: `translations:job:{job_id}`

### Enqueue Process
- Job ID auto-generation if not provided
- Job data stored as JSON with 24-hour TTL
- Uses LPush for FIFO queue (can be changed to priority-based)
- Error handling for Redis unavailability

### Dequeue Process
- Blocking pop with configurable timeout (5 seconds default)
- BRPop for efficient blocking wait
- Job retrieval by ID after dequeue
- Handles nil/empty queue gracefully

### Job Lifecycle
- Enqueue → Processing → Complete/Failed
- Job completion: removes job data from Redis
- Job failure: marks as failed (could implement dead letter queue)
- No retry mechanism at queue level (worker handles retries)

## Implementation Details

### Redis Operations
- `LPush`: Add job ID to queue
- `BRPop`: Blocking pop with timeout
- `Set`: Store job data with TTL
- `Get`: Retrieve job data
- `Del`: Remove job on completion

### Error Handling
- Domain-specific errors (ErrCodeJobQueueFailure, ErrCodeNotFound)
- Redis unavailability handling
- Job not found scenarios

## Potential Improvements
- Add priority queue support (different queues per priority)
- Implement dead letter queue for failed jobs
- Add job retry mechanism at queue level
- Add job status tracking (pending, processing, completed, failed)
- Implement job timeout mechanism
- Add metrics/observability (queue length, processing time)
- Support job cancellation
- Add job batching for bulk translations

