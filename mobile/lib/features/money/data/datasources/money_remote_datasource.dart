import 'dart:convert';
import 'dart:developer';
import 'package:http/http.dart' as http;
import '../../../../core/error/exceptions.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/idea_entity.dart';
import '../../domain/entities/ai_chat_entity.dart';
import '../../domain/entities/idea_node_entity.dart';
import '../models/idea_model.dart';
import '../models/ai_chat_model.dart';
import '../models/idea_node_model.dart';

abstract class MoneyRemoteDataSource {
  // Ideas
  Future<List<IdeaEntity>> getIdeas({
    int? page,
    int? limit,
    bool? featured,
    bool? visible,
    bool? public,
    String? search,
    String? sortBy,
    String? sortOrder,
  });

  Future<IdeaEntity> getIdeaById(String id);
  Future<IdeaEntity> getIdeaBySlug(String slug);
  Future<IdeaEntity> createIdea({
    required String title,
    required String slug,
    required String document,
    String? description,
    required bool featured,
    required bool visible,
    required bool public,
    required int order,
  });
  Future<IdeaEntity> updateIdea({
    required String id,
    String? title,
    String? slug,
    String? document,
    String? description,
    bool? featured,
    bool? visible,
    bool? public,
    int? order,
  });
  Future<void> deleteIdea(String id);

  // AI Chats
  Future<List<AiChatEntity>> getAiChats({
    int? page,
    int? limit,
    String? ideaNodeId,
    String? agent,
    bool? archived,
    bool? visible,
    String? search,
    String? sortBy,
    String? sortOrder,
  });

  Future<AiChatEntity> getAiChatById(String id);
  Future<AiChatEntity> createAiChat({
    required String title,
    String? ideaNodeId,
    required String agent,
    String? model,
    String? systemPrompt,
    required double temperature,
    int? maxTokens,
    required bool visible,
    required bool archived,
  });
  Future<AiChatEntity> updateAiChat({
    required String id,
    String? title,
    String? ideaNodeId,
    String? agent,
    String? model,
    String? systemPrompt,
    double? temperature,
    int? maxTokens,
    bool? visible,
    bool? archived,
  });
  Future<void> deleteAiChat(String id);
  Future<AiChatEntity> sendMessage({
    required String chatId,
    required String message,
  });
  Future<List<ChatMessageEntity>> getChatMessages(String chatId);

  // Idea Nodes
  Future<List<IdeaNodeEntity>> getIdeaNodes({
    required String ideaId,
    String? type,
    bool? visible,
    int? limit,
    int? offset,
  });

  Future<IdeaNodeEntity> getIdeaNodeById(String id);
  Future<IdeaNodeEntity> createIdeaNode({
    required String ideaId,
    required String title,
    String? content,
    required String type,
    required double positionX,
    required double positionY,
    double? width,
    double? height,
    String? color,
    required List<String> connections,
    required bool visible,
  });
  Future<IdeaNodeEntity> updateIdeaNode({
    required String id,
    String? title,
    String? content,
    String? type,
    double? positionX,
    double? positionY,
    double? width,
    double? height,
    String? color,
    List<String>? connections,
    bool? visible,
  });
  Future<void> deleteIdeaNode(String id);
  Future<IdeaNodeEntity> updateNodePosition({
    required String id,
    required double positionX,
    required double positionY,
  });
  Future<void> updateNodePositions(List<Map<String, dynamic>> positions);
  Future<IdeaNodeEntity> updateNodeConnections({
    required String id,
    required List<String> connections,
  });
}

class MoneyRemoteDataSourceImpl implements MoneyRemoteDataSource {
  final ApiClient _apiClient;
  final String _baseUrl;

  // Simple in-memory cache
  static final Map<String, dynamic> _cache = {};
  static final Map<String, DateTime> _cacheTimestamps = {};
  static const Duration _cacheDuration = Duration(minutes: 5);

