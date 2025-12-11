package userpreferences

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserPreferences represents a user's default preferences for job applications.
type UserPreferences struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"column:user_id;type:uuid;uniqueIndex;not null" json:"userId"`
	DefaultLanguage string    `gorm:"column:default_language;size:2;default:'en'" json:"defaultLanguage"` // ISO 639-1
	DefaultCurrency string    `gorm:"column:default_currency;size:3;default:'USD'" json:"defaultCurrency"` // ISO 4217
	DefaultWebsite  string    `gorm:"column:default_website;size:50" json:"defaultWebsite"` // e.g., "linkedin", "glassdoor"
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for UserPreferences.
func (UserPreferences) TableName() string {
	return "user_preferences"
}

// NewUserPreferences creates a new user preferences entity with defaults.
func NewUserPreferences(userID uuid.UUID) (*UserPreferences, error) {
	prefs := &UserPreferences{
		ID:              uuid.New(),
		UserID:          userID,
		DefaultLanguage: "en",
		DefaultCurrency: "USD",
		DefaultWebsite:  "",
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	return prefs, prefs.Validate()
}

// Validate ensures user preferences invariants hold.
func (p *UserPreferences) Validate() error {
	if p == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilPreferences)
	}

	if p.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPreferencesID)
	}

	if p.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	// Validate language code (ISO 639-1: 2 characters)
	if p.DefaultLanguage != "" {
		p.DefaultLanguage = strings.ToLower(strings.TrimSpace(p.DefaultLanguage))
		if len(p.DefaultLanguage) != 2 {
			return NewDomainError(ErrCodeInvalidLanguage, ErrInvalidLanguageCode)
		}
		if !isValidLanguageCode(p.DefaultLanguage) {
			return NewDomainError(ErrCodeInvalidLanguage, ErrInvalidLanguageCode)
		}
	}

	// Validate currency code (ISO 4217: 3 characters, uppercase)
	if p.DefaultCurrency != "" {
		p.DefaultCurrency = strings.ToUpper(strings.TrimSpace(p.DefaultCurrency))
		if len(p.DefaultCurrency) != 3 {
			return NewDomainError(ErrCodeInvalidCurrency, ErrInvalidCurrencyCode)
		}
		if !isValidCurrencyCode(p.DefaultCurrency) {
			return NewDomainError(ErrCodeInvalidCurrency, ErrInvalidCurrencyCode)
		}
	}

	return nil
}

// UpdateLanguage updates the default language.
func (p *UserPreferences) UpdateLanguage(language string) error {
	language = strings.ToLower(strings.TrimSpace(language))
	if language != "" && (len(language) != 2 || !isValidLanguageCode(language)) {
		return NewDomainError(ErrCodeInvalidLanguage, ErrInvalidLanguageCode)
	}
	p.DefaultLanguage = language
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateCurrency updates the default currency.
func (p *UserPreferences) UpdateCurrency(currency string) error {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency != "" && (len(currency) != 3 || !isValidCurrencyCode(currency)) {
		return NewDomainError(ErrCodeInvalidCurrency, ErrInvalidCurrencyCode)
	}
	p.DefaultCurrency = currency
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateWebsite updates the default website.
func (p *UserPreferences) UpdateWebsite(website string) error {
	website = strings.ToLower(strings.TrimSpace(website))
	if website != "" && len(website) > 50 {
		return NewDomainError(ErrCodeInvalidPayload, "website name too long")
	}
	p.DefaultWebsite = website
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// isValidLanguageCode validates ISO 639-1 language codes.
// Common codes: en, pt, es, fr, de, it, ja, zh, ko, ru, ar, hi, etc.
func isValidLanguageCode(code string) bool {
	validCodes := map[string]bool{
		"en": true, "pt": true, "es": true, "fr": true, "de": true,
		"it": true, "ja": true, "zh": true, "ko": true, "ru": true,
		"ar": true, "hi": true, "nl": true, "pl": true, "tr": true,
		"vi": true, "th": true, "cs": true, "sv": true, "da": true,
		"fi": true, "no": true, "he": true, "id": true, "uk": true,
		"ro": true, "hu": true, "el": true, "bg": true, "hr": true,
		"sk": true, "sl": true, "et": true, "lv": true, "lt": true,
		"mt": true, "ga": true, "cy": true,
	}
	return validCodes[code]
}

// isValidCurrencyCode validates ISO 4217 currency codes.
// Common codes: USD, EUR, GBP, JPY, CNY, BRL, INR, CAD, AUD, etc.
func isValidCurrencyCode(code string) bool {
	validCodes := map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true, "CNY": true,
		"BRL": true, "INR": true, "CAD": true, "AUD": true, "CHF": true,
		"MXN": true, "RUB": true, "KRW": true, "SGD": true, "HKD": true,
		"NZD": true, "SEK": true, "NOK": true, "DKK": true, "PLN": true,
		"TRY": true, "ZAR": true, "AED": true, "SAR": true, "THB": true,
		"MYR": true, "IDR": true, "PHP": true, "VND": true, "ILS": true,
		"CLP": true, "ARS": true, "COP": true, "PEN": true, "UAH": true,
		"CZK": true, "HUF": true, "RON": true, "BGN": true, "HRK": true,
		"RSD": true, "ISK": true,
	}
	return validCodes[code]
}

