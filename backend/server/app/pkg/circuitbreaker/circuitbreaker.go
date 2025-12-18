package circuitbreaker

import (
	"log/slog"
	"time"

	"github.com/sony/gobreaker"
)

// Config holds circuit breaker configuration
type Config struct {
	Name                  string
	MaxRequests           uint32        // Half-open: max requests allowed
	Interval              time.Duration // Reset interval for closed state
	Timeout               time.Duration // Timeout before transitioning from open to half-open
	ReadyToTrip           func(counts gobreaker.Counts) bool
	OnStateChange         func(name string, from gobreaker.State, to gobreaker.State)
	Logger                *slog.Logger
}

// DefaultConfig returns a default circuit breaker configuration
func DefaultConfig(name string, logger *slog.Logger) Config {
	return Config{
		Name:        name,
		MaxRequests: 3,                    // Allow 3 requests in half-open state
		Interval:    60 * time.Second,    // Reset interval
		Timeout:     30 * time.Second,    // Timeout before half-open
		ReadyToTrip: defaultReadyToTrip,  // Open after 5 consecutive failures
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if logger != nil {
				logger.Info("circuit breaker state changed",
					slog.String("name", name),
					slog.String("from", from.String()),
					slog.String("to", to.String()),
				)
			}
		},
		Logger: logger,
	}
}

// defaultReadyToTrip opens circuit after 5 consecutive failures
func defaultReadyToTrip(counts gobreaker.Counts) bool {
	return counts.ConsecutiveFailures > 5
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(cfg Config) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: cfg.ReadyToTrip,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if cfg.OnStateChange != nil {
				cfg.OnStateChange(name, from, to)
			}
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}

// Execute wraps a function call with circuit breaker protection
func Execute[T any](cb *gobreaker.CircuitBreaker, fn func() (T, error)) (T, error) {
	var zero T
	result, err := cb.Execute(func() (interface{}, error) {
		return fn()
	})
	if err != nil {
		return zero, err
	}
	return result.(T), nil
}
