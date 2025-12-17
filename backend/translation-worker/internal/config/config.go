package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	URL string
}

// LoadDatabaseConfig reads database settings from environment variables.
func LoadDatabaseConfig() DatabaseConfig {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		// Fallback to individual components
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "woragis"
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "woragis"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "woragis"
		}
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	}

	return DatabaseConfig{
		URL: url,
	}
}

// TranslationProvider represents the translation API provider.
type TranslationProvider string

const (
	ProviderGoogle TranslationProvider = "google"
	ProviderDeepL  TranslationProvider = "deepl"
	ProviderLibre  TranslationProvider = "libre"
)

// TranslationConfig holds translation API configuration.
type TranslationConfig struct {
	Provider TranslationProvider
	// Google Translate API
	GoogleAPIKey string
	GoogleProjectID string
	// DeepL API
	DeepLAPIKey string
	// LibreTranslate API (self-hosted)
	LibreAPIURL string
	LibreAPIKey string
	// Timeout for API requests (seconds)
	Timeout int
	// Retry configuration
	MaxRetries int
	RetryDelay int // milliseconds
}

// LoadTranslationConfig reads translation API settings from environment variables.
func LoadTranslationConfig() TranslationConfig {
	provider := TranslationProvider(strings.ToLower(os.Getenv("TRANSLATION_PROVIDER")))
	if provider == "" {
		provider = ProviderGoogle // Default to Google
	}

	timeout := 30
	if raw := os.Getenv("TRANSLATION_TIMEOUT"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			timeout = parsed
		}
	}

	maxRetries := 3
	if raw := os.Getenv("TRANSLATION_MAX_RETRIES"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			maxRetries = parsed
		}
	}

	retryDelay := 1000
	if raw := os.Getenv("TRANSLATION_RETRY_DELAY"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			retryDelay = parsed
		}
	}

	return TranslationConfig{
		Provider:      provider,
		GoogleAPIKey:  os.Getenv("GOOGLE_TRANSLATE_API_KEY"),
		GoogleProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT_ID"),
		DeepLAPIKey:   os.Getenv("DEEPL_API_KEY"),
		LibreAPIURL:   os.Getenv("LIBRE_TRANSLATE_API_URL"),
		LibreAPIKey:   os.Getenv("LIBRE_TRANSLATE_API_KEY"),
		Timeout:       timeout,
		MaxRetries:    maxRetries,
		RetryDelay:    retryDelay,
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
	queueName := os.Getenv("TRANSLATION_QUEUE_NAME")
	if queueName == "" {
		queueName = "translations.queue"
	}
	exchange := os.Getenv("TRANSLATION_EXCHANGE")
	if exchange == "" {
		exchange = "woragis.tasks"
	}
	routingKey := os.Getenv("TRANSLATION_ROUTING_KEY")
	if routingKey == "" {
		routingKey = "translations.process"
	}
	prefetchCount := 1
	if raw := os.Getenv("TRANSLATION_PREFETCH_COUNT"); raw != "" {
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
