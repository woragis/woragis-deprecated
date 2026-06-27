# Content Management System - README

## Overview
Backend system for managing 400-500 social media posts across multiple platforms, tracking content creation, scheduling, posting, and analytics.

## System Requirements

### Core Features
1. **Content Management**: Store and organize 400-500 posts
2. **Platform Management**: Track posts across 7 platforms
3. **Scheduling**: Schedule posts with platform-specific timing
4. **Analytics**: Track performance metrics per platform
5. **Repurposing**: Track content repurposing across platforms
6. **Backlog Management**: Prioritize and queue content
7. **Status Tracking**: Track content lifecycle (draft → posted → analyzed)

## Database Schema

### Tables

#### `content_posts`
Stores all content posts from backend and other projects.

```sql
CREATE TABLE content_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    content_type VARCHAR(50) NOT NULL, -- 'architecture', 'implementation', 'lessons', 'decision', 'advanced'
    project VARCHAR(100) NOT NULL, -- 'backend', 'project2', 'project3', etc.
    category VARCHAR(100), -- 'architecture', 'technical-decisions', 'implementation', etc.
    priority INTEGER DEFAULT 2, -- 0 (P0) to 3 (P3)
    status VARCHAR(50) DEFAULT 'draft', -- 'draft', 'ready', 'scheduled', 'posted', 'archived'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB -- Additional metadata (tags, keywords, etc.)
);

CREATE INDEX idx_content_posts_type ON content_posts(content_type);
CREATE INDEX idx_content_posts_project ON content_posts(project);
CREATE INDEX idx_content_posts_status ON content_posts(status);
CREATE INDEX idx_content_posts_priority ON content_posts(priority);
```

#### `platforms`
Stores platform information.

```sql
CREATE TABLE platforms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL, -- 'linkedin', 'twitter', 'instagram', etc.
    display_name VARCHAR(100) NOT NULL,
    api_endpoint VARCHAR(255),
    posting_frequency INTEGER, -- Posts per week
    best_days TEXT[], -- ['tuesday', 'wednesday', 'thursday']
    best_times TEXT[], -- ['8-10 AM', '12-1 PM']
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### `social_posts`
Stores platform-specific post instances.

```sql
CREATE TABLE social_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_post_id UUID REFERENCES content_posts(id) ON DELETE CASCADE,
    platform_id UUID REFERENCES platforms(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    format VARCHAR(50), -- 'long-form', 'thread', 'carousel', 'article', 'newsletter'
    word_count INTEGER,
    image_count INTEGER DEFAULT 0,
    scheduled_at TIMESTAMP,
    posted_at TIMESTAMP,
    status VARCHAR(50) DEFAULT 'draft', -- 'draft', 'scheduled', 'posted', 'archived'
    platform_post_id VARCHAR(255), -- External platform post ID
    url VARCHAR(500), -- Link to posted content
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB -- Platform-specific metadata
);

CREATE INDEX idx_social_posts_content ON social_posts(content_post_id);
CREATE INDEX idx_social_posts_platform ON social_posts(platform_id);
CREATE INDEX idx_social_posts_status ON social_posts(status);
CREATE INDEX idx_social_posts_scheduled ON social_posts(scheduled_at);
```

#### `content_assets`
Stores visual assets (diagrams, code snippets, images).

```sql
CREATE TABLE content_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_post_id UUID REFERENCES content_posts(id) ON DELETE CASCADE,
    social_post_id UUID REFERENCES social_posts(id) ON DELETE SET NULL,
    asset_type VARCHAR(50) NOT NULL, -- 'diagram', 'code-snippet', 'image', 'carousel'
    file_path VARCHAR(500) NOT NULL,
    file_url VARCHAR(500),
    alt_text TEXT,
    metadata JSONB, -- Dimensions, format, etc.
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_content_assets_post ON content_assets(content_post_id);
CREATE INDEX idx_content_assets_social ON content_assets(social_post_id);
```

#### `post_analytics`
Stores analytics data for posted content.

```sql
CREATE TABLE post_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    social_post_id UUID REFERENCES social_posts(id) ON DELETE CASCADE,
    platform_id UUID REFERENCES platforms(id) ON DELETE CASCADE,
    metric_date DATE NOT NULL,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0,
    views INTEGER DEFAULT 0,
    clicks INTEGER DEFAULT 0,
    engagement_rate DECIMAL(5,2),
    reach INTEGER DEFAULT 0,
    impressions INTEGER DEFAULT 0,
    metadata JSONB, -- Platform-specific metrics
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(social_post_id, platform_id, metric_date)
);

