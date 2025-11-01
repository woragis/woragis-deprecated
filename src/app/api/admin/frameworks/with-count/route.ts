import { NextRequest } from "next/server";
import { frameworkService } from "@/server/services";
import {
  authMiddleware,
  handleServiceResult,
  withErrorHandling,
  handleAuthError,
} from "@/utils/response-helpers";

// GET /api/admin/frameworks/with-count - Get frameworks with project count
export const GET = withErrorHandling(async (request: NextRequest) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  const { searchParams } = new URL(request.url);
  const limit = parseInt(searchParams.get("limit") || "10");

  const result = await frameworkService.getPopularFrameworks(limit);

  // Transform to match mobile expectations - frameworks as array
  if (result.success && result.data) {
    const frameworks = result.data.map((item) => ({
      ...item.framework,
      projectCount: item.projectCount,
    }));

    return handleServiceResult(
      { success: true, data: frameworks },
      "Frameworks with project count fetched successfully"
    );
  }

  return handleServiceResult(
    result,
    "Frameworks with project count fetched successfully"
  );
});

