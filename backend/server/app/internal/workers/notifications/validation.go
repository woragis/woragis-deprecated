package notifications

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateReportEnvelope validates a report envelope before processing
func ValidateReportEnvelope(envelope *ReportEnvelope) error {
	// Validate user ID (required, UUID)
	if envelope.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if err := validation.ValidateUUID(envelope.UserID); err != nil {
		return fmt.Errorf("user_id: %w", err)
	}

	// Validate destination (required, must be valid email)
	if envelope.Destination == "" {
		return fmt.Errorf("destination is required")
	}
	if err := validation.ValidateEmail(envelope.Destination); err != nil {
		return fmt.Errorf("destination: %w", err)
	}

	// Validate subject (required, 1-200 chars)
	if err := validation.ValidateString(envelope.Subject, 1, 200, "subject"); err != nil {
		return fmt.Errorf("subject: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(envelope.Subject); err != nil {
		return fmt.Errorf("subject: %w", err)
	}
	if err := validation.ValidateNoXSS(envelope.Subject); err != nil {
		return fmt.Errorf("subject: %w", err)
	}

	// Validate text message (optional, but if provided, validate length)
	if envelope.TextMessage != "" {
		if err := validation.ValidateString(envelope.TextMessage, 1, 100000, "text_message"); err != nil {
			return fmt.Errorf("text_message: %w", err)
		}
	}

	// Validate HTML message (optional, but if provided, validate length)
	if envelope.HTMLMessage != "" {
		if err := validation.ValidateString(envelope.HTMLMessage, 1, 100000, "html_message"); err != nil {
			return fmt.Errorf("html_message: %w", err)
		}
	}

	// At least one message body must be provided
	if envelope.TextMessage == "" && envelope.HTMLMessage == "" {
		return fmt.Errorf("at least one of text_message or html_message is required")
	}

	return nil
}

