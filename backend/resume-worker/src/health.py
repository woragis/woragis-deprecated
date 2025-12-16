"""
Health check module for Resume Worker.
Checks service availability and RabbitMQ connection.
"""
import time
from typing import Dict, List, Any
from contextvars import ContextVar

# Cache for health check results
_health_cache: ContextVar[Dict[str, Any]] = ContextVar("health_cache", default={})
_cache_timestamp: ContextVar[float] = ContextVar("cache_timestamp", default=0.0)
_cache_ttl = 5.0  # Cache for 5 seconds

# Global RabbitMQ connection reference (set by worker)
_rabbitmq_connection = None


def set_rabbitmq_connection(conn):
    """Set the RabbitMQ connection for health checks."""
    global _rabbitmq_connection
    _rabbitmq_connection = conn


def check_health() -> Dict[str, Any]:
    """
    Perform health checks for the resume worker.
    
    Returns:
        Dictionary with status and checks
    """
    # Check cache
    cache = _health_cache.get({})
    timestamp = _cache_timestamp.get(0.0)
    
    if cache and (time.time() - timestamp) < _cache_ttl:
        return cache
    
    checks: List[Dict[str, str]] = []
    
    # Check service availability
    checks.append({
        "name": "service",
        "status": "ok"
    })
    
    # Check RabbitMQ connection
    rabbitmq_check = check_rabbitmq()
    checks.append(rabbitmq_check)
    
    # Determine overall status
    has_errors = any(check["status"] == "error" for check in checks)
    
    status = "unhealthy" if has_errors else "healthy"
    
    result = {
        "status": status,
        "checks": checks
    }
    
    # Update cache
    _health_cache.set(result)
    _cache_timestamp.set(time.time())
    
    return result


def check_rabbitmq() -> Dict[str, str]:
    """Check RabbitMQ connection status."""
    global _rabbitmq_connection
    
    if _rabbitmq_connection is None:
        return {
            "name": "rabbitmq",
            "status": "error",
            "message": "not configured"
        }
    
    try:
        # Check if connection is closed
        if hasattr(_rabbitmq_connection, 'is_closed') and _rabbitmq_connection.is_closed:
            return {
                "name": "rabbitmq",
                "status": "error",
                "message": "connection closed"
            }
        
        # Try to check connection (pika specific)
        if hasattr(_rabbitmq_connection, 'is_closing') and _rabbitmq_connection.is_closing:
            return {
                "name": "rabbitmq",
                "status": "error",
                "message": "connection closing"
            }
        
        return {
            "name": "rabbitmq",
            "status": "ok"
        }
    except Exception as e:
        return {
            "name": "rabbitmq",
            "status": "error",
            "message": str(e)
        }
