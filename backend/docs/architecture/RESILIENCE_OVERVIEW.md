# Resilience Overview - Backend Architecture

## General Architecture

Resilience is the ability of a system to handle failures gracefully and continue operating. Our backend implements multiple resilience patterns to ensure high availability and fault tolerance.

### Current State

- ✅ **Dead Letter Queues**: Fully implemented (RabbitMQ DLX)
- ✅ **Retry Policies**: Implemented in workers (translation-worker, etc.)
- ✅ **Graceful Degradation**: Server falls back RabbitMQ → Redis
- ❌ **Circuit Breakers**: Not implemented
- ✅ **Health Checks**: Fully implemented (all components)
- ✅ **Timeout Handling**: Implemented in various components

---

## Resilience Patterns

### 1. **Dead Letter Queues** (✅ Implemented)

**Status**: Fully implemented across all RabbitMQ queues

**Implementation:**
- All queues configured with `x-dead-letter-exchange: woragis.dlx`
- Failed messages automatically routed to DLQ
- DLQ routing key: `{queue-name}.failed`

**Configuration:**
```go
// Example from translation-worker
queueArgs := amqp.Table{
    "x-dead-letter-exchange":    "woragis.dlx",
    "x-dead-letter-routing-key": queueName + ".failed",
}
```

**Benefits:**
- Failed messages don't block queue processing
- Failed messages can be inspected and reprocessed
- Prevents message loss
- Enables failure analysis

**Monitoring:**
- Monitor DLQ size (alert if growing)
- Inspect failed messages for patterns
- Reprocess failed messages after fixing issues

**See**: Runbook for monitoring DLQ (to be created)

---

### 2. **Retry Policies** (✅ Implemented)

**Status**: Implemented in workers (translation-worker, etc.)

**Implementation:**

**Translation Worker:**
- Retry count: 3 attempts
- Retry delay: Exponential backoff (1s, 2s, 4s)
- Retry on transient errors (network, timeouts)
- Don't retry on permanent errors (invalid input, auth failures)

**Other Workers:**
- Similar retry patterns
- Configurable retry count and delay

**Example:**
```go
// Translation worker retry logic
maxRetries := 3
retryDelay := 1000 * time.Millisecond

for attempt := 0; attempt < maxRetries; attempt++ {
    result, err := translator.Translate(ctx, text, targetLang)
    if err == nil {
        return result, nil
    }
    
    if !isTransientError(err) {
        return "", err // Don't retry permanent errors
    }
    
    if attempt < maxRetries-1 {
        time.Sleep(retryDelay * time.Duration(1<<attempt)) // Exponential backoff
    }
}
```

**Benefits:**
- Handles transient failures automatically
- Reduces manual intervention
- Improves success rate for temporary issues

**Limitations:**
- No jitter (could cause thundering herd)
- Fixed retry count (could be adaptive)
- No circuit breaker integration

**Improvements Needed:**
- Add jitter to retry delays
- Implement adaptive retry (back off more on repeated failures)
- Integrate with circuit breakers

---

### 3. **Graceful Degradation** (✅ Implemented)

**Status**: Implemented in server (RabbitMQ → Redis fallback)

**Implementation:**
- Server checks RabbitMQ availability
- Falls back to Redis queue if RabbitMQ unavailable
- Health check shows "degraded" status when using fallback

**Example:**
```go
// Server initialization
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

**Benefits:**
- System continues operating during RabbitMQ outages
- No data loss
- Automatic recovery when RabbitMQ returns

**Limitations:**
- Redis queue has different characteristics (no DLQ, different guarantees)
- Degraded mode may have reduced functionality
- No automatic promotion back to RabbitMQ when it recovers

**Improvements Needed:**
- Automatic promotion back to RabbitMQ when available
- Monitor fallback usage
- Alert when in degraded mode

---

### 4. **Circuit Breakers** (❌ Not Implemented)

**Status**: Not implemented - **High Priority**

**What we need:**
- Circuit breakers for external API calls
- Circuit breakers for service-to-service calls
- Integration with retry policies

**Why it matters:**
- Prevents cascading failures
- Fails fast when downstream services are down
- Reduces load on failing services
- Enables automatic recovery

**Implementation Plan:**

#### Circuit Breaker States

1. **Closed (Normal):**
   - Requests pass through
   - Monitor failure rate
   - If failure rate exceeds threshold → Open

2. **Open (Failing):**
   - Requests fail fast (no calls made)
   - After timeout → Half-open

3. **Half-Open (Testing):**
   - Allow limited requests through
   - If successful → Closed
   - If failing → Open

#### Implementation

**Go Components:**

**Library**: `sony/gobreaker`

```go
import "github.com/sony/gobreaker"

