# Circuit Breaker Implementation Guide

**Purpose:** Guide for implementing circuit breakers across the system to prevent cascading failures

---

## Overview

Circuit breakers are a resilience pattern that prevents cascading failures by detecting when a service is down and failing fast instead of allowing requests to pile up and timeout.

**States:**
- **Closed**: Normal operation, requests pass through
- **Open**: Service failing, requests fail immediately
- **Half-Open**: Testing if service recovered, limited requests allowed

---

## Current Status

### ✅ Already Implemented

1. **Translation Worker**
   - ✅ Google Translate API
   - ✅ DeepL API
   - ✅ LibreTranslate API

2. **Creative Service Client**
   - ✅ Image generation calls

3. **Circuit Breaker Package**
   - ✅ `pkg/circuitbreaker` in both `server/app` and `translation-worker`
   - ✅ Uses `gobreaker` (Sony library)
   - ✅ Generic `Execute[T]` function for type-safe wrapping

### ⏳ Needs Implementation

1. **Langchain Chat Client**
   - ⏳ `GenerateCompletion` method (has circuit breaker field but not used)

2. **Auth Service**
   - ⏳ OAuth provider calls (Google, etc.)

3. **Database Connections**
   - ⏳ Critical database operations (optional, but recommended)

---

## Implementation Pattern

### Step 1: Add Circuit Breaker Field

```go
type Client struct {
    baseURL    string
    httpClient *http.Client
    cb         *gobreaker.CircuitBreaker  // Add this field
    logger     *slog.Logger
}
```

### Step 2: Initialize Circuit Breaker in Constructor

```go
func NewClient(baseURL string, logger *slog.Logger) *Client {
    // Create circuit breaker
    cbConfig := appcircuitbreaker.DefaultConfig("service-name", logger)
    cbConfig.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
        if logger != nil {
            logger.Info("circuit breaker state changed",
                slog.String("name", name),
                slog.String("from", from.String()),
                slog.String("to", to.String()),
            )
        }
    }
    cb := appcircuitbreaker.NewCircuitBreaker(cbConfig)
    
    return &Client{
        baseURL: baseURL,
        httpClient: &http.Client{Timeout: 30 * time.Second},
        cb: cb,
        logger: logger,
    }
}
```

### Step 3: Wrap API Calls with Circuit Breaker

```go
func (c *Client) CallAPI(ctx context.Context, req Request) (Response, error) {
    // Wrap the call with circuit breaker
    result, err := appcircuitbreaker.Execute(c.cb, func() (Response, error) {
        return c.doCallAPI(ctx, req)
    })
    
    if err != nil {
        // Check if circuit is open
        if err == gobreaker.ErrOpenState {
            return Response{}, fmt.Errorf("service circuit breaker is open: service unavailable")
        }
        return Response{}, err
    }
    
    return result, nil
}

func (c *Client) doCallAPI(ctx context.Context, req Request) (Response, error) {
    // Actual HTTP call implementation
    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/endpoint", body)
    if err != nil {
        return Response{}, err
    }
    
    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return Response{}, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Process response...
    return response, nil
}
```

---

## Configuration Options

### Default Configuration

```go
Config{
    Name:        "service-name",
    MaxRequests: 3,                    // Allow 3 requests in half-open state
    Interval:    60 * time.Second,    // Reset interval for closed state
    Timeout:     30 * time.Second,    // Timeout before transitioning from open to half-open
    ReadyToTrip: defaultReadyToTrip,  // Opens after 5 consecutive failures
}
```

### Custom Configuration

```go
cbConfig := appcircuitbreaker.Config{
    Name:        "critical-service",
    MaxRequests: 1,                     // Only 1 request in half-open
    Interval:    30 * time.Second,     // Shorter reset interval
    Timeout:     60 * time.Second,     // Longer timeout before half-open
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // Open circuit after 3 failures OR if 50% of last 100 requests failed
        return counts.ConsecutiveFailures > 3 || 
               (counts.TotalRequests >= 100 && 
                float64(counts.TotalFailures)/float64(counts.TotalRequests) > 0.5)
    },
    Logger: logger,
}
```

---

## Implementation Checklist

### Priority 1: External API Calls

- [ ] **Langchain Chat Client** (`server/app/internal/services/langchain/chatclient.go`)
  - [ ] Add circuit breaker to `GenerateCompletion` method
  - [ ] Wrap HTTP call in `Execute` function
  - [ ] Handle `ErrOpenState` appropriately

- [ ] **Auth Service OAuth** (`server/app/internal/domains/auth/service_oauth.go`)
  - [ ] Add circuit breaker for OAuth provider calls
  - [ ] Wrap user info endpoint calls
  - [ ] Handle failures gracefully

### Priority 2: Internal Service Calls

