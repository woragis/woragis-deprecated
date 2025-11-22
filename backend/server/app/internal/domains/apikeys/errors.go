package apikeys

import "errors"

const (
	ErrCodeInvalidPayload    = 6000
	ErrCodeInvalidName       = 6001
	ErrCodeRepositoryFailure = 6002
	ErrCodeNotFound          = 6003
	ErrCodeInvalidToken      = 6004
	ErrCodeExpiredToken      = 6005
)

const (
	ErrNilAPIKey        = "apikeys: api key entity is nil"
	ErrEmptyAPIKeyID    = "apikeys: api key id cannot be empty"
	ErrEmptyUserID      = "apikeys: user id cannot be empty"
	ErrEmptyAPIKeyName  = "apikeys: api key name cannot be empty"
	ErrEmptyKeyHash     = "apikeys: key hash cannot be empty"
	ErrEmptyPrefix      = "apikeys: prefix cannot be empty"
	ErrAPIKeyNotFound   = "apikeys: api key not found"
	ErrInvalidAPIKey    = "apikeys: invalid api key"
	ErrExpiredAPIKey    = "apikeys: api key has expired"
	ErrUnableToPersist  = "apikeys: unable to persist data"
	ErrUnableToFetch    = "apikeys: unable to fetch data"
	ErrUnableToUpdate   = "apikeys: unable to update data"
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

