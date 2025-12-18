# CI/CD Artifacts Guide

This guide explains what artifacts are saved by CI/CD workflows and how to use them for debugging.

## Overview

All CI/CD workflows save detailed logs and test outputs as artifacts, even when tests fail. These artifacts are available for download from the GitHub Actions workflow run page.

## Artifacts by Workflow

### Integration Tests Workflow

**Artifacts Saved:**

1. **`server-integration-test-logs`**
   - Full test output from server integration tests
   - Includes all test results, errors, and stack traces
   - Location: `backend/server/app/test-output.log`
   - Retention: 30 days

2. **`server-integration-coverage`**
   - Coverage report file
   - Location: `backend/server/app/coverage.out`
   - Retention: 30 days

3. **`email-worker-integration-test-logs`**
   - Full test output from email worker tests
   - Location: `backend/email-worker/test-output.log`
   - Retention: 30 days

4. **`translation-worker-integration-test-logs`**
   - Full test output from translation worker tests
   - Location: `backend/translation-worker/test-output.log`
   - Retention: 30 days

5. **`whatsapp-worker-integration-test-logs`**
   - Full test output from WhatsApp worker tests
   - Location: `backend/whatsapp-worker/test-output.log`
   - Retention: 30 days

6. **`integration-test-failure-report`**
   - Consolidated failure report with all test logs
   - Created when any test suite fails
   - Retention: 30 days

**How to Access:**
1. Go to the workflow run page
2. Scroll to the "Artifacts" section at the bottom
3. Download the relevant artifact ZIP file
4. Extract and view the log files

### Performance Tests Workflow

**Artifacts Saved:**

1. **`performance-test-results`**
   - Performance test output for all workers
   - Includes:
     - `email-worker-performance.txt`
     - `translation-worker-performance.txt`
     - `whatsapp-worker-performance.txt`
     - `test-output.log` files
   - Retention: 30 days

2. **`performance-regression-check-logs`**
   - Regression detection script output
   - Location: `regression-check.log`
   - Retention: 30 days

**What's Included:**
- Throughput measurements
- Latency metrics (average, P95, P99)
- Benchmark results
- Error messages if tests fail
- Regression detection results

### Performance Regression Detection Workflow

**Artifacts Saved:**

1. **`performance-results-{run_number}`**
   - Performance test results for the PR
   - Includes:
     - `email-worker-performance.txt`
     - `translation-worker-performance.txt`
     - `whatsapp-worker-performance.txt`
     - `test-output.log` files
     - `regression-check.log`
   - Retention: 7 days

**What's Included:**
- Current performance metrics
- Baseline comparison results
- Regression detection output
- Detailed test logs

## Using Artifacts for Debugging

### When Tests Fail

1. **Download the relevant artifact:**
   - For integration test failures: Download `{service}-integration-test-logs`
   - For performance test failures: Download `performance-test-results`
   - For regression detection: Download `performance-results-{run_number}`

2. **Extract the ZIP file:**
   ```bash
   unzip artifact-name.zip
   ```

3. **Examine the log files:**
   - Look for error messages
   - Check stack traces
   - Review test output for clues

4. **Share with team:**
   - Upload log files to issue tracker
   - Share in team chat
   - Include in bug reports

### Common Issues and Where to Look

#### Integration Test Failures

**Check:**
- `test-output.log` in the relevant service artifact
- Look for:
  - Database connection errors
  - RabbitMQ connection issues
  - Test timeout errors
  - Assertion failures

**Example:**
```bash
# Download artifact
unzip server-integration-test-logs.zip

# View log
cat backend/server/app/test-output.log | grep -i error
```

#### Performance Test Failures

**Check:**
- `performance-test-results` artifact
- Look for:
  - Benchmark failures
  - Timeout errors
  - Resource exhaustion
  - Regression warnings

**Example:**
```bash
# Download artifact
unzip performance-test-results.zip

# Check for regressions
grep -i "regression\|failed\|error" regression-check.log
```

#### Service Connection Issues

