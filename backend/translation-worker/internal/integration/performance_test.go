//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/woragis/backend/translation-worker/internal/database"
	"github.com/woragis/backend/translation-worker/internal/queue"
)

// BenchmarkTranslationWorkerThroughput benchmarks translation worker throughput
func BenchmarkTranslationWorkerThroughput(b *testing.B) {
	conn := setupRabbitMQConnection(b)
	defer conn.Close()

	dbRepo := setupDatabase(b)

	queueName := fmt.Sprintf("bench.translations.queue.%d", time.Now().Unix())
	exchange := fmt.Sprintf("bench.woragis.translations.%d", time.Now().Unix())
	routingKey := "bench.translations.process"

	// Create queue
	translationQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, nil)
	require.NoError(b, err)

	ch := conn.Channel()
	defer ch.Close()

	// Setup mock translator
	mockTranslator := &mockTranslator{
		translations: make(map[string]string),
		shouldFail:   false,
	}

	// Start consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var processedCount int64
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
			dbRepo.UpdateTranslation(ctx, translation)
			atomic.AddInt64(&processedCount, 1)
			return nil
		})
	}()

	// Give consumer time to start
	time.Sleep(100 * time.Millisecond)

	// Benchmark publishing jobs
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			job := queue.TranslationJob{
				ID:         uuid.New().String(),
				EntityType: "project",
				EntityID:   uuid.New().String(),
				Language:   "es",
				Fields:     []string{"name"},
				SourceText: map[string]string{
					"name": "Benchmark Project",
				},
			}

			body, _ := json.Marshal(job)
			_ = ch.PublishWithContext(
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
		}
	})

	// Wait for processing
	time.Sleep(2 * time.Second)
	b.StopTimer()

	b.Logf("Processed %d jobs", atomic.LoadInt64(&processedCount))
}

// TestTranslationWorkerLoadTest tests translation worker under load
func TestTranslationWorkerLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	dbRepo := setupDatabase(t)

	queueName := fmt.Sprintf("load.translations.queue.%d", time.Now().Unix())
	exchange := fmt.Sprintf("load.woragis.translations.%d", time.Now().Unix())
	routingKey := "load.translations.process"

	// Create queue
	translationQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, nil)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Setup mock translator
	mockTranslator := &mockTranslator{
		translations: map[string]string{
			"Test Project": "Proyecto de Prueba",
		},
		shouldFail: false,
	}

	// Start consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var processedCount int64
	var errorCount int64
	startTime := time.Now()

	go func() {
		_ = translationQueue.Consume(ctx, func(job queue.TranslationJob) error {
			entityUUID, err := uuid.Parse(job.EntityID)
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				return err
			}

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
				translated, err := mockTranslator.Translate(ctx, text, job.Language)
				if err != nil {
					atomic.AddInt64(&errorCount, 1)
					translation.Status = "failed"
					translation.ErrorMessage = err.Error()
					dbRepo.UpdateTranslation(ctx, translation)
					return err
				}
				translatedFields[field] = translated
			}

			translation.SetFields(translatedFields)
			translation.Status = "completed"
			dbRepo.UpdateTranslation(ctx, translation)
			atomic.AddInt64(&processedCount, 1)
			return nil
		})
	}()

	// Give consumer time to start
	time.Sleep(100 * time.Millisecond)

	// Publish jobs concurrently
	numJobs := 50
	numWorkers := 5
	var wg sync.WaitGroup

	publishStart := time.Now()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numJobs/numWorkers; j++ {
				job := queue.TranslationJob{
					ID:         uuid.New().String(),
					EntityType: "project",
					EntityID:   uuid.New().String(),
					Language:   "es",
					Fields:     []string{"name"},
					SourceText: map[string]string{
						"name": fmt.Sprintf("Load Test Project %d", j),
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
		}(i)
	}

	wg.Wait()
	publishDuration := time.Since(publishStart)

	// Wait for processing
	time.Sleep(10 * time.Second)
	totalDuration := time.Since(startTime)

	processed := atomic.LoadInt64(&processedCount)
	errors := atomic.LoadInt64(&errorCount)

	t.Logf("Load Test Results:")
	t.Logf("  Jobs published: %d", numJobs)
	t.Logf("  Jobs processed: %d", processed)
	t.Logf("  Errors: %d", errors)
	t.Logf("  Publish duration: %v", publishDuration)
	t.Logf("  Total duration: %v", totalDuration)
	t.Logf("  Throughput: %.2f jobs/s", float64(processed)/totalDuration.Seconds())

	assert.Equal(t, int64(numJobs), processed, "All jobs should be processed")
	assert.Equal(t, int64(0), errors, "No errors should occur")
}

