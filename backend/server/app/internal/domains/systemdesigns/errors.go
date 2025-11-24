package systemdesigns

import "errors"

const (
	ErrCodeInvalidPayload    = 4200
	ErrCodeInvalidTitle      = 4201
	ErrCodeRepositoryFailure = 4202
	ErrCodeNotFound          = 4203
	ErrCodeUnauthorized      = 4204
	ErrCodeConflict          = 4205
)

const (
	ErrNilSystemDesign      = "systemdesigns: system design entity is nil"
	ErrEmptySystemDesignID  = "systemdesigns: system design id cannot be empty"
	ErrEmptyUserID          = "systemdesigns: user id cannot be empty"
	ErrEmptyTitle           = "systemdesigns: title cannot be empty"
	ErrEmptyDescription     = "systemdesigns: description cannot be empty"
	ErrSystemDesignNotFound = "systemdesigns: system design not found"
	ErrUnableToPersist      = "systemdesigns: unable to persist data"
	ErrUnableToFetch        = "systemdesigns: unable to fetch data"
	ErrUnableToUpdate       = "systemdesigns: unable to update data"
	ErrUnauthorized         = "systemdesigns: unauthorized access"
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

