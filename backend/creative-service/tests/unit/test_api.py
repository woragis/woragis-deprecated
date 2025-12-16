"""
Unit tests for API endpoints.
"""
import pytest
from unittest.mock import Mock, AsyncMock, patch
from fastapi.testclient import TestClient

from app.main import app


class TestHealthCheck:
    """Tests for GET /healthz endpoint."""

    def test_health_check(self, client):
        """Test health check endpoint."""
        response = client.get("/healthz")
        assert response.status_code == 200
        data = response.json()
        assert "status" in data
        assert "checks" in data


class TestListProviders:
    """Tests for provider listing endpoints."""

    def test_list_image_providers(self, client):
        """Test listing image providers."""
        response = client.get("/v1/providers/images")
        assert response.status_code == 200
        data = response.json()
        assert "providers" in data
        assert "descriptions" in data
        assert "openai" in data["providers"]
        assert "stable-diffusion" in data["providers"]
        assert "cipher" in data["providers"]

    def test_list_diagram_providers(self, client):
        """Test listing diagram providers."""
        response = client.get("/v1/providers/diagrams")
        assert response.status_code == 200
        data = response.json()
        assert "providers" in data
        assert "ai_providers" in data
        assert "mermaid" in data["providers"]
        assert "graphviz" in data["providers"]

    def test_list_video_providers(self, client):
        """Test listing video providers."""
        response = client.get("/v1/providers/videos")
        assert response.status_code == 200
        data = response.json()
        assert "providers" in data
        assert "descriptions" in data
        assert "replicate" in data["providers"]
        assert "runway" in data["providers"]


