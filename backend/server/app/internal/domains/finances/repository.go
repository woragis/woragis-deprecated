package finances

import (
	"context"
	"fmt"
	"sort"
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
	CreateTemplate(ctx context.Context, template *RecurringTemplate) error
	UpdateTemplate(ctx context.Context, template *RecurringTemplate) error
	DeleteTemplate(ctx context.Context, userID, templateID uuid.UUID) error
	GetTemplate(ctx context.Context, userID, templateID uuid.UUID) (*RecurringTemplate, error)
	ListTemplates(ctx context.Context, userID uuid.UUID) ([]RecurringTemplate, error)
	MonthlyTotals(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]MonthlyTotal, error)
}

// TransactionFilter represents advanced query filters.
type TransactionFilter struct {
	UserID          uuid.UUID
	Types           []TransactionType
	Categories      []string
	Tags            []string
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
	BaseCurrency      string
}

// MonthlyTotal groups normalized income and expenses per calendar month.
type MonthlyTotal struct {
	Year         int
	Month        int
	IncomeTotal  float64
	ExpenseTotal float64
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
			"type":              tx.Type,
			"category":          tx.Category,
			"description":       tx.Description,
			"amount":            tx.Amount,
			"currency":          tx.Currency,
			"base_currency":     tx.BaseCurrency,
			"normalized_amount": tx.NormalizedAmount,
			"occurred_at":       tx.OccurredAt,
			"is_recurring":      tx.IsRecurring,
			"is_essential":      tx.IsEssential,
			"is_archived":       tx.IsArchived,
			"template_id":       tx.TemplateID,
			"tags":              tx.Tags,
			"updated_at":        tx.UpdatedAt,
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

	applyPagination := len(filter.Tags) == 0
	if applyPagination && filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if applyPagination && filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var transactions []Transaction
	if err := query.Find(&transactions).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	if len(filter.Tags) > 0 {
		transactions = filterTransactionsByTags(transactions, filter.Tags)

		// Manual pagination after filtering.
		start := filter.Offset
		if start > len(transactions) {
			return []Transaction{}, nil
		}

		end := len(transactions)
		if filter.Limit > 0 && start+filter.Limit < end {
			end = start + filter.Limit
		}
		transactions = transactions[start:end]
	}

	return transactions, nil
}

func filterTransactionsByTags(transactions []Transaction, tags []string) []Transaction {
	if len(tags) == 0 {
		return transactions
	}

	normalizedTarget := normalizeTags(tags)
	target := []string(normalizedTarget)

	filtered := make([]Transaction, 0, len(transactions))
	for _, tx := range transactions {
		if tx.Tags.ContainsAll(target) {
			filtered = append(filtered, tx)
		}
	}

	return filtered
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
		Select("type, SUM(normalized_amount) as total").
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

	var baseCurrency string
	if err := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Select("base_currency").
		Where("user_id = ? AND is_archived = ?", userID, false).
		Order("occurred_at DESC").
		Limit(1).
		Scan(&baseCurrency).Error; err == nil && baseCurrency != "" {
		summary.BaseCurrency = baseCurrency
	}

	summary.SavingsAllocation = summary.IncomeTotal * 0.5

	return summary, nil
}

func (r *gormRepository) CreateTemplate(ctx context.Context, template *RecurringTemplate) error {
	if err := template.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(template).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateTemplate(ctx context.Context, template *RecurringTemplate) error {
	if err := template.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Model(&RecurringTemplate{}).
		Where("id = ? AND user_id = ?", template.ID, template.UserID).
		Updates(map[string]any{
			"name":              template.Name,
			"type":              template.Type,
			"category":          template.Category,
			"description":       template.Description,
			"amount":            template.Amount,
			"currency":          template.Currency,
			"base_currency":     template.BaseCurrency,
			"normalized_amount": template.NormalizedAmount,
			"frequency":         template.Frequency,
			"interval":          template.Interval,
			"day_of_month":      template.DayOfMonth,
			"weekday":           template.Weekday,
			"tags":              template.Tags,
			"updated_at":        template.UpdatedAt,
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	return nil
}

func (r *gormRepository) DeleteTemplate(ctx context.Context, userID, templateID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, templateID).
		Delete(&RecurringTemplate{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
	}
	return nil
}

func (r *gormRepository) GetTemplate(ctx context.Context, userID, templateID uuid.UUID) (*RecurringTemplate, error) {
	var template RecurringTemplate
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, templateID).
		First(&template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrTemplateNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &template, nil
}

func (r *gormRepository) ListTemplates(ctx context.Context, userID uuid.UUID) ([]RecurringTemplate, error) {
	var templates []RecurringTemplate
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&templates).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return templates, nil
}

func (r *gormRepository) MonthlyTotals(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]MonthlyTotal, error) {
	if from.After(to) {
		from, to = to, from
	}

	var transactions []Transaction
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND occurred_at >= ? AND occurred_at <= ? AND is_archived = ?", userID, from, to, false).
		Find(&transactions).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	monthly := make(map[string]*MonthlyTotal)
	for _, tx := range transactions {
		occurred := tx.OccurredAt.UTC()
		year, month, _ := occurred.Date()
		key := fmt.Sprintf("%04d-%02d", year, int(month))

		entry, exists := monthly[key]
		if !exists {
			entry = &MonthlyTotal{Year: year, Month: int(month)}
			monthly[key] = entry
		}

		if tx.Type == TransactionTypeIncome {
			entry.IncomeTotal += tx.NormalizedAmount
		} else if tx.Type == TransactionTypeExpense {
			entry.ExpenseTotal += tx.NormalizedAmount
		}
	}

	keys := make([]string, 0, len(monthly))
	for key := range monthly {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]MonthlyTotal, 0, len(keys))
	for _, key := range keys {
		result = append(result, *monthly[key])
	}

	return result, nil
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
