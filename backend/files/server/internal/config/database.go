package config

import (
	"time"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxIdleTime     time.Duration
	ConnMaxLifetime time.Duration
}

const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 25
	defaultMaxIdleTime     = 15 * time.Minute
	defaultConnMaxLifetime = 60 * time.Minute
)

// LoadDatabaseConfig reads database configuration from the environment
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		URL:             getEnvRequired("DATABASE_URL"),
		MaxOpenConns:    getEnvAsInt("DATABASE_MAX_OPEN_CONNS", defaultMaxOpenConns),
		MaxIdleConns:    getEnvAsInt("DATABASE_MAX_IDLE_CONNS", defaultMaxIdleConns),
		MaxIdleTime:     getEnvAsDuration("DATABASE_MAX_IDLE_TIME", defaultMaxIdleTime.String()),
		ConnMaxLifetime: getEnvAsDuration("DATABASE_CONN_MAX_LIFETIME", defaultConnMaxLifetime.String()),
	}
}
