package skills

import "errors"

const (
	ErrCodeInvalidPayload     = 5000
	ErrCodeInvalidName        = 5001
	ErrCodeInvalidCategory    = 5002
	ErrCodeInvalidProficiency = 5003
	ErrCodeRepositoryFailure  = 5004
	ErrCodeNotFound           = 5005
	ErrCodeConflict           = 5006
)

const (
	ErrNilSkill          = "skills: skill entity is nil"
	ErrEmptySkillID      = "skills: skill id cannot be empty"
	ErrEmptySkillName    = "skills: skill name cannot be empty"
	ErrEmptySkillSlug    = "skills: skill slug cannot be empty"
	ErrSkillNotFound         = "skills: skill not found"
	ErrSkillAlreadyExists    = "skills: skill with this name already exists"
	ErrUnsupportedCategory   = "skills: unsupported skill category"
	ErrUnsupportedProficiency = "skills: unsupported proficiency level"
	ErrUnableToPersist       = "skills: unable to persist data"
	ErrUnableToFetch         = "skills: unable to fetch data"
	ErrUnableToUpdate        = "skills: unable to update data"
	ErrUnableToDelete        = "skills: unable to delete data"
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

