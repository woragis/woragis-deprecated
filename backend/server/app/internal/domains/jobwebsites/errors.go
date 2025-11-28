package jobwebsites

import "errors"

const (
	ErrCodeInvalidPayload    = 10100
	ErrCodeRepositoryFailure = 10101
	ErrCodeNotFound          = 10102
)

const (
	ErrNilWebsite          = "jobwebsites: website entity is nil"
	ErrEmptyWebsiteID      = "jobwebsites: website id cannot be empty"
	ErrEmptyWebsiteName    = "jobwebsites: website name cannot be empty"
	ErrEmptyDisplayName    = "jobwebsites: display name cannot be empty"
	ErrWebsiteNotFound     = "jobwebsites: website not found"
	ErrInvalidDailyLimit   = "jobwebsites: daily limit cannot be negative"
	ErrInvalidCurrentCount = "jobwebsites: current count cannot be negative"
	ErrUnableToPersist     = "jobwebsites: unable to persist data"
	ErrUnableToFetch       = "jobwebsites: unable to fetch data"
	ErrUnableToUpdate      = "jobwebsites: unable to update data"
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

