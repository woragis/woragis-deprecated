# Creative Service Integration - How It Works

## Currently Integrated Domains

### ✅ Posts Domain (Fully Integrated)
**Location**: `app/internal/domains/posts/`

**Integration Points**:
1. **Handler Dependency Injection**: The posts handler receives `creativeAssetsService` as a dependency
2. **New Handler Methods**:
   - `GeneratePostThumbnail()` - Generates thumbnail for a post
   - `GeneratePostFeaturedImage()` - Generates featured image
   - `GeneratePostOGImage()` - Generates Open Graph image
   - `GetPostAssets()` - Retrieves all assets for a post

**Endpoints Added**:
- `POST /api/v1/posts/:id/assets/generate/thumbnail`
- `POST /api/v1/posts/:id/assets/generate/featured-image`
- `POST /api/v1/posts/:id/assets/generate/og-image`
- `GET /api/v1/posts/:id/assets`

**How It Works**:
1. User makes request to generate asset (e.g., thumbnail)
2. Handler verifies ownership of the post
3. Handler calls `creativeAssetsService.GenerateAndStoreThumbnail()`
4. Service calls creative-service API to generate image
5. Service stores base64 data in `creative_assets` table
6. Service updates post entity with asset URL (e.g., `FeaturedImage` field)
7. Returns the created asset

## Integration Pattern

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Domain Handler                           │
│  (e.g., posts/handler.go)                                   │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  creativeAssetsService (injected dependency)        │  │
│  │  - GenerateAndStoreThumbnail()                      │  │
│  │  - GenerateAndStoreImage()                          │  │
│  │  - GenerateAndStoreDiagram()                        │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              Creative Assets Service                        │
│  (creativeassets/service.go)                                │
│                                                             │
│  1. Calls creative-service API                             │
│  2. Receives base64 image data                             │
│  3. Stores in creative_assets table                        │
│  4. Returns CreativeAsset entity                           │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│            Creative Service Client                          │
│  (services/creative/client.go)                              │
│                                                             │
│  HTTP calls to: http://creative-service:8000               │
│  - /v1/images/generate                                     │
│  - /v1/images/generate/thumbnail                           │
│  - /v1/diagrams/generate                                   │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│           Creative Service (Python/FastAPI)                 │
│  (backend/creative-service/)                                │
│                                                             │
│  - OpenAI DALL-E                                            │
│  - Stable Diffusion                                         │
│  - Mermaid/Graphviz diagram generation                      │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow

1. **Request Flow**:
   ```
   Client → Domain Handler → Creative Assets Service → Creative Service Client → Creative Service API
   ```

2. **Response Flow**:
   ```
   Creative Service API → Creative Service Client → Creative Assets Service → Database → Domain Handler → Client
   ```

3. **Storage**:
   - Base64 data stored in `creative_assets.b64_data` (TEXT column)
   - Asset metadata stored (provider, format, dimensions, etc.)
   - Link to entity via `entity_type` + `entity_id`
   - Purpose stored (thumbnail, featured_image, og_image, etc.)

### Integration Steps (for new domains)

To integrate creative assets with another domain (e.g., Case Studies):

1. **Update Handler Constructor**:
   ```go
   // In casestudies/handler.go
   type handler struct {
       service Service
       creativeAssetsService creativeassets.Service  // Add this
       logger *slog.Logger
   }
   
   func NewHandler(service Service, creativeAssetsService creativeassets.Service, logger *slog.Logger) Handler {
       return &handler{
           service: service,
           creativeAssetsService: creativeAssetsService,  // Add this
           logger: logger,
       }
   }
   ```

2. **Add Handler Methods**:
   ```go
   func (h *handler) GenerateCaseStudyDiagram(c *fiber.Ctx) error {
       userID, err := authdomain.UserIDFromContext(c)
       // ... verify ownership ...
       
       var payload struct {
           Description string `json:"description"`
           DiagramKind string `json:"diagramKind,omitempty"`
       }
       c.BodyParser(&payload)
       
       asset, err := h.creativeAssetsService.GenerateAndStoreDiagram(
           c.Context(),
           userID,
           creativeassets.EntityTypeCaseStudy,
           caseStudyID,
           payload.Description,
           payload.DiagramKind,
       )
       // ... handle response ...
   }
   ```

3. **Add Routes**:
   ```go
   // In casestudies/routes.go
   api.Post("/:id/assets/generate/diagram", handler.GenerateCaseStudyDiagram)
   ```

4. **Update main.go**:
   ```go
   // Pass creativeAssetsService to handler constructor
   caseStudyHandler := casestudiesdomain.NewHandler(
       caseStudyService,
       creativeAssetsService,  // Add this
       translationEnricher,
       translationService,
       slogLogger,
   )
   ```

## Entity Linking

Creative assets are linked to entities using:
- **entity_type**: Enum value (e.g., `EntityTypePost`, `EntityTypeCaseStudy`)
- **entity_id**: UUID of the entity
- **purpose**: What the asset is used for (thumbnail, featured_image, diagram, etc.)

This allows:
- Multiple assets per entity (e.g., thumbnail + featured_image)
- Querying all assets for an entity
- Querying specific asset by purpose

## Base64 Storage

Base64 data is stored in the database `b64_data` TEXT column:
- **Pros**: Simple, no external storage needed, atomic with entity
- **Cons**: Database size, not optimized for large files
- **Retrieval**: Via `/api/v1/creative-assets/:id/data` endpoint

The endpoint serves the base64 as an image with proper content-type headers, or redirects to URL if available.

## Current Status

- ✅ Creative Assets Domain created
- ✅ Creative Service Client created
- ✅ Posts domain integrated
- ⏳ Case Studies (not yet integrated)
- ⏳ Technical Writings (not yet integrated)
- ⏳ Problem Solutions (not yet integrated)
- ⏳ System Designs (not yet integrated)

