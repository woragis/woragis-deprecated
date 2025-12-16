"""
Integration tests for resume worker.
These tests may require external services (database, RabbitMQ, AI service).
"""
import pytest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))


@pytest.mark.integration
class TestWorkerIntegration:
    """Integration tests for resume worker."""

    def test_health_check_integration(self):
        """Test health check endpoint."""
        from health import check_health
        result = check_health()
        assert "status" in result
        assert "checks" in result
        assert result["status"] in ["healthy", "degraded", "unhealthy"]
