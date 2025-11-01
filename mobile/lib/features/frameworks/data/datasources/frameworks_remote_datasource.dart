import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import '../../../../core/error/exceptions.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/framework_entity.dart';
import '../models/framework_model.dart';

abstract class FrameworksRemoteDataSource {
  Future<List<FrameworkEntity>> getFrameworks({
    int? page,
    int? limit,
    bool? visible,
    bool? public,
    String? type,
    String? search,
    String? sortBy,
    String? sortOrder,
  });

  Future<FrameworkEntity> getFrameworkById(String id);
  Future<FrameworkEntity> getFrameworkBySlug(String slug);
  
  Future<FrameworkEntity> createFramework({
    required String name,
    required String slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    String? type,
    String? proficiencyLevel,
    required int order,
    required bool visible,
    required bool public,
  });

  Future<FrameworkEntity> updateFramework({
    required String id,
    String? name,
    String? slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    String? type,
    String? proficiencyLevel,
    int? order,
    bool? visible,
    bool? public,
  });

  Future<void> deleteFramework(String id);
  Future<void> updateFrameworkOrder(List<Map<String, dynamic>> frameworkOrders);
  Future<int> getFrameworkProjectCount(String frameworkId);
  Future<List<FrameworkEntity>> getFrameworksWithProjectCount();
}

class FrameworksRemoteDataSourceImpl implements FrameworksRemoteDataSource {
  final ApiClient _apiClient;
  final String _baseUrl;

  // Simple in-memory cache
  static final Map<String, dynamic> _cache = {};
  static final Map<String, DateTime> _cacheTimestamps = {};
  static const Duration _cacheDuration = Duration(minutes: 5);

  FrameworksRemoteDataSourceImpl({
    required ApiClient apiClient,
    required String baseUrl,
  })  : _apiClient = apiClient,
        _baseUrl = baseUrl;

  // Helper method to get cached data or fetch fresh
  Future<T> _getCachedOrFetch<T>(
    String cacheKey,
    Future<T> Function() fetcher,
  ) async {
    final now = DateTime.now();
    
    // Check if we have cached data that's still fresh
    if (_cache.containsKey(cacheKey) && 
        _cacheTimestamps.containsKey(cacheKey) &&
        now.difference(_cacheTimestamps[cacheKey]!).compareTo(_cacheDuration) < 0) {
      log('📦 Using cached data for: $cacheKey');
      return _cache[cacheKey] as T;
    }
    
    // Fetch fresh data
    log('🌐 Fetching fresh data for: $cacheKey');
    final data = await fetcher();
    _cache[cacheKey] = data;
    _cacheTimestamps[cacheKey] = now;
    return data;
  }

  // Helper method to invalidate cache
  void _invalidateCache(String pattern) {
    _cache.removeWhere((key, value) => key.contains(pattern));
    _cacheTimestamps.removeWhere((key, value) => key.contains(pattern));
    log('🗑️ Cache invalidated for pattern: $pattern');
  }
  
