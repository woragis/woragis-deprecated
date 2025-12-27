package projects

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationenricher "github.com/woragis/backend/server/app/pkg/translations"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes project endpoints.
type Handler interface {
	CreateProject(c *fiber.Ctx) error
	ListProjects(c *fiber.Ctx) error
	GetProjectBySlug(c *fiber.Ctx) error
	SearchProjectsBySlug(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
	UpdateMetrics(c *fiber.Ctx) error
	DeleteProject(c *fiber.Ctx) error
	AddMilestone(c *fiber.Ctx) error
	ToggleMilestoneCompletion(c *fiber.Ctx) error
	ListMilestones(c *fiber.Ctx) error
	BulkUpdateMilestones(c *fiber.Ctx) error

	CreateKanbanColumn(c *fiber.Ctx) error
	UpdateKanbanColumn(c *fiber.Ctx) error
	ReorderKanbanColumns(c *fiber.Ctx) error
	DeleteKanbanColumn(c *fiber.Ctx) error
	CreateKanbanCard(c *fiber.Ctx) error
	UpdateKanbanCard(c *fiber.Ctx) error
	MoveKanbanCard(c *fiber.Ctx) error
	DeleteKanbanCard(c *fiber.Ctx) error
	GetKanbanBoard(c *fiber.Ctx) error

	CreateDependency(c *fiber.Ctx) error
	ListDependencies(c *fiber.Ctx) error
	DeleteDependency(c *fiber.Ctx) error

	DuplicateProject(c *fiber.Ctx) error

	// Documentation handlers
	CreateDocumentation(c *fiber.Ctx) error
	UpdateDocumentationVisibility(c *fiber.Ctx) error
	GetDocumentation(c *fiber.Ctx) error
	GetPublicDocumentation(c *fiber.Ctx) error
	DeleteDocumentation(c *fiber.Ctx) error

	// Documentation Section handlers
	CreateDocumentationSection(c *fiber.Ctx) error
	UpdateDocumentationSection(c *fiber.Ctx) error
	DeleteDocumentationSection(c *fiber.Ctx) error
	ListDocumentationSections(c *fiber.Ctx) error
	ReorderDocumentationSections(c *fiber.Ctx) error

	// Technology handlers
	CreateTechnology(c *fiber.Ctx) error
	UpdateTechnology(c *fiber.Ctx) error
	DeleteTechnology(c *fiber.Ctx) error
	ListTechnologies(c *fiber.Ctx) error
	BulkCreateTechnologies(c *fiber.Ctx) error
	BulkUpdateTechnologies(c *fiber.Ctx) error

	// File Structure handlers
	CreateFileStructure(c *fiber.Ctx) error
	UpdateFileStructure(c *fiber.Ctx) error
	DeleteFileStructure(c *fiber.Ctx) error
	ListFileStructures(c *fiber.Ctx) error
	BulkCreateFileStructures(c *fiber.Ctx) error
	BulkUpdateFileStructures(c *fiber.Ctx) error

	// Architecture Diagram handlers
	CreateArchitectureDiagram(c *fiber.Ctx) error
	UpdateArchitectureDiagram(c *fiber.Ctx) error
	DeleteArchitectureDiagram(c *fiber.Ctx) error
	ListArchitectureDiagrams(c *fiber.Ctx) error
	GetArchitectureDiagram(c *fiber.Ctx) error
}

type handler struct {
	service          Service
	enricher         *translationenricher.Enricher
	translationService translationsdomain.Service
	logger           *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a project handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createProjectPayload struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	HealthScore int     `json:"health_score"`
	MRR         float64 `json:"mrr"`
	CAC         float64 `json:"cac"`
	LTV         float64 `json:"ltv"`
	ChurnRate   float64 `json:"churn_rate"`
}

type updateStatusPayload struct {
	Status string `json:"status"`
}

type updateMetricsPayload struct {
	HealthScore int     `json:"health_score"`
	MRR         float64 `json:"mrr"`
	CAC         float64 `json:"cac"`
	LTV         float64 `json:"ltv"`
	ChurnRate   float64 `json:"churn_rate"`
}

type addMilestonePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type toggleMilestonePayload struct {
	Completed bool `json:"completed"`
}

type bulkMilestoneUpdatePayload struct {
	Updates []bulkMilestoneUpdateItemPayload `json:"updates"`
}

type bulkMilestoneUpdateItemPayload struct {
	MilestoneID string  `json:"milestone_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"`
	Completed   *bool   `json:"completed"`
}

type createColumnPayload struct {
	Name     string `json:"name"`
	WIPLimit *int   `json:"wip_limit"`
	Position *int   `json:"position"`
}

type updateColumnPayload struct {
	Name     *string `json:"name"`
	WIPLimit *int    `json:"wip_limit"`
}

type reorderColumnsPayload struct {
	ColumnOrder []string `json:"column_order"`
}

type deleteColumnPayload struct{}

type createCardPayload struct {
	ColumnID    string `json:"column_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	MilestoneID string `json:"milestone_id"`
	Position    *int   `json:"position"`
}

type updateCardPayload struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"`
	MilestoneID *string `json:"milestone_id"`
}

type moveCardPayload struct {
	TargetColumnID string `json:"target_column_id"`
	TargetPosition int    `json:"target_position"`
}

type deleteCardPayload struct{}

type dependencyPayload struct {
	DependsOnProjectID string `json:"depends_on_project_id"`
	Type               string `json:"type"`
}

type deleteDependencyPayload struct{}

type duplicateProjectPayload struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Status         *string  `json:"status"`
	HealthScore    *int     `json:"health_score"`
	MRR            *float64 `json:"mrr"`
	CAC            *float64 `json:"cac"`
	LTV            *float64 `json:"ltv"`
	ChurnRate      *float64 `json:"churn_rate"`
	CopyBoard      *bool    `json:"copy_board"`
	CopyMilestones *bool    `json:"copy_milestones"`
	CopyDeps       *bool    `json:"copy_dependencies"`
}