CREATE INDEX idx_post_analytics_post ON post_analytics(social_post_id);
CREATE INDEX idx_post_analytics_platform ON post_analytics(platform_id);
CREATE INDEX idx_post_analytics_date ON post_analytics(metric_date);
```

#### `content_repurposing`
Tracks content repurposing across platforms.

```sql
CREATE TABLE content_repurposing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_content_post_id UUID REFERENCES content_posts(id) ON DELETE CASCADE,
    source_social_post_id UUID REFERENCES social_posts(id) ON DELETE SET NULL,
    target_social_post_id UUID REFERENCES social_posts(id) ON DELETE CASCADE,
    repurposing_type VARCHAR(50), -- 'summary', 'expansion', 'translation', 'visual'
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_repurposing_source ON content_repurposing(source_content_post_id);
CREATE INDEX idx_repurposing_target ON content_repurposing(target_social_post_id);
```

#### `content_schedule`
Manages posting schedule.

```sql
CREATE TABLE content_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_post_id UUID REFERENCES content_posts(id) ON DELETE CASCADE,
    platform_id UUID REFERENCES platforms(id) ON DELETE CASCADE,
    scheduled_date DATE NOT NULL,
    scheduled_time TIME NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled', -- 'scheduled', 'posted', 'cancelled'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_schedule_date ON content_schedule(scheduled_date);
