package translations

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	db        *gorm.DB // GORM DB for fetching entities
}

// NewService constructs a Service.
func NewService(repo Repository, queue Queue, aiClient *langchainservice.Client, db *gorm.DB, logger *slog.Logger) Service {
	return &service{
		repo:     repo,
		queue:    queue,
		aiClient: aiClient,
		logger:   logger,
		db:       db,
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

// fetchSourceTextFromEntity fetches source text from the entity when not provided
func (s *service) fetchSourceTextFromEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, fields []string) (map[string]string, error) {
	sourceText := make(map[string]string)

	switch entityType {
	case EntityTypeTestimonial:
		// Use a generic struct to avoid import cycle
		type Testimonial struct {
			Content        string `gorm:"column:content"`
			AuthorRole     string `gorm:"column:author_role"`
			AuthorCompany  string `gorm:"column:author_company"`
		}
		var testimonial Testimonial
		if err := s.db.WithContext(ctx).Table("testimonials").Where("id = ?", entityID).First(&testimonial).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "content":
				sourceText["content"] = testimonial.Content
			case "authorRole":
				sourceText["authorRole"] = testimonial.AuthorRole
			case "authorCompany":
				sourceText["authorCompany"] = testimonial.AuthorCompany
			}
		}
	case EntityTypePost:
		// Use a generic struct to avoid import cycle
		type Post struct {
			Title          string `gorm:"column:title"`
			Content        string `gorm:"column:content"`
			Excerpt        string `gorm:"column:excerpt"`
			MetaTitle      string `gorm:"column:meta_title"`
			MetaDescription string `gorm:"column:meta_description"`
			OGTitle        string `gorm:"column:og_title"`
			OGDescription  string `gorm:"column:og_description"`
		}
		var post Post
		if err := s.db.WithContext(ctx).Table("posts").Where("id = ?", entityID).First(&post).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "title":
				sourceText["title"] = post.Title
			case "content":
				sourceText["content"] = post.Content
			case "excerpt":
				sourceText["excerpt"] = post.Excerpt
			case "metaTitle":
				sourceText["metaTitle"] = post.MetaTitle
			case "metaDescription":
				sourceText["metaDescription"] = post.MetaDescription
			case "ogTitle":
				sourceText["ogTitle"] = post.OGTitle
			case "ogDescription":
				sourceText["ogDescription"] = post.OGDescription
			}
		}
	case EntityTypeProject:
		// Use a generic struct to avoid import cycle
		type Project struct {
			Name        string `gorm:"column:name"`
			Description string `gorm:"column:description"`
		}
		var project Project
		if err := s.db.WithContext(ctx).Table("projects").Where("id = ?", entityID).First(&project).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "name":
				sourceText["name"] = project.Name
			case "description":
				sourceText["description"] = project.Description
			}
		}
	case EntityTypeCaseStudy:
		// Use a generic struct to avoid import cycle
		type CaseStudy struct {
			Title    string `gorm:"column:title"`
			Problem  string `gorm:"column:problem"`
			Context  string `gorm:"column:context"`
			Solution string `gorm:"column:solution"`
		}
		var caseStudy CaseStudy
		if err := s.db.WithContext(ctx).Table("case_studies").Where("id = ?", entityID).First(&caseStudy).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "title":
				sourceText["title"] = caseStudy.Title
			case "problem":
				sourceText["problem"] = caseStudy.Problem
			case "context":
				sourceText["context"] = caseStudy.Context
			case "solution":
				sourceText["solution"] = caseStudy.Solution
			}
		}
	case EntityTypeAIMLIntegration:
		type AIMLIntegration struct {
			Title       string `gorm:"column:title"`
			Description string `gorm:"column:description"`
			UseCase     string `gorm:"column:use_case"`
			Impact      string `gorm:"column:impact"`
			Architecture string `gorm:"column:architecture"`
		}
		var integration AIMLIntegration
		if err := s.db.WithContext(ctx).Table("aiml_integrations").Where("id = ?", entityID).First(&integration).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "title":
				sourceText["title"] = integration.Title
			case "description":
				sourceText["description"] = integration.Description
			case "useCase":
				sourceText["useCase"] = integration.UseCase
			case "impact":
				sourceText["impact"] = integration.Impact
			case "architecture":
				sourceText["architecture"] = integration.Architecture
			}
		}
	case EntityTypeImpactMetric:
		type ImpactMetric struct {
			Description string `gorm:"column:description"`
		}
		var metric ImpactMetric
		if err := s.db.WithContext(ctx).Table("impact_metrics").Where("id = ?", entityID).First(&metric).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "description":
				sourceText["description"] = metric.Description
			}
		}
	case EntityTypeSocialMediaPost:
		type SocialMediaPost struct {
			Title         string `gorm:"column:title"`
			ContentPreview string `gorm:"column:content_preview"`
		}
		var post SocialMediaPost
		if err := s.db.WithContext(ctx).Table("social_media_posts").Where("id = ?", entityID).First(&post).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "title":
				sourceText["title"] = post.Title
			case "contentPreview":
				sourceText["contentPreview"] = post.ContentPreview
			}
		}
	case EntityTypeTechnicalWriting:
		type TechnicalWriting struct {
			Title       string `gorm:"column:title"`
			Description string `gorm:"column:description"`
			Content     string `gorm:"column:content"`
			Excerpt     string `gorm:"column:excerpt"`
		}
		var writing TechnicalWriting
		if err := s.db.WithContext(ctx).Table("technical_writings").Where("id = ?", entityID).First(&writing).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "title":
				sourceText["title"] = writing.Title
			case "description":
				sourceText["description"] = writing.Description
			case "content":
				sourceText["content"] = writing.Content
			case "excerpt":
				sourceText["excerpt"] = writing.Excerpt
			}
		}
	case EntityTypeInterest:
		type Interest struct {
			Title       string `gorm:"column:title"`
			Description string `gorm:"column:description"`
		}
		var interest Interest
		if err := s.db.WithContext(ctx).Table("interests").Where("id = ?", entityID).First(&interest).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "title":
				sourceText["title"] = interest.Title
			case "description":
				sourceText["description"] = interest.Description
			}
		}
	case EntityTypeCertification:
		type Certification struct {
			Name        string `gorm:"column:name"`
			Issuer      string `gorm:"column:issuer"`
			Description string `gorm:"column:description"`
		}
		var certification Certification
		if err := s.db.WithContext(ctx).Table("certifications").Where("id = ?", entityID).First(&certification).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "name":
				sourceText["name"] = certification.Name
			case "issuer":
				sourceText["issuer"] = certification.Issuer
			case "description":
				sourceText["description"] = certification.Description
			}
		}
	default:
		return nil, fiber.NewError(fiber.StatusBadRequest, "unsupported entity type for auto-fetch")
	}

	return sourceText, nil
}

