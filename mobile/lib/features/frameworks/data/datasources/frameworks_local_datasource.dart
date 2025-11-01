import 'package:sqflite/sqflite.dart';
import '../../../../core/database/database_helper.dart';
import '../../../../core/database/sync_manager.dart';
import '../../../../core/stores/auth_store.dart';
import '../../domain/entities/framework_entity.dart';
import '../models/framework_model.dart';

abstract class FrameworksLocalDataSource {
  Future<List<FrameworkEntity>> getCachedFrameworks({
    int? limit,
    int? offset,
    bool? visible,
    bool? public,
    String? type,
    String? search,
  });

  Future<FrameworkEntity?> getCachedFramework(String id);
  Future<FrameworkEntity?> getCachedFrameworkBySlug(String slug);
  Future<void> cacheFramework(FrameworkEntity framework);
  Future<void> cacheFrameworks(List<FrameworkEntity> frameworks);
  Future<void> updateCachedFramework(FrameworkEntity framework);
  Future<void> removeCachedFramework(String id);

  // Offline operations
  Future<List<FrameworkEntity>> getDirtyFrameworks();
  Future<void> markFrameworkAsDirty(String frameworkId);
  Future<void> markFrameworkAsSynced(String frameworkId);
}

class FrameworksLocalDataSourceImpl implements FrameworksLocalDataSource {
  final DatabaseHelper _dbHelper = DatabaseHelper();
  final SyncManager _syncManager = SyncManager();
  final AuthStoreBloc _authStoreBloc;
  
  FrameworksLocalDataSourceImpl({required AuthStoreBloc authStoreBloc}) 
      : _authStoreBloc = authStoreBloc;

  String? get _currentUserId => _authStoreBloc.state.user?.id;

  @override
  Future<List<FrameworkEntity>> getCachedFrameworks({
    int? limit,
    int? offset,
    bool? visible,
    bool? public,
    String? type,
    String? search,
  }) async {
    final db = await _dbHelper.database;
    
    String whereClause = '1=1';
    List<dynamic> whereArgs = [];

    if (visible != null) {
      whereClause += ' AND visible = ?';
      whereArgs.add(visible ? 1 : 0);
    }

    if (public != null) {
      whereClause += ' AND public = ?';
      whereArgs.add(public ? 1 : 0);
    }

    if (type != null && type.isNotEmpty) {
      whereClause += ' AND type = ?';
      whereArgs.add(type);
    }

    if (search != null && search.isNotEmpty) {
      whereClause += ' AND (name LIKE ? OR description LIKE ?)';
      whereArgs.add('%$search%');
      whereArgs.add('%$search%');
    }

    final result = await db.query(
      'frameworks',
      where: whereClause,
      whereArgs: whereArgs,
      orderBy: '`order` ASC, name ASC',
      limit: limit,
      offset: offset,
    );

    return result.map((frameworkMap) => FrameworkModel.fromJson(frameworkMap, userId: _currentUserId).toEntity()).toList();
  }

  @override
  Future<FrameworkEntity?> getCachedFramework(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'frameworks',
      where: 'id = ?',
      whereArgs: [id],
    );

    if (result.isEmpty) return null;
    return FrameworkModel.fromJson(result.first).toEntity();
  }

  @override
  Future<FrameworkEntity?> getCachedFrameworkBySlug(String slug) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'frameworks',
      where: 'slug = ?',
      whereArgs: [slug],
    );

    if (result.isEmpty) return null;
    return FrameworkModel.fromJson(result.first).toEntity();
  }

  @override
  Future<void> cacheFramework(FrameworkEntity framework) async {
    final db = await _dbHelper.database;
    final frameworkModel = FrameworkModel.fromEntity(framework);
    final frameworkMap = frameworkModel.toJson();
    
    // Ensure userId is set from auth store
    if (_currentUserId != null) {
      frameworkMap['userId'] = _currentUserId;
    }
    
    frameworkMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    frameworkMap['is_dirty'] = 0;

    await db.insert(
      'frameworks',
      frameworkMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheFrameworks(List<FrameworkEntity> frameworks) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final framework in frameworks) {
      final frameworkModel = FrameworkModel.fromEntity(framework);
      final frameworkMap = frameworkModel.toJson();
      
      // Ensure userId is set from auth store
      if (_currentUserId != null) {
        frameworkMap['userId'] = _currentUserId;
      }
      
      frameworkMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      frameworkMap['is_dirty'] = 0;

      batch.insert(
        'frameworks',
        frameworkMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedFramework(FrameworkEntity framework) async {
    final db = await _dbHelper.database;
    final frameworkModel = FrameworkModel.fromEntity(framework);
    final frameworkMap = frameworkModel.toJson();
    
    // Ensure userId is set from auth store
    if (_currentUserId != null) {
      frameworkMap['userId'] = _currentUserId;
    }
    
    frameworkMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    frameworkMap['is_dirty'] = 1;

    await db.update(
      'frameworks',
      frameworkMap,
      where: 'id = ?',
      whereArgs: [framework.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'frameworks',
      recordId: framework.id,
      operation: SyncOperation.update,
      data: frameworkMap,
    );
  }

  @override
  Future<void> removeCachedFramework(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'frameworks',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'frameworks',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  @override
  Future<List<FrameworkEntity>> getDirtyFrameworks() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'frameworks',
      where: 'is_dirty = ?',
      whereArgs: [1],
    );

    return result.map((frameworkMap) => FrameworkModel.fromJson(frameworkMap, userId: _currentUserId).toEntity()).toList();
  }

  @override
  Future<void> markFrameworkAsDirty(String frameworkId) async {
    final db = await _dbHelper.database;
    await db.update(
      'frameworks',
      {'is_dirty': 1},
      where: 'id = ?',
      whereArgs: [frameworkId],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'frameworks',
      recordId: frameworkId,
      operation: SyncOperation.update,
    );
  }

  @override
  Future<void> markFrameworkAsSynced(String frameworkId) async {
    final db = await _dbHelper.database;
    await db.update(
      'frameworks',
      {
        'is_dirty': 0,
        'synced_at': DateTime.now().millisecondsSinceEpoch,
      },
      where: 'id = ?',
      whereArgs: [frameworkId],
    );
  }
}