- [ ] **AI Service Client** (if exists)
  - [ ] Wrap all AI service calls
  - [ ] Configure appropriate timeouts

- [ ] **Other HTTP Clients**
  - [ ] Identify all HTTP clients in codebase
  - [ ] Add circuit breakers where appropriate

### Priority 3: Database Operations (Optional)

- [ ] **Critical Queries**
  - [ ] Wrap critical database operations
  - [ ] Use circuit breaker for connection pool exhaustion scenarios
  - [ ] Consider using connection pool settings instead

---

## Error Handling

### Circuit Open Error

When circuit breaker is open, requests fail immediately with `gobreaker.ErrOpenState`:

```go
result, err := appcircuitbreaker.Execute(cb, fn)
if err != nil {
    if err == gobreaker.ErrOpenState {
        // Circuit is open - service is down
        return Response{}, fmt.Errorf("service unavailable: circuit breaker open")
    }
    // Other errors (network, timeout, etc.)
    return Response{}, err
}
```

### Graceful Degradation

When circuit is open, provide fallback behavior:

```go
result, err := c.CallAPI(ctx, req)
if err != nil {
    if err == gobreaker.ErrOpenState {
        // Fallback to cached data or default response
        return c.getCachedResponse(req), nil
    }
    return Response{}, err
}
```

---

## Monitoring

### Logging State Changes

Circuit breaker logs state changes automatically through `OnStateChange` callback:

```go
cbConfig.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
    logger.Info("circuit breaker state changed",
        slog.String("name", name),
        slog.String("from", from.String()),
        slog.String("to", to.String()),
    )
}
```

### Metrics Integration

Add metrics to track circuit breaker state:

```go
// In OnStateChange callback
if to == gobreaker.StateOpen {
    metrics.RecordCircuitBreakerOpen(name)
} else if from == gobreaker.StateOpen && to == gobreaker.StateClosed {
    metrics.RecordCircuitBreakerClosed(name)
}
```

---

## Testing

### Unit Tests

```go
func TestCircuitBreaker(t *testing.T) {
    cb := appcircuitbreaker.NewCircuitBreaker(
        appcircuitbreaker.DefaultConfig("test", nil),
    )
    
    // Simulate failures to open circuit
    for i := 0; i < 6; i++ {
        _, err := appcircuitbreaker.Execute(cb, func() (string, error) {
            return "", fmt.Errorf("simulated failure")
        })
        assert.Error(t, err)
    }
    
    // Circuit should now be open
    _, err := appcircuitbreaker.Execute(cb, func() (string, error) {
        return "success", nil
    })
    assert.Equal(t, gobreaker.ErrOpenState, err)
}
```

### Integration Tests

Test circuit breaker behavior with actual service:

1. Start service
2. Make requests (should succeed)
3. Stop service
4. Make requests (should fail and open circuit)
5. Circuit should be open (requests fail fast)
6. Wait for timeout
7. Start service
8. Circuit should transition to half-open, then closed

---

## Configuration Recommendations

### For External APIs (Google, DeepL, etc.)

```go
Config{
    Name:        "external-api-name",
    MaxRequests: 3,
    Interval:    60 * time.Second,
    Timeout:     30 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures > 5
    },
}
```

### For Internal Services (AI Service, Creative Service)

```go
Config{
    Name:        "internal-service-name",
    MaxRequests: 2,
    Interval:    30 * time.Second,
    Timeout:     15 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures > 3
    },
}
```

### For Critical Operations

```go
Config{
    Name:        "critical-operation",
    MaxRequests: 1,
    Interval:    15 * time.Second,
    Timeout:     60 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // More aggressive: open after 2 failures
        return counts.ConsecutiveFailures > 2
    },
}
```

---

## Best Practices

1. **Use Circuit Breakers for External Dependencies**
   - External APIs (Google, DeepL, etc.)
   - Internal microservices
   - Third-party services

2. **Don't Use for Internal Operations**
   - Simple function calls
   - In-memory operations
   - Pure computations

3. **Configure Appropriately**
   - More aggressive for critical services
   - Longer timeouts for slow services
   - Monitor and adjust based on metrics

4. **Provide Fallbacks**
   - Cached data when circuit is open
   - Default responses
   - Graceful error messages

5. **Log and Monitor**
   - Log all state changes
   - Track metrics (open/closed duration, failure rates)
   - Set up alerts for circuit breaker open events

---

## Examples

See existing implementations:
- `translation-worker/internal/translator/translator.go` - DeepL/Google Translate
- `server/app/internal/services/creative/client.go` - Creative Service

---

**See Also:**
- `pkg/circuitbreaker/circuitbreaker.go` - Circuit breaker package
- `docs/deployment/performance-optimization.md` - Performance tuning
- `docs/operations/monitoring-alerting.md` - Monitoring setup

---

**Last Updated:** 2025-12-22
