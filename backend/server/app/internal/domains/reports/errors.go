package reports

import "errors"

const (
	ErrCodeInvalidPayload    = 7000
	ErrCodeRepositoryFailure = 7001
	ErrCodeNotFound          = 7002
	ErrCodeInvalidSchedule   = 7003
	ErrCodeInvalidDelivery   = 7004
	ErrCodeInvalidRun        = 7005
)

const (
	ErrUnableToGenerate         = "reports: unable to generate summary"
	ErrUnableToPersist          = "reports: unable to persist data"
	ErrUnableToFetch            = "reports: unable to fetch data"
	ErrNilReportDefinition      = "reports: definition entity is nil"
	ErrEmptyReportDefinitionID  = "reports: definition id cannot be empty"
	ErrEmptyUserID              = "reports: user id cannot be empty"
	ErrEmptyReportName          = "reports: report name cannot be empty"
	ErrNilReportSchedule        = "reports: schedule entity is nil"
	ErrEmptyScheduleID          = "reports: schedule id cannot be empty"
	ErrScheduleNotFound         = "reports: schedule not found"
	ErrNilReportDelivery        = "reports: delivery entity is nil"
	ErrEmptyDeliveryID          = "reports: delivery id cannot be empty"
	ErrEmptyDeliveryChannel     = "reports: delivery channel cannot be empty"
	ErrDeliveryNotFound         = "reports: delivery not found"
	ErrReportDefinitionNotFound = "reports: definition not found"
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
