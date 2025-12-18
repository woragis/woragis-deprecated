# Woragis Docs Service

API service for serving technical documentation. Provides REST endpoints to access markdown documentation files as JSON or rendered HTML.

## Features

- **REST API**: Serve documentation as JSON or HTML
- **Markdown Parsing**: Converts markdown to HTML with syntax highlighting
- **File Listing**: List all available documentation files
- **Category Filtering**: Filter docs by category (architecture, adr, runbooks, etc.)
- **Health Checks**: Built-in health check endpoint with caching
- **Structured Logging**: JSON logging in production, human-readable in development
- **Prometheus Metrics**: Exposes metrics at `/metrics`
- **Request Tracing**: Automatic trace ID generation and propagation

## API Endpoints

### List Documentation Files
```
GET /api/v1/docs/
GET /api/v1/docs/?category=architecture
```

Returns a list of all available documentation files with metadata.

### Get Documentation File
```
GET /api/v1/docs/{path}
GET /api/v1/docs/{path}?format=json
GET /api/v1/docs/{path}?format=html
```

Returns documentation content:
- `format=json` (default): Returns structured JSON with `content` (raw markdown) and `html` (rendered)
- `format=html`: Returns a full HTML page

**Examples:**
- `/api/v1/docs/architecture/system-overview.md`
- `/api/v1/docs/adr/001-rabbitmq-redis-fallback.md`
- `/api/v1/docs/README.md`

### Health Check
```
GET /healthz
```

Returns service health status with checks for:
- Service availability
- Docs directory existence
- Docs directory readability
- Markdown file count

Results are cached for 5 seconds.

### Metrics
```
GET /metrics
```

Prometheus metrics endpoint (exposed via `prometheus-fastapi-instrumentator`).

## Configuration

Environment variables (see `env.sample`):

- `DOCS_ROOT`: Path to docs directory (default: `/app/docs`)
- `ENV`: Environment (development/production)
- `LOG_TO_FILE`: Enable file logging in development (default: `false`)
- `LOG_DIR`: Log directory (default: `logs`)
- `CORS_ENABLED`: Enable CORS (default: `true`)
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed origins (default: `*`)
- `MARKDOWN_EXTENSIONS`: Comma-separated markdown extensions (default: `fenced_code,codehilite,tables,toc,extra`)

## Docker Build

The Dockerfile copies the `docs` directory directly into the image (no symlinks). This ensures docs are available in the built image.

**Important**: The build context must include the `docs` directory. The Dockerfile expects:
- `docs/` directory at the build context root
- `docs-service/app/` directory with the service code
- `docs-service/requirements.txt` at the build context root

## Development

### Local Development

```bash
# Install dependencies
pip install -r requirements.txt

# Set environment variables
cp env.sample .env

# Run the service
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

### Docker

```bash
# Build image
docker build -f Dockerfile.docs-service -t woragis-docs-service .

# Run container
docker run -p 8002:8000 \
  -e DOCS_ROOT=/app/docs \
  -e ENV=development \
  woragis-docs-service
```

### Docker Compose

The service is included in `docker-compose.yml`:

```bash
docker-compose up docs-service
```

## Testing

```bash
# Run all tests
pytest

# Run with coverage
pytest --cov=app --cov-report=html

# Run specific test file
pytest tests/unit/test_api.py

# Run integration tests
pytest tests/integration/
```

## Usage Examples

### Get all docs
```bash
curl http://localhost:8002/api/v1/docs/
```

### Get specific doc as JSON
```bash
curl http://localhost:8002/api/v1/docs/architecture/system-overview.md
```

### Get specific doc as HTML
```bash
curl http://localhost:8002/api/v1/docs/architecture/system-overview.md?format=html
```

### Filter by category
```bash
curl http://localhost:8002/api/v1/docs/?category=adr
```

### Health check
```bash
curl http://localhost:8002/healthz
```

### Metrics
```bash
curl http://localhost:8002/metrics
```

## Response Format

### JSON Response
```json
{
  "path": "architecture/system-overview.md",
  "title": "System Overview",
  "content": "# System Overview\n\n...",
  "html": "<h1>System Overview</h1>...",
  "metadata": null
}
```

### HTML Response
Returns a full HTML page with embedded CSS for styling and code highlighting.

### Health Check Response
```json
{
  "status": "healthy",
  "service": "docs-service",
  "checks": [
    {"name": "service", "status": "ok"},
    {"name": "docs_directory", "status": "ok"},
    {"name": "docs_readable", "status": "ok"},
    {"name": "markdown_files", "status": "ok"}
  ]
}
```

## Architecture

- **FastAPI**: Web framework
- **Markdown**: Markdown parsing library
- **Pygments**: Syntax highlighting for code blocks
- **Structlog**: Structured logging
- **Prometheus**: Metrics collection via `prometheus-fastapi-instrumentator`

## Observability

### Logging
- Structured JSON logs in production
- Human-readable logs in development
- Automatic trace ID generation
- Request/response logging middleware

### Metrics
- HTTP request metrics (automatic via instrumentator)
- Health check metrics
- Custom metrics can be added

### Health Checks
- Service availability check
- Docs directory checks
- Results cached for 5 seconds
- Returns 503 if unhealthy

## Notes

- Docs are copied directly into the Docker image (no symlinks)
- The service reads docs from the filesystem at runtime
- Markdown files are parsed on-demand (not cached)
- Supports standard markdown extensions (tables, code blocks, TOC, etc.)
- Health checks are cached to reduce filesystem load