**Check:**
- Test logs for connection errors
- Look for:
  - PostgreSQL connection failures
  - Redis connection issues
  - RabbitMQ connection problems

**Example:**
```bash
# Search for connection errors
grep -i "connection\|dial\|timeout" test-output.log
```

## Artifact Structure

### Integration Test Artifacts

```
artifact-name.zip
├── backend/
│   ├── server/
│   │   └── app/
│   │       └── test-output.log
│   ├── email-worker/
│   │   └── test-output.log
│   ├── translation-worker/
│   │   └── test-output.log
│   └── whatsapp-worker/
│       └── test-output.log
```

### Performance Test Artifacts

```
performance-test-results.zip
├── backend/
│   ├── email-worker/
│   │   ├── email-worker-performance.txt
│   │   └── test-output.log
│   ├── translation-worker/
│   │   ├── translation-worker-performance.txt
│   │   └── test-output.log
│   └── whatsapp-worker/
│       ├── whatsapp-worker-performance.txt
│       └── test-output.log
└── regression-check.log
```

## Best Practices

1. **Always Check Artifacts on Failure:**
   - Don't rely solely on workflow summary
   - Download and examine logs for details

2. **Share Artifacts When Asking for Help:**
   - Include relevant log files
   - Provide context (commit SHA, workflow run URL)
   - Mention what you've already checked

3. **Keep Artifacts Organized:**
   - Name issues/tickets with artifact names
   - Reference workflow run numbers
   - Link to specific artifacts in discussions

4. **Monitor Artifact Retention:**
   - Artifacts are automatically deleted after retention period
   - Download important artifacts before they expire
   - Archive critical failure logs

## Troubleshooting

### Artifact Not Found

**Possible Causes:**
- Workflow didn't complete
- Artifact upload step failed
- Retention period expired

**Solutions:**
1. Check workflow run status
2. Verify artifact upload step completed
3. Check if retention period has passed

### Artifact Empty or Incomplete

**Possible Causes:**
- Test didn't produce output
- Output redirection failed
- File permissions issue

**Solutions:**
1. Check workflow logs for errors
2. Verify test actually ran
3. Check file paths in artifact

### Can't Download Artifact

**Possible Causes:**
- GitHub Actions permissions
- Artifact too large
- Network issues

**Solutions:**
1. Check repository permissions
2. Try downloading from different network
3. Use GitHub CLI: `gh run download`

## Automation

### Download Artifacts via GitHub CLI

```bash
# List workflow runs
gh run list

# Download artifacts from latest run
gh run download

# Download specific artifact
gh run download <run-id> -n artifact-name
```

### Download Artifacts via API

```bash
# Get artifact download URL
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/owner/repo/actions/artifacts

# Download artifact
curl -L -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3.raw" \
  https://api.github.com/repos/owner/repo/actions/artifacts/{artifact_id}/zip \
  -o artifact.zip
```

## Important: GitHub Actions Logs Are Always Available

**Even if artifacts fail to upload or workflow fails completely, GitHub Actions automatically saves complete workflow logs for every run.**

These logs include:
- All step outputs (stdout/stderr)
- Service startup logs
- Build errors
- Test failures
- Environment details

**How to Access:**
1. Go to Actions → Failed workflow run
2. Click on failed job
3. Click on failed step
4. View or download logs

**These logs are ALWAYS available**, even if:
- ❌ Docker builds fail
- ❌ Services don't start
- ❌ Tests never run
- ❌ Workflow fails completely

See `.github/FAILURE_DEBUGGING.md` for complete failure scenario guide.

## Summary

- ✅ All workflows save detailed logs as artifacts
- ✅ Artifacts are available even when tests fail
- ✅ **GitHub Actions logs are ALWAYS available** (automatic)
- ✅ 30-day retention for integration tests
- ✅ 7-day retention for PR performance tests
- ✅ Easy to download and share for debugging
- ✅ Comprehensive error information included

Use artifacts to debug failures, share with team members, and get help fixing issues!
