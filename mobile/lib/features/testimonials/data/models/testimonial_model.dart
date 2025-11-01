import 'package:json_annotation/json_annotation.dart';
import '../../domain/entities/testimonial_entity.dart';

part 'testimonial_model.g.dart';

// Helper methods for safe type conversion
extension TestimonialModelHelpers on TestimonialModel {
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
class TestimonialModel extends TestimonialEntity {
  const TestimonialModel({
    required super.id,
    required super.userId,
    required super.name,
    required super.position,
    required super.company,
    required super.content,
    super.avatar,
    required super.rating,
    required super.featured,
    required super.order,
    required super.visible,
    required super.public,
    required super.createdAt,
    required super.updatedAt,
  });

  factory TestimonialModel.fromJson(Map<String, dynamic> json) {
    try {
      return _$TestimonialModelFromJson(json);
    } catch (e) {
      // Fallback: manually parse with safe type conversion
      return TestimonialModel(
        id: json['id'] as String,
        userId: json['userId'] as String,
        name: json['name'] as String,
        position: json['position'] as String,
        company: json['company'] as String,
        content: json['content'] as String,
        avatar: json['avatar'] as String?,
        rating: TestimonialModelHelpers._parseInt(json['rating']),
        featured: TestimonialModelHelpers._parseBool(json['featured']),
        order: TestimonialModelHelpers._parseInt(json['order']),
        visible: TestimonialModelHelpers._parseBool(json['visible']),
        public: TestimonialModelHelpers._parseBool(json['public']),
        createdAt: DateTime.parse(json['createdAt'] as String),
        updatedAt: DateTime.parse(json['updatedAt'] as String),
      );
    }
  }

  // Custom methods for API (camelCase) and Local Storage (snake_case) conversion
  factory TestimonialModel.fromApiJson(Map<String, dynamic> json) {
    return TestimonialModel(
      id: json['id'] as String,
      userId: json['userId'] as String,
      name: json['name'] as String,
      position: json['position'] as String,
      company: json['company'] as String,
      content: json['content'] as String,
      avatar: json['avatar'] as String?,
      rating: (json['rating'] as num).toInt(),
      featured: json['featured'] as bool,
      order: (json['order'] as num).toInt(),
      visible: json['visible'] as bool,
      public: json['public'] as bool,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Map<String, dynamic> toApiJson() {
    return {
      'id': id,
      'userId': userId,
      'name': name,
      'position': position,
      'company': company,
      'content': content,
      'avatar': avatar,
      'rating': rating,
      'featured': featured,
      'order': order,
      'visible': visible,
      'public': public,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  /// Creates a TestimonialModel from database JSON with snake_case field names
  factory TestimonialModel.fromDatabaseJson(Map<String, dynamic> json) {
    return TestimonialModel(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      name: json['name'] as String,
      position: json['position'] as String,
      company: json['company'] as String,
      content: json['content'] as String,
      avatar: json['avatar'] as String?,
      rating: json['rating'] as int,
      featured: (json['featured'] as int) == 1,
      order: json['order'] as int,
      visible: (json['visible'] as int) == 1,
      public: (json['public'] as int) == 1,
      createdAt: DateTime.fromMillisecondsSinceEpoch(json['created_at'] as int),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(json['updated_at'] as int),
    );
  }

  Map<String, dynamic> toJson() => _$TestimonialModelToJson(this);

  /// Converts to database format with snake_case field names
  Map<String, dynamic> toDatabaseJson() {
    return {
      'id': id,
      'user_id': userId,
      'name': name,
      'position': position,
      'company': company,
      'content': content,
      'avatar': avatar,
      'rating': rating,
      'featured': featured ? 1 : 0,
      'order': order,
      'visible': visible ? 1 : 0,
      'public': public ? 1 : 0,
      'created_at': createdAt.millisecondsSinceEpoch,
      'updated_at': updatedAt.millisecondsSinceEpoch,
    };
  }

  factory TestimonialModel.fromEntity(TestimonialEntity entity) {
    return TestimonialModel(
      id: entity.id,
      userId: entity.userId,
      name: entity.name,
      position: entity.position,
      company: entity.company,
      content: entity.content,
      avatar: entity.avatar,
      rating: entity.rating,
      featured: entity.featured,
      order: entity.order,
      visible: entity.visible,
      public: entity.public,
      createdAt: entity.createdAt,
      updatedAt: entity.updatedAt,
    );
  }

  TestimonialEntity toEntity() {
    return TestimonialEntity(
      id: id,
      userId: userId,
      name: name,
      position: position,
      company: company,
      content: content,
      avatar: avatar,
      rating: rating,
      featured: featured,
      order: order,
      visible: visible,
      public: public,
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }

  @override
  TestimonialModel copyWith({
    String? id,
    String? userId,
    String? name,
    String? position,
    String? company,
    String? content,
    String? avatar,
    int? rating,
    bool? featured,
    int? order,
    bool? visible,
    bool? public,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return TestimonialModel(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      name: name ?? this.name,
      position: position ?? this.position,
      company: company ?? this.company,
      content: content ?? this.content,
      avatar: avatar ?? this.avatar,
      rating: rating ?? this.rating,
      featured: featured ?? this.featured,
      order: order ?? this.order,
      visible: visible ?? this.visible,
      public: public ?? this.public,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
