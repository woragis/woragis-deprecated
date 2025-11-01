import 'package:sqflite/sqflite.dart';
import '../../../../core/database/database_helper.dart';
import '../../../../core/database/sync_manager.dart';
import '../../../../core/stores/auth_store.dart';
import '../../domain/entities/experience_entity.dart';
import '../models/experience_model.dart';

abstract class ExperienceLocalDataSource {
  Future<List<ExperienceEntity>> getCachedExperienceList();
  Future<ExperienceEntity?> getCachedExperience(String id);
  Future<void> cacheExperience(ExperienceEntity experience);
  Future<void> cacheExperienceList(List<ExperienceEntity> experienceList);
  Future<void> updateCachedExperience(ExperienceEntity experience);
  Future<void> removeCachedExperience(String id);
}

class ExperienceLocalDataSourceImpl implements ExperienceLocalDataSource {
  final DatabaseHelper _dbHelper = DatabaseHelper();
  final SyncManager _syncManager = SyncManager();
  final AuthStoreBloc? _authStore;
  
  ExperienceLocalDataSourceImpl({AuthStoreBloc? authStore}) : _authStore = authStore;

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
  Future<List<ExperienceEntity>> getCachedExperienceList() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'experiences',
      orderBy: '`order` ASC, company ASC',
    );
    return result.map((experienceMap) => ExperienceModel.fromLocalJson(experienceMap).toEntity()).toList();
  }

  @override
  Future<ExperienceEntity?> getCachedExperience(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'experiences',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : ExperienceModel.fromLocalJson(result.first).toEntity();
  }

  @override
  Future<void> cacheExperience(ExperienceEntity experience) async {
    final db = await _dbHelper.database;
    final experienceMap = ExperienceModel.fromEntity(experience).toLocalJson();
    experienceMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    experienceMap['is_dirty'] = 0;

    final safeExperienceMap = _ensureRequiredFields(experienceMap);

    await db.insert(
      'experiences',
      safeExperienceMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheExperienceList(List<ExperienceEntity> experienceList) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final experience in experienceList) {
      final experienceMap = ExperienceModel.fromEntity(experience).toLocalJson();
      experienceMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      experienceMap['is_dirty'] = 0;

      final safeExperienceMap = _ensureRequiredFields(experienceMap);

      batch.insert(
        'experiences',
        safeExperienceMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedExperience(ExperienceEntity experience) async {
    final db = await _dbHelper.database;
    final experienceMap = ExperienceModel.fromEntity(experience).toLocalJson();
    experienceMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    experienceMap['is_dirty'] = 1;

    final safeExperienceMap = _ensureRequiredFields(experienceMap);

    await db.update(
      'experiences',
      safeExperienceMap,
      where: 'id = ?',
      whereArgs: [experience.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'experiences',
      recordId: experience.id,
      operation: SyncOperation.update,
      data: experienceMap,
    );
  }

  @override
  Future<void> removeCachedExperience(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'experiences',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'experiences',
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
