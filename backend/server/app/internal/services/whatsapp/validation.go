package whatsapp

import (
	"fmt"
	"regexp"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidatePhoneNumber validates a phone number format
func ValidatePhoneNumber(phoneNumber string) error {
	if phoneNumber == "" {
		return fmt.Errorf("phone_number is required")
	}

	// Basic phone number validation (E.164 format or local format)
	// Remove common separators
	cleaned := regexp.MustCompile(`[\s\-\(\)]`).ReplaceAllString(phoneNumber, "")
	
	// Check length (minimum 7, maximum 15 for E.164)
	if len(cleaned) < 7 {
		return fmt.Errorf("phone_number: too short (minimum 7 digits)")
	}
	if len(cleaned) > 15 {
		return fmt.Errorf("phone_number: too long (maximum 15 digits)")
	}

	// Check for only digits and optional + prefix
	matched, _ := regexp.MatchString(`^\+?[1-9]\d{6,14}$`, cleaned)
	if !matched {
		return fmt.Errorf("phone_number: invalid format (must be E.164 format or local format)")
	}

	return nil
}

// ValidateMessage validates a WhatsApp message
func ValidateMessage(message string) error {
	if message == "" {
		return fmt.Errorf("message is required")
	}

	// Validate length (1-4096 chars for WhatsApp)
	if err := validation.ValidateString(message, 1, 4096, "message"); err != nil {
		return fmt.Errorf("message: %w", err)
	}

	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(message); err != nil {
		return fmt.Errorf("message: %w", err)
	}
	if err := validation.ValidateNoXSS(message); err != nil {
		return fmt.Errorf("message: %w", err)
	}

	return nil
}