CREATE INDEX idx_schedule_platform ON content_schedule(platform_id);
CREATE INDEX idx_schedule_status ON content_schedule(status);
```

## API Endpoints

### Content Management

#### `GET /api/content/posts`
List all content posts with filtering.

**Query Parameters:**
- `project`: Filter by project
- `content_type`: Filter by type
- `status`: Filter by status
- `priority`: Filter by priority
- `page`: Pagination
- `limit`: Results per page

**Response:**
```json
{
  "posts": [
    {
      "id": "uuid",
      "title": "Microservices Architecture Overview",
      "content_type": "architecture",
      "project": "backend",
      "status": "ready",
      "priority": 0,
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 500,
  "page": 1,
  "limit": 20
}
```

#### `GET /api/content/posts/:id`
Get single content post with related data.

**Response:**
```json
{
  "id": "uuid",
  "title": "Microservices Architecture Overview",
  "content": "...",
  "content_type": "architecture",
  "project": "backend",
  "status": "ready",
  "priority": 0,
  "social_posts": [...],
  "assets": [...],
  "created_at": "2024-01-15T10:00:00Z"
}
```

#### `POST /api/content/posts`
Create new content post.

**Request:**
```json
{
  "title": "New Post Title",
  "slug": "new-post-slug",
  "content": "Post content...",
  "content_type": "architecture",
  "project": "backend",
  "category": "architecture",
  "priority": 0,
  "metadata": {}
}
```

#### `PUT /api/content/posts/:id`
Update content post.

#### `DELETE /api/content/posts/:id`
Delete content post (soft delete recommended).

### Social Posts Management

#### `GET /api/social/posts`
List social posts with filtering.

**Query Parameters:**
- `platform`: Filter by platform
- `status`: Filter by status
- `content_post_id`: Filter by source content
- `scheduled_date`: Filter by scheduled date
- `page`, `limit`: Pagination

#### `POST /api/social/posts`
Create social post from content post.

**Request:**
```json
{
  "content_post_id": "uuid",
  "platform_id": "uuid",
  "title": "Platform-specific title",
  "content": "Platform-adapted content...",
  "format": "long-form",
  "word_count": 600,
  "image_count": 2,
  "scheduled_at": "2024-01-16T08:00:00Z",
  "metadata": {}
}
```

#### `PUT /api/social/posts/:id`
Update social post.

#### `POST /api/social/posts/:id/publish`
Mark social post as published.

**Request:**
```json
{
  "platform_post_id": "external-id",
  "url": "https://platform.com/post/123",
  "posted_at": "2024-01-16T08:00:00Z"
}
```

### Scheduling

#### `GET /api/schedule`
Get schedule for date range.

**Query Parameters:**
- `start_date`: Start date (ISO format)
- `end_date`: End date (ISO format)
- `platform_id`: Filter by platform

**Response:**
```json
{
  "schedule": [
    {
      "date": "2024-01-16",
      "platform": "linkedin",
      "posts": [
        {
          "id": "uuid",
          "title": "Post Title",
          "scheduled_time": "08:00:00",
          "status": "scheduled"
        }
      ]
    }
  ]
}
```

#### `POST /api/schedule`
Schedule content for posting.

**Request:**
```json
{
  "content_post_id": "uuid",
  "platform_id": "uuid",
  "scheduled_date": "2024-01-16",
  "scheduled_time": "08:00:00"
}
```

#### `PUT /api/schedule/:id`
Update schedule.

#### `DELETE /api/schedule/:id`
Cancel scheduled post.

### Analytics

#### `GET /api/analytics/posts/:id`
Get analytics for specific post.

**Query Parameters:**
- `start_date`: Start date
- `end_date`: End date

**Response:**
```json
{
  "post_id": "uuid",
  "platform": "linkedin",
  "metrics": [
    {
      "date": "2024-01-16",
      "likes": 50,
      "comments": 10,
      "shares": 5,
      "views": 500,
      "engagement_rate": 13.0
    }
  ],
  "totals": {
    "likes": 50,
    "comments": 10,
    "shares": 5,
    "views": 500
  }
}
```

#### `GET /api/analytics/summary`
Get analytics summary.

**Query Parameters:**
- `start_date`: Start date
- `end_date`: End date
- `platform_id`: Filter by platform
- `content_type`: Filter by content type

**Response:**
```json
{
  "period": {
    "start": "2024-01-01",
    "end": "2024-01-31"
  },
  "totals": {
    "posts": 20,
    "likes": 1000,
    "comments": 200,
    "shares": 100,
    "views": 10000
  },
  "by_platform": {...},
  "by_content_type": {...},
  "top_posts": [...]
}
```

#### `POST /api/analytics/sync`
Sync analytics from platforms (manual or scheduled).

**Request:**
```json
{
  "platform_id": "uuid",
  "start_date": "2024-01-01",
  "end_date": "2024-01-31"
}
```

### Content Repurposing

#### `POST /api/repurpose`
Repurpose content to another platform.

**Request:**
```json
{
  "source_content_post_id": "uuid",
  "source_social_post_id": "uuid", // Optional
  "target_platform_id": "uuid",
  "repurposing_type": "summary", // 'summary', 'expansion', 'translation', 'visual'
  "customizations": {
    "title": "Custom title",
    "content": "Custom content..."
  }
}
```

#### `GET /api/repurpose/:content_post_id`
Get repurposing history for content post.

### Backlog Management

#### `GET /api/backlog`
Get content backlog with prioritization.

**Query Parameters:**
- `priority`: Filter by priority
- `status`: Filter by status
- `project`: Filter by project
- `limit`: Number of items

**Response:**
```json
{
  "backlog": [
    {
      "id": "uuid",
      "title": "Post Title",
      "content_type": "architecture",
      "priority": 0,
      "status": "ready",
      "estimated_platforms": 5,
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "stats": {
    "total": 500,
    "ready": 50,
    "in_progress": 5,
    "scheduled": 10,
    "posted": 435
  }
}
```

#### `POST /api/backlog/prioritize`
Update priority for content posts.

**Request:**
```json
{
  "content_post_ids": ["uuid1", "uuid2"],
  "priority": 0
}
```

## Domain Models (Go)

### Content Post
```go
type ContentPost struct {
    ID          uuid.UUID
    Title       string
    Slug        string
    Content     string
    ContentType string // architecture, implementation, lessons, decision, advanced
    Project     string
    Category    string
    Priority    int // 0 (P0) to 3 (P3)
    Status      string // draft, ready, scheduled, posted, archived
    Metadata    datatypes.JSON
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Social Post
```go
type SocialPost struct {
    ID            uuid.UUID
    ContentPostID uuid.UUID
    PlatformID    uuid.UUID
    Title         string
    Content       string
    Format        string // long-form, thread, carousel, article, newsletter
    WordCount     int
    ImageCount    int
    ScheduledAt   *time.Time
    PostedAt      *time.Time
    Status        string // draft, scheduled, posted, archived
    PlatformPostID string
    URL           string
    Metadata      datatypes.JSON
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### Post Analytics
```go
type PostAnalytics struct {
    ID            uuid.UUID
    SocialPostID  uuid.UUID
    PlatformID    uuid.UUID
    MetricDate    time.Time
    Likes         int
    Comments      int
    Shares        int
    Views         int
    Clicks        int
    EngagementRate float64
    Reach         int
    Impressions   int
    Metadata      datatypes.JSON
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

## Features to Implement

### Phase 1: Core Content Management
- [ ] Content post CRUD operations
- [ ] Social post creation from content posts
- [ ] Basic scheduling
- [ ] Status tracking
- [ ] Backlog management

### Phase 2: Platform Integration
- [ ] Platform configuration
- [ ] Platform-specific content adaptation
- [ ] Scheduling with platform timing
- [ ] Post status synchronization

### Phase 3: Analytics
- [ ] Analytics data collection
- [ ] Analytics API endpoints
- [ ] Analytics dashboard
- [ ] Performance metrics

### Phase 4: Repurposing
- [ ] Content repurposing logic
- [ ] Platform-specific adaptations
- [ ] Repurposing tracking
- [ ] Automated repurposing suggestions

### Phase 5: Automation
- [ ] Automated scheduling
- [ ] Analytics sync (scheduled)
- [ ] Content recommendations
- [ ] Performance alerts

## Integration Points

### External APIs (Future)
- LinkedIn API (for posting and analytics)
- Twitter API (for posting and analytics)
- Instagram API (for posting and analytics)
- Medium API (for posting)
- Substack API (for newsletter)

### Internal Services
- **Content Service**: Manage content posts
- **Scheduling Service**: Handle post scheduling
- **Analytics Service**: Collect and analyze metrics
- **Repurposing Service**: Handle content adaptation

## Data Migration

### Initial Data Load
1. Import all 140+ backend posts
2. Import posts from other 5 projects
3. Create platform records
4. Generate initial social posts
5. Set up initial schedule

### Migration Script
```sql
-- Example: Import backend posts
INSERT INTO content_posts (title, slug, content, content_type, project, category, priority, status)
SELECT 
    title,
    slug,
    content,
    content_type,
    'backend' as project,
    category,
    CASE 
        WHEN content_type IN ('architecture', 'decision') THEN 0
        WHEN content_type IN ('implementation', 'lessons') THEN 1
        ELSE 2
    END as priority,
    'ready' as status
FROM backend_posts;
```

## Usage Examples

### Create Content Post
```bash
POST /api/content/posts
{
  "title": "Microservices Architecture Overview",
  "slug": "microservices-architecture-overview",
  "content": "...",
  "content_type": "architecture",
  "project": "backend",
  "category": "architecture",
  "priority": 0
}
```

### Create Social Posts from Content
```bash
POST /api/social/posts/repurpose
{
  "content_post_id": "uuid",
  "platforms": ["linkedin", "twitter", "instagram"],
  "auto_adapt": true
}
```

### Schedule Posts
```bash
POST /api/schedule
{
  "content_post_id": "uuid",
  "platform_id": "linkedin-uuid",
  "scheduled_date": "2024-01-16",
  "scheduled_time": "08:00:00"
}
```

### Get Analytics
```bash
GET /api/analytics/summary?start_date=2024-01-01&end_date=2024-01-31
```

## Future Enhancements

### AI Integration
- Content adaptation suggestions
- Optimal posting time recommendations
- Content performance predictions
- Automated content generation

### Advanced Analytics
- Content performance trends
- Platform comparison
- Audience insights
- ROI calculations

### Workflow Automation
- Automated content repurposing
- Smart scheduling based on performance
- Content recommendation engine
- Performance-based prioritization

## Implementation Notes

### Technology Stack
- **Backend**: Go (existing server)
- **Database**: PostgreSQL (existing)
- **ORM**: GORM (existing)
- **API**: Fiber (existing)

### New Domains
- `contentmanagement` domain
- `socialmedia` domain
- `analytics` domain
- `scheduling` domain

### Dependencies
- UUID generation
- JSONB support (existing)
- Time handling
- External API clients (future)

## Testing Strategy

### Unit Tests
- Content post operations
- Social post creation
- Scheduling logic
- Analytics calculations

### Integration Tests
- End-to-end content workflow
- Platform integration
- Analytics collection
- Repurposing logic

## Security Considerations

### Authentication
- JWT authentication (existing)
- Role-based access control
- API key support (for external integrations)

### Data Privacy
- Secure storage of platform credentials
- Encrypted API keys
- Audit logging

## Monitoring

### Metrics to Track
- Content creation rate
- Posting success rate
- Analytics sync success
- API response times
- Error rates

### Alerts
- Failed postings
- Analytics sync failures
- High error rates
- Performance degradation
