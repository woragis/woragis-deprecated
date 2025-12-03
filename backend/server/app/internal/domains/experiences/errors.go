package experiences

import "errors"

const (
	ErrCodeInvalidPayload    = 9000
	ErrCodeInvalidCompany    = 9001
	ErrCodeInvalidPosition   = 9002
	ErrCodeInvalidType       = 9003
	ErrCodeRepositoryFailure = 9004
	ErrCodeNotFound          = 9005
	ErrCodeUnauthorized      = 9006
	ErrCodeConflict          = 9007
)

const (
	ErrNilExperience        = "experiences: experience entity is nil"
	ErrEmptyExperienceID    = "experiences: experience id cannot be empty"
	ErrEmptyUserID          = "experiences: user id cannot be empty"
	ErrEmptyCompany         = "experiences: company cannot be empty"
	ErrEmptyPosition        = "experiences: position cannot be empty"
	ErrExperienceNotFound   = "experiences: experience not found"
	ErrUnsupportedExperienceType = "experiences: unsupported experience type"
	ErrUnableToPersist      = "experiences: unable to persist data"
	ErrUnableToFetch        = "experiences: unable to fetch data"
	ErrUnableToUpdate       = "experiences: unable to update data"
	ErrUnauthorized         = "experiences: unauthorized access"
	ErrExperienceAlreadyExists = "experiences: experience already exists"
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

