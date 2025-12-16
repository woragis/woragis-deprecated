"""
Pytest configuration and shared fixtures.
"""
import os
import pytest
from unittest.mock import Mock, MagicMock, patch
import sys

# Add src directory to path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))


@pytest.fixture
def mock_database():
    """Mock database connection."""
    db = Mock()
    db.connect = Mock()
    db.close = Mock()
    db.get_user_by_id = Mock(return_value={
        'id': 'user123',
        'name': 'Test User',
        'email': 'test@example.com'
    })
    db.get_projects_by_user = Mock(return_value=[])
    db.get_certifications_by_user = Mock(return_value=[])
    return db


@pytest.fixture
def mock_ai_service():
    """Mock AI service."""
    service = Mock()
    service.generate_resume_section = Mock(return_value="<p>Generated content</p>")
    service.base_url = "http://ai-service:8000"
    return service


@pytest.fixture
def mock_translation_helper():
    """Mock translation helper."""
    helper = Mock()
    helper.connect = Mock()
    helper.close = Mock()
    helper.get_translation = Mock(return_value=None)
    helper.translate_projects = Mock(side_effect=lambda projects, lang: projects)
    helper.translate_certifications = Mock(side_effect=lambda certs, lang: certs)
    helper.translate_posts = Mock(side_effect=lambda posts, lang: posts)
    helper.translate_experiences = Mock(side_effect=lambda exps, lang: exps)
    helper.translate_education_status = Mock(side_effect=lambda status, lang: status)
    helper.translate_section_header = Mock(side_effect=lambda header, lang: header)
    return helper


@pytest.fixture
def sample_job_description():
    """Sample job description for testing."""
    return """
    We are looking for a Senior Backend Developer with experience in:
    - Golang and Python
    - Microservices architecture
    - PostgreSQL and Redis
    - Docker and Kubernetes
    - RESTful APIs and gRPC
    """


@pytest.fixture
def sample_projects():
    """Sample projects data."""
    return [
        {
            'id': 'proj1',
            'name': 'Test Project',
            'description': 'A test project',
            'technologies': ['Golang', 'PostgreSQL']
        }
    ]


@pytest.fixture
def sample_certifications():
    """Sample certifications data."""
    return [
        {
            'id': 'cert1',
            'name': 'AWS Certified',
            'issuer': 'Amazon'
        }
    ]


@pytest.fixture(autouse=True)
def reset_env():
    """Reset environment variables before each test."""
    original_env = os.environ.copy()
    yield
    os.environ.clear()
    os.environ.update(original_env)
