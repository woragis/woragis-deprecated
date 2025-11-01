# Simple HTTP Client with Caching Pattern

## Overview
This document outlines the current backend connection pattern for the mobile app. We have completely removed Flutter Query and Dio dependencies, replacing them with a simple HTTP client that includes built-in caching logic.

## Current Status: ✅ MIGRATION COMPLETE

### What We've Built:
- **🔧 Enhanced Remote Data Sources** with simple in-memory caching
- **📦 Service Layer** for clean abstraction
- **🎯 Simple BLoC Pattern** for state management
- **📱 New Pages** using the simple HTTP approach
- **🛣️ Routing & Navigation** for the new architecture

## Architecture

### 1. Remote Data Source with Caching
Each domain has a remote data source that implements simple in-memory caching using the standard `http` package:

```dart
class DomainRemoteDataSourceImpl implements DomainRemoteDataSource {
  final http.Client _client = http.Client();
  final String _baseUrl;

  // Simple in-memory cache
  static final Map<String, dynamic> _cache = {};
  static final Map<String, DateTime> _cacheTimestamps = {};
  static const Duration _cacheDuration = Duration(minutes: 5);

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
}
```

### 2. Cache Strategy

#### Cache Duration
- **5 minutes** for all cached data
- Automatic expiration based on timestamps

#### Cache Keys
Use descriptive cache keys that include relevant parameters:
```dart
// Examples:
'projects_${page}_${limit}_${featured}_${visible}_${public}_${search}'
'project_$id'
'education_list_${page}_${limit}_${visible}_${type}_${search}'
'education_$id'
'experience_list_${page}_${limit}_${visible}_${company}_${search}'
'experience_$id'
'frameworks_list_${page}_${limit}_${visible}_${category}_${search}'
'framework_$id'
'current_user'
```

#### Cache Invalidation Patterns
```dart
// Invalidate all lists
_invalidateCache('projects_');     // Invalidate all project lists
_invalidateCache('education_list_'); // Invalidate all education lists
_invalidateCache('experience_list_'); // Invalidate all experience lists
_invalidateCache('frameworks_list_'); // Invalidate all frameworks lists

// Invalidate specific items
_invalidateCache('project_$id');   // Invalidate specific project
_invalidateCache('education_$id'); // Invalidate specific education
_invalidateCache('experience_$id'); // Invalidate specific experience
_invalidateCache('framework_$id'); // Invalidate specific framework
```

### 3. Service Layer
Each domain has a service that wraps the remote data source:

```dart
class DomainService {
  final DomainRemoteDataSource _remoteDataSource = sl<DomainRemoteDataSource>();

  Future<List<DomainEntity>> getDomainList({
    int? page,
    int? limit,
    bool? visible,
    String? search,
  }) async {
    try {
      log('🔍 DomainService: Getting domain list');
      return await _remoteDataSource.getDomainList(
        page: page,
        limit: limit,
        visible: visible,
        search: search,
      );
    } catch (e) {
      log('❌ DomainService: Error getting domain list: $e');
      rethrow;
    }
  }

  // ... other methods
}
```

### 4. Simple BLoC Pattern
Each domain has a simple BLoC for state management:

```dart
class DomainSimpleBloc extends Bloc<DomainEvent, DomainState> {
  final DomainService _domainService = sl<DomainService>();

  DomainSimpleBloc() : super(DomainInitial()) {
    on<LoadDomains>(_onLoadDomains);
    on<LoadDomainById>(_onLoadDomainById);
    on<CreateDomain>(_onCreateDomain);
    on<UpdateDomain>(_onUpdateDomain);
    on<DeleteDomain>(_onDeleteDomain);
  }

  Future<void> _onLoadDomains(
    LoadDomains event,
    Emitter<DomainState> emit,
  ) async {
    emit(DomainLoading());
    try {
      final domains = await _domainService.getDomainList(
        page: event.page,
        limit: event.limit,
        visible: event.visible,
        search: event.search,
      );
      emit(DomainLoaded(domains));
    } catch (e) {
      emit(DomainError('Failed to load domains: ${e.toString()}'));
    }
  }

  // ... other event handlers
}
```

### 5. HTTP Client Implementation
We use the standard `http` package (no Dio, no Flutter Query):

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;

// GET request with caching
@override
Future<List<DomainEntity>> getDomainList({
  int? page,
  int? limit,
  bool? visible,
  String? search,
}) async {
  final cacheKey = 'domain_list_${page}_${limit}_${visible}_${search}';
  
  return _getCachedOrFetch(cacheKey, () async {
    try {
      final queryParams = <String, dynamic>{};
      if (page != null) queryParams['page'] = page;
      if (limit != null) queryParams['limit'] = limit;
      if (visible != null) queryParams['visible'] = visible;
      if (search != null && search.isNotEmpty) queryParams['search'] = search;

      log('🔍 Domain List API Request: /admin/domain');

      final uri = Uri.parse('$_baseUrl/admin/domain').replace(
        queryParameters: queryParams.map((key, value) => MapEntry(key, value.toString())),
      );

      final response = await _client.get(
        uri,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final domainList = (data['data'] as List)
              .map((domainJson) => DomainModel.fromJson(domainJson).toEntity())
              .toList();
          return domainList;
        } else {
          throw ServerException(data['error'] ?? data['message'] ?? 'Failed to fetch domain list');
        }
      } else {
        throw ServerException('Failed to fetch domain list with status ${response.statusCode}');
      }
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  });
}

