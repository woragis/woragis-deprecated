"""
Unit tests for main module.
"""
import pytest
import sys
import os
import json
from unittest.mock import Mock, patch, MagicMock, mock_open
from datetime import datetime

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))
import main


class TestSaveResult:
    """Tests for save_result function."""

    @patch('main.os.makedirs')
    @patch('main.datetime')
    @patch('builtins.open', new_callable=mock_open)
    def test_save_result(self, mock_file, mock_datetime, mock_makedirs):
        """Test saving result to file."""
        mock_datetime.now.return_value.strftime = Mock(return_value="20240101_120000")

        result = {
            'user_id': 'user123',
            'output_path': '/app/output/resume.pdf',
            'file_size': 12345
        }

        main.save_result(result, "/app/results")

        # Verify directory was created
        mock_makedirs.assert_called_once_with("/app/results", exist_ok=True)

        # Verify file was written
        mock_file.assert_called_once()
        written_data = ''.join(call.args[0] for call in mock_file().write.call_args_list)
        assert 'user123' in written_data

    @patch('main.os.makedirs')
    @patch('main.datetime')
    @patch('builtins.open', side_effect=IOError("Permission denied"))
    def test_save_result_error(self, mock_file, mock_datetime, mock_makedirs):
        """Test saving result when file write fails."""
        mock_datetime.now.return_value.strftime = Mock(return_value="20240101_120000")

        result = {'user_id': 'user123'}

        # Should not raise exception
        main.save_result(result, "/app/results")


class TestProcessResumeJob:
    """Tests for process_resume_job function."""

    @patch('main.Database')
    @patch('main.AIService')
    @patch('main.TranslationHelper')
    @patch('main.ResumeGenerator')
    @patch('main.save_result')
    def test_process_resume_job_success(self, mock_save_result, mock_generator_class,
                                         mock_translation_class, mock_ai_class, mock_db_class):
        """Test successful resume job processing."""
        # Setup mocks
        mock_db = Mock()
        mock_db.connect = Mock()
        mock_db.close = Mock()
        mock_db_class.return_value = mock_db

        mock_ai = Mock()
        mock_ai_class.return_value = mock_ai

        mock_translation = Mock()
        mock_translation.connect = Mock()
        mock_translation.close = Mock()
        mock_translation_class.return_value = mock_translation

        mock_generator = Mock()
        mock_generator.generate_resume = Mock(return_value={
            'output_path': '/app/output/resume.pdf',
            'file_size': 12345,
            'projects_count': 2,
            'certifications_count': 1
        })
        mock_generator_class.return_value = mock_generator

        message = {
            'user_id': 'user123',
            'job_description': 'Backend developer with Golang',
            'job_title': 'Software Engineer',
            'language': 'en'
        }

        result = main.process_resume_job(message)

        assert result is True
        mock_generator.generate_resume.assert_called_once()
        mock_save_result.assert_called_once()
        mock_db.close.assert_called_once()
        mock_translation.close.assert_called_once()

    def test_process_resume_job_missing_fields(self):
        """Test processing job with missing required fields."""
        message = {
            'user_id': 'user123'
            # Missing job_description
        }

        result = main.process_resume_job(message)

        assert result is False

    @patch('main.Database')
    @patch('main.AIService')
    @patch('main.TranslationHelper')
    @patch('main.ResumeGenerator')
    def test_process_resume_job_error(self, mock_generator_class, mock_translation_class,
                                       mock_ai_class, mock_db_class):
        """Test processing job when generation fails."""
        mock_db = Mock()
        mock_db.connect = Mock()
        mock_db.close = Mock()
        mock_db_class.return_value = mock_db

        mock_ai = Mock()
        mock_ai_class.return_value = mock_ai

        mock_translation = Mock()
        mock_translation.connect = Mock()
        mock_translation.close = Mock()
        mock_translation_class.return_value = mock_translation

        mock_generator = Mock()
        mock_generator.generate_resume = Mock(side_effect=Exception("Generation failed"))
        mock_generator_class.return_value = mock_generator

        message = {
            'user_id': 'user123',
            'job_description': 'Test job'
        }

        result = main.process_resume_job(message)

        assert result is False
        mock_db.close.assert_called_once()
        mock_translation.close.assert_called_once()


