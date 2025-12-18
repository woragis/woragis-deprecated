package testutil

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	aimlintegrationsdomain "github.com/woragis/backend/server/app/internal/domains/aimlintegrations"
	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	casestudiesdomain "github.com/woragis/backend/server/app/internal/domains/casestudies"
	certificationsdomain "github.com/woragis/backend/server/app/internal/domains/certifications"
	chatsdomain "github.com/woragis/backend/server/app/internal/domains/chats"
	clientsdomain "github.com/woragis/backend/server/app/internal/domains/clients"
	creativeassetsdomain "github.com/woragis/backend/server/app/internal/domains/creativeassets"
	experiencesdomain "github.com/woragis/backend/server/app/internal/domains/experiences"
	financesdomain "github.com/woragis/backend/server/app/internal/domains/finances"
	ideasdomain "github.com/woragis/backend/server/app/internal/domains/ideas"
	impactmetricsdomain "github.com/woragis/backend/server/app/internal/domains/impactmetrics"
	interestsdomain "github.com/woragis/backend/server/app/internal/domains/interests"
	jobapplicationsdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications"
	jobapplicationresponsesdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications/responses"
	jobapplicationstagesdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications/interviewstages"
	jobwebsitesdomain "github.com/woragis/backend/server/app/internal/domains/jobwebsites"
	languagesdomain "github.com/woragis/backend/server/app/internal/domains/languages"
	postcommentsdomain "github.com/woragis/backend/server/app/internal/domains/posts/comments"
	postsdomain "github.com/woragis/backend/server/app/internal/domains/posts"
	problemsolutionsdomain "github.com/woragis/backend/server/app/internal/domains/problemsolutions"
	projectcasestudiesdomain "github.com/woragis/backend/server/app/internal/domains/projects/projectcasestudies"
	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	reportsdomain "github.com/woragis/backend/server/app/internal/domains/reports"
	resumesdomain "github.com/woragis/backend/server/app/internal/domains/resumes"
	schedulerdomain "github.com/woragis/backend/server/app/internal/domains/scheduler"
	skillsdomain "github.com/woragis/backend/server/app/internal/domains/skills"
	socialmediapostsdomain "github.com/woragis/backend/server/app/internal/domains/socialmediaposts"
	systemdesignsdomain "github.com/woragis/backend/server/app/internal/domains/systemdesigns"
	technicalwritingsdomain "github.com/woragis/backend/server/app/internal/domains/technicalwritings"
	testimonialsdomain "github.com/woragis/backend/server/app/internal/domains/testimonials"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	userprofilesdomain "github.com/woragis/backend/server/app/internal/domains/userprofiles"
	userpreferencesdomain "github.com/woragis/backend/server/app/internal/domains/userpreferences"
	applogger "github.com/woragis/backend/server/app/pkg/logger"
)

// TestConfig holds test configuration
type TestConfig struct {
	DatabaseURL string
	RedisURL    string
	RabbitMQURL string
	JWTSecret   string
}

