#!/usr/bin/env python3
"""
RabbitMQ Queue Consumer for Resume Generation Jobs

This module handles consuming resume generation jobs from RabbitMQ queue
and processing them using the ResumeGenerator.
"""

import os
import json
import logging
import signal
import sys
from typing import Optional, Dict, Any
from datetime import datetime

import pika
from pika.exceptions import AMQPConnectionError, AMQPChannelError

logger = logging.getLogger(__name__)


class ResumeQueueConsumer:
    """Consumer for resume generation jobs from RabbitMQ"""
    
    def __init__(
        self,
        rabbitmq_url: str,
        queue_name: str = "resumes.queue",
        exchange: str = "woragis.tasks",
        routing_key: str = "resumes.queue"
    ):
        """
        Initialize the queue consumer.
        
        Args:
            rabbitmq_url: RabbitMQ connection URL (e.g., amqp://user:pass@host:port/vhost)
            queue_name: Name of the queue to consume from
            exchange: Exchange name to bind the queue to
            routing_key: Routing key for queue binding
        """
        self.rabbitmq_url = rabbitmq_url
        self.queue_name = queue_name
        self.exchange = exchange
        self.routing_key = routing_key
        
        self.connection: Optional[pika.BlockingConnection] = None
        self.channel: Optional[pika.channel.Channel] = None
        self.should_stop = False
        
        # Expose connection for health checks
        self._expose_for_health()
        
        # Setup signal handlers for graceful shutdown
        signal.signal(signal.SIGTERM, self._signal_handler)
        signal.signal(signal.SIGINT, self._signal_handler)
    
    def _signal_handler(self, signum, frame):
        """Handle shutdown signals gracefully"""
        logger.info(f"Received signal {signum}, shutting down...")
        self.should_stop = True
        if self.connection and not self.connection.is_closed:
            self.connection.close()
        sys.exit(0)
    
    def connect(self):
        """Establish connection to RabbitMQ"""
        try:
            logger.info(f"Connecting to RabbitMQ: {self.rabbitmq_url}")
            params = pika.URLParameters(self.rabbitmq_url)
            self.connection = pika.BlockingConnection(params)
            self.channel = self.connection.channel()
            
            # Declare exchange
            self.channel.exchange_declare(
                exchange=self.exchange,
                exchange_type='direct',
                durable=True
            )
            
            # Declare queue with durability and dead letter exchange
            self.channel.queue_declare(
                queue=self.queue_name,
                durable=True,
                arguments={
                    'x-dead-letter-exchange': 'woragis.dlx',
                    'x-dead-letter-routing-key': self.queue_name + '.failed'
                }
            )
            
            # Bind queue to exchange
            self.channel.queue_bind(
                exchange=self.exchange,
                queue=self.queue_name,
                routing_key=self.routing_key
            )
            
            # Set QoS to process one message at a time
            self.channel.basic_qos(prefetch_count=1)
            
            logger.info(f"Connected to RabbitMQ and ready to consume from queue: {self.queue_name}")
            
            # Update health check connection reference
            self._expose_for_health()
            
            return True
        except AMQPConnectionError as e:
            logger.error(f"Failed to connect to RabbitMQ: {e}")
            return False
        except Exception as e:
            logger.error(f"Unexpected error connecting to RabbitMQ: {e}", exc_info=True)
            return False
    
    def _expose_for_health(self):
        """Expose connection for health checks."""
        try:
            from health import set_rabbitmq_connection
            set_rabbitmq_connection(self.connection)
        except ImportError:
            # Health module may not be available
            pass
    
    def start_consuming(self, message_handler):
        """
        Start consuming messages from the queue.
        
        Args:
            message_handler: Callback function that processes messages
                             Should accept (channel, method, properties, body) and return True on success
        """
        if not self.connection or self.connection.is_closed:
            if not self.connect():
                logger.error("Cannot start consuming: connection failed")
                return
        
        try:
            # Set up consumer
            self.channel.basic_consume(
                queue=self.queue_name,
                on_message_callback=self._create_callback(message_handler),
                auto_ack=False  # Manual acknowledgment
            )
            
            logger.info(f"Starting to consume messages from queue: {self.queue_name}")
            logger.info("Waiting for messages. Press CTRL+C to exit.")
            
            # Start consuming (blocking)
            while not self.should_stop:
                try:
                    self.connection.process_data_events(time_limit=1)
                except KeyboardInterrupt:
                    logger.info("Received keyboard interrupt, stopping...")
                    self.should_stop = True
                    break
                except Exception as e:
                    logger.error(f"Error processing data events: {e}", exc_info=True)
                    # Try to reconnect
                    if self.connection.is_closed:
                        logger.info("Connection closed, attempting to reconnect...")
                        if not self.connect():
                            logger.error("Reconnection failed, waiting before retry...")
                            import time
                            time.sleep(5)
            
        except Exception as e:
            logger.error(f"Error in consume loop: {e}", exc_info=True)
        finally:
            self.close()
    
    def _create_callback(self, message_handler):
        """Create a callback wrapper that handles acknowledgment"""
        def callback(channel, method, properties, body):
            try:
                # Parse message
                try:
                    message = json.loads(body.decode('utf-8'))
                except json.JSONDecodeError as e:
                    logger.error(f"Failed to parse message JSON: {e}")
                    channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
                    return
                
                logger.info(f"Received resume generation job: {message.get('id', 'unknown')}")
                
                # Process message
                try:
                    success = message_handler(message)
                    if success:
                        # Acknowledge message on success
                        channel.basic_ack(delivery_tag=method.delivery_tag)
                        logger.info(f"Successfully processed job: {message.get('id', 'unknown')}")
                    else:
                        # Reject and requeue on failure
                        channel.basic_nack(delivery_tag=method.delivery_tag, requeue=True)
                        logger.warning(f"Failed to process job (requeued): {message.get('id', 'unknown')}")
                except Exception as e:
                    logger.error(f"Error processing message: {e}", exc_info=True)
                    # Reject and requeue on exception
                    channel.basic_nack(delivery_tag=method.delivery_tag, requeue=True)
                    
            except Exception as e:
                logger.error(f"Error in callback: {e}", exc_info=True)
                # Reject message without requeue on critical error
                try:
                    channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
                except:
                    pass
        
        return callback
    
    def close(self):
        """Close the connection gracefully"""
        try:
            if self.channel and not self.channel.is_closed:
                self.channel.close()
            if self.connection and not self.connection.is_closed:
                self.connection.close()
            logger.info("RabbitMQ connection closed")
        except Exception as e:
            logger.error(f"Error closing connection: {e}")


