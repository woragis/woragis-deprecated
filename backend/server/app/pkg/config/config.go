package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application level configuration.
type Config struct {
	AppName   string
	Port      int
	Env       string
	PublicURL string
}

// Load reads configuration from environment variables with sane defaults.
func Load() (Config, error) {
	name := getEnv("APP_NAME", "woragis-server")
	env := getEnv("APP_ENV", "development")
	port, err := parsePort(getEnv("APP_PORT", "8080"))
	if err != nil {
		return Config{}, err
	}

	publicURL := getEnv("APP_PUBLIC_URL", "http://localhost:8080")

	return Config{
		AppName:   name,
		Port:      port,
		Env:       env,
		PublicURL: publicURL,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parsePort(value string) (int, error) {
	port, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: %w", value, err)
	}

	if port <= 0 || port > 65535 {
		return 0, fmt.Errorf("port out of range: %d", port)
	}

	return port, nil
}
