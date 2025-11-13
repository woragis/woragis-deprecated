package finances

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for finances.
type Repository interface {
	CreateTransaction(ctx context.Context, tx *Transaction) error
	UpdateTransaction(ctx context.Context, tx *Transaction) error
	GetTransaction(ctx context.Context, userID, id uuid.UUID) (*Transaction, error)
	ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error)
	QueryTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
	BulkCreateTransactions(ctx context.Context, txs []*Transaction) error
	BulkUpdateCategory(ctx context.Context, userID uuid.UUID, ids []uuid.UUID, category string) error
	BulkUpdateType(ctx context.Context, userID uuid.UUID, ids []uuid.UUID, txType TransactionType) error
	BulkDeleteTransactions(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	SetArchived(ctx context.Context, userID, id uuid.UUID, archived bool) error
	SetRecurring(ctx context.Context, userID, id uuid.UUID, recurring bool) error
	SetEssential(ctx context.Context, userID, id uuid.UUID, essential bool) error
	AggregateSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (Summary, error)
}

// TransactionFilter represents advanced query filters.
type TransactionFilter struct {
	UserID          uuid.UUID
	Types           []TransactionType
	Categories      []string
	MinAmount       *float64
	MaxAmount       *float64
	IncludeArchived *bool
	IsRecurring     *bool
	IsEssential     *bool
	Search          string
	From            time.Time
	To              time.Time
	Limit           int
	Offset          int
	Sort            string
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

func (r *gormRepository) UpdateTransaction(ctx context.Context, tx *Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Model(&Transaction{}).
		Where("id = ? AND user_id = ?", tx.ID, tx.UserID).
		Updates(map[string]any{
			"type":         tx.Type,
			"category":     tx.Category,
			"description":  tx.Description,
			"amount":       tx.Amount,
			"currency":     tx.Currency,
			"occurred_at":  tx.OccurredAt,
			"is_recurring": tx.IsRecurring,
			"is_essential": tx.IsEssential,
			"is_archived":  tx.IsArchived,
			"updated_at":   tx.UpdatedAt,
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	return nil
}

func (r *gormRepository) GetTransaction(ctx context.Context, userID, id uuid.UUID) (*Transaction, error) {
	var tx Transaction
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&tx).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrTransactionNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &tx, nil
}

func (r *gormRepository) ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error) {
	filter := TransactionFilter{UserID: userID, From: from, To: to}
	archived := false
	filter.IncludeArchived = &archived
	return r.QueryTransactions(ctx, filter)
}

func (r *gormRepository) QueryTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	query := r.db.WithContext(ctx).Model(&Transaction{}).Where("user_id = ?", filter.UserID)

	if filter.From != (time.Time{}) {
		query = query.Where("occurred_at >= ?", filter.From)
	}
	if filter.To != (time.Time{}) {
		query = query.Where("occurred_at <= ?", filter.To)
	}
	if len(filter.Types) > 0 {
		types := make([]string, 0, len(filter.Types))
		for _, t := range filter.Types {
			types = append(types, string(t))
		}
		query = query.Where("type IN ?", types)
	}
	if len(filter.Categories) > 0 {
		query = query.Where("category IN ?", filter.Categories)
	}
	if filter.MinAmount != nil {
		query = query.Where("amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("amount <= ?", *filter.MaxAmount)
	}
	if filter.IncludeArchived != nil && !*filter.IncludeArchived {
		query = query.Where("is_archived = ?", false)
	}
	if filter.IsRecurring != nil {
		query = query.Where("is_recurring = ?", *filter.IsRecurring)
	}
	if filter.IsEssential != nil {
		query = query.Where("is_essential = ?", *filter.IsEssential)
	}
	if filter.Search != "" {
		pattern := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(category) LIKE ? OR LOWER(description) LIKE ?", pattern, pattern)
	}

	sort := strings.TrimSpace(filter.Sort)
	if sort == "" {
		sort = "occurred_at desc"
	} else {
		sort = normalizeSort(sort)
	}
	query = query.Order(sort)

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var transactions []Transaction
	if err := query.Find(&transactions).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return transactions, nil
}

func (r *gormRepository) BulkCreateTransactions(ctx context.Context, txs []*Transaction) error {
	if len(txs) == 0 {
		return nil
	}

	for _, tx := range txs {
		if err := tx.Validate(); err != nil {
			return err
		}
	}

	if err := r.db.WithContext(ctx).Create(&txs).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToBulkPersist)
	}

	return nil
}

func (r *gormRepository) BulkUpdateCategory(ctx context.Context, userID uuid.UUID, ids []uuid.UUID, category string) error {
	if len(ids) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Updates(map[string]any{
			"category":   category,
			"updated_at": time.Now().UTC(),
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	return nil
}

func (r *gormRepository) BulkUpdateType(ctx context.Context, userID uuid.UUID, ids []uuid.UUID, txType TransactionType) error {
	if len(ids) == 0 {
		return nil
	}

	if txType != TransactionTypeIncome && txType != TransactionTypeExpense {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedTransactionType)
	}

	if err := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Updates(map[string]any{
			"type":       txType,
			"updated_at": time.Now().UTC(),
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	return nil
}

func (r *gormRepository) BulkDeleteTransactions(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND id IN ?", userID, ids).
		Delete(&Transaction{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
	}

	return nil
}

func (r *gormRepository) SetArchived(ctx context.Context, userID, id uuid.UUID, archived bool) error {
	if err := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Where("user_id = ? AND id = ?", userID, id).
		Updates(map[string]any{"is_archived": archived, "updated_at": time.Now().UTC()}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) SetRecurring(ctx context.Context, userID, id uuid.UUID, recurring bool) error {
	if err := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Where("user_id = ? AND id = ?", userID, id).
		Updates(map[string]any{"is_recurring": recurring, "updated_at": time.Now().UTC()}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) SetEssential(ctx context.Context, userID, id uuid.UUID, essential bool) error {
	if err := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Where("user_id = ? AND id = ?", userID, id).
		Updates(map[string]any{"is_essential": essential, "updated_at": time.Now().UTC()}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
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
		Where("user_id = ? AND is_archived = ?", userID, false).
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

	summary.SavingsAllocation = summary.IncomeTotal * 0.5

	return summary, nil
}

func normalizeSort(sort string) string {
	sort = strings.TrimSpace(sort)
	if sort == "" {
		return "occurred_at desc"
	}

	// support leading '-' for desc
	if strings.HasPrefix(sort, "-") {
		field := strings.TrimPrefix(sort, "-")
		return field + " desc"
	}

	return sort + " asc"
}
