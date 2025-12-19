package content

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Platform represents the social media platform (duplicated from parent to avoid import cycle).
type Platform string

const (
	PlatformLinkedIn  Platform = "linkedin"
	PlatformTwitter   Platform = "twitter"
	PlatformInstagram Platform = "instagram"
	PlatformMedium    Platform = "medium"
	PlatformSubstack  Platform = "substack"
	PlatformValete    Platform = "valete"
	PlatformWebsite   Platform = "website"
)

// ContentFormat represents the format of the social media post content (duplicated from parent to avoid import cycle).
type ContentFormat string

const (
	FormatLongForm   ContentFormat = "long-form"
	FormatThread     ContentFormat = "thread"
	FormatCarousel   ContentFormat = "carousel"
	FormatArticle    ContentFormat = "article"
	FormatNewsletter ContentFormat = "newsletter"
	FormatPost       ContentFormat = "post"
)

// Handler exposes content post endpoints.
type Handler interface {
	CreateContentPost(c *fiber.Ctx) error
	GetContentPost(c *fiber.Ctx) error
	ListContentPosts(c *fiber.Ctx) error
	GetContentBacklog(c *fiber.Ctx) error
	UpdateContentPostPriority(c *fiber.Ctx) error
	RepurposeToPlatforms(c *fiber.Ctx) error
	GetRepurposingHistory(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a content handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createContentPostPayload struct {
	PostID      string  `json:"postId"`
	ContentType string  `json:"contentType,omitempty"`
	Project     *string `json:"project,omitempty"`
	Priority    string  `json:"priority,omitempty"`
}

type repurposePayload struct {
	Platforms []repurposePlatformPayload `json:"platforms"`
}

type repurposePlatformPayload struct {
	Platform string `json:"platform"`
	Format   string `json:"format"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type updatePriorityPayload struct {
	Priority string `json:"priority"`
}

// Handlers

func (h *handler) CreateContentPost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload createContentPostPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	postID, err := uuid.Parse(payload.PostID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var priority ContentPostPriority
	if payload.Priority != "" {
		priority = ContentPostPriority(payload.Priority)
	}

	post, err := h.service.CreateContentPostFromBackend(c.Context(), CreateContentPostRequest{
		PostID:      postID,
		ContentType: payload.ContentType,
		Project:     payload.Project,
		Priority:    priority,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toContentPostResponse(post))
}

func (h *handler) GetContentPost(c *fiber.Ctx) error {
	contentPostID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.GetContentPost(c.Context(), contentPostID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toContentPostWithSocialPostsResponse(post))
}

func (h *handler) ListContentPosts(c *fiber.Ctx) error {
	filters := ContentPostFilters{}

	if statusStr := c.Query("status"); statusStr != "" {
		status := ContentPostStatus(statusStr)
		filters.Status = &status
	}

	if priorityStr := c.Query("priority"); priorityStr != "" {
		priority := ContentPostPriority(priorityStr)
		filters.Priority = &priority
	}

	posts, err := h.service.ListContentPosts(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]contentPostResponse, len(posts))
	for i := range posts {
		responses[i] = toContentPostResponse(&posts[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetContentBacklog(c *fiber.Ctx) error {
	posts, err := h.service.GetContentBacklog(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]contentPostResponse, len(posts))
	for i := range posts {
		responses[i] = toContentPostResponse(&posts[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) UpdateContentPostPriority(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	contentPostID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updatePriorityPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	priority := ContentPostPriority(payload.Priority)
	post, err := h.service.UpdateContentPostPriority(c.Context(), contentPostID, priority)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toContentPostResponse(post))
}

func (h *handler) RepurposeToPlatforms(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	contentPostID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload repurposePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Convert payload to service request
	platforms := make([]RepurposePlatform, len(payload.Platforms))
	for i, p := range payload.Platforms {
		platforms[i] = RepurposePlatform{
			Platform: Platform(p.Platform),
			Format:   ContentFormat(p.Format),
			Title:    p.Title,
			Content:  p.Content,
		}
	}

	posts, err := h.service.RepurposeToPlatforms(c.Context(), contentPostID, RepurposeRequest{
		Platforms: platforms,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Convert to response format
	responses := make([]interface{}, len(posts))
	for i, post := range posts {
		responses[i] = map[string]interface{}{
			"id":       post.ID.String(),
			"platform": post.Platform,
			"format":   post.Format,
			"title":    post.Title,
			"status":   post.Status,
		}
	}

	return response.Success(c, fiber.StatusCreated, responses)
}

func (h *handler) GetRepurposingHistory(c *fiber.Ctx) error {
	contentPostID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	history, err := h.service.GetRepurposingHistory(c.Context(), contentPostID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]interface{}, len(history))
	for i, item := range history {
		responses[i] = map[string]interface{}{
			"socialPost": map[string]interface{}{
				"id":       item.SocialPost.ID.String(),
				"platform": item.SocialPost.Platform,
				"format":   item.SocialPost.Format,
				"title":    item.SocialPost.Title,
				"status":   item.SocialPost.Status,
			},
			"repurposedAt": item.RepurposedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// Response helpers

type contentPostResponse struct {
	ID          string             `json:"id"`
	PostID      string              `json:"postId"`
	ContentType string              `json:"contentType,omitempty"`
	Project     *string             `json:"project,omitempty"`
	Priority    ContentPostPriority `json:"priority"`
	Status      ContentPostStatus   `json:"status"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

func toContentPostResponse(post *ContentPost) contentPostResponse {
	return contentPostResponse{
		ID:          post.ID.String(),
		PostID:      post.PostID.String(),
		ContentType: post.ContentType,
		Project:     post.Project,
		Priority:    post.Priority,
		Status:      post.Status,
		CreatedAt:   post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toContentPostWithSocialPostsResponse(post *ContentPostWithSocialPosts) interface{} {
	baseResponse := toContentPostResponse(&post.ContentPost)
	socialPosts := make([]interface{}, len(post.SocialPosts))
	for i, sp := range post.SocialPosts {
		socialPosts[i] = map[string]interface{}{
			"id":       sp.ID.String(),
			"platform": sp.Platform,
			"format":   sp.Format,
			"title":    sp.Title,
			"status":   sp.Status,
		}
	}

	return map[string]interface{}{
		"contentPost": baseResponse,
		"socialPosts": socialPosts,
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
	case ErrCodeInvalidPayload, ErrCodeInvalidStatus, ErrCodeInvalidPriority:
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
