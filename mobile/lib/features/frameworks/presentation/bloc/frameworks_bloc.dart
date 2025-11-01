import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../domain/entities/framework_entity.dart';
import '../../domain/usecases/usecases.dart';
import '../../../../core/injection/injection_container.dart';

// Events
abstract class FrameworksEvent extends Equatable {
  const FrameworksEvent();

  @override
  List<Object?> get props => [];
}

class LoadFrameworks extends FrameworksEvent {
  final int? page;
  final int? limit;
  final bool? visible;
  final bool? public;
  final String? type;
  final String? search;
  final String? sortBy;
  final String? sortOrder;

  const LoadFrameworks({
    this.page,
    this.limit,
    this.visible,
    this.public,
    this.type,
    this.search,
    this.sortBy,
    this.sortOrder,
  });

  @override
  List<Object?> get props => [
        page,
        limit,
        visible,
        public,
        type,
        search,
        sortBy,
        sortOrder,
      ];
}

class RefreshFrameworks extends FrameworksEvent {
  const RefreshFrameworks();
}

class SearchFrameworks extends FrameworksEvent {
  final String query;

  const SearchFrameworks(this.query);

  @override
  List<Object> get props => [query];
}

class FilterFrameworks extends FrameworksEvent {
  final String? type;
  final bool? visible;
  final bool? public;

  const FilterFrameworks({
    this.type,
    this.visible,
    this.public,
  });

  @override
  List<Object?> get props => [type, visible, public];
}

class SortFrameworks extends FrameworksEvent {
  final String sortBy;
  final String sortOrder;

  const SortFrameworks({
    required this.sortBy,
    required this.sortOrder,
  });

  @override
  List<Object> get props => [sortBy, sortOrder];
}

class CreateFramework extends FrameworksEvent {
  final String name;
  final String slug;
  final String? description;
  final String? icon;
  final String? color;
  final String? website;
  final String type;
  final String? proficiencyLevel;
  final int order;
  final bool visible;
  final bool public;

  const CreateFramework({
    required this.name,
    required this.slug,
    this.description,
    this.icon,
    this.color,
    this.website,
    required this.type,
    this.proficiencyLevel,
    required this.order,
    required this.visible,
    required this.public,
  });

  @override
  List<Object?> get props => [name, slug, description, icon, color, website, type, proficiencyLevel, order, visible, public];
}

class UpdateFramework extends FrameworksEvent {
  final String id;
  final String? name;
  final String? slug;
  final String? description;
  final String? icon;
  final String? color;
  final String? website;
  final String? type;
  final String? proficiencyLevel;
  final int? order;
  final bool? visible;
  final bool? public;

  const UpdateFramework({
    required this.id,
    this.name,
    this.slug,
    this.description,
    this.icon,
    this.color,
    this.website,
    this.type,
    this.proficiencyLevel,
    this.order,
    this.visible,
    this.public,
  });

  @override
  List<Object?> get props => [id, name, slug, description, icon, color, website, type, proficiencyLevel, order, visible, public];
}

class DeleteFramework extends FrameworksEvent {
  final String id;

  const DeleteFramework(this.id);

  @override
  List<Object> get props => [id];
}

// States
abstract class FrameworksState extends Equatable {
  const FrameworksState();

  @override
  List<Object?> get props => [];
}

class FrameworksInitial extends FrameworksState {}

class FrameworksLoading extends FrameworksState {}

class FrameworksLoaded extends FrameworksState {
  final List<FrameworkEntity> frameworks;
  final bool hasReachedMax;
  final int currentPage;
  final String? currentSearch;
  final String? currentType;
  final bool? currentVisible;
  final bool? currentPublic;
  final String currentSortBy;
  final String currentSortOrder;

  const FrameworksLoaded({
    required this.frameworks,
    required this.hasReachedMax,
    required this.currentPage,
    this.currentSearch,
    this.currentType,
    this.currentVisible,
    this.currentPublic,
    required this.currentSortBy,
    required this.currentSortOrder,
  });

  @override
  List<Object?> get props => [
        frameworks,
        hasReachedMax,
        currentPage,
        currentSearch,
        currentType,
        currentVisible,
        currentPublic,
        currentSortBy,
        currentSortOrder,
      ];

