package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrations := []string{
		createChatsTable,
		createChatTurnsTable,
		createAIUsageStatsTable,
		createIndexes,
	}

	for _, migration := range migrations {
		if _, err := pool.Exec(ctx, migration); err != nil {
			return err
		}
	}

	return nil
}

const (
	createChatsTable = `
	CREATE TABLE IF NOT EXISTS chats (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		post_id UUID,
		prompt TEXT NOT NULL,
		response TEXT,
		agent VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		error TEXT,
		total_tokens INTEGER,
		estimated_cost DECIMAL(10, 4),
		created_at TIMESTAMP DEFAULT now(),
		completed_at TIMESTAMP,
		updated_at TIMESTAMP DEFAULT now()
	);
	`

	createChatTurnsTable = `
	CREATE TABLE IF NOT EXISTS chat_turns (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
		role VARCHAR(20) NOT NULL,
		content TEXT NOT NULL,
		tokens_used INTEGER,
		created_at TIMESTAMP DEFAULT now()
	);
	`

	createAIUsageStatsTable = `
	CREATE TABLE IF NOT EXISTS ai_usage_stats (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		agent_type VARCHAR(50),
		date DATE,
		total_requests INTEGER,
		total_tokens INTEGER,
		total_cost DECIMAL(10, 4),
		created_at TIMESTAMP DEFAULT now(),
		UNIQUE (user_id, agent_type, date)
	);
	`

	createIndexes = `
	CREATE INDEX IF NOT EXISTS idx_chats_user_id ON chats(user_id);
	CREATE INDEX IF NOT EXISTS idx_chats_post_id ON chats(post_id);
	CREATE INDEX IF NOT EXISTS idx_chats_created_at ON chats(created_at);
	CREATE INDEX IF NOT EXISTS idx_chats_status ON chats(status);
	CREATE INDEX IF NOT EXISTS idx_chat_turns_chat_id ON chat_turns(chat_id);
	CREATE INDEX IF NOT EXISTS idx_ai_usage_stats_user_date ON ai_usage_stats(user_id, date);
	`
)
