"""
Unit tests for translation helper.
"""
import pytest
import sys
import os
from unittest.mock import Mock, patch, MagicMock

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))
from translation_helper import TranslationHelper


class TestTranslationHelper:
    """Tests for TranslationHelper class."""

    def test_init(self):
        """Test TranslationHelper initialization."""
        helper = TranslationHelper("postgresql://test:test@localhost/test")
        assert helper.database_url == "postgresql://test:test@localhost/test"
        assert helper.conn is None

    @patch('translation_helper.psycopg2.connect')
    def test_connect(self, mock_connect):
        """Test database connection."""
        mock_conn = Mock()
        mock_connect.return_value = mock_conn

        helper = TranslationHelper("postgresql://test:test@localhost/test")
        helper.connect()

        assert helper.conn == mock_conn
        mock_connect.assert_called_once_with("postgresql://test:test@localhost/test")

    @patch('translation_helper.psycopg2.connect')
    def test_connect_error(self, mock_connect):
        """Test connection error handling."""
        mock_connect.side_effect = Exception("Connection failed")

        helper = TranslationHelper("postgresql://test:test@localhost/test")
        with pytest.raises(Exception):
            helper.connect()

    def test_close(self):
        """Test closing connection."""
        helper = TranslationHelper("postgresql://test:test@localhost/test")
        mock_conn = Mock()
        helper.conn = mock_conn

        helper.close()

        mock_conn.close.assert_called_once()

    def test_close_no_connection(self):
        """Test closing when no connection exists."""
        helper = TranslationHelper("postgresql://test:test@localhost/test")
        helper.conn = None

        # Should not raise error
        helper.close()

    @patch('translation_helper.psycopg2.connect')
    def test_get_translation(self, mock_connect):
        """Test getting translation for an entity."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        mock_result = Mock()
        mock_result.__getitem__ = Mock(side_effect=lambda key: {'fields': '{"name": "Translated Name"}', 'status': 'completed'}[key])
        mock_cursor.fetchone = Mock(return_value=mock_result)
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        helper = TranslationHelper("postgresql://test:test@localhost/test")
        helper.connect()

        result = helper.get_translation("project", "proj123", "pt")

        assert result is not None
        mock_cursor.execute.assert_called_once()

    @patch('translation_helper.psycopg2.connect')
    def test_get_translation_no_connection(self, mock_connect):
        """Test getting translation without connection."""
        helper = TranslationHelper("postgresql://test:test@localhost/test")
        helper.conn = None

        result = helper.get_translation("project", "proj123", "pt")

        assert result is None

    @patch('translation_helper.psycopg2.connect')
    def test_get_translation_language_mapping(self, mock_connect):
        """Test language code mapping."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        mock_cursor.fetchone = Mock(return_value=None)
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        helper = TranslationHelper("postgresql://test:test@localhost/test")
        helper.connect()

        # Test pt -> pt-BR mapping
        helper.get_translation("project", "proj123", "pt")
        call_args = mock_cursor.execute.call_args[0]
        assert call_args[1][2] == "pt-BR"  # language parameter

    @patch('translation_helper.psycopg2.connect')
    def test_apply_translations_to_projects(self, mock_connect):
        """Test applying translations to projects."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        mock_result = Mock()
        mock_result.__getitem__ = Mock(side_effect=lambda key: {'fields': '{"name": "Translated Project"}', 'status': 'completed'}[key])
        mock_cursor.fetchone = Mock(return_value=mock_result)
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        helper = TranslationHelper("postgresql://test:test@localhost/test")
        helper.connect()

        projects = [
            {'id': 'proj1', 'name': 'Original Project'}
        ]

        result = helper.translate_projects(projects, "pt")

        assert len(result) == 1
        # Translation should be applied if found
        mock_cursor.execute.assert_called()

    @patch('translation_helper.psycopg2.connect')
    def test_apply_translations_to_certifications(self, mock_connect):
        """Test applying translations to certifications."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        mock_result = Mock()
        mock_result.__getitem__ = Mock(side_effect=lambda key: {'fields': '{"name": "Translated Cert"}', 'status': 'completed'}[key])
        mock_cursor.fetchone = Mock(return_value=mock_result)
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        helper = TranslationHelper("postgresql://test:test@localhost/test")
        helper.connect()

        certifications = [
            {'id': 'cert1', 'name': 'Original Cert'}
        ]

        result = helper.translate_certifications(certifications, "pt")

        assert len(result) == 1
        mock_cursor.execute.assert_called()
