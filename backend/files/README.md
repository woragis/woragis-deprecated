# Jobs Service

Job application and resume management service for Woragis platform.

## Overview

The Jobs Service is a standalone microservice responsible for:
- Job application tracking
- Resume management
- Job website/platform tracking
- Interview stage management
- Application response tracking
- Cover letter generation (via AI service)

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
jobs/
├── server/                    # Go application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go       # Application entry point
│   ├── internal/
│   │   ├── config/           # Configuration management
│   │   ├── database/         # Database connection and management
│   │   └── domains/
│   │       ├── jobapplications/  # Job applications domain
│   │       ├── resumes/         # Resumes domain
│   │       └── jobwebsites/     # Job websites domain
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

- `GET /api/v1/job-applications` - List job applications
- `POST /api/v1/job-applications` - Create job application
- `GET /api/v1/job-applications/:id` - Get job application
- `PUT /api/v1/job-applications/:id` - Update job application
- `DELETE /api/v1/job-applications/:id` - Delete job application
- `GET /api/v1/job-applications/:id/interview-stages` - Get interview stages
- `GET /api/v1/job-applications/:id/responses` - Get application responses
- `GET /api/v1/resumes` - List resumes
- `POST /api/v1/resumes` - Create resume
- `GET /api/v1/resumes/:id` - Get resume
- `PUT /api/v1/resumes/:id` - Update resume
- `DELETE /api/v1/resumes/:id` - Delete resume
- `GET /api/v1/job-websites` - List job websites
- `POST /api/v1/job-websites` - Create job website

### System Endpoints

- `GET /healthz` - Health check
- `GET /metrics` - Prometheus metrics

## Environment Variables

```bash
# Application
APP_ENV=development
APP_NAME=jobs-service
APP_PORT=3000

# Database
DATABASE_URL=postgres://user:password@localhost:5432/jobs_service?sslmode=disable
POSTGRES_USER=woragis
POSTGRES_PASSWORD=password
POSTGRES_DB=jobs_service

# Auth Service (for JWT validation)
AUTH_SERVICE_URL=http://auth-service:3000

# Redis
REDIS_URL=redis://localhost:6379/0

# AI Service (for cover letter generation)
AI_SERVICE_URL=http://ai-service:8000

# Creative Service (for resume generation)
CREATIVE_SERVICE_URL=http://creative-service:8000

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
- `job_applications` - Job application records
- `interview_stages` - Interview stage tracking
- `responses` - Application response tracking
- `resumes` - Resume records
- `job_websites` - Job website/platform tracking

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
- **AI Service**: Generates cover letters for job applications
- **Creative Service**: Generates resumes (via resume-worker)
- **RabbitMQ**: Processes job application tasks asynchronously

## License

Proprietary - Woragis Platform

