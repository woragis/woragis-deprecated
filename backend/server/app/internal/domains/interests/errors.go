package interests

import "errors"

const (
	ErrCodeInvalidPayload    = 7000
	ErrCodeInvalidTitle       = 7001
	ErrCodeRepositoryFailure = 7002
	ErrCodeNotFound          = 7003
	ErrCodeConflict          = 7004
)

const (
	ErrNilInterest              = "interests: interest entity is nil"
	ErrEmptyInterestID          = "interests: interest id cannot be empty"
	ErrEmptyInterestTitle       = "interests: interest title cannot be empty"
	ErrEmptyInterestSlug        = "interests: interest slug cannot be empty"
	ErrEmptyInterestDescription = "interests: interest description cannot be empty"
	ErrInterestNotFound         = "interests: interest not found"
	ErrInterestAlreadyExists    = "interests: interest with this title already exists"
	ErrUnableToPersist          = "interests: unable to persist data"
	ErrUnableToFetch            = "interests: unable to fetch data"
	ErrUnableToUpdate            = "interests: unable to update data"
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

