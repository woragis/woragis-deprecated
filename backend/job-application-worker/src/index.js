import 'dotenv/config';
import http from 'http';
import { Worker } from './worker.js';
import { logger } from './utils/logger.js';
import { checkHealth } from './health.js';

const worker = new Worker();

// Start health check HTTP server
const healthServer = http.createServer((req, res) => {
  if (req.url === '/healthz' && req.method === 'GET') {
    const result = checkHealth();
    const statusCode = result.status === 'unhealthy' ? 503 : 200;
    
    res.writeHead(statusCode, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify(result));
  } else {
    res.writeHead(404);
    res.end();
  }
});

healthServer.listen(8080, '0.0.0.0', () => {
  logger.info('Health check server started on port 8080');
});

// Graceful shutdown
process.on('SIGTERM', async () => {
  logger.info('SIGTERM received, shutting down gracefully');
  await worker.stop();
  process.exit(0);
});

process.on('SIGINT', async () => {
  logger.info('SIGINT received, shutting down gracefully');
  await worker.stop();
  process.exit(0);
});

// Start worker
worker.start().catch((error) => {
  logger.error('Worker failed to start', { error: error.message });
  process.exit(1);
});

