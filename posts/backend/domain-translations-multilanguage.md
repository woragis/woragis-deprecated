# Translations Domain - Multilanguage Support

## Overview
How multilanguage support is implemented across the backend using the translations domain.

## Key Points

### Language Detection Middleware
- `LanguageMiddleware()`: Extracts language from requests
- Language detection priority:
  1. Query parameter: `?lang=` or `?language=`
  2. Accept-Language header
  3. Default: English (EN)
- Stores language in Fiber context for downstream handlers

### Supported Languages
- English (EN)
- Portuguese (PT-BR)
- French (FR)
- Spanish (ES)
- German (DE)
- Russian (RU)
- Japanese (JA)
- Korean (KO)
- Chinese Simplified (ZH-CN)
- Greek (EL)
- Latin (LA)

### Translation Architecture
- Entity-based translations (projects, skills, interests, posts, etc.)
- Field-level translations (title, description, content, etc.)
- Status tracking: pending, processing, completed, failed
- Queue-based async processing via translation worker

### Language Context Flow
1. Request arrives → LanguageMiddleware extracts language
2. Language stored in context: `translation_language`
3. Handlers retrieve language: `LanguageFromContext()`
4. Service layer uses language for:
   - Fetching translations from database
   - Requesting new translations if missing
   - Returning translated content

### Translation Enrichment
- `TranslationEnricher`: Enriches entities with translations
- Automatic translation fetching when missing
- Fallback to source language if translation not available
- Transparent to handlers (automatic enrichment)

### Integration Points
- All domain handlers use translation enricher
- Middleware applied to route groups (projects, skills, interests, posts, etc.)
- Database stores translations per entity + language
- AI service generates translations via translation worker

## Implementation Details

### Middleware Implementation
- `LanguageFromRequest()`: Parses Accept-Language header
- `parseAcceptLanguage()`: Handles quality values, language variants
- Context storage: `c.Locals(ContextLanguageKey, lang)`

### Translation Service
- `RequestTranslation()`: Creates translation job
- `GetTranslation()`: Retrieves existing translation
- `ProcessTranslationJob()`: Worker processes job, calls AI service

## Potential Improvements
- Add language preference per user (store in user profile)
- Implement translation caching (Redis)
- Add translation versioning (track changes to source)
- Support fallback chain (en → source language)
- Add translation quality scoring
- Implement batch translation requests
- Add translation memory (reuse similar translations)
- Support partial translations (some fields translated)
- Add translation review/approval workflow
- Implement real-time translation updates
- Add translation analytics (most requested languages)
- Support regional variants (pt-BR vs pt-PT)

