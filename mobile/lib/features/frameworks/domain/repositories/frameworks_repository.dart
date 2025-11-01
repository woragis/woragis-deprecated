import 'package:dartz/dartz.dart';
import '../entities/framework_entity.dart';
import '../../../../core/error/failures.dart';

abstract class FrameworksRepository {
  // Framework methods
  Future<Either<Failure, List<FrameworkEntity>>> getFrameworks({
    int? page,
    int? limit,
    bool? visible,
    bool? public,
    String? type,
    String? search,
    String? sortBy,
    String? sortOrder,
  });

  Future<Either<Failure, FrameworkEntity>> getFrameworkById(String id);
  Future<Either<Failure, FrameworkEntity>> getFrameworkBySlug(String slug);

  Future<Either<Failure, FrameworkEntity>> createFramework({
    required String name,
    required String slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    required String type,
    String? proficiencyLevel,
    required bool visible,
    required bool public,
    required int order,
  });

  Future<Either<Failure, FrameworkEntity>> updateFramework({
    required String id,
    String? name,
    String? slug,
    String? description,
    String? icon,
    String? color,
    String? website,
    String? type,
    String? proficiencyLevel,
    bool? visible,
    bool? public,
    int? order,
  });

  Future<Either<Failure, void>> deleteFramework(String id);

  // Framework ordering methods
  Future<Either<Failure, void>> updateFrameworkOrder(
    List<Map<String, dynamic>> frameworkOrders,
  );

  // Framework statistics methods
  Future<Either<Failure, int>> getFrameworkProjectCount(String frameworkId);
  Future<Either<Failure, List<FrameworkEntity>>> getFrameworksWithProjectCount();

  // Offline/Cache methods
  Future<Either<Failure, List<FrameworkEntity>>> getCachedFrameworks();
  Future<Either<Failure, void>> cacheFrameworks(List<FrameworkEntity> frameworks);
  Future<Either<Failure, FrameworkEntity?>> getCachedFramework(String id);
  Future<Either<Failure, void>> cacheFramework(FrameworkEntity framework);
}
