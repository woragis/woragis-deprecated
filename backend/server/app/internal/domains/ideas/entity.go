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
	ID          uuid.UUID      `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Title       string         `gorm:"column:title;size:160;not null" json:"title"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	PosX        float64        `gorm:"column:pos_x;not null" json:"posX"`
	PosY        float64        `gorm:"column:pos_y;not null" json:"posY"`
	Color       string         `gorm:"column:color;size:16" json:"color"`
	ProjectID   *uuid.UUID     `gorm:"column:project_id;type:uuid;index" json:"projectId,omitempty"`
	Version     int            `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

// IdeaLink represents a relationship between two ideas/projects.
type IdeaLink struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	SourceIdeaID  uuid.UUID `gorm:"column:source_idea_id;type:uuid;index;not null" json:"sourceIdeaId"`
	TargetIdeaID  uuid.UUID `gorm:"column:target_idea_id;type:uuid;index;not null" json:"targetIdeaId"`
	Relation      string    `gorm:"column:relation;size:64;not null" json:"relation"`
	Weight        float64   `gorm:"column:weight;default:1" json:"weight"`
	Bidirectional bool      `gorm:"column:bidirectional;default:false" json:"bidirectional"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"createdAt"`
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
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	IdeaID      uuid.UUID `gorm:"column:idea_id;type:uuid;index;not null" json:"ideaId"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	EditorID    uuid.UUID `gorm:"column:editor_id;type:uuid;index;not null" json:"editorId"`
	Version     int       `gorm:"column:version;index;not null" json:"version"`
	Title       string    `gorm:"column:title;size:160;not null" json:"title"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	PosX        float64   `gorm:"column:pos_x;not null" json:"posX"`
	PosY        float64   `gorm:"column:pos_y;not null" json:"posY"`
	Color       string    `gorm:"column:color;size:16" json:"color"`
	ChangeType  string    `gorm:"column:change_type;size:32;not null" json:"changeType"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
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
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	OwnerID        uuid.UUID `gorm:"column:owner_id;type:uuid;not null;uniqueIndex:idx_owner_collaborator" json:"ownerId"`
	CollaboratorID uuid.UUID `gorm:"column:collaborator_id;type:uuid;not null;uniqueIndex:idx_owner_collaborator" json:"collaboratorId"`
	Role           string    `gorm:"column:role;size:32;not null" json:"role"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
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
