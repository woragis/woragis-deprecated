import 'package:get_it/get_it.dart';
import '../stores/auth_store.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;

// Core
import '../config/env_config.dart';
import '../database/database_helper.dart';
import '../database/sync_manager.dart';
import '../network/network_info.dart';
import '../network/api_client.dart';

// Core BLoCs
import '../presentation/bloc/theme/theme_bloc.dart';
import '../presentation/bloc/language/language_bloc.dart';
import '../presentation/bloc/navigation/navigation_bloc.dart';

// Auth
import '../../features/auth/data/datasources/auth_local_datasource.dart';
import '../../features/auth/data/datasources/auth_remote_datasource.dart';
import '../../features/auth/data/repositories/auth_repository_impl.dart';
import '../../features/auth/domain/repositories/auth_repository.dart';
import '../../features/auth/domain/usecases/usecases.dart';
import '../../features/auth/presentation/bloc/auth_bloc.dart';

// Blog
import '../../features/blog/data/datasources/blog_local_datasource.dart';
import '../../features/blog/data/datasources/blog_remote_datasource.dart';
import '../../features/blog/data/repositories/blog_repository_impl.dart';
import '../../features/blog/domain/repositories/blog_repository.dart';
import '../../features/blog/domain/usecases/usecases.dart';
import '../../features/blog/presentation/bloc/create_blog_post/create_blog_post_bloc.dart';
import '../../features/blog/presentation/bloc/blog_filter/blog_filter_bloc.dart';
import '../../features/blog/presentation/bloc/blog_bloc.dart';

// Projects
import '../../features/projects/data/datasources/projects_local_datasource.dart';
import '../../features/projects/data/datasources/projects_remote_datasource.dart';
import '../../features/projects/data/repositories/projects_repository_impl.dart';
import '../../features/projects/domain/repositories/projects_repository.dart';
import '../../features/projects/domain/usecases/usecases.dart';
import '../../features/projects/presentation/bloc/projects_bloc.dart';

// Frameworks
import '../../features/frameworks/data/datasources/frameworks_local_datasource.dart';
import '../../features/frameworks/data/datasources/frameworks_remote_datasource.dart';
import '../../features/frameworks/data/repositories/frameworks_repository_impl.dart';
import '../../features/frameworks/domain/repositories/frameworks_repository.dart';
import '../../features/frameworks/domain/usecases/usecases.dart';
import '../../features/frameworks/presentation/bloc/frameworks_bloc.dart';

// About
import '../../features/about/data/datasources/about_local_datasource.dart';
import '../../features/about/data/datasources/about_remote_datasource.dart';
import '../../features/about/data/repositories/about_repository_impl.dart';
import '../../features/about/domain/repositories/about_repository.dart';
import '../../features/about/domain/usecases/usecases.dart';
import '../../features/about/presentation/bloc/about_bloc.dart';

// Experience

// Frameworks

// Education
import '../../features/education/data/datasources/education_local_datasource.dart';
import '../../features/education/data/datasources/education_remote_datasource.dart';
import '../../features/education/data/repositories/education_repository_impl.dart';
import '../../features/education/domain/repositories/education_repository.dart';
import '../../features/education/domain/usecases/usecases.dart';
import '../../features/education/presentation/bloc/education_bloc.dart';

// Experience
import '../../features/experience/data/datasources/experience_local_datasource.dart';
import '../../features/experience/data/datasources/experience_remote_datasource.dart';
import '../../features/experience/data/repositories/experience_repository_impl.dart';
import '../../features/experience/domain/repositories/experience_repository.dart';
import '../../features/experience/domain/usecases/usecases.dart';
import '../../features/experience/presentation/bloc/experience_bloc.dart';

// Money
import '../../features/money/data/datasources/money_local_datasource.dart';
import '../../features/money/data/datasources/money_remote_datasource.dart';
import '../../features/money/data/repositories/money_repository_impl.dart';
import '../../features/money/domain/repositories/money_repository.dart';
import '../../features/money/domain/usecases/usecases.dart';
import '../../features/money/presentation/bloc/money_bloc.dart';

// Testimonials
import '../../features/testimonials/data/datasources/testimonials_local_datasource.dart';
import '../../features/testimonials/data/datasources/testimonials_remote_datasource.dart';
import '../../features/testimonials/data/repositories/testimonials_repository_impl.dart';
import '../../features/testimonials/domain/repositories/testimonials_repository.dart';
import '../../features/testimonials/domain/usecases/usecases.dart';
import '../../features/testimonials/presentation/bloc/testimonials_bloc.dart';

