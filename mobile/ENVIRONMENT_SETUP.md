# Environment Configuration Setup

This document explains the environment configuration implementation for the Woragis Flutter mobile app.

## ✅ What Was Implemented

### 1. Environment Configuration Class
- **File:** `lib/core/config/env_config.dart`
- Loads environment variables from `.env` file using `flutter_dotenv`
- Provides type-safe getters for all configuration values
- Includes fallback values for all variables
- Supports development and production environments

### 2. Environment Variables
- **File:** `mobile/.env` (created from `env.example`)
- Contains all configuration variables:
  - API configuration (base URL, timeout)
  - Authentication settings
  - Feature flags
  - UI preferences
  - Security settings
  - And more...

### 3. Integration Points Updated

#### main.dart
- Added `EnvConfig.init()` to load environment variables on startup
- Added debug config printing in development mode

#### injection_container.dart
- Updated to use `EnvConfig.apiBaseUrl` instead of hardcoded URL
- Imports the `env_config.dart`

#### query_client.dart
- Updated to use `EnvConfig.apiBaseUrl` and `EnvConfig.apiTimeout`
- Properly configured with environment-based values

### 4. Git Configuration
- Added `.env` to `.gitignore` to prevent committing sensitive data
- Kept `env.example` for documentation and setup

### 5. Documentation
- Created `lib/core/config/README.md` with usage instructions
- Comprehensive documentation of all available variables

## 🚀 How to Use

### For Development

1. **The `.env` file is already created** from `env.example`
2. **Edit values if needed:**
   ```bash
   cd mobile
   nano .env  # or use your preferred editor
   ```

3. **Access variables in code:**
   ```dart
   import 'package:woragis/core/config/env_config.dart';
   
   // Use environment variables
   final apiUrl = EnvConfig.apiBaseUrl;
   final isDev = EnvConfig.isDevelopment;
   ```

### For Production

1. **Create a production `.env` file:**
   ```env
   API_BASE_URL=https://api.woragis.com/api
   ENVIRONMENT=production
   DEBUG_MODE=false
   FEATURE_ADMIN=true
   ```

2. **Build the app with the production `.env` file**

## 📝 Key Environment Variables

### Most Important Variables

```env
# API Configuration
API_BASE_URL=http://localhost:3000/api  # Change to your backend URL
API_TIMEOUT=30000

# Environment
ENVIRONMENT=development  # or 'production'
DEBUG_MODE=true  # Set to 'false' in production

# Feature Flags
FEATURE_MONEY=true  # Enable money/business features
FEATURE_ADMIN=false  # Enable admin features (set to true for admin users)
```

### All Available Variables

See `env.example` or `lib/core/config/README.md` for a complete list of all available environment variables.

## ✅ Benefits

1. **Type Safety:** All environment variables have typed getters
2. **Fallback Values:** App won't crash if a variable is missing
3. **Security:** Sensitive data kept out of git
4. **Flexibility:** Easy to switch between development/production
5. **Feature Flags:** Enable/disable features per environment
6. **Centralized Config:** Single source of truth for all configuration

## 🔒 Security Notes

- The `.env` file is gitignored and will NOT be committed to version control
- Never commit sensitive API keys, tokens, or passwords
- Use different `.env` files for different environments
- The `env.example` file is safe to commit (it has no real credentials)

## 🐛 Troubleshooting

### If the app can't load environment variables:

1. **Check if `.env` exists:**
   ```bash
   ls -la mobile/.env
   ```

2. **Verify it's loaded in main.dart:**
   ```dart
   await EnvConfig.init();
   ```

3. **Check the pubspec.yaml has the asset:**
   ```yaml
   flutter:
     assets:
       - .env
   ```

4. **Clean and rebuild:**
   ```bash
   cd mobile
   flutter clean
   flutter pub get
   flutter run
   ```

### Debug config values:

```dart
// In main.dart, this prints all config in debug mode
EnvConfig.printConfig();
```

## 📚 Additional Resources

- See `lib/core/config/README.md` for detailed usage documentation
- See `env.example` for all available environment variables
- See `lib/core/config/env_config.dart` for implementation details




