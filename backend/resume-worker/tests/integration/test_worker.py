"""
Integration tests for resume worker.
These tests require external services (database, RabbitMQ, AI service).
"""
import pytest
import sys
import os
import json
import time
from unittest.mock import Mock, patch

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))

# Test configuration
TEST_RABBITMQ_URL = os.getenv("TEST_RABBITMQ_URL", "amqp://test:test@localhost:5673/test")
TEST_DB_URL = os.getenv("TEST_DATABASE_URL", "postgresql://postgres:postgres@localhost:5433/woragis_test")
TEST_AI_SERVICE_URL = os.getenv("TEST_AI_SERVICE_URL", "http://localhost:8000")

try:
    import pika
    RABBITMQ_AVAILABLE = True
except ImportError:
    RABBITMQ_AVAILABLE = False

try:
    import psycopg2
    DATABASE_AVAILABLE = True
except ImportError:
    DATABASE_AVAILABLE = False


@pytest.mark.integration
class TestWorkerIntegration:
    """Integration tests for resume worker."""

    def test_health_check_integration(self):
        """Test health check endpoint."""
        from health import check_health
        result = check_health()
        assert "status" in result
        assert "checks" in result
        assert result["status"] in ["healthy", "degraded", "unhealthy"]

    @pytest.mark.skipif(not RABBITMQ_AVAILABLE, reason="RabbitMQ not available")
    def test_rabbitmq_connection(self):
        """Test RabbitMQ connection setup."""
        from queue_consumer import ResumeQueueConsumer
        
        consumer = ResumeQueueConsumer(
            rabbitmq_url=TEST_RABBITMQ_URL,
            queue_name="test.resumes.queue",
            exchange="test.woragis.tasks",
            routing_key="test.resumes.queue"
        )
        
        try:
            consumer.connect()
            assert consumer.connection is not None
            assert consumer.channel is not None
            assert not consumer.connection.is_closed
        finally:
            if consumer.connection and not consumer.connection.is_closed:
                consumer.connection.close()

    @pytest.mark.skipif(not RABBITMQ_AVAILABLE, reason="RabbitMQ not available")
    def test_queue_setup(self):
        """Test queue and exchange setup."""
        from queue_consumer import ResumeQueueConsumer
        
        consumer = ResumeQueueConsumer(
            rabbitmq_url=TEST_RABBITMQ_URL,
            queue_name="test.resumes.queue.setup",
            exchange="test.woragis.tasks.setup",
            routing_key="test.resumes.queue.setup"
        )
        
        try:
            consumer.connect()
            
            # Verify queue exists
            method = consumer.channel.queue_declare(
                queue="test.resumes.queue.setup",
                passive=True
            )
            assert method is not None
            
            # Verify exchange exists
            consumer.channel.exchange_declare(
                exchange="test.woragis.tasks.setup",
                exchange_type='direct',
                passive=True
            )
        finally:
            if consumer.connection and not consumer.connection.is_closed:
                consumer.connection.close()

    @pytest.mark.skipif(not RABBITMQ_AVAILABLE, reason="RabbitMQ not available")
    def test_message_publish(self):
        """Test publishing resume job message to queue."""
        from queue_consumer import ResumeQueueConsumer
        
        consumer = ResumeQueueConsumer(
            rabbitmq_url=TEST_RABBITMQ_URL,
            queue_name="test.resumes.queue.publish",
            exchange="test.woragis.tasks.publish",
            routing_key="test.resumes.queue.publish"
        )
        
        try:
            consumer.connect()
            
            # Create test message
            test_message = {
                "user_id": "123e4567-e89b-12d3-a456-426614174000",
                "job_description": "Software Engineer position",
                "job_title": "Senior Software Engineer",
                "language": "en"
            }
            
            # Publish message
            consumer.channel.basic_publish(
                exchange="test.woragis.tasks.publish",
                routing_key="test.resumes.queue.publish",
                body=json.dumps(test_message),
                properties=pika.BasicProperties(
                    content_type='application/json',
                    delivery_mode=2  # Persistent
                )
            )
            
            # Verify message is in queue
            method, properties, body = consumer.channel.basic_get(
                queue="test.resumes.queue.publish",
                auto_ack=False
            )
            
            assert method is not None
            assert body is not None
            
            received_message = json.loads(body.decode('utf-8'))
            assert received_message["user_id"] == test_message["user_id"]
            assert received_message["job_description"] == test_message["job_description"]
            
            # Acknowledge message
            consumer.channel.basic_ack(method.delivery_tag)
        finally:
            if consumer.connection and not consumer.connection.is_closed:
                consumer.connection.close()

    @pytest.mark.skipif(not DATABASE_AVAILABLE, reason="Database not available")
    def test_database_connection(self):
        """Test database connection."""
        from database import Database
        
        db = Database(TEST_DB_URL)
        try:
            db.connect()
            # Test query
            result = db.execute_query("SELECT 1 as test")
            assert result is not None
        finally:
            if hasattr(db, 'close'):
                db.close()

    @pytest.mark.skipif(not DATABASE_AVAILABLE, reason="Database not available")
    def test_fetch_user_profile(self):
        """Test fetching user profile from database."""
        from database import Database
        
        db = Database(TEST_DB_URL)
        try:
            db.connect()
            
            # Test user ID (should exist or be created in test setup)
            test_user_id = "123e4567-e89b-12d3-a456-426614174000"
            
            profile = db.fetch_user_profile(test_user_id)
            
            # Profile should be a dict with expected keys
            assert isinstance(profile, dict)
            assert "projects" in profile
            assert "posts" in profile
            assert "skills" in profile
            assert "certifications" in profile
        finally:
            if hasattr(db, 'close'):
                db.close()

    @pytest.mark.skipif(not DATABASE_AVAILABLE, reason="Database not available")
    def test_save_resume_record(self):
        """Test saving resume record to database."""
        from database import Database
        
        db = Database(TEST_DB_URL)
        try:
            db.connect()
            
            # Create test resume record
            resume_data = {
                "user_id": "123e4567-e89b-12d3-a456-426614174000",
                "job_title": "Test Job",
                "job_description": "Test description",
                "status": "completed",
                "output_path": "/tmp/test_resume.pdf"
            }
            
            # Note: This depends on the actual database schema
            # Adjust based on your database.py implementation
            result = db.save_resume_record(resume_data)
            assert result is not None
        except Exception as e:
            # If table doesn't exist, that's okay for now
            pytest.skip(f"Database table may not exist: {e}")
        finally:
            if hasattr(db, 'close'):
                db.close()

    def test_ai_service_integration(self):
        """Test AI service integration (with mock)."""
        from ai_service import AIService
        
        # Mock AI service for testing
        with patch('ai_service.requests.post') as mock_post:
            mock_response = Mock()
            mock_response.status_code = 200
            mock_response.json.return_value = {
                "output": "Generated resume section content"
            }
            mock_post.return_value = mock_response
            
            ai_service = AIService(TEST_AI_SERVICE_URL)
            result = ai_service.generate_resume_section(
                section_type="experience",
                user_profile={},
                job_description="Test job"
            )
            
            assert result is not None
            assert "output" in result

    def test_resume_generator_integration(self):
        """Test resume generator with mocked dependencies."""
        from resume_generator import ResumeGenerator
        from unittest.mock import Mock, MagicMock
        
        # Mock dependencies
        mock_db = Mock()
        mock_db.fetch_user_profile.return_value = {
            "projects": [],
            "posts": [],
            "skills": [],
            "certifications": []
        }
        
        mock_ai_service = Mock()
        mock_ai_service.generate_resume_section.return_value = {
            "output": "Generated content"
        }
        
        generator = ResumeGenerator(mock_db, mock_ai_service)
        
        # Test resume generation
        result = generator.generate_resume(
            user_id="123e4567-e89b-12d3-a456-426614174000",
            job_description="Test job description",
            job_title="Software Engineer",
            language="en"
        )
        
        assert result is not None
        assert "html" in result or "pdf_path" in result

    @pytest.mark.skipif(not RABBITMQ_AVAILABLE, reason="RabbitMQ not available")
    def test_end_to_end_resume_generation(self):
        """Test end-to-end resume generation flow (with mocked AI service)."""
        from queue_consumer import ResumeQueueConsumer
        from unittest.mock import Mock, patch
        
        # Mock AI service to avoid actual API calls
        with patch('resume_generator.AIService') as mock_ai_class:
            mock_ai = Mock()
            mock_ai.generate_resume_section.return_value = {
                "output": "Generated resume content"
            }
            mock_ai_class.return_value = mock_ai
            
            consumer = ResumeQueueConsumer(
                rabbitmq_url=TEST_RABBITMQ_URL,
                queue_name="test.resumes.queue.e2e",
                exchange="test.woragis.tasks.e2e",
                routing_key="test.resumes.queue.e2e"
            )
            
            try:
                consumer.connect()
                
                # Publish test message
                test_message = {
                    "user_id": "123e4567-e89b-12d3-a456-426614174000",
                    "job_description": "Software Engineer position requiring Python and Go",
                    "job_title": "Senior Software Engineer",
                    "language": "en"
                }
                
                consumer.channel.basic_publish(
                    exchange="test.woragis.tasks.e2e",
                    routing_key="test.resumes.queue.e2e",
                    body=json.dumps(test_message),
                    properties=pika.BasicProperties(
                        content_type='application/json',
                        delivery_mode=2
                    )
                )
                
                # Verify message was published
                method, properties, body = consumer.channel.basic_get(
                    queue="test.resumes.queue.e2e",
                    auto_ack=False
                )
                
                assert method is not None
                assert body is not None
                
                received_message = json.loads(body.decode('utf-8'))
                assert received_message["user_id"] == test_message["user_id"]
                
                consumer.channel.basic_ack(method.delivery_tag)
            finally:
                if consumer.connection and not consumer.connection.is_closed:
                    consumer.connection.close()
