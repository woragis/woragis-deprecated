package main

import (
	"context"
	"fmt"
	stdlog "log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	chatsdomain "github.com/woragis/backend/server/app/internal/domains/chats"
	clientsdomain "github.com/woragis/backend/server/app/internal/domains/clients"
	financesdomain "github.com/woragis/backend/server/app/internal/domains/finances"
	ideasdomain "github.com/woragis/backend/server/app/internal/domains/ideas"
	languagesdomain "github.com/woragis/backend/server/app/internal/domains/languages"
	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	reportsdomain "github.com/woragis/backend/server/app/internal/domains/reports"
	schedulerdomain "github.com/woragis/backend/server/app/internal/domains/scheduler"
	skillsdomain "github.com/woragis/backend/server/app/internal/domains/skills"
	"github.com/woragis/backend/server/app/internal/monitoring"
	emailservice "github.com/woragis/backend/server/app/internal/services/email"
	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
	whatsappservice "github.com/woragis/backend/server/app/internal/services/whatsapp"
	notifications "github.com/woragis/backend/server/app/internal/workers/notifications"
	schedulerworker "github.com/woragis/backend/server/app/internal/workers/scheduler"
	appconfig "github.com/woragis/backend/server/app/pkg/config"
	applogger "github.com/woragis/backend/server/app/pkg/logger"
)

