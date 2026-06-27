# Landing Page Multilanguage Support

## Overview
How multilanguage support is implemented in the Svelte landing page.

## Key Points

### i18n Implementation
- Similar to frontend but landing-specific
- JSON translation files in `locales/` folder
- 11 languages supported (en, pt-BR, fr, es, de, ru, ja, ko, zh-CN, el, la)
- Language switcher component

### Language Switcher Component
- Dropdown UI with flag emojis
- Shows current language flag
- Lists all available languages
- Click outside to close
- Updates language store on selection

### Language Store
- Writable Svelte store
- Persists to localStorage
- Browser language detection
- Reactive updates across components

### Translation Files
- One JSON file per language
- Nested key structure
- Parameter substitution (`{{param}}`)
- Fallback to English

### Language Detection
- localStorage preference first
- Browser language detection
- Language code mapping (pt → pt-BR, zh → zh-CN, etc.)
- Default to English

### Supported Languages
- English (🇺🇸)
- Portuguese BR (🇧🇷)
- French (🇫🇷)
- Spanish (🇪🇸)
- German (🇩🇪)
- Russian (🇷🇺)
- Japanese (🇯🇵)
- Korean (🇰🇷)
- Chinese Simplified (🇨🇳)
- Greek (🇬🇷)
- Latin (🏛️)

## Potential Improvements
- Add language-specific SEO (hreflang tags)
- Implement language-specific routes (/en/, /pt-BR/)
- Add language selection to initial visit flow
- Support language-specific content (not just translations)
- Add translation progress indicators
- Implement lazy loading of translation files
- Add translation completeness validation
- Support translation updates without rebuild
- Add language-specific analytics
- Implement A/B testing per language
- Support right-to-left (RTL) layouts
- Add language-specific fonts/typography
- Support regional variants (pt-PT vs pt-BR)

