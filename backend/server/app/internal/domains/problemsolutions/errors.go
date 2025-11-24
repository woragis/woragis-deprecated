package problemsolutions

import "errors"

const (
	ErrCodeInvalidPayload    = 4300
	ErrCodeRepositoryFailure = 4302
	ErrCodeNotFound          = 4303
	ErrCodeUnauthorized      = 4304
	ErrCodeConflict          = 4305
)

const (
	ErrNilProblemSolution      = "problemsolutions: problem solution entity is nil"
	ErrEmptyProblemSolutionID  = "problemsolutions: problem solution id cannot be empty"
	ErrEmptyUserID             = "problemsolutions: user id cannot be empty"
	ErrEmptyProblem             = "problemsolutions: problem cannot be empty"
	ErrEmptyContext             = "problemsolutions: context cannot be empty"
	ErrEmptySolution            = "problemsolutions: solution cannot be empty"
	ErrProblemSolutionNotFound = "problemsolutions: problem solution not found"
	ErrUnableToPersist          = "problemsolutions: unable to persist data"
	ErrUnableToFetch            = "problemsolutions: unable to fetch data"
	ErrUnableToUpdate           = "problemsolutions: unable to update data"
	ErrUnauthorized             = "problemsolutions: unauthorized access"
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

