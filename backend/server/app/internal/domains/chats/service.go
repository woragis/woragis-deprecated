package chats

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	langchain "github.com/woragis/backend/server/app/internal/services/langchain"
)

// Service orchestrates chat operations.
type Service struct {
	repo             Repository
	llmClient        *langchain.Client
	logger           *slog.Logger
	defaultModel     string
	defaultProvider  langchain.ModelProvider
	stream           *StreamHub
	transcriptTTL    time.Duration
	transcriptPrefix string
}

// NewService creates a new chat service.
func NewService(repo Repository, llmClient *langchain.Client, logger *slog.Logger, provider langchain.ModelProvider, model string, stream *StreamHub) *Service {
	return &Service{
		repo:             repo,
		llmClient:        llmClient,
		logger:           logger,
		defaultModel:     model,
		defaultProvider:  provider,
		stream:           stream,
		transcriptTTL:    7 * 24 * time.Hour,
		transcriptPrefix: "woragis-chat",
	}
}

// SetStreamHub allows updating the stream hub after construction.
func (s *Service) SetStreamHub(stream *StreamHub) {
	s.stream = stream
}

// CreateConversationRequest contains data to start a new conversation.
type CreateConversationRequest struct {
	UserID      uuid.UUID
	Title       string
	Description string
	IdeaID      *uuid.UUID
	ProjectID   *uuid.UUID
}

// SearchConversationsRequest encapsulates search parameters.
type SearchConversationsRequest struct {
	UserID          uuid.UUID
	Query           string
	IncludeArchived bool
	Limit           int
}

// AppendMessageRequest contains data to append a user message and optionally request an AI response.
type AppendMessageRequest struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	Content        string
	GenerateReply  bool
	Agent          string
	Provider       langchain.ModelProvider
	Model          string
	MaxTokens      int
	Temperature    float64
}

// BulkUpdateRequest handles archive/delete/restore operations.
type BulkUpdateRequest struct {
	UserID          uuid.UUID
	ConversationIDs []uuid.UUID
}

// ShareTranscriptRequest configures transcript creation.
type ShareTranscriptRequest struct {
	UserID         uuid.UUID
	ConversationID uuid.UUID
	ExpireAfter    time.Duration
}

// AssignConversationRequest handles agent assignment.
type AssignConversationRequest struct {
	UserID         uuid.UUID
	ConversationID uuid.UUID
	AgentID        uuid.UUID
	AgentName      string
	Notes          string
}

// CreateConversation starts a new thread.
func (s *Service) CreateConversation(ctx context.Context, req CreateConversationRequest) (*Conversation, error) {
	conversation, err := NewConversation(req.UserID, req.Title, req.Description, req.IdeaID, req.ProjectID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateConversation(ctx, conversation); err != nil {
		return nil, err
	}

	return conversation, nil
}

// ListConversations returns user threads.
func (s *Service) ListConversations(ctx context.Context, userID uuid.UUID) ([]Conversation, error) {
	return s.repo.ListConversations(ctx, userID)
}

// SearchConversations finds conversations by query.
func (s *Service) SearchConversations(ctx context.Context, req SearchConversationsRequest) ([]Conversation, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrInvalidSearchQuery)
	}

	return s.repo.SearchConversations(ctx, req.UserID, SearchFilters{
		Query:           req.Query,
		IncludeArchived: req.IncludeArchived,
		Limit:           req.Limit,
	})
}

// AppendMessage stores the message and optionally triggers LLM response.
func (s *Service) AppendMessage(ctx context.Context, req AppendMessageRequest) ([]Message, error) {
	conv, err := s.repo.GetConversation(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}
	if conv.IsDeleted() {
		return nil, NewDomainError(ErrCodeConversationAccessDenied, ErrUnauthorizedAccess)
	}

	message, err := NewMessage(req.ConversationID, req.Role, req.Content)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, err
	}
	s.broadcastMessage(conv.ID, message)

	conv.Touch()
	if err := s.repo.UpdateConversation(ctx, conv); err != nil && s.logger != nil {
		s.logger.Warn("chats: failed to persist conversation update", slog.Any("error", err))
	}

	if req.GenerateReply {
		if s.stream != nil && s.llmClient != nil {
			// Stream token deltas to clients, then persist final message
			go s.streamReply(context.Background(), conv.ID, req)
		} else {
			reply, replyErr := s.generateReply(ctx, req)
			if replyErr != nil {
				if s.logger != nil {
					s.logger.Error("chats: failed to generate reply", slog.Any("error", replyErr))
				}
			} else if reply != nil {
				s.broadcastMessage(conv.ID, reply)
			}
		}
	}

	return s.repo.ListMessages(ctx, req.ConversationID, req.UserID)
}

