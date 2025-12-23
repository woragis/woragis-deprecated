# Backend Technical Documentation

Welcome to the backend technical documentation. This directory contains comprehensive documentation for the Woragis backend architecture, decisions, operations, and development guides.

## 📚 Documentation Structure

### [Architecture](./architecture/)
High-level architecture documentation and system overviews:
- [System Overview](./architecture/system-overview.md) - High-level system architecture and component interaction
- [Observability Overview](./architecture/OBSERVABILITY_OVERVIEW.md) - Logging, metrics, and tracing
- [Resilience Overview](./architecture/RESILIENCE_OVERVIEW.md) - Fault tolerance and failure handling
- [Scalability Overview](./architecture/SCALABILITY_OVERVIEW.md) - Scaling strategies and performance
- [Cost Management Overview](./architecture/COST_MANAGEMENT_OVERVIEW.md) - Infrastructure cost optimization
- [Disaster Recovery Overview](./architecture/DISASTER_RECOVERY_OVERVIEW.md) - Backup and recovery procedures
- [Testing Overview](./architecture/TESTING_OVERVIEW.md) - Testing strategy and patterns
- [Structured Logging Overview](./architecture/STRUCTURED_LOGGING_OVERVIEW.md) - Logging implementation
- [Health Checks Overview](./architecture/HEALTH_CHECKS_OVERVIEW.md) - Health check implementation

### [Architecture Decision Records (ADRs)](./adr/)
Recorded architectural decisions and their rationale:
- [ADR-001: RabbitMQ with Redis Fallback](./adr/001-rabbitmq-redis-fallback.md)
- [ADR-002: Standalone Workers Architecture](./adr/002-standalone-workers.md)
- [ADR-003: Structured Logging Implementation](./adr/003-structured-logging.md)
- [ADR-004: Translation Worker Architecture](./adr/004-translation-worker.md)
- [ADR-005: Health Checks Implementation Strategy](./adr/005-health-checks.md)

### [Runbooks](./runbooks/)
Operational procedures and troubleshooting guides:
- [Monitoring Dead Letter Queues](./runbooks/monitoring-dlq.md) - How to monitor and handle DLQ messages
- [Monitoring and Alerting](./runbooks/monitoring-alerting.md) - Monitoring setup and alerting configuration
- [Backup and Restore](./runbooks/backup-restore.md) - Backup and disaster recovery procedures
- [Deploying Services and Workers](./runbooks/deploying-services.md) - Deployment procedures
- [Troubleshooting Common Issues](./runbooks/troubleshooting.md) - Common problems and solutions

### [Components](./components/)
Detailed documentation for each backend component:
- [Server](./components/server.md) - Main API server (Go)
- [Email Worker](./components/email-worker.md) - Email processing worker (Go)
- [WhatsApp Worker](./components/whatsapp-worker.md) - WhatsApp messaging worker (Go)
- [Translation Worker](./components/translation-worker.md) - Translation processing worker (Go)
- [Resume Worker](./components/resume-worker.md) - Resume generation worker (Python)
- [Job Application Worker](./components/job-application-worker.md) - Job application automation worker (Node.js)
- [AI Service](./components/ai-service.md) - AI/LLM service (Python/FastAPI)
- [Creative Service](./components/creative-service.md) - Creative content generation service (Python/FastAPI)
- [Docs Service](./components/docs-service.md) - Documentation serving service (Python/FastAPI)

### [API Documentation](./api/)
API endpoint documentation:
- [Server API](./api/server-api.md) - REST API endpoints
- [AI Service API](./api/ai-service-api.md) - AI service endpoints
- [Creative Service API](./api/creative-service-api.md) - Creative service endpoints
- [Docs Service API](./api/docs-service-api.md) - Documentation service endpoints

