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

    # Defaults
    DEFAULT_TEMPERATURE: float = float(os.getenv("OPENAI_TEMPERATURE", "0.3"))

    # OpenAI (ChatGPT)
    OPENAI_API_KEY: str = os.getenv("OPENAI_API_KEY", "")
    OPENAI_MODEL: str = os.getenv("OPENAI_MODEL", "gpt-4o-mini")

    # Anthropic (Claude)
    ANTHROPIC_API_KEY: str = os.getenv("ANTHROPIC_API_KEY", "")
    ANTHROPIC_MODEL: str = os.getenv("ANTHROPIC_MODEL", "claude-3-5-sonnet-latest")

    # xAI (Grok) via OpenAI-compatible
    XAI_API_KEY: str = os.getenv("XAI_API_KEY", "")
    XAI_BASE_URL: str = os.getenv("XAI_BASE_URL", "https://api.x.ai/v1")
    XAI_MODEL: str = os.getenv("XAI_MODEL", "grok-beta")

    # Manus via OpenAI-compatible
    MANUS_API_KEY: str = os.getenv("MANUS_API_KEY", "")
    MANUS_BASE_URL: str = os.getenv("MANUS_BASE_URL", "")
    MANUS_MODEL: str = os.getenv("MANUS_MODEL", "manus-chat")

    # Cipher (NoFilterGPT)
    CIPHER_API_KEY: str = os.getenv("CIPHER_API_KEY", "")
    CIPHER_BASE_URL: str = os.getenv("CIPHER_BASE_URL", "https://api.nofiltergpt.com/v1/chat/completions")
    CIPHER_MAX_TOKENS: int = int(os.getenv("CIPHER_MAX_TOKENS", "800"))
    CIPHER_TOP_P: float = float(os.getenv("CIPHER_TOP_P", "1.0"))
    # Images
    CIPHER_IMAGE_URL: str = os.getenv("CIPHER_IMAGE_URL", "https://api.nofiltergpt.com/v1/images/generations")
    CIPHER_IMAGE_SIZE: str = os.getenv("CIPHER_IMAGE_SIZE", "1024x1024")
    CIPHER_IMAGE_N: int = int(os.getenv("CIPHER_IMAGE_N", "1"))


settings = Settings()


