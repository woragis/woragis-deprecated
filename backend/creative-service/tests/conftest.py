"""
Pytest configuration and shared fixtures.
"""
import os
import pytest
from unittest.mock import Mock, AsyncMock, patch
from fastapi.testclient import TestClient

from app.main import app
from app.config import settings


@pytest.fixture
def client():
    """FastAPI test client."""
    return TestClient(app)


@pytest.fixture
def mock_image_provider():
    """Mock image provider."""
    provider = Mock()
    provider.generate = AsyncMock(return_value=[
        {"url": "https://example.com/image.png"}
    ])
    provider.generate_with_enhanced_prompt = AsyncMock(return_value=[
        {"url": "https://example.com/image.png"}
    ])
    provider.generate_technical_diagram = AsyncMock(return_value=[
        {"url": "https://example.com/diagram.png"}
    ])
    return provider


@pytest.fixture
def mock_diagram_generator():
    """Mock diagram generator."""
    generator = Mock()
    generator.generate = AsyncMock(return_value={
        "b64_json": "base64encodedstring",
        "code": "graph TD\nA[Start] --> B[End]",
        "format": "png"
    })
    return generator


@pytest.fixture
def mock_video_generator():
    """Mock video generator."""
    generator = Mock()
    generator.generate_from_image = AsyncMock(return_value={
        "video_url": "https://example.com/video.mp4",
        "video_b64": None,
        "format": "mp4"
    })
    return generator


@pytest.fixture
def sample_image_request():
    """Sample image generation request data."""
    return {
        "provider": "openai",
        "prompt": "A futuristic cityscape",
        "size": "1024x1024",
        "style": "technical",
        "n": 1
    }


@pytest.fixture
def sample_diagram_request():
    """Sample diagram generation request data."""
    return {
        "description": "A simple flowchart showing user login process",
        "diagram_type": "mermaid",
        "diagram_kind": "flowchart",
        "output_format": "png",
        "ai_provider": "openai"
    }


@pytest.fixture
def sample_video_request():
    """Sample video generation request data."""
    return {
        "image_url": "https://example.com/image.png",
        "provider": "replicate",
        "motion_bucket_id": 127,
        "num_frames": 25
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
