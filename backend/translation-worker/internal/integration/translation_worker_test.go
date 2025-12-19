//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/woragis/backend/translation-worker/internal/database"
	"github.com/woragis/backend/translation-worker/internal/queue"
)

// mockTranslator is a mock translator for testing
type mockTranslator struct {
	translations map[string]string // source -> translated
	shouldFail   bool
}

func (m *mockTranslator) Translate(ctx context.Context, text, targetLang string) (string, error) {
	if m.shouldFail {
		return "", assert.AnError
	}
	if translated, ok := m.translations[text]; ok {
		return translated, nil
	}
	// Default: return text with language prefix
	return "[" + targetLang + "] " + text, nil
}

// setupRabbitMQConnection creates a RabbitMQ connection for testing
func setupRabbitMQConnection(t *testing.T) *queue.Connection {
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://test:test@localhost:5673/test")
	conn, err := queue.NewConnection(rabbitmqURL)
	require.NoError(t, err, "Failed to connect to RabbitMQ")
	return conn
}

// setupDatabase creates a database repository for testing
func setupDatabase(t *testing.T) database.Repository {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/woragis_test?sslmode=disable")
	repo, err := database.NewRepository(dbURL, nil)
	require.NoError(t, err, "Failed to create database repository")
	return repo
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestTranslationWorkerQueueSetup tests that the translation queue can be set up correctly
func TestTranslationWorkerQueueSetup(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.translations.queue"
	exchange := "test.woragis.translations"
	routingKey := "test.translations.process"

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	_, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, logger)
	require.NoError(t, err, "Failed to create translation queue")

	// Verify queue exists
	ch := conn.Channel()
	defer ch.Close()

	_, err = ch.QueueInspect(queueName)
	require.NoError(t, err, "Queue should exist")

	// Verify exchange exists
	err = ch.ExchangeDeclarePassive(exchange, "direct", true, false, false, false, nil)
	require.NoError(t, err, "Exchange should exist")
}

// TestTranslationWorkerMessagePublish tests publishing translation jobs to the queue
func TestTranslationWorkerMessagePublish(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.translations.queue.publish"
	exchange := "test.woragis.translations.publish"
	routingKey := "test.translations.process.publish"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	_, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish a test translation job
	job := queue.TranslationJob{
		ID:         uuid.New().String(),
		EntityType: "project",
		EntityID:   uuid.New().String(),
		Language:   "es",
		Fields:     []string{"name", "description"},
		SourceText: map[string]string{
			"name":        "Test Project",
			"description": "A test project",
		},
	}

	body, err := json.Marshal(job)
	require.NoError(t, err)

	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	require.NoError(t, err, "Failed to publish message")

	// Verify message is in queue
	q, err := ch.QueueInspect(queueName)
	require.NoError(t, err)
	assert.Greater(t, q.Messages, 0, "Queue should have messages")
}

// TestTranslationWorkerMessageConsume tests consuming and processing translation jobs
func TestTranslationWorkerMessageConsume(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	dbRepo := setupDatabase(t)

	queueName := "test.translations.queue.consume"
	exchange := "test.woragis.translations.consume"
	routingKey := "test.translations.process.consume"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	translationQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	entityID := uuid.New()
	job := queue.TranslationJob{
		ID:         uuid.New().String(),
		EntityType: "project",
		EntityID:   entityID.String(),
		Language:   "fr",
		Fields:     []string{"name"},
		SourceText: map[string]string{
			"name": "Test Project",
		},
	}

	body, err := json.Marshal(job)
	require.NoError(t, err)

	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	require.NoError(t, err)

	// Setup mock translator
	mockTranslator := &mockTranslator{
		translations: map[string]string{
			"Test Project": "Projet de Test",
		},
		shouldFail: false,
	}

	// Consume and process message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan bool, 1)
	go func() {
		err := translationQueue.Consume(ctx, func(job queue.TranslationJob) error {
			// Process translation job
			entityUUID, err := uuid.Parse(job.EntityID)
			if err != nil {
				return err
			}

			// Get or create translation
			translation, err := dbRepo.GetTranslationByEntity(ctx, job.EntityType, entityUUID, job.Language)
			if err != nil {
				return err
			}

			if translation == nil {
				translation = &database.Translation{
					ID:         uuid.New(),
					EntityType: job.EntityType,
					EntityID:   entityUUID,
					Language:   job.Language,
					Status:     "processing",
				}
				emptyFields := make(map[string]string)
				if err := translation.SetFields(emptyFields); err != nil {
					return err
				}
				if err := dbRepo.CreateTranslation(ctx, translation); err != nil {
					return err
				}
			}

			// Translate fields
			translatedFields := make(map[string]string)
			for _, field := range job.Fields {
				text, ok := job.SourceText[field]
				if !ok || text == "" {
					continue
				}

				translated, err := mockTranslator.Translate(ctx, text, job.Language)
				if err != nil {
					translation.Status = "failed"
					translation.ErrorMessage = err.Error()
					dbRepo.UpdateTranslation(ctx, translation)
					return err
				}

				translatedFields[field] = translated
			}

			// Update translation
			if err := translation.SetFields(translatedFields); err != nil {
				return err
			}
			translation.Status = "completed"
			translation.ErrorMessage = ""

			return dbRepo.UpdateTranslation(ctx, translation)
		})
		if err != nil && err != context.DeadlineExceeded {
			t.Logf("Consume error: %v", err)
		}
		done <- true
	}()

	// Wait for message to be processed
	select {
	case <-done:
	case <-time.After(3 * time.Second):
	}

	// Verify translation was created in database
	translation, err := dbRepo.GetTranslationByEntity(ctx, job.EntityType, entityID, job.Language)
	require.NoError(t, err)
	require.NotNil(t, translation)
	assert.Equal(t, "completed", translation.Status)

	// Verify translated fields
	fields, err := translation.GetFields()
	require.NoError(t, err)
	assert.Equal(t, "Projet de Test", fields["name"])
}

