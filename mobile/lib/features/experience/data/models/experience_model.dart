import 'dart:convert';
import 'package:json_annotation/json_annotation.dart';
import '../../domain/entities/experience_entity.dart';

part 'experience_model.g.dart';

// Helper methods for safe type conversion
extension ExperienceModelHelpers on ExperienceModel {
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

@JsonSerializable()
class ExperienceModel extends ExperienceEntity {
  const ExperienceModel({
    required super.id,
    required super.userId,
    required super.title,
    required super.company,
    required super.period,
    required super.location,
    required super.description,
    required super.achievements,
    required super.technologies,
    required super.icon,
    required super.order,
    required super.visible,
    required super.createdAt,
    required super.updatedAt,
  });

  factory ExperienceModel.fromJson(Map<String, dynamic> json) {
    try {
      return _$ExperienceModelFromJson(json);
    } catch (e) {
      // Fallback: manually parse with safe type conversion
      return ExperienceModel(
        id: json['id'] as String,
        userId: json['userId'] as String,
        title: json['title'] as String,
        company: json['company'] as String,
        period: json['period'] as String,
        location: json['location'] as String,
        description: json['description'] as String,
        achievements: List<String>.from(json['achievements'] as List),
        technologies: List<String>.from(json['technologies'] as List),
        icon: json['icon'] as String,
        order: ExperienceModelHelpers._parseInt(json['order']),
        visible: ExperienceModelHelpers._parseBool(json['visible']),
        createdAt: DateTime.parse(json['createdAt'] as String),
        updatedAt: DateTime.parse(json['updatedAt'] as String),
      );
    }
  }

  Map<String, dynamic> toJson() => _$ExperienceModelToJson(this);

  // Custom methods for API (camelCase) and Local Storage (snake_case) conversion
  factory ExperienceModel.fromApiJson(Map<String, dynamic> json) {
    return ExperienceModel(
      id: json['id'] as String,
      userId: json['userId'] as String,
      title: json['title'] as String,
      company: json['company'] as String,
      period: json['period'] as String,
      location: json['location'] as String,
      description: json['description'] as String,
      achievements: List<String>.from(json['achievements'] as List),
      technologies: List<String>.from(json['technologies'] as List),
      icon: json['icon'] as String,
      order: (json['order'] as num).toInt(),
      visible: json['visible'] as bool,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Map<String, dynamic> toApiJson() {
    return {
      'id': id,
      'userId': userId,
      'title': title,
      'company': company,
      'period': period,
      'location': location,
      'description': description,
      'achievements': achievements,
      'technologies': technologies,
      'icon': icon,
      'order': order,
      'visible': visible,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  factory ExperienceModel.fromLocalJson(Map<String, dynamic> json) {
    return ExperienceModel(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      title: json['title'] as String,
      company: json['company'] as String,
      period: json['period'] as String,
      location: json['location'] as String,
      description: json['description'] as String,
      achievements: List<String>.from(json['achievements'] != null 
          ? jsonDecode(json['achievements'] as String) 
          : []),
      technologies: List<String>.from(json['technologies'] != null 
          ? jsonDecode(json['technologies'] as String) 
          : []),
      icon: json['icon'] as String,
      order: json['order'] as int,
      visible: (json['visible'] as int) == 1,
      createdAt: DateTime.fromMillisecondsSinceEpoch(json['created_at'] as int),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(json['updated_at'] as int),
    );
  }

  Map<String, dynamic> toLocalJson() {
    return {
      'id': id,
      'user_id': userId,
      'title': title,
      'company': company,
      'period': period,
      'location': location,
      'description': description,
      'achievements': jsonEncode(achievements),
      'technologies': jsonEncode(technologies),
      'icon': icon,
      'order': order,
      'visible': visible ? 1 : 0,
      'created_at': createdAt.millisecondsSinceEpoch,
      'updated_at': updatedAt.millisecondsSinceEpoch,
    };
  }

  factory ExperienceModel.fromEntity(ExperienceEntity entity) {
    return ExperienceModel(
      id: entity.id,
      userId: entity.userId,
      title: entity.title,
      company: entity.company,
      period: entity.period,
      location: entity.location,
      description: entity.description,
      achievements: entity.achievements,
      technologies: entity.technologies,
      icon: entity.icon,
      order: entity.order,
      visible: entity.visible,
      createdAt: entity.createdAt,
      updatedAt: entity.updatedAt,
    );
  }

  ExperienceEntity toEntity() {
    return ExperienceEntity(
      id: id,
      userId: userId,
      title: title,
      company: company,
      period: period,
      location: location,
      description: description,
      achievements: achievements,
      technologies: technologies,
      icon: icon,
      order: order,
      visible: visible,
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }

  @override
  ExperienceModel copyWith({
    String? id,
    String? userId,
    String? title,
    String? company,
    String? period,
    String? location,
    String? description,
    List<String>? achievements,
    List<String>? technologies,
    String? icon,
    int? order,
    bool? visible,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return ExperienceModel(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      title: title ?? this.title,
      company: company ?? this.company,
      period: period ?? this.period,
      location: location ?? this.location,
      description: description ?? this.description,
      achievements: achievements ?? this.achievements,
      technologies: technologies ?? this.technologies,
      icon: icon ?? this.icon,
      order: order ?? this.order,
      visible: visible ?? this.visible,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
