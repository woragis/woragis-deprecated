package platforms

import "errors"

const (
	ErrCodeInvalidPayload          = 8100
	ErrCodeInvalidPlatform         = 8101
	ErrCodeRepositoryFailure       = 8102
	ErrCodeNotFound                = 8103
	ErrCodeConflict                = 8104
)

const (
	ErrNilPlatformConfig      = "platforms: platform config entity is nil"
	ErrEmptyPlatformConfigID   = "platforms: platform config id cannot be empty"
	ErrEmptyDisplayName        = "platforms: display name cannot be empty"
	ErrPlatformConfigNotFound  = "platforms: platform config not found"
	ErrPlatformConfigExists    = "platforms: platform config already exists"
	ErrUnsupportedPlatform     = "platforms: unsupported platform"
	ErrUnableToPersist          = "platforms: unable to persist data"
	ErrUnableToFetch            = "platforms: unable to fetch data"
	ErrUnableToUpdate           = "platforms: unable to update data"
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