  FrameworksLoaded copyWith({
    List<FrameworkEntity>? frameworks,
    bool? hasReachedMax,
    int? currentPage,
    String? currentSearch,
    String? currentType,
    bool? currentVisible,
    bool? currentPublic,
    String? currentSortBy,
    String? currentSortOrder,
  }) {
    return FrameworksLoaded(
      frameworks: frameworks ?? this.frameworks,
      hasReachedMax: hasReachedMax ?? this.hasReachedMax,
      currentPage: currentPage ?? this.currentPage,
      currentSearch: currentSearch ?? this.currentSearch,
      currentType: currentType ?? this.currentType,
      currentVisible: currentVisible ?? this.currentVisible,
      currentPublic: currentPublic ?? this.currentPublic,
      currentSortBy: currentSortBy ?? this.currentSortBy,
      currentSortOrder: currentSortOrder ?? this.currentSortOrder,
    );
  }
}

class FrameworksError extends FrameworksState {
  final String message;

  const FrameworksError(this.message);

  @override
  List<Object> get props => [message];
}

class FrameworkCreated extends FrameworksState {
  final FrameworkEntity framework;

  const FrameworkCreated(this.framework);

  @override
  List<Object> get props => [framework];
}

class FrameworkUpdated extends FrameworksState {
  final FrameworkEntity framework;

  const FrameworkUpdated(this.framework);

  @override
  List<Object> get props => [framework];
}

class FrameworkDeleted extends FrameworksState {
  final String id;

  const FrameworkDeleted(this.id);

  @override
  List<Object> get props => [id];
}

// BLoC
class FrameworksBloc extends Bloc<FrameworksEvent, FrameworksState> {
  final GetFrameworksUseCase getFrameworksUseCase;
  final CreateFrameworkUseCase createFrameworkUseCase;
  final UpdateFrameworkUseCase updateFrameworkUseCase;
  final DeleteFrameworkUseCase deleteFrameworkUseCase;

  FrameworksBloc({
    required this.getFrameworksUseCase,
    required this.createFrameworkUseCase,
    required this.updateFrameworkUseCase,
    required this.deleteFrameworkUseCase,
  }) : super(FrameworksInitial()) {
    on<LoadFrameworks>(_onLoadFrameworks);
    on<RefreshFrameworks>(_onRefreshFrameworks);
    on<SearchFrameworks>(_onSearchFrameworks);
    on<FilterFrameworks>(_onFilterFrameworks);
    on<SortFrameworks>(_onSortFrameworks);
    on<CreateFramework>(_onCreateFramework);
    on<UpdateFramework>(_onUpdateFramework);
    on<DeleteFramework>(_onDeleteFramework);
  }

  Future<void> _onLoadFrameworks(
    LoadFrameworks event,
    Emitter<FrameworksState> emit,
  ) async {
    if (state is FrameworksLoaded) {
      final currentState = state as FrameworksLoaded;
      if (currentState.hasReachedMax) return;
    }

    emit(FrameworksLoading());

    try {
      final params = GetFrameworksParams(
        page: event.page ?? 1,
        limit: event.limit ?? 20,
        visible: event.visible,
        public: event.public,
        type: event.type,
        search: event.search,
        sortBy: event.sortBy ?? 'order',
        sortOrder: event.sortOrder ?? 'asc',
      );

      final result = await getFrameworksUseCase(params);

      result.fold(
        (failure) => emit(FrameworksError(failure.message)),
        (frameworks) {
          if (state is FrameworksLoaded) {
            final currentState = state as FrameworksLoaded;
            emit(FrameworksLoaded(
              frameworks: [...currentState.frameworks, ...frameworks],
              hasReachedMax: frameworks.length < (event.limit ?? 20),
              currentPage: event.page ?? 1,
              currentSearch: event.search,
              currentType: event.type,
              currentVisible: event.visible,
              currentPublic: event.public,
              currentSortBy: event.sortBy ?? 'order',
              currentSortOrder: event.sortOrder ?? 'asc',
            ));
          } else {
            emit(FrameworksLoaded(
              frameworks: frameworks,
              hasReachedMax: frameworks.length < (event.limit ?? 20),
              currentPage: event.page ?? 1,
              currentSearch: event.search,
              currentType: event.type,
              currentVisible: event.visible,
              currentPublic: event.public,
              currentSortBy: event.sortBy ?? 'order',
              currentSortOrder: event.sortOrder ?? 'asc',
            ));
          }
        },
      );
    } catch (e) {
      emit(FrameworksError('Failed to load frameworks: $e'));
    }
  }

