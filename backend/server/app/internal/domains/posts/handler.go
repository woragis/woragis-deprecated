package posts

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	creativeassets "github.com/woragis/backend/server/app/internal/domains/creativeassets"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationenricher "github.com/woragis/backend/server/app/pkg/translations"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes post endpoints.
type Handler interface {
	// Post handlers
	CreatePost(c *fiber.Ctx) error
	UpdatePost(c *fiber.Ctx) error
	GetPost(c *fiber.Ctx) error
	GetPostBySlug(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
	ListPosts(c *fiber.Ctx) error

	// Category handlers
	CreateCategory(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	GetCategory(c *fiber.Ctx) error
	GetCategoryBySlug(c *fiber.Ctx) error
	ListCategories(c *fiber.Ctx) error

	// Tag handlers
	GetTag(c *fiber.Ctx) error
	GetTagBySlug(c *fiber.Ctx) error
	ListTags(c *fiber.Ctx) error

	// Relationship handlers
	GetPostSkills(c *fiber.Ctx) error
	AttachSkillToPost(c *fiber.Ctx) error
	DetachSkillFromPost(c *fiber.Ctx) error
	GetPostCategories(c *fiber.Ctx) error
	AttachCategoryToPost(c *fiber.Ctx) error
	DetachCategoryFromPost(c *fiber.Ctx) error
	GetPostTags(c *fiber.Ctx) error
	AttachTagToPost(c *fiber.Ctx) error
	DetachTagFromPost(c *fiber.Ctx) error

	// Creative assets integration
	GeneratePostThumbnail(c *fiber.Ctx) error
	GeneratePostFeaturedImage(c *fiber.Ctx) error
	GeneratePostOGImage(c *fiber.Ctx) error
	GetPostAssets(c *fiber.Ctx) error
}

type handler struct {
	service          Service
	enricher         *translationenricher.Enricher
	translationService translationsdomain.Service
	creativeAssetsService creativeassets.Service
	logger           *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a post handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, creativeAssetsService creativeassets.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		creativeAssetsService: creativeAssetsService,
		logger:            logger,
	}
}

// Payloads

type createPostPayload struct {
	Title           string      `json:"title"`
	Content         string      `json:"content"`
	Excerpt         string      `json:"excerpt,omitempty"`
	Status          PostStatus  `json:"status,omitempty"`
	FeaturedImage   string      `json:"featuredImage,omitempty"`
	MetaTitle       string      `json:"metaTitle,omitempty"`
	MetaDescription string      `json:"metaDescription,omitempty"`
	MetaKeywords    string      `json:"metaKeywords,omitempty"`
	OGTitle         string      `json:"ogTitle,omitempty"`
	OGDescription   string      `json:"ogDescription,omitempty"`
	OGImage         string      `json:"ogImage,omitempty"`
	Featured        bool        `json:"featured,omitempty"`
	SkillIDs        []uuid.UUID `json:"skillIds,omitempty"`
	CategoryIDs     []uuid.UUID `json:"categoryIds,omitempty"`
	TagNames        []string    `json:"tagNames,omitempty"`
}

type updatePostPayload struct {
	Title           *string     `json:"title,omitempty"`
	Content         *string     `json:"content,omitempty"`
	Excerpt         *string     `json:"excerpt,omitempty"`
	Status          *PostStatus `json:"status,omitempty"`
	FeaturedImage   *string     `json:"featuredImage,omitempty"`
	MetaTitle       *string     `json:"metaTitle,omitempty"`
	MetaDescription *string     `json:"metaDescription,omitempty"`
	MetaKeywords    *string     `json:"metaKeywords,omitempty"`
	OGTitle         *string     `json:"ogTitle,omitempty"`
	OGDescription   *string     `json:"ogDescription,omitempty"`
	OGImage         *string     `json:"ogImage,omitempty"`
	Featured        *bool       `json:"featured,omitempty"`
	SkillIDs        []uuid.UUID `json:"skillIds,omitempty"`
	CategoryIDs     []uuid.UUID `json:"categoryIds,omitempty"`
	TagNames        []string    `json:"tagNames,omitempty"`
}

type createCategoryPayload struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type updateCategoryPayload struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Handlers

func (h *handler) CreatePost(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createPostPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.CreatePost(c.Context(), userID, CreatePostRequest{
		Title:           payload.Title,
		Content:         payload.Content,
		Excerpt:         payload.Excerpt,
		Status:          payload.Status,
		FeaturedImage:   payload.FeaturedImage,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		MetaKeywords:    payload.MetaKeywords,
		OGTitle:         payload.OGTitle,
		OGDescription:   payload.OGDescription,
		OGImage:         payload.OGImage,
		Featured:        payload.Featured,
		SkillIDs:        payload.SkillIDs,
		CategoryIDs:     payload.CategoryIDs,
		TagNames:        payload.TagNames,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Automatically trigger translations for all supported languages
	if h.translationService != nil {
		// Prepare source text for translation
		sourceText := make(map[string]string)
		if post.Title != "" {
			sourceText["title"] = post.Title
		}
		if post.Content != "" {
			sourceText["content"] = post.Content
		}
		if post.Excerpt != "" {
			sourceText["excerpt"] = post.Excerpt
		}
		if post.MetaTitle != "" {
			sourceText["metaTitle"] = post.MetaTitle
		}
		if post.MetaDescription != "" {
			sourceText["metaDescription"] = post.MetaDescription
		}
		if post.OGTitle != "" {
			sourceText["ogTitle"] = post.OGTitle
		}
		if post.OGDescription != "" {
			sourceText["ogDescription"] = post.OGDescription
		}

		// Fields to translate
		fields := []string{}
		if post.Title != "" {
			fields = append(fields, "title")
		}
		if post.Content != "" {
			fields = append(fields, "content")
		}
		if post.Excerpt != "" {
			fields = append(fields, "excerpt")
		}
		if post.MetaTitle != "" {
			fields = append(fields, "metaTitle")
		}
		if post.MetaDescription != "" {
			fields = append(fields, "metaDescription")
		}
		if post.OGTitle != "" {
			fields = append(fields, "ogTitle")
		}
		if post.OGDescription != "" {
			fields = append(fields, "ogDescription")
		}

		// Queue translations for all supported languages (except English)
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

		// Trigger translations asynchronously (don't block the response)
		// Use background context to avoid cancellation when request completes
		go func() {
			ctx := context.Background()
			for _, lang := range supportedLanguages {
				if err := h.translationService.RequestTranslation(
					ctx,
					translationsdomain.EntityTypePost,
					post.ID,
					lang,
					fields,
					sourceText,
				); err != nil {
					h.logger.Warn("Failed to queue translation",
						slog.String("postId", post.ID.String()),
						slog.String("language", string(lang)),
						slog.Any("error", err),
					)
				}
			}
		}()
	}

	return response.Success(c, fiber.StatusCreated, toPostResponse(post))
}

func (h *handler) UpdatePost(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updatePostPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.UpdatePost(c.Context(), userID, postID, UpdatePostRequest{
		Title:           payload.Title,
		Content:         payload.Content,
		Excerpt:         payload.Excerpt,
		Status:          payload.Status,
		FeaturedImage:   payload.FeaturedImage,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		MetaKeywords:    payload.MetaKeywords,
		OGTitle:         payload.OGTitle,
		OGDescription:   payload.OGDescription,
		OGImage:         payload.OGImage,
		Featured:        payload.Featured,
		SkillIDs:        payload.SkillIDs,
		CategoryIDs:     payload.CategoryIDs,
		TagNames:        payload.TagNames,
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
			"title":          &post.Title,
			"content":        &post.Content,
			"excerpt":        &post.Excerpt,
			"metaTitle":      &post.MetaTitle,
			"metaDescription": &post.MetaDescription,
			"ogTitle":        &post.OGTitle,
			"ogDescription":  &post.OGDescription,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypePost, post.ID, language, fieldMap)
	}

	// Increment views for published posts
	if post.Status == PostStatusPublished {
		go h.service.IncrementPostViews(c.Context(), postID)
	}

	return response.Success(c, fiber.StatusOK, toPostResponse(post))
}

func (h *handler) GetPostBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.GetPostBySlug(c.Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"title":          &post.Title,
			"content":        &post.Content,
			"excerpt":        &post.Excerpt,
			"metaTitle":      &post.MetaTitle,
			"metaDescription": &post.MetaDescription,
			"ogTitle":        &post.OGTitle,
			"ogDescription":  &post.OGDescription,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypePost, post.ID, language, fieldMap)
	}

	// Increment views for published posts
	if post.Status == PostStatusPublished {
		go h.service.IncrementPostViews(c.Context(), post.ID)
	}

	return response.Success(c, fiber.StatusOK, toPostResponse(post))
}

func (h *handler) DeletePost(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeletePost(c.Context(), userID, postID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "post deleted"})
}

