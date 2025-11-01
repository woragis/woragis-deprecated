import 'package:dartz/dartz.dart';
import '../../domain/entities/framework_entity.dart';
import '../../domain/repositories/frameworks_repository.dart';
import '../../../../core/error/failures.dart';
import '../../../../core/error/exceptions.dart';
import '../datasources/frameworks_local_datasource.dart';
import '../datasources/frameworks_remote_datasource.dart';

class FrameworksRepositoryImpl implements FrameworksRepository {
  final FrameworksRemoteDataSource remoteDataSource;
  final FrameworksLocalDataSource localDataSource;

  FrameworksRepositoryImpl({
    required this.remoteDataSource,
    required this.localDataSource,
  });

  @override
  Future<Either<Failure, List<FrameworkEntity>>> getFrameworks({
    int? page,
    int? limit,
    bool? visible,
    bool? public,
    String? type,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) async {
    try {
      final frameworks = await remoteDataSource.getFrameworks(
        page: page,
        limit: limit,
        visible: visible,
        public: public,
        type: type,
        search: search,
        sortBy: sortBy,
        sortOrder: sortOrder,
      );

      // Cache the frameworks locally
      await localDataSource.cacheFrameworks(frameworks);

      return Right(frameworks);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      // Return cached data if available
      try {
        final cachedFrameworks = await localDataSource.getCachedFrameworks();
        return Right(cachedFrameworks);
      } catch (cacheError) {
        return Left(NetworkFailure(e.message));
      }
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, FrameworkEntity>> getFrameworkById(String id) async {
    try {
      final framework = await remoteDataSource.getFrameworkById(id);
      return Right(framework);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      // Try to get from cache
      try {
        final cachedFramework = await localDataSource.getCachedFramework(id);
        if (cachedFramework != null) {
          return Right(cachedFramework);
        }
        return Left(NetworkFailure(e.message));
      } catch (cacheError) {
        return Left(NetworkFailure(e.message));
      }
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, FrameworkEntity>> getFrameworkBySlug(String slug) async {
    try {
      final framework = await remoteDataSource.getFrameworkBySlug(slug);
      return Right(framework);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
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
  }) async {
    try {
      final framework = await remoteDataSource.createFramework(
        name: name,
        slug: slug,
        description: description,
        icon: icon,
        color: color,
        website: website,
        type: type,
        proficiencyLevel: proficiencyLevel,
        visible: visible,
        public: public,
        order: order,
      );

      // Cache the new framework
      await localDataSource.cacheFramework(framework);

      return Right(framework);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
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
  }) async {
    try {
      final framework = await remoteDataSource.updateFramework(
        id: id,
        name: name,
        slug: slug,
        description: description,
        icon: icon,
        color: color,
        website: website,
        type: type,
        proficiencyLevel: proficiencyLevel,
        visible: visible,
        public: public,
        order: order,
      );

      // Update cached framework
      await localDataSource.updateCachedFramework(framework);

      return Right(framework);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> deleteFramework(String id) async {
    try {
      await remoteDataSource.deleteFramework(id);
      await localDataSource.removeCachedFramework(id);
      return const Right(null);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> updateFrameworkOrder(
    List<Map<String, dynamic>> frameworkOrders,
  ) async {
    try {
      await remoteDataSource.updateFrameworkOrder(frameworkOrders);
      return const Right(null);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, int>> getFrameworkProjectCount(String frameworkId) async {
    try {
      final count = await remoteDataSource.getFrameworkProjectCount(frameworkId);
      return Right(count);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<FrameworkEntity>>> getFrameworksWithProjectCount() async {
    try {
      final frameworks = await remoteDataSource.getFrameworksWithProjectCount();
      return Right(frameworks);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<FrameworkEntity>>> getCachedFrameworks() async {
    try {
      final frameworks = await localDataSource.getCachedFrameworks();
      return Right(frameworks);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> cacheFrameworks(List<FrameworkEntity> frameworks) async {
    try {
      await localDataSource.cacheFrameworks(frameworks);
      return const Right(null);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, FrameworkEntity?>> getCachedFramework(String id) async {
    try {
      final framework = await localDataSource.getCachedFramework(id);
      return Right(framework);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> cacheFramework(FrameworkEntity framework) async {
    try {
      await localDataSource.cacheFramework(framework);
      return const Right(null);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }
}
