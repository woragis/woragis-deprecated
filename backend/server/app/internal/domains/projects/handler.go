package projects

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes project endpoints.
type Handler interface {
	CreateProject(c *fiber.Ctx) error
	ListProjects(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
	UpdateMetrics(c *fiber.Ctx) error
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
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a project handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
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

	return response.Success(c, fiber.StatusCreated, toProjectResponse(project))
}

func (h *handler) ListProjects(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
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

func (h *handler) UpdateStatus(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateStatusPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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

func (h *handler) AddMilestone(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload addMilestonePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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
	case ErrCodeInvalidPayload, ErrCodeInvalidName, ErrCodeInvalidStatus, ErrCodeInvalidHealthScore, ErrCodeInvalidMetrics, ErrCodeInvalidDependency:
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
