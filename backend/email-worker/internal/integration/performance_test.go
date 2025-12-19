//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/woragis/backend/email-worker/internal/queue"
	"github.com/woragis/backend/email-worker/internal/sender"
)

// setupRabbitMQConnectionForBenchmark creates a RabbitMQ connection for benchmarks
func setupRabbitMQConnectionForBenchmark(b *testing.B) *queue.Connection {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://test:test@localhost:5673/test"
	}
	conn, err := queue.NewConnection(rabbitmqURL)
	if err != nil {
		b.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	return conn
}

// BenchmarkEmailWorkerThroughput benchmarks email worker throughput
func BenchmarkEmailWorkerThroughput(b *testing.B) {
	conn := setupRabbitMQConnectionForBenchmark(b)
	defer conn.Close()

	queueName := fmt.Sprintf("bench.emails.queue.%d", time.Now().Unix())
	exchange := fmt.Sprintf("bench.woragis.notifications.%d", time.Now().Unix())
	routingKey := "bench.emails.send"

	ch := conn.Channel()
	defer ch.Close()

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(b, err)

	// Setup mock sender
	mockSender := &mockSender{
		sentMessages: make([]sender.Message, 0),
		shouldFail:   false,
	}

	// Start consumer in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var processedCount int64
	go func() {
		_ = emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
			}
			err := mockSender.Send(ctx, msg)
			if err == nil {
				atomic.AddInt64(&processedCount, 1)
			}
			return err
		})
	}()

	// Give consumer time to start
	time.Sleep(100 * time.Millisecond)

	// Benchmark publishing messages
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			envelope := queue.EmailEnvelope{
				UserID:      "bench-user",
				Subject:     "Benchmark Email",
				TextMessage: "Benchmark test message",
				Destination: "bench@example.com",
			}

			body, _ := json.Marshal(envelope)
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

	// Wait for messages to be processed
	time.Sleep(2 * time.Second)
	b.StopTimer()

	b.Logf("Processed %d messages", atomic.LoadInt64(&processedCount))
}

// TestEmailWorkerLoadTest tests email worker under load
func TestEmailWorkerLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := fmt.Sprintf("load.emails.queue.%d", time.Now().Unix())
	exchange := fmt.Sprintf("load.woragis.notifications.%d", time.Now().Unix())
	routingKey := "load.emails.send"

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Setup mock sender
	mockSender := &mockSender{
		sentMessages: make([]sender.Message, 0),
		shouldFail:   false,
	}

	// Start consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var processedCount int64
	var errorCount int64
	startTime := time.Now()

	go func() {
		_ = emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
			}
			err := mockSender.Send(ctx, msg)
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				return err
			}
			atomic.AddInt64(&processedCount, 1)
			return nil
		})
	}()

	// Give consumer time to start
	time.Sleep(100 * time.Millisecond)

	// Publish messages concurrently
	numMessages := 100
	numWorkers := 10
	var wg sync.WaitGroup

	publishStart := time.Now()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numMessages/numWorkers; j++ {
				envelope := queue.EmailEnvelope{
					UserID:      fmt.Sprintf("user-%d-%d", workerID, j),
					Subject:     fmt.Sprintf("Load Test Email %d", j),
					TextMessage: "Load test message",
					Destination: fmt.Sprintf("load%d@example.com", j),
				}

				body, err := json.Marshal(envelope)
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

	// Wait for messages to be processed
	time.Sleep(5 * time.Second)
	totalDuration := time.Since(startTime)

	processed := atomic.LoadInt64(&processedCount)
	errors := atomic.LoadInt64(&errorCount)

	t.Logf("Load Test Results:")
	t.Logf("  Messages published: %d", numMessages)
	t.Logf("  Messages processed: %d", processed)
	t.Logf("  Errors: %d", errors)
	t.Logf("  Publish duration: %v", publishDuration)
	t.Logf("  Total duration: %v", totalDuration)
	t.Logf("  Throughput: %.2f msg/s", float64(processed)/totalDuration.Seconds())

	assert.Equal(t, int64(numMessages), processed, "All messages should be processed")
	assert.Equal(t, int64(0), errors, "No errors should occur")
	assert.Less(t, totalDuration.Seconds(), 10.0, "Processing should complete within 10 seconds")
}

