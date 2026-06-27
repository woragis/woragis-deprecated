# Translations Domain - Database Queries & Joins

## Overview
Database query patterns used in the translations domain.

## Key Points

### Query Patterns

#### Translation Queries
- GetTranslationByEntity: Get translation for entity + language
- ListTranslations: List translations with filters (entity_type, entity_id, language)
- CreateTranslation: Create new translation record
- UpdateTranslation: Update existing translation

#### Entity Queries
- Fetch entities from database for translation
- Generic entity fetching (projects, skills, etc.)
- Extract translatable fields from entities

### Translation Storage
- Translations stored per entity + language
- Field-level translations (JSON map)
- Status tracking (pending, processing, completed, failed)
- Metadata: created_at, updated_at

### Join Patterns
- No explicit joins currently
- Entity fetching separate from translation fetching
- TranslationEnricher combines entities with translations

### Filtering
- By entity_type
- By entity_id
- By language
- By status

## Potential Improvements
- Add indexes for frequently queried fields (entity_type, entity_id, language)
- Optimize entity + translation joins
- Add query result caching
- Implement query result pagination
- Add query logging for slow queries
- Optimize translation lookups
- Add full-text search for translations
- Support translation search across languages
- Implement query result streaming
- Add query explain plans
- Support translation bulk operations
- Add translation statistics queries

