import 'package:json_annotation/json_annotation.dart';
import '../../domain/entities/project_entity.dart';
import '../../../frameworks/data/models/framework_model.dart';

part 'project_model.g.dart';

@JsonSerializable(explicitToJson: true)
class ProjectModel extends ProjectEntity {
  @override
  final List<FrameworkModel>? frameworks;

  const ProjectModel({
    required super.id,
    required super.userId,
    required super.title,
    required super.slug,
    required super.description,
    super.longDescription,
    super.content,
    super.videoUrl,
    required super.image,
    super.githubUrl,
    super.liveUrl,
    required super.featured,
    required super.order,
    required super.visible,
    required super.public,
    this.frameworks,
    required super.createdAt,
    required super.updatedAt,
  }) : super(frameworks: frameworks);

  factory ProjectModel.fromJson(Map<String, dynamic> json) {
    // Handle the case where frameworks might not be present in the response
    final jsonWithFrameworks = Map<String, dynamic>.from(json);
    if (!jsonWithFrameworks.containsKey('frameworks')) {
      jsonWithFrameworks['frameworks'] = null;
    }
    return _$ProjectModelFromJson(jsonWithFrameworks);
  }

  // Custom methods for API (camelCase) and Local Storage (snake_case) conversion
  factory ProjectModel.fromApiJson(Map<String, dynamic> json) {
    return ProjectModel(
      id: json['id'] as String,
      userId: json['userId'] as String,
      title: json['title'] as String,
      slug: json['slug'] as String,
      description: json['description'] as String,
      longDescription: json['longDescription'] as String?,
      content: json['content'] as String?,
      videoUrl: json['videoUrl'] as String?,
      image: json['image'] as String,
      githubUrl: json['githubUrl'] as String?,
      liveUrl: json['liveUrl'] as String?,
      featured: json['featured'] as bool,
      order: (json['order'] as num).toInt(),
      visible: json['visible'] as bool,
      public: json['public'] as bool,
      frameworks: null, // Frameworks are handled separately in API responses
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Map<String, dynamic> toApiJson() {
    return {
      'id': id,
      'userId': userId,
      'title': title,
      'description': description,
      'longDescription': longDescription,
      'content': content,
      'videoUrl': videoUrl,
      'image': image,
      'githubUrl': githubUrl,
      'liveUrl': liveUrl,
      'featured': featured,
      'order': order,
      'visible': visible,
      'public': public,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  /// Creates a ProjectModel from database JSON with snake_case field names
  factory ProjectModel.fromDatabaseJson(Map<String, dynamic> json) {
    return ProjectModel(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      title: json['title'] as String,
      slug: json['slug'] as String,
      description: json['description'] as String,
      longDescription: json['long_description'] as String?,
      content: json['content'] as String?,
      videoUrl: json['video_url'] as String?,
      image: json['image'] as String,
      githubUrl: json['github_url'] as String?,
      liveUrl: json['live_url'] as String?,
      featured: (json['featured'] as int) == 1,
      order: json['order'] as int,
      visible: (json['visible'] as int) == 1,
      public: (json['public'] as int) == 1,
      frameworks: null, // Frameworks are stored in separate project_frameworks table
      createdAt: DateTime.fromMillisecondsSinceEpoch(json['created_at'] as int),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(json['updated_at'] as int),
    );
  }

  Map<String, dynamic> toJson() => _$ProjectModelToJson(this);

  /// Converts to database format with snake_case field names
  Map<String, dynamic> toDatabaseJson() {
    return {
      'id': id,
      'user_id': userId,
      'title': title,
      'slug': slug, // ✅ ADDED MISSING SLUG FIELD
      'description': description,
      'long_description': longDescription,
      'content': content,
      'video_url': videoUrl,
      'image': image,
      'github_url': githubUrl,
      'live_url': liveUrl,
      'featured': featured ? 1 : 0,
      'order': order,
      'visible': visible ? 1 : 0,
      'public': public ? 1 : 0,
      'created_at': createdAt.millisecondsSinceEpoch,
      'updated_at': updatedAt.millisecondsSinceEpoch,
    };
  }

  factory ProjectModel.fromEntity(ProjectEntity entity) {
    return ProjectModel(
      id: entity.id,
      userId: entity.userId,
      title: entity.title,
      slug: entity.slug,
      description: entity.description,
      longDescription: entity.longDescription,
      content: entity.content,
      videoUrl: entity.videoUrl,
      image: entity.image,
      githubUrl: entity.githubUrl,
      liveUrl: entity.liveUrl,
      featured: entity.featured,
      order: entity.order,
      visible: entity.visible,
      public: entity.public,
      frameworks: entity.frameworks?.map<FrameworkModel>((fw) => FrameworkModel.fromEntity(fw)).toList(),
      createdAt: entity.createdAt,
      updatedAt: entity.updatedAt,
    );
  }

  ProjectEntity toEntity() {
    return ProjectEntity(
      id: id,
      userId: userId,
      title: title,
      slug: slug,
      description: description,
      longDescription: longDescription,
      content: content,
      videoUrl: videoUrl,
      image: image,
      githubUrl: githubUrl,
      liveUrl: liveUrl,
      featured: featured,
      order: order,
      visible: visible,
      public: public,
      frameworks: frameworks?.map((fw) => fw.toEntity()).toList(),
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }

  @override
  ProjectModel copyWith({
    String? id,
    String? userId,
    String? title,
    String? slug,
    String? description,
    String? longDescription,
    String? content,
    String? videoUrl,
    String? image,
    String? githubUrl,
    String? liveUrl,
    bool? featured,
    int? order,
    bool? visible,
    bool? public,
    covariant List<FrameworkModel>? frameworks,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return ProjectModel(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      title: title ?? this.title,
      slug: slug ?? this.slug,
      description: description ?? this.description,
      longDescription: longDescription ?? this.longDescription,
      content: content ?? this.content,
      videoUrl: videoUrl ?? this.videoUrl,
      image: image ?? this.image,
      githubUrl: githubUrl ?? this.githubUrl,
      liveUrl: liveUrl ?? this.liveUrl,
      featured: featured ?? this.featured,
      order: order ?? this.order,
      visible: visible ?? this.visible,
      public: public ?? this.public,
      frameworks: frameworks ?? this.frameworks,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
