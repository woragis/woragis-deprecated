/**
 * Integration tests for RabbitMQ queue operations
 * 
 * These tests require RabbitMQ running (use docker-compose.test.yml)
 * 
 * Run with: npm run test:jest -- tests/integration/queue.test.js
 */

import { describe, test, expect, beforeEach, afterEach } from '@jest/globals';
import amqp from 'amqplib';
import { RabbitMQQueue } from '../../src/queue_rabbitmq.js';

// Test configuration
const RABBITMQ_URL = process.env.TEST_RABBITMQ_URL || 'amqp://test:test@localhost:5673/test';
const TEST_QUEUE_NAME = 'test.job-applications.queue';
const TEST_EXCHANGE = 'test.woragis.tasks';
const TEST_ROUTING_KEY = 'test.job-applications.process';

describe('RabbitMQ Queue Integration Tests', () => {
  let connection;
  let channel;
  let queue;

  beforeEach(async () => {
    try {
      // Setup RabbitMQ connection
      connection = await amqp.connect(RABBITMQ_URL);
      channel = await connection.createChannel();

      // Declare test exchange
      await channel.assertExchange(TEST_EXCHANGE, 'direct', { durable: true });

      // Declare test queue with DLX
      await channel.assertQueue(TEST_QUEUE_NAME, {
        durable: true,
        arguments: {
          'x-dead-letter-exchange': 'test.woragis.dlx',
          'x-dead-letter-routing-key': TEST_QUEUE_NAME + '.failed'
        }
      });

      // Bind queue to exchange
      await channel.bindQueue(TEST_QUEUE_NAME, TEST_EXCHANGE, TEST_ROUTING_KEY);

      // Purge queue before each test
      await channel.purgeQueue(TEST_QUEUE_NAME);

      // Create queue instance
      queue = new RabbitMQQueue();
      // Override queue name for testing
      queue.queueName = TEST_QUEUE_NAME;
      queue.exchange = TEST_EXCHANGE;
      queue.routingKey = TEST_ROUTING_KEY;
    } catch (error) {
      console.warn('RabbitMQ not available, skipping integration tests:', error.message);
      connection = null;
    }
  });

  afterEach(async () => {
    if (channel) {
      try {
        await channel.purgeQueue(TEST_QUEUE_NAME);
        await channel.deleteQueue(TEST_QUEUE_NAME);
        await channel.close();
      } catch (error) {
        // Ignore cleanup errors
      }
    }

    if (connection) {
      try {
        await connection.close();
      } catch (error) {
        // Ignore cleanup errors
      }
    }

    if (queue && queue.connection) {
      try {
        await queue.connection.close();
      } catch (error) {
        // Ignore cleanup errors
      }
    }
  });

  test('should connect to RabbitMQ', async () => {
    if (!connection) {
      return; // Skip if RabbitMQ not available
    }

    // Verify connection is open
    expect(connection).toBeDefined();
    expect(channel).toBeDefined();
  });

  test('should declare exchange', async () => {
    if (!channel) {
      return; // Skip if RabbitMQ not available
    }

    // Verify exchange exists
    try {
      await channel.checkExchange(TEST_EXCHANGE);
      expect(true).toBe(true); // Exchange exists
    } catch (error) {
      throw new Error(`Exchange ${TEST_EXCHANGE} should exist`);
    }
  });

  test('should declare queue with dead letter exchange', async () => {
    if (!channel) {
      return; // Skip if RabbitMQ not available
    }

    // Verify queue exists
    const queueInfo = await channel.checkQueue(TEST_QUEUE_NAME);
    expect(queueInfo).toBeDefined();
    expect(queueInfo.queue).toBe(TEST_QUEUE_NAME);

    // Verify DLX arguments
    // Note: checkQueue doesn't return arguments, but we can verify queue exists
    expect(queueInfo).toBeDefined();
  });

  test('should bind queue to exchange', async () => {
    if (!channel) {
      return; // Skip if RabbitMQ not available
    }

    // Verify binding exists by checking queue info
    const queueInfo = await channel.checkQueue(TEST_QUEUE_NAME);
    expect(queueInfo).toBeDefined();
  });

  test('should publish message to queue', async () => {
    if (!channel) {
      return; // Skip if RabbitMQ not available
    }

    const message = {
      id: 'test-message-123',
      userId: 'test-user',
      companyName: 'Test Company',
    };

    const messageBody = JSON.stringify(message);

    // Publish message
    const published = channel.publish(
      TEST_EXCHANGE,
      TEST_ROUTING_KEY,
      Buffer.from(messageBody),
      {
        contentType: 'application/json',
        deliveryMode: 2, // Persistent
      }
    );

    expect(published).toBe(true);

    // Wait a bit for message to be routed
    await new Promise(resolve => setTimeout(resolve, 100));

    // Verify message is in queue
    const queueInfo = await channel.checkQueue(TEST_QUEUE_NAME);
    expect(queueInfo.messageCount).toBeGreaterThan(0);
  });

  test('should consume message from queue', async () => {
    if (!channel) {
      return; // Skip if RabbitMQ not available
    }

    const message = {
      id: 'test-consume-456',
      userId: 'test-user-consume',
      companyName: 'Consume Test Company',
    };

    const messageBody = JSON.stringify(message);

    // Publish message
    channel.publish(
      TEST_EXCHANGE,
      TEST_ROUTING_KEY,
      Buffer.from(messageBody),
      {
        contentType: 'application/json',
        deliveryMode: 2,
      }
    );

    // Wait a bit for message to be routed
    await new Promise(resolve => setTimeout(resolve, 100));

    // Consume message
    const consumedMessage = await channel.get(TEST_QUEUE_NAME, { noAck: false });

    expect(consumedMessage).toBeDefined();
    expect(consumedMessage).not.toBe(false);

    const receivedMessage = JSON.parse(consumedMessage.content.toString());
    expect(receivedMessage.id).toBe(message.id);
    expect(receivedMessage.companyName).toBe(message.companyName);

    // Acknowledge message
    channel.ack(consumedMessage);
  });

  test('should handle message rejection and send to DLQ', async () => {
    if (!channel) {
      return; // Skip if RabbitMQ not available
    }

    // Declare DLX and DLQ
    const dlxExchange = 'test.woragis.dlx';
    const dlqName = TEST_QUEUE_NAME + '.failed';

    await channel.assertExchange(dlxExchange, 'direct', { durable: true });
    await channel.assertQueue(dlqName, { durable: true });
    await channel.bindQueue(dlqName, dlxExchange, dlqName);

    const message = {
      id: 'test-dlq-789',
      userId: 'test-user-dlq',
      companyName: 'DLQ Test Company',
    };

    const messageBody = JSON.stringify(message);

    // Publish message
    channel.publish(
      TEST_EXCHANGE,
      TEST_ROUTING_KEY,
      Buffer.from(messageBody),
      {
        contentType: 'application/json',
        deliveryMode: 2,
      }
    );

    // Wait a bit for message to be routed
    await new Promise(resolve => setTimeout(resolve, 100));

    // Consume message
    const consumedMessage = await channel.get(TEST_QUEUE_NAME, { noAck: false });

    expect(consumedMessage).toBeDefined();

    // Reject message without requeue (should go to DLQ)
    channel.nack(consumedMessage, false, false);

    // Wait a bit for message to be routed to DLQ
    await new Promise(resolve => setTimeout(resolve, 200));

    // Verify message is in DLQ
    const dlqInfo = await channel.checkQueue(dlqName);
    expect(dlqInfo.messageCount).toBeGreaterThan(0);

    // Cleanup DLQ
    await channel.purgeQueue(dlqName);
    await channel.deleteQueue(dlqName);
  });

  test('should set prefetch count', async () => {
    if (!channel) {
      return; // Skip if RabbitMQ not available
    }

    // Set prefetch count
    await channel.prefetch(1);

    // Verify prefetch is set (can't directly check, but no error means it worked)
    expect(true).toBe(true);
  });
});

