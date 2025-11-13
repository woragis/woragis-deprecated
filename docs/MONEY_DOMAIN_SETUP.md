# Money Domain Setup

## Overview
The Money domain has been created with two main features:
1. **Ideas** - Store and manage ideas with markdown documents
2. **Idea Nodes** - Create connected nodes in a 2D canvas for each idea
3. **AI Chats** - AI chat interface that references idea nodes

## Created Files

### Database Schemas
- `src/server/db/schemas/money/ideas.ts` - Ideas table schema
- `src/server/db/schemas/money/idea-nodes.ts` - Idea nodes table schema (for 2D canvas)
- `src/server/db/schemas/money/ai-chats.ts` - AI chats table schema
- `src/server/db/schemas/money/index.ts` - Export file

### Types
- `src/types/money/ideas.ts` - Idea types and interfaces
- `src/types/money/idea-nodes.ts` - Idea node types (canvas positions, connections)
- `src/types/money/ai-chats.ts` - AI chat types and message structures
- `src/types/money/index.ts` - Export file

### Repositories
- `src/server/repositories/money/ideas.repository.ts` - Ideas data access layer
- `src/server/repositories/money/idea-nodes.repository.ts` - Nodes data access with canvas operations
- `src/server/repositories/money/ai-chats.repository.ts` - Chats data access with message management
- `src/server/repositories/money/index.ts` - Export file with instances

### Services
- `src/server/services/money/ideas.service.ts` - Ideas business logic
- `src/server/services/money/idea-nodes.service.ts` - Nodes business logic
- `src/server/services/money/ai-chats.service.ts` - Chats business logic
- `src/server/services/money/index.ts` - Export file with instances

### Hooks (React Query)
- `src/hooks/money/useIdeas.ts` - React hooks for ideas CRUD
- `src/hooks/money/useIdeaNodes.ts` - React hooks for nodes with canvas operations
- `src/hooks/money/useAiChats.ts` - React hooks for chats with messaging
- `src/hooks/money/index.ts` - Export file

### Admin Pages & Components
- `src/app/admin/money/ideas/page.tsx` - Ideas management page
- `src/components/pages/admin/money/IdeaForm.tsx` - Form for creating/editing ideas
- `src/components/pages/admin/money/IdeaCanvas.tsx` - 2D canvas for managing idea nodes
- `src/components/pages/admin/money/AiChatInterface.tsx` - Chat interface component
- `src/components/pages/admin/money/index.ts` - Export file

## Database Schema Details

### Ideas Table
- `id` - UUID primary key
- `userId` - Foreign key to users
- `title` - Idea title
- `slug` - URL-friendly slug
- `document` - Markdown content
- `description` - Short description
- `featured`, `visible`, `public` - Display flags
- `order` - Sort order
- Timestamps

### Idea Nodes Table
- `id` - UUID primary key
- `ideaId` - Foreign key to ideas
- `title`, `content` - Node data
- `type` - Node type (default, note, task, etc.)
- `positionX`, `positionY` - Canvas position
- `width`, `height` - Node dimensions
- `color` - Visual customization
- `connections` - Array of connected node IDs (JSONB)
- `visible` - Display flag
- Timestamps

### AI Chats Table
- `id` - UUID primary key
- `userId` - Foreign key to users
- `ideaNodeId` - Foreign key to idea_nodes (many-to-one)
- `title` - Chat title
- `messages` - Array of chat messages (JSONB)
- `agent` - AI agent/model (gpt-4, claude, etc.)
- `model`, `systemPrompt`, `temperature`, `maxTokens` - AI config
- `visible`, `archived` - Display flags
- Timestamps

## Next Steps

### 1. Database Migration
Create a migration file in `drizzle/` to add the new tables:

```bash
npm run db:generate
npm run db:migrate
```

### 2. Create API Routes
You need to create API routes in `src/app/api/money/`:

#### Ideas API Routes
- `src/app/api/money/ideas/route.ts` - GET (list), POST (create)
- `src/app/api/money/ideas/[id]/route.ts` - GET, PATCH, DELETE
- `src/app/api/money/ideas/[id]/toggle-visibility/route.ts` - PATCH
- `src/app/api/money/ideas/[id]/toggle-featured/route.ts` - PATCH
- `src/app/api/money/ideas/order/route.ts` - PATCH
- `src/app/api/money/ideas/stats/route.ts` - GET