  @override
  Future<List<FrameworkEntity>> getFrameworks({
    int? page,
    int? limit,
    bool? visible,
    bool? public,
    String? type,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) async {
    final cacheKey = 'frameworks_${page}_${limit}_${visible}_${public}_${type}_${search}_${sortBy}_${sortOrder}';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        final queryParams = <String, String>{};
        if (page != null) queryParams['page'] = page.toString();
        if (limit != null) queryParams['limit'] = limit.toString();
        if (visible != null) queryParams['visible'] = visible.toString();
        if (public != null) queryParams['public'] = public.toString();
        if (type != null) queryParams['type'] = type;
        if (search != null && search.isNotEmpty) queryParams['search'] = search;
        if (sortBy != null) queryParams['sortBy'] = sortBy;
        if (sortOrder != null) queryParams['sortOrder'] = sortOrder;

        log('🔍 Frameworks API Request: /admin/frameworks');

        final response = await _apiClient.get('/admin/frameworks', queryParameters: queryParams);

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            final frameworks = (data['data'] as List)
                .map((frameworkJson) => FrameworkModel.fromJson(frameworkJson).toEntity())
                .toList();
            return frameworks;
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch frameworks');
          }
        } else {
          throw ServerException('Failed to fetch frameworks with status ${response.statusCode}');
        }
      } on http.ClientException {
        throw NetworkException('Network error occurred');
      } catch (e) {
        if (e is ServerException || e is NetworkException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }

  @override
  Future<FrameworkEntity> getFrameworkById(String id) async {
    final cacheKey = 'framework_$id';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 Framework by ID API Request: /admin/frameworks/$id');

        final response = await _apiClient.get('/admin/frameworks/$id');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            return FrameworkModel.fromJson(data['data']).toEntity();
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch framework');
          }
        } else if (response.statusCode == 404) {
          throw NotFoundException('Framework not found');
        } else {
          throw ServerException('Failed to fetch framework with status ${response.statusCode}');
        }
      } on http.ClientException {
        throw NetworkException('Network error occurred');
      } catch (e) {
        if (e is ServerException || e is NetworkException || e is NotFoundException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }

  @override
  Future<FrameworkEntity> getFrameworkBySlug(String slug) async {
    final cacheKey = 'framework_slug_$slug';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 Framework by Slug API Request: /admin/frameworks/slug/$slug');

        final response = await _apiClient.get('/admin/frameworks/slug/$slug');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            return FrameworkModel.fromJson(data['data']).toEntity();
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch framework');
          }
        } else if (response.statusCode == 404) {
          throw NotFoundException('Framework not found');
        } else {
          throw ServerException('Failed to fetch framework with status ${response.statusCode}');
        }
      } on http.ClientException {
        throw NetworkException('Network error occurred');
      } catch (e) {
        if (e is ServerException || e is NetworkException || e is NotFoundException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }

  @override
  Future<FrameworkEntity> createFramework({
    required String name,
    required String slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    String? type,
    String? proficiencyLevel,
    required int order,
    required bool visible,
    required bool public,
  }) async {
    try {
      final response = await _apiClient.post(
        '/admin/frameworks',
        headers: {'Content-Type': 'application/json'},
        body: {
          'name': name,
          'slug': slug,
          if (description != null) 'description': description,
          if (icon != null) 'icon': icon,
          if (color != null) 'color': color,
          if (website != null) 'website': website,
          if (type != null) 'type': type,
          if (proficiencyLevel != null) 'proficiencyLevel': proficiencyLevel,
          'order': order,
          'visible': visible,
          'public': public,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final framework = FrameworkModel.fromJson(data['data']).toEntity();
          
          // Invalidate frameworks cache
          _invalidateCache('frameworks_');
          
          return framework;
        } else {
          throw ServerException(data['message'] ?? 'Failed to create framework');
        }
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to create framework with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<FrameworkEntity> updateFramework({
    required String id,
    String? name,
    String? slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    String? type,
    String? proficiencyLevel,
    int? order,
    bool? visible,
    bool? public,
  }) async {
    try {
      final response = await _apiClient.put(
        '/admin/frameworks/$id',
        headers: {'Content-Type': 'application/json'},
        body: {
          if (name != null) 'name': name,
          if (slug != null) 'slug': slug,
          if (description != null) 'description': description,
          if (icon != null) 'icon': icon,
          if (color != null) 'color': color,
          if (website != null) 'website': website,
          if (type != null) 'type': type,
          if (proficiencyLevel != null) 'proficiencyLevel': proficiencyLevel,
          if (order != null) 'order': order,
          if (visible != null) 'visible': visible,
          if (public != null) 'public': public,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final framework = FrameworkModel.fromJson(data['data']).toEntity();
          
          // Invalidate frameworks cache
          _invalidateCache('frameworks_');
          _invalidateCache('framework_$id');
          
          return framework;
        } else {
          throw ServerException(data['message'] ?? 'Failed to update framework');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Framework not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update framework with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> deleteFramework(String id) async {
    try {
      final response = await _apiClient.delete('/admin/frameworks/$id');

      if (response.statusCode == 200 || response.statusCode == 204) {
        // Invalidate frameworks cache
        _invalidateCache('frameworks_');
        _invalidateCache('framework_$id');
      } else {
        if (response.statusCode == 404) {
          throw NotFoundException('Framework not found');
        } else {
          throw ServerException('Failed to delete framework with status ${response.statusCode}');
        }
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> updateFrameworkOrder(List<Map<String, dynamic>> frameworkOrders) async {
    try {
      final response = await _apiClient.put(
        '/admin/frameworks/order',
        body: {'frameworkOrders': frameworkOrders},
      );

      if (response.statusCode != 200) {
        throw ServerException('Failed to update framework order with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<int> getFrameworkProjectCount(String frameworkId) async {
    try {
      final response = await _apiClient.get('/admin/frameworks/$frameworkId/project-count');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return data['data']['count'] as int;
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch project count');
        }
      } else {
        throw ServerException('Failed to fetch project count with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<List<FrameworkEntity>> getFrameworksWithProjectCount() async {
    try {
      final response = await _apiClient.get('/admin/frameworks/with-count');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final frameworks = (data['data'] as List)
              .map((frameworkJson) => FrameworkModel.fromJson(frameworkJson).toEntity())
              .toList();
          return frameworks;
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch frameworks with project count');
        }
      } else {
        throw ServerException('Failed to fetch frameworks with project count with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }
}