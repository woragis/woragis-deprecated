import os
from langchain_openai import ChatOpenAI
from app.config import settings


def build(model_name: str | None, temperature: float | None):
    return ChatOpenAI(
        model=model_name or settings.OPENAI_MODEL,
        temperature=temperature if temperature is not None else settings.DEFAULT_TEMPERATURE,
        timeout=60,
    )


