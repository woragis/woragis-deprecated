# Woragis Backend Microservices Analysis

## Executive Summary

This document analyzes the current monolithic Woragis backend and proposes a microservices architecture. The backend currently contains **30+ domains** organized in a single Go service. This analysis identifies logical service boundaries, domain separation strategies, and potential new microservices.

## Current Domain Inventory

### Core Identity & Access
- **auth** - User authentication, sessions, MFA, OAuth, email tokens, audit logs
- **apikeys** - API key management for public endpoints
- **userprofiles** - User profile information
- **userpreferences** - User preferences and settings

### Job Application & Career Management
- **jobapplications** - Job application tracking with subdomains:
  - `interviewstages` - Interview stage management
  - `responses` - Application response tracking
- **resumes** - Resume management and generation
- **jobwebsites** - Job website/platform tracking

### Project & Idea Management
- **projects** - Project management with:
  - `projectcasestudies` - Project case studies
- **ideas** - Idea management with versioning, relationships, and collaboration
- **chats** - AI-powered chat conversations (linked to projects, ideas, job applications)

### Content & Technical Documentation
- **posts** - Blog posts with:
  - `comments` - Post comments
- **technicalwritings** - Technical writing/articles
- **problemsolutions** - Problem-solving documentation
- **systemdesigns** - System design documentation
- **casestudies** - Case studies
- **impactmetrics** - Impact metrics tracking
- **aimlintegrations** - AI/ML integration documentation
- **certifications** - Certification tracking
- **testimonials** - Testimonials and recommendations
- **experiences** - Professional experiences

### Social Media Management
- **socialmediaposts** - Complex domain with subdomains:
  - `content` - Content creation and repurposing
  - `scheduling` - Post scheduling
  - `platforms` - Platform configuration
  - `analytics` - Post analytics
  - `assets` - Content assets

### Supporting & Utility Domains
- **skills** - Skills management
- **interests** - Interests tracking
- **languages** - Language learning tracking
- **finances** - Financial transaction management
- **clients** - Client management
- **scheduler** - Task scheduling
- **reports** - Report generation and delivery
- **translations** - Multi-language translation management
- **creativeassets** - Creative asset management (integrates with external creative-service)

## Proposed Microservices Architecture

### 1. **Auth Service** (`backend/auth/server/`)

**Responsibilities:**
- User authentication and authorization
- Session management
- Multi-factor authentication (MFA)
- OAuth provider integration (Google, GitHub, Microsoft)
- Email confirmation and password reset
- API key management
- Audit logging
- User profiles and preferences

**Domains to Include:**
- `auth` (User, Session, Device, MFA, OAuth, EmailToken, AuditLog)
- `apikeys`
- `userprofiles`
- `userpreferences`

**Database Tables:**
- `users`
- `sessions`
- `devices`
- `mfa_tokens`
- `audit_logs`
- `oauth_accounts`
- `email_tokens`
- `api_keys`
- `user_profiles`
- `user_preferences`

**External Dependencies:**
- Redis (for token storage and sessions)
- SMTP service (for email confirmations)
- OAuth providers (Google, GitHub, Microsoft)

**API Endpoints:**
- `/api/auth/*` - Authentication endpoints
- `/api/api-keys/*` - API key management
- `/api/user-profiles/*` - User profile management
- `/api/user-preferences/*` - User preferences

**Shared with Other Services:**
- JWT tokens for authentication
- User ID references in other services
- API key validation for public endpoints

---

### 2. **Jobs Service** (`backend/jobs/server/`)

**Responsibilities:**
- Job application tracking
- Resume management
- Job website/platform tracking
- Interview stage management
- Application response tracking
- Cover letter generation (via AI service)

**Domains to Include:**
- `jobapplications` (with subdomains)
  - `interviewstages`
  - `responses`
- `resumes`
- `jobwebsites`

**Database Tables:**
- `job_applications`
- `interview_stages`
- `responses`
- `resumes`
- `job_websites`

**External Dependencies:**
- Auth Service (for user validation)
- AI Service (for cover letter generation)
- Creative Service (for resume generation - via resume-worker)
- RabbitMQ (for job application processing)
- Redis (for queue fallback)

**API Endpoints:**
- `/api/job-applications/*` - Job application CRUD
- `/api/job-applications/:id/interview-stages/*` - Interview stages
- `/api/job-applications/:id/responses/*` - Application responses
- `/api/resumes/*` - Resume management
- `/api/job-websites/*` - Job website management
- `/public/resumes/*` - Public resume endpoints

**Shared with Other Services:**
- Job application IDs referenced in chats
- Resume IDs referenced in job applications

---

### 3. **Management Service** (`backend/management/server/`)

