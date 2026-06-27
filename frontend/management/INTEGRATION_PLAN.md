# Management Backend-Frontend Integration Plan

## Overview

Connect all 12 management backend domains to the frontend with a sidebar navigation structure.

---

## Backend Domains & Routes

### 1. **Projects** (`/projects`)

**Routes:**

- `POST /projects` - Create project
- `GET /projects` - List projects
- `GET /projects/slug/:slug` - Get by slug
- `PATCH /projects/:id/status` - Update status
- `DELETE /projects/:id` - Delete project
- `POST /projects/:id/milestones` - Add milestone
- `GET /projects/:id/milestones` - List milestones
- `GET /projects/:id/technologies` - List technologies
- `GET /projects/:id/kanban` - Get kanban board
- `POST /projects/:id/kanban/columns` - Create kanban column
- `POST /projects/:id/kanban/cards` - Create kanban card
- `POST /projects/:id/documentation` - Create documentation
- `GET /projects/:id/documentation` - Get documentation

**Frontend Pages:**

- `/projects` - List view with create form
- `/projects/[id]` - Detail view with nested tabs (milestones, tech, kanban, docs)

**Data Model:**

```typescript
Project {
  id: string
  name: string
  slug: string
  description: string
  status: 'planning' | 'in_progress' | 'on_hold' | 'completed' | 'archived'
  startDate: string
  endDate?: string
  visibility: 'private' | 'public' | 'shared'
  collaborators: string[]
  metrics: { completed: number, total: number }
  createdAt: string
  updatedAt: string
}
```

---

### 2. **Ideas** (`/ideas`)

**Routes:**

- `POST /ideas` - Create idea
- `GET /ideas` - List ideas
- `GET /ideas/slug/:slug` - Get by slug
- `PATCH /ideas/:id` - Update idea
- `DELETE /ideas/:id` - Delete idea
- `POST /ideas/bulk/move` - Bulk move
- `POST /ideas/bulk/update` - Bulk update
- `POST /ideas/bulk/delete` - Bulk delete
- `POST /ideas/nodes` - Create idea node
- `GET /ideas/:id/nodes` - Get nodes
- `PATCH /ideas/nodes/:id` - Update node

**Frontend Pages:**

- `/ideas` - List with create form
- `/ideas/[id]` - Canvas view with draggable nodes

**Data Model:**

```typescript
Idea {
  id: string
  title: string
  slug: string
  description: string
  category: string
  position: number
  nodes: IdeaNode[]
  links: IdeaLink[]
  createdAt: string
  updatedAt: string
}

IdeaNode {
  id: string
  ideaId: string
  title: string
  x: number
  y: number
  width: number
  height: number
}
```

---

### 3. **Clients** (`/clients`)

**Routes:**

- `POST /clients` - Create client
- `GET /clients` - List clients
- `GET /clients/:id` - Get by ID
- `PATCH /clients/:id` - Update client
- `DELETE /clients/:id` - Delete client
- `PATCH /clients/:id/archive` - Archive client

**Frontend Pages:**

- `/clients` - List with create form

**Data Model:**

```typescript
Client {
  id: string
  name: string
  email: string
  phone: string
  company: string
  address?: string
  isArchived: boolean
  createdAt: string
  updatedAt: string
}
```

---

### 4. **Finances** (`/finance`) [Note: endpoint is `/finance` not `/finances`]

**Routes:**

- `POST /finance` - Create transaction
- `GET /finance` - List transactions
- `GET /finance/:id` - Get by ID
- `PATCH /finance/:id` - Update transaction
- `DELETE /finance/:id` - Delete transaction
- `GET /finance/analytics` - Get analytics
- `GET /finance/reports` - Get reports

**Frontend Pages:**

- `/finances` - List with filters and analytics dashboard

**Data Model:**

```typescript
Transaction {
  id: string
  type: 'income' | 'expense'
  category: string
  description: string
  amount: number
  currency: string
  baseCurrency: string
  normalizedAmount: number
  occurredAt: string
  isRecurring: boolean
  isEssential: boolean
  tags: string[]
  createdAt: string
  updatedAt: string
}
```

---

### 5. **Experiences** (`/experiences`)

**Routes:**

- `POST /experiences` - Create experience
- `GET /experiences` - List experiences
- `GET /experiences/:id` - Get by ID
- `PATCH /experiences/:id` - Update experience
- `DELETE /experiences/:id` - Delete experience

**Frontend Pages:**

- `/experiences` - List with create form
- `/experiences/[id]` - Detail view

**Data Model:**

```typescript
Experience {
  id: string
  title: string
  company: string
  position: string
  description: string
  startDate: string
  endDate?: string
  isCurrentPosition: boolean
  skills: string[]
  createdAt: string
  updatedAt: string
}
```

---

