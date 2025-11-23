package translations

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	ContextLanguageKey = "translation_language"
)

// LanguageFromRequest extracts the preferred language from the request.
// Checks in order:
// 1. Query parameter "lang" or "language"
// 2. Accept-Language header
// 3. Defaults to English
func LanguageFromRequest(c *fiber.Ctx) Language {
	// Check query parameter first
	if lang := c.Query("lang"); lang != "" {
		if isValidLanguage(Language(lang)) {
			return Language(lang)
		}
	}
	if lang := c.Query("language"); lang != "" {
		if isValidLanguage(Language(lang)) {
			return Language(lang)
		}
	}

	// Check Accept-Language header
	acceptLang := c.Get("Accept-Language")
	if acceptLang != "" {
		lang := parseAcceptLanguage(acceptLang)
		if lang != "" && isValidLanguage(lang) {
			return lang
		}
	}

	// Default to English
	return LanguageEN
}

// LanguageMiddleware extracts and stores the language in context.
func LanguageMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := LanguageFromRequest(c)
		c.Locals(ContextLanguageKey, lang)
		return c.Next()
	}
}

// LanguageFromContext retrieves the language from Fiber context.
func LanguageFromContext(c *fiber.Ctx) Language {
	if lang, ok := c.Locals(ContextLanguageKey).(Language); ok {
		return lang
	}
	return LanguageEN
}

// parseAcceptLanguage parses the Accept-Language header and returns the best match.
func parseAcceptLanguage(acceptLang string) Language {
	// Simple parsing: take the first language code
	parts := strings.Split(acceptLang, ",")
	if len(parts) > 0 {
		langPart := strings.TrimSpace(parts[0])
		// Remove quality values (e.g., "pt-BR;q=0.9" -> "pt-BR")
		langPart = strings.Split(langPart, ";")[0]
		langPart = strings.TrimSpace(langPart)

		// Map common language codes
		langPartLower := strings.ToLower(langPart)
		if strings.HasPrefix(langPartLower, "pt") {
			return LanguagePTBR
		}
		if strings.HasPrefix(langPartLower, "fr") {
			return LanguageFR
		}
		if strings.HasPrefix(langPartLower, "es") {
			return LanguageES
		}
		if strings.HasPrefix(langPartLower, "de") {
			return LanguageDE
		}
		if strings.HasPrefix(langPartLower, "ru") {
			return LanguageRU
		}
		if strings.HasPrefix(langPartLower, "ja") {
			return LanguageJA
		}
		if strings.HasPrefix(langPartLower, "ko") {
			return LanguageKO
		}
		if strings.HasPrefix(langPartLower, "zh") {
			return LanguageZHCN
		}
		if strings.HasPrefix(langPartLower, "el") {
			return LanguageEL
		}
		if strings.HasPrefix(langPartLower, "la") {
			return LanguageLA
		}
	}

	return LanguageEN
}

