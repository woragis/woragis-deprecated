# Backend TODO

## Services Overview
- **server** (Go)
- **email-worker** (Go)
- **job-application-worker** (Node.js)
- **resume-worker** (Python)
- **translation-worker** (Go)
- **whatsapp-worker** (Go)

---

## CI/CD

### Tasks
- [ ] CI/CD pipeline for server
- [ ] CI/CD pipeline for email-worker
- [ ] CI/CD pipeline for job-application-worker
- [ ] CI/CD pipeline for resume-worker
- [ ] CI/CD pipeline for translation-worker
- [ ] CI/CD pipeline for whatsapp-worker

### Strategy
**Platform:** GitHub Actions

**Workflow per service:**
1. **Test** - Run unit/integration tests
2. **Build** - Multi-stage Docker builds
3. **Publish** - Push to Docker Hub with tags:
   - `latest` (main branch)
   - `sha-{commit}` (traceability)
   - `v{version}` (releases)
4. **Deploy** - Trigger Railway via CLI/API/webhook

**Implementation:**
- Separate workflow file per service for independent builds
- Parallel execution across services
- Docker layer caching (especially for Node.js/Python)
- Tests on every push/PR
- Builds/deployments on merge to main
- Railway integration: GitHub integration (auto-deploy) or CLI/API (manual control)

---

## Tests

### Tasks
- [ ] Unit tests for server
- [ ] Integration tests for server
- [ ] Unit tests for email-worker
- [ ] Integration tests for email-worker
- [ ] Unit tests for job-application-worker
- [ ] Integration tests for job-application-worker
- [ ] Unit tests for resume-worker
- [ ] Integration tests for resume-worker
- [ ] Unit tests for translation-worker
- [ ] Integration tests for translation-worker
- [ ] Unit tests for whatsapp-worker
- [ ] Integration tests for whatsapp-worker

### Strategy

**Server (Go):**
- **Framework:** `testing` package (standard library) + `testify` for assertions
- **Unit tests:**
  - Test individual functions/methods in isolation
  - Mock external dependencies (database, RabbitMQ, external APIs)
  - Use table-driven tests for multiple scenarios
  - Target: 70-80% code coverage
- **Integration tests:**
  - Test API endpoints with test database
  - Use Docker Compose for dependencies (PostgreSQL, Redis, RabbitMQ)
  - Test authentication/authorization flows
  - Test database migrations
- **Test structure:** `*_test.go` files alongside source code
- **Running tests:** `go test ./...` with coverage `-cover`

**Go Workers (email-worker, translation-worker, whatsapp-worker):**
- **Framework:** Same as server (`testing` + `testify`)
- **Unit tests:**
  - Test queue message processing logic
  - Mock RabbitMQ connections
  - Test error handling and retries
- **Integration tests:**
  - Test with real RabbitMQ (test container)
  - Test message consumption and processing
  - Test connection recovery

**Node.js Worker (job-application-worker):**
- **Framework:** Jest or Mocha + Chai
- **Unit tests:**
  - Test scraping logic, selector finding, cover letter generation
  - Mock external APIs and database
  - Test orchestrator workflow
- **Integration tests:**
  - Test with test RabbitMQ instance
  - Test end-to-end job application flow
  - Test database interactions
- **Coverage:** Use `nyc` or `c8` for coverage reports

**Python Worker (resume-worker):**
- **Framework:** `pytest` (recommended) or `unittest`
- **Unit tests:**
  - Test resume generation logic
  - Test translation helpers
  - Test keyword extraction
  - Mock AI service calls
- **Integration tests:**
  - Test with test RabbitMQ instance
  - Test database operations
  - Test end-to-end resume generation
- **Fixtures:** Use pytest fixtures for test data
- **Coverage:** Use `pytest-cov` for coverage reports

**Test Infrastructure:**
- **Test containers:** Use Docker Compose or testcontainers for dependencies
- **Test database:** Separate test database, reset between tests
- **CI Integration:** Run tests in GitHub Actions on every push/PR
- **Coverage reports:** Generate and track coverage over time
- **Test data:** Use factories/fixtures for consistent test data

