package projectcasestudies

import "errors"

const (
	ErrCodeInvalidPayload    = 4100
	ErrCodeInvalidTitle      = 4101
	ErrCodeRepositoryFailure = 4102
	ErrCodeNotFound          = 4103
	ErrCodeUnauthorized      = 4104
	ErrCodeConflict          = 4105
)

const (
	ErrNilCaseStudy      = "projectcasestudies: case study entity is nil"
	ErrEmptyCaseStudyID  = "projectcasestudies: case study id cannot be empty"
	ErrEmptyProjectID    = "projectcasestudies: project id cannot be empty"
	ErrEmptyTitle        = "projectcasestudies: title cannot be empty"
	ErrEmptyDescription  = "projectcasestudies: description cannot be empty"
	ErrEmptyChallenge    = "projectcasestudies: challenge cannot be empty"
	ErrEmptySolution     = "projectcasestudies: solution cannot be empty"
	ErrCaseStudyNotFound = "projectcasestudies: case study not found"
	ErrUnableToPersist   = "projectcasestudies: unable to persist data"
	ErrUnableToFetch     = "projectcasestudies: unable to fetch data"
	ErrUnableToUpdate    = "projectcasestudies: unable to update data"
	ErrUnauthorized      = "projectcasestudies: unauthorized access"
	ErrCaseStudyAlreadyExists = "projectcasestudies: case study for this project already exists"
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
