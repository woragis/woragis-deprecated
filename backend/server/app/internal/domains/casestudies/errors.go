package casestudies

import "errors"

const (
	ErrCodeInvalidPayload    = 9000
	ErrCodeInvalidTitle      = 9001
	ErrCodeRepositoryFailure = 9002
	ErrCodeNotFound          = 9003
	ErrCodeUnauthorized      = 9004
	ErrCodeConflict          = 9005
)

const (
	ErrNilCaseStudy      = "casestudies: case study entity is nil"
	ErrEmptyCaseStudyID  = "casestudies: case study id cannot be empty"
	ErrEmptyUserID       = "casestudies: user id cannot be empty"
	ErrEmptyProjectID    = "casestudies: project id cannot be empty"
	ErrEmptyProjectSlug  = "casestudies: project slug cannot be empty"
	ErrEmptyTitle        = "casestudies: title cannot be empty"
	ErrEmptyProblem      = "casestudies: problem cannot be empty"
	ErrEmptyContext      = "casestudies: context cannot be empty"
	ErrEmptySolution     = "casestudies: solution cannot be empty"
	ErrCaseStudyNotFound = "casestudies: case study not found"
	ErrUnableToPersist   = "casestudies: unable to persist data"
	ErrUnableToFetch      = "casestudies: unable to fetch data"
	ErrUnableToUpdate    = "casestudies: unable to update data"
	ErrUnauthorized      = "casestudies: unauthorized access"
	ErrCaseStudyAlreadyExists = "casestudies: case study for this project already exists"
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

