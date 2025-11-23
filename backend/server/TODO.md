# Woragis Backend TODO

## Domains Requiring Deeper CRUD / Advanced Workflows
- **Auth**
  - [x] Implement session/device management, multi-factor tokens, audit trails, OAuth provider links, bulk user admin actions.
  - [x] Build email confirmation workflow (token issuance, expiry handling, resend limits, transactional templates).
  - [x] Integrate SMTP provider or notification service for confirmation, password reset, and magic link emails.
  - [x] Create HTML/email templates for account confirmation, password reset, welcome, and session notifications with responsive design.
  - [x] Implement template rendering pipeline (layout partials, localization support, preview tooling).
  - [x] Harden domain services for `RegisterUser`, `ConfirmEmail`, `Login`, `Logout`, `RefreshSession`, `RequestPasswordReset`, `CompletePasswordReset`.
  - [x] Ensure domain validation rules (password policy, unique email, unconfirmed account login handling) emit typed errors and domain events.
  - [x] Document auth use cases with sequence diagrams and acceptance criteria for QA.
  - [ ] Align frontend flows with new auth APIs (confirmation, session management, MFA) and update QA scripts.
- **Finances**
  - [x] Extend reporting (cash-flow projections, tagging).
  - [x] Multi-currency normalization pipeline.
  - [x] Recurring schedule templates for transactions.
  - [x] Bulk upload from CSV/OFX.
  - [ ] Additional forecasting analytics (scenario planning, variance alerts).
- **Languages**
  - [ ] Add spaced-repetition scheduling, bulk vocabulary import/export, proficiency analytics, AI-generated practice sets.
- **Projects**
  - [x] Implement kanban workflows.
  - [x] Build dependency graphs.
  - [x] Bulk milestone updates.
  - [x] Templated project duplication.
  - [ ] **Project Documentation/Architecture Module**
    - [ ] **1. Documentation Storage & Management**
      - [ ] Create `ProjectDocumentation` entity (main container with visibility, versioning metadata).
      - [ ] Create `DocumentationSection` entity (individual sections: Overview/Motivation, Architecture, Tech Stack, File Structure, API Documentation, Deployment Guide, Contributing Guidelines).
      - [ ] Implement versioning system for documentation (track changes over time, similar to ideas domain).
      - [ ] Add CRUD operations for ProjectDocumentation (Create, Read, Update, Delete, List).
      - [ ] Add CRUD operations for DocumentationSection (Create, Read, Update, Delete, Reorder).
      - [ ] Implement section content storage (markdown/text with metadata).
      - [ ] Add validation rules for documentation entities.
      - [ ] Create repository interfaces and implementations for documentation entities.
      - [ ] Create service layer for documentation workflows.
      - [ ] Create HTTP handlers and routes for documentation endpoints.
    - [ ] **2. Architecture Visualization**
      - [ ] Create `ProjectArchitectureDiagram` entity (store diagram metadata, type, content).
      - [ ] Support diagram types: dependency graphs, component diagrams, data flow diagrams, infrastructure diagrams.
      - [ ] Store diagram content (Mermaid, PlantUML syntax, or JSON for custom diagrams).
      - [ ] Add diagram versioning and metadata (title, description, created/updated timestamps).
      - [ ] Implement CRUD operations for architecture diagrams.
      - [ ] Extend existing dependency system to support architecture layer visualization (frontend → backend → database, microservices).
      - [ ] Add repository, service, and handler layers for architecture diagrams.
    - [ ] **3. File Structure Explorer**
      - [ ] Create `ProjectFileStructure` entity (hierarchical file tree storage).
      - [ ] Design tree structure storage (JSON or hierarchical parent-child relationships).
      - [ ] Store file metadata (file purposes, line counts, languages used per file, file paths).
      - [ ] Implement CRUD operations for file structure (Create, Update, Delete tree nodes).
      - [ ] Add bulk import/update operations for file structure.
      - [ ] Create repository, service, and handler layers for file structure.
      - [ ] Add validation for file tree structure integrity.
    - [ ] **4. Technology Stack Tracking**
      - [ ] Create `ProjectTechnology` entity (store technologies with categories, versions, purpose).
      - [ ] Define technology categories (Backend, Database, Frontend, Infrastructure, Monitoring, etc.).
      - [ ] Store technology metadata (name, version, category, purpose/justification, optional link).
      - [ ] Implement CRUD operations for technology stack items.
      - [ ] Add bulk operations for technology stack (add multiple, update, remove).
      - [ ] Create repository, service, and handler layers for technology tracking.
      - [ ] Add validation for technology entries.
    - [ ] **5. Public/Private Visibility**
      - [ ] Add visibility flags to `ProjectDocumentation` (public, authenticated, collaborators only).
      - [ ] Implement access control logic in service layer (check permissions based on visibility).
      - [ ] Create public access endpoints (slug-based URLs like `/projects/{slug}/docs`).
      - [ ] Add middleware or service checks for public documentation access.
      - [ ] Ensure private documentation requires authentication and ownership/collaboration.
      - [ ] Update handlers to respect visibility settings.
    - [ ] **6. Rich Content Support**
      - [ ] Ensure markdown rendering support in DocumentationSection content.
      - [ ] Add support for code snippets with syntax highlighting metadata (language, code block).
      - [ ] Support image/diagram storage or linking (URLs or file references).
      - [ ] Add embedded content support (link to external docs: GitHub README, Notion, etc.).
      - [ ] Create content validation and sanitization for markdown/html content.
    - [ ] **7. Auto-Generated Documentation (Future)**
      - [ ] API endpoint discovery: Automatically document Fiber routes (future enhancement).
      - [ ] Database schema visualization: Generate ER diagrams from GORM models (future enhancement).
      - [ ] Dependency analysis: Show Go module dependencies, Docker image dependencies (future enhancement).
    - [ ] **8. Integration & Testing**
      - [ ] Add error handling and domain errors for documentation module.
      - [ ] Write unit tests for documentation entities and validation.
      - [ ] Write service layer tests for documentation workflows.
      - [ ] Write integration tests for documentation endpoints.
      - [ ] Update project domain routes to include documentation endpoints.
      - [ ] Ensure proper authentication/authorization on all endpoints.
