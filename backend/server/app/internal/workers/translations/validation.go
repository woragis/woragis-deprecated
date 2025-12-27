package translations

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
)

// ValidateTranslationJob validates a translation job before processing
func ValidateTranslationJob(job *translationsdomain.TranslationJob) error {
	// Validate job ID (required, UUID)
	if job.ID == "" {
		return fmt.Errorf("job id is required")
	}
	if err := validation.ValidateUUID(job.ID); err != nil {
		return fmt.Errorf("job id: %w", err)
	}

	// Validate entity ID (required, UUID)
	if job.EntityID == "" {
		return fmt.Errorf("entity_id is required")
	}
	if err := validation.ValidateUUID(job.EntityID); err != nil {
		return fmt.Errorf("entity_id: %w", err)
	}

	// Validate entity type (required)
	if job.EntityType == "" {
		return fmt.Errorf("entity_type is required")
	}
	validEntityTypes := []string{"post", "project", "technical_writing", "case_study"}
	isValid := false
	for _, validType := range validEntityTypes {
		if string(job.EntityType) == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("entity_type: must be one of: post, project, technical_writing, case_study")
	}

	// Validate language (required, ISO 639-1 format)
	if job.Language == "" {
		return fmt.Errorf("language is required")
	}
	if len(string(job.Language)) != 2 {
		return fmt.Errorf("language: must be exactly 2 characters (ISO 639-1 code)")
	}

	// Validate fields (optional, but if provided, validate each)
	if len(job.Fields) > 50 {
		return fmt.Errorf("fields: too many fields (maximum 50)")
	}
	for i, field := range job.Fields {
		if err := validation.ValidateString(field, 1, 100, fmt.Sprintf("fields[%d]", i)); err != nil {
			return fmt.Errorf("fields[%d]: %w", i, err)
		}
	}

	// Validate source text (optional, but if provided, validate each value)
	if len(job.SourceText) > 50 {
		return fmt.Errorf("source_text: too many entries (maximum 50)")
	}
	for key, value := range job.SourceText {
		if err := validation.ValidateString(key, 1, 100, "source_text key"); err != nil {
			return fmt.Errorf("source_text key %q: %w", key, err)
		}
		if err := validation.ValidateString(value, 1, 100000, fmt.Sprintf("source_text[%s]", key)); err != nil {
			return fmt.Errorf("source_text[%s]: %w", key, err)
		}
	}

	return nil
}