class TestRunCLIMode:
    """Tests for run_cli_mode function."""

    @patch('main.Database')
    @patch('main.AIService')
    @patch('main.TranslationHelper')
    @patch('main.ResumeGenerator')
    @patch('main.save_result')
    @patch('sys.argv', ['main.py', 'user123', 'Backend developer'])
    @patch('sys.exit')
    def test_run_cli_mode_success(self, mock_exit, mock_save_result, mock_generator_class,
                                    mock_translation_class, mock_ai_class, mock_db_class):
        """Test successful CLI mode execution."""
        mock_db = Mock()
        mock_db.connect = Mock()
        mock_db.close = Mock()
        mock_db_class.return_value = mock_db

        mock_ai = Mock()
        mock_ai_class.return_value = mock_ai

        mock_translation = Mock()
        mock_translation.connect = Mock()
        mock_translation.close = Mock()
        mock_translation_class.return_value = mock_translation

        mock_generator = Mock()
        mock_generator.generate_resume = Mock(return_value={
            'output_path': '/app/output/resume.pdf',
            'file_size': 12345,
            'projects_count': 2,
            'certifications_count': 1
        })
        mock_generator_class.return_value = mock_generator

        main.run_cli_mode()

        mock_generator.generate_resume.assert_called_once()
        mock_exit.assert_called_once_with(0)

    @patch('sys.argv', ['main.py', 'user123'])
    @patch('sys.exit')
    def test_run_cli_mode_insufficient_args(self, mock_exit):
        """Test CLI mode with insufficient arguments."""
        main.run_cli_mode()
        mock_exit.assert_called_once_with(1)


class TestRunQueueMode:
    """Tests for run_queue_mode function."""

    @patch('main.create_consumer_from_env')
    @patch('http.server.HTTPServer')
    @patch('threading.Thread')
    def test_run_queue_mode(self, mock_thread, mock_httpserver, mock_create_consumer):
        """Test queue mode execution."""
        mock_consumer = Mock()
        mock_consumer.connect = Mock(return_value=True)
        mock_consumer.start_consuming = Mock()
        mock_create_consumer.return_value = mock_consumer

        mock_server = Mock()
        mock_httpserver.return_value = mock_server

        mock_thread_instance = Mock()
        mock_thread.return_value = mock_thread_instance

        # Simulate KeyboardInterrupt to stop the loop
        mock_consumer.start_consuming.side_effect = KeyboardInterrupt()

        try:
            main.run_queue_mode()
        except KeyboardInterrupt:
            pass  # Expected

        mock_consumer.connect.assert_called_once()
        mock_consumer.start_consuming.assert_called_once()

    @patch('main.create_consumer_from_env')
    @patch('http.server.HTTPServer')
    @patch('threading.Thread')
    @patch('sys.exit')
    def test_run_queue_mode_connection_failed(self, mock_exit, mock_thread, mock_httpserver,
                                               mock_create_consumer):
        """Test queue mode when connection fails."""
        mock_consumer = Mock()
        mock_consumer.connect = Mock(return_value=False)
        mock_create_consumer.return_value = mock_consumer

        mock_server = Mock()
        mock_httpserver.return_value = mock_server

        mock_thread_instance = Mock()
        mock_thread.return_value = mock_thread_instance

        main.run_queue_mode()

        mock_exit.assert_called_once_with(1)


class TestMain:
    """Tests for main function."""

    @patch('main.run_cli_mode')
    @patch('sys.argv', ['main.py', 'user123', 'Job description'])
    def test_main_cli_mode(self, mock_run_cli):
        """Test main function in CLI mode."""
        main.main()
        mock_run_cli.assert_called_once()

    @patch('main.run_queue_mode')
    @patch('sys.argv', ['main.py'])
    def test_main_queue_mode(self, mock_run_queue):
        """Test main function in queue mode."""
        main.main()
        mock_run_queue.assert_called_once()

    @patch('sys.argv', ['main.py', 'user123'])
    @patch('sys.exit')
    def test_main_cli_insufficient_args(self, mock_exit):
        """Test main function with insufficient CLI args."""
        main.main()
        mock_exit.assert_called_once_with(1)
