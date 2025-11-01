import 'package:dartz/dartz.dart';
import '../../../../core/error/failures.dart';
import '../entities/idea_node_entity.dart';
import '../repositories/money_repository.dart';

class UpdateIdeaNodeUseCase {
  final MoneyRepository repository;

  UpdateIdeaNodeUseCase(this.repository);

  Future<Either<Failure, IdeaNodeEntity>> call({
    required String id,
    String? title,
    String? content,
    String? type,
    double? positionX,
    double? positionY,
    double? width,
    double? height,
    String? color,
    List<String>? connections,
    bool? visible,
  }) async {
    return await repository.updateIdeaNode(
      id: id,
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

