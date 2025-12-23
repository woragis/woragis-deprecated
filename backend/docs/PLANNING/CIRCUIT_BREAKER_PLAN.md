# Circuit Breaker Implementation Plan

**Date:** 2025-12-22  
**Status:** Partially Implemented, Needs Completion

---

## Current Status

### ✅ Already Implemented

1. **Translation Worker** - ✅ Complete
   - Google Translate API
   - DeepL API
   - LibreTranslate API

2. **Creative Service Client** - ✅ Complete
   - All image generation methods wrapped

3. **Circuit Breaker Package** - ✅ Complete
   - `pkg/circuitbreaker` package with generic `Execute[T]` function
   - Metrics integration (`RecordStateChange`, `RecordRequestAllowed`)
   - Default configurations

### ⏳ Needs Implementation

1. **Langchain Chat Client** - ⏳ Partial
   - Circuit breaker initialized but not used in `GenerateCompletion`
   - `GenerateCompletionStream` already uses circuit breaker

2. **Auth Service OAuth** - ⏳ Not Implemented
   - OAuth provider calls need circuit breaker protection

---

## Implementation Tasks

### Task 1: Fix Langchain Chat Client `GenerateCompletion`

**File:** `server/app/internal/services/langchain/chatclient.go`

**Current Issue:**
- Circuit breaker `aiServiceCB` is initialized (line ~120)
- `GenerateCompletion` method (line ~217) does NOT use circuit breaker
- Direct HTTP call without protection

**Fix Required:**
```go
// Current (line ~217-265)
func (c *Client) GenerateCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
    // Direct HTTP call - no circuit breaker
    httpClient := &http.Client{Timeout: 75 * time.Second}
    // ... HTTP call ...
}

// Should be:
func (c *Client) GenerateCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
    result, err := appcircuitbreaker.Execute(c.aiServiceCB, func() (ChatCompletionResponse, error) {
        appcircuitbreaker.RecordRequestAllowed("ai-service")
        return c.doGenerateCompletion(ctx, req)
    })
    
    if err != nil {
        if err == gobreaker.ErrOpenState {
            return ChatCompletionResponse{}, fmt.Errorf("ai-service circuit breaker is open: service unavailable")
        }
        return ChatCompletionResponse{}, err
    }
    
    return result, nil
}

func (c *Client) doGenerateCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
    // Move existing HTTP call logic here
}
```

**Priority:** High (AI service calls are critical)

---

### Task 2: Add Circuit Breaker to Auth Service OAuth

**File:** `server/app/internal/domains/auth/service_oauth.go`

**Current Issue:**
- `fetchOAuthUserInfo` method makes HTTP calls to OAuth providers
- No circuit breaker protection
- Could fail and block authentication

**Fix Required:**

1. **Add circuit breaker field to Service:**
```go
type Service struct {
    // ... existing fields ...
    oauthCB *gobreaker.CircuitBreaker
}
```

2. **Initialize in constructor:**
```go
func NewService(...) *Service {
    // Create circuit breaker for OAuth calls
    cbConfig := appcircuitbreaker.DefaultConfig("oauth-provider", logger)
    cbConfig.Timeout = 60 * time.Second // Longer timeout for OAuth
    cbConfig.OnStateChange = func(name string, from, to gobreaker.State) {
        appcircuitbreaker.RecordStateChange(name, from, to)
        if logger != nil {
            logger.Info("oauth circuit breaker state changed",
                slog.String("name", name),
                slog.String("from", from.String()),
                slog.String("to", to.String()),
            )
        }
    }
    oauthCB := appcircuitbreaker.NewCircuitBreaker(cbConfig)
    
    return &Service{
        // ... existing fields ...
        oauthCB: oauthCB,
    }
}
```

