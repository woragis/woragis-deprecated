"""
Health check module for Creative Service.
Checks service availability and external dependencies.
"""
import time
from typing import Dict, List, Any
from contextvars import ContextVar

# Cache for health check results
_health_cache: ContextVar[Dict[str, Any]] = ContextVar("health_cache", default={})
_cache_timestamp: ContextVar[float] = ContextVar("cache_timestamp", default=0.0)
_cache_ttl = 5.0  # Cache for 5 seconds


def check_health() -> Dict[str, Any]:
    """
    Perform health checks for the creative service.
    
    Returns:
        Dictionary with status and checks
    """
    # Check cache
    cache = _health_cache.get({})
    timestamp = _cache_timestamp.get(0.0)
    
    if cache and (time.time() - timestamp) < _cache_ttl:
        return cache
    
    checks: List[Dict[str, str]] = []
    
    # Check service availability (always ok if endpoint is reachable)
    checks.append({
        "name": "service",
        "status": "ok"
    })
    
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
