import base64
import httpx
from typing import Optional
from openai import AsyncOpenAI
from app.config import settings


class OpenAIImageProvider:
    """OpenAI DALL-E image generation provider - best for high-quality thumbnails and illustrations."""
    
    def __init__(self):
        self.client = AsyncOpenAI(api_key=settings.OPENAI_API_KEY)
    
    async def generate(
        self,
        prompt: str,
        size: Optional[str] = None,
        quality: Optional[str] = None,
        n: int = 1,
    ) -> list[dict]:
        """
        Generate images using DALL-E.
        
        Returns list of dicts with 'url' or 'b64_json' keys.
        """
        size = size or settings.DALL_E_SIZE
        quality = quality or settings.DALL_E_QUALITY
        
        # DALL-E 3 only supports n=1
        if settings.DALL_E_MODEL == "dall-e-3" and n > 1:
            n = 1
        
        response = await self.client.images.generate(
            model=settings.DALL_E_MODEL,
            prompt=prompt,
            size=size,
            quality=quality,
            n=n,
            response_format="b64_json",  # Get base64 for easier handling
        )
        
        results = []
        for image in response.data:
            results.append({
                "b64_json": image.b64_json,
                "url": None,  # DALL-E 3 doesn't return URLs when using b64_json
            })
        
        return results
    
    async def generate_with_enhanced_prompt(
        self,
        base_prompt: str,
        style: str = "technical",
        context: Optional[str] = None,
        **kwargs
    ) -> list[dict]:
        """
        Generate images with enhanced prompts optimized for technical content.
        
        Args:
            base_prompt: The base description
            style: One of 'technical', 'diagram', 'thumbnail', 'illustration'
            context: Additional context about the technical topic
        """
        # Enhance prompt based on style
        style_prompts = {
            "technical": "professional, clean, modern, technical illustration, high quality, detailed",
            "diagram": "schematic style, clear lines, minimal, professional diagram, technical drawing",
            "thumbnail": "eye-catching thumbnail, bold colors, clean design, professional, social media ready",
            "illustration": "beautiful illustration, detailed, professional, modern style, high quality",
        }
        
        style_text = style_prompts.get(style, style_prompts["technical"])
        
        enhanced_prompt = f"{base_prompt}, {style_text}"
        if context:
            enhanced_prompt = f"{enhanced_prompt}, context: {context}"
        
        return await self.generate(enhanced_prompt, **kwargs)