class TestImageGeneration:
    """Tests for POST /v1/images/generate endpoint."""

    @patch('app.main.ImageProviderFactory')
    def test_generate_images_openai(self, mock_factory, client, mock_image_provider, sample_image_request):
        """Test image generation with OpenAI provider."""
        mock_factory.create.return_value = mock_image_provider
        
        response = client.post("/v1/images/generate", json=sample_image_request)
        assert response.status_code == 200
        data = response.json()
        assert "data" in data
        assert "provider" in data
        assert "prompt" in data
        assert len(data["data"]) > 0

    @patch('app.main.ImageProviderFactory')
    def test_generate_images_stable_diffusion(self, mock_factory, client, mock_image_provider):
        """Test image generation with Stable Diffusion provider."""
        mock_factory.create.return_value = mock_image_provider
        
        request_data = {
            "provider": "stable-diffusion",
            "prompt": "A technical architecture diagram",
            "context": "microservice architecture"
        }
        
        response = client.post("/v1/images/generate", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert data["provider"] == "stable-diffusion"

    @patch('app.main.ImageProviderFactory')
    def test_generate_images_cipher(self, mock_factory, client, mock_image_provider):
        """Test image generation with Cipher provider."""
        mock_factory.create.return_value = mock_image_provider
        
        request_data = {
            "provider": "cipher",
            "prompt": "A test image"
        }
        
        response = client.post("/v1/images/generate", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert data["provider"] == "cipher"

    @patch('app.main.ImageProviderFactory')
    def test_generate_images_with_enhanced_prompt(self, mock_factory, client, mock_image_provider):
        """Test image generation with enhanced prompt for OpenAI."""
        mock_factory.create.return_value = mock_image_provider
        
        request_data = {
            "provider": "openai",
            "prompt": "A technical diagram",
            "style": "technical",
            "context": "System architecture"
        }
        
        response = client.post("/v1/images/generate", json=request_data)
        assert response.status_code == 200
        # Verify enhanced prompt was called
        mock_image_provider.generate_with_enhanced_prompt.assert_called_once()

    @patch('app.main.ImageProviderFactory')
    def test_generate_images_error(self, mock_factory, client):
        """Test image generation error handling."""
        mock_provider = Mock()
        mock_provider.generate = AsyncMock(side_effect=Exception("Generation failed"))
        mock_factory.create.return_value = mock_provider
        
        request_data = {
            "provider": "openai",
            "prompt": "Test prompt"
        }
        
        response = client.post("/v1/images/generate", json=request_data)
        assert response.status_code == 500

    @patch('app.main.ImageProviderFactory')
    def test_generate_thumbnail(self, mock_factory, client, mock_image_provider):
        """Test thumbnail generation endpoint."""
        mock_factory.create.return_value = mock_image_provider
        
        request_data = {
            "provider": "openai",
            "prompt": "A thumbnail image"
        }
        
        response = client.post("/v1/images/generate/thumbnail", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert "data" in data


class TestDiagramGeneration:
    """Tests for POST /v1/diagrams/generate endpoint."""

    @patch('app.main.DiagramProviderFactory')
    def test_generate_diagram_mermaid(self, mock_factory, client, mock_diagram_generator, sample_diagram_request):
        """Test Mermaid diagram generation."""
        mock_factory.create.return_value = mock_diagram_generator
        
        response = client.post("/v1/diagrams/generate", json=sample_diagram_request)
        assert response.status_code == 200
        data = response.json()
        assert "b64_json" in data
        assert "code" in data
        assert "format" in data
        assert "diagram_type" in data
        assert data["diagram_type"] == "mermaid"

    @patch('app.main.DiagramProviderFactory')
    def test_generate_diagram_graphviz(self, mock_factory, client, mock_diagram_generator):
        """Test Graphviz diagram generation."""
        mock_factory.create.return_value = mock_diagram_generator
        
        request_data = {
            "description": "A network graph",
            "diagram_type": "graphviz",
            "output_format": "png",
            "ai_provider": "openai"
        }
        
        response = client.post("/v1/diagrams/generate", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert data["diagram_type"] == "graphviz"

    @patch('app.main.DiagramProviderFactory')
    def test_generate_mermaid_diagram_endpoint(self, mock_factory, client, mock_diagram_generator):
        """Test Mermaid diagram generation via dedicated endpoint."""
        mock_factory.create.return_value = mock_diagram_generator
        
        response = client.post(
            "/v1/diagrams/mermaid",
            params={
                "description": "A simple flowchart",
                "diagram_kind": "flowchart",
                "output_format": "png",
                "ai_provider": "openai"
            }
        )
        assert response.status_code == 200
        data = response.json()
        assert data["diagram_type"] == "mermaid"

    @patch('app.main.DiagramProviderFactory')
    def test_generate_diagram_error(self, mock_factory, client):
        """Test diagram generation error handling."""
        mock_generator = Mock()
        mock_generator.generate = AsyncMock(side_effect=Exception("Generation failed"))
        mock_factory.create.return_value = mock_generator
        
        request_data = {
            "description": "A diagram",
            "diagram_type": "mermaid",
            "output_format": "png",
            "ai_provider": "openai"
        }
        
        response = client.post("/v1/diagrams/generate", json=request_data)
        assert response.status_code == 500


class TestVideoGeneration:
    """Tests for POST /v1/videos/generate endpoint."""

    @patch('app.main.VideoProviderFactory')
    def test_generate_video_with_url(self, mock_factory, client, mock_video_generator, sample_video_request):
        """Test video generation with image URL."""
        mock_factory.create.return_value = mock_video_generator
        
        response = client.post("/v1/videos/generate", json=sample_video_request)
        assert response.status_code == 200
        data = response.json()
        assert "video_url" in data or "video_b64" in data
        assert "format" in data

    @patch('app.main.VideoProviderFactory')
    def test_generate_video_with_b64(self, mock_factory, client, mock_video_generator):
        """Test video generation with base64 image."""
        mock_factory.create.return_value = mock_video_generator
        
        request_data = {
            "image_b64": "base64encodedstring",
            "provider": "replicate"
        }
        
        response = client.post("/v1/videos/generate", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert "format" in data

    def test_generate_video_missing_image(self, client):
        """Test video generation without image URL or base64."""
        request_data = {
            "provider": "replicate"
        }
        
        response = client.post("/v1/videos/generate", json=request_data)
        assert response.status_code == 400

    @patch('app.main.VideoProviderFactory')
    def test_generate_video_replicate(self, mock_factory, client, mock_video_generator):
        """Test video generation with Replicate provider."""
        mock_factory.create.return_value = mock_video_generator
        
        request_data = {
            "image_url": "https://example.com/image.png",
            "provider": "replicate"
        }
        
        response = client.post("/v1/videos/generate", json=request_data)
        assert response.status_code == 200
        mock_factory.create.assert_called_once_with("replicate")

    @patch('app.main.VideoProviderFactory')
    def test_generate_video_runway(self, mock_factory, client, mock_video_generator):
        """Test video generation with Runway provider."""
        mock_factory.create.return_value = mock_video_generator
        
        request_data = {
            "image_url": "https://example.com/image.png",
            "provider": "runway"
        }
        
        response = client.post("/v1/videos/generate", json=request_data)
        assert response.status_code == 200
        mock_factory.create.assert_called_once_with("runway")

    @patch('app.main.VideoProviderFactory')
    def test_animate_image_endpoint(self, mock_factory, client, mock_video_generator):
        """Test animate image endpoint."""
        mock_factory.create.return_value = mock_video_generator
        
        request_data = {
            "image_url": "https://example.com/image.png",
            "provider": "replicate"
        }
        
        response = client.post("/v1/videos/animate", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert "format" in data

    @patch('app.main.VideoProviderFactory')
    def test_generate_video_error(self, mock_factory, client):
        """Test video generation error handling."""
        mock_generator = Mock()
        mock_generator.generate_from_image = AsyncMock(side_effect=Exception("Generation failed"))
        mock_factory.create.return_value = mock_generator
        
        request_data = {
            "image_url": "https://example.com/image.png",
            "provider": "replicate"
        }
        
        response = client.post("/v1/videos/generate", json=request_data)
        assert response.status_code == 500
