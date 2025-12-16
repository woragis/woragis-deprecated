"""
Unit tests for database module.
"""
import pytest
import sys
import os
from unittest.mock import Mock, patch, MagicMock
import json

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))
from database import Database


class TestDatabase:
    """Tests for Database class."""

    def test_init(self):
        """Test Database initialization."""
        db = Database("postgresql://test:test@localhost/test")
        assert db.database_url == "postgresql://test:test@localhost/test"
        assert db.conn is None

    @patch('database.psycopg2.connect')
    def test_connect(self, mock_connect):
        """Test database connection."""
        mock_conn = Mock()
        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        assert db.conn == mock_conn
        mock_connect.assert_called_once_with("postgresql://test:test@localhost/test")

    @patch('database.psycopg2.connect')
    def test_connect_error(self, mock_connect):
        """Test connection error handling."""
        mock_connect.side_effect = Exception("Connection failed")

        db = Database("postgresql://test:test@localhost/test")
        with pytest.raises(Exception):
            db.connect()

    def test_close(self):
        """Test closing connection."""
        db = Database("postgresql://test:test@localhost/test")
        mock_conn = Mock()
        db.conn = mock_conn

        db.close()

        mock_conn.close.assert_called_once()

    def test_close_no_connection(self):
        """Test closing when no connection exists."""
        db = Database("postgresql://test:test@localhost/test")
        db.conn = None

        # Should not raise error
        db.close()

    @patch('database.psycopg2.connect')
    def test_get_user_projects(self, mock_connect):
        """Test getting user projects."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        
        # Mock project result - use a dict-like object that supports dict() conversion
        from collections.abc import Mapping
        class Row(Mapping):
            def __init__(self, **kwargs):
                self._data = kwargs
            def __getitem__(self, key):
                return self._data[key]
            def __iter__(self):
                return iter(self._data)
            def __len__(self):
                return len(self._data)
        
        mock_project = Row(
            id='proj1',
            name='Test Project',
            description='Test Description',
            status='active',
            slug='test-project',
            technologies='[]'
        )
        mock_cursor.fetchall = Mock(return_value=[mock_project])
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        projects = db.get_user_projects("user123")

        assert len(projects) == 1
        assert projects[0]['id'] == 'proj1'
        mock_cursor.execute.assert_called_once()

    @patch('database.psycopg2.connect')
    def test_get_user_projects_with_filters(self, mock_connect):
        """Test getting user projects with category and skill filters."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        mock_cursor.fetchall = Mock(return_value=[])
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        projects = db.get_user_projects(
            "user123",
            tech_categories=['backend'],
            skill_names=['Golang']
        )

        # Verify query includes filters
        call_args = mock_cursor.execute.call_args[0]
        query = call_args[0]
        assert 'ANY' in query  # Should have array filter

    @patch('database.psycopg2.connect')
    def test_get_user_certifications(self, mock_connect):
        """Test getting user certifications."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        
        # Mock cert result
        class Row:
            def __init__(self, **kwargs):
                for k, v in kwargs.items():
                    setattr(self, k, v)
            def __getitem__(self, key):
                return getattr(self, key)
            def keys(self):
                return self.__dict__.keys()
        
        mock_cert = Row(
            id='cert1',
            name='AWS Certified',
            issuer='Amazon',
            status='active'
        )
        mock_cursor.fetchall = Mock(return_value=[mock_cert])
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        certs = db.get_user_certifications("user123")

        assert len(certs) == 1
        assert certs[0]['id'] == 'cert1'
        mock_cursor.execute.assert_called_once()

    @patch('database.psycopg2.connect')
    def test_get_user_certifications_with_categories(self, mock_connect):
        """Test getting user certifications with category filter."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        mock_cursor.fetchall = Mock(return_value=[])
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        certs = db.get_user_certifications("user123", categories=['cloud'])

        # Verify query includes category filter
        call_args = mock_cursor.execute.call_args[0]
        query = call_args[0]
        assert 'category = ANY' in query

    @patch('database.psycopg2.connect')
    def test_get_user_publications(self, mock_connect):
        """Test getting user publications."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        mock_cursor.fetchall = Mock(return_value=[])
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        publications = db.get_user_publications("user123", limit=5)

        # Verify query includes limit
        call_args = mock_cursor.execute.call_args[0]
        assert call_args[1][1] == 5  # limit parameter

    @patch('database.psycopg2.connect')
    def test_get_user_experiences(self, mock_connect):
        """Test getting user experiences."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        
        # Mock experience result
        class Row:
            def __init__(self, **kwargs):
                for k, v in kwargs.items():
                    setattr(self, k, v)
            def __getitem__(self, key):
                return getattr(self, key)
            def keys(self):
                return self.__dict__.keys()
        
        mock_exp = Row(
            id='exp1',
            company='Test Company',
            position='Developer',
            period_start=None,
            period_end=None,
            period_text='2020 - 2022',
            description='Line 1\nLine 2',
            is_current=False
        )
        mock_cursor.fetchall = Mock(return_value=[mock_exp])
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        experiences = db.get_user_experiences("user123")

        assert len(experiences) == 1
        assert experiences[0]['id'] == 'exp1'
        # Verify description was parsed into list
        assert isinstance(experiences[0].get('description'), list)

    @patch('database.psycopg2.connect')
    def test_get_user_info(self, mock_connect):
        """Test getting user info."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        
        # Mock user result
        class Row:
            def __init__(self, **kwargs):
                for k, v in kwargs.items():
                    setattr(self, k, v)
            def __getitem__(self, key):
                return getattr(self, key)
            def keys(self):
                return self.__dict__.keys()
        
        mock_user = Row(id='user123', email='test@example.com')
        mock_cursor.fetchone = Mock(return_value=mock_user)
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        user_info = db.get_user_info("user123")

        assert user_info is not None
        assert user_info['id'] == 'user123'
        assert 'name' in user_info  # Should have generated name

    @patch('database.psycopg2.connect')
    def test_get_user_info_not_found(self, mock_connect):
        """Test getting user info when user doesn't exist."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        mock_cursor.fetchone = Mock(return_value=None)
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        user_info = db.get_user_info("user123")

        assert user_info is None

    @patch('database.psycopg2.connect')
    def test_get_user_info_special_email(self, mock_connect):
        """Test getting user info with special email."""
        mock_conn = Mock()
        mock_cursor = MagicMock()
        mock_cursor.__enter__ = Mock(return_value=mock_cursor)
        mock_cursor.__exit__ = Mock(return_value=None)
        mock_cursor.execute = Mock()
        
        # Mock user result with special email
        class Row:
            def __init__(self, **kwargs):
                for k, v in kwargs.items():
                    setattr(self, k, v)
            def __getitem__(self, key):
                return getattr(self, key)
            def keys(self):
                return self.__dict__.keys()
        
        mock_user = Row(id='user123', email='masteringthecode.woragis@gmail.com')
        mock_cursor.fetchone = Mock(return_value=mock_user)
        mock_conn.cursor = Mock(return_value=mock_cursor)

        mock_connect.return_value = mock_conn

        db = Database("postgresql://test:test@localhost/test")
        db.connect()

        user_info = db.get_user_info("user123")

        assert user_info['name'] == 'Jezreel de Andrade Galvao Veloso'
