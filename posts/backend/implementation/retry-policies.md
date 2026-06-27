# Retry Policies: Handling External API Failures

## Overview
How we implement retry policies for external API calls (translation APIs, AI services) with exponential backoff.

## Key Points

### Problem
- External APIs can fail (network, rate limits, temporary errors)
- Need to retry transient failures
- Don't retry permanent failures
- Need to avoid overwhelming APIs

### Solution
- Retry with exponential backoff
- Max retry attempts (3-5)
- Retry delay (1s, 2s, 4s, ...)
- Error classification (transient vs permanent)

## Implementation Details

### Retry Configuration
```go
type RetryConfig struct {
    MaxAttempts int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
}
```

### Retry Logic
```go
func Retry(ctx context.Context, fn func() error, config RetryConfig) error {
    for attempt := 0; attempt < config.MaxAttempts; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if !isTransient(err) {
            return err // Don't retry permanent errors
        }
        
        delay := calculateBackoff(attempt, config)
        time.Sleep(delay)
    }
    return ErrMaxRetriesExceeded
}
```

### Error Classification
- **Transient**: Network errors, 5xx errors, rate limits
- **Permanent**: 4xx errors (except 429), validation errors

## Benefits
- Handles transient failures
- Reduces API load (backoff)
- Improves reliability
- Better user experience

## Challenges
- Need to classify errors correctly
- Backoff calculation
- Context cancellation
- Logging retry attempts

## Lessons Learned
- Retry policies essential
- Error classification important
- Exponential backoff works well
- Logging helps debug

## Future Improvements
- Circuit breaker integration
- Adaptive retry delays
- Retry metrics
- Retry dashboard
