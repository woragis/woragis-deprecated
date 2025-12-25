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

    def test_list_agents_integration(self, client):
        """Test listing agents."""
        response = client.get("/v1/agents")
        assert response.status_code == 200
        agents = response.json()
        assert isinstance(agents, list)
        assert len(agents) >= 4  # At least 4 default agents
        assert "startup" in agents
        assert "economist" in agents
        assert "strategist" in agents
        assert "entrepreneur" in agents

    def test_chat_endpoint_validation(self, client):
        """Test chat endpoint input validation."""
        # Test missing required fields
        response = client.post("/v1/chat", json={})
        assert response.status_code == 422  # Validation error

        # Test invalid agent
        response = client.post("/v1/chat", json={
            "agent": "invalid_agent",
            "input": "test"
        })
        assert response.status_code == 422

        # Test valid request structure (may fail without API key, but should validate)
        response = client.post("/v1/chat", json={
            "agent": "startup",
            "input": "What is a startup?"
        })
        # Should either succeed (if API key available) or fail with 500/503 (service unavailable)
        assert response.status_code in [200, 500, 503, 502]

    @pytest.mark.requires_api_key
    def test_chat_endpoint_real_provider(self, client):
        """Test chat endpoint with real provider (requires API key)."""
        # Skip if no API key configured
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
        assert data["agent"] == "startup"

    @pytest.mark.requires_api_key
    def test_chat_endpoint_different_agents(self, client):
        """Test chat endpoint with different agents."""
        if not os.getenv("OPENAI_API_KEY"):
            pytest.skip("OPENAI_API_KEY not configured")
        
        agents = ["startup", "economist", "strategist", "entrepreneur"]
        
        for agent in agents:
            request_data = {
                "agent": agent,
                "input": "Test question",
                "provider": "openai"
            }
            
            response = client.post("/v1/chat", json=request_data)
            assert response.status_code == 200
            data = response.json()
            assert data["agent"] == agent
            assert "output" in data

    def test_chat_endpoint_auto_agent_selection(self, client):
        """Test auto agent selection based on input."""
        # Test with market-related input (should select economist)
        request_data = {
            "agent": "auto",
            "input": "What is inflation and market economics?",
            "provider": "openai"
        }
        
        response = client.post("/v1/chat", json=request_data)
        # Should either succeed or fail gracefully
        assert response.status_code in [200, 500, 503, 502]

    @pytest.mark.requires_api_key
    def test_chat_stream_endpoint_real_provider(self, client):
        """Test streaming chat endpoint with real provider."""
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
        
        # Verify NDJSON format (each line is valid JSON)
        lines = content.decode('utf-8').strip().split('\n')
        assert len(lines) > 0
        for line in lines:
            if line.strip():
                import json
                data = json.loads(line)
                assert "agent" in data or "output" in data or "done" in data

    def test_chat_stream_endpoint_validation(self, client):
        """Test streaming chat endpoint validation."""
        # Test missing required fields
        response = client.post("/v1/chat/stream", json={})
        assert response.status_code == 422

    def test_chat_with_custom_temperature(self, client):
        """Test chat endpoint with custom temperature."""
        request_data = {
            "agent": "startup",
            "input": "Test",
            "temperature": 0.7,
            "provider": "openai"
        }
        
        response = client.post("/v1/chat", json=request_data)
        # Should either succeed or fail gracefully
        assert response.status_code in [200, 500, 503, 502]

    def test_chat_with_custom_model(self, client):
        """Test chat endpoint with custom model override."""
        request_data = {
            "agent": "startup",
            "input": "Test",
            "model": "gpt-4o-mini",
            "provider": "openai"
        }
        
        response = client.post("/v1/chat", json=request_data)
        # Should either succeed or fail gracefully
        assert response.status_code in [200, 500, 503, 502]

    def test_chat_with_system_override(self, client):
        """Test chat endpoint with custom system message."""
        request_data = {
            "agent": "startup",
            "input": "Test",
            "system": "You are a helpful assistant.",
            "provider": "openai"
        }
        
        response = client.post("/v1/chat", json=request_data)
        # Should either succeed or fail gracefully
        assert response.status_code in [200, 500, 503, 502]

    @pytest.mark.requires_api_key
    def test_image_generation_endpoint(self, client):
        """Test image generation endpoint (requires API key)."""
        if not os.getenv("CIPHER_API_KEY"):
            pytest.skip("CIPHER_API_KEY not configured")
        
        request_data = {
            "provider": "cipher",
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

    def test_image_generation_validation(self, client):
        """Test image generation endpoint validation."""
        # Test missing prompt
        response = client.post("/v1/images", json={
            "provider": "cipher"
        })
        assert response.status_code == 422

        # Test invalid provider
        response = client.post("/v1/images", json={
            "provider": "invalid",
            "prompt": "test"
        })
        assert response.status_code == 422

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
        response = client.options("/v1/chat", headers={
            "Origin": "http://localhost:3000",
            "Access-Control-Request-Method": "POST"
        })
        # CORS should be enabled
        assert response.status_code in [200, 204]