// Settings
import '../../features/settings/data/datasources/settings_local_datasource.dart';
import '../../features/settings/data/datasources/settings_remote_datasource.dart';
import '../../features/settings/data/repositories/settings_repository_impl.dart';
import '../../features/settings/domain/repositories/settings_repository.dart';
import '../../features/settings/domain/usecases/usecases.dart';
import '../../features/settings/presentation/bloc/settings_bloc.dart';

final sl = GetIt.instance;

Future<void> init() async {
  // Core
  sl.registerLazySingleton<DatabaseHelper>(() => DatabaseHelper());
  sl.registerLazySingleton<SyncManager>(() => SyncManager());
  
  // External
  sl.registerLazySingleton(() => Connectivity());
  sl.registerLazySingleton<http.Client>(() => http.Client());
  
  // API Client (centralized with auto-auth headers)
  sl.registerLazySingleton<ApiClient>(
    () => ApiClient(
      httpClient: sl<http.Client>(),
      authStore: sl<AuthStoreBloc>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );
  
  // SharedPreferences (for BLoCs)
  final sharedPreferences = await SharedPreferences.getInstance();
  sl.registerLazySingleton<SharedPreferences>(() => sharedPreferences);
  
  // Core - Network Info
  sl.registerLazySingleton<NetworkInfo>(() => NetworkInfoImpl(sl()));
  
  // Core - Query Client Manager (Dio wrapper)
  
  
  // Core - Auth Store
  sl.registerLazySingleton<AuthStoreBloc>(() => AuthStoreBloc());
  
  

  // Core BLoCs (for local UI state only)
  _initCoreBLoCs();

  // Auth
  _initAuth();
  
  // Blog
  _initBlog();
  
  // Projects
  _initProjects();
  
  // Frameworks
  _initFrameworks();
  
  // About
  _initAbout();
  
  // Education
  _initEducation();
  
  // Experience
  _initExperience();
  
  // Money
  _initMoney();
  
  // Testimonials
  _initTestimonials();
  
  // Settings
  _initSettings();
}

// Core BLoCs - Local UI State Only
void _initCoreBLoCs() {
  // Theme BLoC - manages dark/light mode (local state)
  sl.registerFactory(() => ThemeBloc(sharedPreferences: sl()));
  
  // Language BLoC - manages app language (local state)
  sl.registerFactory(() => LanguageBloc(sharedPreferences: sl()));
  
  // Navigation BLoC - manages bottom nav and routing (local state)
  sl.registerFactory(() => NavigationBloc());
}

// Auth
void _initAuth() {
  // Data sources
  sl.registerLazySingleton<AuthLocalDataSource>(
    () => AuthLocalDataSourceImpl(),
  );
  
  sl.registerLazySingleton<AuthRemoteDataSource>(
    () => AuthRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
    ),
  );

  // Repository
  sl.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
      authStore: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => LoginUseCase(sl()));
  sl.registerLazySingleton(() => RegisterUseCase(sl()));
  sl.registerLazySingleton(() => GetCurrentUserUseCase(sl()));
  sl.registerLazySingleton(() => LogoutUseCase(sl()));
  sl.registerLazySingleton(() => RefreshTokenUseCase(sl()));
  sl.registerLazySingleton(() => ChangePasswordUseCase(sl()));
  sl.registerLazySingleton(() => UpdateProfileUseCase(sl()));
  sl.registerLazySingleton(() => RestoreAuthStateUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => AuthBloc(
    loginUseCase: sl(),
    registerUseCase: sl(),
    getCurrentUserUseCase: sl(),
    logoutUseCase: sl(),
    refreshTokenUseCase: sl(),
    changePasswordUseCase: sl(),
    updateProfileUseCase: sl(),
    restoreAuthStateUseCase: sl(),
    authStore: sl(),
  ));
}

// Blog
void _initBlog() {
  // Data sources
  sl.registerLazySingleton<BlogLocalDataSource>(
    () => BlogLocalDataSourceImpl(authStore: sl<AuthStoreBloc>()),
  );
  
  sl.registerLazySingleton<BlogRemoteDataSource>(
    () => BlogRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );

  // Repository
  sl.registerLazySingleton<BlogRepository>(
    () => BlogRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetBlogPostsUseCase(sl()));
  sl.registerLazySingleton(() => GetBlogPostByIdUseCase(sl()));
  sl.registerLazySingleton(() => GetBlogPostBySlugUseCase(sl()));
  sl.registerLazySingleton(() => CreateBlogPostUseCase(sl()));
  sl.registerLazySingleton(() => UpdateBlogPostUseCase(sl()));
  sl.registerLazySingleton(() => DeleteBlogPostUseCase(sl()));
  sl.registerLazySingleton(() => IncrementViewCountUseCase(sl()));
  sl.registerLazySingleton(() => IncrementLikeCountUseCase(sl()));
  sl.registerLazySingleton(() => GetBlogTagsUseCase(sl()));
  sl.registerLazySingleton(() => CreateBlogTagUseCase(sl()));
  
  // BLoCs (local UI state only - NOT for server data)
  sl.registerFactory(() => BlogBloc(
    getBlogPostsUseCase: sl(),
  ));             // Main blog state management
  sl.registerFactory(() => CreateBlogPostBloc());   // Form validation & multi-step
  sl.registerFactory(() => BlogFilterBloc());       // Filters, sorting, view mode
}

