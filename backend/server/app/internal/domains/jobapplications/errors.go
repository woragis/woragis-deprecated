package jobapplications

import "errors"

const (
	ErrCodeInvalidPayload    = 10000
	ErrCodeInvalidStatus     = 10001
	ErrCodeRepositoryFailure = 10002
	ErrCodeNotFound          = 10003
	ErrCodeJobQueueFailure   = 10004
	ErrCodeAIServiceFailure  = 10005
	ErrCodePlaywrightFailure = 10006
)

const (
	ErrNilApplication          = "jobapplications: application entity is nil"
	ErrEmptyApplicationID      = "jobapplications: application id cannot be empty"
	ErrEmptyUserID             = "jobapplications: user id cannot be empty"
	ErrEmptyCompanyName        = "jobapplications: company name cannot be empty"
	ErrEmptyJobTitle           = "jobapplications: job title cannot be empty"
	ErrEmptyJobURL             = "jobapplications: job url cannot be empty"
	ErrEmptyWebsite            = "jobapplications: website cannot be empty"
	ErrApplicationNotFound     = "jobapplications: application not found"
	ErrUnsupportedStatus       = "jobapplications: unsupported status"
	ErrUnableToPersist         = "jobapplications: unable to persist data"
	ErrUnableToFetch           = "jobapplications: unable to fetch data"
	ErrUnableToUpdate          = "jobapplications: unable to update data"
	ErrJobQueueUnavailable     = "jobapplications: job queue unavailable"
	ErrAIServiceUnavailable    = "jobapplications: AI service unavailable"
	ErrPlaywrightUnavailable   = "jobapplications: Playwright unavailable"
	ErrJobApplicationFailed    = "jobapplications: job application failed"
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

