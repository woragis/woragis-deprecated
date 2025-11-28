import 'dotenv/config';
import { Worker } from './worker.js';
import { logger } from './utils/logger.js';

const worker = new Worker();

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.info('SIGTERM received, shutting down gracefully');
  worker.stop();
  process.exit(0);
});

process.on('SIGINT', () => {
  logger.info('SIGINT received, shutting down gracefully');
  worker.stop();
  process.exit(0);
});

// Start worker
worker.start().catch((error) => {
  logger.error('Worker failed to start', { error: error.message });
  process.exit(1);
});

