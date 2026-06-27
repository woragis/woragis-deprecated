# CI/CD Integration Guide

This document describes the CI/CD integration for integration and performance tests.

## Overview

GitHub Actions workflows have been set up to:
- Run integration tests on every push and pull request
- Run performance tests daily and on demand
- Detect performance regressions automatically
- Upload test results as artifacts
- Generate test summaries

## Workflows

### 1. Integration Tests (`integration-tests.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- Manual dispatch

**Jobs:**
- `server-integration-tests` - Server API integration tests
- `email-worker-tests` - Email worker integration tests
- `translation-worker-tests` - Translation worker integration tests
- `whatsapp-worker-tests` - WhatsApp worker integration tests
- `test-summary` - Aggregates test results

**Features:**
- Runs all integration tests in parallel
- Uploads coverage reports to Codecov
- Generates test summary in GitHub Actions

### 2. Performance Tests (`performance-tests.yml`)

**Triggers:**
- Daily at 2 AM UTC (scheduled)
- Push to `main` branch (when performance tests change)
- Manual dispatch

**Jobs:**
- `performance-tests` - Runs performance tests for all workers

**Features:**
- Runs load, latency, and benchmark tests
- Uploads performance test results as artifacts
- Checks for performance regressions
- Generates performance test summary

### 3. Performance Regression Detection (`performance-regression.yml`)

**Triggers:**
- Pull requests (when performance tests change)
- Manual dispatch

**Jobs:**
- `detect-regression` - Runs performance tests and detects regressions

**Features:**
- Compares current results with baselines
- Comments on PRs if regressions detected
- Uploads test results as artifacts
- Fails workflow if regressions exceed thresholds

### 4. Update Baselines (`update-baselines.yml`)

**Triggers:**
- Manual dispatch only (requires confirmation)

**Jobs:**
- `update-baselines` - Updates performance baselines from test results

**Features:**
- Runs performance tests
- Extracts metrics from results
- Updates `.github/performance-baselines.json`
- Commits updated baselines

## Performance Regression Detection

### How It Works

1. **Baseline Storage**: Performance baselines are stored in `.github/performance-baselines.json`
2. **Test Execution**: Performance tests run and generate results
3. **Metric Extraction**: Script extracts throughput and latency metrics
4. **Comparison**: Current metrics compared against baselines
5. **Regression Detection**: Flags regressions exceeding thresholds:
   - Throughput degradation > 20%
   - Latency increase > 30%

### Baseline Format

```json
{
  "email-worker": {
    "throughput": 50.0,
    "latency_avg": 100.0,
    "latency_p95": 200.0
  }
}
```

### Updating Baselines

**Manual Update:**
1. Go to Actions → Update Performance Baselines
2. Click "Run workflow"
3. Type "yes" to confirm
4. Workflow will run tests and update baselines

**After Good Performance Run:**
```bash
# Run performance tests locally
cd backend/server/scripts
./run-worker-performance-tests.sh

# Update baselines
python3 .github/scripts/store_performance_baseline.py
git add .github/performance-baselines.json
git commit -m "chore: update performance baselines"
git push
```

## Test Artifacts

All workflows upload test results as artifacts:
- Integration test coverage reports
- Performance test result files
- Test logs and outputs

Artifacts are retained for:
- Integration tests: 30 days
- Performance tests: 30 days
- Regression detection: 7 days

## GitHub Actions Summary

Test summaries are automatically generated in GitHub Actions:
- Integration test status per suite
- Performance test results
- Regression detection results

## Environment Variables

Workflows use GitHub Actions services for dependencies:
- PostgreSQL: `localhost:5433`
- Redis: `localhost:6380`
- RabbitMQ: `localhost:5673`

No additional environment variables needed - services are automatically configured.

## Workflow Status Badges

Add to your README:

```markdown
![Integration Tests](https://github.com/your-org/woragis/workflows/Integration%20Tests/badge.svg)
![Performance Tests](https://github.com/your-org/woragis/workflows/Performance%20Tests/badge.svg)
```

## Troubleshooting

### Tests Fail in CI but Pass Locally

**Possible Causes:**
- Service startup timing
- Resource constraints
- Network issues

**Solutions:**
- Increase service wait times
- Check service health checks
- Review test timeouts

### Performance Regressions False Positives

**Possible Causes:**
- CI environment differences
- Resource contention
- Baseline too strict

**Solutions:**
- Review baseline thresholds
- Update baselines if acceptable
- Check CI runner resources

### Coverage Reports Not Uploading

**Possible Causes:**
- Codecov token not configured
- Coverage file not generated
- Path issues

**Solutions:**
- Configure Codecov token in repository settings
- Check coverage file path
- Review workflow logs

## Best Practices

1. **Review PR Comments**: Check performance regression comments on PRs
2. **Update Baselines**: Update baselines when performance improves
3. **Monitor Trends**: Review performance test artifacts regularly
4. **Fix Regressions**: Address performance regressions promptly
5. **Document Changes**: Document performance-impacting changes

## Next Steps

1. Configure Codecov token for coverage reports
2. Set up performance dashboards (Grafana)
3. Add Slack/email notifications for failures
4. Create performance test reports
5. Add performance budgets
