package technicalwritings

import "errors"

const (
	ErrCodeInvalidPayload    = 12000
	ErrCodeInvalidType       = 12001
	ErrCodeInvalidPlatform   = 12002
	ErrCodeInvalidTitle      = 12003
	ErrCodeRepositoryFailure = 12004
	ErrCodeNotFound          = 12005
	ErrCodeUnauthorized      = 12006
	ErrCodeConflict          = 12007
)

const (
	ErrNilWriting        = "technicalwritings: technical writing entity is nil"
	ErrEmptyWritingID    = "technicalwritings: writing id cannot be empty"
	ErrEmptyUserID       = "technicalwritings: user id cannot be empty"
	ErrEmptyTitle        = "technicalwritings: title cannot be empty"
	ErrEmptyDescription  = "technicalwritings: description cannot be empty"
	ErrEmptyURL          = "technicalwritings: url cannot be empty"
	ErrUnsupportedWritingType = "technicalwritings: unsupported writing type"
	ErrUnsupportedPlatform = "technicalwritings: unsupported publication platform"
	ErrWritingNotFound   = "technicalwritings: writing not found"
	ErrUnableToPersist   = "technicalwritings: unable to persist data"
	ErrUnableToFetch     = "technicalwritings: unable to fetch data"
	ErrUnableToUpdate    = "technicalwritings: unable to update data"
	ErrUnauthorized      = "technicalwritings: unauthorized access"
	ErrWritingAlreadyExists = "technicalwritings: writing already exists"
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

