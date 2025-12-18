"""
Unit tests for API endpoints.
"""
import pytest
from fastapi.testclient import TestClient
from unittest.mock import patch, MagicMock
from pathlib import Path


class TestRootEndpoint:
    """Tests for GET / endpoint."""
    
    def test_root_endpoint(self, client):
        """Test root endpoint returns service info."""
        response = client.get("/")
        assert response.status_code == 200
        data = response.json()
        assert "service" in data
        assert "version" in data
        assert "endpoints" in data
        assert data["service"] == "woragis-docs-service"


class TestHealthCheckEndpoint:
    """Tests for GET /healthz endpoint."""
    
    def test_health_check_endpoint(self, client):
        """Test health check endpoint."""
        response = client.get("/healthz")
        assert response.status_code in [200, 503]  # Can be healthy or unhealthy
        data = response.json()
        assert "status" in data
        assert "service" in data
        assert "checks" in data
        assert data["service"] == "docs-service"
    
    def test_health_check_structure(self, client):
        """Test health check response structure."""
        response = client.get("/healthz")
        data = response.json()
        
        # Check structure
        assert isinstance(data["checks"], list)
        assert len(data["checks"]) > 0
        
        # Check that each check has name and status
        for check in data["checks"]:
            assert "name" in check
            assert "status" in check


class TestDocsEndpoints:
    """Tests for docs API endpoints."""
    
    def test_list_docs_endpoint(self, client):
        """Test listing documentation files."""
        response = client.get("/api/v1/docs/")
        assert response.status_code == 200
        data = response.json()
        assert "files" in data
        assert "total" in data
        assert isinstance(data["files"], list)
        assert isinstance(data["total"], int)
    
    def test_list_docs_with_category_filter(self, client):
        """Test listing docs with category filter."""
        response = client.get("/api/v1/docs/?category=architecture")
        assert response.status_code == 200
        data = response.json()
        assert "files" in data
    
    def test_get_doc_not_found(self, client):
        """Test getting a non-existent doc file."""
        response = client.get("/api/v1/docs/nonexistent.md")
        assert response.status_code == 404
    
    def test_get_doc_invalid_path(self, client):
        """Test getting doc with invalid path (directory traversal)."""
        response = client.get("/api/v1/docs/../../../etc/passwd")
        # Should return 404 (not found) or 400 (bad request)
        assert response.status_code in [400, 404]
    
    def test_get_doc_json_format(self, client, mock_docs_root):
        """Test getting doc in JSON format."""
        response = client.get("/api/v1/docs/test.md?format=json")
        assert response.status_code == 200
        data = response.json()
        assert "path" in data
        assert "title" in data
        assert "content" in data
        assert "html" in data
    
    def test_get_doc_html_format(self, client, mock_docs_root):
        """Test getting doc in HTML format."""
        response = client.get("/api/v1/docs/test.md?format=html")
        assert response.status_code == 200
        assert "text/html" in response.headers["content-type"]
        assert b"<!DOCTYPE html>" in response.content


class TestMetricsEndpoint:
    """Tests for GET /metrics endpoint."""
    
    def test_metrics_endpoint(self, client):
        """Test Prometheus metrics endpoint."""
        response = client.get("/metrics")
        assert response.status_code == 200
        assert "text/plain" in response.headers["content-type"]
        # Prometheus metrics format
        content = response.text
        assert "# HELP" in content or "http" in content.lower()
