import { NextRequest } from "next/server";
import { educationService } from "@/server/services/education.service";
import { authMiddleware } from "@/lib/auth";
import {
  handleServiceResult,
  withErrorHandling,
  handleAuthError,
} from "@/utils/response-helpers";
import type { NewEducation, EducationFilters, EducationType, DegreeLevel } from "@/types/education";

// GET /api/admin/education - Get all education records with optional filtering
export const GET = withErrorHandling(async (request: NextRequest) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  const { searchParams } = new URL(request.url);
  const filters: EducationFilters = {
    visible:
      searchParams.get("visible") === "true"
        ? true
        : searchParams.get("visible") === "false"
        ? false
        : undefined,
    search: searchParams.get("search") || undefined,
    type: (searchParams.get("type") as EducationType) || undefined,
    degreeLevel: (searchParams.get("degreeLevel") as DegreeLevel) || undefined,
    institution: searchParams.get("institution") || undefined,
    limit: searchParams.get("limit")
      ? parseInt(searchParams.get("limit")!)
      : undefined,
    offset: searchParams.get("offset")
      ? parseInt(searchParams.get("offset")!)
      : undefined,
  };

  const result = await educationService.searchEducation(filters);
  return handleServiceResult(result, "Education records fetched successfully");
});

// POST /api/admin/education - Create new education record
export const POST = withErrorHandling(async (request: NextRequest) => {
  const authResult = await authMiddleware(request);
  if (!authResult.success) {
    return handleAuthError("Unauthorized");
  }

  if (!authResult.userId) {
    return handleAuthError("User ID not found");
  }

  const body = await request.json();
  
  // Convert date strings to Date objects for Drizzle
  const educationData: NewEducation = {
    ...body,
    startDate: body.startDate ? new Date(body.startDate) : undefined,
    endDate: body.endDate ? new Date(body.endDate) : undefined,
    completionDate: body.completionDate ? new Date(body.completionDate) : undefined,
  };

  const result = await educationService.createEducation(
    educationData,
    authResult.userId
  );
  return handleServiceResult(result, "Education record created successfully", 201);
});
