package timeout

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware_TimeoutExceeded(t *testing.T) {
	app := fiber.New()
	cfg := DefaultConfig()
	cfg.DefaultTimeout = 100 * time.Millisecond
	app.Use(Middleware(cfg))

	app.Get("/slow", func(c *fiber.Ctx) error {
		// Simulate slow operation
		time.Sleep(200 * time.Millisecond)
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/slow", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusRequestTimeout, resp.StatusCode)
}

func TestMiddleware_WithinTimeout(t *testing.T) {
	app := fiber.New()
	cfg := DefaultConfig()
	cfg.DefaultTimeout = 1 * time.Second
	app.Use(Middleware(cfg))

	app.Get("/fast", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/fast", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestWithTimeout_Success(t *testing.T) {
	err := WithTimeout(context.Background(), 1*time.Second, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
}

func TestWithTimeout_TimeoutExceeded(t *testing.T) {
	err := WithTimeout(context.Background(), 100*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestWithTimeout_ErrorReturned(t *testing.T) {
	expectedErr := errors.New("operation failed")
	err := WithTimeout(context.Background(), 1*time.Second, func(ctx context.Context) error {
		return expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestWithTimeoutValue_Success(t *testing.T) {
	result, err := WithTimeoutValue(context.Background(), 1*time.Second, func(ctx context.Context) (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "success", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "success", result)
}

func TestWithTimeoutValue_TimeoutExceeded(t *testing.T) {
	result, err := WithTimeoutValue(context.Background(), 100*time.Millisecond, func(ctx context.Context) (string, error) {
		time.Sleep(200 * time.Millisecond)
		return "success", nil
	})

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Empty(t, result)
}

func TestWithTimeoutValue_ErrorReturned(t *testing.T) {
	expectedErr := errors.New("operation failed")
	result, err := WithTimeoutValue(context.Background(), 1*time.Second, func(ctx context.Context) (string, error) {
		return "", expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, result)
}

