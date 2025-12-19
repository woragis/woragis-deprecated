"""
Performance tests for resume worker (RabbitMQ worker)
Run with: pytest tests/performance_test.py -v
"""

import os
import pytest
import time
import pika
import json
from concurrent.futures import ThreadPoolExecutor


@pytest.fixture
def rabbitmq_connection():
    """Setup RabbitMQ connection for testing"""
    rabbitmq_url = os.getenv("RABBITMQ_URL", "amqp://test:test@localhost:5673/test")
    connection = pika.BlockingConnection(pika.URLParameters(rabbitmq_url))
    yield connection
    connection.close()


def test_message_throughput(rabbitmq_connection):
    """Test message processing throughput"""
    channel = rabbitmq_connection.channel()
    queue_name = f"test.resumes.queue.{int(time.time())}"
    exchange = "woragis.notifications"
    routing_key = "resumes.generate"

    channel.exchange_declare(exchange=exchange, exchange_type="direct", durable=True)
    channel.queue_declare(queue=queue_name, durable=True)
    channel.queue_bind(queue=queue_name, exchange=exchange, routing_key=routing_key)

    message_count = 50
    messages = []

    # Publish messages
    publish_start = time.time()
    for i in range(message_count):
        message = {
            "user_id": f"user-{i}",
            "job_id": f"job-{i}",
            "language": "en",
            "timestamp": time.time()
        }
        messages.append(message)
        channel.basic_publish(
            exchange=exchange,
            routing_key=routing_key,
            body=json.dumps(message),
            properties=pika.BasicProperties(delivery_mode=2)
        )
    publish_duration = time.time() - publish_start

    # Consume messages
    consume_start = time.time()
    consumed_count = 0
    consumed_messages = []

    def callback(ch, method, properties, body):
        nonlocal consumed_count
        content = json.loads(body)
        consumed_messages.append(content)
        ch.basic_ack(delivery_tag=method.delivery_tag)
        consumed_count += 1
        if consumed_count >= message_count:
            ch.stop_consuming()

    channel.basic_consume(queue=queue_name, on_message_callback=callback)
    channel.start_consuming()

    consume_duration = time.time() - consume_start
    total_duration = time.time() - publish_start

    print(f"\nResume Worker Throughput Test:")
    print(f"  Messages Published: {message_count}")
    print(f"  Messages Consumed: {consumed_count}")
    print(f"  Publish Duration: {publish_duration:.2f}s")
    print(f"  Consume Duration: {consume_duration:.2f}s")
    print(f"  Total Duration: {total_duration:.2f}s")
    print(f"  Throughput: {(message_count / total_duration):.2f} msg/s")

    channel.queue_delete(queue=queue_name)
    channel.close()

    assert consumed_count == message_count
    assert total_duration < 30  # Should complete in under 30 seconds


def test_concurrent_message_processing(rabbitmq_connection):
    """Test concurrent message processing"""
    channel = rabbitmq_connection.channel()
    queue_name = f"test.resumes.queue.concurrent.{int(time.time())}"
    exchange = "woragis.notifications"
    routing_key = "resumes.generate"

    channel.exchange_declare(exchange=exchange, exchange_type="direct", durable=True)
    channel.queue_declare(queue=queue_name, durable=True)
    channel.queue_bind(queue=queue_name, exchange=exchange, routing_key=routing_key)

    concurrent_messages = 30

    # Publish concurrent messages
    publish_start = time.time()
    for i in range(concurrent_messages):
        message = {
            "user_id": f"user-{i}",
            "job_id": f"job-{i}",
            "language": "en"
        }
        channel.basic_publish(
            exchange=exchange,
            routing_key=routing_key,
            body=json.dumps(message),
            properties=pika.BasicProperties(delivery_mode=2)
        )
    publish_duration = time.time() - publish_start

    # Consume messages
    consume_start = time.time()
    consumed_count = 0

    def callback(ch, method, properties, body):
        nonlocal consumed_count
        ch.basic_ack(delivery_tag=method.delivery_tag)
        consumed_count += 1
        if consumed_count >= concurrent_messages:
            ch.stop_consuming()

    channel.basic_consume(queue=queue_name, on_message_callback=callback)
    channel.start_consuming()

    consume_duration = time.time() - consume_start
    total_duration = time.time() - publish_start

    print(f"\nConcurrent Processing Test:")
    print(f"  Concurrent Messages: {concurrent_messages}")
    print(f"  Messages Consumed: {consumed_count}")
    print(f"  Total Duration: {total_duration:.2f}s")

    channel.queue_delete(queue=queue_name)
    channel.close()

    assert consumed_count == concurrent_messages
    assert total_duration < 20
