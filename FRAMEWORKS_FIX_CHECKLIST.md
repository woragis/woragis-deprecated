# Frameworks Backend-Mobile Alignment - Fix Checklist

## 🔴 Critical Issues (Fix Immediately)

### Issue 1: Response Structure Mismatch
**Problem**: Mobile expects nested `data.frameworks`, backend returns array directly in `data`

**Files to Change**:
- [ ] `mobile/lib/features/frameworks/data/datasources/frameworks_remote_datasource.dart`

**Changes**:
```dart
// Line 136-138: Change from:
final frameworks = (data['data']['frameworks'] as List)
    .map((frameworkJson) => FrameworkModel.fromJson(frameworkJson).toEntity())
    .toList();

// To:
final frameworks = (data['data'] as List)
    .map((frameworkJson) => FrameworkModel.fromJson(frameworkJson).toEntity())
    .toList();
```

**Also check**:
- [ ] Line 429 (getFrameworksWithProjectCount) - same issue

---

### Issue 2: Field Name Mismatch - `website` vs `url`
**Problem**: Backend uses `website`, mobile uses `url`

**Decision Required**: Choose one naming convention

**Option A: Use `url` (Mobile-friendly)**
- [ ] Update backend schema: `src/server/db/schemas/frameworks.ts` line 42
- [ ] Create migration to rename column
- [ ] Update TypeScript types

**Option B: Use `website` (Backend current)**  ✅ Recommended
- [ ] Update mobile entity: `mobile/lib/features/frameworks/domain/entities/framework_entity.dart`
- [ ] Update mobile model: `mobile/lib/features/frameworks/data/models/framework_model.dart`
- [ ] Update mobile datasource: `mobile/lib/features/frameworks/data/datasources/frameworks_remote_datasource.dart`
  - Line 30, 247, 292: Change `url` → `website`

---

### Issue 3: Missing API Route - Get Framework by Slug
**Problem**: Mobile calls `/admin/frameworks/slug/:slug` but route doesn't exist

**Files to Create**:
- [ ] `src/app/api/admin/frameworks/slug/[slug]/route.ts`

**Implementation**:
```typescript
import { NextRequest } from "next/server";
import { frameworkService } from "@/server/services";
import { handleServiceResult, withErrorHandling, notFoundResponse } from "@/utils/response-helpers";

export const GET = withErrorHandling(
  async (request: NextRequest, { params }: { params: Promise<{ slug: string }> }) => {
    const { slug } = await params;
    const result = await frameworkService.getFrameworkBySlug(slug);

    if (!result.success || !result.data) {
      return notFoundResponse(result.error || "Framework not found");
    }

    return handleServiceResult(result, "Framework fetched successfully");
  }
);
```

---

### Issue 4: Pagination Mismatch - `page` vs `offset`
**Problem**: Mobile sends `page`, backend expects `offset`

**Files to Change**:
- [ ] `src/app/api/admin/frameworks/route.ts`

**Changes**:
```typescript
// Lines 18-34: Update parameter handling
const { searchParams } = new URL(request.url);

// Add page support
const page = searchParams.get("page");
const limit = parseInt(searchParams.get("limit") || "20");
const offset = page 
  ? (parseInt(page) - 1) * limit 
  : parseInt(searchParams.get("offset") || "0");

const filters: FrameworkFilters = {
  visible: searchParams.get("visible") === "true" ? true 
    : searchParams.get("visible") === "false" ? false 
    : undefined,
  search: searchParams.get("search") || undefined,
  type: (searchParams.get("type") as FrameworkType) || undefined,
  limit,
  offset,
};
```

---

## 🟡 Medium Priority Issues

### Issue 5: Missing Query Parameters in Backend
**Problem**: Backend doesn't support `public`, `sortBy`, `sortOrder` parameters

**Files to Change**:
- [ ] `src/types/framework.ts` (or wherever FrameworkFilters is defined)
- [ ] `src/app/api/admin/frameworks/route.ts`
- [ ] `src/server/repositories/framework.repository.ts`

**Type Updates**:
```typescript
// Add to FrameworkFilters type
export interface FrameworkFilters {
  visible?: boolean;
  public?: boolean;      // ADD THIS
  search?: string;
  type?: FrameworkType;
  sortBy?: string;       // ADD THIS
  sortOrder?: 'asc' | 'desc';  // ADD THIS
  limit?: number;
  offset?: number;
}
```

