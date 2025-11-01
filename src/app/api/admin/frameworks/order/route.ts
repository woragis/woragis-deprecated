import { NextRequest } from "next/server";
import { frameworkService } from "@/server/services";
import {
  handleServiceResult,
  withErrorHandling,
  authMiddleware,
  handleAuthError,
  badRequestResponse,
} from "@/utils/response-helpers";

// PUT /api/admin/frameworks/order - Update framework order
export const PUT = withErrorHandling(async (request: NextRequest) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  const body = await request.json();
  const { frameworkOrders } = body;

  if (!frameworkOrders || !Array.isArray(frameworkOrders)) {
    return badRequestResponse("frameworkOrders array is required");
  }

  const result = await frameworkService.updateFrameworkOrder(frameworkOrders);
  return handleServiceResult(result, "Framework order updated successfully");
});

