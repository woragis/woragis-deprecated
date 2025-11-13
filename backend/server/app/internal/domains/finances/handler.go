package finances

import (
	"log/slog"
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
}

type summaryQueryPayload struct {
	UserID string `query:"user_id"`
	From   string `query:"from"`
	To     string `query:"to"`
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

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var occurredAt time.Time
	if payload.OccurredAt != "" {
		if occurredAt, err = time.Parse(time.RFC3339, payload.OccurredAt); err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	tx, err := h.service.RecordTransaction(c.Context(), RecordTransactionRequest{
		UserID:      userID,
		Type:        TransactionType(payload.Type),
		Category:    payload.Category,
		Description: payload.Description,
		Amount:      payload.Amount,
		Currency:    payload.Currency,
		OccurredAt:  occurredAt,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, transactionResponse{
		ID:          tx.ID.String(),
		UserID:      tx.UserID.String(),
		Type:        string(tx.Type),
		Category:    tx.Category,
		Description: tx.Description,
		Amount:      tx.Amount,
		Currency:    tx.Currency,
		OccurredAt:  tx.OccurredAt,
		CreatedAt:   tx.CreatedAt,
	})
}

// ListTransactions handles GET /finance/transactions
func (h *Handler) ListTransactions(c *fiber.Ctx) error {
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

	transactions, err := h.service.ListTransactions(c.Context(), userID, from, to)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]transactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		resp = append(resp, transactionResponse{
			ID:          tx.ID.String(),
			UserID:      tx.UserID.String(),
			Type:        string(tx.Type),
			Category:    tx.Category,
			Description: tx.Description,
			Amount:      tx.Amount,
			Currency:    tx.Currency,
			OccurredAt:  tx.OccurredAt,
			CreatedAt:   tx.CreatedAt,
		})
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
	case ErrCodeInvalidPayload, ErrCodeInvalidType, ErrCodeInvalidCategory, ErrCodeInvalidAmount, ErrCodeInvalidCurrency:
		return fiber.StatusBadRequest
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
