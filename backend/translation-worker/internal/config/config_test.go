package config

import (
	"os"
	"testing"
)

func TestLoadRabbitMQConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		check func(t *testing.T, cfg RabbitMQConfig)
	}{
		{
			name: "with URL",
			env: map[string]string{
				"RABBITMQ_URL": "amqp://user:pass@host:5672/vhost",
			},
			check: func(t *testing.T, cfg RabbitMQConfig) {
				if cfg.URL != "amqp://user:pass@host:5672/vhost" {
					t.Errorf("URL = %v, want amqp://user:pass@host:5672/vhost", cfg.URL)
				}
			},
		},
		{
			name: "without URL, with individual values",
			env: map[string]string{
				"RABBITMQ_USER":     "testuser",
				"RABBITMQ_PASSWORD": "testpass",
				"RABBITMQ_HOST":     "testhost",
				"RABBITMQ_PORT":     "5673",
				"RABBITMQ_VHOST":    "testvhost",
			},
			check: func(t *testing.T, cfg RabbitMQConfig) {
				if cfg.URL != "amqp://testuser:testpass@testhost:5673/testvhost" {
					t.Errorf("URL = %v, want amqp://testuser:testpass@testhost:5673/testvhost", cfg.URL)
				}
			},
		},
		{
			name: "without URL, defaults",
			env:  map[string]string{},
			check: func(t *testing.T, cfg RabbitMQConfig) {
				if cfg.URL == "" {
					t.Error("URL should be generated from defaults")
				}
				if cfg.URL != "amqp://woragis:woragis@rabbitmq:5672/woragis" {
					t.Errorf("URL = %v, want amqp://woragis:woragis@rabbitmq:5672/woragis", cfg.URL)
				}
			},
		},
		{
			name: "vhost with leading slash",
			env: map[string]string{
				"RABBITMQ_VHOST": "/testvhost",
			},
			check: func(t *testing.T, cfg RabbitMQConfig) {
				if cfg.URL != "amqp://woragis:woragis@rabbitmq:5672/testvhost" {
					t.Errorf("URL should not have leading slash in vhost, got %v", cfg.URL)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			originalEnv := make(map[string]string)
			for k := range tt.env {
				originalEnv[k] = os.Getenv(k)
			}
			// Clear all RabbitMQ env vars first
			os.Unsetenv("RABBITMQ_URL")
			os.Unsetenv("RABBITMQ_USER")
			os.Unsetenv("RABBITMQ_PASSWORD")
			os.Unsetenv("RABBITMQ_HOST")
			os.Unsetenv("RABBITMQ_PORT")
			os.Unsetenv("RABBITMQ_VHOST")

			// Set test env
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			// Cleanup
			defer func() {
				for k, v := range originalEnv {
					if v == "" {
						os.Unsetenv(k)
					} else {
						os.Setenv(k, v)
					}
				}
			}()

			cfg := LoadRabbitMQConfig()
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestLoadDatabaseConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		check func(t *testing.T, cfg DatabaseConfig)
	}{
		{
			name: "with DATABASE_URL",
			env: map[string]string{
				"DATABASE_URL": "postgres://user:pass@host:5432/dbname?sslmode=disable",
			},
			check: func(t *testing.T, cfg DatabaseConfig) {
				if cfg.URL != "postgres://user:pass@host:5432/dbname?sslmode=disable" {
					t.Errorf("URL = %v, want postgres://user:pass@host:5432/dbname?sslmode=disable", cfg.URL)
				}
			},
		},
		{
			name: "without DATABASE_URL, with individual values",
			env: map[string]string{
				"DB_HOST":     "testhost",
				"DB_PORT":     "5433",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_SSLMODE":  "require",
			},
			check: func(t *testing.T, cfg DatabaseConfig) {
				if cfg.URL == "" {
					t.Error("URL should be generated from individual values")
				}
				expected := "host=testhost port=5433 user=testuser password=testpass dbname=testdb sslmode=require"
				if cfg.URL != expected {
					t.Errorf("URL = %v, want %v", cfg.URL, expected)
				}
			},
		},
		{
			name: "without DATABASE_URL, defaults",
			env:  map[string]string{},
			check: func(t *testing.T, cfg DatabaseConfig) {
				if cfg.URL == "" {
					t.Error("URL should be generated from defaults")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			originalEnv := make(map[string]string)
			for k := range tt.env {
				originalEnv[k] = os.Getenv(k)
			}
			// Clear all database env vars first
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_NAME")
			os.Unsetenv("DB_SSLMODE")

			// Set test env
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			// Cleanup
			defer func() {
				for k, v := range originalEnv {
					if v == "" {
						os.Unsetenv(k)
					} else {
						os.Setenv(k, v)
					}
				}
			}()

			cfg := LoadDatabaseConfig()
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestLoadTranslationConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		check func(t *testing.T, cfg TranslationConfig)
	}{
		{
			name: "google provider",
			env: map[string]string{
				"TRANSLATION_PROVIDER":      "google",
				"GOOGLE_TRANSLATE_API_KEY":  "test-key",
				"GOOGLE_CLOUD_PROJECT_ID":   "test-project",
			},
			check: func(t *testing.T, cfg TranslationConfig) {
				if cfg.Provider != ProviderGoogle {
					t.Errorf("Provider = %v, want google", cfg.Provider)
				}
				if cfg.GoogleAPIKey != "test-key" {
					t.Errorf("GoogleAPIKey = %v, want test-key", cfg.GoogleAPIKey)
				}
				if cfg.GoogleProjectID != "test-project" {
					t.Errorf("GoogleProjectID = %v, want test-project", cfg.GoogleProjectID)
				}
			},
		},
		{
			name: "deepl provider",
			env: map[string]string{
				"TRANSLATION_PROVIDER": "deepl",
				"DEEPL_API_KEY":       "test-deepl-key",
			},
			check: func(t *testing.T, cfg TranslationConfig) {
				if cfg.Provider != ProviderDeepL {
					t.Errorf("Provider = %v, want deepl", cfg.Provider)
				}
				if cfg.DeepLAPIKey != "test-deepl-key" {
					t.Errorf("DeepLAPIKey = %v, want test-deepl-key", cfg.DeepLAPIKey)
				}
			},
		},
		{
			name: "libre provider",
			env: map[string]string{
				"TRANSLATION_PROVIDER":    "libre",
				"LIBRE_TRANSLATE_API_URL": "https://custom.libretranslate.com/translate",
				"LIBRE_TRANSLATE_API_KEY": "test-libre-key",
			},
			check: func(t *testing.T, cfg TranslationConfig) {
				if cfg.Provider != ProviderLibre {
					t.Errorf("Provider = %v, want libre", cfg.Provider)
				}
				if cfg.LibreAPIURL != "https://custom.libretranslate.com/translate" {
					t.Errorf("LibreAPIURL = %v, want https://custom.libretranslate.com/translate", cfg.LibreAPIURL)
				}
				if cfg.LibreAPIKey != "test-libre-key" {
					t.Errorf("LibreAPIKey = %v, want test-libre-key", cfg.LibreAPIKey)
				}
			},
		},
		{
			name: "default provider (google)",
			env:  map[string]string{},
			check: func(t *testing.T, cfg TranslationConfig) {
				if cfg.Provider != ProviderGoogle {
					t.Errorf("Provider = %v, want google (default)", cfg.Provider)
				}
			},
		},
		{
			name: "custom timeout and retries",
			env: map[string]string{
				"TRANSLATION_TIMEOUT":     "60",
				"TRANSLATION_MAX_RETRIES": "5",
				"TRANSLATION_RETRY_DELAY": "2000",
			},
			check: func(t *testing.T, cfg TranslationConfig) {
				if cfg.Timeout != 60 {
					t.Errorf("Timeout = %v, want 60", cfg.Timeout)
				}
				if cfg.MaxRetries != 5 {
					t.Errorf("MaxRetries = %v, want 5", cfg.MaxRetries)
				}
				if cfg.RetryDelay != 2000 {
					t.Errorf("RetryDelay = %v, want 2000", cfg.RetryDelay)
				}
			},
		},
		{
			name: "invalid timeout defaults",
			env: map[string]string{
				"TRANSLATION_TIMEOUT": "invalid",
			},
			check: func(t *testing.T, cfg TranslationConfig) {
				if cfg.Timeout != 30 {
					t.Errorf("Timeout = %v, want 30 (default for invalid)", cfg.Timeout)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			originalEnv := make(map[string]string)
			for k := range tt.env {
				originalEnv[k] = os.Getenv(k)
			}
			// Clear all translation env vars first
			os.Unsetenv("TRANSLATION_PROVIDER")
			os.Unsetenv("GOOGLE_TRANSLATE_API_KEY")
			os.Unsetenv("GOOGLE_CLOUD_PROJECT_ID")
			os.Unsetenv("DEEPL_API_KEY")
			os.Unsetenv("LIBRE_TRANSLATE_API_URL")
			os.Unsetenv("LIBRE_TRANSLATE_API_KEY")
			os.Unsetenv("TRANSLATION_TIMEOUT")
			os.Unsetenv("TRANSLATION_MAX_RETRIES")
			os.Unsetenv("TRANSLATION_RETRY_DELAY")

			// Set test env
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			// Cleanup
			defer func() {
				for k, v := range originalEnv {
					if v == "" {
						os.Unsetenv(k)
					} else {
						os.Setenv(k, v)
					}
				}
			}()

			cfg := LoadTranslationConfig()
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestLoadWorkerConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		check func(t *testing.T, cfg WorkerConfig)
	}{
		{
			name: "all values set",
			env: map[string]string{
				"TRANSLATION_QUEUE_NAME":    "test.queue",
				"TRANSLATION_EXCHANGE":      "test.exchange",
				"TRANSLATION_ROUTING_KEY":   "test.routing",
				"TRANSLATION_PREFETCH_COUNT": "5",
			},
			check: func(t *testing.T, cfg WorkerConfig) {
				if cfg.QueueName != "test.queue" {
					t.Errorf("QueueName = %v, want test.queue", cfg.QueueName)
				}
				if cfg.Exchange != "test.exchange" {
					t.Errorf("Exchange = %v, want test.exchange", cfg.Exchange)
				}
				if cfg.RoutingKey != "test.routing" {
					t.Errorf("RoutingKey = %v, want test.routing", cfg.RoutingKey)
				}
				if cfg.PrefetchCount != 5 {
					t.Errorf("PrefetchCount = %v, want 5", cfg.PrefetchCount)
				}
			},
		},
		{
			name: "defaults",
			env:  map[string]string{},
			check: func(t *testing.T, cfg WorkerConfig) {
				if cfg.QueueName != "translations.queue" {
					t.Errorf("QueueName = %v, want translations.queue", cfg.QueueName)
				}
				if cfg.Exchange != "woragis.tasks" {
					t.Errorf("Exchange = %v, want woragis.tasks", cfg.Exchange)
				}
				if cfg.RoutingKey != "translations.process" {
					t.Errorf("RoutingKey = %v, want translations.process", cfg.RoutingKey)
				}
				if cfg.PrefetchCount != 1 {
					t.Errorf("PrefetchCount = %v, want 1 (default)", cfg.PrefetchCount)
				}
			},
		},
		{
			name: "invalid prefetch count",
			env: map[string]string{
				"TRANSLATION_PREFETCH_COUNT": "invalid",
			},
			check: func(t *testing.T, cfg WorkerConfig) {
				if cfg.PrefetchCount != 1 {
					t.Errorf("PrefetchCount = %v, want 1 (default for invalid)", cfg.PrefetchCount)
				}
			},
		},
		{
			name: "zero prefetch count",
			env: map[string]string{
				"TRANSLATION_PREFETCH_COUNT": "0",
			},
			check: func(t *testing.T, cfg WorkerConfig) {
				if cfg.PrefetchCount != 1 {
					t.Errorf("PrefetchCount = %v, want 1 (default for zero)", cfg.PrefetchCount)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			originalEnv := make(map[string]string)
			for k := range tt.env {
				originalEnv[k] = os.Getenv(k)
			}
			// Clear all worker env vars first
			os.Unsetenv("TRANSLATION_QUEUE_NAME")
			os.Unsetenv("TRANSLATION_EXCHANGE")
			os.Unsetenv("TRANSLATION_ROUTING_KEY")
			os.Unsetenv("TRANSLATION_PREFETCH_COUNT")

			// Set test env
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			// Cleanup
			defer func() {
				for k, v := range originalEnv {
					if v == "" {
						os.Unsetenv(k)
					} else {
						os.Setenv(k, v)
					}
				}
			}()

			cfg := LoadWorkerConfig()
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}
