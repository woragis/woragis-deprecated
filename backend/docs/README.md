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

### [Development Guides](./development/)
Guides for developers working on the codebase:
- [Adding a New Domain](./development/adding-domain.md) - How to add a new domain to the server
- [Adding a New Worker](./development/adding-worker.md) - How to create a new worker service
- [Adding a New Service](./development/adding-service.md) - How to create a new microservice
- [Testing Patterns](./development/testing-patterns.md) - Testing conventions and patterns
- [Logging Conventions](./development/logging-conventions.md) - When and how to log
- [Error Handling Patterns](./development/error-handling.md) - Error handling best practices

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

## 🔗 Related Documentation

- [Backend TODO](../TODO.md) - Current tasks and implementation status
- [Testing Guide](../TESTING.md) - Testing instructions
- Component-specific READMEs in each component directory
