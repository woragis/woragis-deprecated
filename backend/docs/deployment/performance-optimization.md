# Performance Optimization Guide

**Last Updated:** 2025-12-22  
**Purpose:** Guide for optimizing system performance

---

## Overview

This guide covers performance optimization strategies for the Woragis backend services.

---

## Performance Baseline

### Current Metrics (Target)

**API Response Times:**
- P50 (median): < 100ms
- P95: < 500ms
- P99: < 1000ms

**Throughput:**
- Requests per second: 1000+
- Concurrent users: 500+

**Resource Usage:**
- CPU: < 70% average
- Memory: < 80% average
- Database connections: < 80% of pool

---

## Optimization Areas

### 1. Database Optimization

#### Query Optimization

**Add Indexes:**
```sql
-- Common indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_jobs_status ON jobs(status) WHERE status = 'pending';

-- Composite indexes
CREATE INDEX idx_projects_user_status ON projects(user_id, status);
```

**Analyze Slow Queries:**
```sql
-- Enable query logging
SET log_min_duration_statement = 1000;  -- Log queries > 1s

-- View slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

**Connection Pooling:**
```go
// Go (database/sql)
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

```python
# Python (SQLAlchemy)
engine = create_engine(
    DATABASE_URL,
    pool_size=20,
    max_overflow=10,
    pool_timeout=30,
    pool_recycle=3600
)
```

#### Query Optimization Tips

1. **Use SELECT specific columns** instead of `SELECT *`
2. **Use LIMIT** for pagination
3. **Avoid N+1 queries** - use JOINs or batch loading
4. **Use EXPLAIN ANALYZE** to understand query plans
5. **Add indexes** for frequently queried columns

---

### 2. Caching Strategy

#### Application-Level Caching

**Redis Caching:**
```go
// Go example
func GetUser(ctx context.Context, userID string) (*User, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("user:%s", userID)
    cached, err := redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // Fetch from database
    user, err := db.GetUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Cache for 1 hour
    userJSON, _ := json.Marshal(user)
    redis.Set(ctx, cacheKey, userJSON, time.Hour)
    
    return user, nil
}
```

**Cache Invalidation:**
```go
func UpdateUser(ctx context.Context, userID string, updates User) error {
    // Update database
    err := db.UpdateUser(ctx, userID, updates)
    if err != nil {
        return err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%s", userID)
    redis.Del(ctx, cacheKey)
    
    return nil
}
```

#### Cache Patterns

1. **Cache-Aside (Lazy Loading)**
   - Check cache, if miss, load from DB and cache
   - Most common pattern

2. **Write-Through**
   - Write to cache and DB simultaneously
   - Ensures cache consistency

3. **Write-Back**
   - Write to cache first, flush to DB later
   - Higher performance, risk of data loss

---

### 3. API Optimization

#### Response Compression

**Go (Fiber):**
```go
import "github.com/gofiber/fiber/v2/middleware/compress"

app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed,
}))
```

**Python (FastAPI):**
```python
from fastapi.middleware.gzip import GZipMiddleware

app.add_middleware(GZipMiddleware, minimum_size=1000)
```

#### Pagination

```go
type PaginationParams struct {
    Page     int `query:"page" validate:"min=1"`
    PageSize int `query:"page_size" validate:"min=1,max=100"`
}

func GetUsers(c *fiber.Ctx) error {
    params := PaginationParams{
        Page:     1,
        PageSize: 20,
    }
    c.QueryParser(&params)
    
    offset := (params.Page - 1) * params.PageSize
    users, err := db.GetUsers(offset, params.PageSize)
    // ...
}
```

#### Response Streaming

**For Large Responses:**
```go
func StreamLargeData(c *fiber.Ctx) error {
    c.Set("Content-Type", "application/json")
    c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
        // Stream data in chunks
        for item := range dataChannel {
            json.NewEncoder(w).Encode(item)
            w.Flush()
        }
    })
    return nil
}
```

---

### 4. Worker Optimization

#### Batch Processing

