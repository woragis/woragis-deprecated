"""
Pytest configuration and shared fixtures.
"""
import os
import pytest
from unittest.mock import Mock, AsyncMock, patch
from fastapi.testclient import TestClient
from langchain_core.language_models import BaseChatModel
from langchain_core.messages import AIMessage

from app.main import app
from app.config import settings


@pytest.fixture
def client():
    """FastAPI test client."""
    return TestClient(app)


@pytest.fixture
def mock_chat_model():
    """Mock LangChain chat model."""
    model = Mock(spec=BaseChatModel)
    model.ainvoke = AsyncMock(return_value=AIMessage(content="Mock response"))
    model.astream_events = AsyncMock()
    return model


@pytest.fixture
def mock_cipher_client():
    """Mock CipherClient."""
    client = Mock()
    client.chat = AsyncMock(return_value="Mock cipher response")
    client.generate_images = AsyncMock(return_value=[
        {"url": "https://example.com/image.png"}
    ])
    return client


@pytest.fixture
def sample_chat_request():
    """Sample chat request data."""
    return {
        "agent": "startup",
        "input": "What is a good MVP strategy?",
        "provider": "openai"
    }


@pytest.fixture
def sample_image_request():
    """Sample image generation request data."""
    return {
        "provider": "cipher",
        "prompt": "A futuristic cityscape",
        "n": 1,
        "size": "1024x1024"
    }


@pytest.fixture(autouse=True)
def reset_env():
    """Reset environment variables before each test."""
    # Store original values
    original_env = os.environ.copy()
    yield
    # Restore original values
    os.environ.clear()
    os.environ.update(original_env)
