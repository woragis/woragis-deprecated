# Adding a New Service

## Overview

This guide explains how to create a new microservice for the Woragis platform. Services are typically Python/FastAPI applications that provide specialized functionality.

## Service Structure

Services follow this structure:

```
{service-name}/
├── app/
│   ├── main.py              # FastAPI application
│   ├── config.py            # Configuration
│   ├── logger.py            # Logging setup
│   ├── health.py            # Health checks
│   └── providers/           # External service providers
├── tests/
│   ├── unit/
│   └── integration/
├── Dockerfile
├── requirements.txt
├── env.sample
└── README.md
```

## Step-by-Step Guide

### 1. Create Service Directory

```bash
mkdir -p backend/{service-name}/app
cd backend/{service-name}
```

### 2. Initialize Python Project

```bash
python -m venv venv
source venv/bin/activate  # or `venv\Scripts\activate` on Windows
pip install fastapi uvicorn[standard] structlog prometheus-fastapi-instrumentator
```

### 3. Create requirements.txt

```txt
fastapi==0.104.1
uvicorn[standard]==0.24.0
structlog==23.2.0
prometheus-fastapi-instrumentator==7.1.0
pydantic==2.5.0
python-dotenv==1.0.0
```

### 4. Create Configuration

**app/config.py**:
```python
from pydantic_settings import BaseSettings
from typing import Literal

class Settings(BaseSettings):
    ENV: Literal["development", "production"] = "development"
    CORS_ENABLED: bool = True
    CORS_ALLOWED_ORIGINS: str = "*"
    
    # Service-specific settings
    API_KEY: str = ""
    
    class Config:
        env_file = ".env"

settings = Settings()
```

### 5. Create Logger

**app/logger.py**:
```python
import structlog
import sys
import os

def configure_logging(env: str = "development", log_to_file: bool = False, log_dir: str = "logs"):
    processors = [
        structlog.stdlib.filter_by_level,
        structlog.stdlib.add_logger_name,
        structlog.stdlib.add_log_level,
        structlog.stdlib.PositionalArgumentsFormatter(),
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.StackInfoRenderer(),
        structlog.processors.format_exc_info,
    ]
    
    if env == "production":
        processors.append(structlog.processors.JSONRenderer())
    else:
        processors.append(structlog.dev.ConsoleRenderer())
    
    structlog.configure(
        processors=processors,
        wrapper_class=structlog.stdlib.BoundLogger,
        context_class=dict,
        logger_factory=structlog.stdlib.LoggerFactory(),
        cache_logger_on_first_use=True,
    )

def get_logger():
    return structlog.get_logger().bind(service="service-name")
```

### 6. Create Health Check

**app/health.py**:
```python
def check_health():
    checks = []
    status = "healthy"
    
    # Check dependencies if needed
    # if not check_dependency():
    #     status = "unhealthy"
    #     checks.append({"name": "dependency", "status": "error"})
    
    return {
        "status": status,
        "checks": checks
    }
```

### 7. Create Main Application

**app/main.py**:
```python
import os
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from dotenv import load_dotenv
from prometheus_fastapi_instrumentator import Instrumentator

from app.config import settings
from app.logger import configure_logging, get_logger
from app.health import check_health

load_dotenv()

# Configure logging
env = os.getenv("ENV", "development")
log_to_file = os.getenv("LOG_TO_FILE", "false").lower() == "true"
log_dir = os.getenv("LOG_DIR", "logs")
configure_logging(env=env, log_to_file=log_to_file, log_dir=log_dir)

logger = get_logger()
logger.info("Service initialized", env=env)

app = FastAPI(title="Woragis {Service Name}", version="0.1.0")

# Add Prometheus metrics
Instrumentator().instrument(app).expose(app)

# CORS
if settings.CORS_ENABLED:
    origins = settings.CORS_ALLOWED_ORIGINS.split(",")
    app.add_middleware(
        CORSMiddleware,
        allow_origins=[o.strip() for o in origins if o.strip()],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

# Request/Response Models
class RequestModel(BaseModel):
    field: str

class ResponseModel(BaseModel):
    result: str

# Endpoints
@app.get("/healthz")
def healthz():
    result = check_health()
    status_code = 200 if result["status"] == "healthy" else 503
    return JSONResponse(content=result, status_code=status_code)

@app.post("/v1/endpoint", response_model=ResponseModel)
async def endpoint(request: RequestModel):
    logger.info("Request received", field=request.field)
    
    try:
        # Process request
        result = process_request(request.field)
        
        return ResponseModel(result=result)
    except Exception as e:
        logger.exception("Request failed", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
```

### 8. Create Dockerfile

```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY . .

EXPOSE 8000

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

### 9. Create env.sample

```env
ENV=development
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
API_KEY=your-api-key-here
```

### 10. Add to docker-compose.yml

```yaml
{service-name}:
  build:
    context: ./{service-name}
    dockerfile: Dockerfile
  environment:
    ENV: ${APP_ENV:-development}
    CORS_ENABLED: ${CORS_ENABLED:-true}
    CORS_ALLOWED_ORIGINS: ${CORS_ALLOWED_ORIGINS:-*}
    API_KEY: ${SERVICE_API_KEY}
  ports:
    - "8000:8000"
  restart: on-failure
```

### 11. Write Tests

**tests/unit/test_main.py**:
```python
import pytest
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_healthz():
    response = client.get("/healthz")
    assert response.status_code == 200
    assert response.json()["status"] == "healthy"

def test_endpoint():
    response = client.post("/v1/endpoint", json={"field": "value"})
    assert response.status_code == 200
    assert "result" in response.json()
```

## Configuration

### Environment Variables

Common variables:
- `ENV` - Environment (development/production)
- `CORS_ENABLED` - Enable CORS
- `CORS_ALLOWED_ORIGINS` - Comma-separated origins

Service-specific:
- Add service-specific environment variables in `config.py`

## Best Practices

1. **API Design**:
   - Use RESTful conventions
   - Version APIs (`/v1/...`)
   - Use consistent response formats

2. **Error Handling**:
   - Return appropriate HTTP status codes
   - Provide clear error messages
   - Log errors with context

3. **Validation**:
   - Use Pydantic models for request/response validation
   - Validate input early
   - Return clear validation errors

4. **Logging**:
   - Use structured logging
   - Include request ID in logs
   - Log important operations

5. **Metrics**:
   - Expose Prometheus metrics
   - Record request rate and latency
   - Record error rates

6. **Testing**:
   - Unit tests for business logic
   - Integration tests for endpoints
   - Mock external dependencies

## Example Services

See existing services for examples:
- `ai-service/` - AI/LLM integration service
- `creative-service/` - Image/diagram/video generation service

## Related Documentation

- [Component Documentation](../components/) - Service component details
- [API Documentation](../api/) - API documentation examples
- [Testing Patterns](./testing-patterns.md) - Testing guidelines
- [Logging Conventions](./logging-conventions.md) - Logging guidelines
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
