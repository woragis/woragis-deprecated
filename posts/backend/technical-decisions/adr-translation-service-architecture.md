# ADR-004: Translation Service: Go + External APIs

## Context
We need a translation service. Options: embedded in server, standalone worker, external service.

## Decision
Standalone Go worker that calls external translation APIs (Google Translate, DeepL, LibreTranslate) and writes directly to database.

## Rationale

### Standalone Worker
- **Performance**: Go's performance for concurrent translations
- **Scalability**: Independent scaling
- **Fault Isolation**: Translation failures don't affect server

### External APIs
- **Quality**: Professional translation services
- **Maintenance**: No need to maintain translation models
- **Cost**: Pay-per-use, no infrastructure

### Direct DB Writes
- **Simplicity**: No need for server to receive translations
- **Performance**: Worker writes directly, no round-trip
- **Consistency**: Single source of truth

## Consequences

### Pros
- High performance (Go)
- Professional translations
- Independent scaling
- Simple architecture

### Cons
- External API dependency
- API costs
- Rate limiting concerns
- Network latency

## Status
Accepted - 2024-01-15

## Alternatives Considered
- Embedded in server: Performance and scaling concerns
- Python worker: Go better for concurrent processing
- Translation back to server: Unnecessary complexity