**Route Updates**:
```typescript
const filters: FrameworkFilters = {
  visible: /* ... */,
  public: searchParams.get("public") === "true" ? true 
    : searchParams.get("public") === "false" ? false 
    : undefined,
  search: /* ... */,
  type: /* ... */,
  sortBy: searchParams.get("sortBy") || undefined,
  sortOrder: (searchParams.get("sortOrder") as 'asc' | 'desc') || undefined,
  limit: /* ... */,
  offset: /* ... */,
};
```

**Repository Updates**:
- [ ] Update `search()` method to handle new parameters
- [ ] Add sorting logic based on `sortBy` and `sortOrder`
- [ ] Add filtering for `public` field

---

### Issue 6: Missing `version` Field in Mobile
**Problem**: Backend has `version` field, mobile doesn't

**Files to Change**:
- [ ] `mobile/lib/features/frameworks/domain/entities/framework_entity.dart`
- [ ] `mobile/lib/features/frameworks/data/models/framework_model.dart`
- [ ] All use cases and BLoC that create/update frameworks

**Entity Update**:
```dart
class FrameworkEntity extends Equatable {
  // ... existing fields ...
  final String? version;  // ADD THIS
  
  const FrameworkEntity({
    // ... existing params ...
    this.version,  // ADD THIS
  });
  
  @override
  List<Object?> get props => [
    // ... existing props ...
    version,  // ADD THIS
  ];
  
  FrameworkEntity copyWith({
    // ... existing params ...
    String? version,  // ADD THIS
  }) {
    return FrameworkEntity(
      // ... existing fields ...
      version: version ?? this.version,  // ADD THIS
    );
  }
}
```

**Model Update**: Similar changes to `FrameworkModel`

---

### Issue 7A: Missing API Route - Update Framework Order
**Problem**: Mobile calls `/admin/frameworks/order` but route doesn't exist

**Files to Create**:
- [ ] `src/app/api/admin/frameworks/order/route.ts`

**Implementation**:
```typescript
import { NextRequest } from "next/server";
import { frameworkService } from "@/server/services";
import { handleServiceResult, withErrorHandling, authMiddleware, handleAuthError } from "@/utils/response-helpers";

export const PUT = withErrorHandling(async (request: NextRequest) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  const body = await request.json();
  const { frameworkOrders } = body;

  const result = await frameworkService.updateFrameworkOrder(frameworkOrders);
  return handleServiceResult(result, "Framework order updated successfully");
});
```

---

### Issue 7B: Missing API Route - Get Framework Project Count
**Problem**: Mobile calls `/admin/frameworks/:id/project-count` but route doesn't exist

**Files to Create**:
- [ ] `src/app/api/admin/frameworks/[id]/project-count/route.ts`

**Implementation**:
```typescript
import { NextRequest } from "next/server";
import { frameworkService } from "@/server/services";
import { handleServiceResult, withErrorHandling, notFoundResponse } from "@/utils/response-helpers";

export const GET = withErrorHandling(
  async (request: NextRequest, { params }: { params: Promise<{ id: string }> }) => {
    const { id } = await params;
    const result = await frameworkService.getFrameworkWithProjectCount(id);

    if (!result.success || !result.data) {
      return notFoundResponse(result.error || "Framework not found");
    }

    // Mobile expects { count: number }
    return handleServiceResult(
      { 
        success: true, 
        data: { count: result.data.projectCount } 
      }, 
      "Project count fetched successfully"
    );
  }
);
```

---

### Issue 7C: Missing API Route - Get Frameworks with Project Count
**Problem**: Mobile calls `/admin/frameworks/with-count` but route doesn't exist

**Files to Create**:
- [ ] `src/app/api/admin/frameworks/with-count/route.ts`

