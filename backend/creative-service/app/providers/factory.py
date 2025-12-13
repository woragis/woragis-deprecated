from typing import Literal, Optional
from fastapi import HTTPException
from app.config import settings
from .openai_image import OpenAIImageProvider
from .stable_diffusion import StableDiffusionProvider
from .cipher_image import CipherImageProvider
from .diagram_generator import DiagramGenerator
from .video_generator import VideoGenerator


class ImageProviderFactory:
    @staticmethod
    def create(provider: Literal["openai", "stable-diffusion", "cipher"]):
        if provider == "openai":
            if not settings.OPENAI_API_KEY:
                raise HTTPException(status_code=500, detail="OPENAI_API_KEY not configured")
            return OpenAIImageProvider()
        elif provider == "stable-diffusion":
            if not settings.REPLICATE_API_KEY:
                raise HTTPException(status_code=500, detail="REPLICATE_API_KEY not configured")
            return StableDiffusionProvider()
        elif provider == "cipher":
            if not settings.CIPHER_API_KEY:
                raise HTTPException(status_code=500, detail="CIPHER_API_KEY not configured")
            return CipherImageProvider()
        else:
            raise HTTPException(status_code=400, detail=f"Unknown image provider: {provider}")


class DiagramProviderFactory:
    @staticmethod
    def create(provider: Literal["openai", "anthropic"] = "openai"):
        if provider == "openai":
            if not settings.OPENAI_API_KEY:
                raise HTTPException(status_code=500, detail="OPENAI_API_KEY not configured")
            return DiagramGenerator(provider="openai")
        elif provider == "anthropic":
            if not settings.ANTHROPIC_API_KEY:
                raise HTTPException(status_code=500, detail="ANTHROPIC_API_KEY not configured")
            return DiagramGenerator(provider="anthropic")
        else:
            raise HTTPException(status_code=400, detail=f"Unknown diagram provider: {provider}")


class VideoProviderFactory:
    @staticmethod
    def create(provider: Literal["replicate", "runway"] = "replicate"):
        if provider == "replicate":
            if not settings.REPLICATE_API_KEY:
                raise HTTPException(status_code=500, detail="REPLICATE_API_KEY not configured")
            return VideoGenerator(provider="replicate")
        elif provider == "runway":
            if not settings.RUNWAY_API_KEY:
                raise HTTPException(status_code=500, detail="RUNWAY_API_KEY not configured")
            return VideoGenerator(provider="runway")
        else:
            raise HTTPException(status_code=400, detail=f"Unknown video provider: {provider}")

