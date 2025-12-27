package email

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateMessage validates an email message before sending
func ValidateMessage(msg Message) error {
	// Validate recipient (required, must be valid email)
	if msg.To == "" {
		return fmt.Errorf("to is required")
	}
	if err := validation.ValidateEmail(msg.To); err != nil {
		return fmt.Errorf("to: %w", err)
	}

	// Validate subject (required, 1-200 chars)
	if err := validation.ValidateString(msg.Subject, 1, 200, "subject"); err != nil {
		return fmt.Errorf("subject: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(msg.Subject); err != nil {
		return fmt.Errorf("subject: %w", err)
	}
	if err := validation.ValidateNoXSS(msg.Subject); err != nil {
		return fmt.Errorf("subject: %w", err)
	}

	// Validate text body (optional, but if provided, validate length)
	if msg.TextBody != "" {
		if err := validation.ValidateString(msg.TextBody, 1, 100000, "text_body"); err != nil {
			return fmt.Errorf("text_body: %w", err)
		}
	}

	// Validate HTML body (optional, but if provided, validate length)
	if msg.HTMLBody != "" {
		if err := validation.ValidateString(msg.HTMLBody, 1, 100000, "html_body"); err != nil {
			return fmt.Errorf("html_body: %w", err)
		}
	}

	// At least one body must be provided
	if msg.TextBody == "" && msg.HTMLBody == "" {
		return fmt.Errorf("at least one of text_body or html_body is required")
	}

	return nil
}

