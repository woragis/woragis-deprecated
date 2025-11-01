import 'package:sqflite/sqflite.dart';
import '../../../../core/database/database_helper.dart';
import '../../../../core/database/sync_manager.dart';
import '../../../../core/stores/auth_store.dart';
import '../../domain/entities/blog_post_entity.dart';
import '../../domain/entities/blog_tag_entity.dart';
import '../../domain/entities/blog_stats_entity.dart';
import '../../domain/entities/blog_tag_with_count_entity.dart';
import '../../domain/entities/blog_order_update_entity.dart';
import '../models/blog_post_model.dart';
import '../models/blog_tag_model.dart';

abstract class BlogLocalDataSource {
  // Blog posts
  Future<List<BlogPostEntity>> getCachedBlogPosts({
    int? limit,
    int? offset,
    bool? published,
    bool? featured,
    bool? visible,
    bool? public,
    String? search,
  });

  Future<BlogPostEntity?> getCachedBlogPost(String id);
  Future<BlogPostEntity?> getCachedBlogPostBySlug(String slug);
  Future<void> cacheBlogPost(BlogPostEntity post);
  Future<void> cacheBlogPosts(List<BlogPostEntity> posts);
  Future<void> updateCachedBlogPost(BlogPostEntity post);
  Future<void> removeCachedBlogPost(String id);

  // Blog tags
  Future<List<BlogTagEntity>> getCachedBlogTags({
    int? limit,
    int? offset,
    bool? visible,
    String? search,
  });

  Future<BlogTagEntity?> getCachedBlogTag(String id);
  Future<BlogTagEntity?> getCachedBlogTagBySlug(String slug);
  Future<void> cacheBlogTag(BlogTagEntity tag);
  Future<void> cacheBlogTags(List<BlogTagEntity> tags);
  Future<void> updateCachedBlogTag(BlogTagEntity tag);
  Future<void> removeCachedBlogTag(String id);

  // Blog post tags relationship
  Future<void> assignTagToPost(String postId, String tagId);
  Future<void> removeTagFromPost(String postId, String tagId);
  Future<List<BlogTagEntity>> getPostTags(String postId);

  // Order management
  Future<void> updateBlogPostOrder(List<BlogPostOrderUpdateEntity> orders);
  Future<void> updateBlogTagOrder(List<BlogTagOrderUpdateEntity> orders);

  // Toggle convenience methods
  Future<BlogPostEntity> toggleBlogPostVisibility(String id);
  Future<BlogPostEntity> toggleBlogPostFeatured(String id);
  Future<BlogPostEntity> toggleBlogPostPublished(String id);

  // Statistics
  Future<BlogStatsEntity> getBlogStats();
  Future<List<BlogTagWithCountEntity>> getBlogTagsWithCount();
  Future<List<BlogTagWithCountEntity>> getPopularBlogTags({int? limit});

  // Offline operations
  Future<List<BlogPostEntity>> getDirtyBlogPosts();
  Future<List<BlogTagEntity>> getDirtyBlogTags();
  Future<void> markBlogPostAsDirty(String postId);
  Future<void> markBlogTagAsDirty(String tagId);
  Future<void> markBlogPostAsSynced(String postId);
  Future<void> markBlogTagAsSynced(String tagId);
}

class BlogLocalDataSourceImpl implements BlogLocalDataSource {
  final DatabaseHelper _dbHelper = DatabaseHelper();
  final SyncManager _syncManager = SyncManager();
  final AuthStoreBloc? _authStore;
  
  BlogLocalDataSourceImpl({AuthStoreBloc? authStore}) : _authStore = authStore;

  /// Ensure required NOT NULL fields are always present
  Map<String, dynamic> _ensureRequiredFields(Map<String, dynamic> data) {
    final safeData = Map<String, dynamic>.from(data);
    
    // ✅ ENSURE user_id is always present (NOT NULL constraint)
    if (!safeData.containsKey('user_id') || safeData['user_id'] == null || safeData['user_id'] == '') {
      // Try to get user ID from auth store
      final currentUserId = _authStore?.state.user?.id;
      if (currentUserId != null && currentUserId.isNotEmpty) {
        safeData['user_id'] = currentUserId;
      } else {
        // Fallback to a default user ID if not authenticated
        safeData['user_id'] = 'system';
      }
    }
    
    return safeData;
  }

