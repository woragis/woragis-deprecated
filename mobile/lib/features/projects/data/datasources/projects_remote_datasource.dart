import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import '../../../../core/error/exceptions.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/project_entity.dart';
import '../models/project_model.dart';

abstract class ProjectsRemoteDataSource {
  Future<List<ProjectEntity>> getProjects({
    int? page,
    int? limit,
    bool? featured,
    bool? visible,
    bool? public,
    String? search,
  });

  Future<ProjectEntity> getProjectById(String id);
  Future<ProjectEntity> getProjectBySlug(String slug);

  Future<ProjectEntity> createProject({
    required String title,
    required String slug,
    required String description,
    String? longDescription,
    String? content,
    String? videoUrl,
    required String image,
    String? githubUrl,
    String? liveUrl,
    required bool featured,
    required int order,
    required bool visible,
    required bool public,
    List<String>? frameworkIds,
  });

  Future<ProjectEntity> updateProject({
    required String id,
    String? title,
    String? slug,
    String? description,
    String? longDescription,
    String? content,
    String? videoUrl,
    String? image,
    String? githubUrl,
    String? liveUrl,
    bool? featured,
    int? order,
    bool? visible,
    bool? public,
    List<String>? frameworkIds,
  });

  Future<void> deleteProject(String id);
  
  // Project interaction methods
  Future<void> incrementViewCount(String id);
  Future<void> incrementLikeCount(String id);


  Future<void> updateProjectOrder({
    required String id,
    required int order,
  });

  // Project frameworks relationship methods
  Future<void> assignFrameworkToProject({
    required String projectId,
    required String frameworkId,
  });

  Future<void> removeFrameworkFromProject({
    required String projectId,
    required String frameworkId,
  });

  Future<List<String>> getProjectFrameworkIds(String projectId);

  // Toggle methods
  Future<void> toggleProjectFeatured(String id);
  Future<void> toggleProjectVisibility(String id);
}

class ProjectsRemoteDataSourceImpl implements ProjectsRemoteDataSource {
  final ApiClient _apiClient;

  // Simple in-memory cache
  static final Map<String, dynamic> _cache = {};
  static final Map<String, DateTime> _cacheTimestamps = {};
  static const Duration _cacheDuration = Duration(minutes: 5);

