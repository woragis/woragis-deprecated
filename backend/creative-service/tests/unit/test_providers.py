"""
Unit tests for provider factories and provider implementations.
"""
import pytest
import os
from unittest.mock import Mock, patch, AsyncMock
from fastapi import HTTPException

from app.providers.factory import (
    ImageProviderFactory,
    DiagramProviderFactory,
    VideoProviderFactory
)


class TestImageProviderFactory:
    """Tests for ImageProviderFactory."""

    @patch('app.providers.factory.OpenAIImageProvider')
    @patch('app.providers.factory.settings')
    def test_create_openai(self, mock_settings, mock_provider_class):
        """Test creating OpenAI image provider."""
        mock_settings.OPENAI_API_KEY = "test-key"
        mock_provider_class.return_value = Mock()
        
        provider = ImageProviderFactory.create("openai")
        assert provider is not None
        mock_provider_class.assert_called_once()

    @patch('app.providers.factory.OpenAIImageProvider')
    @patch('app.providers.factory.settings')
    def test_create_openai_missing_key(self, mock_settings, mock_provider_class):
        """Test creating OpenAI provider without API key."""
        mock_settings.OPENAI_API_KEY = ""
        
        with pytest.raises(HTTPException) as exc_info:
            ImageProviderFactory.create("openai")
        assert exc_info.value.status_code == 500

    @patch('app.providers.factory.StableDiffusionProvider')
    @patch('app.providers.factory.settings')
    def test_create_stable_diffusion(self, mock_settings, mock_provider_class):
        """Test creating Stable Diffusion provider."""
        mock_settings.REPLICATE_API_KEY = "test-key"
        mock_provider_class.return_value = Mock()
        
        provider = ImageProviderFactory.create("stable-diffusion")
        assert provider is not None
        mock_provider_class.assert_called_once()

    @patch('app.providers.factory.CipherImageProvider')
    @patch('app.providers.factory.settings')
    def test_create_cipher(self, mock_settings, mock_provider_class):
        """Test creating Cipher image provider."""
        mock_settings.CIPHER_API_KEY = "test-key"
        mock_provider_class.return_value = Mock()
        
        provider = ImageProviderFactory.create("cipher")
        assert provider is not None
        mock_provider_class.assert_called_once()

    def test_create_invalid_provider(self):
        """Test creating invalid provider raises error."""
        with pytest.raises(HTTPException) as exc_info:
            ImageProviderFactory.create("invalid")
        assert exc_info.value.status_code == 400


class TestDiagramProviderFactory:
    """Tests for DiagramProviderFactory."""

    @patch('app.providers.factory.DiagramGenerator')
    @patch('app.providers.factory.settings')
    def test_create_openai(self, mock_settings, mock_generator_class):
        """Test creating OpenAI diagram generator."""
        mock_settings.OPENAI_API_KEY = "test-key"
        mock_generator_class.return_value = Mock()
        
        generator = DiagramProviderFactory.create("openai")
        assert generator is not None
        mock_generator_class.assert_called_once_with(provider="openai")

    @patch('app.providers.factory.DiagramGenerator')
    @patch('app.providers.factory.settings')
    def test_create_anthropic(self, mock_settings, mock_generator_class):
        """Test creating Anthropic diagram generator."""
        mock_settings.ANTHROPIC_API_KEY = "test-key"
        mock_generator_class.return_value = Mock()
        
        generator = DiagramProviderFactory.create("anthropic")
        assert generator is not None
        mock_generator_class.assert_called_once_with(provider="anthropic")

    @patch('app.providers.factory.DiagramGenerator')
    @patch('app.providers.factory.settings')
    def test_create_openai_missing_key(self, mock_settings, mock_generator_class):
        """Test creating OpenAI generator without API key."""
        mock_settings.OPENAI_API_KEY = ""
        
        with pytest.raises(HTTPException) as exc_info:
            DiagramProviderFactory.create("openai")
        assert exc_info.value.status_code == 500

    def test_create_invalid_provider(self):
        """Test creating invalid provider raises error."""
        with pytest.raises(HTTPException) as exc_info:
            DiagramProviderFactory.create("invalid")
        assert exc_info.value.status_code == 400


class TestVideoProviderFactory:
    """Tests for VideoProviderFactory."""

    @patch('app.providers.factory.VideoGenerator')
    @patch('app.providers.factory.settings')
    def test_create_replicate(self, mock_settings, mock_generator_class):
        """Test creating Replicate video generator."""
        mock_settings.REPLICATE_API_KEY = "test-key"
        mock_generator_class.return_value = Mock()
        
        generator = VideoProviderFactory.create("replicate")
        assert generator is not None
        mock_generator_class.assert_called_once_with(provider="replicate")

    @patch('app.providers.factory.VideoGenerator')
    @patch('app.providers.factory.settings')
    def test_create_runway(self, mock_settings, mock_generator_class):
        """Test creating Runway video generator."""
        mock_settings.RUNWAY_API_KEY = "test-key"
        mock_generator_class.return_value = Mock()
        
        generator = VideoProviderFactory.create("runway")
        assert generator is not None
        mock_generator_class.assert_called_once_with(provider="runway")

    @patch('app.providers.factory.VideoGenerator')
    @patch('app.providers.factory.settings')
    def test_create_replicate_missing_key(self, mock_settings, mock_generator_class):
        """Test creating Replicate generator without API key."""
        mock_settings.REPLICATE_API_KEY = ""
        
        with pytest.raises(HTTPException) as exc_info:
            VideoProviderFactory.create("replicate")
        assert exc_info.value.status_code == 500

    def test_create_invalid_provider(self):
        """Test creating invalid provider raises error."""
        with pytest.raises(HTTPException) as exc_info:
            VideoProviderFactory.create("invalid")
        assert exc_info.value.status_code == 400
