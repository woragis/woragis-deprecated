import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../../../core/injection/injection_container.dart';
import '../../domain/entities/project_entity.dart';
import '../../domain/usecases/usecases.dart';

// Events
abstract class ProjectsEvent extends Equatable {
  const ProjectsEvent();

  @override
  List<Object?> get props => [];
}

class GetProjectsRequested extends ProjectsEvent {
  final int? page;
  final int? limit;
  final bool? featured;
  final bool? visible;
  final bool? public;
  final List<String>? technologies;
  final String? search;
  final String? sortBy;
  final String? sortOrder;

  const GetProjectsRequested({
    this.page,
    this.limit,
    this.featured,
    this.visible,
    this.public,
    this.technologies,
    this.search,
    this.sortBy,
    this.sortOrder,
  });

  @override
  List<Object?> get props => [
        page,
        limit,
        featured,
        visible,
        public,
        technologies,
        search,
        sortBy,
        sortOrder,
      ];
}

class GetProjectByIdRequested extends ProjectsEvent {
  final String id;

  const GetProjectByIdRequested(this.id);

  @override
  List<Object> get props => [id];
}

class CreateProjectRequested extends ProjectsEvent {
  final String title;
  final String slug;
  final String description;
  final String? longDescription;
  final String? content;
  final String image;
  final String? githubUrl;
  final String? liveUrl;
  final String? videoUrl;
  final bool featured;
  final bool visible;
  final bool public;
  final int order;
  final List<String>? frameworkIds;

  const CreateProjectRequested({
    required this.title,
    required this.slug,
    required this.description,
    this.longDescription,
    this.content,
    required this.image,
    this.githubUrl,
    this.liveUrl,
    this.videoUrl,
    required this.featured,
    required this.visible,
    required this.public,
    required this.order,
    this.frameworkIds,
  });

  @override
  List<Object?> get props => [
        title,
        slug,
        description,
        longDescription,
        content,
        image,
        githubUrl,
        liveUrl,
        videoUrl,
        featured,
        visible,
        public,
        order,
        frameworkIds,
      ];
}

class UpdateProjectRequested extends ProjectsEvent {
  final String id;
  final String? title;
  final String? slug;
  final String? description;
  final String? longDescription;
  final String? content;
  final String? image;
  final String? githubUrl;
  final String? liveUrl;
  final String? videoUrl;
  final bool? featured;
  final bool? visible;
  final bool? public;
  final int? order;
  final List<String>? frameworkIds;

  const UpdateProjectRequested({
    required this.id,
    this.title,
    this.slug,
    this.description,
    this.longDescription,
    this.content,
    this.image,
    this.githubUrl,
    this.liveUrl,
    this.videoUrl,
    this.featured,
    this.visible,
    this.public,
    this.order,
    this.frameworkIds,
  });

  @override
  List<Object?> get props => [
        id,
        title,
        slug,
        description,
        longDescription,
        content,
        image,
        githubUrl,
        liveUrl,
        videoUrl,
        featured,
        visible,
        public,
        order,
        frameworkIds,
      ];
}

class DeleteProjectRequested extends ProjectsEvent {
  final String id;

  const DeleteProjectRequested(this.id);

  @override
  List<Object> get props => [id];
}

class UpdateProjectOrderRequested extends ProjectsEvent {
  final List<Map<String, dynamic>> projectOrders;

  const UpdateProjectOrderRequested(this.projectOrders);

  @override
  List<Object> get props => [projectOrders];
}

class IncrementViewCountRequested extends ProjectsEvent {
  final String projectId;

  const IncrementViewCountRequested(this.projectId);

  @override
  List<Object> get props => [projectId];
}

class IncrementLikeCountRequested extends ProjectsEvent {
  final String projectId;

  const IncrementLikeCountRequested(this.projectId);

  @override
  List<Object> get props => [projectId];
}

class ToggleProjectFeaturedRequested extends ProjectsEvent {
  final String projectId;

  const ToggleProjectFeaturedRequested(this.projectId);

  @override
  List<Object> get props => [projectId];
}

class ToggleProjectVisibilityRequested extends ProjectsEvent {
  final String projectId;

  const ToggleProjectVisibilityRequested(this.projectId);

  @override
  List<Object> get props => [projectId];
}

// States
abstract class ProjectsState extends Equatable {
  const ProjectsState();

  @override
  List<Object?> get props => [];
}

class ProjectsInitial extends ProjectsState {}

