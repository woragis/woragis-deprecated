// Reusable middleware functions

// Request logging middleware
export const requestLogger = (req: any, res: any, next: any) => {
  const start = Date.now();
  const { method, url, headers } = req;
  
  console.log(`[${new Date().toISOString()}] ${method} ${url}`);
  console.log('Headers:', headers);
  
  res.on('finish', () => {
    const duration = Date.now() - start;
    console.log(`[${new Date().toISOString()}] ${method} ${url} - ${res.statusCode} (${duration}ms)`);
  });
  
  next();
};

// Error handling middleware
export const errorHandler = (err: any, req: any, res: any, next: any) => {
  console.error('Error:', err);
  
  const status = err.status || err.statusCode || 500;
  const message = err.message || 'Internal Server Error';
  
  res.status(status).json({
    success: false,
    error: message,
    ...(process.env.NODE_ENV === 'development' && { stack: err.stack }),
  });
};

// Enhanced CORS middleware for web and mobile apps
export const corsHandler = (req: any, res: any, next: any) => {
  const origin = req.headers.origin;
  const allowedOrigins = process.env.ALLOWED_ORIGINS?.split(',') || ['http://localhost:3000'];
  const mobileOrigins = process.env.MOBILE_APP_ORIGINS?.split(',') || ['capacitor://localhost', 'ionic://localhost'];
  
  // Combine all allowed origins
  const allAllowedOrigins = [...allowedOrigins, ...mobileOrigins];
  
  // Check if origin is allowed
  if (origin && allAllowedOrigins.includes(origin)) {
    res.setHeader('Access-Control-Allow-Origin', origin);
  } else if (origin && (origin.startsWith('capacitor://') || origin.startsWith('ionic://'))) {
    // Allow all Capacitor/Ionic origins for mobile apps
    res.setHeader('Access-Control-Allow-Origin', origin);
  }
  
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS, PATCH');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-API-Key, X-Requested-With');
  res.setHeader('Access-Control-Allow-Credentials', 'true');
  res.setHeader('Access-Control-Max-Age', '86400'); // 24 hours
  
  if (req.method === 'OPTIONS') {
    res.status(200).end();
    return;
  }
  
  next();
};

// Rate limiting middleware
export const rateLimiter = (maxRequests: number = 100, windowMs: number = 15 * 60 * 1000) => {
  const requests = new Map();
  
  return (req: any, res: any, next: any) => {
    const ip = req.ip || req.connection.remoteAddress;
    const now = Date.now();
    const windowStart = now - windowMs;
    
    // Clean old entries
    for (const [key, value] of requests.entries()) {
      if (value.timestamp < windowStart) {
        requests.delete(key);
      }
    }
    
    const userRequests = requests.get(ip) || { count: 0, timestamp: now };
    
    if (userRequests.timestamp < windowStart) {
      userRequests.count = 1;
      userRequests.timestamp = now;
    } else {
      userRequests.count++;
    }
    
    requests.set(ip, userRequests);
    
    if (userRequests.count > maxRequests) {
      res.status(429).json({
        success: false,
        error: 'Too many requests',
        retryAfter: Math.ceil((userRequests.timestamp + windowMs - now) / 1000),
      });
      return;
    }
    
    next();
  };
};

// Validation middleware factory
export const validateRequest = (schema: any) => {
  return (req: any, res: any, next: any) => {
    try {
      const validatedData = schema.parse(req.body);
      req.body = validatedData;
      next();
    } catch (error: any) {
      res.status(400).json({
        success: false,
        error: 'Validation failed',
        details: error.errors,
      });
    }
  };
};

// Async error wrapper
export const asyncHandler = (fn: (req: any, res: any, next: any) => any) => {
  return (req: any, res: any, next: any) => {
    Promise.resolve(fn(req, res, next)).catch(next);
  };
};

// Security headers middleware
export const securityHeaders = (req: any, res: any, next: any) => {
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'DENY');
  res.setHeader('X-XSS-Protection', '1; mode=block');
  res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
  
  if (process.env.NODE_ENV === 'production') {
    res.setHeader('Strict-Transport-Security', 'max-age=31536000; includeSubDomains');
  }
  
  next();
};
