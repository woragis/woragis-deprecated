package socialmediaposts

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationenricher "github.com/woragis/backend/server/app/pkg/translations"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes social media post endpoints.
type Handler interface {
	// Post operations
	CreatePost(c *fiber.Ctx) error
	UpdatePost(c *fiber.Ctx) error
	GetPost(c *fiber.Ctx) error
	GetPostByURL(c *fiber.Ctx) error
	ListPosts(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
	UpdatePostEngagement(c *fiber.Ctx) error

	// Link operations
	CreateLink(c *fiber.Ctx) error
	UpdateLink(c *fiber.Ctx) error
	DeleteLink(c *fiber.Ctx) error
	GetLinksByPost(c *fiber.Ctx) error
	GetLinksByEntity(c *fiber.Ctx) error
	GetPostsByEntity(c *fiber.Ctx) error
	GetEntitiesByPost(c *fiber.Ctx) error
}

type handler struct {
	service          Service
	enricher         *translationenricher.Enricher
	translationService translationsdomain.Service
	logger           *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a social media post handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createPostPayload struct {
	URL            string     `json:"url"`
	Platform       Platform   `json:"platform"`
	Title          string     `json:"title,omitempty"`
	ContentPreview string     `json:"contentPreview,omitempty"`
	PublishedDate  *time.Time `json:"publishedDate,omitempty"`
}

type updatePostPayload struct {
	Title          string      `json:"title,omitempty"`
	ContentPreview string      `json:"contentPreview,omitempty"`
	PublishedDate  *time.Time  `json:"publishedDate,omitempty"`
	Status         *PostStatus `json:"status,omitempty"`
}

type updateEngagementPayload struct {
	Likes    *int64 `json:"likes,omitempty"`
	Shares   *int64 `json:"shares,omitempty"`
	Comments *int64 `json:"comments,omitempty"`
	Views    *int64 `json:"views,omitempty"`
}

type createLinkPayload struct {
	PostID           uuid.UUID        `json:"postId"`
	EntityType       EntityType       `json:"entityType"`
	EntityID         uuid.UUID        `json:"entityId"`
	RelationshipType RelationshipType `json:"relationshipType"`
}

type updateLinkPayload struct {
	RelationshipType RelationshipType `json:"relationshipType"`
}

// Post handlers

func (h *handler) CreatePost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload createPostPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.CreatePost(c.Context(), CreatePostRequest{
		URL:            payload.URL,
		Platform:       payload.Platform,
		Title:          payload.Title,
		ContentPreview: payload.ContentPreview,
		PublishedDate:  payload.PublishedDate,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toPostResponse(post))
}

func (h *handler) UpdatePost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updatePostPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.UpdatePost(c.Context(), UpdatePostRequest{
		PostID:         postID,
		Title:          payload.Title,
		ContentPreview: payload.ContentPreview,
		PublishedDate:  payload.PublishedDate,
		Status:         payload.Status,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toPostResponse(post))
}

func (h *handler) GetPost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.GetPost(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"title":         &post.Title,
			"contentPreview": &post.ContentPreview,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSocialMediaPost, post.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toPostResponse(post))
}

func (h *handler) GetPostByURL(c *fiber.Ctx) error {
	url := c.Query("url")
	if url == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.GetPostByURL(c.Context(), url)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"title":         &post.Title,
			"contentPreview": &post.ContentPreview,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSocialMediaPost, post.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toPostResponse(post))
}

