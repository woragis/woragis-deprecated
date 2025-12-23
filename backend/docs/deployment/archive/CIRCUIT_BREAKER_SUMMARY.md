# Circuit Breaker Implementation Summary

**Date:** 2025-12-22  
**Status:** ✅ Implemented Across System

---

## Overview

Circuit breakers are implemented using the **Sony gobreaker library** with a custom wrapper package (`pkg/circuitbreaker`) that provides:
- Type-safe generic `Execute[T]` function
- Prometheus metrics integration
- Configurable thresholds and timeouts
- State change logging

---

## Implementation Status

### ✅ Fully Implemented

1. **Translation Worker**
   - ✅ Google Translate API
   - ✅ DeepL API
   - ✅ LibreTranslate API
   - Location: `translation-worker/internal/translator/translator.go`

2. **Creative Service Client**
   - ✅ Image generation
   - ✅ Thumbnail generation
   - ✅ All creative service methods
   - Location: `server/app/internal/services/creative/client.go`

3. **Langchain Chat Client**
   - ✅ AI Service calls (when `AI_SERVICE_URL` is set)
   - ✅ Streaming completions
   - ⚠️ Direct OpenAI calls (bypass circuit breaker when `AI_SERVICE_URL` not set)
   - Location: `server/app/internal/services/langchain/chatclient.go`

### ⏳ Could Be Added (Optional)

1. **Auth Service OAuth**
   - External OAuth provider calls
   - Location: `server/app/internal/domains/auth/service_oauth.go`
   - Priority: Medium (OAuth providers are generally reliable)

---

## How It Works

### 1. Circuit Breaker Package

**Location:** `server/app/pkg/circuitbreaker/circuitbreaker.go`

Provides:
- `Config` struct for configuration
- `DefaultConfig()` for sensible defaults
- `NewCircuitBreaker()` to create circuit breaker
- `Execute[T]()` generic function to wrap calls

**Default Configuration:**
```go
{
    MaxRequests: 3,                    // Allow 3 requests in half-open state
    Interval:    60 * time.Second,    // Reset interval
    Timeout:     30 * time.Second,    // Timeout before half-open
    ReadyToTrip: opens after 5 consecutive failures
}
```

### 2. Metrics Integration

**Location:** `server/app/pkg/circuitbreaker/metrics.go`

Exports Prometheus metrics:
- `circuit_breaker_state` - Current state (0=closed, 1=half-open, 2=open)
- `circuit_breaker_transitions_total` - State transition counter
- `circuit_breaker_requests_rejected_total` - Rejected requests (circuit open)
- `circuit_breaker_requests_allowed_total` - Allowed requests

### 3. Usage Pattern

**Step 1: Initialize Circuit Breaker**
```go
cbConfig := appcircuitbreaker.DefaultConfig("service-name", logger)
cbConfig.OnStateChange = func(name string, from, to gobreaker.State) {
    appcircuitbreaker.RecordStateChange(name, from, to)
    logger.Info("circuit breaker state changed", ...)
}
cb := appcircuitbreaker.NewCircuitBreaker(cbConfig)
```

**Step 2: Wrap API Calls**
```go
result, err := appcircuitbreaker.Execute(cb, func() (Response, error) {
    appcircuitbreaker.RecordRequestAllowed("service-name")
    return actualAPICall(ctx, params)
})

if err != nil {
    if err == gobreaker.ErrOpenState {
        appcircuitbreaker.RecordRequestRejected("service-name")
        return Response{}, fmt.Errorf("service unavailable: circuit breaker open")
    }
    return Response{}, err
}
```

---

## Circuit Breaker States

### Closed (Normal)
- Requests pass through normally
- Failures are counted
- After 5 consecutive failures → transitions to Open

### Open (Service Down)
- Requests fail immediately with `ErrOpenState`
- No actual API calls made
- After 30 seconds → transitions to Half-Open

### Half-Open (Testing Recovery)
- Allows up to 3 requests to test if service recovered
- If all succeed → transitions to Closed
- If any fail → transitions back to Open

---

## Examples from Codebase

### Example 1: DeepL Translator

```go
// Initialize
cbConfig := appcircuitbreaker.DefaultConfig("deepl-translate", logger)
cb := appcircuitbreaker.NewCircuitBreaker(cbConfig)

// Use
func (t *DeepLTranslator) Translate(ctx context.Context, text, lang string) (string, error) {
    result, err := appcircuitbreaker.Execute(t.cb, func() (string, error) {
        return t.doTranslate(ctx, text, lang)
    })
    
    if err == gobreaker.ErrOpenState {
        return "", fmt.Errorf("deepl-translate circuit breaker is open: service unavailable")
    }
    return result, err
}
```

### Example 2: Creative Service

