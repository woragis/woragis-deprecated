import 'package:sqflite/sqflite.dart';
import '../../../../core/database/database_helper.dart';
import '../../../../core/database/sync_manager.dart';
import '../../../../core/stores/auth_store.dart';
import '../../domain/entities/idea_entity.dart';
import '../../domain/entities/ai_chat_entity.dart';
import '../../domain/entities/idea_node_entity.dart';
import '../models/idea_model.dart';
import '../models/ai_chat_model.dart';
import '../models/idea_node_model.dart';

abstract class MoneyLocalDataSource {
  // Ideas methods
  Future<List<IdeaEntity>> getCachedIdeas();
  Future<IdeaEntity?> getCachedIdea(String id);
  Future<void> cacheIdea(IdeaEntity idea);
  Future<void> cacheIdeas(List<IdeaEntity> ideas);
  Future<void> updateCachedIdea(IdeaEntity idea);
  Future<void> removeCachedIdea(String id);

  // AI Chats methods
  Future<List<AiChatEntity>> getCachedAiChats();
  Future<AiChatEntity?> getCachedAiChat(String id);
  Future<void> cacheAiChat(AiChatEntity aiChat);
  Future<void> cacheAiChats(List<AiChatEntity> aiChats);
  Future<void> updateCachedAiChat(AiChatEntity aiChat);
  Future<void> removeCachedAiChat(String id);

  // Idea Nodes methods
  Future<List<IdeaNodeEntity>> getCachedIdeaNodes(String ideaId);
  Future<IdeaNodeEntity?> getCachedIdeaNode(String id);
  Future<void> cacheIdeaNode(IdeaNodeEntity ideaNode);
  Future<void> cacheIdeaNodes(List<IdeaNodeEntity> ideaNodes);
  Future<void> updateCachedIdeaNode(IdeaNodeEntity ideaNode);
  Future<void> removeCachedIdeaNode(String id);
}

class MoneyLocalDataSourceImpl implements MoneyLocalDataSource {
  final DatabaseHelper _dbHelper = DatabaseHelper();
  final SyncManager _syncManager = SyncManager();
  final AuthStoreBloc? _authStore;
  
  MoneyLocalDataSourceImpl({AuthStoreBloc? authStore}) : _authStore = authStore;

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

  // Ideas Implementation
  @override
  Future<List<IdeaEntity>> getCachedIdeas() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'ideas',
      orderBy: '`order` ASC, title ASC',
    );
    return result.map((ideaMap) => IdeaModel.fromLocalJson(ideaMap).toEntity()).toList();
  }

  @override
  Future<IdeaEntity?> getCachedIdea(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'ideas',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : IdeaModel.fromLocalJson(result.first).toEntity();
  }

  @override
  Future<void> cacheIdea(IdeaEntity idea) async {
    final db = await _dbHelper.database;
    final ideaMap = IdeaModel.fromEntity(idea).toLocalJson();
    ideaMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    ideaMap['is_dirty'] = 0;

    final safeIdeaMap = _ensureRequiredFields(ideaMap);

    await db.insert(
      'ideas',
      safeIdeaMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheIdeas(List<IdeaEntity> ideas) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final idea in ideas) {
      final ideaMap = IdeaModel.fromEntity(idea).toLocalJson();
      ideaMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      ideaMap['is_dirty'] = 0;

      final safeIdeaMap = _ensureRequiredFields(ideaMap);

      batch.insert(
        'ideas',
        safeIdeaMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedIdea(IdeaEntity idea) async {
    final db = await _dbHelper.database;
    final ideaMap = IdeaModel.fromEntity(idea).toLocalJson();
    ideaMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    ideaMap['is_dirty'] = 1;

    await db.update(
      'ideas',
      ideaMap,
      where: 'id = ?',
      whereArgs: [idea.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'ideas',
      recordId: idea.id,
      operation: SyncOperation.update,
      data: ideaMap,
    );
  }

  @override
  Future<void> removeCachedIdea(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'ideas',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'ideas',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  // AI Chats Implementation
  @override
  Future<List<AiChatEntity>> getCachedAiChats() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'ai_chats',
      orderBy: 'created_at DESC',
    );
    return result.map((chatMap) => AiChatModel.fromLocalJson(chatMap).toEntity()).toList();
  }

  @override
  Future<AiChatEntity?> getCachedAiChat(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'ai_chats',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : AiChatModel.fromLocalJson(result.first).toEntity();
  }

  @override
  Future<void> cacheAiChat(AiChatEntity aiChat) async {
    final db = await _dbHelper.database;
    final chatMap = AiChatModel.fromEntity(aiChat).toLocalJson();
    chatMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    chatMap['is_dirty'] = 0;

    await db.insert(
      'ai_chats',
      chatMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheAiChats(List<AiChatEntity> aiChats) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final aiChat in aiChats) {
      final chatMap = AiChatModel.fromEntity(aiChat).toLocalJson();
      chatMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      chatMap['is_dirty'] = 0;

      batch.insert(
        'ai_chats',
        chatMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedAiChat(AiChatEntity aiChat) async {
    final db = await _dbHelper.database;
    final chatMap = AiChatModel.fromEntity(aiChat).toLocalJson();
    chatMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    chatMap['is_dirty'] = 1;

    await db.update(
      'ai_chats',
      chatMap,
      where: 'id = ?',
      whereArgs: [aiChat.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'ai_chats',
      recordId: aiChat.id,
      operation: SyncOperation.update,
      data: chatMap,
    );
  }

  @override
  Future<void> removeCachedAiChat(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'ai_chats',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'ai_chats',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  // Idea Nodes Implementation
  @override
  Future<List<IdeaNodeEntity>> getCachedIdeaNodes(String ideaId) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'idea_nodes',
      where: 'idea_id = ?',
      whereArgs: [ideaId],
      orderBy: 'created_at ASC',
    );
    return result.map((nodeMap) => IdeaNodeModel.fromLocalJson(nodeMap).toEntity()).toList();
  }

  @override
  Future<IdeaNodeEntity?> getCachedIdeaNode(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'idea_nodes',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : IdeaNodeModel.fromLocalJson(result.first).toEntity();
  }

  @override
  Future<void> cacheIdeaNode(IdeaNodeEntity ideaNode) async {
    final db = await _dbHelper.database;
    final nodeMap = IdeaNodeModel.fromEntity(ideaNode).toLocalJson();
    nodeMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    nodeMap['is_dirty'] = 0;

    await db.insert(
      'idea_nodes',
      nodeMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheIdeaNodes(List<IdeaNodeEntity> ideaNodes) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final ideaNode in ideaNodes) {
      final nodeMap = IdeaNodeModel.fromEntity(ideaNode).toLocalJson();
      nodeMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      nodeMap['is_dirty'] = 0;

      batch.insert(
        'idea_nodes',
        nodeMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedIdeaNode(IdeaNodeEntity ideaNode) async {
    final db = await _dbHelper.database;
    final nodeMap = IdeaNodeModel.fromEntity(ideaNode).toLocalJson();
    nodeMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    nodeMap['is_dirty'] = 1;

    await db.update(
      'idea_nodes',
      nodeMap,
      where: 'id = ?',
      whereArgs: [ideaNode.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'idea_nodes',
      recordId: ideaNode.id,
      operation: SyncOperation.update,
      data: nodeMap,
    );
  }

  @override
  Future<void> removeCachedIdeaNode(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'idea_nodes',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'idea_nodes',
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
