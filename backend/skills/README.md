# Skills Service

Skills and interests management service for Woragis platform.

## Overview

The Skills Service is a standalone microservice responsible for:
- Skills management
- Interests tracking
- Skill relationships with projects, posts, certifications
- Skill analytics and recommendations

## Architecture

This service follows the same patterns as the main Woragis backend:
- **Go Fiber** web framework
- **GORM** for database operations
- **PostgreSQL** for data persistence
- **Redis** for caching and session storage
- **OpenTelemetry** for distributed tracing
- **Prometheus** for metrics
- **Structured logging** with trace ID support

## Project Structure

```
skills/
├── server/                    # Go application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go       # Application entry point
│   ├── internal/
│   │   ├── config/           # Configuration management
│   │   ├── database/         # Database connection and management
│   │   └── domains/
│   │       ├── skills/          # Skills domain
│   │       └── interests/      # Interests domain
│   └── pkg/                  # Shared packages
│       ├── auth/             # JWT and password utilities
│       ├── health/           # Health check utilities
│       ├── logger/           # Structured logging
│       ├── metrics/          # Prometheus metrics
│       ├── middleware/       # Fiber middleware
│       ├── security/         # Security middleware
│       ├── timeout/          # Timeout utilities
│       ├── tracing/          # OpenTelemetry tracing
│       └── utils/            # Utility functions
├── docker-compose.yml        # Development environment
├── docker-compose.test.yml   # Test environment
└── .github/
    └── workflows/            # CI/CD pipelines
```

## API Endpoints

### Protected Endpoints (Require Authentication via Auth Service)

- `GET /api/v1/skills` - List skills
- `POST /api/v1/skills` - Create skill
- `GET /api/v1/skills/:id` - Get skill
- `PUT /api/v1/skills/:id` - Update skill
- `DELETE /api/v1/skills/:id` - Delete skill
- `GET /api/v1/interests` - List interests
- `POST /api/v1/interests` - Create interest
- `GET /api/v1/interests/:id` - Get interest
- `PUT /api/v1/interests/:id` - Update interest
- `DELETE /api/v1/interests/:id` - Delete interest

### System Endpoints

- `GET /healthz` - Health check
- `GET /metrics` - Prometheus metrics

## Environment Variables

```bash
# Application
APP_ENV=development
APP_NAME=skills-service
APP_PORT=3000

# Database
DATABASE_URL=postgres://user:password@localhost:5432/skills_service?sslmode=disable
POSTGRES_USER=woragis
POSTGRES_PASSWORD=password
POSTGRES_DB=skills_service

# Auth Service (for JWT validation)
AUTH_SERVICE_URL=http://auth-service:3000

# Redis
REDIS_URL=redis://localhost:6379/0

# CORS
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=http://localhost:5173

# Monitoring
OTLP_ENDPOINT=http://jaeger:4318
JAEGER_ENDPOINT=http://jaeger:4318
```

## Development

### Prerequisites

- Go 1.25.1+
- Docker and Docker Compose
- PostgreSQL 15+
- Redis 7+

### Running Locally

1. **Start dependencies:**
   ```bash
   docker-compose up -d database redis
   ```

2. **Run migrations:**
   ```bash
   cd server
   go run cmd/server/main.go
   ```
   (Migrations run automatically on startup)

3. **Run the service:**
   ```bash
   cd server
   go run cmd/server/main.go
   ```

### Running with Docker Compose

```bash
docker-compose up -d
```

The service will be available at `http://localhost:3000`

## Testing

### Run Tests

```bash
cd server
go test ./...
```

### Integration Tests

```bash
docker-compose -f docker-compose.test.yml up -d
cd server
go test -tags=integration ./...
```

## CI/CD

The service has its own CI/CD pipeline:

- **CI**: Runs on push/PR to `main` or `develop` branches
  - Unit tests
  - Integration tests
  - Linting
  - Docker build

- **CD**: Runs on version tag push (e.g., `v1.0.0`)
  - Build and push Docker image
  - Deploy to production

## Database Schema

The service creates the following tables:
- `skills` - Skills records
- `project_skills` - Skills linked to projects
- `interests` - User interests

## Security Features

- **JWT Validation**: Validates JWT tokens via Auth Service
- **Rate Limiting**: 100 requests per minute per IP/user
- **Security Headers**: Helmet middleware for security headers
- **Input Sanitization**: Automatic input sanitization
- **Request Size Limits**: 10MB maximum request size
- **User Isolation**: All operations are scoped to authenticated user

## Monitoring

- **Health Checks**: `/healthz` endpoint
- **Metrics**: `/metrics` endpoint (Prometheus)
- **Tracing**: OpenTelemetry integration with Jaeger
- **Logging**: Structured JSON logging with trace IDs

## Integration with Other Services

This service integrates with:
- **Auth Service**: Validates JWT tokens and gets user information
- **Translation Service**: Provides multi-language support for skills and interests

## License

Proprietary - Woragis Platform

