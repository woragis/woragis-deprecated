package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestUser holds authentication data extracted from the Fiber context.
type RequestUser struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

// ErrNoUserInContext indicates that the request is missing authentication data.
var ErrNoUserInContext = errors.New("auth: user not found in request context")

// UserFromContext extracts the authenticated user from the Fiber context.
func UserFromContext(c *fiber.Ctx) (RequestUser, error) {
	if c == nil {
		return RequestUser{}, ErrNoUserInContext
	}

	rawID := c.Locals(ContextUserIDKey)
	idStr, ok := rawID.(string)
	if !ok || idStr == "" {
		return RequestUser{}, ErrNoUserInContext
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		return RequestUser{}, err
	}

	email, _ := c.Locals(ContextUserEmailKey).(string)

	return RequestUser{
		ID:    userID,
		Email: email,
	}, nil
}

// UserIDFromContext extracts only the authenticated user ID from the Fiber context.
func UserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	user, err := UserFromContext(c)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

// UserEmailFromContext extracts the authenticated user e-mail from the Fiber context.
func UserEmailFromContext(c *fiber.Ctx) (string, error) {
	user, err := UserFromContext(c)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}
