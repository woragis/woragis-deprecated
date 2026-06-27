# Worker Lifecycle: Startup, Processing, Shutdown

## Overview
How workers handle startup, job processing, and graceful shutdown.

## Key Points

### Lifecycle Stages
1. **Startup**: Initialize connections, start health server
2. **Processing**: Consume jobs, process, acknowledge
3. **Shutdown**: Stop consuming, finish current jobs, close connections

### Graceful Shutdown
- Handle SIGTERM/SIGINT
- Stop accepting new jobs
- Finish current jobs
- Close connections
- Exit cleanly

## Implementation Details

### Startup
```go
func main() {
    // Load config
    cfg := config.Load()
    
    // Initialize logger
    logger := logger.New(cfg.Logger)
    
    // Connect to RabbitMQ
    conn := queue.Connect(cfg.RabbitMQ)
    
    // Start health server
    go startHealthServer(conn, logger)
    
    // Start consuming
    worker.Start(conn, logger)
}
```

### Graceful Shutdown
```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

<-sigChan
logger.Info("Shutting down...")

// Stop consuming
worker.Stop()

// Close connections
conn.Close()

logger.Info("Shutdown complete")
```

### Job Processing
```go
for {
    job, err := queue.Dequeue(ctx)
    if err != nil {
        continue
    }
    
    err = processJob(job)
    if err != nil {
        queue.Nack(job)
    } else {
        queue.Ack(job)
    }
}
```

## Benefits
- Clean startup
- Graceful shutdown
- Resource cleanup
- Better reliability

## Challenges
- Signal handling
- Job completion
- Connection cleanup
- Timeout handling

## Lessons Learned
- Graceful shutdown essential
- Signal handling important
- Job completion crucial
- Cleanup important

## Future Improvements
- Shutdown timeout
- Job cancellation
- Health check during shutdown
- Metrics for lifecycle
