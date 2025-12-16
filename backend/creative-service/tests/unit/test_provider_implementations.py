"""
Unit tests for provider implementations.
"""
import pytest
from unittest.mock import Mock, AsyncMock, patch
from fastapi import HTTPException

from app.providers.openai_image import OpenAIImageProvider
from app.providers.cipher_image import CipherImageProvider
from app.providers.stable_diffusion import StableDiffusionProvider
from app.providers.diagram_generator import DiagramGenerator
from app.providers.video_generator import VideoGenerator


class TestOpenAIImageProvider:
    """Tests for OpenAIImageProvider."""

    @patch('app.providers.openai_image.AsyncOpenAI')
    @patch('app.providers.openai_image.settings')
    def test_init(self, mock_settings, mock_openai_class):
        """Test OpenAIImageProvider initialization."""
        mock_settings.OPENAI_API_KEY = "test-key"
        mock_settings.DALL_E_MODEL = "dall-e-3"
        mock_settings.DALL_E_QUALITY = "standard"
        mock_settings.DALL_E_SIZE = "1024x1024"
        
        provider = OpenAIImageProvider()
        assert provider is not None
        mock_openai_class.assert_called_once_with(api_key="test-key")

    @patch('app.providers.openai_image.AsyncOpenAI')
    @patch('app.providers.openai_image.settings')
    @pytest.mark.asyncio
    async def test_generate(self, mock_settings, mock_openai_class):
        """Test image generation."""
        mock_settings.DALL_E_MODEL = "dall-e-3"
        mock_settings.DALL_E_QUALITY = "standard"
        mock_settings.DALL_E_SIZE = "1024x1024"
        
        mock_client = Mock()
        mock_image = Mock()
        mock_image.b64_json = "base64string"
        mock_response = Mock()
        mock_response.data = [mock_image]
        mock_client.images.generate = AsyncMock(return_value=mock_response)
        mock_openai_class.return_value = mock_client
        
        provider = OpenAIImageProvider()
        result = await provider.generate("A test image")
        
        assert len(result) == 1
        assert result[0]["b64_json"] == "base64string"
        assert result[0]["url"] is None

    @patch('app.providers.openai_image.AsyncOpenAI')
    @patch('app.providers.openai_image.settings')
    @pytest.mark.asyncio
    async def test_generate_with_enhanced_prompt(self, mock_settings, mock_openai_class):
        """Test enhanced prompt generation."""
        mock_settings.DALL_E_MODEL = "dall-e-3"
        mock_settings.DALL_E_QUALITY = "standard"
        mock_settings.DALL_E_SIZE = "1024x1024"
        
        mock_client = Mock()
        mock_image = Mock()
        mock_image.b64_json = "base64string"
        mock_response = Mock()
        mock_response.data = [mock_image]
        mock_client.images.generate = AsyncMock(return_value=mock_response)
        mock_openai_class.return_value = mock_client
        
        provider = OpenAIImageProvider()
        result = await provider.generate_with_enhanced_prompt(
            base_prompt="A diagram",
            style="technical",
            context="System architecture"
        )
        
        assert len(result) == 1
        # Verify enhanced prompt was used
        call_args = mock_client.images.generate.call_args
        assert "technical" in call_args[1]["prompt"].lower() or "diagram" in call_args[1]["prompt"].lower()