// Projects
void _initProjects() {
  // Data sources
  sl.registerLazySingleton<ProjectsLocalDataSource>(
    () => ProjectsLocalDataSourceImpl(authStore: sl<AuthStoreBloc>()),
  );
  
  sl.registerLazySingleton<ProjectsRemoteDataSource>(
    () => ProjectsRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
    ),
  );


  // Repository
  sl.registerLazySingleton<ProjectsRepository>(
    () => ProjectsRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetProjectsUseCase(sl()));
  sl.registerLazySingleton(() => GetProjectByIdUseCase(sl()));
  sl.registerLazySingleton(() => GetProjectBySlugUseCase(sl()));
  sl.registerLazySingleton(() => CreateProjectUseCase(sl()));
  sl.registerLazySingleton(() => UpdateProjectUseCase(sl()));
  sl.registerLazySingleton(() => DeleteProjectUseCase(sl()));
  sl.registerLazySingleton(() => UpdateProjectOrderUseCase(sl()));
  sl.registerLazySingleton(() => IncrementProjectViewCountUseCase(sl()));
  sl.registerLazySingleton(() => IncrementProjectLikeCountUseCase(sl()));
  sl.registerLazySingleton(() => ToggleProjectFeaturedUseCase(sl()));
  sl.registerLazySingleton(() => ToggleProjectVisibilityUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => ProjectsBloc(
    getProjectsUseCase: sl(),
    getProjectByIdUseCase: sl(),
    getProjectBySlugUseCase: sl(),
    createProjectUseCase: sl(),
    updateProjectUseCase: sl(),
    deleteProjectUseCase: sl(),
    updateProjectOrderUseCase: sl(),
    incrementViewCountUseCase: sl(),
    incrementLikeCountUseCase: sl(),
    toggleProjectFeaturedUseCase: sl(),
    toggleProjectVisibilityUseCase: sl(),
  ));

}

// Frameworks
void _initFrameworks() {
  // Data sources
  sl.registerLazySingleton<FrameworksLocalDataSource>(
    () => FrameworksLocalDataSourceImpl(
      authStoreBloc: sl<AuthStoreBloc>(),
    ),
  );
  
  sl.registerLazySingleton<FrameworksRemoteDataSource>(
    () => FrameworksRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );


  // Repository
  sl.registerLazySingleton<FrameworksRepository>(
    () => FrameworksRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetFrameworksUseCase(sl()));
  sl.registerLazySingleton(() => CreateFrameworkUseCase(sl()));
  sl.registerLazySingleton(() => UpdateFrameworkUseCase(sl()));
  sl.registerLazySingleton(() => DeleteFrameworkUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => FrameworksBloc(
    getFrameworksUseCase: sl(),
    createFrameworkUseCase: sl(),
    updateFrameworkUseCase: sl(),
    deleteFrameworkUseCase: sl(),
  ));
}

// About
void _initAbout() {
  // Data sources
  sl.registerLazySingleton<AboutLocalDataSource>(
    () => AboutLocalDataSourceImpl(authStore: sl<AuthStoreBloc>()),
  );
  
  sl.registerLazySingleton<AboutRemoteDataSource>(
    () => AboutRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );


  // Repository
  sl.registerLazySingleton<AboutRepository>(
    () => AboutRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetAboutCoreUseCase(sl()));
  sl.registerLazySingleton(() => UpdateAboutCoreUseCase(sl()));
  sl.registerLazySingleton(() => GetBiographyUseCase(sl()));
  sl.registerLazySingleton(() => UpdateBiographyUseCase(sl()));
  sl.registerLazySingleton(() => GetAnimeListUseCase(sl()));
  sl.registerLazySingleton(() => GetAnimeByIdUseCase(sl()));
  sl.registerLazySingleton(() => CreateAnimeUseCase(sl()));
  sl.registerLazySingleton(() => UpdateAnimeUseCase(sl()));
  sl.registerLazySingleton(() => DeleteAnimeUseCase(sl()));
  sl.registerLazySingleton(() => GetMusicGenresUseCase(sl()));
  sl.registerLazySingleton(() => CreateMusicGenreUseCase(sl()));
  sl.registerLazySingleton(() => UpdateMusicGenreUseCase(sl()));
  sl.registerLazySingleton(() => DeleteMusicGenreUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => AboutBloc(
    getAboutCoreUseCase: sl(),
    updateAboutCoreUseCase: sl(),
    getBiographyUseCase: sl(),
    updateBiographyUseCase: sl(),
    getAnimeListUseCase: sl(),
    getAnimeByIdUseCase: sl(),
    createAnimeUseCase: sl(),
    updateAnimeUseCase: sl(),
    deleteAnimeUseCase: sl(),
    getMusicGenresUseCase: sl(),
    createMusicGenreUseCase: sl(),
    updateMusicGenreUseCase: sl(),
    deleteMusicGenreUseCase: sl(),
  ));

}


