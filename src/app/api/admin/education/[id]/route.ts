import { NextRequest } from "next/server";
import { educationService } from "@/server/services/education.service";
import { authMiddleware } from "@/lib/auth";
import {
  handleServiceResult,
  withErrorHandling,
  handleAuthError,
} from "@/utils/response-helpers";
import type { NewEducation } from "@/types/education";

// GET /api/admin/education/[id] - Get education record by ID
export const GET = withErrorHandling(async (
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  const { id } = await params;
  const result = await educationService.getEducationById(id);
  return handleServiceResult(result, "Education record fetched successfully");
});

// PUT /api/admin/education/[id] - Update education record
export const PUT = withErrorHandling(async (
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  const { id } = await params;
  const body = await request.json();
  
  // Convert date strings to Date objects for Drizzle
  const educationData: Partial<NewEducation> = {
    ...body,
    startDate: body.startDate ? new Date(body.startDate) : undefined,
    endDate: body.endDate ? new Date(body.endDate) : undefined,
    completionDate: body.completionDate ? new Date(body.completionDate) : undefined,
  };

  const result = await educationService.updateEducation(id, educationData);
  return handleServiceResult(result, "Education record updated successfully");
});

// DELETE /api/admin/education/[id] - Delete education record
export const DELETE = withErrorHandling(async (
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  const { id } = await params;
  const result = await educationService.deleteEducation(id);
  return handleServiceResult(result, "Education record deleted successfully");
});
