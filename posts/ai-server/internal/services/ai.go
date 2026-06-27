package services

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/woragis/posts-ai-service/internal/models"
)

type AIService struct {
	baseURL string
	client  *http.Client
}

func NewAIService(baseURL string) *AIService {
	return &AIService{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// ChatStream sends a request to AI service and returns streaming response
func (s *AIService) ChatStream(ctx context.Context, agent, input string, options ...func(*ChatRequest)) (io.ReadCloser, error) {
	req := &ChatRequest{
		Agent:    agent,
		Input:    input,
		Provider: "openai",
	}

	// Apply options
	for _, opt := range options {
		opt(req)
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/chat/stream", s.baseURL), strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("AI service error: %d - %s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}

// ParseStreamChunk parses a single NDJSON line into a StreamChunk
func (s *AIService) ParseStreamChunk(line string) (*models.StreamChunk, error) {
	if line == "" {
		return nil, nil
	}

	var chunk models.StreamChunk
	if err := json.Unmarshal([]byte(line), &chunk); err != nil {
		return nil, err
	}

	return &chunk, nil
}

// ScanStream reads from response body and yields chunks
func (s *AIService) ScanStream(body io.ReadCloser) <-chan *models.StreamChunk {
	out := make(chan *models.StreamChunk)

	go func() {
		defer close(out)
		defer body.Close()

		scanner := bufio.NewScanner(body)
		for scanner.Scan() {
			line := scanner.Text()
			chunk, err := s.ParseStreamChunk(line)
			if err != nil {
				out <- &models.StreamChunk{Error: fmt.Sprintf("Parse error: %v", err)}
				continue
			}
			if chunk != nil {
				out <- chunk
			}
		}

		if err := scanner.Err(); err != nil {
			out <- &models.StreamChunk{Error: fmt.Sprintf("Stream error: %v", err)}
		}
	}()

	return out
}

// ChatRequest matches AI service API
type ChatRequest struct {
	Agent       string  `json:"agent"`
	Input       string  `json:"input"`
	System      string  `json:"system,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	Model       string  `json:"model,omitempty"`
	Provider    string  `json:"provider,omitempty"`
}

// WithSystem sets system prompt
func WithSystem(system string) func(*ChatRequest) {
	return func(r *ChatRequest) {
		r.System = system
	}
}

// WithTemperature sets temperature
func WithTemperature(temp float64) func(*ChatRequest) {
	return func(r *ChatRequest) {
		r.Temperature = temp
	}
}

// WithProvider sets provider
func WithProvider(provider string) func(*ChatRequest) {
	return func(r *ChatRequest) {
		r.Provider = provider
	}
}
