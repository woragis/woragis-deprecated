package certifications

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationenricher "github.com/woragis/backend/server/app/pkg/translations"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes certification endpoints.
type Handler interface {
	CreateCertification(c *fiber.Ctx) error
	UpdateCertification(c *fiber.Ctx) error
	GetCertification(c *fiber.Ctx) error
	GetCertificationPublic(c *fiber.Ctx) error
	ListCertifications(c *fiber.Ctx) error
	ListFeaturedCertifications(c *fiber.Ctx) error
	GetCertificationsBySkill(c *fiber.Ctx) error
	DeleteCertification(c *fiber.Ctx) error
	AddCertificationSkill(c *fiber.Ctx) error
	RemoveCertificationSkill(c *fiber.Ctx) error
	GetCertificationSkills(c *fiber.Ctx) error
	// Entity link methods
	CreateCertificationEntityLink(c *fiber.Ctx) error
	GetCertificationEntityLinks(c *fiber.Ctx) error
	GetEntityCertifications(c *fiber.Ctx) error
	DeleteCertificationEntityLink(c *fiber.Ctx) error
	DeleteCertificationEntityLinks(c *fiber.Ctx) error
}

type handler struct {
	service           Service
	enricher          *translationenricher.Enricher
	translationService translationsdomain.Service
	logger            *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a certification handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createCertificationPayload struct {
	Name            string                `json:"name"`
	Issuer          string                `json:"issuer"`
	IssueDate       string                `json:"issueDate"` // ISO 8601 date string
	ExpiryDate      *string               `json:"expiryDate,omitempty"`
	CredentialID    string                `json:"credentialId,omitempty"`
	VerificationURL string                `json:"verificationUrl,omitempty"`
	CertificateURL  string                `json:"certificateUrl,omitempty"`
	Description     string                `json:"description,omitempty"`
	Status          CertificationStatus  `json:"status,omitempty"`
	Category        CertificationCategory `json:"category"`
	Featured        bool                  `json:"featured,omitempty"`
	DisplayOrder    int                   `json:"displayOrder,omitempty"`
	SkillIDs        []string              `json:"skillIds,omitempty"`
}

type updateCertificationPayload struct {
	Name            *string                `json:"name,omitempty"`
	Issuer          *string                 `json:"issuer,omitempty"`
	IssueDate       *string                 `json:"issueDate,omitempty"`
	ExpiryDate      *string                 `json:"expiryDate,omitempty"`
	CredentialID    *string                 `json:"credentialId,omitempty"`
	VerificationURL *string                 `json:"verificationUrl,omitempty"`
	CertificateURL  *string                 `json:"certificateUrl,omitempty"`
	Description     *string                 `json:"description,omitempty"`
	Status          *CertificationStatus    `json:"status,omitempty"`
	Category        *CertificationCategory  `json:"category,omitempty"`
	Featured        *bool                   `json:"featured,omitempty"`
	DisplayOrder    *int                    `json:"displayOrder,omitempty"`
	SkillIDs        []string                `json:"skillIds,omitempty"`
}

// Handlers

func (h *handler) CreateCertification(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createCertificationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Parse dates
	issueDate, err := parseDate(payload.IssueDate)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidDate, fiber.Map{
			"message": "invalid issue date format",
		})
	}

	var expiryDate *time.Time
	if payload.ExpiryDate != nil && *payload.ExpiryDate != "" {
		exp, err := parseDate(*payload.ExpiryDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidDate, fiber.Map{
				"message": "invalid expiry date format",
			})
		}
		expiryDate = &exp
	}

	// Parse skill IDs
	var skillIDs []uuid.UUID
	if len(payload.SkillIDs) > 0 {
		skillIDs = make([]uuid.UUID, 0, len(payload.SkillIDs))
		for _, skillIDStr := range payload.SkillIDs {
			skillID, err := uuid.Parse(skillIDStr)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
					"message": "invalid skill id format",
				})
			}
			skillIDs = append(skillIDs, skillID)
		}
	}

	certification, err := h.service.CreateCertification(c.Context(), userID, CreateCertificationRequest{
		Name:            payload.Name,
		Issuer:          payload.Issuer,
		IssueDate:       issueDate,
		ExpiryDate:      expiryDate,
		CredentialID:    payload.CredentialID,
		VerificationURL: payload.VerificationURL,
		CertificateURL:  payload.CertificateURL,
		Description:     payload.Description,
		Status:          payload.Status,
		Category:        payload.Category,
		Featured:        payload.Featured,
		DisplayOrder:    payload.DisplayOrder,
		SkillIDs:        skillIDs,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Automatically trigger translations for all supported languages
	if h.translationService != nil {
		sourceText := make(map[string]string)
		if certification.Name != "" {
			sourceText["name"] = certification.Name
		}
		if certification.Issuer != "" {
			sourceText["issuer"] = certification.Issuer
		}
		if certification.Description != "" {
			sourceText["description"] = certification.Description
		}

		fields := []string{}
		if certification.Name != "" {
			fields = append(fields, "name")
		}
		if certification.Issuer != "" {
			fields = append(fields, "issuer")
		}
		if certification.Description != "" {
			fields = append(fields, "description")
		}

		supportedLanguages := []translationsdomain.Language{
			translationsdomain.LanguagePTBR,
			translationsdomain.LanguageFR,
			translationsdomain.LanguageES,
			translationsdomain.LanguageDE,
			translationsdomain.LanguageRU,
			translationsdomain.LanguageJA,
			translationsdomain.LanguageKO,
			translationsdomain.LanguageZHCN,
			translationsdomain.LanguageEL,
			translationsdomain.LanguageLA,
		}

		go func() {
			ctx := context.Background()
			for _, lang := range supportedLanguages {
				if err := h.translationService.RequestTranslation(
					ctx,
					translationsdomain.EntityTypeCertification,
					certification.ID,
					lang,
					fields,
					sourceText,
				); err != nil {
					h.logger.Warn("Failed to queue translation",
						slog.String("certificationId", certification.ID.String()),
						slog.String("language", string(lang)),
						slog.Any("error", err),
					)
				}
			}
		}()
	}

	return response.Success(c, fiber.StatusCreated, certification)
}

func (h *handler) UpdateCertification(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	var payload updateCertificationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := UpdateCertificationRequest{}

	if payload.Name != nil {
		req.Name = payload.Name
	}
	if payload.Issuer != nil {
		req.Issuer = payload.Issuer
	}
	if payload.IssueDate != nil {
		issueDate, err := parseDate(*payload.IssueDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidDate, fiber.Map{
				"message": "invalid issue date format",
			})
		}
		req.IssueDate = &issueDate
	}
	if payload.ExpiryDate != nil {
		if *payload.ExpiryDate != "" {
			expiryDate, err := parseDate(*payload.ExpiryDate)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidDate, fiber.Map{
					"message": "invalid expiry date format",
				})
			}
			req.ExpiryDate = &expiryDate
		} else {
			req.ExpiryDate = nil
		}
	}
	if payload.CredentialID != nil {
		req.CredentialID = payload.CredentialID
	}
	if payload.VerificationURL != nil {
		req.VerificationURL = payload.VerificationURL
	}
	if payload.CertificateURL != nil {
		req.CertificateURL = payload.CertificateURL
	}
	if payload.Description != nil {
		req.Description = payload.Description
	}
	if payload.Status != nil {
		req.Status = payload.Status
	}
	if payload.Category != nil {
		req.Category = payload.Category
	}
	if payload.Featured != nil {
		req.Featured = payload.Featured
	}
	if payload.DisplayOrder != nil {
		req.DisplayOrder = payload.DisplayOrder
	}
	if payload.SkillIDs != nil {
		skillIDs := make([]uuid.UUID, 0, len(payload.SkillIDs))
		for _, skillIDStr := range payload.SkillIDs {
			skillID, err := uuid.Parse(skillIDStr)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
					"message": "invalid skill id format",
				})
			}
			skillIDs = append(skillIDs, skillID)
		}
		req.SkillIDs = skillIDs
	}

	certification, err := h.service.UpdateCertification(c.Context(), userID, certificationID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, certification)
}

