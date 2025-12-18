"""
Health check utilities for Docs Service.
Checks service availability and docs directory status.
"""
import os
import time
from pathlib import Path
from typing import Dict, List, Any
from contextvars import ContextVar

from app.config import settings

# Cache for health check results (5 second TTL)
_health_cache: ContextVar[Dict[str, Any]] = ContextVar("health_cache", default={})
_cache_timestamp: ContextVar[float] = ContextVar("cache_timestamp", default=0.0)
_cache_ttl = 5.0  # Cache for 5 seconds


def check_health() -> Dict[str, Any]:
    """
    Perform health checks for the docs service.
    Results are cached for 5 seconds to reduce filesystem load.
    
    Returns:
        Dictionary with status and checks
    """
    # Check cache
    cache = _health_cache.get({})
    timestamp = _cache_timestamp.get(0.0)
    
    if cache and (time.time() - timestamp) < _cache_ttl:
        return cache
    
    checks: List[Dict[str, str]] = []
    status = "healthy"
    
    # Check service availability (always ok if endpoint is reachable)
    checks.append({
        "name": "service",
        "status": "ok"
    })
    
    # Check if docs directory exists
    docs_path = Path(settings.DOCS_ROOT)
    docs_exists = docs_path.exists()
    docs_readable = docs_exists and os.access(docs_path, os.R_OK)
    
    checks.append({
        "name": "docs_directory",
        "status": "ok" if docs_exists else "error"
    })
    
    if not docs_exists:
        status = "unhealthy"
    
    # Check if docs directory is readable
    if docs_exists:
        checks.append({
            "name": "docs_readable",
            "status": "ok" if docs_readable else "error"
        })
        
        if not docs_readable:
            status = "unhealthy"
        
        # Count markdown files (for informational purposes)
        try:
            md_files = list(docs_path.rglob("*.md")) + list(docs_path.rglob("*.markdown"))
            checks.append({
                "name": "markdown_files",
                "status": "ok"
            })
        except Exception:
            # If we can't count files, it's not critical but log it
            checks.append({
                "name": "markdown_files",
                "status": "error"
            })
    
    result = {
        "status": status,
        "service": "docs-service",
        "checks": checks
    }
    
    # Update cache
    _health_cache.set(result)
    _cache_timestamp.set(time.time())
    
    return result
