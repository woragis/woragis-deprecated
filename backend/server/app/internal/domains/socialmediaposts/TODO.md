# Social Media Posts Domain - Complete Rebuild TODO

## Overview
Complete rebuild of the `socialmediaposts` domain with modular subdomains. Since no data exists, we can redesign from scratch with a clean architecture.

## Architecture Decision

### Domain Structure
```
socialmediaposts/
├── entity.go              # Core SocialMediaPost entity
├── errors.go              # Domain errors
├── routes.go              # Main routes (includes subdomain routes)
├── handler.go             # Main handler (minimal, delegates to subdomains)
├── service.go             # Main service (orchestrator)
├── repository.go          # Core repository
│
├── platforms/             # Platform management subdomain
│   ├── entity.go
│   ├── service.go
│   ├── repository.go
│   ├── handler.go
│   ├── routes.go
│   └── errors.go
│
├── content/               # Content management & repurposing subdomain
│   ├── entity.go
│   ├── service.go
│   ├── repository.go
│   ├── handler.go
│   ├── routes.go
│   └── errors.go
│
├── scheduling/            # Post scheduling subdomain
│   ├── entity.go
│   ├── service.go
│   ├── repository.go
│   ├── handler.go
│   ├── routes.go
│   └── errors.go
│
├── analytics/             # Analytics & metrics subdomain
│   ├── entity.go
│   ├── service.go
│   ├── repository.go
│   ├── handler.go
│   ├── routes.go
│   └── errors.go
│
├── assets/                # Content assets subdomain
│   ├── entity.go
│   ├── service.go
│   ├── repository.go
│   ├── handler.go
│   ├── routes.go
│   └── errors.go
│
└── links/                 # Entity linking subdomain (keep existing functionality)
    ├── entity.go
    ├── service.go
    ├── repository.go
    ├── handler.go
    ├── routes.go
    └── errors.go
```

## Implementation Phases

### Phase 1: Core Domain Foundation ✅
- [x] **1.1** Design core SocialMediaPost entity
  - [x] Add platform field (enum: linkedin, twitter, instagram, medium, substack, valete, website)
  - [x] Add content format field (enum: long-form, thread, carousel, article, newsletter, post)
  - [x] Add status workflow (draft, ready, scheduled, posted, analyzed, archived)
  - [x] Add scheduling fields (scheduled_at, posted_at, analyzed_at)
  - [x] Add content fields (title, content, word_count, image_count)
  - [x] Add metadata JSONB field
  - [x] Keep URL field (but make it optional, generated after posting)
  - [x] Keep engagement metrics (likes, shares, comments, views)
  
- [x] **1.2** Create core repository interface and implementation
  - [x] CRUD operations
  - [x] Filtering by platform, status, format
  - [x] Query by content_post_id (for repurposing tracking)
  
- [x] **1.3** Create core service interface
  - [x] Basic CRUD operations
  - [x] Status transitions with validation
  
- [x] **1.4** Create core handler
  - [x] POST / - Create social post
  - [x] GET / - List posts with filters
  - [x] GET /:id - Get post by ID
  - [x] PATCH /:id - Update post
  - [x] DELETE /:id - Delete post
  - [x] PATCH /:id/status - Update status
  - [x] PATCH /:id/engagement - Update engagement metrics

- [x] **1.5** Update routes.go to include subdomain routes
  - [x] Mount platforms routes at /platforms
  - [x] Mount content routes at /content
  - [x] Mount scheduling routes at /scheduling
  - [x] Mount analytics routes at /analytics
  - [x] Mount assets routes at /assets
  - [x] Note: Links functionality remains in main domain (not extracted to separate subdomain)

### Phase 2: Platforms Subdomain ✅
- [x] **2.1** Create platforms entity
  - [x] PlatformConfig entity
  - [x] Fields: name, display_name, posting_frequency, best_days, best_times, supported_formats, metadata
  
- [x] **2.2** Create platforms repository
  - [x] CRUD operations
  - [x] Get platform by name
  
- [x] **2.3** Create platforms service
  - [x] Initialize default platforms (LinkedIn, Twitter, Instagram, Medium, Substack, Valete+, Website)
  - [x] Update platform config
  - [x] Get optimal posting times for platform
  
- [x] **2.4** Create platforms handler
  - [x] GET /platforms - List all platforms
  - [x] GET /platforms/:id - Get platform config
  - [x] GET /platforms/by-name/:name - Get platform by name
  - [x] PATCH /platforms/:id - Update platform config
  - [x] GET /platforms/:name/optimal-times - Get optimal posting times
  
- [x] **2.5** Create platforms routes

