package creative

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client interacts with the creative-service API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new creative service client
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = os.Getenv("CREATIVE_SERVICE_URL")
		if baseURL == "" {
			baseURL = "http://creative-service:8000"
		}
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for image generation
		},
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
	req := ImageGenerationRequest{
		Provider: ProviderOpenAI,
		Prompt:   prompt,
		Style:    StyleThumbnail,
		Context:  context,
		N:        1,
	}
	
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

	return &result, nil
}

// GenerateVideo generates a video from an image
func (c *Client) GenerateVideo(ctx context.Context, req VideoGenerationRequest) (*VideoGenerationResponse, error) {
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

	return &result, nil
}

