import { NextRequest } from "next/server";
import { frameworkService } from "@/server/services";
import {
  handleServiceResult,
  withErrorHandling,
  notFoundResponse,
} from "@/utils/response-helpers";

// GET /api/admin/frameworks/[id]/project-count - Get framework project count
export const GET = withErrorHandling(
  async (
    request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
  ) => {
    const { id } = await params;
    const result = await frameworkService.getFrameworkWithProjectCount(id);

    if (!result.success || !result.data) {
      return notFoundResponse(result.error || "Framework not found");
    }

    // Mobile expects { count: number }
    return handleServiceResult(
      {
        success: true,
        data: { count: result.data.projectCount },
      },
      "Project count fetched successfully"
    );
  }
);

