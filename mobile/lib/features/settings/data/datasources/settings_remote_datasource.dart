import 'dart:convert';
import 'package:http/http.dart' as http;
import '../../../../core/error/exceptions.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/entities/setting_entity.dart';
import '../models/setting_model.dart';

abstract class SettingsRemoteDataSource {
  Future<List<SettingEntity>> getSettings({
    String? category,
    bool? isPublic,
  });

  Future<SettingEntity> getSettingById(String id);
  Future<SettingEntity> getSettingByKey(String key);
  Future<SettingEntity> createSetting({
    required String key,
    required String value,
    String? category,
    String? description,
    String? type,
    bool? isPublic,
    bool? isEditable,
  });
  Future<SettingEntity> updateSetting({
    required String id,
    String? key,
    String? value,
    String? category,
    String? description,
    String? type,
    bool? isPublic,
    bool? isEditable,
  });
  Future<SettingEntity> updateSettingByKey({
    required String key,
    required String value,
  });
  Future<void> deleteSetting(String id);

  // Bulk operations
  Future<Map<String, String>> getSettingsAsMap({
    String? category,
    bool? isPublic,
  });
  Future<void> updateSettingsBulk(Map<String, String> settings);

  // Category-specific settings
  Future<Map<String, String>> getCoreProfileSettings();
  Future<Map<String, String>> getSocialMediaSettings();
  Future<Map<String, String>> getContactSettings();
  Future<Map<String, String>> getSiteSettings();
  Future<void> updateCoreProfileSettings(Map<String, String> settings);
  Future<void> updateSocialMediaSettings(Map<String, String> settings);
  Future<void> updateContactSettings(Map<String, String> settings);
  Future<void> updateSiteSettings(Map<String, String> settings);
}

class SettingsRemoteDataSourceImpl implements SettingsRemoteDataSource {
  final ApiClient _apiClient;
  final String baseUrl;

  SettingsRemoteDataSourceImpl({
    required ApiClient apiClient,
    required this.baseUrl,
  }) : _apiClient = apiClient;

