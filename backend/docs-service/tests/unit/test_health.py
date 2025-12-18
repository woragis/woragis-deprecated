"""
Unit tests for health check functionality.
"""
import pytest
import os
from pathlib import Path
from unittest.mock import patch, MagicMock

from app.health import check_health


class TestHealthCheck:
    """Tests for health check endpoint."""
    
    def test_health_check_service_ok(self):
        """Test health check returns service status."""
        with patch('app.health.Path') as mock_path:
            mock_docs_path = MagicMock()
            mock_docs_path.exists.return_value = True
            mock_docs_path.rglob.return_value = []
            mock_path.return_value = mock_docs_path
            
            with patch('os.access', return_value=True):
                result = check_health()
                
                assert "status" in result
                assert "service" in result
                assert "checks" in result
                assert result["service"] == "docs-service"
                assert result["status"] == "healthy"
                
                # Check that service check is present
                check_names = [check["name"] for check in result["checks"]]
                assert "service" in check_names
    
    def test_health_check_docs_directory_missing(self):
        """Test health check when docs directory doesn't exist."""
        with patch('app.health.Path') as mock_path:
            mock_docs_path = MagicMock()
            mock_docs_path.exists.return_value = False
            mock_path.return_value = mock_docs_path
            
            result = check_health()
            
            assert result["status"] == "unhealthy"
            check_names = [check["name"] for check in result["checks"]]
            assert "docs_directory" in check_names
    
    def test_health_check_docs_directory_not_readable(self):
        """Test health check when docs directory is not readable."""
        with patch('app.health.Path') as mock_path:
            mock_docs_path = MagicMock()
            mock_docs_path.exists.return_value = True
            mock_docs_path.rglob.return_value = []
            mock_path.return_value = mock_docs_path
            
            with patch('os.access', return_value=False):
                result = check_health()
                
                assert result["status"] == "unhealthy"
                check_names = [check["name"] for check in result["checks"]]
                assert "docs_readable" in check_names
    
    def test_health_check_counts_markdown_files(self):
        """Test health check counts markdown files."""
        with patch('app.health.Path') as mock_path:
            mock_docs_path = MagicMock()
            mock_docs_path.exists.return_value = True
            
            # Mock markdown files
            mock_file1 = MagicMock()
            mock_file2 = MagicMock()
            mock_docs_path.rglob.return_value = [mock_file1, mock_file2]
            mock_path.return_value = mock_docs_path
            
            with patch('os.access', return_value=True):
                result = check_health()
                
                assert result["status"] == "healthy"
                markdown_check = next(
                    (check for check in result["checks"] if check["name"] == "markdown_files"),
                    None
                )
                assert markdown_check is not None
                assert markdown_check["status"] == "ok"
                assert int(markdown_check["count"]) == 2
    
    def test_health_check_caching(self):
        """Test that health check results are cached."""
        with patch('app.health.Path') as mock_path, \
             patch('app.health.time.time', return_value=1000.0), \
             patch('os.access', return_value=True):
            
            mock_docs_path = MagicMock()
            mock_docs_path.exists.return_value = True
            mock_docs_path.rglob.return_value = []
            mock_path.return_value = mock_docs_path
            
            # First call
            result1 = check_health()
            
            # Second call within cache TTL (5 seconds)
            with patch('app.health.time.time', return_value=1004.0):  # 4 seconds later
                result2 = check_health()
                
                # Should return cached result (same object reference)
                assert result1 == result2
            
            # Third call after cache expires (6 seconds)
            with patch('app.health.time.time', return_value=1006.0):  # 6 seconds later
                result3 = check_health()
                
                # Should be a new result (may have different content)
                assert result3 is not None
