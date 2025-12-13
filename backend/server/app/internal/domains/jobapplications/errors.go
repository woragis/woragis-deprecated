package jobapplications

import (
	"errors"
	"strings"
)

const (
	ErrCodeInvalidPayload       = 10000
	ErrCodeInvalidStatus        = 10001
	ErrCodeRepositoryFailure    = 10002
	ErrCodeNotFound             = 10003
	ErrCodeJobQueueFailure      = 10004
	ErrCodeAIServiceFailure     = 10005
	ErrCodePlaywrightFailure    = 10006
	ErrCodeAccessDenied         = 10007
	ErrCodeDatabaseConstraint   = 10008
	ErrCodeDatabaseValueTooLong = 10009
	ErrCodeDatabaseUniqueViolation = 10010
	ErrCodeDatabaseForeignKeyViolation = 10011
	ErrCodeDatabaseConnection   = 10012
)

const (
	ErrNilApplication                = "jobapplications: application entity is nil"
	ErrEmptyApplicationID            = "jobapplications: application id cannot be empty"
	ErrEmptyUserID                   = "jobapplications: user id cannot be empty"
	ErrEmptyCompanyName              = "jobapplications: company name cannot be empty"
	ErrEmptyJobTitle                 = "jobapplications: job title cannot be empty"
	ErrEmptyJobURL                   = "jobapplications: job url cannot be empty"
	ErrEmptyWebsite                  = "jobapplications: website cannot be empty"
	ErrApplicationNotFound           = "jobapplications: application not found"
	ErrUnsupportedStatus             = "jobapplications: unsupported status"
	ErrUnableToPersist               = "jobapplications: unable to persist data"
	ErrUnableToFetch                 = "jobapplications: unable to fetch data"
	ErrUnableToUpdate                = "jobapplications: unable to update data"
	ErrJobQueueUnavailable           = "jobapplications: job queue unavailable"
	ErrAIServiceUnavailable          = "jobapplications: AI service unavailable"
	ErrPlaywrightUnavailable         = "jobapplications: Playwright unavailable"
	ErrJobApplicationFailed          = "jobapplications: job application failed"
	ErrDatabaseConstraintViolation   = "jobapplications: database constraint violation"
	ErrValueTooLong                  = "jobapplications: one or more field values exceed the maximum allowed length"
	ErrDatabaseUniqueViolation       = "jobapplications: a record with this information already exists"
	ErrDatabaseForeignKeyViolation   = "jobapplications: referenced record does not exist"
	ErrDatabaseConnectionFailure     = "jobapplications: database connection error"
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

// handleDatabaseError converts database errors to domain errors.
// It checks for PostgreSQL error codes (SQLSTATE) in the error message.
func handleDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	// If it's already a DomainError, return it
	if domainErr, ok := AsDomainError(err); ok {
		return domainErr
	}

	errStr := err.Error()

	// Check for PostgreSQL SQLSTATE codes
	// SQLSTATE 22001: string data right truncated / value too long
	if strings.Contains(errStr, "SQLSTATE 22001") || strings.Contains(errStr, "value too long") {
		return NewDomainError(ErrCodeDatabaseValueTooLong, ErrValueTooLong)
	}

	// SQLSTATE 23505: unique violation
	if strings.Contains(errStr, "SQLSTATE 23505") || strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "unique constraint") {
		return NewDomainError(ErrCodeDatabaseUniqueViolation, ErrDatabaseUniqueViolation)
	}

	// SQLSTATE 23503: foreign key violation
	if strings.Contains(errStr, "SQLSTATE 23503") || strings.Contains(errStr, "foreign key constraint") {
		return NewDomainError(ErrCodeDatabaseForeignKeyViolation, ErrDatabaseForeignKeyViolation)
	}

	// SQLSTATE 23514: check constraint violation
	// SQLSTATE 23502: not null violation
	// SQLSTATE 23XXX: other constraint violations
	if strings.Contains(errStr, "SQLSTATE 23") || strings.Contains(errStr, "constraint") {
		return NewDomainError(ErrCodeDatabaseConstraint, ErrDatabaseConstraintViolation)
	}

	// Connection errors
	if strings.Contains(errStr, "connection") || strings.Contains(errStr, "dial") || strings.Contains(errStr, "network") {
		return NewDomainError(ErrCodeDatabaseConnection, ErrDatabaseConnectionFailure)
	}

	// For any other database error, return a generic repository failure
	return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
}

