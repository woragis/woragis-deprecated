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
	IsRecurring bool
	IsEssential bool
}

// UpdateTransactionRequest handles partial updates.
type UpdateTransactionRequest struct {
	UserID        uuid.UUID
	TransactionID uuid.UUID
	Category      string
	Description   string
	Amount        *float64
	Currency      string
	OccurredAt    *time.Time
	Type          *TransactionType
}

// BulkRecordRequest batches transactions creation.
type BulkRecordRequest struct {
	Transactions []RecordTransactionRequest
}

// BulkCategoryUpdateRequest updates categories for multiple transactions.
type BulkCategoryUpdateRequest struct {
	UserID   uuid.UUID
	IDs      []uuid.UUID
	Category string
}

// BulkTypeUpdateRequest updates type for multiple transactions.
type BulkTypeUpdateRequest struct {
	UserID uuid.UUID
	IDs    []uuid.UUID
	Type   TransactionType
}

// BulkDeleteRequest deletes transactions in bulk.
type BulkDeleteRequest struct {
	UserID uuid.UUID
	IDs    []uuid.UUID
}

// ToggleRequest handles boolean flag toggles.
type ToggleRequest struct {
	UserID uuid.UUID
	ID     uuid.UUID
	Value  bool
}

// QueryRequest carries advanced filters.
type QueryRequest struct {
	Filter TransactionFilter
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

	tx.IsRecurring = req.IsRecurring
	tx.IsEssential = req.IsEssential

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// UpdateTransaction applies partial updates.
func (s *Service) UpdateTransaction(ctx context.Context, req UpdateTransactionRequest) (*Transaction, error) {
	tx, err := s.repo.GetTransaction(ctx, req.UserID, req.TransactionID)
	if err != nil {
		return nil, err
	}

	if req.Type != nil {
		tx.Type = *req.Type
	}

	if err := tx.UpdateMutableFields(req.Category, req.Description, req.Amount, req.Currency, req.OccurredAt); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// BulkRecord persists multiple transactions.
func (s *Service) BulkRecord(ctx context.Context, req BulkRecordRequest) ([]*Transaction, error) {
	if len(req.Transactions) == 0 {
		return []*Transaction{}, nil
	}

	txs := make([]*Transaction, 0, len(req.Transactions))
	for _, payload := range req.Transactions {
		if payload.OccurredAt.IsZero() {
			payload.OccurredAt = time.Now()
		}

		tx, err := NewTransaction(
			payload.UserID,
			payload.Type,
			payload.Category,
			payload.Description,
			payload.Amount,
			payload.Currency,
			payload.OccurredAt,
		)
		if err != nil {
			return nil, err
		}

		tx.IsRecurring = payload.IsRecurring
		tx.IsEssential = payload.IsEssential
		txs = append(txs, tx)
	}

	if err := s.repo.BulkCreateTransactions(ctx, txs); err != nil {
		return nil, err
	}

	return txs, nil
}

// BulkUpdateCategory updates categories in batch.
func (s *Service) BulkUpdateCategory(ctx context.Context, req BulkCategoryUpdateRequest) error {
	return s.repo.BulkUpdateCategory(ctx, req.UserID, req.IDs, req.Category)
}

// BulkUpdateType updates transaction type in batch.
func (s *Service) BulkUpdateType(ctx context.Context, req BulkTypeUpdateRequest) error {
	return s.repo.BulkUpdateType(ctx, req.UserID, req.IDs, req.Type)
}

// BulkDelete removes multiple transactions.
func (s *Service) BulkDelete(ctx context.Context, req BulkDeleteRequest) error {
	return s.repo.BulkDeleteTransactions(ctx, req.UserID, req.IDs)
}

// ToggleArchived sets archive flag.
func (s *Service) ToggleArchived(ctx context.Context, req ToggleRequest) error {
	return s.repo.SetArchived(ctx, req.UserID, req.ID, req.Value)
}

// ToggleRecurring sets recurring flag.
func (s *Service) ToggleRecurring(ctx context.Context, req ToggleRequest) error {
	return s.repo.SetRecurring(ctx, req.UserID, req.ID, req.Value)
}

// ToggleEssential sets essential flag.
func (s *Service) ToggleEssential(ctx context.Context, req ToggleRequest) error {
	return s.repo.SetEssential(ctx, req.UserID, req.ID, req.Value)
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

	s.logSummary(query.UserID, summary)

	return summary, nil
}

// ListTransactions maintains backwards-compatible listing.
func (s *Service) ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error) {
	return s.repo.ListTransactions(ctx, userID, from, to)
}

// QueryTransactions performs advanced filtered queries.
func (s *Service) QueryTransactions(ctx context.Context, req QueryRequest) ([]Transaction, error) {
	if req.Filter.UserID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	return s.repo.QueryTransactions(ctx, req.Filter)
}

func (s *Service) logSummary(userID uuid.UUID, summary Summary) {
	if s.logger == nil {
		return
	}

	s.logger.Debug("finances summary computed",
		slog.String("user_id", userID.String()),
		slog.Float64("income_total", summary.IncomeTotal),
		slog.Float64("expense_total", summary.ExpenseTotal),
		slog.Float64("savings_allocation", summary.SavingsAllocation),
	)
}
