// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'idea_node_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdeaNodeModel _$IdeaNodeModelFromJson(Map<String, dynamic> json) =>
    IdeaNodeModel(
      id: json['id'] as String,
      ideaId: json['ideaId'] as String,
      title: json['title'] as String,
      content: json['content'] as String?,
      type: json['type'] as String,
      positionX: (json['positionX'] as num).toDouble(),
      positionY: (json['positionY'] as num).toDouble(),
      width: (json['width'] as num).toDouble(),
      height: (json['height'] as num).toDouble(),
      color: json['color'] as String?,
      connections: (json['connections'] as List<dynamic>)
          .map((e) => e as String)
          .toList(),
      visible: json['visible'] as bool,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );

Map<String, dynamic> _$IdeaNodeModelToJson(IdeaNodeModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'ideaId': instance.ideaId,
      'title': instance.title,
      'content': instance.content,
      'type': instance.type,
      'positionX': instance.positionX,
      'positionY': instance.positionY,
      'width': instance.width,
      'height': instance.height,
      'color': instance.color,
      'connections': instance.connections,
      'visible': instance.visible,
      'createdAt': instance.createdAt.toIso8601String(),
      'updatedAt': instance.updatedAt.toIso8601String(),
    };
