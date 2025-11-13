package finances

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes finance HTTP endpoints.
type Handler interface {
	RecordTransaction(c *fiber.Ctx) error
	BulkRecord(c *fiber.Ctx) error
	UpdateTransaction(c *fiber.Ctx) error
	BulkUpdateCategory(c *fiber.Ctx) error
	BulkUpdateType(c *fiber.Ctx) error
	BulkDelete(c *fiber.Ctx) error
	ToggleArchived(c *fiber.Ctx) error
	ToggleRecurring(c *fiber.Ctx) error
	ToggleEssential(c *fiber.Ctx) error
	ListTransactions(c *fiber.Ctx) error
	Summary(c *fiber.Ctx) error
	CreateTemplate(c *fiber.Ctx) error
	UpdateTemplate(c *fiber.Ctx) error
	DeleteTemplate(c *fiber.Ctx) error
	ListTemplates(c *fiber.Ctx) error
	Cashflow(c *fiber.Ctx) error
	ImportTransactions(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a Handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{service: service, logger: logger}
}

type recordTransactionPayload struct {
	UserID       string   `json:"user_id"`
	Type         string   `json:"type"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Amount       float64  `json:"amount"`
	Currency     string   `json:"currency"`
	BaseCurrency string   `json:"base_currency"`
	ExchangeRate *float64 `json:"exchange_rate"`
	OccurredAt   string   `json:"occurred_at"`
	IsRecurring  bool     `json:"is_recurring"`
	IsEssential  bool     `json:"is_essential"`
	TemplateID   string   `json:"template_id"`
	Tags         []string `json:"tags"`
}

type bulkRecordPayload struct {
	Transactions []recordTransactionPayload `json:"transactions"`
}

type updateTransactionPayload struct {
	UserID       string   `json:"user_id"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Amount       *float64 `json:"amount"`
	Currency     string   `json:"currency"`
	BaseCurrency string   `json:"base_currency"`
	ExchangeRate *float64 `json:"exchange_rate"`
	OccurredAt   string   `json:"occurred_at"`
	Type         string   `json:"type"`
	TemplateID   string   `json:"template_id"`
	Tags         []string `json:"tags"`
}

type bulkCategoryPayload struct {
	UserID         string   `json:"user_id"`
	TransactionIDs []string `json:"transaction_ids"`
	Category       string   `json:"category"`
}

type bulkTypePayload struct {
	UserID         string   `json:"user_id"`
	TransactionIDs []string `json:"transaction_ids"`
	Type           string   `json:"type"`
}

type bulkDeletePayload struct {
	UserID         string   `json:"user_id"`
	TransactionIDs []string `json:"transaction_ids"`
}

type togglePayload struct {
	UserID string `json:"user_id"`
	Value  bool   `json:"value"`
}

type summaryQueryPayload struct {
	UserID string `query:"user_id"`
	From   string `query:"from"`
	To     string `query:"to"`
}

type transactionQueryPayload struct {
	UserID          string `query:"user_id"`
	Types           string `query:"types"`
	Categories      string `query:"categories"`
	Tags            string `query:"tags"`
	MinAmount       string `query:"min_amount"`
	MaxAmount       string `query:"max_amount"`
	IncludeArchived string `query:"include_archived"`
	IsRecurring     string `query:"is_recurring"`
	IsEssential     string `query:"is_essential"`
	Search          string `query:"search"`
	From            string `query:"from"`
	To              string `query:"to"`
	Limit           string `query:"limit"`
	Offset          string `query:"offset"`
	Sort            string `query:"sort"`
}

type transactionResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Type             string    `json:"type"`
	Category         string    `json:"category"`
	Description      string    `json:"description,omitempty"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	BaseCurrency     string    `json:"base_currency"`
	NormalizedAmount float64   `json:"normalized_amount"`
	OccurredAt       time.Time `json:"occurred_at"`
	IsRecurring      bool      `json:"is_recurring"`
	IsEssential      bool      `json:"is_essential"`
	IsArchived       bool      `json:"is_archived"`
	TemplateID       *string   `json:"template_id,omitempty"`
	Tags             []string  `json:"tags"`
	CreatedAt        time.Time `json:"created_at"`
}

type summaryResponse struct {
	IncomeTotal       float64 `json:"income_total"`
	ExpenseTotal      float64 `json:"expense_total"`
	SavingsAllocation float64 `json:"savings_allocation"`
	BaseCurrency      string  `json:"base_currency"`
}

type templateResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Category         string    `json:"category"`
	Description      string    `json:"description"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	BaseCurrency     string    `json:"base_currency"`
	NormalizedAmount float64   `json:"normalized_amount"`
	Frequency        string    `json:"frequency"`
	Interval         int       `json:"interval"`
	DayOfMonth       *int      `json:"day_of_month,omitempty"`
	Weekday          *int      `json:"weekday,omitempty"`
	Tags             []string  `json:"tags"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type createTemplatePayload struct {
	UserID       string   `json:"user_id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Amount       float64  `json:"amount"`
	Currency     string   `json:"currency"`
	BaseCurrency string   `json:"base_currency"`
	ExchangeRate *float64 `json:"exchange_rate"`
	Frequency    string   `json:"frequency"`
	Interval     int      `json:"interval"`
	DayOfMonth   *int     `json:"day_of_month"`
	Weekday      *int     `json:"weekday"`
	Tags         []string `json:"tags"`
}

