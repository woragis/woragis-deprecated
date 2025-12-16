"""
Unit tests for resume generator module.
"""
import pytest
import sys
import os
from unittest.mock import Mock, patch, MagicMock, mock_open
from datetime import datetime

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))
from resume_generator import ResumeGenerator, sanitize_html, ALLOWED_TAGS, ALLOWED_ATTRIBUTES


class TestSanitizeHTML:
    """Tests for sanitize_html function."""

    def test_sanitize_html_allowed_tags(self):
        """Test that allowed tags are preserved."""
        html = "<p>Test content</p><strong>Bold</strong>"
        result = sanitize_html(html)
        assert "<p>" in result
        assert "<strong>" in result

    def test_sanitize_html_removes_script(self):
        """Test that script tags are removed."""
        html = "<p>Test</p><script>alert('xss')</script>"
        result = sanitize_html(html)
        assert "<script>" not in result
        # Note: bleach removes script tags but may leave text content
        # The important thing is that script tags are removed
        assert "<p>Test</p>" in result

    def test_sanitize_html_empty(self):
        """Test sanitizing empty HTML."""
        assert sanitize_html("") == ""
        assert sanitize_html(None) == ""

    def test_sanitize_html_allowed_attributes(self):
        """Test that allowed attributes are preserved."""
        html = '<div class="test">Content</div>'
        result = sanitize_html(html)
        assert 'class="test"' in result


