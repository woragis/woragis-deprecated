# Docs Service API Documentation

## Overview

The Docs Service provides a REST API for accessing technical documentation. It serves markdown files as JSON or rendered HTML with syntax highlighting.

## Base URL

- **Development**: `http://localhost:8002`
- **Production**: `http://docs-service:8000` (internal) or `https://docs.woragis.com` (if exposed)

## Authentication

Currently, the Docs Service does not require authentication (internal service). In production, consider adding API key authentication if exposing publicly.

## API Endpoints

### Root

**GET** `/`

Get service information and available endpoints.

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

**Status Codes**:
- `200 OK` - Success

**Example**:
```bash
curl http://localhost:8002/
```

---

### List Documentation Files

**GET** `/api/v1/docs/`

List all available documentation files with metadata.

**Query Parameters**:
- `category` (optional) - Filter by category (e.g., `architecture`, `adr`, `runbooks`, `components`, `api`, `development`)

**Response**:
```json
{
  "files": [
    {
      "path": "architecture/system-overview.md",
      "title": "System Overview",
      "size": 12345,
      "category": "architecture"
    },
    {
      "path": "adr/001-rabbitmq-redis-fallback.md",
      "title": "ADR-001: RabbitMQ with Redis Fallback",
      "size": 5678,
      "category": "adr"
    }
  ],
  "total": 2
}
```

**Response Fields**:
- `files`: Array of documentation file objects
  - `path`: Relative path to the file
  - `title`: Title extracted from markdown (first H1)
  - `size`: File size in bytes
  - `category`: Category (subdirectory name)
- `total`: Total number of files (after filtering)

**Status Codes**:
- `200 OK` - Success
- `500 Internal Server Error` - Docs directory not found

**Examples**:
```bash
# List all docs
curl http://localhost:8002/api/v1/docs/

# Filter by category
curl http://localhost:8002/api/v1/docs/?category=architecture

# Filter ADRs
curl http://localhost:8002/api/v1/docs/?category=adr
```

---

### Get Documentation File

**GET** `/api/v1/docs/{path}`

Get documentation content in JSON or HTML format.

**Path Parameters**:
- `path` - Relative path to the documentation file
  - Examples: `architecture/system-overview.md`, `adr/001-rabbitmq-redis-fallback.md`, `README.md`
  - `.md` extension is optional (will be added automatically)
  - Directory paths will look for `README.md` or `index.md`

**Query Parameters**:
- `format` (optional) - Response format: `json` (default) or `html`

**Response (JSON format)**:
```json
{
  "path": "architecture/system-overview.md",
  "title": "System Overview",
  "content": "# System Overview\n\n## Architecture Diagram\n\n...",
  "html": "<h1>System Overview</h1><h2>Architecture Diagram</h2>...",
  "metadata": null
}
```

**Response (HTML format)**:
Returns a full HTML page with:
- Embedded CSS for styling
- Syntax highlighting CSS (Pygments)
- Responsive design
- Proper typography

**Response Fields**:
- `path`: Relative path to the file
- `title`: Title extracted from markdown (first H1)
- `content`: Raw markdown content
- `html`: Rendered HTML content
- `metadata`: Frontmatter metadata (if present in markdown)

**Status Codes**:
- `200 OK` - Success
- `404 Not Found` - Documentation file not found
- `500 Internal Server Error` - Failed to read file

**Examples**:
```bash
# Get doc as JSON (default)
curl http://localhost:8002/api/v1/docs/architecture/system-overview.md

# Get doc as HTML
curl http://localhost:8002/api/v1/docs/architecture/system-overview.md?format=html

# Get ADR
curl http://localhost:8002/api/v1/docs/adr/001-rabbitmq-redis-fallback.md

# Get README (no extension needed)
curl http://localhost:8002/api/v1/docs/README

# Get directory index
curl http://localhost:8002/api/v1/docs/architecture
```

---

### Health Check

**GET** `/healthz`