func (h *handler) ListPosts(c *fiber.Ctx) error {
	filters := PostFilters{}

	if platformStr := c.Query("platform"); platformStr != "" {
		platform := Platform(platformStr)
		filters.Platform = &platform
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := PostStatus(statusStr)
		filters.Status = &status
	}

	posts, err := h.service.ListPosts(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range posts {
			fieldMap := map[string]*string{
				"title":         &posts[i].Title,
				"contentPreview": &posts[i].ContentPreview,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSocialMediaPost, posts[i].ID, language, fieldMap)
		}
	}

	responses := make([]postResponse, len(posts))
	for i := range posts {
		responses[i] = toPostResponse(&posts[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) DeletePost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeletePost(c.Context(), postID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "post deleted"})
}

func (h *handler) UpdatePostEngagement(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateEngagementPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.UpdatePostEngagement(c.Context(), postID, UpdateEngagementRequest{
		Likes:    payload.Likes,
		Shares:   payload.Shares,
		Comments: payload.Comments,
		Views:    payload.Views,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toPostResponse(post))
}

// Link handlers

func (h *handler) CreateLink(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload createLinkPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	link, err := h.service.CreateLink(c.Context(), CreateLinkRequest{
		PostID:           payload.PostID,
		EntityType:       payload.EntityType,
		EntityID:         payload.EntityID,
		RelationshipType: payload.RelationshipType,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toLinkResponse(link))
}

func (h *handler) UpdateLink(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	linkID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateLinkPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	link, err := h.service.UpdateLink(c.Context(), UpdateLinkRequest{
		LinkID:           linkID,
		RelationshipType: payload.RelationshipType,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toLinkResponse(link))
}

func (h *handler) DeleteLink(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	linkID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteLink(c.Context(), linkID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "link deleted"})
}

func (h *handler) GetLinksByPost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("postId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	links, err := h.service.GetLinksByPost(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]linkResponse, len(links))
	for i := range links {
		responses[i] = toLinkResponse(&links[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetLinksByEntity(c *fiber.Ctx) error {
	entityType := EntityType(c.Query("entityType"))
	entityIDStr := c.Query("entityId")

	if entityType == "" || entityIDStr == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	links, err := h.service.GetLinksByEntity(c.Context(), entityType, entityID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]linkResponse, len(links))
	for i := range links {
		responses[i] = toLinkResponse(&links[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetPostsByEntity(c *fiber.Ctx) error {
	entityType := EntityType(c.Query("entityType"))
	entityIDStr := c.Query("entityId")

	if entityType == "" || entityIDStr == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var relationshipType *RelationshipType
	if rtStr := c.Query("relationshipType"); rtStr != "" {
		rt := RelationshipType(rtStr)
		relationshipType = &rt
	}

	posts, err := h.service.GetPostsByEntity(c.Context(), entityType, entityID, relationshipType)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]postWithLinkResponse, len(posts))
	for i := range posts {
		responses[i] = toPostWithLinkResponse(posts[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetEntitiesByPost(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("postId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	entities, err := h.service.GetEntitiesByPost(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, entities)
}

// Response helpers

type postResponse struct {
	ID             string     `json:"id"`
	URL            string     `json:"url"`
	Platform       Platform   `json:"platform"`
	Title          string     `json:"title,omitempty"`
	ContentPreview string     `json:"contentPreview,omitempty"`
	PublishedDate  *string    `json:"publishedDate,omitempty"`
	Likes          *int64     `json:"likes,omitempty"`
	Shares         *int64     `json:"shares,omitempty"`
	Comments       *int64     `json:"comments,omitempty"`
	Views          *int64     `json:"views,omitempty"`
	Status         PostStatus `json:"status"`
	CreatedAt      string     `json:"createdAt"`
	UpdatedAt      string     `json:"updatedAt"`
}

type linkResponse struct {
	ID               string           `json:"id"`
	SocialMediaPostID string          `json:"socialMediaPostId"`
	EntityType       EntityType       `json:"entityType"`
	EntityID         string           `json:"entityId"`
	RelationshipType RelationshipType `json:"relationshipType"`
	CreatedAt        string           `json:"createdAt"`
	UpdatedAt        string           `json:"updatedAt"`
}

type postWithLinkResponse struct {
	postResponse
	RelationshipType RelationshipType `json:"relationshipType"`
	LinkID           string           `json:"linkId"`
}

func toPostResponse(post *SocialMediaPost) postResponse {
	var publishedDate *string
	if post.PublishedDate != nil {
		dateStr := post.PublishedDate.Format("2006-01-02T15:04:05Z07:00")
		publishedDate = &dateStr
	}

	return postResponse{
		ID:             post.ID.String(),
		URL:            post.URL,
		Platform:       post.Platform,
		Title:          post.Title,
		ContentPreview: post.ContentPreview,
		PublishedDate:  publishedDate,
		Likes:          post.Likes,
		Shares:         post.Shares,
		Comments:       post.Comments,
		Views:          post.Views,
		Status:         post.Status,
		CreatedAt:      post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toLinkResponse(link *SocialMediaEntityLink) linkResponse {
	return linkResponse{
		ID:               link.ID.String(),
		SocialMediaPostID: link.SocialMediaPostID.String(),
		EntityType:       link.EntityType,
		EntityID:         link.EntityID.String(),
		RelationshipType: link.RelationshipType,
		CreatedAt:        link.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        link.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPostWithLinkResponse(post SocialMediaPostWithLink) postWithLinkResponse {
	baseResponse := toPostResponse(&post.SocialMediaPost)
	return postWithLinkResponse{
		postResponse:     baseResponse,
		RelationshipType: post.RelationshipType,
		LinkID:           post.LinkID.String(),
	}
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, nil)
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidPlatform, ErrCodeInvalidEntityType,
		ErrCodeInvalidRelationshipType, ErrCodeInvalidStatus:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeConflict:
		statusCode = fiber.StatusConflict
	case ErrCodeRepositoryFailure:
		statusCode = fiber.StatusInternalServerError
	}

	return response.Error(c, statusCode, domainErr.Code, fiber.Map{"message": domainErr.Message})
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "unauthorized"})
}

