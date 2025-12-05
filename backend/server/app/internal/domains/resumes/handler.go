package resumes

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes resume endpoints.
type Handler interface {
	CreateResume(c *fiber.Ctx) error
	UpdateResume(c *fiber.Ctx) error
	DeleteResume(c *fiber.Ctx) error
	GetResume(c *fiber.Ctx) error
	ListResumes(c *fiber.Ctx) error
	MarkAsMain(c *fiber.Ctx) error
	MarkAsFeatured(c *fiber.Ctx) error
	UnmarkAsMain(c *fiber.Ctx) error
	UnmarkAsFeatured(c *fiber.Ctx) error
	DownloadResume(c *fiber.Ctx) error
	PreviewResume(c *fiber.Ctx) error
}

type handler struct {
	service      Service
	logger       *slog.Logger
	baseFilePath string // Base path where resume files are stored
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a resume handler.
func NewHandler(service Service, baseFilePath string, logger *slog.Logger) Handler {
	return &handler{
		service:      service,
		logger:       logger,
		baseFilePath: baseFilePath,
	}
}

// CreateResume creates a new resume.
func (h *handler) CreateResume(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	var req struct {
		Title    string `json:"title"`
		FilePath string `json:"filePath"`
		FileName string `json:"fileName"`
		FileSize int64  `json:"fileSize"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	resume, err := h.service.CreateResume(c.Context(), userID, req.Title, req.FilePath, req.FileName, req.FileSize)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			return response.BadRequest(c, domainErr.Message)
		}
		h.logger.Error("failed to create resume", slog.Any("error", err))
		return response.InternalServerError(c, "failed to create resume")
	}

	return response.Created(c, resume)
}

// UpdateResume updates an existing resume.
func (h *handler) UpdateResume(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid resume ID")
	}

	var req struct {
		Title string `json:"title"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	resume, err := h.service.UpdateResume(c.Context(), userID, resumeID, req.Title)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
			return response.BadRequest(c, domainErr.Message)
		}
		h.logger.Error("failed to update resume", slog.Any("error", err))
		return response.InternalServerError(c, "failed to update resume")
	}

	return response.OK(c, resume)
}

// DeleteResume deletes a resume.
func (h *handler) DeleteResume(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid resume ID")
	}

	if err := h.service.DeleteResume(c.Context(), userID, resumeID); err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
			return response.BadRequest(c, domainErr.Message)
		}
		h.logger.Error("failed to delete resume", slog.Any("error", err))
		return response.InternalServerError(c, "failed to delete resume")
	}

	return response.NoContent(c)
}

// GetResume retrieves a resume by ID.
func (h *handler) GetResume(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid resume ID")
	}

	resume, err := h.service.GetResume(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
			return response.BadRequest(c, domainErr.Message)
		}
		h.logger.Error("failed to get resume", slog.Any("error", err))
		return response.InternalServerError(c, "failed to get resume")
	}

	return response.OK(c, resume)
}

// ListResumes lists all resumes for the authenticated user.
func (h *handler) ListResumes(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	resumes, err := h.service.ListResumes(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to list resumes", slog.Any("error", err))
		return response.InternalServerError(c, "failed to list resumes")
	}

	return response.OK(c, resumes)
}

// MarkAsMain marks a resume as main.
func (h *handler) MarkAsMain(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid resume ID")
	}

	resume, err := h.service.MarkAsMain(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
			return response.BadRequest(c, domainErr.Message)
		}
		h.logger.Error("failed to mark resume as main", slog.Any("error", err))
		return response.InternalServerError(c, "failed to mark resume as main")
	}

	return response.OK(c, resume)
}

// MarkAsFeatured marks a resume as featured.
func (h *handler) MarkAsFeatured(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid resume ID")
	}

	resume, err := h.service.MarkAsFeatured(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
			return response.BadRequest(c, domainErr.Message)
		}
		h.logger.Error("failed to mark resume as featured", slog.Any("error", err))
		return response.InternalServerError(c, "failed to mark resume as featured")
	}

	return response.OK(c, resume)
}