- **Chats**
  - [x] Ship conversation search endpoint with message body filtering.
  - [x] Implement bulk archive/delete/restore workflows with soft-delete metadata.
  - [x] Persist shared transcript exports and assignment history records.
  - [x] Provide websocket streaming hub for live message delivery.
  - [x] Build frontend search UI consuming new chat search API.
  - [x] Wire bulk archive/delete/restore controls in chat list view.
  - [x] Add transcript viewer/share flow in frontend chats module.
  - [x] Integrate websocket client for live message streaming in chat UI.
- **Ideas**
  - [x] Add version history snapshots for edits/moves with history endpoint.
  - [x] Implement bulk move/update/delete/restore APIs with validation.
  - [x] Enable collaborator sharing with role-based access controls.
  - [x] Provide relationship search filters (relation, weight, direction, owner scope).
  - [x] Build Svelte Flow-based canvas UI for ideas graph (frontend).
  - [ ] Ship collaborator & version side panels in frontend with optimistic updates.
  - [ ] Document enhanced ideas workflows and extend regression tests.
- **Scheduler**
  - [x] Extend schedule model with rrule/cron support, priority, channels, and pause state.
  - [x] Implement bulk activation/pausing/resume APIs and authenticated handlers.
  - [x] Track execution history with execution runs and run listing endpoints.
  - [ ] Add alerting hooks for schedule failures and SLA breaches.
- **Reports**
  - [x] Introduce custom report builder workflows (define sections, filters, data sources).
  - [x] Implement scheduling templates for recurring report generation.
  - [x] Enable multi-channel delivery (email, download, share links) with templated messages.
  - [x] Support bulk regeneration/export management for selected reports.
  - [ ] **Implement report definition execution**: Use stored `sections` and `filters` from ReportDefinition during report generation (currently ignored).
  - [ ] **Dynamic data source selection**: Allow report definitions to specify which domains to include (ideas, projects, chats, finances) instead of always querying all.
  - [ ] **Auto-retrieve delivery targets**: Automatically use user's email from auth context for email delivery instead of requiring manual input.
  - [ ] **Add phone number to User entity**: Store phone number in User model or profile table to enable automatic WhatsApp delivery.
  - [ ] **Custom aggregation support**: Allow report definitions to specify custom SQL aggregations, date range filters, and field selections.
  - [ ] **Report generation from definitions**: Create `GenerateReportFromDefinition()` that respects definition's sections/filters and generates custom reports.
  - [ ] **Date range filtering**: Implement date range filtering from definition filters (e.g., "last_30_days", custom ranges) in finance aggregation and other queries.
  - [ ] **Section-based formatting**: Format report output based on definition sections (e.g., only include sections specified in definition).
  - [ ] **Report run tracking**: Link ReportRun entities to specific definitions and store generated output for historical access.
- **Monitoring**
  - [ ] Persist structured events in production DB, alert threshold management, Grafana dashboard provisioning.

## Landing Page Enhancements (Portfolio Optimization)

### High-Impact Features for Hiring
- [x] **Blog Integration & Latest Posts Section** (Complete)
  - [x] Display latest blog posts from backend API (using existing posts domain).
  - [x] Show featured posts prominently in hero or dedicated section.
  - [x] Add "Read More" links to full blog posts.
  - [x] Display post categories/tags for easy filtering.
  - [x] Show reading time and publish date for each post.
  - [ ] Add AI assistant for creating/editing posts: Show an AI helper as a sidebar on CRUD pages to help turn ideas into better, more explanatory posts (the AI translates your thinking into improved content).

