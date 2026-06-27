package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           int
	DatabaseURL    string
	AIServiceURL   string
	CORSOrigins    string
	Environment    string
	LogLevel       string
}

func LoadConfig() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "3014"))
	
	return &Config{
		Port:         port,
		DatabaseURL: getEnv("DATABASE_URL", "postgres://woragis:password@localhost:5432/posts_ai?sslmode=disable"),
		AIServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:8000"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:3000"),
		Environment: getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
