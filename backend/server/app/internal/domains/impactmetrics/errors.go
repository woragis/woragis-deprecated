package impactmetrics

import "errors"

const (
	ErrCodeInvalidPayload    = 10000
	ErrCodeInvalidType       = 10001
	ErrCodeInvalidUnit       = 10002
	ErrCodeInvalidValue      = 10003
	ErrCodeInvalidEntityType = 10004
	ErrCodeInvalidDate       = 10005
	ErrCodeRepositoryFailure = 10006
	ErrCodeNotFound          = 10007
	ErrCodeUnauthorized      = 10008
	ErrCodeConflict          = 10009
)

const (
	ErrNilMetric            = "impactmetrics: impact metric entity is nil"
	ErrEmptyMetricID        = "impactmetrics: metric id cannot be empty"
	ErrEmptyUserID          = "impactmetrics: user id cannot be empty"
	ErrUnsupportedMetricType = "impactmetrics: unsupported metric type"
	ErrUnsupportedMetricUnit = "impactmetrics: unsupported metric unit"
	ErrNegativeValue        = "impactmetrics: metric value cannot be negative"
	ErrUnsupportedEntityType = "impactmetrics: unsupported entity type"
	ErrPeriodEndBeforeStart = "impactmetrics: period end date cannot be before start date"
	ErrMetricNotFound       = "impactmetrics: metric not found"
	ErrUnableToPersist      = "impactmetrics: unable to persist data"
	ErrUnableToFetch        = "impactmetrics: unable to fetch data"
	ErrUnableToUpdate       = "impactmetrics: unable to update data"
	ErrUnauthorized         = "impactmetrics: unauthorized access"
	ErrMetricAlreadyExists  = "impactmetrics: metric already exists"
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

