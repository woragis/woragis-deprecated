# Creative Service Component

## Overview

A Python/FastAPI microservice that provides AI-powered image, diagram, and video generation for technical blog posts. It supports multiple providers for creating thumbnails, architecture diagrams, workflow illustrations, and animated content.

## Architecture

- **Language**: Python 3.11+
- **Framework**: FastAPI
- **Port**: 8000
- **Image Providers**: OpenAI (DALL-E 3), Stable Diffusion XL, Cipher
- **Diagram Providers**: Mermaid, Graphviz
- **Video Providers**: Replicate (Stable Video Diffusion), Runway

## Responsibilities

1. **Image Generation**: Generate images using various AI providers
2. **Diagram Generation**: Generate technical diagrams (Mermaid, Graphviz) from descriptions
3. **Video/GIF Generation**: Animate static images into videos
4. **Thumbnail Generation**: Optimized endpoint for social media thumbnails

## Health Check

**Endpoint**: `GET /healthz`

**Checks**:
- Service availability (no critical dependencies)

**Response**:
```json
{
  "status": "healthy",
  "checks": []
}
```

## Metrics

**Endpoint**: `GET /metrics`

Exposes Prometheus metrics:
- HTTP request rate and latency
- Image generation metrics
- Diagram generation metrics
- Video generation metrics
- Error rates

## Configuration

### Environment Variables

#### Required
- `OPENAI_API_KEY` - OpenAI API key (for DALL-E 3)
- `STABLE_DIFFUSION_API_KEY` - Stable Diffusion API key (if using)
- `CIPHER_API_KEY` - Cipher API key (if using)
- `REPLICATE_API_TOKEN` - Replicate API token (for video generation)

#### Optional
- `CORS_ENABLED` - Enable CORS (default: `true`)
- `CORS_ALLOWED_ORIGINS` - Comma-separated origins
- `ENV` - Environment (development/production)

## API Endpoints

### Generate Images

**POST** `/v1/images/generate`

Generate images using various AI providers.

**Request**:
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
  "prompt": "..."
}
```

### Generate Thumbnails

**POST** `/v1/images/generate/thumbnail`

Optimized endpoint for social media thumbnails.

**Request**: Same as `/v1/images/generate` (style automatically set to "thumbnail")

### Generate Diagrams

**POST** `/v1/diagrams/generate`

Generate technical diagrams from descriptions.

**Request**:
```json
{
  "description": "A microservices architecture with an API gateway, three microservices (user, order, payment), and a Redis cache layer",
  "diagram_type": "mermaid",
  "diagram_kind": "flowchart",
  "output_format": "png",
  "ai_provider": "openai"
}
```

**Response**:
```json
{
  "b64_json": "base64_encoded_diagram_image",
  "code": "graph TD\n    A[API Gateway] --> B[User Service]\n    ...",
  "format": "png",
  "diagram_type": "mermaid"
}
```

### Generate Mermaid Diagrams

**POST** `/v1/diagrams/mermaid`

Quick endpoint for Mermaid diagrams.

**Request**:
```json
{
  "description": "Request flow through microservices",
  "diagram_kind": "sequence",
  "output_format": "png"
}
```

### Generate Videos/GIFs

**POST** `/v1/videos/generate`

Animate a static image into a video.

**Request**:
```json
{
  "image_b64": "base64_encoded_image",
  "motion_bucket_id": 127,
  "num_frames": 25,
  "provider": "replicate"
}
```

**Response**:
```json
{
  "video_url": "https://...",
  "video_b64": "base64_encoded_video",
  "format": "mp4"
}
```

### Animate Image

**POST** `/v1/videos/animate`

Create an animated GIF/video from a static image.

**Request**: Same as `/v1/videos/generate`

### List Providers

**GET** `/v1/providers/images`
**GET** `/v1/providers/diagrams`
**GET** `/v1/providers/videos`

List available providers for each type.

## Image Providers

### OpenAI (DALL-E 3)
- **Best For**: High-quality thumbnails and illustrations
- **Quality**: Excellent
- **Cost**: Pay-per-image
- **Rate Limits**: Per API key tier

### Stable Diffusion XL
- **Best For**: Artistic and diagram-style images
- **Quality**: Very good
- **Cost**: Varies by provider
- **Rate Limits**: Per provider

### Cipher
- **Best For**: Alternative provider
- **Quality**: Good
- **Cost**: Varies
- **Rate Limits**: Per API key tier

## Diagram Types

### Mermaid
- **Types**: flowchart, sequence, er, gantt, pie, etc.
- **Output**: PNG, SVG
- **Best For**: Technical documentation, architecture diagrams

### Graphviz
- **Types**: Network graphs, dependency graphs
- **Output**: PNG, SVG, PDF
- **Best For**: Complex network visualizations

## Video Providers

### Replicate (Stable Video Diffusion)
- **Best For**: Animating static images
- **Quality**: Good
- **Cost**: Pay-per-generation
- **Rate Limits**: Per API token

### Runway
- **Best For**: Advanced video generation (coming soon)
- **Quality**: Excellent
- **Cost**: Varies
- **Rate Limits**: Per API key tier

## Use Cases

### 1. Thumbnails for Blog Posts
Generate eye-catching thumbnails optimized for social media platforms.

**Recommended**: Use `openai` provider with `style: "thumbnail"`

### 2. Architecture Diagrams
Create clear, professional architecture diagrams from descriptions.

**Recommended**: Use diagram generation with `mermaid` and `diagram_kind: "flowchart"`

### 3. Workflow Illustrations
Visualize complex workflows and processes.

**Recommended**: Use diagram generation with `mermaid` and `diagram_kind: "sequence"`

### 4. Animated Content
Transform static diagrams into engaging animations.

**Recommended**: Generate image/diagram first, then use video generation to animate it

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `creative-service`

**Key Log Fields**:
- `provider` - Provider used
- `type` - Generation type (image, diagram, video)
- `style` - Image style (if applicable)
- `status` - Success/failed
- `error` - Error message (if failed)

## Deployment

### Local Development

```bash
cd backend/creative-service
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

