package posts

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates post workflows.
type Service interface {
	// Post operations
	CreatePost(ctx context.Context, userID uuid.UUID, req CreatePostRequest) (*Post, error)
	UpdatePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID, req UpdatePostRequest) (*Post, error)
	GetPost(ctx context.Context, postID uuid.UUID) (*Post, error)
	GetPostBySlug(ctx context.Context, slug string) (*Post, error)
	DeletePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error
	ListPosts(ctx context.Context, filters PostFilters) ([]Post, error)
	IncrementPostViews(ctx context.Context, postID uuid.UUID) error

	// Category operations
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*Category, error)
	UpdateCategory(ctx context.Context, categoryID uuid.UUID, req UpdateCategoryRequest) (*Category, error)
	GetCategory(ctx context.Context, categoryID uuid.UUID) (*Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*Category, error)
	ListCategories(ctx context.Context) ([]Category, error)

	// Tag operations
	GetOrCreateTag(ctx context.Context, name string) (*Tag, error)
	GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error)
	GetTagBySlug(ctx context.Context, slug string) (*Tag, error)
	ListTags(ctx context.Context) ([]Tag, error)

	// Post-Skill relationship operations
	AttachSkillToPost(ctx context.Context, postID, skillID uuid.UUID) error
	DetachSkillFromPost(ctx context.Context, postID, skillID uuid.UUID) error
	GetPostSkills(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error)

	// Post-Category relationship operations
	AttachCategoryToPost(ctx context.Context, postID, categoryID uuid.UUID) error
	DetachCategoryFromPost(ctx context.Context, postID, categoryID uuid.UUID) error
	GetPostCategories(ctx context.Context, postID uuid.UUID) ([]Category, error)

	// Post-Tag relationship operations
	AttachTagToPost(ctx context.Context, postID, tagID uuid.UUID) error
	DetachTagFromPost(ctx context.Context, postID, tagID uuid.UUID) error
	GetPostTags(ctx context.Context, postID uuid.UUID) ([]Tag, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Request payloads

type CreatePostRequest struct {
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Excerpt         string     `json:"excerpt,omitempty"`
	Status          PostStatus `json:"status,omitempty"`
	FeaturedImage   string     `json:"featuredImage,omitempty"`
	MetaTitle       string     `json:"metaTitle,omitempty"`
	MetaDescription string     `json:"metaDescription,omitempty"`
	MetaKeywords    string     `json:"metaKeywords,omitempty"`
	OGTitle         string     `json:"ogTitle,omitempty"`
	OGDescription   string     `json:"ogDescription,omitempty"`
	OGImage         string     `json:"ogImage,omitempty"`
	Featured        bool       `json:"featured,omitempty"`
	SkillIDs        []uuid.UUID `json:"skillIds,omitempty"`
	CategoryIDs     []uuid.UUID `json:"categoryIds,omitempty"`
	TagNames        []string    `json:"tagNames,omitempty"`
}

type UpdatePostRequest struct {
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

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Post operations

func (s *service) CreatePost(ctx context.Context, userID uuid.UUID, req CreatePostRequest) (*Post, error) {
	status := req.Status
	if status == "" {
		status = PostStatusDraft
	}

	post, err := NewPost(userID, req.Title, req.Content, req.Excerpt, status)
	if err != nil {
		return nil, err
	}

	// Check if slug is taken
	taken, err := s.repo.IsPostSlugTaken(ctx, post.Slug, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if taken {
		// Append UUID to make it unique
		post.Slug = post.Slug + "-" + post.ID.String()[:8]
	}

	// Set optional fields
	if req.FeaturedImage != "" {
		post.FeaturedImage = req.FeaturedImage
	}
	if req.MetaTitle != "" || req.MetaDescription != "" || req.MetaKeywords != "" ||
		req.OGTitle != "" || req.OGDescription != "" || req.OGImage != "" {
		post.UpdateSEO(req.MetaTitle, req.MetaDescription, req.MetaKeywords,
			req.OGTitle, req.OGDescription, req.OGImage)
	}
	if req.Featured {
		post.SetFeatured(true)
	}

	// Create post
	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, err
	}

	// Attach relationships
	for _, skillID := range req.SkillIDs {
		if err := s.repo.AttachSkillToPost(ctx, post.ID, skillID); err != nil {
			s.logger.Warn("Failed to attach skill to post", "error", err, "postID", post.ID, "skillID", skillID)
		}
	}

	for _, categoryID := range req.CategoryIDs {
		if err := s.repo.AttachCategoryToPost(ctx, post.ID, categoryID); err != nil {
			s.logger.Warn("Failed to attach category to post", "error", err, "postID", post.ID, "categoryID", categoryID)
		}
	}

	for _, tagName := range req.TagNames {
		tag, err := s.repo.GetOrCreateTag(ctx, tagName)
		if err != nil {
			s.logger.Warn("Failed to get or create tag", "error", err, "tagName", tagName)
			continue
		}
		if err := s.repo.AttachTagToPost(ctx, post.ID, tag.ID); err != nil {
			s.logger.Warn("Failed to attach tag to post", "error", err, "postID", post.ID, "tagID", tag.ID)
		}
	}

	return post, nil
}

func (s *service) UpdatePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID, req UpdatePostRequest) (*Post, error) {
	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if post.UserID != userID {
		return nil, NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	// Update fields
	if req.Title != nil {
		if err := post.UpdateTitle(*req.Title); err != nil {
			return nil, err
		}
		// Check if new slug is taken
		taken, err := s.repo.IsPostSlugTaken(ctx, post.Slug, postID)
		if err != nil {
			return nil, err
		}
		if taken {
			post.Slug = post.Slug + "-" + post.ID.String()[:8]
		}
	}

	if req.Content != nil {
		if err := post.UpdateContent(*req.Content); err != nil {
			return nil, err
		}
	}

	if req.Excerpt != nil {
		post.Excerpt = *req.Excerpt
		post.UpdatedAt = time.Now().UTC()
	}

	if req.Status != nil {
		switch *req.Status {
		case PostStatusPublished:
			if err := post.Publish(); err != nil {
				return nil, err
			}
		case PostStatusDraft:
			post.Unpublish()
		case PostStatusArchived:
			post.Archive()
		}
	}

	if req.FeaturedImage != nil {
		post.FeaturedImage = *req.FeaturedImage
		post.UpdatedAt = time.Now().UTC()
	}

	if req.Featured != nil {
		post.SetFeatured(*req.Featured)
	}

	// Update SEO fields
	if req.MetaTitle != nil || req.MetaDescription != nil || req.MetaKeywords != nil ||
		req.OGTitle != nil || req.OGDescription != nil || req.OGImage != nil {
		metaTitle := ""
		metaDesc := ""
		metaKeywords := ""
		ogTitle := ""
		ogDesc := ""
		ogImage := ""
		if req.MetaTitle != nil {
			metaTitle = *req.MetaTitle
		}
		if req.MetaDescription != nil {
			metaDesc = *req.MetaDescription
		}
		if req.MetaKeywords != nil {
			metaKeywords = *req.MetaKeywords
		}
		if req.OGTitle != nil {
			ogTitle = *req.OGTitle
		}
		if req.OGDescription != nil {
			ogDesc = *req.OGDescription
		}
		if req.OGImage != nil {
			ogImage = *req.OGImage
		}
		post.UpdateSEO(metaTitle, metaDesc, metaKeywords, ogTitle, ogDesc, ogImage)
	}

	// Update relationships if provided
	if req.SkillIDs != nil {
		// Get current skills
		currentSkills, _ := s.repo.GetPostSkills(ctx, postID)
		currentSkillMap := make(map[uuid.UUID]bool)
		for _, skillID := range currentSkills {
			currentSkillMap[skillID] = true
		}

		// Add new skills
		newSkillMap := make(map[uuid.UUID]bool)
		for _, skillID := range req.SkillIDs {
			newSkillMap[skillID] = true
			if !currentSkillMap[skillID] {
				if err := s.repo.AttachSkillToPost(ctx, postID, skillID); err != nil {
					s.logger.Warn("Failed to attach skill to post", "error", err, "postID", postID, "skillID", skillID)
				}
			}
		}

		// Remove skills not in new list
		for skillID := range currentSkillMap {
			if !newSkillMap[skillID] {
				if err := s.repo.DetachSkillFromPost(ctx, postID, skillID); err != nil {
					s.logger.Warn("Failed to detach skill from post", "error", err, "postID", postID, "skillID", skillID)
				}
			}
		}
	}

	if req.CategoryIDs != nil {
		// Similar logic for categories
		currentCategories, _ := s.repo.GetPostCategories(ctx, postID)
		currentCategoryMap := make(map[uuid.UUID]bool)
		for _, cat := range currentCategories {
			currentCategoryMap[cat.ID] = true
		}

		newCategoryMap := make(map[uuid.UUID]bool)
		for _, categoryID := range req.CategoryIDs {
			newCategoryMap[categoryID] = true
			if !currentCategoryMap[categoryID] {
				if err := s.repo.AttachCategoryToPost(ctx, postID, categoryID); err != nil {
					s.logger.Warn("Failed to attach category to post", "error", err, "postID", postID, "categoryID", categoryID)
				}
			}
		}

		for catID := range currentCategoryMap {
			if !newCategoryMap[catID] {
				if err := s.repo.DetachCategoryFromPost(ctx, postID, catID); err != nil {
					s.logger.Warn("Failed to detach category from post", "error", err, "postID", postID, "categoryID", catID)
				}
			}
		}
	}

	if req.TagNames != nil {
		// Similar logic for tags
		currentTags, _ := s.repo.GetPostTags(ctx, postID)
		currentTagMap := make(map[string]uuid.UUID)
		for _, tag := range currentTags {
			currentTagMap[tag.Name] = tag.ID
		}

		newTagMap := make(map[string]bool)
		for _, tagName := range req.TagNames {
			newTagMap[tagName] = true
			if _, exists := currentTagMap[tagName]; !exists {
				tag, err := s.repo.GetOrCreateTag(ctx, tagName)
				if err != nil {
					s.logger.Warn("Failed to get or create tag", "error", err, "tagName", tagName)
					continue
				}
				if err := s.repo.AttachTagToPost(ctx, postID, tag.ID); err != nil {
					s.logger.Warn("Failed to attach tag to post", "error", err, "postID", postID, "tagID", tag.ID)
				}
			}
		}

		for tagName, tagID := range currentTagMap {
			if !newTagMap[tagName] {
				if err := s.repo.DetachTagFromPost(ctx, postID, tagID); err != nil {
					s.logger.Warn("Failed to detach tag from post", "error", err, "postID", postID, "tagID", tagID)
				}
			}
		}
	}

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) GetPost(ctx context.Context, postID uuid.UUID) (*Post, error) {
	return s.repo.GetPost(ctx, postID)
}

func (s *service) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	return s.repo.GetPostBySlug(ctx, slug)
}

func (s *service) DeletePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	return s.repo.DeletePost(ctx, postID, userID)
}

