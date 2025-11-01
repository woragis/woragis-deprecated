import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../entities/framework_entity.dart';
import '../repositories/frameworks_repository.dart';

class CreateFrameworkUseCase {
  final FrameworksRepository repository;

  CreateFrameworkUseCase(this.repository);

  Future<Either<Failure, FrameworkEntity>> call({
    required String name,
    required String slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    required String type,
    String? proficiencyLevel,
    required int order,
    required bool visible,
    required bool public,
  }) async {
    return await repository.createFramework(
      name: name,
      slug: slug,
      description: description,
      icon: icon,
      color: color,
      website: website,
      type: type,
      proficiencyLevel: proficiencyLevel,
      order: order,
      visible: visible,
      public: public,
    );
  }
}