class TestCipherImageProvider:
    """Tests for CipherImageProvider."""

    @patch('app.providers.cipher_image.settings')
    def test_init(self, mock_settings):
        """Test CipherImageProvider initialization."""
        mock_settings.CIPHER_API_KEY = "test-key"
        mock_settings.CIPHER_IMAGE_URL = "https://api.test.com/images"
        mock_settings.CIPHER_IMAGE_SIZE = "1024x1024"
        mock_settings.CIPHER_IMAGE_N = 1
        
        provider = CipherImageProvider()
        assert provider is not None

    @patch('app.providers.cipher_image.httpx')
    @patch('app.providers.cipher_image.settings')
    @pytest.mark.asyncio
    async def test_generate(self, mock_settings, mock_httpx):
        """Test Cipher image generation."""
        mock_settings.CIPHER_API_KEY = "test-key"
        mock_settings.CIPHER_IMAGE_URL = "https://api.test.com/images"
        mock_settings.CIPHER_IMAGE_SIZE = "1024x1024"
        mock_settings.CIPHER_IMAGE_N = 1
        
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "data": [{"url": "https://example.com/image.png"}]
        }
        
        mock_client = Mock()
        mock_client.__aenter__ = AsyncMock(return_value=mock_client)
        mock_client.__aexit__ = AsyncMock(return_value=None)
        mock_client.post = AsyncMock(return_value=mock_response)
        mock_httpx.AsyncClient.return_value = mock_client
        
        provider = CipherImageProvider()
        result = await provider.generate("A test image")
        
        assert len(result) == 1
        assert "url" in result[0] or "b64_json" in result[0]


class TestStableDiffusionProvider:
    """Tests for StableDiffusionProvider."""

    @patch('app.providers.stable_diffusion.settings')
    def test_init(self, mock_settings):
        """Test StableDiffusionProvider initialization."""
        mock_settings.REPLICATE_API_KEY = "test-key"
        mock_settings.STABLE_DIFFUSION_MODEL = "test-model"
        
        provider = StableDiffusionProvider()
        assert provider is not None

    @patch('httpx.AsyncClient')
    @patch('app.providers.stable_diffusion.replicate')
    @patch('app.providers.stable_diffusion.settings')
    @pytest.mark.asyncio
    async def test_generate(self, mock_settings, mock_replicate, mock_httpx_client):
        """Test Stable Diffusion image generation."""
        mock_settings.REPLICATE_API_KEY = "test-key"
        mock_settings.STABLE_DIFFUSION_MODEL = "test-model"
        
        # Mock the client and its run method
        mock_client = Mock()
        mock_output = ["https://example.com/image.png"]
        mock_client.run = Mock(return_value=mock_output)
        mock_replicate.Client.return_value = mock_client
        
        # Mock httpx for image fetching
        mock_http_client = Mock()
        mock_http_client.__aenter__ = AsyncMock(return_value=mock_http_client)
        mock_http_client.__aexit__ = AsyncMock(return_value=None)
        mock_response = Mock()
        mock_response.content = b"image_data"
        mock_http_client.get = AsyncMock(return_value=mock_response)
        mock_httpx_client.return_value = mock_http_client
        
        provider = StableDiffusionProvider()
        result = await provider.generate("A test image")
        
        assert len(result) > 0
        assert "url" in result[0] or "b64_json" in result[0]