  Future<void> _onRefreshFrameworks(
    RefreshFrameworks event,
    Emitter<FrameworksState> emit,
  ) async {
    if (state is FrameworksLoaded) {
      final currentState = state as FrameworksLoaded;
      add(LoadFrameworks(
        page: 1,
        limit: 20,
        visible: currentState.currentVisible,
        public: currentState.currentPublic,
        type: currentState.currentType,
        search: currentState.currentSearch,
        sortBy: currentState.currentSortBy,
        sortOrder: currentState.currentSortOrder,
      ));
    } else {
      add(const LoadFrameworks());
    }
  }

  Future<void> _onSearchFrameworks(
    SearchFrameworks event,
    Emitter<FrameworksState> emit,
  ) async {
    add(LoadFrameworks(
      page: 1,
      limit: 20,
      search: event.query.isEmpty ? null : event.query,
    ));
  }

  Future<void> _onFilterFrameworks(
    FilterFrameworks event,
    Emitter<FrameworksState> emit,
  ) async {
    add(LoadFrameworks(
      page: 1,
      limit: 20,
      type: event.type,
      visible: event.visible,
      public: event.public,
    ));
  }

  Future<void> _onSortFrameworks(
    SortFrameworks event,
    Emitter<FrameworksState> emit,
  ) async {
    add(LoadFrameworks(
      page: 1,
      limit: 20,
      sortBy: event.sortBy,
      sortOrder: event.sortOrder,
    ));
  }

  Future<void> _onCreateFramework(
    CreateFramework event,
    Emitter<FrameworksState> emit,
  ) async {
    emit(FrameworksLoading());

    final result = await createFrameworkUseCase(
      name: event.name,
      slug: event.slug,
      description: event.description,
      icon: event.icon,
      color: event.color,
      website: event.website,
      type: event.type,
      proficiencyLevel: event.proficiencyLevel,
      order: event.order,
      visible: event.visible,
      public: event.public,
    );

    result.fold(
      (failure) => emit(FrameworksError(failure.message)),
      (framework) => emit(FrameworkCreated(framework)),
    );
  }

  Future<void> _onUpdateFramework(
    UpdateFramework event,
    Emitter<FrameworksState> emit,
  ) async {
    emit(FrameworksLoading());

    final result = await updateFrameworkUseCase(
      id: event.id,
      name: event.name,
      slug: event.slug,
      description: event.description,
      icon: event.icon,
      color: event.color,
      website: event.website,
      type: event.type,
      proficiencyLevel: event.proficiencyLevel,
      order: event.order,
      visible: event.visible,
      public: event.public,
    );

    result.fold(
      (failure) => emit(FrameworksError(failure.message)),
      (framework) => emit(FrameworkUpdated(framework)),
    );
  }

  Future<void> _onDeleteFramework(
    DeleteFramework event,
    Emitter<FrameworksState> emit,
  ) async {
    emit(FrameworksLoading());

    final result = await deleteFrameworkUseCase(event.id);

    result.fold(
      (failure) => emit(FrameworksError(failure.message)),
      (_) => emit(FrameworkDeleted(event.id)),
    );
  }
}

// Factory function
FrameworksBloc createFrameworksBloc() {
  return FrameworksBloc(
    getFrameworksUseCase: sl<GetFrameworksUseCase>(),
    createFrameworkUseCase: sl<CreateFrameworkUseCase>(),
    updateFrameworkUseCase: sl<UpdateFrameworkUseCase>(),
    deleteFrameworkUseCase: sl<DeleteFrameworkUseCase>(),
  );
}