### 6. **Languages** (`/languages`)

**Routes:**

- `POST /languages` - Create language
- `GET /languages` - List languages
- `PATCH /languages/:id` - Update language
- `DELETE /languages/:id` - Delete language

**Frontend Pages:**

- `/languages` - List with add/edit inline

**Data Model:**

```typescript
Language {
  id: string
  name: string
  proficiency: 'basic' | 'intermediate' | 'fluent' | 'native'
  yearsOfExperience: number
  createdAt: string
  updatedAt: string
}
```

---

### 7. **Testimonials** (`/testimonials`)

**Routes:**

- `POST /testimonials` - Create testimonial
- `GET /testimonials` - List testimonials
- `PATCH /testimonials/:id` - Update testimonial
- `DELETE /testimonials/:id` - Delete testimonial

**Frontend Pages:**

- `/testimonials` - List with create form

**Data Model:**

```typescript
Testimonial {
  id: string
  author: string
  position: string
  company: string
  content: string
  rating: number
  avatar?: string
  verified: boolean
  createdAt: string
  updatedAt: string
}
```

---

### 8. **User Profiles** (`/user-profiles`)

**Routes:**

- `GET /user-profiles` - Get current user profile
- `PATCH /user-profiles` - Update profile
- `GET /user-profiles/:userId` - Get public profile

**Frontend Pages:**

- `/profile` - User profile view and edit

**Data Model:**

```typescript
UserProfile {
  id: string
  userId: string
  avatar?: string
  bio?: string
  dateOfBirth?: string
  gender?: string
  phone?: string
  location?: string
  website?: string
  socialLinks?: Record<string, string>
  preferences?: Record<string, any>
  createdAt: string
  updatedAt: string
}
```

---

### 9. **User Preferences** (`/user-preferences`)

**Routes:**

- `GET /user-preferences` - Get preferences
- `PATCH /user-preferences` - Update preferences

**Frontend Pages:**

- `/settings/preferences` - Settings dashboard

**Data Model:**

```typescript
UserPreferences {
  id: string
  userId: string
  theme: 'light' | 'dark' | 'auto'
  language: string
  notifications: Record<string, boolean>
  privacy: Record<string, boolean>
  createdAt: string
  updatedAt: string
}
```

---

### 10. **API Keys** (`/api-keys`)

**Routes:**

- `POST /api-keys` - Create API key
- `GET /api-keys` - List API keys
- `DELETE /api-keys/:id` - Delete API key
- `PATCH /api-keys/:id/rotate` - Rotate key

**Frontend Pages:**

- `/settings/api-keys` - API key management

**Data Model:**

```typescript
APIKey {
  id: string
  name: string
  key: string (only on creation)
  lastUsed?: string
  expiresAt?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}
```

---

### 11. **Scheduler** (`/scheduler`)

**Routes:**

- `POST /scheduler` - Schedule task
- `GET /scheduler` - List scheduled tasks
- `PATCH /scheduler/:id` - Update task
- `DELETE /scheduler/:id` - Delete task

**Frontend Pages:**

- `/scheduler` - Calendar view with task management

**Data Model:**

```typescript
ScheduledTask {
  id: string
  title: string
  description: string
  scheduledFor: string
  reminderBefore?: number (minutes)
  isCompleted: boolean
  priority: 'low' | 'medium' | 'high'
  createdAt: string
  updatedAt: string
}
```

---

### 12. **Certifications** (`/certifications`)

**Routes:**

- `POST /certifications` - Create certification
- `GET /certifications` - List certifications
- `PATCH /certifications/:id` - Update certification
- `DELETE /certifications/:id` - Delete certification

**Frontend Pages:**

- `/certifications` - List with create form

**Data Model:**

```typescript
Certification {
  id: string
  name: string
  issuer: string
  issueDate: string
  expiryDate?: string
  credentialUrl?: string
  credentialId?: string
  createdAt: string
  updatedAt: string
}
```

---

### 13. **Chats** (`/chats`) [Requires AI Service Integration]

**Routes:**

- `POST /chats` - Create chat
- `GET /chats` - List chats
- `POST /chats/:id/messages` - Send message
- `GET /chats/:id/messages` - Get chat history
- `DELETE /chats/:id` - Delete chat

**Frontend Pages:**

- `/chats` - Chat list
- `/chats/[id]` - Chat interface with AI responses

**Data Model:**

```typescript
Chat {
  id: string
  title: string
  messages: ChatMessage[]
  createdAt: string
  updatedAt: string
}

ChatMessage {
  id: string
  chatId: string
  role: 'user' | 'assistant'
  content: string
  timestamp: string
}
```

---

## Sidebar Navigation Structure

