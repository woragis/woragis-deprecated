# Database Design: PostgreSQL with JSONB

## Overview
How we use PostgreSQL with JSONB columns for flexible schema design while maintaining relational integrity.

## Key Points

### Database Choice
- PostgreSQL 15+
- JSONB for flexible fields
- Relational tables for structured data
- GORM for ORM (Go)

### JSONB Usage
- Translation fields (map[string]string)
- Flexible metadata storage
- Configuration storage
- Dynamic field storage

## Implementation Details

### JSONB Columns
```go
type Translation struct {
    ID        uuid.UUID
    EntityType string
    EntityID   uuid.UUID
    Language   string
    Fields     datatypes.JSON `gorm:"type:jsonb"` // JSONB column
    Status     string
}
```

### Benefits
- Flexible schema (add fields without migration)
- Fast queries (JSONB indexed)
- Type safety (Go structs)
- Relational integrity (foreign keys)

### Trade-offs
- Less strict schema (validation needed)
- Query complexity (JSONB queries)
- Migration complexity (JSONB changes)

## Use Cases

### Translations
- Store translated fields as JSONB
- Query by language
- Update individual fields
- Maintain relational links

### Metadata
- Store flexible metadata
- No schema changes needed
- Easy to extend

## Query Patterns

### JSONB Queries
```sql
-- Query JSONB field
SELECT * FROM translations 
WHERE fields->>'name' = 'Project Name';

-- Update JSONB field
UPDATE translations 
SET fields = jsonb_set(fields, '{name}', '"New Name"')
WHERE id = '...';
```

## Lessons Learned
- JSONB powerful for flexible data
- Need validation layer
- Indexes important for performance
- GORM handles JSONB well

## Future Improvements
- JSONB query optimization
- Indexing strategies
- Validation layer
- Migration tooling
