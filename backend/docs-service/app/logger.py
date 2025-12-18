"""
Structured logging utility for Docs Service.
Provides JSON logging in production, text logging in development.
"""
import os
import sys
import structlog
from typing import Optional
from contextvars import ContextVar

# Context variable for trace_id (thread-safe)
trace_id_var: ContextVar[Optional[str]] = ContextVar("trace_id", default=None)


def configure_logging(env: str = "development", log_to_file: bool = False, log_dir: str = "logs"):
    """
    Configure structured logging for the service.
    
    Args:
        env: Environment name (development, production, etc.)
        log_to_file: Whether to write logs to files (development only)
        log_dir: Directory for log files (if log_to_file is True)
    """
    is_production = env.lower() in ("production", "prod")
    
    # Configure processors
    processors = [
        structlog.contextvars.merge_contextvars,  # Merge context variables
        structlog.stdlib.add_log_level,            # Add log level
        structlog.stdlib.add_logger_name,         # Add logger name
        structlog.processors.TimeStamper(fmt="iso"),  # ISO 8601 timestamp
        structlog.processors.StackInfoRenderer(),     # Stack traces
        structlog.processors.format_exc_info,          # Exception formatting
    ]
    
    if is_production:
        # JSON output for production
        processors.append(structlog.processors.JSONRenderer())
        log_level = "INFO"
    else:
        # Human-readable output for development
        processors.append(structlog.dev.ConsoleRenderer())
        log_level = "DEBUG"
    
    # Configure structlog
    structlog.configure(
        processors=processors,
        wrapper_class=structlog.stdlib.BoundLogger,
        context_class=dict,
        logger_factory=structlog.stdlib.LoggerFactory(),
        cache_logger_on_first_use=True,
    )
    
    # Configure standard library logging
    import logging
    logging.basicConfig(
        format="%(message)s",
        stream=sys.stdout,
        level=getattr(logging, log_level.upper()),
    )
    
    # If file logging is enabled in development, set up file handler
    if not is_production and log_to_file:
        os.makedirs(log_dir, exist_ok=True)
        log_file = os.path.join(log_dir, "docs-service.log")
        file_handler = logging.FileHandler(log_file)
        file_handler.setFormatter(logging.Formatter("%(message)s"))
        
        root_logger = logging.getLogger()
        root_logger.addHandler(file_handler)


def get_logger(name: str = "woragis.docs-service"):
    """
    Get a structured logger instance.
    
    Args:
        name: Logger name
        
    Returns:
        Bound logger with service name and trace_id support
    """
    logger = structlog.get_logger(name)
    
    # Add service name to all logs
    logger = logger.bind(service="docs-service")
    
    # Add trace_id from context if available
    trace_id = trace_id_var.get()
    if trace_id:
        logger = logger.bind(trace_id=trace_id)
    
    return logger


def set_trace_id(trace_id: str):
    """Set trace_id in context for distributed tracing."""
    trace_id_var.set(trace_id)


def get_trace_id() -> Optional[str]:
    """Get trace_id from context."""
    return trace_id_var.get()
