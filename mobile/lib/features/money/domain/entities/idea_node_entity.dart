import 'package:equatable/equatable.dart';

class IdeaNodeEntity extends Equatable {
  final String id;
  final String ideaId;
  final String title;
  final String? content;
  final String type;
  final double positionX;
  final double positionY;
  final double width;
  final double height;
  final String? color;
  final List<String> connections;
  final bool visible;
  final DateTime createdAt;
  final DateTime updatedAt;

  const IdeaNodeEntity({
    required this.id,
    required this.ideaId,
    required this.title,
    this.content,
    required this.type,
    required this.positionX,
    required this.positionY,
    required this.width,
    required this.height,
    this.color,
    required this.connections,
    required this.visible,
    required this.createdAt,
    required this.updatedAt,
  });

  @override
  List<Object?> get props => [
        id,
        ideaId,
        title,
        content,
        type,
        positionX,
        positionY,
        width,
        height,
        color,
        connections,
        visible,
        createdAt,
        updatedAt,
      ];

  IdeaNodeEntity copyWith({
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
    return IdeaNodeEntity(
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

