package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	whatsappservice "github.com/woragis/backend/server/app/internal/services/whatsapp"
	"github.com/woragis/backend/server/app/internal/workers/notifications"
	appconfig "github.com/woragis/backend/server/app/pkg/config"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	whatsappCfg := appconfig.LoadWhatsAppConfig()
	if !whatsappCfg.Enabled() {
		logger.Info("WhatsApp worker disabled, exiting")
		os.Exit(0)
	}

	// Connect to Redis
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		logger.Error("failed to parse redis url", slog.Any("error", err))
		os.Exit(1)
	}

	redisClient := redis.NewClient(opts)
	defer redisClient.Close()

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("failed to connect to redis", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("connected to redis", slog.String("url", redisURL))

	// Initialize WhatsApp notifier
	notifier, err := whatsappservice.NewWhatsmeowNotifier(whatsappCfg.SessionPath, logger)
	if err != nil {
		logger.Error("failed to initialize whatsapp notifier", slog.Any("error", err))
		os.Exit(1)
	}
	defer notifier.Disconnect()

	// Check if leader election is enabled
	leaderElectionEnabled := os.Getenv("LEADER_ELECTION_ENABLED") == "true"

	if leaderElectionEnabled {
		runWithLeaderElection(ctx, logger, redisClient, notifier)
	} else {
		runStandalone(ctx, logger, redisClient, notifier)
	}
}

func runStandalone(ctx context.Context, logger *slog.Logger, redisClient *redis.Client, notifier whatsappservice.Notifier) {
	logger.Info("starting whatsapp worker (standalone mode)")

	// Connect to WhatsApp
	connectCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := notifier.(*whatsappservice.WhatsmeowNotifier).Connect(connectCtx); err != nil {
		logger.Error("failed to connect whatsapp", slog.Any("error", err))
		os.Exit(1)
	}

	// Start worker
	if err := notifications.StartWhatsAppWorker(ctx, redisClient, notifier, logger); err != nil {
		logger.Error("failed to start whatsapp worker", slog.Any("error", err))
		os.Exit(1)
	}

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info("shutting down whatsapp worker")
}

func runWithLeaderElection(ctx context.Context, logger *slog.Logger, redisClient *redis.Client, notifier whatsappservice.Notifier) {
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		podName = "whatsapp-worker-unknown"
	}

	leaseName := os.Getenv("LEADER_ELECTION_LEASE_NAME")
	if leaseName == "" {
		leaseName = "whatsapp-worker-leader"
	}

	namespace := os.Getenv("LEADER_ELECTION_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	// Create Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Error("failed to get kubernetes config", slog.Any("error", err))
		logger.Info("falling back to standalone mode")
		runStandalone(ctx, logger, redisClient, notifier)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error("failed to create kubernetes client", slog.Any("error", err))
		logger.Info("falling back to standalone mode")
		runStandalone(ctx, logger, redisClient, notifier)
		return
	}

	// Create leader election lock
	lock, err := resourcelock.New(
		resourcelock.LeasesResourceLock,
		namespace,
		leaseName,
		clientset.CoreV1(),
		clientset.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: podName,
		},
	)
	if err != nil {
		logger.Error("failed to create leader election lock", slog.Any("error", err))
		logger.Info("falling back to standalone mode")
		runStandalone(ctx, logger, redisClient, notifier)
		return
	}

	// Leader election configuration
	leaderElectionConfig := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				logger.Info("became leader, starting whatsapp worker", slog.String("pod", podName))

				// Connect to WhatsApp (only leader connects)
				connectCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if err := notifier.(*whatsappservice.WhatsmeowNotifier).Connect(connectCtx); err != nil {
					logger.Error("failed to connect whatsapp", slog.Any("error", err))
					return
				}

				// Start worker
				if err := notifications.StartWhatsAppWorker(ctx, redisClient, notifier, logger); err != nil {
					logger.Error("failed to start whatsapp worker", slog.Any("error", err))
					return
				}

				// Keep running until context is cancelled
				<-ctx.Done()
				logger.Info("leader context cancelled, stopping worker")
			},
			OnStoppedLeading: func() {
				logger.Info("stopped leading, disconnecting whatsapp", slog.String("pod", podName))
				notifier.(*whatsappservice.WhatsmeowNotifier).Disconnect()
			},
			OnNewLeader: func(identity string) {
				if identity == podName {
					return
				}
				logger.Info("new leader elected", slog.String("leader", identity), slog.String("pod", podName))
			},
		},
	}

	// Start leader election
	leaderelection.RunOrDie(ctx, leaderElectionConfig)
}