// Responses

type projectResponse struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Slug        string        `json:"slug"`
	Status      ProjectStatus `json:"status"`
	HealthScore int           `json:"health_score"`
	MRR         float64       `json:"mrr"`
	CAC         float64       `json:"cac"`
	LTV         float64       `json:"ltv"`
	ChurnRate   float64       `json:"churn_rate"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type milestoneResponse struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type kanbanBoardResponse struct {
	ProjectID string                 `json:"project_id"`
	Columns   []kanbanColumnResponse `json:"columns"`
}

type kanbanColumnResponse struct {
	ID        string               `json:"id"`
	ProjectID string               `json:"project_id"`
	Name      string               `json:"name"`
	WIPLimit  int                  `json:"wip_limit"`
	Position  int                  `json:"position"`
	Cards     []kanbanCardResponse `json:"cards"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type kanbanCardResponse struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	ColumnID    string     `json:"column_id"`
	MilestoneID *string    `json:"milestone_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type dependencyResponse struct {
	ID                 string         `json:"id"`
	ProjectID          string         `json:"project_id"`
	DependsOnProjectID string         `json:"depends_on_project_id"`
	Type               DependencyType `json:"type"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// Handlers

func (h *handler) CreateProject(c *fiber.Ctx) error {
	var payload createProjectPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateCreateProjectPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	project, err := h.service.CreateProject(c.Context(), CreateProjectRequest{
		UserID:      userID,
		Name:        payload.Name,
		Description: payload.Description,
		Status:      ProjectStatus(payload.Status),
		HealthScore: payload.HealthScore,
		MRR:         payload.MRR,
		CAC:         payload.CAC,
		LTV:         payload.LTV,
		ChurnRate:   payload.ChurnRate,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Automatically trigger translations for all supported languages
	if h.translationService != nil {
		// Prepare source text for translation
		sourceText := make(map[string]string)
		if project.Name != "" {
			sourceText["name"] = project.Name
		}
		if project.Description != "" {
			sourceText["description"] = project.Description
		}

		// Fields to translate
		fields := []string{}
		if project.Name != "" {
			fields = append(fields, "name")
		}
		if project.Description != "" {
			fields = append(fields, "description")
		}

		// Queue translations for all supported languages (except English)
		supportedLanguages := []translationsdomain.Language{
			translationsdomain.LanguagePTBR,
			translationsdomain.LanguageFR,
			translationsdomain.LanguageES,
			translationsdomain.LanguageDE,
			translationsdomain.LanguageRU,
			translationsdomain.LanguageJA,
			translationsdomain.LanguageKO,
			translationsdomain.LanguageZHCN,
			translationsdomain.LanguageEL,
			translationsdomain.LanguageLA,
		}

		// Trigger translations asynchronously (don't block the response)
		// Use background context to avoid cancellation when request completes
		go func() {
			ctx := context.Background()
			for _, lang := range supportedLanguages {
				if err := h.translationService.RequestTranslation(
					ctx,
					translationsdomain.EntityTypeProject,
					project.ID,
					lang,
					fields,
					sourceText,
				); err != nil {
					h.logger.Warn("Failed to queue translation",
						slog.String("projectId", project.ID.String()),
						slog.String("language", string(lang)),
						slog.Any("error", err),
					)
				}
			}
		}()
	}

	return response.Success(c, fiber.StatusCreated, toProjectResponse(project))
}

func (h *handler) ListProjects(c *fiber.Ctx) error {
	var userID uuid.UUID
	var err error
	
	// Check if request is authenticated via API key
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		// Use the API key's user ID
		userID = apiKey.UserID
	} else {
		// Fall back to JWT user ID
		userID, err = authdomain.UserIDFromContext(c)
		if err != nil {
			return unauthorizedResponse(c)
		}
	}

	projects, err := h.service.ListProjects(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range projects {
			fieldMap := map[string]*string{
				"name":        &projects[i].Name,
				"description": &projects[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProject, projects[i].ID, language, fieldMap)
		}
	}

	resp := make([]projectResponse, 0, len(projects))
	for _, project := range projects {
		p := project
		resp = append(resp, toProjectResponse(&p))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) GetProjectBySlug(c *fiber.Ctx) error {
	slug := strings.TrimSpace(c.Params("slug"))
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var userID uuid.UUID
	var err error
	
	// Check if request is authenticated via API key
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		// Use the API key's user ID
		userID = apiKey.UserID
	} else {
		// Fall back to JWT user ID
		userID, err = authdomain.UserIDFromContext(c)
		if err != nil {
			return unauthorizedResponse(c)
		}
	}

	project, err := h.service.GetProjectBySlug(c.Context(), userID, slug)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"name":        &project.Name,
			"description": &project.Description,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProject, project.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toProjectResponse(project))
}

func (h *handler) SearchProjectsBySlug(c *fiber.Ctx) error {
	slug := strings.TrimSpace(c.Query("slug"))

	var userID uuid.UUID
	var err error
	
	// Check if request is authenticated via API key
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		// Use the API key's user ID
		userID = apiKey.UserID
	} else {
		// Fall back to JWT user ID
		userID, err = authdomain.UserIDFromContext(c)
		if err != nil {
			return unauthorizedResponse(c)
		}
	}

	if slug == "" {
		return response.Success(c, fiber.StatusOK, []projectResponse{})
	}

	projects, err := h.service.SearchProjectsBySlug(c.Context(), userID, slug)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range projects {
			fieldMap := map[string]*string{
				"name":        &projects[i].Name,
				"description": &projects[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProject, projects[i].ID, language, fieldMap)
		}
	}

	resp := make([]projectResponse, 0, len(projects))
	for _, project := range projects {
		p := project
		resp = append(resp, toProjectResponse(&p))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) UpdateStatus(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateStatusPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateUpdateStatusPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	project, err := h.service.UpdateProjectStatus(c.Context(), UpdateStatusRequest{
		ProjectID: projectID,
		UserID:    userID,
		Status:    ProjectStatus(payload.Status),
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toProjectResponse(project))
}

func (h *handler) UpdateMetrics(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateMetricsPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateUpdateMetricsPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	project, err := h.service.UpdateProjectMetrics(c.Context(), UpdateMetricsRequest{
		ProjectID:   projectID,
		UserID:      userID,
		HealthScore: payload.HealthScore,
		MRR:         payload.MRR,
		CAC:         payload.CAC,
		LTV:         payload.LTV,
		ChurnRate:   payload.ChurnRate,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toProjectResponse(project))
}

func (h *handler) DeleteProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteProject(c.Context(), projectID, userID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, map[string]string{"message": "Project deleted successfully"})
}

func (h *handler) AddMilestone(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload addMilestonePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateAddMilestonePayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var due time.Time
	if payload.DueDate != "" {
		if due, err = time.Parse(time.RFC3339, payload.DueDate); err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	milestone, err := h.service.AddMilestone(c.Context(), AddMilestoneRequest{
		ProjectID:   projectID,
		UserID:      userID,
		Title:       payload.Title,
		Description: payload.Description,
		DueDate:     due,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toMilestoneResponse(milestone))
}

func (h *handler) ToggleMilestoneCompletion(c *fiber.Ctx) error {
	milestoneID, err := uuid.Parse(c.Params("milestoneID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload toggleMilestonePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	milestone, err := h.service.ToggleMilestone(c.Context(), ToggleMilestoneRequest{
		MilestoneID: milestoneID,
		UserID:      userID,
		Completed:   payload.Completed,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toMilestoneResponse(milestone))
}

func (h *handler) ListMilestones(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	milestones, err := h.service.ListMilestones(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]milestoneResponse, 0, len(milestones))
	for _, milestone := range milestones {
		m := milestone
		resp = append(resp, toMilestoneResponse(&m))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) BulkUpdateMilestones(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload bulkMilestoneUpdatePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	updates := make([]MilestoneUpdate, 0, len(payload.Updates))
	for _, item := range payload.Updates {
		milestoneID, err := uuid.Parse(item.MilestoneID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}

		var due *time.Time
		if item.DueDate != nil && *item.DueDate != "" {
			t, err := time.Parse(time.RFC3339, *item.DueDate)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
			}
			due = &t
		}

		updates = append(updates, MilestoneUpdate{
			MilestoneID: milestoneID,
			Title:       item.Title,
			Description: item.Description,
			DueDate:     due,
			Completed:   item.Completed,
		})
	}

	updated, err := h.service.BulkUpdateMilestones(c.Context(), BulkUpdateMilestonesRequest{
		ProjectID: projectID,
		UserID:    userID,
		Updates:   updates,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]milestoneResponse, 0, len(updated))
	for _, milestone := range updated {
		resp = append(resp, toMilestoneResponse(milestone))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// Kanban

func (h *handler) CreateKanbanColumn(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload createColumnPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateCreateColumnPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	board, err := h.service.CreateKanbanColumn(c.Context(), CreateKanbanColumnRequest{
		ProjectID: projectID,
		UserID:    userID,
		Name:      payload.Name,
		WIPLimit:  payload.WIPLimit,
		Position:  payload.Position,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toKanbanBoardResponse(board))
}

func (h *handler) UpdateKanbanColumn(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}
	columnID, err := uuid.Parse(c.Params("columnID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateColumnPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	board, err := h.service.UpdateKanbanColumn(c.Context(), UpdateKanbanColumnRequest{
		ProjectID: projectID,
		UserID:    userID,
		ColumnID:  columnID,
		Name:      payload.Name,
		WIPLimit:  payload.WIPLimit,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toKanbanBoardResponse(board))
}

func (h *handler) ReorderKanbanColumns(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload reorderColumnsPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	columnOrder := make([]uuid.UUID, 0, len(payload.ColumnOrder))
	for _, raw := range payload.ColumnOrder {
		id, err := uuid.Parse(raw)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		columnOrder = append(columnOrder, id)
	}

	board, err := h.service.ReorderKanbanColumns(c.Context(), ReorderKanbanColumnsRequest{
		ProjectID:   projectID,
		UserID:      userID,
		ColumnOrder: columnOrder,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toKanbanBoardResponse(board))
}

func (h *handler) DeleteKanbanColumn(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}
	columnID, err := uuid.Parse(c.Params("columnID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload deleteColumnPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	board, err := h.service.DeleteKanbanColumn(c.Context(), DeleteKanbanColumnRequest{
		ProjectID: projectID,
		UserID:    userID,
		ColumnID:  columnID,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toKanbanBoardResponse(board))
}

func (h *handler) CreateKanbanCard(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload createCardPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateCreateCardPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	columnID, err := uuid.Parse(payload.ColumnID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var due *time.Time
	if payload.DueDate != "" {
		t, err := time.Parse(time.RFC3339, payload.DueDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		due = &t
	}

	var milestoneID *uuid.UUID
	if payload.MilestoneID != "" {
		id, err := uuid.Parse(payload.MilestoneID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		milestoneID = &id
	}

	board, err := h.service.CreateKanbanCard(c.Context(), CreateKanbanCardRequest{
		ProjectID:   projectID,
		UserID:      userID,
		ColumnID:    columnID,
		Title:       payload.Title,
		Description: payload.Description,
		DueDate:     due,
		MilestoneID: milestoneID,
		Position:    payload.Position,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toKanbanBoardResponse(board))
}

func (h *handler) UpdateKanbanCard(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}
	cardID, err := uuid.Parse(c.Params("cardID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateCardPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var due *time.Time
	if payload.DueDate != nil && *payload.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *payload.DueDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		due = &t
	}

	var milestoneID *uuid.UUID
	if payload.MilestoneID != nil && *payload.MilestoneID != "" {
		id, err := uuid.Parse(*payload.MilestoneID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		milestoneID = &id
	}

	board, err := h.service.UpdateKanbanCard(c.Context(), UpdateKanbanCardRequest{
		ProjectID:   projectID,
		UserID:      userID,
		CardID:      cardID,
		Title:       payload.Title,
		Description: payload.Description,
		DueDate:     due,
		MilestoneID: milestoneID,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toKanbanBoardResponse(board))
}

func (h *handler) MoveKanbanCard(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}
	cardID, err := uuid.Parse(c.Params("cardID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload moveCardPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateMoveCardPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	targetColumnID, err := uuid.Parse(payload.TargetColumnID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	board, err := h.service.MoveKanbanCard(c.Context(), MoveKanbanCardRequest{
		ProjectID:      projectID,
		UserID:         userID,
		CardID:         cardID,
		TargetColumnID: targetColumnID,
		TargetPosition: payload.TargetPosition,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toKanbanBoardResponse(board))
}

func (h *handler) DeleteKanbanCard(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}
	cardID, err := uuid.Parse(c.Params("cardID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload deleteCardPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	board, err := h.service.DeleteKanbanCard(c.Context(), DeleteKanbanCardRequest{
		ProjectID: projectID,
		UserID:    userID,
		CardID:    cardID,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toKanbanBoardResponse(board))
}

func (h *handler) GetKanbanBoard(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	board, err := h.service.GetKanbanBoard(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toKanbanBoardResponse(board))
}

// Dependencies

func (h *handler) CreateDependency(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload dependencyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateDependencyPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	dependsOn, err := uuid.Parse(payload.DependsOnProjectID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	dependency, err := h.service.CreateDependency(c.Context(), CreateDependencyRequest{
		ProjectID:          projectID,
		UserID:             userID,
		DependsOnProjectID: dependsOn,
		Type:               DependencyType(payload.Type),
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toDependencyResponse(dependency))
}

func (h *handler) ListDependencies(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	dependencies, err := h.service.ListDependencies(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]dependencyResponse, 0, len(dependencies))
	for _, dep := range dependencies {
		copy := dep
		resp = append(resp, toDependencyResponse(&copy))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) DeleteDependency(c *fiber.Ctx) error {
	dependencyID, err := uuid.Parse(c.Params("dependencyID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload deleteDependencyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteDependency(c.Context(), DeleteDependencyRequest{
		DependencyID: dependencyID,
		UserID:       userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": dependencyID.String()})
}

// Duplication

func (h *handler) DuplicateProject(c *fiber.Ctx) error {
	templateID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload duplicateProjectPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var status *ProjectStatus
	if payload.Status != nil {
		s := ProjectStatus(*payload.Status)
		status = &s
	}

	copyBoard := true
	if payload.CopyBoard != nil {
		copyBoard = *payload.CopyBoard
	}

	copyMilestones := true
	if payload.CopyMilestones != nil {
		copyMilestones = *payload.CopyMilestones
	}

	copyDeps := false
	if payload.CopyDeps != nil {
		copyDeps = *payload.CopyDeps
	}

	project, err := h.service.DuplicateProject(c.Context(), DuplicateProjectRequest{
		TemplateProjectID: templateID,
		UserID:            userID,
		Name:              payload.Name,
		Description:       payload.Description,
		Status:            status,
		HealthScore:       payload.HealthScore,
		MRR:               payload.MRR,
		CAC:               payload.CAC,
		LTV:               payload.LTV,
		ChurnRate:         payload.ChurnRate,
		CopyBoard:         copyBoard,
		CopyMilestones:    copyMilestones,
		CopyDependencies:  copyDeps,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toProjectResponse(project))
}

// Helpers

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromErrorCode(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("projects: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, authdomain.ErrCodeInvalidToken, fiber.Map{
		"message": "authentication required",
	})
}

func statusFromErrorCode(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidName, ErrCodeInvalidStatus, ErrCodeInvalidHealthScore, ErrCodeInvalidMetrics, ErrCodeInvalidDependency, ErrCodeInvalidVisibility, ErrCodeInvalidSectionType, ErrCodeInvalidTechCategory, ErrCodeInvalidDiagramType:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	case ErrCodeConflict:
		return fiber.StatusConflict
	case ErrCodeRepositoryFailure:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusInternalServerError
	}
}

func (h *handler) logWarn(message string) {
	if h.logger != nil {
		h.logger.Warn(message)
	}
}

func (h *handler) logError(message string, err error) {
	if h.logger != nil {
		h.logger.Error(message, slog.Any("error", err))
	}
}

func toProjectResponse(project *Project) projectResponse {
	return projectResponse{
		ID:          project.ID.String(),
		UserID:      project.UserID.String(),
		Name:        project.Name,
		Description: project.Description,
		Slug:        project.Slug,
		Status:      project.Status,
		HealthScore: project.HealthScore,
		MRR:         project.MRR,
		CAC:         project.CAC,
		LTV:         project.LTV,
		ChurnRate:   project.ChurnRate,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}

func toMilestoneResponse(m *Milestone) milestoneResponse {
	return milestoneResponse{
		ID:          m.ID.String(),
		ProjectID:   m.ProjectID.String(),
		Title:       m.Title,
		Description: m.Description,
		DueDate:     m.DueDate,
		Completed:   m.Completed,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toKanbanBoardResponse(board KanbanBoard) kanbanBoardResponse {
	resp := kanbanBoardResponse{
		ProjectID: board.ProjectID.String(),
	}

	for _, column := range board.Columns {
		colResp := kanbanColumnResponse{
			ID:        column.Column.ID.String(),
			ProjectID: column.Column.ProjectID.String(),
			Name:      column.Column.Name,
			WIPLimit:  column.Column.WIPLimit,
			Position:  column.Column.Position,
			CreatedAt: column.Column.CreatedAt,
			UpdatedAt: column.Column.UpdatedAt,
		}

		for _, card := range column.Cards {
			colResp.Cards = append(colResp.Cards, toKanbanCardResponse(card))
		}

		resp.Columns = append(resp.Columns, colResp)
	}

	return resp
}

func toKanbanCardResponse(card KanbanCard) kanbanCardResponse {
	var milestoneID *string
	if card.MilestoneID != nil {
		id := card.MilestoneID.String()
		milestoneID = &id
	}

	var due *time.Time
	if card.DueDate != nil {
		d := card.DueDate.UTC()
		due = &d
	}

	return kanbanCardResponse{
		ID:          card.ID.String(),
		ProjectID:   card.ProjectID.String(),
		ColumnID:    card.ColumnID.String(),
		MilestoneID: milestoneID,
		Title:       card.Title,
		Description: card.Description,
		DueDate:     due,
		Position:    card.Position,
		CreatedAt:   card.CreatedAt,
		UpdatedAt:   card.UpdatedAt,
	}
}

func toDependencyResponse(dep *ProjectDependency) dependencyResponse {
	return dependencyResponse{
		ID:                 dep.ID.String(),
		ProjectID:          dep.ProjectID.String(),
		DependsOnProjectID: dep.DependsOnProjectID.String(),
		Type:               dep.Type,
		CreatedAt:          dep.CreatedAt,
		UpdatedAt:          dep.UpdatedAt,
	}
}

// Documentation payloads

type createDocumentationPayload struct {
	Visibility string `json:"visibility"`
}

type updateDocumentationVisibilityPayload struct {
	Visibility string `json:"visibility"`
}

type createDocumentationSectionPayload struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Position *int   `json:"position"`
}

type updateDocumentationSectionPayload struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	Position *int    `json:"position"`
}

type reorderDocumentationSectionsPayload struct {
	SectionOrder []string `json:"section_order"`
}

// Technology payloads

type createTechnologyPayload struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Category string `json:"category"`
	Purpose  string `json:"purpose"`
	Link     string `json:"link"`
}

type updateTechnologyPayload struct {
	Name     *string `json:"name"`
	Version  *string `json:"version"`
	Category *string `json:"category"`
	Purpose  *string `json:"purpose"`
	Link     *string `json:"link"`
}

type bulkCreateTechnologiesPayload struct {
	Technologies []createTechnologyPayload `json:"technologies"`
}

type bulkUpdateTechnologiesPayload struct {
	Technologies []struct {
		TechID   string  `json:"tech_id"`
		Name     *string `json:"name"`
		Version  *string `json:"version"`
		Category *string `json:"category"`
		Purpose  *string `json:"purpose"`
		Link     *string `json:"link"`
	} `json:"technologies"`
}

// File Structure payloads

type createFileStructurePayload struct {
	Path        string  `json:"path"`
	Name        string  `json:"name"`
	IsDirectory bool    `json:"is_directory"`
	ParentID    *string `json:"parent_id"`
	Language    string  `json:"language"`
	LineCount   int     `json:"line_count"`
	Purpose     string  `json:"purpose"`
	Position    *int    `json:"position"`
}

type updateFileStructurePayload struct {
	Purpose   *string `json:"purpose"`
	LineCount *int    `json:"line_count"`
	Language  *string `json:"language"`
	Position  *int    `json:"position"`
}

type bulkCreateFileStructuresPayload struct {
	Structures []createFileStructurePayload `json:"structures"`
}

type bulkUpdateFileStructuresPayload struct {
	Structures []struct {
		FileStructureID string  `json:"file_structure_id"`
		Purpose         *string `json:"purpose"`
		LineCount       *int    `json:"line_count"`
		Language        *string `json:"language"`
		Position        *int    `json:"position"`
	} `json:"structures"`
}

// Architecture Diagram payloads

type createArchitectureDiagramPayload struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Format      string `json:"format"`
	ImageURL    string `json:"image_url"`
}

type updateArchitectureDiagramPayload struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Content     *string `json:"content"`
	ImageURL    *string `json:"image_url"`
}

// Documentation responses

type documentationResponse struct {
	ID         string                 `json:"id"`
	ProjectID string                 `json:"project_id"`
	Visibility DocumentationVisibility `json:"visibility"`
	Version    int                    `json:"version"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type documentationSectionResponse struct {
	ID             string                 `json:"id"`
	DocumentationID string                 `json:"documentation_id"`
	Type           DocumentationSectionType `json:"type"`
	Title          string                 `json:"title"`
	Content        string                 `json:"content"`
	Position       int                    `json:"position"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type technologyResponse struct {
	ID        string             `json:"id"`
	ProjectID string             `json:"project_id"`
	Name      string             `json:"name"`
	Version   string             `json:"version"`
	Category  TechnologyCategory `json:"category"`
	Purpose   string             `json:"purpose"`
	Link      string             `json:"link"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type fileStructureResponse struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	ParentID    *string    `json:"parent_id"`
	Path        string     `json:"path"`
	Name        string     `json:"name"`
	IsDirectory bool       `json:"is_directory"`
	Language    string     `json:"language"`
	LineCount   int        `json:"line_count"`
	Purpose     string     `json:"purpose"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type architectureDiagramResponse struct {
	ID          string                 `json:"id"`
	ProjectID   string                 `json:"project_id"`
	Type        ArchitectureDiagramType `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Format      string                 `json:"format"`
	ImageURL    string                 `json:"image_url"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Documentation handlers

func (h *handler) CreateDocumentation(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload createDocumentationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	visibility := DocumentationVisibility(payload.Visibility)
	if visibility == "" {
		visibility = VisibilityCollaborators
	}

	doc, err := h.service.CreateDocumentation(c.Context(), CreateDocumentationRequest{
		ProjectID:  projectID,
		UserID:     userID,
		Visibility: visibility,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toDocumentationResponse(doc))
}

func (h *handler) UpdateDocumentationVisibility(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateDocumentationVisibilityPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	doc, err := h.service.UpdateDocumentationVisibility(c.Context(), UpdateDocumentationVisibilityRequest{
		ProjectID:  projectID,
		UserID:     userID,
		Visibility: DocumentationVisibility(payload.Visibility),
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toDocumentationResponse(doc))
}

func (h *handler) GetDocumentation(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	doc, err := h.service.GetDocumentation(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toDocumentationResponse(doc))
}

func (h *handler) GetPublicDocumentation(c *fiber.Ctx) error {
	slug := strings.TrimSpace(c.Params("slug"))
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	doc, err := h.service.GetPublicDocumentation(c.Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toDocumentationResponse(doc))
}

func (h *handler) DeleteDocumentation(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteDocumentation(c.Context(), DeleteDocumentationRequest{
		ProjectID: projectID,
		UserID:    userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": projectID.String()})
}

// Documentation Section handlers

func (h *handler) CreateDocumentationSection(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload createDocumentationSectionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	section, err := h.service.CreateDocumentationSection(c.Context(), CreateDocumentationSectionRequest{
		ProjectID: projectID,
		UserID:    userID,
		Type:      DocumentationSectionType(payload.Type),
		Title:     payload.Title,
		Content:   payload.Content,
		Position:  payload.Position,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toDocumentationSectionResponse(section))
}

func (h *handler) UpdateDocumentationSection(c *fiber.Ctx) error {
	sectionID, err := uuid.Parse(c.Params("sectionID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateDocumentationSectionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	section, err := h.service.UpdateDocumentationSection(c.Context(), UpdateDocumentationSectionRequest{
		SectionID: sectionID,
		UserID:    userID,
		Title:     payload.Title,
		Content:   payload.Content,
		Position:  payload.Position,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toDocumentationSectionResponse(section))
}

func (h *handler) DeleteDocumentationSection(c *fiber.Ctx) error {
	sectionID, err := uuid.Parse(c.Params("sectionID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteDocumentationSection(c.Context(), DeleteDocumentationSectionRequest{
		SectionID: sectionID,
		UserID:    userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": sectionID.String()})
}

func (h *handler) ListDocumentationSections(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	sections, err := h.service.ListDocumentationSections(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]documentationSectionResponse, 0, len(sections))
	for _, section := range sections {
		s := section
		resp = append(resp, toDocumentationSectionResponse(&s))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) ReorderDocumentationSections(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload reorderDocumentationSectionsPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	sectionOrder := make([]uuid.UUID, 0, len(payload.SectionOrder))
	for _, raw := range payload.SectionOrder {
		id, err := uuid.Parse(raw)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		sectionOrder = append(sectionOrder, id)
	}

	sections, err := h.service.ReorderDocumentationSections(c.Context(), ReorderDocumentationSectionsRequest{
		ProjectID:    projectID,
		UserID:       userID,
		SectionOrder: sectionOrder,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]documentationSectionResponse, 0, len(sections))
	for _, section := range sections {
		s := section
		resp = append(resp, toDocumentationSectionResponse(&s))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// Technology handlers

func (h *handler) CreateTechnology(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload createTechnologyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	tech, err := h.service.CreateTechnology(c.Context(), CreateTechnologyRequest{
		ProjectID: projectID,
		UserID:    userID,
		Name:      payload.Name,
		Version:   payload.Version,
		Category:  TechnologyCategory(payload.Category),
		Purpose:   payload.Purpose,
		Link:      payload.Link,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toTechnologyResponse(tech))
}

func (h *handler) UpdateTechnology(c *fiber.Ctx) error {
	techID, err := uuid.Parse(c.Params("techID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateTechnologyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var category *TechnologyCategory
	if payload.Category != nil {
		c := TechnologyCategory(*payload.Category)
		category = &c
	}

	tech, err := h.service.UpdateTechnology(c.Context(), UpdateTechnologyRequest{
		TechID:   techID,
		UserID:   userID,
		Name:     payload.Name,
		Version:  payload.Version,
		Category: category,
		Purpose:  payload.Purpose,
		Link:     payload.Link,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toTechnologyResponse(tech))
}

func (h *handler) DeleteTechnology(c *fiber.Ctx) error {
	techID, err := uuid.Parse(c.Params("techID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteTechnology(c.Context(), DeleteTechnologyRequest{
		TechID: techID,
		UserID: userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": techID.String()})
}

func (h *handler) ListTechnologies(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	technologies, err := h.service.ListTechnologies(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]technologyResponse, 0, len(technologies))
	for _, tech := range technologies {
		t := tech
		resp = append(resp, toTechnologyResponse(&t))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) BulkCreateTechnologies(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload bulkCreateTechnologiesPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	techReqs := make([]CreateTechnologyRequest, 0, len(payload.Technologies))
	for _, tech := range payload.Technologies {
		techReqs = append(techReqs, CreateTechnologyRequest{
			ProjectID: projectID,
			UserID:    userID,
			Name:      tech.Name,
			Version:   tech.Version,
			Category:  TechnologyCategory(tech.Category),
			Purpose:   tech.Purpose,
			Link:      tech.Link,
		})
	}

	technologies, err := h.service.BulkCreateTechnologies(c.Context(), BulkCreateTechnologiesRequest{
		ProjectID:    projectID,
		UserID:       userID,
		Technologies: techReqs,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]technologyResponse, 0, len(technologies))
	for _, tech := range technologies {
		t := tech
		resp = append(resp, toTechnologyResponse(&t))
	}

	return response.Success(c, fiber.StatusCreated, resp)
}

func (h *handler) BulkUpdateTechnologies(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload bulkUpdateTechnologiesPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	techReqs := make([]UpdateTechnologyRequest, 0, len(payload.Technologies))
	for _, tech := range payload.Technologies {
		techID, err := uuid.Parse(tech.TechID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}

		var category *TechnologyCategory
		if tech.Category != nil {
			c := TechnologyCategory(*tech.Category)
			category = &c
		}

		techReqs = append(techReqs, UpdateTechnologyRequest{
			TechID:   techID,
			UserID:   userID,
			Name:     tech.Name,
			Version:  tech.Version,
			Category: category,
			Purpose:  tech.Purpose,
			Link:     tech.Link,
		})
	}

	technologies, err := h.service.BulkUpdateTechnologies(c.Context(), BulkUpdateTechnologiesRequest{
		ProjectID:    projectID,
		UserID:       userID,
		Technologies: techReqs,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]technologyResponse, 0, len(technologies))
	for _, tech := range technologies {
		t := tech
		resp = append(resp, toTechnologyResponse(&t))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// File Structure handlers

func (h *handler) CreateFileStructure(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload createFileStructurePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var parentID *uuid.UUID
	if payload.ParentID != nil && *payload.ParentID != "" {
		id, err := uuid.Parse(*payload.ParentID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		parentID = &id
	}

	fs, err := h.service.CreateFileStructure(c.Context(), CreateFileStructureRequest{
		ProjectID:   projectID,
		UserID:     userID,
		Path:       payload.Path,
		Name:       payload.Name,
		IsDirectory: payload.IsDirectory,
		ParentID:   parentID,
		Language:   payload.Language,
		LineCount:  payload.LineCount,
		Purpose:    payload.Purpose,
		Position:   payload.Position,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toFileStructureResponse(fs))
}

func (h *handler) UpdateFileStructure(c *fiber.Ctx) error {
	fsID, err := uuid.Parse(c.Params("fileStructureID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateFileStructurePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	fs, err := h.service.UpdateFileStructure(c.Context(), UpdateFileStructureRequest{
		FileStructureID: fsID,
		UserID:          userID,
		Purpose:         payload.Purpose,
		LineCount:       payload.LineCount,
		Language:        payload.Language,
		Position:        payload.Position,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toFileStructureResponse(fs))
}

func (h *handler) DeleteFileStructure(c *fiber.Ctx) error {
	fsID, err := uuid.Parse(c.Params("fileStructureID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteFileStructure(c.Context(), DeleteFileStructureRequest{
		FileStructureID: fsID,
		UserID:          userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": fsID.String()})
}

func (h *handler) ListFileStructures(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	structures, err := h.service.ListFileStructures(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]fileStructureResponse, 0, len(structures))
	for _, fs := range structures {
		f := fs
		resp = append(resp, toFileStructureResponse(&f))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) BulkCreateFileStructures(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload bulkCreateFileStructuresPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	fsReqs := make([]CreateFileStructureRequest, 0, len(payload.Structures))
	for _, fs := range payload.Structures {
		var parentID *uuid.UUID
		if fs.ParentID != nil && *fs.ParentID != "" {
			id, err := uuid.Parse(*fs.ParentID)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
			}
			parentID = &id
		}

		fsReqs = append(fsReqs, CreateFileStructureRequest{
			ProjectID:   projectID,
			UserID:      userID,
			Path:        fs.Path,
			Name:        fs.Name,
			IsDirectory: fs.IsDirectory,
			ParentID:    parentID,
			Language:    fs.Language,
			LineCount:   fs.LineCount,
			Purpose:     fs.Purpose,
			Position:    fs.Position,
		})
	}

	structures, err := h.service.BulkCreateFileStructures(c.Context(), BulkCreateFileStructuresRequest{
		ProjectID: projectID,
		UserID:    userID,
		Structures: fsReqs,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]fileStructureResponse, 0, len(structures))
	for _, fs := range structures {
		f := fs
		resp = append(resp, toFileStructureResponse(&f))
	}

	return response.Success(c, fiber.StatusCreated, resp)
}

func (h *handler) BulkUpdateFileStructures(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload bulkUpdateFileStructuresPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	fsReqs := make([]UpdateFileStructureRequest, 0, len(payload.Structures))
	for _, fs := range payload.Structures {
		fsID, err := uuid.Parse(fs.FileStructureID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}

		fsReqs = append(fsReqs, UpdateFileStructureRequest{
			FileStructureID: fsID,
			UserID:          userID,
			Purpose:         fs.Purpose,
			LineCount:       fs.LineCount,
			Language:        fs.Language,
			Position:        fs.Position,
		})
	}

	structures, err := h.service.BulkUpdateFileStructures(c.Context(), BulkUpdateFileStructuresRequest{
		ProjectID: projectID,
		UserID:    userID,
		Structures: fsReqs,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]fileStructureResponse, 0, len(structures))
	for _, fs := range structures {
		f := fs
		resp = append(resp, toFileStructureResponse(&f))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// Architecture Diagram handlers

func (h *handler) CreateArchitectureDiagram(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload createArchitectureDiagramPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	diagram, err := h.service.CreateArchitectureDiagram(c.Context(), CreateArchitectureDiagramRequest{
		ProjectID:   projectID,
		UserID:      userID,
		Type:        ArchitectureDiagramType(payload.Type),
		Title:       payload.Title,
		Description: payload.Description,
		Content:     payload.Content,
		Format:      payload.Format,
		ImageURL:    payload.ImageURL,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toArchitectureDiagramResponse(diagram))
}

func (h *handler) UpdateArchitectureDiagram(c *fiber.Ctx) error {
	diagramID, err := uuid.Parse(c.Params("diagramID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateArchitectureDiagramPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	diagram, err := h.service.UpdateArchitectureDiagram(c.Context(), UpdateArchitectureDiagramRequest{
		DiagramID:   diagramID,
		UserID:      userID,
		Title:       payload.Title,
		Description: payload.Description,
		Content:     payload.Content,
		ImageURL:    payload.ImageURL,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toArchitectureDiagramResponse(diagram))
}

func (h *handler) DeleteArchitectureDiagram(c *fiber.Ctx) error {
	diagramID, err := uuid.Parse(c.Params("diagramID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteArchitectureDiagram(c.Context(), DeleteArchitectureDiagramRequest{
		DiagramID: diagramID,
		UserID:    userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": diagramID.String()})
}

func (h *handler) ListArchitectureDiagrams(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	diagrams, err := h.service.ListArchitectureDiagrams(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]architectureDiagramResponse, 0, len(diagrams))
	for _, diagram := range diagrams {
		d := diagram
		resp = append(resp, toArchitectureDiagramResponse(&d))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) GetArchitectureDiagram(c *fiber.Ctx) error {
	diagramID, err := uuid.Parse(c.Params("diagramID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	diagram, err := h.service.GetArchitectureDiagram(c.Context(), diagramID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toArchitectureDiagramResponse(diagram))
}

// Response converters

func toDocumentationResponse(doc *ProjectDocumentation) documentationResponse {
	return documentationResponse{
		ID:         doc.ID.String(),
		ProjectID:  doc.ProjectID.String(),
		Visibility: doc.Visibility,
		Version:    doc.Version,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
	}
}

func toDocumentationSectionResponse(section *DocumentationSection) documentationSectionResponse {
	return documentationSectionResponse{
		ID:             section.ID.String(),
		DocumentationID: section.DocumentationID.String(),
		Type:           section.Type,
		Title:          section.Title,
		Content:        section.Content,
		Position:       section.Position,
		CreatedAt:      section.CreatedAt,
		UpdatedAt:      section.UpdatedAt,
	}
}

func toTechnologyResponse(tech *ProjectTechnology) technologyResponse {
	return technologyResponse{
		ID:        tech.ID.String(),
		ProjectID: tech.ProjectID.String(),
		Name:      tech.Name,
		Version:   tech.Version,
		Category:  tech.Category,
		Purpose:   tech.Purpose,
		Link:      tech.Link,
		CreatedAt: tech.CreatedAt,
		UpdatedAt: tech.UpdatedAt,
	}
}

func toFileStructureResponse(fs *ProjectFileStructure) fileStructureResponse {
	var parentID *string
	if fs.ParentID != nil {
		id := fs.ParentID.String()
		parentID = &id
	}

	return fileStructureResponse{
		ID:          fs.ID.String(),
		ProjectID:   fs.ProjectID.String(),
		ParentID:    parentID,
		Path:        fs.Path,
		Name:        fs.Name,
		IsDirectory: fs.IsDirectory,
		Language:    fs.Language,
		LineCount:   fs.LineCount,
		Purpose:     fs.Purpose,
		Position:    fs.Position,
		CreatedAt:   fs.CreatedAt,
		UpdatedAt:   fs.UpdatedAt,
	}
}

func toArchitectureDiagramResponse(diagram *ProjectArchitectureDiagram) architectureDiagramResponse {
	return architectureDiagramResponse{
		ID:          diagram.ID.String(),
		ProjectID:   diagram.ProjectID.String(),
		Type:        diagram.Type,
		Title:       diagram.Title,
		Description: diagram.Description,
		Content:     diagram.Content,
		Format:      diagram.Format,
		ImageURL:    diagram.ImageURL,
		CreatedAt:   diagram.CreatedAt,
		UpdatedAt:   diagram.UpdatedAt,
	}
}
