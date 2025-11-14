package chats

import "errors"

const (
	ErrCodeInvalidPayload           = 5000
	ErrCodeInvalidTitle             = 5001
	ErrCodeInvalidRole              = 5002
	ErrCodeInvalidContent           = 5003
	ErrCodeRepositoryFailure        = 5004
	ErrCodeNotFound                 = 5005
	ErrCodeLLMFailure               = 5006
	ErrCodeConversationAccessDenied = 5007
	ErrCodeTranscriptNotFound       = 5008
	ErrCodeAssignmentFailure        = 5009
	ErrCodeStreamFailure            = 5010
)

const (
	ErrNilConversation      = "chats: conversation entity is nil"
	ErrNilMessage           = "chats: message entity is nil"
	ErrEmptyConversationID  = "chats: conversation id cannot be empty"
	ErrEmptyMessageID       = "chats: message id cannot be empty"
	ErrEmptyUserID          = "chats: user id cannot be empty"
	ErrEmptyTitle           = "chats: conversation title cannot be empty"
	ErrEmptyRole            = "chats: role cannot be empty"
	ErrEmptyContent         = "chats: content cannot be empty"
	ErrConversationNotFound = "chats: conversation not found"
	ErrMessageNotFound      = "chats: message not found"
	ErrUnableToPersist      = "chats: unable to persist data"
	ErrUnableToFetch        = "chats: unable to fetch data"
	ErrUnableToInvokeLLM    = "chats: unable to invoke LLM"
	ErrUnauthorizedAccess   = "chats: conversation access denied"
	ErrTranscriptNotFound   = "chats: transcript not found"
	ErrAssignmentNotFound   = "chats: assignment not found"
	ErrUnableToStream       = "chats: unable to stream messages"
	ErrInvalidSearchQuery   = "chats: invalid search query"
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