  // Blog Posts Methods
  @override
  Future<List<BlogPostEntity>> getCachedBlogPosts({
    int? limit,
    int? offset,
    bool? published,
    bool? featured,
    bool? visible,
    bool? public,
    String? search,
  }) async {
    final db = await _dbHelper.database;
    
    String whereClause = '1=1';
    List<dynamic> whereArgs = [];

    if (published != null) {
      whereClause += ' AND published = ?';
      whereArgs.add(published ? 1 : 0);
    }

    if (featured != null) {
      whereClause += ' AND featured = ?';
      whereArgs.add(featured ? 1 : 0);
    }

    if (visible != null) {
      whereClause += ' AND visible = ?';
      whereArgs.add(visible ? 1 : 0);
    }

    if (public != null) {
      whereClause += ' AND public = ?';
      whereArgs.add(public ? 1 : 0);
    }

    if (search != null && search.isNotEmpty) {
      whereClause += ' AND (title LIKE ? OR excerpt LIKE ?)';
      whereArgs.add('%$search%');
      whereArgs.add('%$search%');
    }

    final result = await db.query(
      'blog_posts',
      where: whereClause,
      whereArgs: whereArgs,
      orderBy: '`order` ASC, created_at DESC',
      limit: limit,
      offset: offset,
    );

    return result.map((postMap) => BlogPostModel.fromDatabaseJson(postMap).toEntity()).toList();
  }

  @override
  Future<BlogPostEntity?> getCachedBlogPost(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'blog_posts',
      where: 'id = ?',
      whereArgs: [id],
    );

