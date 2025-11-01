import { NextRequest } from "next/server";
import { authService } from "../../../../server/services";
import {
  successResponse,
  unauthorizedResponse,
  notFoundResponse,
  withErrorHandling,
} from "@/utils/response-helpers";

export const GET = withErrorHandling(async (request: NextRequest) => {
  // Check for API key first (mobile app)
  const apiKey = request.headers.get('X-API-Key');
  if (apiKey) {
    // For mobile API requests with API key, we can return a system user
    // or validate JWT token if present
    const authHeader = request.headers.get("authorization");
    if (authHeader && authHeader.startsWith("Bearer ")) {
      const token = authHeader.substring(7);
      const tokenData = await authService.verifyToken(token);

      if (tokenData) {
        const user = await authService.getUserById(tokenData.userId);
        if (user) {
          const userData = {
            id: user.id,
            email: user.email,
            username: user.username,
            firstName: user.firstName,
            lastName: user.lastName,
            avatar: user.avatar,
            role: user.role,
            isActive: user.isActive,
            emailVerified: user.emailVerified,
            lastLoginAt: user.lastLoginAt,
            createdAt: user.createdAt,
            updatedAt: user.updatedAt,
          };
          return successResponse(userData, "User profile fetched successfully");
        }
      }
    }
    
    // If no valid JWT token, return 401 for mobile API requests too
    return unauthorizedResponse("Valid JWT token required for user data");
  }

  // For web requests, require JWT token
  const authHeader = request.headers.get("authorization");
  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    return unauthorizedResponse("No authorization token provided");
  }

  const token = authHeader.substring(7);
  const tokenData = await authService.verifyToken(token);

  if (!tokenData) {
    return unauthorizedResponse("Invalid or expired token");
  }

  const user = await authService.getUserById(tokenData.userId);
  if (!user) {
    return notFoundResponse("User not found");
  }

  const userData = {
    id: user.id,
    email: user.email,
    username: user.username,
    firstName: user.firstName,
    lastName: user.lastName,
    avatar: user.avatar,
    role: user.role,
    isActive: user.isActive,
    emailVerified: user.emailVerified,
    lastLoginAt: user.lastLoginAt,
    createdAt: user.createdAt,
    updatedAt: user.updatedAt,
  };

  return successResponse(userData, "User profile fetched successfully");
});
