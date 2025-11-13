package reports

import "errors"

const (
	ErrCodeInvalidPayload    = 7000
	ErrCodeRepositoryFailure = 7001
	ErrCodeNotFound          = 7002
)

const (
	ErrUnableToGenerate = "reports: unable to generate summary"
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
