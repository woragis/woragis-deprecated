# ADR-006: Testing Strategies: Go, Python, Node.js

## Context
We have three languages (Go, Python, Node.js) and need consistent testing strategies while using language-appropriate tools.

## Decision
Use language-standard testing frameworks with consistent patterns: Go (`testing`), Python (`pytest`), Node.js (`Jest`).

## Rationale

### Language-Standard Tools
- **Go**: `testing` package (standard library)
- **Python**: `pytest` (industry standard)
- **Node.js**: `Jest` (popular, well-supported)

### Consistent Patterns
- **Coverage Target**: 70% minimum
- **Test Structure**: Unit tests alongside source, integration tests separate
- **Docker Support**: `Dockerfile.test` for all components
- **Makefile**: Standardized `make test`, `make test-cov`

## Consequences

### Pros
- Right tool for each language
- Consistent patterns
- Easy to maintain
- Good coverage

### Cons
- Need to know three testing frameworks
- Different syntax
- Tooling differences

## Status
Accepted - 2024-01-15

## Alternatives Considered
- Single testing framework: Not practical across languages
- No testing: Unacceptable
- Different patterns per language: Too inconsistent