type updateTemplatePayload struct {
	UserID       string   `json:"user_id"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Amount       *float64 `json:"amount"`
	Currency     string   `json:"currency"`
	BaseCurrency string   `json:"base_currency"`
	ExchangeRate *float64 `json:"exchange_rate"`
	Frequency    string   `json:"frequency"`
	Interval     *int     `json:"interval"`
	DayOfMonth   *int     `json:"day_of_month"`
	Weekday      *int     `json:"weekday"`
	Tags         []string `json:"tags"`
}

type templateQueryPayload struct {
	UserID string `query:"user_id"`
}

type cashflowQueryPayload struct {
	UserID       string `query:"user_id"`
	PastMonths   string `query:"past_months"`
	FutureMonths string `query:"future_months"`
}

type importTransactionsPayload struct {
	UserID          string             `json:"user_id"`
	Format          string             `json:"format"`
	BaseCurrency    string             `json:"base_currency"`
	DefaultCurrency string             `json:"default_currency"`
	ExchangeRates   map[string]float64 `json:"exchange_rates"`
	Contents        string             `json:"contents"`
}

// RecordTransaction handles POST /finance/transactions
func (h *handler) RecordTransaction(c *fiber.Ctx) error {
	var payload recordTransactionPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrUnableToPersist, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toRecordRequest(payload)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	tx, err := h.service.RecordTransaction(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toTransactionResponse(tx))
}

// BulkRecord handles POST /finance/transactions/bulk
func (h *handler) BulkRecord(c *fiber.Ctx) error {
	var payload bulkRecordPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := BulkRecordRequest{}
	for _, p := range payload.Transactions {
		recordReq, err := h.toRecordRequest(p)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		req.Transactions = append(req.Transactions, recordReq)
	}

	txs, err := h.service.BulkRecord(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]transactionResponse, 0, len(txs))
	for _, tx := range txs {
		resp = append(resp, toTransactionResponse(tx))
	}

	return response.Success(c, fiber.StatusCreated, resp)
}

// UpdateTransaction handles PATCH /finance/transactions/:id
func (h *handler) UpdateTransaction(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateTransactionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toUpdateRequest(id, payload)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	tx, err := h.service.UpdateTransaction(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toTransactionResponse(tx))
}

// BulkUpdateCategory handles PATCH /finance/transactions/bulk/category
func (h *handler) BulkUpdateCategory(c *fiber.Ctx) error {
	var payload bulkCategoryPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toBulkCategoryRequest(payload)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.BulkUpdateCategory(c.Context(), req); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"updated": len(req.IDs)})
}

// BulkUpdateType handles PATCH /finance/transactions/bulk/type
func (h *handler) BulkUpdateType(c *fiber.Ctx) error {
	var payload bulkTypePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toBulkTypeRequest(payload)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.BulkUpdateType(c.Context(), req); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"updated": len(req.IDs)})
}

// BulkDelete handles DELETE /finance/transactions/bulk
func (h *handler) BulkDelete(c *fiber.Ctx) error {
	var payload bulkDeletePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toBulkDeleteRequest(payload)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.BulkDelete(c.Context(), req); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"deleted": len(req.IDs)})
}

// ToggleArchived handles PATCH /finance/transactions/:id/archive
func (h *handler) ToggleArchived(c *fiber.Ctx) error {
	return h.handleToggle(c, h.service.ToggleArchived)
}

// ToggleRecurring handles PATCH /finance/transactions/:id/recurring
func (h *handler) ToggleRecurring(c *fiber.Ctx) error {
	return h.handleToggle(c, h.service.ToggleRecurring)
}

// ToggleEssential handles PATCH /finance/transactions/:id/essential
func (h *handler) ToggleEssential(c *fiber.Ctx) error {
	return h.handleToggle(c, h.service.ToggleEssential)
}

// ListTransactions handles GET /finance/transactions
func (h *handler) ListTransactions(c *fiber.Ctx) error {
	var query transactionQueryPayload
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	filter, err := h.toTransactionFilter(query)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidQuery, nil)
	}

	txs, err := h.service.QueryTransactions(c.Context(), QueryRequest{Filter: filter})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]transactionResponse, 0, len(txs))
	for _, tx := range txs {
		copy := tx
		resp = append(resp, toTransactionResponse(&copy))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// Summary handles GET /finance/summary
func (h *handler) Summary(c *fiber.Ctx) error {
	var query summaryQueryPayload
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	from, to, err := parseRange(query.From, query.To)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	summary, err := h.service.GetSummary(c.Context(), SummaryQuery{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, summaryResponse{
		IncomeTotal:       summary.IncomeTotal,
		ExpenseTotal:      summary.ExpenseTotal,
		SavingsAllocation: summary.SavingsAllocation,
		BaseCurrency:      summary.BaseCurrency,
	})
}

// CreateTemplate handles POST /finance/templates
func (h *handler) CreateTemplate(c *fiber.Ctx) error {
	var payload createTemplatePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toCreateTemplateRequest(payload)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	template, err := h.service.CreateTemplate(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toTemplateResponse(template))
}

// UpdateTemplate handles PUT /finance/templates/:id
func (h *handler) UpdateTemplate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateTemplatePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toUpdateTemplateRequest(id, payload)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	template, err := h.service.UpdateTemplate(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toTemplateResponse(template))
}

// DeleteTemplate handles DELETE /finance/templates/:id
func (h *handler) DeleteTemplate(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var query templateQueryPayload
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteTemplate(c.Context(), userID, id); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": id})
}

// ListTemplates handles GET /finance/templates
func (h *handler) ListTemplates(c *fiber.Ctx) error {
	var query templateQueryPayload
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	templates, err := h.service.ListTemplates(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]templateResponse, 0, len(templates))
	for _, tpl := range templates {
		copy := tpl
		resp = append(resp, toTemplateResponse(&copy))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// Cashflow handles GET /finance/reports/cashflow
func (h *handler) Cashflow(c *fiber.Ctx) error {
	var query cashflowQueryPayload
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toCashflowRequest(query)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	projection, err := h.service.CashflowProjection(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, projection)
}

// ImportTransactions handles POST /finance/transactions/import
func (h *handler) ImportTransactions(c *fiber.Ctx) error {
	var payload importTransactionsPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req, err := h.toImportRequest(payload)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	txs, err := h.service.ImportTransactions(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]transactionResponse, 0, len(txs))
	for _, tx := range txs {
		resp = append(resp, toTransactionResponse(tx))
	}

	return response.Success(c, fiber.StatusCreated, resp)
}

// Helpers

func (h *handler) toRecordRequest(payload recordTransactionPayload) (RecordTransactionRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return RecordTransactionRequest{}, err
	}

	var occurredAt time.Time
	if payload.OccurredAt != "" {
		if occurredAt, err = time.Parse(time.RFC3339, payload.OccurredAt); err != nil {
			return RecordTransactionRequest{}, err
		}
	}

	templateID, err := parseTemplateID(payload.TemplateID)
	if err != nil {
		return RecordTransactionRequest{}, err
	}

	exchangeRate := 0.0
	if payload.ExchangeRate != nil {
		exchangeRate = *payload.ExchangeRate
	}

	return RecordTransactionRequest{
		UserID:       userID,
		Type:         TransactionType(payload.Type),
		Category:     payload.Category,
		Description:  payload.Description,
		Amount:       payload.Amount,
		Currency:     payload.Currency,
		BaseCurrency: payload.BaseCurrency,
		ExchangeRate: exchangeRate,
		OccurredAt:   occurredAt,
		IsRecurring:  payload.IsRecurring,
		IsEssential:  payload.IsEssential,
		TemplateID:   templateID,
		Tags:         payload.Tags,
	}, nil
}

func (h *handler) toUpdateRequest(id uuid.UUID, payload updateTransactionPayload) (UpdateTransactionRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return UpdateTransactionRequest{}, err
	}

	var occurredAt *time.Time
	if payload.OccurredAt != "" {
		parsed, err := time.Parse(time.RFC3339, payload.OccurredAt)
		if err != nil {
			return UpdateTransactionRequest{}, err
		}
		occurredAt = &parsed
	}

	var txType *TransactionType
	if payload.Type != "" {
		t := TransactionType(payload.Type)
		txType = &t
	}

	templateID, err := parseTemplateID(payload.TemplateID)
	if err != nil {
		return UpdateTransactionRequest{}, err
	}

	return UpdateTransactionRequest{
		UserID:        userID,
		TransactionID: id,
		Category:      payload.Category,
		Description:   payload.Description,
		Amount:        payload.Amount,
		Currency:      payload.Currency,
		BaseCurrency:  payload.BaseCurrency,
		ExchangeRate:  payload.ExchangeRate,
		OccurredAt:    occurredAt,
		Type:          txType,
		TemplateID:    templateID,
		Tags:          payload.Tags,
		TagsProvided:  payload.Tags != nil,
	}, nil
}

func (h *handler) toBulkCategoryRequest(payload bulkCategoryPayload) (BulkCategoryUpdateRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return BulkCategoryUpdateRequest{}, err
	}

	ids, err := parseUUIDList(payload.TransactionIDs)
	if err != nil {
		return BulkCategoryUpdateRequest{}, err
	}

	return BulkCategoryUpdateRequest{UserID: userID, IDs: ids, Category: payload.Category}, nil
}

func (h *handler) toBulkTypeRequest(payload bulkTypePayload) (BulkTypeUpdateRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return BulkTypeUpdateRequest{}, err
	}

	ids, err := parseUUIDList(payload.TransactionIDs)
	if err != nil {
		return BulkTypeUpdateRequest{}, err
	}

	return BulkTypeUpdateRequest{UserID: userID, IDs: ids, Type: TransactionType(payload.Type)}, nil
}

func (h *handler) toBulkDeleteRequest(payload bulkDeletePayload) (BulkDeleteRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return BulkDeleteRequest{}, err
	}

	ids, err := parseUUIDList(payload.TransactionIDs)
	if err != nil {
		return BulkDeleteRequest{}, err
	}

	return BulkDeleteRequest{UserID: userID, IDs: ids}, nil
}

func (h *handler) handleToggle(c *fiber.Ctx, fn func(context.Context, ToggleRequest) error) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload togglePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := ToggleRequest{UserID: userID, ID: id, Value: payload.Value}

	if err := fn(c.Context(), req); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": id})
}

func (h *handler) toTransactionFilter(query transactionQueryPayload) (TransactionFilter, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return TransactionFilter{}, err
	}

	filter := TransactionFilter{UserID: userID}

	if query.From != "" {
		from, err := time.Parse(time.RFC3339, query.From)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.From = from
	}
	if query.To != "" {
		to, err := time.Parse(time.RFC3339, query.To)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.To = to
	}

	if query.Types != "" {
		for _, part := range strings.Split(query.Types, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			filter.Types = append(filter.Types, TransactionType(part))
		}
	}

	if query.Categories != "" {
		for _, part := range strings.Split(query.Categories, ",") {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				filter.Categories = append(filter.Categories, trimmed)
			}
		}
	}

	if query.Tags != "" {
		for _, part := range strings.Split(query.Tags, ",") {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				filter.Tags = append(filter.Tags, trimmed)
			}
		}
	}

	if query.MinAmount != "" {
		min, err := strconv.ParseFloat(query.MinAmount, 64)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.MinAmount = &min
	}
	if query.MaxAmount != "" {
		max, err := strconv.ParseFloat(query.MaxAmount, 64)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.MaxAmount = &max
	}

	if query.IncludeArchived != "" {
		archived, err := strconv.ParseBool(query.IncludeArchived)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.IncludeArchived = &archived
	}

	if query.IsRecurring != "" {
		recurring, err := strconv.ParseBool(query.IsRecurring)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.IsRecurring = &recurring
	}

	if query.IsEssential != "" {
		es, err := strconv.ParseBool(query.IsEssential)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.IsEssential = &es
	}

	if query.Search != "" {
		filter.Search = query.Search
	}

	if query.Limit != "" {
		limit, err := strconv.Atoi(query.Limit)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.Limit = limit
	}
	if query.Offset != "" {
		offset, err := strconv.Atoi(query.Offset)
		if err != nil {
			return TransactionFilter{}, err
		}
		filter.Offset = offset
	}

	filter.Sort = query.Sort

	return filter, nil
}

func parseRange(fromRaw, toRaw string) (time.Time, time.Time, error) {
	var (
		from time.Time
		to   time.Time
		err  error
	)

	if fromRaw != "" {
		if from, err = time.Parse(time.RFC3339, fromRaw); err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if toRaw != "" {
		if to, err = time.Parse(time.RFC3339, toRaw); err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	return from, to, nil
}

func parseUUIDList(values []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(values))
	for _, raw := range values {
		if raw == "" {
			continue
		}
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func parseTemplateID(raw string) (*uuid.UUID, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func toTransactionResponse(tx *Transaction) transactionResponse {
	var templateID *string
	if tx.TemplateID != nil {
		value := tx.TemplateID.String()
		templateID = &value
	}

	return transactionResponse{
		ID:               tx.ID.String(),
		UserID:           tx.UserID.String(),
		Type:             string(tx.Type),
		Category:         tx.Category,
		Description:      tx.Description,
		Amount:           tx.Amount,
		Currency:         tx.Currency,
		BaseCurrency:     tx.BaseCurrency,
		NormalizedAmount: tx.NormalizedAmount,
		OccurredAt:       tx.OccurredAt,
		IsRecurring:      tx.IsRecurring,
		IsEssential:      tx.IsEssential,
		IsArchived:       tx.IsArchived,
		TemplateID:       templateID,
		Tags:             tx.Tags.AsSlice(),
		CreatedAt:        tx.CreatedAt,
	}
}

func toTemplateResponse(template *RecurringTemplate) templateResponse {
	tags := template.Tags.AsSlice()
	if tags == nil {
		tags = []string{}
	}

	return templateResponse{
		ID:               template.ID.String(),
		UserID:           template.UserID.String(),
		Name:             template.Name,
		Type:             string(template.Type),
		Category:         template.Category,
		Description:      template.Description,
		Amount:           template.Amount,
		Currency:         template.Currency,
		BaseCurrency:     template.BaseCurrency,
		NormalizedAmount: template.NormalizedAmount,
		Frequency:        string(template.Frequency),
		Interval:         template.Interval,
		DayOfMonth:       template.DayOfMonth,
		Weekday:          template.Weekday,
		Tags:             tags,
		CreatedAt:        template.CreatedAt,
		UpdatedAt:        template.UpdatedAt,
	}
}

func (h *handler) toCreateTemplateRequest(payload createTemplatePayload) (CreateTemplateRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return CreateTemplateRequest{}, err
	}

	frequency, err := parseFrequency(payload.Frequency)
	if err != nil {
		return CreateTemplateRequest{}, err
	}

	exchangeRate := 0.0
	if payload.ExchangeRate != nil {
		exchangeRate = *payload.ExchangeRate
	}

	interval := payload.Interval
	if interval <= 0 {
		interval = 1
	}

	return CreateTemplateRequest{
		UserID:       userID,
		Name:         payload.Name,
		Type:         TransactionType(payload.Type),
		Category:     payload.Category,
		Description:  payload.Description,
		Amount:       payload.Amount,
		Currency:     payload.Currency,
		BaseCurrency: payload.BaseCurrency,
		ExchangeRate: exchangeRate,
		Frequency:    frequency,
		Interval:     interval,
		DayOfMonth:   payload.DayOfMonth,
		Weekday:      payload.Weekday,
		Tags:         payload.Tags,
	}, nil
}

func (h *handler) toUpdateTemplateRequest(id uuid.UUID, payload updateTemplatePayload) (UpdateTemplateRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}

	var frequency *RecurringFrequency
	if strings.TrimSpace(payload.Frequency) != "" {
		freq, err := parseFrequency(payload.Frequency)
		if err != nil {
			return UpdateTemplateRequest{}, err
		}
		frequency = &freq
	}

	return UpdateTemplateRequest{
		UserID:       userID,
		TemplateID:   id,
		Name:         payload.Name,
		Category:     payload.Category,
		Description:  payload.Description,
		Amount:       payload.Amount,
		Currency:     payload.Currency,
		BaseCurrency: payload.BaseCurrency,
		ExchangeRate: payload.ExchangeRate,
		Frequency:    frequency,
		Interval:     payload.Interval,
		DayOfMonth:   payload.DayOfMonth,
		Weekday:      payload.Weekday,
		Tags:         payload.Tags,
		TagsProvided: payload.Tags != nil,
	}, nil
}

func (h *handler) toCashflowRequest(payload cashflowQueryPayload) (CashflowProjectionRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return CashflowProjectionRequest{}, err
	}

	pastMonths := parseIntDefault(payload.PastMonths, 6)
	futureMonths := parseIntDefault(payload.FutureMonths, 6)

	return CashflowProjectionRequest{
		UserID:       userID,
		PastMonths:   pastMonths,
		FutureMonths: futureMonths,
	}, nil
}

func (h *handler) toImportRequest(payload importTransactionsPayload) (ImportTransactionsRequest, error) {
	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return ImportTransactionsRequest{}, err
	}

	format := ImportFormat(strings.ToLower(strings.TrimSpace(payload.Format)))
	if format != ImportFormatCSV && format != ImportFormatOFX {
		return ImportTransactionsRequest{}, fmt.Errorf("unsupported import format")
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(payload.Contents))
	if err != nil {
		return ImportTransactionsRequest{}, err
	}

	baseCurrency := payload.BaseCurrency
	defaultCurrency := payload.DefaultCurrency
	if defaultCurrency == "" {
		defaultCurrency = baseCurrency
	}

	return ImportTransactionsRequest{
		UserID:          userID,
		Format:          format,
		BaseCurrency:    baseCurrency,
		DefaultCurrency: defaultCurrency,
		ExchangeRates:   payload.ExchangeRates,
		Contents:        decoded,
	}, nil
}

func parseFrequency(raw string) (RecurringFrequency, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(FrequencyWeekly):
		return FrequencyWeekly, nil
	case string(FrequencyBiWeekly):
		return FrequencyBiWeekly, nil
	case string(FrequencyMonthly):
		return FrequencyMonthly, nil
	case string(FrequencyQuarterly):
		return FrequencyQuarterly, nil
	default:
		return "", fmt.Errorf("unsupported frequency")
	}
}

func parseIntDefault(value string, def int) int {
	if value == "" {
		return def
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return parsed
}

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromErrorCode(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("finances: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromErrorCode(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidType, ErrCodeInvalidCategory, ErrCodeInvalidAmount, ErrCodeInvalidCurrency, ErrCodeInvalidQuery:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	case ErrCodeRepositoryFailure:
		return fiber.StatusInternalServerError
	case ErrCodeSummaryFailure:
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
