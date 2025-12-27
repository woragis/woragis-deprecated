package langchain

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateChatCompletionRequest validates a chat completion request
func ValidateChatCompletionRequest(req ChatCompletionRequest) error {
	// Validate provider (required)
	if req.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	validProviders := []ModelProvider{
		ProviderOpenAI,
		ProviderAnthropic,
		ProviderXAI,
		ProviderManus,
		ProviderCipher,
	}
	isValid := false
	for _, validProvider := range validProviders {
		if req.Provider == validProvider {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("provider: must be one of: openai, anthropic, xai, manus, cipher")
	}

	// Validate model (optional, but if provided, validate)
	if req.Model != "" {
		if err := validation.ValidateString(req.Model, 1, 100, "model"); err != nil {
			return fmt.Errorf("model: %w", err)
		}
	}

	// Validate messages (required, at least one)
	if len(req.Messages) == 0 {
		return fmt.Errorf("messages: at least one message is required")
	}
	if len(req.Messages) > 100 {
		return fmt.Errorf("messages: too many messages (maximum 100)")
	}

	// Validate each message
	for i, msg := range req.Messages {
		// Validate role (required)
		if msg.Role == "" {
			return fmt.Errorf("messages[%d].role is required", i)
		}
		validRoles := []string{"user", "assistant", "system", "ai", "human"}
		isValid := false
		for _, validRole := range validRoles {
			if msg.Role == validRole {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("messages[%d].role: must be one of: user, assistant, system, ai, human", i)
		}

		// Validate content (required, 1-100000 chars)
		if err := validation.ValidateString(msg.Content, 1, 100000, fmt.Sprintf("messages[%d].content", i)); err != nil {
			return fmt.Errorf("messages[%d].content: %w", i, err)
		}
		// Check for SQL injection and XSS (for user/system messages)
		if msg.Role == "user" || msg.Role == "system" {
			if err := validation.ValidateNoSQLInjection(msg.Content); err != nil {
				return fmt.Errorf("messages[%d].content: %w", i, err)
			}
			if err := validation.ValidateNoXSS(msg.Content); err != nil {
				return fmt.Errorf("messages[%d].content: %w", i, err)
			}
		}
	}

	// Validate max tokens (optional, but if provided, validate range)
	if req.MaxTokens > 0 {
		if req.MaxTokens < 1 {
			return fmt.Errorf("max_tokens: must be at least 1")
		}
		if req.MaxTokens > 100000 {
			return fmt.Errorf("max_tokens: must be at most 100,000")
		}
	}

	// Validate temperature (optional, but if provided, validate range)
	if req.Temperature != 0 {
		if req.Temperature < 0 {
			return fmt.Errorf("temperature: must be at least 0")
		}
		if req.Temperature > 2 {
			return fmt.Errorf("temperature: must be at most 2")
		}
	}

	// Validate agent (optional, but if provided, validate)
	if req.Agent != "" {
		if err := validation.ValidateString(req.Agent, 1, 100, "agent"); err != nil {
			return fmt.Errorf("agent: %w", err)
		}
	}

	return nil
}

