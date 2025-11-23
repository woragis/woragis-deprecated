package translations

import "errors"

const (
	ErrCodeInvalidPayload      = 9000
	ErrCodeInvalidEntityType   = 9001
	ErrCodeInvalidLanguage     = 9002
	ErrCodeRepositoryFailure   = 9003
	ErrCodeNotFound            = 9004
	ErrCodeTranslationFailed   = 9005
	ErrCodeAIServiceFailure    = 9006
	ErrCodeJobQueueFailure     = 9007
)

const (
	ErrNilTranslation          = "translations: translation entity is nil"
	ErrEmptyTranslationID      = "translations: translation id cannot be empty"
	ErrEmptyEntityType         = "translations: entity type cannot be empty"
	ErrEmptyEntityID           = "translations: entity id cannot be empty"
	ErrEmptyLanguage           = "translations: language cannot be empty"
	ErrTranslationNotFound     = "translations: translation not found"
	ErrUnsupportedEntityType   = "translations: unsupported entity type"
	ErrUnsupportedLanguage     = "translations: unsupported language"
	ErrUnableToPersist         = "translations: unable to persist data"
	ErrUnableToFetch           = "translations: unable to fetch data"
	ErrUnableToUpdate          = "translations: unable to update data"
	ErrTranslationFailed       = "translations: translation failed"
	ErrAIServiceUnavailable    = "translations: AI service unavailable"
	ErrJobQueueUnavailable     = "translations: job queue unavailable"
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

