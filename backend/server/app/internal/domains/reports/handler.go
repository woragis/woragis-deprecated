package reports

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes report-related endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler builds a Handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type generateSummaryPayload struct {
	SendEmail    bool   `json:"send_email"`
	EmailAddress string `json:"email_address"`
	SendWhatsApp bool   `json:"send_whatsapp"`
	PhoneNumber  string `json:"phone_number"`
	AgentAlias   string `json:"agent_alias"`
}

type createDefinitionPayload struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sections    map[string]any `json:"sections"`
	Filters     map[string]any `json:"filters"`
	Favorite    bool           `json:"favorite"`
}

type updateDefinitionPayload = createDefinitionPayload

type bulkDefinitionPayload struct {
	DefinitionIDs []string `json:"definition_ids"`
}

type toggleFavoritePayload struct {
	DefinitionID string `json:"definition_id"`
	Favorite     bool   `json:"favorite"`
}

type schedulePayload struct {
	Cron      string         `json:"cron"`
	Frequency string         `json:"frequency"`
	Timezone  string         `json:"timezone"`
	NextRun   string         `json:"next_run"`
	Enabled   *bool          `json:"enabled"`
	Meta      map[string]any `json:"meta"`
}

type togglePayload struct {
	Enabled bool `json:"enabled"`
}

type deliveryPayload struct {
	Channel  string         `json:"channel"`
	Target   string         `json:"target"`
	Template map[string]any `json:"template"`
	Enabled  *bool          `json:"enabled"`
}

type bulkRunPayload struct {
	DefinitionIDs []string       `json:"definition_ids"`
	Metadata      map[string]any `json:"metadata"`
}

type definitionResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	IsFavorite  bool           `json:"is_favorite"`
	Sections    map[string]any `json:"sections"`
	Filters     map[string]any `json:"filters"`
	ArchivedAt  *string        `json:"archived_at,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type scheduleResponse struct {
	ID        string         `json:"id"`
	ReportID  string         `json:"report_id"`
	Cron      string         `json:"cron"`
	Frequency string         `json:"frequency"`
	Timezone  string         `json:"timezone"`
	NextRun   *string        `json:"next_run,omitempty"`
	LastRunAt *string        `json:"last_run_at,omitempty"`
	Enabled   bool           `json:"enabled"`
	Meta      map[string]any `json:"meta,omitempty"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

type deliveryResponse struct {
	ID        string         `json:"id"`
	ReportID  string         `json:"report_id"`
	Channel   string         `json:"channel"`
	Target    string         `json:"target"`
	Template  map[string]any `json:"template,omitempty"`
	Enabled   bool           `json:"enabled"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

type runResponse struct {
	ID             string         `json:"id"`
	ReportID       string         `json:"report_id"`
	Status         string         `json:"status"`
	StartedAt      *string        `json:"started_at,omitempty"`
	CompletedAt    *string        `json:"completed_at,omitempty"`
	OutputLocation string         `json:"output_location,omitempty"`
	ErrorMessage   string         `json:"error_message,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
}

type definitionDetailResponse struct {
	Definition definitionResponse `json:"definition"`
	Schedules  []scheduleResponse `json:"schedules"`
	Deliveries []deliveryResponse `json:"deliveries"`
}

// PostSummary generates a summary and dispatches it.
func (h *Handler) PostSummary(c *fiber.Ctx) error {
	var payload generateSummaryPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	summary, err := h.service.GenerateSummary(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	opts := DispatchOptions{
		SendEmail:    payload.SendEmail,
		EmailAddress: payload.EmailAddress,
		SendWhatsApp: payload.SendWhatsApp,
		PhoneNumber:  payload.PhoneNumber,
		AgentAlias:   payload.AgentAlias,
	}

	if err := h.service.DispatchSummary(c.Context(), summary, opts); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, summary)
}

