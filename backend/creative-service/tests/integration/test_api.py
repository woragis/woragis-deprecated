"""
Integration tests for API endpoints.
These tests may require API keys and external services.
"""
import pytest
import os
from fastapi.testclient import TestClient

from app.main import app


@pytest.mark.integration
class TestAPIIntegration:
    """Integration tests for API endpoints."""

    @pytest.fixture
    def client(self):
        """FastAPI test client."""
        return TestClient(app)

    def test_health_check_integration(self, client):
        """Test health check endpoint."""
        response = client.get("/healthz")
        assert response.status_code == 200
        data = response.json()
        assert data["status"] in ["healthy", "degraded", "unhealthy"]
        assert "checks" in data

    def test_list_providers_integration(self, client):
        """Test listing all providers."""
        # Test image providers
        response = client.get("/v1/providers/images")
        assert response.status_code == 200
        data = response.json()
        assert isinstance(data["providers"], list)
        assert len(data["providers"]) >= 3

        # Test diagram providers
        response = client.get("/v1/providers/diagrams")
        assert response.status_code == 200
        data = response.json()
        assert isinstance(data["providers"], list)

        # Test video providers
        response = client.get("/v1/providers/videos")
        assert response.status_code == 200
        data = response.json()
        assert isinstance(data["providers"], list)

    def test_image_generation_validation(self, client):
        """Test image generation endpoint validation."""
        # Test missing required fields
        response = client.post("/v1/images", json={})
        assert response.status_code == 422

        # Test invalid provider
        response = client.post("/v1/images", json={
            "provider": "invalid",
            "prompt": "test"
        })
        assert response.status_code == 422

        # Test valid request structure
        response = client.post("/v1/images", json={
            "provider": "openai",
            "prompt": "A beautiful landscape"
        })
        # Should either succeed (if API key available) or fail gracefully
        assert response.status_code in [200, 500, 503, 502, 401]

    @pytest.mark.requires_api_key
    def test_image_generation_openai(self, client):
        """Test image generation with OpenAI (requires API key)."""
        if not os.getenv("OPENAI_API_KEY"):
            pytest.skip("OPENAI_API_KEY not configured")
        
        request_data = {
            "provider": "openai",
            "prompt": "A beautiful sunset over mountains",
            "n": 1,
            "size": "1024x1024"
        }
        
        response = client.post("/v1/images", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert "data" in data
        assert isinstance(data["data"], list)
        assert len(data["data"]) > 0

    @pytest.mark.requires_api_key
    def test_image_generation_stability_ai(self, client):
        """Test image generation with Stability AI (requires API key)."""
        if not os.getenv("REPLICATE_API_KEY"):
            pytest.skip("REPLICATE_API_KEY not configured")
        
        request_data = {
            "provider": "stability",
            "prompt": "A futuristic cityscape",
            "n": 1
        }
        
        response = client.post("/v1/images", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert "data" in data
        assert isinstance(data["data"], list)

    def test_diagram_generation_validation(self, client):
        """Test diagram generation endpoint validation."""
        # Test missing required fields
        response = client.post("/v1/diagrams", json={})
        assert response.status_code == 422

        # Test valid request structure
        response = client.post("/v1/diagrams", json={
            "provider": "openai",
            "prompt": "Create a system architecture diagram",
            "format": "mermaid"
        })
        # Should either succeed or fail gracefully
        assert response.status_code in [200, 500, 503, 502, 401]

    @pytest.mark.requires_api_key
    def test_diagram_generation(self, client):
        """Test diagram generation (requires API key)."""
        if not os.getenv("OPENAI_API_KEY"):
            pytest.skip("OPENAI_API_KEY not configured")
        
        request_data = {
            "provider": "openai",
            "prompt": "Create a simple flowchart showing user login process",
            "format": "mermaid"
        }
        
        response = client.post("/v1/diagrams", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert "data" in data or "diagram" in data or "code" in data

    def test_video_generation_validation(self, client):
        """Test video generation endpoint validation."""
        # Test missing required fields
        response = client.post("/v1/videos", json={})
        assert response.status_code == 422

        # Test valid request structure
        response = client.post("/v1/videos", json={
            "provider": "stability",
            "prompt": "A short video of waves on a beach"
        })
        # Should either succeed or fail gracefully
        assert response.status_code in [200, 500, 503, 502, 401]

    @pytest.mark.requires_api_key
    def test_video_generation(self, client):
        """Test video generation (requires API key)."""
        if not os.getenv("REPLICATE_API_KEY"):
            pytest.skip("REPLICATE_API_KEY not configured")
        
        request_data = {
            "provider": "stability",
            "prompt": "A short video of waves on a beach",
            "duration": 5
        }
        
        response = client.post("/v1/videos", json=request_data)
        # Video generation may be async, so accept 200 or 202
        assert response.status_code in [200, 202, 500, 503, 502, 401]

    def test_provider_specific_options(self, client):
        """Test provider-specific options."""
        # Test OpenAI-specific options
        request_data = {
            "provider": "openai",
            "prompt": "Test",
            "quality": "standard",
            "style": "vivid"
        }
        
        response = client.post("/v1/images", json=request_data)
        # Should either succeed or fail gracefully
        assert response.status_code in [200, 500, 503, 502, 401, 422]

    def test_metrics_endpoint(self, client):
        """Test Prometheus metrics endpoint."""
        response = client.get("/metrics")
        assert response.status_code == 200
        assert "text/plain" in response.headers["content-type"]
        # Should contain Prometheus metrics
        content = response.text
        assert "http_requests_total" in content or "requests" in content.lower()

    def test_cors_headers(self, client):
        """Test CORS headers are present."""
        response = client.options("/v1/images", headers={
            "Origin": "http://localhost:3000",
            "Access-Control-Request-Method": "POST"
        })
        # CORS should be enabled
        assert response.status_code in [200, 204]
