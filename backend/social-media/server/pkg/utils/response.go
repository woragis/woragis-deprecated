package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Response
	Pagination Pagination `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// SendJSON sends a JSON response
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendSuccess sends a success response
func SendSuccess(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	SendJSON(w, http.StatusOK, response)
}

// SendError sends an error response
func SendError(w http.ResponseWriter, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
		Code:    statusCode,
	}
	
	if err != nil {
		response.Error = err.Error()
	}
	
	SendJSON(w, statusCode, response)
}

// SendValidationError sends a validation error response
func SendValidationError(w http.ResponseWriter, message string, errors interface{}) {
	response := Response{
		Success: false,
		Message: message,
		Data:    errors,
		Code:    http.StatusBadRequest,
	}
	SendJSON(w, http.StatusBadRequest, response)
}

// SendPaginatedResponse sends a paginated response
func SendPaginatedResponse(w http.ResponseWriter, data interface{}, pagination Pagination) {
	response := PaginatedResponse{
		Response: Response{
			Success: true,
			Data:    data,
		},
		Pagination: pagination,
	}
	SendJSON(w, http.StatusOK, response)
}

// SendUnauthorized sends an unauthorized response
func SendUnauthorized(w http.ResponseWriter, message string) {
	SendError(w, http.StatusUnauthorized, message, nil)
}

// SendForbidden sends a forbidden response
func SendForbidden(w http.ResponseWriter, message string) {
	SendError(w, http.StatusForbidden, message, nil)
}

// SendNotFound sends a not found response
func SendNotFound(w http.ResponseWriter, message string) {
	SendError(w, http.StatusNotFound, message, nil)
}

// SendInternalError sends an internal server error response
func SendInternalError(w http.ResponseWriter, message string, err error) {
	SendError(w, http.StatusInternalServerError, message, err)
}

// SendBadRequest sends a bad request response
func SendBadRequest(w http.ResponseWriter, message string) {
	SendError(w, http.StatusBadRequest, message, nil)
}

// SendCreated sends a created response
func SendCreated(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	SendJSON(w, http.StatusCreated, response)
}

// SendNoContent sends a no content response
func SendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, limit int, total int64) Pagination {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	
	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// Fiber Response Functions

// SuccessResponse sends a success response using Fiber
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// CreatedResponse sends a created response using Fiber
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// BadRequestResponse sends a bad request response using Fiber
func BadRequestResponse(c *fiber.Ctx, message string) error {
	response := Response{
		Success: false,
		Message: message,
		Code:    fiber.StatusBadRequest,
	}
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

// UnauthorizedResponse sends an unauthorized response using Fiber
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	response := Response{
		Success: false,
		Message: message,
		Code:    fiber.StatusUnauthorized,
	}
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

// ForbiddenResponse sends a forbidden response using Fiber
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	response := Response{
		Success: false,
		Message: message,
		Code:    fiber.StatusForbidden,
	}
	return c.Status(fiber.StatusForbidden).JSON(response)
}

// NotFoundResponse sends a not found response using Fiber
func NotFoundResponse(c *fiber.Ctx, message string) error {
	response := Response{
		Success: false,
		Message: message,
		Code:    fiber.StatusNotFound,
	}
	return c.Status(fiber.StatusNotFound).JSON(response)
}

// ConflictResponse sends a conflict response using Fiber
func ConflictResponse(c *fiber.Ctx, message string) error {
	response := Response{
		Success: false,
		Message: message,
		Code:    fiber.StatusConflict,
	}
	return c.Status(fiber.StatusConflict).JSON(response)
}

// InternalServerErrorResponse sends an internal server error response using Fiber
func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	response := Response{
		Success: false,
		Message: message,
		Code:    fiber.StatusInternalServerError,
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

// ValidationErrorResponse sends a validation error response using Fiber
func ValidationErrorResponse(c *fiber.Ctx, err error) error {
	response := Response{
		Success: false,
		Message: "Validation failed",
		Error:   err.Error(),
		Code:    fiber.StatusBadRequest,
	}
	return c.Status(fiber.StatusBadRequest).JSON(response)
}


// ValidateStruct validates a struct using basic validation
func ValidateStruct(s interface{}) error {
	// This is a simplified validation - in a real implementation,
	// you would use a validation library like go-playground/validator
	return nil
}
