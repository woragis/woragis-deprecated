package chats

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Conversation represents a chat thread persisted by Woragis.
type Conversation struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID  `gorm:"type:uuid;index;not null"`
	Title       string     `gorm:"size:120;not null"`
	Description string     `gorm:"size:255"`
	IdeaID      *uuid.UUID `gorm:"type:uuid;index"`
	ProjectID   *uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Message represents a single message in a conversation.
type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	ConversationID uuid.UUID `gorm:"type:uuid;index;not null"`
	Role           string    `gorm:"size:32;not null"`
	Content        string    `gorm:"type:text;not null"`
	CreatedAt      time.Time
}

// NewConversation creates a new conversation.
func NewConversation(userID uuid.UUID, title, description string, ideaID, projectID *uuid.UUID) (*Conversation, error) {
	conv := &Conversation{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		IdeaID:      ideaID,
		ProjectID:   projectID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return conv, conv.Validate()
}

// Validate ensures conversation invariants.
func (c *Conversation) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilConversation)
	}

	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyConversationID)
	}

	if c.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if c.Title == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyTitle)
	}

	return nil
}

// Touch updates the conversation's timestamp.
func (c *Conversation) Touch() {
	c.UpdatedAt = time.Now().UTC()
}

// NewMessage constructs a new message entity.
func NewMessage(conversationID uuid.UUID, role, content string) (*Message, error) {
	msg := &Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		Role:           strings.ToLower(strings.TrimSpace(role)),
		Content:        strings.TrimSpace(content),
		CreatedAt:      time.Now().UTC(),
	}

	return msg, msg.Validate()
}

// Validate message invariants.
func (m *Message) Validate() error {
	if m == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilMessage)
	}

	if m.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyMessageID)
	}

	if m.ConversationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyConversationID)
	}

	if m.Role == "" {
		return NewDomainError(ErrCodeInvalidRole, ErrEmptyRole)
	}

	if m.Content == "" {
		return NewDomainError(ErrCodeInvalidContent, ErrEmptyContent)
	}

	return nil
}
