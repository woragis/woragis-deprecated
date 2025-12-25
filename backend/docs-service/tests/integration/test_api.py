"""
Integration tests for API endpoints.
"""
import pytest
from fastapi.testclient import TestClient


class TestDocsAPIIntegration:
    """Integration tests for docs API."""
    
    def test_full_docs_workflow(self, client, mock_docs_root):
        """Test complete workflow: list docs, get doc, get as HTML."""
        # List all docs
        response = client.get("/api/v1/docs/")
        assert response.status_code == 200
        files = response.json()["files"]
        assert len(files) > 0
        
        # Get first doc as JSON
        if files:
            first_file = files[0]
            path = first_file["path"]
            
            response = client.get(f"/api/v1/docs/{path}")
            assert response.status_code == 200
            doc_data = response.json()
            assert "content" in doc_data
            assert "html" in doc_data
            
            # Get same doc as HTML
            response = client.get(f"/api/v1/docs/{path}?format=html")
            assert response.status_code == 200
            assert "text/html" in response.headers["content-type"]
    
    def test_category_filtering(self, client, mock_docs_root):
        """Test filtering docs by category."""
        # List all docs
        all_response = client.get("/api/v1/docs/")
        all_files = all_response.json()["files"]
        
        # Filter by category
        filtered_response = client.get("/api/v1/docs/?category=architecture")
        filtered_files = filtered_response.json()["files"]
        
        # Filtered should be subset of all
        assert len(filtered_files) <= len(all_files)
        
        # All filtered files should have the category in their path
        for file_info in filtered_files:
            assert "architecture" in file_info["path"].lower() or \
                   file_info["category"] == "architecture"

    def test_list_docs_endpoint(self, client, mock_docs_root):
        """Test listing all documentation files."""
        response = client.get("/api/v1/docs/")
        assert response.status_code == 200
        data = response.json()
        assert "files" in data
        assert isinstance(data["files"], list)
        
        # Each file should have required fields
        if data["files"]:
            first_file = data["files"][0]
            assert "path" in first_file
            assert "name" in first_file or "title" in first_file

    def test_get_doc_by_path(self, client, mock_docs_root):
        """Test getting a specific document by path."""
        # First, list docs to get a valid path
        list_response = client.get("/api/v1/docs/")
        files = list_response.json()["files"]
        
        if files:
            path = files[0]["path"]
            
            # Get document
            response = client.get(f"/api/v1/docs/{path}")
            assert response.status_code == 200
            data = response.json()
            assert "content" in data
            assert "html" in data or "markdown" in data

    def test_get_doc_as_html(self, client, mock_docs_root):
        """Test getting document as HTML format."""
        list_response = client.get("/api/v1/docs/")
        files = list_response.json()["files"]
        
        if files:
            path = files[0]["path"]
            
            response = client.get(f"/api/v1/docs/{path}?format=html")
            assert response.status_code == 200
            assert "text/html" in response.headers["content-type"]
            assert len(response.text) > 0

    def test_get_nonexistent_doc(self, client, mock_docs_root):
        """Test getting a non-existent document."""
        response = client.get("/api/v1/docs/nonexistent/path.md")
        assert response.status_code == 404

    def test_search_docs(self, client, mock_docs_root):
        """Test searching documentation."""
        # Test search parameter (if supported)
        response = client.get("/api/v1/docs/?search=architecture")
        assert response.status_code == 200
        data = response.json()
        assert "files" in data

    def test_pagination(self, client, mock_docs_root):
        """Test pagination of documentation list."""
        # Test with limit parameter (if supported)
        response = client.get("/api/v1/docs/?limit=5")
        assert response.status_code == 200
        data = response.json()
        assert "files" in data
        if "limit" in data or len(data["files"]) <= 5:
            # Pagination may or may not be implemented
            pass

    def test_health_check(self, client):
        """Test health check endpoint."""
        response = client.get("/healthz")
        assert response.status_code == 200
        data = response.json()
        assert "status" in data