import 'package:dartz/dartz.dart';
import '../entities/project_entity.dart';
import '../../../../core/error/failures.dart';

abstract class ProjectsRepository {
  // Project methods
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
  });

  Future<Either<Failure, ProjectEntity>> getProjectById(String id);
  Future<Either<Failure, ProjectEntity>> getProjectBySlug(String slug);

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
  });

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
  });

  Future<Either<Failure, void>> deleteProject(String id);
  
  // Project interaction methods
  Future<Either<Failure, void>> incrementViewCount(String id);
  Future<Either<Failure, void>> incrementLikeCount(String id);

  // Project ordering methods
  Future<Either<Failure, void>> updateProjectOrder(
    List<Map<String, dynamic>> projectOrders,
  );

  // Project frameworks relationship methods
  Future<Either<Failure, void>> assignFrameworkToProject({
    required String projectId,
    required String frameworkId,
  });

  Future<Either<Failure, void>> removeFrameworkFromProject({
    required String projectId,
    required String frameworkId,
  });

  Future<Either<Failure, List<String>>> getProjectFrameworkIds(String projectId);

  // Toggle methods
  Future<Either<Failure, void>> toggleProjectFeatured(String id);
  Future<Either<Failure, void>> toggleProjectVisibility(String id);

  // Offline/Cache methods
  Future<Either<Failure, List<ProjectEntity>>> getCachedProjects();
  Future<Either<Failure, void>> cacheProjects(List<ProjectEntity> projects);
  Future<Either<Failure, ProjectEntity?>> getCachedProject(String id);
  Future<Either<Failure, void>> cacheProject(ProjectEntity project);
}
