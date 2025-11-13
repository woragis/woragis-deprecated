package config

import "os"

// RedisConfig holds connection details for Redis.
type RedisConfig struct {
	URL string
}

// LoadRedisConfig reads Redis configuration from environment variables.
func LoadRedisConfig() RedisConfig {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379/0"
	}

	return RedisConfig{
		URL: url,
	}
}

