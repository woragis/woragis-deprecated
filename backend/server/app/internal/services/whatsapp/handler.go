package whatsapp

import (
	"encoding/base64"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes WhatsApp-related endpoints.
type Handler struct {
	notifier     *WhatsmeowNotifier
	service      Service
	aiClient     *langchainservice.Client
	logger       *slog.Logger
	defaultModel string
}

// NewHandler builds a Handler.
func NewHandler(notifier *WhatsmeowNotifier, service Service, aiClient *langchainservice.Client, defaultModel string, logger *slog.Logger) *Handler {
	return &Handler{
		notifier:     notifier,
		service:      service,
		aiClient:     aiClient,
		logger:       logger,
		defaultModel: defaultModel,
	}
}

// GetQRCode returns the WhatsApp QR code as a base64-encoded PNG image.
func (h *Handler) GetQRCode(c *fiber.Ctx) error {
	if h.notifier == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "WhatsApp service not available",
		})
	}

	qrCode := h.notifier.GetQRCode()
	if qrCode == "" {
		// Check if already connected
		if h.notifier.IsConnected() {
			return response.Success(c, fiber.StatusOK, fiber.Map{
				"connected": true,
				"qr_code":   nil,
				"message":   "WhatsApp is already connected",
			})
		}

		return response.Error(c, fiber.StatusNotFound, 404, fiber.Map{
			"message": "QR code has not been generated yet or has expired. Please ensure WhatsApp is enabled and the service is running.",
		})
	}

	// Generate QR code image
	png, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("whatsapp: failed to generate QR code image", slog.Any("error", err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate QR code image",
		})
	}

	// Convert to base64
	base64Image := base64.StdEncoding.EncodeToString(png)

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"connected": false,
		"qr_code":   "data:image/png;base64," + base64Image,
		"qr_text":   qrCode,
		"message":   "Scan this QR code with WhatsApp on your phone",
	})
}

// GetStatus returns the WhatsApp connection status.
func (h *Handler) GetStatus(c *fiber.Ctx) error {
	if h.notifier == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"connected": false,
			"error":     "WhatsApp service not available",
		})
	}

	connected := h.notifier.IsConnected()
	qrCode := h.notifier.GetQRCode()

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"connected":   connected,
		"has_qr_code": qrCode != "",
		"message": func() string {
			if connected {
				return "WhatsApp is connected"
			}
			if qrCode != "" {
				return "QR code available - scan to connect"
			}
			return "Not connected - waiting for QR code generation"
		}(),
	})
}

type sendMessagePayload struct {
	ClientID      string  `json:"client_id"`
	Message       string  `json:"message,omitempty"`
	UseAI         bool    `json:"use_ai,omitempty"`
	Template      string  `json:"template,omitempty"`
	Instructions  string  `json:"instructions,omitempty"`
	ClientContext string  `json:"client_context,omitempty"` // Additional context about the client
}

// SendMessage sends a WhatsApp message to a client.
func (h *Handler) SendMessage(c *fiber.Ctx) error {
	if h.service == nil {
		return h.handleError(c, NewDomainError(ErrCodeServiceNotConfigured, ErrServiceNotConfigured))
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "Authentication required",
		})
	}

	var payload sendMessagePayload
	if err := c.BodyParser(&payload); err != nil {
		return h.handleError(c, NewDomainError(ErrCodeInvalidPayload, "whatsapp: invalid request payload"))
	}

	if payload.ClientID == "" {
		return h.handleError(c, NewDomainError(ErrCodeInvalidPayload, ErrEmptyClientID))
	}

	clientUUID, err := uuid.Parse(payload.ClientID)
	if err != nil {
		return h.handleError(c, NewDomainError(ErrCodeInvalidPayload, ErrInvalidClientID))
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return h.handleError(c, NewDomainError(ErrCodeInvalidPayload, "whatsapp: invalid user_id format"))
	}

	message := strings.TrimSpace(payload.Message)

	// If AI is requested, generate the message
	if payload.UseAI && (payload.Template != "" || payload.Instructions != "") {
		generated, err := h.generateAIMessage(payload, c)
		if err != nil {
			if h.logger != nil {
				h.logger.Error("failed to generate AI message", slog.Any("error", err))
			}
			return h.handleError(c, NewDomainError(ErrCodeAIGenerationFailure, ErrAIGenerationFailure))
		}
		message = generated
	}

	if message == "" {
		return h.handleError(c, NewDomainError(ErrCodeInvalidPayload, ErrEmptyMessage))
	}

	if err := h.service.SendToClient(c.Context(), userUUID, clientUUID, message); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message":      "Message sent successfully",
		"sent_message": message,
	})
}

