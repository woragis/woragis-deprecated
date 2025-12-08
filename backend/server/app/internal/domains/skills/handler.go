package skills

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationenricher "github.com/woragis/backend/server/app/pkg/translations"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes skill endpoints.
type Handler interface {
	CreateSkill(c *fiber.Ctx) error
	UpdateSkill(c *fiber.Ctx) error
	DeleteSkill(c *fiber.Ctx) error
	GetSkill(c *fiber.Ctx) error
	GetSkillBySlug(c *fiber.Ctx) error
	ListSkills(c *fiber.Ctx) error
	ListSkillsByCategory(c *fiber.Ctx) error
	SearchSkills(c *fiber.Ctx) error
	GetAllSkillsWithProjectCounts(c *fiber.Ctx) error

	// Project-Skill relationship handlers
	AttachSkillToProject(c *fiber.Ctx) error
	DetachSkillFromProject(c *fiber.Ctx) error
	GetProjectSkills(c *fiber.Ctx) error
	GetProjectsBySkill(c *fiber.Ctx) error
	
	// Timeline operations
	GetSkillsTimeline(c *fiber.Ctx) error
}

type handler struct {
	service           Service
	enricher          *translationenricher.Enricher
	translationService translationsdomain.Service
	logger            *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a skill handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createSkillPayload struct {
	Name              string          `json:"name"`
	Description       string          `json:"description,omitempty"`
	Icon              string          `json:"icon,omitempty"`
	Color             string          `json:"color,omitempty"`
	BgGradient        string          `json:"bgGradient,omitempty"`
	BorderColor       string          `json:"borderColor,omitempty"`
	HoverBorderColor  string          `json:"hoverBorderColor,omitempty"`
	ShadowColor       string          `json:"shadowColor,omitempty"`
	Category          SkillCategory   `json:"category"`
	ProficiencyLevel  ProficiencyLevel `json:"proficiencyLevel,omitempty"`
	YearsOfExperience *int            `json:"yearsOfExperience,omitempty"`
	FirstUsedDate     *string         `json:"firstUsedDate,omitempty"` // ISO date string
	LastUsedDate      *string         `json:"lastUsedDate,omitempty"`  // ISO date string
}

type updateSkillPayload struct {
	Name              string          `json:"name,omitempty"`
	Description       string          `json:"description,omitempty"`
	Icon              string          `json:"icon,omitempty"`
	Color             string          `json:"color,omitempty"`
	BgGradient        string          `json:"bgGradient,omitempty"`
	BorderColor       string          `json:"borderColor,omitempty"`
	HoverBorderColor  string          `json:"hoverBorderColor,omitempty"`
	ShadowColor       string          `json:"shadowColor,omitempty"`
	Category          SkillCategory   `json:"category,omitempty"`
	ProficiencyLevel  ProficiencyLevel `json:"proficiencyLevel,omitempty"`
	YearsOfExperience *int            `json:"yearsOfExperience,omitempty"`
	FirstUsedDate     *string         `json:"firstUsedDate,omitempty"` // ISO date string
	LastUsedDate      *string         `json:"lastUsedDate,omitempty"`  // ISO date string
}

// Handlers

func (h *handler) CreateSkill(c *fiber.Ctx) error {
	var payload createSkillPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var firstUsedDate, lastUsedDate *time.Time
	if payload.FirstUsedDate != nil {
		parsed, err := time.Parse("2006-01-02", *payload.FirstUsedDate)
		if err == nil {
			firstUsedDate = &parsed
		}
	}
	if payload.LastUsedDate != nil {
		parsed, err := time.Parse("2006-01-02", *payload.LastUsedDate)
		if err == nil {
			lastUsedDate = &parsed
		}
	}

	skill, err := h.service.CreateSkill(c.Context(), CreateSkillRequest{
		Name:              payload.Name,
		Description:       payload.Description,
		Icon:              payload.Icon,
		Color:             payload.Color,
		BgGradient:        payload.BgGradient,
		BorderColor:       payload.BorderColor,
		HoverBorderColor:  payload.HoverBorderColor,
		ShadowColor:       payload.ShadowColor,
		Category:          payload.Category,
		ProficiencyLevel:  payload.ProficiencyLevel,
		YearsOfExperience: payload.YearsOfExperience,
		FirstUsedDate:     firstUsedDate,
		LastUsedDate:      lastUsedDate,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toSkillResponse(skill))
}

func (h *handler) UpdateSkill(c *fiber.Ctx) error {
	skillID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateSkillPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var firstUsedDate, lastUsedDate *time.Time
	if payload.FirstUsedDate != nil {
		parsed, err := time.Parse("2006-01-02", *payload.FirstUsedDate)
		if err == nil {
			firstUsedDate = &parsed
		}
	}
	if payload.LastUsedDate != nil {
		parsed, err := time.Parse("2006-01-02", *payload.LastUsedDate)
		if err == nil {
			lastUsedDate = &parsed
		}
	}

	skill, err := h.service.UpdateSkill(c.Context(), UpdateSkillRequest{
		SkillID:           skillID,
		Name:              payload.Name,
		Description:       payload.Description,
		Icon:              payload.Icon,
		Color:             payload.Color,
		BgGradient:        payload.BgGradient,
		BorderColor:       payload.BorderColor,
		HoverBorderColor:  payload.HoverBorderColor,
		ShadowColor:       payload.ShadowColor,
		Category:          payload.Category,
		ProficiencyLevel:  payload.ProficiencyLevel,
		YearsOfExperience: payload.YearsOfExperience,
		FirstUsedDate:     firstUsedDate,
		LastUsedDate:      lastUsedDate,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toSkillResponse(skill))
}

func (h *handler) DeleteSkill(c *fiber.Ctx) error {
	skillIDStr := c.Params("id")
	if skillIDStr == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skillID, err := uuid.Parse(skillIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	_, err = authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, nil)
	}

	if err := h.service.DeleteSkill(c.Context(), skillID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, map[string]string{"message": "Skill deleted successfully"})
}

func (h *handler) GetSkill(c *fiber.Ctx) error {
	skillID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skill, err := h.service.GetSkill(c.Context(), skillID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"name":        &skill.Name,
			"description": &skill.Description,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSkill, skill.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toSkillResponse(skill))
}

func (h *handler) GetSkillBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skill, err := h.service.GetSkillBySlug(c.Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"name":        &skill.Name,
			"description": &skill.Description,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSkill, skill.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toSkillResponse(skill))
}

func (h *handler) ListSkills(c *fiber.Ctx) error {
	skills, err := h.service.ListSkills(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each skill
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range skills {
			fieldMap := map[string]*string{
				"name":        &skills[i].Name,
				"description": &skills[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSkill, skills[i].ID, language, fieldMap)
		}
	}

	responses := make([]skillResponse, len(skills))
	for i := range skills {
		responses[i] = toSkillResponse(&skills[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) ListSkillsByCategory(c *fiber.Ctx) error {
	category := SkillCategory(c.Query("category"))
	if category == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidCategory, nil)
	}

	skills, err := h.service.ListSkillsByCategory(c.Context(), category)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each skill
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range skills {
			fieldMap := map[string]*string{
				"name":        &skills[i].Name,
				"description": &skills[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSkill, skills[i].ID, language, fieldMap)
		}
	}

	responses := make([]skillResponse, len(skills))
	for i := range skills {
		responses[i] = toSkillResponse(&skills[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) SearchSkills(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skills, err := h.service.SearchSkills(c.Context(), query)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each skill
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range skills {
			fieldMap := map[string]*string{
				"name":        &skills[i].Name,
				"description": &skills[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSkill, skills[i].ID, language, fieldMap)
		}
	}

	responses := make([]skillResponse, len(skills))
	for i := range skills {
		responses[i] = toSkillResponse(&skills[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetAllSkillsWithProjectCounts(c *fiber.Ctx) error {
	skills, err := h.service.GetAllSkillsWithProjectCounts(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each skill
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range skills {
			fieldMap := map[string]*string{
				"name":        &skills[i].Name,
				"description": &skills[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSkill, skills[i].ID, language, fieldMap)
		}
	}

	responses := make([]skillWithCountResponse, len(skills))
	for i, skill := range skills {
		responses[i] = toSkillWithCountResponse(skill)
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// Project-Skill relationship handlers

func (h *handler) AttachSkillToProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skillID, err := uuid.Parse(c.Params("skillId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	_, err = authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	// Verify project ownership (we need to check this through projects domain)
	// For now, we'll trust the middleware, but ideally we'd inject a projects service
	// to verify ownership. This is a simplified version.

	if err := h.service.AttachSkillToProject(c.Context(), projectID, skillID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "skill attached to project"})
}

func (h *handler) DetachSkillFromProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skillID, err := uuid.Parse(c.Params("skillId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	_, err = authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	// Verify project ownership (same note as above)

	if err := h.service.DetachSkillFromProject(c.Context(), projectID, skillID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "skill detached from project"})
}

func (h *handler) GetProjectSkills(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skills, err := h.service.GetProjectSkills(c.Context(), projectID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each skill
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range skills {
			fieldMap := map[string]*string{
				"name":        &skills[i].Name,
				"description": &skills[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSkill, skills[i].ID, language, fieldMap)
		}
	}

	responses := make([]skillResponse, len(skills))
	for i := range skills {
		responses[i] = toSkillResponse(&skills[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetProjectsBySkill(c *fiber.Ctx) error {
	skillID, err := uuid.Parse(c.Params("skillId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	projectIDs, err := h.service.GetProjectsBySkill(c.Context(), skillID)
	if err != nil {
		return h.handleError(c, err)
	}

	projectIDStrings := make([]string, len(projectIDs))
	for i, id := range projectIDs {
		projectIDStrings[i] = id.String()
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"project_ids": projectIDStrings})
}

func (h *handler) GetSkillsTimeline(c *fiber.Ctx) error {
	skills, err := h.service.GetSkillsTimeline(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each skill
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range skills {
			fieldMap := map[string]*string{
				"name":        &skills[i].Name,
				"description": &skills[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSkill, skills[i].ID, language, fieldMap)
		}
	}

	responses := make([]skillResponse, len(skills))
	for i := range skills {
		responses[i] = toSkillResponse(&skills[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// Response helpers

type skillResponse struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Slug              string          `json:"slug"`
	Category          SkillCategory   `json:"category"`
	Description       string          `json:"description,omitempty"`
	Icon              string          `json:"icon,omitempty"`
	Color             string          `json:"color,omitempty"`
	BgGradient        string          `json:"bgGradient,omitempty"`
	BorderColor       string          `json:"borderColor,omitempty"`
	HoverBorderColor  string          `json:"hoverBorderColor,omitempty"`
	ShadowColor       string          `json:"shadowColor,omitempty"`
	ProficiencyLevel  ProficiencyLevel `json:"proficiencyLevel,omitempty"`
	YearsOfExperience *int            `json:"yearsOfExperience,omitempty"`
	FirstUsedDate     *string         `json:"firstUsedDate,omitempty"`
	LastUsedDate      *string         `json:"lastUsedDate,omitempty"`
	CreatedAt         string          `json:"createdAt"`
	UpdatedAt         string          `json:"updatedAt"`
}

type skillWithCountResponse struct {
	skillResponse
	ProjectCount int64 `json:"projectCount"`
}

func toSkillResponse(skill *Skill) skillResponse {
	var firstUsedDate, lastUsedDate *string
	if skill.FirstUsedDate != nil {
		dateStr := skill.FirstUsedDate.Format("2006-01-02")
		firstUsedDate = &dateStr
	}
	if skill.LastUsedDate != nil {
		dateStr := skill.LastUsedDate.Format("2006-01-02")
		lastUsedDate = &dateStr
	}

	return skillResponse{
		ID:                skill.ID.String(),
		Name:              skill.Name,
		Slug:              skill.Slug,
		Category:          skill.Category,
		Description:       skill.Description,
		Icon:              skill.Icon,
		Color:             skill.Color,
		BgGradient:        skill.BgGradient,
		BorderColor:       skill.BorderColor,
		HoverBorderColor:  skill.HoverBorderColor,
		ShadowColor:       skill.ShadowColor,
		ProficiencyLevel:  skill.ProficiencyLevel,
		YearsOfExperience: skill.YearsOfExperience,
		FirstUsedDate:     firstUsedDate,
		LastUsedDate:      lastUsedDate,
		CreatedAt:         skill.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         skill.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toSkillWithCountResponse(skill SkillWithCount) skillWithCountResponse {
	return skillWithCountResponse{
		skillResponse: toSkillResponse(&skill.Skill),
		ProjectCount:  skill.ProjectCount,
	}
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, nil)
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidName, ErrCodeInvalidCategory:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeConflict:
		statusCode = fiber.StatusConflict
	}

	return response.Error(c, statusCode, domainErr.Code, domainErr.Message)
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, 0, "unauthorized")
}