#### Idea Nodes API Routes
- `src/app/api/money/idea-nodes/route.ts` - GET, POST
- `src/app/api/money/idea-nodes/[id]/route.ts` - GET, PATCH, DELETE
- `src/app/api/money/idea-nodes/[id]/position/route.ts` - PATCH
- `src/app/api/money/idea-nodes/[id]/connections/route.ts` - PATCH
- `src/app/api/money/idea-nodes/[id]/toggle-visibility/route.ts` - PATCH
- `src/app/api/money/idea-nodes/positions/route.ts` - PATCH (bulk update)
- `src/app/api/money/idea-nodes/stats/route.ts` - GET

#### AI Chats API Routes
- `src/app/api/money/ai-chats/route.ts` - GET, POST
- `src/app/api/money/ai-chats/[id]/route.ts` - GET, PATCH, DELETE
- `src/app/api/money/ai-chats/[id]/message/route.ts` - POST (send message)
- `src/app/api/money/ai-chats/[id]/messages/route.ts` - DELETE (clear)
- `src/app/api/money/ai-chats/[id]/toggle-visibility/route.ts` - PATCH
- `src/app/api/money/ai-chats/[id]/toggle-archived/route.ts` - PATCH
- `src/app/api/money/ai-chats/stats/route.ts` - GET

### 3. Example API Route Structure

```typescript
// src/app/api/money/ideas/route.ts
import { NextRequest, NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { ideasService } from "@/server/services/money";

export async function GET(request: NextRequest) {
  const session = await getServerSession();
  if (!session?.user?.id) {
    return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  }

  const result = await ideasService.getAllIdeas(session.user.id);
  
  if (!result.success) {
    return NextResponse.json({ error: result.error }, { status: 400 });
  }

  return NextResponse.json({ success: true, data: result.data });
}

export async function POST(request: NextRequest) {
  const session = await getServerSession();
  if (!session?.user?.id) {
    return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  }

  const data = await request.json();
  const result = await ideasService.createIdea(data, session.user.id);
  
  if (!result.success) {
    return NextResponse.json({ error: result.error }, { status: 400 });
  }

  return NextResponse.json({ success: true, data: result.data });
}
```

### 4. Update Schema Exports
Make sure `src/server/db/schema.ts` exports the new tables for Drizzle migrations.

### 5. AI Integration
For the AI chat feature, you'll need to:
- Integrate with an AI service (OpenAI, Anthropic, etc.)
- Create a message handler that:
  1. Takes user message
  2. Sends to AI with system prompt and config
  3. Stores both user and AI messages
  4. Returns response

Example location: `src/lib/ai/chat-handler.ts`

### 6. Canvas Enhancements
The canvas component can be enhanced with:
- Drag and drop for nodes
- Connection lines between nodes
- Zoom controls
- Pan controls
- Node resizing
- Different node types with visual variations

Consider using libraries like:
- `react-flow` or `xyflow` for advanced canvas features
- `react-markdown` for rendering markdown in nodes

### 7. Testing
Create tests for:
- Repository CRUD operations
- Service business logic
- API endpoints
- React hooks

## Features Overview

### Ideas Feature
- Create, read, update, delete ideas
- Markdown document support
- Slug-based URLs
- Featured/visible/public flags
- Order management
- Search and filtering

### Idea Nodes Feature
- Create nodes on a 2D canvas
- Position tracking (x, y coordinates)
- Node connections (relationships)
- Node types and colors
- Canvas viewport (pan, zoom)
- Bulk position updates

### AI Chats Feature
- Multiple chats per node (many-to-one)
- Message history stored as JSON
- Configurable AI agents (GPT-4, Claude, etc.)
- Custom system prompts
- Temperature and token controls
- Archive functionality

## Relations
```
User (1) ----< Ideas (many)
Idea (1) ----< IdeaNodes (many)
IdeaNode (1) ----< AiChats (many)
User (1) ----< AiChats (many)
```

## File Structure Summary
```
src/
├── server/
│   ├── db/schemas/money/          # Database schemas
│   ├── repositories/money/        # Data access layer
│   └── services/money/            # Business logic
├── types/money/                   # TypeScript types
├── hooks/money/                   # React Query hooks
├── app/admin/money/               # Admin pages
└── components/pages/admin/money/  # Admin components
```
