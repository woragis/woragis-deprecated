import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../repositories/projects_repository.dart';

class ToggleProjectVisibilityUseCase {
  final ProjectsRepository repository;

  ToggleProjectVisibilityUseCase(this.repository);

  Future<Either<Failure, void>> call(String id) async {
    return await repository.toggleProjectVisibility(id);
  }
}
