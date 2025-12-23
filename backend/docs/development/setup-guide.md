# Development Setup Guide

**Last Updated:** 2025-12-22  
**Purpose:** Step-by-step guide to set up the Woragis backend for local development

---

## Prerequisites

Before you begin, ensure you have the following installed:

### Required Software

- **Docker** 20.10+ and **Docker Compose** 2.0+
  - Download: https://www.docker.com/products/docker-desktop
  - Verify: `docker --version` and `docker-compose --version`

- **Go** 1.21+ (for Go services)
  - Download: https://go.dev/dl/
  - Verify: `go version`

- **Python** 3.11+ (for Python services)
  - Download: https://www.python.org/downloads/
  - Verify: `python --version` or `python3 --version`

- **Node.js** 18+ and **npm** (for Node.js workers)
  - Download: https://nodejs.org/
  - Verify: `node --version` and `npm --version`

- **Git** (for version control)
  - Download: https://git-scm.com/downloads
  - Verify: `git --version`

### Optional Tools

- **Make** (for running build commands)
- **Postman** or **curl** (for API testing)
- **VS Code** or your preferred IDE

---

## Step 1: Clone the Repository

```bash
# Clone the repository
git clone <repository-url>
cd woragis/backend

# Verify you're in the correct directory
ls -la
# You should see: docker-compose.yml, README.md, docs/, etc.
```

---

## Step 2: Environment Configuration

### 2.1 Create Environment File

Create a `.env` file in the `backend/` directory:

```bash
# Copy from example (if exists) or create new
cp .env.example .env  # If example exists
# OR
touch .env
```

### 2.2 Configure Environment Variables

Edit `.env` and set the following **required** variables:

```bash
# Application
APP_NAME=woragis-server
APP_ENV=development
APP_PORT=8080

# Database
DATABASE_URL=postgres://postgres:postgres@database:5432/woragis?sslmode=disable

# Redis
REDIS_URL=redis://redis:6379/0

# RabbitMQ
RABBITMQ_URL=amqp://woragis:woragis@rabbitmq:5672/
RABBITMQ_USER=woragis
RABBITMQ_PASSWORD=woragis
RABBITMQ_VHOST=/

# CORS (for local development)
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173

# AI Service (choose one provider)
CHAT_PROVIDER=openai
OPENAI_API_KEY=your-openai-api-key-here
OPENAI_MODEL=gpt-4o-mini

# Email (optional for local dev)
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=your-email@example.com
SMTP_PASSWORD=your-password
SMTP_FROM=noreply@example.com
SMTP_TLS=true
```

**Note:** See [Configuration Reference](../deployment/configuration.md) for all available environment variables.

---

## Step 3: Start Infrastructure Services

Start the core infrastructure services (database, Redis, RabbitMQ):

```bash
# Start infrastructure only
docker-compose up -d database redis rabbitmq

# Verify they're running
docker-compose ps

# Check health
docker-compose ps database redis rabbitmq
# All should show "healthy" status
```

**Expected output:**
```
NAME                 STATUS
woragis-database     Up X seconds (healthy)
woragis-redis        Up X seconds (healthy)
woragis-rabbitmq     Up X seconds (healthy)
```

---

## Step 4: Start Application Services

### 4.1 Start All Services

```bash
# Start all services and workers
docker-compose up -d

# View logs
docker-compose logs -f

# Or view logs for specific service
docker-compose logs -f app
docker-compose logs -f ai-service
```

### 4.2 Verify Services Are Running

```bash
# Check all services
docker-compose ps

# Test main API
curl http://localhost:8080/healthz

# Test AI service
curl http://localhost:8000/health

# Test creative service
curl http://localhost:8001/health

# Test docs service
curl http://localhost:8002/health
```

**Expected responses:**
- All health endpoints should return `200 OK` or JSON with status

---

## Step 5: Verify Setup

### 5.1 Check Service Health

```bash
# Main server
curl http://localhost:8080/healthz

# AI Service
curl http://localhost:8000/health

# Creative Service
curl http://localhost:8001/health

# Docs Service
curl http://localhost:8002/health
```

### 5.2 Check Database Connection

