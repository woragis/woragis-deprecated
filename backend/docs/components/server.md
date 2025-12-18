# Server Component

## Overview

The main API server for the Woragis platform. Built with Go and Fiber framework, it handles all HTTP requests, authentication, database operations, and job publishing to message queues.

## Architecture

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Port**: 8080
- **Database**: PostgreSQL (via GORM)
- **Cache**: Redis
- **Message Queue**: RabbitMQ (primary), Redis (fallback)

## Responsibilities

1. **REST API Endpoints**: Handle HTTP requests from frontend and mobile clients
2. **Authentication/Authorization**: JWT-based auth, OAuth integration (GitHub, Google, Microsoft)
3. **Database Operations**: CRUD operations for all domains
4. **Job Publishing**: Publish jobs to RabbitMQ for async processing
5. **Caching**: Redis caching for frequently accessed data
6. **WebSocket Support**: Real-time communication for chats

## Health Check

**Endpoint**: `GET /healthz`

**Checks**:
- Database (PostgreSQL) - CRITICAL
- Redis - CRITICAL
- RabbitMQ - OPTIONAL (degraded if failed)

**Response**:
```json
{
  "status": "healthy|degraded|unhealthy",
  "checks": [
    {"name": "database", "status": "ok"},
    {"name": "redis", "status": "ok"},
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

## Metrics

**Endpoint**: `GET /metrics`

Exposes Prometheus metrics:
- HTTP request rate and latency
- Database query duration
- External API call metrics
- Connection pool metrics

## Configuration

### Environment Variables

#### Required
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection URL
- `JWT_SECRET` - JWT signing secret
- `ENV` - Environment (development/production)

#### Optional
- `RABBITMQ_URL` - RabbitMQ connection URL (optional, falls back to Redis)
- `PORT` - Server port (default: 8080)
- `CORS_ENABLED` - Enable CORS (default: true)
- `CORS_ALLOWED_ORIGINS` - Comma-separated origins

### Database Configuration
- Connection pool: 20-50 connections (configurable)
- Query timeout: 5 seconds
- Migration: Auto-migrate on startup (development)

## API Structure

### Domains

The server is organized into domains, each with:
- **Repository**: Database operations
- **Service**: Business logic
- **Handler**: HTTP handlers
- **Routes**: Route definitions

**Domains**:
- `auth` - Authentication and authorization
- `projects` - Project management
- `resumes` - Resume management
- `jobapplications` - Job application tracking
- `chats` - Chat functionality
- `posts` - Blog posts
- `translations` - Translation management
- `skills` - Skills management
- `certifications` - Certifications
- `experiences` - Work experiences
- `testimonials` - Testimonials
- `clients` - Client management
- `finances` - Financial tracking
- `reports` - Reporting
- `socialmediaposts` - Social media posts
- `creativeassets` - Creative assets
- And more...

### API Endpoints

**Base URL**: `http://localhost:8080` (development)

**Common Patterns**:
- `GET /api/{domain}` - List resources
- `GET /api/{domain}/:id` - Get resource
- `POST /api/{domain}` - Create resource
- `PUT /api/{domain}/:id` - Update resource
- `DELETE /api/{domain}/:id` - Delete resource

**Authentication**:
- Most endpoints require JWT token in `Authorization` header
- `Bearer {token}` format

## Message Queue Integration

### RabbitMQ (Primary)

**Exchanges**:
- `woragis.tasks` - Main task exchange

**Queues**:
- `emails.queue` - Email sending jobs
- `whatsapp.queue` - WhatsApp messaging jobs
- `translations.queue` - Translation jobs
- `resumes.queue` - Resume generation jobs
- `job-applications.queue` - Job application jobs

### Redis Fallback

If RabbitMQ is unavailable:
- Falls back to Redis queue
- Health check shows "degraded" status
- System continues operating

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `server`

**Log Levels**:
- `DEBUG` - Detailed debugging (development only)
- `INFO` - General information
- `WARN` - Warnings
- `ERROR` - Errors
- `FATAL` - Fatal errors

**Trace IDs**: Automatically included in logs for request correlation

## Deployment

### Local Development

```bash
cd backend/server
go run app/cmd/server/main.go
```

### Docker

```bash
docker build -t woragis/server .
docker run -p 8080:8080 --env-file .env woragis/server
```

### Kubernetes

See deployment manifests in `server/k8s/` (if available)

## Scaling

### Horizontal Scaling
- Stateless design allows multiple replicas
- Load balancer distributes requests
- Database connection pool per replica

### Resource Requirements
- **CPU**: 500m-1000m (0.5-1 core)
- **Memory**: 512Mi-1Gi
- **Database Connections**: 20-50 per replica

## Security

### Authentication
- JWT tokens
- OAuth integration (GitHub, Google, Microsoft)
- Session management

### Authorization
- Role-based access control (RBAC)
- Resource-level permissions

### Security Headers
- CORS configuration
- Rate limiting (if implemented)

## Monitoring

### Metrics
- Request rate (req/s)
- Latency (p50, p95, p99)
- Error rate
- Database connection pool usage

### Logs
- Structured JSON logs
- Trace ID correlation
- Error tracking

### Health Checks
- `/healthz` endpoint
- Dependency checks
- Kubernetes liveness/readiness probes

## Troubleshooting

### Common Issues

#### High Latency
- Check database query performance
- Check external API calls
- Check connection pool usage

#### Database Connection Errors
- Check connection pool limits
- Check database health
- Verify DATABASE_URL

#### RabbitMQ Fallback Active
- Check RabbitMQ connection
- Check RabbitMQ health
- System continues operating in degraded mode

## Related Documentation

- [Architecture Decision Records](../adr/) - Architectural decisions
- [API Documentation](../api/server-api.md) - API endpoint details
- [Development Guides](../development/) - How to extend the server
