/**
 * Structured logging utility for Job Application Worker.
 * Provides JSON logging with service name and trace_id support.
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Context for trace_id (simple in-memory for now, can be enhanced with AsyncLocalStorage)
let traceId = null;

// File logging setup
let logFileStream = null;
const env = process.env.ENV || 'development';
const isProduction = env === 'production' || env === 'prod';
const logToFile = process.env.LOG_TO_FILE === 'true';
const logDir = process.env.LOG_DIR || 'logs';

// Initialize file logging if enabled (development only)
if (!isProduction && logToFile) {
  try {
    // Create log directory if it doesn't exist
    const logDirPath = path.resolve(process.cwd(), logDir);
    if (!fs.existsSync(logDirPath)) {
      fs.mkdirSync(logDirPath, { recursive: true });
    }
    
    const logFile = path.join(logDirPath, 'job-application-worker.log');
    logFileStream = fs.createWriteStream(logFile, { flags: 'a' });
  } catch (err) {
    console.error('Failed to create log file, falling back to stdout only', err);
  }
}

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
 * Write log entry to both stdout and file (if enabled)
 */
const writeLog = (level, entry) => {
  const logLine = JSON.stringify(entry) + '\n';
  
  // Always write to stdout
  if (level === 'error') {
    console.error(logLine);
  } else if (level === 'warn') {
    console.warn(logLine);
  } else if (level === 'debug') {
    console.debug(logLine);
  } else {
    console.log(logLine);
  }
  
  // Also write to file if enabled
  if (logFileStream) {
    logFileStream.write(logLine);
  }
};

/**
 * Create a log entry with standard fields
 */
const createLogEntry = (level, message, meta = {}) => {
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
  
  return entry;
};

export const logger = {
  info: (message, meta = {}) => {
    const entry = createLogEntry('info', message, meta);
    writeLog('info', entry);
  },

  warn: (message, meta = {}) => {
    const entry = createLogEntry('warn', message, meta);
    writeLog('warn', entry);
  },

  error: (message, meta = {}) => {
    const entry = createLogEntry('error', message, meta);
    writeLog('error', entry);
  },
  
  debug: (message, meta = {}) => {
    if (!isProduction) {
      const entry = createLogEntry('debug', message, meta);
      writeLog('debug', entry);
    }
  },
  
  /**
   * Close log file stream (call on shutdown)
   */
  close: () => {
    if (logFileStream) {
      logFileStream.end();
      logFileStream = null;
    }
  },
};