- [x] **Interactive Project Showcase** (Complete)
  - [x] Add project filtering by tech stack, type, status (active/completed).
  - [x] Implement search functionality for projects.
  - [x] Add "View Live" / "View Code" buttons with status indicators.
  - [x] Show project metrics (health score displayed).
  - [x] Add project preview cards with hover effects and animations.
  - [x] Display project tech stack badges with icons.
  - [ ] Add AI assistant for creating/editing posts: Show an AI helper as a sidebar on CRUD pages to help turn ideas into better, more explanatory posts (the AI translates your thinking into improved content).

- [x] **Case Studies & Deep Dives** (Complete)
  - [x] Create detailed case study pages for major projects.
  - [x] Include problem statement, solution approach, architecture diagrams.
  - [x] Show before/after metrics, impact, lessons learned.
  - [x] Add interactive architecture diagrams (Mermaid/PlantUML rendering).
  - [x] Link case studies from project cards.
  - [ ] Add AI assistant for creating/editing posts: Show an AI helper as a sidebar on CRUD pages to help turn ideas into better, more explanatory posts (the AI translates your thinking into improved content).
- [ ] - Fazer o pipeline ser automatico e traduzir para todas as 11 linguas. seja para testimonials, projects, posts, ou qualquer outra coisa que ficara publica, tipo descricao de skills or qualquer coisa
- [x] **Testimonials & Recommendations** (Complete)
  - [x] Add testimonials section from colleagues, clients, or mentors (backend domain created).
  - [x] Include LinkedIn recommendations integration (LinkedInURL field in entity).
  - [x] Display recommendation cards with photos, names, roles (AuthorPhoto, AuthorName, AuthorRole, AuthorCompany fields).
  - [x] Add rotating testimonials carousel (frontend implementation complete with autoplay, navigation, and rating display).
  - [ ] Add AI assistant for creating/editing posts: Show an AI helper as a sidebar on CRUD pages to help turn ideas into better, more explanatory posts (the AI translates your thinking into improved content).

- [ ] **Resume/CV Download & Print View**
  - [ ] Add "Download Resume" button (PDF generation or static PDF).
  - [ ] Create print-optimized resume view (`/resume` route).
  - [ ] Ensure resume is ATS-friendly (text-based, structured).
  - [ ] Include QR code linking to portfolio.

- [ ] **Contact Form & Easy Outreach**
  - [ ] Implement contact form with validation.
  - [ ] Add email integration (send emails via backend notification service).
  - [ ] Include social media links (LinkedIn, GitHub, Twitter/X).
  - [ ] Add calendar booking link (Calendly integration).
  - [ ] Display response time expectations.

### Medium-Impact Features
- [ ] **GitHub Activity Integration**
  - [ ] Display GitHub contribution graph.
  - [ ] Show recent commits, repositories, languages used.
  - [ ] Add "View on GitHub" links to projects.
  - [ ] Show GitHub stars, forks, contributions count.

- [ ] **Skills Visualization**
  - [ ] Create interactive skills radar chart or skill tree.
  - [ ] Show proficiency levels (beginner/intermediate/advanced/expert).
  - [ ] Add skill categories (Backend, Frontend, DevOps, etc.).
  - [ ] Link skills to projects that use them.

- [ ] **Achievement Badges & Certifications**
  - [ ] Display certifications (AWS, Google Cloud, etc.).
  - [ ] Show achievement badges (hackathons, open source contributions).
  - [ ] Add certification verification links.
  - [ ] Include completion dates and expiration dates.

- [ ] **Video Introduction**
  - [ ] Add short video introduction (30-60 seconds).
  - [ ] Embed YouTube/Vimeo video or self-hosted.
  - [ ] Include transcript for accessibility.
  - [ ] Add play button overlay on hero section.

- [ ] **Newsletter Signup**
  - [ ] Add newsletter subscription form.
  - [ ] Integrate with backend notification service.
  - [ ] Send welcome email on subscription.
  - [ ] Show subscriber count (if public).

- [ ] **Search Functionality**
  - [ ] Implement site-wide search (projects, blog posts, skills).
  - [ ] Add search bar in header/navigation.
  - [ ] Show search results with categories.
  - [ ] Add keyboard shortcuts (Ctrl+K for search).

### UX/UI Enhancements
- [ ] **Dark Mode Toggle**
  - [ ] Implement dark/light theme switcher.
  - [ ] Persist theme preference in localStorage.
  - [ ] Add smooth theme transition animations.
  - [ ] Ensure all components support both themes.

