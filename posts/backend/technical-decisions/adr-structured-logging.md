# ADR-003: Structured Logging: JSON Format

## Context
We need logging that's both human-readable (development) and machine-parseable (production, log aggregation).

## Decision
Use structured JSON logging with configurable format (JSON for production, text for development).

## Rationale

### JSON Format (Production)
- **Machine Parseable**: Easy log aggregation (ELK, Grafana)
- **Structured Fields**: Consistent field names
- **Queryable**: Search and filter by fields
- **Trace ID**: Distributed tracing support

### Text Format (Development)
- **Human Readable**: Easier to read during development
- **Faster Debugging**: Quick visual inspection
- **Less Verbose**: Cleaner output

## Implementation
- Go: `log/slog` with JSON/text handlers
- Python: `structlog` with JSON/text formatters
- Node.js: Custom logger with JSON/text output

## Consequences

### Pros
- Best of both worlds (readable + parseable)
- Easy log aggregation
- Distributed tracing support
- Consistent across services

### Cons
- More verbose than plain text
- Need log aggregation tools
- Configuration complexity

## Status
Accepted - 2024-01-15

## Alternatives Considered
- Plain text only: Hard to aggregate
- JSON only: Hard to read in development
- Binary format: Too complex
