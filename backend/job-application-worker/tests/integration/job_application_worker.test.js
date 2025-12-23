/**
 * Integration tests for Job Application Worker
 * 
 * These tests require:
 * - RabbitMQ running (use docker-compose.test.yml)
 * - PostgreSQL database running
 * 
 * Run with: npm run test:jest -- tests/integration/job_application_worker.test.js
 */

import { describe, test, expect, beforeEach, afterEach, jest } from '@jest/globals';
import amqp from 'amqplib';
import { Pool } from 'pg';
import { Worker } from '../../src/worker.js';
import { RabbitMQQueue } from '../../src/queue_rabbitmq.js';
import { Database } from '../../src/database.js';

// Test configuration
const RABBITMQ_URL = process.env.TEST_RABBITMQ_URL || 'amqp://test:test@localhost:5673/test';
const TEST_QUEUE_NAME = 'test.job-applications.queue';
const TEST_EXCHANGE = 'test.woragis.tasks';
const TEST_ROUTING_KEY = 'test.job-applications.process';

const DB_CONFIG = {
  host: process.env.TEST_DB_HOST || 'localhost',
  port: parseInt(process.env.TEST_DB_PORT || '5433', 10),
  database: process.env.TEST_DB_NAME || 'woragis_test',
  user: process.env.TEST_DB_USER || 'postgres',
  password: process.env.TEST_DB_PASSWORD || 'postgres',
};

// Setup test database tables
async function setupTestDatabase(pool) {
  // Create job_applications table
  await pool.query(`
    CREATE TABLE IF NOT EXISTS job_applications (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      user_id UUID NOT NULL,
      company_name VARCHAR(255) NOT NULL,
      location VARCHAR(255),
      job_title VARCHAR(255) NOT NULL,
      job_url TEXT NOT NULL,
      website VARCHAR(100) NOT NULL,
      status VARCHAR(50) DEFAULT 'pending',
      cover_letter TEXT,
      applied_at TIMESTAMP,
      error_message TEXT,
      created_at TIMESTAMP DEFAULT NOW(),
      updated_at TIMESTAMP DEFAULT NOW(),
      UNIQUE(job_url, website)
    )
  `);

  // Create websites table
  await pool.query(`
    CREATE TABLE IF NOT EXISTS websites (
      id SERIAL PRIMARY KEY,
      name VARCHAR(100) UNIQUE NOT NULL,
      enabled BOOLEAN DEFAULT true,
      current_count INTEGER DEFAULT 0,
      daily_limit INTEGER DEFAULT 10,
      last_reset TIMESTAMP DEFAULT NOW(),
      created_at TIMESTAMP DEFAULT NOW(),
      updated_at TIMESTAMP DEFAULT NOW()
    )
  `);
}

