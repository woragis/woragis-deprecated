# Frontend Multilanguage Support

## Overview
How multilanguage support is implemented in the Svelte frontend application.

## Key Points

### i18n Architecture
- Svelte stores for language state management
- JSON translation files per language
- Reactive translation functions
- Browser language detection
- localStorage persistence

### Translation Storage
- JSON files for each language (en, pt-BR, fr, es, de, ru, ja, ko, zh-CN, el, la)
- Nested key structure (e.g., `hero.title`)
- Parameter substitution support (`{{param}}`)
- Fallback to English if translation missing

### Language Store
- Writable store for current language
- Reactive to language changes
- Persists to localStorage
- Initializes from localStorage or browser language

### Translation Functions
- `t(key, params)`: Get translation for key with optional parameters
- `$translationsStore`: Reactive store for translations in templates
- `getTranslation()`: Internal translation getter with fallback

### Language Detection
- Checks localStorage first
- Falls back to browser language (`navigator.language`)
- Maps browser language codes to supported languages
- Defaults to English

### Integration Points
- API calls include language parameter
- Resume download/preview with language selection
- Date formatting with locale

## Potential Improvements
- Add language preference sync with backend (user profile)
- Implement lazy loading of translation files
- Add translation loading states
- Support translation hot-reloading in development
- Add translation key validation
- Implement translation completeness checks
- Add translation context/locale formatting
- Support pluralization rules
- Add translation interpolation for complex strings
- Implement translation caching
- Support translation versioning
- Add translation analytics (most used keys)
- Support RTL languages
- Add translation memory/proofreading tools

