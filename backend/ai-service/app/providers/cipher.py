import os
import httpx
from fastapi import HTTPException
from app.config import settings


class CipherClient:
    def __init__(self, base_url: str, api_key: str, timeout: int = 60, image_url: str | None = None):
        self.base_url = base_url.rstrip("/")
        self.api_key = api_key
        self.timeout = timeout
        self.image_url = (image_url or settings.CIPHER_IMAGE_URL).rstrip("/")

    @classmethod
    def from_env(cls) -> "CipherClient":
        base_url = settings.CIPHER_BASE_URL
        api_key = settings.CIPHER_API_KEY
        if not api_key:
            raise HTTPException(status_code=500, detail="CIPHER_API_KEY not configured")
        return cls(base_url=base_url, api_key=api_key, image_url=settings.CIPHER_IMAGE_URL)

    async def chat(self, messages: list[dict], temperature: float, max_tokens: int, top_p: float) -> str:
        url = f"{self.base_url}?api_key={self.api_key}"
        payload = {
            "messages": messages,
            "temperature": temperature,
            "max_tokens": max_tokens,
            "top_p": top_p,
        }
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            r = await client.post(url, json=payload, headers={"Content-Type": "application/json"})
            if r.status_code >= 400:
                raise HTTPException(status_code=r.status_code, detail=f"cipher error: {r.text}")
            data = r.json()
            # Try OpenAI-like shape
            text = None
            try:
                text = data["choices"][0]["message"]["content"]
            except Exception:
                text = data.get("text") or data.get("content") or str(data)
            return text

    async def generate_images(self, prompt: str, n: int, size: str) -> list[dict]:
        url = f"{self.image_url}?api_key={self.api_key}"
        payload = {
            "prompt": prompt,
            "n": n,
            "size": size,
        }
        async with httpx.AsyncClient(timeout=self.timeout) as client:
            r = await client.post(url, json=payload, headers={"Content-Type": "application/json"})
            if r.status_code >= 400:
                raise HTTPException(status_code=r.status_code, detail=f"cipher image error: {r.text}")
            data = r.json()
            # OpenAI-like images response: { data: [ {url? or b64_json?}, ... ] }
            items = data.get("data") or []
            results: list[dict] = []
            for item in items:
                if "b64_json" in item:
                    results.append({"b64_json": item["b64_json"]})
                elif "url" in item:
                    results.append({"url": item["url"]})
                else:
                    results.append(item)
            return results


