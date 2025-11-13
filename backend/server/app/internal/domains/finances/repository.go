package finances

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for finances.
type Repository interface {
	CreateTransaction(ctx context.Context, tx *Transaction) error
	ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error)
	AggregateSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (Summary, error)
}

// Summary represents aggregated values for the finance module.
type Summary struct {
	IncomeTotal       float64
	ExpenseTotal      float64
	SavingsAllocation float64
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository instantiates a GORM-backed Repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateTransaction(ctx context.Context, tx *Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error) {
	var transactions []Transaction

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if !from.IsZero() {
		query = query.Where("occurred_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("occurred_at <= ?", to)
	}

	if err := query.Order("occurred_at desc").Find(&transactions).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return transactions, nil
}

func (r *gormRepository) AggregateSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (Summary, error) {
	type result struct {
		Type  string
		Total float64
	}

	var rows []result

	query := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Select("type, SUM(amount) as total").
		Where("user_id = ?", userID).
		Group("type")

	if !from.IsZero() {
		query = query.Where("occurred_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("occurred_at <= ?", to)
	}

	if err := query.Scan(&rows).Error; err != nil {
		return Summary{}, NewDomainError(ErrCodeSummaryFailure, ErrUnableToSummarize)
	}

	summary := Summary{}
	for _, row := range rows {
		switch TransactionType(row.Type) {
		case TransactionTypeIncome:
			summary.IncomeTotal = row.Total
		case TransactionTypeExpense:
			summary.ExpenseTotal = row.Total
		}
	}

	// Woragis 50/50 rule: allocate 50% of income to intentional spending / savings.
	summary.SavingsAllocation = summary.IncomeTotal * 0.5

	return summary, nil
}
