# ADR-008: Database: PostgreSQL vs Alternatives

## Context
We need a relational database with JSON support for flexible schema. Options: PostgreSQL, MySQL, MongoDB.

## Decision
PostgreSQL 15+ with JSONB columns for flexible fields, relational tables for structured data.

## Rationale

### PostgreSQL
- **JSONB Support**: Native JSONB with indexing
- **Relational**: ACID transactions, foreign keys
- **Performance**: Excellent query performance
- **Ecosystem**: Good Go (GORM) and Python support

### JSONB Columns
- **Flexibility**: Add fields without migration
- **Performance**: Indexed JSONB queries
- **Type Safety**: Go structs with JSONB

## Consequences

### Pros
- Best of both worlds (relational + flexible)
- Excellent performance
- Good tooling
- JSONB indexing

### Cons
- Need to validate JSONB
- Query complexity
- Migration complexity

## Status
Accepted - 2024-01-15

## Alternatives Considered
- MySQL: Less JSON support
- MongoDB: No relational features
- Both: Too complex
