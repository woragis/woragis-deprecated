# Trace ID Propagation: Distributed Tracing Basics

## Overview
How we implement trace ID propagation across services and workers for distributed tracing.

## Key Points

### Problem
- Requests span multiple services
- Need to correlate logs across services
- Debugging distributed systems is hard
- Need request tracing

### Solution
- Generate trace ID on request entry
- Propagate via context (Go) or headers (HTTP)
- Include in all logs
- Use for log correlation

## Implementation Details

### Trace ID Generation
```go
func GenerateTraceID() string {
    return uuid.New().String()
}
```

### Context Propagation (Go)
```go
ctx := context.WithValue(ctx, traceIDKey, traceID)
logger := logger.With("trace_id", traceID)
```

### HTTP Header Propagation
```go
req.Header.Set("X-Trace-ID", traceID)
```

### Logging with Trace ID
```go
logger.Info("Processing job",
    "trace_id", traceID,
    "job_id", jobID,
    "status", "processing",
)
```

## Benefits
- Log correlation
- Request tracing
- Debugging help
- Distributed system visibility

## Challenges
- Need to propagate everywhere
- Context management
- HTTP header propagation
- Log format consistency

## Lessons Learned
- Trace ID essential for distributed systems
- Context propagation works well
- HTTP headers for cross-service
- Logging consistency important

## Future Improvements
- OpenTelemetry integration
- Jaeger/Zipkin integration
- Trace visualization
- Performance tracing
