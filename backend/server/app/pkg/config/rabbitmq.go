package config

import "os"

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
		url = "amqp://guest:guest@localhost:5672/"
	}

	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest"
	}

	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		password = "guest"
	}

	vhost := os.Getenv("RABBITMQ_VHOST")
	if vhost == "" {
		vhost = "/"
	}

	return RabbitMQConfig{
		URL:      url,
		User:     user,
		Password: password,
		VHost:    vhost,
	}
}
