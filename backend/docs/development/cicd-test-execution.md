# CI/CD Test Execution Documentation

**Last Updated:** 2025-12-23  
**Purpose:** Guide for understanding and using CI/CD test workflows

---

## Overview

This document explains how tests are executed in CI/CD pipelines, what workflows exist, and how to interpret results.

---

## Test Workflows

### 1. Integration Tests Workflow

**File:** `.github/workflows/integration-tests.yml`

**Purpose:** Run integration tests for all services

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Manual workflow dispatch

**Services Tested:**
- Server (Go)
- Email Worker (Go)
- Translation Worker (Go)
- WhatsApp Worker (Go)
- Job Application Worker (Node.js)
- Resume Worker (Python)
- AI Service (Python)
- Creative Service (Python)
- Docs Service (Python)

**Test Services:**
- PostgreSQL (port 5433)
- Redis (port 6380)
- RabbitMQ (port 5673)

**Coverage:**
- All services upload coverage to Codecov
- Coverage reports in PR comments

---

### 2. Build All Services Workflow

**File:** `.github/workflows/build-all.yml`

**Purpose:** Build and test all services before deployment

**Triggers:**
- Push to tags matching `v*` (e.g., `v1.0.0`)

**Steps:**
1. Extract tag name
2. Run tests for all services
3. Build Docker images
4. Push to Docker Hub
5. Deploy to Railway (optional)

---

### 3. Performance Tests Workflow

**File:** `.github/workflows/performance-tests.yml`

**Purpose:** Run performance tests and detect regressions

**Triggers:**
- Daily at 2 AM UTC (scheduled)
- Manual workflow dispatch
- Push to `main` (if performance files changed)

**Tests:**
- Load tests
- Latency tests
- Rate limiting tests
- Benchmark tests

---

## Workflow Execution

### Integration Tests Workflow

#### Job Structure

**Go Services:**
- `server-integration-tests`
- `email-worker-tests`
- `translation-worker-tests`
- `whatsapp-worker-tests`

**Node.js Services:**
- `job-application-worker-tests`

**Python Services:**
- `resume-worker-tests`
- `ai-service-tests`
- `creative-service-tests`
- `docs-service-tests`

**Summary:**
- `test-summary` - Aggregates all test results

#### Execution Steps

1. **Checkout Code**
   ```yaml
   - uses: actions/checkout@v4
   ```

2. **Set Up Language**
   ```yaml
   # Go
   - uses: actions/setup-go@v5
   
   # Python
   - uses: actions/setup-python@v5
   
   # Node.js
   - uses: actions/setup-node@v4
   ```

3. **Start Test Services**
   ```yaml
   services:
     postgres:
       image: postgres:15-alpine
     redis:
       image: redis:7-alpine
     rabbitmq:
       image: rabbitmq:3.13-management-alpine
   ```

4. **Install Dependencies**
   ```yaml
   # Go
   - run: go mod download
   
   # Python
   - run: pip install -r requirements.txt
   
   # Node.js
   - run: npm ci
   ```

5. **Run Tests**
   ```yaml
   # Go
   - run: go test -tags=integration -v ./...
   
   # Python
   - run: pytest tests/integration/ -v
   
   # Node.js
   - run: npm run test:integration
   ```

6. **Upload Coverage**
   ```yaml
   - uses: codecov/codecov-action@v4
     with:
       files: ./coverage.out
   ```

7. **Upload Artifacts**
   ```yaml
   - uses: actions/upload-artifact@v4
     with:
       name: test-results
       path: test-results/
   ```

---

## Viewing Test Results

### GitHub Actions

1. **Go to Actions Tab**
   - Repository → Actions

2. **Select Workflow**
   - Click on "Integration Tests" workflow

3. **View Run**
   - Click on specific run
   - View job status
   - Expand job to see steps

4. **View Logs**
   - Click on job
   - Expand step to see logs
   - Download logs if needed

### Test Summary

The `test-summary` job creates a summary of all test results:

```markdown
## Test Results Summary

| Service | Status | Tests | Coverage |
|---------|--------|-------|----------|
| Server | ✅ Pass | 30 | 75% |
| Email Worker | ✅ Pass | 15 | 80% |
| ... | ... | ... | ... |
```

### Coverage Reports

**Codecov:**
- View in PR comments
- View on Codecov dashboard
- Compare coverage trends

