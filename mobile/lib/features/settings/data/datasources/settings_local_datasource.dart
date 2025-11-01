import 'package:sqflite/sqflite.dart';
import '../../../../core/database/database_helper.dart';
import '../../../../core/database/sync_manager.dart';
import '../../../../core/stores/auth_store.dart';
import '../../domain/entities/setting_entity.dart';
import '../models/setting_model.dart';

abstract class SettingsLocalDataSource {
  Future<List<SettingEntity>> getCachedSettings();
  Future<SettingEntity?> getCachedSetting(String id);
  Future<SettingEntity?> getCachedSettingByKey(String key);
  Future<List<SettingEntity>> getCachedSettingsByCategory(String category);
  Future<List<SettingEntity>> getCachedPublicSettings();
  Future<void> cacheSetting(SettingEntity setting);
  Future<void> cacheSettings(List<SettingEntity> settings);
  Future<void> updateCachedSetting(SettingEntity setting);
  Future<void> removeCachedSetting(String id);
  Future<void> removeCachedSettingByKey(String key);
}

class SettingsLocalDataSourceImpl implements SettingsLocalDataSource {
  final DatabaseHelper _dbHelper = DatabaseHelper();
  final SyncManager _syncManager = SyncManager();
  final AuthStoreBloc? _authStore;
  
  SettingsLocalDataSourceImpl({AuthStoreBloc? authStore}) : _authStore = authStore;

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
  Future<List<SettingEntity>> getCachedSettings() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'settings',
      orderBy: 'category ASC, key ASC',
    );
    return result.map((settingMap) => SettingModel.fromJson(settingMap)).toList();
  }

  @override
  Future<SettingEntity?> getCachedSetting(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'settings',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : SettingModel.fromJson(result.first);
  }

  @override
  Future<SettingEntity?> getCachedSettingByKey(String key) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'settings',
      where: 'key = ?',
      whereArgs: [key],
    );
    return result.isEmpty ? null : SettingModel.fromJson(result.first);
  }

  @override
  Future<List<SettingEntity>> getCachedSettingsByCategory(String category) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'settings',
      where: 'category = ?',
      whereArgs: [category],
      orderBy: 'key ASC',
    );
    return result.map((settingMap) => SettingModel.fromJson(settingMap)).toList();
  }

  @override
  Future<List<SettingEntity>> getCachedPublicSettings() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'settings',
      where: 'is_public = ?',
      whereArgs: [1],
      orderBy: 'category ASC, key ASC',
    );
    return result.map((settingMap) => SettingModel.fromJson(settingMap)).toList();
  }

  @override
  Future<void> cacheSetting(SettingEntity setting) async {
    final db = await _dbHelper.database;
    final settingMap = SettingModel.fromEntity(setting).toJson();
    settingMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    settingMap['is_dirty'] = 0;

    final safeSettingMap = _ensureRequiredFields(settingMap);

    await db.insert(
      'settings',
      safeSettingMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheSettings(List<SettingEntity> settings) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final setting in settings) {
      final settingMap = SettingModel.fromEntity(setting).toJson();
      settingMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      settingMap['is_dirty'] = 0;

      final safeSettingMap = _ensureRequiredFields(settingMap);

      batch.insert(
        'settings',
        safeSettingMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedSetting(SettingEntity setting) async {
    final db = await _dbHelper.database;
    final settingMap = SettingModel.fromEntity(setting).toJson();
    settingMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    settingMap['is_dirty'] = 1;

    final safeSettingMap = _ensureRequiredFields(settingMap);

    await db.update(
      'settings',
      safeSettingMap,
      where: 'id = ?',
      whereArgs: [setting.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'settings',
      recordId: setting.id,
      operation: SyncOperation.update,
      data: settingMap,
    );
  }

  @override
  Future<void> removeCachedSetting(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'settings',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'settings',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  @override
  Future<void> removeCachedSettingByKey(String key) async {
    final db = await _dbHelper.database;
    
    // First get the setting to get its ID for sync queue
    final setting = await getCachedSettingByKey(key);
    if (setting != null) {
      await _syncManager.addToSyncQueue(
        tableName: 'settings',
        recordId: setting.id,
        operation: SyncOperation.delete,
      );
    }

    await db.delete(
      'settings',
      where: 'key = ?',
      whereArgs: [key],
    );
  }
}
