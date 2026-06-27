package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	UserID         uuid.UUID  `db:"user_id" json:"user_id"`
	PostID         *uuid.UUID `db:"post_id" json:"post_id,omitempty"`
	Prompt         string     `db:"prompt" json:"prompt"`
	Response       *string    `db:"response" json:"response,omitempty"`
	Agent          string     `db:"agent" json:"agent"`
	Status         string     `db:"status" json:"status"`
	Error          *string    `db:"error" json:"error,omitempty"`
	TotalTokens    *int       `db:"total_tokens" json:"total_tokens,omitempty"`
	EstimatedCost  *float64   `db:"estimated_cost" json:"estimated_cost,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	CompletedAt    *time.Time `db:"completed_at" json:"completed_at,omitempty"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

type ChatTurn struct {
	ID         uuid.UUID `db:"id" json:"id"`
	ChatID     uuid.UUID `db:"chat_id" json:"chat_id"`
	Role       string    `db:"role" json:"role"` // "user" or "assistant"
	Content    string    `db:"content" json:"content"`
	TokensUsed *int      `db:"tokens_used" json:"tokens_used,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type AIUsageStat struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	UserID        uuid.UUID  `db:"user_id" json:"user_id"`
	AgentType     *string    `db:"agent_type" json:"agent_type,omitempty"`
	Date          time.Time  `db:"date" json:"date"`
	TotalRequests int        `db:"total_requests" json:"total_requests"`
	TotalTokens   int        `db:"total_tokens" json:"total_tokens"`
	TotalCost     float64    `db:"total_cost" json:"total_cost"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}

type StreamChunk struct {
	Delta  string `json:"delta,omitempty"`
	Done   bool   `json:"done,omitempty"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}
