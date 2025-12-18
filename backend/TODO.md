# Backend TODO

## Services Overview
- **server** (Go)
- **email-worker** (Go)
- **job-application-worker** (Node.js)
- **resume-worker** (Python)
- **translation-worker** (Go)
- **whatsapp-worker** (Go)
- **ai-service** (Python)
- **creative-service** (Python)

---

## CI/CD

### Tasks
- [x] CI/CD pipeline for server (via docker-build-push.yml)
- [x] CI/CD pipeline for email-worker (via docker-build-push.yml)
- [x] CI/CD pipeline for job-application-worker (via docker-build-push.yml)
- [x] CI/CD pipeline for resume-worker (via docker-build-push.yml)
- [x] CI/CD pipeline for translation-worker (via docker-build-push.yml)
- [x] CI/CD pipeline for whatsapp-worker (via docker-build-push.yml)

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
- ~~Separate workflow file per service for independent builds~~ (Currently using combined workflows: docker-build-push.yml and test-only.yml)
- Parallel execution across services ✅
- Docker layer caching (especially for Node.js/Python) ✅
- Tests on every push/PR ✅ (test-only.yml)
- Builds/deployments on merge to main ✅ (docker-build-push.yml on tags)
- Railway integration: GitHub integration (auto-deploy) or CLI/API (manual control) ✅ (railway-deploy.yml)

---

## Tests

### Tasks
- [x] Unit tests for server (partial - service_test.go files exist)
- [ ] Integration tests for server
- [x] Unit tests for email-worker (config, queue, sender, logger, health tests exist)
- [ ] Integration tests for email-worker
- [x] Unit tests for job-application-worker (coverLetter, health, orchestrator tests exist)
- [ ] Integration tests for job-application-worker
- [x] Unit tests for resume-worker (comprehensive unit tests exist)
- [x] Integration tests for resume-worker (test_worker.py exists)
- [x] Unit tests for translation-worker (comprehensive unit tests exist)
- [ ] Integration tests for translation-worker
- [x] Unit tests for whatsapp-worker (config, queue, notifier, logger, health tests exist)
- [ ] Integration tests for whatsapp-worker
- [x] Unit tests for ai-service (agents, api, providers tests exist)
- [x] Integration tests for ai-service (test_api.py exists)
- [x] Unit tests for creative-service (api, provider implementations, providers tests exist)
- [x] Integration tests for creative-service (test_api.py exists)

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
- [x] Health check endpoint for server
- [x] Health check for email-worker
- [x] Health check for job-application-worker
- [x] Health check for resume-worker
- [x] Health check for translation-worker
- [x] Health check for whatsapp-worker
- [x] Health check for ai-service
- [x] Health check for creative-service

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
- [x] Kubernetes deployment for whatsapp-worker (StatefulSet and deployment-leader-election manifests exist)

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
- [x] Structured logging for server
- [x] Structured logging for email-worker
- [x] Structured logging for job-application-worker
- [x] Structured logging for resume-worker
- [x] Structured logging for translation-worker
- [x] Structured logging for whatsapp-worker
- [x] Structured logging for ai-service
- [x] Structured logging for creative-service

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

**Log Storage:**
- **Development:** stdout by default (optional file logging with `LOG_TO_FILE=true`)
- **Production:** stdout only (Kubernetes will collect automatically)
- **File logging:** Only available in development, writes to `logs/` directory (dual output: file + stdout)

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

---

## Technical Documentation

### Tasks
- [ ] Architecture overview document with high-level diagram
- [ ] ADR-001: RabbitMQ with Redis fallback
- [ ] ADR-002: Standalone workers architecture
- [ ] ADR-003: Structured logging implementation
- [ ] ADR-004: Translation worker architecture (Go + external APIs)
- [ ] ADR-005: Health checks implementation strategy
- [ ] Runbook: Monitoring Dead Letter Queues
- [ ] Runbook: Deploying services and workers
- [ ] Runbook: Troubleshooting common issues
- [ ] Component documentation (one page per component)
- [ ] API documentation for server endpoints
- [ ] API documentation for AI service
- [ ] API documentation for Creative service
- [ ] Development guide (how to add new domain, worker, etc.)

### Strategy

**Architecture Decision Records (ADRs):**
- **Purpose:** Capture why decisions were made, context, and trade-offs
- **Format:** Markdown files in `docs/adr/` directory
- **Structure:**
  - Context: What problem are we solving?
  - Decision: What did we choose?
  - Consequences: What are the trade-offs? (pros/cons)
  - Status: Accepted/Deprecated/Superseded
- **Priority ADRs:**
  1. RabbitMQ + Redis fallback (high availability strategy)
  2. Standalone workers (scalability and separation of concerns)
  3. Structured logging (observability)
  4. Translation worker architecture (hybrid approach: Go + external APIs)
  5. Health checks implementation (monitoring and reliability)

