package ideas

import "errors"

const (
	ErrCodeInvalidPayload       = 6000
	ErrCodeInvalidTitle         = 6001
	ErrCodeInvalidRelation      = 6002
	ErrCodeRepositoryFailure    = 6003
	ErrCodeNotFound             = 6004
	ErrCodeInvalidCollaborator  = 6005
	ErrCodeCollaboratorConflict = 6006
)

const (
	ErrNilIdea                  = "ideas: idea entity is nil"
	ErrNilLink                  = "ideas: link entity is nil"
	ErrNilCollaborator          = "ideas: collaborator entity is nil"
	ErrEmptyIdeaID              = "ideas: idea id cannot be empty"
	ErrEmptyIdeaSlug            = "ideas: idea slug cannot be empty"
	ErrEmptyLinkID              = "ideas: link id cannot be empty"
	ErrEmptyUserID              = "ideas: user id cannot be empty"
	ErrEmptyTitle               = "ideas: title cannot be empty"
	ErrEmptyRelationNodes       = "ideas: relation requires source and target"
	ErrEmptyRelationLabel       = "ideas: relation label cannot be empty"
	ErrSelfRelation             = "ideas: source and target cannot be the same"
	ErrSelfCollaborator         = "ideas: cannot add yourself as collaborator"
	ErrIdeaNotFound             = "ideas: idea not found"
	ErrLinkNotFound             = "ideas: link not found"
	ErrCollaboratorExists       = "ideas: collaborator already granted"
	ErrCollaboratorUnauthorized = "ideas: collaborator is not permitted for this board"
	ErrUnableToPersist          = "ideas: unable to persist data"
	ErrUnableToFetch            = "ideas: unable to fetch data"
	ErrUnableToUpdate           = "ideas: unable to update data"
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
