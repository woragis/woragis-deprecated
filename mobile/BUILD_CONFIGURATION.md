# Build Configuration Guide

This document explains how environment variables work in different build scenarios and how to properly configure your app for production.

## 🔧 **Environment Variable Behavior**

### **Development (Debug Mode)**
- ✅ **`.env` file is read** from the project directory
- ✅ **Hot reload** picks up changes to `.env`
- ✅ **Flexible configuration** for testing

### **Release Builds (APK/IPA)**
- ❌ **`.env` file is NOT included** in the APK
- ✅ **Build-time variables** are embedded in the app
- ✅ **Secure configuration** for production

## 🚀 **Build Methods**

### **Method 1: Using Build Scripts (Recommended)**

#### **Development Build**
```bash
cd mobile
./scripts/build-dev.sh
```

#### **Production Build**
```bash
# Set production environment variables
export PROD_API_BASE_URL=https://yourdomain.com/api
export PROD_API_KEY=your-production-api-key

# Build production APK
cd mobile
./scripts/build-prod.sh
```

### **Method 2: Manual Flutter Commands**

#### **Development Build**
```bash
cd mobile
flutter build apk --debug \
  --dart-define=ENVIRONMENT=development \
  --dart-define=API_BASE_URL=http://localhost:3000/api \
  --dart-define=API_KEY=mobile-key-1 \
  --dart-define=API_KEY_HEADER=X-API-Key
```

#### **Production Build**
```bash
cd mobile
flutter build apk --release \
  --dart-define=ENVIRONMENT=production \
  --dart-define=API_BASE_URL=https://yourdomain.com/api \
  --dart-define=API_KEY=your-production-api-key \
  --dart-define=API_KEY_HEADER=X-API-Key
```

## 🔐 **Security Considerations**

### **API Keys in Production**

1. **Never commit production API keys** to version control
2. **Use environment variables** for production builds
3. **Rotate API keys** regularly
4. **Use different keys** for different environments

### **Environment Separation**

```bash
# Development
API_KEY=mobile-key-1
API_BASE_URL=http://localhost:3000/api

# Staging
API_KEY=staging-mobile-key-123
API_BASE_URL=https://staging.yourdomain.com/api

# Production
API_KEY=prod-mobile-key-456
API_BASE_URL=https://yourdomain.com/api
```

## 📱 **Build Configuration Priority**

The app uses this priority order for configuration:

1. **`.env` file** (development only)
2. **Build-time variables** (`--dart-define`)
3. **Default values** (fallback)

```dart
// This is how the configuration works:
static String get apiKey => dotenv.get(
  'API_KEY',  // 1. Try .env file first
  fallback: const String.fromEnvironment(  // 2. Then build-time variable
    'API_KEY', 
    defaultValue: 'mobile-key-1',  // 3. Finally default value
  ),
);
```

## 🛠️ **CI/CD Integration**

### **GitHub Actions Example**

```yaml
name: Build Production APK

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Flutter
        uses: subosito/flutter-action@v2
        with:
          flutter-version: '3.16.0'
          
      - name: Build Production APK
        run: |
          cd mobile
          flutter build apk --release \
            --dart-define=ENVIRONMENT=production \
            --dart-define=API_BASE_URL=${{ secrets.PROD_API_BASE_URL }} \
            --dart-define=API_KEY=${{ secrets.PROD_API_KEY }} \
            --dart-define=API_KEY_HEADER=X-API-Key
```

### **Environment Variables in CI/CD**

Set these as secrets in your CI/CD platform:

- `PROD_API_BASE_URL` - Production API URL
- `PROD_API_KEY` - Production API key
- `STAGING_API_BASE_URL` - Staging API URL
- `STAGING_API_KEY` - Staging API key

## 🧪 **Testing Different Configurations**

### **Local Testing**

```bash
# Test with different API URL
flutter run --dart-define=API_BASE_URL=https://staging.yourdomain.com/api

# Test with different API key
flutter run --dart-define=API_KEY=test-key-123
```

### **Build Testing**

```bash
# Test production build locally
export PROD_API_BASE_URL=https://yourdomain.com/api
export PROD_API_KEY=your-production-key
./scripts/build-prod.sh
```

## 📋 **Build Checklist**

### **Before Building Production APK**

- [ ] Set production environment variables
- [ ] Verify API key is correct
- [ ] Test API connectivity
- [ ] Check API endpoints are accessible
- [ ] Verify CORS configuration (if needed)
- [ ] Test the APK on a device

### **After Building**

- [ ] Test APK installation
- [ ] Verify API calls work
- [ ] Check error handling
- [ ] Test offline scenarios
- [ ] Verify security (no sensitive data in logs)

## 🔍 **Debugging Build Issues**

### **Check Build Configuration**

Add this to your app to verify configuration:

```dart
// In your main.dart or test page
void printBuildConfig() {
  print('=== Build Configuration ===');
  print('Environment: ${EnvConfig.environment}');
  print('API Base URL: ${EnvConfig.apiBaseUrl}');
  print('API Key: ${EnvConfig.apiKey}');
  print('==========================');
}
```

### **Common Issues**

1. **API calls failing in production**
   - Check if API key is correctly set
   - Verify API base URL is accessible
   - Check network connectivity

2. **Build-time variables not working**
   - Ensure `--dart-define` syntax is correct
   - Check for typos in variable names
   - Verify quotes around values

3. **Environment variables not loading**
   - Check `.env` file exists in mobile directory
   - Verify `EnvConfig.init()` is called
   - Check file permissions

## 📚 **Additional Resources**

- [Flutter Build Configuration](https://docs.flutter.dev/deployment/flutter-for-android)
- [Environment Variables in Flutter](https://docs.flutter.dev/deployment/environment-variables)
- [Dart Define Documentation](https://dart.dev/tools/dart-compile#--dart-define)
