package finances

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates finance domain use-cases.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService builds a finance domain service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// RecordTransactionRequest transports API payloads to the domain layer.
type RecordTransactionRequest struct {
	UserID      uuid.UUID
	Type        TransactionType
	Category    string
	Description string
	Amount      float64
	Currency    string
	OccurredAt  time.Time
}

// RecordTransaction registers a new transaction and returns it.
func (s *Service) RecordTransaction(ctx context.Context, req RecordTransactionRequest) (*Transaction, error) {
	if req.OccurredAt.IsZero() {
		req.OccurredAt = time.Now()
	}

	tx, err := NewTransaction(
		req.UserID,
		req.Type,
		req.Category,
		req.Description,
		req.Amount,
		req.Currency,
		req.OccurredAt,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// SummaryQuery represents filters for summary retrieval.
type SummaryQuery struct {
	UserID uuid.UUID
	From   time.Time
	To     time.Time
}

// GetSummary returns aggregated totals honoring the 50/50 rule.
func (s *Service) GetSummary(ctx context.Context, query SummaryQuery) (Summary, error) {
	summary, err := s.repo.AggregateSummary(ctx, query.UserID, query.From, query.To)
	if err != nil {
		return Summary{}, err
	}

	if s.logger != nil {
		s.logger.Debug("finances summary computed",
			slog.String("user_id", query.UserID.String()),
			slog.Float64("income_total", summary.IncomeTotal),
			slog.Float64("expense_total", summary.ExpenseTotal),
			slog.Float64("savings_allocation", summary.SavingsAllocation),
		)
	}

	return summary, nil
}

// ListTransactions returns the set of transactions for a user within the window.
func (s *Service) ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error) {
	return s.repo.ListTransactions(ctx, userID, from, to)
}
