# Event-Driven Architecture: Message Queue Patterns

## Overview
How we use message queues to implement event-driven patterns for asynchronous job processing.

## Key Points

### Event-Driven Patterns
- Job publishing (fire-and-forget)
- Job consumption (workers)
- Event sourcing (future)
- Event streaming (future)

### Message Queue Patterns
- Producer-Consumer
- Work queues
- Dead letter queues
- Priority queues (future)

## Implementation Details

### Job Publishing
```go
// Server publishes job
err := queue.Enqueue(ctx, Job{
    ID: uuid.New(),
    Type: "translation",
    Data: jobData,
})
```

### Job Consumption
```go
// Worker consumes job
job, err := queue.Dequeue(ctx)
if err != nil {
    return err
}
// Process job
```

### Dead Letter Queues
- Failed jobs → DLQ
- Manual reprocessing
- Error analysis
- Retry logic

## Benefits
- Decoupling
- Scalability
- Fault tolerance
- Async processing

## Challenges
- Eventual consistency
- Error handling
- Message ordering
- Duplicate handling

## Lessons Learned
- Message queues enable event-driven architecture
- Dead letter queues crucial
- Error handling important
- Monitoring essential

## Future Improvements
- Event sourcing
- Event streaming
- Saga pattern
- Event replay
