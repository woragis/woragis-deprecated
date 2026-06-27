# Backend Posts Index

Navigation index for all backend-related posts and articles.

## Architecture

### System Architecture
- [Microservices Architecture Overview](./architecture/microservices-architecture-overview.md)
- [Message Queue Patterns](./architecture/message-queue-patterns.md)
- [Worker Architecture: Standalone](./architecture/worker-architecture-standalone.md)
- [Service Communication Patterns](./architecture/service-communication-patterns.md)
- [Database Design: PostgreSQL with JSONB](./architecture/database-design-postgresql-jsonb.md)
- [API Design: RESTful Services](./architecture/api-design-restful.md)
- [Domain-Driven Design](./architecture/domain-driven-design.md)
- [Event-Driven Architecture](./architecture/event-driven-architecture.md)
- [Service Boundaries](./architecture/service-boundaries.md)
- [Data Consistency Patterns](./architecture/data-consistency-patterns.md)

## Technical Decisions (ADRs)

- [ADR-001: Go vs Python](./technical-decisions/adr-go-vs-python.md)
- [ADR-002: RabbitMQ + Redis Fallback](./technical-decisions/adr-rabbitmq-redis-fallback.md)
- [ADR-003: Structured Logging](./technical-decisions/adr-structured-logging.md)
- [ADR-004: Translation Service Architecture](./technical-decisions/adr-translation-service-architecture.md)
- [ADR-005: Health Checks Patterns](./technical-decisions/adr-health-checks-patterns.md)
- [ADR-006: Testing Strategies](./technical-decisions/adr-testing-strategies.md)
- [ADR-007: Docker Compose Architecture](./technical-decisions/adr-docker-compose-architecture.md)
- [ADR-008: Database PostgreSQL](./technical-decisions/adr-database-postgresql.md)

## Implementation

- [Dead Letter Queues](./implementation/dead-letter-queues.md)
- [Retry Policies](./implementation/retry-policies.md)
- [Graceful Degradation](./implementation/graceful-degradation.md)
- [Health Check Caching](./implementation/health-check-caching.md)
- [Trace ID Propagation](./implementation/trace-id-propagation.md)
- [Error Handling Patterns](./implementation/error-handling-patterns.md)
- [Configuration Management](./implementation/configuration-management.md)
- [Connection Pooling](./implementation/connection-pooling.md)
- [Request ID Middleware](./implementation/request-id-middleware.md)
- [Structured Logging: JSON](./implementation/structured-logging-json.md)
- [Queue Declarations](./implementation/queue-declarations.md)
- [Worker Lifecycle](./implementation/worker-lifecycle.md)

## Cross-Cutting Concerns

- [Observability Overview](./cross-cutting/observability-overview.md)
- [Metrics: Prometheus](./cross-cutting/metrics-prometheus.md)
- [Tracing: OpenTelemetry](./cross-cutting/tracing-opentelemetry.md)
- [Circuit Breakers](./cross-cutting/circuit-breakers.md)
- [Rate Limiting](./cross-cutting/rate-limiting.md)
- [Timeout Strategies](./cross-cutting/timeout-strategies.md)
- [Error Classification](./cross-cutting/error-classification.md)
- [Log Aggregation](./cross-cutting/log-aggregation.md)
- [Monitoring: Health + Metrics](./cross-cutting/monitoring-health-metrics.md)
- [Alerting Strategies](./cross-cutting/alerting-strategies.md)

## Lessons Learned

- [5 Things I Learned Building Microservices](./lessons-learned/five-things-microservices.md)
- [Observability: What Works](./lessons-learned/observability-what-works.md)
- [Testing Distributed Systems](./lessons-learned/testing-distributed-systems.md)
- [Operational Patterns](./lessons-learned/operational-patterns.md)
- [Message Queue Pitfalls](./lessons-learned/message-queue-pitfalls.md)
- [Worker Patterns](./lessons-learned/worker-patterns.md)
- [Database Query Performance](./lessons-learned/database-query-performance.md)
- [Deployment Lessons](./lessons-learned/deployment-lessons.md)

## Advanced Topics

- [Scaling Strategies](./advanced/scaling-strategies.md)
- [Performance Optimization](./advanced/performance-optimization.md)
- [Security Patterns](./advanced/security-patterns.md)
- [Data Migration](./advanced/data-migration.md)
- [Disaster Recovery](./advanced/disaster-recovery.md)
- [Cost Optimization](./advanced/cost-optimization.md)
- [Load Testing](./advanced/load-testing.md)
- [Capacity Planning](./advanced/capacity-planning.md)
- [Multi-Region](./advanced/multi-region.md)
- [Service Mesh](./advanced/service-mesh.md)

## Meta/Documentation

- [Documentation Strategy](./meta/documentation-strategy.md)
- [Technical Writing](./meta/technical-writing.md)
- [Architecture Diagrams](./meta/architecture-diagrams.md)
- [Runbooks](./meta/runbooks.md)
- [ADRs: Capturing Decisions](./meta/adrs-capturing-decisions.md)

## Existing Content

### Domains
- [Auth Domain](../domain-auth-architecture.md)
- [Translations Domain](../domain-translations-multilanguage.md)
- [Resumes Domain](../domain-resumes-generation.md)
- [Job Applications Domain](../domain-jobapplications-workflow.md)
- [Chats Domain](../domain-chats-realtime.md)
- [Projects Domain](../domain-projects-management.md)
- [Finances Domain](../domain-finances-operations.md)

### Workers
- [Translation Worker](../worker-translation-queue.md)
- [Resume Worker](../worker-resume-queue.md)
- [Job Application Worker](../worker-jobapplication-queue.md)
- [Email Worker](../worker-email-queue.md)
- [WhatsApp Worker](../worker-whatsapp-queue.md)
- [Scheduler Worker](../worker-scheduler-queue.md)

### Infrastructure
- [Logging Strategy](../backend-logging-strategy.md)
- [Caching Strategies](../backend-caching-strategies.md)
- [Database Queries](../backend-database-queries.md)

### Testing
- [Tests Index](./tests/index.md)

### CI/CD
- [CI/CD Index](./cicd/index.md)

### Health Checks
- [Health Checks Index](./health-checks/index.md)
