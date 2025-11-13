package config

import (
	"os"
)

// AIConfig holds API keys and defaults for AI providers.
type AIConfig struct {
	OpenAIKey     string
	AnthropicKey  string
	GrokKey       string
	ManusKey      string
	CipherKey     string
	DefaultModel  string
	DefaultAlias  string
	ProviderAlias string
}

// LoadAIConfig reads AI-related environment variables.
func LoadAIConfig() AIConfig {
	return AIConfig{
		OpenAIKey:     os.Getenv("OPENAI_API_KEY"),
		AnthropicKey:  os.Getenv("ANTHROPIC_API_KEY"),
		GrokKey:       os.Getenv("GROK_API_KEY"),
		ManusKey:      os.Getenv("MANUS_API_KEY"),
		CipherKey:     os.Getenv("CIPHER_API_KEY"),
		DefaultModel:  os.Getenv("CHAT_MODEL"),
		DefaultAlias:  os.Getenv("CHAT_ALIAS"),
		ProviderAlias: os.Getenv("CHAT_PROVIDER"),
	}
}

