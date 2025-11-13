package chats

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	langchain "github.com/woragis/backend/server/app/internal/services/langchain"
)

// Service orchestrates chat operations.
type Service struct {
	repo            Repository
	llmClient       *langchain.Client
	logger          *slog.Logger
	defaultModel    string
	defaultProvider langchain.ModelProvider
}

// NewService creates a new chat service.
func NewService(repo Repository, llmClient *langchain.Client, logger *slog.Logger, provider langchain.ModelProvider, model string) *Service {
	return &Service{
		repo:            repo,
		llmClient:       llmClient,
		logger:          logger,
		defaultModel:    model,
		defaultProvider: provider,
	}
}

// CreateConversationRequest contains data to start a new conversation.
type CreateConversationRequest struct {
	UserID      uuid.UUID
	Title       string
	Description string
}

// AppendMessageRequest contains data to append a user message and optionally request an AI response.
type AppendMessageRequest struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	Content        string
	GenerateReply  bool
	Provider       langchain.ModelProvider
	Model          string
	MaxTokens      int
	Temperature    float64
}

// CreateConversation starts a new thread.
func (s *Service) CreateConversation(ctx context.Context, req CreateConversationRequest) (*Conversation, error) {
	conversation, err := NewConversation(req.UserID, req.Title, req.Description)
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

// AppendMessage stores the message and optionally triggers LLM response.
func (s *Service) AppendMessage(ctx context.Context, req AppendMessageRequest) ([]Message, error) {
	conv, err := s.repo.GetConversation(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, err
	}

	message, err := NewMessage(req.ConversationID, req.Role, req.Content)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, err
	}

	conv.Touch()
	if err := s.repo.UpdateConversation(ctx, conv); err != nil && s.logger != nil {
		s.logger.Warn("failed to persist conversation update", slog.Any("error", err))
	}

	if req.GenerateReply {
		if err := s.generateReply(ctx, req); err != nil {
			if s.logger != nil {
				s.logger.Error("failed to generate reply", slog.Any("error", err))
			}
		}
	}

	return s.repo.ListMessages(ctx, req.ConversationID, req.UserID)
}

func (s *Service) generateReply(ctx context.Context, req AppendMessageRequest) error {
	if s.llmClient == nil {
		return NewDomainError(ErrCodeLLMFailure, ErrUnableToInvokeLLM)
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
		return err
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
	})
	if err != nil {
		return NewDomainError(ErrCodeLLMFailure, ErrUnableToInvokeLLM)
	}

	reply, err := NewMessage(req.ConversationID, resp.Message.Role, resp.Message.Content)
	if err != nil {
		return err
	}

	return s.repo.CreateMessage(ctx, reply)
}

// ListMessages returns conversation messages.
func (s *Service) ListMessages(ctx context.Context, conversationID, userID uuid.UUID) ([]Message, error) {
	return s.repo.ListMessages(ctx, conversationID, userID)
}
