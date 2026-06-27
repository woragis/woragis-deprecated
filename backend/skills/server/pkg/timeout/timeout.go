package timeout

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	appmetrics "woragis-skills-service/pkg/metrics"
)

// Config holds timeout configuration
type Config struct {
	// DefaultTimeout is the default timeout for requests
	DefaultTimeout time.Duration
	// HandlerTimeoutMsg is the message returned when timeout is exceeded
	HandlerTimeoutMsg string
	// OnTimeout is called when a timeout occurs
	OnTimeout func(*fiber.Ctx)
}

// DefaultConfig returns default timeout configuration
func DefaultConfig() Config {
	return Config{
		DefaultTimeout:    30 * time.Second,
		HandlerTimeoutMsg: "Request timeout",
		OnTimeout:         nil,
	}
}

// Middleware creates a timeout middleware that sets a context timeout for requests
func Middleware(cfg Config) fiber.Handler {
	if cfg.DefaultTimeout <= 0 {
		cfg.DefaultTimeout = 30 * time.Second
	}
	if cfg.HandlerTimeoutMsg == "" {
		cfg.HandlerTimeoutMsg = "Request timeout"
	}

	return func(c *fiber.Ctx) error {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.UserContext(), cfg.DefaultTimeout)
		defer cancel()

		// Set the context in the request
		c.SetUserContext(ctx)

		// Continue with the request
		err := c.Next()

		// Check if context was cancelled due to timeout
		if ctx.Err() == context.DeadlineExceeded {
			// Record timeout metric
			endpoint := c.Route().Path
			if endpoint == "" {
				endpoint = c.Path()
			}
			appmetrics.RecordRequestTimeout(endpoint)

			if cfg.OnTimeout != nil {
				cfg.OnTimeout(c)
			}
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
				"error":   cfg.HandlerTimeoutMsg,
				"message": "The request took longer than the allowed time",
			})
		}

		return err
	}
}

// WithTimeout executes a function with a timeout context
func WithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(ctx)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WithTimeoutValue executes a function with a timeout and returns a value
func WithTimeoutValue[T any](ctx context.Context, timeout time.Duration, fn func(context.Context) (T, error)) (T, error) {
	var zero T
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	type result struct {
		value T
		err   error
	}

	done := make(chan result, 1)
	go func() {
		val, err := fn(ctx)
		done <- result{value: val, err: err}
	}()

	select {
	case res := <-done:
		return res.value, res.err
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}

