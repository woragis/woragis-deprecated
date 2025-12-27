package creative

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateImageGenerationRequest validates image generation request
func ValidateImageGenerationRequest(req ImageGenerationRequest) error {
	// Validate provider (required)
	if req.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	validProviders := []ImageProvider{
		ProviderOpenAI,
		ProviderStableDiffusion,
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
		return fmt.Errorf("provider: must be one of: openai, stable-diffusion, cipher")
	}

	// Validate prompt (required, 1-1000 chars)
	if err := validation.ValidateString(req.Prompt, 1, 1000, "prompt"); err != nil {
		return fmt.Errorf("prompt: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(req.Prompt); err != nil {
		return fmt.Errorf("prompt: %w", err)
	}
	if err := validation.ValidateNoXSS(req.Prompt); err != nil {
		return fmt.Errorf("prompt: %w", err)
	}

	// Validate size (optional, but if provided, validate)
	if req.Size != "" {
		validSizes := []string{"256x256", "512x512", "1024x1024", "1792x1024", "1024x1792"}
		isValid := false
		for _, validSize := range validSizes {
			if req.Size == validSize {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("size: must be one of: 256x256, 512x512, 1024x1024, 1792x1024, 1024x1792")
		}
	}

	// Validate style (optional, but if provided, validate)
	if req.Style != "" {
		validStyles := []ImageStyle{
			StyleTechnical,
			StyleDiagram,
			StyleThumbnail,
			StyleIllustration,
		}
		isValid := false
		for _, validStyle := range validStyles {
			if req.Style == validStyle {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("style: must be one of: technical, diagram, thumbnail, illustration")
		}
	}

	// Validate context (optional, but if provided, validate)
	if req.Context != "" {
		if err := validation.ValidateString(req.Context, 1, 2000, "context"); err != nil {
			return fmt.Errorf("context: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(req.Context); err != nil {
			return fmt.Errorf("context: %w", err)
		}
		if err := validation.ValidateNoXSS(req.Context); err != nil {
			return fmt.Errorf("context: %w", err)
		}
	}

	// Validate N (optional, but if provided, validate range)
	if req.N > 0 {
		if req.N < 1 {
			return fmt.Errorf("n: must be at least 1")
		}
		if req.N > 10 {
			return fmt.Errorf("n: must be at most 10")
		}
	}

	return nil
}

// ValidateDiagramGenerationRequest validates diagram generation request
func ValidateDiagramGenerationRequest(req DiagramGenerationRequest) error {
	// Validate description (required, 1-2000 chars)
	if err := validation.ValidateString(req.Description, 1, 2000, "description"); err != nil {
		return fmt.Errorf("description: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(req.Description); err != nil {
		return fmt.Errorf("description: %w", err)
	}
	if err := validation.ValidateNoXSS(req.Description); err != nil {
		return fmt.Errorf("description: %w", err)
	}

	// Validate diagram type (required)
	if req.DiagramType == "" {
		return fmt.Errorf("diagram_type is required")
	}
	validTypes := []DiagramType{
		DiagramTypeMermaid,
		DiagramTypeGraphviz,
	}
	isValid := false
	for _, validType := range validTypes {
		if req.DiagramType == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("diagram_type: must be one of: mermaid, graphviz")
	}

	// Validate diagram kind (optional, but if provided, validate)
	if req.DiagramKind != "" {
		if err := validation.ValidateString(req.DiagramKind, 1, 100, "diagram_kind"); err != nil {
			return fmt.Errorf("diagram_kind: %w", err)
		}
	}

	// Validate output format (optional, but if provided, validate)
	if req.OutputFormat != "" {
		validFormats := []string{"png", "svg", "pdf", "json"}
		isValid := false
		for _, validFormat := range validFormats {
			if req.OutputFormat == validFormat {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("output_format: must be one of: png, svg, pdf, json")
		}
	}

	// Validate AI provider (optional, but if provided, validate)
	if req.AIProvider != "" {
		if err := validation.ValidateString(req.AIProvider, 1, 100, "ai_provider"); err != nil {
			return fmt.Errorf("ai_provider: %w", err)
		}
	}

	return nil
}

// ValidateVideoGenerationRequest validates video generation request
func ValidateVideoGenerationRequest(req VideoGenerationRequest) error {
	// At least one image source must be provided
	if req.ImageURL == "" && req.ImageB64 == "" {
		return fmt.Errorf("at least one of image_url or image_b64 is required")
	}

	// Validate image URL (optional, but if provided, validate URL)
	if req.ImageURL != "" {
		if err := validation.ValidateURL(req.ImageURL); err != nil {
			return fmt.Errorf("image_url: %w", err)
		}
	}

	// Validate image B64 (optional, but if provided, validate length)
	if req.ImageB64 != "" {
		if err := validation.ValidateString(req.ImageB64, 1, 10000000, "image_b64"); err != nil {
			return fmt.Errorf("image_b64: %w", err)
		}
	}

	// Validate motion bucket ID (optional, but if provided, validate range)
	if req.MotionBucketID > 0 {
		if req.MotionBucketID < 1 || req.MotionBucketID > 255 {
			return fmt.Errorf("motion_bucket_id: must be between 1 and 255")
		}
	}

	// Validate num frames (optional, but if provided, validate range)
	if req.NumFrames > 0 {
		if req.NumFrames < 1 {
			return fmt.Errorf("num_frames: must be at least 1")
		}
		if req.NumFrames > 1000 {
			return fmt.Errorf("num_frames: must be at most 1000")
		}
	}

	// Validate provider (optional, but if provided, validate)
	if req.Provider != "" {
		if err := validation.ValidateString(req.Provider, 1, 100, "provider"); err != nil {
			return fmt.Errorf("provider: %w", err)
		}
	}

	return nil
}

// ValidateImageGenerationResponse validates image generation response
func ValidateImageGenerationResponse(resp *ImageGenerationResponse) error {
	// Validate data (required, at least one image)
	if len(resp.Data) == 0 {
		return fmt.Errorf("data: at least one image is required")
	}
	if len(resp.Data) > 10 {
		return fmt.Errorf("data: too many images (maximum 10)")
	}

	// Validate each image
	for i, img := range resp.Data {
		// At least one of URL or B64JSON must be provided
		if img.URL == "" && img.B64JSON == "" {
			return fmt.Errorf("data[%d]: at least one of url or b64_json is required", i)
		}

		// Validate URL if provided
		if img.URL != "" {
			if err := validation.ValidateURL(img.URL); err != nil {
				return fmt.Errorf("data[%d].url: %w", i, err)
			}
		}

		// Validate B64JSON if provided
		if img.B64JSON != "" {
			if err := validation.ValidateString(img.B64JSON, 1, 10000000, fmt.Sprintf("data[%d].b64_json", i)); err != nil {
				return fmt.Errorf("data[%d].b64_json: %w", i, err)
			}
		}
	}

	// Validate provider (optional, but if provided, validate)
	if resp.Provider != "" {
		if err := validation.ValidateString(resp.Provider, 1, 100, "provider"); err != nil {
			return fmt.Errorf("provider: %w", err)
		}
	}

	// Validate prompt (optional, but if provided, validate)
	if resp.Prompt != "" {
		if err := validation.ValidateString(resp.Prompt, 1, 1000, "prompt"); err != nil {
			return fmt.Errorf("prompt: %w", err)
		}
	}

	return nil
}

// ValidateDiagramGenerationResponse validates diagram generation response
func ValidateDiagramGenerationResponse(resp *DiagramGenerationResponse) error {
	// Validate B64JSON (required)
	if resp.B64JSON == "" {
		return fmt.Errorf("b64_json is required")
	}
	if err := validation.ValidateString(resp.B64JSON, 1, 10000000, "b64_json"); err != nil {
		return fmt.Errorf("b64_json: %w", err)
	}

	// Validate code (optional, but if provided, validate)
	if resp.Code != "" {
		if err := validation.ValidateString(resp.Code, 1, 50000, "code"); err != nil {
			return fmt.Errorf("code: %w", err)
		}
	}

	// Validate format (optional, but if provided, validate)
	if resp.Format != "" {
		validFormats := []string{"png", "svg", "pdf", "json"}
		isValid := false
		for _, validFormat := range validFormats {
			if resp.Format == validFormat {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("format: must be one of: png, svg, pdf, json")
		}
	}

	// Validate diagram type (optional, but if provided, validate)
	if resp.DiagramType != "" {
		validTypes := []string{"mermaid", "graphviz"}
		isValid := false
		for _, validType := range validTypes {
			if resp.DiagramType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("diagram_type: must be one of: mermaid, graphviz")
		}
	}

	return nil
}

// ValidateVideoGenerationResponse validates video generation response
func ValidateVideoGenerationResponse(resp *VideoGenerationResponse) error {
	// At least one video source must be provided
	if resp.VideoURL == "" && resp.VideoB64 == "" {
		return fmt.Errorf("at least one of video_url or video_b64 is required")
	}

	// Validate video URL if provided
	if resp.VideoURL != "" {
		if err := validation.ValidateURL(resp.VideoURL); err != nil {
			return fmt.Errorf("video_url: %w", err)
		}
	}

	// Validate video B64 if provided
	if resp.VideoB64 != "" {
		if err := validation.ValidateString(resp.VideoB64, 1, 100000000, "video_b64"); err != nil {
			return fmt.Errorf("video_b64: %w", err)
		}
	}

	// Validate format (required)
	if resp.Format == "" {
		return fmt.Errorf("format is required")
	}
	validFormats := []string{"mp4", "webm", "gif"}
	isValid := false
	for _, validFormat := range validFormats {
		if resp.Format == validFormat {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("format: must be one of: mp4, webm, gif")
	}

	return nil
}

