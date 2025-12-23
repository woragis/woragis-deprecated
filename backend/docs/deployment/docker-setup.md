# Docker Setup Guide

**Last Updated:** 2025-12-22  
**Purpose:** Comprehensive guide for Docker and Docker Compose setup

---

## Overview

The Woragis backend uses Docker and Docker Compose for containerization and orchestration. This guide covers:
- Docker installation and setup
- Docker Compose configuration
- Container management
- Volume and network configuration
- Troubleshooting

---

## Prerequisites

- **Docker Desktop** (Windows/Mac) or **Docker Engine** (Linux) 20.10+
- **Docker Compose** 2.0+
- At least **4GB RAM** available for Docker
- At least **10GB disk space** for images and volumes

---

## Installation

### Windows

1. Download Docker Desktop: https://www.docker.com/products/docker-desktop
2. Install and restart your computer
3. Start Docker Desktop
4. Verify installation:
   ```powershell
   docker --version
   docker-compose --version
   ```

### macOS

1. Download Docker Desktop: https://www.docker.com/products/docker-desktop
2. Install and start Docker Desktop
3. Verify installation:
   ```bash
   docker --version
   docker-compose --version
   ```

### Linux

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version

# Add user to docker group (optional, to run without sudo)
sudo usermod -aG docker $USER
# Log out and back in for changes to take effect
```

---

## Docker Compose Overview

The `docker-compose.yml` file defines all services, networks, and volumes.

### Services

**Infrastructure:**
- `database` - PostgreSQL 15
- `redis` - Redis 7
- `rabbitmq` - RabbitMQ 3.13

**Application Services:**
- `app` - Main API server (Go)
- `ai-service` - AI service (Python/FastAPI)
- `creative-service` - Creative service (Python/FastAPI)
- `docs-service` - Docs service (Python/FastAPI)

**Workers:**
- `resume-worker` - Resume generation (Python)
- `translation-worker` - Translation (Go)
- `email-worker` - Email sending (Go)
- `whatsapp-worker` - WhatsApp messaging (Go)
- `job-application-worker` - Job applications (Node.js)

**Monitoring:**
- `prometheus` - Metrics collection
- `loki` - Log aggregation
- `promtail` - Log shipper
- `grafana` - Visualization
- `jaeger` - Distributed tracing

---

## Basic Commands

### Starting Services

```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d app

# Start with build
docker-compose up -d --build

# Start and view logs
docker-compose up
```

### Stopping Services

```bash
# Stop all services
docker-compose down

# Stop specific service
docker-compose stop app

# Stop and remove volumes (⚠️ deletes data)
docker-compose down -v
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f app

# Last 100 lines
docker-compose logs --tail=100 app

# Since timestamp
docker-compose logs --since 10m app
```

### Checking Status

```bash
# List all services
docker-compose ps

# Check specific service
docker-compose ps app

# Show resource usage
docker stats
```

---

## Container Management

### Executing Commands in Containers

```bash
# Execute command in running container
docker-compose exec app sh

# Run one-off command
docker-compose run --rm app go test ./...

# Access database
docker-compose exec database psql -U postgres -d woragis

# Access Redis CLI
docker-compose exec redis redis-cli
```

### Rebuilding Containers

```bash
# Rebuild specific service
docker-compose build app

# Rebuild without cache
docker-compose build --no-cache app

# Rebuild and restart
docker-compose up -d --build app
```

### Restarting Services

```bash
# Restart all services
docker-compose restart

# Restart specific service
docker-compose restart app

# Force recreate (picks up config changes)
docker-compose up -d --force-recreate app
```

---

## Volume Management

### Viewing Volumes

```bash
# List all volumes
docker volume ls

# Inspect volume
docker volume inspect backend_postgres-data

# View volume usage
docker system df -v
```

### Backup Volumes

```bash
# Backup PostgreSQL data
docker run --rm -v backend_postgres-data:/data -v $(pwd):/backup alpine tar czf /backup/postgres-backup.tar.gz /data

# Backup Redis data
docker run --rm -v backend_redis-data:/data -v $(pwd):/backup alpine tar czf /backup/redis-backup.tar.gz /data
```

### Restore Volumes

```bash
# Restore PostgreSQL data
docker run --rm -v backend_postgres-data:/data -v $(pwd):/backup alpine sh -c "cd /data && tar xzf /backup/postgres-backup.tar.gz"
```

### Cleaning Volumes

```bash
# Remove unused volumes
docker volume prune

# Remove specific volume (⚠️ deletes data)
docker volume rm backend_postgres-data
```

---

## Network Configuration

### Viewing Networks

```bash
# List networks
docker network ls

# Inspect network
docker network inspect backend_default

# View network details
docker-compose config | grep -A 10 networks
```

### Connecting to Network

```bash
# Connect container to network
docker network connect backend_default my-container

