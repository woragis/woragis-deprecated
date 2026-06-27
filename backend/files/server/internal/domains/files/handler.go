package files

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"woragis-files-service/pkg/middleware"
	"woragis-files-service/pkg/utils"
)

// Handler exposes file endpoints.
type Handler interface {
	UploadFile(c *fiber.Ctx) error
	GetFile(c *fiber.Ctx) error
	ListFiles(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
	DownloadFile(c *fiber.Ctx) error
	GetFileURL(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a file handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// UploadFile handles file upload.
// @Summary Upload a file
// @Description Upload a file to the storage system
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /files [post]
func (h *handler) UploadFile(c *fiber.Ctx) error {
	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, 0, err.Error(), nil)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "invalid user ID", nil)
	}

	// Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "file is required", nil)
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "failed to open file", nil)
	}
	defer file.Close()

	// Create upload request
	req := UploadFileRequest{
		FileName: fileHeader.Filename,
		MimeType: fileHeader.Header.Get("Content-Type"),
		Size:     fileHeader.Size,
		Reader:   file,
	}

	// Upload file
	uploadedFile, err := h.service.UploadFile(c.Context(), userUUID, req)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			return utils.Error(c, fiber.StatusBadRequest, 0, domainErr.Message, nil)
		}
		h.logger.Error("failed to upload file", "error", err)
		return utils.Error(c, fiber.StatusInternalServerError, 0, "failed to upload file", nil)
	}

	return utils.Success(c, fiber.StatusCreated, "File uploaded successfully", uploadedFile)
}

// GetFile retrieves file metadata.
// @Summary Get file metadata
// @Description Get metadata for a specific file
// @Tags files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /files/:id [get]
func (h *handler) GetFile(c *fiber.Ctx) error {
	fileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "invalid file ID", nil)
	}

	file, err := h.service.GetFile(c.Context(), fileID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok && domainErr.Code == ErrCodeNotFound {
			return utils.Error(c, fiber.StatusNotFound, 0, domainErr.Message, nil)
		}
		return utils.Error(c, fiber.StatusInternalServerError, 0, "failed to get file", nil)
	}

	return utils.Success(c, fiber.StatusOK, "File retrieved successfully", file)
}

// ListFiles lists files with optional filters.
// @Summary List files
// @Description List files with optional filters
// @Tags files
// @Produce json
// @Param file_type query string false "File type filter"
// @Param status query string false "Status filter"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} utils.SuccessResponse
// @Router /files [get]
func (h *handler) ListFiles(c *fiber.Ctx) error {
	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, 0, err.Error(), nil)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "invalid user ID", nil)
	}

	filters := FileFilters{
		UserID: &userUUID,
	}

	// Parse query parameters
	if fileTypeStr := c.Query("file_type"); fileTypeStr != "" {
		fileType := FileType(fileTypeStr)
		filters.FileType = &fileType
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := FileStatus(statusStr)
		filters.Status = &status
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	if orderBy := c.Query("order_by"); orderBy != "" {
		filters.OrderBy = orderBy
	}

	if order := c.Query("order"); order != "" {
		filters.Order = order
	}

	files, err := h.service.ListFiles(c.Context(), filters)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, 0, "failed to list files", nil)
	}

	return utils.Success(c, fiber.StatusOK, "Files retrieved successfully", files)
}

// DeleteFile deletes a file.
// @Summary Delete a file
// @Description Delete a file by ID
// @Tags files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /files/:id [delete]
func (h *handler) DeleteFile(c *fiber.Ctx) error {
	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, 0, err.Error(), nil)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "invalid user ID", nil)
	}

	fileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "invalid file ID", nil)
	}

	if err := h.service.DeleteFile(c.Context(), fileID, userUUID); err != nil {
		if domainErr, ok := err.(*DomainError); ok && domainErr.Code == ErrCodeNotFound {
			return utils.Error(c, fiber.StatusNotFound, 0, domainErr.Message, nil)
		}
		return utils.Error(c, fiber.StatusInternalServerError, 0, "failed to delete file", nil)
	}

	return utils.Success(c, fiber.StatusOK, "File deleted successfully", nil)
}

// DownloadFile downloads a file.
// @Summary Download a file
// @Description Download a file by ID
// @Tags files
// @Produce application/octet-stream
// @Param id path string true "File ID"
// @Success 200
// @Failure 404 {object} utils.ErrorResponse
// @Router /files/:id/download [get]
func (h *handler) DownloadFile(c *fiber.Ctx) error {
	fileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "invalid file ID", nil)
	}

	reader, file, err := h.service.DownloadFile(c.Context(), fileID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok && domainErr.Code == ErrCodeNotFound {
			return utils.Error(c, fiber.StatusNotFound, 0, domainErr.Message, nil)
		}
		return utils.Error(c, fiber.StatusInternalServerError, 0, "failed to download file", nil)
	}
	defer reader.Close()

	// Set headers
	c.Set("Content-Type", file.MimeType)
	c.Set("Content-Disposition", "attachment; filename=\""+file.OriginalName+"\"")

	// Stream file
	return c.SendStream(reader, int(file.Size))
}

// GetFileURL gets a file URL.
// @Summary Get file URL
// @Description Get a URL for accessing a file
// @Tags files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /files/:id/url [get]
func (h *handler) GetFileURL(c *fiber.Ctx) error {
	fileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, 0, "invalid file ID", nil)
	}

	url, err := h.service.GetFileURL(c.Context(), fileID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok && domainErr.Code == ErrCodeNotFound {
			return utils.Error(c, fiber.StatusNotFound, 0, domainErr.Message, nil)
		}
		return utils.Error(c, fiber.StatusInternalServerError, 0, "failed to get file URL", nil)
	}

	return utils.Success(c, fiber.StatusOK, "File URL retrieved successfully", map[string]string{
		"url": url,
	})
}

