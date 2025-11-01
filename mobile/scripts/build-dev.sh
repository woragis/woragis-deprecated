#!/bin/bash

# Development build script
# Uses .env file for configuration

echo "🔧 Building development APK..."

flutter build apk --debug \
  --dart-define=ENVIRONMENT=development \
  --dart-define=API_BASE_URL=http://localhost:3000/api \
  --dart-define=API_KEY=mobile-key-1 \
  --dart-define=API_KEY_HEADER=X-API-Key

echo "✅ Development APK built successfully!"
echo "📱 APK location: build/app/outputs/flutter-apk/app-debug.apk"
