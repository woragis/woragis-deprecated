import httpx
from typing import Optional
from app.config import settings


class CipherImageProvider:
    """Cipher (NoFilterGPT) image generation - uses existing infrastructure."""
    
    def __init__(self):
        self.api_key = settings.CIPHER_API_KEY
        self.base_url = settings.CIPHER_IMAGE_URL.rstrip("/")
        self.timeout = 60
    
    async def generate(
        self,
        prompt: str,
        size: Optional[str] = None,
        n: Optional[int] = None,
    ) -> list[dict]:
        """
        Generate images using Cipher API.
        
        Returns list of dicts with 'url' or 'b64_json' keys.
        """
        size = size or settings.CIPHER_IMAGE_SIZE
        n = n or settings.CIPHER_IMAGE_N
        
        url = f"{self.base_url}?api_key={self.api_key}"
        payload = {
            "prompt": prompt,
            "n": n,
            "size": size,
        }
        
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            r = await client.post(url, json=payload, headers={"Content-Type": "application/json"})
            if r.status_code >= 400:
                raise Exception(f"Cipher image error: {r.text}")
            data = r.json()
            
            # OpenAI-like images response: { data: [ {url? or b64_json?}, ... ] }
            items = data.get("data") or []
            results = []
            for item in items:
                if "b64_json" in item:
                    results.append({"b64_json": item["b64_json"]})
                elif "url" in item:
                    results.append({"url": item["url"]})
                else:
                    results.append(item)
            return results

