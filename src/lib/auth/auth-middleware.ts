import { NextRequest, NextResponse } from "next/server";
import { cookies } from "next/headers";
import { env } from "@/lib/config/env";
import { authRepository } from "@/server/repositories";
import { authService } from "@/server/services";

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  role: "admin" | "user";
  avatar?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthenticatedUser {
  userId: string;
  id: string;
  email: string;
  name: string;
  role: "admin" | "user";
  avatar?: string;
  createdAt: string;
  updatedAt: string;
}

// Helper function to get a valid user ID for mobile API requests
async function getMobileApiUserId(): Promise<string | null> {
  try {
    // Get the first user (which should be an admin user)
    const firstUser = await authRepository.getFirstUser();
    return firstUser?.id || null;
  } catch (error) {
    console.error("Error getting mobile API user ID:", error);
    return null;
  }
}

export async function getServerAuth(): Promise<AuthUser | null> {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("auth-token")?.value;

    if (!token) {
      return null;
    }

    // Verify token directly using authService to avoid circular dependency
    const tokenData = await authService.verifyToken(token);
    if (!tokenData) {
      return null;
    }

    const user = await authService.getUserById(tokenData.userId);
    if (!user) {
      return null;
    }

    // Map user data to AuthUser interface
    const authUser: AuthUser = {
      id: user.id,
      email: user.email,
      name:
        user.firstName && user.lastName
          ? `${user.firstName} ${user.lastName}`.trim()
          : user.username || user.email,
      role: user.role,
      avatar: user.avatar,
      createdAt: user.createdAt.toISOString(),
      updatedAt: user.updatedAt.toISOString(),
    };

    return authUser;
  } catch (error) {
    console.error("Auth verification failed:", error);
    return null;
  }
}

export function requireAuth(
  handler: (
    request: NextRequest,
    user: AuthenticatedUser
  ) => Promise<NextResponse> | NextResponse
) {
  return async (request: NextRequest) => {
    // Check if this is an API request with X-API-Key header (mobile app)
    const apiKey = request.headers.get('X-API-Key');
    if (apiKey) {
      // Validate API key
      const allowedApiKeys = env.MOBILE_API_KEYS.split(',').map(key => key.trim());
      if (allowedApiKeys.includes(apiKey)) {
        // For mobile API requests, we can either:
        // 1. Use the JWT token if present (preferred for user-specific data)
        // 2. Fall back to a system user if no JWT token
        
        const authHeader = request.headers.get('Authorization');
        if (authHeader && authHeader.startsWith('Bearer ')) {
          // Mobile app is sending both API key and JWT token
          // Try to validate the JWT token first using authService directly
          try {
            const token = authHeader.substring(7);
            const tokenData = await authService.verifyToken(token);

            if (tokenData) {
              const user = await authService.getUserById(tokenData.userId);
              if (user) {
                // Use the actual authenticated user
                const authenticatedUser: AuthenticatedUser = {
                  userId: user.id,
                  id: user.id,
                  email: user.email,
                  name: user.firstName && user.lastName
                    ? `${user.firstName} ${user.lastName}`.trim()
                    : user.username || user.email,
                  role: user.role,
                  avatar: user.avatar,
                  createdAt: user.createdAt.toISOString(),
                  updatedAt: user.updatedAt.toISOString(),
                };
                return handler(request, authenticatedUser);
              }
            }
          } catch (error) {
            console.error("JWT validation failed for mobile API:", error);
            // Fall through to use system user
          }
        }
        
        // Fall back to system user for mobile API requests
        const mobileUserId = await getMobileApiUserId();
        if (!mobileUserId) {
          return NextResponse.json(
            {
              success: false,
              error: 'No valid user found for mobile API',
            },
            { status: 500 }
          );
        }
        
        // Create an authenticated user for API key requests using real user data
        const authenticatedUser: AuthenticatedUser = {
          userId: mobileUserId,
          id: mobileUserId,
          email: 'mobile@api.com',
          name: 'Mobile API User',
          role: 'admin', // Allow admin access for mobile API
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };
        return handler(request, authenticatedUser);
      } else {
        return NextResponse.json(
          {
            success: false,
            error: 'Invalid API key',
          },
          { status: 401 }
        );
      }
    }

    // Fall back to cookie-based authentication for web requests
    const authUser = await getServerAuth();

    if (!authUser) {
      // Check if this is an API request (not a page request)
      if (request.nextUrl.pathname.startsWith('/api/')) {
        return NextResponse.json(
          {
            success: false,
            error: 'Authentication required',
          },
          { status: 401 }
        );
      }
      // For page requests, redirect to login
      return NextResponse.redirect(new URL("/login", request.url));
    }

    // Convert AuthUser to AuthenticatedUser format
    const authenticatedUser: AuthenticatedUser = {
      ...authUser,
      userId: authUser.id,
    };

    return handler(request, authenticatedUser);
  };
}

