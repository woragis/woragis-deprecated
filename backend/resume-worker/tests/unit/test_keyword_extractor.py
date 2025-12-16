"""
Unit tests for keyword extraction.
"""
import pytest
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))
from keyword_extractor import extract_keywords, TECH_CATEGORY_MAPPING, CERT_CATEGORY_MAPPING


class TestExtractKeywords:
    """Tests for extract_keywords function."""

    def test_extract_backend_keywords(self):
        """Test extraction of backend-related keywords."""
        job_desc = "Looking for a backend developer with Go and Python experience"
        result = extract_keywords(job_desc)
        
        assert 'backend' in result['tech_categories']
        assert 'Golang' in result['skill_names'] or 'Python' in result['skill_names']

    def test_extract_frontend_keywords(self):
        """Test extraction of frontend-related keywords."""
        job_desc = "Frontend developer needed with React and TypeScript"
        result = extract_keywords(job_desc)
        
        assert 'frontend' in result['tech_categories']
        assert 'TypeScript' in result['skill_names']

    def test_extract_devops_keywords(self):
        """Test extraction of DevOps-related keywords."""
        job_desc = "DevOps engineer with Docker and Kubernetes experience"
        result = extract_keywords(job_desc)
        
        assert 'devops' in result['tech_categories']
        assert 'Docker' in result['skill_names']
        assert 'Kubernetes' in result['skill_names']

    def test_extract_database_keywords(self):
        """Test extraction of database-related keywords."""
        job_desc = "Database administrator with PostgreSQL and Redis knowledge"
        result = extract_keywords(job_desc)
        
        assert 'database' in result['tech_categories']
        assert 'PostgreSQL' in result['skill_names']
        assert 'Redis' in result['skill_names']

    def test_extract_certification_keywords(self):
        """Test extraction of certification-related keywords."""
        job_desc = "AWS certified cloud architect with Kubernetes expertise"
        result = extract_keywords(job_desc)
        
        assert 'cloud' in result['cert_categories']
        assert 'devops' in result['cert_categories'] or 'devops' in result['tech_categories']

    def test_extract_multiple_categories(self):
        """Test extraction of multiple categories."""
        job_desc = "Full-stack developer with backend, frontend, and DevOps experience"
        result = extract_keywords(job_desc)
        
        assert 'backend' in result['tech_categories']
        assert 'frontend' in result['tech_categories']
        assert 'devops' in result['tech_categories']

    def test_extract_no_keywords(self):
        """Test extraction with no matching keywords."""
        job_desc = "General position with no specific technology requirements"
        result = extract_keywords(job_desc)
        
        assert result['tech_categories'] is None or len(result['tech_categories']) == 0
        assert result['cert_categories'] is None or len(result['cert_categories']) == 0
        assert result['skill_names'] is None or len(result['skill_names']) == 0

    def test_extract_case_insensitive(self):
        """Test that keyword extraction is case-insensitive."""
        job_desc = "GOLANG DEVELOPER WITH PYTHON"
        result = extract_keywords(job_desc)
        
        assert 'Golang' in result['skill_names']
        assert 'Python' in result['skill_names']

    def test_extract_golang_variations(self):
        """Test extraction of Golang keyword variations."""
        variations = [
            "golang developer",
            "go programming expert",
            "go language specialist"
        ]
        
        for job_desc in variations:
            result = extract_keywords(job_desc)
            assert result['skill_names'] is not None, f"skill_names should not be None for '{job_desc}'"
            assert 'Golang' in result['skill_names'], f"'Golang' should be in skill_names for '{job_desc}'"
            assert result['tech_categories'] is not None, f"tech_categories should not be None for '{job_desc}'"
            assert 'backend' in result['tech_categories'] or 'language' in result['tech_categories']

    def test_extract_grpc_websocket(self):
        """Test extraction of gRPC and WebSocket keywords."""
        job_desc = "API developer with gRPC and WebSocket experience"
        result = extract_keywords(job_desc)
        
        assert 'gRPC' in result['skill_names']
        assert 'WebSocket' in result['skill_names']
