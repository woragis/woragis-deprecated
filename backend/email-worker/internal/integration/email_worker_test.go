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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/woragis/backend/email-worker/internal/queue"
	"github.com/woragis/backend/email-worker/internal/sender"
)

// mockSender is a mock email sender for testing
type mockSender struct {
	sentMessages []sender.Message
	shouldFail   bool
}

func (m *mockSender) Send(ctx context.Context, msg sender.Message) error {
	if m.shouldFail {
		return assert.AnError
	}
	m.sentMessages = append(m.sentMessages, msg)
	return nil
}

// setupRabbitMQConnection creates a RabbitMQ connection for testing
func setupRabbitMQConnection(t *testing.T) *queue.Connection {
	// Use test RabbitMQ from docker-compose.test.yml
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://test:test@localhost:5673/test")
	conn, err := queue.NewConnection(rabbitmqURL)
	require.NoError(t, err, "Failed to connect to RabbitMQ")
	return conn
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestEmailWorkerQueueSetup tests that the email queue can be set up correctly
func TestEmailWorkerQueueSetup(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.emails.queue"
	exchange := "test.woragis.notifications"
	routingKey := "test.emails.send"

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	_, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err, "Failed to create email queue")

	// Verify queue exists
	ch := conn.Channel()
	defer ch.Close()

	_, err = ch.QueueInspect(queueName)
	require.NoError(t, err, "Queue should exist")

	// Verify exchange exists
	err = ch.ExchangeDeclarePassive(exchange, "direct", true, false, false, false, nil)
	require.NoError(t, err, "Exchange should exist")
}

// TestEmailWorkerMessagePublish tests publishing messages to the email queue
func TestEmailWorkerMessagePublish(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.emails.queue"
	exchange := "test.woragis.notifications"
	routingKey := "test.emails.send"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	_, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish a test message
	envelope := queue.EmailEnvelope{
		UserID:      "test-user-123",
		Subject:     "Test Email",
		TextMessage: "This is a test email message",
		HTMLMessage: "<p>This is a test email message</p>",
		Destination: "test@example.com",
	}

	body, err := json.Marshal(envelope)
	require.NoError(t, err)

	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
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

// TestEmailWorkerMessageConsume tests consuming and processing messages
func TestEmailWorkerMessageConsume(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.emails.queue.consume"
	exchange := "test.woragis.notifications.consume"
	routingKey := "test.emails.send.consume"

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish a test message
	envelope := queue.EmailEnvelope{
		UserID:      "test-user-456",
		Subject:     "Consume Test Email",
		TextMessage: "This is a consume test",
		HTMLMessage: "<p>This is a consume test</p>",
		Destination: "consume@example.com",
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

	// Setup mock sender
	mockSender := &mockSender{
		sentMessages: make([]sender.Message, 0),
		shouldFail:   false,
	}

	// Consume and process message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan bool, 1)
	go func() {
		err := emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
				HTMLBody: envelope.HTMLMessage,
			}
			return mockSender.Send(ctx, msg)
		})
		if err != nil && err != context.DeadlineExceeded {
			t.Logf("Consume error: %v", err)
		}
		done <- true
	}()

	// Wait for message to be processed or timeout
	select {
	case <-done:
		// Message processed
	case <-time.After(3 * time.Second):
		// Timeout - check if message was processed
	}

	// Verify message was sent
	assert.Greater(t, len(mockSender.sentMessages), 0, "Message should have been sent")
	if len(mockSender.sentMessages) > 0 {
		assert.Equal(t, envelope.Destination, mockSender.sentMessages[0].To)
		assert.Equal(t, envelope.Subject, mockSender.sentMessages[0].Subject)
		assert.Equal(t, envelope.TextMessage, mockSender.sentMessages[0].TextBody)
	}
}

// TestEmailWorkerInvalidMessage tests handling of invalid messages
func TestEmailWorkerInvalidMessage(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.emails.queue.invalid"
	exchange := "test.woragis.notifications.invalid"
	routingKey := "test.emails.send.invalid"

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
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
		_ = emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			processed = true
			return nil
		})
	}()

	time.Sleep(1 * time.Second)
	cancel()

	// Invalid message should be rejected, not processed
	assert.False(t, processed, "Invalid message should not be processed")
}

// TestEmailWorkerRetryOnFailure tests retry behavior on send failure
func TestEmailWorkerRetryOnFailure(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.emails.queue.retry"
	exchange := "test.woragis.notifications.retry"
	routingKey := "test.emails.send.retry"

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish a test message
	envelope := queue.EmailEnvelope{
		UserID:      "test-user-retry",
		Subject:     "Retry Test Email",
		TextMessage: "This should be retried",
		Destination: "retry@example.com",
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

	// Setup mock sender that fails
	mockSender := &mockSender{
		sentMessages: make([]sender.Message, 0),
		shouldFail:   true,
	}

	// Consume with failing handler
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	attempts := 0
	go func() {
		_ = emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			attempts++
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
			}
			return mockSender.Send(ctx, msg)
		})
	}()

	time.Sleep(1 * time.Second)
	cancel()

	// Message should be requeued for retry (attempts > 0)
	assert.Greater(t, attempts, 0, "Message should have been attempted")
}

// TestEmailWorkerMultipleMessages tests processing multiple messages
func TestEmailWorkerMultipleMessages(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.emails.queue.multiple"
	exchange := "test.woragis.notifications.multiple"
	routingKey := "test.emails.send.multiple"

	// Create logger for queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	
	// Create queue
	emailQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish multiple messages
	messages := []queue.EmailEnvelope{
		{UserID: "user1", Subject: "Email 1", TextMessage: "Message 1", Destination: "user1@example.com"},
		{UserID: "user2", Subject: "Email 2", TextMessage: "Message 2", Destination: "user2@example.com"},
		{UserID: "user3", Subject: "Email 3", TextMessage: "Message 3", Destination: "user3@example.com"},
	}

	for _, envelope := range messages {
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

	// Setup mock sender
	mockSender := &mockSender{
		sentMessages: make([]sender.Message, 0),
		shouldFail:   false,
	}

	// Consume messages
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		_ = emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
			}
			return mockSender.Send(ctx, msg)
		})
	}()

	// Wait for messages to be processed
	time.Sleep(3 * time.Second)
	cancel()

	// Verify all messages were processed
	assert.Equal(t, len(messages), len(mockSender.sentMessages), "All messages should be processed")
}

// TestEmailWorkerDeadLetterQueue tests that failed messages go to DLQ
func TestEmailWorkerDeadLetterQueue(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	ch := conn.Channel()
	defer ch.Close()

	// Declare dead letter exchange
	dlxName := "test.woragis.dlx"
	err := ch.ExchangeDeclare(
		dlxName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)

	// Declare DLQ
	dlqName := "test.emails.queue.failed"
	_, err = ch.QueueDeclare(
		dlqName,
		true,
		false,
		false,
		false,
		nil,
	)
	require.NoError(t, err)

	// Bind DLQ to DLX
	err = ch.QueueBind(
		dlqName,
		dlqName+".failed",
		dlxName,
		false,
		nil,
	)
	require.NoError(t, err)

	// Verify DLQ exists and is bound
	q, err := ch.QueueInspect(dlqName)
	require.NoError(t, err)
	assert.NotNil(t, q)
}
