import base64
from typing import Optional
import replicate
from app.config import settings


class StableDiffusionProvider:
    """Stable Diffusion via Replicate - good for artistic and diagram-style images."""
    
    def __init__(self):
        self.client = replicate.Client(api_token=settings.REPLICATE_API_KEY)
    
    async def generate(
        self,
        prompt: str,
        negative_prompt: Optional[str] = None,
        width: int = 1024,
        height: int = 1024,
        num_outputs: int = 1,
        guidance_scale: float = 7.5,
        num_inference_steps: int = 50,
    ) -> list[dict]:
        """
        Generate images using Stable Diffusion.
        
        Returns list of dicts with 'url' keys (Replicate returns URLs).
        """
        if negative_prompt is None:
            negative_prompt = "blurry, low quality, distorted, watermark, text"
        
        # Run in executor since replicate.run may block
        import asyncio
        loop = asyncio.get_event_loop()
        output = await loop.run_in_executor(
            None,
            lambda: self.client.run(
                settings.STABLE_DIFFUSION_MODEL,
                input={
                    "prompt": prompt,
                    "negative_prompt": negative_prompt,
                    "width": width,
                    "height": height,
                    "num_outputs": num_outputs,
                    "guidance_scale": guidance_scale,
                    "num_inference_steps": num_inference_steps,
                }
            )
        )
        
        # Replicate returns a list of URLs or a single URL
        urls = output if isinstance(output, list) else [output]
        
        results = []
        for url in urls:
            # Fetch the image and convert to base64 for consistency
            import httpx
            async with httpx.AsyncClient() as client:
                response = await client.get(url)
                b64_data = base64.b64encode(response.content).decode('utf-8')
                results.append({
                    "url": url,
                    "b64_json": b64_data,
                })
        
        return results
    
    async def generate_technical_diagram(
        self,
        description: str,
        diagram_type: str = "architecture",
        **kwargs
    ) -> list[dict]:
        """
        Generate technical diagrams optimized for architecture, microservices, etc.
        
        Args:
            description: Description of the diagram content
            diagram_type: One of 'architecture', 'flowchart', 'database', 'network', 'microservices'
        """
        type_prompts = {
            "architecture": "software architecture diagram, boxes and arrows, clean lines, professional technical drawing",
            "flowchart": "flowchart, process flow diagram, decision boxes, arrows connecting components",
            "database": "database schema diagram, tables, relationships, ER diagram style",
            "network": "network diagram, servers, connections, infrastructure, professional network diagram",
            "microservices": "microservices architecture, service boxes, API connections, distributed system diagram",
        }
        
        diagram_prompt = type_prompts.get(diagram_type, type_prompts["architecture"])
        
        prompt = f"{description}, {diagram_prompt}, schematic style, technical illustration, high detail, clean design"
        
        # Optimize for diagrams
        diagram_kwargs = {
            "negative_prompt": "photorealistic, photograph, blurry, low quality, watermark",
            "guidance_scale": 9.0,  # Higher guidance for more precise diagrams
            **kwargs
        }
        
        return await self.generate(prompt, **diagram_kwargs)

