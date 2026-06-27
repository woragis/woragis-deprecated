package analytics

import "errors"

const (
	ErrCodeInvalidPayload    = 8400
	ErrCodeRepositoryFailure = 8401
	ErrCodeNotFound          = 8402
	ErrCodeConflict          = 8403
)

const (
	ErrNilPostAnalytics      = "analytics: post analytics entity is nil"
	ErrEmptyPostAnalyticsID  = "analytics: post analytics id cannot be empty"
	ErrEmptySocialPostID     = "analytics: social post id cannot be empty"
	ErrPostAnalyticsNotFound = "analytics: post analytics not found"
	ErrPostAnalyticsExists   = "analytics: post analytics already exists"
	ErrUnableToPersist       = "analytics: unable to persist data"
	ErrUnableToFetch         = "analytics: unable to fetch data"
	ErrUnableToUpdate        = "analytics: unable to update data"
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
