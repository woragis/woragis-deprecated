# Testing Guide

This guide explains how to set up and run tests for the backend services.

## Test Structure

The CI/CD pipeline runs tests for all backend components:

- **Go services** (main server, translation-worker, whatsapp-worker)
- **Python services** (resume-worker, ai-service)
- **Node.js services** (job-application-worker)

## Setting Up Tests

### Go Tests

Go tests are located alongside the code with `_test.go` suffix.

**Example test file**: `backend/server/app/internal/domains/jobapplications/service_test.go`

```go
package jobapplications

import (
    "testing"
    "context"
    // ... imports
)

func TestService_RequestJobApplication(t *testing.T) {
    // Test implementation
}
```

**Run tests locally**:
```bash
cd backend/server
go test ./...
go test -v -race ./...  # With race detection
go test -cover ./...   # With coverage
```

### Python Tests (Resume Worker)

Create tests in `backend/server/resume-worker/tests/` or alongside source files.

**Example test file**: `backend/server/resume-worker/tests/test_resume_generator.py`

```python
import pytest
from src.resume_generator import ResumeGenerator

def test_resume_generation():
    # Test implementation
    pass
```

**Run tests locally**:
```bash
cd backend/server/resume-worker
pip install -r requirements.txt
pip install pytest pytest-cov
pytest -v
pytest --cov=src --cov-report=html
```

### Python Tests (AI Service)

Create tests in `backend/ai-service/tests/` or alongside source files.

**Example test file**: `backend/ai-service/tests/test_providers.py`

```python
import pytest
from app.providers import ProviderFactory

def test_provider_factory():
    # Test implementation
    pass
```

**Run tests locally**:
```bash
cd backend/ai-service
pip install -r requirements.txt
pip install pytest pytest-cov
pytest -v
pytest --cov=app --cov-report=html
```

### Node.js Tests (Job Application Worker)

Create tests in `backend/server/job-application-worker/tests/` or use Jest/Mocha.

**Example test file**: `backend/server/job-application-worker/tests/worker.test.js`

```javascript
const { describe, it, expect } = require('@jest/globals');
const worker = require('../src/worker');

describe('Worker', () => {
  it('should process job applications', () => {
    // Test implementation
  });
});
```

**Add test script to package.json**:
```json
{
  "scripts": {
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage"
  }
}
```

**Run tests locally**:
```bash
cd backend/server/job-application-worker
npm install
npm test
```

## Test Database Setup

For integration tests, you'll need a test database. The CI/CD pipeline uses PostgreSQL 15.

**Local setup**:
```bash
# Using Docker
docker run -d \
  --name test-postgres \
  -e POSTGRES_USER=test \
  -e POSTGRES_PASSWORD=test \
  -e POSTGRES_DB=test_db \
  -p 5432:5432 \
  postgres:15

# Connection string
DATABASE_URL=postgres://test:test@localhost:5432/test_db?sslmode=disable
```

## Test Environment Variables

Set these environment variables for tests:

```bash
# Database
DATABASE_URL=postgres://test:test@localhost:5432/test_db?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# Other services (if needed)
AI_SERVICE_URL=http://localhost:8000
```

## Running All Tests

**From project root**:
```bash
# Go tests
cd backend/server && go test ./... && cd ../..

# Python tests (Resume Worker)
cd backend/server/resume-worker && pytest && cd ../../..

# Python tests (AI Service)
cd backend/ai-service && pytest && cd ../..

# Node.js tests
cd backend/server/job-application-worker && npm test && cd ../../..
```

## Test Coverage

The CI/CD pipeline generates coverage reports. To view coverage locally:

**Go**:
```bash
cd backend/server
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Python**:
```bash
cd backend/server/resume-worker
pytest --cov=src --cov-report=html
# Open htmlcov/index.html
```

## Best Practices

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **Test Isolation**: Each test should be independent
4. **Mock External Services**: Don't call real APIs in tests
5. **Use Test Databases**: Never use production databases
6. **Clean Up**: Clean up test data after each test
7. **Test Coverage**: Aim for >80% coverage on critical paths

## CI/CD Integration

Tests run automatically:
- On every push to `main` or `develop` branches
- On every pull request
- Before building Docker images (on tag push)

If tests fail, Docker images will **not** be built or pushed.

## Troubleshooting

**Tests fail in CI but pass locally**:
- Check environment variables
- Verify database/Redis connections
- Check for race conditions (use `-race` flag for Go)

**Coverage not uploading**:
- Coverage upload is optional (won't fail the build)
- Requires Codecov token (optional setup)

**No tests found**:
- The workflow will skip gracefully if no tests exist
- Add tests gradually as you implement features

