# Flutter Mobile App - Architecture Documentation

## 🏗️ Clean Architecture with Simple HTTP Client + Caching

This Flutter application follows **Clean Architecture** principles with a **simple HTTP client** that includes built-in caching logic. We have completely removed Flutter Query and Dio dependencies, replacing them with a lightweight, efficient solution.

## 📐 Architecture Layers

```
lib/
├── core/                           # Core functionality
│   ├── database/                   # SQLite offline storage
│   ├── network/                    # Network connectivity
│   ├── error/                      # Error handling
│   ├── injection/                  # Dependency injection
│   └── presentation/
│       └── bloc/                   # Core BLoCs (local UI state)
│           ├── theme/              # Theme management
│           ├── language/           # Localization
│           └── navigation/         # Navigation state
│
└── features/                       # Feature modules
    └── [feature]/
        ├── domain/                 # Business logic layer
        │   ├── entities/           # Domain models
        │   ├── repositories/       # Abstract interfaces
        │   └── usecases/           # Business operations
        ├── data/                   # Data layer
        │   ├── models/             # Data models (JSON)
        │   ├── datasources/
        │   │   ├── local/          # SQLite (offline)
        │   │   └── remote/         # HTTP with caching (server)
        │   └── services/           # Service layer abstraction
        └── presentation/           # UI layer
            ├── pages/              # Screens (use Simple BLoC)
            ├── widgets/            # Reusable components
            └── bloc/               # Simple BLoCs for state management
```

## 🎯 **Key Architecture Decisions**

### **1. Simple HTTP Client with Built-in Caching** ⭐

**We use the standard `http` package with custom caching logic for:**
- ✅ Data fetching from API
- ✅ 5-minute in-memory caching
- ✅ Automatic cache invalidation
- ✅ Clear logging for debugging
- ✅ No heavy dependencies

**Example:**
```dart
class DomainRemoteDataSourceImpl implements DomainRemoteDataSource {
  final http.Client _client = http.Client();
  final String _baseUrl;

  // Simple in-memory cache
  static final Map<String, dynamic> _cache = {};
  static final Map<String, DateTime> _cacheTimestamps = {};
  static const Duration _cacheDuration = Duration(minutes: 5);

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

### **2. Service Layer for Clean Abstraction** 🎨

**Each domain has a service that wraps the remote data source:**

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

  Future<DomainEntity> createDomain({
    required String title,
    required String description,
    required bool visible,
  }) async {
    try {
      log('🔍 DomainService: Creating domain: $title');
      return await _remoteDataSource.createDomain(
        title: title,
        description: description,
        visible: visible,
      );
    } catch (e) {
      log('❌ DomainService: Error creating domain: $e');
      rethrow;
    }
  }
}
```

### **3. Simple BLoC for State Management** 🎯

**Each domain has a simple BLoC that uses the service:**

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

  Future<void> _onCreateDomain(
    CreateDomain event,
    Emitter<DomainState> emit,
  ) async {
    emit(DomainLoading());
    try {
      await _domainService.createDomain(
        title: event.title,
        description: event.description,
        visible: event.visible,
      );
      emit(DomainSuccess('Domain created successfully'));
    } catch (e) {
      emit(DomainError('Failed to create domain: ${e.toString()}'));
    }
  }
}
```

## 🔄 **Data Flow**

### **Server Data (Simple HTTP + Caching)**
```
Widget → Simple BLoC → Service → Remote DataSource → HTTP Client → API
                ↓
           Cache (5-minute in-memory)
                ↓
           SQLite (offline backup)
```

### **Local UI State (BLoC)**
```
User Action → Event → BLoC → State → Widget Update
                              ↓
                    SharedPreferences (persistence)