func (s *service) ListPosts(ctx context.Context, filters PostFilters) ([]Post, error) {
	return s.repo.ListPosts(ctx, filters)
}

func (s *service) IncrementPostViews(ctx context.Context, postID uuid.UUID) error {
	return s.repo.IncrementPostViews(ctx, postID)
}

// Category operations

func (s *service) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*Category, error) {
	category, err := NewCategory(req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	// Check if slug is taken
	taken, err := s.repo.IsCategorySlugTaken(ctx, category.Slug, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if taken {
		category.Slug = category.Slug + "-" + category.ID.String()[:8]
	}

	if err := s.repo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *service) UpdateCategory(ctx context.Context, categoryID uuid.UUID, req UpdateCategoryRequest) (*Category, error) {
	category, err := s.repo.GetCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
		category.Slug = generateCategorySlug(req.Name)
		// Check if new slug is taken
		taken, err := s.repo.IsCategorySlugTaken(ctx, category.Slug, categoryID)
		if err != nil {
			return nil, err
		}
		if taken {
			category.Slug = category.Slug + "-" + category.ID.String()[:8]
		}
	}

	if req.Description != "" {
		category.Description = req.Description
	}

	category.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *service) GetCategory(ctx context.Context, categoryID uuid.UUID) (*Category, error) {
	return s.repo.GetCategory(ctx, categoryID)
}

func (s *service) GetCategoryBySlug(ctx context.Context, slug string) (*Category, error) {
	return s.repo.GetCategoryBySlug(ctx, slug)
}

func (s *service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.repo.ListCategories(ctx)
}

// Tag operations

func (s *service) GetOrCreateTag(ctx context.Context, name string) (*Tag, error) {
	return s.repo.GetOrCreateTag(ctx, name)
}

func (s *service) GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error) {
	return s.repo.GetTag(ctx, tagID)
}