// Create circuit breaker
cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "ai-service",
    MaxRequests: 3,              // Half-open: allow 3 requests
    Interval:    60 * time.Second, // Reset interval
    Timeout:     30 * time.Second, // Timeout before half-open
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures > 5
    },
})

// Use circuit breaker
result, err := cb.Execute(func() (interface{}, error) {
    return aiService.Call(ctx, request)
})
```

**Python Components:**

**Library**: `circuitbreaker`

```python
from circuitbreaker import circuit

@circuit(failure_threshold=5, recovery_timeout=30)
def call_ai_service(request):
    return ai_service.call(request)
```

**Integration Points:**

1. **Server → AI Service:**
   - Wrap AI service calls with circuit breaker
   - Fail fast if AI service is down
   - Return error response to client

2. **Server → Creative Service:**
   - Wrap Creative service calls with circuit breaker
   - Fail fast if Creative service is down
   - Return error response to client

3. **Translation Worker → Translation APIs:**
   - Wrap translation API calls with circuit breaker
   - Fail fast if all APIs are down
   - Route to DLQ if circuit is open

4. **Workers → Database:**
   - Wrap database calls with circuit breaker (optional)
   - Fail fast if database is down
   - Route to DLQ

**Configuration:**

```go
type CircuitBreakerConfig struct {
    Name                  string
    MaxRequests           uint32        // Half-open requests
    Interval              time.Duration // Reset interval
    Timeout               time.Duration // Timeout before half-open
    FailureThreshold      int           // Failures before opening
    SuccessThreshold      int           // Successes before closing (half-open)
}
```

**Metrics:**
- Circuit breaker state (gauge)
- Circuit breaker transitions (counter)
- Requests rejected (counter)
- Requests allowed (counter)

---

### 5. **Health Checks** (✅ Implemented)

**Status**: Fully implemented across all components

**Implementation:**
- All components expose `/healthz` endpoint
- Health checks verify dependencies
- Caching (5 seconds) to reduce load
- HTTP status codes: 200 (healthy/degraded), 503 (unhealthy)

**See**: `HEALTH_CHECKS_OVERVIEW.md` for detailed implementation

**Benefits:**
- Kubernetes can restart unhealthy pods
- Load balancers can route away from unhealthy instances
- Monitoring can alert on health check failures

---

### 6. **Timeout Handling** (✅ Implemented)

**Status**: Implemented in various components

**Implementation:**
- HTTP client timeouts
- Database query timeouts
- External API call timeouts
- Health check timeouts (2 seconds)

**Example:**
```go
// HTTP client timeout
client := &http.Client{
    Timeout: 30 * time.Second,
}

