import { NextRequest } from "next/server";
import { frameworkService } from "@/server/services";
import {
  handleServiceResult,
  withErrorHandling,
  notFoundResponse,
} from "@/utils/response-helpers";

// GET /api/admin/frameworks/slug/[slug] - Get framework by slug
export const GET = withErrorHandling(
  async (
    request: NextRequest,
    { params }: { params: Promise<{ slug: string }> }
  ) => {
    const { slug } = await params;
    const result = await frameworkService.getFrameworkBySlug(slug);

    if (!result.success || !result.data) {
      return notFoundResponse(result.error || "Framework not found");
    }

    return handleServiceResult(result, "Framework fetched successfully");
  }
);

