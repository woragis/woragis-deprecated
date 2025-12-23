# Performance Regression Detection Guide

This guide explains how performance regression detection works and how to use it.

## Overview

Performance regression detection automatically compares current performance test results against baseline metrics to identify performance degradations.

## How It Works

### 1. Baseline Storage

Performance baselines are stored in `.github/performance-baselines.json`:

```json
{
  "email-worker": {
    "throughput": 50.0,
    "latency_avg": 100.0,
    "latency_p95": 200.0
  }
}
```

### 2. Test Execution

Performance tests run and generate results with metrics:
- Throughput (messages/jobs per second)
- Average latency (milliseconds)
- P95 latency (milliseconds)

### 3. Metric Extraction

The regression detection script (`check_performance_regression.py`) extracts metrics from test logs:
- Parses "Throughput: X.XX msg/s" patterns
- Parses "Avg latency: X.XXms" patterns
- Parses benchmark results

### 4. Comparison

Current metrics are compared against baselines:
- **Throughput**: Degradation > 20% triggers regression
- **Latency**: Increase > 30% triggers regression

### 5. Reporting

If regressions detected:
- Workflow fails
- PR comments added (if PR workflow)
- Summary written to GitHub Actions
- Artifacts uploaded

## Regression Thresholds

### Throughput Threshold: 20%

If current throughput is 20% lower than baseline:
```
Baseline: 50 msg/s
Current:  40 msg/s (20% degradation)
Result:   ⚠️ Regression detected
```

### Latency Threshold: 30%

If current latency is 30% higher than baseline:
```
Baseline: 100ms
Current:  130ms (30% increase)
Result:   ⚠️ Regression detected
```

## Updating Baselines

### Method 1: GitHub Actions Workflow

1. Go to Actions → Update Performance Baselines
2. Click "Run workflow"
3. Type "yes" in confirmation field
4. Workflow runs tests and updates baselines
5. Baselines are committed automatically

### Method 2: Manual Update

1. Run performance tests locally:
   ```bash
   cd backend/server/scripts
   ./run-worker-performance-tests.sh
   ```

2. Store baselines:
   ```bash
   python3 .github/scripts/store_performance_baseline.py
   ```

3. Review and commit:
   ```bash
   git add .github/performance-baselines.json
   git commit -m "chore: update performance baselines"
   git push
   ```

## Interpreting Results

### No Regression
```
✅ Throughput OK: 52.34 ops/s (baseline: 50.00 ops/s)
✅ Average latency OK: 95.23ms (baseline: 100.00ms)
✅ No performance regressions detected!
```

### Regression Detected
```
⚠️ Throughput regression: 35.67 ops/s (baseline: 50.00 ops/s, degradation: 28.7%)
⚠️ Average latency regression: 145.23ms (baseline: 100.00ms, increase: 45.2%)
```

## Handling Regressions

### 1. Investigate

Check what changed:
- Code changes in the PR
- Dependencies updated
- Infrastructure changes
- Test environment differences

### 2. Fix or Accept

**Fix Performance:**
- Optimize code
- Fix bottlenecks
- Improve algorithms

**Accept Regression:**
- If regression is acceptable (e.g., added features)
- Update baselines to reflect new performance
- Document the reason

### 3. Update Baselines

If regression is acceptable:
1. Run Update Baselines workflow
2. Or manually update `.github/performance-baselines.json`
3. Commit the changes

## Best Practices

1. **Regular Updates**: Update baselines when performance improves
2. **Document Changes**: Document why baselines changed
3. **Review Trends**: Monitor performance over time
4. **Fix Promptly**: Address regressions quickly
5. **Set Realistic Baselines**: Baselines should reflect actual good performance

## Troubleshooting

### False Positives

**Cause**: CI environment differences, resource contention

**Solution**: 
- Review baseline thresholds
- Check CI runner resources
- Update baselines if acceptable

### Missing Metrics

**Cause**: Test output format changed, parsing issues

**Solution**:
- Check test output format
- Update parsing regex in script
- Verify test logs

### Baselines Not Found

**Cause**: Baseline file missing or corrupted

**Solution**:
- Check `.github/performance-baselines.json` exists
- Use default baselines if file missing
- Re-run baseline update workflow

## Advanced Configuration

### Custom Thresholds

Edit `.github/scripts/check_performance_regression.py`:

```python
THROUGHPUT_THRESHOLD = 0.20  # 20% degradation
LATENCY_THRESHOLD = 0.30     # 30% increase
```

### Custom Baselines

Edit `.github/performance-baselines.json`:

```json
{
  "email-worker": {
    "throughput": 60.0,
    "latency_avg": 80.0,
    "latency_p95": 150.0
  }
}
```

## Integration with CI/CD

### Automatic Detection

- Runs on PRs with performance test changes
- Comments on PRs if regressions detected
- Fails workflow if thresholds exceeded

### Daily Monitoring

- Performance tests run daily at 2 AM UTC
- Results stored as artifacts
- Trends can be tracked over time

## Example Workflow

1. **Developer makes changes** → Opens PR
2. **CI runs integration tests** → All pass ✅
3. **CI runs performance regression check** → Regression detected ⚠️
4. **PR comment added** → "Performance regression detected"
5. **Developer investigates** → Finds bottleneck
6. **Developer fixes** → Optimizes code
7. **CI re-runs** → No regression ✅
8. **PR merged** → Performance maintained

## Metrics Tracked

### Email Worker
- Throughput: Messages per second
- Latency: Average processing time
- P95 Latency: 95th percentile latency

### Translation Worker
- Throughput: Jobs per second
- Latency: Average processing time
- P95 Latency: 95th percentile latency
- Database Performance: Operations per second

### WhatsApp Worker
- Throughput: Messages per second
- Latency: Average processing time
- P95 Latency: 95th percentile latency

## Next Steps

1. Run initial baseline collection
2. Monitor performance trends
3. Set up alerts for regressions
4. Create performance dashboards
5. Document performance improvements
