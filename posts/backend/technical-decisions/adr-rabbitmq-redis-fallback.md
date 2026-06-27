# ADR-002: RabbitMQ with Redis Fallback

## Context
We need a reliable message queue for async job processing. Workers must be able to consume jobs even if RabbitMQ is down.

## Decision
Use RabbitMQ as primary queue, Redis as fallback. Server checks RabbitMQ availability and falls back to Redis if needed.

## Rationale

### RabbitMQ Primary
- **Features**: Dead letter queues, exchanges, routing
- **Reliability**: Enterprise-grade message broker
- **Scalability**: Handles high throughput
- **Management**: Web UI for monitoring

### Redis Fallback
- **Availability**: Already in infrastructure
- **Simplicity**: LPush/BRPop pattern
- **Performance**: Fast for simple queue operations
- **Graceful Degradation**: System continues working

## Consequences

### Pros
- High availability (works even if RabbitMQ down)
- Workers can still process jobs
- Graceful degradation
- No single point of failure

### Cons
- More complex code (two queue implementations)
- Need to handle both systems
- Potential inconsistency if both used simultaneously
- Redis less feature-rich than RabbitMQ

## Status
Accepted - 2024-01-15

## Alternatives Considered
- RabbitMQ only: Single point of failure
- Redis only: Missing advanced features
- Kafka: Overkill for our use case
