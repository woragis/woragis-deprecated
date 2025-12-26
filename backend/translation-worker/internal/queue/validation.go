package queue

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateTranslationJob validates a TranslationJob message.
func ValidateTranslationJob(job TranslationJob) error {
	// Validate ID
	if job.ID == "" {
		return fmt.Errorf("job ID is required")
	}
	if err := validateUUID(job.ID); err != nil {
		return fmt.Errorf("job ID: %w", err)
	}

	// Validate EntityType
	if job.EntityType == "" {
		return fmt.Errorf("entity type is required")
	}
	validEntityTypes := map[string]bool{
		"testimonial":        true,
		"project":           true,
		"certification":     true,
		"skill":             true,
		"idea":              true,
		"job_application":   true,
		"resume":            true,
		"technical_writing": true,
		"case_study":        true,
		"post":              true,
	}
	if !validEntityTypes[strings.ToLower(job.EntityType)] {
		return fmt.Errorf("invalid entity type: %s", job.EntityType)
	}

	// Validate EntityID
	if job.EntityID == "" {
		return fmt.Errorf("entity ID is required")
	}
	if err := validateUUID(job.EntityID); err != nil {
		return fmt.Errorf("entity ID: %w", err)
	}

	// Validate Language
	if job.Language == "" {
		return fmt.Errorf("target language is required")
	}
	if err := validateLanguageCode(job.Language); err != nil {
		return fmt.Errorf("target language: %w", err)
	}

	// Validate Fields
	if len(job.Fields) == 0 {
		return fmt.Errorf("fields to translate cannot be empty")
	}
	if len(job.Fields) > 50 {
		return fmt.Errorf("too many fields (max 50)")
	}
	for i, field := range job.Fields {
		if err := validateString(field, 1, 100, fmt.Sprintf("fields[%d]", i)); err != nil {
			return fmt.Errorf("fields[%d]: %w", i, err)
		}
		if err := validateNoSQLInjection(field); err != nil {
			return fmt.Errorf("fields[%d]: %w", i, err)
		}
		if err := validateNoXSS(field); err != nil {
			return fmt.Errorf("fields[%d]: %w", i, err)
		}
	}

	// Validate SourceText if provided
	if len(job.SourceText) > 50 {
		return fmt.Errorf("too many source text entries (max 50)")
	}
	for key, text := range job.SourceText {
		if err := validateString(key, 1, 100, fmt.Sprintf("sourceText key '%s'", key)); err != nil {
			return fmt.Errorf("sourceText key '%s': %w", key, err)
		}
		if err := validateString(text, 1, 10000, fmt.Sprintf("sourceText value for '%s'", key)); err != nil {
			return fmt.Errorf("sourceText value for '%s': %w", key, err)
		}
		if err := validateNoSQLInjection(text); err != nil {
			return fmt.Errorf("sourceText value for '%s': %w", key, err)
		}
		if err := validateNoXSS(text); err != nil {
			return fmt.Errorf("sourceText value for '%s': %w", key, err)
		}
	}

	return nil
}

// validateUUID validates UUID format
func validateUUID(value string) error {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(strings.ToLower(value)) {
		return fmt.Errorf("must be a valid UUID")
	}
	return nil
}

// validateString validates string length
func validateString(value string, minLength, maxLength int, fieldName string) error {
	if len(value) < minLength {
		return fmt.Errorf("too short (minimum %d characters)", minLength)
	}
	if len(value) > maxLength {
		return fmt.Errorf("too long (maximum %d characters)", maxLength)
	}
	return nil
}

// validateLanguageCode validates ISO 639-1 language code
func validateLanguageCode(language string) error {
	if len(language) < 2 || len(language) > 10 {
		return fmt.Errorf("must be 2-10 characters (ISO 639-1 or locale code)")
	}
	// Allow language codes like "en", "pt-BR", "zh-CN"
	langRegex := regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)
	if !langRegex.MatchString(strings.ToLower(language)) {
		return fmt.Errorf("invalid language code format")
	}
	return nil
}

// validateNoSQLInjection checks for potential SQL injection patterns
func validateNoSQLInjection(value string) error {
	dangerousPatterns := []string{
		`(?i)\b(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|EXEC|EXECUTE)\b`,
		`(--|#|/\*|\*/)`,
		`(?i)\b(UNION|OR|AND)\s+\d+\s*=\s*\d+`,
		`('|(\\')|(--)|(;)|(\|)|(\*))`,
	}
	for _, pattern := range dangerousPatterns {
		matched, _ := regexp.MatchString(pattern, value)
		if matched {
			return fmt.Errorf("contains potentially dangerous content")
		}
	}
	return nil
}

// validateNoXSS checks for potential XSS patterns
func validateNoXSS(value string) error {
	dangerousPatterns := []string{
		`(?i)<script[^>]*>`,
		`(?i)javascript:`,
		`(?i)on\w+\s*=`,
		`(?i)<iframe[^>]*>`,
		`(?i)<object[^>]*>`,
		`(?i)<embed[^>]*>`,
	}
	for _, pattern := range dangerousPatterns {
		matched, _ := regexp.MatchString(pattern, value)
		if matched {
			return fmt.Errorf("contains potentially dangerous content")
		}
	}
	return nil
}

