"""
Unit tests for provider factory and provider implementations.
"""
import pytest
import os
from unittest.mock import Mock, patch, MagicMock
from langchain_core.language_models import BaseChatModel

from app.providers import make_model
from app.providers.factory import make_model as factory_make_model


class TestProviderFactory:
    """Tests for provider factory."""

    @patch('app.providers.openai.build')
    def test_make_model_openai(self, mock_openai_build, mock_chat_model):
        """Test creating OpenAI model."""
        mock_openai_build.return_value = mock_chat_model
        model = make_model("openai", None, None)
        assert model is not None
        mock_openai_build.assert_called_once_with(None, None)

    @patch('app.providers.openai.build')
    def test_make_model_openai_with_params(self, mock_openai_build, mock_chat_model):
        """Test creating OpenAI model with parameters."""
        mock_openai_build.return_value = mock_chat_model
        model = make_model("openai", "gpt-4", 0.7)
        assert model is not None
        mock_openai_build.assert_called_once_with("gpt-4", 0.7)

    @patch('app.providers.anthropic.build')
    def test_make_model_anthropic(self, mock_anthropic_build, mock_chat_model):
        """Test creating Anthropic model."""
        mock_anthropic_build.return_value = mock_chat_model
        model = make_model("anthropic", None, None)
        assert model is not None
        mock_anthropic_build.assert_called_once_with(None, None)

    @patch('app.providers.xai.build')
    def test_make_model_xai(self, mock_xai_build, mock_chat_model):
        """Test creating xAI model."""
        mock_xai_build.return_value = mock_chat_model
        model = make_model("xai", None, None)
        assert model is not None
        mock_xai_build.assert_called_once_with(None, None)

    @patch('app.providers.manus.build')
    def test_make_model_manus(self, mock_manus_build, mock_chat_model):
        """Test creating Manus model."""
        mock_manus_build.return_value = mock_chat_model
        model = make_model("manus", None, None)
        assert model is not None
        mock_manus_build.assert_called_once_with(None, None)

    def test_make_model_case_insensitive(self, mock_chat_model):
        """Test that provider names are case-insensitive."""
        with patch('app.providers.openai.build', return_value=mock_chat_model):
            model1 = make_model("OPENAI", None, None)
            model2 = make_model("openai", None, None)
            assert model1 is not None
            assert model2 is not None

    def test_make_model_unsupported(self):
        """Test creating unsupported provider raises error."""
        with pytest.raises(ValueError, match="Unsupported provider"):
            make_model("unsupported", None, None)

    @pytest.mark.parametrize("provider", ["openai", "anthropic", "xai", "manus"])
    def test_make_model_all_providers(self, provider, mock_chat_model):
        """Test creating all supported providers."""
        with patch(f'app.providers.{provider}.build', return_value=mock_chat_model):
            model = make_model(provider, None, None)
            assert model is not None


class TestCipherClient:
    """Tests for CipherClient."""

    @pytest.fixture
    def cipher_client(self):
        """Create a CipherClient instance."""
        from app.providers.cipher import CipherClient
        return CipherClient(
            base_url="https://api.test.com/v1/chat/completions",
            api_key="test-key",
            timeout=60
        )

    def test_cipher_client_init(self, cipher_client):
        """Test CipherClient initialization."""
        assert cipher_client.base_url == "https://api.test.com/v1/chat/completions"
        assert cipher_client.api_key == "test-key"
        assert cipher_client.timeout == 60

    def test_cipher_client_url_stripping(self):
        """Test that base_url is stripped of trailing slashes."""
        from app.providers.cipher import CipherClient
        client = CipherClient(
            base_url="https://api.test.com/",
            api_key="test-key"
        )
        assert not client.base_url.endswith("/")

    @patch.dict(os.environ, {"CIPHER_API_KEY": "test-key"})
    @patch('app.providers.cipher.settings')
    def test_cipher_client_from_env(self, mock_settings):
        """Test creating CipherClient from environment."""
        from app.providers.cipher import CipherClient
        mock_settings.CIPHER_BASE_URL = "https://api.test.com"
        mock_settings.CIPHER_API_KEY = "test-key"
        mock_settings.CIPHER_IMAGE_URL = "https://api.test.com/images"
        
        client = CipherClient.from_env()
        assert client.api_key == "test-key"
        assert client.base_url == "https://api.test.com"

    @patch.dict(os.environ, {}, clear=True)
    @patch('app.providers.cipher.settings')
    def test_cipher_client_from_env_missing_key(self, mock_settings):
        """Test that missing API key raises error."""
        from app.providers.cipher import CipherClient
        from fastapi import HTTPException
        
        mock_settings.CIPHER_API_KEY = ""
        mock_settings.CIPHER_BASE_URL = "https://api.test.com"
        
        with pytest.raises(HTTPException) as exc_info:
            CipherClient.from_env()
        assert exc_info.value.status_code == 500

    @pytest.mark.asyncio
    async def test_cipher_client_chat_success(self, cipher_client):
        """Test successful chat request."""
        import httpx
        from unittest.mock import AsyncMock
        
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "choices": [{"message": {"content": "Test response"}}]
        }
        
        with patch('httpx.AsyncClient') as mock_client:
            mock_client.return_value.__aenter__.return_value.post = AsyncMock(return_value=mock_response)
            result = await cipher_client.chat(
                messages=[{"role": "user", "content": "Hello"}],
                temperature=0.7,
                max_tokens=100,
                top_p=1.0
            )
            assert result == "Test response"

    @pytest.mark.asyncio
    async def test_cipher_client_chat_error(self, cipher_client):
        """Test chat request with error response."""
        import httpx
        from unittest.mock import AsyncMock
        from fastapi import HTTPException
        
        mock_response = Mock()
        mock_response.status_code = 400
        mock_response.text = "Bad request"
        
        with patch('httpx.AsyncClient') as mock_client:
            mock_client.return_value.__aenter__.return_value.post = AsyncMock(return_value=mock_response)
            with pytest.raises(HTTPException) as exc_info:
                await cipher_client.chat(
                    messages=[{"role": "user", "content": "Hello"}],
                    temperature=0.7,
                    max_tokens=100,
                    top_p=1.0
                )
            assert exc_info.value.status_code == 400

    @pytest.mark.asyncio
    async def test_cipher_client_generate_images(self, cipher_client):
        """Test image generation."""
        import httpx
        from unittest.mock import AsyncMock
        
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "data": [{"url": "https://example.com/image.png"}]
        }
        
        with patch('httpx.AsyncClient') as mock_client:
            mock_client.return_value.__aenter__.return_value.post = AsyncMock(return_value=mock_response)
            result = await cipher_client.generate_images(
                prompt="A test image",
                n=1,
                size="1024x1024"
            )
            assert len(result) == 1
            assert "url" in result[0]
