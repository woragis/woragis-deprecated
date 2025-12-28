# Integration Tests Documentation

**Last Updated:** 2025-12-23  
**Purpose:** Complete guide for integration tests across all services

---

## Overview

Integration tests verify that services work correctly with real dependencies (PostgreSQL, RabbitMQ, Redis). All services have comprehensive integration test suites.

---

## Test Coverage

### Services with Integration Tests

| Service | Language | Test Count | Status |
|---------|----------|------------|--------|
| **Server** | Go | 30+ tests | ✅ Complete |
| **Email Worker** | Go | 15+ tests | ✅ Complete |
| **Translation Worker** | Go | 10+ tests | ✅ Complete |
| **WhatsApp Worker** | Go | 12+ tests | ✅ Complete |
| **Job Application Worker** | Node.js | 16 tests | ✅ Complete |
| **Resume Worker** | Python | 20+ tests | ✅ Complete |
| **AI Service** | Python | 15+ tests | ✅ Complete |
| **Creative Service** | Python | 10+ tests | ✅ Complete |
| **Docs Service** | Python | 12+ tests | ✅ Complete |

**Total:** 129+ integration test functions across all services

---

## Running Integration Tests

### Prerequisites

1. **Test Services Running**
   - PostgreSQL (port 5433)
   - Redis (port 6380)
   - RabbitMQ (port 5673)

2. **Start Test Services**
   ```bash
   cd backend
   docker-compose -f docker-compose.test.yml up -d
   ```

### Go Services

**Server:**
```bash
cd backend/server/app
go test -tags=integration -v ./internal/integration/...
```

**Email Worker:**
```bash
cd backend/email-worker
go test -tags=integration -v ./internal/integration/...
```

**Translation Worker:**
```bash
cd backend/translation-worker
go test -tags=integration -v ./internal/integration/...
```

**WhatsApp Worker:**
```bash
cd backend/whatsapp-worker
go test -tags=integration -v ./internal/integration/...
```

### Node.js Services

**Job Application Worker:**
```bash
cd backend/job-application-worker
npm run test:integration
```

### Python Services

**Resume Worker:**
```bash
cd backend/resume-worker
pytest tests/integration/ -v -m integration
```

**AI Service:**
```bash
cd backend/ai-service
pytest tests/integration/ -v -m integration
```

**Creative Service:**
```bash
cd backend/creative-service
pytest tests/integration/ -v -m integration
```

**Docs Service:**
```bash
cd backend/docs-service
pytest tests/integration/ -v
```

---

## Test Structure

### Go Integration Tests

**Location:** `{service}/internal/integration/`

**Build Tags:**
```go
//go:build integration
// +build integration
```

**Test Files:**
- `{service}_test.go` - Main integration tests
- `performance_test.go` - Performance tests (if applicable)
- `security_test.go` - Security tests (server only)

**Example:**
```go
func TestServiceFunctionality(t *testing.T) {
    db := testutil.SetupTestDB(t)
    defer testutil.CleanupTestDB(t, db)
    
    redis := testutil.SetupTestRedis(t)
    defer testutil.CleanupTestRedis(t, redis)
    
    // Test implementation
}
```

### Node.js Integration Tests

**Location:** `{service}/tests/integration/`

**Test Framework:** Jest

**Example:**
```javascript
describe('Integration Tests', () => {
    beforeAll(async () => {
        // Setup
    });
    
    afterAll(async () => {
        // Cleanup
    });
    
    test('should work correctly', async () => {
        // Test implementation
    });
});
```

### Python Integration Tests

**Location:** `{service}/tests/integration/`

**Test Framework:** pytest

**Markers:**
```python
@pytest.mark.integration
def test_service_functionality():
    # Test implementation
    pass
```

---

## Test Categories

### 1. Database Tests
- CRUD operations
- Transactions
- Migrations
- Query performance

### 2. Queue Tests
- Message publishing
- Message consumption
- Error handling
- Dead letter queues

### 3. Cache Tests
- Cache set/get
- Cache invalidation
- Cache expiration

### 4. API Tests
- Endpoint functionality
- Authentication/authorization
- Request/response validation
- Error handling

### 5. Worker Tests
- Job processing
- Retry logic
- Rate limiting
- Error handling

### 6. Service Integration Tests
- Service-to-service communication
- External API integration
- End-to-end workflows

---

## CI/CD Integration

### GitHub Actions Workflow