### Phase 3: Content Subdomain ✅
- [x] **3.1** Create content entities
  - [x] ContentPost entity (links to backend Post)
  - [x] Fields: post_id, content_type, project, priority, status, metadata
  - [x] ContentRepurposing entity (tracks repurposing relationships)
  
- [x] **3.2** Create content repository
  - [x] Create content post from backend post
  - [x] Track repurposing relationships
  - [x] Get repurposing history
  
- [x] **3.3** Create content service
  - [x] CreateContentPostFromBackend - Create from existing Post
  - [x] RepurposeToPlatforms - Create social posts for multiple platforms
  - [x] AdaptContentForPlatform - Transform content for platform/format
  - [x] GetRepurposingHistory - Get all social posts from a content post
  - [x] GetContentBacklog - Get posts ready for repurposing
  
- [x] **3.4** Create content handler
  - [x] POST /content/posts - Create content post from backend post
  - [x] GET /content/posts - List content posts
  - [x] GET /content/posts/:id - Get content post with social posts
  - [x] POST /content/posts/:id/repurpose - Repurpose to platforms
  - [x] GET /content/backlog - Get content backlog
  - [x] PATCH /content/posts/:id/priority - Update priority
  
- [x] **3.5** Create content routes

### Phase 4: Scheduling Subdomain ✅
- [x] **4.1** Create scheduling entity
  - [x] ScheduledPost entity
  - [x] Fields: social_post_id, platform_id, scheduled_date, scheduled_time, status, metadata
  
- [x] **4.2** Create scheduling repository
  - [x] Schedule post
  - [x] Get schedule for date range
  - [x] Get upcoming posts
  - [x] Check conflicts
  - [x] Update/cancel schedule
  
- [x] **4.3** Create scheduling service
  - [x] SchedulePost - Schedule a post with conflict checking
  - [x] GetSchedule - Get schedule for date range
  - [x] GetUpcomingPosts - Get posts scheduled for today/this week
  - [x] CancelSchedule - Cancel scheduled post
  - [x] AutoSchedule - Suggest optimal times based on platform config
  - [x] CheckConflicts - Check if time slot is available
  
- [x] **4.4** Create scheduling handler
  - [x] POST /scheduling - Schedule a post
  - [x] GET /scheduling - Get schedule (with date range filters)
  - [x] GET /scheduling/upcoming - Get upcoming scheduled posts
  - [x] GET /scheduling/:id - Get scheduled post
  - [x] PATCH /scheduling/:id - Update schedule
  - [x] DELETE /scheduling/:id - Cancel schedule
  - [x] POST /scheduling/:id/auto - Auto-schedule with optimal time
  
- [x] **4.5** Create scheduling routes

### Phase 5: Analytics Subdomain ✅
- [x] **5.1** Create analytics entity
  - [x] PostAnalytics entity
  - [x] Fields: social_post_id, metric_date, likes, comments, shares, views, clicks, engagement_rate, reach, impressions, metadata
  
- [x] **5.2** Create analytics repository
  - [x] Record analytics
  - [x] Get analytics for post (date range)
  - [x] Get analytics summary (aggregations)
  - [x] Get platform analytics (placeholder)
  
- [x] **5.3** Create analytics service
  - [x] RecordAnalytics - Record metrics for a post
  - [x] GetPostAnalytics - Get analytics for specific post
  - [x] GetAnalyticsSummary - Get aggregated analytics
  - [x] GetTopPosts - Get top performing posts
  - [x] CalculateEngagementRate - Calculate engagement rate
  
- [x] **5.4** Create analytics handler
  - [x] POST /analytics - Record analytics
  - [x] GET /analytics/posts/:id - Get post analytics
  - [x] GET /analytics/summary - Get analytics summary (with filters)
  - [x] GET /analytics/top-posts - Get top performing posts
  
- [x] **5.5** Create analytics routes

### Phase 6: Assets Subdomain ✅
- [x] **6.1** Create assets entity
  - [x] ContentAsset entity
  - [x] Fields: content_post_id, social_post_id (optional), asset_type, file_path, file_url, alt_text, metadata
  
- [x] **6.2** Create assets repository
  - [x] CRUD operations
  - [x] Get assets by content post
  - [x] Get assets by social post
  
- [x] **6.3** Create assets service
  - [x] Create/Update/Delete assets
  - [x] Link assets to content or social posts
  
- [x] **6.4** Create assets handler
  - [x] POST /assets - Create asset
  - [x] GET /assets/:id - Get asset
  - [x] GET /assets/content-posts/:contentPostId - Get assets by content post
  - [x] GET /assets/social-posts/:socialPostId - Get assets by social post
  - [x] PATCH /assets/:id - Update asset
  - [x] DELETE /assets/:id - Delete asset
  
- [x] **6.5** Create assets routes

