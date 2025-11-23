package translations

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes translation endpoints.
type Handler interface {
	RequestTranslation(c *fiber.Ctx) error
	GetTranslation(c *fiber.Ctx) error
	ListTranslations(c *fiber.Ctx) error
	TranslateEntity(c *fiber.Ctx) error // Convenience endpoint to translate an entity to all languages
}

type handler struct {
	service Service
	logger  *slog.Logger
	db      *gorm.DB // For fetching entity content
}

// NewHandler constructs a translation handler.
func NewHandler(service Service, db *gorm.DB, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
		db:      db,
	}
}

type requestTranslationPayload struct {
	EntityType EntityType           `json:"entityType"`
	EntityID   uuid.UUID            `json:"entityId"`
	Language   Language             `json:"language"`
	Fields     []string             `json:"fields"`
	SourceText map[string]string   `json:"sourceText"`
}

type translateEntityPayload struct {
	EntityType EntityType `json:"entityType"`
	EntityID   uuid.UUID  `json:"entityId"`
	Languages  []Language `json:"languages,omitempty"` // If empty, translates to all languages
}

func (h *handler) RequestTranslation(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "authentication required",
		})
	}
	_ = userID // User ID available for future authorization checks

	var payload requestTranslationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// If sourceText is empty, fetch it from the entity
	if len(payload.SourceText) == 0 {
		sourceText, err := h.fetchSourceTextFromEntity(c.Context(), payload.EntityType, payload.EntityID, payload.Fields)
		if err != nil {
			h.logger.Warn("Failed to fetch source text from entity", slog.Any("error", err), slog.String("entityType", string(payload.EntityType)), slog.String("entityId", payload.EntityID.String()))
			// Continue anyway - the worker will handle empty sourceText
		} else {
			payload.SourceText = sourceText
		}
	}

	if err := h.service.RequestTranslation(c.Context(), payload.EntityType, payload.EntityID, payload.Language, payload.Fields, payload.SourceText); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusAccepted, fiber.Map{
		"message": "Translation job queued successfully",
	})
}

// fetchSourceTextFromEntity fetches source text from the entity when not provided
func (h *handler) fetchSourceTextFromEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, fields []string) (map[string]string, error) {
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
		if err := h.db.WithContext(ctx).Table("testimonials").Where("id = ?", entityID).First(&testimonial).Error; err != nil {
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
		if err := h.db.WithContext(ctx).Table("posts").Where("id = ?", entityID).First(&post).Error; err != nil {
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
		if err := h.db.WithContext(ctx).Table("projects").Where("id = ?", entityID).First(&project).Error; err != nil {
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
		if err := h.db.WithContext(ctx).Table("case_studies").Where("id = ?", entityID).First(&caseStudy).Error; err != nil {
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
	default:
		return nil, fiber.NewError(fiber.StatusBadRequest, "unsupported entity type for auto-fetch")
	}

	return sourceText, nil
}

func (h *handler) GetTranslation(c *fiber.Ctx) error {
	entityTypeStr := c.Query("entityType")
	entityIDStr := c.Query("entityId")
	languageStr := c.Query("language")

	if entityTypeStr == "" || entityIDStr == "" || languageStr == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "entityType, entityId, and language are required",
		})
	}

	entityType := EntityType(entityTypeStr)
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid entityId",
		})
	}
	language := Language(languageStr)

	translation, err := h.service.GetTranslation(c.Context(), entityType, entityID, language)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, translation)
}

func (h *handler) ListTranslations(c *fiber.Ctx) error {
	filters := TranslationFilters{}

	if entityTypeStr := c.Query("entityType"); entityTypeStr != "" {
		et := EntityType(entityTypeStr)
		filters.EntityType = &et
	}

	if entityIDStr := c.Query("entityId"); entityIDStr != "" {
		entityID, err := uuid.Parse(entityIDStr)
		if err == nil {
			filters.EntityID = &entityID
		}
	}

	if languageStr := c.Query("language"); languageStr != "" {
		lang := Language(languageStr)
		filters.Language = &lang
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := TranslationStatus(statusStr)
		filters.Status = &status
	}

	translations, err := h.service.ListTranslations(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, translations)
}

func (h *handler) TranslateEntity(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "authentication required",
		})
	}
	_ = userID

	var payload translateEntityPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// If no languages specified, translate to all supported languages
	languages := payload.Languages
	if len(languages) == 0 {
		languages = []Language{
			LanguagePTBR, LanguageFR, LanguageES, LanguageDE, LanguageRU,
			LanguageJA, LanguageKO, LanguageZHCN, LanguageEL, LanguageLA,
		}
	}

	// This is a convenience endpoint - actual implementation would need to fetch source text
	// from the entity and queue jobs for each language
	// For now, return success and let the caller use RequestTranslation for each language
	return response.Success(c, fiber.StatusAccepted, fiber.Map{
		"message": "Use /translations/request endpoint for each language",
		"languages": languages,
	})
}

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error in translation handler", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
			"message": "internal server error",
		})
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidEntityType, ErrCodeInvalidLanguage:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeJobQueueFailure, ErrCodeAIServiceFailure:
		statusCode = fiber.StatusServiceUnavailable
	}

	return response.Error(c, statusCode, domainErr.Code, fiber.Map{
		"message": domainErr.Message,
	})
}

