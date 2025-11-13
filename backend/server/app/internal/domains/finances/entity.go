package finances

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// TransactionType defines allowed finance transaction categories.
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

// Transaction represents a financial movement.
type Transaction struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID       `gorm:"type:uuid;index;not null"`
	Type        TransactionType `gorm:"type:varchar(16);not null"`
	Category    string          `gorm:"size:120;not null"`
	Description string          `gorm:"size:255"`
	Amount      float64         `gorm:"not null"`
	Currency    string          `gorm:"size:8;not null"`
	OccurredAt  time.Time       `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTransaction creates a new Transaction with the supplied fields.
func NewTransaction(userID uuid.UUID, txType TransactionType, category, description string, amount float64, currency string, occurredAt time.Time) (*Transaction, error) {
	t := &Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        TransactionType(strings.ToLower(string(txType))),
		Category:    strings.TrimSpace(category),
		Description: strings.TrimSpace(description),
		Amount:      amount,
		Currency:    strings.ToUpper(strings.TrimSpace(currency)),
		OccurredAt:  occurredAt.UTC(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return t, t.Validate()
}

// Validate enforces invariants for the transaction domain entity.
func (t *Transaction) Validate() error {
	if t == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilTransaction)
	}

	if t.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTransactionID)
	}

	if t.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if t.Type != TransactionTypeIncome && t.Type != TransactionTypeExpense {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedTransactionType)
	}

	if t.Category == "" {
		return NewDomainError(ErrCodeInvalidCategory, ErrEmptyCategory)
	}

	if t.Amount <= 0 {
		return NewDomainError(ErrCodeInvalidAmount, ErrAmountMustBePositive)
	}

	if len(t.Currency) == 0 {
		return NewDomainError(ErrCodeInvalidCurrency, ErrEmptyCurrency)
	}

	if len(t.Currency) != 3 {
		return NewDomainError(ErrCodeInvalidCurrency, ErrCurrencyMustBeISO)
	}

	return nil
}