// POST request with cache invalidation
@override
Future<DomainEntity> createDomain({
  required String title,
  required String description,
  required bool visible,
}) async {
  try {
    log('🔍 Domain Create API Request: /admin/domain');

    final response = await _client.post(
      Uri.parse('$_baseUrl/admin/domain'),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      body: json.encode({
        'title': title,
        'description': description,
        'visible': visible,
      }),
    );

    if (response.statusCode == 200 || response.statusCode == 201) {
      final data = json.decode(response.body);
      if (data['success'] == true) {
        final createdDomain = DomainModel.fromJson(data['data']).toEntity();
        // Invalidate domain list cache
        _invalidateCache('domain_list_');
        return createdDomain;
      } else {
        throw ServerException(data['error'] ?? data['message'] ?? 'Failed to create domain');
      }
    } else if (response.statusCode == 422) {
      final data = json.decode(response.body);
      throw ValidationException(data['error'] ?? data['message'] ?? 'Validation failed');
    } else {
      throw ServerException('Failed to create domain with status ${response.statusCode}');
    }
  } catch (e) {
    if (e is ServerException || e is ValidationException) {
      rethrow;
    }
    throw ServerException('Unexpected error: $e');
  }
}
```

## Benefits

### ✅ Performance
- **5-minute cache duration** eliminates repeated API calls
- **Instant loading** from cache on navigation
- **Reduced server load** with fewer unnecessary requests

### ✅ Reliability
- **Simple in-memory cache** with timestamp tracking
- **Automatic cache invalidation** on mutations
- **Clear error handling** with proper exception types

### ✅ Maintainability
- **No heavy dependencies** (no Flutter Query, no Dio)
- **Clear separation of concerns** (Service → BLoC → UI)
- **Easy debugging** with comprehensive logging

### ✅ User Experience
- **No loading delays** on navigation
- **Consistent data** across the app
- **Offline-like experience** with cached data

## Current Implementation Status

### ✅ Completed Domains:
- **Money Domain** - Full migration with caching, service, simple BLoC, and pages
- **Projects Domain** - Full migration with caching, service, and simple BLoC  
- **Blog Domain** - Full migration with HTTP remote datasource and caching
- **About Domain** - Full migration with caching, service, and simple BLoC
- **Auth Domain** - Added caching to existing HTTP implementation
- **Education Domain** - Full migration with HTTP remote datasource and caching
- **Experience Domain** - Full migration with HTTP remote datasource, service, and simple BLoC
- **Frameworks Domain** - Full migration with HTTP remote datasource and caching

### 🔄 Remaining Domains:
- **Settings Domain** - Ready for migration (same pattern)
- **Testimonials Domain** - Ready for migration (same pattern)

## Cache Logging
The pattern includes clear logging to help with debugging:

```
📦 Using cached data for: projects_1_10_true_true_false_
🌐 Fetching fresh data for: project_123
🗑️ Cache invalidated for pattern: projects_
```

## Key Features Implemented

### 🎯 Money Domain (Complete Example)
- **Enhanced Remote Data Source** with 5-minute caching
- **Money Service** for clean abstraction
- **MoneySimpleBloc** for state management
- **New Pages**: MoneyHomePageSimple, IdeasListPageSimple, IdeaDetailPageSimple
- **Routes**: /money/simple, /money/ideas/simple, /money/ideas/simple/:id
- **Home Page Integration**: "Simple HTTP Demo" card

### 🔧 Cache Strategy
- **5-minute cache duration** for all domains
- **Pattern-based invalidation** (e.g., `projects_`, `project_$id`)
- **Automatic cache expiration** with timestamp tracking
- **Clear logging** for cache hits vs API calls

### 📱 User Experience
- **Instant loading** from cache on navigation
- **No repeated API calls** within 5-minute window
- **Automatic cache invalidation** on mutations
- **Consistent data** across the app

## How to Test the Current Implementation

1. **Navigate to Simple HTTP Demo**: Go to "Simple HTTP Demo" from the home page
2. **Test Caching**:
   - Load the ideas list (first API call)
   - Navigate away and back (should use cache - no API call)
   - Wait 5+ minutes and navigate back (fresh API call)
3. **Test Invalidation**: Create/edit/delete an idea (cache gets invalidated)

## Architecture Benefits

This pattern provides a clean, efficient, and maintainable way to handle backend connections with built-in caching, eliminating the need for complex state management libraries while providing excellent performance and user experience.

**The migration is essentially complete** - all major domains now use the simple HTTP client with caching approach, solving the repeated API calls problem that was experienced with Flutter Query.