type streamDeltaEvent struct {
	Type           string `json:"type"`
	ConversationID string `json:"conversationId"`
	Delta          string `json:"delta"`
}

func (s *Service) streamReply(ctx context.Context, conversationID uuid.UUID, req AppendMessageRequest) {
	// Load full message history to build proper prompt
	messages, err := s.repo.ListMessages(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return
	}
	lcMessages := make([]langchain.ChatMessage, 0, len(messages))
	for _, msg := range messages {
		lcMessages = append(lcMessages, langchain.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	resp, err := s.llmClient.GenerateCompletionStream(ctx, langchain.ChatCompletionRequest{
		Provider:    req.Provider,
		Model:       req.Model,
		Messages:    lcMessages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Agent:       req.Agent,
	}, func(delta string) {
		if s.stream != nil && delta != "" {
			s.stream.Broadcast(conversationID, streamDeltaEvent{
				Type:           "delta",
				ConversationID: conversationID.String(),
				Delta:          delta,
			})
		}
	})
	if err != nil {
		if s.logger != nil {
			s.logger.Error("chats: stream reply failed", slog.Any("error", err))
		}
		return
	}

	// Persist final assistant message
	reply, err := NewMessage(conversationID, "assistant", resp.Message.Content)
	if err != nil {
		return
	}
	if err := s.repo.CreateMessage(ctx, reply); err != nil {
		return
	}
	s.broadcastMessage(conversationID, reply)
}

func (s *Service) generateReply(ctx context.Context, req AppendMessageRequest) (*Message, error) {
	if s.llmClient == nil {
		return nil, NewDomainError(ErrCodeLLMFailure, ErrUnableToInvokeLLM)
	}

	provider := req.Provider
	if provider == "" {
		provider = s.defaultProvider
	}

	model := req.Model
	if model == "" {
		model = s.defaultModel
	}

	if spec, ok := langchain.LookupModel(model); ok {
		if provider == "" {
			provider = spec.Provider
		}
		model = spec.Model
	}

	if provider == "" {
		provider = s.defaultProvider
	}

	messages, err := s.repo.ListMessages(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}

	lcMessages := make([]langchain.ChatMessage, 0, len(messages))
	for _, msg := range messages {
		lcMessages = append(lcMessages, langchain.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	resp, err := s.llmClient.GenerateCompletion(ctx, langchain.ChatCompletionRequest{
		Provider:    provider,
		Model:       model,
		Messages:    lcMessages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Agent:       req.Agent,
	})
	if err != nil {
		return nil, NewDomainError(ErrCodeLLMFailure, ErrUnableToInvokeLLM)
	}

	reply, err := NewMessage(req.ConversationID, resp.Message.Role, resp.Message.Content)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateMessage(ctx, reply); err != nil {
		return nil, err
	}

	return reply, nil
}

// ListMessages returns conversation messages.
func (s *Service) ListMessages(ctx context.Context, conversationID, userID uuid.UUID) ([]Message, error) {
	return s.repo.ListMessages(ctx, conversationID, userID)
}

// ArchiveConversations archives conversations for the user.
func (s *Service) ArchiveConversations(ctx context.Context, req BulkUpdateRequest) error {
	if len(req.ConversationIDs) == 0 {
		return nil
	}
	return s.repo.BulkArchive(ctx, req.UserID, req.ConversationIDs)
}

// DeleteConversations soft deletes conversations for the user.
func (s *Service) DeleteConversations(ctx context.Context, req BulkUpdateRequest) error {
	if len(req.ConversationIDs) == 0 {
		return nil
	}
	return s.repo.BulkDelete(ctx, req.UserID, req.ConversationIDs)
}

// RestoreConversations restores archived/deleted conversations.
func (s *Service) RestoreConversations(ctx context.Context, req BulkUpdateRequest) error {
	if len(req.ConversationIDs) == 0 {
		return nil
	}
	return s.repo.BulkRestore(ctx, req.UserID, req.ConversationIDs)
}

// ShareTranscript generates a shareable transcript snapshot.
func (s *Service) ShareTranscript(ctx context.Context, req ShareTranscriptRequest) (*ConversationTranscript, error) {
	conv, err := s.repo.GetConversation(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}

	messages, err := s.repo.ListMessages(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}

	contentBytes, err := json.Marshal(messages)
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	shareCode := s.generateShareCode()
	transcript := NewTranscript(conv.ID, shareCode, string(contentBytes), s.resolveTranscriptTTL(req.ExpireAfter))

	if err := s.repo.CreateTranscript(ctx, transcript); err != nil {
		return nil, err
	}

	conv.SharedTranscript = shareCode
	conv.Touch()
	if err := s.repo.UpdateConversation(ctx, conv); err != nil && s.logger != nil {
		s.logger.Warn("chats: failed to update conversation after transcript", slog.Any("error", err))
	}

	return transcript, nil
}

// GetTranscript retrieves a transcript by share code.
func (s *Service) GetTranscript(ctx context.Context, shareCode string) (*ConversationTranscript, error) {
	return s.repo.GetTranscriptByCode(ctx, shareCode)
}

// ListTranscripts lists existing transcripts for a conversation.
func (s *Service) ListTranscripts(ctx context.Context, userID, conversationID uuid.UUID) ([]ConversationTranscript, error) {
	return s.repo.ListTranscripts(ctx, conversationID, userID)
}

// AssignConversation stores assignment history and updates conversation.
func (s *Service) AssignConversation(ctx context.Context, req AssignConversationRequest) error {
	conv, err := s.repo.GetConversation(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return err
	}

	assignment := &ConversationAssignment{
		ID:             uuid.New(),
		ConversationID: conv.ID,
		AgentID:        req.AgentID,
		AgentName:      strings.TrimSpace(req.AgentName),
		AssignedAt:     time.Now().UTC(),
		Notes:          strings.TrimSpace(req.Notes),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	if err := s.repo.CloseAssignments(ctx, conv.ID); err != nil {
		return err
	}

	if err := s.repo.CreateAssignment(ctx, assignment); err != nil {
		return err
	}

	now := time.Now().UTC()
	conv.AssignedAgentID = &req.AgentID
	conv.LastAssignedAt = &now
	conv.Touch()

	return s.repo.UpdateConversation(ctx, conv)
}

// UnassignConversation closes active assignments.
func (s *Service) UnassignConversation(ctx context.Context, userID, conversationID uuid.UUID) error {
	if err := s.repo.CloseAssignments(ctx, conversationID); err != nil {
		return err
	}
	conv, err := s.repo.GetConversation(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	conv.AssignedAgentID = nil
	conv.Touch()
	return s.repo.UpdateConversation(ctx, conv)
}

// ListAssignments returns assignment history.
func (s *Service) ListAssignments(ctx context.Context, conversationID uuid.UUID) ([]ConversationAssignment, error) {
	return s.repo.ListAssignments(ctx, conversationID)
}

func (s *Service) broadcastMessage(conversationID uuid.UUID, message *Message) {
	if s.stream == nil || message == nil {
		return
	}
	s.stream.Broadcast(conversationID, messageResponse{
		ID:             message.ID.String(),
		ConversationID: message.ConversationID.String(),
		Role:           message.Role,
		Content:        message.Content,
		CreatedAt:      message.CreatedAt.Format(time.RFC3339),
	})
}

func (s *Service) generateShareCode() string {
	token := uuid.NewString()
	return strings.ReplaceAll(token, "-", "")
}

func (s *Service) resolveTranscriptTTL(value time.Duration) time.Duration {
	if value <= 0 {
		return s.transcriptTTL
	}
	return value
}
