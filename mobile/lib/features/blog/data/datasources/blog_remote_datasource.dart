import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import '../../../../core/error/exceptions.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/blog_post_entity.dart';
import '../../domain/entities/blog_tag_entity.dart';
import '../../domain/entities/blog_stats_entity.dart';
import '../../domain/entities/blog_tag_with_count_entity.dart';
import '../../domain/entities/blog_order_update_entity.dart';
import '../models/blog_post_model.dart';
import '../models/blog_tag_model.dart';
import '../models/blog_stats_model.dart';
import '../models/blog_tag_with_count_model.dart';
import '../models/blog_order_update_model.dart';

abstract class BlogRemoteDataSource {
  // Blog posts
  Future<List<BlogPostEntity>> getBlogPosts({
    int? page,
    int? limit,
    bool? published,
    bool? featured,
    bool? visible,
    bool? public,
    String? search,
    String? sortBy,
    String? sortOrder,
  });

  Future<BlogPostEntity> getBlogPostById(String id);
  Future<BlogPostEntity> getBlogPostBySlug(String slug);
  Future<BlogPostEntity> createBlogPost({
    required String title,
    required String slug,
    required String excerpt,
    required String content,
    String? featuredImage,
    int? readingTime,
    required bool featured,
    required bool published,
    DateTime? publishedAt,
    required int order,
    required bool visible,
    required bool public,
  });
  Future<BlogPostEntity> updateBlogPost({
    required String id,
    String? title,
    String? slug,
    String? excerpt,
    String? content,
    String? featuredImage,
    int? readingTime,
    bool? featured,
    bool? published,
    DateTime? publishedAt,
    int? order,
    bool? visible,
    bool? public,
  });
  Future<void> deleteBlogPost(String id);

  // Blog tags
  Future<List<BlogTagEntity>> getBlogTags({
    int? page,
    int? limit,
    bool? visible,
    String? search,
  });
  Future<BlogTagEntity> getBlogTagById(String id);
  Future<BlogTagEntity> getBlogTagBySlug(String slug);
  Future<BlogTagEntity> createBlogTag({
    required String name,
    required String slug,
    String? description,
    String? color,
  });
  Future<BlogTagEntity> updateBlogTag({
    required String id,
    String? name,
    String? slug,
    String? description,
    String? color,
  });
  Future<void> deleteBlogTag(String id);

  // Blog post tags relationship
  Future<void> assignTagToPost({required String postId, required String tagId});
  Future<void> removeTagFromPost({required String postId, required String tagId});
  Future<List<BlogTagEntity>> getPostTags(String postId);

  // Order management
  Future<void> updateBlogPostOrder(List<BlogPostOrderUpdateEntity> orders);
  Future<void> updateBlogTagOrder(List<BlogTagOrderUpdateEntity> orders);

  // Toggle convenience methods
  Future<BlogPostEntity> toggleBlogPostVisibility(String id);
  Future<BlogPostEntity> toggleBlogPostFeatured(String id);
  Future<BlogPostEntity> toggleBlogPostPublished(String id);

  // Statistics
  Future<void> incrementViewCount(String postId);
  Future<void> incrementLikeCount(String postId);
  Future<BlogStatsEntity> getBlogStats();
  Future<BlogStatsEntity> getPublicBlogStats();

  // Enhanced tag methods
  Future<List<BlogTagWithCountEntity>> getBlogTagsWithCount();
  Future<List<BlogTagWithCountEntity>> getPopularBlogTags({int? limit});
}

class BlogRemoteDataSourceImpl implements BlogRemoteDataSource {
  final ApiClient _apiClient;
  final String _baseUrl;

  // Simple in-memory cache
  static final Map<String, dynamic> _cache = {};
  static final Map<String, DateTime> _cacheTimestamps = {};
  static const Duration _cacheDuration = Duration(minutes: 5);

