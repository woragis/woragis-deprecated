package config

import (
	"fmt"
	"strconv"
)

// EmailConfig holds SMTP configuration for transactional emails
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

// LoadEmailConfig reads SMTP settings from environment variables
func LoadEmailConfig() (*EmailConfig, error) {
	host := getEnv("SMTP_HOST", "")
	port := 587
	if raw := getEnv("SMTP_PORT", ""); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid SMTP_PORT %q: %w", raw, err)
		}
		port = parsed
	}

	username := getEnv("SMTP_USERNAME", "")
	password := getEnv("SMTP_PASSWORD", "")
	from := getEnv("EMAIL_FROM", "")
	useTLS := getEnv("SMTP_TLS", "true") != "false"

	return &EmailConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		UseTLS:   useTLS,
	}, nil
}

// Enabled returns true when SMTP is configured
func (c *EmailConfig) Enabled() bool {
	return c.Host != "" && c.From != ""
}

// Address returns host:port combination
func (c *EmailConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
