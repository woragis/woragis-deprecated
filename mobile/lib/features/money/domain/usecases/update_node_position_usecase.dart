import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../entities/idea_node_entity.dart';
import '../repositories/money_repository.dart';

class UpdateNodePositionUseCase {
  final MoneyRepository repository;

  UpdateNodePositionUseCase(this.repository);

  Future<Either<Failure, IdeaNodeEntity>> call({
    required String id,
    required double positionX,
    required double positionY,
  }) async {
    return await repository.updateNodePosition(
      id: id,
      positionX: positionX,
      positionY: positionY,
    );
  }
}

