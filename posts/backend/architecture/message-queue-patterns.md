# Message Queue Patterns: RabbitMQ + Redis Fallback

## Overview
How we implemented a resilient message queue system using RabbitMQ as primary and Redis as fallback for graceful degradation.

## Key Points

### Problem
- Need reliable message queue for async job processing
- Workers must process jobs even if RabbitMQ is down
- System should degrade gracefully, not fail completely

### Solution
- RabbitMQ as primary message queue
- Redis as fallback queue
- Server checks RabbitMQ availability
- Automatic fallback to Redis if RabbitMQ unavailable

### Implementation
- Queue interface abstraction
- RabbitMQ implementation (primary)
- Redis implementation (fallback)
- Health check integration
- Logging for queue system in use

## Implementation Details

### Queue Interface
```go
type Queue interface {
    Enqueue(ctx context.Context, job Job) error
    Dequeue(ctx context.Context) (Job, error)
    // ...
}
```

### RabbitMQ Queue
- AMQP 0.9.1 protocol
- Exchange: `woragis.tasks`
- Queues: `{domain}.queue`
- Dead letter exchange: `woragis.dlx`

### Redis Queue
- LPush/BRPop pattern
- Queue keys: `{domain}:queue`
- Job storage: `{domain}:job:{job_id}`
- TTL: 24 hours

### Fallback Logic
1. Check RabbitMQ connection
2. If available → use RabbitMQ
3. If unavailable → use Redis
4. Log which system is active

## Benefits
- High availability (99.9%+)
- Graceful degradation
- No single point of failure
- Workers continue processing

## Trade-offs
- More complex code (two implementations)
- Need to handle both systems
- Potential inconsistency if both used
- Redis less feature-rich than RabbitMQ

## Lessons Learned
- Interface abstraction enables flexibility
- Health checks crucial for fallback decisions
- Logging helps debug queue system selection
- Workers don't need to know which queue system

## Future Improvements
- Queue system metrics
- Automatic failover back to RabbitMQ
- Queue system comparison dashboard
- Message deduplication across systems