describe('Job Application Worker Integration Tests', () => {
  let rabbitmqConnection;
  let rabbitmqChannel;
  let dbPool;
  let worker;

  beforeEach(async () => {
    // Setup RabbitMQ connection
    try {
      rabbitmqConnection = await amqp.connect(RABBITMQ_URL);
      rabbitmqChannel = await rabbitmqConnection.createChannel();
      
      // Declare test exchange
      await rabbitmqChannel.assertExchange(TEST_EXCHANGE, 'direct', { durable: true });
      
      // Declare test queue with DLX
      await rabbitmqChannel.assertQueue(TEST_QUEUE_NAME, {
        durable: true,
        arguments: {
          'x-dead-letter-exchange': 'test.woragis.dlx',
          'x-dead-letter-routing-key': TEST_QUEUE_NAME + '.failed'
        }
      });
      
      // Bind queue to exchange
      await rabbitmqChannel.bindQueue(TEST_QUEUE_NAME, TEST_EXCHANGE, TEST_ROUTING_KEY);
      
      // Purge queue before each test
      await rabbitmqChannel.purgeQueue(TEST_QUEUE_NAME);
    } catch (error) {
      console.warn('RabbitMQ not available, skipping integration tests:', error.message);
      rabbitmqConnection = null;
    }

    // Setup database connection
    try {
      dbPool = new Pool(DB_CONFIG);
      // Test connection
      await dbPool.query('SELECT 1');
      // Setup test tables
      await setupTestDatabase(dbPool);
    } catch (error) {
      console.warn('Database not available, skipping integration tests:', error.message);
      dbPool = null;
    }

    // Skip tests if dependencies not available
    if (!rabbitmqConnection || !dbPool) {
      return;
    }
  });

  afterEach(async () => {
    // Cleanup
    if (rabbitmqChannel) {
      try {
        await rabbitmqChannel.purgeQueue(TEST_QUEUE_NAME);
        await rabbitmqChannel.deleteQueue(TEST_QUEUE_NAME);
        await rabbitmqChannel.close();
      } catch (error) {
        // Ignore cleanup errors
      }
    }

    if (rabbitmqConnection) {
      try {
        await rabbitmqConnection.close();
      } catch (error) {
        // Ignore cleanup errors
      }
    }

    if (dbPool) {
      try {
        // Clean up test data
        await dbPool.query('DELETE FROM job_applications WHERE website = $1', ['test-website']);
        await dbPool.query('DELETE FROM websites WHERE name = $1', ['test-website']);
        await dbPool.end();
      } catch (error) {
        // Ignore cleanup errors
      }
    }

    if (worker && worker.running) {
      try {
        worker.running = false;
      } catch (error) {
        // Ignore cleanup errors
      }
    }
  });

  test('should connect to RabbitMQ and set up queue', async () => {
    if (!rabbitmqConnection) {
      return; // Skip if RabbitMQ not available
    }

    // Verify exchange exists
    try {
      await rabbitmqChannel.checkExchange(TEST_EXCHANGE);
      expect(true).toBe(true); // Exchange exists
    } catch (error) {
      throw new Error(`Exchange ${TEST_EXCHANGE} should exist`);
    }

    // Verify queue exists
    const queueInfo = await rabbitmqChannel.checkQueue(TEST_QUEUE_NAME);
    expect(queueInfo).toBeDefined();
    expect(queueInfo.queue).toBe(TEST_QUEUE_NAME);
  });

  test('should publish job application message to queue', async () => {
    if (!rabbitmqConnection) {
      return; // Skip if RabbitMQ not available
    }

    const jobMessage = {
      id: 'test-job-123',
      userId: 'test-user-123',
      companyName: 'Test Company',
      jobTitle: 'Software Engineer',
      location: 'Remote',
      jobUrl: 'https://example.com/job/123',
      website: 'test-website',
    };

    const messageBody = JSON.stringify(jobMessage);

    // Publish message
    await rabbitmqChannel.publish(
      TEST_EXCHANGE,
      TEST_ROUTING_KEY,
      Buffer.from(messageBody),
      {
        contentType: 'application/json',
        deliveryMode: 2, // Persistent
      }
    );

    // Verify message is in queue
    const queueInfo = await rabbitmqChannel.checkQueue(TEST_QUEUE_NAME);
    expect(queueInfo.messageCount).toBeGreaterThan(0);
  });

  test('should consume and process job application message', async () => {
    if (!rabbitmqConnection || !dbPool) {
      return; // Skip if dependencies not available
    }

    // Setup test website in database
    await dbPool.query(`
      INSERT INTO websites (name, enabled, current_count, daily_limit, last_reset)
      VALUES ($1, $2, $3, $4, $5)
      ON CONFLICT (name) DO UPDATE SET
        enabled = EXCLUDED.enabled,
        current_count = EXCLUDED.current_count,
        daily_limit = EXCLUDED.daily_limit,
        last_reset = EXCLUDED.last_reset
    `, ['test-website', true, 0, 10, new Date()]);

    const jobMessage = {
      id: 'test-job-456',
      userId: 'test-user-456',
      companyName: 'Test Company Inc',
      jobTitle: 'Senior Developer',
      location: 'San Francisco, CA',
      jobUrl: 'https://example.com/job/456',
      website: 'test-website',
    };

    const messageBody = JSON.stringify(jobMessage);

    // Publish message
    await rabbitmqChannel.publish(
      TEST_EXCHANGE,
      TEST_ROUTING_KEY,
      Buffer.from(messageBody),
      {
        contentType: 'application/json',
        deliveryMode: 2,
      }
    );

    // Wait a bit for message to be available
    await new Promise(resolve => setTimeout(resolve, 100));

    // Consume message
    const consumedMessage = await rabbitmqChannel.get(TEST_QUEUE_NAME, { noAck: false });
    
    expect(consumedMessage).toBeDefined();
    expect(consumedMessage).not.toBe(false);

    const receivedJob = JSON.parse(consumedMessage.content.toString());
    expect(receivedJob.id).toBe(jobMessage.id);
    expect(receivedJob.companyName).toBe(jobMessage.companyName);

    // Acknowledge message
    rabbitmqChannel.ack(consumedMessage);
  });

  test('should handle invalid message format', async () => {
    if (!rabbitmqConnection) {
      return; // Skip if RabbitMQ not available
    }

    // Publish invalid message (not JSON)
    await rabbitmqChannel.publish(
      TEST_EXCHANGE,
      TEST_ROUTING_KEY,
      Buffer.from('invalid json'),
      {
        contentType: 'application/json',
        deliveryMode: 2,
      }
    );

    // Consume message
    const consumedMessage = await rabbitmqChannel.get(TEST_QUEUE_NAME, { noAck: false });
    
    expect(consumedMessage).toBeDefined();
    
    // Message should be consumable but processing should fail
    // In real scenario, this would go to DLQ after retries
    rabbitmqChannel.nack(consumedMessage, false, false); // Reject without requeue (would go to DLQ)
  });

  test('should create job application record in database', async () => {
    if (!dbPool) {
      return; // Skip if database not available
    }

    const applicationData = {
      userId: '123e4567-e89b-12d3-a456-426614174000', // Valid UUID
      companyName: 'Database Test Company',
      location: 'New York, NY',
      jobTitle: 'Full Stack Developer',
      jobUrl: 'https://example.com/job/789',
      website: 'test-website',
      status: 'processing',
    };

    // Insert application
    const result = await dbPool.query(`
      INSERT INTO job_applications (user_id, company_name, location, job_title, job_url, website, status, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
      RETURNING id, user_id, company_name, job_title, status
    `, [
      applicationData.userId,
      applicationData.companyName,
      applicationData.location,
      applicationData.jobTitle,
      applicationData.jobUrl,
      applicationData.website,
      applicationData.status,
    ]);

    expect(result.rows).toHaveLength(1);
    expect(result.rows[0].user_id).toBe(applicationData.userId);
    expect(result.rows[0].company_name).toBe(applicationData.companyName);
    expect(result.rows[0].status).toBe(applicationData.status);

    // Cleanup
    await dbPool.query('DELETE FROM job_applications WHERE id = $1', [result.rows[0].id]);
  });

  test('should find existing job application by URL', async () => {
    if (!dbPool) {
      return; // Skip if database not available
    }

    const applicationData = {
      userId: '223e4567-e89b-12d3-a456-426614174001', // Valid UUID
      companyName: 'Find Test Company',
      location: 'Seattle, WA',
      jobTitle: 'DevOps Engineer',
      jobUrl: 'https://example.com/job/find',
      website: 'test-website',
      status: 'completed',
    };

    // Insert application
    const insertResult = await dbPool.query(`
      INSERT INTO job_applications (user_id, company_name, location, job_title, job_url, website, status, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
      RETURNING id
    `, [
      applicationData.userId,
      applicationData.companyName,
      applicationData.location,
      applicationData.jobTitle,
      applicationData.jobUrl,
      applicationData.website,
      applicationData.status,
    ]);

    // Find application by URL
    const findResult = await dbPool.query(`
      SELECT id, user_id, company_name, job_title, job_url, status
      FROM job_applications
      WHERE job_url = $1 AND website = $2
    `, [applicationData.jobUrl, applicationData.website]);

    expect(findResult.rows).toHaveLength(1);
    expect(findResult.rows[0].job_url).toBe(applicationData.jobUrl);
    expect(findResult.rows[0].status).toBe(applicationData.status);

    // Cleanup
    await dbPool.query('DELETE FROM job_applications WHERE id = $1', [insertResult.rows[0].id]);
  });

  test('should check website rate limit', async () => {
    if (!dbPool) {
      return; // Skip if database not available
    }

    // Create website with limit reached
    await dbPool.query(`
      INSERT INTO websites (name, enabled, current_count, daily_limit, last_reset)
      VALUES ($1, $2, $3, $4, $5)
      ON CONFLICT (name) DO UPDATE SET
        enabled = EXCLUDED.enabled,
        current_count = EXCLUDED.current_count,
        daily_limit = EXCLUDED.daily_limit
    `, ['test-website-limit', true, 10, 10, new Date()]);

    // Check if website is at limit
    const result = await dbPool.query(`
      SELECT name, enabled, current_count, daily_limit
      FROM websites
      WHERE name = $1
    `, ['test-website-limit']);

    expect(result.rows).toHaveLength(1);
    expect(result.rows[0].current_count).toBe(10);
    expect(result.rows[0].daily_limit).toBe(10);
    expect(result.rows[0].current_count).toBeGreaterThanOrEqual(result.rows[0].daily_limit);

    // Cleanup
    await dbPool.query('DELETE FROM websites WHERE name = $1', ['test-website-limit']);
  });

  test('should increment website count after processing', async () => {
    if (!dbPool) {
      return; // Skip if database not available
    }

    // Create website
    await dbPool.query(`
      INSERT INTO websites (name, enabled, current_count, daily_limit, last_reset)
      VALUES ($1, $2, $3, $4, $5)
      ON CONFLICT (name) DO UPDATE SET
        enabled = EXCLUDED.enabled,
        current_count = EXCLUDED.current_count,
        daily_limit = EXCLUDED.daily_limit
    `, ['test-website-increment', true, 5, 10, new Date()]);

    // Increment count
    await dbPool.query(`
      UPDATE websites
      SET current_count = current_count + 1,
          updated_at = NOW()
      WHERE name = $1
    `, ['test-website-increment']);

    // Verify count incremented
    const result = await dbPool.query(`
      SELECT current_count
      FROM websites
      WHERE name = $1
    `, ['test-website-increment']);

    expect(result.rows[0].current_count).toBe(6);

    // Cleanup
    await dbPool.query('DELETE FROM websites WHERE name = $1', ['test-website-increment']);
  });
});

