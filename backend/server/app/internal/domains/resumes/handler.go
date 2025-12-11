package resumes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	ListResumeTags(c *fiber.Ctx) error
	DownloadResume(c *fiber.Ctx) error
	PreviewResume(c *fiber.Ctx) error
	GenerateResume(c *fiber.Ctx) error
	MarkAsMain(c *fiber.Ctx) error
	MarkAsFeatured(c *fiber.Ctx) error
	UnmarkAsMain(c *fiber.Ctx) error
	UnmarkAsFeatured(c *fiber.Ctx) error
	RecalculateMetrics(c *fiber.Ctx) error
}

// JobApplicationService is an interface to avoid circular dependencies
type JobApplicationService interface {
	GetJobApplication(ctx context.Context, applicationID uuid.UUID) (*JobApplication, error)
}

// JobApplication represents a job application (minimal interface)
type JobApplication struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	JobTitle       string
	JobDescription string
	Language       string
	CompanyName    string
}

type handler struct {
	service              Service
	jobApplicationService JobApplicationService // Optional: for generating resumes
	logger               *slog.Logger
	baseFilePath         string                 // Base path where resume files are stored
	resumeWorkerPath     string                 // Path to resume worker script or Docker container
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a resume handler.
func NewHandler(service Service, baseFilePath string, logger *slog.Logger) Handler {
	return &handler{
		service:      service,
		logger:       logger,
		baseFilePath: baseFilePath,
		resumeWorkerPath: "woragis-resume-worker", // Docker container name
	}
}

// NewHandlerWithJobApplicationService constructs a resume handler with job application service.
func NewHandlerWithJobApplicationService(service Service, jobApplicationService JobApplicationService, baseFilePath string, logger *slog.Logger) Handler {
	return &handler{
		service:              service,
		jobApplicationService: jobApplicationService,
		logger:               logger,
		baseFilePath:         baseFilePath,
		resumeWorkerPath:     "woragis-resume-worker", // Docker container name
	}
}

// CreateResume creates a new resume.
func (h *handler) CreateResume(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	var req struct {
		Title    string   `json:"title"`
		FilePath string   `json:"filePath"`
		FileName string   `json:"fileName"`
		FileSize int64    `json:"fileSize"`
		Tags     []string `json:"tags"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid request body"})
	}

	resume, err := h.service.CreateResume(c.Context(), userID, req.Title, req.FilePath, req.FileName, req.FileSize, JSONArray(req.Tags))
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

	// Create resume entry (tags can be added later via update)
	resume, err := h.service.CreateResume(c.Context(), userID, title, filePath, fileHeader.Filename, fileInfo.Size(), JSONArray{})
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
		Title string   `json:"title"`
		Tags  []string `json:"tags"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid request body"})
	}

	var tags JSONArray
	if req.Tags != nil {
		tags = JSONArray(req.Tags)
	}

	resume, err := h.service.UpdateResume(c.Context(), userID, resumeID, req.Title, tags)
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

// ListResumes lists all resumes for the authenticated user, optionally filtered by tags.
func (h *handler) ListResumes(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	// Check for tag filter query parameter
	tagFilter := c.Query("tags")
	var resumes []Resume
	if tagFilter != "" {
		// Parse comma-separated tags
		tags := strings.Split(tagFilter, ",")
		normalizedTags := make([]string, 0, len(tags))
		for _, tag := range tags {
			normalized := strings.ToLower(strings.TrimSpace(tag))
			if normalized != "" {
				normalizedTags = append(normalizedTags, normalized)
			}
		}
		if len(normalizedTags) > 0 {
			resumes, err = h.service.ListResumesByTags(c.Context(), userID, normalizedTags)
		} else {
			resumes, err = h.service.ListResumes(c.Context(), userID)
		}
	} else {
		resumes, err = h.service.ListResumes(c.Context(), userID)
	}

	if err != nil {
		h.logger.Error("failed to list resumes", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to list resumes"})
	}

	return response.Success(c, fiber.StatusOK, resumes)
}

// ListResumeTags returns all unique tags from all resumes for the authenticated user (for autocomplete).
func (h *handler) ListResumeTags(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumes, err := h.service.ListResumes(c.Context(), userID)
	if err != nil {
		h.logger.Error("failed to list resumes", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to list resumes"})
	}

	// Collect all unique tags
	tagSet := make(map[string]bool)
	for _, resume := range resumes {
		for _, tag := range resume.Tags {
			if tag != "" {
				tagSet[tag] = true
			}
		}
	}

	// Convert to slice
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	return response.Success(c, fiber.StatusOK, tags)
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

// GenerateResume generates a resume for a job application using the resume worker.
func (h *handler) GenerateResume(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	var req struct {
		JobApplicationID string `json:"jobApplicationId"`
		Language         string `json:"language"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid request body"})
	}

	if req.JobApplicationID == "" {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "jobApplicationId is required"})
	}

	jobAppID, err := uuid.Parse(req.JobApplicationID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid job application ID"})
	}

	// Get job application details
	if h.jobApplicationService == nil {
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "job application service not configured"})
	}

	jobApp, err := h.jobApplicationService.GetJobApplication(c.Context(), jobAppID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": "job application not found"})
	}
	
	// Verify the job application belongs to the user
	if jobApp.UserID != userID {
		return response.Error(c, fiber.StatusForbidden, 0, fiber.Map{"message": "access denied"})
	}

	// Prepare job description
	jobDescription := jobApp.JobDescription
	if jobDescription == "" {
		jobDescription = fmt.Sprintf("Position: %s at %s", jobApp.JobTitle, jobApp.CompanyName)
	}

	language := req.Language
	if language == "" {
		language = jobApp.Language
	}
	if language == "" {
		language = "en"
	}

	// Execute resume worker via Docker
	// Format: python main.py <user_id> <job_description> [job_title] [output_filename] [language]
	cmd := exec.Command("docker", "exec", h.resumeWorkerPath, "python", "/app/src/main.py",
		userID.String(),
		jobDescription,
		jobApp.JobTitle,
		"", // output_filename - let worker generate it
		language,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		h.logger.Error("failed to execute resume worker", slog.Any("error", err), slog.String("output", string(output)))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{
			"message": "failed to generate resume",
			"error":   err.Error(),
		})
	}

	// Parse output to get the generated file path
	// The worker prints: "Resume generated: <path>"
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	var generatedPath string
	for _, line := range lines {
		if strings.HasPrefix(line, "Resume generated:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				generatedPath = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if generatedPath == "" {
		// Try to find the latest result JSON file
		resultsDir := filepath.Join(h.baseFilePath, "../results")
		if resultsDir == "" {
			resultsDir = "/app/results" // Default in Docker
		}
		
		// Look for the most recent result file
		files, err := filepath.Glob(filepath.Join(resultsDir, fmt.Sprintf("resume_result_%s_*.json", userID.String())))
		if err == nil && len(files) > 0 {
			// Get the most recent file
			var latestFile string
			var latestTime time.Time
			for _, file := range files {
				info, err := os.Stat(file)
				if err == nil && info.ModTime().After(latestTime) {
					latestTime = info.ModTime()
					latestFile = file
				}
			}
			
			if latestFile != "" {
				// Read the JSON to get the output path
				data, err := os.ReadFile(latestFile)
				if err == nil {
					var result struct {
						OutputPath string `json:"output_path"`
						FileName   string `json:"file_name"`
						FileSize   int64  `json:"file_size"`
					}
					if json.Unmarshal(data, &result) == nil && result.OutputPath != "" {
						generatedPath = result.OutputPath
					}
				}
			}
		}
	}

	if generatedPath == "" {
		h.logger.Error("could not determine generated resume path", slog.String("output", outputStr))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to determine generated resume path"})
	}

	// Extract relative path from absolute path
	// The path from worker is like /app/output/resume_Backend_Engineer_20241129_120000.pdf
	// We need to convert it to a relative path like "output/resume_Backend_Engineer_20241129_120000.pdf"
	relativePath := generatedPath
	if strings.HasPrefix(generatedPath, "/app/") {
		relativePath = strings.TrimPrefix(generatedPath, "/app/")
	} else if strings.HasPrefix(generatedPath, h.baseFilePath) {
		relativePath, _ = filepath.Rel(h.baseFilePath, generatedPath)
	}

	// Get file info
	fileInfo, err := os.Stat(generatedPath)
	if err != nil {
		// Try with base path
		fullPath := filepath.Join(h.baseFilePath, relativePath)
		fileInfo, err = os.Stat(fullPath)
		if err != nil {
			h.logger.Error("failed to get file info", slog.String("path", generatedPath), slog.Any("error", err))
			return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "generated resume file not found"})
		}
		generatedPath = fullPath
	}

	fileName := filepath.Base(generatedPath)
	title := fmt.Sprintf("%s - %s", jobApp.JobTitle, jobApp.CompanyName)

	// Extract tags from result JSON if available
	var tags JSONArray
	if generatedPath != "" {
		// Try to read tags from the result file
		resultsDir := filepath.Join(h.baseFilePath, "../results")
		if resultsDir == "" {
			resultsDir = "/app/results"
		}
		files, err := filepath.Glob(filepath.Join(resultsDir, fmt.Sprintf("resume_result_%s_*.json", userID.String())))
		if err == nil && len(files) > 0 {
			var latestFile string
			var latestTime time.Time
			for _, file := range files {
				info, err := os.Stat(file)
				if err == nil && info.ModTime().After(latestTime) {
					latestTime = info.ModTime()
					latestFile = file
				}
			}
			if latestFile != "" {
				data, err := os.ReadFile(latestFile)
				if err == nil {
					var result struct {
						Tags []string `json:"tags"`
					}
					if json.Unmarshal(data, &result) == nil && len(result.Tags) > 0 {
						tags = JSONArray(result.Tags)
					}
				}
			}
		}
	}

	// Create resume entry with tags
	resume, err := h.service.CreateResume(c.Context(), userID, title, relativePath, fileName, fileInfo.Size(), tags)
	if err != nil {
		h.logger.Error("failed to create resume entry", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to create resume entry"})
	}

	return response.Success(c, fiber.StatusCreated, resume)
}

// RecalculateMetrics manually recalculates metrics for a resume.
func (h *handler) RecalculateMetrics(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid resume ID"})
	}

	// Verify the resume belongs to the user
	_, err = h.service.GetResume(c.Context(), userID, resumeID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			if domainErr.Code == ErrCodeNotFound {
				return response.Error(c, fiber.StatusNotFound, 0, fiber.Map{"message": domainErr.Message})
			}
		}
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to get resume"})
	}

	// Recalculate metrics
	if err := h.service.RecalculateResumeMetrics(c.Context(), resumeID); err != nil {
		h.logger.Error("failed to recalculate resume metrics", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to recalculate metrics"})
	}

	// Return updated resume
	resume, err := h.service.GetResume(c.Context(), userID, resumeID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to get updated resume"})
	}

	return response.Success(c, fiber.StatusOK, resume)
}