// Database query timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
db.WithContext(ctx).Find(&results)
```

**Benefits:**
- Prevents hanging requests
- Fails fast on slow dependencies
- Prevents resource exhaustion

---

## Resilience Patterns by Component

### Server

**Implemented:**
- ✅ Graceful degradation (RabbitMQ → Redis)
- ✅ Health checks (database, Redis, RabbitMQ)
- ✅ Timeout handling
- ✅ Error handling and logging

**Missing:**
- ❌ Circuit breakers for AI/Creative services
- ❌ Rate limiting (mentioned but not verified)
- ❌ Request queuing/throttling

### Workers

**Implemented:**
- ✅ Dead letter queues
- ✅ Retry policies
- ✅ Health checks
- ✅ Timeout handling
- ✅ Error handling and logging

**Missing:**
- ❌ Circuit breakers for external APIs
- ❌ Adaptive retry policies
- ❌ Backpressure handling

### Services (AI, Creative)

**Implemented:**
- ✅ Health checks
- ✅ Timeout handling
- ✅ Error handling and logging

**Missing:**
- ❌ Rate limiting
- ❌ Request queuing
- ❌ Circuit breakers for external APIs (if they call other services)

---

## Failure Scenarios and Responses

### Scenario 1: RabbitMQ Down

**Current Response:**
- Server falls back to Redis queue ✅
- Health check shows "degraded" ✅
- Workers continue processing (if already connected) ✅

**Improvements:**
- Alert when in degraded mode
- Automatic promotion back to RabbitMQ when available
- Monitor Redis queue depth

### Scenario 2: Database Down

**Current Response:**
- Health check shows "unhealthy" ✅
- Requests fail with error ✅
- Workers can't process jobs ✅

**Improvements:**
- Circuit breaker for database calls
- Queue jobs for later processing
- Read-only mode (if applicable)

### Scenario 3: External API Down (AI Service, Translation APIs)

**Current Response:**
- Requests fail ✅
- Workers retry (if transient error) ✅
- Jobs go to DLQ after max retries ✅

**Improvements:**
- Circuit breaker to fail fast
- Fallback to alternative APIs (if available)
- Graceful degradation (return cached/default responses)

### Scenario 4: High Load

**Current Response:**
- Requests queue in RabbitMQ ✅
- Workers process at their rate ✅
- Health checks may slow down ✅

**Improvements:**
- Rate limiting
- Request throttling
- Auto-scaling based on queue depth
- Backpressure handling

---

## Best Practices

### 1. **Fail Fast**

- Use timeouts on all external calls
- Use circuit breakers to fail fast when services are down
- Don't wait indefinitely for responses

### 2. **Retry Wisely**

- Only retry transient errors
- Use exponential backoff with jitter
- Limit retry count
- Integrate with circuit breakers

### 3. **Degrade Gracefully**

- Provide fallback functionality when possible
- Return cached data when services are down
- Reduce functionality rather than failing completely

### 4. **Monitor and Alert**

- Monitor failure rates
- Monitor circuit breaker states
- Monitor DLQ sizes
- Alert on degradation

### 5. **Test Failure Scenarios**

- Test service failures
- Test network partitions
- Test high load scenarios
- Test recovery procedures

---

## Implementation Roadmap

### Phase 1: Circuit Breakers (High Priority) - 1-2 weeks

**Week 1:**
- Implement circuit breakers for AI Service calls
- Implement circuit breakers for Creative Service calls
- Implement circuit breakers for Translation API calls

**Week 2:**
- Add circuit breaker metrics
- Integrate with retry policies
- Test failure scenarios

### Phase 2: Enhanced Retry (Medium Priority) - 1 week

- Add jitter to retry delays
- Implement adaptive retry
- Integrate with circuit breakers

### Phase 3: Rate Limiting (Medium Priority) - 1 week

- Implement rate limiting in server
- Implement rate limiting in services
- Configure rate limits per endpoint/user

### Phase 4: Backpressure (Low Priority) - 1-2 weeks

- Implement backpressure handling
- Queue management during high load
- Auto-scaling based on queue depth

---

## Metrics for Resilience

### Key Metrics

1. **Failure Rate:**
   - `requests_failed_total` / `requests_total`
   - Target: < 1%

2. **Circuit Breaker State:**
   - `circuit_breaker_state` (gauge)
   - Monitor transitions

3. **Retry Rate:**
   - `requests_retried_total` / `requests_total`
   - Monitor retry patterns

4. **DLQ Size:**
   - `queue_dlq_size` (gauge)
   - Alert if growing

5. **Degradation Events:**
   - `degradation_events_total` (counter)
   - Monitor fallback usage

### Alerts

**Critical:**
- Circuit breaker open for > 5 minutes
- DLQ size > 1000
- Failure rate > 5%

**Warning:**
- Circuit breaker open
- DLQ size > 100
- Failure rate > 1%
- Degradation mode active

---

## Summary

**Current State:**
- ✅ Dead letter queues (complete)
- ✅ Retry policies (implemented)
- ✅ Graceful degradation (implemented)
- ✅ Health checks (complete)
- ✅ Timeout handling (implemented)
- ❌ Circuit breakers (not implemented)

**Priority:**
1. **Circuit Breakers** - High priority (prevents cascading failures)
2. **Enhanced Retry** - Medium priority (improves reliability)
3. **Rate Limiting** - Medium priority (prevents overload)
4. **Backpressure** - Low priority (handles high load)

**Impact:**
- Circuit breakers: Prevents cascading failures, fails fast
- Enhanced retry: Improves success rate for transient failures
- Rate limiting: Prevents overload, ensures fair usage
- Backpressure: Handles high load gracefully
