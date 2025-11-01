import 'package:sqflite/sqflite.dart';
import '../../../../core/database/database_helper.dart';
import '../../../../core/database/sync_manager.dart';
import '../../../../core/stores/auth_store.dart';
import '../../domain/entities/about_core_entity.dart';
import '../../domain/entities/biography_entity.dart';
import '../../domain/entities/anime_entity.dart';
import '../../domain/entities/music_genre_entity.dart';
import '../models/about_core_model.dart';
import '../models/biography_model.dart';
import '../models/anime_model.dart';
import '../models/music_genre_model.dart';

abstract class AboutLocalDataSource {
  // About core methods
  Future<AboutCoreEntity?> getCachedAboutCore();
  Future<void> cacheAboutCore(AboutCoreEntity aboutCore);
  Future<void> updateCachedAboutCore(AboutCoreEntity aboutCore);

  // Biography methods
  Future<BiographyEntity?> getCachedBiography();
  Future<void> cacheBiography(BiographyEntity biography);
  Future<void> updateCachedBiography(BiographyEntity biography);

  // Anime methods
  Future<List<AnimeEntity>> getCachedAnimeList();
  Future<AnimeEntity?> getCachedAnime(String id);
  Future<void> cacheAnime(AnimeEntity anime);
  Future<void> cacheAnimeList(List<AnimeEntity> animeList);
  Future<void> updateCachedAnime(AnimeEntity anime);
  Future<void> removeCachedAnime(String id);

  // Music genre methods
  Future<List<MusicGenreEntity>> getCachedMusicGenres();
  Future<MusicGenreEntity?> getCachedMusicGenre(String id);
  Future<void> cacheMusicGenre(MusicGenreEntity musicGenre);
  Future<void> cacheMusicGenres(List<MusicGenreEntity> musicGenres);
  Future<void> updateCachedMusicGenre(MusicGenreEntity musicGenre);
  Future<void> removeCachedMusicGenre(String id);
}

class AboutLocalDataSourceImpl implements AboutLocalDataSource {
  final DatabaseHelper _dbHelper = DatabaseHelper();
  final SyncManager _syncManager = SyncManager();
  final AuthStoreBloc? _authStore;
  
  AboutLocalDataSourceImpl({AuthStoreBloc? authStore}) : _authStore = authStore;

  /// Ensure required NOT NULL fields are always present
  Map<String, dynamic> _ensureRequiredFields(Map<String, dynamic> data) {
    final safeData = Map<String, dynamic>.from(data);
    
    // ✅ ENSURE user_id is always present (NOT NULL constraint)
    if (!safeData.containsKey('user_id') || safeData['user_id'] == null || safeData['user_id'] == '') {
      // Try to get user ID from auth store
      final currentUserId = _authStore?.state.user?.id;
      if (currentUserId != null && currentUserId.isNotEmpty) {
        safeData['user_id'] = currentUserId;
      } else {
        // Fallback to a default user ID if not authenticated
        safeData['user_id'] = 'system';
      }
    }
    
    return safeData;
  }

  // About Core Implementation
  @override
  Future<AboutCoreEntity?> getCachedAboutCore() async {
    final db = await _dbHelper.database;
    final result = await db.query('about_core', limit: 1);
    return result.isEmpty ? null : AboutCoreModel.fromJson(result.first);
  }

