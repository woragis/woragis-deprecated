#!/bin/bash

# Production build script
# Uses build-time environment variables

# Check if required environment variables are set
if [ -z "$PROD_API_BASE_URL" ]; then
    echo "❌ Error: PROD_API_BASE_URL environment variable is required"
    echo "Example: export PROD_API_BASE_URL=https://yourdomain.com/api"
    exit 1
fi

if [ -z "$PROD_API_KEY" ]; then
    echo "❌ Error: PROD_API_KEY environment variable is required"
    echo "Example: export PROD_API_KEY=your-production-api-key"
    exit 1
fi

echo "🚀 Building production APK..."

flutter build apk --release \
  --dart-define=ENVIRONMENT=production \
  --dart-define=API_BASE_URL="$PROD_API_BASE_URL" \
  --dart-define=API_KEY="$PROD_API_KEY" \
  --dart-define=API_KEY_HEADER=X-API-Key

echo "✅ Production APK built successfully!"
echo "📱 APK location: build/app/outputs/flutter-apk/app-release.apk"
echo "🔐 API Key: $PROD_API_KEY"
echo "🌐 API URL: $PROD_API_BASE_URL"
