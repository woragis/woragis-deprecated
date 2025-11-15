package config

import (
	"os"
)

// WhatsAppConfig holds WhatsApp configuration for whatsmeow.
type WhatsAppConfig struct {
	SessionPath string
	enabled     bool
}

// LoadWhatsAppConfig reads WhatsApp settings from environment variables.
func LoadWhatsAppConfig() WhatsAppConfig {
	sessionPath := os.Getenv("WHATSAPP_SESSION_PATH")
	if sessionPath == "" {
		sessionPath = "./whatsapp-session"
	}

	enabled := os.Getenv("WHATSAPP_ENABLED") == "true"

	return WhatsAppConfig{
		SessionPath: sessionPath,
		enabled:     enabled,
	}
}

// Enabled returns true when WhatsApp is enabled.
func (c WhatsAppConfig) Enabled() bool {
	return c.enabled
}