  BlogRemoteDataSourceImpl({
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
  Future<List<BlogPostEntity>> getBlogPosts({
    int? page,
    int? limit,
    bool? published,
    bool? featured,
    bool? visible,
    bool? public,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) async {
    final cacheKey = 'blog_posts_${page}_${limit}_${published}_${featured}_${visible}_${public}_${search}_${sortBy}_${sortOrder}';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        final queryParams = <String, String>{};
        if (page != null) queryParams['page'] = page.toString();
        if (limit != null) queryParams['limit'] = limit.toString();
        if (published != null) queryParams['published'] = published.toString();
        if (featured != null) queryParams['featured'] = featured.toString();
        if (visible != null) queryParams['visible'] = visible.toString();
        if (public != null) queryParams['public'] = public.toString();
        if (search != null && search.isNotEmpty) queryParams['search'] = search;
        if (sortBy != null) queryParams['sortBy'] = sortBy;
        if (sortOrder != null) queryParams['sortOrder'] = sortOrder;

        log('🔍 Blog Posts API Request: /admin/blog');

        final response = await _apiClient.get('/admin/blog', queryParameters: queryParams);

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            final posts = (data['data']['posts'] as List)
                .map((postJson) => BlogPostModel.fromJson(postJson).toEntity())
                .toList();
            return posts;
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch blog posts');
          }
        } else {
          throw ServerException('Failed to fetch blog posts with status ${response.statusCode}');
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
  Future<BlogPostEntity> getBlogPostById(String id) async {
    final cacheKey = 'blog_post_$id';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 Blog Post by ID API Request: /admin/blog/$id');

        final response = await _apiClient.get('/admin/blog/$id');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            return BlogPostModel.fromJson(data['data']).toEntity();
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch blog post');
          }
        } else if (response.statusCode == 404) {
          throw NotFoundException('Blog post not found');
        } else {
          throw ServerException('Failed to fetch blog post with status ${response.statusCode}');
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
  Future<BlogPostEntity> getBlogPostBySlug(String slug) async {
    final cacheKey = 'blog_post_slug_$slug';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 Blog Post by Slug API Request: /admin/blog/slug/$slug');

        final response = await _apiClient.get('/admin/blog/slug/$slug');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            return BlogPostModel.fromJson(data['data']).toEntity();
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch blog post');
          }
        } else if (response.statusCode == 404) {
          throw NotFoundException('Blog post not found');
        } else {
          throw ServerException('Failed to fetch blog post with status ${response.statusCode}');
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
  Future<BlogPostEntity> createBlogPost({
    required String title,
    required String slug,
    required String excerpt,
    required String content,
    String? featuredImage,
    int? readingTime,
    required bool featured,
    required bool published,
    DateTime? publishedAt,
    required int order,
    required bool visible,
    required bool public,
  }) async {
    try {
      final response = await _apiClient.post(
        '/admin/blog',
        headers: {'Content-Type': 'application/json'},
        body: {
          'title': title,
          'slug': slug,
          'excerpt': excerpt,
          'content': content,
          if (featuredImage != null) 'featuredImage': featuredImage,
          if (readingTime != null) 'readingTime': readingTime,
          'featured': featured,
          'published': published,
          if (publishedAt != null) 'publishedAt': publishedAt.toIso8601String(),
          'order': order,
          'visible': visible,
          'public': public,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final blogPost = BlogPostModel.fromJson(data['data']).toEntity();
          
          // Invalidate blog posts cache
          _invalidateCache('blog_posts_');
          
          return blogPost;
        } else {
          throw ServerException(data['message'] ?? 'Failed to create blog post');
        }
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to create blog post with status ${response.statusCode}');
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
  Future<BlogPostEntity> updateBlogPost({
    required String id,
    String? title,
    String? slug,
    String? excerpt,
    String? content,
    String? featuredImage,
    int? readingTime,
    bool? featured,
    bool? published,
    DateTime? publishedAt,
    int? order,
    bool? visible,
    bool? public,
  }) async {
    try {
      final response = await _apiClient.put(
        '/admin/blog/$id',
        headers: {'Content-Type': 'application/json'},
        body: {
          if (title != null) 'title': title,
          if (slug != null) 'slug': slug,
          if (excerpt != null) 'excerpt': excerpt,
          if (content != null) 'content': content,
          if (featuredImage != null) 'featuredImage': featuredImage,
          if (readingTime != null) 'readingTime': readingTime,
          if (featured != null) 'featured': featured,
          if (published != null) 'published': published,
          if (publishedAt != null) 'publishedAt': publishedAt.toIso8601String(),
          if (order != null) 'order': order,
          if (visible != null) 'visible': visible,
          if (public != null) 'public': public,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final blogPost = BlogPostModel.fromJson(data['data']).toEntity();
          
          // Invalidate blog posts cache
          _invalidateCache('blog_posts_');
          _invalidateCache('blog_post_$id');
          
          return blogPost;
        } else {
          throw ServerException(data['message'] ?? 'Failed to update blog post');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Blog post not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update blog post with status ${response.statusCode}');
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
  Future<void> deleteBlogPost(String id) async {
    try {
      final response = await _apiClient.delete('/admin/blog/$id');

      if (response.statusCode == 200 || response.statusCode == 204) {
        // Invalidate blog posts cache
        _invalidateCache('blog_posts_');
        _invalidateCache('blog_post_$id');
      } else {
        if (response.statusCode == 404) {
          throw NotFoundException('Blog post not found');
        } else {
          throw ServerException('Failed to delete blog post with status ${response.statusCode}');
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

  // Blog Tags Implementation
  @override
  Future<List<BlogTagEntity>> getBlogTags({
    int? page,
    int? limit,
    bool? visible,
    String? search,
  }) async {
    try {
      final queryParams = <String, String>{};
      if (page != null) queryParams['page'] = page.toString();
      if (limit != null) queryParams['limit'] = limit.toString();
      if (visible != null) queryParams['visible'] = visible.toString();
      if (search != null && search.isNotEmpty) queryParams['search'] = search;

      final response = await _apiClient.get('/admin/blog-tags', queryParameters: queryParams);

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final tags = (data['data']['tags'] as List)
              .map((tagJson) => BlogTagModel.fromJson(tagJson).toEntity())
              .toList();
          return tags;
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch blog tags');
        }
      } else {
        throw ServerException('Failed to fetch blog tags with status ${response.statusCode}');
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
  Future<BlogTagEntity> getBlogTagById(String id) async {
    try {
      final response = await _apiClient.get('/admin/blog-tags/$id');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return BlogTagModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch blog tag');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Blog tag not found');
      } else {
        throw ServerException('Failed to fetch blog tag with status ${response.statusCode}');
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
  Future<BlogTagEntity> getBlogTagBySlug(String slug) async {
    try {
      final response = await _apiClient.get('/admin/blog-tags/slug/$slug');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return BlogTagModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch blog tag');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Blog tag not found');
      } else {
        throw ServerException('Failed to fetch blog tag with status ${response.statusCode}');
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
  Future<BlogTagEntity> createBlogTag({
    required String name,
    required String slug,
    String? description,
    String? color,
  }) async {
    try {
      final response = await _apiClient.post(
        '/admin/blog-tags',
        headers: {'Content-Type': 'application/json'},
        body: {
          'name': name,
          'slug': slug,
          if (description != null) 'description': description,
          if (color != null) 'color': color,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return BlogTagModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to create blog tag');
        }
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to create blog tag with status ${response.statusCode}');
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
  Future<BlogTagEntity> updateBlogTag({
    required String id,
    String? name,
    String? slug,
    String? description,
    String? color,
  }) async {
    try {
      final response = await _apiClient.put(
        '/admin/blog-tags/$id',
        headers: {'Content-Type': 'application/json'},
        body: {
          if (name != null) 'name': name,
          if (slug != null) 'slug': slug,
          if (description != null) 'description': description,
          if (color != null) 'color': color,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return BlogTagModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to update blog tag');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Blog tag not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update blog tag with status ${response.statusCode}');
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
  Future<void> deleteBlogTag(String id) async {
    try {
      final response = await _apiClient.delete('/admin/blog-tags/$id');

      if (response.statusCode != 200 && response.statusCode != 204) {
        if (response.statusCode == 404) {
          throw NotFoundException('Blog tag not found');
        } else {
          throw ServerException('Failed to delete blog tag with status ${response.statusCode}');
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

  // Blog Post Tags Relationship Methods
  @override
  Future<void> assignTagToPost({required String postId, required String tagId}) async {
    try {
      final response = await _apiClient.post(
        '/admin/blog/$postId/tags',
        body: {'tagId': tagId},
      );

      if (response.statusCode != 200 && response.statusCode != 201) {
        throw ServerException('Failed to assign tag to post with status ${response.statusCode}');
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
  Future<void> removeTagFromPost({required String postId, required String tagId}) async {
    try {
      final response = await _apiClient.delete(
        '/admin/blog/$postId/tags/$tagId',
      );

      if (response.statusCode != 200 && response.statusCode != 204) {
        throw ServerException('Failed to remove tag from post with status ${response.statusCode}');
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
  Future<List<BlogTagEntity>> getPostTags(String postId) async {
    try {
      final response = await _apiClient.get('/admin/blog/$postId/tags');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final tags = (data['data'] as List)
              .map((tagJson) => BlogTagModel.fromJson(tagJson).toEntity())
              .toList();
          return tags;
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch post tags');
        }
      } else {
        throw ServerException('Failed to fetch post tags with status ${response.statusCode}');
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
  Future<void> incrementViewCount(String postId) async {
    try {
      final response = await _apiClient.post(
        '/admin/blog/$postId/view',
      );

      if (response.statusCode != 200) {
        throw ServerException('Failed to increment view count with status ${response.statusCode}');
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
  Future<void> incrementLikeCount(String postId) async {
    try {
      final response = await _apiClient.post(
        '/admin/blog/$postId/like',
      );

      if (response.statusCode != 200) {
        throw ServerException('Failed to increment like count with status ${response.statusCode}');
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
  Future<void> updateBlogPostOrder(List<BlogPostOrderUpdateEntity> orders) async {
    try {
      final orderModels = orders.map((order) => BlogPostOrderUpdateModel.fromEntity(order)).toList();
      final response = await _apiClient.put(
        '/admin/blog/order',
        body: {'orders': orderModels.map((model) => model.toJson()).toList()},
      );

      if (response.statusCode != 200) {
        throw ServerException('Failed to update blog post order with status ${response.statusCode}');
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
  Future<void> updateBlogTagOrder(List<BlogTagOrderUpdateEntity> orders) async {
    try {
      final orderModels = orders.map((order) => BlogTagOrderUpdateModel.fromEntity(order)).toList();
      final response = await _apiClient.put(
        '/admin/blog-tags/order',
        body: {'orders': orderModels.map((model) => model.toJson()).toList()},
      );

      if (response.statusCode != 200) {
        throw ServerException('Failed to update blog tag order with status ${response.statusCode}');
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
  Future<BlogPostEntity> toggleBlogPostVisibility(String id) async {
    try {
      final response = await _apiClient.post(
        '/admin/blog/$id/toggle-visibility',
      );

      if (response.statusCode == 200) {
        final jsonData = json.decode(response.body);
        return BlogPostModel.fromJson(jsonData['data']).toEntity();
      } else {
        throw ServerException('Failed to toggle blog post visibility with status ${response.statusCode}');
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
  Future<BlogPostEntity> toggleBlogPostFeatured(String id) async {
    try {
      final response = await _apiClient.post(
        '/admin/blog/$id/toggle-featured',
      );

      if (response.statusCode == 200) {
        final jsonData = json.decode(response.body);
        return BlogPostModel.fromJson(jsonData['data']).toEntity();
      } else {
        throw ServerException('Failed to toggle blog post featured with status ${response.statusCode}');
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
  Future<BlogPostEntity> toggleBlogPostPublished(String id) async {
    try {
      final response = await _apiClient.post(
        '/admin/blog/$id/toggle-published',
      );

      if (response.statusCode == 200) {
        final jsonData = json.decode(response.body);
        return BlogPostModel.fromJson(jsonData['data']).toEntity();
      } else {
        throw ServerException('Failed to toggle blog post published with status ${response.statusCode}');
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
  Future<BlogStatsEntity> getBlogStats() async {
    try {
      final response = await _apiClient.get(
        '/admin/blog/stats',
      );

      if (response.statusCode == 200) {
        final jsonData = json.decode(response.body);
        return BlogStatsModel.fromJson(jsonData['data']).toEntity();
      } else {
        throw ServerException('Failed to get blog stats with status ${response.statusCode}');
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
  Future<BlogStatsEntity> getPublicBlogStats() async {
    try {
      final response = await _apiClient.get(
        '/blog/stats',
      );

      if (response.statusCode == 200) {
        final jsonData = json.decode(response.body);
        return BlogStatsModel.fromJson(jsonData['data']).toEntity();
      } else {
        throw ServerException('Failed to get public blog stats with status ${response.statusCode}');
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
  Future<List<BlogTagWithCountEntity>> getBlogTagsWithCount() async {
    try {
      final response = await _apiClient.get(
        '/admin/blog-tags/with-count',
      );

      if (response.statusCode == 200) {
        final jsonData = json.decode(response.body);
        final List<dynamic> tagsJson = jsonData['data'];
        return tagsJson.map((json) => BlogTagWithCountModel.fromJson(json).toEntity()).toList();
      } else {
        throw ServerException('Failed to get blog tags with count with status ${response.statusCode}');
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
  Future<List<BlogTagWithCountEntity>> getPopularBlogTags({int? limit}) async {
    try {
      final queryParams = limit != null ? {'limit': limit.toString()} : <String, String>{};
      final response = await _apiClient.get('/blog-tags/popular', queryParameters: queryParams);

      if (response.statusCode == 200) {
        final jsonData = json.decode(response.body);
        final List<dynamic> tagsJson = jsonData['data'];
        return tagsJson.map((json) => BlogTagWithCountModel.fromJson(json).toEntity()).toList();
      } else {
        throw ServerException('Failed to get popular blog tags with status ${response.statusCode}');
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
