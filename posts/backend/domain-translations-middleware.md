# Translations Domain - Language Detection Middleware

## Overview
How the language detection middleware extracts and manages request languages.

## Key Points

### Language Detection Priority
1. Query parameter: `?lang=` or `?language=`
2. Accept-Language header
3. Default: English (EN)

### Middleware Implementation
- LanguageMiddleware() extracts language from request
- Stores language in Fiber context (ContextLanguageKey)
- Transparent to handlers (automatic)
- Applied to route groups

### Language Parsing
- parseAcceptLanguage handles Accept-Language header
- Parses quality values (e.g., `pt-BR;q=0.9`)
- Maps language codes to supported languages
- Handles language variants (pt → pt-BR)

### Supported Languages
- English, Portuguese (BR), French, Spanish, German
- Russian, Japanese, Korean, Chinese (Simplified)
- Greek, Latin

### Context Storage
- Language stored in Fiber context
- LanguageFromContext retrieves language
- Used by TranslationEnricher
- Fallback to English if not found

### Integration
- Applied to route groups (projects, skills, interests, posts, etc.)
- Works with TranslationEnricher
- Automatic translation fetching based on language

## Potential Improvements
- Add language preference per user (database)
- Support language variants (pt-PT vs pt-BR)
- Add language detection caching
- Implement language fallback chains
- Add language detection metrics
- Support browser language detection
- Add language switching without reload
- Implement language cookie persistence
- Support regional language detection
- Add language detection API
- Support custom language mappings
- Add language detection debugging

