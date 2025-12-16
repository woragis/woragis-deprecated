import amqp from 'amqplib';
import { logger } from './utils/logger.js';

export class RabbitMQQueue {
  constructor() {
    this.connection = null;
    this.channel = null;
    this.queueName = process.env.JOB_APPLICATION_QUEUE_NAME || 'job-applications.queue';
    this.exchange = process.env.JOB_APPLICATION_EXCHANGE || 'woragis.tasks';
    this.routingKey = process.env.JOB_APPLICATION_ROUTING_KEY || 'job-applications.process';
  }

  async connect() {
    const rabbitmqUrl = process.env.RABBITMQ_URL || 
      `amqp://${process.env.RABBITMQ_USER || 'woragis'}:${process.env.RABBITMQ_PASSWORD || 'woragis'}@${process.env.RABBITMQ_HOST || 'rabbitmq'}:${process.env.RABBITMQ_PORT || '5672'}/${process.env.RABBITMQ_VHOST || 'woragis'}`;
    
    try {
      this.connection = await amqp.connect(rabbitmqUrl);
      
      this.connection.on('error', (err) => {
        logger.error('RabbitMQ connection error', { error: err.message });
      });

      this.connection.on('close', () => {
        logger.warn('RabbitMQ connection closed');
      });

      this.channel = await this.connection.createChannel();

      // Declare exchange
      await this.channel.assertExchange(this.exchange, 'direct', {
        durable: true
      });

      // Declare queue with dead letter exchange
      await this.channel.assertQueue(this.queueName, {
        durable: true,
        arguments: {
          'x-dead-letter-exchange': 'woragis.dlx',
          'x-dead-letter-routing-key': this.queueName + '.failed'
        }
      });

      // Bind queue to exchange
      await this.channel.bindQueue(this.queueName, this.exchange, this.routingKey);

      // Set QoS to process one message at a time
      const prefetchCount = parseInt(process.env.JOB_APPLICATION_PREFETCH_COUNT || '1', 10);
      await this.channel.prefetch(prefetchCount);

      logger.info('Connected to RabbitMQ', {
        exchange: this.exchange,
        queue: this.queueName,
        routingKey: this.routingKey
      });
      
      // Expose connection for health checks
      try {
        const { setRabbitMQConnection } = await import('./health.js');
        setRabbitMQConnection(this.connection);
      } catch (err) {
        // Health module may not be available
      }
    } catch (error) {
      logger.error('Failed to connect to RabbitMQ', { error: error.message });
      throw error;
    }
  }

  async consume(handler) {
    if (!this.channel) {
      throw new Error('Not connected to RabbitMQ');
    }

    logger.info('Starting to consume job application messages', {
      queue: this.queueName
    });

    await this.channel.consume(
      this.queueName,
      async (msg) => {
        if (!msg) {
          return;
        }

        try {
          const job = JSON.parse(msg.content.toString());
          logger.info('Received job application', {
            jobId: job.id,
            company: job.companyName,
            website: job.website
          });

          // Process the job
          await handler(job);

          // Acknowledge successful processing
          this.channel.ack(msg);
          logger.info('Job application processed successfully', {
            jobId: job.id
          });
        } catch (error) {
          logger.error('Error processing job application', {
            error: error.message,
            stack: error.stack
          });

          // Reject and requeue for retry
          this.channel.nack(msg, false, true);
        }
      },
      {
        noAck: false // Manual acknowledgment
      }
    );
  }

  async disconnect() {
    if (this.channel) {
      await this.channel.close();
      this.channel = null;
    }
    if (this.connection) {
      await this.connection.close();
      this.connection = null;
    }
    logger.info('Disconnected from RabbitMQ');
  }

  // Legacy methods for compatibility (not used with RabbitMQ)
  async enqueueJob(job) {
    // Not used - server publishes directly
    throw new Error('enqueueJob not supported - server publishes directly to RabbitMQ');
  }

  async dequeueJob(timeout = 5000) {
    // Not used - use consume() instead
    throw new Error('dequeueJob not supported - use consume() instead');
  }

  async getJob(jobId) {
    // Not used - job data is in message body
    throw new Error('getJob not supported - job data is in message body');
  }

  async markJobComplete(jobId) {
    // Not used - acknowledgment is handled in consume()
    // This is a no-op for compatibility
    return;
  }

  async markJobFailed(jobId, errorMsg) {
    // Not used - failed jobs go to DLQ automatically
    // This is a no-op for compatibility
    logger.error('Job failed', { jobId, errorMsg });
    return;
  }
}