class ProjectsLoading extends ProjectsState {}

class ProjectsError extends ProjectsState {
  final String message;

  const ProjectsError(this.message);

  @override
  List<Object> get props => [message];
}

class ProjectsLoaded extends ProjectsState {
  final List<ProjectEntity> projects;

  const ProjectsLoaded(this.projects);

  @override
  List<Object> get props => [projects];
}

class ProjectLoaded extends ProjectsState {
  final ProjectEntity project;

  const ProjectLoaded(this.project);

  @override
  List<Object> get props => [project];
}

class ProjectCreated extends ProjectsState {
  final ProjectEntity project;

  const ProjectCreated(this.project);

  @override
  List<Object> get props => [project];
}

class ProjectUpdated extends ProjectsState {
  final ProjectEntity project;

  const ProjectUpdated(this.project);

  @override
  List<Object> get props => [project];
}

class ProjectDeleted extends ProjectsState {
  final String id;

  const ProjectDeleted(this.id);

  @override
  List<Object> get props => [id];
}

class ProjectOrderUpdated extends ProjectsState {
  final List<Map<String, dynamic>> projectOrders;

  const ProjectOrderUpdated(this.projectOrders);

  @override
  List<Object> get props => [projectOrders];
}

class ViewCountIncremented extends ProjectsState {
  final String projectId;

  const ViewCountIncremented(this.projectId);

  @override
  List<Object> get props => [projectId];
}

class LikeCountIncremented extends ProjectsState {
  final String projectId;

  const LikeCountIncremented(this.projectId);

  @override
  List<Object> get props => [projectId];
}

class ProjectFeaturedToggled extends ProjectsState {
  final String projectId;

  const ProjectFeaturedToggled(this.projectId);

  @override
  List<Object> get props => [projectId];
}

class ProjectVisibilityToggled extends ProjectsState {
  final String projectId;

  const ProjectVisibilityToggled(this.projectId);

  @override
  List<Object> get props => [projectId];
}

// BLoC
class ProjectsBloc extends Bloc<ProjectsEvent, ProjectsState> {
  final GetProjectsUseCase getProjectsUseCase;
  final GetProjectByIdUseCase getProjectByIdUseCase;
  final GetProjectBySlugUseCase getProjectBySlugUseCase;
  final CreateProjectUseCase createProjectUseCase;
  final UpdateProjectUseCase updateProjectUseCase;
  final DeleteProjectUseCase deleteProjectUseCase;
  final UpdateProjectOrderUseCase updateProjectOrderUseCase;
  final IncrementProjectViewCountUseCase incrementViewCountUseCase;
  final IncrementProjectLikeCountUseCase incrementLikeCountUseCase;
  final ToggleProjectFeaturedUseCase toggleProjectFeaturedUseCase;
  final ToggleProjectVisibilityUseCase toggleProjectVisibilityUseCase;

  ProjectsBloc({
    required this.getProjectsUseCase,
    required this.getProjectByIdUseCase,
    required this.getProjectBySlugUseCase,
    required this.createProjectUseCase,
    required this.updateProjectUseCase,
    required this.deleteProjectUseCase,
    required this.updateProjectOrderUseCase,
    required this.incrementViewCountUseCase,
    required this.incrementLikeCountUseCase,
    required this.toggleProjectFeaturedUseCase,
    required this.toggleProjectVisibilityUseCase,
  }) : super(ProjectsInitial()) {
    on<GetProjectsRequested>(_onGetProjectsRequested);
    on<GetProjectByIdRequested>(_onGetProjectByIdRequested);
    on<CreateProjectRequested>(_onCreateProjectRequested);
    on<UpdateProjectRequested>(_onUpdateProjectRequested);
    on<DeleteProjectRequested>(_onDeleteProjectRequested);
    on<UpdateProjectOrderRequested>(_onUpdateProjectOrderRequested);
    on<IncrementViewCountRequested>(_onIncrementViewCountRequested);
    on<IncrementLikeCountRequested>(_onIncrementLikeCountRequested);
    on<ToggleProjectFeaturedRequested>(_onToggleProjectFeaturedRequested);
    on<ToggleProjectVisibilityRequested>(_onToggleProjectVisibilityRequested);
  }

