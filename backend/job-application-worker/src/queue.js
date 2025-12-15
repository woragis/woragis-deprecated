import Redis from 'ioredis';
import { logger } from './utils/logger.js';

const QUEUE_KEY = 'job-applications:queue';
const JOB_PREFIX = 'job-applications:job:';

export class Queue {
  constructor() {
    this.client = null;
  }

  async connect() {
    const redisUrl = process.env.REDIS_URL || 'redis://localhost:6379/0';
    this.client = new Redis(redisUrl);
    
    this.client.on('error', (err) => {
      logger.error('Redis error', { error: err.message });
    });

    this.client.on('connect', () => {
      logger.info('Connected to Redis', { url: redisUrl });
    });

    // Test connection
    await this.client.ping();
  }

  async enqueueJob(job) {
    if (!job.id) {
      job.id = this.generateId();
    }

    const jobKey = JOB_PREFIX + job.id;
    const jobData = JSON.stringify(job);

    // Store job data with 24 hour TTL
    await this.client.setex(jobKey, 86400, jobData);

    // Add to queue
    await this.client.lpush(QUEUE_KEY, job.id);
  }

  async dequeueJob(timeout = 5000) {
    // Blocking pop from queue
    const result = await this.client.brpop(QUEUE_KEY, timeout / 1000);
    
    if (!result || result.length < 2) {
      return null;
    }

    const jobId = result[1];
    return this.getJob(jobId);
  }

  async getJob(jobId) {
    const jobKey = JOB_PREFIX + jobId;
    const jobData = await this.client.get(jobKey);
    
    if (!jobData) {
      return null;
    }

    return JSON.parse(jobData);
  }

  async markJobComplete(jobId) {
    const jobKey = JOB_PREFIX + jobId;
    await this.client.del(jobKey);
  }

  async markJobFailed(jobId, errorMsg) {
    // For now, just log the failure
    // Could store in a failed jobs list
    logger.error('Job failed', { jobId, errorMsg });
  }

  async disconnect() {
    if (this.client) {
      await this.client.quit();
    }
  }

  generateId() {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }
}

