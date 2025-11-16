## AI Service Environment Variables

Copy these into a `.env` file at `backend/ai-service/.env` or export them in your shell. The service already loads `.env`.

- CORS
  - CORS_ENABLED (default: true)
  - CORS_ALLOWED_ORIGINS (default: "*")

- Defaults
  - OPENAI_TEMPERATURE (default: 0.3)

- OpenAI (ChatGPT)
  - OPENAI_API_KEY (required for provider "openai")
  - OPENAI_MODEL (default: gpt-4o-mini)

- Anthropic (Claude)
  - ANTHROPIC_API_KEY (required for provider "anthropic")
  - ANTHROPIC_MODEL (default: claude-3-5-sonnet-latest)

- xAI (Grok)
  - XAI_API_KEY (required for provider "xai")
  - XAI_BASE_URL (default: https://api.x.ai/v1)
  - XAI_MODEL (default: grok-beta)

- Manus (OpenAI-compatible)
  - MANUS_API_KEY (required for provider "manus")
  - MANUS_BASE_URL (required for provider "manus")
  - MANUS_MODEL (default: manus-chat)

- Cipher (NoFilterGPT)
  - CIPHER_API_KEY (required for provider "cipher")
  - CIPHER_BASE_URL (default: https://api.nofiltergpt.com/v1/chat/completions)
  - CIPHER_MAX_TOKENS (default: 800)
  - CIPHER_TOP_P (default: 1.0)
  - CIPHER_IMAGE_URL (default: https://api.nofiltergpt.com/v1/images/generations)
  - CIPHER_IMAGE_SIZE (default: 1024x1024)
  - CIPHER_IMAGE_N (default: 1)