```bash
# Connect to database
docker-compose exec database psql -U postgres -d woragis -c "SELECT version();"
```

### 5.3 Check RabbitMQ Management UI

Open in browser: http://localhost:15672
- Username: `woragis`
- Password: `woragis`

### 5.4 Check Redis

```bash
# Test Redis connection
docker-compose exec redis redis-cli ping
# Should return: PONG
```

---

## Step 6: Running Services Locally (Without Docker)

If you prefer to run services locally (for debugging):

### 6.1 Setup Database Locally

```bash
# Install PostgreSQL locally or use Docker
# Create database
createdb woragis

# Run migrations (if applicable)
# cd server && go run cmd/migrate/main.go
```

### 6.2 Run Main Server (Go)

```bash
cd server

# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/woragis?sslmode=disable
export REDIS_URL=redis://localhost:6379/0
export RABBITMQ_URL=amqp://woragis:woragis@localhost:5672/

# Run server
go run cmd/server/main.go
```

### 6.3 Run AI Service (Python)

```bash
cd ai-service

# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Set environment variables
export OPENAI_API_KEY=your-key-here
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/woragis?sslmode=disable

# Run service
uvicorn app.main:app --reload --port 8000
```

### 6.4 Run Workers

**Resume Worker (Python):**
```bash
cd resume-worker
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python src/main.py
```

**Translation Worker (Go):**
```bash
cd translation-worker
go mod download
go run cmd/worker/main.go
```

---

## Step 7: Running Tests

### 7.1 Run All Tests

```bash
# Using Make (if available)
make test-all

# Or manually
cd server && go test ./...
cd ai-service && pytest
cd resume-worker && pytest
```

### 7.2 Run Integration Tests

```bash
# Start services first
docker-compose up -d

# Run integration tests
cd server && go test -tags=integration ./...
```

---

## Common Issues & Solutions

### Issue: Port Already in Use

**Error:** `Bind for 0.0.0.0:8080 failed: port is already allocated`

**Solution:**
```bash
# Find process using port
# On Linux/Mac:
lsof -i :8080
# On Windows:
netstat -ano | findstr :8080

# Kill process or change port in docker-compose.yml
```

### Issue: Database Connection Failed

**Error:** `connection refused` or `database does not exist`

**Solution:**
```bash
# Check database is running
docker-compose ps database

# Check database logs
docker-compose logs database

# Restart database
docker-compose restart database

# Wait for health check
docker-compose ps database  # Should show "healthy"
```

### Issue: RabbitMQ Connection Failed

**Error:** `NOT_ALLOWED - vhost not found`

**Solution:**
```bash
# Check RabbitMQ vhost configuration
# In docker-compose.yml, ensure RABBITMQ_VHOST=/
# Restart RabbitMQ
docker-compose restart rabbitmq
```

### Issue: Services Not Starting

**Error:** Services crash or fail to start

**Solution:**
```bash
# Check logs
docker-compose logs <service-name>

# Check environment variables
docker-compose config

# Rebuild containers
docker-compose build <service-name>
docker-compose up -d <service-name>
```

---

## Next Steps

Once setup is complete:

1. **Read the Documentation:**
   - [Architecture Overview](../architecture/system-overview.md)
   - [Development Guides](./)
   - [API Documentation](../api/)

2. **Explore the Codebase:**
   - Start with `server/` (main API)
   - Check `docs/components/` for service details

3. **Run Tests:**
   - Ensure all tests pass
   - Write tests for new features

4. **Start Developing:**
   - Follow [Coding Standards](./coding-standards.md) (when created)
   - Check [Contributing Guide](./contributing.md) (when created)

---

## Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Go Documentation](https://go.dev/doc/)
- [Python Documentation](https://docs.python.org/3/)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [Fiber Documentation](https://docs.gofiber.io/)

---

## Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](../runbooks/troubleshooting.md)
2. Review service logs: `docker-compose logs <service-name>`
3. Check [GitHub Issues](<repository-url>/issues)
4. Ask in team chat or create an issue

---

**Setup Complete!** 🎉

You should now have a fully functional Woragis backend running locally.
