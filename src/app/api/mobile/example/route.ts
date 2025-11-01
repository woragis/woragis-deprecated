import { NextRequest, NextResponse } from "next/server";

function validateApiKey(request: NextRequest): boolean {
  const apiKey = request.headers.get('X-API-Key');
  
  if (!apiKey) {
    return false;
  }

  // Valid API keys
  const allowedApiKeys = ['mobile-key-1', 'mobile-key-2', 'mobile-key-3'];
  return allowedApiKeys.includes(apiKey);
}

/**
 * Example protected API route for mobile apps
 * Requires X-API-Key header
 */
export async function GET(request: NextRequest) {
  // Validate API key
  if (!validateApiKey(request)) {
    return NextResponse.json(
      {
        success: false,
        error: 'API key required for this endpoint',
      },
      { status: 401 }
    );
  }

  return NextResponse.json({
    success: true,
    data: {
      message: "Hello from mobile API!",
      authenticated: true,
      authSource: 'api-key',
      userId: 'mobile-user',
      timestamp: new Date().toISOString(),
    },
    message: "Mobile API access successful"
  });
}

/**
 * Example public API route for mobile apps
 * No authentication required
 */
export async function POST(request: NextRequest) {
  const body = await request.json();
  
  return NextResponse.json({
    success: true,
    data: {
      message: "Public mobile API endpoint",
      data: body,
      authenticated: false,
      authSource: 'none',
      timestamp: new Date().toISOString(),
    },
    message: "Public mobile API access successful"
  });
}