### [Deployment Guides](./deployment/)
Production deployment and configuration guides:
- [Configuration](./deployment/configuration.md) - Environment variables and configuration
- [Docker Setup](./deployment/docker-setup.md) - Docker deployment guide
- [Secrets Management](./deployment/secrets-management.md) - Managing secrets securely
- [SSL/TLS Configuration](./deployment/ssl-tls-configuration.md) - SSL/TLS setup
- [Authentication & Authorization](./deployment/authentication-authorization.md) - Auth configuration
- [Input Validation](./deployment/input-validation.md) - Input validation setup
- [Performance Optimization](./deployment/performance-optimization.md) - Performance tuning
- [Circuit Breaker Implementation](./deployment/circuit-breaker-implementation.md) - Circuit breaker setup

### [Development Guides](./development/)
Guides for developers working on the codebase:
- [Setup Guide](./development/setup-guide.md) - Local development environment setup
- [Contributing](./development/contributing.md) - Contribution guidelines
- [Coding Standards](./development/coding-standards.md) - Code style and standards
- [Code Review Process](./development/code-review-process.md) - Code review workflow
- [Adding a New Domain](./development/adding-domain.md) - How to add a new domain to the server
- [Adding a New Worker](./development/adding-worker.md) - How to create a new worker service
- [Adding a New Service](./development/adding-service.md) - How to create a new microservice
- [Testing Patterns](./development/testing-patterns.md) - Testing conventions and patterns
- [Logging Conventions](./development/logging-conventions.md) - When and how to log
- [Error Handling Patterns](./development/error-handling.md) - Error handling best practices

### [Planning Documents](./PLANNING/)
Planning documents and implementation plans:
- [01: Logging Aggregation Plan](./PLANNING/01-logging-aggregation-plan.md) - Logging infrastructure planning
- [02: Documentation Plan](./PLANNING/02-documentation-plan.md) - Documentation strategy
- [03: Production Readiness Plan](./PLANNING/03-production-readiness-plan.md) - Production deployment planning
- [04: Development Workflow Plan](./PLANNING/04-development-workflow-plan.md) - Development process planning
- [05: Distributed Tracing Plan](./PLANNING/05-distributed-tracing-plan.md) - Distributed tracing implementation
- [Circuit Breaker Plan](./PLANNING/CIRCUIT_BREAKER_PLAN.md) - Circuit breaker implementation plan
- [Documentation Audit](./PLANNING/DOCUMENTATION_AUDIT.md) - Documentation audit and gap analysis

**Note:** Completed implementation summaries and status reports are archived in [PLANNING/archive/](./PLANNING/archive/) for historical reference.

## 🚀 Quick Start

1. **New to the project?** Start with [System Overview](./architecture/system-overview.md)
2. **Understanding decisions?** Check [Architecture Decision Records](./adr/)
3. **Deploying?** See [Deploying Services and Workers](./runbooks/deploying-services.md)
4. **Adding features?** Read [Development Guides](./development/)
5. **Troubleshooting?** Check [Troubleshooting Guide](./runbooks/troubleshooting.md)

## 📖 Documentation Standards

- **ADRs**: Follow the standard ADR format (Context, Decision, Consequences)
- **Runbooks**: Step-by-step, actionable procedures
- **Component Docs**: Architecture, configuration, usage, examples
- **API Docs**: Endpoint descriptions, request/response examples, error codes
- **Development Guides**: Clear instructions with code examples

## 🔄 Keeping Documentation Updated

- Update ADRs when making significant architectural decisions
- Update runbooks when procedures change
- Update component docs when features change
- Review documentation quarterly for accuracy

## 📝 Contributing

When adding new documentation:
1. Follow the existing structure and format
2. Use Markdown with clear headings
3. Include code examples where relevant
4. Add diagrams (Mermaid) for complex concepts
5. Link related documents
6. Update this README if adding new sections

## 📁 Archive Folders

- **[PLANNING/archive/](./PLANNING/archive/)** - Historical implementation summaries and status reports
- **[deployment/archive/](./deployment/archive/)** - Archived deployment implementation summaries

These folders contain completed implementation documentation for historical reference. Active planning documents remain in their respective parent folders.

## 🔗 Related Documentation

- [Backend TODO](../TODO.md) - Current tasks and implementation status
- [Testing Guide](../TESTING.md) - Testing instructions
- [Logging Format Specification](./LOGGING_FORMAT_SPECIFICATION.md) - Structured logging format specification
- Component-specific READMEs in each component directory
