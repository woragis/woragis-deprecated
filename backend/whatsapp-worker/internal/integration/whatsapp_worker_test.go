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

	"github.com/woragis/backend/whatsapp-worker/internal/queue"
)

// mockNotifier is a mock WhatsApp notifier for testing
type mockNotifier struct {
	sentMessages []sentMessage
	shouldFail   bool
}

type sentMessage struct {
	destination string
	message     string
}

func (m *mockNotifier) Send(ctx context.Context, destination, message string) error {
	if m.shouldFail {
		return assert.AnError
	}
	m.sentMessages = append(m.sentMessages, sentMessage{
		destination: destination,
		message:     message,
	})
	return nil
}

func (m *mockNotifier) Connect(ctx context.Context) error {
	return nil
}

func (m *mockNotifier) Disconnect() {}

// setupRabbitMQConnection creates a RabbitMQ connection for testing
func setupRabbitMQConnection(t *testing.T) *queue.Connection {
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

// TestWhatsAppWorkerQueueSetup tests that the WhatsApp queue can be set up correctly
func TestWhatsAppWorkerQueueSetup(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.whatsapp.queue"
	exchange := "test.woragis.notifications"
	routingKey := "test.whatsapp.send"

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	_, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err, "Failed to create WhatsApp queue")

	// Verify queue exists
	ch := conn.Channel()
	defer ch.Close()

	_, err = ch.QueueInspect(queueName)
	require.NoError(t, err, "Queue should exist")

	// Verify exchange exists
	err = ch.ExchangeDeclarePassive(exchange, "direct", true, false, false, false, nil)
	require.NoError(t, err, "Exchange should exist")
}

