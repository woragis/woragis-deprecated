# Job Application Worker Integration Tests

This directory contains integration tests for the Job Application Worker that test real interactions with RabbitMQ and PostgreSQL.

## Prerequisites

### Required Services

1. **RabbitMQ** - Message queue for job processing
2. **PostgreSQL** - Database for job applications and website tracking

### Test Configuration

The tests use separate test instances of RabbitMQ and PostgreSQL to avoid affecting production data.

**RabbitMQ:**
- URL: `amqp://test:test@localhost:5673/test`
- Port: `5673` (different from production `5672`)
- VHost: `test`

**PostgreSQL:**
- Host: `localhost`
- Port: `5433` (different from production `5432`)
- Database: `woragis_test`
- User: `postgres`
- Password: `postgres`

### Environment Variables

Set these environment variables to override defaults:

```bash
export TEST_RABBITMQ_URL="amqp://test:test@localhost:5673/test"
export TEST_DB_HOST="localhost"
export TEST_DB_PORT="5433"
export TEST_DB_NAME="woragis_test"
export TEST_DB_USER="postgres"
export TEST_DB_PASSWORD="postgres"
```

## Running Tests

### Start Test Services

Use `docker-compose.test.yml` to start test services:

```bash
docker-compose -f docker-compose.test.yml up -d
```

### Run Integration Tests

```bash
# Run all integration tests
npm run test:jest -- tests/integration/

# Run specific test file
npm run test:jest -- tests/integration/job_application_worker.test.js
npm run test:jest -- tests/integration/queue.test.js
```

### Run with Coverage

```bash
npm run test:jest:coverage -- tests/integration/
```

## Test Files

### `job_application_worker.test.js`

Tests the full job application worker flow:
- RabbitMQ queue setup and message consumption
- Database operations (create, find, update job applications)
- Website rate limiting
- End-to-end job processing

### `queue.test.js`

Tests RabbitMQ queue operations:
- Connection and channel setup
- Exchange and queue declaration
- Message publishing and consumption
- Dead letter queue handling
- Prefetch configuration

## Test Structure

Each test:
1. Sets up test dependencies (RabbitMQ, database)
2. Runs the test
3. Cleans up test data

Tests are isolated and can run independently.

## Skipping Tests

If RabbitMQ or PostgreSQL are not available, tests will automatically skip with a warning message. This allows tests to run in CI/CD even if services are temporarily unavailable.

## CI/CD Integration

Integration tests should run:
- On pull requests (to verify changes)
- On merge to main (to ensure stability)
- Before production deployments

Add to CI/CD pipeline:

```yaml
- name: Run Integration Tests
  run: |
    docker-compose -f docker-compose.test.yml up -d
    sleep 10  # Wait for services to be ready
    npm run test:jest -- tests/integration/
  env:
    TEST_RABBITMQ_URL: "amqp://test:test@localhost:5673/test"
    TEST_DB_HOST: "localhost"
    TEST_DB_PORT: "5433"
```

## Troubleshooting

### Tests Skipping

If tests are skipping, check:
1. RabbitMQ is running on port 5673
2. PostgreSQL is running on port 5433
3. Test database exists: `woragis_test`
4. Test user has proper permissions

### Connection Errors

- Verify services are running: `docker-compose -f docker-compose.test.yml ps`
- Check service logs: `docker-compose -f docker-compose.test.yml logs`
- Verify ports are not in use: `netstat -an | grep 5673` and `netstat -an | grep 5433`

### Database Errors

- Ensure test database exists: `CREATE DATABASE woragis_test;`
- The docker-compose.test.yml automatically creates the database with user `postgres`
- If running manually, ensure PostgreSQL is accessible with user `postgres` and password `postgres`

## Notes

- Tests use separate test queues and databases to avoid affecting production
- Tests clean up after themselves, but manual cleanup may be needed if tests fail
- Mock external services (Playwright, AI service) are not tested in integration tests (use unit tests)

