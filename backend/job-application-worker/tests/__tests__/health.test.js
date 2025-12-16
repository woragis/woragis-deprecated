import { describe, test, expect, beforeEach } from '@jest/globals';
import { checkHealth, setRabbitMQConnection } from '../../src/health.js';

describe('Health Check', () => {
  beforeEach(() => {
    // Reset connection before each test
    setRabbitMQConnection(null);
  });

  test('should return healthy status when RabbitMQ is connected', () => {
    const mockConnection = {
      connection: { readyState: 'open' }
    };
    setRabbitMQConnection(mockConnection);

    const result = checkHealth();

    expect(result.status).toBe('healthy');
    expect(result.checks).toHaveLength(2);
    expect(result.checks[0].name).toBe('service');
    expect(result.checks[0].status).toBe('ok');
    expect(result.checks[1].name).toBe('rabbitmq');
    expect(result.checks[1].status).toBe('ok');
  });

  test('should return unhealthy status when RabbitMQ is not configured', () => {
    setRabbitMQConnection(null);

    const result = checkHealth();

    expect(result.status).toBe('unhealthy');
    expect(result.checks[1].name).toBe('rabbitmq');
    expect(result.checks[1].status).toBe('error');
    expect(result.checks[1].message).toBe('not configured');
  });

  test('should return unhealthy status when RabbitMQ connection is not initialized', () => {
    const mockConnection = {
      connection: null
    };
    setRabbitMQConnection(mockConnection);

    const result = checkHealth();

    expect(result.status).toBe('unhealthy');
    expect(result.checks[1].name).toBe('rabbitmq');
    expect(result.checks[1].status).toBe('error');
    expect(result.checks[1].message).toBe('connection not initialized');
  });

  test('should include service check', () => {
    const result = checkHealth();

    expect(result.checks).toHaveLength(2);
    expect(result.checks[0].name).toBe('service');
    expect(result.checks[0].status).toBe('ok');
  });
});
