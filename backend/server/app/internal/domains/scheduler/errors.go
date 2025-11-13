package scheduler

import "errors"

const (
	ErrCodeInvalidPayload    = 8000
	ErrCodeInvalidReport     = 8001
	ErrCodeInvalidAgent      = 8002
	ErrCodeInvalidFrequency  = 8003
	ErrCodeRepositoryFailure = 8004
	ErrCodeNotFound          = 8005
)

const (
	ErrNilSchedule          = "scheduler: schedule entity is nil"
	ErrEmptyScheduleID      = "scheduler: schedule id cannot be empty"
	ErrEmptyUserID          = "scheduler: user id cannot be empty"
	ErrEmptyReportType      = "scheduler: report type cannot be empty"
	ErrEmptyAgentAlias      = "scheduler: agent alias cannot be empty"
	ErrUnsupportedFrequency = "scheduler: frequency must be daily or weekly"
	ErrWeekdayRequired      = "scheduler: weekday required for weekly schedules"
	ErrTimeRequired         = "scheduler: time of day must be provided"
	ErrScheduleNotFound     = "scheduler: schedule not found"
	ErrUnableToPersist      = "scheduler: unable to persist schedule"
	ErrUnableToFetch        = "scheduler: unable to fetch schedules"
	ErrUnableToUpdate       = "scheduler: unable to update schedule"
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