```
Management Portal
├── Dashboard
│   └── /dashboard
├── Projects
│   └── /projects
├── Ideas
│   └── /ideas
├── Clients
│   └── /clients
├── Finances
│   └── /finances
├── Career
│   ├── Experiences
│   │   └── /experiences
│   ├── Languages
│   │   └── /languages
│   ├── Certifications
│   │   └── /certifications
│   └── Testimonials
│       └── /testimonials
├── Planning
│   ├── Scheduler
│   │   └── /scheduler
│   └── Chats (AI Assistant)
│       └── /chats
├── Account
│   ├── Profile
│   │   └── /profile
│   └── Settings
│       ├── Preferences → /settings/preferences
│       └── API Keys → /settings/api-keys
└── Logout
```

---

## Frontend File Structure to Create

```
src/routes/
├── dashboard/
│   └── +page.svelte
├── projects/
│   ├── +page.svelte (list)
│   └── [id]/
│       └── +page.svelte (detail with tabs)
├── ideas/
│   ├── +page.svelte (list)
│   └── [id]/
│       └── +page.svelte (canvas view)
├── clients/
│   └── +page.svelte (list)
├── finances/
│   └── +page.svelte (list with analytics)
├── experiences/
│   ├── +page.svelte (list)
│   └── [id]/
│       └── +page.svelte (detail)
├── languages/
│   └── +page.svelte (list)
├── certifications/
│   └── +page.svelte (list)
├── testimonials/
│   └── +page.svelte (list)
├── scheduler/
│   └── +page.svelte (calendar)
├── chats/
│   ├── +page.svelte (list)
│   └── [id]/
│       └── +page.svelte (chat interface)
├── profile/
│   └── +page.svelte
└── settings/
    ├── preferences/
    │   └── +page.svelte
    └── api-keys/
        └── +page.svelte

src/lib/api/
├── dashboard/ (new)
├── experiences/ (new)
├── finances/ (new)
├── languages/ (new)
├── certifications/ (new)
├── testimonials/ (new)
├── chats/ (new)
├── scheduler/ (new)
├── api-keys/ (new)
├── user-preferences/ (already exists)
├── user-profiles/ (already exists)
├── projects/ (already exists)
├── ideas/ (already exists)
├── clients/ (already exists)

src/lib/components/
├── Sidebar.svelte (new)
├── Navigation.svelte (new - top nav)
├── Breadcrumbs.svelte (new)
```

---

## API Clients to Create/Update

### Existing (to verify/enhance):

- `projects` - ✅ Exists, update if needed
- `ideas` - ✅ Exists, update if needed
- `clients` - ✅ Exists, update if needed

### New to Create:

1. `dashboard` - For dashboard data aggregation
2. `experiences` - CRUD operations
3. `finances` - Transactions with filtering/analytics
4. `languages` - CRUD operations
5. `certifications` - CRUD operations
6. `testimonials` - CRUD operations
7. `chats` - Messages and chat management
8. `scheduler` - Task scheduling
9. `api-keys` - Key management

---

## Implementation Phases

### Phase 1: Layout & Navigation

- [ ] Create Sidebar component
- [ ] Create Navigation component
- [ ] Update main layout structure
- [ ] Implement responsive sidebar toggle

### Phase 2: Core Domains (CRUD)

- [ ] Experiences
- [ ] Languages
- [ ] Certifications
- [ ] Testimonials
- [ ] Clients (enhance if needed)

### Phase 3: Complex Domains

- [ ] Finances (with filtering & analytics)
- [ ] Scheduler (calendar integration)
- [ ] Dashboard (data aggregation)

### Phase 4: Advanced Features

- [ ] Chats (AI integration)
- [ ] Settings (Preferences & API Keys)
- [ ] Profile management

---

## Request/Response Patterns

### List Endpoint

```typescript
// Request
GET /api/v1/{domain}?page=1&limit=10

// Response
{
  success: true,
  data: [...],
  pagination: {
    page: 1,
    limit: 10,
    total: 100,
    totalPages: 10
  }
}
```

### Create Endpoint

```typescript
// Request
POST /api/v1/{domain}
X-CSRF-Token: <token>
{...payload}

// Response
{
  success: true,
  data: {...created_resource}
}
```

### Update Endpoint

```typescript
// Request
PATCH /api/v1/{domain}/:id
X-CSRF-Token: <token>
{...payload}

// Response
{
  success: true,
  data: {...updated_resource}
}
```

### Delete Endpoint

```typescript
// Request
DELETE /api/v1/{domain}/:id
X-CSRF-Token: <token>

// Response
{
  success: true,
  message: "Resource deleted"
}
```

---

## Notes

- All state-changing requests require CSRF token in header `X-CSRF-Token`
- Token refresh happens automatically on 401
- Base client automatically handles CSRF for POST, PATCH, PUT, DELETE
- All domains use snake_case for response field names (consistent with auth)
