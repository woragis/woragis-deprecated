"""
Integration tests for API endpoints.
These tests may require API keys and external services.
"""
import pytest
from fastapi.testclient import TestClient

from app.main import app


@pytest.mark.integration
class TestAPIIntegration:
    """Integration tests for API endpoints."""

    @pytest.fixture
    def client(self):
        """FastAPI test client."""
        return TestClient(app)

    @pytest.mark.requires_api_key
    def test_chat_endpoint_real_provider(self, client):
        """Test chat endpoint with real provider (requires API key)."""
        # Skip if no API key configured
        import os
        if not os.getenv("OPENAI_API_KEY"):
            pytest.skip("OPENAI_API_KEY not configured")
        
        request_data = {
            "agent": "startup",
            "input": "What is a startup?",
            "provider": "openai"
        }
        
        response = client.post("/v1/chat", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert "agent" in data
        assert "output" in data
        assert len(data["output"]) > 0

    @pytest.mark.requires_api_key
    def test_chat_stream_endpoint_real_provider(self, client):
        """Test streaming chat endpoint with real provider."""
        import os
        if not os.getenv("OPENAI_API_KEY"):
            pytest.skip("OPENAI_API_KEY not configured")
        
        request_data = {
            "agent": "startup",
            "input": "What is a startup?",
            "provider": "openai"
        }
        
        response = client.post("/v1/chat/stream", json=request_data)
        assert response.status_code == 200
        assert response.headers["content-type"] == "application/x-ndjson"
        
        # Read streaming response
        content = b""
        for chunk in response.iter_bytes():
            content += chunk
        
        assert len(content) > 0

    def test_health_check_integration(self, client):
        """Test health check endpoint."""
        response = client.get("/healthz")
        assert response.status_code == 200
        data = response.json()
        assert data["status"] in ["healthy", "degraded", "unhealthy"]
        assert "checks" in data

    def test_list_agents_integration(self, client):
        """Test listing agents."""
        response = client.get("/v1/agents")
        assert response.status_code == 200
        agents = response.json()
        assert isinstance(agents, list)
        assert len(agents) >= 4  # At least 4 default agents
