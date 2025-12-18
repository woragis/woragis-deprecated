# AI Service Component

## Overview

A Python/FastAPI microservice that provides AI/LLM integration for the Woragis platform. It supports multiple LLM providers (OpenAI, Anthropic, xAI, Manus, Cipher) and provides chat completion, streaming, and image generation capabilities.

## Architecture

- **Language**: Python 3.11+
- **Framework**: FastAPI
- **Port**: 8000
- **LLM Providers**: OpenAI, Anthropic, xAI, Manus, Cipher
- **Image Provider**: Cipher (for images)

## Responsibilities

1. **Chat Completion**: Provide chat completion with various agent personas
2. **Streaming**: Support streaming responses for real-time chat
3. **Image Generation**: Generate images using Cipher API
4. **Agent Management**: Manage different agent personas (economist, strategist, entrepreneur, startup, auto)

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
- LLM API call metrics
- Error rates

## Configuration

### Environment Variables

#### Required
- `OPENAI_API_KEY` - OpenAI API key (if using OpenAI)
- `ANTHROPIC_API_KEY` - Anthropic API key (if using Anthropic)
- `CIPHER_API_KEY` - Cipher API key (if using Cipher)

#### Optional
- `OPENAI_MODEL` - OpenAI model (default: `gpt-4o-mini`)
- `OPENAI_TEMPERATURE` - Temperature (default: `0.3`)
- `CIPHER_MAX_TOKENS` - Max tokens for Cipher (default: `4000`)
- `CIPHER_TOP_P` - Top-p for Cipher (default: `0.9`)
- `CIPHER_IMAGE_N` - Number of images (default: `1`)
- `CIPHER_IMAGE_SIZE` - Image size (default: `1024x1024`)
- `CORS_ENABLED` - Enable CORS (default: `true`)
- `CORS_ALLOWED_ORIGINS` - Comma-separated origins
- `ENV` - Environment (development/production)

## API Endpoints

### List Agents

**GET** `/v1/agents`

Returns list of available agent personas.

**Response**:
```json
["economist", "strategist", "entrepreneur", "startup", "auto"]
```

### Chat Completion

**POST** `/v1/chat`

Chat completion with agent persona.

**Request**:
```json
{
  "agent": "economist",
  "input": "What are the economic implications of microservices?",
  "system": "Optional additional system instruction",
  "temperature": 0.7,
  "model": "gpt-4o-mini",
  "provider": "openai"
}
```

**Response**:
```json
{
  "agent": "economist",
  "output": "Microservices architecture has several economic implications..."
}
```

### Streaming Chat

**POST** `/v1/chat/stream`

Streaming chat completion.

**Request**: Same as `/v1/chat`

**Response**: NDJSON stream
```json
{"delta": "Microservices"}
{"delta": " architecture"}
{"done": true, "output": "Full output text"}
```

### Image Generation

**POST** `/v1/images`

Generate images using Cipher API.

**Request**:
```json
{
  "provider": "cipher",
  "prompt": "A futuristic microservices architecture diagram",
  "n": 1,
  "size": "1024x1024"
}
```

**Response**:
```json
{
  "data": [
    {
      "url": "https://...",
      "b64_json": "base64_encoded_image"
    }
  ]
}
```

## Agent Personas

### Economist
- **Focus**: Economic analysis, market trends, pricing, unit economics
- **Keywords**: market, inflation, macro, econom, unit economics, pricing

### Strategist
- **Focus**: Business strategy, positioning, go-to-market, competitive analysis
- **Keywords**: strategy, positioning, go-to-market, gtm, competitor, moat

### Entrepreneur
- **Focus**: MVP development, launching, prototyping, validation
- **Keywords**: mvp, launch, prototype, hack, validate, scrappy

### Startup
- **Focus**: General startup advice and guidance
- **Default**: Used when no specific persona matches

### Auto
- **Focus**: Automatically selects persona based on input keywords
- **Logic**: Analyzes input to determine best persona

## LLM Providers

### OpenAI
- **Models**: gpt-4o-mini, gpt-4, gpt-3.5-turbo
- **Features**: Chat completion, streaming
- **Rate Limits**: Per API key tier

### Anthropic
- **Models**: claude-3-opus, claude-3-sonnet, claude-3-haiku
- **Features**: Chat completion, streaming
- **Rate Limits**: Per API key tier

### xAI
- **Models**: grok (if available)
- **Features**: Chat completion
- **Rate Limits**: Per API key tier

### Manus
- **Models**: Various (if available)
- **Features**: Chat completion
- **Rate Limits**: Per API key tier

### Cipher
- **Models**: Cipher-specific models
- **Features**: Chat completion, image generation
- **Rate Limits**: Per API key tier

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `ai-service`

**Key Log Fields**:
- `agent` - Agent persona
- `provider` - LLM provider
- `model` - Model used
- `input_length` - Input text length
- `output_length` - Output text length
- `duration_ms` - Processing duration

## Deployment

### Local Development

```bash
cd backend/ai-service
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

### Docker

```bash
docker build -t woragis/ai-service .
docker run -p 8000:8000 --env-file .env woragis/ai-service
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

LLM providers have rate limits:
- **OpenAI**: Per API key tier
- **Anthropic**: Per API key tier
- **Cipher**: Per API key tier

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
- LLM API call duration
- Rate limit usage
- Token usage (if tracked)

### Alerts
- Error rate > 5%
- Latency p95 > 5 seconds
- Rate limit approaching
- Service unavailable

## Troubleshooting

### Common Issues

#### High Latency
- Check LLM provider status
- Check network connectivity
- Verify API keys are valid
- Check rate limits

#### High Error Rate
- Check LLM provider status
- Verify API keys are valid
- Check rate limits
- Review error logs for specific errors

#### Rate Limit Errors
- Check rate limit usage
- Implement exponential backoff
- Consider using multiple API keys
- Monitor rate limit metrics

#### Service Unavailable
- Check service health (`/healthz`)
- Check logs for errors
- Verify environment variables
- Check resource usage (CPU, memory)

## Related Documentation

- [API Documentation](../api/ai-service-api.md) - Detailed API documentation
- [Architecture Decision Records](../adr/) - Architectural decisions
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
