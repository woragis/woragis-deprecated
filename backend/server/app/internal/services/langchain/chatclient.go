package langchain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// ModelProvider enumerates supported chat providers.
type ModelProvider string

const (
	ProviderOpenAI      ModelProvider = "openai"
	ProviderAnthropic   ModelProvider = "anthropic"
	ProviderXAI         ModelProvider = "xai"
	ProviderManus       ModelProvider = "manus"
	ProviderCipher      ModelProvider = "cipher"
	ProviderUnspecified ModelProvider = ""
)

// ModelSpec describes an LLM preset.
type ModelSpec struct {
	Alias       string
	Provider    ModelProvider
	Model       string
	Description string
	Specialty   string
}

var modelCatalog = map[string]ModelSpec{
	"grok": {
		Alias:       "grok",
		Provider:    ProviderXAI,
		Model:       "grok-1",
		Description: "xAI Grok model tuned for real-time, trending, and political insights.",
		Specialty:   "recency, world-events, politics",
	},
	"claude": {
		Alias:       "claude",
		Provider:    ProviderAnthropic,
		Model:       "claude-3-5-sonnet",
		Description: "Anthropic Claude for nuanced, thoughtful reasoning and structured responses.",
		Specialty:   "analytical reasoning, structured plans",
	},
	"manus": {
		Alias:       "manus",
		Provider:    ProviderManus,
		Model:       "manus-intellect",
		Description: "Experimental Manus intelligence model for deep strategic insights.",
		Specialty:   "strategy, advanced analysis",
	},
	"chatgpt": {
		Alias:       "chatgpt",
		Provider:    ProviderOpenAI,
		Model:       "gpt-4o-mini",
		Description: "OpenAI ChatGPT for general-purpose assistance.",
		Specialty:   "general assistance, drafting, Q&A",
	},
	"cipher": {
		Alias:       "cipher",
		Provider:    ProviderCipher,
		Model:       "cipher-shadow",
		Description: "Private Cipher endpoint for confidential or sensitive brainstorming.",
		Specialty:   "discreet discussions, dark-mode ideation",
	},
}

// LookupModel retrieves a model specification by alias.
func LookupModel(alias string) (ModelSpec, bool) {
	if alias == "" {
		return ModelSpec{}, false
	}
	spec, ok := modelCatalog[strings.ToLower(alias)]
	return spec, ok
}

// ChatMessage represents a chat interaction stored by the service.
type ChatMessage struct {
	Role      string
	Content   string
	Timestamp time.Time
}

// ChatCompletionRequest groups request parameters.
type ChatCompletionRequest struct {
	Provider    ModelProvider
	Model       string
	Messages    []ChatMessage
	MaxTokens   int
	Temperature float64
}

// ChatCompletionResponse wraps the AI response.
type ChatCompletionResponse struct {
	Message ChatMessage
	Raw     any
}

// Client exposes LangChain backed LLM interactions.
type Client struct {
	logger *slog.Logger
}

// NewClient creates a new Client.
func NewClient(logger *slog.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

// GenerateCompletion executes a chat completion call using the requested provider.
func (c *Client) GenerateCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
	switch req.Provider {
	case ProviderOpenAI:
		return c.callOpenAI(ctx, req)
	default:
		return ChatCompletionResponse{}, fmt.Errorf("provider %q not implemented", req.Provider)
	}
}

func (c *Client) callOpenAI(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return ChatCompletionResponse{}, errors.New("OPENAI_API_KEY not set")
	}

	client, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithModel(req.Model),
	)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("create openai client: %w", err)
	}

	lcMessages := make([]llms.MessageContent, 0, len(req.Messages))
	for _, msg := range req.Messages {
		role := mapRole(msg.Role)
		lcMessages = append(lcMessages, llms.MessageContent{
			Role:  role,
			Parts: []llms.ContentPart{llms.TextPart(msg.Content)},
		})
	}

	options := []llms.CallOption{}
	if req.MaxTokens > 0 {
		options = append(options, llms.WithMaxTokens(req.MaxTokens))
	}
	if req.Temperature != 0 {
		options = append(options, llms.WithTemperature(req.Temperature))
	}

	response, err := client.GenerateContent(ctx, lcMessages, options...)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("openai completion: %w", err)
	}

	if response == nil || len(response.Choices) == 0 {
		return ChatCompletionResponse{}, errors.New("empty response from provider")
	}

	choice := response.Choices[0]

	msg := ChatMessage{
		Role:      "assistant",
		Content:   choice.Content,
		Timestamp: time.Now().UTC(),
	}

	return ChatCompletionResponse{
		Message: msg,
		Raw:     response,
	}, nil
}

func mapRole(role string) llms.ChatMessageType {
	switch strings.ToLower(role) {
	case "system":
		return llms.ChatMessageTypeSystem
	case "assistant", "ai":
		return llms.ChatMessageTypeAI
	case "user", "human":
		return llms.ChatMessageTypeHuman
	default:
		return llms.ChatMessageTypeGeneric
	}
}
