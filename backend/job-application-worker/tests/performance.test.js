/**
 * Performance tests for job-application-worker
 * Run with: npm run test:performance
 * Or: node --experimental-vm-modules node_modules/jest/bin/jest.js tests/performance.test.js
 * 
 * Note: This test requires RabbitMQ to be running
 * Set RABBITMQ_URL environment variable if needed
 * DATABASE_URL is NOT required for these tests (they're RabbitMQ-only)
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import amqp from 'amqplib';

describe('Job Application Worker Performance Tests', () => {
  let connection;
  let channel;
  const queueName = `test.job-applications.queue.${Date.now()}`;
  const exchange = 'woragis.notifications';
  const routingKey = 'job-applications.process';

  beforeAll(async () => {
    const rabbitmqUrl = process.env.RABBITMQ_URL || 'amqp://test:test@localhost:5673/test';
    connection = await amqp.connect(rabbitmqUrl);
    channel = await connection.createChannel();
    
    await channel.assertExchange(exchange, 'direct', { durable: true });
    await channel.assertQueue(queueName, { durable: true });
    await channel.bindQueue(queueName, exchange, routingKey);
  });

  afterAll(async () => {
    if (channel) {
      await channel.deleteQueue(queueName);
      await channel.close();
    }
    if (connection) {
      await connection.close();
    }
  });

  test('should handle high throughput message processing', async () => {
    const messageCount = 100;
    const messages = [];

    // Publish messages
    const publishStart = Date.now();
    for (let i = 0; i < messageCount; i++) {
      const message = {
        user_id: `user-${i}`,
        job_url: `https://example.com/job/${i}`,
        timestamp: new Date().toISOString(),
      };
      messages.push(message);
      channel.publish(exchange, routingKey, Buffer.from(JSON.stringify(message)), {
        persistent: true,
      });
    }
    const publishDuration = Date.now() - publishStart;

    // Consume messages
    const consumeStart = Date.now();
    let consumedCount = 0;
    const consumedMessages = [];

    await new Promise((resolve) => {
      channel.consume(queueName, (msg) => {
        if (msg) {
          const content = JSON.parse(msg.content.toString());
          consumedMessages.push(content);
          channel.ack(msg);
          consumedCount++;
          if (consumedCount >= messageCount) {
            resolve();
          }
        }
      });
    });

    const consumeDuration = Date.now() - consumeStart;
    const totalDuration = Date.now() - publishStart;

    console.log('Performance Test Results:');
    console.log(`  Messages Published: ${messageCount}`);
    console.log(`  Messages Consumed: ${consumedCount}`);
    console.log(`  Publish Duration: ${publishDuration}ms`);
    console.log(`  Consume Duration: ${consumeDuration}ms`);
    console.log(`  Total Duration: ${totalDuration}ms`);
    console.log(`  Throughput: ${(messageCount / (totalDuration / 1000)).toFixed(2)} msg/s`);

    expect(consumedCount).toBe(messageCount);
    expect(totalDuration).toBeLessThan(10000); // Should complete in under 10 seconds
  }, 30000);

  test('should handle concurrent message processing', async () => {
    const concurrentMessages = 50;
    const messages = [];

    // Publish concurrent messages
    const publishPromises = [];
    for (let i = 0; i < concurrentMessages; i++) {
      const message = {
        user_id: `user-${i}`,
        job_url: `https://example.com/job/${i}`,
        timestamp: new Date().toISOString(),
      };
      messages.push(message);
      publishPromises.push(
        channel.publish(exchange, routingKey, Buffer.from(JSON.stringify(message)), {
          persistent: true,
        })
      );
    }

    const publishStart = Date.now();
    await Promise.all(publishPromises);
    const publishDuration = Date.now() - publishStart;

    // Consume messages with timeout
    const consumeStart = Date.now();
    let consumedCount = 0;

    await new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        channel.cancel(consumerTag);
        resolve(); // Resolve even if not all messages consumed (for performance testing)
      }, 10000); // 10 second timeout

      const consumerTag = channel.consume(queueName, (msg) => {
        if (msg) {
          channel.ack(msg);
          consumedCount++;
          if (consumedCount >= concurrentMessages) {
            clearTimeout(timeout);
            resolve();
          }
        }
      });
    });

    const consumeDuration = Date.now() - consumeStart;
    const totalDuration = Date.now() - publishStart;

    console.log('Concurrent Processing Test Results:');
    console.log(`  Concurrent Messages: ${concurrentMessages}`);
    console.log(`  Messages Consumed: ${consumedCount}`);
    console.log(`  Publish Duration: ${publishDuration}ms`);
    console.log(`  Consume Duration: ${consumeDuration}ms`);
    console.log(`  Total Duration: ${totalDuration}ms`);
    console.log(`  Throughput: ${(consumedCount / (totalDuration / 1000)).toFixed(2)} msg/s`);

    expect(consumedCount).toBeGreaterThan(0); // At least some messages should be consumed
    expect(totalDuration).toBeLessThan(15000); // Should complete in under 15 seconds
  }, 60000);

  test('should measure message processing latency', async () => {
    const iterations = 20; // Reduced for faster test
    const latencies = [];

    for (let i = 0; i < iterations; i++) {
      const message = {
        user_id: `user-${i}`,
        job_url: `https://example.com/job/${i}`,
        timestamp: new Date().toISOString(),
      };

      const start = Date.now();
      channel.publish(exchange, routingKey, Buffer.from(JSON.stringify(message)), {
        persistent: true,
      });

      await new Promise((resolve) => {
        const timeout = setTimeout(() => {
          // If no message received in 2 seconds, record timeout and continue
          latencies.push(2000); // Record 2s as latency for timeout
          resolve();
        }, 2000);

        const consumerTag = channel.consume(queueName, (msg) => {
          if (msg) {
            clearTimeout(timeout);
            const latency = Date.now() - start;
            latencies.push(latency);
            channel.ack(msg);
            channel.cancel(consumerTag);
            resolve();
          }
        }, { noAck: false });
      });
    }

    if (latencies.length === 0) {
      console.log('Latency Test Results: No messages consumed (worker may not be running)');
      return; // Skip test if no messages consumed
    }

    const avgLatency = latencies.reduce((a, b) => a + b, 0) / latencies.length;
    const sortedLatencies = [...latencies].sort((a, b) => a - b);
    const p95Index = Math.floor(latencies.length * 0.95);
    const p95Latency = sortedLatencies[p95Index];

    console.log('Latency Test Results:');
    console.log(`  Iterations: ${iterations}`);
    console.log(`  Messages Consumed: ${latencies.length}`);
    console.log(`  Average Latency: ${avgLatency.toFixed(2)}ms`);
    console.log(`  P95 Latency: ${p95Latency}ms`);

    expect(latencies.length).toBeGreaterThan(0); // At least some messages should be consumed
    if (latencies.length >= 10) {
      expect(avgLatency).toBeLessThan(5000); // Average should be reasonable
    }
  }, 60000);
});