func main() {
	cfg, err := appconfig.Load()
	if err != nil {
		stdlog.Fatalf("config: %v", err)
	}

	authCfg, err := appconfig.LoadAuthConfig()
	if err != nil {
		stdlog.Fatalf("auth config: %v", err)
	}

	oauthCfg := appconfig.LoadOAuthConfig(cfg.PublicURL)

	emailCfg, _ := appconfig.LoadEmailConfig()
	whatsappCfg := appconfig.LoadWhatsAppConfig()
	monitoringCfg := appconfig.LoadMonitoringConfig()
	aiCfg := appconfig.LoadAIConfig()
	redisCfg := appconfig.LoadRedisConfig()
	corsCfg := appconfig.LoadCORSConfig()

	slogLogger := applogger.New(cfg.Env)
	slogLogger.Info("starting woragis backend",
		slog.String("env", cfg.Env),
		slog.Int("port", cfg.Port),
	)

	db, err := connectDatabase(slogLogger)
	if err != nil {
		slogLogger.Error("database connection failed", slog.Any("error", err))
		os.Exit(1)
	}

	if err := migrate(db); err != nil {
		slogLogger.Error("database migration failed", slog.Any("error", err))
		os.Exit(1)
	}

	redisOpts, err := redis.ParseURL(redisCfg.URL)
	if err != nil {
		slogLogger.Error("redis configuration invalid", slog.Any("error", err))
		os.Exit(1)
	}

	redisClient := redis.NewClient(redisOpts)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		slogLogger.Error("redis connection failed", slog.Any("error", err))
		os.Exit(1)
	}

	workerCtx, workerCancel := context.WithCancel(context.Background())

	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	if corsCfg.Enabled {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     corsCfg.AllowedOrigins,
			AllowMethods:     corsCfg.AllowedMethods,
			AllowHeaders:     corsCfg.AllowedHeaders,
			ExposeHeaders:    corsCfg.ExposedHeaders,
			AllowCredentials: corsCfg.AllowCredentials,
			MaxAge:           corsCfg.MaxAge,
		}))
	}

	app.Use(recover.New())
	app.Use(fiberlogger.New())

	var monitoringRepo monitoring.Repository
	if monitoringCfg.Enabled && cfg.Env == "production" && monitoringCfg.DBURL != "" {
		if monitorDB, err := connectAuxDatabase(monitoringCfg.DBURL, slogLogger); err != nil {
			slogLogger.Warn("monitoring database disabled", slog.Any("error", err))
		} else {
			if err := monitorDB.AutoMigrate(&monitoring.Event{}); err != nil {
				slogLogger.Warn("monitoring migration failed", slog.Any("error", err))
			} else {
				monitoringRepo = monitoring.NewGormRepository(monitorDB)
			}
		}
	}

	monitoringService := monitoring.NewService(monitoringCfg, monitoringRepo, slogLogger)
	app.Use(monitoringService.MetricsMiddleware())
	app.Get("/metrics", adaptor.HTTPHandler(monitoringService.MetricsHandler()))

	app.Use("/metrics/stream", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Use("/api/chats/conversations/:id/stream", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/metrics/stream", websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		sendSnapshot := func() error {
			snapshot, err := monitoringService.MetricsSnapshot()
			if err != nil {
				if conn.WriteMessage(websocket.TextMessage, []byte("# error: "+err.Error())) != nil {
					return err
				}
				return err
			}
			return conn.WriteMessage(websocket.TextMessage, []byte(snapshot))
		}

		if err := sendSnapshot(); err != nil {
			return
		}

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := sendSnapshot(); err != nil {
				return
			}
		}
	}))

	api := app.Group("/api")

	var emailSender emailservice.Sender = emailservice.NewNoopSender(slogLogger)
	if emailCfg.Enabled() {
		if sender, err := emailservice.NewSMTPSender(emailCfg, slogLogger); err != nil {
			slogLogger.Warn("smtp initialization failed", slog.Any("error", err))
		} else {
			emailSender = sender
		}
	}

	var whatsappNotifier whatsappservice.Notifier = whatsappservice.NewNoopNotifier(slogLogger)
	var whatsmeowNotifier *whatsappservice.WhatsmeowNotifier
	var whatsappHandler *whatsappservice.Handler
	var whatsappService whatsappservice.Service
	if whatsappCfg.Enabled() {
		if notifier, err := whatsappservice.NewWhatsmeowNotifier(whatsappCfg.SessionPath, slogLogger); err != nil {
			slogLogger.Warn("whatsapp initialization failed, using noop", slog.Any("error", err))
		} else {
			whatsappNotifier = notifier
			whatsmeowNotifier = notifier
			// Connect to WhatsApp in a goroutine
			go func() {
				if err := notifier.Connect(workerCtx); err != nil {
					slogLogger.Error("failed to connect whatsapp", slog.Any("error", err))
				}
			}()
		}
	}
	publisher := notifications.NewPublisher(redisClient)

	tokenStore := authdomain.NewRedisTokenStore(redisClient)

	if err := notifications.StartEmailWorker(workerCtx, redisClient, emailSender, slogLogger); err != nil && slogLogger != nil {
		slogLogger.Error("failed to start email worker", slog.Any("error", err))
	}
	if err := notifications.StartWhatsAppWorker(workerCtx, redisClient, whatsappNotifier, slogLogger); err != nil && slogLogger != nil {
		slogLogger.Error("failed to start whatsapp worker", slog.Any("error", err))
	}
	authRepo := authdomain.NewGormRepository(db)
	jwtManager, err := authdomain.NewJWTManager(authCfg.JWTSecret, authCfg.JWTTTL, cfg.AppName)
	if err != nil {
		slogLogger.Error("failed to initialize jwt manager", slog.Any("error", err))
		os.Exit(1)
	}

	authService := authdomain.NewService(authRepo, emailSender, tokenStore, monitoringService, cfg.PublicURL, jwtManager, slogLogger)
	if len(oauthCfg.Providers) > 0 {
		providerSettings := make(map[authdomain.OAuthProvider]authdomain.OAuthProviderSettings)
		for key, providerCfg := range oauthCfg.Providers {
			providerID := authdomain.OAuthProvider(key)
			oauthConfig := &oauth2.Config{
				ClientID:     providerCfg.ClientID,
				ClientSecret: providerCfg.ClientSecret,
				RedirectURL:  providerCfg.RedirectURL,
				Scopes:       providerCfg.Scopes,
			}

			var userInfoURL string

			switch providerID {
			case authdomain.OAuthProviderGoogle:
				oauthConfig.Endpoint = google.Endpoint
				userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
			case authdomain.OAuthProviderGithub:
				oauthConfig.Endpoint = github.Endpoint
				userInfoURL = "https://api.github.com/user"
			case authdomain.OAuthProviderMicrosoft:
				oauthConfig.Endpoint = microsoft.AzureADEndpoint("common")
				userInfoURL = "https://graph.microsoft.com/v1.0/me"
			default:
				continue
			}

			providerSettings[providerID] = authdomain.OAuthProviderSettings{
				Name:        providerCfg.Name,
				Config:      oauthConfig,
				UserInfoURL: userInfoURL,
			}
		}

		if len(providerSettings) > 0 {
			authService.ConfigureOAuthProviders(providerSettings)
		}
	}
	authHandler := authdomain.NewHandler(authService, slogLogger)
	authdomain.SetupRoutes(api, authHandler)

	// API Key service for public GET requests
	apiKeyRepo := apikeysdomain.NewGormRepository(db)
	apiKeyService := apikeysdomain.NewService(apiKeyRepo, slogLogger)

	// IMPORTANT: Register projects and skills groups BEFORE protectedAPI
	// to ensure their middleware runs first for matching routes
	projectRepo := projectsdomain.NewGormRepository(db)
	projectService := projectsdomain.NewService(projectRepo, slogLogger)
	projectHandler := projectsdomain.NewHandler(projectService, slogLogger)
	// Projects: GET endpoints support API keys, POST/PATCH/DELETE require JWT
	projectsGroup := api.Group("/projects")
	projectsGroup.Use(apikeysdomain.RequireAPIKeyOrAuth(
		apiKeyService,
		authdomain.NewAuthMiddleware(jwtManager, slogLogger),
		slogLogger,
	))
	projectsdomain.SetupRoutes(projectsGroup, projectHandler)

	skillRepo := skillsdomain.NewGormRepository(db)
	skillService := skillsdomain.NewService(skillRepo, slogLogger)
	skillHandler := skillsdomain.NewHandler(skillService, slogLogger)
	// Skills: GET endpoints support API keys, POST/PATCH/DELETE require JWT
	skillsGroup := api.Group("/skills")
	skillsGroup.Use(apikeysdomain.RequireAPIKeyOrAuth(
		apiKeyService,
		authdomain.NewAuthMiddleware(jwtManager, slogLogger),
		slogLogger,
	))
	skillsdomain.SetupRoutes(skillsGroup, skillHandler)

	// Protected API group - requires JWT for all operations
	// Create protected routes group and apply JWT middleware
	protectedAPI := api.Group("")
	protectedAPI.Use(authdomain.NewAuthMiddleware(jwtManager, slogLogger))
	authdomain.SetupProtectedRoutes(protectedAPI, authHandler)

	if monitoringRepo != nil {
		monitoringHandler := monitoring.NewHandler(monitoringService)
		monitoring.SetupRoutes(api, monitoringHandler)
	}

	financeRepo := financesdomain.NewGormRepository(db)
	financeService := financesdomain.NewService(financeRepo, slogLogger)
	financeHandler := financesdomain.NewHandler(financeService, slogLogger)
	financesdomain.SetupRoutes(protectedAPI, financeHandler)

	languageRepo := languagesdomain.NewGormRepository(db)
	languageService := languagesdomain.NewService(languageRepo, slogLogger)
	languageHandler := languagesdomain.NewHandler(languageService, slogLogger)
	languagesdomain.SetupRoutes(protectedAPI, languageHandler)

	apiKeyHandler := apikeysdomain.NewHandler(apiKeyService, slogLogger)
	// API key management requires JWT (admin only)
	// Register routes directly on api group with middleware applied via Use
	apiKeyGroup := api.Group("/api-keys")
	apiKeyGroup.Use(authdomain.NewAuthMiddleware(jwtManager, slogLogger))
	apiKeyGroup.Post("/", apiKeyHandler.CreateAPIKey)
	apiKeyGroup.Get("/", apiKeyHandler.ListAPIKeys)
	apiKeyGroup.Get("/:id", apiKeyHandler.GetAPIKey)
	apiKeyGroup.Patch("/:id", apiKeyHandler.UpdateAPIKey)
	apiKeyGroup.Delete("/:id", apiKeyHandler.DeleteAPIKey)

	ideaRepo := ideasdomain.NewGormRepository(db)
	ideaService := ideasdomain.NewService(ideaRepo, slogLogger)
	ideaHandler := ideasdomain.NewHandler(ideaService, slogLogger)
	ideasdomain.SetupRoutes(protectedAPI, ideaHandler)

	langchainClient := langchainservice.NewClient(slogLogger)
	defaultProvider := langchainservice.ModelProvider(strings.ToLower(aiCfg.ProviderAlias))
	if defaultProvider == "" {
		defaultProvider = langchainservice.ProviderOpenAI
	}
	defaultModel := aiCfg.DefaultAlias
	if defaultModel == "" {
		defaultModel = "chatgpt"
	}

	chatsRepo := chatsdomain.NewGormRepository(db)
	chatsStream := chatsdomain.NewStreamHub()
	chatsService := chatsdomain.NewService(chatsRepo, langchainClient, slogLogger, defaultProvider, defaultModel, chatsStream)
	chatsHandler := chatsdomain.NewHandler(chatsService, slogLogger, chatsStream)
	chatsdomain.SetupRoutes(protectedAPI, chatsHandler)

	clientsRepo := clientsdomain.NewGormRepository(db)
	clientsService := clientsdomain.NewService(clientsRepo, slogLogger)
	clientsHandler := clientsdomain.NewHandler(clientsService, slogLogger)
	clientsdomain.SetupRoutes(protectedAPI, clientsHandler)

	reportsService := reportsdomain.NewService(
		reportsdomain.NewGormRepository(db),
		ideaRepo,
		projectRepo,
		financeRepo,
		chatsRepo,
		publisher,
		slogLogger,
	)
	reportsHandler := reportsdomain.NewHandler(reportsService, slogLogger)
	reportsdomain.SetupRoutes(protectedAPI, reportsHandler)

	schedulerRepo := schedulerdomain.NewGormRepository(db)
	schedulerService := schedulerdomain.NewService(schedulerRepo, reportsService, slogLogger)
	schedulerHandler := schedulerdomain.NewHandler(schedulerService, slogLogger)
	schedulerdomain.SetupRoutes(protectedAPI, schedulerHandler)

	// Initialize WhatsApp service with repository adapters (after clientsRepo and langchainClient are created)
	// Initialize even if WhatsApp isn't fully enabled, so the send endpoint is available
	userRepoAdapter := whatsappservice.NewUserRepositoryAdapter(authRepo)
	clientRepoAdapter := whatsappservice.NewClientRepositoryAdapter(clientsRepo)
	whatsappService = whatsappservice.NewService(whatsappNotifier, userRepoAdapter, clientRepoAdapter, slogLogger)
	if whatsmeowNotifier != nil {
		whatsappHandler = whatsappservice.NewHandler(whatsmeowNotifier, whatsappService, langchainClient, defaultModel, slogLogger)
	} else {
		// Create handler with nil notifier if whatsmeow isn't available (will use noop internally)
		whatsappHandler = whatsappservice.NewHandler(nil, whatsappService, langchainClient, defaultModel, slogLogger)
	}

	// Setup WhatsApp routes (handler is now initialized)
	whatsappservice.SetupRoutes(protectedAPI, whatsappHandler)

	schedulerRunner := schedulerworker.NewRunner(schedulerService, slogLogger, time.Minute)
	go schedulerRunner.Start(workerCtx)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		slogLogger.Info("http server listening", slog.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			slogLogger.Error("fiber shutdown", slog.Any("error", err))
		}
	}()

	waitForShutdown(slogLogger, app, workerCancel, redisClient, whatsmeowNotifier)
}

