package config

import (
	"fmt"
	"os"
	"strconv"
)

// EmailConfig holds SMTP configuration for transactional emails.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

// LoadEmailConfig reads SMTP settings from environment variables.
func LoadEmailConfig() (EmailConfig, error) {
	host := os.Getenv("SMTP_HOST")
	port := 587
	if raw := os.Getenv("SMTP_PORT"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return EmailConfig{}, fmt.Errorf("invalid SMTP_PORT %q: %w", raw, err)
		}
		port = parsed
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	useTLS := os.Getenv("SMTP_TLS") != "false"

	return EmailConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		UseTLS:   useTLS,
	}, nil
}

// Enabled returns true when SMTP is configured.
func (c EmailConfig) Enabled() bool {
	return c.Host != "" && c.From != ""
}

// Address returns host:port combination.
func (c EmailConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
