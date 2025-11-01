# Mobile API Integration Setup

This document explains how to set up and use the API key authentication for the mobile app.

## 🔧 Setup

### 1. Environment Configuration

Create a `.env` file in the `mobile/` directory based on `env.example`:

```bash
# Copy the example file
cp mobile/env.example mobile/.env
```

### 2. Configure API Key

Edit `mobile/.env` and set your API key:

```bash
# API Key Authentication (for mobile app access)
API_KEY=mobile-key-1
API_KEY_HEADER=X-API-Key
```

### 3. Backend Configuration

Make sure your backend `.env` file includes the same API key:

```bash
# API Key Configuration (for mobile apps)
MOBILE_API_KEYS=mobile-key-1,mobile-key-2,mobile-key-3
```

## 📱 Usage

### Automatic API Key Authentication

The mobile app automatically includes the API key in all requests:

```dart
// The API key is automatically added to all requests
final response = await _apiService.testMobileApi();
```

### Manual API Calls

For custom API calls, the API key is automatically included:

```dart
final response = await _queryClient.dio.get('/mobile/example');
// X-API-Key: mobile-key-1 is automatically added
```

### Public Endpoints

Public endpoints (like `/projects`, `/experience`) work without authentication:

```dart
final response = await _apiService.getProjects();
// No authentication required
```

## 🧪 Testing

### Test Page

Use the `ApiTestPage` to test API integration:

```dart
Navigator.push(
  context,
  MaterialPageRoute(
    builder: (context) => const ApiTestPage(),
  ),
);
```

### Available Tests

1. **Mobile API (Protected)** - Tests `/api/mobile/example` with API key
2. **Public Mobile API** - Tests POST to `/api/mobile/example` without auth
3. **Health API** - Tests `/api/health` endpoint
4. **Projects API** - Tests `/api/projects` endpoint

## 🔐 Authentication Flow

### Protected Endpoints (`/api/mobile/*`, `/api/admin/*`)

1. Mobile app automatically adds `X-API-Key` header
2. Backend validates the API key
3. If valid, request proceeds
4. If invalid, returns 401 Unauthorized

### Public Endpoints

- `/api/health`
- `/api/projects`
- `/api/experience`
- `/api/education`
- `/api/frameworks`
- `/api/testimonials`
- `/api/blog`

These endpoints work without authentication.

### JWT Endpoints (Future)

For user-specific endpoints, JWT tokens will be used:

```dart
// JWT authentication (when implemented)
final response = await _queryClient.dio.get('/api/user/profile');
// Authorization: Bearer <jwt-token> is automatically added
```

## 🚨 Error Handling

The mobile app handles common API errors:

- **401 Unauthorized** - Invalid API key or expired JWT
- **403 Forbidden** - Insufficient permissions
- **404 Not Found** - Endpoint not found
- **500 Server Error** - Backend error
- **Connection Timeout** - Network issues

## 🔄 Development vs Production

### Development

```bash
# mobile/.env
API_BASE_URL=http://localhost:3000/api
API_KEY=mobile-key-1
```

### Production

```bash
# mobile/.env
API_BASE_URL=https://yourdomain.com/api
API_KEY=your-production-api-key
```

## 📋 API Endpoints Reference

### Protected Endpoints (Require API Key)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/mobile/example` | GET | Test mobile API |
| `/api/mobile/example` | POST | Test public mobile API |
| `/api/admin/*` | ALL | Admin endpoints |

### Public Endpoints (No Auth Required)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/projects` | GET | Get projects |
| `/api/experience` | GET | Get experience |
| `/api/education` | GET | Get education |
| `/api/frameworks` | GET | Get frameworks |
| `/api/testimonials` | GET | Get testimonials |
| `/api/blog` | GET | Get blog posts |

## 🛠️ Troubleshooting

### API Key Not Working

1. Check that `API_KEY` is set in `mobile/.env`
2. Verify the same key is in backend `MOBILE_API_KEYS`
3. Check network logs for the `X-API-Key` header
4. Test with the `ApiTestPage`

### Connection Issues

1. Verify `API_BASE_URL` is correct
2. Check if backend is running
3. Test with `/api/health` endpoint
4. Check network connectivity

### CORS Issues

CORS is handled automatically by the backend middleware. For native mobile apps, CORS is not needed.

## 📚 Additional Resources

- [Backend API Documentation](../MOBILE_API_INTEGRATION.md)
- [Flutter Dio Documentation](https://pub.dev/packages/dio)
- [Environment Configuration](../lib/core/config/env_config.dart)
