package reports

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

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
	ideasRepo    IdeaRepository
	projectsRepo ProjectRepository
	financeRepo  FinanceRepository
	chatsRepo    ChatRepository
	publisher    Publisher
	logger       *slog.Logger
}

// NewService builds a new reports service.
func NewService(
	ideasRepo IdeaRepository,
	projectsRepo ProjectRepository,
	financeRepo FinanceRepository,
	chatsRepo ChatRepository,
	publisher Publisher,
	logger *slog.Logger,
) *Service {
	return &Service{
		ideasRepo:    ideasRepo,
		projectsRepo: projectsRepo,
		financeRepo:  financeRepo,
		chatsRepo:    chatsRepo,
		publisher:    publisher,
		logger:       logger,
	}
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

// DispatchSummary sends the summary through configured channels.
func (s *Service) DispatchSummary(ctx context.Context, summary Summary, opts DispatchOptions) error {
	if s.publisher == nil {
		return nil
	}

	message := formatSummary(summary)
	if opts.SendEmail {
		env := notifications.ReportEnvelope{
			UserID:      summary.UserID.String(),
			Subject:     "Woragis Daily Insights",
			Message:     message,
			Destination: opts.EmailAddress,
		}
		if err := s.publisher.PublishEmailReport(ctx, env); err != nil && s.logger != nil {
			s.logger.Error("reports: publish email failed", slog.Any("error", err))
		}
	}

	if opts.SendWhatsApp {
		env := notifications.ReportEnvelope{
			UserID:      summary.UserID.String(),
			Message:     message,
			Destination: opts.PhoneNumber,
		}
		if err := s.publisher.PublishWhatsAppReport(ctx, env); err != nil && s.logger != nil {
			s.logger.Error("reports: publish whatsapp failed", slog.Any("error", err))
		}
	}

	return nil
}

func formatSummary(summary Summary) string {
	return fmt.Sprintf(
		"Insights (%s)\nIdeas: %d\nProjects: %d\nChats: %d\nIncome: %.2f\nExpenses: %.2f\n50/50 Savings: %.2f",
		summary.GeneratedAt.Format(time.RFC822),
		summary.IdeaCount,
		summary.ProjectCount,
		summary.ConversationCount,
		summary.IncomeTotal,
		summary.ExpenseTotal,
		summary.SavingsAllocation,
	)
}