**Responsibilities:**
- Project management
- Idea management with relationships
- Chat conversations (AI-powered)
- Client management
- Task scheduling
- Report generation
- Financial transaction tracking

**Domains to Include:**
- `projects`
  - `projectcasestudies`
- `ideas`
- `chats`
- `clients`
- `scheduler`
- `reports`
- `finances`

**Database Tables:**
- `projects`
- `milestones`
- `kanban_columns`
- `kanban_cards`
- `project_dependencies`
- `project_case_studies`
- `ideas`
- `idea_nodes`
- `idea_node_connections`
- `idea_links`
- `idea_versions`
- `idea_collaborators`
- `documents`
- `conversations`
- `messages`
- `conversation_transcripts`
- `conversation_assignments`
- `clients`
- `schedules`
- `execution_runs`
- `report_definitions`
- `report_schedules`
- `report_deliveries`
- `report_runs`
- `transactions`
- `recurring_templates`

**External Dependencies:**
- Auth Service (for user validation)
- AI Service (LangChain for chat conversations)
- Translation Service (for multi-language support)
- RabbitMQ (for report delivery)
- Redis (for queue fallback)

**API Endpoints:**
- `/api/projects/*` - Project management
- `/api/projects/:id/case-studies/*` - Project case studies
- `/api/ideas/*` - Idea management
- `/api/chats/*` - Chat conversations
- `/api/chats/conversations/:id/stream` - WebSocket streaming
- `/api/clients/*` - Client management
- `/api/scheduler/*` - Task scheduling
- `/api/reports/*` - Report generation
- `/api/finances/*` - Financial transactions

**Shared with Other Services:**
- Project IDs referenced in posts, skills, testimonials
- Idea IDs referenced in chats
- Chat conversations linked to job applications

**Special Considerations:**
- Chat service requires context from multiple domains (projects, skills, resumes, etc.)
- Reports service aggregates data from multiple domains
- Scheduler triggers report generation

---

### 4. **Posts Service** (`backend/posts/server/`)

**Responsibilities:**
- Blog posts and articles
- Technical writing
- Problem solutions
- System designs
- Case studies
- Impact metrics
- AI/ML integrations documentation
- Certifications
- Testimonials
- Professional experiences
- Post comments

**Domains to Include:**
- `posts`
  - `comments`
- `technicalwritings`
- `problemsolutions`
- `systemdesigns`
- `casestudies`
- `impactmetrics`
- `aimlintegrations`
- `certifications`
- `testimonials`
- `experiences`

**Database Tables:**
- `posts`
- `post_skills`
- `post_categories`
- `post_tags`
- `categories`
- `tags`
- `comments`
- `technical_writings`
- `problem_solutions`
- `system_designs`
- `case_studies`
- `impact_metrics`
- `aiml_integrations`
- `certifications`
- `certification_skills`
- `certification_entity_links`
- `testimonials`
- `testimonial_entity_links`
- `experiences`
- `experience_technologies`
- `experience_projects`
- `experience_achievements`

**External Dependencies:**
- Auth Service (for user validation and API keys)
- Translation Service (for multi-language support)
- Creative Service (for asset management)
- Skills Service (for skill relationships - see below)

**API Endpoints:**
- `/api/posts/*` - Blog posts
- `/api/posts/:id/comments/*` - Post comments
- `/api/technical-writings/*` - Technical articles
- `/api/problem-solutions/*` - Problem solutions
- `/api/system-designs/*` - System designs
- `/api/case-studies/*` - Case studies
- `/api/impact-metrics/*` - Impact metrics
- `/api/aiml-integrations/*` - AI/ML integrations
- `/api/certifications/*` - Certifications
- `/api/testimonials/*` - Testimonials
- `/api/experiences/*` - Professional experiences

**Shared with Other Services:**
- Posts linked to projects, skills
- Testimonials linked to projects, skills
- Certifications linked to skills

---

### 5. **Social Media Service** (`backend/social-media/server/`)

**Responsibilities:**
- Social media post creation and management
- Content repurposing
- Post scheduling
- Platform configuration
- Analytics tracking
- Asset management

**Domains to Include:**
- `socialmediaposts` (with all subdomains)
  - `content`
  - `scheduling`
  - `platforms`
  - `analytics`
  - `assets`

**Database Tables:**
- `social_media_posts`
- `social_media_entity_links`
- `content_posts`
- `content_repurposings`
- `platform_configs`
- `scheduled_posts`
- `post_analytics`
- `content_assets`

**External Dependencies:**
- Auth Service (for user validation)
- Translation Service (for multi-language support)
- Creative Service (for asset generation)
- RabbitMQ (for scheduled post publishing)
- External social media APIs (LinkedIn, Twitter, Facebook, etc.)

