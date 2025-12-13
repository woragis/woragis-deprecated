# Creative Service

AI-powered image, diagram, and video generation service for technical blog posts. This service provides APIs for creating thumbnails, architecture diagrams, workflow illustrations, and animated content for platforms like LinkedIn, Twitter, Substack, Medium, and others.

## Features

### 🎨 Image Generation
- **DALL-E 3** (OpenAI) - High-quality thumbnails and illustrations
- **Stable Diffusion XL** - Artistic and diagram-style images
- **Cipher/NoFilterGPT** - Alternative image provider
- Optimized prompts for technical content (architecture, microservices, workflows)

### 📊 Diagram Generation
- **Mermaid Diagrams** - Flowcharts, sequence diagrams, ER diagrams, etc.
- **Graphviz** - Network graphs, dependency graphs, etc.
- AI-powered code generation from natural language descriptions
- Automatic rendering to PNG, SVG, or PDF

### 🎬 Video/GIF Generation
- **Stable Video Diffusion** - Animate static images into videos
- Generate animated content from technical diagrams or illustrations
- Support for MP4 and GIF formats

## Quick Start

### Prerequisites
- Python 3.11+
- Docker (optional, for containerized deployment)
- API keys for desired providers (see [ENVIRONMENT.md](./ENVIRONMENT.md))

### Installation

1. **Clone and navigate to the service:**
```bash
cd backend/creative-service
```

2. **Install dependencies:**
```bash
pip install -r requirements.txt
```

3. **Set up environment variables:**
```bash
cp env.sample .env
# Edit .env with your API keys
```

4. **Run the service:**
```bash
uvicorn app.main:app --reload --port 8000
```

### Docker

```bash
docker build -t creative-service .
docker run -p 8000:8000 --env-file .env creative-service
```

## API Usage

### Generate Images

**POST** `/v1/images/generate`

Generate images using various AI providers.

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

**Response:**
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

```json
{
  "provider": "openai",
  "prompt": "Modern software architecture with microservices",
  "context": "Technical blog post thumbnail"
}
```

### Generate Diagrams

**POST** `/v1/diagrams/generate`

Generate technical diagrams from descriptions.

```json
{
  "description": "A microservices architecture with an API gateway, three microservices (user, order, payment), and a Redis cache layer",
  "diagram_type": "mermaid",
  "diagram_kind": "flowchart",
  "output_format": "png",
  "ai_provider": "openai"
}
```

**Response:**
```json
{
  "b64_json": "base64_encoded_diagram_image",
  "code": "graph TD\n    A[API Gateway] --> B[User Service]\n    ...",
  "format": "png",
  "diagram_type": "mermaid"
}
```

**POST** `/v1/diagrams/mermaid`

Quick endpoint for Mermaid diagrams.

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

```json
{
  "image_b64": "base64_encoded_image",
  "motion_bucket_id": 127,
  "num_frames": 25,
  "provider": "replicate"
}
```

**Response:**
```json
{
  "video_url": "https://...",
  "video_b64": "base64_encoded_video",
  "format": "mp4"
}
```

## Use Cases

### 1. Thumbnails for Blog Posts
Generate eye-catching thumbnails optimized for social media platforms.

```python
POST /v1/images/generate/thumbnail
{
  "provider": "openai",
  "prompt": "Kubernetes orchestration architecture",
  "style": "thumbnail"
}
```

### 2. Architecture Diagrams
Create clear, professional architecture diagrams from descriptions.

```python
POST /v1/diagrams/generate
{
  "description": "Three-tier architecture: web server, application server, database",
  "diagram_type": "mermaid",
  "diagram_kind": "flowchart"
}
```

### 3. Workflow Illustrations
Visualize complex workflows and processes.

```python
POST /v1/diagrams/mermaid
{
  "description": "User authentication flow with OAuth, JWT tokens, and session management",
  "diagram_kind": "sequence"
}
```

### 4. Animated Content
Transform static diagrams into engaging animations.

```python
# First generate an image/diagram
image_response = POST /v1/images/generate/...

# Then animate it
POST /v1/videos/animate
{
  "image_b64": image_response.data[0].b64_json,
  "motion_bucket_id": 100  # Subtle animation
}
```

## Supported Providers

### Image Generation
- **openai** - DALL-E 3 (best for thumbnails)
- **stable-diffusion** - Stable Diffusion XL (good for diagrams and artistic content)
- **cipher** - Cipher/NoFilterGPT (alternative)

### Diagram Generation
- **mermaid** - Flowcharts, sequence diagrams, ER diagrams, etc.
- **graphviz** - Network graphs, dependency graphs, etc.

### Video Generation
- **replicate** - Stable Video Diffusion
- **runway** - Runway Gen-2 (coming soon)

## Recommendations

### For Technical Blog Posts:

1. **Thumbnails**: Use `openai` provider with `style: "thumbnail"`
   - Best quality and consistency
   - Optimized for social media platforms

2. **Architecture Diagrams**: Use `stable-diffusion` with `generate_technical_diagram` OR use diagram generation with `mermaid`
   - Mermaid diagrams are more precise and editable
   - Stable Diffusion gives more artistic flexibility

3. **Process Flows**: Use diagram generation with `mermaid` and `diagram_kind: "sequence"` or `"flowchart"`
   - Clean, professional, and easy to modify
   - Perfect for technical documentation

4. **Animations**: Generate a static image first, then use video generation to animate it
   - Great for demonstrating workflows or system interactions
   - Keep motion_bucket_id moderate (100-150) for professional content

## Environment Setup

See [ENVIRONMENT.md](./ENVIRONMENT.md) for detailed environment variable configuration.

## Architecture

```
creative-service/
├── app/
│   ├── main.py              # FastAPI application and endpoints
│   ├── config.py            # Configuration and settings
│   └── providers/
│       ├── factory.py       # Provider factory patterns
│       ├── openai_image.py  # DALL-E image generation
│       ├── stable_diffusion.py  # Stable Diffusion images
│       ├── cipher_image.py  # Cipher images
│       ├── diagram_generator.py  # Mermaid/Graphviz diagram generation
│       └── video_generator.py    # Video/GIF generation
├── Dockerfile
├── requirements.txt
└── env.sample
```

## Development

### Running Tests
```bash
# Install test dependencies
pip install pytest pytest-asyncio httpx

# Run tests
pytest
```

### Local Development
```bash
# Install with dev dependencies
pip install -r requirements.txt

# Run with auto-reload
uvicorn app.main:app --reload --port 8000
```

## License

Part of the Woragis project.

