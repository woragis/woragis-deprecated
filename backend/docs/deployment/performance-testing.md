# Performance Testing Guide

**Last Updated:** 2025-12-23  
**Purpose:** Guide for running performance tests and establishing baselines

---

## Overview

Performance testing ensures services can handle expected load and meet performance requirements. This guide covers running performance tests, interpreting results, and setting up performance regression detection.

---

## Performance Test Coverage

### Services with Performance Tests

| Service | Test Types | Status |
|---------|-----------|--------|
| **Server** | Load, Latency, Concurrent Users | ✅ Available (uses `performance_test` build tag) |
| **Email Worker** | Load, Latency, Rate Limiting | ✅ Complete |
| **Translation Worker** | Load, Latency, Multi-language, Database | ✅ Complete |
| **WhatsApp Worker** | Load, Latency, Rate Limiting | ✅ Complete |
| **Job Application Worker** | Load, Concurrent, Latency | ✅ Available (needs DATABASE_URL) |
| **Resume Worker** | Throughput, Concurrent | ✅ Available |
| **AI Service** | Health, Concurrent, Load | ✅ Available (AI endpoint test skipped - needs API keys) |
| **Creative Service** | Health, Concurrent | ✅ Available |
| **Docs Service** | Health, Concurrent | ✅ Available |

---

## Running Performance Tests

### Prerequisites

1. **Test Services Running**
   ```bash
   cd backend
   docker-compose -f docker-compose.test.yml up -d
   ```