// CreateDefinition handles POST /reports
func (h *Handler) CreateDefinition(c *fiber.Ctx) error {
	var payload createDefinitionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	def, err := h.service.CreateDefinition(c.Context(), CreateDefinitionRequest{
		UserID:      userID,
		Name:        payload.Name,
		Description: payload.Description,
		Sections:    payload.Sections,
		Filters:     payload.Filters,
		Favorite:    payload.Favorite,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toDefinitionResponse(*def))
}

// UpdateDefinition handles PUT /reports/:id
func (h *Handler) UpdateDefinition(c *fiber.Ctx) error {
	defID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateDefinitionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	def, err := h.service.UpdateDefinition(c.Context(), UpdateDefinitionRequest{
		UserID:       userID,
		DefinitionID: defID,
		Name:         payload.Name,
		Description:  payload.Description,
		Sections:     payload.Sections,
		Filters:      payload.Filters,
		Favorite:     payload.Favorite,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toDefinitionResponse(*def))
}

// ListDefinitions handles GET /reports
func (h *Handler) ListDefinitions(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	limit := c.QueryInt("limit", 25)
	offset := c.QueryInt("offset", 0)
	includeArchived := strings.ToLower(c.Query("include_archived")) == "true"
	onlyFavorites := strings.ToLower(c.Query("favorites")) == "true"

	defs, err := h.service.ListDefinitions(c.Context(), ListDefinitionsRequest{
		UserID:          userID,
		Search:          c.Query("search"),
		IncludeArchived: includeArchived,
		OnlyFavorites:   onlyFavorites,
		Channel:         c.Query("channel"),
		Limit:           limit,
		Offset:          offset,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]definitionResponse, 0, len(defs))
	for _, def := range defs {
		resp = append(resp, toDefinitionResponse(def))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// GetDefinition handles GET /reports/:id
func (h *Handler) GetDefinition(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	defID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	detail, err := h.service.GetDefinition(c.Context(), userID, defID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toDefinitionDetailResponse(detail))
}

// ArchiveDefinitions handles POST /reports/archive
func (h *Handler) ArchiveDefinitions(c *fiber.Ctx) error {
	req, err := h.parseBulkDefinitionPayload(c)
	if err != nil {
		return err
	}
	if err := h.service.ArchiveDefinitions(c.Context(), req); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "archived"})
}

// RestoreDefinitions handles POST /reports/restore
func (h *Handler) RestoreDefinitions(c *fiber.Ctx) error {
	req, err := h.parseBulkDefinitionPayload(c)
	if err != nil {
		return err
	}
	if err := h.service.RestoreDefinitions(c.Context(), req); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "restored"})
}

// DeleteDefinitions handles POST /reports/delete
func (h *Handler) DeleteDefinitions(c *fiber.Ctx) error {
	req, err := h.parseBulkDefinitionPayload(c)
	if err != nil {
		return err
	}
	if err := h.service.DeleteDefinitions(c.Context(), req); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deleted"})
}

// ToggleFavorite handles POST /reports/favorite
func (h *Handler) ToggleFavorite(c *fiber.Ctx) error {
	var payload toggleFavoritePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	defID, err := uuid.Parse(payload.DefinitionID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.ToggleFavorite(c.Context(), ToggleFavoriteRequest{
		UserID:       userID,
		DefinitionID: defID,
		Favorite:     payload.Favorite,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "updated"})
}

// CreateSchedule handles POST /reports/:id/schedules
func (h *Handler) CreateSchedule(c *fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload schedulePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var nextRun *time.Time
	if payload.NextRun != "" {
		parsed, parseErr := time.Parse(time.RFC3339, payload.NextRun)
		if parseErr != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid next_run"})
		}
		nextRun = &parsed
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	enabled := true
	if payload.Enabled != nil {
		enabled = *payload.Enabled
	}

	schedule, err := h.service.CreateSchedule(c.Context(), CreateScheduleRequest{
		UserID:    userID,
		ReportID:  reportID,
		Cron:      payload.Cron,
		Frequency: payload.Frequency,
		Timezone:  payload.Timezone,
		NextRun:   nextRun,
		Enabled:   enabled,
		Meta:      payload.Meta,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toScheduleResponse(*schedule))
}

// UpdateSchedule handles PUT /reports/schedules/:scheduleID
func (h *Handler) UpdateSchedule(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("scheduleID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload schedulePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	var nextRun *time.Time
	if payload.NextRun != "" {
		parsed, parseErr := time.Parse(time.RFC3339, payload.NextRun)
		if parseErr != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid next_run"})
		}
		nextRun = &parsed
	}

	enabled := true
	if payload.Enabled != nil {
		enabled = *payload.Enabled
	}

	schedule, err := h.service.UpdateSchedule(c.Context(), UpdateScheduleRequest{
		UserID:     userID,
		ScheduleID: scheduleID,
		Cron:       payload.Cron,
		Frequency:  payload.Frequency,
		Timezone:   payload.Timezone,
		NextRun:    nextRun,
		Enabled:    enabled,
		Meta:       payload.Meta,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toScheduleResponse(*schedule))
}

// ToggleSchedule handles POST /reports/schedules/:scheduleID/toggle
func (h *Handler) ToggleSchedule(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("scheduleID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload togglePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.ToggleSchedule(c.Context(), ToggleScheduleRequest{
		UserID:     userID,
		ScheduleID: scheduleID,
		Enabled:    payload.Enabled,
	}); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "updated"})
}

