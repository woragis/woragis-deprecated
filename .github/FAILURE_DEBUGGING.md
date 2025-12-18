# Failure Debugging Guide

This guide explains what happens when CI/CD fails and how to get logs for debugging.

## What Happens When Everything Fails?

### Scenario: Complete Workflow Failure

If you do `git push` and everything starts failing (Docker builds, services, tests, etc.), here's what you'll have available:

## ✅ Always Available: GitHub Actions Workflow Logs

**GitHub Actions automatically saves complete workflow logs** for every run, regardless of success or failure. These logs include:

- ✅ All step outputs (stdout/stderr)
- ✅ Service startup logs
- ✅ Build errors
- ✅ Test failures
- ✅ Environment details
- ✅ Service health checks

**How to Access:**
1. Go to your repository on GitHub
2. Click "Actions" tab
3. Click on the failed workflow run
4. Click on any failed job
5. Click on any failed step to see its logs
6. **Download all logs:** Click the "..." menu → "Download log"

**These logs are ALWAYS available**, even if:
- Docker builds fail
- Services don't start
- Tests never run
- Workflow fails completely

## ✅ Artifacts We Save

In addition to GitHub's automatic logs, our workflows save specific artifacts:

### Even If Tests Don't Run

1. **Setup Logs**
   - Go dependency installation
   - Module downloads
   - Build setup errors
   - Saved even if setup fails

2. **Service Health Logs**
   - PostgreSQL connection status
   - Redis connection status
   - RabbitMQ connection status
   - Service startup errors

3. **Build Logs**
   - Compilation errors
   - Dependency resolution issues
   - Build failures

### If Tests Run But Fail

4. **Test Output Logs**
   - Full test output
   - Error messages
   - Stack traces
   - Test failures

5. **Coverage Reports**
   - Code coverage data
   - Coverage files

## Complete Failure Scenario

### Example: Everything Fails

```
git push
  ↓
GitHub Actions starts workflow
  ↓
Step 1: Checkout code ✅
  ↓
Step 2: Set up Go ✅
  ↓
Step 3: Install dependencies ❌ FAILS
  ↓
Step 4: Wait for services ❌ (never runs)
  ↓
Step 5: Run tests ❌ (never runs)
```

**What You'll Have:**

1. ✅ **GitHub Actions Logs** (automatic)
   - Step 3 output showing dependency failure
   - Error messages
   - Stack traces

2. ✅ **Setup Logs Artifact** (if we got that far)
   - `setup.log` with dependency errors

3. ✅ **Workflow Summary**
   - Which steps failed
   - Error messages
   - Job status

## How to Get Logs for Debugging

### Method 1: GitHub Actions UI (Recommended)

1. **Go to Actions tab**
2. **Click failed workflow run**
3. **Click failed job**
4. **Click failed step**
5. **Copy logs** or **Download log**

**Pro Tip:** You can copy the entire log and paste it into a file or share it directly.

### Method 2: Download Artifacts

1. **Go to Actions tab**
2. **Click failed workflow run**
3. **Scroll to "Artifacts" section**
4. **Download artifact ZIP**
5. **Extract and view logs**

### Method 3: GitHub CLI

```bash
# List workflow runs
gh run list

# View logs for specific run
gh run view <run-id> --log

# Download artifacts
gh run download <run-id>
```

### Method 4: API

```bash
# Get workflow run logs
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/owner/repo/actions/runs/{run_id}/logs
```

## What to Share When Asking for Help

When everything fails and you need help, share:

1. **Workflow Run URL**
   ```
   https://github.com/owner/repo/actions/runs/123456789
   ```

2. **Failed Step Logs**
   - Copy the log output from the failed step
   - Include error messages
   - Include stack traces

3. **Artifacts** (if available)
   - Download and attach artifact files
   - Or share artifact names

4. **Context**
   - What you were trying to do
   - What changed recently
   - Commit SHA

### Example Request

```
Workflow failed: https://github.com/owner/repo/actions/runs/123456789

Failed at: "Install dependencies" step

Error:
[copy error log here]

Commit: abc123def
Branch: main
```

## Failure Points and What's Captured

### 1. Checkout Failure
- **GitHub Logs:** ✅ Full error
- **Artifacts:** None (too early)

### 2. Go Setup Failure
- **GitHub Logs:** ✅ Setup error
- **Artifacts:** None (too early)

### 3. Dependency Installation Failure
- **GitHub Logs:** ✅ Full error output
- **Artifacts:** ✅ `setup.log` (if step ran)

### 4. Service Startup Failure
- **GitHub Logs:** ✅ Service logs
- **Artifacts:** ✅ `service-health.log` (if step ran)

### 5. Test Execution Failure
- **GitHub Logs:** ✅ Test output
- **Artifacts:** ✅ `test-output.log`

### 6. Complete Workflow Failure
- **GitHub Logs:** ✅ All step logs
- **Artifacts:** ✅ All available artifacts
- **Summary:** ✅ Failure details in workflow summary

## Important Notes

### GitHub Actions Logs Are Always Available

**Even if:**
- ❌ Artifacts fail to upload
- ❌ Workflow fails completely
- ❌ Services don't start
- ❌ Tests never run

**You still have:**
- ✅ Complete workflow logs
- ✅ All step outputs
- ✅ Error messages
- ✅ Environment details

### Artifacts Are Bonus

Artifacts provide:
- ✅ Organized log files
- ✅ Easy to download and share
- ✅ Long retention (30 days)
- ✅ Consolidated reports

But **GitHub Actions logs are the primary source** and are always available.

## Best Practices

1. **Always Check GitHub Actions Logs First**
   - They're always available
   - They have complete information
   - They're easy to access

2. **Download Artifacts for Sharing**
   - Easier to share with team
   - Organized files
   - Can be archived

3. **Include Context**
   - Workflow run URL
   - Commit SHA
   - What you were doing

4. **Share Early**
   - Don't wait to debug
   - Share logs immediately
   - Get help faster

## Summary

**Yes, artifacts will be available** for you to get logs and send to fix everything, because:

1. ✅ **GitHub Actions automatically saves all workflow logs** - always available
2. ✅ **Our workflows save artifacts** - even if tests don't run
3. ✅ **Setup logs are captured** - dependency/build errors
4. ✅ **Service health is checked** - connection issues
5. ✅ **Test logs are saved** - if tests run
6. ✅ **Failure summaries are created** - what failed and why

**Even in complete failure scenarios, you'll have:**
- Complete workflow logs (automatic)
- Setup/build logs (artifacts)
- Service health logs (artifacts)
- Error messages and stack traces

**You can always get the logs you need to debug and fix issues!**
