package config

import (
	"fmt"
	"os"
)

// RabbitMQConfig holds RabbitMQ connection configuration
type RabbitMQConfig struct {
	URL      string
	User     string
	Password string
	Host     string
	Port     string
	VHost    string
}

// LoadRabbitMQConfig loads RabbitMQ configuration from environment variables
func LoadRabbitMQConfig() *RabbitMQConfig {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		// Build URL from components
		user := getEnv("RABBITMQ_USER", "woragis")
		password := getEnv("RABBITMQ_PASSWORD", "woragis")
		host := getEnv("RABBITMQ_HOST", "rabbitmq")
		port := getEnv("RABBITMQ_PORT", "5672")
		vhost := getEnv("RABBITMQ_VHOST", "/")
		
		// Remove leading slash if present
		if len(vhost) > 0 && vhost[0] == '/' {
			vhost = vhost[1:]
		}
		
		url = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost)
	}

	return &RabbitMQConfig{
		URL:      url,
		User:     getEnv("RABBITMQ_USER", "woragis"),
		Password: getEnv("RABBITMQ_PASSWORD", "woragis"),
		Host:     getEnv("RABBITMQ_HOST", "rabbitmq"),
		Port:     getEnv("RABBITMQ_PORT", "5672"),
		VHost:    getEnv("RABBITMQ_VHOST", "/"),
	}
}