// Education
void _initEducation() {
  // Data sources
  sl.registerLazySingleton<EducationLocalDataSource>(
    () => EducationLocalDataSourceImpl(authStore: sl<AuthStoreBloc>()),
  );
  
  sl.registerLazySingleton<EducationRemoteDataSource>(
    () => EducationRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );


  // Repository
  sl.registerLazySingleton<EducationRepository>(
    () => EducationRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetEducationListUseCase(sl()));
  sl.registerLazySingleton(() => GetEducationByIdUseCase(sl()));
  sl.registerLazySingleton(() => CreateEducationUseCase(sl()));
  sl.registerLazySingleton(() => UpdateEducationUseCase(sl()));
  sl.registerLazySingleton(() => DeleteEducationUseCase(sl()));
  sl.registerLazySingleton(() => UpdateEducationOrderUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => EducationBloc(
    getEducationListUseCase: sl(),
    getEducationByIdUseCase: sl(),
    createEducationUseCase: sl(),
    updateEducationUseCase: sl(),
    deleteEducationUseCase: sl(),
    updateEducationOrderUseCase: sl(),
  ));
}

// Experience
void _initExperience() {
  // Data sources
  sl.registerLazySingleton<ExperienceLocalDataSource>(
    () => ExperienceLocalDataSourceImpl(authStore: sl<AuthStoreBloc>()),
  );
  
  sl.registerLazySingleton<ExperienceRemoteDataSource>(
    () => ExperienceRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );

  // Repository
  sl.registerLazySingleton<ExperienceRepository>(
    () => ExperienceRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetExperienceListUseCase(sl()));
  sl.registerLazySingleton(() => GetExperienceByIdUseCase(sl()));
  sl.registerLazySingleton(() => CreateExperienceUseCase(sl()));
  sl.registerLazySingleton(() => UpdateExperienceUseCase(sl()));
  sl.registerLazySingleton(() => DeleteExperienceUseCase(sl()));
  sl.registerLazySingleton(() => UpdateExperienceOrderUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => ExperienceBloc(
    getExperienceListUseCase: sl(),
    getExperienceByIdUseCase: sl(),
    createExperienceUseCase: sl(),
    updateExperienceUseCase: sl(),
    deleteExperienceUseCase: sl(),
    updateExperienceOrderUseCase: sl(),
  ));
}

// Money
void _initMoney() {
  // Data sources
  sl.registerLazySingleton<MoneyLocalDataSource>(
    () => MoneyLocalDataSourceImpl(authStore: sl<AuthStoreBloc>()),
  );
  
  sl.registerLazySingleton<MoneyRemoteDataSource>(
    () => MoneyRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );

  // Repository
  sl.registerLazySingleton<MoneyRepository>(
    () => MoneyRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetIdeasUseCase(sl()));
  sl.registerLazySingleton(() => GetIdeaByIdUseCase(sl()));
  sl.registerLazySingleton(() => GetIdeaBySlugUseCase(sl()));
  sl.registerLazySingleton(() => CreateIdeaUseCase(sl()));
  sl.registerLazySingleton(() => UpdateIdeaUseCase(sl()));
  sl.registerLazySingleton(() => DeleteIdeaUseCase(sl()));
  sl.registerLazySingleton(() => GetAiChatsUseCase(sl()));
  sl.registerLazySingleton(() => GetAiChatByIdUseCase(sl()));
  sl.registerLazySingleton(() => CreateAiChatUseCase(sl()));
  sl.registerLazySingleton(() => UpdateAiChatUseCase(sl()));
  sl.registerLazySingleton(() => DeleteAiChatUseCase(sl()));
  sl.registerLazySingleton(() => SendMessageUseCase(sl()));
  sl.registerLazySingleton(() => GetChatMessagesUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => MoneyBloc(
    getIdeasUseCase: sl(),
    getIdeaByIdUseCase: sl(),
    getIdeaBySlugUseCase: sl(),
    createIdeaUseCase: sl(),
    updateIdeaUseCase: sl(),
    deleteIdeaUseCase: sl(),
    getAiChatsUseCase: sl(),
    getAiChatByIdUseCase: sl(),
    createAiChatUseCase: sl(),
    updateAiChatUseCase: sl(),
    deleteAiChatUseCase: sl(),
    sendMessageUseCase: sl(),
    getChatMessagesUseCase: sl(),
  ));

}

