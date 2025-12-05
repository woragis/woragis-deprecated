package scheduler

import "errors"

const (
	ErrCodeInvalidPayload    = 8000
	ErrCodeInvalidReport     = 8001
	ErrCodeInvalidAgent      = 8002
	ErrCodeInvalidFrequency  = 8003
	ErrCodeRepositoryFailure = 8004
	ErrCodeNotFound          = 8005
	ErrCodeInvalidChannel    = 8006
)

const (
	ErrNilSchedule            = "scheduler: schedule entity is nil"
	ErrEmptyScheduleID        = "scheduler: schedule id cannot be empty"
	ErrEmptyUserID            = "scheduler: user id cannot be empty"
	ErrEmptyReportType        = "scheduler: report type cannot be empty"
	ErrEmptyAgentAlias        = "scheduler: agent alias cannot be empty"
	ErrUnsupportedFrequency   = "scheduler: frequency must be daily, weekly, or custom"
	ErrWeekdayRequired        = "scheduler: weekday required for weekly schedules"
	ErrTimeRequired           = "scheduler: time of day must be provided"
	ErrRRuleRequired          = "scheduler: rrule is required for custom schedules"
	ErrUnableToComputeNextRun = "scheduler: unable to compute next run"
	ErrScheduleNotFound       = "scheduler: schedule not found"
	ErrUnableToPersist        = "scheduler: unable to persist schedule"
	ErrUnableToFetch          = "scheduler: unable to fetch schedules"
	ErrUnableToUpdate         = "scheduler: unable to update schedule"
	ErrUnableToDelete         = "scheduler: unable to delete schedule"
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