**Artifacts:**
- Download coverage reports from artifacts
- View HTML coverage reports locally

---

## Test Artifacts

### Available Artifacts

1. **Test Results**
   - Test output logs
   - Test reports (JUnit XML, etc.)
   - Coverage reports

2. **Performance Results**
   - Performance test output
   - Benchmark results
   - Regression check logs

3. **Build Artifacts**
   - Docker images (pushed to Docker Hub)
   - Build logs

### Downloading Artifacts

1. Go to workflow run
2. Scroll to "Artifacts" section
3. Click artifact name
4. Download ZIP file
5. Extract and view files

---

## Interpreting Results

### Test Status

**✅ Success:**
- All tests passed
- Coverage meets threshold
- No errors

**❌ Failure:**
- One or more tests failed
- Check logs for details
- Fix issues and re-run

**⚠️ Flaky:**
- Tests pass/fail inconsistently
- Review test code
- Add retries if appropriate

### Coverage Status

**Good:**
- Coverage > 70%
- Coverage increased or maintained
- All critical paths covered

**Needs Improvement:**
- Coverage < 70%
- Coverage decreased
- Missing critical path coverage

---

## Troubleshooting CI/CD Tests

### Tests Fail in CI but Pass Locally

**Possible Causes:**
1. Environment differences
2. Timing issues
3. Resource constraints
4. Service availability

**Solutions:**
1. Add retries for flaky tests
2. Increase timeouts
3. Check service health checks
4. Review test isolation

### Slow Tests

**Optimize:**
1. Run tests in parallel
2. Use test matrix for parallelization
3. Cache dependencies
4. Optimize test code

### Coverage Decreased

**Check:**
1. New code not covered
2. Tests removed
3. Coverage calculation changed

**Fix:**
1. Add tests for new code
2. Review removed tests
3. Verify coverage configuration

---

## Best Practices

### 1. Fast Feedback
- Keep test execution time < 30 minutes
- Run critical tests first
- Use parallel execution

### 2. Reliable Tests
- Avoid flaky tests
- Use proper test isolation
- Clean up after tests

### 3. Good Coverage
- Maintain > 70% coverage
- Cover critical paths
- Test error scenarios

### 4. Clear Results
- Use descriptive test names
- Provide clear error messages
- Include context in logs

### 5. Regular Maintenance
- Review failing tests weekly
- Update test data monthly
- Refactor slow tests quarterly

---

## Workflow Configuration

### Environment Variables

**Test Services:**
```yaml
env:
  POSTGRES_VERSION: '15-alpine'
  REDIS_VERSION: '7-alpine'
  RABBITMQ_VERSION: '3.13-management-alpine'
  GO_VERSION: '1.21'
```

**Service URLs:**
```yaml
env:
  DATABASE_URL: postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable
  REDIS_URL: redis://localhost:6380/0
  RABBITMQ_URL: amqp://test:test@localhost:5673/test
```

### Secrets

**Required:**
- `DOCKER_HUB_USERNAME` - Docker Hub username
- `DOCKER_HUB_TOKEN` - Docker Hub access token

**Optional:**
- `OPENAI_API_KEY` - For AI service tests
- `ANTHROPIC_API_KEY` - For AI service tests
- `RAILWAY_TOKEN` - For Railway deployment
- `RAILWAY_PROJECT_ID` - For Railway deployment

---

## Manual Test Execution

### Trigger Workflow Manually

1. Go to GitHub Actions
2. Select workflow
3. Click "Run workflow"
4. Select branch
5. Click "Run"

### Run Specific Tests

**In Workflow:**
- Modify workflow to run specific tests
- Use `-run` flag for Go tests
- Use `-k` flag for pytest
- Use `--testNamePattern` for Jest

**Locally:**
```bash
# Go
go test -tags=integration -run TestSpecificTest ./...

# Python
pytest tests/integration/test_specific.py

# Node.js
npm test -- --testNamePattern="Specific Test"
```

---

## Related Documentation

- **Integration Tests:** `docs/development/integration-tests.md`
- **Testing Patterns:** `docs/development/testing-patterns.md`
- **Build Workflow:** `docs/deployment/build-workflow.md`
- **GitHub Actions README:** `.github/workflows/README.md`

---

**Last Updated:** 2025-12-23

