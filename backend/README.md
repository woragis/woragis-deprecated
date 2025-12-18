# Woragis Backend

Backend services and workers for the Woragis platform.

## Overview

The Woragis backend is a microservices architecture consisting of:
- **Server**: Main API server (Go/Fiber)
- **Workers**: Asynchronous job processors (Go, Python, Node.js)
- **Services**: Specialized microservices (Python/FastAPI)

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for Go components)
- Python 3.11+ (for Python components)
- Node.js 18+ (for Node.js components)
- PostgreSQL 14+
- RabbitMQ 3.12+
- Redis 7+

### Local Development

```bash
# Start all services
docker-compose up -d

# Check health
curl http://localhost:8080/healthz

# View logs
docker-compose logs -f
```

## Architecture

See [System Overview](./docs/architecture/system-overview.md) for detailed architecture documentation.

## Documentation

Comprehensive documentation is available in the [`docs/`](./docs/) directory:

- **[Technical Documentation](./docs/README.md)** - Complete technical documentation index
- **[Architecture Decision Records](./docs/adr/)** - Key architectural decisions
- **[Runbooks](./docs/runbooks/)** - Operational procedures
- **[Component Documentation](./docs/components/)** - Detailed component docs
- **[API Documentation](./docs/api/)** - API endpoint documentation
- **[Development Guides](./docs/development/)** - How to extend the system

## Components

### Server
- **Language**: Go
- **Framework**: Fiber
- **Port**: 8080
- **Health**: `/healthz`
- **Metrics**: `/metrics`

### Workers
- **Email Worker** (Go) - Email sending
- **WhatsApp Worker** (Go) - WhatsApp messaging
- **Translation Worker** (Go) - Content translation
- **Resume Worker** (Python) - Resume generation
- **Job Application Worker** (Node.js) - Job application automation

### Services
- **AI Service** (Python/FastAPI) - AI/LLM integration
- **Creative Service** (Python/FastAPI) - Creative content generation
- **Docs Service** (Python/FastAPI) - Technical documentation serving

## Development

### Running Tests
```bash
# Run all tests
make test-all

# Run tests for specific component
cd server && make test
cd email-worker && make test
```

### Building
```bash
# Build all components
make build-all

# Build specific component
cd server && make build
```

### Code Quality
```bash
# Format code
make format

# Lint code
make lint
```

## Deployment

See [Deploying Services and Workers](./docs/runbooks/deploying-services.md) for deployment procedures.

## Monitoring

- **Health Checks**: All components expose `/healthz`
- **Metrics**: All components expose `/metrics` (Prometheus format)
- **Logging**: Structured JSON logging in production

## Contributing

1. Read [Development Guides](./docs/development/)
2. Follow coding standards
3. Write tests
4. Update documentation
5. Submit pull request

## License

[Your License Here]

## Related Links

- [Backend TODO](./TODO.md) - Current tasks and implementation status
- [Testing Guide](./TESTING.md) - Testing instructions
- [Architecture Overview](./docs/architecture/system-overview.md) - System architecture
