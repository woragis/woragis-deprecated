"""
Unit tests for queue consumer module.
"""
import pytest
import sys
import os
import json
import signal
from unittest.mock import Mock, patch, MagicMock

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', 'src'))
from queue_consumer import ResumeQueueConsumer, create_consumer_from_env


class TestResumeQueueConsumer:
    """Tests for ResumeQueueConsumer class."""

    def test_init(self):
        """Test ResumeQueueConsumer initialization."""
        consumer = ResumeQueueConsumer(
            rabbitmq_url="amqp://user:pass@localhost:5672/vhost",
            queue_name="test.queue",
            exchange="test.exchange",
            routing_key="test.key"
        )

        assert consumer.rabbitmq_url == "amqp://user:pass@localhost:5672/vhost"
        assert consumer.queue_name == "test.queue"
        assert consumer.exchange == "test.exchange"
        assert consumer.routing_key == "test.key"
        assert consumer.connection is None
        assert consumer.channel is None
        assert consumer.should_stop is False

    @patch('queue_consumer.pika.BlockingConnection')
    @patch('queue_consumer.pika.URLParameters')
    def test_connect(self, mock_url_params, mock_blocking_connection):
        """Test connecting to RabbitMQ."""
        mock_connection = Mock()
        mock_channel = Mock()
        mock_connection.channel.return_value = mock_channel
        mock_blocking_connection.return_value = mock_connection

        mock_params = Mock()
        mock_url_params.return_value = mock_params

        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")

        result = consumer.connect()

        assert result is True
        assert consumer.connection == mock_connection
        assert consumer.channel == mock_channel
        mock_channel.exchange_declare.assert_called_once()
        mock_channel.queue_declare.assert_called_once()
        mock_channel.queue_bind.assert_called_once()
        mock_channel.basic_qos.assert_called_once()

    @patch('queue_consumer.pika.BlockingConnection')
    @patch('queue_consumer.pika.URLParameters')
    def test_connect_error(self, mock_url_params, mock_blocking_connection):
        """Test connection error handling."""
        from pika.exceptions import AMQPConnectionError
        mock_blocking_connection.side_effect = AMQPConnectionError("Connection failed")

        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")

        result = consumer.connect()

        assert result is False

    @patch('queue_consumer.pika.BlockingConnection')
    @patch('queue_consumer.pika.URLParameters')
    def test_close(self, mock_url_params, mock_blocking_connection):
        """Test closing connection."""
        mock_connection = Mock()
        mock_connection.is_closed = False
        mock_channel = Mock()
        mock_channel.is_closed = False
        mock_connection.channel.return_value = mock_channel
        mock_blocking_connection.return_value = mock_connection

        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")
        consumer.connect()
        consumer.close()

        mock_channel.close.assert_called_once()
        mock_connection.close.assert_called_once()

    @patch('queue_consumer.pika.BlockingConnection')
    @patch('queue_consumer.pika.URLParameters')
    def test_close_already_closed(self, mock_url_params, mock_blocking_connection):
        """Test closing when already closed."""
        mock_connection = Mock()
        mock_connection.is_closed = True
        mock_channel = Mock()
        mock_channel.is_closed = True
        mock_connection.channel.return_value = mock_channel
        mock_blocking_connection.return_value = mock_connection

        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")
        consumer.connect()
        consumer.close()

        # Should not raise error
        assert True

    @patch('queue_consumer.pika.BlockingConnection')
    @patch('queue_consumer.pika.URLParameters')
    def test_start_consuming(self, mock_url_params, mock_blocking_connection):
        """Test starting message consumption."""
        mock_connection = Mock()
        mock_connection.is_closed = False
        mock_connection.process_data_events = Mock()
        mock_channel = Mock()
        mock_connection.channel.return_value = mock_channel
        mock_blocking_connection.return_value = mock_connection

        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")
        consumer.should_stop = True  # Stop immediately

        def message_handler(msg):
            return True

        consumer.start_consuming(message_handler)

        mock_channel.basic_consume.assert_called_once()

    @patch('queue_consumer.pika.BlockingConnection')
    @patch('queue_consumer.pika.URLParameters')
    def test_start_consuming_not_connected(self, mock_url_params, mock_blocking_connection):
        """Test starting consumption when not connected."""
        mock_connection = Mock()
        mock_connection.is_closed = True
        mock_connection.channel.return_value = Mock()
        mock_blocking_connection.return_value = mock_connection

        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")
        consumer.connection = None

        def message_handler(msg):
            return True

        # Should attempt to connect
        with patch.object(consumer, 'connect', return_value=False):
            consumer.start_consuming(message_handler)
            # Should return early without starting consumption

    def test_create_callback_success(self):
        """Test callback creation with successful message processing."""
        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")

        def message_handler(msg):
            return True

        callback = consumer._create_callback(message_handler)

        mock_channel = Mock()
        mock_method = Mock()
        mock_method.delivery_tag = 123
        mock_properties = Mock()
        message_body = json.dumps({'id': 'job123', 'user_id': 'user123'}).encode('utf-8')

        callback(mock_channel, mock_method, mock_properties, message_body)

        mock_channel.basic_ack.assert_called_once_with(delivery_tag=123)

    def test_create_callback_failure(self):
        """Test callback creation with failed message processing."""
        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")

        def message_handler(msg):
            return False

        callback = consumer._create_callback(message_handler)

        mock_channel = Mock()
        mock_method = Mock()
        mock_method.delivery_tag = 123
        mock_properties = Mock()
        message_body = json.dumps({'id': 'job123'}).encode('utf-8')

        callback(mock_channel, mock_method, mock_properties, message_body)

        mock_channel.basic_nack.assert_called_once_with(delivery_tag=123, requeue=True)

    def test_create_callback_invalid_json(self):
        """Test callback with invalid JSON message."""
        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")

        def message_handler(msg):
            return True

        callback = consumer._create_callback(message_handler)

        mock_channel = Mock()
        mock_method = Mock()
        mock_method.delivery_tag = 123
        mock_properties = Mock()
        message_body = b"invalid json"

        callback(mock_channel, mock_method, mock_properties, message_body)

        mock_channel.basic_nack.assert_called_once_with(delivery_tag=123, requeue=False)

    def test_create_callback_exception(self):
        """Test callback when message handler raises exception."""
        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")

        def message_handler(msg):
            raise Exception("Processing error")

        callback = consumer._create_callback(message_handler)

        mock_channel = Mock()
        mock_method = Mock()
        mock_method.delivery_tag = 123
        mock_properties = Mock()
        message_body = json.dumps({'id': 'job123'}).encode('utf-8')

        callback(mock_channel, mock_method, mock_properties, message_body)

        mock_channel.basic_nack.assert_called_once_with(delivery_tag=123, requeue=True)

    @patch('queue_consumer.signal.signal')
    def test_signal_handler(self, mock_signal):
        """Test signal handler setup."""
        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")

        # Verify signal handlers were set up
        assert mock_signal.call_count == 2

    def test_expose_for_health(self):
        """Test exposing connection for health checks."""
        consumer = ResumeQueueConsumer("amqp://user:pass@localhost:5672/vhost")
        consumer.connection = Mock()

        # Should not raise even if health module is not available
        consumer._expose_for_health()
        
        # Verify the method completes without error
        assert True