// LoadTestConfig loads test configuration from environment or uses defaults
func LoadTestConfig() TestConfig {
	return TestConfig{
		DatabaseURL: getEnv("TEST_DATABASE_URL", "postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable"),
		RedisURL:    getEnv("TEST_REDIS_URL", "redis://localhost:6380/0"),
		RabbitMQURL: getEnv("TEST_RABBITMQ_URL", "amqp://test:test@localhost:5673/test"),
		JWTSecret:   getEnv("TEST_JWT_SECRET", "test-secret-key-for-integration-tests"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *gorm.DB {
	cfg := LoadTestConfig()
	
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormlogger.Default.LogLevel(gormlogger.Silent), // Silence GORM logs in tests
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Set connection pool settings for tests
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Clean database before each test
	if err := cleanDatabase(db); err != nil {
		t.Fatalf("failed to clean database: %v", err)
	}

	// Run migrations
	MigrateTestDB(t, db)

	return db
}

// SetupTestRedis creates a test Redis connection
func SetupTestRedis(t *testing.T) *redis.Client {
	cfg := LoadTestConfig()
	
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		t.Fatalf("failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opts)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("failed to connect to test Redis: %v", err)
	}

	// Clean Redis before each test
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("failed to flush Redis: %v", err)
	}

	return client
}

// SetupTestApp creates a Fiber app instance for testing
func SetupTestApp(t *testing.T, db *gorm.DB, redisClient *redis.Client) *fiber.App {
	logger := applogger.New("test")
	
	app := fiber.New(fiber.Config{
		AppName:       "woragis-test",
		ErrorHandler:  defaultErrorHandler,
		DisableStartupMessage: true,
	})

	// Add middleware
	app.Use(applogger.RequestIDMiddleware(logger))
	app.Use(applogger.RequestLoggerMiddleware(logger))

	return app
}

// CleanupTestDB cleans up test database
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	if err := cleanDatabase(db); err != nil {
		t.Logf("warning: failed to clean database: %v", err)
	}
	
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// CleanupTestRedis cleans up test Redis
func CleanupTestRedis(t *testing.T, client *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Logf("warning: failed to flush Redis: %v", err)
	}
	
	client.Close()
}

// cleanDatabase drops all tables and recreates schema
func cleanDatabase(db *gorm.DB) error {
	// Get all table names
	var tables []string
	if err := db.Raw(`
		SELECT tablename FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename NOT LIKE 'pg_%'
	`).Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Drop all tables
	if len(tables) > 0 {
		if err := db.Exec("DROP TABLE IF EXISTS " + joinTables(tables) + " CASCADE").Error; err != nil {
			return fmt.Errorf("failed to drop tables: %w", err)
		}
	}

	return nil
}

func joinTables(tables []string) string {
	result := ""
	for i, table := range tables {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf(`"%s"`, table)
	}
	return result
}

func defaultErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}

// CreateTestUser creates a test user in the database and returns the user ID
func CreateTestUser(t *testing.T, db *gorm.DB, email, password string) uuid.UUID {
	cfg := LoadTestConfig()
	
	// Create auth repository
	authRepo := authdomain.NewGormRepository(db)
	
	// Create JWT manager
	jwtManager, err := authdomain.NewJWTManager(cfg.JWTSecret, 24*time.Hour, "woragis-test")
	if err != nil {
		t.Fatalf("failed to create JWT manager: %v", err)
	}
	
	// Create auth service (with noop email sender for tests)
	emailSender := emailservice.NewNoopSender(applogger.New("test"))
	// TokenStore is for password reset tokens - can use nil Redis for tests
	tokenStore := authdomain.NewRedisTokenStore(nil)
	logger := applogger.New("test")
	authService := authdomain.NewService(authRepo, emailSender, tokenStore, "http://localhost:8080", jwtManager, logger)
	
	// Register user
	ctx := context.Background()
	user, err := authService.RegisterUser(ctx, authdomain.RegisterRequest{
		Email:       email,
		Password:    password,
		Locale:      "en",
		DisplayName: "Test User",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	
	// Confirm email (for tests, we'll manually confirm)
	user.ConfirmEmail()
	if err := authRepo.Update(ctx, user); err != nil {
		t.Fatalf("failed to confirm user email: %v", err)
	}
	
	return user.ID
}

// GenerateTestJWT generates a test JWT token for a user
func GenerateTestJWT(t *testing.T, userID uuid.UUID, email string) string {
	cfg := LoadTestConfig()
	
	jwtManager, err := authdomain.NewJWTManager(cfg.JWTSecret, 24*time.Hour, "woragis-test")
	if err != nil {
		t.Fatalf("failed to create JWT manager: %v", err)
	}
	
	token, err := jwtManager.Generate(userID, email)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	
	return token
}