  ProjectsRemoteDataSourceImpl({
    required ApiClient apiClient,
  }) : _apiClient = apiClient;

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
  Future<List<ProjectEntity>> getProjects({
    int? page,
    int? limit,
    bool? featured,
    bool? visible,
    bool? public,
    String? search,
  }) async {
    final cacheKey = 'projects_${page}_${limit}_${featured}_${visible}_${public}_${search}';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        final queryParams = <String, dynamic>{};
        if (page != null) queryParams['page'] = page;
        if (limit != null) queryParams['limit'] = limit;
        if (featured != null) queryParams['featured'] = featured;
        if (visible != null) queryParams['visible'] = visible;
        if (public != null) queryParams['public'] = public;
        if (search != null && search.isNotEmpty) queryParams['search'] = search;

        log('🔍 Projects List API Request: /admin/projects');

        final response = await _apiClient.get(
          'admin/projects',
          queryParameters: queryParams.map((key, value) => MapEntry(key, value.toString())),
        );

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            final projectsList = (data['data'] as List)
                .map((projectJson) => ProjectModel.fromApiJson(projectJson).toEntity())
                .toList();
            return projectsList;
          } else {
            throw ServerException(data['error'] ?? data['message'] ?? 'Failed to fetch projects');
          }
        } else {
          throw ServerException('Failed to fetch projects with status ${response.statusCode}');
        }
      } catch (e) {
        if (e is ServerException || e is NetworkException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }

  @override
  Future<ProjectEntity> getProjectById(String id) async {
    final cacheKey = 'project_$id';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 Project by ID API Request: /admin/projects/$id');

        final response = await _apiClient.get('/admin/projects/$id');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            final project = ProjectModel.fromApiJson(data['data']);
            return project.toEntity();
          } else {
            throw ServerException(data['error'] ?? data['message'] ?? 'Failed to fetch project');
          }
        } else {
          throw ServerException('Failed to fetch project with status ${response.statusCode}');
        }
      } catch (e) {
        if (e is ServerException || e is NetworkException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }


  @override
  Future<ProjectEntity> createProject({
    required String title,
    required String slug,
    required String description,
    String? longDescription,
    String? content,
    String? videoUrl,
    required String image,
    String? githubUrl,
    String? liveUrl,
    required bool featured,
    required int order,
    required bool visible,
    required bool public,
    List<String>? frameworkIds,
  }) async {
    try {
      log('🔍 Create Project API Request: /admin/projects');

      final response = await _apiClient.post(
        '/admin/projects',
        body: {
          'title': title,
          'slug': slug,
          'description': description,
          'longDescription': longDescription,
          'content': content,
          'videoUrl': videoUrl,
          'image': image,
          'githubUrl': githubUrl,
          'liveUrl': liveUrl,
          'featured': featured,
          'order': order,
          'visible': visible,
          'public': public,
          if (frameworkIds != null) 'frameworkIds': frameworkIds,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final project = ProjectModel.fromApiJson(data['data']);
          
          // Invalidate projects cache
          _invalidateCache('projects_');
          
          return project.toEntity();
        } else {
          throw ServerException(data['error'] ?? data['message'] ?? 'Failed to create project');
        }
      } else {
        throw ServerException('Failed to create project with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<ProjectEntity> updateProject({
    required String id,
    String? title,
    String? slug,
    String? description,
    String? longDescription,
    String? content,
    String? videoUrl,
    String? image,
    String? githubUrl,
    String? liveUrl,
    bool? featured,
    int? order,
    bool? visible,
    bool? public,
    List<String>? frameworkIds,
  }) async {
    try {
      log('🔍 Update Project API Request: /admin/projects/$id');

      final response = await _apiClient.put(
        '/admin/projects/$id',
        body: {
          if (title != null) 'title': title,
          if (slug != null) 'slug': slug,
          if (description != null) 'description': description,
          if (longDescription != null) 'longDescription': longDescription,
          if (content != null) 'content': content,
          if (videoUrl != null) 'videoUrl': videoUrl,
          if (image != null) 'image': image,
          if (githubUrl != null) 'githubUrl': githubUrl,
          if (liveUrl != null) 'liveUrl': liveUrl,
          if (featured != null) 'featured': featured,
          if (order != null) 'order': order,
          if (visible != null) 'visible': visible,
          if (public != null) 'public': public,
          if (frameworkIds != null) 'frameworkIds': frameworkIds,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final project = ProjectModel.fromApiJson(data['data']);
          
          // Invalidate projects cache
          _invalidateCache('projects_');
          _invalidateCache('project_$id');
          
          return project.toEntity();
        } else {
          throw ServerException(data['error'] ?? data['message'] ?? 'Failed to update project');
        }
      } else {
        throw ServerException('Failed to update project with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> deleteProject(String id) async {
    try {
      log('🔍 Delete Project API Request: /admin/projects/$id');

      final response = await _apiClient.delete('/admin/projects/$id');

      if (response.statusCode == 200) {
        // Invalidate projects cache
        _invalidateCache('projects_');
        _invalidateCache('project_$id');
        
        log('✅ Project deleted successfully');
      } else {
        throw ServerException('Failed to delete project with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }


  @override
  Future<void> updateProjectOrder({
    required String id,
    required int order,
  }) async {
    try {
      log('🔍 Update Project Order API Request: /admin/projects/$id/order');

      final response = await _apiClient.put(
        '/admin/projects/$id/order',
        body: {
          'order': order,
        },
      );

      if (response.statusCode == 200) {
        // Invalidate projects cache
        _invalidateCache('projects_');
        _invalidateCache('project_$id');
        
        log('✅ Project order updated successfully');
      } else {
        throw ServerException('Failed to update project order with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  // Project frameworks relationship methods
  @override
  Future<void> assignFrameworkToProject({
    required String projectId,
    required String frameworkId,
  }) async {
    try {
      log('🔍 Assign Framework to Project API Request: /api/admin/projects/$projectId/frameworks');

      final response = await _apiClient.post(
        '/api/admin/projects/$projectId/frameworks',
        body: {
          'frameworkId': frameworkId,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        log('✅ Framework assigned to project successfully');
        // Invalidate cache for this project
        _invalidateCache('project_$projectId');
      } else {
        throw ServerException('Failed to assign framework to project with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> removeFrameworkFromProject({
    required String projectId,
    required String frameworkId,
  }) async {
    try {
      log('🔍 Remove Framework from Project API Request: /api/admin/projects/$projectId/frameworks/$frameworkId');

      final response = await _apiClient.delete('/api/admin/projects/$projectId/frameworks/$frameworkId');

      if (response.statusCode == 200 || response.statusCode == 204) {
        log('✅ Framework removed from project successfully');
        // Invalidate cache for this project
        _invalidateCache('project_$projectId');
      } else {
        throw ServerException('Failed to remove framework from project with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<List<String>> getProjectFrameworkIds(String projectId) async {
    try {
      log('🔍 Get Project Framework IDs API Request: /api/admin/projects/$projectId/frameworks');

      final response = await _apiClient.get('/api/admin/projects/$projectId/frameworks');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final frameworks = (data['data'] as List)
              .map((framework) => framework['id'] as String)
              .toList();
          return frameworks;
        } else {
          throw ServerException(data['error'] ?? data['message'] ?? 'Failed to get project frameworks');
        }
      } else {
        throw ServerException('Failed to get project frameworks with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<ProjectEntity> getProjectBySlug(String slug) async {
    final cacheKey = 'project_slug_$slug';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 Project by Slug API Request: /admin/projects/slug/$slug');

        final response = await _apiClient.get('/admin/projects/slug/$slug');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            final project = ProjectModel.fromApiJson(data['data']);
            return project.toEntity();
          } else {
            throw ServerException(data['error'] ?? data['message'] ?? 'Failed to fetch project');
          }
        } else {
          throw ServerException('Failed to fetch project with status ${response.statusCode}');
        }
      } catch (e) {
        if (e is ServerException || e is NetworkException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }

  @override
  Future<void> incrementViewCount(String id) async {
    try {
      log('🔍 Increment View Count API Request: /admin/projects/$id/view');

      final response = await _apiClient.post('/admin/projects/$id/view');

      if (response.statusCode == 200) {
        log('✅ Project view count incremented successfully');
      } else {
        throw ServerException('Failed to increment view count with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> incrementLikeCount(String id) async {
    try {
      log('🔍 Increment Like Count API Request: /admin/projects/$id/like');

      final response = await _apiClient.post('/admin/projects/$id/like');

      if (response.statusCode == 200) {
        log('✅ Project like count incremented successfully');
      } else {
        throw ServerException('Failed to increment like count with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> toggleProjectFeatured(String id) async {
    try {
      log('🔍 Toggle Project Featured API Request: /admin/projects/$id/toggle-featured');

      final response = await _apiClient.post('/admin/projects/$id/toggle-featured');

      if (response.statusCode == 200) {
        log('✅ Project featured status toggled successfully');
        // Invalidate cache for this project
        _invalidateCache('project_$id');
        _invalidateCache('projects_');
      } else {
        throw ServerException('Failed to toggle project featured status with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> toggleProjectVisibility(String id) async {
    try {
      log('🔍 Toggle Project Visibility API Request: /admin/projects/$id/toggle-visibility');

      final response = await _apiClient.post('/admin/projects/$id/toggle-visibility');

      if (response.statusCode == 200) {
        log('✅ Project visibility toggled successfully');
        // Invalidate cache for this project
        _invalidateCache('project_$id');
        _invalidateCache('projects_');
      } else {
        throw ServerException('Failed to toggle project visibility with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

}
