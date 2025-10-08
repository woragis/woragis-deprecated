import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import '../config/env_config.dart';
import '../stores/auth_store.dart';

/// Centralized API client that handles all HTTP requests with automatic
/// authentication headers (API key and JWT token)
class ApiClient {
  final http.Client _httpClient;
  final AuthStoreBloc? _authStore;
  final String _baseUrl;

  ApiClient({
    http.Client? httpClient,
    AuthStoreBloc? authStore,
    String? baseUrl,
  })  : _httpClient = httpClient ?? http.Client(),
        _authStore = authStore,
        _baseUrl = baseUrl ?? EnvConfig.apiBaseUrl;

  /// Build headers with API key and optional auth token
  Map<String, String> _buildHeaders({Map<String, String>? additionalHeaders}) {
    final headers = <String, String>{
      'Content-Type': 'application/json',
      EnvConfig.apiKeyHeader: EnvConfig.apiKey,
    };

    // Add Authorization header if user is authenticated
    if (_authStore != null && _authStore.state.accessToken != null) {
      headers['Authorization'] = 'Bearer ${_authStore.state.accessToken}';
      log('🔐 Authorization header added: Bearer ${_authStore.state.accessToken?.substring(0, 20)}...');
    } else {
      log('⚠️ No auth token available, using API key only');
    }

    // Add any additional headers
    if (additionalHeaders != null) {
      headers.addAll(additionalHeaders);
    }

    if (EnvConfig.networkLogging) {
      log('📤 Request Headers: ${headers.keys.join(", ")}');
    }

    return headers;
  }

  /// Build full URL from path
  String _buildUrl(String path) {
    // Remove leading slash if present
    final cleanPath = path.startsWith('/') ? path.substring(1) : path;
    return '$_baseUrl/$cleanPath';
  }

  /// GET request
  Future<http.Response> get(
    String path, {
    Map<String, String>? queryParameters,
    Map<String, String>? headers,
  }) async {
    final uri = Uri.parse(_buildUrl(path));
    final uriWithParams = queryParameters != null
        ? uri.replace(queryParameters: queryParameters)
        : uri;

    if (EnvConfig.networkLogging) {
      log('📡 GET ${uriWithParams.toString()}');
    }

    try {
      final response = await _httpClient.get(
        uriWithParams,
        headers: _buildHeaders(additionalHeaders: headers),
      );

      _logResponse(response, 'GET', uriWithParams.toString());
      return response;
    } catch (e) {
      log('❌ GET request failed: $e');
      rethrow;
    }
  }

  /// POST request
  Future<http.Response> post(
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? headers,
  }) async {
    final uri = Uri.parse(_buildUrl(path));

    if (EnvConfig.networkLogging) {
      log('📡 POST ${uri.toString()}');
      if (body != null) {
        log('📦 Request Body: ${json.encode(body)}');
      }
    }

    try {
      final response = await _httpClient.post(
        uri,
        headers: _buildHeaders(additionalHeaders: headers),
        body: body != null ? json.encode(body) : null,
      );

      _logResponse(response, 'POST', uri.toString());
      return response;
    } catch (e) {
      log('❌ POST request failed: $e');
      rethrow;
    }
  }

  /// PUT request
  Future<http.Response> put(
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? headers,
  }) async {
    final uri = Uri.parse(_buildUrl(path));

    if (EnvConfig.networkLogging) {
      log('📡 PUT ${uri.toString()}');
      if (body != null) {
        log('📦 Request Body: ${json.encode(body)}');
      }
    }

    try {
      final response = await _httpClient.put(
        uri,
        headers: _buildHeaders(additionalHeaders: headers),
        body: body != null ? json.encode(body) : null,
      );

      _logResponse(response, 'PUT', uri.toString());
      return response;
    } catch (e) {
      log('❌ PUT request failed: $e');
      rethrow;
    }
  }

  /// PATCH request
  Future<http.Response> patch(
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? headers,
  }) async {
    final uri = Uri.parse(_buildUrl(path));

    if (EnvConfig.networkLogging) {
      log('📡 PATCH ${uri.toString()}');
      if (body != null) {
        log('📦 Request Body: ${json.encode(body)}');
      }
    }

    try {
      final response = await _httpClient.patch(
        uri,
        headers: _buildHeaders(additionalHeaders: headers),
        body: body != null ? json.encode(body) : null,
      );

      _logResponse(response, 'PATCH', uri.toString());
      return response;
    } catch (e) {
      log('❌ PATCH request failed: $e');
      rethrow;
    }
  }

  /// DELETE request
  Future<http.Response> delete(
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? headers,
  }) async {
    final uri = Uri.parse(_buildUrl(path));

    if (EnvConfig.networkLogging) {
      log('📡 DELETE ${uri.toString()}');
    }

    try {
      final response = await _httpClient.delete(
        uri,
        headers: _buildHeaders(additionalHeaders: headers),
        body: body != null ? json.encode(body) : null,
      );

      _logResponse(response, 'DELETE', uri.toString());
      return response;
    } catch (e) {
      log('❌ DELETE request failed: $e');
      rethrow;
    }
  }

  /// Log response details
  void _logResponse(http.Response response, String method, String url) {
    if (!EnvConfig.networkLogging) return;

    final statusCode = response.statusCode;
    final emoji = statusCode >= 200 && statusCode < 300
        ? '✅'
        : statusCode >= 400 && statusCode < 500
            ? '⚠️'
            : '❌';

    log('$emoji $method $url - Status: $statusCode');

    // Log response body for errors
    if (statusCode >= 400) {
      try {
        final data = json.decode(response.body);
        log('📥 Error Response: ${json.encode(data)}');
      } catch (e) {
        log('📥 Error Response (raw): ${response.body}');
      }
    }
  }

  /// Close the HTTP client
  void close() {
    _httpClient.close();
  }

  /// Dispose resources
  void dispose() {
    close();
  }
}

