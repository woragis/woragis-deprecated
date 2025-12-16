/**
 * Structured logging utility for Job Application Worker.
 * Provides JSON logging with service name and trace_id support.
 */

// Context for trace_id (simple in-memory for now, can be enhanced with AsyncLocalStorage)
let traceId = null;

/**
 * Set trace_id for distributed tracing
 */
export const setTraceId = (id) => {
  traceId = id;
};

/**
 * Get trace_id from context
 */
export const getTraceId = () => {
  return traceId;
};

/**
 * Create a log entry with standard fields
 */
const createLogEntry = (level, message, meta = {}) => {
  const env = process.env.ENV || 'development';
  const isProduction = env === 'production' || env === 'prod';
  
  const entry = {
    timestamp: new Date().toISOString(),
    level,
    service: 'job-application-worker',
    message,
    ...meta,
  };
  
  // Add trace_id if available
  if (traceId) {
    entry.trace_id = traceId;
  }
  
  // In production, always output JSON
  // In development, output JSON for consistency (can be parsed by log aggregators)
  return JSON.stringify(entry);
};

export const logger = {
  info: (message, meta = {}) => {
    console.log(createLogEntry('info', message, meta));
  },

  warn: (message, meta = {}) => {
    console.warn(createLogEntry('warn', message, meta));
  },

  error: (message, meta = {}) => {
    console.error(createLogEntry('error', message, meta));
  },
  
  debug: (message, meta = {}) => {
    const env = process.env.ENV || 'development';
    if (env !== 'production' && env !== 'prod') {
      console.debug(createLogEntry('debug', message, meta));
    }
  },
};

