/**
 * Performance tests for job-application-worker
 * Run with: npm test -- tests/performance.test.js
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import amqp from 'amqplib';

// Note: This test requires RabbitMQ to be running
// Set RABBITMQ_URL environment variable if needed

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

    // Consume messages
    const consumeStart = Date.now();
    let consumedCount = 0;

    await new Promise((resolve) => {
      channel.consume(queueName, (msg) => {
        if (msg) {
          channel.ack(msg);
          consumedCount++;
          if (consumedCount >= concurrentMessages) {
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

    expect(consumedCount).toBe(concurrentMessages);
    expect(totalDuration).toBeLessThan(5000); // Should complete in under 5 seconds
  }, 30000);

  test('should measure message processing latency', async () => {
    const iterations = 50;
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
        channel.consume(queueName, (msg) => {
          if (msg) {
            const latency = Date.now() - start;
            latencies.push(latency);
            channel.ack(msg);
            resolve();
          }
        }, { noAck: false });
      });
    }

    const avgLatency = latencies.reduce((a, b) => a + b, 0) / latencies.length;
    const sortedLatencies = [...latencies].sort((a, b) => a - b);
    const p95Index = Math.floor(latencies.length * 0.95);
    const p95Latency = sortedLatencies[p95Index];

    console.log('Latency Test Results:');
    console.log(`  Iterations: ${iterations}`);
    console.log(`  Average Latency: ${avgLatency.toFixed(2)}ms`);
    console.log(`  P95 Latency: ${p95Latency}ms`);

    expect(avgLatency).toBeLessThan(100); // Average should be under 100ms
    expect(p95Latency).toBeLessThan(200); // P95 should be under 200ms
  }, 30000);
});