```go
// Initialize in constructor
cbConfig := appcircuitbreaker.DefaultConfig("creative-service", logger)
cbConfig.OnStateChange = func(name string, from, to gobreaker.State) {
    appcircuitbreaker.RecordStateChange(name, from, to)
    logger.Info("circuit breaker state changed", ...)
}
cb := appcircuitbreaker.NewCircuitBreaker(cbConfig)

// Use in method
func (c *Client) GenerateImage(ctx context.Context, req ImageGenerationRequest) (*ImageGenerationResponse, error) {
    result, err := appcircuitbreaker.Execute(c.cb, func() (*ImageGenerationResponse, error) {
        appcircuitbreaker.RecordRequestAllowed("creative-service")
        return c.doGenerateImage(ctx, req)
    })
    
    if err == gobreaker.ErrOpenState {
        appcircuitbreaker.RecordRequestRejected("creative-service")
        return nil, fmt.Errorf("creative-service circuit breaker is open: service unavailable")
    }
    return result, err
}
```

### Example 3: Langchain AI Service

```go
// Initialize in constructor
aiServiceCB := appcircuitbreaker.NewCircuitBreaker(
    appcircuitbreaker.DefaultConfig("ai-service", logger),
)

// Use
func (c *Client) callAIService(ctx context.Context, baseURL string, req ChatCompletionRequest) (ChatCompletionResponse, error) {
    result, err := appcircuitbreaker.Execute(c.aiServiceCB, func() (ChatCompletionResponse, error) {
        appcircuitbreaker.RecordRequestAllowed("ai-service")
        return c.doAIServiceCall(ctx, baseURL, req)
    })
    
    if err == gobreaker.ErrOpenState {
        appcircuitbreaker.RecordRequestRejected("ai-service")
        return ChatCompletionResponse{}, fmt.Errorf("ai-service circuit breaker is open: service unavailable")
    }
    return result, err
}
```

---

## Configuration Guidelines

### External APIs (Google, DeepL, etc.)
```go
Config{
    MaxRequests: 3,
    Interval:    60 * time.Second,
    Timeout:     30 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures > 5
    },
}
```

### Internal Services (AI Service, Creative Service)
```go
Config{
    MaxRequests: 2,                    // Lower for internal services
    Interval:    60 * time.Second,
    Timeout:     30 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures > 3  // More sensitive
    },
}
```

### Critical Operations
```go
Config{
    MaxRequests: 1,                    // Very conservative
    Interval:    30 * time.Second,
    Timeout:     60 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures > 2  // Very sensitive
    },
}
```

---

## Monitoring

### Metrics Available

1. **Circuit Breaker State**
   ```
   circuit_breaker_state{name="ai-service"} 0  # 0=closed, 1=half-open, 2=open
   ```

2. **State Transitions**
   ```
   circuit_breaker_transitions_total{name="ai-service",from="closed",to="open"} 5
   ```

3. **Rejected Requests**
   ```
   circuit_breaker_requests_rejected_total{name="ai-service"} 42
   ```

4. **Allowed Requests**
   ```
   circuit_breaker_requests_allowed_total{name="ai-service"} 1234
   ```

### Grafana Queries

**Circuit Breaker Open Duration:**
```promql
time() - circuit_breaker_state{state="2"} * (time() - circuit_breaker_transitions_total{to="open"})
```

**Rejection Rate:**
```promql
rate(circuit_breaker_requests_rejected_total[5m]) / 
  (rate(circuit_breaker_requests_allowed_total[5m]) + 
   rate(circuit_breaker_requests_rejected_total[5m]))
```

---

## Best Practices

1. ✅ **Use for External Dependencies**
   - External APIs
   - Internal microservices
   - Third-party services

2. ❌ **Don't Use for Internal Operations**
   - Simple function calls
   - In-memory operations
   - Pure computations

3. ✅ **Configure Appropriately**
   - More aggressive for critical services
   - Longer timeouts for slow services
   - Monitor and adjust based on metrics

4. ✅ **Provide Fallbacks**
   - Cached data when circuit is open
   - Default responses
   - Graceful error messages

5. ✅ **Log and Monitor**
   - Log all state changes
   - Track metrics
   - Set up alerts

---

## When to Add Circuit Breakers

### ✅ Should Add

- External API calls (Google, DeepL, etc.)
- Internal microservice calls (AI Service, Creative Service)
- Third-party service calls
- OAuth provider calls (optional but recommended)

### ❌ Don't Need

- Database queries (use connection pool settings)
- Redis operations (use connection pool settings)
- Simple function calls
- In-memory operations

---

## Additional Resources

- **Implementation Guide:** `docs/deployment/circuit-breaker-implementation.md`
- **Implementation Plan:** `docs/PLANNING/CIRCUIT_BREAKER_PLAN.md`
- **Package Code:** `pkg/circuitbreaker/circuitbreaker.go`
- **Metrics Code:** `pkg/circuitbreaker/metrics.go`
- **Example:** `server/app/internal/services/creative/client.go`

---

**Circuit breakers are fully implemented and working across the system!** ✅

---

**Last Updated:** 2025-12-22
