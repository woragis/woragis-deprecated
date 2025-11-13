# Backend Feature Suggestions

This document outlines the recommended backend features and improvements to ensure full alignment with the mobile implementation and enhance the overall functionality.

## 📊 **Projects Domain**

### **1. Add Missing Fields to Projects Schema**

**Current Issue**: Mobile implementation includes `view_count` and `like_count` operations that don't exist in the backend schema.

**Recommended Changes**:

```typescript
// In src/server/db/schemas/projects.ts
export const projects = pgTable("projects", {
  // ... existing fields ...
  viewCount: integer("view_count").default(0),
  likeCount: integer("like_count").default(0),
  // ... rest of fields ...
});
```

**Migration Required**:
```sql
ALTER TABLE projects 
ADD COLUMN view_count INTEGER DEFAULT 0,
ADD COLUMN like_count INTEGER DEFAULT 0;
```

**API Endpoints to Add**:
```typescript
// src/app/api/admin/projects/[id]/view/route.ts
export const POST = withErrorHandling(async (request, { params }) => {
  const { id } = await params;
  await projectRepository.incrementViewCount(id);
  return handleServiceResult({ success: true }, "View count incremented");
});

// src/app/api/admin/projects/[id]/like/route.ts  
export const POST = withErrorHandling(async (request, { params }) => {
  const { id } = await params;
  await projectRepository.incrementLikeCount(id);
  return handleServiceResult({ success: true }, "Like count incremented");
});
```

**Repository Methods to Add**:
```typescript
// In src/server/repositories/project.repository.ts
async incrementViewCount(id: string): Promise<void> {
  await db
    .update(projects)
    .set({
      viewCount: sql`${projects.viewCount} + 1`,
      updatedAt: new Date(),
    })
    .where(eq(projects.id, id));
}

async incrementLikeCount(id: string): Promise<void> {
  await db
    .update(projects)
    .set({
      likeCount: sql`${projects.likeCount} + 1`,
      updatedAt: new Date(),
    })
    .where(eq(projects.id, id));
}
```

### **2. Add Slug Support for Projects**

**Current Issue**: Mobile implementation supports slug-based project retrieval, but backend doesn't have slug endpoints.

**Recommended Changes**:

**Add Slug Endpoints**:
```typescript
// src/app/api/projects/slug/[slug]/route.ts
export const GET = withErrorHandling(async (request, { params }) => {
  const { slug } = await params;
  const result = await projectService.getPublicProjectBySlug(slug);
  return handleServiceResult(result, "Project fetched successfully");
});

// src/app/api/admin/projects/slug/[slug]/route.ts
export const GET = requireAuth(async (request, { params }, user) => {
  const { slug } = await params;
  const result = await projectService.getProjectBySlug(slug, user.userId);
  return handleServiceResult(result, "Project fetched successfully");
});
```

**Add Repository Methods**:
```typescript
// In src/server/repositories/project.repository.ts
async findBySlug(slug: string, userId?: string): Promise<Project | null> {
  const conditions = [eq(projects.slug, slug)];
  if (userId) {
    conditions.push(eq(projects.userId, userId));
  }
  const result = await db
    .select()
    .from(projects)
    .where(and(...conditions));
  return result[0] || null;
}

async findPublicBySlug(slug: string): Promise<Project | null> {
  const result = await db
    .select()
    .from(projects)
    .where(
      and(
        eq(projects.slug, slug),
        eq(projects.visible, true),
        eq(projects.public, true)
      )
    );
  return result[0] || null;
}
```

**Add Service Methods**:
```typescript
// In src/server/services/project.service.ts
async getProjectBySlug(slug: string, userId?: string): Promise<ApiResponse<Project | null>> {
  try {
    if (!slug || slug.trim().length === 0) {
      return {
        success: false,
        error: "Invalid project slug",
      };
    }
    const project = await projectRepository.findBySlug(slug, userId);
    return this.success(project);
  } catch (error) {
    return this.handleError(error, "getProjectBySlug");
  }
}

async getPublicProjectBySlug(slug: string): Promise<ApiResponse<Project | null>> {
  try {
    if (!slug || slug.trim().length === 0) {
      return {
        success: false,
        error: "Invalid project slug",
      };
    }
    const project = await projectRepository.findPublicBySlug(slug);
    return this.success(project);
  } catch (error) {
    return this.handleError(error, "getPublicProjectBySlug");
  }
}
```

