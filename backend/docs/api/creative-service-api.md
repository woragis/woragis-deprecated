# Creative Service API Documentation

## Overview

The Creative Service provides AI-powered image, diagram, and video generation for technical blog posts. It supports multiple providers for creating thumbnails, architecture diagrams, workflow illustrations, and animated content.

## Base URL

- **Development**: `http://localhost:8000`
- **Production**: `http://creative-service:8000` (internal) or `https://creative.woragis.com` (if exposed)

## Authentication

Currently, the Creative Service does not require authentication (internal service). In production, consider adding API key authentication.

## API Endpoints

### Generate Images

**POST** `/v1/images/generate`

Generate images using various AI providers.

**Request Body**:
```json
{
  "provider": "openai",
  "prompt": "Microservices architecture with connecting gears and boxes",
  "style": "technical",
  "context": "Software architecture for a distributed system",
  "size": "1024x1024",
  "n": 1
}
```

**Request Fields**:
- `provider` (required): Image provider - `"openai"`, `"stable-diffusion"`, or `"cipher"` (default: `"openai"`)
- `prompt` (required): Image generation prompt
- `style` (optional): Image style - `"technical"`, `"diagram"`, `"thumbnail"`, or `"illustration"` (default: `"technical"`)
- `context` (optional): Additional context for the image
- `size` (optional): Image size, e.g., `"1024x1024"` (default: provider-specific)
- `n` (optional): Number of images to generate (default: 1)

**Response**:
```json
{
  "data": [
    {
      "b64_json": "base64_encoded_image",
      "url": null
    }
  ],
  "provider": "openai",
  "prompt": "Microservices architecture with connecting gears and boxes"
}
```

**Response Fields**:
- `data`: Array of image objects
  - `b64_json`: Base64-encoded image
  - `url`: Image URL (if available)
- `provider`: Provider used
- `prompt`: Prompt used

**Status Codes**:
- `200 OK` - Success
- `400 Bad Request` - Invalid request
- `500 Internal Server Error` - Image generation error

**Example**:
```bash
curl -X POST http://localhost:8000/v1/images/generate \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "prompt": "Microservices architecture with connecting gears and boxes",
    "style": "technical",
    "size": "1024x1024"
  }'
```

---

### Generate Thumbnails

**POST** `/v1/images/generate/thumbnail`

Optimized endpoint for social media thumbnails. Automatically sets `style` to `"thumbnail"`.

**Request Body**: Same as `/v1/images/generate`

**Response**: Same as `/v1/images/generate`

**Example**:
```bash
curl -X POST http://localhost:8000/v1/images/generate/thumbnail \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "prompt": "Kubernetes orchestration architecture",
    "context": "Technical blog post thumbnail"
  }'
```

---

### Generate Diagrams

**POST** `/v1/diagrams/generate`

Generate technical diagrams from descriptions using AI-generated code (Mermaid/Graphviz).

**Request Body**:
```json
{
  "description": "A microservices architecture with an API gateway, three microservices (user, order, payment), and a Redis cache layer",
  "diagram_type": "mermaid",
  "diagram_kind": "flowchart",
  "output_format": "png",
  "ai_provider": "openai"
}
```

**Request Fields**:
- `description` (required): Description of the diagram to generate
- `diagram_type` (required): Type of diagram code - `"mermaid"` or `"graphviz"` (default: `"mermaid"`)
- `diagram_kind` (optional): Kind of diagram - `"flowchart"`, `"sequence"`, `"er"`, etc. (default: `"flowchart"`)
- `output_format` (required): Output format - `"png"`, `"svg"`, or `"pdf"` (default: `"png"`)
- `ai_provider` (required): AI provider for code generation - `"openai"` or `"anthropic"` (default: `"openai"`)

**Response**:
```json
{
  "b64_json": "base64_encoded_diagram_image",
  "code": "graph TD\n    A[API Gateway] --> B[User Service]\n    A --> C[Order Service]\n    A --> D[Payment Service]\n    B --> E[Redis Cache]",
  "format": "png",
  "diagram_type": "mermaid"
}
```

