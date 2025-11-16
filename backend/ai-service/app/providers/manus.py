import os
from fastapi import HTTPException
from langchain_openai import ChatOpenAI
from app.config import settings


def build(model_name: str | None, temperature: float | None):
    base_url = settings.MANUS_BASE_URL
    api_key = settings.MANUS_API_KEY
    if not base_url:
        raise HTTPException(status_code=500, detail="MANUS_BASE_URL not configured")
    return ChatOpenAI(
        model=model_name or settings.MANUS_MODEL,
        temperature=temperature if temperature is not None else settings.DEFAULT_TEMPERATURE,
        timeout=60,
        base_url=base_url,
        api_key=api_key,
    )


