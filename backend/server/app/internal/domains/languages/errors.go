package languages

import "errors"

const (
	ErrCodeInvalidPayload     = 3000
	ErrCodeInvalidLanguage    = 3001
	ErrCodeInvalidDuration    = 3002
	ErrCodeInvalidCompletedAt = 3003
	ErrCodeInvalidVocabulary  = 3004
	ErrCodeInvalidReviewAt    = 3005
	ErrCodeRepositoryFailure  = 3006
	ErrCodeSummaryFailure     = 3007
)

const (
	ErrNilStudySession        = "languages: study session entity is nil"
	ErrEmptySessionID         = "languages: study session id cannot be empty"
	ErrEmptyUserID            = "languages: user id cannot be empty"
	ErrEmptyLanguageCode      = "languages: language code cannot be empty"
	ErrInvalidLanguageCode    = "languages: language code must be between 2 and 8 characters"
	ErrDurationMustBePositive = "languages: duration must be positive"
	ErrCompletedAtRequired    = "languages: completed_at must be specified"

	ErrNilVocabularyEntry = "languages: vocabulary entry is nil"
	ErrEmptyVocabularyID  = "languages: vocabulary id cannot be empty"
	ErrEmptyTerm          = "languages: vocabulary term cannot be empty"
	ErrEmptyTranslation   = "languages: vocabulary translation cannot be empty"
	ErrReviewAtRequired   = "languages: review_at must be specified"

	ErrUnableToPersist   = "languages: unable to persist record"
	ErrUnableToFetch     = "languages: unable to fetch records"
	ErrUnableToSummarize = "languages: unable to build summary"
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
