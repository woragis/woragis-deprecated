# Health Check Caching: 5-Second TTL Pattern

## Overview
How we implement caching for health checks to reduce load on dependencies while maintaining freshness.

## Key Points

### Problem
- Health checks called frequently (Kubernetes probes every 10s)
- Each check hits database, Redis, RabbitMQ
- Can cause load on dependencies
- Need fast responses

### Solution
- Cache health check results for 5 seconds
- Use `sync.RWMutex` for thread safety
- Check cache before hitting dependencies
- Invalidate cache after TTL

## Implementation Details

### Caching Structure
```go
type HealthCache struct {
    mu          sync.RWMutex
    result      *HealthResult
    lastCheck   time.Time
    ttl         time.Duration
}
```

### Cache Check
```go
func (c *HealthCache) Get() *HealthResult {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if time.Since(c.lastCheck) < c.ttl {
        return c.result // Return cached result
    }
    
    // Cache expired, need to check
    return nil
}
```

### Cache Update
```go
func (c *HealthCache) Update(result *HealthResult) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.result = result
    c.lastCheck = time.Now()
}
```

## Benefits
- Reduces dependency load
- Fast responses
- Still fresh (5s TTL)
- Thread-safe

## Challenges
- Cache invalidation
- Thread safety
- TTL tuning
- Stale data risk

## Lessons Learned
- Caching essential for health checks
- 5s TTL good balance
- Thread safety crucial
- Monitoring helps

## Future Improvements
- Adaptive TTL
- Cache metrics
- Per-dependency caching
- Cache warming