// Testimonials
void _initTestimonials() {
  // Data sources
  sl.registerLazySingleton<TestimonialsLocalDataSource>(
    () => TestimonialsLocalDataSourceImpl(authStore: sl<AuthStoreBloc>()),
  );
  
  sl.registerLazySingleton<TestimonialsRemoteDataSource>(
    () => TestimonialsRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );

  // Repository
  sl.registerLazySingleton<TestimonialsRepository>(
    () => TestimonialsRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetTestimonialsUseCase(sl()));
  sl.registerLazySingleton(() => GetTestimonialByIdUseCase(sl()));
  sl.registerLazySingleton(() => CreateTestimonialUseCase(sl()));
  sl.registerLazySingleton(() => UpdateTestimonialUseCase(sl()));
  sl.registerLazySingleton(() => DeleteTestimonialUseCase(sl()));
  sl.registerLazySingleton(() => UpdateTestimonialOrderUseCase(sl()));
  sl.registerLazySingleton(() => IncrementTestimonialViewCountUseCase(sl()));
  sl.registerLazySingleton(() => IncrementTestimonialLikeCountUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => TestimonialsBloc(
    getTestimonialsUseCase: sl(),
    getTestimonialByIdUseCase: sl(),
    createTestimonialUseCase: sl(),
    updateTestimonialUseCase: sl(),
    deleteTestimonialUseCase: sl(),
    updateTestimonialOrderUseCase: sl(),
    incrementViewCountUseCase: sl(),
    incrementLikeCountUseCase: sl(),
  ));
}

// Settings
void _initSettings() {
  // Data sources
  sl.registerLazySingleton<SettingsLocalDataSource>(
    () => SettingsLocalDataSourceImpl(authStore: sl<AuthStoreBloc>()),
  );
  
  sl.registerLazySingleton<SettingsRemoteDataSource>(
    () => SettingsRemoteDataSourceImpl(
      apiClient: sl<ApiClient>(),
      baseUrl: EnvConfig.apiBaseUrl,
    ),
  );

  // Repository
  sl.registerLazySingleton<SettingsRepository>(
    () => SettingsRepositoryImpl(
      remoteDataSource: sl(),
      localDataSource: sl(),
    ),
  );

  // Use cases
  sl.registerLazySingleton(() => GetSettingsUseCase(sl()));
  sl.registerLazySingleton(() => GetSettingByKeyUseCase(sl()));
  sl.registerLazySingleton(() => UpdateSettingUseCase(sl()));
  sl.registerLazySingleton(() => UpdateSettingByKeyUseCase(sl()));
  sl.registerLazySingleton(() => UpdateSettingsBulkUseCase(sl()));
  sl.registerLazySingleton(() => GetCoreProfileSettingsUseCase(sl()));
  sl.registerLazySingleton(() => GetSocialMediaSettingsUseCase(sl()));
  sl.registerLazySingleton(() => GetContactSettingsUseCase(sl()));
  sl.registerLazySingleton(() => UpdateCoreProfileSettingsUseCase(sl()));
  sl.registerLazySingleton(() => UpdateSocialMediaSettingsUseCase(sl()));
  sl.registerLazySingleton(() => UpdateContactSettingsUseCase(sl()));

  // BLoCs (local UI state only)
  sl.registerFactory(() => SettingsBloc(
    getSettingsUseCase: sl(),
    getSettingByKeyUseCase: sl(),
    updateSettingUseCase: sl(),
    updateSettingByKeyUseCase: sl(),
    updateSettingsBulkUseCase: sl(),
    getCoreProfileSettingsUseCase: sl(),
    getSocialMediaSettingsUseCase: sl(),
    getContactSettingsUseCase: sl(),
    updateCoreProfileSettingsUseCase: sl(),
    updateSocialMediaSettingsUseCase: sl(),
    updateContactSettingsUseCase: sl(),
  ));
}