# ADR-001: Why Go for Workers, Python for Services

## Context
We need to choose programming languages for workers and services. Workers handle async job processing, services handle AI/LLM operations.

## Decision
- **Workers**: Go (email, whatsapp, translation, resume, job-application)
- **Services**: Python (AI service, creative service)

## Rationale

### Go for Workers
- **Performance**: Fast compilation, low memory footprint
- **Concurrency**: Excellent goroutine support for concurrent job processing
- **Consistency**: All workers in same language simplifies maintenance
- **Deployment**: Single binary, easy containerization
- **Ecosystem**: Good RabbitMQ, database libraries

### Python for Services
- **Ecosystem**: Rich AI/ML libraries (LangChain, OpenAI SDK)
- **Rapid Development**: Faster iteration for AI features
- **Community**: Large AI/ML community, more examples
- **Flexibility**: Easy to integrate new AI providers

## Consequences

### Pros
- Right tool for right job
- Performance where needed (workers)
- Flexibility where needed (services)
- Technology diversity

### Cons
- Need to maintain two languages
- Different deployment strategies
- Team needs both skills
- Code sharing limitations

## Status
Accepted - 2024-01-15

## Alternatives Considered
- All Go: Missing AI ecosystem
- All Python: Performance concerns for workers
- Node.js: Not ideal for either use case
