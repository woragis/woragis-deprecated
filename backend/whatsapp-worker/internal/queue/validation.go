package queue

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateWhatsAppEnvelope validates a WhatsAppEnvelope message.
func ValidateWhatsAppEnvelope(envelope WhatsAppEnvelope) error {
	// Validate UserID
	if envelope.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if err := validateUUID(envelope.UserID); err != nil {
		return fmt.Errorf("user ID: %w", err)
	}

	// Validate Destination (phone number)
	if envelope.Destination == "" {
		return fmt.Errorf("destination phone number is required")
	}
	if err := validatePhoneNumber(envelope.Destination); err != nil {
		return fmt.Errorf("destination: %w", err)
	}

	// Validate TextMessage
	if envelope.TextMessage == "" {
		return fmt.Errorf("text message is required")
	}
	if err := validateString(envelope.TextMessage, 1, 4096, "text_message"); err != nil {
		return fmt.Errorf("text_message: %w", err)
	}
	if err := validateNoSQLInjection(envelope.TextMessage); err != nil {
		return fmt.Errorf("text_message: %w", err)
	}
	if err := validateNoXSS(envelope.TextMessage); err != nil {
		return fmt.Errorf("text_message: %w", err)
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

// validatePhoneNumber validates phone number format (international format)
func validatePhoneNumber(phone string) error {
	// Remove common separators
	cleaned := regexp.MustCompile(`[\s\-\(\)]`).ReplaceAllString(phone, "")
	
	// Check if it starts with + and has digits
	if !strings.HasPrefix(cleaned, "+") {
		return fmt.Errorf("phone number must start with + (international format)")
	}
	
	// Check length (minimum 8 digits after +, maximum 15)
	digits := strings.TrimPrefix(cleaned, "+")
	if len(digits) < 8 || len(digits) > 15 {
		return fmt.Errorf("phone number must have 8-15 digits after country code")
	}
	
	// Check if all characters after + are digits
	digitRegex := regexp.MustCompile(`^\d+$`)
	if !digitRegex.MatchString(digits) {
		return fmt.Errorf("phone number must contain only digits after country code")
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