  @override
  Future<void> cacheAboutCore(AboutCoreEntity aboutCore) async {
    final db = await _dbHelper.database;
    final aboutCoreMap = AboutCoreModel.fromEntity(aboutCore).toJson();
    aboutCoreMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    aboutCoreMap['is_dirty'] = 0;

    final safeAboutCoreMap = _ensureRequiredFields(aboutCoreMap);

    await db.insert(
      'about_core',
      safeAboutCoreMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> updateCachedAboutCore(AboutCoreEntity aboutCore) async {
    final db = await _dbHelper.database;
    final aboutCoreMap = AboutCoreModel.fromEntity(aboutCore).toJson();
    aboutCoreMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    aboutCoreMap['is_dirty'] = 1;

    final safeAboutCoreMap = _ensureRequiredFields(aboutCoreMap);

    await db.update(
      'about_core',
      safeAboutCoreMap,
      where: 'id = ?',
      whereArgs: [aboutCore.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'about_core',
      recordId: aboutCore.id,
      operation: SyncOperation.update,
      data: aboutCoreMap,
    );
  }

  // Biography Implementation
  @override
  Future<BiographyEntity?> getCachedBiography() async {
    final db = await _dbHelper.database;
    final result = await db.query('biographies', limit: 1);
    return result.isEmpty ? null : BiographyModel.fromJson(result.first);
  }

  @override
  Future<void> cacheBiography(BiographyEntity biography) async {
    final db = await _dbHelper.database;
    final biographyMap = BiographyModel.fromEntity(biography).toJson();
    biographyMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    biographyMap['is_dirty'] = 0;

    await db.insert(
      'biographies',
      biographyMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> updateCachedBiography(BiographyEntity biography) async {
    final db = await _dbHelper.database;
    final biographyMap = BiographyModel.fromEntity(biography).toJson();
    biographyMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    biographyMap['is_dirty'] = 1;

    await db.update(
      'biographies',
      biographyMap,
      where: 'id = ?',
      whereArgs: [biography.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'biographies',
      recordId: biography.id,
      operation: SyncOperation.update,
      data: biographyMap,
    );
  }

  // Anime Implementation
  @override
  Future<List<AnimeEntity>> getCachedAnimeList() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'anime_list',
      orderBy: '`order` ASC, title ASC',
    );
    return result.map((animeMap) => AnimeModel.fromJson(animeMap)).toList();
  }

  @override
  Future<AnimeEntity?> getCachedAnime(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'anime_list',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : AnimeModel.fromJson(result.first);
  }

  @override
  Future<void> cacheAnime(AnimeEntity anime) async {
    final db = await _dbHelper.database;
    final animeMap = AnimeModel.fromEntity(anime).toJson();
    animeMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    animeMap['is_dirty'] = 0;

    await db.insert(
      'anime_list',
      animeMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheAnimeList(List<AnimeEntity> animeList) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final anime in animeList) {
      final animeMap = AnimeModel.fromEntity(anime).toJson();
      animeMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      animeMap['is_dirty'] = 0;

      batch.insert(
        'anime_list',
        animeMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedAnime(AnimeEntity anime) async {
    final db = await _dbHelper.database;
    final animeMap = AnimeModel.fromEntity(anime).toJson();
    animeMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    animeMap['is_dirty'] = 1;

    await db.update(
      'anime_list',
      animeMap,
      where: 'id = ?',
      whereArgs: [anime.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'anime_list',
      recordId: anime.id,
      operation: SyncOperation.update,
      data: animeMap,
    );
  }

  @override
  Future<void> removeCachedAnime(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'anime_list',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'anime_list',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  // Music Genre Implementation
  @override
  Future<List<MusicGenreEntity>> getCachedMusicGenres() async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'music_genres',
      orderBy: '`order` ASC, name ASC',
    );
    return result.map((genreMap) => MusicGenreModel.fromJson(genreMap)).toList();
  }

  @override
  Future<MusicGenreEntity?> getCachedMusicGenre(String id) async {
    final db = await _dbHelper.database;
    final result = await db.query(
      'music_genres',
      where: 'id = ?',
      whereArgs: [id],
    );
    return result.isEmpty ? null : MusicGenreModel.fromJson(result.first);
  }

  @override
  Future<void> cacheMusicGenre(MusicGenreEntity musicGenre) async {
    final db = await _dbHelper.database;
      final genreMap = MusicGenreModel.fromEntity(musicGenre).toJson();
    genreMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
    genreMap['is_dirty'] = 0;

    await db.insert(
      'music_genres',
      genreMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
  }

  @override
  Future<void> cacheMusicGenres(List<MusicGenreEntity> musicGenres) async {
    final db = await _dbHelper.database;
    final batch = db.batch();

    for (final genre in musicGenres) {
      final genreMap = MusicGenreModel.fromEntity(genre).toJson();
      genreMap['synced_at'] = DateTime.now().millisecondsSinceEpoch;
      genreMap['is_dirty'] = 0;

      batch.insert(
        'music_genres',
        genreMap,
        conflictAlgorithm: ConflictAlgorithm.replace,
      );
    }

    await batch.commit();
  }

  @override
  Future<void> updateCachedMusicGenre(MusicGenreEntity musicGenre) async {
    final db = await _dbHelper.database;
      final genreMap = MusicGenreModel.fromEntity(musicGenre).toJson();
    genreMap['updated_at'] = DateTime.now().millisecondsSinceEpoch;
    genreMap['is_dirty'] = 1;

    await db.update(
      'music_genres',
      genreMap,
      where: 'id = ?',
      whereArgs: [musicGenre.id],
    );

    await _syncManager.addToSyncQueue(
      tableName: 'music_genres',
      recordId: musicGenre.id,
      operation: SyncOperation.update,
      data: genreMap,
    );
  }

  @override
  Future<void> removeCachedMusicGenre(String id) async {
    final db = await _dbHelper.database;
    
    await _syncManager.addToSyncQueue(
      tableName: 'music_genres',
      recordId: id,
      operation: SyncOperation.delete,
    );

    await db.delete(
      'music_genres',
      where: 'id = ?',
      whereArgs: [id],
    );
  }
}
