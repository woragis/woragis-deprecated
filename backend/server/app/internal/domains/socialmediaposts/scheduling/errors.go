package scheduling

import "errors"

const (
	ErrCodeInvalidPayload          = 8300
	ErrCodeInvalidStatus           = 8301
	ErrCodeInvalidStatusTransition = 8302
	ErrCodeRepositoryFailure       = 8303
	ErrCodeNotFound                = 8304
	ErrCodeConflict                = 8305
)

const (
	ErrNilScheduledPost        = "scheduling: scheduled post entity is nil"
	ErrEmptyScheduledPostID    = "scheduling: scheduled post id cannot be empty"
	ErrEmptySocialPostID       = "scheduling: social post id cannot be empty"
	ErrScheduledPostNotFound   = "scheduling: scheduled post not found"
	ErrScheduledPostExists     = "scheduling: scheduled post already exists"
	ErrUnsupportedStatus       = "scheduling: unsupported status"
	ErrInvalidStatusTransition = "scheduling: invalid status transition"
	ErrScheduleConflict        = "scheduling: schedule conflict detected"
	ErrUnableToPersist         = "scheduling: unable to persist data"
	ErrUnableToFetch           = "scheduling: unable to fetch data"
	ErrUnableToUpdate          = "scheduling: unable to update data"
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