### Docker

```bash
docker build -t woragis/creative-service .
docker run -p 8000:8000 --env-file .env woragis/creative-service
```

### Kubernetes

Deploy as a Deployment:
- Health check probe on `/healthz`
- No critical dependencies (stateless)
- Can scale horizontally

## Scaling

### Horizontal Scaling
- Stateless design allows multiple replicas
- Load balancer distributes requests
- Each replica has its own API key connections

### Resource Requirements
- **CPU**: 300m-1000m (0.3-1 core)
- **Memory**: 512Mi-1Gi
- **External API Connections**: 10-50 per replica

## Rate Limiting

Providers have rate limits:
- **OpenAI**: Per API key tier
- **Stable Diffusion**: Per provider
- **Replicate**: Per API token

**Handling**:
- Respect rate limits
- Implement exponential backoff on rate limit errors
- Monitor rate limit usage
- Consider using multiple API keys for scaling

## Monitoring

### Key Metrics
- Request rate (requests/second)
- Latency (p50, p95, p99)
- Error rate
- Generation success rate (by type)
- Provider-specific metrics

### Alerts
- Error rate > 5%
- Latency p95 > 10 seconds (generation can be slow)
- Rate limit approaching
- Service unavailable

## Troubleshooting

### Common Issues

#### High Latency
- Check provider status
- Check network connectivity
- Verify API keys are valid
- Check rate limits
- **Note**: Image/diagram/video generation can be inherently slow (10-30 seconds)

#### High Error Rate
- Check provider status
- Verify API keys are valid
- Check rate limits
- Review error logs for specific errors
- Verify input format is correct

#### Rate Limit Errors
- Check rate limit usage
- Implement exponential backoff
- Consider using multiple API keys
- Monitor rate limit metrics

#### Generation Failures
- Check provider status
- Verify prompts are valid
- Check input image format (for video generation)
- Review provider-specific error messages

## Related Documentation

- [API Documentation](../api/creative-service-api.md) - Detailed API documentation
- [Architecture Decision Records](../adr/) - Architectural decisions
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