func (h *handler) GetCertification(c *fiber.Ctx) error {
	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certification, err := h.service.GetCertification(c.Context(), certificationID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"name":        &certification.Name,
			"issuer":     &certification.Issuer,
			"description": &certification.Description,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeCertification, certification.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, certification)
}

func (h *handler) GetCertificationPublic(c *fiber.Ctx) error {
	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	certification, err := h.service.GetCertificationPublic(c.Context(), certificationID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"name":        &certification.Name,
			"issuer":     &certification.Issuer,
			"description": &certification.Description,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeCertification, certification.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, certification)
}

func (h *handler) ListCertifications(c *fiber.Ctx) error {
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		uid := apiKey.UserID
		userID = &uid
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := ListCertificationsFilters{
		UserID: userID,
	}

	// Parse query parameters
	if statusStr := c.Query("status"); statusStr != "" {
		status := CertificationStatus(statusStr)
		filters.Status = &status
	} else {
		// Default to active for public access
		active := CertificationStatusActive
		filters.Status = &active
	}

	if categoryStr := c.Query("category"); categoryStr != "" {
		category := CertificationCategory(categoryStr)
		filters.Category = &category
	}

	if issuer := c.Query("issuer"); issuer != "" {
		filters.Issuer = &issuer
	}

	if featuredStr := c.Query("featured"); featuredStr != "" {
		featured := featuredStr == "true"
		filters.Featured = &featured
	}

	if expiringSoonStr := c.Query("expiringSoon"); expiringSoonStr != "" {
		expiringSoon := expiringSoonStr == "true"
		filters.ExpiringSoon = &expiringSoon
	}

	if skillIDStr := c.Query("skillId"); skillIDStr != "" {
		skillID, err := uuid.Parse(skillIDStr)
		if err == nil {
			filters.SkillID = &skillID
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	if orderBy := c.Query("orderBy"); orderBy != "" {
		filters.OrderBy = orderBy
	}

	if order := c.Query("order"); order != "" {
		filters.Order = order
	}

	certifications, err := h.service.ListCertifications(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range certifications {
			fieldMap := map[string]*string{
				"name":        &certifications[i].Name,
				"issuer":     &certifications[i].Issuer,
				"description": &certifications[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeCertification, certifications[i].ID, language, fieldMap)
		}
	}

	return response.Success(c, fiber.StatusOK, certifications)
}

func (h *handler) ListFeaturedCertifications(c *fiber.Ctx) error {
	certifications, err := h.service.ListFeaturedCertifications(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range certifications {
			fieldMap := map[string]*string{
				"name":        &certifications[i].Name,
				"issuer":     &certifications[i].Issuer,
				"description": &certifications[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeCertification, certifications[i].ID, language, fieldMap)
		}
	}

	return response.Success(c, fiber.StatusOK, certifications)
}

func (h *handler) GetCertificationsBySkill(c *fiber.Ctx) error {
	skillID, err := uuid.Parse(c.Params("skillId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid skill id",
		})
	}

	certifications, err := h.service.GetCertificationsBySkill(c.Context(), skillID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range certifications {
			fieldMap := map[string]*string{
				"name":        &certifications[i].Name,
				"issuer":     &certifications[i].Issuer,
				"description": &certifications[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeCertification, certifications[i].ID, language, fieldMap)
		}
	}

	return response.Success(c, fiber.StatusOK, certifications)
}

func (h *handler) DeleteCertification(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	if err := h.service.DeleteCertification(c.Context(), userID, certificationID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "certification deleted successfully",
	})
}

func (h *handler) AddCertificationSkill(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	skillID, err := uuid.Parse(c.Params("skillId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid skill id",
		})
	}

	if err := h.service.AddCertificationSkill(c.Context(), userID, certificationID, skillID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "skill linked to certification successfully",
	})
}

func (h *handler) RemoveCertificationSkill(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	skillID, err := uuid.Parse(c.Params("skillId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid skill id",
		})
	}

	if err := h.service.RemoveCertificationSkill(c.Context(), userID, certificationID, skillID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "skill unlinked from certification successfully",
	})
}