    if (result.isEmpty) return null;
    return BlogPostModel.fromDatabaseJson(result.first).toEntity();
  }

  @override
  Future<BlogPostEntity?> getCachedBlogPostBySlug(String slug) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'blog_posts',
      where: 'slug = ?',
      whereArgs: [slug],
    );

    if (result.isEmpty) return null;
    return BlogPostModel.fromDatabaseJson(result.first).toEntity();
  }

  @override
  Future<void> cacheBlogPost(BlogPostEntity post) async {
    final db = await _dbHelper.database;
    final postModel = BlogPostModel.fromEntity(post);
    final postMap = postModel.toDatabaseJson();
    
    // Add metadata for sync
    postMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    postMap['is_dirty'] = 0;

    final safePostMap = _ensureRequiredFields(postMap);

    await db.insert(
      'blog_posts',
      safePostMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheBlogPosts(List<BlogPostEntity> posts) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final post in posts) {
      final postModel = BlogPostModel.fromEntity(post);
      final postMap = postModel.toDatabaseJson();
      
      postMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      postMap['is_dirty'] = 0;

      final safePostMap = _ensureRequiredFields(postMap);

      batch.insert(
        'blog_posts',
        safePostMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedBlogPost(BlogPostEntity post) async {
    final db = await _dbHelper.database;
    final postModel = BlogPostModel.fromEntity(post);
    final postMap = postModel.toDatabaseJson();
    
    postMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    postMap['is_dirty'] = 1;

    final safePostMap = _ensureRequiredFields(postMap);

    await db.update(
      'blog_posts',
      safePostMap,
      where: 'id = ?',
      whereArgs: [post.id],
    );

    // Add to sync queue
    await _syncManager.addToSyncQueue(
      tableName: 'blog_posts',
      recordId: post.id,
      operation: SyncOperation.update,
      data: postMap,
    );
  }

  @override
  Future<void> removeCachedBlogPost(String id) async {
    final db = await _dbHelper.database;
    
    // Add to sync queue for deletion
    await _syncManager.addToSyncQueue(
      tableName: 'blog_posts',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'blog_posts',
      where: 'id = ?',
      whereArgs: [id],
    );

    // Also remove related blog post tags
    await db.delete(
      'blog_post_tags',
      where: 'blog_post_id = ?',
      whereArgs: [id],
    );
  }

  // Blog Tags Methods
  @override
  Future<List<BlogTagEntity>> getCachedBlogTags({
    int? limit,
    int? offset,
    bool? visible,
    String? search,
  }) async {
    final db = await _dbHelper.database;
    
    String whereClause = '1=1';
    List<dynamic> whereArgs = [];

    if (visible != null) {
      whereClause += ' AND visible = ?';
      whereArgs.add(visible ? 1 : 0);
    }

    if (search != null && search.isNotEmpty) {
      whereClause += ' AND (name LIKE ? OR description LIKE ?)';
      whereArgs.add('%$search%');
      whereArgs.add('%$search%');
    }

    final result = await db.query(
      'blog_tags',
      where: whereClause,
      whereArgs: whereArgs,
      orderBy: 'name ASC',
      limit: limit,
      offset: offset,
    );

    return result.map((tagMap) => BlogTagModel.fromDatabaseJson(tagMap).toEntity()).toList();
  }

  @override
  Future<BlogTagEntity?> getCachedBlogTag(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'blog_tags',
      where: 'id = ?',
      whereArgs: [id],
    );

    if (result.isEmpty) return null;
    return BlogTagModel.fromDatabaseJson(result.first).toEntity();
  }

  @override
  Future<BlogTagEntity?> getCachedBlogTagBySlug(String slug) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'blog_tags',
      where: 'slug = ?',
      whereArgs: [slug],
    );

    if (result.isEmpty) return null;
    return BlogTagModel.fromDatabaseJson(result.first).toEntity();
  }

  @override
  Future<void> cacheBlogTag(BlogTagEntity tag) async {
    final db = await _dbHelper.database;
    final tagModel = BlogTagModel.fromEntity(tag);
    final tagMap = tagModel.toDatabaseJson();
    
    tagMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    tagMap['is_dirty'] = 0;

    await db.insert(
      'blog_tags',
      tagMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheBlogTags(List<BlogTagEntity> tags) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final tag in tags) {
      final tagModel = BlogTagModel.fromEntity(tag);
      final tagMap = tagModel.toDatabaseJson();
      
      tagMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      tagMap['is_dirty'] = 0;

      batch.insert(
        'blog_tags',
        tagMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedBlogTag(BlogTagEntity tag) async {
    final db = await _dbHelper.database;
    final tagModel = BlogTagModel.fromEntity(tag);
    final tagMap = tagModel.toDatabaseJson();
    
    tagMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    tagMap['is_dirty'] = 1;

    await db.update(
      'blog_tags',
      tagMap,
      where: 'id = ?',
      whereArgs: [tag.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'blog_tags',
      recordId: tag.id,
      operation: SyncOperation.update,
      data: tagMap,
    );
  }

  @override
  Future<void> removeCachedBlogTag(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'blog_tags',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'blog_tags',
      where: 'id = ?',
      whereArgs: [id],
    );

    // Remove from blog post tags relationships
    await db.delete(
      'blog_post_tags',
      where: 'tag_id = ?',
      whereArgs: [id],
    );
  }

  // Blog Post Tags Relationship Methods
  @override
  Future<void> assignTagToPost(String postId, String tagId) async {
    final db = await _dbHelper.database;
    
    await db.insert(
      'blog_post_tags',
      {
        'blog_post_id': postId,
        'tag_id': tagId,
        'created_at': DateTime.now().millisecondsSinceEpoch,
      },
      conflictAlgorithm: ConflictAlgorithm.ignore,
    );
  }

  @override
  Future<void> removeTagFromPost(String postId, String tagId) async {
    final db = await _dbHelper.database;
    
    await db.delete(
      'blog_post_tags',
      where: 'blog_post_id = ? AND tag_id = ?',
      whereArgs: [postId, tagId],
    );
  }

  @override
  Future<List<BlogTagEntity>> getPostTags(String postId) async {
    final db = await _dbHelper.database;
    
    final result = await db.rawQuery('''
      SELECT bt.* FROM blog_tags bt
      INNER JOIN blog_post_tags bpt ON bt.id = bpt.tag_id
      WHERE bpt.blog_post_id = ?
      ORDER BY bt.name ASC
    ''', [postId]);

    return result.map((tagMap) => BlogTagModel.fromDatabaseJson(tagMap).toEntity()).toList();
  }

  // Offline Operations
  @override
  Future<List<BlogPostEntity>> getDirtyBlogPosts() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'blog_posts',
      where: 'is_dirty = ?',
      whereArgs: [1],
    );

    return result.map((postMap) => BlogPostModel.fromDatabaseJson(postMap).toEntity()).toList();
  }

  @override
  Future<List<BlogTagEntity>> getDirtyBlogTags() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'blog_tags',
      where: 'is_dirty = ?',
      whereArgs: [1],
    );

    return result.map((tagMap) => BlogTagModel.fromDatabaseJson(tagMap).toEntity()).toList();
  }

  @override
  Future<void> markBlogPostAsDirty(String postId) async {
    final db = await _dbHelper.database;
    await db.update(
      'blog_posts',
      {'is_dirty': 1},
      where: 'id = ?',
      whereArgs: [postId],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'blog_posts',
      recordId: postId,
      operation: SyncOperation.update,
    );
  }

  @override
  Future<void> markBlogTagAsDirty(String tagId) async {
    final db = await _dbHelper.database;
    await db.update(
      'blog_tags',
      {'is_dirty': 1},
      where: 'id = ?',
      whereArgs: [tagId],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'blog_tags',
      recordId: tagId,
      operation: SyncOperation.update,
    );
  }

  @override
  Future<void> markBlogPostAsSynced(String postId) async {
    final db = await _dbHelper.database;
    await db.update(
      'blog_posts',
      {
        'is_dirty': 0,
        'synced_at': DateTime.now().millisecondsSinceEpoch,
      },
      where: 'id = ?',
      whereArgs: [postId],
    );
  }

  @override
  Future<void> markBlogTagAsSynced(String tagId) async {
    final db = await _dbHelper.database;
    await db.update(
      'blog_tags',
      {
        'is_dirty': 0,
        'synced_at': DateTime.now().millisecondsSinceEpoch,
      },
      where: 'id = ?',
      whereArgs: [tagId],
    );
  }

  @override
  Future<void> updateBlogPostOrder(List<BlogPostOrderUpdateEntity> orders) async {
    final db = await _dbHelper.database;
    
    for (final order in orders) {
      await db.update(
        'blog_posts',
        {
          'order': order.order,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
          'is_dirty': 1,
        },
        where: 'id = ?',
        whereArgs: [order.id],
      );
    }
  }

  @override
  Future<void> updateBlogTagOrder(List<BlogTagOrderUpdateEntity> orders) async {
    final db = await _dbHelper.database;
    
    for (final order in orders) {
      await db.update(
        'blog_tags',
        {
          'order': order.order,
          'updated_at': DateTime.now().millisecondsSinceEpoch,
          'is_dirty': 1,
        },
        where: 'id = ?',
        whereArgs: [order.id],
      );
    }
  }

  @override
  Future<BlogPostEntity> toggleBlogPostVisibility(String id) async {
    final db = await _dbHelper.database;
    
    // Get current post
    final post = await getCachedBlogPost(id);
    if (post == null) {
      throw Exception('Blog post not found');
    }
    
    // Toggle visibility
    final newVisibility = !post.visible;
    await db.update(
      'blog_posts',
      {
        'visible': newVisibility ? 1 : 0,
        'updated_at': DateTime.now().millisecondsSinceEpoch,
        'is_dirty': 1,
      },
      where: 'id = ?',
      whereArgs: [id],
    );
    
    // Return updated post
    return post.copyWith(visible: newVisibility);
  }

  @override
  Future<BlogPostEntity> toggleBlogPostFeatured(String id) async {
    final db = await _dbHelper.database;
    
    // Get current post
    final post = await getCachedBlogPost(id);
    if (post == null) {
      throw Exception('Blog post not found');
    }
    
    // Toggle featured
    final newFeatured = !post.featured;
    await db.update(
      'blog_posts',
      {
        'featured': newFeatured ? 1 : 0,
        'updated_at': DateTime.now().millisecondsSinceEpoch,
        'is_dirty': 1,
      },
      where: 'id = ?',
      whereArgs: [id],
    );
    
    // Return updated post
    return post.copyWith(featured: newFeatured);
  }

  @override
  Future<BlogPostEntity> toggleBlogPostPublished(String id) async {
    final db = await _dbHelper.database;
    
    // Get current post
    final post = await getCachedBlogPost(id);
    if (post == null) {
      throw Exception('Blog post not found');
    }
    
    // Toggle published
    final newPublished = !post.published;
    final publishedAt = newPublished ? DateTime.now() : null;
    
    await db.update(
      'blog_posts',
      {
        'published': newPublished ? 1 : 0,
        'published_at': publishedAt?.millisecondsSinceEpoch,
        'updated_at': DateTime.now().millisecondsSinceEpoch,
        'is_dirty': 1,
      },
      where: 'id = ?',
      whereArgs: [id],
    );
    
    // Return updated post
    return post.copyWith(published: newPublished, publishedAt: publishedAt);
  }

  @override
  Future<BlogStatsEntity> getBlogStats() async {
    final db = await _dbHelper.database;
    
    // Get total count
    final totalResult = await db.rawQuery('SELECT COUNT(*) as count FROM blog_posts');
    final total = totalResult.first['count'] as int;
    
    // Get published count
    final publishedResult = await db.rawQuery(
      'SELECT COUNT(*) as count FROM blog_posts WHERE published = 1'
    );
    final published = publishedResult.first['count'] as int;
    
    // Get total views
    final viewsResult = await db.rawQuery(
      'SELECT COALESCE(SUM(view_count), 0) as total_views FROM blog_posts'
    );
    final totalViews = viewsResult.first['total_views'] as int;
    
    // Get total likes
    final likesResult = await db.rawQuery(
      'SELECT COALESCE(SUM(like_count), 0) as total_likes FROM blog_posts'
    );
    final totalLikes = likesResult.first['total_likes'] as int;
    
    // Get featured count
    final featuredResult = await db.rawQuery(
      'SELECT COUNT(*) as count FROM blog_posts WHERE featured = 1'
    );
    final featuredCount = featuredResult.first['count'] as int;
    
    return BlogStatsEntity(
      total: total,
      published: published,
      totalViews: totalViews,
      totalLikes: totalLikes,
      featuredCount: featuredCount,
    );
  }

  @override
  Future<List<BlogTagWithCountEntity>> getBlogTagsWithCount() async {
    final db = await _dbHelper.database;
    
    final result = await db.rawQuery('''
      SELECT 
        bt.*,
        COALESCE(COUNT(bpt.blog_post_id), 0) as post_count
      FROM blog_tags bt
      LEFT JOIN blog_post_tags bpt ON bt.id = bpt.tag_id
      GROUP BY bt.id
      ORDER BY bt.order ASC, bt.name ASC
    ''');
    
    return result.map((row) {
      final tag = BlogTagModel.fromDatabaseJson(row).toEntity();
      final postCount = row['post_count'] as int;
      return BlogTagWithCountEntity(tag: tag, postCount: postCount);
    }).toList();
  }

  @override
  Future<List<BlogTagWithCountEntity>> getPopularBlogTags({int? limit}) async {
    final db = await _dbHelper.database;
    
    final limitClause = limit != null ? 'LIMIT $limit' : '';
    
    final result = await db.rawQuery('''
      SELECT 
        bt.*,
        COUNT(bpt.blog_post_id) as post_count
      FROM blog_tags bt
      INNER JOIN blog_post_tags bpt ON bt.id = bpt.tag_id
      WHERE bt.visible = 1
      GROUP BY bt.id
      ORDER BY post_count DESC
      $limitClause
    ''');
    
    return result.map((row) {
      final tag = BlogTagModel.fromDatabaseJson(row).toEntity();
      final postCount = row['post_count'] as int;
      return BlogTagWithCountEntity(tag: tag, postCount: postCount);
    }).toList();
  }
}