class TestCreateConsumerFromEnv:
    """Tests for create_consumer_from_env function."""

    @patch.dict(os.environ, {
        'RABBITMQ_URL': 'amqp://user:pass@localhost:5672/vhost',
        'RESUME_QUEUE_NAME': 'test.queue',
        'RESUME_EXCHANGE': 'test.exchange',
        'RESUME_ROUTING_KEY': 'test.key'
    })
    def test_create_consumer_from_env_with_url(self):
        """Test creating consumer from environment with full URL."""
        consumer = create_consumer_from_env()

        assert consumer.rabbitmq_url == 'amqp://user:pass@localhost:5672/vhost'
        assert consumer.queue_name == 'test.queue'
        assert consumer.exchange == 'test.exchange'
        assert consumer.routing_key == 'test.key'

    @patch.dict(os.environ, {
        'RABBITMQ_USER': 'testuser',
        'RABBITMQ_PASSWORD': 'testpass',
        'RABBITMQ_HOST': 'rabbitmq',
        'RABBITMQ_PORT': '5672',
        'RABBITMQ_VHOST': 'woragis',
        'RESUME_QUEUE_NAME': 'resumes.queue'
    }, clear=True)
    def test_create_consumer_from_env_with_components(self):
        """Test creating consumer from environment with components."""
        consumer = create_consumer_from_env()

        assert 'amqp://testuser:testpass@rabbitmq:5672/woragis' == consumer.rabbitmq_url
        assert consumer.queue_name == 'resumes.queue'

    @patch.dict(os.environ, {}, clear=True)
    def test_create_consumer_from_env_defaults(self):
        """Test creating consumer with default values."""
        consumer = create_consumer_from_env()

        assert 'amqp://woragis:woragis@rabbitmq:5672/woragis' == consumer.rabbitmq_url
        assert consumer.queue_name == 'resumes.queue'
        assert consumer.exchange == 'woragis.tasks'
        assert consumer.routing_key == 'resumes.queue'
