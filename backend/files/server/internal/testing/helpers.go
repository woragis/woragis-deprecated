package testhelpers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"woragis-files-service/pkg/auth"
)

// TestConfig holds test configuration
type TestConfig struct {
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
}

// LoadTestConfig loads test configuration from environment
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL: getEnv("TEST_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/auth_service_test?sslmode=disable"),
		RedisURL:    getEnv("TEST_REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:   getEnv("TEST_JWT_SECRET", "test-secret-key-for-integration-tests"),
	}
}

// SetupTestDB creates a test database connection
func SetupTestDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// SetupTestRedis creates a test Redis connection
func SetupTestRedis(url string) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// SetupTestJWTManager creates a test JWT manager
func SetupTestJWTManager(secret string, redisClient *redis.Client) *auth.JWTManager {
	jwtManager := auth.NewJWTManager(
		secret,
		"woragis-files-service-test",
		1*time.Hour,  // Access token expiry
		24*time.Hour, // Refresh token expiry
	)
	if redisClient != nil {
		jwtManager.SetRedisClient(redisClient)
	}
	return jwtManager
}

// CleanupTestDB cleans up test database
func CleanupTestDB(db *gorm.DB) error {
	// Drop all tables
	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		return fmt.Errorf("failed to clean test database: %w", err)
	}
	return nil
}

// CleanupTestRedis cleans up test Redis data
func CleanupTestRedis(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.FlushDB(ctx).Err()
}

// GenerateTestUserID generates a test user ID
func GenerateTestUserID() uuid.UUID {
	return uuid.New()
}

// GetTestContext returns a test context with timeout
func GetTestContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return ctx
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

