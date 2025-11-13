package ideas

import (
	"strings"
	"time"

	"github.com/google/uuid"
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
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