```

## 📦 **Current Implementation Status**

### **✅ Completed Domains:**
- **Money Domain** - Full migration with caching, service, simple BLoC, and pages
- **Projects Domain** - Full migration with caching, service, and simple BLoC  
- **Blog Domain** - Full migration with HTTP remote datasource and caching
- **About Domain** - Full migration with caching, service, and simple BLoC
- **Auth Domain** - Added caching to existing HTTP implementation
- **Education Domain** - Full migration with HTTP remote datasource and caching
- **Experience Domain** - Full migration with HTTP remote datasource, service, and simple BLoC
- **Frameworks Domain** - Full migration with HTTP remote datasource and caching

### **🔄 Remaining Domains:**
- **Settings Domain** - Ready for migration (same pattern)
- **Testimonials Domain** - Ready for migration (same pattern)

## 🚀 **State Management Strategy**

| State Type | Solution | Example |
|------------|----------|---------|
| Server Data | Simple HTTP + Caching | Blog posts, projects, user data |
| Local UI | Simple BLoC | Theme, language, filters |
| Complex Forms | Simple BLoC | Multi-step wizards, validation |
| Simple Forms | Local state (useState) | Search bar, simple inputs |
| Global App State | Simple BLoC | Navigation, theme, language |

## 💾 **Caching Strategy**

### **Cache Duration**
- **5 minutes** for all cached data
- Automatic expiration based on timestamps

### **Cache Keys**
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

### **Cache Invalidation Patterns**
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

## 🧪 **Testing Strategy**

### **Simple BLoCs**
```dart
blocTest<DomainSimpleBloc, DomainState>(
  'emits loaded state when domains are fetched successfully',
  build: () => DomainSimpleBloc(),
  act: (bloc) => bloc.add(LoadDomains()),
  expect: () => [
    DomainLoading(),
    DomainLoaded([mockDomain]),
  ],
);
```

### **Services**
```dart
test('should return domain list from remote data source', () async {
  // Arrange
  when(mockRemoteDataSource.getDomainList(any))
      .thenAnswer((_) async => [mockDomain]);
  
  // Act
  final result = await domainService.getDomainList();
  
  // Assert
  expect(result, [mockDomain]);
});
```

## 📚 **Key Packages**

### **State Management**
- `flutter_bloc: ^8.1.4` - Simple BLoC for state management
- `flutter_hooks: ^0.20.5` - React-like hooks for Flutter

### **Architecture**
- `get_it: ^7.6.7` - Dependency injection
- `injectable: ^2.3.2` - Code generation for DI
- `dartz: ^0.10.1` - Functional programming (Either, Option)

### **Data**
- `sqflite: ^2.3.0` - SQLite database
- `http: ^1.1.0` - Simple HTTP client
- `shared_preferences: ^2.2.2` - Key-value storage

### **Offline & Sync**
- `connectivity_plus: ^5.0.2` - Network connectivity
- Custom `SyncManager` for background sync
- Custom `DatabaseHelper` for SQLite management

## 🎓 **Best Practices**

### **DO ✅**
- Use simple HTTP client with built-in caching for all server data
- Use Simple BLoC for state management
- Keep BLoCs simple and focused
- Use service layer for clean abstraction
- Implement cache invalidation on mutations
- Handle errors gracefully
- Write tests for BLoCs and services
- Use clear logging for debugging

### **DON'T ❌**
- Don't use heavy dependencies like Flutter Query or Dio
- Don't manually manage complex caching logic
- Don't duplicate state between different layers
- Don't forget to dispose BLoCs properly
- Don't ignore linter warnings
- Don't forget to invalidate cache on mutations

## 🔧 **Development Workflow**

### **Adding a New Feature**

1. **Define Domain Models** (entities)
2. **Create Data Models** (with JSON serialization)
3. **Define Repository Interface**
4. **Implement Local DataSource** (SQLite)
5. **Implement Remote DataSource** (HTTP with caching)
6. **Create Service Layer** (abstraction)
7. **Create Simple BLoC** (state management)
8. **Build UI** (with BLoC)
9. **Register in DI Container**
10. **Test the implementation**

## 📖 **Examples**

### **Money Domain (Complete Example)**
- **Enhanced Remote Data Source** with 5-minute caching
- **Money Service** for clean abstraction
- **MoneySimpleBloc** for state management
- **New Pages**: MoneyHomePageSimple, IdeasListPageSimple, IdeaDetailPageSimple
- **Routes**: /money/simple, /money/ideas/simple, /money/ideas/simple/:id
- **Home Page Integration**: "Simple HTTP Demo" card

### **Cache Logging**
The pattern includes clear logging to help with debugging:

```
📦 Using cached data for: projects_1_10_true_true_false_
🌐 Fetching fresh data for: project_123
🗑️ Cache invalidated for pattern: projects_
```

## 🚀 **Getting Started**

1. **Install dependencies:**
   ```bash
   flutter pub get
   ```

2. **Generate code:**
   ```bash
   flutter pub run build_runner build
   ```

3. **Run the app:**
   ```bash
   flutter run
   ```

4. **Test the Simple HTTP Demo:**
   - Navigate to "Simple HTTP Demo" from the home page
   - Test caching behavior
   - Test cache invalidation

## 🤝 **Contributing**

When contributing, please:
1. Follow the simple HTTP client pattern
2. Use Simple BLoC for state management
3. Implement proper caching with invalidation
4. Write tests for your code
5. Update this documentation if needed

---

**Remember**: The goal is **simplicity** and **maintainability**. We've eliminated heavy dependencies and complex state management while providing excellent performance and user experience with our simple HTTP client and caching approach! 🎉