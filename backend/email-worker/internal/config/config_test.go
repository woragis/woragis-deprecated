package config

import (
	"os"
	"testing"
)

func TestLoadEmailConfig(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
		check   func(t *testing.T, cfg EmailConfig)
	}{
		{
			name: "valid config",
			env: map[string]string{
				"SMTP_HOST":     "smtp.example.com",
				"SMTP_PORT":     "587",
				"SMTP_USERNAME": "user",
				"SMTP_PASSWORD": "pass",
				"SMTP_FROM":     "noreply@example.com",
				"SMTP_TLS":      "true",
			},
			wantErr: false,
			check: func(t *testing.T, cfg EmailConfig) {
				if cfg.Host != "smtp.example.com" {
					t.Errorf("Host = %v, want smtp.example.com", cfg.Host)
				}
				if cfg.Port != 587 {
					t.Errorf("Port = %v, want 587", cfg.Port)
				}
				if cfg.Username != "user" {
					t.Errorf("Username = %v, want user", cfg.Username)
				}
				if cfg.Password != "pass" {
					t.Errorf("Password = %v, want pass", cfg.Password)
				}
				if cfg.From != "noreply@example.com" {
					t.Errorf("From = %v, want noreply@example.com", cfg.From)
				}
				if !cfg.UseTLS {
					t.Errorf("UseTLS = %v, want true", cfg.UseTLS)
				}
			},
		},
		{
			name: "default port",
			env: map[string]string{
				"SMTP_HOST": "smtp.example.com",
				"SMTP_FROM": "noreply@example.com",
			},
			wantErr: false,
			check: func(t *testing.T, cfg EmailConfig) {
				if cfg.Port != 587 {
					t.Errorf("Port = %v, want 587 (default)", cfg.Port)
				}
			},
		},
		{
			name: "invalid port",
			env: map[string]string{
				"SMTP_HOST": "smtp.example.com",
				"SMTP_PORT": "invalid",
				"SMTP_FROM": "noreply@example.com",
			},
			wantErr: true,
		},
		{
			name: "TLS disabled",
			env: map[string]string{
				"SMTP_HOST": "smtp.example.com",
				"SMTP_FROM": "noreply@example.com",
				"SMTP_TLS":  "false",
			},
			wantErr: false,
			check: func(t *testing.T, cfg EmailConfig) {
				if cfg.UseTLS {
					t.Errorf("UseTLS = %v, want false", cfg.UseTLS)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			originalEnv := make(map[string]string)
			for k, v := range tt.env {
				originalEnv[k] = os.Getenv(k)
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

			cfg, err := LoadEmailConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadEmailConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestEmailConfig_Enabled(t *testing.T) {
	tests := []struct {
		name string
		cfg  EmailConfig
		want bool
	}{
		{
			name: "enabled with host and from",
			cfg: EmailConfig{
				Host: "smtp.example.com",
				From: "noreply@example.com",
			},
			want: true,
		},
		{
			name: "disabled without host",
			cfg: EmailConfig{
				Host: "",
				From: "noreply@example.com",
			},
			want: false,
		},
		{
			name: "disabled without from",
			cfg: EmailConfig{
				Host: "smtp.example.com",
				From: "",
			},
			want: false,
		},
		{
			name: "disabled without both",
			cfg: EmailConfig{
				Host: "",
				From: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.Enabled(); got != tt.want {
				t.Errorf("EmailConfig.Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailConfig_Address(t *testing.T) {
	cfg := EmailConfig{
		Host: "smtp.example.com",
		Port: 587,
	}
	want := "smtp.example.com:587"
	if got := cfg.Address(); got != want {
		t.Errorf("EmailConfig.Address() = %v, want %v", got, want)
	}
}

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
				"EMAIL_QUEUE_NAME":    "test.queue",
				"EMAIL_EXCHANGE":      "test.exchange",
				"EMAIL_ROUTING_KEY":   "test.routing",
				"EMAIL_PREFETCH_COUNT": "5",
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
				if cfg.QueueName != "emails.queue" {
					t.Errorf("QueueName = %v, want emails.queue", cfg.QueueName)
				}
				if cfg.Exchange != "woragis.notifications" {
					t.Errorf("Exchange = %v, want woragis.notifications", cfg.Exchange)
				}
				if cfg.RoutingKey != "emails.send" {
					t.Errorf("RoutingKey = %v, want emails.send", cfg.RoutingKey)
				}
				if cfg.PrefetchCount != 1 {
					t.Errorf("PrefetchCount = %v, want 1 (default)", cfg.PrefetchCount)
				}
			},
		},
		{
			name: "invalid prefetch count",
			env: map[string]string{
				"EMAIL_PREFETCH_COUNT": "invalid",
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
				"EMAIL_PREFETCH_COUNT": "0",
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
			os.Unsetenv("EMAIL_QUEUE_NAME")
			os.Unsetenv("EMAIL_EXCHANGE")
			os.Unsetenv("EMAIL_ROUTING_KEY")
			os.Unsetenv("EMAIL_PREFETCH_COUNT")

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
