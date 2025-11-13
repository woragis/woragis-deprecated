package finances

import (
	"database/sql/driver"
	"encoding/json"
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
	ID               uuid.UUID       `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID       `gorm:"type:uuid;index;not null"`
	Type             TransactionType `gorm:"type:varchar(16);not null"`
	Category         string          `gorm:"size:120;not null"`
	Description      string          `gorm:"size:255"`
	Amount           float64         `gorm:"not null"`
	Currency         string          `gorm:"size:8;not null"`
	BaseCurrency     string          `gorm:"size:8;not null"`
	NormalizedAmount float64         `gorm:"not null"`
	OccurredAt       time.Time       `gorm:"not null"`
	IsRecurring      bool            `gorm:"not null;default:false"`
	IsEssential      bool            `gorm:"not null;default:false"`
	IsArchived       bool            `gorm:"not null;default:false;index"`
	TemplateID       *uuid.UUID      `gorm:"type:uuid;index"`
	Tags             TagList         `gorm:"type:jsonb"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TagList represents a normalized list of tags persisted as JSON.
type TagList []string

// Value implements the driver.Valuer interface for GORM.
func (tl TagList) Value() (driver.Value, error) {
	normalized := normalizeTags([]string(tl))
	data, err := json.Marshal([]string(normalized))
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

// Scan implements the sql.Scanner interface for GORM.
func (tl *TagList) Scan(value any) error {
	if value == nil {
		*tl = TagList{}
		return nil
	}

	switch raw := value.(type) {
	case []byte:
		return tl.fromJSON(raw)
	case string:
		return tl.fromJSON([]byte(raw))
	default:
		return NewDomainError(ErrCodeInvalidPayload, ErrUnsupportedTagEncoding)
	}
}

func (tl *TagList) fromJSON(raw []byte) error {
	if len(raw) == 0 {
		*tl = TagList{}
		return nil
	}

	var tags []string
	if err := json.Unmarshal(raw, &tags); err != nil {
		return err
	}
	*tl = normalizeTags(tags)
	return nil
}

// AsSlice returns the tag list as a copy slice.
func (tl TagList) AsSlice() []string {
	out := make([]string, len(tl))
	copy(out, tl)
	return out
}

func normalizeTags(tags []string) TagList {
	if len(tags) == 0 {
		return TagList{}
	}

	seen := make(map[string]struct{})
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag == "" {
			continue
		}
		if _, exists := seen[tag]; exists {
			continue
		}
		seen[tag] = struct{}{}
		normalized = append(normalized, tag)
	}

	return TagList(normalized)
}

func normalizeCurrency(currency string) string {
	return strings.ToUpper(strings.TrimSpace(currency))
}

// NewTransaction creates a new Transaction with the supplied fields.
func NewTransaction(userID uuid.UUID, txType TransactionType, category, description string, amount float64, currency string, occurredAt time.Time, baseCurrency string) (*Transaction, error) {
	currency = normalizeCurrency(currency)
	baseCurrency = normalizeCurrency(baseCurrency)
	if baseCurrency == "" {
		baseCurrency = currency
	}

	t := &Transaction{
		ID:               uuid.New(),
		UserID:           userID,
		Type:             TransactionType(strings.ToLower(string(txType))),
		Category:         strings.TrimSpace(category),
		Description:      strings.TrimSpace(description),
		Amount:           amount,
		Currency:         currency,
		BaseCurrency:     baseCurrency,
		NormalizedAmount: amount,
		OccurredAt:       occurredAt.UTC(),
		IsRecurring:      false,
		IsEssential:      false,
		IsArchived:       false,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
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

	if t.BaseCurrency == "" {
		return NewDomainError(ErrCodeInvalidCurrency, ErrEmptyCurrency)
	}

	if len(t.BaseCurrency) != 3 {
		return NewDomainError(ErrCodeInvalidCurrency, ErrCurrencyMustBeISO)
	}

	if t.NormalizedAmount <= 0 {
		return NewDomainError(ErrCodeInvalidAmount, ErrAmountMustBePositive)
	}

	return nil
}

// UpdateMutableFields updates mutable attributes and preserves invariants.
func (t *Transaction) UpdateMutableFields(category, description string, amount *float64, currency string, occurredAt *time.Time) error {
	if category != "" {
		t.Category = strings.TrimSpace(category)
	}
	if description != "" {
		t.Description = strings.TrimSpace(description)
	}
	if amount != nil {
		t.Amount = *amount
	}
	if currency != "" {
		t.Currency = normalizeCurrency(currency)
	}
	if occurredAt != nil && !occurredAt.IsZero() {
		t.OccurredAt = occurredAt.UTC()
	}
	t.UpdatedAt = time.Now().UTC()
	return t.Validate()
}

// UpdateNormalization sets the normalization data for the transaction.
func (t *Transaction) UpdateNormalization(baseCurrency string, normalizedAmount float64) error {
	if baseCurrency != "" {
		t.BaseCurrency = normalizeCurrency(baseCurrency)
	}
	if normalizedAmount > 0 {
		t.NormalizedAmount = normalizedAmount
	}
	t.UpdatedAt = time.Now().UTC()
	return t.Validate()
}

// AttachTemplate links the transaction to a recurring template.
func (t *Transaction) AttachTemplate(templateID *uuid.UUID) {
	if templateID != nil {
		copyID := *templateID
		t.TemplateID = &copyID
	} else {
		t.TemplateID = nil
	}
}

// ApplyTags normalizes and sets the tags slice.
func (t *Transaction) ApplyTags(tags []string) {
	t.Tags = normalizeTags(tags)
	if t.Tags == nil {
		t.Tags = TagList{}
	}
	t.UpdatedAt = time.Now().UTC()
}

// ToggleArchived sets archived flag.
func (t *Transaction) ToggleArchived(archived bool) {
	t.IsArchived = archived
	t.UpdatedAt = time.Now().UTC()
}

// ToggleRecurring sets recurring flag.
func (t *Transaction) ToggleRecurring(recurring bool) {
	t.IsRecurring = recurring
	t.UpdatedAt = time.Now().UTC()
}

// ToggleEssential sets essential flag.
func (t *Transaction) ToggleEssential(essential bool) {
	t.IsEssential = essential
	t.UpdatedAt = time.Now().UTC()
}

// ContainsAll reports whether the tag list contains all provided tags after normalization.
func (tl TagList) ContainsAll(tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	if len(tl) == 0 {
		return false
	}

	normalized := normalizeTags(tags)
	set := make(map[string]struct{}, len(tl))
	for _, tag := range tl {
		set[tag] = struct{}{}
	}

	for _, tag := range normalized {
		if _, ok := set[tag]; !ok {
			return false
		}
	}

	return true
}

// RecurringFrequency indicates how often a template should generate transactions.
type RecurringFrequency string

const (
	FrequencyWeekly    RecurringFrequency = "weekly"
	FrequencyBiWeekly  RecurringFrequency = "biweekly"
	FrequencyMonthly   RecurringFrequency = "monthly"
	FrequencyQuarterly RecurringFrequency = "quarterly"
)

// RecurringTemplate represents a reusable blueprint for scheduled transactions.
type RecurringTemplate struct {
	ID               uuid.UUID          `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID          `gorm:"type:uuid;index;not null"`
	Name             string             `gorm:"size:120;not null"`
	Type             TransactionType    `gorm:"type:varchar(16);not null"`
	Category         string             `gorm:"size:120;not null"`
	Description      string             `gorm:"size:255"`
	Amount           float64            `gorm:"not null"`
	Currency         string             `gorm:"size:8;not null"`
	BaseCurrency     string             `gorm:"size:8;not null"`
	NormalizedAmount float64            `gorm:"not null"`
	Frequency        RecurringFrequency `gorm:"size:32;not null"`
	Interval         int                `gorm:"not null;default:1"`
	DayOfMonth       *int               `gorm:"type:int"`
	Weekday          *int               `gorm:"type:int"`
	Tags             TagList            `gorm:"type:jsonb"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewRecurringTemplate builds a recurring template and validates invariants.
func NewRecurringTemplate(userID uuid.UUID, name string, txType TransactionType, category, description string, amount float64, currency, baseCurrency string, frequency RecurringFrequency, interval int, dayOfMonth *int, weekday *int) (*RecurringTemplate, error) {
	template := &RecurringTemplate{
		ID:               uuid.New(),
		UserID:           userID,
		Name:             strings.TrimSpace(name),
		Type:             TransactionType(strings.ToLower(string(txType))),
		Category:         strings.TrimSpace(category),
		Description:      strings.TrimSpace(description),
		Amount:           amount,
		Currency:         normalizeCurrency(currency),
		BaseCurrency:     normalizeCurrency(baseCurrency),
		NormalizedAmount: amount,
		Frequency:        frequency,
		Interval:         interval,
		DayOfMonth:       copyOptionalInt(dayOfMonth),
		Weekday:          copyOptionalInt(weekday),
		Tags:             TagList{},
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	return template, template.Validate()
}

// Validate ensures template invariants.
func (rt *RecurringTemplate) Validate() error {
	if rt == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilTransaction)
	}

	if rt.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTransactionID)
	}

	if rt.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if rt.Name == "" {
		return NewDomainError(ErrCodeInvalidPayload, "finances: template name cannot be empty")
	}

	if rt.Type != TransactionTypeIncome && rt.Type != TransactionTypeExpense {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedTransactionType)
	}

	if rt.Category == "" {
		return NewDomainError(ErrCodeInvalidCategory, ErrEmptyCategory)
	}

	if rt.Amount <= 0 {
		return NewDomainError(ErrCodeInvalidAmount, ErrAmountMustBePositive)
	}

	if rt.Currency == "" || len(rt.Currency) != 3 {
		return NewDomainError(ErrCodeInvalidCurrency, ErrCurrencyMustBeISO)
	}

	if rt.BaseCurrency == "" {
		rt.BaseCurrency = rt.Currency
	}

	if len(rt.BaseCurrency) != 3 {
		return NewDomainError(ErrCodeInvalidCurrency, ErrCurrencyMustBeISO)
	}

	if rt.NormalizedAmount <= 0 {
		return NewDomainError(ErrCodeInvalidAmount, ErrAmountMustBePositive)
	}

	switch rt.Frequency {
	case FrequencyWeekly, FrequencyBiWeekly, FrequencyMonthly, FrequencyQuarterly:
	default:
		return NewDomainError(ErrCodeInvalidPayload, "finances: unsupported recurring frequency")
	}

	if rt.Interval <= 0 {
		return NewDomainError(ErrCodeInvalidPayload, "finances: interval must be positive")
	}

	if rt.DayOfMonth != nil {
		if *rt.DayOfMonth < 1 || *rt.DayOfMonth > 31 {
			return NewDomainError(ErrCodeInvalidPayload, "finances: day_of_month must be between 1 and 31")
		}
	}

	if rt.Weekday != nil {
		if *rt.Weekday < 0 || *rt.Weekday > 6 {
			return NewDomainError(ErrCodeInvalidPayload, "finances: weekday must be between 0 (Sunday) and 6 (Saturday)")
		}
	}

	return nil
}