```go
func ProcessBatch(jobs []Job) error {
    batchSize := 100
    
    for i := 0; i < len(jobs); i += batchSize {
        end := i + batchSize
        if end > len(jobs) {
            end = len(jobs)
        }
        
        batch := jobs[i:end]
        if err := processBatch(batch); err != nil {
            return err
        }
    }
    
    return nil
}
```

#### Worker Scaling

**Horizontal Scaling:**
- Run multiple worker instances
- Use queue prefetch to distribute work
- Monitor queue depth

**Vertical Scaling:**
- Increase worker resources (CPU, memory)
- Adjust worker concurrency
- Optimize worker code

---

### 5. Code Performance

#### Profiling

**Go Profiling:**
```go
import _ "net/http/pprof"

// In main.go
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// Profile with:
// go tool pprof http://localhost:6060/debug/pprof/profile
```

**Python Profiling:**
```python
import cProfile
import pstats

profiler = cProfile.Profile()
profiler.enable()

# Your code here

profiler.disable()
stats = pstats.Stats(profiler)
stats.sort_stats('cumulative')
stats.print_stats(10)  # Top 10 functions
```

#### Optimization Tips

1. **Avoid premature optimization** - Profile first
2. **Optimize hot paths** - Focus on frequently executed code
3. **Use appropriate data structures** - Maps vs slices
4. **Minimize allocations** - Reuse buffers, pools
5. **Use goroutines/pooling** - For concurrent operations

---

## Load Testing

### Tools

**k6 (Recommended):**
```javascript
// load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '2m', target: 100 },  // Ramp up to 100 users
        { duration: '5m', target: 100 },  // Stay at 100 users
        { duration: '2m', target: 200 },  // Ramp up to 200 users
        { duration: '5m', target: 200 },  // Stay at 200 users
        { duration: '2m', target: 0 },    // Ramp down
    ],
};

export default function() {
    const response = http.get('http://localhost:8080/api/users');
    check(response, {
        'status is 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
    });
    sleep(1);
}
```

**Run:**
```bash
k6 run load-test.js
```

### Load Testing Scenarios

1. **Baseline Test** - Normal load
2. **Stress Test** - Maximum capacity
3. **Spike Test** - Sudden load increase
4. **Endurance Test** - Sustained load
5. **Volume Test** - Large data volumes

---

## Performance Monitoring

### Key Metrics

**Track:**
- Request rate (req/s)
- Response time (p50, p95, p99)
- Error rate (%)
- Throughput (bytes/s)
- Resource usage (CPU, memory, disk, network)

**Grafana Queries:**
```promql
# Request rate
rate(http_requests_total[5m])

# Response time (p95)
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])
```

---

## Performance Checklist

### Database
- [ ] Add indexes for frequently queried columns
- [ ] Optimize slow queries
- [ ] Configure connection pooling
- [ ] Enable query logging
- [ ] Monitor query performance

### Caching
- [ ] Implement Redis caching
- [ ] Cache frequently accessed data
- [ ] Set appropriate TTLs
- [ ] Implement cache invalidation
- [ ] Monitor cache hit rates

### API
- [ ] Enable response compression
- [ ] Implement pagination
- [ ] Optimize response sizes
- [ ] Use streaming for large responses
- [ ] Add rate limiting

### Workers
- [ ] Optimize batch processing
- [ ] Scale workers horizontally
- [ ] Monitor queue processing times
- [ ] Optimize worker code
- [ ] Balance worker load

### Code
- [ ] Profile application
- [ ] Optimize hot paths
- [ ] Minimize allocations
- [ ] Use appropriate data structures
- [ ] Optimize serialization

---

## Related Documentation

- [Monitoring & Alerting](./monitoring-alerting.md)
- [Database Configuration](../deployment/configuration.md)
- [Caching Strategy](../architecture/system-overview.md)

---

**Next Steps:**
1. Establish performance baselines
2. Profile application
3. Identify bottlenecks
4. Implement optimizations
5. Run load tests
6. Monitor performance improvements

---

**Last Updated:** 2025-12-22
