package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

// NewRedis creates a new Redis client connection
func NewRedis(config RedisConfig) (*redis.Client, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("redis URL is required")
	}

	// Parse Redis URL
	opt, err := redis.ParseURL(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Override password and DB if provided separately
	if config.Password != "" {
		opt.Password = config.Password
	}
	if config.DB != 0 {
		opt.DB = config.DB
	}

	// Create Redis client
	client := redis.NewClient(opt)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("Successfully connected to Redis: %s", pong)
	return client, nil
}

// CloseRedis closes the Redis client connection
func CloseRedis(client *redis.Client) error {
	if client == nil {
		return nil
	}

	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	log.Println("Redis connection closed")
	return nil
}

// HealthCheck performs a health check on the Redis connection
func RedisHealthCheck(client *redis.Client) error {
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

// GetRedisContext returns a context with timeout for Redis operations
func GetRedisContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
