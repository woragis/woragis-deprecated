# Performance and Load Testing Guide

This document provides guidance on running and interpreting performance tests for backend workers.

## Overview

Performance tests have been implemented for all Go-based workers to measure:
- **Throughput**: Messages/jobs processed per second
- **Latency**: Time to process individual messages
- **Concurrency**: Behavior with multiple consumers
- **Load Handling**: Performance under high message rates
- **Database Performance**: Database operations under load (translation worker)

## Running Performance Tests

### Prerequisites

1. Docker and Docker Compose installed
2. Go 1.21+ installed
3. Test dependencies running (RabbitMQ, PostgreSQL)

### Quick Start

**Linux/Mac:**
```bash
cd backend/server/scripts
./run-worker-performance-tests.sh
```

**Windows:**
```cmd
cd backend\server\scripts
run-worker-performance-tests.bat
```

### Run Specific Worker Tests

**Linux/Mac:**
```bash
./run-worker-performance-tests.sh email-worker
./run-worker-performance-tests.sh translation-worker
./run-worker-performance-tests.sh whatsapp-worker
```

**Windows:**
```cmd
run-worker-performance-tests.bat email-worker
run-worker-performance-tests.bat translation-worker
run-worker-performance-tests.bat whatsapp-worker
```

### Manual Execution

1. Start test dependencies:
   ```bash
   cd backend/server
   docker-compose -f docker-compose.test.yml up -d
   ```

2. Run performance tests for a specific worker:
   ```bash
   cd backend/email-worker
   go test -tags=integration -run "Test.*Load|Test.*Latency|Benchmark" ./internal/integration/... -v
   ```

3. Run benchmarks:
   ```bash
   cd backend/email-worker
   go test -tags=integration -bench=Benchmark -benchmem ./internal/integration/...
   ```

## Test Types

### 1. Load Tests

Load tests measure worker performance under normal to high load conditions.

**Email/WhatsApp Workers:**
- Publishes 100 messages with 10 concurrent publishers
- Measures throughput and processing time
- Validates all messages are processed successfully

**Translation Worker:**
- Publishes 50 translation jobs with 5 concurrent publishers
- Tests database operations under load
- Validates translation record creation/updates

**Expected Results:**
- All messages/jobs processed successfully
- Processing completes within reasonable time (< 10 seconds)
- Throughput > 10 messages/jobs per second

### 2. Latency Tests

Latency tests measure the time to process individual messages.

**Test Details:**
- Publishes 50 messages sequentially
- Measures processing time for each message
- Calculates min, max, and average latency

**Expected Results:**
- Average latency < 100ms
- Consistent latency across messages

### 3. Concurrent Consumer Tests

Tests worker behavior with multiple concurrent consumers.

**Test Details:**
- Starts 3 concurrent consumers
- Publishes 30 messages
- Validates load balancing across consumers

**Expected Results:**
- All messages processed
- Even distribution across consumers

### 4. Rate Limiting Tests

Tests worker behavior under high message rates.

**Test Details:**
- Publishes 200 messages at high rate
- Measures publish rate vs processing rate
- Validates queue buffering and processing

**Expected Results:**
- All messages eventually processed
- Queue handles burst traffic gracefully

### 5. Benchmarks

Go benchmarks measure raw throughput.

**Running Benchmarks:**
```bash
go test -tags=integration -bench=Benchmark -benchmem ./internal/integration/...
```

**Output Format:**
```
BenchmarkEmailWorkerThroughput-8    1000    1234 ns/op    512 B/op    2 allocs/op
```

- `1000`: Number of iterations
- `1234 ns/op`: Nanoseconds per operation
- `512 B/op`: Bytes allocated per operation
- `2 allocs/op`: Number of allocations per operation

## Interpreting Results

### Throughput

**Good Performance:**
- Email/WhatsApp: > 50 messages/second
- Translation: > 10 jobs/second (includes database operations)

**Needs Optimization:**
- Email/WhatsApp: < 20 messages/second
- Translation: < 5 jobs/second

### Latency

**Good Performance:**
- Average latency < 100ms
- P95 latency < 200ms
- P99 latency < 500ms

**Needs Optimization:**
- Average latency > 500ms
- High variance in latency

### Resource Usage

Monitor during tests:
- CPU usage
- Memory usage
- RabbitMQ queue depth
- Database connection pool usage

## Performance Baselines

### Email Worker
- **Throughput**: 50-100 msg/s
- **Latency**: 50-100ms average
- **Concurrent Consumers**: 3+ consumers handle load well

### Translation Worker
- **Throughput**: 10-20 jobs/s (includes DB operations)
- **Latency**: 100-200ms average
- **Multi-language**: Handles 5+ languages concurrently

### WhatsApp Worker
- **Throughput**: 50-100 msg/s
- **Latency**: 50-100ms average
- **Concurrent Consumers**: 3+ consumers handle load well

## Troubleshooting

### Tests Timeout

**Possible Causes:**
- RabbitMQ not ready
- Database connection issues
- High system load

**Solutions:**
- Increase timeout values in tests
- Check service health: `docker-compose -f docker-compose.test.yml ps`
- Reduce concurrent load in tests

### Low Throughput

**Possible Causes:**
- Network latency
- Database bottlenecks
- Mock service overhead

**Solutions:**
- Run tests on faster hardware
- Optimize database queries
- Reduce mock service delays

### High Latency

**Possible Causes:**
- Queue depth too high
- Consumer processing slowly
- Network delays

**Solutions:**
- Increase consumer count
- Optimize message processing
- Check network connectivity

## Continuous Performance Monitoring

### CI/CD Integration

Add performance tests to CI/CD pipeline:

```yaml
# .github/workflows/performance-tests.yml
name: Performance Tests
on:
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM
  workflow_dispatch:

jobs:
  performance-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Performance Tests
        run: |
          cd backend/server/scripts
          ./run-worker-performance-tests.sh
```

### Metrics Collection

Consider adding Prometheus metrics to performance tests:
- Messages processed per second
- Average processing latency
- Error rates
- Queue depth

## Best Practices

1. **Run tests in isolation**: Use separate test queues/exchanges
2. **Warm up**: Allow consumers to start before measuring
3. **Multiple runs**: Run tests multiple times and average results
4. **Monitor resources**: Watch CPU, memory, and network during tests
5. **Document baselines**: Record expected performance metrics
6. **Compare over time**: Track performance trends

## Next Steps

1. Add Prometheus metrics collection to performance tests
2. Create performance dashboards in Grafana
3. Set up automated performance regression detection
4. Add performance tests to CI/CD pipeline
5. Create performance test reports