**API Endpoints:**
- `/api/social-media-posts/*` - Post management
- `/api/social-media-posts/content/*` - Content management
- `/api/social-media-posts/scheduling/*` - Scheduling
- `/api/social-media-posts/platforms/*` - Platform configuration
- `/api/social-media-posts/analytics/*` - Analytics
- `/api/social-media-posts/assets/*` - Asset management

**Shared with Other Services:**
- Can link to posts, projects, skills for content repurposing

---

## New Microservices to Consider

### 6. **Skills Service** (`backend/skills/server/`)

**Why Separate:**
- Skills are referenced across multiple services (posts, projects, certifications, testimonials)
- Centralized skill management enables better consistency
- Can provide skill recommendations and analytics

**Responsibilities:**
- Skill CRUD operations
- Skill relationships with projects, posts, certifications
- Skill analytics and recommendations

**Domains to Include:**
- `skills`
- `interests` (closely related to skills)

**Database Tables:**
- `skills`
- `project_skills`
- `interests`

**External Dependencies:**
- Auth Service (for user validation and API keys)
- Translation Service (for multi-language support)

**API Endpoints:**
- `/api/skills/*` - Skill management
- `/api/interests/*` - Interest management

**Shared with Other Services:**
- Skills referenced in:
  - Posts Service (post_skills)
  - Management Service (projects)
  - Posts Service (certifications)

---

### 7. **Translation Service** (`backend/translation/server/`)

**Note:** This already exists as a standalone worker, but could be expanded into a full service.

**Current State:**
- Standalone translation-worker exists
- Uses RabbitMQ for job processing
- Integrates with LangChain for AI translations

**Proposed Enhancement:**
- Convert to full microservice with REST API
- Manage translation metadata and history
- Provide translation status and retrieval endpoints

**Responsibilities:**
- Translation job management
- Translation history and metadata
- Multi-language content management

**Domains to Include:**
- `translations` (enhanced)

**Database Tables:**
- `translations`

**External Dependencies:**
- AI Service (LangChain for translations)
- RabbitMQ (for translation jobs)

**API Endpoints:**
- `/api/translations/*` - Translation management
- `/api/translations/status/:id` - Translation status
- `/api/translations/history` - Translation history

**Shared with Other Services:**
- All services that need multi-language support

---

### 8. **Languages Service** (`backend/languages/server/`)

**Why Separate:**
- Language learning is a distinct domain
- Can be used independently of other features
- May have different scaling requirements

**Responsibilities:**
- Language learning session tracking
- Vocabulary management
- Proficiency tracking

**Domains to Include:**
- `languages`

**Database Tables:**
- `study_sessions`
- `vocabulary_entries`

**External Dependencies:**
- Auth Service (for user validation)

**API Endpoints:**
- `/api/languages/*` - Language learning management

**Shared with Other Services:**
- Can be referenced in user profiles

---

## Cross-Service Communication Strategy

### 1. **Synchronous Communication (HTTP/gRPC)**
- **Auth Service** → Other services: User validation, API key validation
- **Skills Service** → Posts/Management: Skill lookups
- **Translation Service** → All services: Translation requests

### 2. **Asynchronous Communication (RabbitMQ)**
- Job application processing
- Translation jobs
- Report generation
- Scheduled social media posts
- Email/WhatsApp notifications

### 3. **Shared Data Patterns**
- **User IDs**: All services reference users by ID (no user data duplication)
- **Foreign Key References**: Use UUIDs for cross-service references
- **Event Sourcing**: Consider events for cross-service updates (e.g., user profile changes)

## Database Separation Strategy

### Option 1: Database per Service (Recommended)
Each microservice has its own database:
- `auth_db` - Auth Service
- `jobs_db` - Jobs Service
- `management_db` - Management Service
- `posts_db` - Posts Service
- `social_media_db` - Social Media Service
- `skills_db` - Skills Service
- `translation_db` - Translation Service
- `languages_db` - Languages Service

**Benefits:**
- True service isolation
- Independent scaling
- Technology flexibility per service

**Challenges:**
- Cross-service queries require API calls
- Distributed transactions complexity
- Data consistency across services

### Option 2: Shared Database with Schema Separation
Single database with separate schemas:
- `auth` schema
- `jobs` schema
- `management` schema
- `posts` schema
- `social_media` schema
- `skills` schema
- `translation` schema
- `languages` schema

**Benefits:**
- Easier migration path
- Simpler cross-service queries (if needed)
- Single backup/restore process

**Challenges:**
- Less isolation
- Potential for tight coupling
- Harder to scale individual services

## Migration Path

