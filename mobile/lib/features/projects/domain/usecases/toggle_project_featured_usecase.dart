import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../repositories/projects_repository.dart';

class ToggleProjectFeaturedUseCase {
  final ProjectsRepository repository;

  ToggleProjectFeaturedUseCase(this.repository);

  Future<Either<Failure, void>> call(String id) async {
    return await repository.toggleProjectFeatured(id);
  }
}
