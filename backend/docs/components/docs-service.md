# Docs Service Component

## Overview

A Python/FastAPI microservice that provides a REST API for serving technical documentation. It converts markdown files to HTML with syntax highlighting and provides endpoints to list and retrieve documentation files.

## Architecture

- **Language**: Python 3.11+
- **Framework**: FastAPI
- **Port**: 8000 (8002 in docker-compose)
- **Dependencies**: None (stateless service)
- **Storage**: Reads from filesystem (docs directory)

## Responsibilities

1. **Serve Documentation**: Provide REST API endpoints to access documentation files
2. **Markdown Rendering**: Convert markdown to HTML with syntax highlighting
3. **File Listing**: List all available documentation files with metadata
4. **Category Filtering**: Filter documentation by category (architecture, adr, runbooks, etc.)

## Health Check

**Endpoint**: `GET /healthz`

**Checks**:
- Service availability
- Docs directory existence
- Docs directory readability
- Markdown file count

**Response**:
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

**Caching**: Results cached for 5 seconds

## Metrics

**Endpoint**: `GET /metrics`

Exposes Prometheus metrics:
- HTTP request rate and latency (automatic via `prometheus-fastapi-instrumentator`)
- Health check metrics

## Configuration

### Environment Variables

#### Required
- None (service has sensible defaults)

#### Optional
- `DOCS_ROOT` - Path to docs directory (default: `/app/docs`)
- `ENV` - Environment (development/production)
- `CORS_ENABLED` - Enable CORS (default: `true`)
- `CORS_ALLOWED_ORIGINS` - Comma-separated origins (default: `*`)
- `MARKDOWN_EXTENSIONS` - Comma-separated markdown extensions (default: `fenced_code,codehilite,tables,toc,extra`)
- `LOG_TO_FILE` - Enable file logging in development (default: `false`)
- `LOG_DIR` - Log directory (default: `logs`)

## API Endpoints

### Root

**GET** `/`

Service information and available endpoints.

**Response**:
```json
{
  "service": "woragis-docs-service",
  "version": "0.1.0",
  "endpoints": {
    "docs": "/api/v1/docs",
    "health": "/healthz",
    "metrics": "/metrics"
  }
}
```

### List Documentation Files

**GET** `/api/v1/docs/`

List all available documentation files.

**Query Parameters**:
- `category` (optional) - Filter by category (e.g., `architecture`, `adr`, `runbooks`)

**Response**:
```json
{
  "files": [
    {
      "path": "architecture/system-overview.md",
      "title": "System Overview",
      "size": 12345,
      "category": "architecture"
    }
  ],
  "total": 1
}
```

### Get Documentation File

**GET** `/api/v1/docs/{path}`

Get documentation content.

**Path Parameters**:
- `path` - Relative path to doc file (e.g., `architecture/system-overview.md`)

**Query Parameters**:
- `format` (optional) - Response format: `json` (default) or `html`

**Response (JSON format)**:
```json
{
  "path": "architecture/system-overview.md",
  "title": "System Overview",
  "content": "# System Overview\n\n...",
  "html": "<h1>System Overview</h1>...",
  "metadata": null
}
```

**Response (HTML format)**:
Returns a full HTML page with embedded CSS.

## Features

### Markdown Processing

- **Extensions Supported**:
  - `fenced_code` - Fenced code blocks
  - `codehilite` - Syntax highlighting (via Pygments)
  - `tables` - Markdown tables
  - `toc` - Table of contents
  - `extra` - Additional markdown features

- **Syntax Highlighting**: Uses Pygments with GitHub style
- **Frontmatter Support**: Extracts YAML frontmatter if present

### File Discovery

The service automatically finds documentation files:
1. Tries exact path
2. Tries with `.md` extension if not provided
3. Tries directory with `README.md` or `index.md` if path is a directory

### Security

- Prevents directory traversal attacks
- Validates file paths
- Only serves markdown files from docs directory

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `docs-service`

**Key Log Fields**:
- `path` - Document path
- `format` - Response format (json/html)
- `category` - Category filter (if used)
- `trace_id` - Request trace ID

## Deployment

### Local Development

```bash
cd backend/docs-service
pip install -r requirements.txt
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

### Docker

```bash
docker build -f Dockerfile.docs-service -t woragis-docs-service .
docker run -p 8002:8000 \
  -e DOCS_ROOT=/app/docs \
  -e ENV=development \
  woragis-docs-service
```

### Docker Compose

The service is included in `docker-compose.yml`:

```yaml
docs-service:
  build:
    context: .
    dockerfile: Dockerfile.docs-service
  ports:
    - "8002:8000"
  environment:
    DOCS_ROOT: /app/docs
    ENV: ${APP_ENV:-development}
```

## Scaling

### Horizontal Scaling
- Stateless design allows multiple replicas
- Load balancer distributes requests
- Each replica reads from same docs directory (shared volume)

### Resource Requirements
- **CPU**: 100m-300m (0.1-0.3 core)
- **Memory**: 256Mi-512Mi
- **Storage**: Read-only access to docs directory

## Performance Considerations

1. **File Reading**: Files are read on-demand (not cached in memory)
2. **Markdown Parsing**: Parsing happens on each request (consider caching for high traffic)
3. **Health Check Caching**: Health checks cached for 5 seconds
4. **File Discovery**: File listing scans filesystem (consider caching for large doc sets)

## Monitoring

### Key Metrics
- Request rate (requests/second)
- Latency (p50, p95, p99)
- Error rate
- File read errors
- Health check status

### Alerts
- Error rate > 5%
- Latency p95 > 1 second
- Health check failures
- Docs directory not accessible

## Troubleshooting

### Common Issues

#### Documentation Not Found
- Check `DOCS_ROOT` environment variable
- Verify docs directory exists in container
- Check file permissions
- Verify file path is correct

#### High Latency
- Consider caching parsed markdown
- Optimize file discovery (cache file list)
- Check filesystem performance

#### Health Check Failures
- Check docs directory exists
- Verify directory permissions
- Check for filesystem errors

## Related Documentation

- [API Documentation](../api/docs-service-api.md) - Detailed API documentation
- [Architecture Decision Records](../adr/) - Architectural decisions
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