**Best Practices:**
- Write tests before or alongside code (TDD when possible)
- Keep tests fast (unit tests < 1s, integration tests < 30s)
- Use meaningful test names describing what is tested
- Isolate tests (no shared state between tests)
- Test error cases and edge cases
- Mock external services to avoid flakiness

---

## Health Checks

### Tasks
- [ ] Health check endpoint for server
- [ ] Health check for email-worker
- [ ] Health check for job-application-worker
- [ ] Health check for resume-worker
- [ ] Health check for translation-worker
- [ ] Health check for whatsapp-worker

### Strategy

**Server (HTTP endpoint):**
- **Endpoint:** `GET /healthz` or `GET /health`
- **Liveness probe:** Basic service availability (200 OK)
- **Readiness probe:** Check dependencies:
  - Database connection
  - RabbitMQ connection
  - External API availability (if critical)
- **Response format:** JSON with status and component health
  ```json
  {
    "status": "healthy|degraded|unhealthy",
    "checks": {
      "database": "ok|error",
      "rabbitmq": "ok|error",
      "external_apis": "ok|error"
    }
  }
  ```

**Workers (HTTP endpoint or signal-based):**
- **Option 1:** HTTP server with `/healthz` endpoint (recommended)
- **Option 2:** File-based health check (write/update health file)
- **Checks:**
  - RabbitMQ connection status
  - Database connection (if used)
  - Worker processing capability (can consume messages)
- **Response:** Simple 200 OK or detailed component status

**Kubernetes Integration:**
- **Liveness probe:** Restart container if unhealthy
- **Readiness probe:** Remove from service endpoints if not ready
- **Startup probe:** Allow time for initialization
- **Configuration:**
  - Initial delay: 10-30s (depending on startup time)
  - Period: 10s
  - Timeout: 5s
  - Failure threshold: 3

**Implementation Notes:**
- Use lightweight checks (avoid heavy operations)
- Cache dependency status (update every 5-10s, not on every request)
- Return appropriate HTTP status codes (200 healthy, 503 unhealthy)
- Log health check failures for monitoring

---

## Kubernetes

### Tasks
- [ ] Kubernetes deployment for server
- [ ] Kubernetes deployment for email-worker
- [ ] Kubernetes deployment for job-application-worker
- [ ] Kubernetes deployment for resume-worker
- [ ] Kubernetes deployment for translation-worker
- [ ] Kubernetes deployment for whatsapp-worker

### Strategy

**Manifest Structure:**
- Separate manifests per service (deployment, service, configmap, secret)
- Organize in `k8s/` directory with subdirectories per service
- Use Kustomize or Helm for environment-specific configs (dev/staging/prod)

**Deployment Configuration:**
- **Replicas:**
  - Server: 2-3 (high availability)
  - Workers: 1-2 per worker type (scale based on queue depth)
- **Resource limits/requests:**
  - Server: 500m CPU, 512Mi memory (requests), 1000m CPU, 1Gi memory (limits)
  - Go workers: 200m CPU, 256Mi memory (requests), 500m CPU, 512Mi memory (limits)
  - Node.js/Python workers: 300m CPU, 512Mi memory (requests), 1000m CPU, 1Gi memory (limits)
- **Image pull policy:** `IfNotPresent` (or `Always` for `latest` tag)

**Service Definitions:**
- **Server:** ClusterIP or LoadBalancer (if external access needed)
- **Workers:** ClusterIP (internal only, for health checks)
- **Ports:** Standardize on 8080 for HTTP health endpoints

**Configuration Management:**
- **ConfigMaps:** Environment variables, non-sensitive config
- **Secrets:** Database credentials, API keys, RabbitMQ credentials
- **External secrets operator:** Consider for secret management (optional)

**Scaling:**
- **Horizontal Pod Autoscaler (HPA):**
  - Server: Based on CPU/memory usage (target 70%)
  - Workers: Based on queue depth (custom metrics) or CPU
- **Vertical Pod Autoscaler (VPA):** Optional, for right-sizing resources

**Health Probes Integration:**
- **Liveness probe:** HTTP GET `/healthz` (from health check implementation)
- **Readiness probe:** HTTP GET `/healthz` (from health check implementation)
- **Startup probe:** HTTP GET `/healthz` with longer initial delay
- **Probe settings:** Period 10s, timeout 5s, failure threshold 3