func (h *handler) GetCertificationSkills(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	skillIDs, err := h.service.GetCertificationSkills(c.Context(), certificationID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"skillIds": skillIDs,
	})
}

// Entity link handlers

func (h *handler) CreateCertificationEntityLink(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	var payload struct {
		EntityType string `json:"entityType"`
		EntityID   string `json:"entityId"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	entityType := EntityType(payload.EntityType)
	entityID, err := uuid.Parse(payload.EntityID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid entity id format",
		})
	}

	if err := h.service.CreateCertificationEntityLink(c.Context(), userID, certificationID, entityType, entityID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, fiber.Map{
		"message": "entity linked to certification successfully",
	})
}

func (h *handler) GetCertificationEntityLinks(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	links, err := h.service.GetCertificationEntityLinks(c.Context(), certificationID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, links)
}

func (h *handler) GetEntityCertifications(c *fiber.Ctx) error {
	entityTypeStr := c.Params("entityType")
	entityIDStr := c.Params("entityId")

	entityType := EntityType(entityTypeStr)
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid entity id",
		})
	}

	certifications, err := h.service.GetEntityCertifications(c.Context(), entityType, entityID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range certifications {
			fieldMap := map[string]*string{
				"name":        &certifications[i].Name,
				"issuer":     &certifications[i].Issuer,
				"description": &certifications[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeCertification, certifications[i].ID, language, fieldMap)
		}
	}

	return response.Success(c, fiber.StatusOK, certifications)
}

func (h *handler) DeleteCertificationEntityLink(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	linkID, err := uuid.Parse(c.Params("linkId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid link id",
		})
	}

	if err := h.service.DeleteCertificationEntityLink(c.Context(), userID, linkID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "entity link deleted successfully",
	})
}

func (h *handler) DeleteCertificationEntityLinks(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	certificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid certification id",
		})
	}

	if err := h.service.DeleteCertificationEntityLinks(c.Context(), userID, certificationID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "all entity links deleted successfully",
	})
}

// Helper functions

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error in certification handler", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
			"message": "internal server error",
		})
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidName, ErrCodeInvalidIssuer, ErrCodeInvalidDate, ErrCodeInvalidStatus, ErrCodeInvalidCategory:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeUnauthorized:
		statusCode = fiber.StatusUnauthorized
	case ErrCodeConflict:
		statusCode = fiber.StatusConflict
	}

	return response.Error(c, statusCode, domainErr.Code, fiber.Map{
		"message": domainErr.Message,
	})
}

func parseDate(dateStr string) (time.Time, error) {
	// Try ISO 8601 format first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}
	// Try RFC3339 format
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}
	return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid date format")
}