  Future<void> _onGetProjectsRequested(
    GetProjectsRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    emit(ProjectsLoading());

    final result = await getProjectsUseCase(GetProjectsParams(
      page: event.page,
      limit: event.limit,
      featured: event.featured,
      visible: event.visible,
      public: event.public,
      technologies: event.technologies,
      search: event.search,
      sortBy: event.sortBy,
      sortOrder: event.sortOrder,
    ));

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (projects) => emit(ProjectsLoaded(projects)),
    );
  }

  Future<void> _onGetProjectByIdRequested(
    GetProjectByIdRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    emit(ProjectsLoading());

    final result = await getProjectByIdUseCase(event.id);

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (project) => emit(ProjectLoaded(project)),
    );
  }

  Future<void> _onCreateProjectRequested(
    CreateProjectRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    emit(ProjectsLoading());

    final result = await createProjectUseCase(CreateProjectParams(
      title: event.title,
      slug: event.slug,
      description: event.description,
      longDescription: event.longDescription,
      content: event.content,
      videoUrl: event.videoUrl,
      image: event.image,
      githubUrl: event.githubUrl,
      liveUrl: event.liveUrl,
      frameworkIds: event.frameworkIds,
      featured: event.featured,
      order: event.order,
      visible: event.visible,
      public: event.public,
    ));

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (project) => emit(ProjectCreated(project)),
    );
  }

  Future<void> _onUpdateProjectRequested(
    UpdateProjectRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    emit(ProjectsLoading());

    final result = await updateProjectUseCase(
      id: event.id,
      title: event.title,
      slug: event.slug,
      description: event.description,
      longDescription: event.longDescription,
      content: event.content,
      videoUrl: event.videoUrl,
      image: event.image,
      githubUrl: event.githubUrl,
      liveUrl: event.liveUrl,
      frameworkIds: event.frameworkIds,
      featured: event.featured,
      order: event.order,
      visible: event.visible,
      public: event.public,
    );

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (project) => emit(ProjectUpdated(project)),
    );
  }

  Future<void> _onDeleteProjectRequested(
    DeleteProjectRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    emit(ProjectsLoading());

    final result = await deleteProjectUseCase(event.id);

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (_) => emit(ProjectDeleted(event.id)),
    );
  }

  Future<void> _onUpdateProjectOrderRequested(
    UpdateProjectOrderRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    emit(ProjectsLoading());

    final result = await updateProjectOrderUseCase(event.projectOrders);

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (_) => emit(ProjectOrderUpdated(event.projectOrders)),
    );
  }

  Future<void> _onIncrementViewCountRequested(
    IncrementViewCountRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    final result = await incrementViewCountUseCase(event.projectId);

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (_) => emit(ViewCountIncremented(event.projectId)),
    );
  }

  Future<void> _onIncrementLikeCountRequested(
    IncrementLikeCountRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    final result = await incrementLikeCountUseCase(event.projectId);

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (_) => emit(LikeCountIncremented(event.projectId)),
    );
  }

  Future<void> _onToggleProjectFeaturedRequested(
    ToggleProjectFeaturedRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    final result = await toggleProjectFeaturedUseCase(event.projectId);

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (_) => emit(ProjectFeaturedToggled(event.projectId)),
    );
  }

  Future<void> _onToggleProjectVisibilityRequested(
    ToggleProjectVisibilityRequested event,
    Emitter<ProjectsState> emit,
  ) async {
    final result = await toggleProjectVisibilityUseCase(event.projectId);

    result.fold(
      (failure) => emit(ProjectsError(failure.message)),
      (_) => emit(ProjectVisibilityToggled(event.projectId)),
    );
  }
}

// Factory function for creating ProjectsBloc with dependency injection
ProjectsBloc createProjectsBloc() {
  return ProjectsBloc(
    getProjectsUseCase: sl<GetProjectsUseCase>(),
    getProjectByIdUseCase: sl<GetProjectByIdUseCase>(),
    getProjectBySlugUseCase: sl<GetProjectBySlugUseCase>(),
    createProjectUseCase: sl<CreateProjectUseCase>(),
    updateProjectUseCase: sl<UpdateProjectUseCase>(),
    deleteProjectUseCase: sl<DeleteProjectUseCase>(),
    updateProjectOrderUseCase: sl<UpdateProjectOrderUseCase>(),
    incrementViewCountUseCase: sl<IncrementProjectViewCountUseCase>(),
    incrementLikeCountUseCase: sl<IncrementProjectLikeCountUseCase>(),
    toggleProjectFeaturedUseCase: sl<ToggleProjectFeaturedUseCase>(),
    toggleProjectVisibilityUseCase: sl<ToggleProjectVisibilityUseCase>(),
  );
}
