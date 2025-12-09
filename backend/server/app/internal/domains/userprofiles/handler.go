package userprofiles

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes user profile endpoints.
type Handler interface {
	GetProfile(c *fiber.Ctx) error
	UpsertProfile(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

// NewHandler constructs a user profile handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// GetProfile retrieves the current user's profile.
func (h *handler) GetProfile(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	profile, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				// Return empty profile if not found
				return response.Success(c, fiber.StatusOK, &UserProfile{
					UserID:  userID,
					AboutMe: "",
				})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to get profile", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to get profile"})
	}

	return response.Success(c, fiber.StatusOK, profile)
}

// UpsertProfile creates or updates the current user's profile.
func (h *handler) UpsertProfile(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	var req struct {
		AboutMe string `json:"aboutMe"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid request body"})
	}

	profile, err := h.service.UpsertProfile(c.Context(), userID, req.AboutMe)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to upsert profile", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to save profile"})
	}

	return response.Success(c, fiber.StatusOK, profile)
}