### Phase 1: Extract Auth Service
1. Create `backend/auth/server/`
2. Move auth, apikeys, userprofiles, userpreferences domains
3. Update other services to call Auth Service API
4. Deploy Auth Service
5. Update monolith to use Auth Service

### Phase 2: Extract Jobs Service
1. Create `backend/jobs/server/`
2. Move jobapplications, resumes, jobwebsites domains
3. Update Management Service (chats) to reference job applications
4. Deploy Jobs Service

### Phase 3: Extract Management Service
1. Create `backend/management/server/`
2. Move projects, ideas, chats, clients, scheduler, reports, finances
3. Update chat context builder to use external services
4. Deploy Management Service

### Phase 4: Extract Posts Service
1. Create `backend/posts/server/`
2. Move all content-related domains
3. Update references in other services
4. Deploy Posts Service

### Phase 5: Extract Social Media Service
1. Create `backend/social-media/server/`
2. Move socialmediaposts domain
3. Deploy Social Media Service

### Phase 6: Extract Supporting Services
1. Create Skills Service
2. Create Languages Service
3. Enhance Translation Service
4. Deploy all supporting services

## Service Dependencies Graph

```
Auth Service (No dependencies)
    ↓
    ├─→ Jobs Service
    ├─→ Management Service
    ├─→ Posts Service
    ├─→ Social Media Service
    ├─→ Skills Service
    ├─→ Translation Service
    └─→ Languages Service

Jobs Service
    ↓
    ├─→ Auth Service
    └─→ AI Service (external)

Management Service
    ↓
    ├─→ Auth Service
    ├─→ AI Service (external)
    └─→ Translation Service

Posts Service
    ↓
    ├─→ Auth Service
    ├─→ Skills Service
    └─→ Translation Service

Social Media Service
    ↓
    ├─→ Auth Service
    └─→ Translation Service

Skills Service
    ↓
    └─→ Auth Service

Translation Service
    ↓
    └─→ AI Service (external)

Languages Service
    ↓
    └─→ Auth Service
```

## Recommended Folder Structure

```
woragis/
  backend/
    auth/
      server/
        app/
          cmd/
            server/
              main.go
          internal/
            domains/
              auth/
              apikeys/
              userprofiles/
              userpreferences/
          pkg/
        go.mod
        Dockerfile
    jobs/
      server/
        app/
          cmd/
            server/
              main.go
          internal/
            domains/
              jobapplications/
              resumes/
              jobwebsites/
          pkg/
        go.mod
        Dockerfile
    management/
      server/
        app/
          cmd/
            server/
              main.go
          internal/
            domains/
              projects/
              ideas/
              chats/
              clients/
              scheduler/
              reports/
              finances/
          pkg/
        go.mod
        Dockerfile
    posts/
      server/
        app/
          cmd/
            server/
              main.go
          internal/
            domains/
              posts/
              technicalwritings/
              problemsolutions/
              systemdesigns/
              casestudies/
              impactmetrics/
              aimlintegrations/
              certifications/
              testimonials/
              experiences/
          pkg/
        go.mod
        Dockerfile
    social-media/
      server/
        app/
          cmd/
            server/
              main.go
          internal/
            domains/
              socialmediaposts/
          pkg/
        go.mod
        Dockerfile
    skills/
      server/
        app/
          cmd/
            server/
              main.go
          internal/
            domains/
              skills/
              interests/
          pkg/
        go.mod
        Dockerfile
    translation/
      server/
        app/
          cmd/
            server/
              main.go
          internal/
            domains/
              translations/
          pkg/
        go.mod
        Dockerfile
    languages/
      server/
        app/
          cmd/
            server/
              main.go
          internal/
            domains/
              languages/
          pkg/
        go.mod
        Dockerfile
```

## Shared Packages Strategy

Consider creating a shared packages repository or module:
- `pkg/auth` - JWT validation, user context
- `pkg/logger` - Structured logging
- `pkg/metrics` - Prometheus metrics
- `pkg/tracing` - OpenTelemetry tracing
- `pkg/validation` - Request validation
- `pkg/response` - Standardized API responses

## Summary

### Core Microservices (5)
1. **Auth Service** - Identity and access management
2. **Jobs Service** - Job applications and resumes
3. **Management Service** - Projects, ideas, chats, clients, scheduling, reports, finances
4. **Posts Service** - All content types (posts, technical writing, case studies, etc.)
5. **Social Media Service** - Social media post management

### Supporting Microservices (3)
6. **Skills Service** - Skills and interests management
7. **Translation Service** - Multi-language translation (enhance existing worker)
8. **Languages Service** - Language learning tracking

### Total: 8 Microservices

This architecture provides:
- Clear service boundaries
- Independent scaling
- Technology flexibility
- Easier maintenance
- Better team organization

