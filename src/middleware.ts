import { NextRequest, NextResponse } from 'next/server';
import { env } from '@/lib/config/env';

// Define which routes require API key authentication
const PROTECTED_ROUTES = [
  '/api/mobile/',
  '/api/admin/',
];

function isProtectedRoute(pathname: string): boolean {
  return PROTECTED_ROUTES.some(route => pathname.startsWith(route));
}

function validateApiKey(request: NextRequest): boolean {
  const apiKey = request.headers.get('X-API-Key');
  
  if (!apiKey) {
    return false;
  }

  // Get API keys from environment configuration
  const allowedApiKeys = env.MOBILE_API_KEYS.split(',').map(key => key.trim());
  return allowedApiKeys.includes(apiKey);
}

export function middleware(request: NextRequest) {
  // Handle CORS and API key validation for API routes
  if (request.nextUrl.pathname.startsWith('/api/')) {
    const origin = request.headers.get('origin');
    
    // Set CORS headers
    const response = NextResponse.next();
    
    // Allow common origins for development
    if (origin && (
      origin.includes('localhost') || 
      origin.includes('127.0.0.1') ||
      origin.startsWith('capacitor://') || 
      origin.startsWith('ionic://')
    )) {
      response.headers.set('Access-Control-Allow-Origin', origin);
    }

    response.headers.set('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS, PATCH');
    response.headers.set('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-API-Key, X-Requested-With');
    response.headers.set('Access-Control-Allow-Credentials', 'true');
    response.headers.set('Access-Control-Max-Age', '86400');

    // Handle preflight requests
    if (request.method === 'OPTIONS') {
      return new NextResponse(null, { status: 200, headers: response.headers });
    }

    // API Key Authentication for protected routes
    if (isProtectedRoute(request.nextUrl.pathname)) {
      if (!validateApiKey(request)) {
        return NextResponse.json(
          {
            success: false,
            error: 'API key required for this endpoint',
          },
          { 
            status: 401,
            headers: response.headers
          }
        );
      }
    }

    return response;
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/api/:path*',
  ],
};
