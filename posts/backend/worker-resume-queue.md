# Resume Worker Queue Strategy

## Overview
How the resume worker manages job queuing using Redis with advanced features like status tracking and events.

## Key Points

### Queue Architecture
- Redis-based queue with job status tracking
- Job storage with 7-day TTL (longer than translation worker)
- Queue key: `resumes:queue`
- Job prefix: `resumes:job:{job_id}`
- Dead letter queue: `resumes:dead-letter:queue`
- Events channel: `resumes:events` (Pub/Sub)

### Enqueue Process
- Auto-generates job ID if not provided
- Sets initial status to "pending"
- Default max retries: 3 (configurable)
- Stores job data as JSON with 7-day TTL
- Uses LPush for FIFO queue

### Job Status Management
- Status transitions: pending → processing → completed/failed/retrying/dead_letter
- UpdateJobStatus method for status changes
- Tracks retry count, error messages, error types
- Stores result data (output path, file name, file size, tags, duration)

### Event Publishing
- Publishes job status events to Redis Pub/Sub channel
- Events include: job_id, status, timestamp, error (if failed), result (if completed)
- Enables real-time status tracking for frontend

### Dead Letter Queue
- Failed jobs after max retries go to dead letter queue
- Separate tracking for permanently failed jobs
- Enables manual inspection and retry

## Implementation Details

### Redis Operations
- `LPush`: Add job ID to queue
- `Set`: Store job data with 7-day TTL
- `Get`: Retrieve job data
- `Del`: Remove job (not used, relies on TTL)
- `Publish`: Publish status events

### Job Data Structure
- ID, UserID, JobApplicationID
- JobDescription, JobTitle, Language
- Status, RetryCount, MaxRetries
- LastError, LastErrorType, LastErrorAt
- Result (output path, file name, file size, tags, duration)
- CreatedAt, UpdatedAt timestamps

### Error Classification
- Error type tracking (transient vs permanent)
- Retry logic with exponential backoff
- Error message preservation

## Potential Improvements
- Implement actual dequeue operation (currently worker polls differently)
- Add priority queue support
- Implement job timeout mechanism
- Add job cancellation support
- Implement job result caching
- Add metrics/observability (queue length, processing time, success rate)
- Support job scheduling (delayed resume generation)
- Add job deduplication
- Implement job batching
- Add job preview/snapshot before full generation
- Implement partial resume generation (incremental updates)

