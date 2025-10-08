import 'package:json_annotation/json_annotation.dart';
import '../../domain/entities/education_entity.dart';

part 'education_model.g.dart';

enum EducationTypeModel {
  @JsonValue('degree')
  degree,
  @JsonValue('certificate')
  certificate,
}

enum DegreeLevelModel {
  @JsonValue('associate')
  associate,
  @JsonValue('bachelor')
  bachelor,
  @JsonValue('master')
  master,
  @JsonValue('doctorate')
  doctorate,
  @JsonValue('diploma')
  diploma,
}

@JsonSerializable(
  fieldRename: FieldRename.snake,
  explicitToJson: true,
  includeIfNull: false,
)
class EducationModel extends EducationEntity {
  const EducationModel({
    required super.id,
    required super.userId,
    required super.title,
    required super.institution,
    super.description,
    required super.type,
    super.degreeLevel,
    super.fieldOfStudy,
    super.startDate,
    super.endDate,
    super.completionDate,
    super.grade,
    super.credits,
    super.certificateId,
    super.issuer,
    super.validityPeriod,
    super.pdfDocument,
    super.verificationUrl,
    super.skills,
    required super.order,
    required super.visible,
    super.createdAt,
    super.updatedAt,
  });

  factory EducationModel.fromJson(Map<String, dynamic> json) {
    // Handle potential null values with safe defaults
    return EducationModel(
      id: json['id']?.toString() ?? '',
      userId: json['user_id']?.toString() ?? '',
      title: json['title']?.toString() ?? '',
      institution: json['institution']?.toString() ?? '',
      description: json['description']?.toString(),
      type: _parseEducationType(json['type']),
      degreeLevel: _parseDegreeLevel(json['degree_level']),
      fieldOfStudy: json['field_of_study']?.toString(),
      startDate: json['start_date'] != null 
          ? DateTime.tryParse(json['start_date'].toString())
          : null,
      endDate: json['end_date'] != null 
          ? DateTime.tryParse(json['end_date'].toString())
          : null,
      completionDate: json['completion_date'] != null 
          ? DateTime.tryParse(json['completion_date'].toString())
          : null,
      grade: json['grade']?.toString(),
      credits: json['credits'] != null 
          ? int.tryParse(json['credits'].toString())
          : null,
      certificateId: json['certificate_id']?.toString(),
      issuer: json['issuer']?.toString(),
      validityPeriod: json['validity_period']?.toString(),
      pdfDocument: json['pdf_document']?.toString(),
      verificationUrl: json['verification_url']?.toString(),
      skills: json['skills'] is List 
          ? (json['skills'] as List).map((e) => e.toString()).toList()
          : null,
      order: json['order'] != null 
          ? int.tryParse(json['order'].toString()) ?? 0
          : 0,
      visible: json['visible'] is bool 
          ? json['visible'] as bool
          : true,
      createdAt: json['created_at'] != null 
          ? DateTime.tryParse(json['created_at'].toString())
          : null,
      updatedAt: json['updated_at'] != null 
          ? DateTime.tryParse(json['updated_at'].toString())
          : null,
    );
  }

  static EducationType _parseEducationType(dynamic value) {
    if (value == null) return EducationType.degree;
    final stringValue = value.toString().toLowerCase();
    switch (stringValue) {
      case 'certificate':
        return EducationType.certificate;
      case 'degree':
      default:
        return EducationType.degree;
    }
  }

  static DegreeLevel? _parseDegreeLevel(dynamic value) {
    if (value == null) return null;
    final stringValue = value.toString().toLowerCase();
    switch (stringValue) {
      case 'associate':
        return DegreeLevel.associate;
      case 'bachelor':
        return DegreeLevel.bachelor;
      case 'master':
        return DegreeLevel.master;
      case 'doctorate':
        return DegreeLevel.doctorate;
      case 'diploma':
        return DegreeLevel.diploma;
      default:
        return null;
    }
  }

  Map<String, dynamic> toJson() => _$EducationModelToJson(this);

  factory EducationModel.fromEntity(EducationEntity entity) {
    return EducationModel(
      id: entity.id,
      userId: entity.userId,
      title: entity.title,
      institution: entity.institution,
      description: entity.description,
      type: entity.type,
      degreeLevel: entity.degreeLevel,
      fieldOfStudy: entity.fieldOfStudy,
      startDate: entity.startDate,
      endDate: entity.endDate,
      completionDate: entity.completionDate,
      grade: entity.grade,
      credits: entity.credits,
      certificateId: entity.certificateId,
      issuer: entity.issuer,
      validityPeriod: entity.validityPeriod,
      pdfDocument: entity.pdfDocument,
      verificationUrl: entity.verificationUrl,
      skills: entity.skills,
      order: entity.order,
      visible: entity.visible,
      createdAt: entity.createdAt,
      updatedAt: entity.updatedAt,
    );
  }

  EducationEntity toEntity() {
    return EducationEntity(
      id: id,
      userId: userId,
      title: title,
      institution: institution,
      description: description,
      type: type,
      degreeLevel: degreeLevel,
      fieldOfStudy: fieldOfStudy,
      startDate: startDate,
      endDate: endDate,
      completionDate: completionDate,
      grade: grade,
      credits: credits,
      certificateId: certificateId,
      issuer: issuer,
      validityPeriod: validityPeriod,
      pdfDocument: pdfDocument,
      verificationUrl: verificationUrl,
      skills: skills,
      order: order,
      visible: visible,
      createdAt: createdAt,
      updatedAt: updatedAt,
    );
  }

  @override
  EducationModel copyWith({
    String? id,
    String? userId,
    String? title,
    String? institution,
    String? description,
    EducationType? type,
    DegreeLevel? degreeLevel,
    String? fieldOfStudy,
    DateTime? startDate,
    DateTime? endDate,
    DateTime? completionDate,
    String? grade,
    int? credits,
    String? certificateId,
    String? issuer,
    String? validityPeriod,
    String? pdfDocument,
    String? verificationUrl,
    List<String>? skills,
    int? order,
    bool? visible,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return EducationModel(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      title: title ?? this.title,
      institution: institution ?? this.institution,
      description: description ?? this.description,
      type: type ?? this.type,
      degreeLevel: degreeLevel ?? this.degreeLevel,
      fieldOfStudy: fieldOfStudy ?? this.fieldOfStudy,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      completionDate: completionDate ?? this.completionDate,
      grade: grade ?? this.grade,
      credits: credits ?? this.credits,
      certificateId: certificateId ?? this.certificateId,
      issuer: issuer ?? this.issuer,
      validityPeriod: validityPeriod ?? this.validityPeriod,
      pdfDocument: pdfDocument ?? this.pdfDocument,
      verificationUrl: verificationUrl ?? this.verificationUrl,
      skills: skills ?? this.skills,
      order: order ?? this.order,
      visible: visible ?? this.visible,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}