Check service health and documentation availability.

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

**Response Fields**:
- `status`: Overall status (`healthy` or `unhealthy`)
- `service`: Service name
- `checks`: Array of health checks
  - `service`: Service availability
  - `docs_directory`: Docs directory exists
  - `docs_readable`: Docs directory is readable
  - `markdown_files`: Markdown files found

**Status Codes**:
- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service is unhealthy

**Caching**: Results are cached for 5 seconds

**Example**:
```bash
curl http://localhost:8002/healthz
```

---

### Metrics

**GET** `/metrics`

Prometheus metrics endpoint (internal use only).

**Content-Type**: `text/plain`

**Example**:
```bash
curl http://localhost:8002/metrics
```

---

## Path Resolution

The service uses intelligent path resolution:

1. **Exact Match**: Tries the exact path provided
2. **Add Extension**: If no extension, tries adding `.md`
3. **Directory Index**: If path is a directory, looks for `README.md` or `index.md`

**Examples**:
- `architecture/system-overview` → `architecture/system-overview.md`
- `architecture` → `architecture/README.md` (if exists)
- `README` → `README.md`

## Markdown Features

### Supported Extensions

- **fenced_code**: Fenced code blocks with language specification
- **codehilite**: Syntax highlighting via Pygments
- **tables**: Markdown tables
- **toc**: Table of contents generation
- **extra**: Additional markdown features (abbreviations, footnotes, etc.)

### Syntax Highlighting

Code blocks are automatically highlighted using Pygments with GitHub style:

````markdown
```go
func main() {
    fmt.Println("Hello, World!")
}
```
````

### Frontmatter Support

If a markdown file starts with YAML frontmatter, it will be extracted:

```markdown
---
title: "Document Title"
author: "Author Name"
---

# Content
```

The frontmatter will be available in the `metadata` field of the JSON response.

## Error Responses

### Error Format

```json
{
  "detail": "Error message"
}
```

### Common Errors

**404 Not Found**:
```json
{
  "detail": "Documentation file not found: {path}"
}
```

**500 Internal Server Error**:
```json
{
  "detail": "Failed to read documentation file: {error}"
}
```

or

```json
{
  "detail": "Docs directory not found"
}
```

## Rate Limiting

Currently, no rate limiting is implemented. Consider adding rate limiting if exposing publicly.

## Examples

### Complete Flow: Browse Documentation

```bash
# 1. List all documentation
curl http://localhost:8002/api/v1/docs/

# 2. Filter by category
curl http://localhost:8002/api/v1/docs/?category=adr

# 3. Get specific document as JSON
curl http://localhost:8002/api/v1/docs/adr/001-rabbitmq-redis-fallback.md

# 4. Get same document as HTML (for viewing in browser)
curl http://localhost:8002/api/v1/docs/adr/001-rabbitmq-redis-fallback.md?format=html

# 5. Check service health
curl http://localhost:8002/healthz
```

### Using in Frontend

```javascript
// Fetch documentation list
const response = await fetch('http://localhost:8002/api/v1/docs/?category=architecture');
const data = await response.json();
console.log(data.files); // Array of documentation files

// Fetch specific document
const docResponse = await fetch('http://localhost:8002/api/v1/docs/architecture/system-overview.md?format=html');
const html = await docResponse.text();
document.getElementById('content').innerHTML = html;
```

### Using in CLI Scripts

```bash
#!/bin/bash
# Fetch and display documentation

DOC_PATH="architecture/system-overview.md"
BASE_URL="http://localhost:8002"

# Get as HTML and open in browser
curl -s "${BASE_URL}/api/v1/docs/${DOC_PATH}?format=html" > /tmp/doc.html
open /tmp/doc.html  # macOS
# xdg-open /tmp/doc.html  # Linux
```

## Related Documentation

- [Component Documentation](../components/docs-service.md) - Docs Service component details
- [Architecture Decision Records](../adr/) - Architectural decisions
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