**Ingress (Server only):**
- **Ingress controller:** nginx or traefik
- **TLS:** Cert-manager for automatic SSL certificates
- **Routes:** Domain-based or path-based routing

**Namespace Organization:**
- **Default:** `woragis` or `woragis-production`
- **Alternative:** Separate namespaces per environment (dev/staging/prod)
- **Resource quotas:** Set per namespace to prevent resource exhaustion

**Storage:**
- **Persistent volumes:** Only if workers need local storage (unlikely)
- **Ephemeral storage:** Default for logs/temp files

**Networking:**
- **Service discovery:** Use Kubernetes DNS (`service-name.namespace.svc.cluster.local`)
- **Network policies:** Optional, restrict inter-pod communication if needed

**Implementation Notes:**
- Use deployment strategy: `RollingUpdate` with maxSurge=1, maxUnavailable=0
- Set pod disruption budgets for high availability
- Use node selectors/affinity if specific node requirements exist
- Consider pod security policies/standards for security
- Integrate with monitoring (Prometheus metrics, Grafana dashboards)

---

## Structured Logs

### Tasks
- [ ] Structured logging for server
- [ ] Structured logging for email-worker
- [ ] Structured logging for job-application-worker
- [ ] Structured logging for resume-worker
- [ ] Structured logging for translation-worker
- [ ] Structured logging for whatsapp-worker

### Strategy

**Format:** JSON (production) or human-readable (development)

**Required Fields (all logs):**
- `timestamp` - ISO 8601 format
- `level` - `debug|info|warn|error|fatal`
- `service` - Service name (server, email-worker, etc.)
- `message` - Human-readable message
- `trace_id` - Request/operation correlation ID (for distributed tracing)

**Server (Go):**
- **Library:** `log/slog` (standard library, already in use)
- **Configuration:**
  - JSON handler for production
  - Text handler for development
  - Set log level via environment variable
- **Context fields:**
  - Request ID, user ID, HTTP method, path, status code
  - Database query duration, error details
  - External API call details
- **Structured fields:** Use `slog.String()`, `slog.Int()`, `slog.Any()` for key-value pairs

**Go Workers (email-worker, translation-worker, whatsapp-worker):**
- **Library:** `log/slog` (standard library)
- **Context fields:**
  - Job ID, message ID, queue name
  - Processing duration, retry count
  - Error details with stack traces
- **Log levels:**
  - `info`: Job started/completed, connection established
  - `warn`: Retries, connection issues
  - `error`: Processing failures, critical errors

**Node.js Worker (job-application-worker):**
- **Library:** `pino` (recommended) or `winston` with JSON formatter
- **Current:** Custom logger utility (migrate to structured format)
- **Context fields:**
  - Job ID, user ID, website URL
  - Scraping attempts, selector cache hits
  - Processing duration, error details
- **Configuration:**
  - JSON output for production
  - Pretty print for development
  - Log level from environment variable

**Python Worker (resume-worker):**
- **Library:** `structlog` (recommended) or `python-json-logger`
- **Current:** Standard `logging` module (migrate to structured)
- **Context fields:**
  - Resume ID, user ID, language
  - Generation duration, AI service used
  - Error details with stack traces
- **Configuration:**
  - JSON formatter for production
  - Console formatter for development
  - Log level from environment variable

**Log Levels:**
- **debug:** Detailed diagnostic information (development only)
- **info:** General operational messages (startup, normal operations)
- **warn:** Warning messages (retries, degraded functionality)
- **error:** Error messages (failures, exceptions)
- **fatal:** Critical errors (application cannot continue)

**Log Aggregation:**
- **Output:** stdout/stderr (Kubernetes will collect)
- **Aggregation:** Use Fluentd, Fluent Bit, or similar log forwarder
- **Storage:** Send to centralized logging (ELK stack, Loki, CloudWatch, etc.)
- **Search:** Enable full-text search on structured fields

**Best Practices:**
- Never log sensitive data (passwords, tokens, PII)
- Use consistent field names across services
- Include correlation IDs for tracing requests across services
- Log at appropriate levels (don't spam info logs)
- Use structured fields instead of string interpolation
- Log errors with full context (stack traces, request details)
- Set appropriate log levels per environment (debug in dev, info in prod)
