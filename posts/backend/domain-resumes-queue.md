# Resumes Domain - Queue & Job Management

## Overview
How resume generation jobs are managed via Redis queue.

## Key Points

### Queue Architecture
- Redis-based queue with advanced features
- Job status tracking (pending, processing, completed, failed, retrying, dead_letter)
- Event publishing via Pub/Sub
- 7-day TTL for job data
- Dead letter queue support

### Job Lifecycle
- **Enqueue**: Job created, status set to "pending", stored in Redis
- **Dequeue**: Worker picks up job (not standard dequeue, worker polls)
- **Processing**: Status updated to "processing"
- **Completion**: Status updated to "completed" with result data
- **Failure**: Status updated to "failed" with error details, retry logic

### Status Management
- UpdateJobStatus method handles status transitions
- Tracks retry count, error messages, error types
- Stores result data (output path, file info, tags, duration)
- Event publishing on status changes

### Event Publishing
- Publishes events to Redis Pub/Sub channel: `resumes:events`
- Events include: job_id, status, timestamp, error (if failed), result (if completed)
- Enables real-time status tracking for frontend

### Error Handling
- Error classification (transient vs permanent)
- Retry logic with exponential backoff
- Max retries: 3 (configurable)
- Dead letter queue for permanently failed jobs

## Potential Improvements
- Implement standard dequeue operation
- Add priority queue support
- Implement job timeout mechanism
- Add job cancellation support
- Implement job result caching
- Add metrics/observability
- Support job scheduling
- Add job deduplication
- Implement job batching
- Add job preview/snapshot

