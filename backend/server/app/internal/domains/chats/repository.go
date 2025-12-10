package chats

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for chat conversations and messages.
type Repository interface {
	CreateConversation(ctx context.Context, conversation *Conversation) error
	UpdateConversation(ctx context.Context, conversation *Conversation) error
	GetConversation(ctx context.Context, id, userID uuid.UUID) (*Conversation, error)
	ListConversations(ctx context.Context, userID uuid.UUID) ([]Conversation, error)
	SearchConversations(ctx context.Context, userID uuid.UUID, filters SearchFilters) ([]Conversation, error)
	BulkArchive(ctx context.Context, userID uuid.UUID, conversationIDs []uuid.UUID) error
	BulkDelete(ctx context.Context, userID uuid.UUID, conversationIDs []uuid.UUID) error
	BulkRestore(ctx context.Context, userID uuid.UUID, conversationIDs []uuid.UUID) error
	CreateMessage(ctx context.Context, message *Message) error
	ListMessages(ctx context.Context, conversationID, userID uuid.UUID) ([]Message, error)
	CreateTranscript(ctx context.Context, transcript *ConversationTranscript) error
	GetTranscriptByCode(ctx context.Context, shareCode string) (*ConversationTranscript, error)
	ListTranscripts(ctx context.Context, conversationID, userID uuid.UUID) ([]ConversationTranscript, error)
	CreateAssignment(ctx context.Context, assignment *ConversationAssignment) error
	CloseAssignments(ctx context.Context, conversationID uuid.UUID) error
	ListAssignments(ctx context.Context, conversationID uuid.UUID) ([]ConversationAssignment, error)
	UnlinkFromJobApplication(ctx context.Context, jobApplicationID uuid.UUID) error
}

type gormRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewGormRepository instantiates the repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{
		db:     db,
		logger: slog.Default(),
	}
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
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("updated_at desc").
		Find(&conversations).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return conversations, nil
}

// SearchFilters describes conversation filtering options.
type SearchFilters struct {
	Query            string
	IncludeArchived  bool
	JobApplicationID *uuid.UUID
	Limit            int
}

func (r *gormRepository) SearchConversations(ctx context.Context, userID uuid.UUID, filters SearchFilters) ([]Conversation, error) {
	var conversations []Conversation

	// Build query using Model to ensure proper struct mapping
	db := r.db.WithContext(ctx).Model(&Conversation{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if !filters.IncludeArchived {
		db = db.Where("archived_at IS NULL")
	}

	if filters.JobApplicationID != nil {
		r.logger.Info("SearchConversations: filtering by job_application_id",
			"job_application_id", filters.JobApplicationID.String(),
			"user_id", userID.String())
		// GORM should handle UUID comparison automatically
		db = db.Where("job_application_id = ?", *filters.JobApplicationID)
	}

	if strings.TrimSpace(filters.Query) != "" {
		pattern := "%" + strings.ToLower(strings.TrimSpace(filters.Query)) + "%"
		// Reference the conversations table in the EXISTS subquery
		db = db.Where(`
			LOWER(title) LIKE ? OR
			LOWER(description) LIKE ? OR
			EXISTS (
				SELECT 1 FROM messages
				WHERE messages.conversation_id = conversations.id
				AND LOWER(messages.content) LIKE ?
			)
		`, pattern, pattern, pattern)
	}

	if filters.Limit > 0 {
		db = db.Limit(filters.Limit)
	}

	if err := db.Order("updated_at DESC").Find(&conversations).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	r.logger.Info("SearchConversations: query completed",
		"user_id", userID.String(),
		"job_application_id", func() string {
			if filters.JobApplicationID != nil {
				return filters.JobApplicationID.String()
			}
			return "nil"
		}(),
		"result_count", len(conversations))

	return conversations, nil
}

func (r *gormRepository) BulkArchive(ctx context.Context, userID uuid.UUID, conversationIDs []uuid.UUID) error {
	now := time.Now().UTC()
	if err := r.db.WithContext(ctx).Model(&Conversation{}).
		Where("user_id = ? AND id IN ?", userID, conversationIDs).
		Update("archived_at", now).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) BulkDelete(ctx context.Context, userID uuid.UUID, conversationIDs []uuid.UUID) error {
	now := time.Now().UTC()
	if err := r.db.WithContext(ctx).Model(&Conversation{}).
		Where("user_id = ? AND id IN ?", userID, conversationIDs).
		Update("deleted_at", now).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) BulkRestore(ctx context.Context, userID uuid.UUID, conversationIDs []uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&Conversation{}).
		Where("user_id = ? AND id IN ?", userID, conversationIDs).
		Updates(map[string]any{
			"archived_at": nil,
			"deleted_at":  nil,
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
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
		Where("messages.conversation_id = ? AND conversations.user_id = ? AND conversations.deleted_at IS NULL", conversationID, userID).
		Order("messages.created_at asc").
		Find(&messages).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return messages, nil
}

func (r *gormRepository) CreateTranscript(ctx context.Context, transcript *ConversationTranscript) error {
	if transcript == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrUnableToPersist)
	}
	if err := r.db.WithContext(ctx).Create(transcript).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) GetTranscriptByCode(ctx context.Context, shareCode string) (*ConversationTranscript, error) {
	var transcript ConversationTranscript
	if err := r.db.WithContext(ctx).
		Where("share_code = ?", shareCode).
		First(&transcript).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrUnableToFetch)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	if transcript.ExpiresAt != nil && transcript.ExpiresAt.Before(time.Now().UTC()) {
		return nil, NewDomainError(ErrCodeNotFound, ErrUnableToFetch)
	}
	return &transcript, nil
}

func (r *gormRepository) ListTranscripts(ctx context.Context, conversationID, userID uuid.UUID) ([]ConversationTranscript, error) {
	var transcripts []ConversationTranscript
	if err := r.db.WithContext(ctx).
		Joins("JOIN conversations ON conversations.id = conversation_transcripts.conversation_id").
		Where("conversation_id = ? AND conversations.user_id = ?", conversationID, userID).
		Order("created_at DESC").
		Find(&transcripts).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return transcripts, nil
}

func (r *gormRepository) CreateAssignment(ctx context.Context, assignment *ConversationAssignment) error {
	if assignment == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrUnableToPersist)
	}
	if err := r.db.WithContext(ctx).Create(assignment).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) CloseAssignments(ctx context.Context, conversationID uuid.UUID) error {
	now := time.Now().UTC()
	if err := r.db.WithContext(ctx).Model(&ConversationAssignment{}).
		Where("conversation_id = ? AND unassigned_at IS NULL", conversationID).
		Updates(map[string]any{
			"unassigned_at": now,
			"updated_at":    now,
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) ListAssignments(ctx context.Context, conversationID uuid.UUID) ([]ConversationAssignment, error) {
	var assignments []ConversationAssignment
	if err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("assigned_at DESC").
		Find(&assignments).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return assignments, nil
}

func (r *gormRepository) UnlinkFromJobApplication(ctx context.Context, jobApplicationID uuid.UUID) error {
	// Set job_application_id to NULL for all conversations linked to this job application
	// This preserves chat history while unlinking from the deleted application
	if err := r.db.WithContext(ctx).Model(&Conversation{}).
		Where("job_application_id = ?", jobApplicationID).
		Update("job_application_id", nil).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}
