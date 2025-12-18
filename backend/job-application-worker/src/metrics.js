/**
 * Prometheus metrics for Job Application Worker.
 */
import promClient from 'prom-client';

// Create a Registry to register the metrics
const register = new promClient.Registry();

// Add default metrics (CPU, memory, etc.)
promClient.collectDefaultMetrics({ register });

// Worker metrics
const workerJobsProcessedTotal = new promClient.Counter({
  name: 'worker_jobs_processed_total',
  help: 'Total number of jobs processed',
  labelNames: ['worker', 'status'], // status: success, failed
  registers: [register],
});

const workerJobDurationSeconds = new promClient.Histogram({
  name: 'worker_job_duration_seconds',
  help: 'Job processing duration in seconds',
  labelNames: ['worker'],
  buckets: [.01, .05, .1, .25, .5, 1, 2.5, 5, 10, 30],
  registers: [register],
});

const workerJobsFailedTotal = new promClient.Counter({
  name: 'worker_jobs_failed_total',
  help: 'Total number of failed jobs',
  labelNames: ['worker', 'error_type'],
  registers: [register],
});

const workerJobsRetriedTotal = new promClient.Counter({
  name: 'worker_jobs_retried_total',
  help: 'Total number of retried jobs',
  labelNames: ['worker'],
  registers: [register],
});

// Queue metrics
const queueDepth = new promClient.Gauge({
  name: 'queue_depth',
  help: 'Current queue depth',
  labelNames: ['queue_name'],
  registers: [register],
});

const queueDLQSize = new promClient.Gauge({
  name: 'queue_dlq_size',
  help: 'Dead letter queue size',
  labelNames: ['queue_name'],
  registers: [register],
});

/**
 * Record a processed job metric
 */
export const recordJobProcessed = (worker, status, duration) => {
  workerJobsProcessedTotal.labels({ worker, status }).inc();
  workerJobDurationSeconds.labels({ worker }).observe(duration);
};

/**
 * Record a failed job metric
 */
export const recordJobFailed = (worker, errorType) => {
  workerJobsFailedTotal.labels({ worker, error_type: errorType }).inc();
};

/**
 * Record a retried job metric
 */
export const recordJobRetried = (worker) => {
  workerJobsRetriedTotal.labels({ worker }).inc();
};

/**
 * Set the queue depth gauge
 */
export const setQueueDepth = (queueName, depth) => {
  queueDepth.labels({ queue_name: queueName }).set(depth);
};

/**
 * Set the DLQ size gauge
 */
export const setQueueDLQSize = (queueName, size) => {
  queueDLQSize.labels({ queue_name: queueName }).set(size);
};

/**
 * Get metrics in Prometheus format
 */
export const getMetrics = async () => {
  return await register.metrics();
};
