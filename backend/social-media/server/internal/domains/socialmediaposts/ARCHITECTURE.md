# Social Media Posts Domain - Architecture Overview

## 🏗️ Domain Architecture

The `socialmediaposts` domain is built using a **modular subdomain architecture** that separates concerns into distinct, focused modules. This design allows for better maintainability, scalability, and clear separation of responsibilities.

## 📐 Overall Structure

```
socialmediaposts/
├── Core Domain (Main)
│   ├── entity.go          # SocialMediaPost & SocialMediaEntityLink entities
│   ├── repository.go       # Core CRUD operations
│   ├── service.go          # Business logic orchestration
│   ├── handler.go          # HTTP handlers
│   ├── routes.go           # Route mounting & subdomain integration
│   └── errors.go          # Domain-specific errors
│
└── Subdomains (Specialized Modules)
    ├── platforms/         # Platform configuration & management
    ├── content/           # Content repurposing workflow
    ├── scheduling/        # Post scheduling & conflict management
    ├── analytics/         # Metrics tracking & reporting
    └── assets/            # Media asset management
```

## 🔄 How It Works

### 1. Core Domain - The Foundation

The **core domain** manages the fundamental `SocialMediaPost` entity and entity linking functionality.

#### Key Entities:
- **SocialMediaPost**: Represents a single social media post with:
  - Platform (LinkedIn, Twitter, Instagram, etc.)
  - Format (long-form, thread, carousel, article, etc.)
  - Status workflow (draft → ready → scheduled → posted → analyzed → archived)
  - Content (title, content, word count, image count)
  - Scheduling timestamps (scheduled_at, posted_at, analyzed_at)
  - Engagement metrics (likes, shares, comments, views)
  - Optional URL (generated after posting)

- **SocialMediaEntityLink**: Links social posts to other entities (posts, projects, skills, etc.) with relationship types (main_topic, secondary_topic, etc.)

#### Workflow:
```
Create Post → Draft → Ready → Scheduled → Posted → Analyzed → Archived
```

### 2. Platforms Subdomain - Platform Configuration

**Purpose**: Manages platform-specific configurations and optimal posting strategies.

#### Key Features:
- **PlatformConfig** entity stores:
  - Platform name and display name
  - Posting frequency recommendations
  - Best days and times for posting
  - Supported content formats per platform
  - Active/inactive status

#### Use Cases:
- Initialize default platforms (LinkedIn, Twitter, Instagram, Medium, Substack, Valete+, Website)
- Configure optimal posting times per platform
- Query platform capabilities (what formats are supported)
- Get optimal posting times for auto-scheduling

#### Example Flow:
```go
// Get platform config
platform := platformsService.GetConfigByName(ctx, "linkedin")

// Get optimal times
optimalTimes := platformsService.GetOptimalTimes(ctx, "linkedin")
// Returns: best days, best times, posting frequency
```

### 3. Content Subdomain - Content Repurposing

**Purpose**: Manages the workflow of repurposing backend blog posts into social media content.

#### Key Entities:
- **ContentPost**: Links to a backend `Post` entity, tracks:
  - Priority (low, medium, high)
  - Status (pending, in_progress, completed, archived)
  - Project association
  - Content type

- **ContentRepurposing**: Tracks which social posts were created from which content post

#### Workflow:
```
Backend Post → ContentPost (pending) → Repurpose → Multiple SocialMediaPosts → Completed
```

#### Example Flow:
```go
// 1. Create content post from backend post
contentPost := contentService.CreateContentPostFromBackend(ctx, CreateContentPostRequest{
    PostID: backendPostID,
    Priority: "high",
})

// 2. Repurpose to multiple platforms
socialPosts := contentService.RepurposeToPlatforms(ctx, contentPostID, RepurposeRequest{
    Platforms: [
        {Platform: "linkedin", Format: "long-form", Title: "...", Content: "..."},
        {Platform: "twitter", Format: "thread", Title: "...", Content: "..."},
    ],
})
// Creates multiple SocialMediaPost entities and links them via ContentRepurposing

// 3. Get repurposing history
history := contentService.GetRepurposingHistory(ctx, contentPostID)
// Returns all social posts created from this content post
```