**Response Fields**:
- `b64_json`: Base64-encoded diagram image
- `code`: Generated diagram code (Mermaid or Graphviz)
- `format`: Output format used
- `diagram_type`: Diagram type used

**Status Codes**:
- `200 OK` - Success
- `400 Bad Request` - Invalid request
- `500 Internal Server Error` - Diagram generation error

**Example**:
```bash
curl -X POST http://localhost:8000/v1/diagrams/generate \
  -H "Content-Type: application/json" \
  -d '{
    "description": "A microservices architecture with an API gateway, three microservices, and a Redis cache",
    "diagram_type": "mermaid",
    "diagram_kind": "flowchart",
    "output_format": "png",
    "ai_provider": "openai"
  }'
```

---

### Generate Mermaid Diagrams

**POST** `/v1/diagrams/mermaid`

Quick endpoint for Mermaid diagrams.

**Request Parameters** (query or JSON body):
- `description` (required): Description of the diagram
- `diagram_kind` (optional): Kind of diagram (default: `"flowchart"`)
- `output_format` (optional): Output format - `"png"` or `"svg"` (default: `"png"`)
- `ai_provider` (optional): AI provider (default: `"openai"`)

**Response**: Same as `/v1/diagrams/generate`

**Example**:
```bash
curl -X POST "http://localhost:8000/v1/diagrams/mermaid?description=Request%20flow%20through%20microservices&diagram_kind=sequence&output_format=png"
```

---

### Generate Videos/GIFs

**POST** `/v1/videos/generate`

Animate a static image into a video/GIF.

**Request Body**:
```json
{
  "image_url": "https://example.com/image.jpg",
  "image_b64": null,
  "motion_bucket_id": 127,
  "num_frames": 25,
  "provider": "replicate"
}
```

**Request Fields**:
- `image_url` (optional): URL of the input image
- `image_b64` (optional): Base64-encoded input image
- `motion_bucket_id` (optional): Motion intensity (1-255, default: 127)
- `num_frames` (optional): Number of frames (14-25, default: 25)
- `provider` (required): Video provider - `"replicate"` or `"runway"` (default: `"replicate"`)

**Note**: Either `image_url` or `image_b64` must be provided.

**Response**:
```json
{
  "video_url": "https://...",
  "video_b64": "base64_encoded_video",
  "format": "mp4"
}
```

**Response Fields**:
- `video_url`: Video URL (if available)
- `video_b64`: Base64-encoded video (if available)
- `format`: Video format (`"mp4"` or `"gif"`)

**Status Codes**:
- `200 OK` - Success
- `400 Bad Request` - Invalid request (missing image, invalid parameters)
- `500 Internal Server Error` - Video generation error

**Example**:
```bash
curl -X POST http://localhost:8000/v1/videos/generate \
  -H "Content-Type: application/json" \
  -d '{
    "image_b64": "base64_encoded_image",
    "motion_bucket_id": 127,
    "num_frames": 25,
    "provider": "replicate"
  }'
```

---

### Animate Image

**POST** `/v1/videos/animate`

Create an animated GIF/video from a static image. Alias for `/v1/videos/generate`.

**Request Body**: Same as `/v1/videos/generate`

**Response**: Same as `/v1/videos/generate`

---

### List Providers

**GET** `/v1/providers/images`

List available image generation providers.

**Response**:
```json
{
  "providers": ["openai", "stable-diffusion", "cipher"],
  "descriptions": {
    "openai": "DALL-E 3 - Best for high-quality thumbnails and illustrations",
    "stable-diffusion": "Stable Diffusion XL - Good for artistic and diagram-style images",
    "cipher": "Cipher/NoFilterGPT - Alternative provider"
  }
}
```

---

**GET** `/v1/providers/diagrams`

List available diagram generation providers.

