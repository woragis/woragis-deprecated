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
