# Scheduler Worker Queue Strategy

## Overview
How the scheduler worker processes scheduled tasks without a traditional queue.

## Key Points

### Scheduler Pattern (Not Queue-Based)
- Time-based polling instead of queue
- Periodically checks for due schedules
- Executes schedules directly (no enqueue/dequeue)
- Interval-based: checks every minute (configurable)

### Processing Flow
1. Start ticker with interval (default: 1 minute)
2. Process due schedules immediately
3. Wait for ticker interval
4. Repeat: check and process due schedules
5. Stop on context cancellation

### Schedule Execution
- Query database for due schedules (ListDue)
- Execute each schedule sequentially
- Log errors but continue processing
- No queue involved (direct execution)

### Error Handling
- ListDue errors: log and return (skip this cycle)
- Execute errors: log per schedule, continue with others
- Graceful shutdown via context cancellation

## Implementation Details

### Database Queries
- `ListDue(now)`: Get schedules due before/at now
- Execute schedules via service

### Integration Points
- Scheduler domain service
- Context for cancellation
- Ticker for interval management

## Potential Improvements
- Add queue-based scheduling (enqueue jobs to execute)
- Implement priority scheduling
- Add schedule retry mechanism
- Implement schedule deduplication
- Add schedule locking (prevent concurrent execution)
- Support distributed scheduling (multiple workers)
- Add schedule execution timeout
- Implement schedule result tracking
- Add schedule execution history
- Support recurring schedules (cron-like)
- Add schedule dependency management
- Implement schedule batching
- Add schedule execution metrics
- Support schedule cancellation
- Add schedule execution queue (for long-running tasks)

