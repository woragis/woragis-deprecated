package config

import (
	"fmt"
	"os"
)

// RabbitMQConfig holds connection details for RabbitMQ.
type RabbitMQConfig struct {
	URL      string
	User     string
	Password string
	VHost    string
}

// LoadRabbitMQConfig reads RabbitMQ configuration from environment variables.
func LoadRabbitMQConfig() RabbitMQConfig {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		// Build URL from components
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

	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "woragis"
	}

	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		password = "woragis"
	}

	vhost := os.Getenv("RABBITMQ_VHOST")
	if vhost == "" {
		vhost = "woragis"
	}

	return RabbitMQConfig{
		URL:      url,
		User:     user,
		Password: password,
		VHost:    vhost,
	}
}