### 4. Scheduling Subdomain - Post Scheduling

**Purpose**: Manages when posts should be published, with conflict detection and auto-scheduling.

#### Key Entity:
- **ScheduledPost**: Links a social post to a specific date/time
  - Social post ID
  - Scheduled date/time
  - Platform config ID (optional)
  - Status (pending, scheduled, posted, cancelled)

#### Features:
- **Conflict Detection**: Prevents scheduling multiple posts within 15 minutes
- **Auto-Scheduling**: Suggests optimal times based on platform config
- **Date Range Queries**: Get schedule for specific periods
- **Upcoming Posts**: Get posts scheduled for today/this week

#### Example Flow:
```go
// 1. Manual scheduling with conflict check
schedule := schedulingService.SchedulePost(ctx, SchedulePostRequest{
    SocialPostID: postID,
    ScheduledAt: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
})
// Automatically checks for conflicts

// 2. Auto-schedule based on platform
schedule := schedulingService.AutoSchedule(ctx, postID, "linkedin")
// Finds next optimal time slot based on platform config

// 3. Get upcoming posts
upcoming := schedulingService.GetUpcomingPosts(ctx, 10)
// Returns next 10 scheduled posts
```

### 5. Analytics Subdomain - Metrics Tracking

**Purpose**: Tracks and analyzes post performance over time.

#### Key Entity:
- **PostAnalytics**: Time-series metrics for each post
  - Metric date (allows tracking over time)
  - Engagement metrics (likes, comments, shares, views, clicks)
  - Reach and impressions
  - Calculated engagement rate

#### Features:
- **Time-Series Tracking**: Record metrics for different dates
- **Engagement Rate Calculation**: Automatic calculation based on impressions
- **Analytics Summaries**: Aggregate metrics across posts
- **Top Posts**: Find best performing posts by metric

#### Example Flow:
```go
// 1. Record analytics
analytics := analyticsService.RecordAnalytics(ctx, RecordAnalyticsRequest{
    SocialPostID: postID,
    MetricDate: time.Now(),
    Likes: 150,
    Comments: 20,
    Shares: 10,
    Views: 1000,
    Impressions: 1200,
})
// Automatically calculates engagement rate

// 2. Get analytics for a post
postAnalytics := analyticsService.GetPostAnalytics(ctx, postID, startDate, endDate)
// Returns time-series data

// 3. Get summary
summary := analyticsService.GetAnalyticsSummary(ctx, AnalyticsFilters{
    StartDate: startDate,
    EndDate: endDate,
})
// Returns aggregated totals and averages
```

### 6. Assets Subdomain - Media Management

**Purpose**: Manages media assets (images, videos, documents) associated with content or social posts.

#### Key Entity:
- **ContentAsset**: Represents a media file
  - Asset type (image, video, document, other)
  - File path and URL
  - Alt text for accessibility
  - Links to either content post or social post (or both)

#### Use Cases:
- Attach images to content posts before repurposing
- Link assets to specific social media posts
- Manage asset metadata (alt text, URLs)

#### Example Flow:
```go
// 1. Create asset for content post
asset := assetsService.CreateAsset(ctx, CreateAssetRequest{
    ContentPostID: contentPostID,
    AssetType: "image",
    FilePath: "/uploads/image.jpg",
    FileURL: "https://cdn.example.com/image.jpg",
    AltText: "Screenshot of application",
})

// 2. Get assets for a social post
assets := assetsService.GetAssetsBySocialPost(ctx, socialPostID)
```

## 🔗 How Subdomains Interact

### Typical Workflow Example:

```
1. Content Repurposing Flow:
   Backend Post
   ↓
   ContentPost (content subdomain)
   ↓
   Repurpose to Platforms (content subdomain)
   ↓
   Multiple SocialMediaPosts (core domain)
   ↓
   Schedule Posts (scheduling subdomain)
   ↓
   Post Published
   ↓
   Track Analytics (analytics subdomain)

2. Auto-Scheduling Flow:
   SocialMediaPost (core domain)
   ↓
   Get Platform Config (platforms subdomain)
   ↓
   Find Optimal Time (scheduling subdomain)
   ↓
   Schedule Post (scheduling subdomain)
```