func (h *handler) ListPosts(c *fiber.Ctx) error {
	filters := PostFilters{}

	// Get user ID if authenticated (for filtering own posts)
	if userID, err := authdomain.UserIDFromContext(c); err == nil {
		filters.UserID = &userID
	}

	// Query parameters
	if statusStr := c.Query("status"); statusStr != "" {
		status := PostStatus(statusStr)
		filters.Status = &status
	} else {
		// Default to published for public access
		published := PostStatusPublished
		filters.Status = &published
	}

	if featuredStr := c.Query("featured"); featuredStr != "" {
		featured := featuredStr == "true"
		filters.Featured = &featured
	}

	if categoryIDStr := c.Query("categoryId"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			filters.CategoryID = &categoryID
		}
	}

	if tagIDStr := c.Query("tagId"); tagIDStr != "" {
		if tagID, err := uuid.Parse(tagIDStr); err == nil {
			filters.TagID = &tagID
		}
	}

	if skillIDStr := c.Query("skillId"); skillIDStr != "" {
		if skillID, err := uuid.Parse(skillIDStr); err == nil {
			filters.SkillID = &skillID
		}
	}

	if search := c.Query("search"); search != "" {
		filters.Search = search
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
	} else {
		filters.OrderBy = "created_at"
	}

	if order := c.Query("order"); order != "" {
		filters.Order = order
	} else {
		filters.Order = "desc"
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
				"title":          &posts[i].Title,
				"content":        &posts[i].Content,
				"excerpt":        &posts[i].Excerpt,
				"metaTitle":      &posts[i].MetaTitle,
				"metaDescription": &posts[i].MetaDescription,
				"ogTitle":        &posts[i].OGTitle,
				"ogDescription":  &posts[i].OGDescription,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypePost, posts[i].ID, language, fieldMap)
		}
	}

	responses := make([]postResponse, len(posts))
	for i := range posts {
		responses[i] = toPostResponse(&posts[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// Category handlers

func (h *handler) CreateCategory(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createCategoryPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	category, err := h.service.CreateCategory(c.Context(), CreateCategoryRequest{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toCategoryResponse(category))
}

func (h *handler) UpdateCategory(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateCategoryPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	category, err := h.service.UpdateCategory(c.Context(), categoryID, UpdateCategoryRequest{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toCategoryResponse(category))
}

func (h *handler) GetCategory(c *fiber.Ctx) error {
	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	category, err := h.service.GetCategory(c.Context(), categoryID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toCategoryResponse(category))
}

func (h *handler) GetCategoryBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	category, err := h.service.GetCategoryBySlug(c.Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toCategoryResponse(category))
}

func (h *handler) ListCategories(c *fiber.Ctx) error {
	categories, err := h.service.ListCategories(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]categoryResponse, len(categories))
	for i := range categories {
		responses[i] = toCategoryResponse(&categories[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// Tag handlers

func (h *handler) GetTag(c *fiber.Ctx) error {
	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	tag, err := h.service.GetTag(c.Context(), tagID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toTagResponse(tag))
}

func (h *handler) GetTagBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	tag, err := h.service.GetTagBySlug(c.Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toTagResponse(tag))
}

func (h *handler) ListTags(c *fiber.Ctx) error {
	tags, err := h.service.ListTags(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]tagResponse, len(tags))
	for i := range tags {
		responses[i] = toTagResponse(&tags[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// Relationship handlers

func (h *handler) GetPostSkills(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skillIDs, err := h.service.GetPostSkills(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"skillIds": skillIDs})
}

func (h *handler) AttachSkillToPost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload struct {
		SkillID uuid.UUID `json:"skillId"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.AttachSkillToPost(c.Context(), postID, payload.SkillID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "skill attached"})
}

func (h *handler) DetachSkillFromPost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	skillID, err := uuid.Parse(c.Params("skillId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DetachSkillFromPost(c.Context(), postID, skillID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "skill detached"})
}

func (h *handler) GetPostCategories(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	categories, err := h.service.GetPostCategories(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]categoryResponse, len(categories))
	for i := range categories {
		responses[i] = toCategoryResponse(&categories[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) AttachCategoryToPost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload struct {
		CategoryID uuid.UUID `json:"categoryId"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.AttachCategoryToPost(c.Context(), postID, payload.CategoryID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "category attached"})
}

func (h *handler) DetachCategoryFromPost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	categoryID, err := uuid.Parse(c.Params("categoryId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DetachCategoryFromPost(c.Context(), postID, categoryID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "category detached"})
}

func (h *handler) GetPostTags(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	tags, err := h.service.GetPostTags(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]tagResponse, len(tags))
	for i := range tags {
		responses[i] = toTagResponse(&tags[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) AttachTagToPost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload struct {
		TagID uuid.UUID `json:"tagId"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.AttachTagToPost(c.Context(), postID, payload.TagID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "tag attached"})
}

func (h *handler) DetachTagFromPost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	tagID, err := uuid.Parse(c.Params("tagId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DetachTagFromPost(c.Context(), postID, tagID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "tag detached"})
}

// Creative assets integration

type generateThumbnailPayload struct {
	Prompt  string `json:"prompt"`
	Context string `json:"context,omitempty"`
}

func (h *handler) GeneratePostThumbnail(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Verify post ownership
	post, err := h.service.GetPost(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}
	if post.UserID != userID {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "unauthorized",
		})
	}

	var payload generateThumbnailPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if payload.Prompt == "" {
		payload.Prompt = post.Title // Use post title as default prompt
	}

	asset, err := h.creativeAssetsService.GenerateAndStoreThumbnail(
		c.Context(),
		userID,
		creativeassets.EntityTypePost,
		postID,
		payload.Prompt,
		payload.Context,
	)
	if err != nil {
		h.logger.Error("failed to generate thumbnail", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to generate thumbnail",
		})
	}

	// Update post with asset URL if we have one
	assetURL := fmt.Sprintf("/api/v1/creative-assets/%s/data", asset.ID.String())
	updateReq := UpdatePostRequest{FeaturedImage: &assetURL}
	h.service.UpdatePost(c.Context(), userID, postID, updateReq)

	return response.Success(c, fiber.StatusCreated, asset)
}

func (h *handler) GeneratePostFeaturedImage(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.GetPost(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}
	if post.UserID != userID {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "unauthorized",
		})
	}

	var payload generateThumbnailPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if payload.Prompt == "" {
		payload.Prompt = post.Title
	}

	asset, err := h.creativeAssetsService.GenerateAndStoreImage(
		c.Context(),
		userID,
		creativeassets.EntityTypePost,
		postID,
		creativeassets.PurposeFeaturedImage,
		payload.Prompt,
		payload.Context,
	)
	if err != nil {
		h.logger.Error("failed to generate featured image", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to generate featured image",
		})
	}

	assetURL := fmt.Sprintf("/api/v1/creative-assets/%s/data", asset.ID.String())
	updateReq := UpdatePostRequest{FeaturedImage: &assetURL}
	h.service.UpdatePost(c.Context(), userID, postID, updateReq)

	return response.Success(c, fiber.StatusCreated, asset)
}

func (h *handler) GeneratePostOGImage(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	post, err := h.service.GetPost(c.Context(), postID)
	if err != nil {
		return h.handleError(c, err)
	}
	if post.UserID != userID {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "unauthorized",
		})
	}

	var payload generateThumbnailPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if payload.Prompt == "" {
		payload.Prompt = post.Title
	}

	asset, err := h.creativeAssetsService.GenerateAndStoreImage(
		c.Context(),
		userID,
		creativeassets.EntityTypePost,
		postID,
		creativeassets.PurposeOGImage,
		payload.Prompt,
		payload.Context,
	)
	if err != nil {
		h.logger.Error("failed to generate OG image", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to generate OG image",
		})
	}

	assetURL := fmt.Sprintf("/api/v1/creative-assets/%s/data", asset.ID.String())
	updateReq := UpdatePostRequest{OGImage: &assetURL}
	h.service.UpdatePost(c.Context(), userID, postID, updateReq)

	return response.Success(c, fiber.StatusCreated, asset)
}

func (h *handler) GetPostAssets(c *fiber.Ctx) error {
	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	assets, err := h.creativeAssetsService.GetAssetsByEntity(
		c.Context(),
		creativeassets.EntityTypePost,
		postID,
	)
	if err != nil {
		h.logger.Error("failed to get post assets", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to get assets",
		})
	}

	return response.Success(c, fiber.StatusOK, assets)
}

// Response helpers

type postResponse struct {
	ID              string     `json:"id"`
	UserID          string     `json:"userId"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Content         string     `json:"content"`
	Excerpt         string     `json:"excerpt,omitempty"`
	Status          PostStatus `json:"status"`
	PublishedAt     *string    `json:"publishedAt,omitempty"`
	FeaturedImage   string     `json:"featuredImage,omitempty"`
	MetaTitle       string     `json:"metaTitle,omitempty"`
	MetaDescription string     `json:"metaDescription,omitempty"`
	MetaKeywords    string     `json:"metaKeywords,omitempty"`
	OGTitle         string     `json:"ogTitle,omitempty"`
	OGDescription   string     `json:"ogDescription,omitempty"`
	OGImage         string     `json:"ogImage,omitempty"`
	Featured        bool       `json:"featured"`
	ViewsCount      int64      `json:"viewsCount"`
	CreatedAt       string     `json:"createdAt"`
	UpdatedAt       string     `json:"updatedAt"`
}

func toPostResponse(post *Post) postResponse {
	resp := postResponse{
		ID:              post.ID.String(),
		UserID:          post.UserID.String(),
		Title:           post.Title,
		Slug:            post.Slug,
		Content:         post.Content,
		Excerpt:         post.Excerpt,
		Status:          post.Status,
		FeaturedImage:   post.FeaturedImage,
		MetaTitle:       post.MetaTitle,
		MetaDescription: post.MetaDescription,
		MetaKeywords:    post.MetaKeywords,
		OGTitle:         post.OGTitle,
		OGDescription:   post.OGDescription,
		OGImage:         post.OGImage,
		Featured:        post.Featured,
		ViewsCount:      post.ViewsCount,
		CreatedAt:       post.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       post.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if post.PublishedAt != nil {
		publishedAt := post.PublishedAt.Format("2006-01-02T15:04:05Z")
		resp.PublishedAt = &publishedAt
	}

	return resp
}

type categoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func toCategoryResponse(category *Category) categoryResponse {
	return categoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type tagResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func toTagResponse(tag *Tag) tagResponse {
	return tagResponse{
		ID:        tag.ID.String(),
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: tag.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		statusCode := fiber.StatusInternalServerError
		switch domainErr.Code {
		case ErrCodePostNotFound, ErrCodeCategoryNotFound, ErrCodeTagNotFound:
			statusCode = fiber.StatusNotFound
		case ErrCodeInvalidPayload, ErrCodeInvalidTitle, ErrCodeInvalidContent, ErrCodeInvalidStatus:
			statusCode = fiber.StatusBadRequest
		case ErrCodeUnauthorized:
			statusCode = fiber.StatusUnauthorized
		case ErrCodeDuplicateSlug:
			statusCode = fiber.StatusConflict
		}
		return response.Error(c, statusCode, domainErr.Code, fiber.Map{
			"message": domainErr.Message,
		})
	}

	h.logger.Error("unexpected error in post handler", "error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
		"message": "internal server error",
	})
}