class TestResumeGenerator:
    """Tests for ResumeGenerator class."""

    @patch('resume_generator.os.path.exists')
    @patch('resume_generator.os.path.isdir')
    @patch('resume_generator.os.path.abspath')
    @patch('resume_generator.os.makedirs')
    def test_init(self, mock_makedirs, mock_abspath, mock_isdir, mock_exists):
        """Test ResumeGenerator initialization."""
        mock_exists.return_value = True
        mock_isdir.return_value = True
        mock_abspath.return_value = "/app/templates"

        mock_db = Mock()
        mock_ai = Mock()
        mock_translation = Mock()

        generator = ResumeGenerator(mock_db, mock_ai, "/app/output", mock_translation)

        assert generator.db == mock_db
        assert generator.ai_service == mock_ai
        assert generator.output_dir == "/app/output"
        assert generator.translation_helper == mock_translation
        mock_makedirs.assert_called_once_with("/app/output", exist_ok=True)

    @patch('resume_generator.os.path.exists')
    @patch('resume_generator.os.path.isdir')
    @patch('resume_generator.os.path.abspath')
    def test_init_template_dir_not_found(self, mock_abspath, mock_isdir, mock_exists):
        """Test initialization when template directory is not found."""
        mock_exists.return_value = False
        mock_isdir.return_value = False

        mock_db = Mock()
        mock_ai = Mock()

        with pytest.raises(FileNotFoundError):
            ResumeGenerator(mock_db, mock_ai, "/app/output")

    @patch('resume_generator.os.path.exists')
    @patch('resume_generator.os.path.isdir')
    @patch('resume_generator.os.path.abspath')
    @patch('resume_generator.os.makedirs')
    def test_parse_hard_skills(self, mock_makedirs, mock_abspath, mock_isdir, mock_exists):
        """Test parsing hard skills text."""
        mock_exists.return_value = True
        mock_isdir.return_value = True
        mock_abspath.return_value = "/app/templates"

        mock_db = Mock()
        mock_ai = Mock()
        generator = ResumeGenerator(mock_db, mock_ai, "/app/output")

        skills_text = "Backend\n• Golang • Python\nFrontend\n• React • Vue"
        result = generator._parse_hard_skills(skills_text)

        assert len(result) == 2
        assert result[0]['title'] == 'Backend'
        assert 'Golang' in result[0]['skills']
        assert result[1]['title'] == 'Frontend'

    @patch('resume_generator.os.path.exists')
    @patch('resume_generator.os.path.isdir')
    @patch('resume_generator.os.path.abspath')
    @patch('resume_generator.os.makedirs')
    def test_parse_hard_skills_empty(self, mock_makedirs, mock_abspath, mock_isdir, mock_exists):
        """Test parsing empty hard skills."""
        mock_exists.return_value = True
        mock_isdir.return_value = True
        mock_abspath.return_value = "/app/templates"

        mock_db = Mock()
        mock_ai = Mock()
        generator = ResumeGenerator(mock_db, mock_ai, "/app/output")

        result = generator._parse_hard_skills("")
        assert result == []

    @patch('resume_generator.os.path.exists')
    @patch('resume_generator.os.path.isdir')
    @patch('resume_generator.os.path.abspath')
    @patch('resume_generator.os.makedirs')
    @patch('resume_generator.os.path.getsize')
    @patch('resume_generator.HTML')
    @patch('resume_generator.CSS')
    @patch('resume_generator.datetime')
    def test_generate_resume(self, mock_datetime, mock_css, mock_html, mock_getsize,
                             mock_makedirs, mock_abspath, mock_isdir, mock_exists):
        """Test resume generation."""
        mock_exists.return_value = True
        mock_isdir.return_value = True
        mock_abspath.return_value = "/app/templates"
        mock_getsize.return_value = 12345
        mock_datetime.now.return_value.strftime = Mock(return_value="20240101_120000")
        mock_datetime.now.return_value.isoformat = Mock(return_value="2024-01-01T12:00:00")

        # Mock database
        mock_db = Mock()
        mock_db.get_user_info = Mock(return_value={
            'id': 'user123',
            'email': 'test@example.com',
            'name': 'Test User'
        })
        mock_db.get_user_projects = Mock(return_value=[])
        mock_db.get_user_certifications = Mock(return_value=[])
        mock_db.get_user_publications = Mock(return_value=[])
        mock_db.get_user_experiences = Mock(return_value=[])

        # Mock AI service
        mock_ai = Mock()
        mock_ai.generate_resume_section = Mock(return_value="<p>Generated content</p>")
        mock_ai.generate_tags = Mock(return_value=["tag1", "tag2"])

        # Mock translation helper
        mock_translation = Mock()
        mock_translation.translate_projects = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_certifications = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_posts = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_experiences = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_education_status = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_section_header = Mock(side_effect=lambda x, y: x)

        # Mock template rendering
        mock_template_env = Mock()
        mock_template = Mock()
        mock_template.render = Mock(return_value="<html>Resume HTML</html>")
        mock_template_env.get_template = Mock(return_value=mock_template)

        # Mock WeasyPrint HTML
        mock_html_instance = Mock()
        mock_html.return_value = mock_html_instance

        generator = ResumeGenerator(mock_db, mock_ai, "/app/output", mock_translation)
        generator.template_env = mock_template_env

        result = generator.generate_resume(
            user_id="user123",
            job_description="Backend developer with Golang",
            job_title="Software Engineer",
            language="en"
        )

        assert result['user_id'] == 'user123'
        assert result['job_title'] == 'Software Engineer'
        assert 'output_path' in result
        assert 'file_size' in result
        assert 'generated_at' in result
        assert 'tags' in result

        # Verify AI service was called for all sections
        assert mock_ai.generate_resume_section.call_count >= 5
        mock_ai.generate_tags.assert_called_once()

    @patch('resume_generator.os.path.exists')
    @patch('resume_generator.os.path.isdir')
    @patch('resume_generator.os.path.abspath')
    @patch('resume_generator.os.makedirs')
    def test_generate_resume_user_not_found(self, mock_makedirs, mock_abspath, mock_isdir, mock_exists):
        """Test resume generation when user is not found."""
        mock_exists.return_value = True
        mock_isdir.return_value = True
        mock_abspath.return_value = "/app/templates"

        mock_db = Mock()
        mock_db.get_user_info = Mock(return_value=None)

        mock_ai = Mock()
        generator = ResumeGenerator(mock_db, mock_ai, "/app/output")

        with pytest.raises(ValueError, match="User .* not found"):
            generator.generate_resume(
                user_id="user123",
                job_description="Test job"
            )

    @patch('resume_generator.os.path.exists')
    @patch('resume_generator.os.path.isdir')
    @patch('resume_generator.os.path.abspath')
    @patch('resume_generator.os.makedirs')
    @patch('resume_generator.os.path.getsize')
    @patch('resume_generator.HTML')
    @patch('resume_generator.CSS')
    @patch('resume_generator.datetime')
    def test_generate_resume_portuguese(self, mock_datetime, mock_css, mock_html, mock_getsize,
                                         mock_makedirs, mock_abspath, mock_isdir, mock_exists):
        """Test resume generation in Portuguese."""
        mock_exists.return_value = True
        mock_isdir.return_value = True
        mock_abspath.return_value = "/app/templates"
        mock_getsize.return_value = 12345
        mock_datetime.now.return_value.strftime = Mock(return_value="20240101_120000")
        mock_datetime.now.return_value.isoformat = Mock(return_value="2024-01-01T12:00:00")

        mock_db = Mock()
        mock_db.get_user_info = Mock(return_value={
            'id': 'user123',
            'email': 'test@example.com'
        })
        mock_db.get_user_projects = Mock(return_value=[])
        mock_db.get_user_certifications = Mock(return_value=[])
        mock_db.get_user_publications = Mock(return_value=[])
        mock_db.get_user_experiences = Mock(return_value=[])

        mock_ai = Mock()
        mock_ai.generate_resume_section = Mock(return_value="<p>Conteúdo</p>")
        mock_ai.generate_tags = Mock(return_value=["tag1"])

        mock_translation = Mock()
        mock_translation.translate_projects = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_certifications = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_posts = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_experiences = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_education_status = Mock(side_effect=lambda x, y: x)
        mock_translation.translate_section_header = Mock(side_effect=lambda x, y: x)

        mock_template_env = Mock()
        mock_template = Mock()
        mock_template.render = Mock(return_value="<html>Resume HTML</html>")
        mock_template_env.get_template = Mock(return_value=mock_template)

        mock_html_instance = Mock()
        mock_html.return_value = mock_html_instance

        generator = ResumeGenerator(mock_db, mock_ai, "/app/output", mock_translation)
        generator.template_env = mock_template_env

        result = generator.generate_resume(
            user_id="user123",
            job_description="Desenvolvedor backend",
            language="pt"
        )

        # Verify Portuguese template was used
        mock_template_env.get_template.assert_called_with('resume_pt.html')
        assert result is not None
