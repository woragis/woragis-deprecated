import 'package:equatable/equatable.dart';

enum EducationType { degree, certificate }
enum DegreeLevel { associate, bachelor, master, doctorate, diploma }

class EducationEntity extends Equatable {
  final String id;
  final String userId;
  final String title;
  final String institution;
  final String? description;
  final EducationType type;
  final DegreeLevel? degreeLevel;
  final String? fieldOfStudy;
  final DateTime? startDate;
  final DateTime? endDate;
  final DateTime? completionDate;
  final String? grade;
  final int? credits;
  final String? certificateId;
  final String? issuer;
  final String? validityPeriod;
  final String? pdfDocument;
  final String? verificationUrl;
  final List<String>? skills;
  final int order;
  final bool visible;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  const EducationEntity({
    required this.id,
    required this.userId,
    required this.title,
    required this.institution,
    this.description,
    required this.type,
    this.degreeLevel,
    this.fieldOfStudy,
    this.startDate,
    this.endDate,
    this.completionDate,
    this.grade,
    this.credits,
    this.certificateId,
    this.issuer,
    this.validityPeriod,
    this.pdfDocument,
    this.verificationUrl,
    this.skills,
    required this.order,
    required this.visible,
    this.createdAt,
    this.updatedAt,
  });

  @override
  List<Object?> get props => [
        id,
        userId,
        title,
        institution,
        description,
        type,
        degreeLevel,
        fieldOfStudy,
        startDate,
        endDate,
        completionDate,
        grade,
        credits,
        certificateId,
        issuer,
        validityPeriod,
        pdfDocument,
        verificationUrl,
        skills,
        order,
        visible,
        createdAt,
        updatedAt,
      ];

  EducationEntity copyWith({
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
    return EducationEntity(
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
