package aimlintegrations

import "errors"

const (
	ErrCodeInvalidPayload    = 11000
	ErrCodeInvalidType       = 11001
	ErrCodeInvalidFramework  = 11002
	ErrCodeInvalidTitle      = 11003
	ErrCodeRepositoryFailure = 11004
	ErrCodeNotFound          = 11005
	ErrCodeUnauthorized      = 11006
	ErrCodeConflict          = 11007
)

const (
	ErrNilIntegration        = "aimlintegrations: AI/ML integration entity is nil"
	ErrEmptyIntegrationID    = "aimlintegrations: integration id cannot be empty"
	ErrEmptyUserID           = "aimlintegrations: user id cannot be empty"
	ErrEmptyTitle            = "aimlintegrations: title cannot be empty"
	ErrEmptyDescription      = "aimlintegrations: description cannot be empty"
	ErrUnsupportedIntegrationType = "aimlintegrations: unsupported integration type"
	ErrUnsupportedFramework  = "aimlintegrations: unsupported framework"
	ErrIntegrationNotFound   = "aimlintegrations: integration not found"
	ErrUnableToPersist       = "aimlintegrations: unable to persist data"
	ErrUnableToFetch         = "aimlintegrations: unable to fetch data"
	ErrUnableToUpdate        = "aimlintegrations: unable to update data"
	ErrUnauthorized          = "aimlintegrations: unauthorized access"
	ErrIntegrationAlreadyExists = "aimlintegrations: integration already exists"
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

