# Service Communication: HTTP vs Message Queues

## Overview
When to use HTTP/REST for synchronous communication vs message queues for asynchronous communication in our microservices architecture.

## Key Points

### Communication Patterns

#### HTTP/REST (Synchronous)
- Server ↔ Services (AI, Creative)
- Workers ↔ Services (for AI/Creative calls)
- Client ↔ Server (API requests)
- Real-time, request-response

#### Message Queues (Asynchronous)
- Server ↔ Workers
- Job processing
- Event-driven workflows
- Fire-and-forget, eventual consistency

## Implementation Details

### HTTP Communication
- FastAPI services (AI, Creative)
- RESTful endpoints
- JSON request/response
- Error handling and retries
- Timeout management

### Message Queue Communication
- RabbitMQ (primary)
- Redis (fallback)
- Job-based messaging
- Acknowledgment patterns
- Dead letter queues

## Decision Criteria

### Use HTTP When:
- Need immediate response
- Synchronous operation required
- Request-response pattern
- Low latency critical
- Simple request/response

### Use Message Queues When:
- Async processing acceptable
- Long-running operations
- Decoupling needed
- Fault tolerance required
- Event-driven architecture

## Examples

### HTTP: AI Service Call
```go
// Server calls AI Service synchronously
response, err := aiService.GenerateContent(ctx, request)
if err != nil {
    return err
}
// Use response immediately
```

### Message Queue: Translation Job
```go
// Server publishes job asynchronously
err := translationQueue.Enqueue(ctx, job)
if err != nil {
    return err
}
// Job processed later by worker
```

## Benefits
- Right tool for right job
- Performance optimization
- Fault tolerance
- Scalability

## Challenges
- Need to choose correctly
- Debugging distributed systems
- Error handling complexity
- Consistency management

## Lessons Learned
- HTTP for real-time, queues for async
- Message queues enable better scalability
- HTTP simpler for request-response
- Both patterns needed in microservices

## Future Improvements
- Service mesh for HTTP communication
- Event sourcing for complex workflows
- API gateway for HTTP routing
- Message queue monitoring
