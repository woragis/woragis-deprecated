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
	ID               uuid.UUID       `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID       `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Type             TransactionType `gorm:"column:type;type:varchar(16);not null" json:"type"`
	Category         string          `gorm:"column:category;size:120;not null" json:"category"`
	Description      string          `gorm:"column:description;size:255" json:"description"`
	Amount           float64         `gorm:"column:amount;not null" json:"amount"`
	Currency         string          `gorm:"column:currency;size:8;not null" json:"currency"`
	BaseCurrency     string          `gorm:"column:base_currency;size:8;not null" json:"baseCurrency"`
	NormalizedAmount float64         `gorm:"column:normalized_amount;not null" json:"normalizedAmount"`
	OccurredAt       time.Time       `gorm:"column:occurred_at;not null" json:"occurredAt"`
	IsRecurring      bool            `gorm:"column:is_recurring;not null;default:false" json:"isRecurring"`
	IsEssential      bool            `gorm:"column:is_essential;not null;default:false" json:"isEssential"`
	IsArchived       bool            `gorm:"column:is_archived;not null;default:false;index" json:"isArchived"`
	TemplateID       *uuid.UUID      `gorm:"column:template_id;type:uuid;index" json:"templateId,omitempty"`
	Tags             TagList         `gorm:"column:tags;type:jsonb" json:"tags"`
	CreatedAt        time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time       `gorm:"column:updated_at" json:"updatedAt"`
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
	ID               uuid.UUID          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID          `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Name             string             `gorm:"column:name;size:120;not null" json:"name"`
	Type             TransactionType    `gorm:"column:type;type:varchar(16);not null" json:"type"`
	Category         string             `gorm:"column:category;size:120;not null" json:"category"`
	Description      string             `gorm:"column:description;size:255" json:"description"`
	Amount           float64            `gorm:"column:amount;not null" json:"amount"`
	Currency         string             `gorm:"column:currency;size:8;not null" json:"currency"`
	BaseCurrency     string             `gorm:"column:base_currency;size:8;not null" json:"baseCurrency"`
	NormalizedAmount float64            `gorm:"column:normalized_amount;not null" json:"normalizedAmount"`
	Frequency        RecurringFrequency `gorm:"column:frequency;size:32;not null" json:"frequency"`
	Interval         int                `gorm:"column:interval;not null;default:1" json:"interval"`
	DayOfMonth       *int               `gorm:"column:day_of_month;type:int" json:"dayOfMonth,omitempty"`
	Weekday          *int               `gorm:"column:weekday;type:int" json:"weekday,omitempty"`
	Tags             TagList            `gorm:"column:tags;type:jsonb" json:"tags"`
	CreatedAt        time.Time          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time          `gorm:"column:updated_at" json:"updatedAt"`
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