// DeleteSchedule handles DELETE /reports/schedules/:scheduleID
func (h *Handler) DeleteSchedule(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("scheduleID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteSchedule(c.Context(), userID, scheduleID); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deleted"})
}

// ListSchedules handles GET /reports/:id/schedules
func (h *Handler) ListSchedules(c *fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	schedules, err := h.service.ListSchedules(c.Context(), userID, reportID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]scheduleResponse, 0, len(schedules))
	for _, schedule := range schedules {
		resp = append(resp, toScheduleResponse(schedule))
	}
	return response.Success(c, fiber.StatusOK, resp)
}

// CreateDelivery handles POST /reports/:id/deliveries
func (h *Handler) CreateDelivery(c *fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload deliveryPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	enabled := true
	if payload.Enabled != nil {
		enabled = *payload.Enabled
	}

	delivery, err := h.service.CreateDelivery(c.Context(), CreateDeliveryRequest{
		UserID:   userID,
		ReportID: reportID,
		Channel:  payload.Channel,
		Target:   payload.Target,
		Template: payload.Template,
		Enabled:  enabled,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toDeliveryResponse(*delivery))
}

// UpdateDelivery handles PUT /reports/deliveries/:deliveryID
func (h *Handler) UpdateDelivery(c *fiber.Ctx) error {
	deliveryID, err := uuid.Parse(c.Params("deliveryID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload deliveryPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	enabled := true
	if payload.Enabled != nil {
		enabled = *payload.Enabled
	}

	delivery, err := h.service.UpdateDelivery(c.Context(), UpdateDeliveryRequest{
		UserID:     userID,
		DeliveryID: deliveryID,
		Channel:    payload.Channel,
		Target:     payload.Target,
		Template:   payload.Template,
		Enabled:    enabled,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toDeliveryResponse(*delivery))
}

// ToggleDelivery handles POST /reports/deliveries/:deliveryID/toggle
func (h *Handler) ToggleDelivery(c *fiber.Ctx) error {
	deliveryID, err := uuid.Parse(c.Params("deliveryID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload togglePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.ToggleDelivery(c.Context(), ToggleDeliveryRequest{
		UserID:     userID,
		DeliveryID: deliveryID,
		Enabled:    payload.Enabled,
	}); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "updated"})
}

// DeleteDelivery handles DELETE /reports/deliveries/:deliveryID
func (h *Handler) DeleteDelivery(c *fiber.Ctx) error {
	deliveryID, err := uuid.Parse(c.Params("deliveryID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteDelivery(c.Context(), userID, deliveryID); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deleted"})
}

// ListDeliveries handles GET /reports/:id/deliveries
func (h *Handler) ListDeliveries(c *fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	deliveries, err := h.service.ListDeliveries(c.Context(), userID, reportID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]deliveryResponse, 0, len(deliveries))
	for _, delivery := range deliveries {
		resp = append(resp, toDeliveryResponse(delivery))
	}
	return response.Success(c, fiber.StatusOK, resp)
}

// QueueRuns handles POST /reports/runs/bulk
func (h *Handler) QueueRuns(c *fiber.Ctx) error {
	var payload bulkRunPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	defIDs, err := parseUUIDList(payload.DefinitionIDs)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid definition_ids"})
	}

	runs, err := h.service.QueueRuns(c.Context(), BulkRunRequest{
		UserID:        userID,
		DefinitionIDs: defIDs,
		Metadata:      payload.Metadata,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]runResponse, 0, len(runs))
	for _, run := range runs {
		resp = append(resp, toRunResponse(run))
	}
	return response.Success(c, fiber.StatusAccepted, resp)
}

