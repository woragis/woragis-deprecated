package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application level configuration
type Config struct {
	AppName   string
	Port      string
	Env       string
	PublicURL string
}

// Load reads configuration from environment variables with sane defaults
func Load() *Config {
	// Load .env file in development
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	return &Config{
		AppName:   getEnv("APP_NAME", "woragis-social-media-service"),
		Port:      getEnv("PORT", "3000"),
		Env:       getEnv("ENV", "development"),
		PublicURL: getEnv("APP_PUBLIC_URL", "http://localhost:3000"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvRequired gets a required environment variable, panics if not found
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer with a fallback
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsInt64 gets an environment variable as an int64 with a fallback
func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as a duration with a fallback
func getEnvAsDuration(key string, defaultValue string) time.Duration {
	valueStr := getEnv(key, defaultValue)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

// getEnvAsSlice gets an environment variable as a slice of strings
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}