## 🏗️ **Frameworks Domain**

### **3. Add Public Field to Frameworks Schema**

**Current Issue**: Mobile implementation includes a `public` field for frameworks, but backend schema doesn't have this field.

**Recommended Changes**:

```typescript
// In src/server/db/schemas/frameworks.ts
export const frameworks = pgTable("frameworks", {
  // ... existing fields ...
  public: boolean("public").default(true),
  // ... rest of fields ...
});
```

**Migration Required**:
```sql
ALTER TABLE frameworks 
ADD COLUMN public BOOLEAN DEFAULT true;
```

**Update Repository Methods**:
```typescript
// Add public filtering to existing methods
async findPublic(userId?: string): Promise<Framework[]> {
  const conditions = [
    eq(frameworks.visible, true),
    eq(frameworks.public, true)
  ];
  if (userId) {
    conditions.push(eq(frameworks.userId, userId));
  }
  return await db
    .select()
    .from(frameworks)
    .where(and(...conditions))
    .orderBy(asc(frameworks.order), asc(frameworks.name));
}
```

### **4. Standardize Proficiency Levels**

**Current Issue**: Mobile uses enum for proficiency levels, but backend uses strings.

**Recommended Changes**:

**Add Proficiency Enum**:
```typescript
// In src/server/db/schemas/frameworks.ts
export const proficiencyEnum = pgEnum("proficiency_level", [
  "beginner",
  "intermediate", 
  "advanced",
  "expert",
]);

// Update projectFrameworks table
export const projectFrameworks = pgTable("project_frameworks", {
  // ... existing fields ...
  proficiency: proficiencyEnum("proficiency"), // Change from text to enum
  // ... rest of fields ...
});
```

**Migration Required**:
```sql
-- Create enum type
CREATE TYPE proficiency_level AS ENUM ('beginner', 'intermediate', 'advanced', 'expert');

-- Update column
ALTER TABLE project_frameworks 
ALTER COLUMN proficiency TYPE proficiency_level 
USING proficiency::proficiency_level;
```

## 📈 **Enhanced Features**

### **5. Add Project Statistics API**

**Recommended New Endpoint**:
```typescript
// src/app/api/admin/projects/stats/route.ts
export const GET = requireAuth(async (request, {}, user) => {
  const stats = await projectRepository.getStats(user.userId);
  return handleServiceResult(stats, "Project statistics fetched successfully");
});

// Repository method
async getStats(userId?: string): Promise<{
  total: number;
  featured: number;
  totalViews: number;
  totalLikes: number;
}> {
  const conditions = [];
  if (userId) {
    conditions.push(eq(projects.userId, userId));
  }

  const [total, featured, totalViews, totalLikes] = await Promise.all([
    this.getTotalCount(userId),
    this.getFeaturedCount(userId),
    this.getTotalViews(userId),
    this.getTotalLikes(userId),
  ]);

  return { total, featured, totalViews, totalLikes };
}
```

### **6. Add Bulk Order Update API**

**Current Issue**: Mobile supports bulk order updates, but backend only supports individual updates.

**Recommended Changes**:

```typescript
// src/app/api/admin/projects/order/route.ts
export const PUT = requireAuth(async (request, {}, user) => {
  const body = await request.json();
  const orders: ProjectOrderUpdate[] = body;
  
  await projectRepository.updateOrder(orders, user.userId);
  return handleServiceResult({ success: true }, "Project order updated successfully");
});
```

### **7. Add Framework Statistics API**

**Recommended New Endpoints**:
```typescript
// src/app/api/admin/frameworks/stats/route.ts
export const GET = requireAuth(async (request, {}, user) => {
  const stats = await frameworkRepository.getStats(user.userId);
  return handleServiceResult(stats, "Framework statistics fetched successfully");
});
```

