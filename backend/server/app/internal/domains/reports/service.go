package reports

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	chatsdomain "github.com/woragis/backend/server/app/internal/domains/chats"
	financesdomain "github.com/woragis/backend/server/app/internal/domains/finances"
	ideasdomain "github.com/woragis/backend/server/app/internal/domains/ideas"
	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	"github.com/woragis/backend/server/app/internal/workers/notifications"
)

// IdeaRepository describes the subset of methods needed from ideas.
type IdeaRepository interface {
	ListIdeas(ctx context.Context, userID uuid.UUID) ([]ideasdomain.Idea, error)
}

// ProjectRepository describes the subset needed from projects.
type ProjectRepository interface {
	ListProjects(ctx context.Context, userID uuid.UUID) ([]projectsdomain.Project, error)
}

// FinanceRepository describes the subset needed from finances.
type FinanceRepository interface {
	AggregateSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (financesdomain.Summary, error)
}

// ChatRepository describes the subset from chats domain.
type ChatRepository interface {
	ListConversations(ctx context.Context, userID uuid.UUID) ([]chatsdomain.Conversation, error)
}

// Publisher defines the notification publisher contract.
type Publisher interface {
	PublishEmailReport(ctx context.Context, env notifications.ReportEnvelope) error
	PublishWhatsAppReport(ctx context.Context, env notifications.ReportEnvelope) error
}

// Service orchestrates report generation and dispatch.
type Service struct {
	repo         Repository
	ideasRepo    IdeaRepository
	projectsRepo ProjectRepository
	financeRepo  FinanceRepository
	chatsRepo    ChatRepository
	publisher    Publisher
	logger       *slog.Logger
}

// NewService builds a new reports service.
func NewService(
	repo Repository,
	ideasRepo IdeaRepository,
	projectsRepo ProjectRepository,
	financeRepo FinanceRepository,
	chatsRepo ChatRepository,
	publisher Publisher,
	logger *slog.Logger,
) *Service {
	return &Service{
		repo:         repo,
		ideasRepo:    ideasRepo,
		projectsRepo: projectsRepo,
		financeRepo:  financeRepo,
		chatsRepo:    chatsRepo,
		publisher:    publisher,
		logger:       logger,
	}
}

const (
	RunStatusPending   = "pending"
	RunStatusRunning   = "running"
	RunStatusCompleted = "completed"
	RunStatusFailed    = "failed"
)

// DefinitionDetail aggregates a definition and its related entities.
type DefinitionDetail struct {
	Definition ReportDefinition `json:"definition"`
	Schedules  []ReportSchedule `json:"schedules"`
	Deliveries []ReportDelivery `json:"deliveries"`
}

// CreateDefinitionRequest defines inputs required to create a report definition.
type CreateDefinitionRequest struct {
	UserID      uuid.UUID
	Name        string
	Description string
	Sections    map[string]any
	Filters     map[string]any
	Favorite    bool
}

// UpdateDefinitionRequest defines inputs to update a definition.
type UpdateDefinitionRequest struct {
	UserID       uuid.UUID
	DefinitionID uuid.UUID
	Name         string
	Description  string
	Sections     map[string]any
	Filters      map[string]any
	Favorite     bool
}

// ListDefinitionsRequest contains filters for listing definitions.
type ListDefinitionsRequest struct {
	UserID          uuid.UUID
	Search          string
	IncludeArchived bool
	OnlyFavorites   bool
	Channel         string
	Limit           int
	Offset          int
}

// BulkDefinitionRequest is used for batch operations.
type BulkDefinitionRequest struct {
	UserID        uuid.UUID
	DefinitionIDs []uuid.UUID
}

// ToggleFavoriteRequest toggles the favorite flag.
type ToggleFavoriteRequest struct {
	UserID       uuid.UUID
	DefinitionID uuid.UUID
	Favorite     bool
}

// CreateScheduleRequest defines schedule creation input.
type CreateScheduleRequest struct {
	UserID    uuid.UUID
	ReportID  uuid.UUID
	Cron      string
	Frequency string
	Timezone  string
	NextRun   *time.Time
	Enabled   bool
	Meta      map[string]any
}