// getCurrentUserID extracts the user ID from the Fiber context (set by auth middleware).
func getCurrentUserID(c *fiber.Ctx) (string, error) {
	// The auth middleware stores user_id using ContextUserIDKey
	// We'll check both possible keys for compatibility
	userID := c.Locals("user_id")
	if userID == nil {
		// Try the auth domain's context key
		userID = c.Locals(authdomain.ContextUserIDKey)
	}
	
	if userID == nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusInternalServerError, "invalid user_id type")
	}

	return userIDStr, nil
}

// handleError maps domain errors to HTTP responses.
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromErrorCode(domainErr.Code)
		return response.Error(c, status, domainErr.Code, fiber.Map{
			"message": domainErr.Message,
		})
	}

	// Fallback for non-domain errors
	if h.logger != nil {
		h.logger.Error("unexpected error in handler", slog.Any("error", err))
	}
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeSendFailure, fiber.Map{
		"message": "An unexpected error occurred",
	})
}

// statusFromErrorCode maps error codes to HTTP status codes.
func statusFromErrorCode(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidPhoneNumber:
		return fiber.StatusBadRequest
	case ErrCodeServiceNotConfigured, ErrCodeNotConnected:
		return fiber.StatusServiceUnavailable
	case ErrCodeClientNotFound, ErrCodeUserNotFound:
		return fiber.StatusNotFound
	case ErrCodeNoPhoneNumber:
		return fiber.StatusBadRequest
	case ErrCodeSendFailure, ErrCodeAIGenerationFailure, ErrCodeRepositoryFailure:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusInternalServerError
	}
}

// generateAIMessage generates a WhatsApp message using AI based on template or instructions.
func (h *Handler) generateAIMessage(payload sendMessagePayload, c *fiber.Ctx) (string, error) {
	if h.aiClient == nil {
		return "", NewDomainError(ErrCodeAIGenerationFailure, ErrAIClientNotConfigured)
	}

	ctx := c.Context()
	model := h.defaultModel
	if model == "" {
		model = "chatgpt"
	}

	// Build the system prompt
	systemPrompt := "You are a professional business communication assistant. Generate concise, friendly WhatsApp messages for client communication. Keep messages brief, professional, and appropriate for WhatsApp (under 200 words)."

	// Build user prompt based on template or instructions
	var userPrompt strings.Builder

	if payload.Template != "" {
		userPrompt.WriteString("Use this template as a guide:\n")
		userPrompt.WriteString(payload.Template)
		userPrompt.WriteString("\n\n")
	}

	if payload.Instructions != "" {
		userPrompt.WriteString("Follow these instructions:\n")
		userPrompt.WriteString(payload.Instructions)
		userPrompt.WriteString("\n\n")
	}

	if payload.ClientContext != "" {
		userPrompt.WriteString("Client context:\n")
		userPrompt.WriteString(payload.ClientContext)
		userPrompt.WriteString("\n\n")
	}

	userPrompt.WriteString("Generate a professional WhatsApp message based on the above information. Only return the message text, nothing else.")

	messages := []langchainservice.ChatMessage{
		{
			Role:      "system",
			Content:   systemPrompt,
			Timestamp: time.Now().UTC(),
		},
		{
			Role:      "user",
			Content:   userPrompt.String(),
			Timestamp: time.Now().UTC(),
		},
	}

	req := langchainservice.ChatCompletionRequest{
		Provider:    langchainservice.ProviderOpenAI,
		Model:       model,
		Messages:    messages,
		MaxTokens:   200,
		Temperature: 0.7,
	}

	resp, err := h.aiClient.GenerateCompletion(ctx, req)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("AI generation failed", slog.Any("error", err))
		}
		return "", NewDomainError(ErrCodeAIGenerationFailure, ErrAIGenerationFailure)
	}

	generated := strings.TrimSpace(resp.Message.Content)
	if generated == "" {
		return "", NewDomainError(ErrCodeAIGenerationFailure, ErrAIEmptyResponse)
	}

	return generated, nil
}
