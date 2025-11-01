import { NextRequest } from "next/server";
import { projectService } from "@/server/services";
import { authMiddleware } from "@/lib/auth";
import {
  handleServiceResult,
  withErrorHandling,
  handleAuthError,
  notFoundResponse,
} from "@/utils/response-helpers";

// GET /api/admin/projects/slug/[slug] - Get project by slug
export const GET = withErrorHandling(
  async (
    request: NextRequest,
    { params }: { params: Promise<{ slug: string }> }
  ) => {
    const authResult = await authMiddleware(request);
    if (!authResult.success) {
      return handleAuthError(authResult.error);
    }

    const { slug } = await params;
    const result = await projectService.getProjectBySlug(slug, authResult.userId!);

    if (!result.success) {
      return notFoundResponse(result.error || "Project not found");
    }

    return handleServiceResult(result, "Project fetched successfully");
  }
);
