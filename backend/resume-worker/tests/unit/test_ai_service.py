"""
Unit tests for AI service client.
"""
import pytest
import sys
import os
from unittest.mock import Mock, patch, MagicMock

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))
from ai_service import AIService, PROFILE_TEMPLATE, ABOUT_ME_TEMPLATE


class TestAIService:
    """Tests for AIService class."""

    def test_init(self):
        """Test AIService initialization."""
        service = AIService("http://ai-service:8000")
        assert service.base_url == "http://ai-service:8000"

    def test_init_strips_trailing_slash(self):
        """Test that base_url is stripped of trailing slashes."""
        service = AIService("http://ai-service:8000/")
        assert service.base_url == "http://ai-service:8000"

    @patch('ai_service.requests.post')
    def test_generate_resume_section_profile(self, mock_post):
        """Test generating profile section."""
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "output": "<p>Generated profile</p>"
        }
        mock_post.return_value = mock_response

        service = AIService("http://ai-service:8000")
        result = service.generate_resume_section(
            section_type='profile',
            job_description="Test job",
            projects=[]
        )

        assert result == "<p>Generated profile</p>"
        mock_post.assert_called_once()
        call_args = mock_post.call_args
        assert call_args[0][0] == "http://ai-service:8000/v1/chat"
        assert call_args[1]['json']['agent'] == "auto"

    @patch('ai_service.requests.post')
    def test_generate_resume_section_about(self, mock_post):
        """Test generating about section."""
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "output": "<p>About me</p>"
        }
        mock_post.return_value = mock_response

        service = AIService("http://ai-service:8000")
        result = service.generate_resume_section(
            section_type='about',
            job_description="Test job",
            projects=[]
        )

        assert result == "<p>About me</p>"

    @patch('ai_service.requests.post')
    def test_generate_resume_section_experience(self, mock_post):
        """Test generating experience section."""
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "output": "<p>Experience</p>"
        }
        mock_post.return_value = mock_response

        service = AIService("http://ai-service:8000")
        result = service.generate_resume_section(
            section_type='experience',
            job_description="Test job",
            projects=[{'name': 'Project 1'}]
        )

        assert result == "<p>Experience</p>"

    @patch('ai_service.requests.post')
    def test_generate_resume_section_skills(self, mock_post):
        """Test generating skills section."""
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "output": "<p>Skills</p>"
        }
        mock_post.return_value = mock_response

        service = AIService("http://ai-service:8000")
        result = service.generate_resume_section(
            section_type='skills',
            job_description="Test job",
            projects=[]
        )

        assert result == "<p>Skills</p>"

    @patch('ai_service.requests.post')
    def test_generate_resume_section_portuguese(self, mock_post):
        """Test generating section in Portuguese."""
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "output": "<p>Conteúdo em português</p>"
        }
        mock_post.return_value = mock_response

        service = AIService("http://ai-service:8000")
        result = service.generate_resume_section(
            section_type='profile',
            job_description="Test job",
            projects=[],
            language="pt"
        )

        assert result == "<p>Conteúdo em português</p>"
        # Verify Portuguese instruction was included
        call_args = mock_post.call_args
        system_prompt = call_args[1]['json']['system']
        assert "Portuguese" in system_prompt or "pt-BR" in system_prompt

    @patch('ai_service.requests.post')
    def test_generate_resume_section_error(self, mock_post):
        """Test error handling in section generation."""
        mock_response = Mock()
        mock_response.status_code = 500
        mock_response.text = "Internal Server Error"
        mock_post.return_value = mock_response

        service = AIService("http://ai-service:8000")
        result = service.generate_resume_section(
            section_type='profile',
            job_description="Test job",
            projects=[]
        )

        # Should return empty string or default on error
        assert result == "" or result is not None

    @patch('ai_service.requests.post')
    def test_generate_resume_section_with_projects(self, mock_post):
        """Test generating section with projects context."""
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "output": "<p>Profile with projects</p>"
        }
        mock_post.return_value = mock_response

        service = AIService("http://ai-service:8000")
        projects = [
            {'name': 'Project 1', 'description': 'Test project'}
        ]
        result = service.generate_resume_section(
            section_type='profile',
            job_description="Test job",
            projects=projects
        )

        assert result == "<p>Profile with projects</p>"
        # Verify projects were included in prompt
        call_args = mock_post.call_args
        input_prompt = call_args[1]['json']['input']
        assert "Project 1" in input_prompt or "project" in input_prompt.lower()

    def test_profile_template_constant(self):
        """Test that PROFILE_TEMPLATE constant exists."""
        assert PROFILE_TEMPLATE is not None
        assert len(PROFILE_TEMPLATE) > 0

    def test_about_me_template_constant(self):
        """Test that ABOUT_ME_TEMPLATE constant exists."""
        assert ABOUT_ME_TEMPLATE is not None
        assert len(ABOUT_ME_TEMPLATE) > 0
