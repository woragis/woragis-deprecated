package chats

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for chat conversations and messages.
type Repository interface {
	CreateConversation(ctx context.Context, conversation *Conversation) error
	UpdateConversation(ctx context.Context, conversation *Conversation) error
	GetConversation(ctx context.Context, id, userID uuid.UUID) (*Conversation, error)
	ListConversations(ctx context.Context, userID uuid.UUID) ([]Conversation, error)
	CreateMessage(ctx context.Context, message *Message) error
	ListMessages(ctx context.Context, conversationID, userID uuid.UUID) ([]Message, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository instantiates the repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateConversation(ctx context.Context, conversation *Conversation) error {
	if err := conversation.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(conversation).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateConversation(ctx context.Context, conversation *Conversation) error {
	if err := conversation.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(conversation).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) GetConversation(ctx context.Context, id, userID uuid.UUID) (*Conversation, error) {
	var conversation Conversation
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&conversation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrConversationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &conversation, nil
}

func (r *gormRepository) ListConversations(ctx context.Context, userID uuid.UUID) ([]Conversation, error) {
	var conversations []Conversation
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at desc").
		Find(&conversations).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return conversations, nil
}

func (r *gormRepository) CreateMessage(ctx context.Context, message *Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(message).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) ListMessages(ctx context.Context, conversationID, userID uuid.UUID) ([]Message, error) {
	var messages []Message

	err := r.db.WithContext(ctx).
		Joins("JOIN conversations ON conversations.id = messages.conversation_id").
		Where("messages.conversation_id = ? AND conversations.user_id = ?", conversationID, userID).
		Order("messages.created_at asc").
		Find(&messages).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return messages, nil
}
