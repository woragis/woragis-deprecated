# Woragis Codebase Posts Index

Navigation index for potential articles about different aspects of the Woragis codebase.

## Workers

### Translation Worker
- [Queue Strategy](./worker-translation-queue.md)
- [Logging Strategy](./worker-translation-logging.md)
- [Orchestration Pattern](./worker-translation-orchestration.md)

### Resume Worker
- [Queue Strategy](./worker-resume-queue.md)
- [Logging Strategy](./worker-resume-logging.md)
- [Orchestration Pattern](./worker-resume-orchestration.md)

### Job Application Worker
- [Queue Strategy](./worker-jobapplication-queue.md)
- [Logging Strategy](./worker-jobapplication-logging.md)
- [Orchestration Pattern](./worker-jobapplication-orchestration.md)

### WhatsApp Worker
- [Queue Strategy](./worker-whatsapp-queue.md)
- [Logging Strategy](./worker-whatsapp-logging.md)

### Email Worker
- [Queue Strategy](./worker-email-queue.md)
- [Logging Strategy](./worker-email-logging.md)

### Scheduler Worker
- [Queue Strategy](./worker-scheduler-queue.md)
- [Logging Strategy](./worker-scheduler-logging.md)

## Domains

### Auth Domain
- [Architecture & Flow](./domain-auth-architecture.md)
- [JWT Management](./domain-auth-jwt.md)
- [OAuth Integration](./domain-auth-oauth.md)
- [MFA Implementation](./domain-auth-mfa.md)
- [Session Management](./domain-auth-sessions.md)
- [Audit Logging](./domain-auth-audit.md)

### Translations Domain
- [Multilanguage Support](./domain-translations-multilanguage.md)
- [Translation Queue & Processing](./domain-translations-queue.md)
- [Language Detection Middleware](./domain-translations-middleware.md)
- [Database Joins & Queries](./domain-translations-queries.md)

### Resumes Domain
- [Resume Generation Flow](./domain-resumes-generation.md)
- [Queue & Job Management](./domain-resumes-queue.md)
- [Database Queries & Joins](./domain-resumes-queries.md)
- [File Management](./domain-resumes-files.md)

### Job Applications Domain
- [Application Workflow](./domain-jobapplications-workflow.md)
- [Queue Strategy](./domain-jobapplications-queue.md)
- [Database Queries & Joins](./domain-jobapplications-queries.md)
- [Interview Stages](./domain-jobapplications-interviews.md)

### Chats Domain
- [Real-time Communication](./domain-chats-realtime.md)
- [Stream Handling](./domain-chats-streams.md)
- [Database Queries](./domain-chats-queries.md)

### Projects Domain
- [Project Management](./domain-projects-management.md)
- [Translation Integration](./domain-projects-translations.md)
- [Database Queries](./domain-projects-queries.md)

### Finances Domain
- [Financial Operations](./domain-finances-operations.md)
- [Database Queries](./domain-finances-queries.md)

## Backend Infrastructure

### Logging
- [Backend Logging Strategy](./backend-logging-strategy.md)
- [Backend Logging Visualization](./backend-logging-visualization.md) *(To be implemented)*

### Caching
- [Caching Strategies](./backend-caching-strategies.md)
- [Redis Integration](./backend-caching-redis.md)

### Database
- [Query Patterns](./backend-database-queries.md)
- [Join Strategies](./backend-database-joins.md)
- [Query Optimization](./backend-database-optimization.md)

### Multilanguage Support
- [Backend Multilanguage](./backend-multilanguage.md)
- [Frontend Multilanguage](./frontend-multilanguage.md)
- [Landing Page Multilanguage](./landing-multilanguage.md)

## Testing

- [Tests Index](./tests/index.md)
- Unit Tests (Backend, Frontend, Workers)
- Integration Tests (API, Database, Workers)
- E2E Tests
- Test Infrastructure

## CI/CD

- [CI/CD Index](./cicd/index.md)
- CI/CD Strategy & Pipelines
- Docker Build Strategy
- Kubernetes Deployment
- Testing in CI/CD

## Health Checks

- [Health Checks Index](./health-checks/index.md)
- Health Check Strategy
- API Health Endpoints
- Container Health Checks
- Prometheus Integration

## Status

### ✅ Completed Skeleton Files

#### Workers (All Complete)
- Translation Worker (Queue, Logging, Orchestration)
- Resume Worker (Queue, Logging, Orchestration)
- Job Application Worker (Queue, Logging, Orchestration)
- WhatsApp Worker (Queue, Logging)
- Email Worker (Queue, Logging)
- Scheduler Worker (Queue, Logging)

#### Domains (Complete)
- Translations Domain (Multilanguage, Queue, Middleware, Queries)
- Auth Domain (Architecture, JWT, OAuth, MFA, Sessions, Audit)
- Resumes Domain (Generation, Queue, Queries, Files)
- Job Applications Domain (Workflow, Queue, Queries, Interviews)
- Chats Domain (Realtime, Streams, Queries)
- Projects Domain (Management, Translations, Queries)
- Finances Domain (Operations, Queries)

#### Backend Infrastructure
- Backend Logging Strategy
- Backend Logging Visualization
- Caching Strategies
- Redis Integration
- Database Queries
- Database Joins
- Database Optimization

#### Frontend & Landing
- Frontend Multilanguage
- Landing Page Multilanguage
- Backend Multilanguage

#### Testing
- Unit Tests (Go Backend, Domains, Services)
- Integration Tests (API, Workers)
- E2E Test Strategy
- Test Infrastructure Setup

#### CI/CD
- CI/CD Strategy
- GitHub Actions Pipeline
- Docker Build Strategy
- Kubernetes Deployment
- Database Migrations

#### Health Checks
- Health Check Strategy
- API Health Endpoints
- Container Health Checks
- Prometheus Integration
- Database & Redis Health Checks
- Health Check Monitoring

### 📝 To Be Created
- All major skeletons completed! 🎉
- Some minor topics may be added as codebase evolves

## Improvements & Future Work

### Translation Worker
- **Queue**: Priority support, dead letter queue, retry mechanism, job timeout, metrics
- **Logging**: Request ID/trace ID, performance metrics, log rotation, error context
- **Orchestration**: Worker pool, rate limiting, circuit breaker, job batching, scheduling

### Resume Worker
- **Queue**: Dequeue operation, priority queue, job timeout, cancellation, result caching
- **Logging**: Trace IDs, correlation IDs, log aggregation, performance profiling, metrics
- **Orchestration**: Worker pool, prioritization, rate limiting, caching, partial generation, template selection

### Job Application Worker
- **Queue**: Priority queue, dead letter queue, status tracking, timeout, deduplication
- **Logging**: Trace IDs, log rotation, scraping metrics, rate limiting logs, structured errors
- **Orchestration**: Worker pool, prioritization, circuit breaker, browser pooling, proxy rotation, CAPTCHA solving

### Backend Logging
- **Strategy**: File-based logging implementation, log rotation, request ID middleware, log aggregation
- **Visualization**: ELK/Grafana integration, real-time dashboard, error pattern analysis, performance metrics

### Caching
- **Strategies**: Translation caching, user profile caching, resume caching, API response caching, AI response caching
- **Implementation**: Cache warming, hit/miss metrics, compression, versioning, multi-level caching

### Database
- **Queries**: Performance monitoring, slow query logging, pagination, read replicas, query caching
- **Joins**: Index optimization, query analysis, denormalization, view creation, join optimization hints

