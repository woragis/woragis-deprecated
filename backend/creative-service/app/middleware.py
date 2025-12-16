"""
FastAPI middleware for structured logging and trace ID handling.
"""
import uuid
from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware

from app.logger import get_logger, set_trace_id, get_trace_id

logger = get_logger()


class RequestIDMiddleware(BaseHTTPMiddleware):
    """Middleware to generate and propagate trace IDs (request IDs)."""
    
    async def dispatch(self, request: Request, call_next):
        # Check if trace_id already exists in header (for distributed tracing)
        trace_id = request.headers.get("X-Trace-ID")
        if not trace_id:
            # Generate new trace_id if not present
            trace_id = str(uuid.uuid4())
        
        # Set trace_id in context
        set_trace_id(trace_id)
        
        # Add trace_id to response header
        response = await call_next(request)
        response.headers["X-Trace-ID"] = trace_id
        
        return response


class RequestLoggerMiddleware(BaseHTTPMiddleware):
    """Middleware to log HTTP requests with structured fields."""
    
    async def dispatch(self, request: Request, call_next):
        import time
        start_time = time.time()
        
        # Get request body size (if available)
        body_size = 0
        if hasattr(request, "_body"):
            body_size = len(request._body) if request._body else 0
        
        # Process request
        response = await call_next(request)
        
        # Calculate duration
        duration = time.time() - start_time
        
        # Get trace_id
        trace_id = get_trace_id()
        
        # Build log attributes
        log_data = {
            "method": request.method,
            "path": request.url.path,
            "ip": request.client.host if request.client else None,
            "status": response.status_code,
            "duration_ms": round(duration * 1000, 2),
        }
        
        if body_size > 0:
            log_data["bytes_in"] = body_size
        
        if trace_id:
            log_data["trace_id"] = trace_id
        
        # Add query parameters if present
        if request.query_params:
            log_data["query"] = str(request.query_params)
        
        # Log based on status code
        status = response.status_code
        if status >= 500:
            logger.error("http request", **log_data)
        elif status >= 400:
            logger.warn("http request", **log_data)
        else:
            logger.info("http request", **log_data)
        
        return response