  MoneyRemoteDataSourceImpl({
    required ApiClient apiClient,
    required String baseUrl,
  })  : _apiClient = apiClient,
        _baseUrl = baseUrl;

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

  // Ideas Implementation
  @override
  Future<List<IdeaEntity>> getIdeas({
    int? page,
    int? limit,
    bool? featured,
    bool? visible,
    bool? public,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) async {
    final cacheKey = 'ideas_${page}_${limit}_${featured}_${visible}_${public}_${search}_${sortBy}_${sortOrder}';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        final queryParams = <String, String>{};
        if (page != null) queryParams['page'] = page.toString();
        if (limit != null) queryParams['limit'] = limit.toString();
        if (featured != null) queryParams['featured'] = featured.toString();
        if (visible != null) queryParams['visible'] = visible.toString();
        if (public != null) queryParams['public'] = public.toString();
        if (search != null && search.isNotEmpty) queryParams['search'] = search;
        if (sortBy != null) queryParams['sortBy'] = sortBy;
        if (sortOrder != null) queryParams['sortOrder'] = sortOrder;

        log('🔍 Ideas List API Request: /api/money/ideas');

        final response = await _apiClient.get('/api/money/ideas', queryParameters: queryParams);

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            final ideas = (data['data'] as List)
                .map((ideaJson) => IdeaModel.fromApiJson(ideaJson).toEntity())
                .toList();
            return ideas;
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch ideas');
          }
        } else {
          throw ServerException('Failed to fetch ideas with status ${response.statusCode}');
        }
      } on http.ClientException {
        throw NetworkException('Network error occurred');
      } catch (e) {
        if (e is ServerException || e is NetworkException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }

  @override
  Future<IdeaEntity> getIdeaById(String id) async {
    final cacheKey = 'idea_$id';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 Idea by ID API Request: /api/money/ideas/$id');

        final response = await _apiClient.get('/api/money/ideas/$id');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            return IdeaModel.fromApiJson(data['data']).toEntity();
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch idea');
          }
        } else if (response.statusCode == 404) {
          throw NotFoundException('Idea not found');
        } else {
          throw ServerException('Failed to fetch idea with status ${response.statusCode}');
        }
      } on http.ClientException {
        throw NetworkException('Network error occurred');
      } catch (e) {
        if (e is ServerException || e is NetworkException || e is NotFoundException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }

  @override
  Future<IdeaEntity> getIdeaBySlug(String slug) async {
    try {
      final response = await _apiClient.get('/api/money/ideas/slug/$slug');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return IdeaModel.fromApiJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch idea');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Idea not found');
      } else {
        throw ServerException('Failed to fetch idea with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<IdeaEntity> createIdea({
    required String title,
    required String slug,
    required String document,
    String? description,
    required bool featured,
    required bool visible,
    required bool public,
    required int order,
  }) async {
    try {
      final response = await _apiClient.post(
        '/api/money/ideas',
        headers: {'Content-Type': 'application/json'},
        body: {
          'title': title,
          'slug': slug,
          'document': document,
          if (description != null) 'description': description,
          'featured': featured,
          'visible': visible,
          'public': public,
          'order': order,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final idea = IdeaModel.fromApiJson(data['data']).toEntity();
          
          // Invalidate ideas cache
          _invalidateCache('ideas_');
          
          return idea;
        } else {
          throw ServerException(data['message'] ?? 'Failed to create idea');
        }
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to create idea with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<IdeaEntity> updateIdea({
    required String id,
    String? title,
    String? slug,
    String? document,
    String? description,
    bool? featured,
    bool? visible,
    bool? public,
    int? order,
  }) async {
    try {
      final response = await _apiClient.put(
        '/api/money/ideas/$id',
        headers: {'Content-Type': 'application/json'},
        body: {
          if (title != null) 'title': title,
          if (slug != null) 'slug': slug,
          if (document != null) 'document': document,
          if (description != null) 'description': description,
          if (featured != null) 'featured': featured,
          if (visible != null) 'visible': visible,
          if (public != null) 'public': public,
          if (order != null) 'order': order,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final idea = IdeaModel.fromApiJson(data['data']).toEntity();
          
          // Invalidate ideas cache
          _invalidateCache('ideas_');
          _invalidateCache('idea_$id');
          
          return idea;
        } else {
          throw ServerException(data['message'] ?? 'Failed to update idea');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Idea not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update idea with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> deleteIdea(String id) async {
    try {
      final response = await _apiClient.delete('/api/money/ideas/$id');

      if (response.statusCode == 200 || response.statusCode == 204) {
        // Invalidate ideas cache
        _invalidateCache('ideas_');
        _invalidateCache('idea_$id');
      } else {
        if (response.statusCode == 404) {
          throw NotFoundException('Idea not found');
        } else {
          throw ServerException('Failed to delete idea with status ${response.statusCode}');
        }
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  // AI Chats Implementation
  @override
  Future<List<AiChatEntity>> getAiChats({
    int? page,
    int? limit,
    String? ideaNodeId,
    String? agent,
    bool? archived,
    bool? visible,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) async {
    try {
      final queryParams = <String, String>{};
      if (page != null) queryParams['page'] = page.toString();
      if (limit != null) queryParams['limit'] = limit.toString();
      if (ideaNodeId != null) queryParams['ideaNodeId'] = ideaNodeId;
      if (agent != null) queryParams['agent'] = agent;
      if (archived != null) queryParams['archived'] = archived.toString();
      if (visible != null) queryParams['visible'] = visible.toString();
      if (search != null && search.isNotEmpty) queryParams['search'] = search;
      if (sortBy != null) queryParams['sortBy'] = sortBy;
      if (sortOrder != null) queryParams['sortOrder'] = sortOrder;

      final response = await _apiClient.get('/api/money/ai-chats', queryParameters: queryParams);

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final chats = (data['data'] as List)
              .map((chatJson) => AiChatModel.fromApiJson(chatJson).toEntity())
              .toList();
          return chats;
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch AI chats');
        }
      } else {
        throw ServerException('Failed to fetch AI chats with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<AiChatEntity> getAiChatById(String id) async {
    final cacheKey = 'ai_chat_$id';
    
    return _getCachedOrFetch(cacheKey, () async {
      try {
        log('🔍 AI Chat by ID API Request: /api/money/ai-chats/$id');

        final response = await _apiClient.get('/api/money/ai-chats/$id');

        if (response.statusCode == 200) {
          final data = json.decode(response.body);
          if (data['success'] == true) {
            return AiChatModel.fromApiJson(data['data']).toEntity();
          } else {
            throw ServerException(data['message'] ?? 'Failed to fetch AI chat');
          }
        } else if (response.statusCode == 404) {
          throw NotFoundException('AI chat not found');
        } else {
          throw ServerException('Failed to fetch AI chat with status ${response.statusCode}');
        }
      } on http.ClientException {
        throw NetworkException('Network error occurred');
      } catch (e) {
        if (e is ServerException || e is NetworkException || e is NotFoundException) {
          rethrow;
        }
        throw ServerException('Unexpected error: $e');
      }
    });
  }

  @override
  Future<AiChatEntity> createAiChat({
    required String title,
    String? ideaNodeId,
    required String agent,
    String? model,
    String? systemPrompt,
    required double temperature,
    int? maxTokens,
    required bool visible,
    required bool archived,
  }) async {
    try {
      final response = await _apiClient.post(
        '/api/money/ai-chats',
        headers: {'Content-Type': 'application/json'},
        body: {
          'title': title,
          if (ideaNodeId != null) 'ideaNodeId': ideaNodeId,
          'agent': agent,
          if (model != null) 'model': model,
          if (systemPrompt != null) 'systemPrompt': systemPrompt,
          'temperature': temperature,
          if (maxTokens != null) 'maxTokens': maxTokens,
          'visible': visible,
          'archived': archived,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final aiChat = AiChatModel.fromApiJson(data['data']).toEntity();
          
          // Invalidate AI chats cache
          _invalidateCache('ai_chats_');
          
          return aiChat;
        } else {
          throw ServerException(data['message'] ?? 'Failed to create AI chat');
        }
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to create AI chat with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<AiChatEntity> updateAiChat({
    required String id,
    String? title,
    String? ideaNodeId,
    String? agent,
    String? model,
    String? systemPrompt,
    double? temperature,
    int? maxTokens,
    bool? visible,
    bool? archived,
  }) async {
    try {
      final response = await _apiClient.put(
        '/api/money/ai-chats/$id',
        headers: {'Content-Type': 'application/json'},
        body: {
          if (title != null) 'title': title,
          if (ideaNodeId != null) 'ideaNodeId': ideaNodeId,
          if (agent != null) 'agent': agent,
          if (model != null) 'model': model,
          if (systemPrompt != null) 'systemPrompt': systemPrompt,
          if (temperature != null) 'temperature': temperature,
          if (maxTokens != null) 'maxTokens': maxTokens,
          if (visible != null) 'visible': visible,
          if (archived != null) 'archived': archived,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final aiChat = AiChatModel.fromApiJson(data['data']).toEntity();
          
          // Invalidate AI chats cache
          _invalidateCache('ai_chats_');
          _invalidateCache('ai_chat_$id');
          
          return aiChat;
        } else {
          throw ServerException(data['message'] ?? 'Failed to update AI chat');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('AI chat not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update AI chat with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> deleteAiChat(String id) async {
    try {
      final response = await _apiClient.delete('/api/money/ai-chats/$id');

      if (response.statusCode == 200 || response.statusCode == 204) {
        // Invalidate AI chats cache
        _invalidateCache('ai_chats_');
        _invalidateCache('ai_chat_$id');
      } else {
        if (response.statusCode == 404) {
          throw NotFoundException('AI chat not found');
        } else {
          throw ServerException('Failed to delete AI chat with status ${response.statusCode}');
        }
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<AiChatEntity> sendMessage({
    required String chatId,
    required String message,
  }) async {
    try {
      final response = await _apiClient.post(
        '/api/money/ai-chats/$chatId/messages',
        headers: {'Content-Type': 'application/json'},
        body: {
          'message': message,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return AiChatModel.fromApiJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to send message');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('AI chat not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to send message with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<List<ChatMessageEntity>> getChatMessages(String chatId) async {
    try {
      final response = await _apiClient.get('/api/money/ai-chats/$chatId/messages');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          // For now, return an empty list since we don't have a ChatMessageModel yet
          // TODO: Implement ChatMessageModel and proper deserialization
          return [];
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch messages');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('AI chat not found');
      } else {
        throw ServerException('Failed to fetch messages with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  // Idea Nodes Implementation
  @override
  Future<List<IdeaNodeEntity>> getIdeaNodes({
    required String ideaId,
    String? type,
    bool? visible,
    int? limit,
    int? offset,
  }) async {
    try {
      final queryParams = <String, String>{
        'ideaId': ideaId,
      };
      if (type != null) queryParams['type'] = type;
      if (visible != null) queryParams['visible'] = visible.toString();
      if (limit != null) queryParams['limit'] = limit.toString();
      if (offset != null) queryParams['offset'] = offset.toString();

      final response = await _apiClient.get('/api/money/idea-nodes', queryParameters: queryParams);

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final nodes = (data['data'] as List)
              .map((nodeJson) => IdeaNodeModel.fromApiJson(nodeJson).toEntity())
              .toList();
          return nodes;
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch idea nodes');
        }
      } else {
        throw ServerException('Failed to fetch idea nodes with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<IdeaNodeEntity> getIdeaNodeById(String id) async {
    try {
      final response = await _apiClient.get('/api/money/idea-nodes/$id');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return IdeaNodeModel.fromApiJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch idea node');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Idea node not found');
      } else {
        throw ServerException('Failed to fetch idea node with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<IdeaNodeEntity> createIdeaNode({
    required String ideaId,
    required String title,
    String? content,
    required String type,
    required double positionX,
    required double positionY,
    double? width,
    double? height,
    String? color,
    required List<String> connections,
    required bool visible,
  }) async {
    try {
      final response = await _apiClient.post(
        '/api/money/idea-nodes',
        headers: {'Content-Type': 'application/json'},
        body: {
          'ideaId': ideaId,
          'title': title,
          if (content != null) 'content': content,
          'type': type,
          'positionX': positionX,
          'positionY': positionY,
          if (width != null) 'width': width,
          if (height != null) 'height': height,
          if (color != null) 'color': color,
          'connections': connections,
          'visible': visible,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return IdeaNodeModel.fromApiJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to create idea node');
        }
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to create idea node with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<IdeaNodeEntity> updateIdeaNode({
    required String id,
    String? title,
    String? content,
    String? type,
    double? positionX,
    double? positionY,
    double? width,
    double? height,
    String? color,
    List<String>? connections,
    bool? visible,
  }) async {
    try {
      final response = await _apiClient.put(
        '/api/money/idea-nodes/$id',
        headers: {'Content-Type': 'application/json'},
        body: {
          if (title != null) 'title': title,
          if (content != null) 'content': content,
          if (type != null) 'type': type,
          if (positionX != null) 'positionX': positionX,
          if (positionY != null) 'positionY': positionY,
          if (width != null) 'width': width,
          if (height != null) 'height': height,
          if (color != null) 'color': color,
          if (connections != null) 'connections': connections,
          if (visible != null) 'visible': visible,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return IdeaNodeModel.fromApiJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to update idea node');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Idea node not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update idea node with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException || e is ValidationException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> deleteIdeaNode(String id) async {
    try {
      final response = await _apiClient.delete('/api/money/idea-nodes/$id');

      if (response.statusCode == 200 || response.statusCode == 204) {
        return;
      } else {
        if (response.statusCode == 404) {
          throw NotFoundException('Idea node not found');
        } else {
          throw ServerException('Failed to delete idea node with status ${response.statusCode}');
        }
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<IdeaNodeEntity> updateNodePosition({
    required String id,
    required double positionX,
    required double positionY,
  }) async {
    try {
      final response = await _apiClient.patch(
        '/api/money/idea-nodes/$id/position',
        headers: {'Content-Type': 'application/json'},
        body: {
          'positionX': positionX,
          'positionY': positionY,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return IdeaNodeModel.fromApiJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to update node position');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Idea node not found');
      } else {
        throw ServerException('Failed to update node position with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> updateNodePositions(List<Map<String, dynamic>> positions) async {
    try {
      final response = await _apiClient.patch(
        '/api/money/idea-nodes/positions',
        headers: {'Content-Type': 'application/json'},
        body: {
          'positions': positions,
        },
      );

      if (response.statusCode == 200) {
        return;
      } else {
        throw ServerException('Failed to update node positions with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<IdeaNodeEntity> updateNodeConnections({
    required String id,
    required List<String> connections,
  }) async {
    try {
      final response = await _apiClient.patch(
        '/api/money/idea-nodes/$id/connections',
        headers: {'Content-Type': 'application/json'},
        body: {
          'connections': connections,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return IdeaNodeModel.fromApiJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to update node connections');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Idea node not found');
      } else {
        throw ServerException('Failed to update node connections with status ${response.statusCode}');
      }
    } on http.ClientException {
      throw NetworkException('Network error occurred');
    } catch (e) {
      if (e is ServerException || e is NetworkException || e is NotFoundException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }
}
