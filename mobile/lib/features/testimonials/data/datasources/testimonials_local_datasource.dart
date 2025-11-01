import 'package:sqflite/sqflite.dart';
import '../../../../core/database/database_helper.dart';
import '../../../../core/database/sync_manager.dart';
import '../../../../core/stores/auth_store.dart';
import '../../domain/entities/testimonial_entity.dart';
import '../models/testimonial_model.dart';

abstract class TestimonialsLocalDataSource {
  Future<List<TestimonialEntity>> getCachedTestimonials();
  Future<TestimonialEntity?> getCachedTestimonial(String id);
  Future<List<TestimonialEntity>> getCachedFeaturedTestimonials();
  Future<void> cacheTestimonial(TestimonialEntity testimonial);
  Future<void> cacheTestimonials(List<TestimonialEntity> testimonials);
  Future<void> updateCachedTestimonial(TestimonialEntity testimonial);
  Future<void> removeCachedTestimonial(String id);
}

class TestimonialsLocalDataSourceImpl implements TestimonialsLocalDataSource {
  final DatabaseHelper _dbHelper = DatabaseHelper();
  final SyncManager _syncManager = SyncManager();
  final AuthStoreBloc? _authStore;
  
  TestimonialsLocalDataSourceImpl({AuthStoreBloc? authStore}) : _authStore = authStore;

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

  @override
  Future<List<TestimonialEntity>> getCachedTestimonials() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'testimonials',
      orderBy: '`order` ASC, name ASC',
    );
    return result.map((testimonialMap) => TestimonialModel.fromDatabaseJson(testimonialMap).toEntity()).toList();
  }

  @override
  Future<TestimonialEntity?> getCachedTestimonial(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'testimonials',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : TestimonialModel.fromDatabaseJson(result.first).toEntity();
  }

  @override
  Future<List<TestimonialEntity>> getCachedFeaturedTestimonials() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'testimonials',
      where: 'featured = ? AND visible = ?',
      whereArgs: [1, 1],
      orderBy: '`order` ASC, name ASC',
    );
    return result.map((testimonialMap) => TestimonialModel.fromDatabaseJson(testimonialMap).toEntity()).toList();
  }

  @override
  Future<void> cacheTestimonial(TestimonialEntity testimonial) async {
    final db = await _dbHelper.database;
    final testimonialMap = TestimonialModel.fromEntity(testimonial).toDatabaseJson();
    testimonialMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    testimonialMap['is_dirty'] = 0;

    final safeTestimonialMap = _ensureRequiredFields(testimonialMap);

    await db.insert(
      'testimonials',
      safeTestimonialMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheTestimonials(List<TestimonialEntity> testimonials) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final testimonial in testimonials) {
      final testimonialMap = TestimonialModel.fromEntity(testimonial).toDatabaseJson();
      testimonialMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      testimonialMap['is_dirty'] = 0;

      final safeTestimonialMap = _ensureRequiredFields(testimonialMap);

      batch.insert(
        'testimonials',
        safeTestimonialMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedTestimonial(TestimonialEntity testimonial) async {
    final db = await _dbHelper.database;
    final testimonialMap = TestimonialModel.fromEntity(testimonial).toDatabaseJson();
    testimonialMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    testimonialMap['is_dirty'] = 1;

    await db.update(
      'testimonials',
      testimonialMap,
      where: 'id = ?',
      whereArgs: [testimonial.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'testimonials',
      recordId: testimonial.id,
      operation: SyncOperation.update,
      data: testimonialMap,
    );
  }

  @override
  Future<void> removeCachedTestimonial(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'testimonials',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'testimonials',
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
