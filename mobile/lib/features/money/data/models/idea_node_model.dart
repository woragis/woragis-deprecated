import 'dart:convert';
import 'package:json_annotation/json_annotation.dart';
import '../../domain/entities/idea_node_entity.dart';

part 'idea_node_model.g.dart';

@JsonSerializable()
class IdeaNodeModel extends IdeaNodeEntity {
  const IdeaNodeModel({
    required super.id,
    required super.ideaId,
    required super.title,
    super.content,
    required super.type,
    required super.positionX,
    required super.positionY,
    required super.width,
    required super.height,
    super.color,
    required super.connections,
    required super.visible,
    required super.createdAt,
    required super.updatedAt,
  });

  factory IdeaNodeModel.fromJson(Map<String, dynamic> json) =>
      _$IdeaNodeModelFromJson(json);

  Map<String, dynamic> toJson() => _$IdeaNodeModelToJson(this);

  // Custom methods for API (camelCase) and Local Storage (snake_case) conversion
  factory IdeaNodeModel.fromApiJson(Map<String, dynamic> json) {
    return IdeaNodeModel(
      id: json['id'] as String,
      ideaId: json['ideaId'] as String,
      title: json['title'] as String,
      content: json['content'] as String?,
      type: json['type'] as String? ?? 'default',
      positionX: (json['positionX'] as num).toDouble(),
      positionY: (json['positionY'] as num).toDouble(),
      width: (json['width'] as num?)?.toDouble() ?? 200.0,
      height: (json['height'] as num?)?.toDouble() ?? 100.0,
      color: json['color'] as String?,
      connections: (json['connections'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList() ?? [],
      visible: json['visible'] as bool,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Map<String, dynamic> toApiJson() {
    return {
      'id': id,
      'ideaId': ideaId,
      'title': title,
      'content': content,
      'type': type,
      'positionX': positionX,
      'positionY': positionY,
      'width': width,
      'height': height,
      'color': color,
      'connections': connections,
      'visible': visible,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  factory IdeaNodeModel.fromLocalJson(Map<String, dynamic> json) {
    return IdeaNodeModel(
      id: json['id'] as String,
      ideaId: json['idea_id'] as String,
      title: json['title'] as String,
      content: json['content'] as String?,
      type: json['type'] as String? ?? 'default',
      positionX: (json['position_x'] as num).toDouble(),
      positionY: (json['position_y'] as num).toDouble(),
      width: (json['width'] as num?)?.toDouble() ?? 200.0,
      height: (json['height'] as num?)?.toDouble() ?? 100.0,
      color: json['color'] as String?,
      connections: (json['connections'] as String?)
          ?.let((connectionsJson) => (jsonDecode(connectionsJson) as List<dynamic>)
              .map((e) => e as String)
              .toList()) ?? [],
      visible: (json['visible'] as int) == 1,
      createdAt: DateTime.fromMillisecondsSinceEpoch(json['created_at'] as int),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(json['updated_at'] as int),
    );
  }

  Map<String, dynamic> toLocalJson() {
    return {
      'id': id,
      'idea_id': ideaId,
      'title': title,
      'content': content,
      'type': type,
      'position_x': positionX,
      'position_y': positionY,
      'width': width,
      'height': height,
      'color': color,
      'connections': jsonEncode(connections),
      'visible': visible ? 1 : 0,
      'created_at': createdAt.millisecondsSinceEpoch,
      'updated_at': updatedAt.millisecondsSinceEpoch,
    };
  }

  factory IdeaNodeModel.fromEntity(IdeaNodeEntity entity) {
    return IdeaNodeModel(
      id: entity.id,
      ideaId: entity.ideaId,
      title: entity.title,
      content: entity.content,
      type: entity.type,
      positionX: entity.positionX,
      positionY: entity.positionY,
      width: entity.width,
      height: entity.height,
      color: entity.color,
      connections: entity.connections,
      visible: entity.visible,
      createdAt: entity.createdAt,
      updatedAt: entity.updatedAt,
    );
  }

  IdeaNodeEntity toEntity() {
    return IdeaNodeEntity(
      id: id,
      ideaId: ideaId,
      title: title,
      content: content,
      type: type,
      positionX: positionX,
      positionY: positionY,
      width: width,
      height: height,
      color: color,
      connections: connections,
      visible: visible,
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }

  @override
  IdeaNodeModel copyWith({
    String? id,
    String? ideaId,
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
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return IdeaNodeModel(
      id: id ?? this.id,
      ideaId: ideaId ?? this.ideaId,
      title: title ?? this.title,
      content: content ?? this.content,
      type: type ?? this.type,
      positionX: positionX ?? this.positionX,
      positionY: positionY ?? this.positionY,
      width: width ?? this.width,
      height: height ?? this.height,
      color: color ?? this.color,
      connections: connections ?? this.connections,
      visible: visible ?? this.visible,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

// Helper extension for null-safe operations
extension NullableExtension<T> on T? {
  R? let<R>(R Function(T) transform) {
    if (this != null) {
      return transform(this as T);
    }
    return null;
  }
}