**Response**:
```json
{
  "providers": ["mermaid", "graphviz"],
  "ai_providers": ["openai", "anthropic"],
  "descriptions": {
    "mermaid": "Mermaid diagrams - Flowcharts, sequence diagrams, ER diagrams, etc.",
    "graphviz": "Graphviz DOT - Network graphs, dependency graphs, etc."
  }
}
```

---

**GET** `/v1/providers/videos`

List available video generation providers.

**Response**:
```json
{
  "providers": ["replicate", "runway"],
  "descriptions": {
    "replicate": "Stable Video Diffusion - Generate videos from images",
    "runway": "Runway Gen-2 - Advanced video generation (coming soon)"
  }
}
```

---

### Health Check

**GET** `/healthz`

Check service health.

**Response**:
```json
{
  "status": "healthy",
  "checks": []
}
```

**Status Codes**:
- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service is unhealthy

---

### Metrics

**GET** `/metrics`

Prometheus metrics endpoint (internal use only).

**Content-Type**: `text/plain`

---

## Use Cases

### 1. Thumbnails for Blog Posts

Generate eye-catching thumbnails optimized for social media platforms.

**Recommended**: Use `openai` provider with `style: "thumbnail"`

```bash
curl -X POST http://localhost:8000/v1/images/generate/thumbnail \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "prompt": "Kubernetes orchestration architecture",
    "style": "thumbnail"
  }'
```

### 2. Architecture Diagrams

Create clear, professional architecture diagrams from descriptions.

**Recommended**: Use diagram generation with `mermaid` and `diagram_kind: "flowchart"`

```bash
curl -X POST http://localhost:8000/v1/diagrams/generate \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Three-tier architecture: web server, application server, database",
    "diagram_type": "mermaid",
    "diagram_kind": "flowchart",
    "output_format": "png"
  }'
```

### 3. Workflow Illustrations

Visualize complex workflows and processes.

**Recommended**: Use diagram generation with `mermaid` and `diagram_kind: "sequence"`

```bash
curl -X POST http://localhost:8000/v1/diagrams/mermaid \
  -H "Content-Type: application/json" \
  -d '{
    "description": "User authentication flow with OAuth, JWT tokens, and session management",
    "diagram_kind": "sequence",
    "output_format": "png"
  }'
```

### 4. Animated Content

Transform static diagrams into engaging animations.

**Recommended**: Generate image/diagram first, then use video generation to animate it

```bash
# First generate an image/diagram
IMAGE_RESPONSE=$(curl -X POST http://localhost:8000/v1/images/generate ...)

# Then animate it
curl -X POST http://localhost:8000/v1/videos/animate \
  -H "Content-Type: application/json" \
  -d "{
    \"image_b64\": \"$(echo $IMAGE_RESPONSE | jq -r '.data[0].b64_json')\",
    \"motion_bucket_id\": 100,
    \"provider\": \"replicate\"
  }"
```

---

## Error Responses

### Error Format

```json
{
  "detail": "Error message"
}
```

### Common Errors

**400 Bad Request**:
- Missing required fields
- Invalid provider
- Invalid format
- Missing image (for video generation)

**500 Internal Server Error**:
- Provider API error
- Generation timeout
- Service error

---

## Rate Limiting

Rate limiting is handled by providers. The service respects provider rate limits and may return errors if limits are exceeded.

**Note**: Image/diagram/video generation can be slow (10-30 seconds). Consider implementing request timeouts and async processing.

---

## Examples

### Complete Flow: Generate Thumbnail and Diagram

```bash
# 1. Generate thumbnail
curl -X POST http://localhost:8000/v1/images/generate/thumbnail \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "prompt": "Microservices architecture",
    "context": "Technical blog post"
  }'

# 2. Generate architecture diagram
curl -X POST http://localhost:8000/v1/diagrams/generate \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Microservices architecture with API gateway, services, and database",
    "diagram_type": "mermaid",
    "diagram_kind": "flowchart",
    "output_format": "png"
  }'
```

---

## Related Documentation

- [Component Documentation](../components/creative-service.md) - Creative Service component details
- [Architecture Decision Records](../adr/) - Architectural decisions
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
