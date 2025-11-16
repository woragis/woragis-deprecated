package langchain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"bytes"
	"encoding/json"
	"bufio"
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
	// Route all providers through the dedicated AI service to standardize behavior,
	// falling back to direct OpenAI if AI service URL is not configured.
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		switch req.Provider {
		case ProviderOpenAI:
			return c.callOpenAI(ctx, req)
		default:
			return ChatCompletionResponse{}, fmt.Errorf("AI_SERVICE_URL not set and provider %q not supported directly", req.Provider)
		}
	}
	return c.callAIService(ctx, aiURL, req)
}

type aiServiceChatRequest struct {
	Agent      string  `json:"agent"`
	Input      string  `json:"input"`
	System     string  `json:"system,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	Model      string  `json:"model,omitempty"`
	Provider   string  `json:"provider,omitempty"`
}

type aiServiceChatResponse struct {
	Agent  string `json:"agent"`
	Output string `json:"output"`
}

func (c *Client) callAIService(ctx context.Context, baseURL string, req ChatCompletionRequest) (ChatCompletionResponse, error) {
	// Concatenate messages into a single input preserving roles.
	var b strings.Builder
	for _, m := range req.Messages {
		role := strings.ToLower(m.Role)
		switch role {
		case "system":
			// Prepend as system guidance
			// We'll include it in the input as a labeled line to give model context.
			b.WriteString("System: ")
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		case "assistant", "ai":
			b.WriteString("Assistant: ")
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		default:
			b.WriteString("User: ")
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		}
	}
	input := strings.TrimSpace(b.String())

	var tempPtr *float64
	if req.Temperature != 0 {
		t := req.Temperature
		tempPtr = &t
	}

	payload := aiServiceChatRequest{
		Agent:       "startup", // default agent persona
		Input:       input,
		Temperature: tempPtr,
		Model:       req.Model,
		Provider:    string(req.Provider),
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	httpClient := &http.Client{Timeout: 75 * time.Second}
	url := strings.TrimRight(baseURL, "/") + "/v1/chat"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(httpReq)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("ai-service request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return ChatCompletionResponse{}, fmt.Errorf("ai-service error: %s", res.Status)
	}

	var data aiServiceChatResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("ai-service decode: %w", err)
	}

	msg := ChatMessage{
		Role:      "assistant",
		Content:   data.Output,
		Timestamp: time.Now().UTC(),
	}
	return ChatCompletionResponse{Message: msg, Raw: data}, nil
}

// GenerateCompletionStream streams deltas from the AI service and invokes onDelta for each.
func (c *Client) GenerateCompletionStream(ctx context.Context, req ChatCompletionRequest, onDelta func(delta string)) (ChatCompletionResponse, error) {
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		// Fallback: no streaming support directly; do a normal call
		return c.GenerateCompletion(ctx, req)
	}

	// Build input text same as non-streaming path
	var b strings.Builder
	for _, m := range req.Messages {
		role := strings.ToLower(m.Role)
		switch role {
		case "system":
			b.WriteString("System: ")
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		case "assistant", "ai":
			b.WriteString("Assistant: ")
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		default:
			b.WriteString("User: ")
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		}
	}
	input := strings.TrimSpace(b.String())

	var tempPtr *float64
	if req.Temperature != 0 {
		t := req.Temperature
		tempPtr = &t
	}
	payload := aiServiceChatRequest{
		Agent:       "startup",
		Input:       input,
		Temperature: tempPtr,
		Model:       req.Model,
		Provider:    string(req.Provider),
	}
	buf, err := json.Marshal(&payload)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	httpClient := &http.Client{Timeout: 0} // stream until server closes
	url := strings.TrimRight(aiURL, "/") + "/v1/chat/stream"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(httpReq)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("ai-service stream request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return ChatCompletionResponse{}, fmt.Errorf("ai-service stream error: %s", res.Status)
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var full strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var evt map[string]any
		if err := json.Unmarshal([]byte(line), &evt); err != nil {
			continue
		}
		if delta, ok := evt["delta"].(string); ok {
			full.WriteString(delta)
			if onDelta != nil {
				onDelta(delta)
			}
			continue
		}
		if done, ok := evt["done"].(bool); ok && done {
			if out, ok := evt["output"].(string); ok && out != "" {
				full.Reset()
				full.WriteString(out)
			}
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("ai-service stream read: %w", err)
	}

	msg := ChatMessage{
		Role:      "assistant",
		Content:   full.String(),
		Timestamp: time.Now().UTC(),
	}
	return ChatCompletionResponse{Message: msg, Raw: nil}, nil
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