  @override
  Future<List<SettingEntity>> getSettings({
    String? category,
    bool? isPublic,
  }) async {
    try {
      final queryParams = <String, String>{};
      if (category != null) queryParams['category'] = category;
      if (isPublic != null) queryParams['isPublic'] = isPublic.toString();

      final response = await _apiClient.get('/admin/settings', queryParameters: queryParams);

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          final settings = (data['data']['settings'] as List)
              .map((settingJson) => SettingModel.fromJson(settingJson).toEntity())
              .toList();
          return settings;
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch settings');
        }
      } else {
        throw ServerException('Failed to fetch settings with status ${response.statusCode}');
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
  Future<SettingEntity> getSettingById(String id) async {
    try {
      final response = await _apiClient.get('/admin/settings/$id');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return SettingModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch setting');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Setting not found');
      } else {
        throw ServerException('Failed to fetch setting with status ${response.statusCode}');
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
  Future<SettingEntity> getSettingByKey(String key) async {
    try {
      final response = await _apiClient.get('/admin/settings/key/$key');

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return SettingModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to fetch setting');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Setting not found');
      } else {
        throw ServerException('Failed to fetch setting with status ${response.statusCode}');
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
  Future<SettingEntity> createSetting({
    required String key,
    required String value,
    String? category,
    String? description,
    String? type,
    bool? isPublic,
    bool? isEditable,
  }) async {
    try {
      final response = await _apiClient.post(
        '/admin/settings',
        headers: {'Content-Type': 'application/json'},
        body: {
          'key': key,
          'value': value,
          if (category != null) 'category': category,
          if (description != null) 'description': description,
          if (type != null) 'type': type,
          if (isPublic != null) 'isPublic': isPublic,
          if (isEditable != null) 'isEditable': isEditable,
        },
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return SettingModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to create setting');
        }
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to create setting with status ${response.statusCode}');
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
  Future<SettingEntity> updateSetting({
    required String id,
    String? key,
    String? value,
    String? category,
    String? description,
    String? type,
    bool? isPublic,
    bool? isEditable,
  }) async {
    try {
      final response = await _apiClient.put(
        '/admin/settings/$id',
        headers: {'Content-Type': 'application/json'},
        body: {
          if (key != null) 'key': key,
          if (value != null) 'value': value,
          if (category != null) 'category': category,
          if (description != null) 'description': description,
          if (type != null) 'type': type,
          if (isPublic != null) 'isPublic': isPublic,
          if (isEditable != null) 'isEditable': isEditable,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return SettingModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to update setting');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Setting not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update setting with status ${response.statusCode}');
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
  Future<SettingEntity> updateSettingByKey({
    required String key,
    required String value,
  }) async {
    try {
      final response = await _apiClient.put(
        '/admin/settings/key/$key',
        headers: {'Content-Type': 'application/json'},
        body: {
          'value': value,
        },
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        if (data['success'] == true) {
          return SettingModel.fromJson(data['data']).toEntity();
        } else {
          throw ServerException(data['message'] ?? 'Failed to update setting');
        }
      } else if (response.statusCode == 404) {
        throw NotFoundException('Setting not found');
      } else if (response.statusCode == 422) {
        final data = json.decode(response.body);
        throw ValidationException(data['message'] ?? 'Validation failed');
      } else {
        throw ServerException('Failed to update setting with status ${response.statusCode}');
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
  Future<void> deleteSetting(String id) async {
    try {
      final response = await _apiClient.delete('/admin/settings/$id');

      if (response.statusCode != 200 && response.statusCode != 204) {
        if (response.statusCode == 404) {
          throw NotFoundException('Setting not found');
        } else {
          throw ServerException('Failed to delete setting with status ${response.statusCode}');
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

  Future<void> deleteSettingByKey(String key) async {
    try {
      final response = await _apiClient.delete('/admin/settings/key/$key');

      if (response.statusCode != 200 && response.statusCode != 204) {
        if (response.statusCode == 404) {
          throw NotFoundException('Setting not found');
        } else {
          throw ServerException('Failed to delete setting with status ${response.statusCode}');
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

  // Bulk operations
  @override
  Future<Map<String, String>> getSettingsAsMap({
    String? category,
    bool? isPublic,
  }) async {
    try {
      final settings = await getSettings(category: category, isPublic: isPublic);
      return Map.fromEntries(
        settings.map((setting) => MapEntry(setting.key, setting.value)),
      );
    } catch (e) {
      if (e is ServerException || e is NetworkException) {
        rethrow;
      }
      throw ServerException('Unexpected error: $e');
    }
  }

  @override
  Future<void> updateSettingsBulk(Map<String, String> settings) async {
    try {
      final response = await _apiClient.put(
        '/admin/settings/bulk',
        body: {'settings': settings},
      );

      if (response.statusCode != 200) {
        throw ServerException('Failed to update settings bulk with status ${response.statusCode}');
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

  // Category-specific settings
  @override
  Future<Map<String, String>> getCoreProfileSettings() async {
    return getSettingsAsMap(category: 'core_profile');
  }

  @override
  Future<Map<String, String>> getSocialMediaSettings() async {
    return getSettingsAsMap(category: 'social_media');
  }

  @override
  Future<Map<String, String>> getContactSettings() async {
    return getSettingsAsMap(category: 'contact');
  }

  @override
  Future<Map<String, String>> getSiteSettings() async {
    return getSettingsAsMap(category: 'site');
  }

  @override
  Future<void> updateCoreProfileSettings(Map<String, String> settings) async {
    return updateSettingsBulk(settings);
  }

  @override
  Future<void> updateSocialMediaSettings(Map<String, String> settings) async {
    return updateSettingsBulk(settings);
  }

  @override
  Future<void> updateContactSettings(Map<String, String> settings) async {
    return updateSettingsBulk(settings);
  }

  @override
  Future<void> updateSiteSettings(Map<String, String> settings) async {
    return updateSettingsBulk(settings);
  }
}
