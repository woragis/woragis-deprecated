# Translation Worker Integration Tests

This directory contains integration tests for the translation worker service.

## Prerequisites

1. **Docker and Docker Compose** - Required to run test dependencies
2. **RabbitMQ** - Message queue for translation job processing
3. **PostgreSQL** - Database for storing translation records

## Running Tests

### Using Docker Compose (Recommended)

From the `backend/server` directory:

```bash
# Start test dependencies
docker-compose -f docker-compose.test.yml up -d

# Run tests
cd ../../translation-worker
go test -tags=integration ./internal/integration/...

# Cleanup
cd ../server
docker-compose -f docker-compose.test.yml down
```

### Manual Setup

1. Start RabbitMQ and PostgreSQL locally:
   ```bash
   docker run -d --name test-rabbitmq \
     -p 5673:5672 \
     -e RABBITMQ_DEFAULT_USER=test \
     -e RABBITMQ_DEFAULT_PASS=test \
     -e RABBITMQ_DEFAULT_VHOST=test \
     rabbitmq:3.13-management-alpine

   docker run -d --name test-postgres \
     -p 5433:5432 \
     -e POSTGRES_DB=woragis_test \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     postgres:15-alpine
   ```

2. Set environment variables:
   ```bash
   export RABBITMQ_URL="amqp://test:test@localhost:5673/test"
   export DATABASE_URL="postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable"
   ```

3. Run tests:
   ```bash
   go test -tags=integration ./internal/integration/...
   ```

## Test Coverage

### Queue Setup Tests
- ✅ `TestTranslationWorkerQueueSetup` - Tests queue and exchange creation
- ✅ `TestTranslationWorkerMessagePublish` - Tests publishing translation jobs

### Message Processing Tests
- ✅ `TestTranslationWorkerMessageConsume` - Tests consuming and processing jobs with database persistence
- ✅ `TestTranslationWorkerInvalidMessage` - Tests handling of invalid messages
- ✅ `TestTranslationWorkerRetryOnFailure` - Tests retry behavior on translation failure
- ✅ `TestTranslationWorkerMultipleLanguages` - Tests processing translations for multiple languages

## Test Structure

Tests use a mock translator to avoid actually calling translation APIs during testing. The mock translator:
- Simulates translation API calls
- Can simulate failures for retry testing
- Validates translation job processing

Tests also integrate with PostgreSQL to verify that translation records are correctly created and updated in the database.

## Environment Variables

- `RABBITMQ_URL` - RabbitMQ connection URL (default: `amqp://test:test@localhost:5673/test`)
- `DATABASE_URL` - PostgreSQL connection URL (default: `postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable`)