// UpdateScheduleRequest updates schedule properties.
type UpdateScheduleRequest struct {
	UserID     uuid.UUID
	ScheduleID uuid.UUID
	Cron       string
	Frequency  string
	Timezone   string
	NextRun    *time.Time
	Enabled    bool
	Meta       map[string]any
}

// ToggleScheduleRequest toggles schedule enabled state.
type ToggleScheduleRequest struct {
	UserID     uuid.UUID
	ScheduleID uuid.UUID
	Enabled    bool
}

// CreateDeliveryRequest defines delivery configuration.
type CreateDeliveryRequest struct {
	UserID   uuid.UUID
	ReportID uuid.UUID
	Channel  string
	Target   string
	Template map[string]any
	Enabled  bool
}

// UpdateDeliveryRequest updates delivery values.
type UpdateDeliveryRequest struct {
	UserID     uuid.UUID
	DeliveryID uuid.UUID
	Channel    string
	Target     string
	Template   map[string]any
	Enabled    bool
}

// ToggleDeliveryRequest toggles delivery state.
type ToggleDeliveryRequest struct {
	UserID     uuid.UUID
	DeliveryID uuid.UUID
	Enabled    bool
}

// BulkRunRequest queues regeneration for multiple reports.
type BulkRunRequest struct {
	UserID        uuid.UUID
	DefinitionIDs []uuid.UUID
	Metadata      map[string]any
}

// ListRunsRequest filters run history.
type ListRunsRequest struct {
	UserID   uuid.UUID
	ReportID uuid.UUID
	Status   string
	Limit    int
	Offset   int
}

// Summary aggregates insights for a user.
type Summary struct {
	UserID            uuid.UUID `json:"user_id"`
	GeneratedAt       time.Time `json:"generated_at"`
	IdeaCount         int       `json:"idea_count"`
	ProjectCount      int       `json:"project_count"`
	ConversationCount int       `json:"conversation_count"`
	IncomeTotal       float64   `json:"income_total"`
	ExpenseTotal      float64   `json:"expense_total"`
	SavingsAllocation float64   `json:"savings_allocation"`
}

// DispatchOptions controls notification channels.
type DispatchOptions struct {
	SendEmail    bool
	EmailAddress string
	SendWhatsApp bool
	PhoneNumber  string
	AgentAlias   string
}

// GenerateSummary compiles a report snapshot.
func (s *Service) GenerateSummary(ctx context.Context, userID uuid.UUID) (Summary, error) {
	var (
		ideas    []ideasdomain.Idea
		projects []projectsdomain.Project
		chats    []chatsdomain.Conversation
		finances financesdomain.Summary
		err      error
	)

	if ideas, err = s.ideasRepo.ListIdeas(ctx, userID); err != nil {
		return Summary{}, err
	}

	if projects, err = s.projectsRepo.ListProjects(ctx, userID); err != nil {
		return Summary{}, err
	}

	if s.chatsRepo != nil {
		if chats, err = s.chatsRepo.ListConversations(ctx, userID); err != nil {
			return Summary{}, err
		}
	}

	if s.financeRepo != nil {
		if finances, err = s.financeRepo.AggregateSummary(ctx, userID, time.Time{}, time.Time{}); err != nil {
			return Summary{}, err
		}
	}

	return Summary{
		UserID:            userID,
		GeneratedAt:       time.Now().UTC(),
		IdeaCount:         len(ideas),
		ProjectCount:      len(projects),
		ConversationCount: len(chats),
		IncomeTotal:       finances.IncomeTotal,
		ExpenseTotal:      finances.ExpenseTotal,
		SavingsAllocation: finances.SavingsAllocation,
	}, nil
}

// CreateDefinition stores a new report definition.
func (s *Service) CreateDefinition(ctx context.Context, req CreateDefinitionRequest) (*ReportDefinition, error) {
	def, err := NewReportDefinition(req.UserID, req.Name, req.Description, toJSONMap(req.Sections), toJSONMap(req.Filters), req.Favorite)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateDefinition(ctx, def); err != nil {
		return nil, err
	}

	return def, nil
}

