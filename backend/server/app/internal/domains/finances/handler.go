package finances

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes finance HTTP endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler constructs a Handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type recordTransactionPayload struct {
	UserID      string  `json:"user_id"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	OccurredAt  string  `json:"occurred_at"`
	IsRecurring bool    `json:"is_recurring"`
	IsEssential bool    `json:"is_essential"`
}

type bulkRecordPayload struct {
	Transactions []recordTransactionPayload `json:"transactions"`
}

type updateTransactionPayload struct {
	UserID      string   `json:"user_id"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Amount      *float64 `json:"amount"`
	Currency    string   `json:"currency"`
	OccurredAt  string   `json:"occurred_at"`
	Type        string   `json:"type"`
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
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	OccurredAt  time.Time `json:"occurred_at"`
	IsRecurring bool      `json:"is_recurring"`
	IsEssential bool      `json:"is_essential"`
	IsArchived  bool      `json:"is_archived"`
	CreatedAt   time.Time `json:"created_at"`
}

type summaryResponse struct {
	IncomeTotal       float64 `json:"income_total"`
	ExpenseTotal      float64 `json:"expense_total"`
	SavingsAllocation float64 `json:"savings_allocation"`
}

// RecordTransaction handles POST /finance/transactions
func (h *Handler) RecordTransaction(c *fiber.Ctx) error {
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
func (h *Handler) BulkRecord(c *fiber.Ctx) error {
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
func (h *Handler) UpdateTransaction(c *fiber.Ctx) error {
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
func (h *Handler) BulkUpdateCategory(c *fiber.Ctx) error {
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
func (h *Handler) BulkUpdateType(c *fiber.Ctx) error {
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
func (h *Handler) BulkDelete(c *fiber.Ctx) error {
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
func (h *Handler) ToggleArchived(c *fiber.Ctx) error {
	return h.handleToggle(c, h.service.ToggleArchived)
}

// ToggleRecurring handles PATCH /finance/transactions/:id/recurring
func (h *Handler) ToggleRecurring(c *fiber.Ctx) error {
	return h.handleToggle(c, h.service.ToggleRecurring)
}

// ToggleEssential handles PATCH /finance/transactions/:id/essential
func (h *Handler) ToggleEssential(c *fiber.Ctx) error {
	return h.handleToggle(c, h.service.ToggleEssential)
}

// ListTransactions handles GET /finance/transactions
func (h *Handler) ListTransactions(c *fiber.Ctx) error {
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
func (h *Handler) Summary(c *fiber.Ctx) error {
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
	})
}

// Helpers

func (h *Handler) toRecordRequest(payload recordTransactionPayload) (RecordTransactionRequest, error) {
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

	return RecordTransactionRequest{
		UserID:      userID,
		Type:        TransactionType(payload.Type),
		Category:    payload.Category,
		Description: payload.Description,
		Amount:      payload.Amount,
		Currency:    payload.Currency,
		OccurredAt:  occurredAt,
		IsRecurring: payload.IsRecurring,
		IsEssential: payload.IsEssential,
	}, nil
}

func (h *Handler) toUpdateRequest(id uuid.UUID, payload updateTransactionPayload) (UpdateTransactionRequest, error) {
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

	return UpdateTransactionRequest{
		UserID:        userID,
		TransactionID: id,
		Category:      payload.Category,
		Description:   payload.Description,
		Amount:        payload.Amount,
		Currency:      payload.Currency,
		OccurredAt:    occurredAt,
		Type:          txType,
	}, nil
}

func (h *Handler) toBulkCategoryRequest(payload bulkCategoryPayload) (BulkCategoryUpdateRequest, error) {
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

func (h *Handler) toBulkTypeRequest(payload bulkTypePayload) (BulkTypeUpdateRequest, error) {
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

func (h *Handler) toBulkDeleteRequest(payload bulkDeletePayload) (BulkDeleteRequest, error) {
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

func (h *Handler) handleToggle(c *fiber.Ctx, fn func(context.Context, ToggleRequest) error) error {
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

func (h *Handler) toTransactionFilter(query transactionQueryPayload) (TransactionFilter, error) {
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

func toTransactionResponse(tx *Transaction) transactionResponse {
	return transactionResponse{
		ID:          tx.ID.String(),
		UserID:      tx.UserID.String(),
		Type:        string(tx.Type),
		Category:    tx.Category,
		Description: tx.Description,
		Amount:      tx.Amount,
		Currency:    tx.Currency,
		OccurredAt:  tx.OccurredAt,
		IsRecurring: tx.IsRecurring,
		IsEssential: tx.IsEssential,
		IsArchived:  tx.IsArchived,
		CreatedAt:   tx.CreatedAt,
	}
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
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