// TestEmailWorkerConcurrentConsumers tests multiple concurrent consumers
func TestEmailWorkerConcurrentConsumers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent consumer test in short mode")
	}

	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := fmt.Sprintf("concurrent.emails.queue.%d", time.Now().Unix())
	exchange := fmt.Sprintf("concurrent.woragis.notifications.%d", time.Now().Unix())
	routingKey := "concurrent.emails.send"

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Setup mock sender
	mockSender := &mockSender{
		sentMessages: make([]sender.Message, 0),
		shouldFail:   false,
	}

	// Start multiple consumers
	numConsumers := 3
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var processedCount int64
	var wg sync.WaitGroup

	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			_ = emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
				msg := sender.Message{
					To:       envelope.Destination,
					Subject:  envelope.Subject,
					TextBody: envelope.TextMessage,
				}
				err := mockSender.Send(ctx, msg)
				if err == nil {
					atomic.AddInt64(&processedCount, 1)
				}
				return err
			})
		}(i)
	}

	// Give consumers time to start
	time.Sleep(200 * time.Millisecond)

	// Publish messages
	numMessages := 30
	for i := 0; i < numMessages; i++ {
		envelope := queue.EmailEnvelope{
			UserID:      fmt.Sprintf("user-%d", i),
			Subject:     "Concurrent Test Email",
			TextMessage: "Concurrent test message",
			Destination: fmt.Sprintf("concurrent%d@example.com", i),
		}

		body, err := json.Marshal(envelope)
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

	// Wait for processing
	time.Sleep(3 * time.Second)
	cancel()

	processed := atomic.LoadInt64(&processedCount)
	t.Logf("Concurrent Consumers Test:")
	t.Logf("  Consumers: %d", numConsumers)
	t.Logf("  Messages: %d", numMessages)
	t.Logf("  Processed: %d", processed)

	assert.Equal(t, int64(numMessages), processed, "All messages should be processed")
}

// TestEmailWorkerLatency tests message processing latency
func TestEmailWorkerLatency(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := fmt.Sprintf("latency.emails.queue.%d", time.Now().Unix())
	exchange := fmt.Sprintf("latency.woragis.notifications.%d", time.Now().Unix())
	routingKey := "latency.emails.send"

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Setup mock sender with minimal delay
	mockSender := &mockSender{
		sentMessages: make([]sender.Message, 0),
		shouldFail:   false,
	}

	// Track latencies
	var latencies []time.Duration
	var mu sync.Mutex

	// Start consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			start := time.Now()
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
			}
			err := mockSender.Send(ctx, msg)
			latency := time.Since(start)

			mu.Lock()
			latencies = append(latencies, latency)
			mu.Unlock()

			return err
		})
	}()

	// Give consumer time to start
	time.Sleep(100 * time.Millisecond)

	// Publish messages and measure latency
	numMessages := 50
	for i := 0; i < numMessages; i++ {
		publishStart := time.Now()
		envelope := queue.EmailEnvelope{
			UserID:      fmt.Sprintf("user-%d", i),
			Subject:     "Latency Test Email",
			TextMessage: "Latency test message",
			Destination: fmt.Sprintf("latency%d@example.com", i),
		}

		body, err := json.Marshal(envelope)
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
		_ = publishStart
	}

	// Wait for processing
	time.Sleep(3 * time.Second)
	cancel()

	mu.Lock()
	defer mu.Unlock()

	if len(latencies) > 0 {
		var total time.Duration
		var min, max time.Duration = latencies[0], latencies[0]
		for _, lat := range latencies {
			total += lat
			if lat < min {
				min = lat
			}
			if lat > max {
				max = lat
			}
		}
		avg := total / time.Duration(len(latencies))

		t.Logf("Latency Test Results:")
		t.Logf("  Messages: %d", len(latencies))
		t.Logf("  Min latency: %v", min)
		t.Logf("  Max latency: %v", max)
		t.Logf("  Avg latency: %v", avg)

		assert.Less(t, avg, 100*time.Millisecond, "Average latency should be less than 100ms")
	}
}

// TestEmailWorkerRateLimiting tests behavior under high message rate
func TestEmailWorkerRateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rate limiting test in short mode")
	}

	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := fmt.Sprintf("ratelimit.emails.queue.%d", time.Now().Unix())
	exchange := fmt.Sprintf("ratelimit.woragis.notifications.%d", time.Now().Unix())
	routingKey := "ratelimit.emails.send"

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Setup mock sender with delay to simulate rate limiting
	mockSender := &mockSender{
		sentMessages: make([]sender.Message, 0),
		shouldFail:   false,
	}

	// Start consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var processedCount int64
	go func() {
		_ = emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			// Simulate processing delay
			time.Sleep(10 * time.Millisecond)
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
			}
			err := mockSender.Send(ctx, msg)
			if err == nil {
				atomic.AddInt64(&processedCount, 1)
			}
			return err
		})
	}()

	// Give consumer time to start
	time.Sleep(100 * time.Millisecond)

	// Publish messages at high rate
	numMessages := 200
	startTime := time.Now()
	for i := 0; i < numMessages; i++ {
		envelope := queue.EmailEnvelope{
			UserID:      fmt.Sprintf("user-%d", i),
			Subject:     "Rate Limit Test Email",
			TextMessage: "Rate limit test message",
			Destination: fmt.Sprintf("ratelimit%d@example.com", i),
		}

		body, err := json.Marshal(envelope)
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
	publishDuration := time.Since(startTime)

	// Wait for processing
	time.Sleep(10 * time.Second)
	cancel()

	processed := atomic.LoadInt64(&processedCount)
	totalDuration := time.Since(startTime)

	t.Logf("Rate Limiting Test Results:")
	t.Logf("  Messages published: %d", numMessages)
	t.Logf("  Messages processed: %d", processed)
	t.Logf("  Publish rate: %.2f msg/s", float64(numMessages)/publishDuration.Seconds())
	t.Logf("  Processing rate: %.2f msg/s", float64(processed)/totalDuration.Seconds())

	assert.Equal(t, int64(numMessages), processed, "All messages should eventually be processed")
}
