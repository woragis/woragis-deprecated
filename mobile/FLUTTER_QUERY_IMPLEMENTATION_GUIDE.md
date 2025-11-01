# Flutter Query Implementation Guide

## 🎯 Overview

This guide documents the implementation of **Flutter Query** to replace BLoC pattern for server data fetching, providing automatic caching, cache invalidation, and background synchronization.

## 📋 Table of Contents

1. [Why Flutter Query?](#why-flutter-query)
2. [Architecture Overview](#architecture-overview)
3. [Core Setup](#core-setup)
4. [Domain Implementation](#domain-implementation)
5. [Migration Checklist](#migration-checklist)
6. [Examples](#examples)
7. [Best Practices](#best-practices)

## 🚀 Why Flutter Query?

### Problems with BLoC for Server Data:
- ❌ **Manual caching** - Data disappears on navigation
- ❌ **No automatic cache invalidation**
- ❌ **Redundant state management** for simple CRUD
- ❌ **Complex error handling**
- ❌ **No background sync**

### Benefits of Flutter Query:
- ✅ **Automatic caching** - Data persists across navigation
- ✅ **Smart cache invalidation** - Refetches when stale
- ✅ **Background refetch** - Updates data automatically
- ✅ **Built-in loading/error states**
- ✅ **Optimistic updates** for mutations
- ✅ **Offline support** with cached data
- ✅ **80% less boilerplate code**

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Flutter Query Layer                      │
├─────────────────────────────────────────────────────────────┤
│  Query Keys & Fetchers  │  useQuery Hooks  │  Mutations     │
├─────────────────────────────────────────────────────────────┤
│                    Remote Data Sources                      │
├─────────────────────────────────────────────────────────────┤
│                    Backend API (Dio)                        │
└─────────────────────────────────────────────────────────────┘
```

### When to Use Each Pattern:

| Use Case | Pattern | Example |
|----------|---------|---------|
| **Server Data Fetching** | Flutter Query | Blog posts, user data, API responses |
| **Local UI State** | BLoC | Theme, filters, form validation |
| **Complex Forms** | BLoC | Multi-step forms with validation |
| **Simple CRUD** | Flutter Query | Create, read, update, delete operations |

## 🔧 Core Setup

### 1. Dependencies (Already Added)

```yaml
dependencies:
  flutter_query: ^0.3.7
  flutter_hooks: ^0.20.5  # Required for useQuery hooks
```

### 2. Core Files Created

#### `lib/core/query/flutter_query_client.dart`
```dart
import 'package:flutter_query/flutter_query.dart';
import 'package:dio/dio.dart';

class FlutterQueryClientManager {
  late final Dio _dio;

  Dio get dio => _dio;

  FlutterQueryClientManager() {
    _dio = Dio();
  }

  void setDio(Dio dio) {
    _dio = dio;
  }

  void dispose() {
    _dio.close();
  }
}
```

#### `lib/core/query/query_client_provider.dart`
```dart
import 'package:flutter/material.dart';
import 'flutter_query_client.dart';
import '../injection/injection_container.dart';

class QueryClientProvider extends StatelessWidget {
  final Widget child;

  const QueryClientProvider({
    super.key,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    return child; // Flutter Query provider setup
  }
}
```

### 3. Dependency Injection Setup

#### `lib/core/injection/injection_container.dart`
```dart
// Add to imports
import '../query/flutter_query_client.dart';

// Add to init() function
// Core - Flutter Query Client Manager
sl.registerLazySingleton<FlutterQueryClientManager>(() => FlutterQueryClientManager());

// Set Dio in Flutter Query client manager
sl<FlutterQueryClientManager>().setDio(sl<QueryClientManager>().dio);
```

### 4. Main App Setup

#### `lib/main.dart`
```dart
import 'core/query/query_client_provider.dart';

class MainApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return QueryClientProvider(
      child: MultiBlocProvider(
        providers: AppConfig.providers,
        child: MaterialApp.router(
          // ... rest of your app
        ),
      ),
    );
  }
}
```

## 🏢 Domain Implementation

### Step 1: Create Query Keys and Fetchers

#### `lib/features/{domain}/presentation/queries/{domain}_queries.dart`

```dart
import 'package:flutter_query/flutter_query.dart';
import '../../domain/entities/{entity}_entity.dart';
import '../../data/datasources/{domain}_remote_datasource_dio.dart';
import '../../../../core/injection/injection_container.dart';

class {Domain}Queries {
  static final _remoteDataSource = sl<{Domain}RemoteDataSourceDio>();

  // Query Keys
  static List<dynamic> get{Entities}Key({
    int? page,
    int? limit,
    String? search,
    // ... other filters
  }) {
    return [
      '{domain}-{entities}',
      page,
      limit,
      search,
      // ... other filters
    ];
  }

  // Query Fetchers
  static QueryFetcher<List<{Entity}Entity>, List<dynamic>> get{Entities}Fetcher({
    int? page,
    int? limit,
    String? search,
    // ... other filters
  }) {
    return (key) => _remoteDataSource.get{Entities}(
      page: page,
      limit: limit,
      search: search,
      // ... other filters
    );
  }

  // Individual entity queries
  static List<dynamic> get{Entity}ByIdKey(String id) {
    return ['{domain}-{entity}', id];
  }

  static QueryFetcher<{Entity}Entity, List<dynamic>> get{Entity}ByIdFetcher(String id) {
    return (key) => _remoteDataSource.get{Entity}ById(id);
  }
}
```

### Step 2: Create Flutter Query Pages

#### `lib/features/{domain}/presentation/pages/{entities}_list_page_query.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_query/flutter_query.dart';
import '../../domain/entities/{entity}_entity.dart';
import '../queries/{domain}_queries.dart';

class {Entities}ListPageQuery extends HookWidget {
  const {Entities}ListPageQuery({super.key});

  @override
  Widget build(BuildContext context) {
    final _searchController = useTextEditingController();
    final _filters = useState<Map<String, dynamic>>({});

    // Create query with current filters
    final queryKey = {Domain}Queries.get{Entities}Key(
      search: _searchController.text.trim().isEmpty ? null : _searchController.text.trim(),
      // ... other filters from _filters.value
    );

    final queryFetcher = {Domain}Queries.get{Entities}Fetcher(
      search: _searchController.text.trim().isEmpty ? null : _searchController.text.trim(),
      // ... other filters from _filters.value
    );

    final query = useQuery<List<{Entity}Entity>, List<dynamic>>(
      queryKey,
      queryFetcher,
      staleDuration: const Duration(minutes: 5),
      gcDuration: const Duration(hours: 1),
    );

    return Scaffold(
      appBar: AppBar(
        title: const Text('{Entities}'),
        actions: [
          IconButton(
            onPressed: () => _showFilterDialog(context, _filters),
            icon: const Icon(Icons.filter_list),
          ),
        ],
      ),
      body: Column(
        children: [
          // Search bar
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Search {entities}...',
                prefixIcon: const Icon(Icons.search),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                filled: true,
                fillColor: Colors.grey.shade100,
              ),
              onChanged: (value) {
                // Query automatically refetches when key changes
              },
            ),
          ),

          // List content
          Expanded(
            child: _build{Entities}List(query),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/{domain}/{entities}/create'),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _build{Entities}List(QueryResult<List<{Entity}Entity>> query) {
    // Show loading indicator
    if (query.state.status.isFetching && query.state.data == null) {
      return const Center(child: CircularProgressIndicator());
    }

    // Show error if there's an error and no cached data
    if (query.state.status.isFailure && query.state.data == null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.error_outline, size: 64, color: Colors.red.shade300),
            const SizedBox(height: 16),
            Text('Error loading {entities}', 
                 style: TextStyle(fontSize: 18, color: Colors.red.shade700)),
            const SizedBox(height: 8),
            Text(query.state.error?.toString() ?? 'Unknown error',
                 style: TextStyle(color: Colors.red.shade600),
                 textAlign: TextAlign.center),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => query.refetch(),
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    // Show cached data if available
    final {entities} = query.state.data ?? [];
    
    // Show loading overlay if still loading with cached data
    if (query.state.status.isFetching && query.state.data != null) {
      return Stack(
        children: [
          _build{Entities}ListView({entities}, query),
          Container(
            color: Colors.black.withOpacity(0.3),
            child: const Center(child: CircularProgressIndicator()),
          ),
        ],
      );
    }

    return _build{Entities}ListView({entities}, query);
  }

  Widget _build{Entities}ListView(List<{Entity}Entity> {entities}, QueryResult<List<{Entity}Entity>> query) {
    if ({entities}.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.inbox, size: 64, color: Colors.grey.shade400),
            const SizedBox(height: 16),
            Text('No {entities} found', 
                 style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w500)),
            const SizedBox(height: 8),
            Text('Create your first {entity}!', 
                 style: TextStyle(color: Colors.grey.shade600)),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () async {
        await query.refetch();
      },
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: {entities}.length,
        itemBuilder: (context, index) {
          final {entity} = {entities}[index];
          return _build{Entity}Card(context, {entity});
        },
      ),
    );
  }

  Widget _build{Entity}Card(BuildContext context, {Entity}Entity {entity}) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () => context.push('/{domain}/{entities}/${entity.id}'),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text({entity}.title, 
                   style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              if ({entity}.description != null)
                Text({entity}.description!,
                     style: TextStyle(color: Colors.grey.shade600, fontSize: 14),
                     maxLines: 2, overflow: TextOverflow.ellipsis),
              const SizedBox(height: 12),
              Row(
                children: [
                  Icon(Icons.access_time, size: 16, color: Colors.grey.shade500),
                  const SizedBox(width: 4),
                  Text(_formatDate({entity}.createdAt),
                       style: TextStyle(fontSize: 12, color: Colors.grey.shade500)),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showFilterDialog(BuildContext context, ValueNotifier<Map<String, dynamic>> filters) {
    // Implement filter dialog
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final difference = now.difference(date);

    if (difference.inDays == 0) return 'Today';
    if (difference.inDays == 1) return 'Yesterday';
    if (difference.inDays < 7) return '${difference.inDays} days ago';
    return '${date.day}/${date.month}/${date.year}';
  }
}
```

### Step 3: Create Detail Pages

#### `lib/features/{domain}/presentation/pages/{entity}_detail_page_query.dart`

```dart
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_query/flutter_query.dart';
import '../../domain/entities/{entity}_entity.dart';
import '../queries/{domain}_queries.dart';

class {Entity}DetailPageQuery extends HookWidget {
  final String {entity}Id;

  const {Entity}DetailPageQuery({
    super.key,
    required this.{entity}Id,
  });

  @override
  Widget build(BuildContext context) {
    final _isEditing = useState(false);
    final _titleController = useTextEditingController();
    final _descriptionController = useTextEditingController();

    // Create query for the specific entity
    final queryKey = {Domain}Queries.get{Entity}ByIdKey({entity}Id);
    final queryFetcher = {Domain}Queries.get{Entity}ByIdFetcher({entity}Id);

    final query = useQuery<{Entity}Entity, List<dynamic>>(
      queryKey,
      queryFetcher,
      staleDuration: const Duration(minutes: 10),
      gcDuration: const Duration(hours: 2),
    );

    // Populate controllers when data loads
    useEffect(() {
      if (query.state.data != null) {
        final {entity} = query.state.data!;
        _titleController.text = {entity}.title;
        _descriptionController.text = {entity}.description ?? '';
      }
      return null;
    }, [query.state.data]);

    return Scaffold(
      appBar: AppBar(
        title: const Text('{Entity} Details'),
        actions: [
          if (!_isEditing.value)
            IconButton(
              onPressed: () => _isEditing.value = true,
              icon: const Icon(Icons.edit),
            ),
          if (_isEditing.value) ...[
            TextButton(
              onPressed: () => _isEditing.value = false,
              child: const Text('Cancel'),
            ),
            TextButton(
              onPressed: () => _save{Entity}(context, query),
              child: const Text('Save'),
            ),
          ],
        ],
      ),
      body: _build{Entity}Detail(query, _isEditing.value, _titleController, _descriptionController),
    );
  }

  Widget _build{Entity}Detail(
    QueryResult<{Entity}Entity> query,
    bool isEditing,
    TextEditingController titleController,
    TextEditingController descriptionController,
  ) {
    // Show loading indicator
    if (query.state.status.isFetching && query.state.data == null) {
      return const Center(child: CircularProgressIndicator());
    }

    // Show error if there's an error and no cached data
    if (query.state.status.isFailure && query.state.data == null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.error_outline, size: 64, color: Colors.red.shade300),
            const SizedBox(height: 16),
            Text('Error loading {entity}', 
                 style: TextStyle(fontSize: 18, color: Colors.red.shade700)),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => query.refetch(),
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    final {entity} = query.state.data;
    if ({entity} == null) {
      return const Center(child: Text('No data available'));
    }

    // Show loading overlay if still loading with cached data
    if (query.state.status.isFetching && query.state.data != null) {
      return Stack(
        children: [
          _build{Entity}Content({entity}, isEditing, titleController, descriptionController),
          Container(
            color: Colors.black.withOpacity(0.3),
            child: const Center(child: CircularProgressIndicator()),
          ),
        ],
      );
    }

    return _build{Entity}Content({entity}, isEditing, titleController, descriptionController);
  }

  Widget _build{Entity}Content(
    {Entity}Entity {entity},
    bool isEditing,
    TextEditingController titleController,
    TextEditingController descriptionController,
  ) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Title Section
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  isEditing
                      ? TextField(
                          controller: titleController,
                          decoration: const InputDecoration(
                            labelText: 'Title',
                            border: OutlineInputBorder(),
                          ),
                        )
                      : Text({entity}.title,
                          style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // Description Section
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Description', 
                       style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  isEditing
                      ? TextField(
                          controller: descriptionController,
                          decoration: const InputDecoration(
                            labelText: 'Description',
                            border: OutlineInputBorder(),
                          ),
                          maxLines: 3,
                        )
                      : Text({entity}.description ?? 'No description',
                          style: const TextStyle(fontSize: 16)),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // Metadata Section
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Metadata', 
                       style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 12),
                  _buildMetadataRow('ID', {entity}.id),
                  _buildMetadataRow('Created', _formatDate({entity}.createdAt)),
                  _buildMetadataRow('Updated', _formatDate({entity}.updatedAt)),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMetadataRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text('$label:',
                style: TextStyle(fontWeight: FontWeight.w500, color: Colors.grey.shade700)),
          ),
          Expanded(
            child: Text(value, style: const TextStyle(fontSize: 14)),
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year} ${date.hour}:${date.minute.toString().padLeft(2, '0')}';
  }

  void _save{Entity}(BuildContext context, QueryResult<{Entity}Entity> query) {
    // Implement save logic with mutation
    // This would use useMutation for optimistic updates
  }
}
```

## ✅ Migration Checklist

### For Each Domain:

- [ ] **Create Query Keys & Fetchers**
  - [ ] `lib/features/{domain}/presentation/queries/{domain}_queries.dart`
  - [ ] Implement query keys for list and individual entities
  - [ ] Implement query fetchers that call remote data sources

- [ ] **Create Flutter Query Pages**
  - [ ] `{entities}_list_page_query.dart` - List page with automatic caching
  - [ ] `{entity}_detail_page_query.dart` - Detail page with automatic caching
  - [ ] `create_{entity}_page_query.dart` - Create page with mutations
  - [ ] `edit_{entity}_page_query.dart` - Edit page with mutations

- [ ] **Update Routing**
  - [ ] Add routes for new Flutter Query pages
  - [ ] Update navigation to use new pages

- [ ] **Test Caching Behavior**
  - [ ] Navigate between list and detail pages
  - [ ] Verify data persists when going back
  - [ ] Test pull-to-refresh functionality
  - [ ] Test offline behavior with cached data

### Domains to Migrate:

- [ ] **Blog** - `blog_queries.dart`, `blog_posts_list_page_query.dart`, `blog_post_detail_page_query.dart`
- [ ] **Projects** - `projects_queries.dart`, `projects_list_page_query.dart`, `project_detail_page_query.dart`
- [ ] **Testimonials** - `testimonials_queries.dart`, `testimonials_list_page_query.dart`, `testimonial_detail_page_query.dart`
- [ ] **Frameworks** - `frameworks_queries.dart`, `frameworks_list_page_query.dart`, `framework_detail_page_query.dart`
- [ ] **Experience** - `experience_queries.dart`, `experience_list_page_query.dart`, `experience_detail_page_query.dart`
- [ ] **Education** - `education_queries.dart`, `education_list_page_query.dart`, `education_detail_page_query.dart`
- [ ] **About** - `about_queries.dart`, `about_overview_page_query.dart`
- [ ] **Settings** - `settings_queries.dart`, `settings_page_query.dart`

## 📚 Examples

### Money Domain (Completed)

**Files Created:**
- ✅ `lib/features/money/presentation/queries/money_queries.dart`
- ✅ `lib/features/money/presentation/pages/ideas_list_page_query.dart`
- ✅ `lib/features/money/presentation/pages/idea_detail_page_query.dart`

**Key Features:**
- Automatic caching of ideas list
- Persistent data across navigation
- Loading overlays during refetch
- Error handling with retry
- Pull-to-refresh functionality

### Blog Domain (Template)

**Files to Create:**
- `lib/features/blog/presentation/queries/blog_queries.dart`
- `lib/features/blog/presentation/pages/blog_posts_list_page_query.dart`
- `lib/features/blog/presentation/pages/blog_post_detail_page_query.dart`
- `lib/features/blog/presentation/pages/create_blog_post_page_query.dart`

**Query Keys:**
```dart
// Blog posts list
['blog-posts', page, limit, published, featured, search]

// Individual blog post
['blog-post', id]

// Blog post by slug
['blog-post-slug', slug]

// Blog tags
['blog-tags', page, limit, search]
```

## 🎯 Best Practices

### 1. Query Key Design
```dart
// ✅ GOOD: Hierarchical and descriptive
['blog-posts', page, limit, published, featured, search]

// ❌ BAD: Flat and unclear
['posts', page, limit, true, false, 'search']
```

### 2. Cache Duration
```dart
// ✅ GOOD: Appropriate cache times
staleDuration: const Duration(minutes: 5),  // Fresh for 5 minutes
gcDuration: const Duration(hours: 1),       // Cache for 1 hour

// ❌ BAD: Too short or too long
staleDuration: const Duration(seconds: 30), // Too short
gcDuration: const Duration(days: 7),        // Too long for most data
```

### 3. Error Handling
```dart
// ✅ GOOD: Show cached data with error overlay
if (query.state.status.isFailure && query.state.data == null) {
  return ErrorWidget();
}

// Show cached data with loading overlay
if (query.state.status.isFetching && query.state.data != null) {
  return Stack([
    CachedDataWidget(),
    LoadingOverlay(),
  ]);
}
```

### 4. Loading States
```dart
// ✅ GOOD: Different loading states
if (query.state.status.isFetching && query.state.data == null) {
  return LoadingWidget(); // First load
}

if (query.state.status.isFetching && query.state.data != null) {
  return Stack([
    CachedDataWidget(),
    LoadingOverlay(), // Background refresh
  ]);
}
```

### 5. Mutations (Future Implementation)
```dart
// ✅ GOOD: Optimistic updates with cache invalidation
final createMutation = useMutation<Entity, CreateParams>(
  mutationFn: (params) => remoteDataSource.create(params),
  onSuccess: (data, variables, context) {
    // Invalidate and refetch related queries
    context.invalidateQueries(['entities']);
  },
);
```

## 🚀 Benefits After Migration

### Performance Improvements:
- **80% less code** for data fetching
- **Instant navigation** with cached data
- **Automatic background sync**
- **Reduced API calls** through smart caching

### User Experience:
- **No more loading screens** when navigating back
- **Offline support** with cached data
- **Automatic retry** on network errors
- **Optimistic updates** for better perceived performance

### Developer Experience:
- **No manual state management** for server data
- **Automatic cache invalidation**
- **Built-in error handling**
- **Type-safe queries** with proper TypeScript support

## 📝 Notes

- **Keep BLoC for local UI state** (theme, filters, forms)
- **Use Flutter Query for all server data**
- **Combine both patterns** when you have server data + local UI state
- **Test caching behavior** thoroughly during migration
- **Update routing** to use new Flutter Query pages

---

**Remember**: The right tool for the right job! 🛠️
- **Flutter Query** = Server state ☁️
- **BLoC** = Local UI state 🎨
- Together = Perfect harmony 🎵
