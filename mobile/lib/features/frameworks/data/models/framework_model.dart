import 'package:json_annotation/json_annotation.dart';
import '../../domain/entities/framework_entity.dart';
import '../../../../core/stores/auth_store.dart';

part 'framework_model.g.dart';

// Helper methods for safe type conversion
extension FrameworkModelHelpers on FrameworkModel {
  static int _parseInt(dynamic value) {
    if (value == null) return 0;
    if (value is int) return value;
    if (value is String) return int.tryParse(value) ?? 0;
    if (value is num) return value.toInt();
    return 0;
  }

  static bool _parseBool(dynamic value) {
    if (value == null) return false;
    if (value is bool) return value;
    if (value is int) return value == 1;
    if (value is String) return value.toLowerCase() == 'true' || value == '1';
    return false;
  }
}

// Helper functions for enum conversion
FrameworkType _convertToEntityType(FrameworkTypeModel modelType) {
  switch (modelType) {
    case FrameworkTypeModel.language:
      return FrameworkType.language;
    case FrameworkTypeModel.framework:
      return FrameworkType.framework;
    case FrameworkTypeModel.library:
      return FrameworkType.library;
    case FrameworkTypeModel.tool:
      return FrameworkType.tool;
    case FrameworkTypeModel.database:
      return FrameworkType.database;
    case FrameworkTypeModel.other:
      return FrameworkType.other;
  }
}

ProficiencyLevel _convertToEntityProficiencyLevel(ProficiencyLevelModel modelLevel) {
  switch (modelLevel) {
    case ProficiencyLevelModel.beginner:
      return ProficiencyLevel.beginner;
    case ProficiencyLevelModel.intermediate:
      return ProficiencyLevel.intermediate;
    case ProficiencyLevelModel.advanced:
      return ProficiencyLevel.advanced;
    case ProficiencyLevelModel.expert:
      return ProficiencyLevel.expert;
  }
}

enum FrameworkTypeModel {
  @JsonValue('language')
  language,
  @JsonValue('framework')
  framework,
  @JsonValue('library')
  library,
  @JsonValue('tool')
  tool,
  @JsonValue('database')
  database,
  @JsonValue('other')
  other,
}

enum ProficiencyLevelModel {
  @JsonValue('beginner')
  beginner,
  @JsonValue('intermediate')
  intermediate,
  @JsonValue('advanced')
  advanced,
  @JsonValue('expert')
  expert,
}

@JsonSerializable()
class FrameworkModel extends FrameworkEntity {
  const FrameworkModel({
    required super.id,
    required super.userId,
    required super.name,
    required super.slug,
    super.description,
    super.icon,
    super.color,
    super.website,
    required super.type,
    super.proficiencyLevel,
    super.version,
    required super.visible,
    required super.public,
    required super.order,
    required super.createdAt,
    required super.updatedAt,
  });

  factory FrameworkModel.fromJson(Map<String, dynamic> json, {String? userId}) {
    try {
      final model = _$FrameworkModelFromJson(json);
      // Override userId with the one from auth store if provided
      if (userId != null) {
        return model.copyWith(userId: userId);
      }
      return model;
    } catch (e) {
      // Fallback: manually parse with safe type conversion
      return FrameworkModel(
        id: json['id'] as String,
        userId: userId ?? json['userId'] as String, // Use auth store userId if provided
        name: json['name'] as String,
        slug: json['slug'] as String,
        description: json['description'] as String?,
        icon: json['icon'] as String?,
        color: json['color'] as String?,
        website: json['website'] as String?,
        type: _convertToEntityType(FrameworkTypeModel.values.firstWhere(
          (e) => e.name == json['type'],
          orElse: () => FrameworkTypeModel.other,
        )),
        proficiencyLevel: json['proficiencyLevel'] != null
            ? _convertToEntityProficiencyLevel(ProficiencyLevelModel.values.firstWhere(
                (e) => e.name == json['proficiencyLevel'],
                orElse: () => ProficiencyLevelModel.beginner,
              ))
            : null,
        version: json['version'] as String?,
        visible: FrameworkModelHelpers._parseBool(json['visible']),
        public: FrameworkModelHelpers._parseBool(json['public']),
        order: FrameworkModelHelpers._parseInt(json['order']),
        createdAt: DateTime.parse(json['createdAt'] as String),
        updatedAt: DateTime.parse(json['updatedAt'] as String),
      );
    }
  }

  Map<String, dynamic> toJson() {
    final json = _$FrameworkModelToJson(this);
    // Convert booleans to integers for SQLite compatibility
    json['visible'] = json['visible'] ? 1 : 0;
    json['public'] = json['public'] ? 1 : 0;
    return json;
  }

  factory FrameworkModel.fromEntity(FrameworkEntity entity) {
    return FrameworkModel(
      id: entity.id,
      userId: entity.userId,
      name: entity.name,
      slug: entity.slug,
      description: entity.description,
      icon: entity.icon,
      color: entity.color,
      website: entity.website,
      type: entity.type,
      proficiencyLevel: entity.proficiencyLevel,
      version: entity.version,
      visible: entity.visible,
      public: entity.public,
      order: entity.order,
      createdAt: entity.createdAt,
      updatedAt: entity.updatedAt,
    );
  }

  FrameworkEntity toEntity() {
    return FrameworkEntity(
      id: id,
      userId: userId,
      name: name,
      slug: slug,
      description: description,
      icon: icon,
      color: color,
      website: website,
      type: type,
      proficiencyLevel: proficiencyLevel,
      version: version,
      visible: visible,
      public: public,
      order: order,
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }

  @override
  FrameworkModel copyWith({
    String? id,
    String? userId,
    String? name,
    String? slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    FrameworkType? type,
    ProficiencyLevel? proficiencyLevel,
    String? version,
    bool? visible,
    bool? public,
    int? order,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return FrameworkModel(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      name: name ?? this.name,
      slug: slug ?? this.slug,
      description: description ?? this.description,
      icon: icon ?? this.icon,
      color: color ?? this.color,
      website: website ?? this.website,
      type: type ?? this.type,
      proficiencyLevel: proficiencyLevel ?? this.proficiencyLevel,
      version: version ?? this.version,
      visible: visible ?? this.visible,
      public: public ?? this.public,
      order: order ?? this.order,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
