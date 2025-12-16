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

func TestLoadWorkerConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		check func(t *testing.T, cfg WorkerConfig)
	}{
		{
			name: "all values set",
			env: map[string]string{
				"WHATSAPP_QUEUE_NAME":    "test.queue",
				"WHATSAPP_EXCHANGE":      "test.exchange",
				"WHATSAPP_ROUTING_KEY":   "test.routing",
				"WHATSAPP_PREFETCH_COUNT": "5",
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
				if cfg.QueueName != "whatsapp.queue" {
					t.Errorf("QueueName = %v, want whatsapp.queue", cfg.QueueName)
				}
				if cfg.Exchange != "woragis.notifications" {
					t.Errorf("Exchange = %v, want woragis.notifications", cfg.Exchange)
				}
				if cfg.RoutingKey != "whatsapp.send" {
					t.Errorf("RoutingKey = %v, want whatsapp.send", cfg.RoutingKey)
				}
				if cfg.PrefetchCount != 1 {
					t.Errorf("PrefetchCount = %v, want 1 (default)", cfg.PrefetchCount)
				}
			},
		},
		{
			name: "invalid prefetch count",
			env: map[string]string{
				"WHATSAPP_PREFETCH_COUNT": "invalid",
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
				"WHATSAPP_PREFETCH_COUNT": "0",
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
			os.Unsetenv("WHATSAPP_QUEUE_NAME")
			os.Unsetenv("WHATSAPP_EXCHANGE")
			os.Unsetenv("WHATSAPP_ROUTING_KEY")
			os.Unsetenv("WHATSAPP_PREFETCH_COUNT")

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

func TestLoadWhatsAppConfig(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		check func(t *testing.T, cfg WhatsAppConfig)
	}{
		{
			name: "with session path",
			env: map[string]string{
				"WHATSAPP_SESSION_PATH": "/custom/session",
			},
			check: func(t *testing.T, cfg WhatsAppConfig) {
				if cfg.SessionPath != "/custom/session" {
					t.Errorf("SessionPath = %v, want /custom/session", cfg.SessionPath)
				}
			},
		},
		{
			name: "default session path",
			env:  map[string]string{},
			check: func(t *testing.T, cfg WhatsAppConfig) {
				if cfg.SessionPath != "/app/session" {
					t.Errorf("SessionPath = %v, want /app/session (default)", cfg.SessionPath)
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
			// Clear WhatsApp env var first
			os.Unsetenv("WHATSAPP_SESSION_PATH")

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

			cfg := LoadWhatsAppConfig()
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}