// UnmarkAsMain removes the main flag from a resume.
func (h *handler) UnmarkAsMain(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid resume ID")
	}

	resume, err := h.service.UnmarkAsMain(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
			return response.BadRequest(c, domainErr.Message)
		}
		h.logger.Error("failed to unmark resume as main", slog.Any("error", err))
		return response.InternalServerError(c, "failed to unmark resume as main")
	}

	return response.OK(c, resume)
}

// UnmarkAsFeatured removes the featured flag from a resume.
func (h *handler) UnmarkAsFeatured(c *fiber.Ctx) error {
	userID, err := authdomain.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid resume ID")
	}

	resume, err := h.service.UnmarkAsFeatured(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
			return response.BadRequest(c, domainErr.Message)
		}
		h.logger.Error("failed to unmark resume as featured", slog.Any("error", err))
		return response.InternalServerError(c, "failed to unmark resume as featured")
	}

	return response.OK(c, resume)
}

// DownloadResume downloads the best resume (public endpoint).
func (h *handler) DownloadResume(c *fiber.Ctx) error {
	// Get user ID from query param or use default user
	// For now, we'll get it from a query param or use a default
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		// Use a default user ID - you might want to configure this
		return response.BadRequest(c, "userId query parameter is required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid userId")
	}

	// Get the best resume (main > featured > most recent)
	resume, err := h.service.GetBestResume(c.Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
		}
		h.logger.Error("failed to get resume for download", slog.Any("error", err))
		return response.InternalServerError(c, "failed to get resume")
	}

	// Build full file path
	fullPath := filepath.Join(h.baseFilePath, resume.FilePath)
	if !filepath.IsAbs(resume.FilePath) {
		fullPath = filepath.Join(h.baseFilePath, resume.FilePath)
	} else {
		fullPath = resume.FilePath
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		h.logger.Error("resume file not found", slog.String("path", fullPath))
		return response.NotFound(c, ErrFileNotFound)
	}

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		h.logger.Error("failed to open resume file", slog.Any("error", err))
		return response.InternalServerError(c, ErrFileReadError)
	}
	defer file.Close()

	// Set headers for PDF download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", `attachment; filename="`+resume.FileName+`"`)

	// Stream file to response
	_, err = io.Copy(c.Response().BodyWriter(), file)
	if err != nil {
		h.logger.Error("failed to stream resume file", slog.Any("error", err))
		return response.InternalServerError(c, "failed to stream file")
	}

	return nil
}

// PreviewResume serves the resume for preview (public endpoint).
func (h *handler) PreviewResume(c *fiber.Ctx) error {
	// Get user ID from query param
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return response.BadRequest(c, "userId query parameter is required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid userId")
	}

	// Get the best resume
	resume, err := h.service.GetBestResume(c.Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.NotFound(c, domainErr.Message)
			}
		}
		h.logger.Error("failed to get resume for preview", slog.Any("error", err))
		return response.InternalServerError(c, "failed to get resume")
	}

	// Build full file path
	fullPath := filepath.Join(h.baseFilePath, resume.FilePath)
	if !filepath.IsAbs(resume.FilePath) {
		fullPath = filepath.Join(h.baseFilePath, resume.FilePath)
	} else {
		fullPath = resume.FilePath
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		h.logger.Error("resume file not found", slog.String("path", fullPath))
		return response.NotFound(c, ErrFileNotFound)
	}

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		h.logger.Error("failed to open resume file", slog.Any("error", err))
		return response.InternalServerError(c, ErrFileReadError)
	}
	defer file.Close()

	// Set headers for PDF preview (inline display)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", `inline; filename="`+resume.FileName+`"`)

	// Stream file to response
	_, err = io.Copy(c.Response().BodyWriter(), file)
	if err != nil {
		h.logger.Error("failed to stream resume file", slog.Any("error", err))
		return response.InternalServerError(c, "failed to stream file")
	}

	return nil
}

