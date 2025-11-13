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
	ErrCodeInvalidToken         = 1009
	ErrCodeTokenIssuanceFailure = 1010
	ErrCodeEmailNotConfirmed    = 1011
	ErrCodeSessionRevoked       = 1012
	ErrCodeSessionExpired       = 1013
	ErrCodeMFARequired          = 1014
	ErrCodeMFAInvalid           = 1015
	ErrCodeMFAGenerationFailure = 1016
	ErrCodeOAuthLinkFailure     = 1017
	ErrCodeRateLimited          = 1018
)

const (
	// Log-oriented error messages that enrich observability.
	ErrNilUser                   = "auth: user entity is nil"
	ErrEmptyUserID               = "auth: user id cannot be empty"
	ErrEmptyEmail                = "auth: email cannot be empty"
	ErrInvalidEmailFormat        = "auth: email format is invalid"
	ErrEmptyPasswordHash         = "auth: password hash cannot be empty"
	ErrPasswordTooShort          = "auth: password must be at least 8 characters"
	ErrPasswordTooLong           = "auth: password must be at most 72 characters"
	ErrUserAlreadyExists         = "auth: user already exists"
	ErrUserNotFound              = "auth: user not found"
	ErrCredentialsMismatch       = "auth: credentials mismatch"
	ErrUnableToPersist           = "auth: unable to persist user"
	ErrUnableToSendEmail         = "auth: unable to send email"
	ErrHashPassword              = "auth: unable to hash password"
	ErrMalformedPayload          = "auth: malformed payload"
	ErrInvalidResetToken         = "auth: reset token invalid or expired"
	ErrUnableToIssueToken        = "auth: unable to issue access token"
	ErrEmptyDeviceID             = "auth: device id cannot be empty"
	ErrEmptySessionID            = "auth: session id cannot be empty"
	ErrEmptyRefreshTokenHash     = "auth: refresh token hash cannot be empty"
	ErrInvalidExpiry             = "auth: expiry cannot be before issued time"
	ErrNilSession                = "auth: session entity is nil"
	ErrInvalidIPAddress          = "auth: ip address is invalid"
	ErrEmptyDeviceFingerprint    = "auth: device fingerprint cannot be empty"
	ErrEmptyMFAType              = "auth: mfa type cannot be empty"
	ErrUnableToGenerateMFASecret = "auth: unable to generate mfa secret"
	ErrEmptyAuditAction          = "auth: audit action cannot be empty"
	ErrEmptyOAuthProvider        = "auth: oauth provider cannot be empty"
	ErrEmptyOAuthSubject         = "auth: oauth provider user id cannot be empty"
	ErrEmptyEmailTokenType       = "auth: email token type cannot be empty"
	ErrEmptyEmailToken           = "auth: email token cannot be empty"
	ErrEmailAlreadyConfirmed     = "auth: email already confirmed"
	ErrEmailNotConfirmed         = "auth: email not confirmed"
	ErrSessionNotFound           = "auth: session not found"
	ErrSessionExpired            = "auth: session expired"
	ErrSessionRevoked            = "auth: session revoked"
	ErrMFAEnrollmentRequired     = "auth: multi-factor authentication required"
	ErrMFAInvalidCode            = "auth: multi-factor code invalid"
	ErrRateLimited               = "auth: operation temporarily rate limited"
	ErrNilDevice                 = "auth: device entity is nil"
	ErrNilMFAToken               = "auth: mfa token entity is nil"
	ErrMFATokenNotFound          = "auth: mfa token not found"
	ErrNilAuditLog               = "auth: audit log entry is nil"
	ErrNilOAuthAccount           = "auth: oauth account entity is nil"
	ErrNilEmailToken             = "auth: email token entity is nil"
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
