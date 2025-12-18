# AI Service API Documentation

## Overview

The AI Service provides AI/LLM integration for the Woragis platform. It supports multiple LLM providers (OpenAI, Anthropic, xAI, Manus, Cipher) and provides chat completion, streaming, and image generation capabilities.

## Base URL

- **Development**: `http://localhost:8000`
- **Production**: `http://ai-service:8000` (internal) or `https://ai.woragis.com` (if exposed)

## Authentication

Currently, the AI Service does not require authentication (internal service). In production, consider adding API key authentication.

## API Endpoints

### List Agents

**GET** `/v1/agents`

Returns list of available agent personas.

**Response**:
```json
["economist", "strategist", "entrepreneur", "startup", "auto"]
```

**Status Codes**:
- `200 OK` - Success

---

### Chat Completion

**POST** `/v1/chat`

Chat completion with agent persona.

**Request Body**:
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

**Request Fields**:
- `agent` (required): Agent persona - `"economist"`, `"strategist"`, `"entrepreneur"`, `"startup"`, or `"auto"`
- `input` (required): User input or question
- `system` (optional): Additional system instruction
- `temperature` (optional): Temperature override (0.0-2.0)
- `model` (optional): Model override (provider-specific)
- `provider` (optional): LLM provider - `"openai"`, `"anthropic"`, `"xai"`, `"manus"`, or `"cipher"` (default: `"openai"`)

**Response**:
```json
{
  "agent": "economist",
  "output": "Microservices architecture has several economic implications..."
}
```

**Response Fields**:
- `agent`: Agent persona used
- `output`: Generated response text

**Status Codes**:
- `200 OK` - Success
- `400 Bad Request` - Invalid request (invalid provider, model, etc.)
- `404 Not Found` - Unknown agent
- `500 Internal Server Error` - LLM API error

**Example**:
```bash
curl -X POST http://localhost:8000/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "agent": "economist",
    "input": "What are the economic implications of microservices?",
    "provider": "openai"
  }'
```

---

### Streaming Chat

**POST** `/v1/chat/stream`

Streaming chat completion (NDJSON format).

**Request Body**: Same as `/v1/chat`

**Response**: NDJSON stream
```json
{"delta": "Microservices"}
{"delta": " architecture"}
{"delta": " has"}
{"done": true, "output": "Full output text"}
```

**Response Fields**:
- `delta`: Incremental text chunk
- `done`: Boolean indicating completion
- `output`: Full output text (when done)

**Content-Type**: `application/x-ndjson`

**Status Codes**:
- `200 OK` - Success (streaming)
- `400 Bad Request` - Invalid request
- `404 Not Found` - Unknown agent
- `500 Internal Server Error` - LLM API error

**Example**:
```bash
curl -X POST http://localhost:8000/v1/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "agent": "strategist",
    "input": "What is a go-to-market strategy?",
    "provider": "anthropic"
  }'
```

---

### Image Generation

**POST** `/v1/images`

Generate images using Cipher API.

**Request Body**:
```json
{
  "provider": "cipher",
  "prompt": "A futuristic microservices architecture diagram",
  "n": 1,
  "size": "1024x1024"
}
```

**Request Fields**:
- `provider` (required): Image provider - `"cipher"` (only supported currently)
- `prompt` (required): Image generation prompt
- `n` (optional): Number of images to generate (default: 1)
- `size` (optional): Image size, e.g., `"1024x1024"` (default: `"1024x1024"`)

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

**Response Fields**:
- `data`: Array of image objects
  - `url`: Image URL (if available)
  - `b64_json`: Base64-encoded image (if available)

**Status Codes**:
- `200 OK` - Success
- `400 Bad Request` - Invalid request (unsupported provider, invalid prompt, etc.)
- `500 Internal Server Error` - Image generation API error

**Example**:
```bash
curl -X POST http://localhost:8000/v1/images \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "cipher",
    "prompt": "A futuristic microservices architecture diagram",
    "n": 1,
    "size": "1024x1024"
  }'
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

## Agent Personas

### Economist

**Focus**: Economic analysis, market trends, pricing, unit economics

**Keywords**: market, inflation, macro, econom, unit economics, pricing

**Use Case**: Economic analysis, pricing strategies, market research

### Strategist

**Focus**: Business strategy, positioning, go-to-market, competitive analysis

**Keywords**: strategy, positioning, go-to-market, gtm, competitor, moat

**Use Case**: Business strategy, competitive analysis, positioning

### Entrepreneur

**Focus**: MVP development, launching, prototyping, validation

**Keywords**: mvp, launch, prototype, hack, validate, scrappy

**Use Case**: Product development, MVP planning, rapid prototyping

### Startup

**Focus**: General startup advice and guidance

**Use Case**: General startup questions, business advice

### Auto

**Focus**: Automatically selects persona based on input keywords

**Logic**: Analyzes input to determine best persona

**Use Case**: When unsure which persona to use

---

## LLM Providers

### OpenAI

**Models**: `gpt-4o-mini`, `gpt-4`, `gpt-3.5-turbo`

**Features**: Chat completion, streaming

**Rate Limits**: Per API key tier

**Configuration**: Set `OPENAI_API_KEY` environment variable

### Anthropic

**Models**: `claude-3-opus`, `claude-3-sonnet`, `claude-3-haiku`

**Features**: Chat completion, streaming

**Rate Limits**: Per API key tier

**Configuration**: Set `ANTHROPIC_API_KEY` environment variable

### Cipher

**Models**: Cipher-specific models

**Features**: Chat completion, image generation

**Rate Limits**: Per API key tier

**Configuration**: Set `CIPHER_API_KEY` environment variable

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
- Invalid provider
- Invalid model
- Invalid request format

**404 Not Found**:
- Unknown agent persona

**500 Internal Server Error**:
- LLM API error
- Network error
- Service error

---

## Rate Limiting

Rate limiting is handled by LLM providers. The service respects provider rate limits and may return errors if limits are exceeded.

**Rate Limit Headers** (if implemented):
- `X-RateLimit-Limit` - Request limit per window
- `X-RateLimit-Remaining` - Remaining requests
- `X-RateLimit-Reset` - Reset time (Unix timestamp)

---

## Examples

### Complete Flow: Chat with Auto Agent Selection

```bash
# 1. List available agents
curl http://localhost:8000/v1/agents

# Response: ["economist", "strategist", "entrepreneur", "startup", "auto"]

# 2. Chat with auto agent selection
curl -X POST http://localhost:8000/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "agent": "auto",
    "input": "What are the unit economics of a SaaS business?",
    "provider": "openai"
  }'

# Response: {"agent": "economist", "output": "..."}
```

### Streaming Chat Example

```bash
curl -X POST http://localhost:8000/v1/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "agent": "strategist",
    "input": "What is a go-to-market strategy?",
    "provider": "anthropic"
  }' \
  --no-buffer
```

### Image Generation Example

```bash
curl -X POST http://localhost:8000/v1/images \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "cipher",
    "prompt": "A futuristic microservices architecture diagram with connecting services",
    "n": 1,
    "size": "1024x1024"
  }'
```

---

## Related Documentation

- [Component Documentation](../components/ai-service.md) - AI Service component details
- [Architecture Decision Records](../adr/) - Architectural decisions
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
