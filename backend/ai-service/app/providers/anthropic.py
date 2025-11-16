import os
from langchain_anthropic import ChatAnthropic
from app.config import settings


def build(model_name: str | None, temperature: float | None):
    return ChatAnthropic(
        model=model_name or settings.ANTHROPIC_MODEL,
        temperature=temperature if temperature is not None else settings.DEFAULT_TEMPERATURE,
        timeout=60,
    )


