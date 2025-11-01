import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../entities/idea_node_entity.dart';
import '../repositories/money_repository.dart';

class CreateIdeaNodeUseCase {
  final MoneyRepository repository;

  CreateIdeaNodeUseCase(this.repository);

  Future<Either<Failure, IdeaNodeEntity>> call({
    required String ideaId,
    required String title,
    String? content,
    required String type,
    required double positionX,
    required double positionY,
    double? width,
    double? height,
    String? color,
    required List<String> connections,
    required bool visible,
  }) async {
    return await repository.createIdeaNode(
      ideaId: ideaId,
      title: title,
      content: content,
      type: type,
      positionX: positionX,
      positionY: positionY,
      width: width,
      height: height,
      color: color,
      connections: connections,
      visible: visible,
    );
  }
}

