package chats

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Conversation represents a chat thread persisted by Woragis.
type Conversation struct {
	ID                uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID            uuid.UUID  `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Title             string     `gorm:"column:title;size:120;not null" json:"title"`
	Description       string     `gorm:"column:description;size:255" json:"description"`
	IdeaID            *uuid.UUID `gorm:"column:idea_id;type:uuid;index" json:"ideaId,omitempty"`
	ProjectID         *uuid.UUID `gorm:"column:project_id;type:uuid;index" json:"projectId,omitempty"`
	JobApplicationID  *uuid.UUID `gorm:"column:job_application_id;type:uuid;index" json:"jobApplicationId,omitempty"`
	AssignedAgentID   *uuid.UUID `gorm:"column:assigned_agent_id;type:uuid;index" json:"assignedAgentId,omitempty"`
	SharedTranscript  string     `gorm:"column:shared_transcript;size:255" json:"sharedTranscript"`
	ArchivedAt        *time.Time `gorm:"column:archived_at;index" json:"archivedAt,omitempty"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	LastAssignedAt    *time.Time `gorm:"column:last_assigned_at" json:"lastAssignedAt,omitempty"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// Message represents a single message in a conversation.
type Message struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ConversationID uuid.UUID `gorm:"column:conversation_id;type:uuid;index;not null" json:"conversationId"`
	Role           string    `gorm:"column:role;size:32;not null" json:"role"`
	Content        string    `gorm:"column:content;type:text;not null" json:"content"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
}

// NewConversation creates a new conversation.
func NewConversation(userID uuid.UUID, title, description string, ideaID, projectID, jobApplicationID *uuid.UUID) (*Conversation, error) {
	conv := &Conversation{
		ID:               uuid.New(),
		UserID:           userID,
		Title:            strings.TrimSpace(title),
		Description:      strings.TrimSpace(description),
		IdeaID:           ideaID,
		ProjectID:       projectID,
		JobApplicationID: jobApplicationID,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
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

// IsArchived reports if the conversation is archived.
func (c *Conversation) IsArchived() bool {
	return c.ArchivedAt != nil
}

// IsDeleted reports if the conversation is soft deleted.
func (c *Conversation) IsDeleted() bool {
	return c.DeletedAt != nil
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

// ConversationTranscript represents a shared transcript snapshot.
type ConversationTranscript struct {
	ID             uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ConversationID uuid.UUID  `gorm:"column:conversation_id;type:uuid;index;not null" json:"conversationId"`
	ShareCode      string     `gorm:"column:share_code;size:64;uniqueIndex;not null" json:"shareCode"`
	Content        string     `gorm:"column:content;type:text;not null" json:"content"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	ExpiresAt      *time.Time `gorm:"column:expires_at" json:"expiresAt,omitempty"`
}

// NewTranscript constructs a transcript entity.
func NewTranscript(conversationID uuid.UUID, shareCode, content string, ttl time.Duration) *ConversationTranscript {
	now := time.Now().UTC()
	var expiresAt *time.Time
	if ttl > 0 {
		exp := now.Add(ttl)
		expiresAt = &exp
	}
	return &ConversationTranscript{
		ID:             uuid.New(),
		ConversationID: conversationID,
		ShareCode:      shareCode,
		Content:        content,
		CreatedAt:      now,
		ExpiresAt:      expiresAt,
	}
}

// ConversationAssignment stores assignment history for agents.
type ConversationAssignment struct {
	ID             uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ConversationID uuid.UUID  `gorm:"column:conversation_id;type:uuid;index;not null" json:"conversationId"`
	AgentID        uuid.UUID  `gorm:"column:agent_id;type:uuid;index;not null" json:"agentId"`
	AgentName      string     `gorm:"column:agent_name;size:120" json:"agentName"`
	AssignedAt     time.Time  `gorm:"column:assigned_at;index" json:"assignedAt"`
	UnassignedAt   *time.Time `gorm:"column:unassigned_at;index" json:"unassignedAt,omitempty"`
	Notes          string     `gorm:"column:notes;size:255" json:"notes"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// Close marks the assignment as closed.
func (a *ConversationAssignment) Close() {
	now := time.Now().UTC()
	a.UnassignedAt = &now
	a.UpdatedAt = now
}