func (s *service) GetTagBySlug(ctx context.Context, slug string) (*Tag, error) {
	return s.repo.GetTagBySlug(ctx, slug)
}

func (s *service) ListTags(ctx context.Context) ([]Tag, error) {
	return s.repo.ListTags(ctx)
}

// Post-Skill relationship operations

func (s *service) AttachSkillToPost(ctx context.Context, postID, skillID uuid.UUID) error {
	return s.repo.AttachSkillToPost(ctx, postID, skillID)
}

func (s *service) DetachSkillFromPost(ctx context.Context, postID, skillID uuid.UUID) error {
	return s.repo.DetachSkillFromPost(ctx, postID, skillID)
}

func (s *service) GetPostSkills(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.GetPostSkills(ctx, postID)
}

// Post-Category relationship operations

func (s *service) AttachCategoryToPost(ctx context.Context, postID, categoryID uuid.UUID) error {
	return s.repo.AttachCategoryToPost(ctx, postID, categoryID)
}

func (s *service) DetachCategoryFromPost(ctx context.Context, postID, categoryID uuid.UUID) error {
	return s.repo.DetachCategoryFromPost(ctx, postID, categoryID)
}

func (s *service) GetPostCategories(ctx context.Context, postID uuid.UUID) ([]Category, error) {
	return s.repo.GetPostCategories(ctx, postID)
}

// Post-Tag relationship operations

func (s *service) AttachTagToPost(ctx context.Context, postID, tagID uuid.UUID) error {
	return s.repo.AttachTagToPost(ctx, postID, tagID)
}

func (s *service) DetachTagFromPost(ctx context.Context, postID, tagID uuid.UUID) error {
	return s.repo.DetachTagFromPost(ctx, postID, tagID)
}

func (s *service) GetPostTags(ctx context.Context, postID uuid.UUID) ([]Tag, error) {
	return s.repo.GetPostTags(ctx, postID)
}

