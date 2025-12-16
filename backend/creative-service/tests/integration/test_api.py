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