# Disconnect from network
docker network disconnect backend_default my-container
```

---

## Health Checks

All services have health checks configured. Check health status:

```bash
# Check all services
docker-compose ps

# Check specific service health
docker inspect woragis-app | grep -A 10 Health

# Wait for service to be healthy
docker-compose up -d
docker-compose ps  # Check status
```

**Health Check Endpoints:**
- Main app: `http://localhost:8080/healthz`
- AI service: `http://localhost:8000/health`
- Database: `pg_isready -U postgres`
- Redis: `redis-cli ping`
- RabbitMQ: `rabbitmq-diagnostics ping`

---

## Resource Limits

### Setting Resource Limits

Edit `docker-compose.yml`:

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### Monitoring Resource Usage

```bash
# Real-time stats
docker stats

# Container resource usage
docker stats woragis-app --no-stream

# System-wide usage
docker system df
```

---

## Environment Variables

### Using .env File

Create `.env` file in `backend/` directory:

```bash
# .env
DATABASE_URL=postgres://postgres:postgres@database:5432/woragis?sslmode=disable
REDIS_URL=redis://redis:6379/0
```

Docker Compose automatically loads `.env` file.

### Overriding Environment Variables

```bash
# Set in command
docker-compose up -d -e DATABASE_URL=postgres://...

# Or in docker-compose.yml
services:
  app:
    environment:
      - DATABASE_URL=${DATABASE_URL:-default-value}
```

---

## Troubleshooting

### Issue: Container Won't Start

**Check logs:**
```bash
docker-compose logs app
```

**Check if port is in use:**
```bash
# Linux/Mac
lsof -i :8080

# Windows
netstat -ano | findstr :8080
```

**Check container status:**
```bash
docker-compose ps app
docker inspect woragis-app
```

### Issue: Database Connection Failed

**Check database is running:**
```bash
docker-compose ps database
```

**Check database logs:**
```bash
docker-compose logs database
```

**Test connection:**
```bash
docker-compose exec database psql -U postgres -c "SELECT 1"
```

### Issue: Services Can't Communicate

**Check network:**
```bash
docker network inspect backend_default
```

**Verify service names:**
- Use service names (e.g., `database`, `redis`) not `localhost`
- Check `docker-compose.yml` for correct service names

**Test connectivity:**
```bash
docker-compose exec app ping database
docker-compose exec app nc -zv redis 6379
```

### Issue: Volume Permission Errors

**Fix permissions (Linux):**
```bash
sudo chown -R $USER:$USER /var/lib/docker/volumes/
```

**Or recreate volume:**
```bash
docker-compose down -v
docker-compose up -d
```

### Issue: Out of Disk Space

**Clean up:**
```bash
# Remove unused containers, networks, images
docker system prune

# Remove unused volumes
docker volume prune

# Remove unused images
docker image prune -a
```

### Issue: Container Keeps Restarting

**Check logs:**
```bash
docker-compose logs --tail=50 app
```

**Check health status:**
```bash
docker inspect woragis-app | grep -A 20 Health
```

**Common causes:**
- Missing environment variables
- Database not ready
- Port conflicts
- Configuration errors

---

## Production Considerations

### Security

1. **Use secrets management:**
   ```yaml
   services:
     app:
       secrets:
         - db_password
   secrets:
     db_password:
       external: true
   ```

2. **Limit container resources**
3. **Use read-only filesystems where possible**
4. **Scan images for vulnerabilities**
5. **Use non-root users in containers**

### Performance

1. **Set resource limits** to prevent resource exhaustion
2. **Use health checks** for automatic recovery
3. **Configure log rotation** to prevent disk fill
4. **Use persistent volumes** for data that needs to survive
5. **Monitor resource usage** regularly

### High Availability

1. **Use Docker Swarm or Kubernetes** for production
2. **Set up service replicas**
3. **Configure health checks**
4. **Use load balancers**
5. **Implement backup strategies**

---

## Best Practices

1. **Use `.env` files** for environment-specific configs
2. **Never commit `.env` files** to version control
3. **Use health checks** for all services
4. **Set resource limits** to prevent resource exhaustion
5. **Use named volumes** for persistent data
6. **Keep images updated** for security patches
7. **Monitor container logs** regularly
8. **Use Docker Compose profiles** for different environments
9. **Document all custom configurations**
10. **Test locally before deploying**

---

## Advanced Topics

### Docker Compose Profiles

```yaml
services:
  monitoring:
    profiles:
      - monitoring
```

Start with profile:
```bash
docker-compose --profile monitoring up -d
```

### Multi-Environment Compose

```bash
# docker-compose.dev.yml
# docker-compose.prod.yml

docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Custom Networks

```yaml
networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
```

---

## Related Documentation

- [Development Setup Guide](../development/setup-guide.md)
- [Configuration Reference](./configuration.md)
- [Deployment Procedures](./deployment-procedures.md) (when created)
- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

---

**Last Updated:** 2025-12-22