## Implementation Status: ✅ COMPLETE

All phases have been successfully implemented. The domain is ready for integration and testing.

---

## Integration Phases

### Phase 7: Creative Service Integration 🎨

**Purpose**: Integrate with creative-service for AI-powered media generation, uploads, and asset management.

- [ ] **7.1** Media Generation Integration
  - [ ] Generate images for posts (DALL-E, Midjourney, Stable Diffusion)
  - [ ] Generate GIFs for posts (animated content)
  - [ ] Generate videos for posts (short-form content)
  - [ ] Style consistency across platforms (brand guidelines)
  - [ ] Batch media generation for content repurposing
  - [ ] Media generation queue integration
  
- [ ] **7.2** Asset Upload & Management
  - [ ] Upload generated media to CDN/storage
  - [ ] Automatic asset linking to ContentAsset entities
  - [ ] Media optimization per platform (sizes, formats)
  - [ ] Thumbnail generation for videos
  - [ ] Image compression and format conversion
  - [ ] Asset versioning and rollback
  
- [ ] **7.3** Creative Service Client
  - [ ] Create creative-service client interface
  - [ ] Implement media generation requests
  - [ ] Handle async generation callbacks
  - [ ] Error handling and retry logic
  - [ ] Rate limiting and quota management
  
- [ ] **7.4** Content Subdomain Integration
  - [ ] Auto-generate media when repurposing content
  - [ ] Platform-specific media generation (LinkedIn carousel, Instagram stories, etc.)
  - [ ] Media suggestions based on content analysis
  - [ ] Generate alt text using AI vision models
  
- [ ] **7.5** Assets Subdomain Enhancement
  - [ ] Link generated assets to ContentAsset entities
  - [ ] Track generation metadata (model used, prompt, cost)
  - [ ] Store generation history and versions
  - [ ] Regenerate assets on demand

### Phase 8: AI Service Integration 🤖

**Purpose**: Integrate with ai-service for intelligent content generation, planning, and optimization.

- [ ] **8.1** Content Generation & Repurposing
  - [ ] AI-powered content repurposing (adapt content for each platform)
  - [ ] Generate post captions and titles
  - [ ] Generate hashtags per platform
  - [ ] Content tone/style adaptation (professional for LinkedIn, casual for Twitter)
  - [ ] Content length optimization per platform
  - [ ] Thread generation for Twitter/X
  - [ ] Article expansion for Medium/Substack
  
- [ ] **8.2** Content Planning & Strategy
  - [ ] Generate weekly content plans
  - [ ] Generate monthly content calendars
  - [ ] Platform-specific content plans (LinkedIn plan, Twitter plan, etc.)
  - [ ] Content theme suggestions
  - [ ] Content mix recommendations (educational, promotional, personal)
  - [ ] Optimal posting frequency suggestions
  - [ ] Content gap analysis
  
- [ ] **8.3** AI-Powered Scheduling
  - [ ] AI-suggested optimal posting times (beyond platform config)
  - [ ] Engagement prediction for different time slots
  - [ ] Content-performance correlation analysis
  - [ ] Auto-reschedule based on predicted performance
  
- [ ] **8.4** Analytics & Insights
  - [ ] AI-powered analytics insights
  - [ ] Content performance predictions
  - [ ] Audience engagement analysis
  - [ ] Content recommendations based on performance
  - [ ] Trend analysis and topic suggestions
  - [ ] Competitor content analysis (if data available)
  
- [ ] **8.5** Content Optimization
  - [ ] A/B testing suggestions for content
  - [ ] Content improvement recommendations
  - [ ] SEO optimization for posts
  - [ ] Engagement optimization suggestions
  - [ ] Content refresh recommendations (when to repost/update)
  
- [ ] **8.6** AI Service Client
  - [ ] Create ai-service client interface
  - [ ] Implement chat completion requests
  - [ ] Handle streaming responses for long content
  - [ ] Context management for multi-turn conversations
  - [ ] Token usage tracking and cost management
  - [ ] Model selection (GPT-4, Claude, etc.)

### Phase 9: Worker Integrations ⚙️

**Purpose**: Integrate background workers for async processing, notifications, and automation.

- [ ] **9.1** Email Worker Integration
  - [ ] Send notifications for scheduled posts (reminders)
  - [ ] Weekly analytics reports via email
  - [ ] Monthly content performance summaries
  - [ ] Content suggestions and recommendations
  - [ ] Post publishing confirmations
  - [ ] Engagement milestone notifications
  - [ ] Content calendar reminders
  - [ ] Backlog notifications (pending repurposing)
  
