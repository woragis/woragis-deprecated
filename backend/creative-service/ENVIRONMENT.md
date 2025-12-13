## Creative Service Environment Variables

Copy these into a `.env` file at `backend/creative-service/.env` or export them in your shell.

### CORS
- `CORS_ENABLED` (default: true)
- `CORS_ALLOWED_ORIGINS` (default: "*")

### Image Generation Providers

#### OpenAI (DALL-E)
- `OPENAI_API_KEY` (required for provider "openai")
- `DALL_E_MODEL` (default: dall-e-3)
- `DALL_E_QUALITY` (default: standard) - Options: standard, hd
- `DALL_E_SIZE` (default: 1024x1024) - Options: 1024x1024, 1792x1024, 1024x1792

#### Stable Diffusion (via Replicate)
- `REPLICATE_API_KEY` (required for provider "stable-diffusion")
- `STABLE_DIFFUSION_MODEL` (default: stability-ai/sdxl:...)

#### Cipher (NoFilterGPT)
- `CIPHER_API_KEY` (required for provider "cipher")
- `CIPHER_IMAGE_URL` (default: https://api.nofiltergpt.com/v1/images/generations)
- `CIPHER_IMAGE_SIZE` (default: 1024x1024)
- `CIPHER_IMAGE_N` (default: 1)

### Diagram Generation (AI Code Generation)

#### OpenAI
- `OPENAI_API_KEY` (same as above)
- `OPENAI_MODEL` (default: gpt-4o-mini) - Used for generating diagram code

#### Anthropic (Claude)
- `ANTHROPIC_API_KEY` (required for ai_provider "anthropic")
- `ANTHROPIC_MODEL` (default: claude-3-5-sonnet-latest)

### Video/GIF Generation

#### Replicate (Stable Video Diffusion)
- `REPLICATE_API_KEY` (same as above)
- `REPLICATE_VIDEO_MODEL` (default: stability-ai/stable-video-diffusion:...)

#### Runway (optional, future support)
- `RUNWAY_API_KEY` (for Runway Gen-2 video generation)

### Defaults
- `DEFAULT_IMAGE_FORMAT` (default: png) - Options: png, jpeg, webp
- `DEFAULT_DIAGRAM_FORMAT` (default: png) - Options: png, svg, pdf

## Setup Notes

1. **Graphviz**: The Dockerfile includes graphviz, but for local development, install:
   - macOS: `brew install graphviz`
   - Ubuntu/Debian: `apt-get install graphviz`
   - Windows: Download from https://graphviz.org/download/

2. **Mermaid**: For Mermaid diagrams, the service uses the online API (mermaid.ink) by default. For offline/local rendering, install:
   - `npm install -g @mermaid-js/mermaid-cli`

3. **Recommended Setup**:
   - At minimum, set up `OPENAI_API_KEY` for DALL-E images and diagram code generation
   - Add `REPLICATE_API_KEY` for Stable Diffusion images and video generation
   - Optional: Add `ANTHROPIC_API_KEY` for alternative diagram code generation

