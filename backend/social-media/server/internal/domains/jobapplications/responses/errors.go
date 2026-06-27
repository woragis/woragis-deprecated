package responses

import "errors"

const (
	ErrCodeInvalidPayload        = 10100
	ErrCodeInvalidResponseType   = 10101
	ErrCodeRepositoryFailure     = 10102
	ErrCodeNotFound              = 10103
)

const (
	ErrNilResponse                = "responses: response entity is nil"
	ErrEmptyResponseID            = "responses: response id cannot be empty"
	ErrEmptyJobApplicationID      = "responses: job application id cannot be empty"
	ErrResponseNotFound           = "responses: response not found"
	ErrUnsupportedResponseType    = "responses: unsupported response type"
	ErrUnableToPersist            = "responses: unable to persist data"
	ErrUnableToFetch              = "responses: unable to fetch data"
	ErrUnableToUpdate             = "responses: unable to update data"
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

