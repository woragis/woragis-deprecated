package projects

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes project endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler constructs a project handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type createProjectPayload struct {
	UserID      string  `json:"user_id"`
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
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type updateMetricsPayload struct {
	UserID      string  `json:"user_id"`
	HealthScore int     `json:"health_score"`
	MRR         float64 `json:"mrr"`
	CAC         float64 `json:"cac"`
	LTV         float64 `json:"ltv"`
	ChurnRate   float64 `json:"churn_rate"`
}

type addMilestonePayload struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type toggleMilestonePayload struct {
	UserID    string `json:"user_id"`
	Completed bool   `json:"completed"`
}

type projectResponse struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
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

// CreateProject handles creation.
func (h *Handler) CreateProject(c *fiber.Ctx) error {
	var payload createProjectPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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

	return response.Success(c, fiber.StatusCreated, toProjectResponse(project))
}

// ListProjects handles listing.
func (h *Handler) ListProjects(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	projects, err := h.service.ListProjects(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]projectResponse, 0, len(projects))
	for _, project := range projects {
		p := project
		resp = append(resp, toProjectResponse(&p))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// UpdateStatus handles status updates.
func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateStatusPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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

// UpdateMetrics handles KPI updates.
func (h *Handler) UpdateMetrics(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateMetricsPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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

// AddMilestone handles milestone creation.
func (h *Handler) AddMilestone(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload addMilestonePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var dueDate time.Time
	if payload.DueDate != "" {
		if dueDate, err = time.Parse(time.RFC3339, payload.DueDate); err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	milestone, err := h.service.AddMilestone(c.Context(), AddMilestoneRequest{
		ProjectID:   projectID,
		UserID:      userID,
		Title:       payload.Title,
		Description: payload.Description,
		DueDate:     dueDate,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toMilestoneResponse(milestone))
}

// ToggleMilestoneCompletion toggles milestone completion.
func (h *Handler) ToggleMilestoneCompletion(c *fiber.Ctx) error {
	milestoneID, err := uuid.Parse(c.Params("milestoneID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload toggleMilestonePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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

// ListMilestones handles milestone listing.
func (h *Handler) ListMilestones(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromErrorCode(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("projects: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromErrorCode(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidName, ErrCodeInvalidStatus, ErrCodeInvalidHealthScore, ErrCodeInvalidMetrics:
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

func (h *Handler) logWarn(message string) {
	if h.logger != nil {
		h.logger.Warn(message)
	}
}

func (h *Handler) logError(message string, err error) {
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