## 🔧 **Technical Improvements**

### **8. Add Request Validation Middleware**

**Recommended Implementation**:
```typescript
// src/lib/validation/project-validation.ts
import { z } from 'zod';

export const createProjectSchema = z.object({
  title: z.string().min(1).max(255),
  slug: z.string().min(1).max(255).regex(/^[a-z0-9-]+$/),
  description: z.string().min(1),
  longDescription: z.string().optional(),
  content: z.string().optional(),
  videoUrl: z.string().url().optional(),
  image: z.string().min(1),
  githubUrl: z.string().url().optional(),
  liveUrl: z.string().url().optional(),
  featured: z.boolean().default(false),
  order: z.number().int().min(0).default(0),
  visible: z.boolean().default(true),
  public: z.boolean().default(true),
  frameworkIds: z.array(z.string().uuid()).optional(),
});

// Usage in route handlers
export const POST = requireAuth(
  withValidation(createProjectSchema, async (request, {}, user, validatedData) => {
    const result = await projectService.createProject(validatedData, user.userId);
    return handleServiceResult(result, "Project created successfully", 201);
  })
);
```

### **9. Add Database Indexes for Performance**

**Recommended Indexes**:
```sql
-- Projects table indexes
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_slug ON projects(slug);
CREATE INDEX idx_projects_featured ON projects(featured);
CREATE INDEX idx_projects_visible ON projects(visible);
CREATE INDEX idx_projects_public ON projects(public);
CREATE INDEX idx_projects_order ON projects("order");

-- Frameworks table indexes
CREATE INDEX idx_frameworks_user_id ON frameworks(user_id);
CREATE INDEX idx_frameworks_slug ON frameworks(slug);
CREATE INDEX idx_frameworks_type ON frameworks(type);
CREATE INDEX idx_frameworks_visible ON frameworks(visible);
CREATE INDEX idx_frameworks_order ON frameworks("order");

-- Project frameworks junction table indexes
CREATE INDEX idx_project_frameworks_project_id ON project_frameworks(project_id);
CREATE INDEX idx_project_frameworks_framework_id ON project_frameworks(framework_id);
```

### **10. Add API Rate Limiting**

**Recommended Implementation**:
```typescript
// src/lib/rate-limiting.ts
import rateLimit from 'express-rate-limit';

export const createRateLimit = (windowMs: number, max: number) =>
  rateLimit({
    windowMs,
    max,
    message: 'Too many requests from this IP, please try again later.',
    standardHeaders: true,
    legacyHeaders: false,
  });

// Usage
export const POST = requireAuth(
  withRateLimit(createRateLimit(15 * 60 * 1000, 100)), // 100 requests per 15 minutes
  async (request, {}, user) => {
    // Route handler
  }
);
```

## 📋 **Implementation Priority**

### **High Priority** (Required for Mobile Alignment)
1. ✅ Add `view_count` and `like_count` to projects schema
2. ✅ Add slug support for projects
3. ✅ Add `public` field to frameworks schema
4. ✅ Standardize proficiency levels

### **Medium Priority** (Enhanced Functionality)
5. Add project statistics API
6. Add bulk order update API
7. Add framework statistics API
8. Add request validation middleware

### **Low Priority** (Performance & Security)
9. Add database indexes
10. Add API rate limiting

## 🚀 **Implementation Steps**

1. **Database Migrations**: Create and run migrations for schema changes
2. **Repository Updates**: Add new methods to repositories
3. **Service Updates**: Add new service methods
4. **API Endpoints**: Create new route handlers
5. **Testing**: Add comprehensive tests for new functionality
6. **Documentation**: Update API documentation
7. **Mobile Sync**: Update mobile implementation to use new endpoints

## 📝 **Notes**

- All changes should maintain backward compatibility where possible
- Database migrations should be reversible
- New API endpoints should follow existing patterns and conventions
- Consider adding proper error handling and logging for new functionality
- Update TypeScript types accordingly for all changes

---

**Last Updated**: December 2024  
**Review Status**: Pending Implementation