- [ ] **9.2** Translation Worker Integration
  - [ ] Auto-translate posts for multi-language audiences
  - [ ] Platform-specific translation (LinkedIn in English, Twitter in Spanish)
  - [ ] Translation quality checks
  - [ ] Cultural adaptation (not just literal translation)
  - [ ] Multi-language content repurposing
  
- [ ] **9.3** Media Processing Worker
  - [ ] Background image optimization
  - [ ] Video transcoding and compression
  - [ ] Format conversion (PNG to JPG, MP4 to GIF)
  - [ ] Platform-specific media processing
  - [ ] Batch processing for bulk uploads
  
- [ ] **9.4** Analytics Collection Worker
  - [ ] Periodic analytics collection from platforms (daily/hourly)
  - [ ] Automatic engagement rate calculation
  - [ ] Trend detection and alerting
  - [ ] Performance anomaly detection
  - [ ] Automated reporting generation
  
- [ ] **9.5** Publishing Worker
  - [ ] Background post publishing to platforms
  - [ ] Retry logic for failed publishes
  - [ ] Platform API rate limit handling
  - [ ] Post status updates after publishing
  - [ ] URL and platform post ID capture
  
- [ ] **9.6** Content Generation Worker
  - [ ] Background content generation (async)
  - [ ] Batch content repurposing
  - [ ] Scheduled content plan generation
  - [ ] AI-powered content suggestions queue
  - [ ] Media generation queue processing

### Phase 10: Advanced Features & Integrations 🚀

**Purpose**: Advanced features leveraging all integrated services.

- [ ] **10.1** Smart Content Repurposing
  - [ ] AI analyzes backend post and suggests platforms
  - [ ] Auto-generates platform-specific content variations
  - [ ] Auto-generates platform-specific media
  - [ ] Suggests optimal posting times per platform
  - [ ] One-click repurpose to all platforms
  
- [ ] **10.2** Content Calendar Intelligence
  - [ ] AI-generated content calendar with themes
  - [ ] Automatic content suggestions based on calendar
  - [ ] Content mix balancing (educational, promotional, etc.)
  - [ ] Holiday and event-aware content suggestions
  - [ ] Content series planning (multi-part posts)
  
- [ ] **10.3** Performance-Based Optimization
  - [ ] Auto-adjust posting times based on performance
  - [ ] Content style recommendations based on analytics
  - [ ] Platform performance comparison
  - [ ] Content refresh suggestions for underperforming posts
  - [ ] Best-performing content replication
  
- [ ] **10.4** Multi-Channel Campaigns
  - [ ] Create campaigns spanning multiple platforms
  - [ ] Coordinated posting across platforms
  - [ ] Campaign performance tracking
  - [ ] Cross-platform engagement analysis
  
- [ ] **10.5** Content Templates & Automation
  - [ ] Save successful content as templates
  - [ ] Template-based content generation
  - [ ] Recurring content automation (weekly tips, monthly updates)
  - [ ] Content series automation
  
- [ ] **10.6** Collaboration Features
  - [ ] Content approval workflows
  - [ ] Team member assignments
  - [ ] Content review and feedback
  - [ ] Version control for content drafts

## Integration Architecture

### Service Communication Pattern

```
socialmediaposts domain
    ↓
    ├──→ creative-service (HTTP/gRPC)
    │    ├── Media generation requests
    │    ├── Asset uploads
    │    └── Media optimization
    │
    ├──→ ai-service (HTTP/gRPC)
    │    ├── Content generation
    │    ├── Content planning
    │    ├── Analytics insights
    │    └── Optimization suggestions
    │
    └──→ Workers (Message Queue - RabbitMQ/Kafka)
         ├── email-worker (notifications, reports)
         ├── translation-worker (multi-language)
         ├── media-worker (processing)
         ├── analytics-worker (collection)
         └── publishing-worker (platform APIs)
```

### Integration Points

1. **Content Subdomain → AI Service**
   - Content repurposing with AI adaptation
   - Platform-specific content generation
   - Content planning and calendar generation

2. **Assets Subdomain → Creative Service**
   - Media generation requests
   - Asset upload and management
   - Media optimization

3. **Scheduling Subdomain → AI Service**
   - AI-powered optimal time suggestions
   - Engagement prediction

4. **Analytics Subdomain → AI Service**
   - Performance insights
   - Content recommendations
   - Trend analysis

5. **All Subdomains → Email Worker**
   - Notifications and reports
   - Content suggestions

6. **Core Domain → Publishing Worker**
   - Background post publishing
   - Status updates

### Queue Integration Points

- **Content Generation Queue**: AI service requests for content generation
- **Media Generation Queue**: Creative service requests for media
- **Publishing Queue**: Posts ready to be published
- **Analytics Queue**: Analytics collection tasks
- **Notification Queue**: Email notifications
- **Translation Queue**: Content translation tasks
