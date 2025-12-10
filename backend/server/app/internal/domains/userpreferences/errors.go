package userpreferences

import "errors"

const (
	ErrCodeInvalidPayload    = 20000
	ErrCodeRepositoryFailure = 20001
	ErrCodeNotFound          = 20002
	ErrCodeInvalidLanguage   = 20003
	ErrCodeInvalidCurrency   = 20004
)

const (
	ErrNilPreferences        = "userpreferences: preferences entity is nil"
	ErrEmptyPreferencesID    = "userpreferences: preferences id cannot be empty"
	ErrEmptyUserID           = "userpreferences: user id cannot be empty"
	ErrPreferencesNotFound   = "userpreferences: preferences not found"
	ErrInvalidLanguageCode   = "userpreferences: invalid language code (must be ISO 639-1, 2 characters)"
	ErrInvalidCurrencyCode   = "userpreferences: invalid currency code (must be ISO 4217, 3 characters)"
	ErrUnableToPersist       = "userpreferences: unable to persist data"
	ErrUnableToFetch         = "userpreferences: unable to fetch data"
	ErrUnableToUpdate        = "userpreferences: unable to update data"
)

type DomainError struct {
	Code    int
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(code int, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

func AsDomainError(err error) (*DomainError, bool) {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr, true
	}
	return nil, false
}

