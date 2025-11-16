from typing import Optional

from . import openai as openai_builder
from . import anthropic as anthropic_builder
from . import xai as xai_builder
from . import manus as manus_builder
from .cipher import CipherClient


def make_model(provider: str, model_name: Optional[str], temperature: Optional[float]):
    provider = provider.lower()
    if provider == "openai":
        return openai_builder.build(model_name, temperature)
    if provider == "anthropic":
        return anthropic_builder.build(model_name, temperature)
    if provider == "xai":
        return xai_builder.build(model_name, temperature)
    if provider == "manus":
        return manus_builder.build(model_name, temperature)
    raise ValueError(f"Unsupported provider '{provider}'")