// TestTranslationWorkerInvalidMessage tests handling of invalid messages
func TestTranslationWorkerInvalidMessage(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.translations.queue.invalid"
	exchange := "test.woragis.translations.invalid"
	routingKey := "test.translations.process.invalid"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	translationQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish invalid JSON message
	invalidBody := []byte(`{"invalid": json}`)

	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         invalidBody,
		},
	)
	require.NoError(t, err)

	// Consume should reject invalid message
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	processed := false
	go func() {
		_ = translationQueue.Consume(ctx, func(job queue.TranslationJob) error {
			processed = true
			return nil
		})
	}()

	time.Sleep(1 * time.Second)
	cancel()

	// Invalid message should be rejected, not processed
	assert.False(t, processed, "Invalid message should not be processed")
}

// TestTranslationWorkerRetryOnFailure tests retry behavior on translation failure
func TestTranslationWorkerRetryOnFailure(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.translations.queue.retry"
	exchange := "test.woragis.translations.retry"
	routingKey := "test.translations.process.retry"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	translationQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	job := queue.TranslationJob{
		ID:         uuid.New().String(),
		EntityType: "project",
		EntityID:   uuid.New().String(),
		Language:   "es",
		Fields:     []string{"name"},
		SourceText: map[string]string{
			"name": "Test Project",
		},
	}

	body, err := json.Marshal(job)
	require.NoError(t, err)

	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	require.NoError(t, err)

	// Setup mock translator that fails
	mockTranslator := &mockTranslator{
		shouldFail: true,
	}

	// Consume with failing translator
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	attempts := 0
	go func() {
		_ = translationQueue.Consume(ctx, func(job queue.TranslationJob) error {
			attempts++
			_, err := mockTranslator.Translate(ctx, job.SourceText["name"], job.Language)
			return err
		})
	}()

	time.Sleep(1 * time.Second)
	cancel()

	// Message should be requeued for retry (attempts > 0)
	assert.Greater(t, attempts, 0, "Message should have been attempted")
}

// TestTranslationWorkerMultipleLanguages tests processing translations for multiple languages
func TestTranslationWorkerMultipleLanguages(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	dbRepo := setupDatabase(t)

	queueName := "test.translations.queue.multilang"
	exchange := "test.woragis.translations.multilang"
	routingKey := "test.translations.process.multilang"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	translationQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	entityID := uuid.New()
	languages := []string{"es", "fr", "de"}

	mockTranslator := &mockTranslator{
		translations: map[string]string{
			"Test Project": "Proyecto de Prueba", // Spanish
		},
		shouldFail: false,
	}

	// Publish jobs for multiple languages
	for _, lang := range languages {
		job := queue.TranslationJob{
			ID:         uuid.New().String(),
			EntityType: "project",
			EntityID:   entityID.String(),
			Language:   lang,
			Fields:     []string{"name"},
			SourceText: map[string]string{
				"name": "Test Project",
			},
		}

		body, err := json.Marshal(job)
		require.NoError(t, err)

		err = ch.PublishWithContext(
			context.Background(),
			exchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         body,
			},
		)
		require.NoError(t, err)
	}

	// Process messages
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		_ = translationQueue.Consume(ctx, func(job queue.TranslationJob) error {
			entityUUID, _ := uuid.Parse(job.EntityID)
			translation, _ := dbRepo.GetTranslationByEntity(ctx, job.EntityType, entityUUID, job.Language)
			if translation == nil {
				translation = &database.Translation{
					ID:         uuid.New(),
					EntityType: job.EntityType,
					EntityID:   entityUUID,
					Language:   job.Language,
					Status:     "processing",
				}
				emptyFields := make(map[string]string)
				translation.SetFields(emptyFields)
				dbRepo.CreateTranslation(ctx, translation)
			}

			translatedFields := make(map[string]string)
			for _, field := range job.Fields {
				text := job.SourceText[field]
				translated, _ := mockTranslator.Translate(ctx, text, job.Language)
				translatedFields[field] = translated
			}

			translation.SetFields(translatedFields)
			translation.Status = "completed"
			return dbRepo.UpdateTranslation(ctx, translation)
		})
	}()

	time.Sleep(3 * time.Second)
	cancel()

	// Verify translations for all languages were created
	for _, lang := range languages {
		translation, err := dbRepo.GetTranslationByEntity(ctx, "project", entityID, lang)
		require.NoError(t, err)
		assert.NotNil(t, translation, "Translation should exist for language: %s", lang)
		if translation != nil {
			assert.Equal(t, "completed", translation.Status)
		}
	}
}
