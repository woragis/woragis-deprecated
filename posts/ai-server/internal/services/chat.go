package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/woragis/posts-ai-service/internal/models"
)

type ChatService struct {
	db *pgxpool.Pool
	ai *AIService
}

func NewChatService(db *pgxpool.Pool, ai *AIService) *ChatService {
	return &ChatService{
		db: db,
		ai: ai,
	}
}

// GenerateDraft streams AI-generated draft and persists to database
func (s *ChatService) GenerateDraft(ctx context.Context, userID uuid.UUID, postID *uuid.UUID, prompt, agent string) (string, error) {
	// Insert chat record
	var chatID uuid.UUID
	err := s.db.QueryRow(ctx,
		`INSERT INTO chats (user_id, post_id, prompt, agent, status) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		userID, postID, prompt, agent, "pending",
	).Scan(&chatID)
	if err != nil {
		return "", fmt.Errorf("failed to create chat: %w", err)
	}

	// Get streaming response from AI service
	body, err := s.ai.ChatStream(ctx, agent, prompt, WithTemperature(0.7))
	if err != nil {
		s.updateChatError(ctx, chatID, err.Error())
		return "", err
	}

	var fullResponse strings.Builder
	var tokenCount int

	// Process stream
	for chunk := range s.ai.ScanStream(body) {
		if chunk.Error != "" {
			s.updateChatError(ctx, chatID, chunk.Error)
			return "", fmt.Errorf("AI error: %s", chunk.Error)
		}

		if chunk.Delta != "" {
			fullResponse.WriteString(chunk.Delta)
			tokenCount++

			// Save turn (optional - could batch for performance)
			_ = s.saveChatTurn(ctx, chatID, "assistant", chunk.Delta, 1)
		}

		if chunk.Done {
			// Update chat as completed
			s.updateChatCompleted(ctx, chatID, fullResponse.String(), tokenCount)
		}
	}

	return fullResponse.String(), nil
}

// ImproveContent improves existing content via AI
func (s *ChatService) ImproveContent(ctx context.Context, userID, postID uuid.UUID, improvement, agent string) (string, error) {
	// Get post content first (in real implementation, would fetch from posts service)
	// For now, assume content is passed in improvement parameter

	prompt := fmt.Sprintf("Improvement request: %s", improvement)
	systemPrompt := "You are an expert content editor. Improve the article based on the user request. Return only the improved content, no explanation."

	var chatID uuid.UUID
	err := s.db.QueryRow(ctx,
		`INSERT INTO chats (user_id, post_id, prompt, agent, status) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		userID, postID, prompt, agent, "pending",
	).Scan(&chatID)
	if err != nil {
		return "", fmt.Errorf("failed to create chat: %w", err)
	}

	body, err := s.ai.ChatStream(ctx, agent, prompt, WithSystem(systemPrompt), WithTemperature(0.7))
	if err != nil {
		s.updateChatError(ctx, chatID, err.Error())
		return "", err
	}

	var fullResponse strings.Builder
	var tokenCount int

	for chunk := range s.ai.ScanStream(body) {
		if chunk.Error != "" {
			s.updateChatError(ctx, chatID, chunk.Error)
			return "", fmt.Errorf("AI error: %s", chunk.Error)
		}

		if chunk.Delta != "" {
			fullResponse.WriteString(chunk.Delta)
			tokenCount++
			_ = s.saveChatTurn(ctx, chatID, "assistant", chunk.Delta, 1)
		}

		if chunk.Done {
			s.updateChatCompleted(ctx, chatID, fullResponse.String(), tokenCount)
		}
	}

	return fullResponse.String(), nil
}

// GetChat retrieves a chat and its turns
func (s *ChatService) GetChat(ctx context.Context, chatID uuid.UUID) (*models.Chat, []models.ChatTurn, error) {
	chat := &models.Chat{}

	err := s.db.QueryRow(ctx,
		`SELECT id, user_id, post_id, prompt, response, agent, status, error, total_tokens, estimated_cost, created_at, completed_at, updated_at
		 FROM chats WHERE id = $1`,
		chatID,
	).Scan(
		&chat.ID, &chat.UserID, &chat.PostID, &chat.Prompt, &chat.Response, &chat.Agent, &chat.Status,
		&chat.Error, &chat.TotalTokens, &chat.EstimatedCost, &chat.CreatedAt, &chat.CompletedAt, &chat.UpdatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chat: %w", err)
	}

	// Get turns
	rows, err := s.db.Query(ctx,
		`SELECT id, chat_id, role, content, tokens_used, created_at FROM chat_turns WHERE chat_id = $1 ORDER BY created_at ASC`,
		chatID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chat turns: %w", err)
	}
	defer rows.Close()

	var turns []models.ChatTurn
	for rows.Next() {
		turn := models.ChatTurn{}
		err := rows.Scan(&turn.ID, &turn.ChatID, &turn.Role, &turn.Content, &turn.TokensUsed, &turn.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan turn: %w", err)
		}
		turns = append(turns, turn)
	}

	return chat, turns, nil
}

// ListChats lists chats for a user
func (s *ChatService) ListChats(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Chat, int, error) {
	// Get total count
	var total int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM chats WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count chats: %w", err)
	}

	// Get paginated results
	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, post_id, prompt, response, agent, status, error, total_tokens, estimated_cost, created_at, completed_at, updated_at
		 FROM chats WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list chats: %w", err)
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		chat := models.Chat{}
		err := rows.Scan(&chat.ID, &chat.UserID, &chat.PostID, &chat.Prompt, &chat.Response, &chat.Agent, &chat.Status,
			&chat.Error, &chat.TotalTokens, &chat.EstimatedCost, &chat.CreatedAt, &chat.CompletedAt, &chat.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan chat: %w", err)
		}
		chats = append(chats, chat)
	}

	return chats, total, nil
}

// GetUsageStats retrieves usage statistics
func (s *ChatService) GetUsageStats(ctx context.Context, userID uuid.UUID, days int) ([]models.AIUsageStat, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, agent_type, date, total_requests, total_tokens, total_cost, created_at
		 FROM ai_usage_stats WHERE user_id = $1 AND date >= $2 ORDER BY date DESC`,
		userID, startDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage stats: %w", err)
	}
	defer rows.Close()

	var stats []models.AIUsageStat
	for rows.Next() {
		stat := models.AIUsageStat{}
		err := rows.Scan(&stat.ID, &stat.UserID, &stat.AgentType, &stat.Date, &stat.TotalRequests, &stat.TotalTokens, &stat.TotalCost, &stat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stat: %w", err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// Helper functions

func (s *ChatService) saveChatTurn(ctx context.Context, chatID uuid.UUID, role, content string, tokens int) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO chat_turns (chat_id, role, content, tokens_used) VALUES ($1, $2, $3, $4)`,
		chatID, role, content, tokens,
	)
	return err
}

func (s *ChatService) updateChatCompleted(ctx context.Context, chatID uuid.UUID, response string, tokenCount int) error {
	_, err := s.db.Exec(ctx,
		`UPDATE chats SET status = $1, response = $2, total_tokens = $3, completed_at = now(), updated_at = now() WHERE id = $4`,
		"completed", response, tokenCount, chatID,
	)
	return err
}

func (s *ChatService) updateChatError(ctx context.Context, chatID uuid.UUID, errMsg string) error {
	_, err := s.db.Exec(ctx,
		`UPDATE chats SET status = $1, error = $2, updated_at = now() WHERE id = $3`,
		"error", errMsg, chatID,
	)
	return err
}
