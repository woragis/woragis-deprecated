import os
from dataclasses import dataclass
from dotenv import load_dotenv


# Load variables from .env if present
load_dotenv()


@dataclass(frozen=True)
class Settings:
    # General
    CORS_ENABLED: bool = (os.getenv("CORS_ENABLED", "true").lower() == "true")
    CORS_ALLOWED_ORIGINS: str = os.getenv("CORS_ALLOWED_ORIGINS", "*")

    # OpenAI DALL-E
    OPENAI_API_KEY: str = os.getenv("OPENAI_API_KEY", "")
    DALL_E_MODEL: str = os.getenv("DALL_E_MODEL", "dall-e-3")
    DALL_E_QUALITY: str = os.getenv("DALL_E_QUALITY", "standard")  # standard or hd
    DALL_E_SIZE: str = os.getenv("DALL_E_SIZE", "1024x1024")  # 1024x1024, 1792x1024, or 1024x1792

    # Stable Diffusion (via Replicate)
    REPLICATE_API_KEY: str = os.getenv("REPLICATE_API_KEY", "")
    STABLE_DIFFUSION_MODEL: str = os.getenv("STABLE_DIFFUSION_MODEL", "stability-ai/sdxl:39ed52f2a78e934b3ba6e2a89f5b1c712de7dfea535525255b1aa35c5565e08b")

    # Cipher (NoFilterGPT) - already available from ai-service
    CIPHER_API_KEY: str = os.getenv("CIPHER_API_KEY", "")
    CIPHER_IMAGE_URL: str = os.getenv("CIPHER_IMAGE_URL", "https://api.nofiltergpt.com/v1/images/generations")
    CIPHER_IMAGE_SIZE: str = os.getenv("CIPHER_IMAGE_SIZE", "1024x1024")
    CIPHER_IMAGE_N: int = int(os.getenv("CIPHER_IMAGE_N", "1"))

    # Diagram Generation (AI + Rendering)
    OPENAI_MODEL: str = os.getenv("OPENAI_MODEL", "gpt-4o-mini")  # For generating diagram code
    ANTHROPIC_API_KEY: str = os.getenv("ANTHROPIC_API_KEY", "")
    ANTHROPIC_MODEL: str = os.getenv("ANTHROPIC_MODEL", "claude-3-5-sonnet-latest")

    # Video/GIF Generation
    REPLICATE_VIDEO_MODEL: str = os.getenv("REPLICATE_VIDEO_MODEL", "stability-ai/stable-video-diffusion:3f0457e4619daac51203dedb472816fd4af51f3149fa7a9e0b5ffcf1b8172438")
    RUNWAY_API_KEY: str = os.getenv("RUNWAY_API_KEY", "")
    
    # Default outputs
    DEFAULT_IMAGE_FORMAT: str = os.getenv("DEFAULT_IMAGE_FORMAT", "png")  # png, jpeg, webp
    DEFAULT_DIAGRAM_FORMAT: str = os.getenv("DEFAULT_DIAGRAM_FORMAT", "png")  # png, svg, pdf


settings = Settings()

