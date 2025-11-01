import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../entities/idea_node_entity.dart';
import '../repositories/money_repository.dart';

class GetIdeaNodesUseCase {
  final MoneyRepository repository;

  GetIdeaNodesUseCase(this.repository);

  Future<Either<Failure, List<IdeaNodeEntity>>> call({
    required String ideaId,
    String? type,
    bool? visible,
    int? limit,
    int? offset,
  }) async {
    return await repository.getIdeaNodes(
      ideaId: ideaId,
      type: type,
      visible: visible,
      limit: limit,
      offset: offset,
    );
  }
}

