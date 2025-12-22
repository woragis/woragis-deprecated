import os
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from dotenv import load_dotenv
from prometheus_fastapi_instrumentator import Instrumentator

from app.config import settings
from app.logger import configure_logging, get_logger
from app.middleware import RequestIDMiddleware, RequestLoggerMiddleware
from app.tracing import init_tracing
from app.health import check_health
from app.routes import docs

load_dotenv()

# Configure structured logging
env = os.getenv("ENV", "development")
log_to_file = os.getenv("LOG_TO_FILE", "false").lower() == "true"
log_dir = os.getenv("LOG_DIR", "logs")
configure_logging(env=env, log_to_file=log_to_file, log_dir=log_dir)

logger = get_logger()
logger.info("Docs service initialized", env=env, docs_root=settings.DOCS_ROOT)

# Initialize OpenTelemetry tracing
try:
    init_tracing(
        service_name="docs-service",
        service_version="0.1.0",
        environment=env,
    )
    logger.info("Tracing initialized")
except Exception as e:
    logger.warn("Failed to initialize tracing", error=str(e))

app = FastAPI(
    title="Woragis Docs Service",
    description="API service for serving technical documentation",
    version="0.1.0",
)

# Add middleware for request ID and logging
app.add_middleware(RequestIDMiddleware)
app.add_middleware(RequestLoggerMiddleware)

# Add Prometheus metrics instrumentation
Instrumentator().instrument(app).expose(app)

# Configure CORS
if settings.CORS_ENABLED:
    origins = settings.CORS_ALLOWED_ORIGINS.split(",") if settings.CORS_ALLOWED_ORIGINS != "*" else ["*"]
    app.add_middleware(
        CORSMiddleware,
        allow_origins=[o.strip() for o in origins if o.strip()] if origins != ["*"] else ["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

# Include routers
app.include_router(docs.router)


@app.get("/")
async def root():
    """Root endpoint with service information."""
    return {
        "service": "woragis-docs-service",
        "version": "0.1.0",
        "endpoints": {
            "docs": "/api/v1/docs",
            "health": "/healthz",
            "metrics": "/metrics",
        },
    }


@app.get("/healthz")
def healthz():
    """
    Health check endpoint.
    Returns service availability and dependency status.
    """
    result = check_health()
    
    # Determine HTTP status code
    status_code = 200
    if result["status"] == "unhealthy":
        status_code = 503
    
    return JSONResponse(content=result, status_code=status_code)