**Implementation**:
```typescript
import { NextRequest } from "next/server";
import { frameworkService } from "@/server/services";
import { authMiddleware, handleServiceResult, withErrorHandling, handleAuthError } from "@/utils/response-helpers";

export const GET = withErrorHandling(async (request: NextRequest) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  const { searchParams } = new URL(request.url);
  const limit = parseInt(searchParams.get("limit") || "10");

  const result = await frameworkService.getPopularFrameworks(limit);
  
  // Transform to match mobile expectations
  if (result.success && result.data) {
    const frameworks = result.data.map(item => ({
      ...item.framework,
      projectCount: item.projectCount
    }));
    
    return handleServiceResult(
      { success: true, data: { frameworks } },
      "Frameworks with project count fetched successfully"
    );
  }
  
  return handleServiceResult(result, "Frameworks with project count fetched successfully");
});
```

---

## 🟢 Low Priority / Future Enhancements

### Enhancement 1: Expose Analytics to Mobile
- [ ] Create routes for version distribution
- [ ] Create routes for proficiency distribution
- [ ] Update mobile to consume these endpoints

### Enhancement 2: Add Tests
**Backend**:
- [ ] Unit tests for `FrameworkService`
- [ ] Integration tests for API routes
- [ ] Repository tests with test database

**Mobile**:
- [ ] Unit tests for use cases
- [ ] BLoC tests
- [ ] Widget tests for framework screens
- [ ] Integration tests with mock API

### Enhancement 3: Documentation
- [ ] OpenAPI/Swagger specification
- [ ] API documentation
- [ ] Mobile SDK documentation
- [ ] Architecture decision records (ADRs)

---

## Verification Checklist

After making changes, verify:

### Backend
- [ ] Run TypeScript compiler: `npm run type-check` or `tsc --noEmit`
- [ ] Run linter: `npm run lint`
- [ ] Test all new routes with Postman/curl
- [ ] Verify database migrations if schema changed
- [ ] Check all routes return correct response structure

### Mobile
- [ ] Run Dart analyzer: `flutter analyze`
- [ ] Run tests: `flutter test`
- [ ] Test API calls with running backend
- [ ] Verify offline caching still works
- [ ] Check BLoC state transitions
- [ ] Test pagination behavior
- [ ] Verify sorting and filtering

### Integration
- [ ] Test full CRUD flow from mobile → backend → database
- [ ] Verify authentication on all routes
- [ ] Test error scenarios (network errors, validation errors)
- [ ] Verify response structure matches on all endpoints
- [ ] Test pagination with various page sizes
- [ ] Verify sorting and filtering work correctly

---

## Quick Reference: Files Changed

### Backend (7-9 files)
1. `src/app/api/admin/frameworks/route.ts` - Add pagination, sorting, filtering
2. `src/app/api/admin/frameworks/slug/[slug]/route.ts` - NEW
3. `src/app/api/admin/frameworks/order/route.ts` - NEW
4. `src/app/api/admin/frameworks/[id]/project-count/route.ts` - NEW
5. `src/app/api/admin/frameworks/with-count/route.ts` - NEW
6. `src/server/repositories/framework.repository.ts` - Add sorting/filtering
7. `src/types/framework.ts` (or similar) - Update FrameworkFilters type
8. `src/server/db/schemas/frameworks.ts` - OPTIONAL: rename website→url or vice versa
9. Migration file - OPTIONAL: if schema changed

### Mobile (2-5 files)
1. `mobile/lib/features/frameworks/data/datasources/frameworks_remote_datasource.dart` - Fix response parsing, rename url→website
2. `mobile/lib/features/frameworks/domain/entities/framework_entity.dart` - OPTIONAL: add version field
3. `mobile/lib/features/frameworks/data/models/framework_model.dart` - OPTIONAL: add version field
4. BLoC/UseCase files - OPTIONAL: if adding version field
5. Any UI files using frameworks - OPTIONAL: if exposing new features

---

## Estimated Effort

| Priority | Tasks | Effort | Time Estimate |
|----------|-------|--------|---------------|
| 🔴 Critical | 4 issues | Medium | 2-3 hours |
| 🟡 Medium | 4 issues | High | 4-6 hours |
| 🟢 Low | 3 enhancements | High | 8+ hours |

**Total for Critical + Medium**: ~6-9 hours of development + testing

---

## Progress Tracking

**Started**: ___________  
**Critical Issues Complete**: ___________  
**Medium Issues Complete**: ___________  
**All Issues Resolved**: ___________  
**Testing Complete**: ___________  
**Production Deployed**: ___________

