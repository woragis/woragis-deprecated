package resumes

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes resume endpoints.
type Handler interface {
	CreateResume(c *fiber.Ctx) error
	UploadResume(c *fiber.Ctx) error
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
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	var req struct {
		Title    string `json:"title"`
		FilePath string `json:"filePath"`
		FileName string `json:"fileName"`
		FileSize int64  `json:"fileSize"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid request body"})
	}

	resume, err := h.service.CreateResume(c.Context(), userID, req.Title, req.FilePath, req.FileName, req.FileSize)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to create resume", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to create resume"})
	}

	return response.Success(c, fiber.StatusCreated, resume)
}

// UploadResume handles file upload and creates a new resume.
func (h *handler) UploadResume(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid multipart form"})
	}

	// Get title from form
	titleValues := form.Value["title"]
	if len(titleValues) == 0 || titleValues[0] == "" {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "title is required"})
	}
	title := titleValues[0]

	// Get file from form
	files := form.File["file"]
	if len(files) == 0 {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "file is required"})
	}
	fileHeader := files[0]

	// Validate file type (should be PDF)
	if fileHeader.Header.Get("Content-Type") != "application/pdf" {
		// Check extension as fallback
		ext := filepath.Ext(fileHeader.Filename)
		if ext != ".pdf" {
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "only PDF files are allowed"})
		}
	}

	// Create upload directory if it doesn't exist
	uploadDir := filepath.Join(h.baseFilePath, "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		h.logger.Error("failed to create upload directory", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to create upload directory"})
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	safeFilename := fmt.Sprintf("%d_%s", timestamp, fileHeader.Filename)
	filePath := filepath.Join("uploads", safeFilename)
	fullPath := filepath.Join(h.baseFilePath, filePath)

	// Save file
	if err := c.SaveFile(fileHeader, fullPath); err != nil {
		h.logger.Error("failed to save file", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to save file"})
	}

	// Get file size
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		h.logger.Error("failed to get file info", slog.Any("error", err))
		// Clean up file
		os.Remove(fullPath)
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to get file info"})
	}

	// Create resume entry
	resume, err := h.service.CreateResume(c.Context(), userID, title, filePath, fileHeader.Filename, fileInfo.Size())
	if err != nil {
		// Clean up file if resume creation fails
		os.Remove(fullPath)
		if domainErr, ok := err.(*DomainError); ok {
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to create resume", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to create resume"})
	}

	return response.Success(c, fiber.StatusCreated, resume)
}

// UpdateResume updates an existing resume.
func (h *handler) UpdateResume(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid resume ID"})
	}

	var req struct {
		Title string `json:"title"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid request body"})
	}

	resume, err := h.service.UpdateResume(c.Context(), userID, resumeID, req.Title)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to update resume", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to update resume"})
	}

	return response.Success(c, fiber.StatusOK, resume)
}

// DeleteResume deletes a resume.
func (h *handler) DeleteResume(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid resume ID"})
	}

	// Get resume first to get file path for deletion
	resume, err := h.service.GetResume(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to get resume for deletion", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to get resume"})
	}

	// Delete the resume from database
	if err := h.service.DeleteResume(c.Context(), userID, resumeID); err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to delete resume", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to delete resume"})
	}

	// Delete the file
	fullPath := filepath.Join(h.baseFilePath, resume.FilePath)
	if !filepath.IsAbs(resume.FilePath) {
		fullPath = filepath.Join(h.baseFilePath, resume.FilePath)
	} else {
		fullPath = resume.FilePath
	}

	// Try to delete the file (ignore error if file doesn't exist)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		h.logger.Warn("failed to delete resume file", slog.String("path", fullPath), slog.Any("error", err))
		// Don't fail the request if file deletion fails
	}

	return response.Success(c, fiber.StatusNoContent, nil)
}

// GetResume retrieves a resume by ID.
func (h *handler) GetResume(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid resume ID"})
	}

	resume, err := h.service.GetResume(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to get resume", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to get resume"})
	}

	return response.Success(c, fiber.StatusOK, resume)
}

// ListResumes lists all resumes for the authenticated user.
func (h *handler) ListResumes(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumes, err := h.service.ListResumes(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to list resumes", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to list resumes"})
	}

	return response.Success(c, fiber.StatusOK, resumes)
}

// MarkAsMain marks a resume as main.
func (h *handler) MarkAsMain(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid resume ID"})
	}

	resume, err := h.service.MarkAsMain(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to mark resume as main", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to mark resume as main"})
	}

	return response.Success(c, fiber.StatusOK, resume)
}

// MarkAsFeatured marks a resume as featured.
func (h *handler) MarkAsFeatured(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid resume ID"})
	}

	resume, err := h.service.MarkAsFeatured(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to mark resume as featured", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to mark resume as featured"})
	}

	return response.Success(c, fiber.StatusOK, resume)
}

// UnmarkAsMain removes the main flag from a resume.
func (h *handler) UnmarkAsMain(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid resume ID"})
	}

	resume, err := h.service.UnmarkAsMain(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to unmark resume as main", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to unmark resume as main"})
	}

	return response.Success(c, fiber.StatusOK, resume)
}

// UnmarkAsFeatured removes the featured flag from a resume.
func (h *handler) UnmarkAsFeatured(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid resume ID"})
	}

	resume, err := h.service.UnmarkAsFeatured(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
			return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": domainErr.Message})
		}
		h.logger.Error("failed to unmark resume as featured", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to unmark resume as featured"})
	}

	return response.Success(c, fiber.StatusOK, resume)
}

// DownloadResume downloads the best resume (public endpoint).
func (h *handler) DownloadResume(c *fiber.Ctx) error {
	// Get user ID from query param or use default user
	// For now, we'll get it from a query param or use a default
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		// Use a default user ID - you might want to configure this
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "userId query parameter is required"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid userId"})
	}

	// Get the best resume (main > featured > most recent)
	resume, err := h.service.GetBestResume(c.Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
		}
		h.logger.Error("failed to get resume for download", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to get resume"})
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
		return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": ErrFileNotFound})
	}

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		h.logger.Error("failed to open resume file", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": ErrFileReadError})
	}
	defer file.Close()

	// Set headers for PDF download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", `attachment; filename="`+resume.FileName+`"`)

	// Stream file to response
	_, err = io.Copy(c.Response().BodyWriter(), file)
	if err != nil {
		h.logger.Error("failed to stream resume file", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to stream file"})
	}

	return nil
}

// PreviewResume serves the resume for preview (public endpoint).
func (h *handler) PreviewResume(c *fiber.Ctx) error {
	// Get user ID from query param
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "userId query parameter is required"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid userId"})
	}

	// Get the best resume
	resume, err := h.service.GetBestResume(c.Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
		}
		h.logger.Error("failed to get resume for preview", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to get resume"})
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
		return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": ErrFileNotFound})
	}

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		h.logger.Error("failed to open resume file", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": ErrFileReadError})
	}
	defer file.Close()

	// Set headers for PDF preview (inline display)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", `inline; filename="`+resume.FileName+`"`)

	// Stream file to response
	_, err = io.Copy(c.Response().BodyWriter(), file)
	if err != nil {
		h.logger.Error("failed to stream resume file", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to stream file"})
	}

	return nil
}

