package interviewstages

import "errors"

const (
	ErrCodeInvalidPayload    = 10200
	ErrCodeInvalidStageType   = 10201
	ErrCodeInvalidOutcome     = 10202
	ErrCodeRepositoryFailure  = 10203
	ErrCodeNotFound           = 10204
)

const (
	ErrNilStage                 = "interviewstages: stage entity is nil"
	ErrEmptyStageID            = "interviewstages: stage id cannot be empty"
	ErrEmptyJobApplicationID   = "interviewstages: job application id cannot be empty"
	ErrStageNotFound           = "interviewstages: stage not found"
	ErrUnsupportedStageType    = "interviewstages: unsupported stage type"
	ErrUnsupportedOutcome      = "interviewstages: unsupported outcome"
	ErrUnableToPersist         = "interviewstages: unable to persist data"
	ErrUnableToFetch           = "interviewstages: unable to fetch data"
	ErrUnableToUpdate          = "interviewstages: unable to update data"
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

