package database

import (
	"woragis-social-media-service/internal/config"
)

// NewFromConfig creates a new database manager from application config
func NewFromConfig(dbCfg *config.DatabaseConfig, redisCfg *config.RedisConfig) (*Manager, error) {
	dbConfig := Config{
		Postgres: PostgresConfig{
			DSN:             dbCfg.URL,
			MaxOpenConns:    dbCfg.MaxOpenConns,
			MaxIdleConns:    dbCfg.MaxIdleConns,
			ConnMaxIdleTime: dbCfg.MaxIdleTime,
			ConnMaxLifetime: dbCfg.ConnMaxLifetime,
		},
		Redis: RedisConfig{
			URL:      redisCfg.URL,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		},
	}

	return NewManager(dbConfig)
}
