# Test Utilities

This package provides utilities for integration testing.

## Setup

Before running integration tests, start the test dependencies:

```bash
cd backend/server
docker-compose -f docker-compose.test.yml up -d
```

Wait for all services to be healthy, then run tests:

```bash
go test ./app/internal/... -tags=integration -v
```

## Usage

```go
import "github.com/woragis/backend/server/app/internal/testutil"

func TestMyIntegration(t *testing.T) {
    db := testutil.SetupTestDB(t)
    defer testutil.CleanupTestDB(t, db)
    
    redis := testutil.SetupTestRedis(t)
    defer testutil.CleanupTestRedis(t, redis)
    
    app := testutil.SetupTestApp(t, db, redis)
    
    // Your test code here
}
```

## Environment Variables

- `TEST_DATABASE_URL` - PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable`)
- `TEST_REDIS_URL` - Redis connection string (default: `redis://localhost:6380/0`)
- `TEST_RABBITMQ_URL` - RabbitMQ connection string (default: `amqp://test:test@localhost:5673/test`)
- `TEST_JWT_SECRET` - JWT secret for tests (default: `test-secret-key-for-integration-tests`)