// UpdateDefinition updates an existing definition.
func (s *Service) UpdateDefinition(ctx context.Context, req UpdateDefinitionRequest) (*ReportDefinition, error) {
	def, err := s.repo.GetDefinition(ctx, req.DefinitionID, req.UserID)
	if err != nil {
		return nil, err
	}

	def.Name = strings.TrimSpace(req.Name)
	def.Description = strings.TrimSpace(req.Description)
	def.Sections = toJSONMap(req.Sections)
	def.Filters = toJSONMap(req.Filters)
	def.IsFavorite = req.Favorite
	def.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateDefinition(ctx, def); err != nil {
		return nil, err
	}

	return def, nil
}

// ListDefinitions returns definitions for the user with filters applied.
func (s *Service) ListDefinitions(ctx context.Context, req ListDefinitionsRequest) ([]ReportDefinition, error) {
	filters := DefinitionFilters{
		Search:          req.Search,
		IncludeArchived: req.IncludeArchived,
		OnlyFavorites:   req.OnlyFavorites,
		Channel:         req.Channel,
		Limit:           req.Limit,
		Offset:          req.Offset,
	}

	return s.repo.ListDefinitions(ctx, req.UserID, filters)
}

// GetDefinition returns a definition with its schedules and deliveries.
func (s *Service) GetDefinition(ctx context.Context, userID, definitionID uuid.UUID) (DefinitionDetail, error) {
	def, err := s.repo.GetDefinition(ctx, definitionID, userID)
	if err != nil {
		return DefinitionDetail{}, err
	}

	schedules, err := s.repo.ListSchedules(ctx, definitionID)
	if err != nil {
		return DefinitionDetail{}, err
	}

	deliveries, err := s.repo.ListDeliveries(ctx, definitionID)
	if err != nil {
		return DefinitionDetail{}, err
	}

	return DefinitionDetail{
		Definition: *def,
		Schedules:  schedules,
		Deliveries: deliveries,
	}, nil
}

// ArchiveDefinitions archives the provided definitions.
func (s *Service) ArchiveDefinitions(ctx context.Context, req BulkDefinitionRequest) error {
	return s.repo.BulkArchiveDefinitions(ctx, req.UserID, req.DefinitionIDs)
}

// RestoreDefinitions clears archived flags for provided definitions.
func (s *Service) RestoreDefinitions(ctx context.Context, req BulkDefinitionRequest) error {
	return s.repo.BulkRestoreDefinitions(ctx, req.UserID, req.DefinitionIDs)
}

// DeleteDefinitions soft deletes definitions.
func (s *Service) DeleteDefinitions(ctx context.Context, req BulkDefinitionRequest) error {
	return s.repo.BulkDeleteDefinitions(ctx, req.UserID, req.DefinitionIDs)
}

// ToggleFavorite updates favorite state.
func (s *Service) ToggleFavorite(ctx context.Context, req ToggleFavoriteRequest) error {
	return s.repo.SetFavorite(ctx, req.UserID, req.DefinitionID, req.Favorite)
}

