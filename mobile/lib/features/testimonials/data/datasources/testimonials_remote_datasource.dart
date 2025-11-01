import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import '../../../../core/error/exceptions.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/testimonial_entity.dart';
import '../models/testimonial_model.dart';

abstract class TestimonialsRemoteDataSource {
  Future<List<TestimonialEntity>> getTestimonials({
    int? page,
    int? limit,
    bool? featured,
    bool? visible,
    bool? public,
    int? rating,
    String? search,
    String? sortBy,
    String? sortOrder,
  });

  Future<TestimonialEntity> getTestimonialById(String id);
  Future<TestimonialEntity> createTestimonial({
    required String name,
    required String position,
    required String company,
    required String content,
    String? avatar,
    required int rating,
    required bool featured,
    required int order,
    required bool visible,
    required bool public,
  });
  Future<TestimonialEntity> updateTestimonial({
    required String id,
    String? name,
    String? position,
    String? company,
    String? content,
    String? avatar,
    int? rating,
    bool? featured,
    int? order,
    bool? visible,
    bool? public,
  });
  Future<void> deleteTestimonial(String id);
  Future<void> updateTestimonialOrder(List<Map<String, dynamic>> testimonialOrders);
  Future<void> incrementViewCount(String testimonialId);
  Future<void> incrementLikeCount(String testimonialId);
}

class TestimonialsRemoteDataSourceImpl implements TestimonialsRemoteDataSource {
  final ApiClient _apiClient;
  final String _baseUrl;

  // Simple in-memory cache
  static final Map<String, dynamic> _cache = {};
  static final Map<String, DateTime> _cacheTimestamps = {};
  static const Duration _cacheDuration = Duration(minutes: 5);

  TestimonialsRemoteDataSourceImpl({
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
  Future<List<TestimonialEntity>> getTestimonials({
    int? page,
    int? limit,
    bool? featured,
    bool? visible,
    bool? public,
    int? rating,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) async {
    final cacheKey = 'testimonials_${page}_${limit}_${featured}_${visible}_${public}_${rating}_${search}_${sortBy}_${sortOrder}';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        final queryParams = <String, String>{};
        if (page != null) queryParams['page'] = page.toString();
        if (limit != null) queryParams['limit'] = limit.toString();
        if (featured != null) queryParams['featured'] = featured.toString();
        if (visible != null) queryParams['visible'] = visible.toString();
        if (public != null) queryParams['public'] = public.toString();
        if (rating != null) queryParams['rating'] = rating.toString();
        if (search != null && search.isNotEmpty) queryParams['search'] = search;
        if (sortBy != null) queryParams['sortBy'] = sortBy;
        if (sortOrder != null) queryParams['sortOrder'] = sortOrder;

        log('🔍 Testimonials API Request: /admin/testimonials');

        final response = await _apiClient.get('/admin/testimonials', queryParameters: queryParams);

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            final testimonials = (data['data']['testimonials'] as List)
                .map((testimonialJson) => TestimonialModel.fromApiJson(testimonialJson).toEntity())
                .toList();
            return testimonials;
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch testimonials');
          }
        } else {
          throw ServerException('Failed to fetch testimonials with status ${response.statusCode}');
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
  Future<TestimonialEntity> getTestimonialById(String id) async {
    final cacheKey = 'testimonial_$id';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 Testimonial by ID API Request: /admin/testimonials/$id');

        final response = await _apiClient.get('/admin/testimonials/$id');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            return TestimonialModel.fromApiJson(data['data']).toEntity();
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch testimonial');
          }
        } else if (response.statusCode == 404) {
          throw NotFoundException('Testimonial not found');
        } else {
          throw ServerException('Failed to fetch testimonial with status ${response.statusCode}');
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
  Future<TestimonialEntity> createTestimonial({
    required String name,
    required String position,
    required String company,
    required String content,
    String? avatar,
    required int rating,
    required bool featured,
    required int order,
    required bool visible,
    required bool public,
  }) async {
    try {
      final response = await _apiClient.post(
        '/admin/testimonials',
        headers: {'Content-Type': 'application/json'},
        body: {
          'name': name,
          'position': position,
          'company': company,
          'content': content,
          if (avatar != null) 'avatar': avatar,
          'rating': rating,
          'featured': featured,
          'order': order,
          'visible': visible,
          'public': public,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final testimonial = TestimonialModel.fromJson(data['data']).toEntity();
          
          // Invalidate testimonials cache
          _invalidateCache('testimonials_');
          
          return testimonial;
        } else {
          throw ServerException(data['message'] ?? 'Failed to create testimonial');
        }
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to create testimonial with status ${response.statusCode}');
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
  Future<TestimonialEntity> updateTestimonial({
    required String id,
    String? name,
    String? position,
    String? company,
    String? content,
    String? avatar,
    int? rating,
    bool? featured,
    int? order,
    bool? visible,
    bool? public,
  }) async {
    try {
      final response = await _apiClient.put(
        '/admin/testimonials/$id',
        headers: {'Content-Type': 'application/json'},
        body: {
          if (name != null) 'name': name,
          if (position != null) 'position': position,
          if (company != null) 'company': company,
          if (content != null) 'content': content,
          if (avatar != null) 'avatar': avatar,
          if (rating != null) 'rating': rating,
          if (featured != null) 'featured': featured,
          if (order != null) 'order': order,
          if (visible != null) 'visible': visible,
          if (public != null) 'public': public,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final testimonial = TestimonialModel.fromJson(data['data']).toEntity();
          
          // Invalidate testimonials cache
          _invalidateCache('testimonials_');
          _invalidateCache('testimonial_$id');
          
          return testimonial;
        } else {
          throw ServerException(data['message'] ?? 'Failed to update testimonial');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Testimonial not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update testimonial with status ${response.statusCode}');
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
  Future<void> deleteTestimonial(String id) async {
    try {
      final response = await _apiClient.delete('/admin/testimonials/$id');

      if (response.statusCode == 200 || response.statusCode == 204) {
        // Invalidate testimonials cache
        _invalidateCache('testimonials_');
        _invalidateCache('testimonial_$id');
      } else {
        if (response.statusCode == 404) {
          throw NotFoundException('Testimonial not found');
        } else {
          throw ServerException('Failed to delete testimonial with status ${response.statusCode}');
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
  Future<void> updateTestimonialOrder(List<Map<String, dynamic>> testimonialOrders) async {
    try {
      final response = await _apiClient.put(
        '/admin/testimonials/order',
        body: {'testimonialOrders': testimonialOrders},
      );

      if (response.statusCode != 200) {
        throw ServerException('Failed to update testimonial order with status ${response.statusCode}');
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
  Future<void> incrementViewCount(String testimonialId) async {
    try {
      final response = await _apiClient.post(
        '/admin/testimonials/$testimonialId/view',
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
  Future<void> incrementLikeCount(String testimonialId) async {
    try {
      final response = await _apiClient.post(
        '/admin/testimonials/$testimonialId/like',
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
}
