package auth

import "errors"

const (
	// Public error codes (shared with HTTP responses).
	ErrCodeInvalidPayload       = 1000
	ErrCodeInvalidEmail         = 1001
	ErrCodeInvalidPassword      = 1002
	ErrCodeEmailAlreadyExists   = 1003
	ErrCodeUserNotFound         = 1004
	ErrCodeInvalidCredentials   = 1005
	ErrCodeRepositoryFailure    = 1006
	ErrCodeEmailDispatchFailure = 1007
	ErrCodePasswordHashFailure  = 1008
)

const (
	// Log-oriented error messages that enrich observability.
	ErrNilUser             = "auth: user entity is nil"
	ErrEmptyUserID         = "auth: user id cannot be empty"
	ErrEmptyEmail          = "auth: email cannot be empty"
	ErrInvalidEmailFormat  = "auth: email format is invalid"
	ErrEmptyPasswordHash   = "auth: password hash cannot be empty"
	ErrPasswordTooShort    = "auth: password must be at least 8 characters"
	ErrPasswordTooLong     = "auth: password must be at most 72 characters"
	ErrUserAlreadyExists   = "auth: user already exists"
	ErrUserNotFound        = "auth: user not found"
	ErrCredentialsMismatch = "auth: credentials mismatch"
	ErrUnableToPersist     = "auth: unable to persist user"
	ErrUnableToSendEmail   = "auth: unable to send email"
	ErrHashPassword        = "auth: unable to hash password"
	ErrMalformedPayload    = "auth: malformed payload"
)

// DomainError wraps an error message with a machine-readable code.
type DomainError struct {
	Code    int
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError constructs a new DomainError.
func NewDomainError(code int, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// AsDomainError attempts to cast an error into a DomainError.
func AsDomainError(err error) (*DomainError, bool) {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr, true
	}
	return nil, false
}
