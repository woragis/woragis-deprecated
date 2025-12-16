package config

import (
	"fmt"
	"os"
	"strconv"
)

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
	QueueName     string
	Exchange      string
	RoutingKey    string
	PrefetchCount int
}

// LoadWorkerConfig reads worker settings from environment variables.
func LoadWorkerConfig() WorkerConfig {
	queueName := os.Getenv("WHATSAPP_QUEUE_NAME")
	if queueName == "" {
		queueName = "whatsapp.queue"
	}
	exchange := os.Getenv("WHATSAPP_EXCHANGE")
	if exchange == "" {
		exchange = "woragis.notifications"
	}
	routingKey := os.Getenv("WHATSAPP_ROUTING_KEY")
	if routingKey == "" {
		routingKey = "whatsapp.send"
	}
	prefetchCount := 1
	if raw := os.Getenv("WHATSAPP_PREFETCH_COUNT"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			prefetchCount = parsed
		}
	}

	return WorkerConfig{
		QueueName:     queueName,
		Exchange:      exchange,
		RoutingKey:    routingKey,
		PrefetchCount: prefetchCount,
	}
}

// WhatsAppConfig holds WhatsApp session configuration.
type WhatsAppConfig struct {
	SessionPath string
}

// LoadWhatsAppConfig reads WhatsApp settings from environment variables.
func LoadWhatsAppConfig() WhatsAppConfig {
	sessionPath := os.Getenv("WHATSAPP_SESSION_PATH")
	if sessionPath == "" {
		sessionPath = "/app/session"
	}

	return WhatsAppConfig{
		SessionPath: sessionPath,
	}
}
