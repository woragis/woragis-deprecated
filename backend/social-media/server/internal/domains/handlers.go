package socialmedia

import (
	"strconv"

	"woragis-social-media-service/pkg/middleware"
	"woragis-social-media-service/pkg/utils"
	apptracing "woragis-social-media-service/pkg/tracing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

// Handler handles HTTP requests for auth domain
type Handler struct {
	service Service
}

// NewHandler creates a new auth handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Add custom span attributes for business context
	ctx := c.UserContext()
	apptracing.SetSpanAttributes(ctx,
		attribute.String("auth.operation", "register"),
		attribute.String("auth.email", req.Email),
		attribute.String("auth.username", req.Username),
	)

	response, err := h.service.register(&req)
	if err != nil {
		switch err {
		case ErrUserAlreadyExists:
			return utils.ConflictResponse(c, "User already exists")
		case ErrPasswordTooWeak:
			return utils.BadRequestResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "Failed to register user")
		}
	}

	// Add user ID to span after successful registration
	if response != nil && response.User != nil {
		apptracing.SetSpanAttributes(ctx,
			attribute.String("auth.user_id", response.User.ID.String()),
		)
	}

	return utils.CreatedResponse(c, "User registered successfully", response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Get client info
	userAgent := c.Get("User-Agent", "")
	ipAddress := c.IP()

	// Add custom span attributes for business context
	ctx := c.UserContext()
	apptracing.SetSpanAttributes(ctx,
		attribute.String("auth.operation", "login"),
		attribute.String("auth.email", req.Email),
		attribute.String("auth.username", req.Username),
	)

	response, err := h.service.login(&req, userAgent, ipAddress)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			return utils.UnauthorizedResponse(c, "Invalid credentials")
		case ErrUserInactive:
			return utils.UnauthorizedResponse(c, "Account is inactive")
		default:
			return utils.InternalServerErrorResponse(c, "Failed to login")
		}
	}

	// Add user ID to span after successful login
	if response != nil && response.User != nil {
		apptracing.SetSpanAttributes(ctx,
			attribute.String("auth.user_id", response.User.ID.String()),
		)
	}

	return utils.SuccessResponse(c, "Login successful", response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Add custom span attributes for business context
	ctx := c.UserContext()
	apptracing.SetSpanAttributes(ctx,
		attribute.String("auth.operation", "refresh_token"),
	)

	response, err := h.service.refreshToken(req.RefreshToken)
	if err != nil {
		switch err {
		case ErrSessionExpired:
			return utils.UnauthorizedResponse(c, "Session expired")
		default:
			return utils.InternalServerErrorResponse(c, "Failed to refresh token")
		}
	}

	// Add user ID to span after successful refresh
	if response != nil && response.User != nil {
		apptracing.SetSpanAttributes(ctx,
			attribute.String("auth.user_id", response.User.ID.String()),
		)
	}

	return utils.SuccessResponse(c, "Token refreshed successfully", response)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate session
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	if err := h.service.logout(req.RefreshToken); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to logout")
	}

	return utils.SuccessResponse(c, "Logout successful", nil)
}