// TestTranslationWorkerMultiLanguageLoad tests processing multiple languages concurrently
func TestTranslationWorkerMultiLanguageLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping multi-language load test in short mode")
	}

	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	dbRepo := setupDatabase(t)

	queueName := fmt.Sprintf("multilang.translations.queue.%d", time.Now().Unix())
	exchange := fmt.Sprintf("multilang.woragis.translations.%d", time.Now().Unix())
	routingKey := "multilang.translations.process"

	// Create queue
	translationQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, 1, nil)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Setup mock translator
	mockTranslator := &mockTranslator{
		translations: make(map[string]string),
		shouldFail:   false,
	}

	// Start consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var processedCount int64
	languages := []string{"es", "fr", "de", "it", "pt"}

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
			dbRepo.UpdateTranslation(ctx, translation)
			atomic.AddInt64(&processedCount, 1)
			return nil
		})
	}()

	// Give consumer time to start
	time.Sleep(100 * time.Millisecond)

	// Publish jobs for multiple languages
	numJobsPerLanguage := 20
	entityID := uuid.New()
	startTime := time.Now()

	for _, lang := range languages {
		for i := 0; i < numJobsPerLanguage; i++ {
			job := queue.TranslationJob{
				ID:         uuid.New().String(),
				EntityType: "project",
				EntityID:   entityID.String(),
				Language:   lang,
				Fields:     []string{"name", "description"},
				SourceText: map[string]string{
					"name":        "Multi-language Project",
					"description": "A project for multi-language testing",
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
	}

	// Wait for processing
	time.Sleep(15 * time.Second)
	cancel()

	processed := atomic.LoadInt64(&processedCount)
	totalDuration := time.Since(startTime)
	expectedJobs := len(languages) * numJobsPerLanguage

	t.Logf("Multi-Language Load Test Results:")
	t.Logf("  Languages: %d", len(languages))
	t.Logf("  Jobs per language: %d", numJobsPerLanguage)
	t.Logf("  Total jobs: %d", expectedJobs)
	t.Logf("  Jobs processed: %d", processed)
	t.Logf("  Total duration: %v", totalDuration)
	t.Logf("  Throughput: %.2f jobs/s", float64(processed)/totalDuration.Seconds())

	assert.Equal(t, int64(expectedJobs), processed, "All jobs should be processed")
}

// TestTranslationWorkerDatabaseLoad tests database performance under load
func TestTranslationWorkerDatabaseLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database load test in short mode")
	}

	dbRepo := setupDatabase(t)

	// Test creating and updating many translations
	numTranslations := 100
	startTime := time.Now()

	for i := 0; i < numTranslations; i++ {
		translation := &database.Translation{
			ID:         uuid.New(),
			EntityType: "project",
			EntityID:   uuid.New(),
			Language:   "es",
			Status:     "processing",
		}
		fields := map[string]string{
			"name": fmt.Sprintf("DB Load Test Project %d", i),
		}
		translation.SetFields(fields)
		err := dbRepo.CreateTranslation(context.Background(), translation)
		require.NoError(t, err)

		// Update translation
		translation.Status = "completed"
		err = dbRepo.UpdateTranslation(context.Background(), translation)
		require.NoError(t, err)
	}

	duration := time.Since(startTime)

	t.Logf("Database Load Test Results:")
	t.Logf("  Translations created/updated: %d", numTranslations)
	t.Logf("  Duration: %v", duration)
	t.Logf("  Throughput: %.2f ops/s", float64(numTranslations*2)/duration.Seconds())

	assert.Less(t, duration.Seconds(), 5.0, "Database operations should complete within 5 seconds")
}