func connectDatabase(log *slog.Logger) (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	gormLogger := gormlogger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	config := &gorm.Config{
		Logger: gormLogger,
	}

	const maxAttempts = 5
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := openDatabase(dsn, config)
		if err == nil {
			if err := configurePool(dsn, db); err != nil {
				return nil, err
			}
			return db, nil
		}

		lastErr = err
		if log != nil {
			log.Warn("database connection attempt failed",
				slog.Int("attempt", attempt),
				slog.Int("max_attempts", maxAttempts),
				slog.Any("error", err),
			)
		}

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return nil, fmt.Errorf("connect database: %w", lastErr)
}

func openDatabase(dsn string, config *gorm.Config) (*gorm.DB, error) {
	switch {
	case strings.HasPrefix(dsn, "sqlite://"):
		sqliteDSN := strings.TrimPrefix(dsn, "sqlite://")
		return gorm.Open(sqlite.Open(sqliteDSN), config)
	default:
		return gorm.Open(postgres.Open(dsn), config)
	}
}

func configurePool(dsn string, db *gorm.DB) error {
	if strings.HasPrefix(dsn, "sqlite://") {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	return nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&authdomain.User{},
		&authdomain.Session{},
		&authdomain.Device{},
		&authdomain.MFAToken{},
		&authdomain.AuditLog{},
		&authdomain.OAuthAccount{},
		&authdomain.EmailToken{},
		&ideasdomain.Idea{},
		&ideasdomain.IdeaNode{},
		&ideasdomain.IdeaNodeConnection{},
		&ideasdomain.Document{},
		&ideasdomain.IdeaLink{},
		&ideasdomain.IdeaVersion{},
		&ideasdomain.IdeaCollaborator{},
		&chatsdomain.Conversation{},
		&chatsdomain.Message{},
		&chatsdomain.ConversationTranscript{},
		&chatsdomain.ConversationAssignment{},
		&financesdomain.Transaction{},
		&financesdomain.RecurringTemplate{},
		&languagesdomain.StudySession{},
		&languagesdomain.VocabularyEntry{},
		&projectsdomain.Project{},
		&projectsdomain.Milestone{},
		&projectsdomain.KanbanColumn{},
		&projectsdomain.KanbanCard{},
		&projectsdomain.ProjectDependency{},
		&reportsdomain.ReportDefinition{},
		&reportsdomain.ReportSchedule{},
		&reportsdomain.ReportDelivery{},
		&reportsdomain.ReportRun{},
		&schedulerdomain.Schedule{},
		&schedulerdomain.ExecutionRun{},
		&clientsdomain.Client{},
		&skillsdomain.Skill{},
		&skillsdomain.ProjectSkill{},
		&apikeysdomain.APIKey{},
	)
}

func connectAuxDatabase(dsn string, log *slog.Logger) (*gorm.DB, error) {
	gormLogger := gormlogger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	config := &gorm.Config{Logger: gormLogger}

	db, err := openDatabase(dsn, config)
	if err != nil {
		return nil, err
	}

	if err := configurePool(dsn, db); err != nil {
		return nil, err
	}

	return db, nil
}

func waitForShutdown(log *slog.Logger, app *fiber.App, cancel context.CancelFunc, redisClient *redis.Client, whatsmeowNotifier *whatsappservice.WhatsmeowNotifier) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Info("shutdown signal received")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error("error during fiber shutdown", slog.Any("error", err))
	}

	if cancel != nil {
		cancel()
	}

	if whatsmeowNotifier != nil {
		whatsmeowNotifier.Disconnect()
	}

	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Error("error closing redis client", slog.Any("error", err))
		}
	}
}