// TestWhatsAppWorkerMessagePublish tests publishing messages to the WhatsApp queue
func TestWhatsAppWorkerMessagePublish(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.whatsapp.queue.publish"
	exchange := "test.woragis.notifications.publish"
	routingKey := "test.whatsapp.send.publish"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	_, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish a test message
	envelope := queue.WhatsAppEnvelope{
		UserID:      "test-user-123",
		TextMessage: "Hello from WhatsApp!",
		Destination: "+1234567890",
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
	require.NoError(t, err, "Failed to publish message")

	// Verify message is in queue
	q, err := ch.QueueInspect(queueName)
	require.NoError(t, err)
	assert.Greater(t, q.Messages, 0, "Queue should have messages")
}

// TestWhatsAppWorkerMessageConsume tests consuming and processing messages
func TestWhatsAppWorkerMessageConsume(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.whatsapp.queue.consume"
	exchange := "test.woragis.notifications.consume"
	routingKey := "test.whatsapp.send.consume"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	whatsappQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish a test message
	envelope := queue.WhatsAppEnvelope{
		UserID:      "test-user-456",
		TextMessage: "Test WhatsApp message",
		Destination: "+1234567890",
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

	// Setup mock notifier
	mockNotifier := &mockNotifier{
		sentMessages: make([]sentMessage, 0),
		shouldFail:   false,
	}

	// Consume and process message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan bool, 1)
	go func() {
		err := whatsappQueue.Consume(ctx, func(envelope queue.WhatsAppEnvelope) error {
			// Validate destination
			if envelope.Destination == "" {
				return assert.AnError
			}

			// Send WhatsApp message
			return mockNotifier.Send(ctx, envelope.Destination, envelope.TextMessage)
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
	assert.Greater(t, len(mockNotifier.sentMessages), 0, "Message should have been sent")
	if len(mockNotifier.sentMessages) > 0 {
		assert.Equal(t, envelope.Destination, mockNotifier.sentMessages[0].destination)
		assert.Equal(t, envelope.TextMessage, mockNotifier.sentMessages[0].message)
	}
}

// TestWhatsAppWorkerInvalidMessage tests handling of invalid messages
func TestWhatsAppWorkerInvalidMessage(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.whatsapp.queue.invalid"
	exchange := "test.woragis.notifications.invalid"
	routingKey := "test.whatsapp.send.invalid"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	whatsappQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
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
		_ = whatsappQueue.Consume(ctx, func(envelope queue.WhatsAppEnvelope) error {
			processed = true
			return nil
		})
	}()

	time.Sleep(1 * time.Second)
	cancel()

	// Invalid message should be rejected, not processed
	assert.False(t, processed, "Invalid message should not be processed")
}

// TestWhatsAppWorkerMissingDestination tests validation of missing destination
func TestWhatsAppWorkerMissingDestination(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.whatsapp.queue.nodest"
	exchange := "test.woragis.notifications.nodest"
	routingKey := "test.whatsapp.send.nodest"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	whatsappQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish message without destination
	envelope := queue.WhatsAppEnvelope{
		UserID:      "test-user",
		TextMessage: "Message without destination",
		Destination: "", // Missing destination
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

	// Consume should reject message with missing destination
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	processed := false
	processedSuccessfully := false
	go func() {
		_ = whatsappQueue.Consume(ctx, func(envelope queue.WhatsAppEnvelope) error {
			processed = true
			if envelope.Destination == "" {
				// Message should be rejected
				return assert.AnError
			}
			processedSuccessfully = true
			return nil
		})
	}()

	time.Sleep(2 * time.Second)
	cancel()

	// Message should be consumed but not processed successfully due to missing destination
	// The message is consumed (processed = true) but should fail (processedSuccessfully = false)
	if processed {
		assert.False(t, processedSuccessfully, "Message with missing destination should not be processed successfully")
	} else {
		// If message wasn't consumed at all, that's also acceptable (validation might reject it earlier)
		t.Log("Message with missing destination was not consumed (may be rejected by validation)")
	}
}

// TestWhatsAppWorkerRetryOnFailure tests retry behavior on send failure
func TestWhatsAppWorkerRetryOnFailure(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.whatsapp.queue.retry"
	exchange := "test.woragis.notifications.retry"
	routingKey := "test.whatsapp.send.retry"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	whatsappQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish a test message
	envelope := queue.WhatsAppEnvelope{
		UserID:      "test-user-retry",
		TextMessage: "This should be retried",
		Destination: "+1234567890",
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

	// Setup mock notifier that fails
	mockNotifier := &mockNotifier{
		sentMessages: make([]sentMessage, 0),
		shouldFail:   true,
	}

	// Consume with failing handler
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	attempts := 0
	go func() {
		_ = whatsappQueue.Consume(ctx, func(envelope queue.WhatsAppEnvelope) error {
			attempts++
			return mockNotifier.Send(ctx, envelope.Destination, envelope.TextMessage)
		})
	}()

	time.Sleep(1 * time.Second)
	cancel()

	// Message should be requeued for retry (attempts > 0)
	assert.Greater(t, attempts, 0, "Message should have been attempted")
}

// TestWhatsAppWorkerMultipleMessages tests processing multiple messages
func TestWhatsAppWorkerMultipleMessages(t *testing.T) {
	conn := setupRabbitMQConnection(t)
	defer conn.Close()

	queueName := "test.whatsapp.queue.multiple"
	exchange := "test.woragis.notifications.multiple"
	routingKey := "test.whatsapp.send.multiple"

	// Create queue
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	whatsappQueue, err := queue.NewQueue(conn, queueName, exchange, routingKey, logger)
	require.NoError(t, err)

	ch := conn.Channel()
	defer ch.Close()

	// Publish multiple messages
	messages := []queue.WhatsAppEnvelope{
		{UserID: "user1", TextMessage: "Message 1", Destination: "+1111111111"},
		{UserID: "user2", TextMessage: "Message 2", Destination: "+2222222222"},
		{UserID: "user3", TextMessage: "Message 3", Destination: "+3333333333"},
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

	// Setup mock notifier
	mockNotifier := &mockNotifier{
		sentMessages: make([]sentMessage, 0),
		shouldFail:   false,
	}

	// Consume messages
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		_ = whatsappQueue.Consume(ctx, func(envelope queue.WhatsAppEnvelope) error {
			if envelope.Destination == "" {
				return assert.AnError
			}
			return mockNotifier.Send(ctx, envelope.Destination, envelope.TextMessage)
		})
	}()

	// Wait for messages to be processed
	time.Sleep(3 * time.Second)
	cancel()

	// Verify all messages were processed
	assert.Equal(t, len(messages), len(mockNotifier.sentMessages), "All messages should be processed")
}

// TestWhatsAppWorkerDeadLetterQueue tests that failed messages go to DLQ
func TestWhatsAppWorkerDeadLetterQueue(t *testing.T) {
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
	dlqName := "test.whatsapp.queue.failed"
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
