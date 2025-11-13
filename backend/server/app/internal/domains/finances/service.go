package finances

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service defines finance domain operations consumed by transport layers.
type Service interface {
	RecordTransaction(ctx context.Context, req RecordTransactionRequest) (*Transaction, error)
	UpdateTransaction(ctx context.Context, req UpdateTransactionRequest) (*Transaction, error)
	BulkRecord(ctx context.Context, req BulkRecordRequest) ([]*Transaction, error)
	BulkUpdateCategory(ctx context.Context, req BulkCategoryUpdateRequest) error
	BulkUpdateType(ctx context.Context, req BulkTypeUpdateRequest) error
	BulkDelete(ctx context.Context, req BulkDeleteRequest) error
	ToggleArchived(ctx context.Context, req ToggleRequest) error
	ToggleRecurring(ctx context.Context, req ToggleRequest) error
	ToggleEssential(ctx context.Context, req ToggleRequest) error
	GetSummary(ctx context.Context, query SummaryQuery) (Summary, error)
	ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error)
	QueryTransactions(ctx context.Context, req QueryRequest) ([]Transaction, error)
	CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*RecurringTemplate, error)
	UpdateTemplate(ctx context.Context, req UpdateTemplateRequest) (*RecurringTemplate, error)
	DeleteTemplate(ctx context.Context, userID, templateID uuid.UUID) error
	ListTemplates(ctx context.Context, userID uuid.UUID) ([]RecurringTemplate, error)
	CashflowProjection(ctx context.Context, req CashflowProjectionRequest) (CashflowProjection, error)
	ImportTransactions(ctx context.Context, req ImportTransactionsRequest) ([]*Transaction, error)
}

// service orchestrates finance domain use-cases.
type service struct {
	repo   Repository
	logger *slog.Logger
}

// Ensure service implements Service.
var _ Service = (*service)(nil)

// NewService builds a finance domain service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// RecordTransactionRequest transports API payloads to the domain layer.
type RecordTransactionRequest struct {
	UserID       uuid.UUID
	Type         TransactionType
	Category     string
	Description  string
	Amount       float64
	Currency     string
	BaseCurrency string
	ExchangeRate float64
	OccurredAt   time.Time
	IsRecurring  bool
	IsEssential  bool
	TemplateID   *uuid.UUID
	Tags         []string
}

// UpdateTransactionRequest handles partial updates.
type UpdateTransactionRequest struct {
	UserID        uuid.UUID
	TransactionID uuid.UUID
	Category      string
	Description   string
	Amount        *float64
	Currency      string
	BaseCurrency  string
	ExchangeRate  *float64
	OccurredAt    *time.Time
	Type          *TransactionType
	TemplateID    *uuid.UUID
	Tags          []string
	TagsProvided  bool
}

// BulkRecordRequest batches transactions creation.
type BulkRecordRequest struct {
	Transactions []RecordTransactionRequest
}

// CreateTemplateRequest represents inputs for creating recurring templates.
type CreateTemplateRequest struct {
	UserID       uuid.UUID
	Name         string
	Type         TransactionType
	Category     string
	Description  string
	Amount       float64
	Currency     string
	BaseCurrency string
	ExchangeRate float64
	Frequency    RecurringFrequency
	Interval     int
	DayOfMonth   *int
	Weekday      *int
	Tags         []string
}

// UpdateTemplateRequest handles template mutation.
type UpdateTemplateRequest struct {
	UserID       uuid.UUID
	TemplateID   uuid.UUID
	Name         string
	Category     string
	Description  string
	Amount       *float64
	Currency     string
	BaseCurrency string
	ExchangeRate *float64
	Frequency    *RecurringFrequency
	Interval     *int
	DayOfMonth   *int
	Weekday      *int
	Tags         []string
	TagsProvided bool
}

// CashflowProjectionRequest defines cashflow calculation params.
type CashflowProjectionRequest struct {
	UserID       uuid.UUID
	PastMonths   int
	FutureMonths int
}

// CashflowMonth describes actual and projected totals for a month.
type CashflowMonth struct {
	Year             int     `json:"year"`
	Month            int     `json:"month"`
	ActualIncome     float64 `json:"actual_income"`
	ActualExpense    float64 `json:"actual_expense"`
	ProjectedIncome  float64 `json:"projected_income"`
	ProjectedExpense float64 `json:"projected_expense"`
	NetActual        float64 `json:"net_actual"`
	NetProjected     float64 `json:"net_projected"`
}