// UpdateNormalization updates template normalization data.
func (rt *RecurringTemplate) UpdateNormalization(baseCurrency string, normalizedAmount float64) error {
	if baseCurrency != "" {
		rt.BaseCurrency = normalizeCurrency(baseCurrency)
	}
	if normalizedAmount > 0 {
		rt.NormalizedAmount = normalizedAmount
	}
	rt.UpdatedAt = time.Now().UTC()
	return rt.Validate()
}

// ApplyTags assigns tags to the template.
func (rt *RecurringTemplate) ApplyTags(tags []string) {
	rt.Tags = normalizeTags(tags)
	rt.UpdatedAt = time.Now().UTC()
}

// UpdateMutableFields updates template mutable fields.
func (rt *RecurringTemplate) UpdateMutableFields(name, category, description string, amount *float64, currency string, frequency *RecurringFrequency, interval *int, dayOfMonth *int, weekday *int) error {
	if name != "" {
		rt.Name = strings.TrimSpace(name)
	}
	if category != "" {
		rt.Category = strings.TrimSpace(category)
	}
	if description != "" {
		rt.Description = strings.TrimSpace(description)
	}
	if amount != nil {
		rt.Amount = *amount
	}
	if currency != "" {
		rt.Currency = normalizeCurrency(currency)
	}
	if frequency != nil {
		rt.Frequency = *frequency
	}
	if interval != nil {
		rt.Interval = *interval
	}
	if dayOfMonth != nil {
		rt.DayOfMonth = copyOptionalInt(dayOfMonth)
	}
	if weekday != nil {
		rt.Weekday = copyOptionalInt(weekday)
	}
	rt.UpdatedAt = time.Now().UTC()
	return rt.Validate()
}

func copyOptionalInt(value *int) *int {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
