# Woragis Microservices Port Scheme

This document outlines the port allocation for all Woragis microservices and their dependencies.

## Service Ports (HTTP)

All services run internally on port `3000` and are mapped to external ports starting at `3010`:

| Service | External Port | Internal Port | URL |
|---------|--------------|---------------|-----|
| Auth Service | 3010 | 3000 | http://localhost:3010 |
| Jobs Service | 3011 | 3000 | http://localhost:3011 |
| Skills Service | 3012 | 3000 | http://localhost:3012 |
| Management Service | 3013 | 3000 | http://localhost:3013 |
| Posts Service | 3014 | 3000 | http://localhost:3014 |
| Social Media Service | 3015 | 3000 | http://localhost:3015 |
| Files Service | 3016 | 3000 | http://localhost:3016 |

## PostgreSQL Ports

Each service has its own PostgreSQL database:

| Service | External Port | Internal Port | Database Name |
|---------|--------------|---------------|---------------|
| Auth | 5443 | 5432 | auth_service |
| Jobs | 5444 | 5432 | jobs_service |
| Skills | 5445 | 5432 | skills_service |
| Management | 5446 | 5432 | management_service |
| Posts | 5447 | 5432 | posts_service |
| Social Media | 5448 | 5432 | social_media_service |
| Files | 5449 | 5432 | files_service |

**Connection String Format:**
```
postgres://woragis:password@localhost:<EXTERNAL_PORT>/<DATABASE_NAME>?sslmode=disable
```

Example for Auth Service:
```
postgres://woragis:password@localhost:5443/auth_service?sslmode=disable
```

## Redis Ports

Each service has its own Redis instance:

| Service | External Port | Internal Port |
|---------|--------------|---------------|
| Auth | 6391 | 6379 |
| Jobs | 6392 | 6379 |
| Skills | 6393 | 6379 |
| Management | 6394 | 6379 |
| Posts | 6395 | 6379 |
| Social Media | 6396 | 6379 |
| Files | 6397 | 6379 |

**Connection String Format:**
```
redis://localhost:<EXTERNAL_PORT>/0
```

Example for Auth Service:
```
redis://localhost:6391/0
```

## RabbitMQ Ports

Most services have their own RabbitMQ instance (except Files service):

| Service | AMQP Port (External) | Management UI (External) | AMQP Port (Internal) | Management UI (Internal) |
|---------|---------------------|-------------------------|---------------------|-------------------------|
| Auth | 5681 | 15681 | 5672 | 15672 |
| Jobs | 5682 | 15682 | 5672 | 15672 |
| Skills | 5683 | 15683 | 5672 | 15672 |
| Management | 5684 | 15684 | 5672 | 15672 |
| Posts | 5685 | 15685 | 5672 | 15672 |
| Social Media | 5686 | 15686 | 5672 | 15672 |

**Connection String Format:**
```
amqp://woragis:woragis@localhost:<AMQP_EXTERNAL_PORT>/
```

**Management UI URLs:**
```
http://localhost:<MGMT_EXTERNAL_PORT>
```

Example for Auth Service:
- AMQP: `amqp://woragis:woragis@localhost:5681/`
- Management UI: `http://localhost:15681`

## Inter-Service Communication

Services communicate with each other using Docker service names on the `woragis-network`:

### Auth Service URL
All services use this URL to communicate with the Auth service:
```
http://auth-service:3000
```

This is set via the `AUTH_SERVICE_URL` environment variable.

### Internal Network Communication

Within Docker, services use service names and internal ports:
- Auth Service: `http://auth-service:3000`
- Jobs Service: `http://jobs-service:3000`
- Skills Service: `http://skills-service:3000`
- Management Service: `http://management-service:3000`
- Posts Service: `http://posts-service:3000`
- Social Media Service: `http://social-media-service:3000`
- Files Service: `http://files-service:3000`

## Running Services

### Start All Services
```bash
docker-compose -f docker-compose.services.yml up -d
```

### Start Specific Service
```bash
docker-compose -f docker-compose.services.yml up -d auth-service auth-database auth-redis auth-rabbitmq
```

### View Logs
```bash
# All services
docker-compose -f docker-compose.services.yml logs -f

# Specific service
docker-compose -f docker-compose.services.yml logs -f auth-service
```

### Stop All Services
```bash
docker-compose -f docker-compose.services.yml down
```

### Stop and Remove Volumes
```bash
docker-compose -f docker-compose.services.yml down -v
```

## Health Checks

All services expose health check endpoints:
- `GET http://localhost:<SERVICE_PORT>/healthz` - Combined health check
- `GET http://localhost:<SERVICE_PORT>/healthz/live` - Liveness probe
- `GET http://localhost:<SERVICE_PORT>/healthz/ready` - Readiness probe

## API Endpoints

All services use the `/api/v1` prefix for their API endpoints:

- Auth: `http://localhost:3010/api/v1/auth/...`
- Jobs: `http://localhost:3011/api/v1/job-applications/...`
- Skills: `http://localhost:3012/api/v1/skills/...`
- Management: `http://localhost:3013/api/v1/projects/...`
- Posts: `http://localhost:3014/api/v1/posts/...`
- Social Media: `http://localhost:3015/api/v1/social-media-posts/...`
- Files: `http://localhost:3016/api/v1/files/...`

## Environment Variables

Key environment variables that can be set:

- `POSTGRES_PASSWORD` - PostgreSQL password (default: `password`)
- `POSTGRES_USER` - PostgreSQL user (default: `woragis`)
- `RABBITMQ_USER` - RabbitMQ user (default: `woragis`)
- `RABBITMQ_PASSWORD` - RabbitMQ password (default: `woragis`)
- `AUTH_JWT_SECRET` - JWT secret for auth service (default: `dev-secret-change-me`)
- `APP_ENV` - Application environment (default: `development`)
- `CORS_ALLOWED_ORIGINS` - Comma-separated list of allowed CORS origins

## Notes

1. All services share the same `woragis-network` Docker network for inter-service communication
2. Each service has its own isolated database, Redis, and RabbitMQ instances
3. Services communicate internally using service names (e.g., `auth-service:3000`)
4. External access uses mapped ports (e.g., `localhost:3010`)
5. The Files service does not require RabbitMQ as it doesn't use message queues