def create_consumer_from_env() -> ResumeQueueConsumer:
    """
    Create a ResumeQueueConsumer from environment variables.
    
    Environment variables:
        RABBITMQ_URL: Full RabbitMQ connection URL
        RABBITMQ_USER: RabbitMQ username (if URL not provided)
        RABBITMQ_PASSWORD: RabbitMQ password (if URL not provided)
        RABBITMQ_HOST: RabbitMQ host (if URL not provided)
        RABBITMQ_PORT: RabbitMQ port (if URL not provided, default: 5672)
        RABBITMQ_VHOST: RabbitMQ vhost (if URL not provided, default: /woragis)
        RESUME_QUEUE_NAME: Queue name (default: resumes.queue)
        RESUME_EXCHANGE: Exchange name (default: woragis.tasks)
        RESUME_ROUTING_KEY: Routing key (default: resumes.queue)
    """
    # Try to get full URL first
    rabbitmq_url = os.getenv('RABBITMQ_URL')
    
    if not rabbitmq_url:
        # Build URL from components
        user = os.getenv('RABBITMQ_USER', 'woragis')
        password = os.getenv('RABBITMQ_PASSWORD', 'woragis')
        host = os.getenv('RABBITMQ_HOST', 'rabbitmq')
        port = os.getenv('RABBITMQ_PORT', '5672')
        vhost = os.getenv('RABBITMQ_VHOST', 'woragis')
        # Remove leading slash if present
        vhost = vhost.lstrip('/')
        rabbitmq_url = f"amqp://{user}:{password}@{host}:{port}/{vhost}"
    
    queue_name = os.getenv('RESUME_QUEUE_NAME', 'resumes.queue')
    exchange = os.getenv('RESUME_EXCHANGE', 'woragis.tasks')
    routing_key = os.getenv('RESUME_ROUTING_KEY', 'resumes.queue')
    
    return ResumeQueueConsumer(
        rabbitmq_url=rabbitmq_url,
        queue_name=queue_name,
        exchange=exchange,
        routing_key=routing_key
    )