**Architecture Diagrams:**
- **High-level system architecture:**
  - Server → Services (AI, Creative)
  - Server → Workers (Email, WhatsApp, Translation, Resume, Job App)
  - RabbitMQ as message broker
  - Redis for caching/fallback
  - PostgreSQL for data
- **Component interaction diagrams:**
  - Request flow: API → Server → Queue → Worker → DB
  - Translation flow: Server → RabbitMQ → Translation Worker → Translation API → DB
  - Data flow through the system
- **Deployment architecture:**
  - Docker containers and networks
  - Dependencies between services
- **Tools:** Mermaid diagrams (in Markdown) or draw.io/Excalidraw
- **Location:** `docs/architecture/` directory

**Runbooks:**
- **Purpose:** Step-by-step guides for common operational tasks
- **Format:** Markdown files in `docs/runbooks/` directory
- **Priority Runbooks:**
  1. Monitoring Dead Letter Queues
     - How to check DLQ in RabbitMQ
     - What to look for (normal vs abnormal)
     - How to reprocess failed jobs
     - When to alert
  2. Deploying services and workers
     - Steps to deploy new version
     - Rollback procedure
     - Verification steps
  3. Troubleshooting common issues
     - Translation failures
     - Queue backlogs
     - Database connection issues
     - Worker not processing jobs
- **Structure:**
  - When to use this runbook
  - Prerequisites
  - Step-by-step instructions
  - Expected outcomes
  - Troubleshooting tips

**Component Documentation:**
- **Purpose:** Explain what each component does and how it works
- **Format:** One Markdown file per component in `docs/components/` directory
- **Content per component:**
  - Purpose and responsibilities
  - Dependencies (what it needs)
  - Configuration options
  - How it interacts with other components
  - Common issues and solutions
- **Components to document:**
  - Server (main API)
  - Email Worker
  - WhatsApp Worker
  - Translation Worker
  - Resume Worker
  - Job Application Worker
  - AI Service
  - Creative Service

**API Documentation:**
- **Purpose:** How to use the APIs
- **Format:** OpenAPI/Swagger specs or Markdown with examples
- **Location:** `docs/api/` directory
- **APIs to document:**
  - Server API endpoints (REST)
  - AI Service API (FastAPI)
  - Creative Service API (FastAPI)
- **Content:**
  - Endpoint descriptions
  - Request/response examples
  - Error codes and meanings
  - Authentication requirements
  - Rate limiting (if applicable)

**Development Guides:**
- **Purpose:** How to work with the codebase
- **Format:** Markdown files in `docs/development/` directory
- **Guides to create:**
  1. How to add a new domain
     - File structure
     - Required components (repository, service, handler)
     - Testing requirements
  2. How to add a new worker
     - Worker structure
     - Queue setup
     - Configuration
     - Testing
  3. How to add a new service
     - Service structure
     - API endpoints
     - Testing
  4. Testing patterns
     - Unit test patterns
     - Integration test patterns
     - Mocking strategies
  5. Logging conventions
     - When to log
     - What to log
     - Log levels
  6. Error handling patterns
     - Error types
     - Error propagation
     - Error responses

**Documentation Structure:**
```
docs/
├── adr/                    # Architecture Decision Records
│   ├── 001-rabbitmq-redis-fallback.md
│   ├── 002-standalone-workers.md
│   ├── 003-structured-logging.md
│   ├── 004-translation-worker.md
│   └── 005-health-checks.md
├── architecture/            # Architecture diagrams and descriptions
│   ├── system-overview.md
│   ├── component-interaction.md
│   └── deployment.md
├── runbooks/               # Operational runbooks
│   ├── monitoring-dlq.md
│   ├── deploying-services.md
│   └── troubleshooting.md
├── components/             # Component documentation
│   ├── server.md
│   ├── email-worker.md
│   ├── whatsapp-worker.md
│   ├── translation-worker.md
│   ├── resume-worker.md
│   ├── job-application-worker.md
│   ├── ai-service.md
│   └── creative-service.md
├── api/                    # API documentation
│   ├── server-api.md
│   ├── ai-service-api.md
│   └── creative-service-api.md
└── development/           # Development guides
    ├── adding-domain.md
    ├── adding-worker.md
    ├── adding-service.md
    ├── testing-patterns.md
    ├── logging-conventions.md
    └── error-handling.md
```

**Best Practices:**
- Keep documentation close to code (in repository)
- Use Markdown for easy editing and version control
- Include diagrams (Mermaid or images)
- Keep it simple and focused
- Update documentation when making significant changes
- Link from README.md for discoverability
- Use consistent formatting and structure
- Include examples and code snippets
- Keep runbooks actionable (step-by-step)
- Review and update periodically

**Implementation Priority:**
1. **Phase 1 (Quick wins):**
   - Architecture overview document
   - 3-5 priority ADRs
   - 2-3 critical runbooks
2. **Phase 2 (Expansion):**
   - Component documentation
   - API documentation
   - Development guides
3. **Phase 3 (Maintenance):**
   - Keep documentation updated
   - Add new ADRs as decisions are made
   - Expand runbooks as needed