2. **Environment Variables**
   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable"
   export REDIS_URL="redis://localhost:6380/0"
   export RABBITMQ_URL="amqp://test:test@localhost:5673/test"
   ```

### Go Services

#### Email Worker

**Load Test:**
```bash
cd backend/email-worker
go test -tags=integration -v ./internal/integration/... -run TestEmailWorkerLoadTest
```

**Latency Test:**
```bash
go test -tags=integration -v ./internal/integration/... -run TestEmailWorkerLatency
```

**Rate Limiting Test:**
```bash
go test -tags=integration -v ./internal/integration/... -run TestEmailWorkerRateLimiting
```

**Benchmark:**
```bash
go test -tags=integration -bench=BenchmarkEmailWorkerThroughput -benchmem ./internal/integration/...
```

#### WhatsApp Worker

**Load Test:**
```bash
cd backend/whatsapp-worker
go test -tags=integration -v ./internal/integration/... -run TestWhatsAppWorkerLoadTest
```

**Latency Test:**
```bash
go test -tags=integration -v ./internal/integration/... -run TestWhatsAppWorkerLatency
```

**Rate Limiting Test:**
```bash
go test -tags=integration -v ./internal/integration/... -run TestWhatsAppWorkerRateLimiting
```

#### Translation Worker

**Load Test:**
```bash
cd backend/translation-worker
go test -tags=integration -v ./internal/integration/... -run TestTranslationWorkerLoadTest
```

**Latency Test:**
```bash
go test -tags=integration -v ./internal/integration/... -run TestTranslationWorkerLatency
```

#### Server

**Note:** Server performance tests use a separate build tag (`performance_test`) to avoid conflicts with regular integration tests.

**Load Test:**
```bash
cd backend/server/app
go test -tags=integration,performance_test -v ./internal/integration/... -run TestServerAPILoadTest
```

**Latency Test:**
```bash
go test -tags=integration,performance_test -v ./internal/integration/... -run TestServerAPILatency
```

**Concurrent Users Test:**
```bash
go test -tags=integration,performance_test -v ./internal/integration/... -run TestServerConcurrentUsers
```

**Benchmark:**
```bash
go test -tags=integration,performance_test -bench=BenchmarkServerHealthEndpoint -benchmem ./internal/integration/...
```

#### Resume Worker (Python)

**Throughput Test:**
```bash
cd backend/resume-worker
export RABBITMQ_URL="amqp://test:test@localhost:5673/test"
pytest tests/performance_test.py::test_message_throughput -v
```

**Concurrent Processing Test:**
```bash
pytest tests/performance_test.py::test_concurrent_message_processing -v
```

#### AI Service (Python)

**Health Endpoint Performance:**
```bash
cd backend/ai-service
pytest app/tests/performance_test.py::test_health_endpoint_performance -v
```

**Concurrent Requests:**
```bash
pytest app/tests/performance_test.py::test_concurrent_requests -v
```

#### Creative Service (Python)

**Health Endpoint Performance:**
```bash
cd backend/creative-service
pytest app/tests/performance_test.py::test_health_endpoint_performance -v
```

**Concurrent Requests:**
```bash
pytest app/tests/performance_test.py::test_concurrent_requests -v
```

#### Docs Service (Python)

**Health Endpoint Performance:**
```bash
cd backend/docs-service
pytest tests/performance_test.py::test_health_endpoint_performance -v
```

**Concurrent Requests:**
```bash
pytest tests/performance_test.py::test_concurrent_requests -v
```

---

## Performance Metrics

### Key Metrics

1. **Throughput (RPS)**
   - Requests per second
   - Messages per second (for workers)
   - Target: Varies by service

2. **Latency**
   - **P50 (Median):** 50% of requests
   - **P95:** 95% of requests
   - **P99:** 99% of requests
   - Target: < 200ms for P95

3. **Error Rate**
   - Percentage of failed requests
   - Target: < 1%

4. **Concurrent Users**
   - Number of simultaneous users
   - Target: Varies by service

### Performance Baselines

#### Server API

**Load Test (500 requests, 50 concurrent):**
- Throughput: > 100 req/s
- Success Rate: > 90%
- Average Latency: < 100ms
- P95 Latency: < 200ms

**Latency Test (100 iterations):**
- Average Latency: < 100ms
- P95 Latency: < 200ms

**Concurrent Users (20 users, 5 requests each):**
- Success Rate: > 80%
- Throughput: > 50 req/s

#### Email Worker

**Load Test:**
- Throughput: ~20 msg/s
- Success Rate: > 90%

**Latency Test:**
- Average Latency: < 6ms
- Max Latency: < 7ms

**Rate Limiting:**
- Processing Rate: ~20 msg/s
- Respects rate limits

#### WhatsApp Worker

**Load Test:**
- Throughput: ~20 msg/s
- Success Rate: > 90%

**Latency Test:**
- Average Latency: < 10ms
- Max Latency: < 15ms

**Rate Limiting:**
- Processing Rate: ~20 msg/s
- Respects rate limits

---

## Performance Test Types

### 1. Load Tests

**Purpose:** Test system under expected load

**Configuration:**
- Concurrent requests: 50-100
- Total requests: 500-1000
- Duration: 30-60 seconds

**Metrics:**
- Throughput (req/s)
- Success rate
- Error rate
- Response times

### 2. Latency Tests

**Purpose:** Measure response time distribution

**Configuration:**
- Iterations: 100-1000
- Single-threaded
- Measure p50, p95, p99

**Metrics:**
- Average latency
- P50 latency
- P95 latency
- P99 latency
- Max latency

### 3. Stress Tests

**Purpose:** Find system limits

**Configuration:**
- Gradually increase load
- Find breaking point
- Monitor resource usage

**Metrics:**
- Maximum throughput
- Breaking point
- Resource usage at limit

### 4. Endurance Tests

**Purpose:** Test system stability over time

**Configuration:**
- Moderate load
- Long duration (hours)
- Monitor for memory leaks

**Metrics:**
- Memory usage over time
- CPU usage over time
- Error rate over time

---

## Performance Regression Detection

### CI/CD Integration

**Workflow:** `.github/workflows/performance-tests.yml`

**Schedule:**
- Daily at 2 AM UTC
- On push to main (if performance files changed)
- Manual trigger

**Process:**
1. Run performance tests
2. Compare with baseline
3. Detect regressions (> 20% degradation)
4. Alert on regression
5. Update baseline if improvement

### Baseline Management

**Baseline Storage:**
- Stored in `.github/baselines/performance-baselines.json`
- Updated on successful performance test runs
- Version controlled

**Baseline Format:**
```json
{
  "server": {
    "load_test": {
      "throughput": 120.5,
      "p95_latency_ms": 180,
      "success_rate": 0.95
    },
    "latency_test": {
      "avg_latency_ms": 85,
      "p95_latency_ms": 175
    }
  },
  "email_worker": {
    "load_test": {
      "throughput": 19.58,
      "success_rate": 1.0
    }
  }
}
```

### Regression Detection

**Thresholds:**
- **Throughput:** > 20% decrease = regression
- **Latency:** > 20% increase = regression
- **Error Rate:** > 5% = regression

**Actions on Regression:**
1. Alert team
2. Block merge (optional)
3. Create issue
4. Investigate cause

---

## Load Testing Tools

### Built-in Tests

**Go Services:**
- Use Go's testing package
- `go test -bench` for benchmarks
- Custom load test functions

**Python Services:**
- Use pytest with custom fixtures
- `pytest-benchmark` for benchmarks
- Custom load test functions

**Node.js Services:**
- Use Jest for testing
- Custom load test functions

### External Tools

**k6:**
```bash
# Install
brew install k6  # macOS
# or download from k6.io

# Run test
k6 run load-test.js
```

**Apache Bench (ab):**
```bash
# Install
# macOS: brew install httpd
# Linux: apt-get install apache2-utils

# Run test
ab -n 1000 -c 50 https://api.woragis.com/healthz
```

**hey:**
```bash
# Install
go install github.com/rakyll/hey@latest

# Run test
hey -n 1000 -c 50 https://api.woragis.com/healthz
```

---

## Performance Test Scripts

### Run All Performance Tests

```bash
#!/bin/bash
# scripts/run-performance-tests.sh

