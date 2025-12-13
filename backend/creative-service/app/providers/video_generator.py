import base64
from typing import Optional
import replicate
from app.config import settings


class VideoGenerator:
    """
    Generate GIFs and videos from images or prompts.
    Supports Stable Video Diffusion and Runway.
    """
    
    def __init__(self, provider: str = "replicate"):
        self.provider = provider
        if provider == "replicate":
            self.client = replicate.Client(api_token=settings.REPLICATE_API_KEY)
        elif provider == "runway":
            # Runway API setup would go here
            self.api_key = settings.RUNWAY_API_KEY
    
    async def generate_from_image(
        self,
        image_url: Optional[str] = None,
        image_b64: Optional[str] = None,
        motion_bucket_id: int = 127,  # 1-255, higher = more motion
        cond_aug: float = 0.02,  # Conditioning augmentation
        num_frames: int = 25,
        num_inference_steps: int = 25,
    ) -> dict:
        """
        Generate a video from an existing image using Stable Video Diffusion.
        
        Args:
            image_url: URL of the input image
            image_b64: Base64 encoded image (alternative to URL)
            motion_bucket_id: Motion intensity (1-255)
            cond_aug: Conditioning augmentation (0.01-1.0)
            num_frames: Number of frames (14-25)
            num_inference_steps: Inference steps (25-50)
        
        Returns dict with 'video_url' or 'video_b64'.
        """
        if self.provider != "replicate":
            raise NotImplementedError("Only Replicate/Stable Video Diffusion is currently supported")
        
        # Convert b64 to file upload if needed
        if image_b64:
            # For Replicate, upload the base64 image as a file
            import tempfile
            import os
            import io
            from PIL import Image
            
            image_bytes = base64.b64decode(image_b64)
            # Upload to Replicate using file upload
            image_file = io.BytesIO(image_bytes)
            uploaded = self.client.files.upload(image_file)
            image_url = uploaded.url
            
            # Clean up
            image_file.close()
        
        if not image_url:
            raise ValueError("Either image_url or image_b64 must be provided")
        
        # Run in executor since replicate.run may block
        import asyncio
        loop = asyncio.get_event_loop()
        output = await loop.run_in_executor(
            None,
            lambda: self.client.run(
                settings.REPLICATE_VIDEO_MODEL,
                input={
                    "image": image_url,
                    "motion_bucket_id": motion_bucket_id,
                    "cond_aug": cond_aug,
                    "num_frames": num_frames,
                    "num_inference_steps": num_inference_steps,
                }
            )
        )
        
        # Replicate returns a URL or list of URLs
        video_url = output if isinstance(output, str) else output[0] if output else None
        
        # Download and convert to base64
        import httpx
        async with httpx.AsyncClient() as client:
            response = await client.get(video_url)
            video_b64 = base64.b64encode(response.content).decode('utf-8')
        
        return {
            "video_url": video_url,
            "video_b64": video_b64,
            "format": "mp4",
        }
    
    async def generate_animated_gif(
        self,
        image_url: Optional[str] = None,
        image_b64: Optional[str] = None,
        duration: float = 2.0,  # seconds
        loop: bool = True,
        **kwargs
    ) -> dict:
        """
        Generate an animated GIF from an image.
        Creates a video first, then converts to GIF.
        """
        video_result = await self.generate_from_image(image_url, image_b64, **kwargs)
        
        # Convert video to GIF
        # This would require ffmpeg or similar
        # For now, return the video and let the client convert if needed
        return {
            **video_result,
            "format": "mp4",  # Return as MP4, can be converted to GIF client-side
            "note": "Convert to GIF using ffmpeg or imageio",
        }
    
    async def generate_from_prompt(
        self,
        prompt: str,
        duration: int = 5,  # seconds
        **kwargs
    ) -> dict:
        """
        Generate a video directly from a text prompt.
        This is more advanced and may require different models.
        """
        # For now, this would need a text-to-video model
        # Runway Gen-2 or similar could be integrated here
        raise NotImplementedError("Text-to-video from prompt not yet implemented. Use generate_from_image instead.")

