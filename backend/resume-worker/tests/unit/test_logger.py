"""
Unit tests for logger module.
"""
import pytest
import sys
import os
from unittest.mock import patch, Mock, MagicMock

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))
from logger import configure_logging, get_logger, set_trace_id, get_trace_id, trace_id_var


class TestLogger:
    """Tests for logger functions."""

    def test_configure_logging_development(self):
        """Test logging configuration for development environment."""
        with patch('logger.structlog.configure') as mock_configure, \
             patch('logging.basicConfig') as mock_basic_config, \
             patch('os.makedirs'):
            configure_logging(env="development", log_to_file=False)
            
            # Verify structlog was configured
            mock_configure.assert_called_once()
            call_args = mock_configure.call_args
            assert 'processors' in call_args.kwargs
            assert 'wrapper_class' in call_args.kwargs
            
            # Verify basicConfig was called
            mock_basic_config.assert_called_once()

    def test_configure_logging_production(self):
        """Test logging configuration for production environment."""
        with patch('logger.structlog.configure') as mock_configure, \
             patch('logging.basicConfig') as mock_basic_config:
            configure_logging(env="production", log_to_file=False)
            
            # Verify structlog was configured
            mock_configure.assert_called_once()
            call_args = mock_configure.call_args
            processors = call_args.kwargs['processors']
            
            # Check that JSONRenderer is in processors for production
            processor_names = [p.__name__ if hasattr(p, '__name__') else str(p) for p in processors]
            assert any('JSONRenderer' in str(p) for p in processors)

    def test_configure_logging_with_file_logging(self):
        """Test logging configuration with file logging enabled."""
        with patch('logger.structlog.configure'), \
             patch('logging.basicConfig'), \
             patch('logging.getLogger') as mock_get_logger, \
             patch('logging.FileHandler') as mock_file_handler, \
             patch('os.makedirs') as mock_makedirs, \
             patch('os.path.join', return_value='/app/logs/resume-worker.log'):
            mock_logger = Mock()
            mock_get_logger.return_value = mock_logger
            
            configure_logging(env="development", log_to_file=True, log_dir="/app/logs")
            
            # Verify directory was created
            mock_makedirs.assert_called_once_with("/app/logs", exist_ok=True)
            
            # Verify file handler was added
            mock_file_handler.assert_called_once()

    def test_get_logger(self):
        """Test getting a logger instance."""
        with patch('logger.structlog.get_logger') as mock_get_logger:
            mock_logger = Mock()
            mock_get_logger.return_value = mock_logger
            
            result = get_logger("test.logger")
            
            # Verify structlog.get_logger was called
            mock_get_logger.assert_called_once_with("test.logger")
            
            # Verify logger was bound with service name
            assert mock_logger.bind.called

    def test_get_logger_with_trace_id(self):
        """Test getting a logger with trace_id in context."""
        with patch('logger.structlog.get_logger') as mock_get_logger:
            mock_logger = Mock()
            mock_get_logger.return_value = mock_logger
            
            # Set trace_id in context
            set_trace_id("test-trace-123")
            
            result = get_logger("test.logger")
            
            # Verify logger was bound with trace_id
            assert mock_logger.bind.called

    def test_set_trace_id(self):
        """Test setting trace_id in context."""
        set_trace_id("test-trace-456")
        assert get_trace_id() == "test-trace-456"

    def test_get_trace_id(self):
        """Test getting trace_id from context."""
        # Clear trace_id first
        trace_id_var.set(None)
        assert get_trace_id() is None
        
        # Set and get trace_id
        set_trace_id("test-trace-789")
        assert get_trace_id() == "test-trace-789"

    def test_get_trace_id_none(self):
        """Test getting trace_id when none is set."""
        trace_id_var.set(None)
        assert get_trace_id() is None
