package ideas

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ChangeTypeCreated = "created"
	ChangeTypeEdited  = "edited"
	ChangeTypeMoved   = "moved"
	ChangeTypeBulk    = "bulk"
)

// Idea represents a graph node in the ideas canvas.
type Idea struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID  `gorm:"type:uuid;index;not null"`
	Title       string     `gorm:"size:160;not null"`
	Description string     `gorm:"type:text"`
	PosX        float64    `gorm:"not null"`
	PosY        float64    `gorm:"not null"`
	Color       string     `gorm:"size:16"`
	ProjectID   *uuid.UUID `gorm:"type:uuid;index"`
	Version     int        `gorm:"not null;default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// IdeaLink represents a relationship between two ideas/projects.
type IdeaLink struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;index;not null"`
	SourceIdeaID  uuid.UUID `gorm:"type:uuid;index;not null"`
	TargetIdeaID  uuid.UUID `gorm:"type:uuid;index;not null"`
	Relation      string    `gorm:"size:64;not null"`
	Weight        float64   `gorm:"default:1"`
	Bidirectional bool      `gorm:"default:false"`
	CreatedAt     time.Time
}

// NewIdea constructs a new idea node.
func NewIdea(userID uuid.UUID, title, description string, posX, posY float64, color string, projectID *uuid.UUID) (*Idea, error) {
	idea := &Idea{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		PosX:        posX,
		PosY:        posY,
		Color:       strings.TrimSpace(color),
		ProjectID:   projectID,
		Version:     1,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return idea, idea.Validate()
}

// Validate ensures idea invariants.
func (i *Idea) Validate() error {
	if i == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilIdea)
	}

	if i.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyIdeaID)
	}

	if i.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if strings.TrimSpace(i.Title) == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyTitle)
	}

	return nil
}

// Touch updates the updated at timestamp.
func (i *Idea) Touch() {
	i.UpdatedAt = time.Now().UTC()
}

// Move updates canvas coordinates.
func (i *Idea) Move(posX, posY float64) {
	i.PosX = posX
	i.PosY = posY
	i.Version++
	i.Touch()
}

// UpdateDetails updates textual metadata.
func (i *Idea) UpdateDetails(title, description, color string) error {
	if title != "" {
		i.Title = strings.TrimSpace(title)
	}
	if description != "" {
		i.Description = strings.TrimSpace(description)
	}
	if color != "" {
		i.Color = strings.TrimSpace(color)
	}
	i.Version++
	return i.Validate()
}

// NewIdeaLink creates a connection between ideas.
func NewIdeaLink(userID, sourceID, targetID uuid.UUID, relation string, weight float64, bidirectional bool) (*IdeaLink, error) {
	link := &IdeaLink{
		ID:            uuid.New(),
		UserID:        userID,
		SourceIdeaID:  sourceID,
		TargetIdeaID:  targetID,
		Relation:      strings.TrimSpace(relation),
		Weight:        weight,
		Bidirectional: bidirectional,
		CreatedAt:     time.Now().UTC(),
	}

	return link, link.Validate()
}

// Validate ensures IdeaLink invariants.
func (l *IdeaLink) Validate() error {
	if l == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilLink)
	}

	if l.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}

	if l.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if l.SourceIdeaID == uuid.Nil || l.TargetIdeaID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidRelation, ErrEmptyRelationNodes)
	}

	if l.SourceIdeaID == l.TargetIdeaID {
		return NewDomainError(ErrCodeInvalidRelation, ErrSelfRelation)
	}

	if strings.TrimSpace(l.Relation) == "" {
		return NewDomainError(ErrCodeInvalidRelation, ErrEmptyRelationLabel)
	}

	return nil
}

// IdeaVersion captures a snapshot of idea metadata.
type IdeaVersion struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	IdeaID      uuid.UUID `gorm:"type:uuid;index;not null"`
	UserID      uuid.UUID `gorm:"type:uuid;index;not null"`
	EditorID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Version     int       `gorm:"index;not null"`
	Title       string    `gorm:"size:160;not null"`
	Description string    `gorm:"type:text"`
	PosX        float64   `gorm:"not null"`
	PosY        float64   `gorm:"not null"`
	Color       string    `gorm:"size:16"`
	ChangeType  string    `gorm:"size:32;not null"`
	CreatedAt   time.Time
}

// NewIdeaVersion builds a snapshot representation.
func NewIdeaVersion(idea *Idea, editorID uuid.UUID, changeType string) *IdeaVersion {
	if changeType == "" {
		changeType = ChangeTypeEdited
	}
	return &IdeaVersion{
		ID:          uuid.New(),
		IdeaID:      idea.ID,
		UserID:      idea.UserID,
		EditorID:    editorID,
		Version:     idea.Version,
		Title:       idea.Title,
		Description: idea.Description,
		PosX:        idea.PosX,
		PosY:        idea.PosY,
		Color:       idea.Color,
		ChangeType:  strings.ToLower(strings.TrimSpace(changeType)),
		CreatedAt:   time.Now().UTC(),
	}
}

// IdeaCollaborator tracks shared access to an idea canvas.
type IdeaCollaborator struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_owner_collaborator"`
	CollaboratorID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_owner_collaborator"`
	Role           string    `gorm:"size:32;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewIdeaCollaborator constructs a collaborator entry.
func NewIdeaCollaborator(ownerID, collaboratorID uuid.UUID, role string) (*IdeaCollaborator, error) {
	entry := &IdeaCollaborator{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		CollaboratorID: collaboratorID,
		Role:           strings.ToLower(strings.TrimSpace(role)),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	return entry, entry.Validate()
}

// Validate ensures collaborator invariants.
func (c *IdeaCollaborator) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCollaborator)
	}
	if c.OwnerID == uuid.Nil || c.CollaboratorID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if c.OwnerID == c.CollaboratorID {
		return NewDomainError(ErrCodeInvalidCollaborator, ErrSelfCollaborator)
	}
	if c.Role == "" {
		c.Role = "editor"
	}
	return nil
}
