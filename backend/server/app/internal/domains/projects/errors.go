package projects

import "errors"

const (
	ErrCodeInvalidPayload     = 4000
	ErrCodeInvalidName        = 4001
	ErrCodeInvalidStatus      = 4002
	ErrCodeInvalidHealthScore = 4003
	ErrCodeInvalidMetrics     = 4004
	ErrCodeRepositoryFailure  = 4005
	ErrCodeNotFound           = 4006
	ErrCodeConflict           = 4007
)

const (
	ErrNilProject            = "projects: project entity is nil"
	ErrNilMilestone          = "projects: milestone entity is nil"
	ErrEmptyProjectID        = "projects: project id cannot be empty"
	ErrEmptyMilestoneID      = "projects: milestone id cannot be empty"
	ErrEmptyUserID           = "projects: user id cannot be empty"
	ErrEmptyProjectName      = "projects: project name cannot be empty"
	ErrUnsupportedStatus     = "projects: unsupported status transition"
	ErrHealthScoreOutOfRange = "projects: health score must be between 0 and 100"
	ErrMetricsMustBePositive = "projects: metrics must be non-negative"
	ErrEmptyMilestoneTitle   = "projects: milestone title cannot be empty"
	ErrProjectNotFound       = "projects: project not found"
	ErrMilestoneNotFound     = "projects: milestone not found"
	ErrUnableToPersist       = "projects: unable to persist data"
	ErrUnableToFetch         = "projects: unable to fetch data"
	ErrUnableToUpdate        = "projects: unable to update data"
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
