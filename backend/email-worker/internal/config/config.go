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

// RabbitMQConfig holds RabbitMQ connection configuration.
type RabbitMQConfig struct {
	URL      string
	User     string
	Password string
	Host     string
	Port     string
	VHost    string
}

// LoadRabbitMQConfig reads RabbitMQ settings from environment variables.
func LoadRabbitMQConfig() RabbitMQConfig {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		user := os.Getenv("RABBITMQ_USER")
		if user == "" {
			user = "woragis"
		}
		password := os.Getenv("RABBITMQ_PASSWORD")
		if password == "" {
			password = "woragis"
		}
		host := os.Getenv("RABBITMQ_HOST")
		if host == "" {
			host = "rabbitmq"
		}
		port := os.Getenv("RABBITMQ_PORT")
		if port == "" {
			port = "5672"
		}
		vhost := os.Getenv("RABBITMQ_VHOST")
		if vhost == "" {
			vhost = "woragis"
		}
		// Remove leading slash if present
		if len(vhost) > 0 && vhost[0] == '/' {
			vhost = vhost[1:]
		}
		url = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost)
	}

	return RabbitMQConfig{
		URL:      url,
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		VHost:    os.Getenv("RABBITMQ_VHOST"),
	}
}

// WorkerConfig holds worker-specific configuration.
type WorkerConfig struct {
	QueueName    string
	Exchange     string
	RoutingKey   string
	PrefetchCount int
}

// LoadWorkerConfig reads worker settings from environment variables.
func LoadWorkerConfig() WorkerConfig {
	queueName := os.Getenv("EMAIL_QUEUE_NAME")
	if queueName == "" {
		queueName = "emails.queue"
	}
	exchange := os.Getenv("EMAIL_EXCHANGE")
	if exchange == "" {
		exchange = "woragis.notifications"
	}
	routingKey := os.Getenv("EMAIL_ROUTING_KEY")
	if routingKey == "" {
		routingKey = "emails.send"
	}
	prefetchCount := 1
	if raw := os.Getenv("EMAIL_PREFETCH_COUNT"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			prefetchCount = parsed
		}
	}

	return WorkerConfig{
		QueueName:    queueName,
		Exchange:     exchange,
		RoutingKey:   routingKey,
		PrefetchCount: prefetchCount,
	}
}
