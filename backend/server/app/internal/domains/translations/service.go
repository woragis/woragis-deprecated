package translations

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
)

// Service orchestrates translation workflows.
type Service interface {
	RequestTranslation(ctx context.Context, entityType EntityType, entityID uuid.UUID, language Language, fields []string, sourceText map[string]string) error
	GetTranslation(ctx context.Context, entityType EntityType, entityID uuid.UUID, language Language) (*Translation, error)
	ListTranslations(ctx context.Context, filters TranslationFilters) ([]Translation, error)
	ProcessTranslationJob(ctx context.Context, job *TranslationJob) error
}

type service struct {
	repo      Repository
	queue     Queue
	aiClient  *langchainservice.Client
	logger    *slog.Logger
	db        interface{} // GORM DB for fetching entities (temporary solution)
}

// NewService constructs a Service.
func NewService(repo Repository, queue Queue, aiClient *langchainservice.Client, logger *slog.Logger) Service {
	return &service{
		repo:     repo,
		queue:    queue,
		aiClient: aiClient,
		logger:   logger,
	}
}

func (s *service) RequestTranslation(ctx context.Context, entityType EntityType, entityID uuid.UUID, language Language, fields []string, sourceText map[string]string) error {
	// Check if translation already exists
	existing, err := s.repo.GetTranslationByEntity(ctx, entityType, entityID, language)
	if err == nil && existing != nil && existing.Status == TranslationStatusCompleted {
		// Translation already exists and is completed
		return nil
	}

	// Create translation job
	job := &TranslationJob{
		ID:         uuid.New().String(),
		EntityType: entityType,
		EntityID:   entityID.String(),
		Language:   language,
		Fields:     fields,
		SourceText: sourceText,
	}

	// Enqueue job
	if err := s.queue.EnqueueJob(ctx, job); err != nil {
		return err
	}

	// Create pending translation record
	translation, err := NewTranslation(entityType, entityID, language, make(map[string]string))
	if err != nil {
		return err
	}
	translation.Status = TranslationStatusPending

	// If translation exists but is not completed, update it
	if existing != nil {
		translation.ID = existing.ID
		translation.Status = TranslationStatusPending
		return s.repo.UpdateTranslation(ctx, translation)
	}

	return s.repo.CreateTranslation(ctx, translation)
}

func (s *service) GetTranslation(ctx context.Context, entityType EntityType, entityID uuid.UUID, language Language) (*Translation, error) {
	return s.repo.GetTranslationByEntity(ctx, entityType, entityID, language)
}

func (s *service) ListTranslations(ctx context.Context, filters TranslationFilters) ([]Translation, error) {
	return s.repo.ListTranslations(ctx, filters)
}

func (s *service) ProcessTranslationJob(ctx context.Context, job *TranslationJob) error {
	entityID, err := uuid.Parse(job.EntityID)
	if err != nil {
		return fmt.Errorf("invalid entity ID: %w", err)
	}

	// Get or create translation record
	translation, err := s.repo.GetTranslationByEntity(ctx, job.EntityType, entityID, job.Language)
	if err != nil {
		// Create new translation if it doesn't exist
		translation, err = NewTranslation(job.EntityType, entityID, job.Language, make(map[string]string))
		if err != nil {
			return err
		}
		translation.Status = TranslationStatusProcessing
		if err := s.repo.CreateTranslation(ctx, translation); err != nil {
			return err
		}
	} else {
		translation.Status = TranslationStatusProcessing
		if err := s.repo.UpdateTranslation(ctx, translation); err != nil {
			return err
		}
	}

	// Translate each field using AI service
	translatedFields := make(map[string]string)
	
	// If sourceText is empty, fetch it from the entity
	if len(job.SourceText) == 0 {
		sourceText, err := s.fetchSourceTextFromEntity(ctx, job.EntityType, entityID, job.Fields)
		if err != nil {
			translation.Status = TranslationStatusFailed
			translation.ErrorMessage = fmt.Sprintf("Failed to fetch source text: %v", err)
			s.repo.UpdateTranslation(ctx, translation)
			return fmt.Errorf("failed to fetch source text: %w", err)
		}
		job.SourceText = sourceText
	}
	
	for _, field := range job.Fields {
		sourceText, ok := job.SourceText[field]
		if !ok || sourceText == "" {
			if s.logger != nil {
				s.logger.Warn("Skipping field with empty source text", slog.String("field", field), slog.String("entityId", job.EntityID))
			}
			continue
		}

		translated, err := s.translateText(ctx, sourceText, job.Language)
		if err != nil {
			translation.Status = TranslationStatusFailed
			translation.ErrorMessage = fmt.Sprintf("Failed to translate field %s: %v", field, err)
			s.repo.UpdateTranslation(ctx, translation)
			return fmt.Errorf("failed to translate field %s: %w", field, err)
		}

		translatedFields[field] = translated
	}

	// Update translation with results
	if err := translation.SetFields(translatedFields); err != nil {
		return err
	}
	translation.Status = TranslationStatusCompleted
	translation.ErrorMessage = ""

	return s.repo.UpdateTranslation(ctx, translation)
}

func (s *service) translateText(ctx context.Context, text string, targetLanguage Language) (string, error) {
	if s.aiClient == nil {
		return "", NewDomainError(ErrCodeAIServiceFailure, ErrAIServiceUnavailable)
	}

	// Get language name for prompt
	languageName := s.getLanguageName(targetLanguage)
	
	// Create translation prompt
	prompt := fmt.Sprintf(`Translate the following text to %s. 
Return only the translated text, without any explanations, quotes, or additional formatting.
Preserve any markdown formatting, code blocks, or special characters.
Original text:
%s`, languageName, text)

	// Call AI service
	req := langchainservice.ChatCompletionRequest{
		Provider:    langchainservice.ProviderOpenAI,
		Model:       "gpt-4o-mini",
		Temperature: 0.3, // Lower temperature for more consistent translations
		Messages: []langchainservice.ChatMessage{
			{
				Role:      "user",
				Content:   prompt,
				Timestamp: time.Now(),
			},
		},
		MaxTokens: 2000,
	}

	resp, err := s.aiClient.GenerateCompletion(ctx, req)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("AI translation failed", slog.Any("error", err), slog.String("language", string(targetLanguage)))
		}
		return "", NewDomainError(ErrCodeAIServiceFailure, ErrAIServiceUnavailable)
	}

	translated := strings.TrimSpace(resp.Message.Content)
	
	// Remove quotes if the AI wrapped the response
	translated = strings.Trim(translated, `"'`)
	
	return translated, nil
}

func (s *service) getLanguageName(lang Language) string {
	names := map[Language]string{
		LanguageEN:   "English (US)",
		LanguagePTBR: "Portuguese (Brazil)",
		LanguageFR:   "French",
		LanguageES:   "Spanish",
		LanguageDE:   "German",
		LanguageRU:   "Russian",
		LanguageJA:   "Japanese",
		LanguageKO:   "Korean",
		LanguageZHCN: "Chinese (Simplified)",
		LanguageEL:   "Greek",
		LanguageLA:   "Latin",
	}
	if name, ok := names[lang]; ok {
		return name
	}
	return string(lang)
}

