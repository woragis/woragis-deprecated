package creative

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sony/gobreaker"
	appcircuitbreaker "github.com/woragis/backend/server/app/pkg/circuitbreaker"
	"github.com/woragis/backend/server/app/pkg/validation"
)

// Client interacts with the creative-service API
type Client struct {
	baseURL    string
	httpClient *http.Client
	cb         *gobreaker.CircuitBreaker
	logger     *slog.Logger
}

// NewClient creates a new creative service client
func NewClient(baseURL string, logger *slog.Logger) *Client {
	if baseURL == "" {
		baseURL = os.Getenv("CREATIVE_SERVICE_URL")
		if baseURL == "" {
			baseURL = "http://creative-service:8000"
		}
	}
	
	// Create circuit breaker for Creative Service calls
	cbConfig := appcircuitbreaker.DefaultConfig("creative-service", logger)
	cbConfig.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		appcircuitbreaker.RecordStateChange(name, from, to)
		if logger != nil {
			logger.Info("circuit breaker state changed",
				slog.String("name", name),
				slog.String("from", from.String()),
				slog.String("to", to.String()),
			)
		}
	}
	cb := appcircuitbreaker.NewCircuitBreaker(cbConfig)
	
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for image generation
		},
		cb:     cb,
		logger: logger,
	}
}

// ImageProvider represents supported image providers
type ImageProvider string

const (
	ProviderOpenAI         ImageProvider = "openai"
	ProviderStableDiffusion ImageProvider = "stable-diffusion"
	ProviderCipher         ImageProvider = "cipher"
)

// ImageStyle represents image generation styles
type ImageStyle string

const (
	StyleTechnical     ImageStyle = "technical"
	StyleDiagram       ImageStyle = "diagram"
	StyleThumbnail     ImageStyle = "thumbnail"
	StyleIllustration  ImageStyle = "illustration"
)

// ImageGenerationRequest represents a request to generate images
type ImageGenerationRequest struct {
	Provider ImageProvider `json:"provider"`
	Prompt   string        `json:"prompt"`
	Size     string        `json:"size,omitempty"`
	Style    ImageStyle    `json:"style,omitempty"`
	Context  string        `json:"context,omitempty"`
	N        int           `json:"n,omitempty"`
}

// ImageData represents generated image data
type ImageData struct {
	URL      string `json:"url,omitempty"`
	B64JSON  string `json:"b64_json,omitempty"`
}

// ImageGenerationResponse represents the response from image generation
type ImageGenerationResponse struct {
	Data     []ImageData `json:"data"`
	Provider string      `json:"provider"`
	Prompt   string      `json:"prompt"`
}

// DiagramType represents supported diagram types
type DiagramType string

const (
	DiagramTypeMermaid  DiagramType = "mermaid"
	DiagramTypeGraphviz DiagramType = "graphviz"
)

// DiagramGenerationRequest represents a request to generate diagrams
type DiagramGenerationRequest struct {
	Description  string      `json:"description"`
	DiagramType  DiagramType `json:"diagram_type"`
	DiagramKind  string      `json:"diagram_kind,omitempty"`
	OutputFormat string      `json:"output_format,omitempty"`
	AIProvider   string      `json:"ai_provider,omitempty"`
}

// DiagramGenerationResponse represents the response from diagram generation
type DiagramGenerationResponse struct {
	B64JSON     string `json:"b64_json"`
	Code        string `json:"code"`
	Format      string `json:"format"`
	DiagramType string `json:"diagram_type"`
}

// VideoGenerationRequest represents a request to generate videos
type VideoGenerationRequest struct {
	ImageURL       string `json:"image_url,omitempty"`
	ImageB64       string `json:"image_b64,omitempty"`
	MotionBucketID int    `json:"motion_bucket_id,omitempty"`
	NumFrames      int    `json:"num_frames,omitempty"`
	Provider       string `json:"provider,omitempty"`
}

// VideoGenerationResponse represents the response from video generation
type VideoGenerationResponse struct {
	VideoURL  string `json:"video_url,omitempty"`
	VideoB64  string `json:"video_b64,omitempty"`
	Format    string `json:"format"`
}

// GenerateImage generates an image using the creative service
func (c *Client) GenerateImage(ctx context.Context, req ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// Validate request
	if err := ValidateImageGenerationRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Wrap the API call with circuit breaker
	result, err := appcircuitbreaker.Execute(c.cb, func() (*ImageGenerationResponse, error) {
		appcircuitbreaker.RecordRequestAllowed("creative-service")
		return c.doGenerateImage(ctx, req)
	})
	
	if err != nil {
		// Check if error is due to circuit breaker being open
		if err == gobreaker.ErrOpenState {
			appcircuitbreaker.RecordRequestRejected("creative-service")
			return nil, fmt.Errorf("creative-service circuit breaker is open: service unavailable")
		}
		return nil, err
	}
	
	return result, nil
}