- [ ] **Accessibility Improvements**
  - [ ] Add ARIA labels to all interactive elements.
  - [ ] Ensure keyboard navigation works throughout.
  - [ ] Add skip-to-content link.
  - [ ] Test with screen readers (NVDA, JAWS).
  - [ ] Ensure color contrast meets WCAG AA standards.
  - [ ] Add focus indicators for keyboard users.

- [ ] **Performance Optimization**
  - [ ] Implement lazy loading for images and components.
  - [ ] Add loading skeletons for API data.
  - [ ] Optimize bundle size (code splitting).
  - [ ] Add service worker for offline support.
  - [ ] Display Lighthouse scores (performance, accessibility, SEO).

- [ ] **Multi-language Support**
  - [ ] Add language switcher (English, Portuguese, etc.).
  - [ ] Implement i18n with SvelteKit.
  - [ ] Translate all content (projects, blog posts, skills).
  - [ ] Store language preference in localStorage.

- [ ] **Responsive Design Enhancements**
  - [ ] Ensure mobile-first design throughout.
  - [ ] Test on various screen sizes (mobile, tablet, desktop).
  - [ ] Add mobile navigation menu (hamburger menu).
  - [ ] Optimize touch targets for mobile.

- [ ] **Animations & Micro-interactions**
  - [ ] Add smooth scroll animations (fade-in, slide-up).
  - [ ] Implement hover effects on cards and buttons.
  - [ ] Add loading animations for API calls.
  - [ ] Include progress indicators for long pages.
  - [ ] Add confetti animation on form submission success.

### SEO & Discoverability
- [ ] **SEO Optimization**
  - [ ] Add meta tags for all pages (title, description, keywords).
  - [ ] Implement Open Graph tags for social sharing.
  - [ ] Add structured data (JSON-LD) for Person, Organization.
  - [ ] Create sitemap.xml and robots.txt.
  - [ ] Add canonical URLs to prevent duplicate content.

- [ ] **Social Media Integration**
  - [ ] Add social sharing buttons (LinkedIn, Twitter, Facebook).
  - [ ] Implement Open Graph preview cards.
  - [ ] Add Twitter Card meta tags.
  - [ ] Create shareable project/blog post links.

- [ ] **Analytics & Tracking**
  - [ ] Integrate Google Analytics or Plausible Analytics.
  - [ ] Track page views, user interactions, form submissions.
  - [ ] Monitor bounce rate, time on page, conversion rates.
  - [ ] Add heatmap tracking (optional, Hotjar/Clarity).

### Content & Engagement
- [ ] **About Me Section Enhancement**
  - [ ] Add personal story, background, journey.
  - [ ] Include professional photo.
  - [ ] Show career timeline or journey.
  - [ ] Add "Why I Code" or mission statement.

- [ ] **Call-to-Action (CTA) Optimization**
  - [ ] Add clear CTAs throughout page (Hire Me, View Projects, Contact).
  - [ ] Use action-oriented button text.
  - [ ] Place CTAs strategically (above fold, after sections).
  - [ ] A/B test different CTA copy and placements.

- [ ] **Content Freshness**
  - [ ] Add "Last Updated" timestamps to projects.
  - [ ] Show recent activity or "What I'm Working On" section.
  - [ ] Display current availability status (Available for Hire, Open to Opportunities).
  - [ ] Add blog post update notifications.

### Technical Demonstrations
- [ ] **Live Code Examples**
  - [ ] Add interactive code playground (CodeSandbox/StackBlitz embed).
  - [ ] Show code snippets with syntax highlighting.
  - [ ] Display GitHub Gists for quick code examples.
  - [ ] Add "Run Code" buttons for executable examples.

- [ ] **API Documentation Integration**
  - [ ] Link to backend API documentation (if public).
  - [ ] Show API endpoint examples.
  - [ ] Display API usage statistics (if public API).

- [ ] **System Architecture Visualization**
  - [ ] Add interactive architecture diagrams.
  - [ ] Show system components and data flow.
  - [ ] Include infrastructure diagrams (AWS, Docker, etc.).
  - [ ] Link to detailed technical blog posts.

## Observability / Platform Enhancements
- Add **WebSocket/Server-Sent Events** pipeline for live metrics and chat streaming.
- Evaluate **gRPC** endpoint layer for Flutter clients; define protobuf schema and gateway.
- Introduce centralized tracing (OpenTelemetry) and log aggregation for cross-service diagnostics.

## Infrastructure Follow-ups
- Dedicated monitoring database connection string for production.
- CI pipeline to build/push Docker images and validate Grafana/Prometheus configs.
- Secrets management for SMTP, AI providers, and future gRPC credentials.
- Comprehensive automated tests (unit/integration) covering new domain workflows before release.
- Kubernetes deployment plan (Helm charts, manifests, GitOps workflow) for production-like environments.
