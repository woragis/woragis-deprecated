# Worker Error Logs

**Date:** 2025-12-22  
**Services Analyzed:** `woragis-resume-worker`, `woragis-job-application-worker`

---

## Summary

- **resume-worker**: ✅ No critical errors found (only verbose debug output from pika/RabbitMQ client)
- **job-application-worker**: ❌ **Critical Redis connection errors** - Worker fails to start due to Redis connection issues

---

## 1. Job Application Worker - Redis Connection Errors

**Service:** `woragis-job-application-worker`  
**Status:** ❌ **FAILING** - Worker crashes and restarts continuously  
**Error Type:** Redis Connection Error

### Error Messages

```
{"timestamp":"2025-12-22T06:36:22.140Z","level":"error","service":"job-application-worker","message":"SelectorCache Redis error","error":"connect ECONNREFUSED 127.0.0.1:6379"}
```

This error repeats continuously (20 retries) until the worker fails:

```
{"timestamp":"2025-12-22T06:36:39.899Z","level":"error","service":"job-application-worker","message":"Worker failed to start","error":"Reached the max retries per request limit (which is 20). Refer to \"maxRetriesPerRequest\" option for details."}
```

### Root Cause

The `SelectorCache` class in `backend/job-application-worker/src/selectorCache.js` is trying to connect to Redis at `127.0.0.1:6379` (localhost), but in Docker it should connect to the Redis service using the service name `redis:6379`.

**Problematic Code:**
```javascript
// Line 13 in selectorCache.js
const redisUrl = process.env.REDIS_URL || 'redis://localhost:6379/0';
```

The default fallback `redis://localhost:6379/0` doesn't work in Docker because:
- The worker container doesn't have Redis running on localhost
- Redis is running in a separate container named `redis`
- The `REDIS_URL` environment variable is not set in `docker-compose.yml` for the job-application-worker service

### Current Configuration

Looking at `docker-compose.yml`, the `job-application-worker` service does NOT have a `REDIS_URL` environment variable set, so it falls back to `localhost:6379`.

### Fix Required

1. **Option 1 (Recommended):** Add `REDIS_URL` environment variable to `job-application-worker` service in `docker-compose.yml`:
   ```yaml
   job-application-worker:
     environment:
       REDIS_URL: redis://redis:6379/0
       # ... other env vars
   ```

2. **Option 2:** Update the default in `selectorCache.js` to use the Docker service name:
   ```javascript
   const redisUrl = process.env.REDIS_URL || 'redis://redis:6379/0';
   ```
   (Note: This would break local development, so Option 1 is better)

### Impact

- Worker cannot start successfully
- Worker restarts continuously in a crash loop
- Selector caching functionality is unavailable
- Worker may still process jobs but without selector caching optimization

### Log Sample

Full error pattern (repeats 20 times before worker fails):
```json
{"timestamp":"2025-12-22T06:37:02.620Z","level":"error","service":"job-application-worker","message":"SelectorCache Redis error","error":"connect ECONNREFUSED 127.0.0.1:6379"}
{"timestamp":"2025-12-22T06:37:02.675Z","level":"error","service":"job-application-worker","message":"SelectorCache Redis error","error":"connect ECONNREFUSED 127.0.0.1:6379"}
{"timestamp":"2025-12-22T06:37:02.780Z","level":"error","service":"job-application-worker","message":"SelectorCache Redis error","error":"connect ECONNREFUSED 127.0.0.1:6379"}
...
{"timestamp":"2025-12-22T06:37:13.287Z","level":"error","service":"job-application-worker","message":"Worker failed to start","error":"Reached the max retries per request limit (which is 20). Refer to \"maxRetriesPerRequest\" option for details."}
```

---

## 2. Resume Worker - No Critical Errors

**Service:** `woragis-resume-worker`  
**Status:** ✅ **Running** - No critical errors detected  
**Error Type:** None (only verbose debug output)

### Observations

The resume-worker logs contain extensive verbose debug output from the `pika` RabbitMQ client library, but no actual errors. The messages include:

- Connection lifecycle messages (heartbeat frames, connection state changes)
- Normal shutdown messages (when connections are closed gracefully)
- Verbose logging from `SelectorIOServicesAdapter` and `HeartbeatChecker`

### Example Log Output

```
Sending heartbeat frame
SelectorIOServicesAdapter.set_writer(7, <bound method _AsyncPlaintextTransport._on_socket_writable...>)
Received heartbeat frame
Received 886 heartbeat frames, sent 886, idle intervals 0
```

### Status

These are **informational/debug messages**, not errors. The worker appears to be:
- ✅ Connected to RabbitMQ successfully
- ✅ Processing messages (heartbeat frames indicate active connection)
- ✅ Running without critical issues

### Recommendation

If the verbose logging is undesirable, consider:
1. Reducing pika log level in the Python code
2. Filtering out debug messages in log aggregation
3. These messages don't indicate any problems and can be safely ignored

---

## Next Steps

### Immediate Actions Required

1. **Fix job-application-worker Redis connection:**
   - Add `REDIS_URL: redis://redis:6379/0` to `job-application-worker` environment variables in `docker-compose.yml`
   - Restart the worker: `docker-compose restart job-application-worker`
   - Verify connection: Check logs for successful Redis connection

2. **Verify fix:**
   ```bash
   docker logs woragis-job-application-worker --tail 50
   ```
   Should see successful connection messages instead of connection errors.

### Optional Actions

1. **Reduce resume-worker verbosity:**
   - Configure pika logging level if needed
   - Not critical - worker is functioning correctly

---

## Files Referenced

- `backend/docker-compose.yml` - Docker service configuration
- `backend/job-application-worker/src/selectorCache.js` - Redis connection code (line 13)
- `backend/job-application-worker/src/index.js` - Worker entry point
- `backend/resume-worker/src/main.py` - Resume worker entry point

---

## Environment Variables Reference

### Current (Missing)
```yaml
job-application-worker:
  environment:
    # REDIS_URL is missing - should be added
```

### Should Be
```yaml
job-application-worker:
  environment:
    REDIS_URL: redis://redis:6379/0  # Add this line
    # ... other existing env vars
```

---

**Note:** All logs were collected on 2025-12-22. Error patterns may change after fixes are applied.
