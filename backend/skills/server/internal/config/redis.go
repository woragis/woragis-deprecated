package config

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

const defaultRedisDB = 0

// LoadRedisConfig reads Redis configuration from the environment
func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		URL:      getEnvRequired("REDIS_URL"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", defaultRedisDB),
	}
}