func (c *Client) doGenerateImage(ctx context.Context, req ImageGenerationRequest) (*ImageGenerationResponse, error) {
	url := fmt.Sprintf("%s/v1/images/generate", c.baseURL)
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call creative service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("creative service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result ImageGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GenerateThumbnail generates a thumbnail optimized for social media
func (c *Client) GenerateThumbnail(ctx context.Context, prompt, context string) (*ImageGenerationResponse, error) {
	// Validate inputs
	if err := validation.ValidateString(prompt, 1, 1000, "prompt"); err != nil {
		return nil, fmt.Errorf("prompt: %w", err)
	}
	if context != "" {
		if err := validation.ValidateString(context, 1, 2000, "context"); err != nil {
			return nil, fmt.Errorf("context: %w", err)
		}
	}

	req := ImageGenerationRequest{
		Provider: ProviderOpenAI,
		Prompt:   prompt,
		Style:    StyleThumbnail,
		Context:  context,
		N:        1,
	}
	
	// Wrap the API call with circuit breaker
	result, err := appcircuitbreaker.Execute(c.cb, func() (*ImageGenerationResponse, error) {
		appcircuitbreaker.RecordRequestAllowed("creative-service")
		return c.doGenerateThumbnail(ctx, req)
	})
	
	if err != nil {
		// Check if error is due to circuit breaker being open
		if err == gobreaker.ErrOpenState {
			appcircuitbreaker.RecordRequestRejected("creative-service")
			return nil, fmt.Errorf("creative-service circuit breaker is open: service unavailable")
		}
		return nil, err
	}
	
	return result, nil
}

func (c *Client) doGenerateThumbnail(ctx context.Context, req ImageGenerationRequest) (*ImageGenerationResponse, error) {
	url := fmt.Sprintf("%s/v1/images/generate/thumbnail", c.baseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call creative service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("creative service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result ImageGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GenerateDiagram generates a technical diagram
func (c *Client) GenerateDiagram(ctx context.Context, req DiagramGenerationRequest) (*DiagramGenerationResponse, error) {
	// Validate request
	if err := ValidateDiagramGenerationRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Wrap the API call with circuit breaker
	result, err := appcircuitbreaker.Execute(c.cb, func() (*DiagramGenerationResponse, error) {
		appcircuitbreaker.RecordRequestAllowed("creative-service")
		return c.doGenerateDiagram(ctx, req)
	})
	
	if err != nil {
		// Check if error is due to circuit breaker being open
		if err == gobreaker.ErrOpenState {
			appcircuitbreaker.RecordRequestRejected("creative-service")
			return nil, fmt.Errorf("creative-service circuit breaker is open: service unavailable")
		}
		return nil, err
	}
	
	return result, nil
}

func (c *Client) doGenerateDiagram(ctx context.Context, req DiagramGenerationRequest) (*DiagramGenerationResponse, error) {
	url := fmt.Sprintf("%s/v1/diagrams/generate", c.baseURL)
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call creative service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("creative service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result DiagramGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate response
	if err := ValidateDiagramGenerationResponse(&result); err != nil {
		return nil, fmt.Errorf("invalid response: %w", err)
	}

	return &result, nil
}

// GenerateVideo generates a video from an image
func (c *Client) GenerateVideo(ctx context.Context, req VideoGenerationRequest) (*VideoGenerationResponse, error) {
	// Validate request
	if err := ValidateVideoGenerationRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Wrap the API call with circuit breaker
	result, err := appcircuitbreaker.Execute(c.cb, func() (*VideoGenerationResponse, error) {
		appcircuitbreaker.RecordRequestAllowed("creative-service")
		return c.doGenerateVideo(ctx, req)
	})
	
	if err != nil {
		// Check if error is due to circuit breaker being open
		if err == gobreaker.ErrOpenState {
			appcircuitbreaker.RecordRequestRejected("creative-service")
			return nil, fmt.Errorf("creative-service circuit breaker is open: service unavailable")
		}
		return nil, err
	}
	
	return result, nil
}

func (c *Client) doGenerateVideo(ctx context.Context, req VideoGenerationRequest) (*VideoGenerationResponse, error) {
	url := fmt.Sprintf("%s/v1/videos/generate", c.baseURL)
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call creative service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("creative service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result VideoGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validate response
	if err := ValidateVideoGenerationResponse(&result); err != nil {
		return nil, fmt.Errorf("invalid response: %w", err)
	}

	return &result, nil
}