// CashflowProjection aggregates monthly cashflow results.
type CashflowProjection struct {
	BaseCurrency string          `json:"base_currency"`
	Months       []CashflowMonth `json:"months"`
}

// ImportFormat enumerates supported import types.
type ImportFormat string

const (
	ImportFormatCSV ImportFormat = "csv"
	ImportFormatOFX ImportFormat = "ofx"
)

// ImportTransactionsRequest describes a bulk import payload.
type ImportTransactionsRequest struct {
	UserID          uuid.UUID
	Format          ImportFormat
	BaseCurrency    string
	DefaultCurrency string
	ExchangeRates   map[string]float64
	Contents        []byte
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
func (s *service) RecordTransaction(ctx context.Context, req RecordTransactionRequest) (*Transaction, error) {
	if req.OccurredAt.IsZero() {
		req.OccurredAt = time.Now()
	}

	baseCurrency, normalizedAmount, err := s.normalizeAmount(req.Amount, req.Currency, req.BaseCurrency, req.ExchangeRate)
	if err != nil {
		return nil, err
	}

	tx, err := NewTransaction(
		req.UserID,
		req.Type,
		req.Category,
		req.Description,
		req.Amount,
		req.Currency,
		req.OccurredAt,
		baseCurrency,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.UpdateNormalization(baseCurrency, normalizedAmount); err != nil {
		return nil, err
	}

	tx.IsRecurring = req.IsRecurring
	tx.IsEssential = req.IsEssential
	tx.AttachTemplate(req.TemplateID)
	tx.ApplyTags(req.Tags)

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// UpdateTransaction applies partial updates.
func (s *service) UpdateTransaction(ctx context.Context, req UpdateTransactionRequest) (*Transaction, error) {
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

	baseCurrency := tx.BaseCurrency
	if req.BaseCurrency != "" {
		baseCurrency = req.BaseCurrency
	}

	var exchangeRate float64
	if req.ExchangeRate != nil {
		exchangeRate = *req.ExchangeRate
	}

	fallbackRate := 0.0
	if tx.Amount > 0 {
		fallbackRate = tx.NormalizedAmount / tx.Amount
	}

	baseCurrency, normalizedAmount, err := s.normalizeAmountWithFallback(tx.Amount, tx.Currency, baseCurrency, exchangeRate, fallbackRate)
	if err != nil {
		return nil, err
	}

	if err := tx.UpdateNormalization(baseCurrency, normalizedAmount); err != nil {
		return nil, err
	}

	if req.TemplateID != nil {
		tx.AttachTemplate(req.TemplateID)
	}

	if req.TagsProvided {
		tx.ApplyTags(req.Tags)
	}

	if err := s.repo.UpdateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// BulkRecord persists multiple transactions.
func (s *service) BulkRecord(ctx context.Context, req BulkRecordRequest) ([]*Transaction, error) {
	if len(req.Transactions) == 0 {
		return []*Transaction{}, nil
	}

	txs := make([]*Transaction, 0, len(req.Transactions))
	for _, payload := range req.Transactions {
		if payload.OccurredAt.IsZero() {
			payload.OccurredAt = time.Now()
		}

		baseCurrency, normalizedAmount, err := s.normalizeAmount(payload.Amount, payload.Currency, payload.BaseCurrency, payload.ExchangeRate)
		if err != nil {
			return nil, err
		}

		tx, err := NewTransaction(
			payload.UserID,
			payload.Type,
			payload.Category,
			payload.Description,
			payload.Amount,
			payload.Currency,
			payload.OccurredAt,
			baseCurrency,
		)
		if err != nil {
			return nil, err
		}

		if err := tx.UpdateNormalization(baseCurrency, normalizedAmount); err != nil {
			return nil, err
		}

		tx.IsRecurring = payload.IsRecurring
		tx.IsEssential = payload.IsEssential
		tx.AttachTemplate(payload.TemplateID)
		tx.ApplyTags(payload.Tags)
		txs = append(txs, tx)
	}

	if err := s.repo.BulkCreateTransactions(ctx, txs); err != nil {
		return nil, err
	}

	return txs, nil
}

// BulkUpdateCategory updates categories in batch.
func (s *service) BulkUpdateCategory(ctx context.Context, req BulkCategoryUpdateRequest) error {
	return s.repo.BulkUpdateCategory(ctx, req.UserID, req.IDs, req.Category)
}

// BulkUpdateType updates transaction type in batch.
func (s *service) BulkUpdateType(ctx context.Context, req BulkTypeUpdateRequest) error {
	return s.repo.BulkUpdateType(ctx, req.UserID, req.IDs, req.Type)
}

// BulkDelete removes multiple transactions.
func (s *service) BulkDelete(ctx context.Context, req BulkDeleteRequest) error {
	return s.repo.BulkDeleteTransactions(ctx, req.UserID, req.IDs)
}

// ToggleArchived sets archive flag.
func (s *service) ToggleArchived(ctx context.Context, req ToggleRequest) error {
	return s.repo.SetArchived(ctx, req.UserID, req.ID, req.Value)
}

// ToggleRecurring sets recurring flag.
func (s *service) ToggleRecurring(ctx context.Context, req ToggleRequest) error {
	return s.repo.SetRecurring(ctx, req.UserID, req.ID, req.Value)
}

// ToggleEssential sets essential flag.
func (s *service) ToggleEssential(ctx context.Context, req ToggleRequest) error {
	return s.repo.SetEssential(ctx, req.UserID, req.ID, req.Value)
}

// CreateTemplate persists a recurring template blueprint.
func (s *service) CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*RecurringTemplate, error) {
	baseCurrency, normalizedAmount, err := s.normalizeAmount(req.Amount, req.Currency, req.BaseCurrency, req.ExchangeRate)
	if err != nil {
		return nil, err
	}

	template, err := NewRecurringTemplate(req.UserID, req.Name, req.Type, req.Category, req.Description, req.Amount, req.Currency, baseCurrency, req.Frequency, req.Interval, req.DayOfMonth, req.Weekday)
	if err != nil {
		return nil, err
	}

	if err := template.UpdateNormalization(baseCurrency, normalizedAmount); err != nil {
		return nil, err
	}

	template.ApplyTags(req.Tags)

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

// UpdateTemplate mutates an existing recurring template.
func (s *service) UpdateTemplate(ctx context.Context, req UpdateTemplateRequest) (*RecurringTemplate, error) {
	template, err := s.repo.GetTemplate(ctx, req.UserID, req.TemplateID)
	if err != nil {
		return nil, err
	}

	if err := template.UpdateMutableFields(req.Name, req.Category, req.Description, req.Amount, req.Currency, req.Frequency, req.Interval, req.DayOfMonth, req.Weekday); err != nil {
		return nil, err
	}

	baseCurrency := template.BaseCurrency
	if req.BaseCurrency != "" {
		baseCurrency = req.BaseCurrency
	}

	exchangeRate := 0.0
	if req.ExchangeRate != nil {
		exchangeRate = *req.ExchangeRate
	}

	fallbackRate := 0.0
	if template.Amount > 0 {
		fallbackRate = template.NormalizedAmount / template.Amount
	}

	baseCurrency, normalizedAmount, err := s.normalizeAmountWithFallback(template.Amount, template.Currency, baseCurrency, exchangeRate, fallbackRate)
	if err != nil {
		return nil, err
	}

	if err := template.UpdateNormalization(baseCurrency, normalizedAmount); err != nil {
		return nil, err
	}

	if req.TagsProvided {
		template.ApplyTags(req.Tags)
	}

	if err := s.repo.UpdateTemplate(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

// DeleteTemplate removes a recurring template.
func (s *service) DeleteTemplate(ctx context.Context, userID, templateID uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, userID, templateID)
}

// ListTemplates fetches all recurring templates for a user.
func (s *service) ListTemplates(ctx context.Context, userID uuid.UUID) ([]RecurringTemplate, error) {
	return s.repo.ListTemplates(ctx, userID)
}

// SummaryQuery represents filters for summary retrieval.
type SummaryQuery struct {
	UserID uuid.UUID
	From   time.Time
	To     time.Time
}

// GetSummary returns aggregated totals honoring the 50/50 rule.
func (s *service) GetSummary(ctx context.Context, query SummaryQuery) (Summary, error) {
	summary, err := s.repo.AggregateSummary(ctx, query.UserID, query.From, query.To)
	if err != nil {
		return Summary{}, err
	}

	s.logSummary(query.UserID, summary)

	return summary, nil
}

// ListTransactions maintains backwards-compatible listing.
func (s *service) ListTransactions(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]Transaction, error) {
	return s.repo.ListTransactions(ctx, userID, from, to)
}

// QueryTransactions performs advanced filtered queries.
func (s *service) QueryTransactions(ctx context.Context, req QueryRequest) ([]Transaction, error) {
	if req.Filter.UserID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	return s.repo.QueryTransactions(ctx, req.Filter)
}

// CashflowProjection builds an actual vs projected cashflow timeline.
func (s *service) CashflowProjection(ctx context.Context, req CashflowProjectionRequest) (CashflowProjection, error) {
	if req.PastMonths <= 0 {
		req.PastMonths = 6
	}
	if req.FutureMonths < 0 {
		req.FutureMonths = 0
	}

	now := time.Now().UTC()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	pastStart := currentMonthStart.AddDate(0, -(req.PastMonths - 1), 0)
	pastEnd := currentMonthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	monthlyTotals, err := s.repo.MonthlyTotals(ctx, req.UserID, pastStart, pastEnd)
	if err != nil {
		return CashflowProjection{}, err
	}

	templates, err := s.repo.ListTemplates(ctx, req.UserID)
	if err != nil {
		return CashflowProjection{}, err
	}

	summary, err := s.repo.AggregateSummary(ctx, req.UserID, time.Time{}, time.Time{})
	if err != nil {
		return CashflowProjection{}, err
	}

	baseCurrency := summary.BaseCurrency
	if baseCurrency == "" && len(templates) > 0 {
		baseCurrency = templates[0].BaseCurrency
	}

	monthlyMap := make(map[string]MonthlyTotal)
	for _, mt := range monthlyTotals {
		key := monthlyKey(mt.Year, mt.Month)
		monthlyMap[key] = mt
	}

	projection := CashflowProjection{BaseCurrency: baseCurrency}

	for offset := -(req.PastMonths - 1); offset <= 0; offset++ {
		month := currentMonthStart.AddDate(0, offset, 0)
		year, monthIdx := month.Year(), int(month.Month())
		key := monthlyKey(year, monthIdx)

		data := monthlyMap[key]
		monthEntry := CashflowMonth{
			Year:             year,
			Month:            monthIdx,
			ActualIncome:     data.IncomeTotal,
			ActualExpense:    data.ExpenseTotal,
			ProjectedIncome:  data.IncomeTotal,
			ProjectedExpense: data.ExpenseTotal,
		}
		monthEntry.NetActual = monthEntry.ActualIncome - monthEntry.ActualExpense
		monthEntry.NetProjected = monthEntry.ProjectedIncome - monthEntry.ProjectedExpense
		projection.Months = append(projection.Months, monthEntry)
	}

	if req.FutureMonths > 0 {
		incomeContribution, expenseContribution := aggregateTemplateMonthlyContribution(templates)

		for offset := 1; offset <= req.FutureMonths; offset++ {
			month := currentMonthStart.AddDate(0, offset, 0)
			entry := CashflowMonth{
				Year:             month.Year(),
				Month:            int(month.Month()),
				ActualIncome:     0,
				ActualExpense:    0,
				ProjectedIncome:  incomeContribution,
				ProjectedExpense: expenseContribution,
			}
			entry.NetActual = entry.ActualIncome - entry.ActualExpense
			entry.NetProjected = entry.ProjectedIncome - entry.ProjectedExpense
			projection.Months = append(projection.Months, entry)
		}
	}

	return projection, nil
}

// ImportTransactions parses and persists transactions from external formats.
func (s *service) ImportTransactions(ctx context.Context, req ImportTransactionsRequest) ([]*Transaction, error) {
	var records []RecordTransactionRequest
	var err error

	switch req.Format {
	case ImportFormatCSV:
		records, err = parseCSVTransactions(req)
	case ImportFormatOFX:
		records, err = parseOFXTransactions(req)
	default:
		return nil, NewDomainError(ErrCodeInvalidPayload, "finances: unsupported import format")
	}
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []*Transaction{}, nil
	}

	batch := BulkRecordRequest{Transactions: records}
	return s.BulkRecord(ctx, batch)
}

func (s *service) logSummary(userID uuid.UUID, summary Summary) {
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

func (s *service) normalizeAmount(amount float64, currency, baseCurrency string, exchangeRate float64) (string, float64, error) {
	return s.normalizeAmountWithFallback(amount, currency, baseCurrency, exchangeRate, 0)
}

func (s *service) normalizeAmountWithFallback(amount float64, currency, baseCurrency string, exchangeRate float64, fallbackRate float64) (string, float64, error) {
	currency = normalizeCurrency(currency)
	baseCurrency = normalizeCurrency(baseCurrency)

	if currency == "" {
		return "", 0, NewDomainError(ErrCodeInvalidCurrency, ErrEmptyCurrency)
	}

	if baseCurrency == "" {
		baseCurrency = currency
	}

	rate := exchangeRate
	if baseCurrency == currency {
		rate = 1
	}

	if rate <= 0 {
		rate = fallbackRate
	}

	if baseCurrency != currency && rate <= 0 {
		return "", 0, NewDomainError(ErrCodeInvalidCurrency, ErrMissingExchangeRate)
	}

	normalized := amount
	if baseCurrency != currency {
		normalized = amount * rate
	}

	return baseCurrency, normalized, nil
}

func aggregateTemplateMonthlyContribution(templates []RecurringTemplate) (income float64, expense float64) {
	for _, tpl := range templates {
		amount := templateMonthlyAverage(tpl)
		if tpl.Type == TransactionTypeIncome {
			income += amount
		} else {
			expense += amount
		}
	}
	return
}

func templateMonthlyAverage(tpl RecurringTemplate) float64 {
	if tpl.Interval <= 0 {
		return tpl.NormalizedAmount
	}

	switch tpl.Frequency {
	case FrequencyWeekly:
		weeksPerMonth := 52.0 / 12.0
		return tpl.NormalizedAmount * (weeksPerMonth / float64(tpl.Interval))
	case FrequencyBiWeekly:
		paymentsPerMonth := 26.0 / 12.0
		return tpl.NormalizedAmount * (paymentsPerMonth / float64(tpl.Interval))
	case FrequencyMonthly:
		return tpl.NormalizedAmount / float64(tpl.Interval)
	case FrequencyQuarterly:
		return tpl.NormalizedAmount / float64(tpl.Interval*3)
	default:
		return tpl.NormalizedAmount
	}
}

func monthlyKey(year, month int) string {
	return fmt.Sprintf("%04d-%02d", year, month)
}

func parseCSVTransactions(req ImportTransactionsRequest) ([]RecordTransactionRequest, error) {
	reader := csv.NewReader(bytes.NewReader(req.Contents))
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, fmt.Sprintf("finances: unable to parse csv: %v", err))
	}

	if len(records) == 0 {
		return []RecordTransactionRequest{}, nil
	}

	header := map[string]int{}
	for idx, column := range records[0] {
		header[strings.ToLower(strings.TrimSpace(column))] = idx
	}

	rows := records[1:]
	result := make([]RecordTransactionRequest, 0, len(rows))

	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		get := func(key string) string {
			if idx, ok := header[key]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		amountStr := get("amount")
		if amountStr == "" {
			continue
		}
		amount, err := strconv.ParseFloat(strings.ReplaceAll(amountStr, ",", ""), 64)
		if err != nil {
			return nil, NewDomainError(ErrCodeInvalidPayload, fmt.Sprintf("finances: invalid amount %q", amountStr))
		}

		txType := TransactionType(strings.ToLower(get("type")))
		if txType != TransactionTypeIncome && txType != TransactionTypeExpense {
			if amount >= 0 {
				txType = TransactionTypeIncome
			} else {
				txType = TransactionTypeExpense
			}
		}

		if amount < 0 {
			amount = -amount
		}

		currency := get("currency")
		if currency == "" {
			currency = req.DefaultCurrency
		}

		baseCurrency := get("base_currency")
		if baseCurrency == "" {
			baseCurrency = req.BaseCurrency
		}

		exchangeRate := 0.0
		if rateStr := get("exchange_rate"); rateStr != "" {
			if rate, err := strconv.ParseFloat(rateStr, 64); err == nil {
				exchangeRate = rate
			}
		}
		if exchangeRate == 0 {
			exchangeRate = lookupExchangeRate(req.ExchangeRates, currency)
		}

		occurredAt, _ := parseFlexibleDate(get("occurred_at"))

		templateID, err := parseTemplateID(get("template_id"))
		if err != nil {
			return nil, err
		}

		tags := splitTags(get("tags"))

		isRecurring := parseBool(get("is_recurring"))
		isEssential := parseBool(get("is_essential"))

		result = append(result, RecordTransactionRequest{
			UserID:       req.UserID,
			Type:         txType,
			Category:     get("category"),
			Description:  get("description"),
			Amount:       amount,
			Currency:     currency,
			BaseCurrency: baseCurrency,
			ExchangeRate: exchangeRate,
			OccurredAt:   occurredAt,
			IsRecurring:  isRecurring,
			IsEssential:  isEssential,
			TemplateID:   templateID,
			Tags:         tags,
		})
	}

	return result, nil
}

func parseOFXTransactions(req ImportTransactionsRequest) ([]RecordTransactionRequest, error) {
	scanner := bufio.NewScanner(bytes.NewReader(req.Contents))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var (
		inTxn      bool
		currentReq RecordTransactionRequest
	)

	transactions := make([]RecordTransactionRequest, 0)

	resetTxn := func() {
		inTxn = true
		currentReq = RecordTransactionRequest{
			UserID:       req.UserID,
			Type:         TransactionTypeExpense,
			Category:     "import",
			Description:  "",
			Currency:     req.DefaultCurrency,
			BaseCurrency: req.BaseCurrency,
			ExchangeRate: 0,
			IsRecurring:  false,
			IsEssential:  false,
			OccurredAt:   time.Time{},
			TemplateID:   nil,
			Tags:         nil,
		}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		upper := strings.ToUpper(line)
		if upper == "<STMTTRN>" {
			resetTxn()
			continue
		}

		if upper == "</STMTTRN>" {
			if !inTxn {
				continue
			}

			if currentReq.Currency == "" {
				currentReq.Currency = req.BaseCurrency
			}
			if currentReq.BaseCurrency == "" {
				currentReq.BaseCurrency = currentReq.Currency
			}

			if currentReq.ExchangeRate == 0 {
				currentReq.ExchangeRate = lookupExchangeRate(req.ExchangeRates, currentReq.Currency)
			}

			if currentReq.Amount <= 0 {
				inTxn = false
				continue
			}

			if currentReq.OccurredAt.IsZero() {
				currentReq.OccurredAt = time.Now().UTC()
			}

			transactions = append(transactions, currentReq)
			inTxn = false
			continue
		}

		if !inTxn {
			continue
		}

		switch {
		case strings.HasPrefix(upper, "<TRNTYPE>"):
			typeValue := strings.TrimSpace(line[len("<TRNTYPE>"):])
			switch strings.ToUpper(typeValue) {
			case "CREDIT":
				currentReq.Type = TransactionTypeIncome
			case "DEBIT", "PAYMENT":
				currentReq.Type = TransactionTypeExpense
			}
		case strings.HasPrefix(upper, "<TRNAMT>"):
			amountStr := strings.TrimSpace(line[len("<TRNAMT>"):])
			if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
				if amount < 0 {
					amount = -amount
					currentReq.Type = TransactionTypeExpense
				} else if currentReq.Type != TransactionTypeExpense {
					currentReq.Type = TransactionTypeIncome
				}
				currentReq.Amount = amount
			}
		case strings.HasPrefix(upper, "<DTPOSTED>"):
			dateValue := strings.TrimSpace(line[len("<DTPOSTED>"):])
			if parsed, err := parseOFXDate(dateValue); err == nil {
				currentReq.OccurredAt = parsed
			}
		case strings.HasPrefix(upper, "<NAME>"):
			currentReq.Category = strings.TrimSpace(line[len("<NAME>"):])
		case strings.HasPrefix(upper, "<MEMO>"):
			currentReq.Description = strings.TrimSpace(line[len("<MEMO>"):])
		case strings.HasPrefix(upper, "<CURRENCY>"):
			currentReq.Currency = normalizeCurrency(strings.TrimSpace(line[len("<CURRENCY>"):]))
		case strings.HasPrefix(upper, "<EXCHRATE>"):
			if rate, err := strconv.ParseFloat(strings.TrimSpace(line[len("<EXCHRATE>"):]), 64); err == nil {
				currentReq.ExchangeRate = rate
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, fmt.Sprintf("finances: ofx parse error: %v", err))
	}

	return transactions, nil
}

func lookupExchangeRate(rates map[string]float64, currency string) float64 {
	if rates == nil {
		return 0
	}
	if rate, ok := rates[normalizeCurrency(currency)]; ok {
		return rate
	}
	return 0
}

func parseFlexibleDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}

	layouts := []string{time.RFC3339, "2006-01-02", "02/01/2006"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %s", value)
}

func parseOFXDate(value string) (time.Time, error) {
	if len(value) < 8 {
		return time.Time{}, fmt.Errorf("invalid ofx date")
	}
	datePart := value[:8]
	parsed, err := time.Parse("20060102", datePart)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

func splitTags(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';'
	})
	if len(parts) == 0 {
		return nil
	}
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}

func parseBool(value string) bool {
	if value == "" {
		return false
	}
	boolean, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return boolean
}