// ListRuns handles GET /reports/:id/runs
func (h *Handler) ListRuns(c *fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	limit := c.QueryInt("limit", 25)
	offset := c.QueryInt("offset", 0)

	runs, err := h.service.ListRuns(c.Context(), ListRunsRequest{
		UserID:   userID,
		ReportID: reportID,
		Status:   c.Query("status"),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]runResponse, 0, len(runs))
	for _, run := range runs {
		resp = append(resp, toRunResponse(run))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *Handler) parseBulkDefinitionPayload(c *fiber.Ctx) (BulkDefinitionRequest, error) {
	var payload bulkDefinitionPayload
	if err := c.BodyParser(&payload); err != nil {
		return BulkDefinitionRequest{}, response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}
	if len(payload.DefinitionIDs) == 0 {
		return BulkDefinitionRequest{}, response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "definition_ids required"})
	}

	defIDs, err := parseUUIDList(payload.DefinitionIDs)
	if err != nil {
		return BulkDefinitionRequest{}, response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid definition_ids"})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return BulkDefinitionRequest{}, response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	return BulkDefinitionRequest{
		UserID:        userID,
		DefinitionIDs: defIDs,
	}, nil
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("reports: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromError(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidSchedule, ErrCodeInvalidDelivery:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
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

func toDefinitionResponse(def ReportDefinition) definitionResponse {
	var archivedAt *string
	if def.ArchivedAt != nil {
		str := def.ArchivedAt.Format(time.RFC3339)
		archivedAt = &str
	}

	return definitionResponse{
		ID:          def.ID.String(),
		Name:        def.Name,
		Description: def.Description,
		IsFavorite:  def.IsFavorite,
		Sections:    def.Sections.Data,
		Filters:     def.Filters.Data,
		ArchivedAt:  archivedAt,
		CreatedAt:   def.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   def.UpdatedAt.Format(time.RFC3339),
	}
}

func toScheduleResponse(schedule ReportSchedule) scheduleResponse {
	var nextRun, lastRun *string
	if schedule.NextRun != nil {
		str := schedule.NextRun.Format(time.RFC3339)
		nextRun = &str
	}
	if schedule.LastRunAt != nil {
		str := schedule.LastRunAt.Format(time.RFC3339)
		lastRun = &str
	}

	return scheduleResponse{
		ID:        schedule.ID.String(),
		ReportID:  schedule.ReportID.String(),
		Cron:      schedule.Cron,
		Frequency: schedule.Frequency,
		Timezone:  schedule.Timezone,
		NextRun:   nextRun,
		LastRunAt: lastRun,
		Enabled:   schedule.Enabled,
		Meta:      schedule.Meta.Data,
		CreatedAt: schedule.CreatedAt.Format(time.RFC3339),
		UpdatedAt: schedule.UpdatedAt.Format(time.RFC3339),
	}
}

func toDeliveryResponse(delivery ReportDelivery) deliveryResponse {
	return deliveryResponse{
		ID:        delivery.ID.String(),
		ReportID:  delivery.ReportID.String(),
		Channel:   delivery.Channel,
		Target:    delivery.Target,
		Template:  delivery.Template.Data,
		Enabled:   delivery.Enabled,
		CreatedAt: delivery.CreatedAt.Format(time.RFC3339),
		UpdatedAt: delivery.UpdatedAt.Format(time.RFC3339),
	}
}

func toRunResponse(run ReportRun) runResponse {
	var started, completed *string
	if run.StartedAt != nil {
		str := run.StartedAt.Format(time.RFC3339)
		started = &str
	}
	if run.CompletedAt != nil {
		str := run.CompletedAt.Format(time.RFC3339)
		completed = &str
	}

	return runResponse{
		ID:             run.ID.String(),
		ReportID:       run.ReportID.String(),
		Status:         run.Status,
		StartedAt:      started,
		CompletedAt:    completed,
		OutputLocation: run.OutputLocation,
		ErrorMessage:   run.ErrorMessage,
		Metadata:       run.Metadata.Data,
		CreatedAt:      run.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      run.UpdatedAt.Format(time.RFC3339),
	}
}

func toDefinitionDetailResponse(detail DefinitionDetail) definitionDetailResponse {
	schedules := make([]scheduleResponse, 0, len(detail.Schedules))
	for _, schedule := range detail.Schedules {
		schedules = append(schedules, toScheduleResponse(schedule))
	}

	deliveries := make([]deliveryResponse, 0, len(detail.Deliveries))
	for _, delivery := range detail.Deliveries {
		deliveries = append(deliveries, toDeliveryResponse(delivery))
	}

	return definitionDetailResponse{
		Definition: toDefinitionResponse(detail.Definition),
		Schedules:  schedules,
		Deliveries: deliveries,
	}
}

func parseUUIDList(values []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(values))
	for _, raw := range values {
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
