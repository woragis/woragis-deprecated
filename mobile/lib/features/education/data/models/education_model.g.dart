// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'education_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

EducationModel _$EducationModelFromJson(
  Map<String, dynamic> json,
) => EducationModel(
  id: json['id'] as String,
  userId: json['user_id'] as String,
  title: json['title'] as String,
  institution: json['institution'] as String,
  description: json['description'] as String?,
  type: $enumDecode(_$EducationTypeEnumMap, json['type']),
  degreeLevel: $enumDecodeNullable(_$DegreeLevelEnumMap, json['degree_level']),
  fieldOfStudy: json['field_of_study'] as String?,
  startDate: json['start_date'] == null
      ? null
      : DateTime.parse(json['start_date'] as String),
  endDate: json['end_date'] == null
      ? null
      : DateTime.parse(json['end_date'] as String),
  completionDate: json['completion_date'] == null
      ? null
      : DateTime.parse(json['completion_date'] as String),
  grade: json['grade'] as String?,
  credits: (json['credits'] as num?)?.toInt(),
  certificateId: json['certificate_id'] as String?,
  issuer: json['issuer'] as String?,
  validityPeriod: json['validity_period'] as String?,
  pdfDocument: json['pdf_document'] as String?,
  verificationUrl: json['verification_url'] as String?,
  skills: (json['skills'] as List<dynamic>?)?.map((e) => e as String).toList(),
  order: (json['order'] as num).toInt(),
  visible: json['visible'] as bool,
  createdAt: json['created_at'] == null
      ? null
      : DateTime.parse(json['created_at'] as String),
  updatedAt: json['updated_at'] == null
      ? null
      : DateTime.parse(json['updated_at'] as String),
);

Map<String, dynamic> _$EducationModelToJson(
  EducationModel instance,
) => <String, dynamic>{
  'id': instance.id,
  'user_id': instance.userId,
  'title': instance.title,
  'institution': instance.institution,
  if (instance.description case final value?) 'description': value,
  'type': _$EducationTypeEnumMap[instance.type]!,
  if (_$DegreeLevelEnumMap[instance.degreeLevel] case final value?)
    'degree_level': value,
  if (instance.fieldOfStudy case final value?) 'field_of_study': value,
  if (instance.startDate?.toIso8601String() case final value?)
    'start_date': value,
  if (instance.endDate?.toIso8601String() case final value?) 'end_date': value,
  if (instance.completionDate?.toIso8601String() case final value?)
    'completion_date': value,
  if (instance.grade case final value?) 'grade': value,
  if (instance.credits case final value?) 'credits': value,
  if (instance.certificateId case final value?) 'certificate_id': value,
  if (instance.issuer case final value?) 'issuer': value,
  if (instance.validityPeriod case final value?) 'validity_period': value,
  if (instance.pdfDocument case final value?) 'pdf_document': value,
  if (instance.verificationUrl case final value?) 'verification_url': value,
  if (instance.skills case final value?) 'skills': value,
  'order': instance.order,
  'visible': instance.visible,
  if (instance.createdAt?.toIso8601String() case final value?)
    'created_at': value,
  if (instance.updatedAt?.toIso8601String() case final value?)
    'updated_at': value,
};

const _$EducationTypeEnumMap = {
  EducationType.degree: 'degree',
  EducationType.certificate: 'certificate',
};

const _$DegreeLevelEnumMap = {
  DegreeLevel.associate: 'associate',
  DegreeLevel.bachelor: 'bachelor',
  DegreeLevel.master: 'master',
  DegreeLevel.doctorate: 'doctorate',
  DegreeLevel.diploma: 'diploma',
};