class TestDiagramGenerator:
    """Tests for DiagramGenerator."""

    @patch('app.providers.diagram_generator.AsyncOpenAI')
    @patch('app.providers.diagram_generator.settings')
    def test_init_openai(self, mock_settings, mock_openai_class):
        """Test DiagramGenerator initialization with OpenAI."""
        mock_settings.OPENAI_API_KEY = "test-key"
        mock_settings.OPENAI_MODEL = "gpt-4o-mini"
        
        generator = DiagramGenerator(provider="openai")
        assert generator is not None
        assert generator.provider == "openai"

    @patch('app.providers.diagram_generator.anthropic')
    @patch('app.providers.diagram_generator.settings')
    def test_init_anthropic(self, mock_settings, mock_anthropic):
        """Test DiagramGenerator initialization with Anthropic."""
        mock_settings.ANTHROPIC_API_KEY = "test-key"
        mock_settings.ANTHROPIC_MODEL = "claude-3-5-sonnet-latest"
        mock_anthropic.AsyncAnthropic = Mock(return_value=Mock())
        
        generator = DiagramGenerator(provider="anthropic")
        assert generator is not None
        assert generator.provider == "anthropic"

    @patch('app.providers.diagram_generator.tempfile')
    @patch('app.providers.diagram_generator.subprocess')
    @patch('app.providers.diagram_generator.AsyncOpenAI')
    @patch('app.providers.diagram_generator.settings')
    @pytest.mark.asyncio
    async def test_generate_mermaid(self, mock_settings, mock_openai_class, mock_subprocess, mock_tempfile):
        """Test Mermaid diagram generation."""
        mock_settings.OPENAI_API_KEY = "test-key"
        mock_settings.OPENAI_MODEL = "gpt-4o-mini"
        
        # Mock OpenAI response
        mock_client = Mock()
        mock_message = Mock()
        mock_message.content = "graph TD\nA[Start] --> B[End]"
        mock_choice = Mock()
        mock_choice.message = mock_message
        mock_response = Mock()
        mock_response.choices = [mock_choice]
        mock_client.chat.completions.create = AsyncMock(return_value=mock_response)
        mock_openai_class.return_value = mock_client
        
        # Mock tempfile
        mock_tempfile.NamedTemporaryFile.return_value.__enter__.return_value.name = "/tmp/test.mmd"
        
        # Mock subprocess for rendering
        mock_subprocess.run.return_value = Mock(returncode=0)
        
        # Mock file operations
        mock_file_content = b"image_data"
        with patch('builtins.open', create=True) as mock_open, \
             patch('base64.b64encode', return_value=b'base64string'), \
             patch('os.path.exists', return_value=True), \
             patch('os.remove'):
            mock_file = Mock()
            mock_file.read.return_value = mock_file_content
            mock_open.return_value.__enter__.return_value = mock_file
            
            generator = DiagramGenerator(provider="openai")
            result = await generator.generate(
                description="A simple flowchart",
                diagram_type="mermaid",
                diagram_kind="flowchart",
                output_format="png"
            )
            
            assert "b64_json" in result
            assert "code" in result
            assert "format" in result


class TestVideoGenerator:
    """Tests for VideoGenerator."""

    @patch('app.providers.video_generator.settings')
    def test_init_replicate(self, mock_settings):
        """Test VideoGenerator initialization with Replicate."""
        mock_settings.REPLICATE_API_KEY = "test-key"
        mock_settings.REPLICATE_VIDEO_MODEL = "test-model"
        
        generator = VideoGenerator(provider="replicate")
        assert generator is not None
        assert generator.provider == "replicate"

    @patch('httpx.AsyncClient')
    @patch('app.providers.video_generator.replicate')
    @patch('app.providers.video_generator.settings')
    @pytest.mark.asyncio
    async def test_generate_from_image_url(self, mock_settings, mock_replicate, mock_httpx_client):
        """Test video generation from image URL."""
        mock_settings.REPLICATE_API_KEY = "test-key"
        mock_settings.REPLICATE_VIDEO_MODEL = "test-model"
        
        # Mock the client and its run method
        mock_client = Mock()
        mock_output = "https://example.com/video.mp4"
        mock_client.run = Mock(return_value=mock_output)
        mock_replicate.Client.return_value = mock_client
        
        # Mock httpx for video fetching
        mock_http_client = Mock()
        mock_http_client.__aenter__ = AsyncMock(return_value=mock_http_client)
        mock_http_client.__aexit__ = AsyncMock(return_value=None)
        mock_response = Mock()
        mock_response.content = b"video_data"
        mock_http_client.get = AsyncMock(return_value=mock_response)
        mock_httpx_client.return_value = mock_http_client
        
        generator = VideoGenerator(provider="replicate")
        result = await generator.generate_from_image(
            image_url="https://example.com/image.png",
            motion_bucket_id=127,
            num_frames=25
        )
        
        assert "video_url" in result or "video_b64" in result
        assert "format" in result
