# Mobile API Integration Guide

This guide explains how to configure and use the backend API with mobile applications (Flutter, React Native, etc.).

## 🚀 Quick Start

### 1. Environment Configuration

Add these variables to your `.env` file:

```bash
# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001,https://yourdomain.com
MOBILE_APP_ORIGINS=capacitor://localhost,ionic://localhost,http://localhost:8100

# API Key Configuration (for mobile apps)
API_KEY_SECRET=your-secure-api-key-secret-here
MOBILE_API_KEYS=mobile-key-1,mobile-key-2,mobile-key-3
```

### 2. CORS Configuration

**Important**: CORS is only needed for web views and PWAs, NOT for native mobile apps!

- **Web Origins**: Configured via `ALLOWED_ORIGINS` (for web browsers)
- **Mobile Web Views**: Configured via `MOBILE_APP_ORIGINS` (for Capacitor/Ionic web views)
- **Native Mobile Apps**: No CORS needed - they can connect from anywhere in the world

### 3. API Key Authentication

Mobile apps can authenticate using API keys instead of JWT tokens:

```http
X-API-Key: mobile-key-1
```

## 📱 Mobile App Integration

### Flutter Example

```dart
import 'package:http/http.dart' as http;
import 'dart:convert';

class ApiService {
  static const String baseUrl = 'http://localhost:3000/api';
  static const String apiKey = 'mobile-key-1';
  
  static Map<String, String> get headers => {
    'Content-Type': 'application/json',
    'X-API-Key': apiKey,
  };
  
  static Future<Map<String, dynamic>> getExample() async {
    final response = await http.get(
      Uri.parse('$baseUrl/mobile/example'),
      headers: headers,
    );
    
    if (response.statusCode == 200) {
      return json.decode(response.body);
    } else {
      throw Exception('Failed to load data');
    }
  }
  
  static Future<Map<String, dynamic>> postExample(Map<String, dynamic> data) async {
    final response = await http.post(
      Uri.parse('$baseUrl/mobile/example'),
      headers: headers,
      body: json.encode(data),
    );
    
    if (response.statusCode == 200) {
      return json.decode(response.body);
    } else {
      throw Exception('Failed to post data');
    }
  }
}
```

### React Native Example

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:3000/api',
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': 'mobile-key-1',
  },
});

export const mobileApi = {
  getExample: () => api.get('/mobile/example'),
  postExample: (data) => api.post('/mobile/example', data),
};
```

## 🔐 Authentication Methods

### 1. API Key Authentication (Recommended for Mobile)

```http
GET /api/mobile/example
X-API-Key: mobile-key-1
```

**Pros:**
- Simple to implement
- No token refresh needed
- Works well with mobile apps
- Stateless

**Cons:**
- Less secure than JWT
- Harder to revoke individual sessions

### 2. JWT Authentication (For Web Apps)

```http
GET /api/some-endpoint
Authorization: Bearer <jwt-token>
```

**Pros:**
- More secure
- Can include user information
- Can be revoked
- Stateless

**Cons:**
- Requires token refresh logic
- More complex for mobile apps

## 🛠️ API Endpoints

### Public Endpoints (No Authentication Required)

```http
GET /api/health
GET /api/settings/public
GET /api/projects
GET /api/experience
GET /api/education
GET /api/frameworks
GET /api/testimonials
GET /api/blog
```

### Protected Endpoints (Authentication Required)

**These routes are protected by middleware and require `X-API-Key` header:**

```http
# Mobile API Key Authentication
GET /api/mobile/example
X-API-Key: mobile-key-1

# Admin Endpoints (API Key or JWT)
GET /api/admin/projects
X-API-Key: mobile-key-1
# OR
Authorization: Bearer <jwt-token>
```

**Protected Routes:**
- `/api/mobile/*` - All mobile-specific endpoints
- `/api/admin/*` - All admin endpoints

## 🔧 Configuration Options

### Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `ALLOWED_ORIGINS` | Web browser origins | `http://localhost:3000` | `http://localhost:3000,https://app.com` |
| `MOBILE_APP_ORIGINS` | Mobile web view origins | `capacitor://localhost` | `capacitor://localhost,ionic://localhost` |
| `API_KEY_SECRET` | Secret for API key generation | `your-api-key-secret` | `super-secret-key-123` |
| `MOBILE_API_KEYS` | Valid API keys for mobile apps | `mobile-key-1,mobile-key-2` | `key1,key2,key3` |

### 🌍 **Global Mobile App Access**

**Native mobile apps (Flutter, React Native) can connect from anywhere in the world without CORS restrictions!**

- ✅ **No origin restrictions** for native apps
- ✅ **Works globally** - users from any country can use your app
- ✅ **Only API key authentication** is needed
- ✅ **No CORS configuration** required for native apps

### CORS Headers

The backend automatically sets these CORS headers:

```
Access-Control-Allow-Origin: <origin>
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
Access-Control-Allow-Headers: Content-Type, Authorization, X-API-Key, X-Requested-With
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 86400
```

## 🧪 Testing

### Test API Key Authentication

```bash
# Test protected endpoint
curl -X GET http://localhost:3000/api/mobile/example \
  -H "X-API-Key: mobile-key-1" \
  -H "Content-Type: application/json"

# Test public endpoint
curl -X POST http://localhost:3000/api/mobile/example \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'
```

### Test CORS

```bash
# Test from different origin
curl -X GET http://localhost:3000/api/health \
  -H "Origin: http://localhost:8100" \
  -H "X-API-Key: mobile-key-1" \
  -v
```

## 🚨 Security Considerations

### API Key Security

1. **Rotate Keys Regularly**: Change API keys periodically
2. **Use HTTPS**: Always use HTTPS in production
3. **Limit Permissions**: Different API keys for different access levels
4. **Monitor Usage**: Log API key usage for monitoring

### CORS Security

1. **Specific Origins**: Only allow necessary origins
2. **No Wildcards**: Avoid using `*` in production
3. **HTTPS Only**: Use HTTPS origins in production
4. **Regular Review**: Review allowed origins regularly

## 🔄 Migration from JWT to API Key

If you're migrating from JWT to API key authentication:

1. **Update Environment**: Add API key configuration
2. **Update Mobile App**: Replace JWT logic with API key headers
3. **Test Thoroughly**: Test all endpoints with new authentication
4. **Gradual Rollout**: Deploy gradually to avoid breaking changes

## 📚 Additional Resources

- [Next.js Middleware Documentation](https://nextjs.org/docs/app/building-your-application/routing/middleware)
- [CORS Best Practices](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [API Key Authentication Best Practices](https://swagger.io/docs/specification/authentication/api-keys/)

## 🆘 Troubleshooting

### Common Issues

1. **CORS Errors**: Check `ALLOWED_ORIGINS` and `MOBILE_APP_ORIGINS`
2. **401 Unauthorized**: Verify API key is correct and in `MOBILE_API_KEYS`
3. **Network Errors**: Ensure backend is running and accessible
4. **Preflight Failures**: Check OPTIONS request handling

### Debug Mode

Enable debug logging by setting:

```bash
NODE_ENV=development
```

This will provide detailed error messages and request logging.
