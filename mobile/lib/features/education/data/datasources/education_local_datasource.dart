import 'dart:convert';
import 'package:sqflite/sqflite.dart';
import '../../../../core/database/database_helper.dart';
import '../../../../core/stores/auth_store.dart';
import '../../domain/entities/education_entity.dart';
import '../models/education_model.dart';

abstract class EducationLocalDataSource {
  Future<List<EducationEntity>> getCachedEducationList();
  Future<EducationEntity?> getCachedEducation(String id);
  Future<void> cacheEducation(EducationEntity education);
  Future<void> cacheEducationList(List<EducationEntity> educationList);
  Future<void> updateCachedEducation(EducationEntity education);
  Future<void> removeCachedEducation(String id);
}

class EducationLocalDataSourceImpl implements EducationLocalDataSource {
  final DatabaseHelper _dbHelper = DatabaseHelper();
  final AuthStoreBloc? _authStore;
  
  EducationLocalDataSourceImpl({AuthStoreBloc? authStore}) : _authStore = authStore;

  @override
  Future<List<EducationEntity>> getCachedEducationList() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'education',
      orderBy: '`order` ASC, institution ASC',
    );
    return result.map((educationMap) => EducationModel.fromJson(_convertFromSqliteSafe(educationMap))).toList();
  }

  @override
  Future<EducationEntity?> getCachedEducation(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'education',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : EducationModel.fromJson(_convertFromSqliteSafe(result.first));
  }

  @override
  Future<void> cacheEducation(EducationEntity education) async {
    final db = await _dbHelper.database;
    final educationMap = _convertToSqliteSafe(EducationModel.fromEntity(education).toJson());

    await db.insert(
      'education',
      educationMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheEducationList(List<EducationEntity> educationList) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final education in educationList) {
      final educationMap = _convertToSqliteSafe(EducationModel.fromEntity(education).toJson());

      batch.insert(
        'education',
        educationMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  /// Convert data types to SQLite-safe types
  Map<String, dynamic> _convertToSqliteSafe(Map<String, dynamic> data) {
    final safeData = <String, dynamic>{};
    final now = DateTime.now();
    
    for (final entry in data.entries) {
      final key = entry.key;
      final value = entry.value;
      
      if (value == null) {
        safeData[key] = null;
      } else if (value is List<String>) {
        // Convert List<String> to JSON string
        safeData[key] = json.encode(value);
      } else if (value is bool) {
        // Convert bool to int (0 or 1)
        safeData[key] = value ? 1 : 0;
      } else if (value is DateTime) {
        // Convert DateTime to ISO string
        safeData[key] = value.toIso8601String();
      } else {
        // Keep other types as-is (String, int, etc.)
        safeData[key] = value;
      }
    }
    
    // ✅ ENSURE required NOT NULL fields are always present
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
    
    if (!safeData.containsKey('created_at') || safeData['created_at'] == null) {
      safeData['created_at'] = now.toIso8601String();
    }
    if (!safeData.containsKey('updated_at') || safeData['updated_at'] == null) {
      safeData['updated_at'] = now.toIso8601String();
    }
    
    return safeData;
  }

  /// Convert SQLite-safe types back to normal types
  Map<String, dynamic> _convertFromSqliteSafe(Map<String, dynamic> data) {
    final normalData = <String, dynamic>{};
    
    for (final entry in data.entries) {
      final key = entry.key;
      final value = entry.value;
      
      if (value == null) {
        normalData[key] = null;
      } else if (key == 'skills' && value is String) {
        // Convert JSON string back to List<String>
        try {
          final decoded = json.decode(value);
          if (decoded is List) {
            normalData[key] = decoded.map((e) => e.toString()).toList();
          } else {
            normalData[key] = null;
          }
        } catch (e) {
          normalData[key] = null;
        }
      } else if (key == 'visible' && value is int) {
        // Convert int back to bool
        normalData[key] = value == 1;
      } else if ((key == 'start_date' || key == 'end_date' || key == 'completion_date' || 
                  key == 'created_at' || key == 'updated_at') && value is String) {
        // Convert ISO string back to DateTime
        try {
          normalData[key] = DateTime.parse(value);
        } catch (e) {
          normalData[key] = null;
        }
      } else {
        // Keep other types as-is
        normalData[key] = value;
      }
    }
    
    return normalData;
  }

  @override
  Future<void> updateCachedEducation(EducationEntity education) async {
    final db = await _dbHelper.database;
    final educationMap = EducationModel.fromEntity(education).toJson();
    educationMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;

    final safeEducationMap = _convertToSqliteSafe(educationMap);

    await db.update(
      'education',
      safeEducationMap,
      where: 'id = ?',
      whereArgs: [education.id],
    );
  }

  @override
  Future<void> removeCachedEducation(String id) async {
    final db = await _dbHelper.database;
    
    await db.delete(
      'education',
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