echo "Running performance tests..."

# Go workers
cd backend/email-worker && go test -tags=integration -run "Test.*Load|Test.*Latency|Test.*Rate" ./internal/integration/...
cd ../whatsapp-worker && go test -tags=integration -run "Test.*Load|Test.*Latency|Test.*Rate" ./internal/integration/...
cd ../translation-worker && go test -tags=integration -run "Test.*Load|Test.*Latency" ./internal/integration/...

# Server
cd ../server/app && go test -tags=integration,performance_test -run "Test.*Load|Test.*Latency" ./internal/integration/...

echo "Performance tests completed"
```

### Generate Performance Report

```bash
#!/bin/bash
# scripts/generate-performance-report.sh

echo "Generating performance report..."

# Run tests and capture output
cd backend/email-worker
go test -tags=integration -run "Test.*Load|Test.*Latency" ./internal/integration/... > /tmp/email-worker-perf.txt 2>&1

# Parse results and generate report
# (Implementation depends on output format)

echo "Performance report generated: /tmp/performance-report.json"
```

---

## Performance Characteristics

### Server API

**Expected Performance:**
- **Throughput:** 100-200 req/s (single instance)
- **P50 Latency:** < 50ms
- **P95 Latency:** < 200ms
- **P99 Latency:** < 500ms
- **Concurrent Users:** 50-100

**Bottlenecks:**
- Database queries
- External API calls
- Authentication/authorization

### Email Worker

**Expected Performance:**
- **Throughput:** 15-25 msg/s
- **Average Latency:** < 10ms
- **P95 Latency:** < 20ms

**Bottlenecks:**
- SMTP server response time
- RabbitMQ message processing
- Rate limiting

### Translation Worker

**Expected Performance:**
- **Throughput:** 10-20 translations/s
- **Average Latency:** < 100ms
- **P95 Latency:** < 200ms

**Bottlenecks:**
- Translation API response time
- Database operations
- Message processing

---

## Performance Optimization

### Identified Bottlenecks

1. **Database Queries**
   - Add indexes
   - Optimize queries
   - Use connection pooling
   - Implement caching

2. **External API Calls**
   - Add caching
   - Use async processing
   - Implement retries with backoff
   - Batch requests

3. **Authentication/Authorization**
   - Cache JWT validation
   - Optimize token verification
   - Use Redis for session storage

4. **Message Processing**
   - Increase worker concurrency
   - Optimize message parsing
   - Batch processing

### Optimization Checklist

- [ ] Database indexes added
- [ ] Queries optimized
- [ ] Caching implemented
- [ ] Connection pooling configured
- [ ] Async processing used
- [ ] Rate limiting optimized
- [ ] Resource limits adjusted

---

## Performance Monitoring

### Metrics to Monitor

1. **Request Metrics**
   - Request rate
   - Response time (p50, p95, p99)
   - Error rate
   - Timeout rate

2. **Resource Metrics**
   - CPU usage
   - Memory usage
   - Network I/O
   - Disk I/O

3. **Queue Metrics**
   - Queue depth
   - Processing rate
   - Message age
   - Dead letter queue size

### Alerting Thresholds

- **High Error Rate:** > 5%
- **High Latency:** P95 > 500ms
- **Low Throughput:** < 50% of baseline
- **High Resource Usage:** CPU > 80%, Memory > 90%
- **Queue Backlog:** > 1000 messages

---

## Performance Test Results

### Latest Results (2025-12-23)

#### Email Worker

**Load Test:**
- Throughput: 19.58 msg/s
- Success Rate: 100%
- Duration: 5.42s

**Latency Test:**
- Average Latency: < 6ms
- Max Latency: 6.94ms

**Rate Limiting:**
- Processing Rate: 19.99 msg/s
- Respects rate limits

#### WhatsApp Worker

**Load Test:**
- Throughput: ~20 msg/s
- Success Rate: > 90%

**Latency Test:**
- Average Latency: < 10ms
- Max Latency: < 15ms

---

## Next Steps

1. **Complete Performance Tests**
   - Add performance tests for Python services
   - Add performance tests for Node.js worker
   - Fix server performance test build issues

2. **Establish Baselines**
   - Run comprehensive performance tests
   - Document baselines
   - Set up regression detection

3. **Optimize Performance**
   - Identify bottlenecks
   - Implement optimizations
   - Re-test and update baselines

4. **Set Up CI/CD**
   - Configure performance regression tests
   - Set up baseline management
   - Configure alerts

---

## Related Documentation

- **Performance Tests Workflow:** `.github/workflows/performance-tests.yml`
- **Performance Regression:** `.github/workflows/performance-regression.yml`
- **Testing Patterns:** `docs/development/testing-patterns.md`
- **Integration Tests:** `docs/development/integration-tests.md`

---

**Last Updated:** 2025-12-23

