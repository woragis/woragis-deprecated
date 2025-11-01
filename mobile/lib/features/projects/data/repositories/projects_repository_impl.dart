import 'package:dartz/dartz.dart';
import '../../domain/entities/project_entity.dart';
import '../../domain/repositories/projects_repository.dart';
import '../../../../core/error/failures.dart';
import '../../../../core/error/exceptions.dart';
import '../datasources/projects_local_datasource.dart';
import '../datasources/projects_remote_datasource.dart';

class ProjectsRepositoryImpl implements ProjectsRepository {
  final ProjectsRemoteDataSource remoteDataSource;
  final ProjectsLocalDataSource localDataSource;

  ProjectsRepositoryImpl({
    required this.remoteDataSource,
    required this.localDataSource,
  });

  @override
  Future<Either<Failure, List<ProjectEntity>>> getProjects({
    int? page,
    int? limit,
    bool? featured,
    bool? visible,
    bool? public,
    List<String>? technologies,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) async {
    try {
      final projects = await remoteDataSource.getProjects(
        page: page,
        limit: limit,
        featured: featured,
        visible: visible,
        public: public,
        search: search,
      );

      // Cache the projects locally
      await localDataSource.cacheProjects(projects);

      return Right(projects);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      // Return cached data if available
      try {
        final cachedProjects = await localDataSource.getCachedProjects(
          limit: limit,
          featured: featured,
          visible: visible,
          public: public,
          search: search,
        );
        return Right(cachedProjects);
      } catch (cacheError) {
        return Left(NetworkFailure(e.message));
      }
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, ProjectEntity>> getProjectById(String id) async {
    try {
      final project = await remoteDataSource.getProjectById(id);
      // Cache the project locally
      await localDataSource.cacheProject(project);
      return Right(project);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      // Try to get from cache
      try {
        final cachedProject = await localDataSource.getCachedProject(id);
        if (cachedProject != null) {
          return Right(cachedProject);
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
  Future<Either<Failure, ProjectEntity>> getProjectBySlug(String slug) async {
    try {
      final project = await remoteDataSource.getProjectBySlug(slug);
      // Cache the project locally
      await localDataSource.cacheProject(project);
      return Right(project);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }


  @override
  Future<Either<Failure, ProjectEntity>> createProject({
    required String title,
    required String slug,
    required String description,
    String? longDescription,
    String? content,
    String? videoUrl,
    required String image,
    String? githubUrl,
    String? liveUrl,
    required bool featured,
    required int order,
    required bool visible,
    required bool public,
    List<String>? frameworkIds,
  }) async {
    try {
      final project = await remoteDataSource.createProject(
        title: title,
        slug: slug,
        description: description,
        longDescription: longDescription,
        content: content,
        videoUrl: videoUrl,
        image: image,
        githubUrl: githubUrl,
        liveUrl: liveUrl,
        featured: featured,
        order: order,
        visible: visible,
        public: public,
        frameworkIds: frameworkIds,
      );

      // Cache the new project
      await localDataSource.cacheProject(project);

      return Right(project);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, ProjectEntity>> updateProject({
    required String id,
    String? title,
    String? slug,
    String? description,
    String? longDescription,
    String? content,
    String? videoUrl,
    String? image,
    String? githubUrl,
    String? liveUrl,
    bool? featured,
    int? order,
    bool? visible,
    bool? public,
    List<String>? frameworkIds,
  }) async {
    try {
      final project = await remoteDataSource.updateProject(
        id: id,
        title: title,
        slug: slug,
        description: description,
        longDescription: longDescription,
        content: content,
        videoUrl: videoUrl,
        image: image,
        githubUrl: githubUrl,
        liveUrl: liveUrl,
        featured: featured,
        order: order,
        visible: visible,
        public: public,
        frameworkIds: frameworkIds,
      );

      // Update cached project
      await localDataSource.updateCachedProject(project);

      return Right(project);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> deleteProject(String id) async {
    try {
      await remoteDataSource.deleteProject(id);
      await localDataSource.removeCachedProject(id);
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
  Future<Either<Failure, void>> incrementViewCount(String id) async {
    try {
      await remoteDataSource.incrementViewCount(id);
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
  Future<Either<Failure, void>> incrementLikeCount(String id) async {
    try {
      await remoteDataSource.incrementLikeCount(id);
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
  Future<Either<Failure, void>> updateProjectOrder(
    List<Map<String, dynamic>> projectOrders,
  ) async {
    try {
      // Update each project order individually
      for (final orderData in projectOrders) {
        await remoteDataSource.updateProjectOrder(
          id: orderData['id'] as String,
          order: orderData['order'] as int,
        );
      }
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
  Future<Either<Failure, void>> assignFrameworkToProject({
    required String projectId,
    required String frameworkId,
  }) async {
    try {
      await remoteDataSource.assignFrameworkToProject(
        projectId: projectId,
        frameworkId: frameworkId,
      );
      
      // Also update local cache
      await localDataSource.assignFrameworkToProject(projectId, frameworkId);
      
      return const Right(null);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      // Still update local cache for offline support
      try {
        await localDataSource.assignFrameworkToProject(projectId, frameworkId);
        return const Right(null);
      } catch (cacheError) {
        return Left(NetworkFailure(e.message));
      }
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> removeFrameworkFromProject({
    required String projectId,
    required String frameworkId,
  }) async {
    try {
      await remoteDataSource.removeFrameworkFromProject(
        projectId: projectId,
        frameworkId: frameworkId,
      );
      
      // Also update local cache
      await localDataSource.removeFrameworkFromProject(projectId, frameworkId);
      
      return const Right(null);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      // Still update local cache for offline support
      try {
        await localDataSource.removeFrameworkFromProject(projectId, frameworkId);
        return const Right(null);
      } catch (cacheError) {
        return Left(NetworkFailure(e.message));
      }
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<String>>> getProjectFrameworkIds(String projectId) async {
    try {
      final frameworkIds = await remoteDataSource.getProjectFrameworkIds(projectId);
      return Right(frameworkIds);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      // Try to get from local cache
      try {
        final cachedFrameworkIds = await localDataSource.getProjectFrameworkIds(projectId);
        return Right(cachedFrameworkIds);
      } catch (cacheError) {
        return Left(NetworkFailure(e.message));
      }
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }


  @override
  Future<Either<Failure, List<ProjectEntity>>> getCachedProjects() async {
    try {
      final projects = await localDataSource.getCachedProjects();
      return Right(projects);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> cacheProjects(List<ProjectEntity> projects) async {
    try {
      await localDataSource.cacheProjects(projects);
      return const Right(null);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, ProjectEntity?>> getCachedProject(String id) async {
    try {
      final project = await localDataSource.getCachedProject(id);
      return Right(project);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> cacheProject(ProjectEntity project) async {
    try {
      await localDataSource.cacheProject(project);
      return const Right(null);
    } catch (e) {
      return Left(CacheFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> toggleProjectFeatured(String id) async {
    try {
      await remoteDataSource.toggleProjectFeatured(id);
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
  Future<Either<Failure, void>> toggleProjectVisibility(String id) async {
    try {
      await remoteDataSource.toggleProjectVisibility(id);
      return const Right(null);
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Left(NetworkFailure(e.message));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }
}