// CreateSchedule adds a schedule for a definition.
func (s *Service) CreateSchedule(ctx context.Context, req CreateScheduleRequest) (*ReportSchedule, error) {
	if _, err := s.repo.GetDefinition(ctx, req.ReportID, req.UserID); err != nil {
		return nil, err
	}

	schedule, err := NewReportSchedule(req.ReportID, req.Cron, req.Frequency, req.Timezone, req.NextRun, req.Enabled, toJSONMap(req.Meta))
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

// UpdateSchedule updates an existing schedule.
func (s *Service) UpdateSchedule(ctx context.Context, req UpdateScheduleRequest) (*ReportSchedule, error) {
	schedule, err := s.repo.GetSchedule(ctx, req.ScheduleID)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetDefinition(ctx, schedule.ReportID, req.UserID); err != nil {
		return nil, err
	}

	schedule.Cron = strings.TrimSpace(req.Cron)
	schedule.Frequency = strings.ToLower(strings.TrimSpace(req.Frequency))
	schedule.Timezone = strings.TrimSpace(req.Timezone)
	schedule.NextRun = req.NextRun
	schedule.Enabled = req.Enabled
	schedule.Meta = toJSONMap(req.Meta)
	schedule.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateSchedule(ctx, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

// ToggleSchedule enables or disables a schedule.
func (s *Service) ToggleSchedule(ctx context.Context, req ToggleScheduleRequest) error {
	schedule, err := s.repo.GetSchedule(ctx, req.ScheduleID)
	if err != nil {
		return err
	}
	if _, err := s.repo.GetDefinition(ctx, schedule.ReportID, req.UserID); err != nil {
		return err
	}
	schedule.Toggle(req.Enabled)
	return s.repo.UpdateSchedule(ctx, schedule)
}

// DeleteSchedule removes a schedule.
func (s *Service) DeleteSchedule(ctx context.Context, userID, scheduleID uuid.UUID) error {
	schedule, err := s.repo.GetSchedule(ctx, scheduleID)
	if err != nil {
		return err
	}
	if _, err := s.repo.GetDefinition(ctx, schedule.ReportID, userID); err != nil {
		return err
	}
	return s.repo.DeleteSchedule(ctx, scheduleID)
}

// ListSchedules returns schedules for a report.
func (s *Service) ListSchedules(ctx context.Context, userID, reportID uuid.UUID) ([]ReportSchedule, error) {
	if _, err := s.repo.GetDefinition(ctx, reportID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListSchedules(ctx, reportID)
}

// CreateDelivery stores delivery configuration.
func (s *Service) CreateDelivery(ctx context.Context, req CreateDeliveryRequest) (*ReportDelivery, error) {
	if _, err := s.repo.GetDefinition(ctx, req.ReportID, req.UserID); err != nil {
		return nil, err
	}

	delivery, err := NewReportDelivery(req.ReportID, req.Channel, req.Target, toJSONMap(req.Template), req.Enabled)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateDelivery(ctx, delivery); err != nil {
		return nil, err
	}
	return delivery, nil
}

// UpdateDelivery updates delivery configuration.
func (s *Service) UpdateDelivery(ctx context.Context, req UpdateDeliveryRequest) (*ReportDelivery, error) {
	delivery, err := s.repo.GetDelivery(ctx, req.DeliveryID)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetDefinition(ctx, delivery.ReportID, req.UserID); err != nil {
		return nil, err
	}

	delivery.Channel = strings.ToLower(strings.TrimSpace(req.Channel))
	delivery.Target = strings.TrimSpace(req.Target)
	delivery.Template = toJSONMap(req.Template)
	delivery.Template = toJSONMap(req.Template)
	delivery.Enabled = req.Enabled
	delivery.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateDelivery(ctx, delivery); err != nil {
		return nil, err
	}
	return delivery, nil
}

// ToggleDelivery toggles delivery enabled state.
func (s *Service) ToggleDelivery(ctx context.Context, req ToggleDeliveryRequest) error {
	delivery, err := s.repo.GetDelivery(ctx, req.DeliveryID)
	if err != nil {
		return err
	}
	if _, err := s.repo.GetDefinition(ctx, delivery.ReportID, req.UserID); err != nil {
		return err
	}
	delivery.Toggle(req.Enabled)
	return s.repo.UpdateDelivery(ctx, delivery)
}

// DeleteDelivery deletes a delivery.
func (s *Service) DeleteDelivery(ctx context.Context, userID, deliveryID uuid.UUID) error {
	delivery, err := s.repo.GetDelivery(ctx, deliveryID)
	if err != nil {
		return err
	}
	if _, err := s.repo.GetDefinition(ctx, delivery.ReportID, userID); err != nil {
		return err
	}
	return s.repo.DeleteDelivery(ctx, deliveryID)
}

// ListDeliveries returns deliveries for a report.
func (s *Service) ListDeliveries(ctx context.Context, userID, reportID uuid.UUID) ([]ReportDelivery, error) {
	if _, err := s.repo.GetDefinition(ctx, reportID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListDeliveries(ctx, reportID)
}

// QueueRuns creates run entries for downstream processing.
func (s *Service) QueueRuns(ctx context.Context, req BulkRunRequest) ([]ReportRun, error) {
	metadata := toJSONMap(req.Metadata)
	runs := make([]ReportRun, 0, len(req.DefinitionIDs))

	for _, id := range req.DefinitionIDs {
		if _, err := s.repo.GetDefinition(ctx, id, req.UserID); err != nil {
			return nil, err
		}
		run := NewReportRun(id, req.UserID, metadata)
		if err := s.repo.CreateRun(ctx, run); err != nil {
			return nil, err
		}
		runs = append(runs, *run)
	}
	return runs, nil
}

func toJSONMap(data map[string]any) datatypes.JSONType[map[string]any] {
	if data == nil {
		return datatypes.JSONType[map[string]any]{Data: map[string]any{}}
	}
	return datatypes.JSONType[map[string]any]{Data: data}
}

// ListRuns lists run history.
func (s *Service) ListRuns(ctx context.Context, req ListRunsRequest) ([]ReportRun, error) {
	if _, err := s.repo.GetDefinition(ctx, req.ReportID, req.UserID); err != nil {
		return nil, err
	}
	filters := RunFilters{
		Status: req.Status,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	return s.repo.ListRuns(ctx, req.ReportID, filters)
}

// DispatchSummary sends the summary through configured channels.
func (s *Service) DispatchSummary(ctx context.Context, summary Summary, opts DispatchOptions) error {
	if s.publisher == nil {
		return nil
	}

	message := formatSummary(summary, opts.AgentAlias)
	subject := formatSubject(opts.AgentAlias)
	if opts.SendEmail {
		env := notifications.ReportEnvelope{
			UserID:      summary.UserID.String(),
			Subject:     subject,
			TextMessage: message,
			Destination: opts.EmailAddress,
		}
		if err := s.publisher.PublishEmailReport(ctx, env); err != nil && s.logger != nil {
			s.logger.Error("reports: publish email failed", slog.Any("error", err))
		}
	}

	if opts.SendWhatsApp {
		env := notifications.ReportEnvelope{
			UserID:      summary.UserID.String(),
			TextMessage: message,
			Destination: opts.PhoneNumber,
		}
		if err := s.publisher.PublishWhatsAppReport(ctx, env); err != nil && s.logger != nil {
			s.logger.Error("reports: publish whatsapp failed", slog.Any("error", err))
		}
	}

	return nil
}

type agentProfile struct {
	Name    string
	Persona string
	Signoff string
	Subject string
}

var agentProfiles = map[string]agentProfile{
	"chatgpt": {
		Name:    "Atlas",
		Persona: "Here is your general Woragis status update. I applied balanced, data-driven insights.",
		Signoff: "— Atlas, your Woragis co-pilot",
		Subject: "Woragis Daily Insights",
	},
	"grok": {
		Name:    "Grok Analyst",
		Persona: "Snapshot with emphasis on recent trends and headlines.",
		Signoff: "— Grok Analyst",
		Subject: "Woragis Real-Time Briefing",
	},
	"claude": {
		Name:    "Claude Strategist",
		Persona: "Thoughtful summary to guide decisions.",
		Signoff: "— Claude Strategist",
		Subject: "Woragis Strategic Digest",
	},
	"manus": {
		Name:    "Manus",
		Persona: "Advanced strategic advisor weighing probabilities.",
		Signoff: "— Manus Strategic Intelligence",
		Subject: "Woragis Deep Strategy Report",
	},
	"cipher": {
		Name:    "Cipher",
		Persona: "Quietly sharing candid insights. Keep this between us.",
		Signoff: "— Cipher",
		Subject: "Woragis Confidential Brief",
	},
}

func formatSubject(agentAlias string) string {
	if profile, ok := agentProfiles[strings.ToLower(agentAlias)]; ok && profile.Subject != "" {
		return profile.Subject
	}
	return "Woragis Daily Insights"
}

func formatSummary(summary Summary, agentAlias string) string {
	profile, ok := agentProfiles[strings.ToLower(agentAlias)]
	if !ok {
		profile = agentProfiles["chatgpt"]
	}

	return fmt.Sprintf(
		"%s\n\nGenerated: %s\nIdeas: %d\nProjects: %d\nChats: %d\nIncome: %.2f\nExpenses: %.2f\n50/50 Savings: %.2f\n\n%s",
		profile.Persona,
		summary.GeneratedAt.Format(time.RFC822),
		summary.IdeaCount,
		summary.ProjectCount,
		summary.ConversationCount,
		summary.IncomeTotal,
		summary.ExpenseTotal,
		summary.SavingsAllocation,
		profile.Signoff,
	)
}