3. **Wrap OAuth calls:**
```go
func (s *Service) fetchOAuthUserInfo(ctx context.Context, provider OAuthProvider, cfg *oauthProviderConfig, token *oauth2.Token) (*oauthUserInfo, error) {
    result, err := appcircuitbreaker.Execute(s.oauthCB, func() (*oauthUserInfo, error) {
        return s.doFetchOAuthUserInfo(ctx, provider, cfg, token)
    })
    
    if err != nil {
        if err == gobreaker.ErrOpenState {
            return nil, fmt.Errorf("oauth provider circuit breaker is open: service unavailable")
        }
        return nil, err
    }
    
    return result, nil
}

func (s *Service) doFetchOAuthUserInfo(ctx context.Context, provider OAuthProvider, cfg *oauthProviderConfig, token *oauth2.Token) (*oauthUserInfo, error) {
    // Move existing HTTP call logic here
}
```

**Priority:** Medium (OAuth is important but has retry logic)

---

## Configuration Recommendations

### For AI Service (Langchain)

```go
Config{
    Name:        "ai-service",
    MaxRequests: 2,                    // Lower for AI service
    Interval:    60 * time.Second,
    Timeout:     30 * time.Second,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // Open after 3 failures (more sensitive)
        return counts.ConsecutiveFailures > 3
    },
}
```

**Reasoning:**
- AI service calls are expensive and slow
- Want to fail fast if service is down
- Lower threshold prevents wasting resources

### For OAuth Providers

```go
Config{
    Name:        "oauth-provider",
    MaxRequests: 3,
    Interval:    60 * time.Second,
    Timeout:     60 * time.Second,     // Longer timeout (OAuth can be slow)
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        // Open after 5 failures (less sensitive)
        return counts.ConsecutiveFailures > 5
    },
}
```

**Reasoning:**
- OAuth providers are generally reliable
- Longer timeout accounts for network latency
- Higher threshold prevents false positives

---

## Testing Strategy

### Unit Tests

1. **Test Circuit Breaker Opening**
   - Simulate failures
   - Verify circuit opens after threshold
   - Verify requests fail fast when open

2. **Test Circuit Breaker Recovery**
   - Open circuit
   - Wait for timeout
   - Verify transition to half-open
   - Verify successful request closes circuit

### Integration Tests

1. **Test with Real Service Down**
   - Stop AI service
   - Make requests
   - Verify circuit opens
   - Verify fast failure

2. **Test Recovery**
   - Start AI service after circuit opens
   - Wait for timeout
   - Verify circuit closes after successful request

---

## Monitoring

### Metrics to Track

1. **Circuit Breaker State Changes**
   - Open events (service down)
   - Close events (service recovered)
   - Half-open attempts

2. **Request Metrics**
   - Total requests
   - Successful requests
   - Failed requests
   - Rejected requests (circuit open)

### Alerts

1. **Circuit Breaker Open Alert**
   - Alert when circuit opens
   - Include service name
   - Include failure count

2. **Prolonged Open State**
   - Alert if circuit stays open > 5 minutes
   - Indicates persistent service issue

---

## Implementation Priority

1. **High Priority:**
   - ✅ Fix Langchain `GenerateCompletion` (AI service calls are critical)

2. **Medium Priority:**
   - ⏳ Add OAuth circuit breaker (important but less critical)

3. **Low Priority:**
   - Optional: Database circuit breakers (use connection pool settings instead)
   - Optional: Redis circuit breakers (use connection pool settings instead)

---

## Success Criteria

- [ ] All external API calls protected by circuit breakers
- [ ] All circuit breakers have appropriate configurations
- [ ] Circuit breaker state changes logged and monitored
- [ ] Metrics exported for Prometheus
- [ ] Alerts configured for circuit breaker open events
- [ ] Tests written and passing

---

**See Also:**
- `docs/deployment/circuit-breaker-implementation.md` - Implementation guide
- `pkg/circuitbreaker/circuitbreaker.go` - Circuit breaker package
- `pkg/circuitbreaker/metrics.go` - Metrics integration

---

**Last Updated:** 2025-12-22
