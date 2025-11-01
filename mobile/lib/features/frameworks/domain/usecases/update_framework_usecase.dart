import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../entities/framework_entity.dart';
import '../repositories/frameworks_repository.dart';

class UpdateFrameworkUseCase {
  final FrameworksRepository repository;

  UpdateFrameworkUseCase(this.repository);

  Future<Either<Failure, FrameworkEntity>> call({
    required String id,
    String? name,
    String? slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    String? type,
    String? proficiencyLevel,
    int? order,
    bool? visible,
    bool? public,
  }) async {
    return await repository.updateFramework(
      id: id,
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
