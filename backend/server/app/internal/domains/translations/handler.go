package translations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
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
	// For POST requests, try to get userID from API key context first, then JWT
	var userID uuid.UUID
	var err error
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		userID = apiKey.UserID
	} else {
		userID, err = authdomain.UserIDFromContext(c)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
				"message": "authentication required",
			})
		}
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
	case EntityTypeAIMLIntegration:
		type AIMLIntegration struct {
			Title       string `gorm:"column:title"`
			Description string `gorm:"column:description"`
			UseCase     string `gorm:"column:use_case"`
			Impact      string `gorm:"column:impact"`
			Architecture string `gorm:"column:architecture"`
		}
		var integration AIMLIntegration
		if err := h.db.WithContext(ctx).Table("aiml_integrations").Where("id = ?", entityID).First(&integration).Error; err != nil {
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
		if err := h.db.WithContext(ctx).Table("impact_metrics").Where("id = ?", entityID).First(&metric).Error; err != nil {
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
		if err := h.db.WithContext(ctx).Table("social_media_posts").Where("id = ?", entityID).First(&post).Error; err != nil {
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
		if err := h.db.WithContext(ctx).Table("technical_writings").Where("id = ?", entityID).First(&writing).Error; err != nil {
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
		if err := h.db.WithContext(ctx).Table("interests").Where("id = ?", entityID).First(&interest).Error; err != nil {
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
		if err := h.db.WithContext(ctx).Table("certifications").Where("id = ?", entityID).First(&certification).Error; err != nil {
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
	case EntityTypeSkill:
		type Skill struct {
			Name        string `gorm:"column:name"`
			Description string `gorm:"column:description"`
		}
		var skill Skill
		if err := h.db.WithContext(ctx).Table("skills").Where("id = ?", entityID).First(&skill).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "name":
				sourceText["name"] = skill.Name
			case "description":
				sourceText["description"] = skill.Description
			}
		}
	case EntityTypeSystemDesign:
		type SystemDesign struct {
			Title       string `gorm:"column:title"`
			Description string `gorm:"column:description"`
			DataFlow    string `gorm:"column:data_flow"`
			Scalability string `gorm:"column:scalability"`
			Reliability string `gorm:"column:reliability"`
		}
		var systemDesign SystemDesign
		if err := h.db.WithContext(ctx).Table("system_designs").Where("id = ?", entityID).First(&systemDesign).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "title":
				sourceText["title"] = systemDesign.Title
			case "description":
				sourceText["description"] = systemDesign.Description
			case "dataFlow":
				sourceText["dataFlow"] = systemDesign.DataFlow
			case "scalability":
				sourceText["scalability"] = systemDesign.Scalability
			case "reliability":
				sourceText["reliability"] = systemDesign.Reliability
			}
		}
	case EntityTypeProblemSolution:
		type ProblemSolution struct {
			Problem  string `gorm:"column:problem"`
			Context  string `gorm:"column:context"`
			Solution string `gorm:"column:solution"`
			Impact   string `gorm:"column:impact"`
		}
		var problemSolution ProblemSolution
		if err := h.db.WithContext(ctx).Table("problem_solutions").Where("id = ?", entityID).First(&problemSolution).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "problem":
				sourceText["problem"] = problemSolution.Problem
			case "context":
				sourceText["context"] = problemSolution.Context
			case "solution":
				sourceText["solution"] = problemSolution.Solution
			case "impact":
				sourceText["impact"] = problemSolution.Impact
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
	// For POST requests, try to get userID from API key context first, then JWT
	var userID uuid.UUID
	var err error
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		userID = apiKey.UserID
	} else {
		userID, err = authdomain.UserIDFromContext(c)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
				"message": "authentication required",
			})
		}
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

	// Determine fields to translate based on entity type
	var fields []string
	switch payload.EntityType {
	case EntityTypeTestimonial:
		fields = []string{"content", "context", "authorRole", "authorCompany"}
	case EntityTypePost:
		fields = []string{"title", "content", "excerpt", "metaTitle", "metaDescription", "ogTitle", "ogDescription"}
	case EntityTypeProject:
		fields = []string{"name", "description"}
	case EntityTypeCaseStudy:
		fields = []string{"title", "problem", "context", "solution"}
	case EntityTypeProjectCaseStudy:
		fields = []string{"title", "description", "challenge", "solution", "architecture"}
	case EntityTypeSystemDesign:
		fields = []string{"title", "description", "dataFlow", "scalability", "reliability"}
	case EntityTypeProblemSolution:
		fields = []string{"problem", "context", "solution", "impact"}
	case EntityTypeCertification:
		fields = []string{"name", "issuer", "description"}
	case EntityTypeAIMLIntegration:
		fields = []string{"title", "description", "useCase", "impact", "architecture"}
	case EntityTypeImpactMetric:
		fields = []string{"description"}
	case EntityTypeSocialMediaPost:
		fields = []string{"title", "contentPreview"}
	case EntityTypeTechnicalWriting:
		fields = []string{"title", "description", "content", "excerpt"}
	case EntityTypeInterest:
		fields = []string{"title", "description"}
	case EntityTypeSkill:
		fields = []string{"name", "description"}
	case EntityTypeExperience:
		fields = []string{"position", "description"}
	default:
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidEntityType, fiber.Map{
			"message": "unsupported entity type",
		})
	}

	// Fetch source text from entity
	sourceText, err := h.fetchSourceTextFromEntity(c.Context(), payload.EntityType, payload.EntityID, fields)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": fmt.Sprintf("failed to fetch entity: %v", err),
		})
	}

	// Filter out empty fields
	filteredFields := []string{}
	filteredSourceText := make(map[string]string)
	for _, field := range fields {
		if text, ok := sourceText[field]; ok && text != "" {
			filteredFields = append(filteredFields, field)
			filteredSourceText[field] = text
		}
	}

	if len(filteredFields) == 0 {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "no translatable fields found in entity",
		})
	}

	// Queue translation jobs for each language
	queuedCount := 0
	for _, language := range languages {
		if err := h.service.RequestTranslation(c.Context(), payload.EntityType, payload.EntityID, language, filteredFields, filteredSourceText); err != nil {
			h.logger.Warn("Failed to queue translation",
				slog.String("entityType", string(payload.EntityType)),
				slog.String("entityId", payload.EntityID.String()),
				slog.String("language", string(language)),
				slog.Any("error", err),
			)
			continue
		}
		queuedCount++
	}

	return response.Success(c, fiber.StatusAccepted, fiber.Map{
		"message":     fmt.Sprintf("Queued %d translation jobs for entity", queuedCount),
		"entityType":  payload.EntityType,
		"entityId":    payload.EntityID,
		"languages":   languages,
		"fields":      filteredFields,
		"queuedCount": queuedCount,
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

