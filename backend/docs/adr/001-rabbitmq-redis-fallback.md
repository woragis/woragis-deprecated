# ADR-001: RabbitMQ with Redis Fallback

## Status
Accepted

## Context
The backend requires asynchronous job processing for various operations (email sending, translations, resume generation, etc.). We need a reliable message queue system that can handle failures gracefully and ensure system availability even when the primary message queue is unavailable.

**Requirements:**
- Asynchronous job processing
- High availability
- Graceful degradation during failures
- No message loss
- Simple fallback mechanism

**Constraints:**
- Already using Redis for caching
- Need to minimize infrastructure complexity
- Must work with existing codebase patterns

## Decision
We will use **RabbitMQ as the primary message queue** with **Redis as a fallback queue** when RabbitMQ is unavailable.

**Implementation:**
- Server checks RabbitMQ availability on startup
- If RabbitMQ is available, use RabbitMQ queue
- If RabbitMQ is unavailable, fall back to Redis queue
- Health check reflects "degraded" status when using fallback
- System continues operating in degraded mode

**Queue Interface:**
- Abstract queue interface (`Queue` interface)
- Two implementations: `RabbitMQQueue` and `RedisQueue`
- Server selects implementation based on availability

## Consequences

### Positive
- ✅ **High Availability**: System continues operating during RabbitMQ outages
- ✅ **No Message Loss**: Messages are queued even when RabbitMQ is down
- ✅ **Simple Fallback**: Uses existing Redis infrastructure
- ✅ **Graceful Degradation**: Health check clearly indicates degraded state
- ✅ **Transparent to Workers**: Workers continue processing (if already connected to RabbitMQ)

### Negative
- ⚠️ **Different Characteristics**: Redis queue has different guarantees than RabbitMQ
  - No Dead Letter Queue (DLX) in Redis fallback
  - No message acknowledgments in Redis fallback
  - Different message ordering guarantees
- ⚠️ **Manual Promotion**: No automatic promotion back to RabbitMQ when it recovers
- ⚠️ **Monitoring**: Need to monitor fallback usage and alert when in degraded mode

### Neutral
- Redis is already in use for caching, so no additional infrastructure needed
- Code complexity increases slightly (queue abstraction, fallback logic)

## Implementation Details

### Queue Interface
```go
type Queue interface {
    EnqueueJob(ctx context.Context, job *Job) error
    // ... other methods
}
```

### Fallback Logic
```go
var translationQueue translationsdomain.Queue
if rabbitmqConn != nil {
    rabbitmqQueue, err := translationsdomain.NewRabbitMQQueue(rabbitmqConn)
    if err != nil {
        slogLogger.Warn("failed to create RabbitMQ queue, falling back to Redis")
        translationQueue = translationsdomain.NewRedisQueue(redisClient)
    } else {
        translationQueue = rabbitmqQueue
    }
} else {
    translationQueue = translationsdomain.NewRedisQueue(redisClient)
}
```

### Health Check
- Health check shows "degraded" status when using Redis fallback
- RabbitMQ check is optional (non-critical) in server health check
- Workers still require RabbitMQ (no fallback for workers)

## Alternatives Considered

### 1. RabbitMQ Only (No Fallback)
- **Pros**: Simpler, consistent behavior
- **Cons**: System unavailable during RabbitMQ outages
- **Rejected**: Doesn't meet high availability requirement

### 2. Redis Only (No RabbitMQ)
- **Pros**: Simpler infrastructure
- **Cons**: Missing RabbitMQ features (DLX, acknowledgments, routing)
- **Rejected**: Need RabbitMQ features for production reliability

### 3. Multiple RabbitMQ Instances (Clustering)
- **Pros**: High availability with RabbitMQ features
- **Cons**: More complex infrastructure, higher cost
- **Rejected**: Overkill for current scale, adds complexity

### 4. Kafka
- **Pros**: High throughput, built-in replication
- **Cons**: More complex, different paradigm, overkill for current needs
- **Rejected**: Too complex for current requirements

## Notes
- This pattern is used specifically for the server's job publishing
- Workers still require RabbitMQ (no fallback implemented for workers)
- Future enhancement: Automatic promotion back to RabbitMQ when it recovers
- Future enhancement: Monitor fallback usage and alert when in degraded mode

## Related ADRs
- [ADR-002: Standalone Workers Architecture](./002-standalone-workers.md) - Workers use RabbitMQ directly
- [ADR-005: Health Checks Implementation Strategy](./005-health-checks.md) - Health checks reflect degraded state