**File:** `.github/workflows/integration-tests.yml`

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Manual workflow dispatch

**Services Tested:**
- All 9 backend services
- Go, Python, and Node.js services

**Test Services:**
- PostgreSQL (port 5433)
- Redis (port 6380)
- RabbitMQ (port 5673)

### Running Tests in CI/CD

Tests run automatically on:
- Every push to main/develop
- Every pull request
- Manual trigger via workflow dispatch

**View Results:**
1. Go to GitHub Actions tab
2. Select "Integration Tests" workflow
3. View individual job results
4. Download test artifacts

### Coverage Reporting

**Codecov Integration:**
- All services upload coverage to Codecov
- Coverage reports available in PR comments
- Coverage trends tracked over time

**View Coverage:**
- PR comments (automatic)
- Codecov dashboard
- GitHub Actions artifacts

---

## Test Environment Variables

### Go Services

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable
REDIS_URL=redis://localhost:6380/0
RABBITMQ_URL=amqp://test:test@localhost:5673/test
JWT_SECRET=test-secret-key-for-integration-tests
```

### Python Services

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable
REDIS_URL=redis://localhost:6380/0
RABBITMQ_URL=amqp://test:test@localhost:5673/test
OPENAI_API_KEY=test-key  # For AI/Creative services
```

### Node.js Services

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable
RABBITMQ_URL=amqp://test:test@localhost:5673/test
```

---

## Test Utilities

### Go Test Utilities

**Location:** `backend/server/app/internal/testutil/`

**Functions:**
- `SetupTestDB(t)` - Create test database connection
- `CleanupTestDB(t, db)` - Clean up test database
- `SetupTestRedis(t)` - Create test Redis connection
- `CleanupTestRedis(t, redis)` - Clean up test Redis
- `CreateTestUser(t, db, email, password)` - Create test user
- `GenerateTestJWT(t, userID, email)` - Generate test JWT token

### Python Test Utilities

**Location:** `{service}/tests/conftest.py` or `tests/fixtures.py`

**Fixtures:**
- `test_db` - Database fixture
- `test_redis` - Redis fixture
- `test_rabbitmq` - RabbitMQ fixture

### Node.js Test Utilities

**Location:** `{service}/tests/setup.js` or `tests/helpers.js`

**Helpers:**
- Database setup/teardown
- RabbitMQ connection helpers
- Test data generators

---

## Best Practices

### 1. Isolation
- Each test should be independent
- Clean up after each test
- Use unique test data

### 2. Real Dependencies
- Use real PostgreSQL, Redis, RabbitMQ
- Don't mock dependencies in integration tests
- Test actual interactions

### 3. Test Data
- Use realistic test data
- Clean up test data after tests
- Use factories for test data generation

### 4. Error Scenarios
- Test error cases
- Test edge cases
- Test failure scenarios

### 5. Performance
- Keep tests fast (< 30 seconds per test)
- Use parallel execution when possible
- Avoid unnecessary waits

---

## Troubleshooting

### Tests Fail Locally

**Check:**
1. Test services are running
2. Environment variables are set
3. Ports are not in use
4. Database is accessible

**Common Issues:**
- Test services not running
- Port conflicts
- Database connection issues
- Missing environment variables

### Tests Fail in CI/CD

**Check:**
1. GitHub Actions logs
2. Service health checks
3. Test timeout settings
4. Resource limits

**Common Issues:**
- Services not ready (timing)
- Resource exhaustion
- Network issues
- Test flakiness

### Slow Tests

**Optimize:**
1. Reduce test data size
2. Use parallel execution
3. Optimize database queries
4. Cache test data when possible

---

## Test Maintenance

### Regular Tasks

- **Weekly:** Review test failures
- **Monthly:** Update test data
- **Quarterly:** Review test coverage
- **Annually:** Refactor slow tests

### Adding New Tests

1. Follow existing patterns
2. Use test utilities
3. Document test purpose
4. Add to CI/CD workflow

---

## Related Documentation

- **Testing Patterns:** `docs/development/testing-patterns.md`
- **CI/CD Test Execution:** `docs/development/cicd-test-execution.md`
- **Integration Tests Status:** `docs/PLANNING/INTEGRATION_TESTS_COMPLETE.md`
- **Local Test Results:** `docs/PLANNING/LOCAL_TEST_RESULTS_FINAL.md`

---

**Last Updated:** 2025-12-23

