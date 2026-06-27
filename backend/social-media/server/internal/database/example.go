package database

import (
	"log"

	"woragis-social-media-service/internal/config"
)

// ExampleUsage demonstrates how to use the database manager
func ExampleUsage() {
	// Load configuration
	dbCfg := config.LoadDatabaseConfig()
	redisCfg := config.LoadRedisConfig()

	// Create database manager
	dbManager, err := NewFromConfig(dbCfg, redisCfg)
	if err != nil {
		log.Fatalf("Failed to initialize database manager: %v", err)
	}
	defer func() {
		if err := dbManager.Close(); err != nil {
			log.Printf("Error closing database connections: %v", err)
		}
	}()

	// Perform health check
	if err := dbManager.HealthCheck(); err != nil {
		log.Printf("Database health check failed: %v", err)
	} else {
		log.Println("All database connections are healthy")
	}

	// Example PostgreSQL usage
	postgresDB := dbManager.GetPostgres()
	if postgresDB != nil {
		// Example: Create a simple table
		// db.AutoMigrate(&YourModel{})
		log.Println("PostgreSQL connection ready")
	}

	// Example Redis usage
	redisClient := dbManager.GetRedis()
	if redisClient != nil {
		ctx, cancel := GetRedisContext()
		defer cancel()

		// Example: Set a key-value pair
		err := redisClient.Set(ctx, "example_key", "example_value", 0).Err()
		if err != nil {
			log.Printf("Redis set operation failed: %v", err)
		} else {
			log.Println("Redis set operation successful")
		}

		// Example: Get a value
		val, err := redisClient.Get(ctx, "example_key").Result()
		if err != nil {
			log.Printf("Redis get operation failed: %v", err)
		} else {
			log.Printf("Redis get operation successful: %s", val)
		}
	}
}

// ExampleModel is a simple example model for GORM
type ExampleModel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
	UpdatedAt int64  `gorm:"autoUpdateTime"`
}

// TableName returns the table name for the ExampleModel
func (ExampleModel) TableName() string {
	return "example_models"
}