export function requireAdmin(
  handler: (
    request: NextRequest,
    user: AuthenticatedUser
  ) => Promise<NextResponse> | NextResponse
) {
  return async (request: NextRequest) => {
    // Check if this is an API request with X-API-Key header (mobile app)
    const apiKey = request.headers.get('X-API-Key');
    if (apiKey) {
      // Validate API key
      const allowedApiKeys = env.MOBILE_API_KEYS.split(',').map(key => key.trim());
      if (allowedApiKeys.includes(apiKey)) {
        // For mobile API requests, we can either:
        // 1. Use the JWT token if present (preferred for user-specific data)
        // 2. Fall back to a system user if no JWT token
        
        const authHeader = request.headers.get('Authorization');
        if (authHeader && authHeader.startsWith('Bearer ')) {
          // Mobile app is sending both API key and JWT token
          // Try to validate the JWT token first using authService directly
          try {
            const token = authHeader.substring(7);
            const tokenData = await authService.verifyToken(token);

            if (tokenData) {
              const user = await authService.getUserById(tokenData.userId);
              if (user) {
                // Check if user has admin role
                if (user.role !== 'admin' && user.role !== 'super_admin') {
                  return NextResponse.json(
                    {
                      success: false,
                      error: 'Admin access required',
                    },
                    { status: 403 }
                  );
                }
                // Use the actual authenticated user
                const authenticatedUser: AuthenticatedUser = {
                  userId: user.id,
                  id: user.id,
                  email: user.email,
                  name: user.firstName && user.lastName
                    ? `${user.firstName} ${user.lastName}`.trim()
                    : user.username || user.email,
                  role: user.role,
                  avatar: user.avatar,
                  createdAt: user.createdAt.toISOString(),
                  updatedAt: user.updatedAt.toISOString(),
                };
                return handler(request, authenticatedUser);
              }
            }
          } catch (error) {
            console.error("JWT validation failed for mobile API:", error);
            // Fall through to use system user
          }
        }
        
        // Fall back to system user for mobile API requests (with admin role)
        const mobileUserId = await getMobileApiUserId();
        if (!mobileUserId) {
          return NextResponse.json(
            {
              success: false,
              error: 'No valid user found for mobile API',
            },
            { status: 500 }
          );
        }
        
        // Create an authenticated user for API key requests with admin role using real user data
        const authenticatedUser: AuthenticatedUser = {
          userId: mobileUserId,
          id: mobileUserId,
          email: 'mobile@api.com',
          name: 'Mobile API User',
          role: 'admin', // Allow admin access for mobile API
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };
        return handler(request, authenticatedUser);
      } else {
        return NextResponse.json(
          {
            success: false,
            error: 'Invalid API key',
          },
          { status: 401 }
        );
      }
    }

    // Fall back to cookie-based authentication for web requests
    const authUser = await getServerAuth();

    if (!authUser) {
      // Check if this is an API request (not a page request)
      if (request.nextUrl.pathname.startsWith('/api/')) {
        return NextResponse.json(
          {
            success: false,
            error: 'Authentication required',
          },
          { status: 401 }
        );
      }
      // For page requests, redirect to login
      return NextResponse.redirect(new URL("/login", request.url));
    }

    if (authUser.role !== "admin") {
      // Check if this is an API request (not a page request)
      if (request.nextUrl.pathname.startsWith('/api/')) {
        return NextResponse.json(
          {
            success: false,
            error: 'Admin access required',
          },
          { status: 403 }
        );
      }
      // For page requests, redirect to home
      return NextResponse.redirect(new URL("/", request.url));
    }

    // Convert AuthUser to AuthenticatedUser format
    const authenticatedUser: AuthenticatedUser = {
      ...authUser,
      userId: authUser.id,
    };

    return handler(request, authenticatedUser);
  };
}

// Auth middleware that returns a result object instead of redirecting
export async function authMiddleware(request: NextRequest): Promise<{
  success: boolean;
  userId?: string;
  error?: string;
}> {
  try {
    // Check if this is an API request with X-API-Key header (mobile app)
    const apiKey = request.headers.get('X-API-Key');
    if (apiKey) {
      // Validate API key
      const allowedApiKeys = env.MOBILE_API_KEYS.split(',').map(key => key.trim());
      if (allowedApiKeys.includes(apiKey)) {
        // For mobile API requests, we can either:
        // 1. Use the JWT token if present (preferred for user-specific data)
        // 2. Fall back to a system user if no JWT token
        
        const authHeader = request.headers.get('Authorization');
        if (authHeader && authHeader.startsWith('Bearer ')) {
          // Mobile app is sending both API key and JWT token
          // Try to validate the JWT token first using authService directly
          try {
            const token = authHeader.substring(7);
            const tokenData = await authService.verifyToken(token);

            if (tokenData) {
              const user = await authService.getUserById(tokenData.userId);
              if (user) {
                // Use the actual authenticated user
                return {
                  success: true,
                  userId: user.id,
                };
              }
            }
          } catch (error) {
            console.error("JWT validation failed for mobile API:", error);
            // Fall through to use system user
          }
        }
        
        // Fall back to system user for mobile API requests
        const mobileUserId = await getMobileApiUserId();
        if (!mobileUserId) {
          return {
            success: false,
            error: 'No valid user found for mobile API',
          };
        }
        
        return {
          success: true,
          userId: mobileUserId,
        };
      } else {
        return {
          success: false,
          error: 'Invalid API key',
        };
      }
    }

    // Fall back to cookie-based authentication for web requests
    const authUser = await getServerAuth();

    if (!authUser) {
      return {
        success: false,
        error: "Authentication required",
      };
    }

    return {
      success: true,
      userId: authUser.id,
    };
  } catch (error) {
    console.error("Auth middleware error:", error);
    return {
      success: false,
      error: "Authentication failed",
    };
  }
}
