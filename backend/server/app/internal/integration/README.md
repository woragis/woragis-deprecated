# Integration Tests

Integration tests for the server API endpoints.

## Prerequisites

1. Docker and Docker Compose installed
2. Test dependencies running (PostgreSQL, Redis, RabbitMQ)

## Setup

Start test dependencies:

```bash
cd backend/server
docker-compose -f docker-compose.test.yml up -d
```

Wait for services to be healthy, then run tests:

```bash
# Run all integration tests
go test ./app/internal/integration/... -tags=integration -v

# Run specific test
go test ./app/internal/integration/... -tags=integration -v -run TestServerHealthCheck
```

## Test Structure

- `server_test.go` - Server API endpoint tests
- Tests use real database, Redis, and RabbitMQ connections
- Each test cleans up after itself
- Tests are isolated and can run in parallel

## Environment Variables

See `testutil/README.md` for test configuration.

## Notes

- Tests use a separate test database (`woragis_test`)
- Tests use different ports to avoid conflicts with development services
- Database is cleaned before each test
- Redis is flushed before each test
