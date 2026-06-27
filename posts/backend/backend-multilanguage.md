# Backend Multilanguage Support

## Overview
Complete multilanguage system across the backend (combines translations domain + middleware).

## Key Points

### System Components
- Translations Domain: Core translation management
- Language Middleware: Request language detection
- Translation Enricher: Automatic translation fetching
- Translation Worker: Async translation processing
- AI Service Integration: Translation generation

### Language Detection Flow
1. Request arrives with language preferences
2. LanguageMiddleware extracts language (query param → header → default)
3. Language stored in Fiber context
4. Handlers use language from context
5. TranslationEnricher fetches translations
6. Translated content returned

### Translation Workflow
1. Entity created/updated in source language
2. Translation job enqueued for target languages
3. Worker processes translation jobs
4. AI service generates translations
5. Translations stored in database
6. Future requests return translated content

### Integration Points
- All domain handlers support translations
- Middleware applied to route groups
- Automatic enrichment transparent to handlers
- Database stores translations per entity + language

### Supported Languages
- Matches frontend languages (11 languages)
- Same language codes (en, pt-BR, fr, es, de, ru, ja, ko, zh-CN, el, la)

## Potential Improvements
- Add translation caching layer
- Implement translation fallback chains
- Add translation versioning
- Support translation review/approval
- Implement translation quality scoring
- Add translation memory/reuse
- Support translation templates
- Add translation analytics
- Implement translation batch operations
- Add translation export/import
- Support translation collaboration
- Add translation API for external tools

