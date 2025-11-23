package translations

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
)

// Enricher enriches entities with translations.
type Enricher struct {
	repo   translationsdomain.Repository
	logger *slog.Logger
}

// NewEnricher creates a new translation enricher.
func NewEnricher(repo translationsdomain.Repository, logger *slog.Logger) *Enricher {
	return &Enricher{
		repo:   repo,
		logger: logger,
	}
}

// EnrichTestimonial enriches a testimonial with translations.
func (e *Enricher) EnrichTestimonial(ctx context.Context, testimonial interface{}, language translationsdomain.Language) error {
	if language == translationsdomain.LanguageEN {
		return nil
	}

	// Use reflection or type assertion to get ID and fields
	// For now, we'll use a simpler approach with a callback
	return nil
}

// EnrichEntityFields enriches entity fields with translations.
// entityType: the type of entity
// entityID: the UUID of the entity
// language: target language
// fieldMap: map of field names to pointers where translations should be applied
func (e *Enricher) EnrichEntityFields(
	ctx context.Context,
	entityType translationsdomain.EntityType,
	entityID uuid.UUID,
	language translationsdomain.Language,
	fieldMap map[string]*string,
) error {
	if language == translationsdomain.LanguageEN {
		return nil // No translation needed for English
	}

	translation, err := e.repo.GetTranslationByEntity(ctx, entityType, entityID, language)
	if err != nil {
		// Translation not found - not an error, just return
		return nil
	}

	if translation.Status != translationsdomain.TranslationStatusCompleted {
		// Translation not ready
		return nil
	}

	translatedFields, err := translation.GetFields()
	if err != nil {
		return err
	}

	// Apply translations to the provided field pointers
	for fieldName, fieldPtr := range fieldMap {
		if translated, ok := translatedFields[fieldName]; ok && translated != "" && fieldPtr != nil {
			*fieldPtr = translated
		}
	}

	return nil
}

