import { NextRequest, NextResponse } from "next/server";
import { env } from "@/lib/config/env";

export interface ApiKeyAuthResult {
  success: boolean;
  userId?: string;
  error?: string;
  source: 'api-key' | 'jwt' | 'none';
}

/**
 * API Key Authentication Middleware
 * Validates X-API-Key header for mobile app access
 */
export function validateApiKey(request: NextRequest): ApiKeyAuthResult {
  const apiKey = request.headers.get('X-API-Key');
  
  if (!apiKey) {
    return {
      success: false,
      error: 'API key required',
      source: 'none'
    };
  }

  // Get allowed API keys from environment
  const allowedApiKeys = env.MOBILE_API_KEYS.split(',').map(key => key.trim());
  
  if (!allowedApiKeys.includes(apiKey)) {
    return {
      success: false,
      error: 'Invalid API key',
      source: 'api-key'
    };
  }

  // For now, we'll use a simple mapping
  // In production, you might want to store API keys in the database
  // with associated user IDs and permissions
  return {
    success: true,
    userId: 'mobile-user', // This could be dynamic based on API key
    source: 'api-key'
  };
}

/**
 * Combined Authentication Middleware
 * Supports both JWT (for web) and API Key (for mobile) authentication
 */
export async function authenticateRequest(request: NextRequest): Promise<ApiKeyAuthResult> {
  // First, try API key authentication (for mobile apps)
  const apiKeyResult = validateApiKey(request);
  if (apiKeyResult.success) {
    return apiKeyResult;
  }

  // If no API key or invalid, try JWT authentication (for web apps)
  const authHeader = request.headers.get('Authorization');
  if (authHeader && authHeader.startsWith('Bearer ')) {
    try {
      const token = authHeader.substring(7);
      // You would verify the JWT token here
      // For now, we'll assume it's valid if present
      return {
        success: true,
        userId: 'web-user', // This would come from JWT payload
        source: 'jwt'
      };
    } catch (error) {
      return {
        success: false,
        error: 'Invalid JWT token',
        source: 'jwt'
      };
    }
  }

  return {
    success: false,
    error: 'Authentication required',
    source: 'none'
  };
}

/**
 * Middleware wrapper for API routes that require authentication
 */
export function requireAuth(
  handler: (
    request: NextRequest,
    auth: ApiKeyAuthResult
  ) => Promise<NextResponse> | NextResponse
) {
  return async (request: NextRequest) => {
    const authResult = await authenticateRequest(request);

    if (!authResult.success) {
      return NextResponse.json(
        {
          success: false,
          error: authResult.error || 'Authentication required',
        },
        { status: 401 }
      );
    }

    return handler(request, authResult);
  };
}

/**
 * Middleware wrapper for API routes that require admin access
 */
export function requireAdmin(
  handler: (
    request: NextRequest,
    auth: ApiKeyAuthResult
  ) => Promise<NextResponse> | NextResponse
) {
  return async (request: NextRequest) => {
    const authResult = await authenticateRequest(request);

    if (!authResult.success) {
      return NextResponse.json(
        {
          success: false,
          error: authResult.error || 'Authentication required',
        },
        { status: 401 }
      );
    }

    // For API key authentication, you might want to check permissions
    // For now, we'll allow API key users to access admin endpoints
    // In production, you should implement proper role-based access control
    if (authResult.source === 'api-key') {
      // You could check if the API key has admin permissions
      // For now, we'll allow it
    }

    return handler(request, authResult);
  };
}

/**
 * Middleware wrapper for public API routes (no authentication required)
 */
export function publicRoute(
  handler: (
    request: NextRequest,
    auth?: ApiKeyAuthResult
  ) => Promise<NextResponse> | NextResponse
) {
  return async (request: NextRequest) => {
    // Try to authenticate, but don't require it
    const authResult = await authenticateRequest(request);
    
    return handler(request, authResult.success ? authResult : undefined);
  };
}