### Integration Points:

1. **Content → Core**: Creates `SocialMediaPost` entities via `socialPostsService.CreatePost()`
2. **Scheduling → Core**: Updates `SocialMediaPost.status` to "scheduled"
3. **Scheduling → Platforms**: Queries platform config for optimal times
4. **Analytics → Core**: Links analytics to `SocialMediaPost` entities
5. **Assets → Content/Social**: Links assets to either content posts or social posts

## 📊 Data Flow Examples

### Example 1: Complete Post Lifecycle

```
1. User creates a blog post in backend
   → Post entity created

2. User marks post for repurposing
   → ContentPost created (content subdomain)
   → Status: pending

3. User repurposes to LinkedIn and Twitter
   → Two SocialMediaPost entities created (core domain)
   → ContentRepurposing records created (content subdomain)
   → ContentPost status: in_progress

4. System auto-schedules LinkedIn post
   → Queries platform config (platforms subdomain)
   → Finds optimal time (scheduling subdomain)
   → ScheduledPost created
   → SocialMediaPost status: scheduled

5. Post is published
   → SocialMediaPost status: posted
   → ScheduledPost status: posted
   → URL and platform post ID stored

6. Analytics tracked daily
   → PostAnalytics records created (analytics subdomain)
   → Engagement metrics updated
   → Engagement rate calculated

7. After analysis period
   → SocialMediaPost status: analyzed
   → ContentPost status: completed
```

### Example 2: Querying and Reporting

```
1. Get content backlog
   → contentService.GetContentBacklog()
   → Returns ContentPosts with status: pending

2. Get upcoming scheduled posts
   → schedulingService.GetUpcomingPosts(limit: 10)
   → Returns ScheduledPosts with status: scheduled

3. Get analytics summary
   → analyticsService.GetAnalyticsSummary(filters)
   → Returns aggregated metrics across all posts

4. Get repurposing history
   → contentService.GetRepurposingHistory(contentPostID)
   → Returns all SocialMediaPosts created from content post
```

## 🎯 Key Design Principles

1. **Separation of Concerns**: Each subdomain handles one specific aspect
2. **Single Responsibility**: Each service has a clear, focused purpose
3. **Domain-Driven Design**: Entities and services reflect business concepts
4. **Status Workflows**: Clear state transitions with validation
5. **Conflict Prevention**: Scheduling prevents overlapping posts
6. **Time-Series Analytics**: Track metrics over time, not just current state
7. **Flexible Linking**: Assets and links can connect to multiple entity types

## 🔐 Error Handling

Each subdomain has its own error codes:
- Core: 8000-8011
- Platforms: 8100-8104
- Content: 8200-8205
- Scheduling: 8300-8305
- Analytics: 8400-8403
- Assets: 8500-8504

Errors are domain-specific and provide clear context about what went wrong.

## 🚀 API Structure

All routes are mounted under `/api/social-media-posts` (or similar base path):

```
POST   /                           # Create social post
GET    /                           # List posts
GET    /:id                        # Get post
PATCH  /:id                        # Update post
PATCH  /:id/status                 # Update status
PATCH  /:id/engagement             # Update engagement
DELETE /:id                        # Delete post

# Subdomain routes
GET    /platforms                  # List platforms
GET    /platforms/:id              # Get platform config
PATCH  /platforms/:id              # Update platform config

POST   /content/posts              # Create content post
POST   /content/posts/:id/repurpose # Repurpose to platforms
GET    /content/backlog            # Get backlog

POST   /scheduling                 # Schedule post
GET    /scheduling/upcoming        # Get upcoming
POST   /scheduling/:id/auto        # Auto-schedule

POST   /analytics                  # Record analytics
GET    /analytics/posts/:id       # Get post analytics
GET    /analytics/summary          # Get summary

POST   /assets                     # Create asset
GET    /assets/content-posts/:id  # Get assets by content post
```

This architecture provides a clean, maintainable, and scalable foundation for managing social media content across multiple platforms.
