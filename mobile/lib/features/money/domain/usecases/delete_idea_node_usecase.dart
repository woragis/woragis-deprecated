import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../repositories/money_repository.dart';

class DeleteIdeaNodeUseCase {
  final MoneyRepository repository;

  DeleteIdeaNodeUseCase(this.repository);

  Future<Either<Failure, void>> call(String id) async {
    return await repository.deleteIdeaNode(id);
  }
}