// LogoutAll godoc
// @Summary Logout from all devices
// @Description Logout user from all devices and invalidate all sessions
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/logout-all [post]
func (h *Handler) LogoutAll(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user context")
	}

	// Add custom span attributes for business context
	ctx := c.UserContext()
	apptracing.SetSpanAttributes(ctx,
		attribute.String("auth.operation", "logout_all"),
		attribute.String("auth.user_id", userID.String()),
	)

	if err := h.service.logoutAll(userID); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to logout from all devices")
	}

	return utils.SuccessResponse(c, "Logged out from all devices", nil)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Profile
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/profile [get]
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user context")
	}

	// Add custom span attributes for business context
	ctx := c.UserContext()
	apptracing.SetSpanAttributes(ctx,
		attribute.String("auth.operation", "get_profile"),
		attribute.String("auth.user_id", userID.String()),
	)

	profile, err := h.service.getUserProfile(userID)
	if err != nil {
		switch err {
		case ErrProfileNotFound:
			return utils.NotFoundResponse(c, "Profile not found")
		default:
			return utils.InternalServerErrorResponse(c, "Failed to get profile")
		}
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ProfileUpdateRequest true "Profile update data"
// @Success 200 {object} Profile
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/profile [put]
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user context")
	}

	var req ProfileUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Add custom span attributes for business context
	ctx := c.UserContext()
	apptracing.SetSpanAttributes(ctx,
		attribute.String("auth.operation", "update_profile"),
		attribute.String("auth.user_id", userID.String()),
	)

	profile, err := h.service.updateUserProfile(userID, &req)
	if err != nil {
		switch err {
		case ErrProfileNotFound:
			return utils.NotFoundResponse(c, "Profile not found")
		default:
			return utils.InternalServerErrorResponse(c, "Failed to update profile")
		}
	}

	return utils.SuccessResponse(c, "Profile updated successfully", profile)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PasswordChangeRequest true "Password change data"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/change-password [post]
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid user context")
	}

	var req PasswordChangeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err)
	}

	// Add custom span attributes for business context
	ctx := c.UserContext()
	apptracing.SetSpanAttributes(ctx,
		attribute.String("auth.operation", "change_password"),
		attribute.String("auth.user_id", userID.String()),
	)

	if err := h.service.changePassword(userID, &req); err != nil {
		switch err {
		case ErrInvalidPassword:
			return utils.BadRequestResponse(c, "Current password is incorrect")
		case ErrPasswordTooWeak:
			return utils.BadRequestResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "Failed to change password")
		}
	}

	return utils.SuccessResponse(c, "Password changed successfully", nil)
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify user email using verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/verify-email [get]
func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return utils.BadRequestResponse(c, "Verification token is required")
	}

	// Add custom span attributes for business context
	ctx := c.UserContext()
	apptracing.SetSpanAttributes(ctx,
		attribute.String("auth.operation", "verify_email"),
	)

	if err := h.service.verifyEmail(token); err != nil {
		switch err {
		case ErrTokenExpired, ErrTokenAlreadyUsed:
			return utils.BadRequestResponse(c, err.Error())
		default:
			return utils.InternalServerErrorResponse(c, "Failed to verify email")
		}
	}

	return utils.SuccessResponse(c, "Email verified successfully", nil)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user information by ID (admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} User
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/users/{id} [get]
func (h *Handler) GetUser(c *fiber.Ctx) error {
	// Check if user is admin
	userRole, err := middleware.GetUserRoleFromFiberContext(c)
	if err != nil || userRole != "admin" {
		return utils.ForbiddenResponse(c, "Admin access required")
	}

	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID")
	}

	user, err := h.service.getUserByID(userID)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			return utils.NotFoundResponse(c, "User not found")
		default:
			return utils.InternalServerErrorResponse(c, "Failed to get user")
		}
	}

	return utils.SuccessResponse(c, "User retrieved successfully", user)
}

// ListUsers godoc
// @Summary List users
// @Description List all users with pagination (admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/users [get]
func (h *Handler) ListUsers(c *fiber.Ctx) error {
	// Check if user is admin
	userRole, err := middleware.GetUserRoleFromFiberContext(c)
	if err != nil || userRole != "admin" {
		return utils.ForbiddenResponse(c, "Admin access required")
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.service.listUsers(page, limit, search)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to list users")
	}

	response := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	return utils.SuccessResponse(c, "Users retrieved successfully", response)
}

// CleanupExpiredSessions godoc
// @Summary Cleanup expired sessions
// @Description Remove expired sessions and tokens (admin only)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/cleanup [post]
func (h *Handler) CleanupExpiredSessions(c *fiber.Ctx) error {
	// Check if user is admin
	userRole, err := middleware.GetUserRoleFromFiberContext(c)
	if err != nil || userRole != "admin" {
		return utils.ForbiddenResponse(c, "Admin access required")
	}

	if err := h.service.cleanupExpiredSessions(); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to cleanup expired sessions")
	}

	return utils.SuccessResponse(c, "Expired sessions cleaned up successfully", nil)
}
