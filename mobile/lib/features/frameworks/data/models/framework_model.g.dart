// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'framework_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

FrameworkModel _$FrameworkModelFromJson(Map<String, dynamic> json) =>
    FrameworkModel(
      id: json['id'] as String,
      userId: json['userId'] as String,
      name: json['name'] as String,
      slug: json['slug'] as String,
      description: json['description'] as String?,
      icon: json['icon'] as String?,
      color: json['color'] as String?,
      website: json['website'] as String?,
      type: $enumDecode(_$FrameworkTypeEnumMap, json['type']),
      proficiencyLevel: $enumDecodeNullable(
        _$ProficiencyLevelEnumMap,
        json['proficiencyLevel'],
      ),
      version: json['version'] as String?,
      visible: json['visible'] as bool,
      public: json['public'] as bool,
      order: (json['order'] as num).toInt(),
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );

Map<String, dynamic> _$FrameworkModelToJson(FrameworkModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'userId': instance.userId,
      'name': instance.name,
      'slug': instance.slug,
      'description': instance.description,
      'icon': instance.icon,
      'color': instance.color,
      'website': instance.website,
      'type': _$FrameworkTypeEnumMap[instance.type]!,
      'proficiencyLevel': _$ProficiencyLevelEnumMap[instance.proficiencyLevel],
      'version': instance.version,
      'visible': instance.visible,
      'public': instance.public,
      'order': instance.order,
      'createdAt': instance.createdAt.toIso8601String(),
      'updatedAt': instance.updatedAt.toIso8601String(),
    };

const _$FrameworkTypeEnumMap = {
  FrameworkType.language: 'language',
  FrameworkType.framework: 'framework',
  FrameworkType.library: 'library',
  FrameworkType.tool: 'tool',
  FrameworkType.database: 'database',
  FrameworkType.other: 'other',
};

const _$ProficiencyLevelEnumMap = {
  ProficiencyLevel.beginner: 'beginner',
  ProficiencyLevel.intermediate: 'intermediate',
  ProficiencyLevel.advanced: 'advanced',
  ProficiencyLevel.expert: 'expert',
};
