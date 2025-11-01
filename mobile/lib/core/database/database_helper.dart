import 'dart:async';
import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart';

class DatabaseHelper {
  static final DatabaseHelper _instance = DatabaseHelper._internal();
  static Database? _database;

  DatabaseHelper._internal();

  factory DatabaseHelper() => _instance;

  Future<Database> get database async {
    if (_database != null) return _database!;
    _database = await _initDatabase();
    return _database!;
  }

  Future<Database> _initDatabase() async {
    String path = join(await getDatabasesPath(), 'woragis.db');
    return await openDatabase(
      path,
      version: 2, // Increment version for schema migration
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
      onConfigure: _onConfigure,
    );
  }

  Future<void> _onConfigure(Database db) async {
    await db.execute('PRAGMA foreign_keys = ON');
    await db.execute('PRAGMA journal_mode = WAL');
  }

  Future<void> _onCreate(Database db, int version) async {
    // Create all tables
    await _createAuthTables(db);
    await _createBlogTables(db);
    await _createProjectsTables(db);
    await _createFrameworksTables(db);
    await _createAboutTables(db);
    await _createMusicTables(db);
    await _createEducationTables(db);
    await _createExperienceTables(db);
    await _createMoneyTables(db);
    await _createTestimonialsTables(db);
    await _createSettingsTables(db);
    await _createSyncTables(db);
  }

  Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    // Handle database migrations here
    if (oldVersion < 2) {
      // Migration for frameworks table schema update
      await _migrateFrameworksTable(db);
    }
  }

  Future<void> _migrateFrameworksTable(Database db) async {
    try {
      // Check if we need to migrate the frameworks table
      final tableInfo = await db.rawQuery('PRAGMA table_info(frameworks)');
      final columnNames = tableInfo.map((column) => column['name'] as String).toList();
      
      // Add missing columns if they don't exist
      if (!columnNames.contains('proficiencyLevel')) {
        await db.execute('ALTER TABLE frameworks ADD COLUMN proficiencyLevel TEXT');
      }
      if (!columnNames.contains('public')) {
        await db.execute('ALTER TABLE frameworks ADD COLUMN public INTEGER NOT NULL DEFAULT 1');
      }
      
      // Rename columns if they exist with old names
      if (columnNames.contains('user_id') && !columnNames.contains('userId')) {
        // SQLite doesn't support RENAME COLUMN in older versions, so we'll recreate
        await db.execute('''
          CREATE TABLE frameworks_new (
            id TEXT PRIMARY KEY,
            userId TEXT NOT NULL,
            name TEXT NOT NULL UNIQUE,
            slug TEXT NOT NULL UNIQUE,
            description TEXT,
            icon TEXT,
            color TEXT,
            website TEXT,
            type TEXT NOT NULL DEFAULT 'framework',
            proficiencyLevel TEXT,
            version TEXT,
            `order` INTEGER DEFAULT 0,
            visible INTEGER NOT NULL DEFAULT 1,
            public INTEGER NOT NULL DEFAULT 1,
            createdAt INTEGER NOT NULL,
            updatedAt INTEGER NOT NULL,
            synced_at INTEGER,
            is_dirty INTEGER NOT NULL DEFAULT 0
          )
        ''');
        
        await db.execute('''
          INSERT INTO frameworks_new (
            id, userId, name, slug, description, icon, color, website, type, 
            proficiencyLevel, version, `order`, visible, public, createdAt, updatedAt, synced_at, is_dirty
          )
          SELECT 
            id, user_id, name, slug, description, icon, color, website, type,
            proficiencyLevel, version, `order`, visible, 1, created_at, updated_at, synced_at, is_dirty
          FROM frameworks
        ''');
        
        await db.execute('DROP TABLE frameworks');
        await db.execute('ALTER TABLE frameworks_new RENAME TO frameworks');
      }
    } catch (e) {
      // If migration fails, recreate the table
      await db.execute('DROP TABLE IF EXISTS frameworks');
      await _createFrameworksTables(db);
    }
  }

  // Auth tables
  Future<void> _createAuthTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT,
        first_name TEXT,
        last_name TEXT,
        avatar TEXT,
        role TEXT NOT NULL DEFAULT 'user',
        is_active INTEGER NOT NULL DEFAULT 1,
        email_verified INTEGER NOT NULL DEFAULT 0,
        last_login_at INTEGER,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS user_sessions (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        access_token TEXT,
        refresh_token TEXT,
        token_hash TEXT UNIQUE,
        expires_at INTEGER NOT NULL,
        is_active INTEGER NOT NULL DEFAULT 1,
        user_agent TEXT,
        ip_address TEXT,
        created_at INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
      )
    ''');
  }

  // Blog tables
  Future<void> _createBlogTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS blog_posts (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        slug TEXT NOT NULL UNIQUE,
        excerpt TEXT NOT NULL,
        content TEXT NOT NULL,
        featured_image TEXT,
        reading_time INTEGER,
        featured INTEGER NOT NULL DEFAULT 0,
        published INTEGER NOT NULL DEFAULT 0,
        published_at INTEGER,
        `order` INTEGER NOT NULL DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        public INTEGER NOT NULL DEFAULT 1,
        view_count INTEGER DEFAULT 0,
        like_count INTEGER DEFAULT 0,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS blog_tags (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        name TEXT NOT NULL,
        slug TEXT NOT NULL UNIQUE,
        description TEXT,
        color TEXT,
        visible INTEGER NOT NULL DEFAULT 1,
        `order` INTEGER NOT NULL DEFAULT 0,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS blog_post_tags (
        blog_post_id TEXT NOT NULL,
        tag_id TEXT NOT NULL,
        created_at INTEGER NOT NULL,
        FOREIGN KEY (blog_post_id) REFERENCES blog_posts (id) ON DELETE CASCADE,
        FOREIGN KEY (tag_id) REFERENCES blog_tags (id) ON DELETE CASCADE,
        PRIMARY KEY (blog_post_id, tag_id)
      )
    ''');
  }

  // Projects tables
  Future<void> _createProjectsTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS projects (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        slug TEXT NOT NULL,
        description TEXT NOT NULL,
        long_description TEXT,
        content TEXT,
        video_url TEXT,
        image TEXT NOT NULL,
        github_url TEXT,
        live_url TEXT,
        featured INTEGER NOT NULL DEFAULT 1,
        `order` INTEGER NOT NULL DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        public INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS project_frameworks (
        id TEXT PRIMARY KEY,
        project_id TEXT NOT NULL,
        framework_id TEXT NOT NULL,
        version TEXT,
        proficiency TEXT,
        created_at INTEGER NOT NULL,
        FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
        FOREIGN KEY (framework_id) REFERENCES frameworks (id) ON DELETE CASCADE,
        UNIQUE(project_id, framework_id)
      )
    ''');
  }

  // Frameworks tables
  Future<void> _createFrameworksTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS frameworks (
        id TEXT PRIMARY KEY,
        userId TEXT NOT NULL,
        name TEXT NOT NULL UNIQUE,
        slug TEXT NOT NULL UNIQUE,
        description TEXT,
        icon TEXT,
        color TEXT,
        website TEXT,
        type TEXT NOT NULL DEFAULT 'framework',
        proficiencyLevel TEXT,
        version TEXT,
        `order` INTEGER DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        public INTEGER NOT NULL DEFAULT 1,
        createdAt INTEGER NOT NULL,
        updatedAt INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');
  }

  // About tables
  Future<void> _createAboutTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS biography (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        featured_biography TEXT,
        full_biography TEXT,
        visible INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS anime_list (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        status TEXT NOT NULL DEFAULT 'want_to_watch',
        myanimelist_id TEXT,
        cover_image TEXT,
        genres TEXT, -- JSON array
        episodes INTEGER,
        current_episode INTEGER DEFAULT 0,
        rating INTEGER,
        notes TEXT,
        started_at INTEGER,
        completed_at INTEGER,
        `order` INTEGER DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS music_genres (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        name TEXT NOT NULL,
        description TEXT,
        `order` INTEGER DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS last_listened_songs (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        artist TEXT NOT NULL,
        album TEXT,
        spotify_url TEXT,
        youtube_url TEXT,
        listened_at INTEGER,
        `order` INTEGER DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');
  }

  // Education tables
  Future<void> _createEducationTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS education (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        institution TEXT NOT NULL,
        description TEXT,
        type TEXT NOT NULL,
        degree_level TEXT,
        field_of_study TEXT,
        start_date INTEGER,
        end_date INTEGER,
        completion_date INTEGER,
        grade TEXT,
        credits INTEGER,
        certificate_id TEXT,
        issuer TEXT,
        validity_period TEXT,
        pdf_document TEXT,
        verification_url TEXT,
        skills TEXT, -- JSON array
        `order` INTEGER DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
      )
    ''');
  }

  // Experience tables
  Future<void> _createExperienceTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS experiences (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        company TEXT NOT NULL,
        period TEXT NOT NULL,
        location TEXT NOT NULL,
        description TEXT NOT NULL,
        achievements TEXT, -- JSON array
        technologies TEXT, -- JSON array
        icon TEXT NOT NULL DEFAULT '💼',
        `order` INTEGER NOT NULL DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');
  }

  // Money tables
  Future<void> _createMoneyTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS ideas (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        slug TEXT NOT NULL UNIQUE,
        document TEXT NOT NULL,
        description TEXT,
        featured INTEGER NOT NULL DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        public INTEGER NOT NULL DEFAULT 1,
        `order` INTEGER NOT NULL DEFAULT 0,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS idea_nodes (
        id TEXT PRIMARY KEY,
        idea_id TEXT NOT NULL,
        title TEXT NOT NULL,
        content TEXT,
        type TEXT DEFAULT 'default',
        position_x REAL NOT NULL DEFAULT 0,
        position_y REAL NOT NULL DEFAULT 0,
        width REAL DEFAULT 200,
        height REAL DEFAULT 100,
        color TEXT,
        connections TEXT, -- JSON array
        visible INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS ai_chats (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        idea_node_id TEXT NOT NULL,
        title TEXT NOT NULL,
        messages TEXT NOT NULL, -- JSON array
        agent TEXT NOT NULL DEFAULT 'gpt-4',
        model TEXT,
        system_prompt TEXT,
        temperature TEXT DEFAULT '0.7',
        max_tokens TEXT,
        visible INTEGER NOT NULL DEFAULT 1,
        archived INTEGER NOT NULL DEFAULT 0,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');
  }

  // Testimonials tables
  Future<void> _createTestimonialsTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS testimonials (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        name TEXT NOT NULL,
        position TEXT NOT NULL,
        company TEXT NOT NULL,
        content TEXT NOT NULL,
        avatar TEXT,
        rating INTEGER NOT NULL DEFAULT 5,
        featured INTEGER NOT NULL DEFAULT 1,
        `order` INTEGER NOT NULL DEFAULT 0,
        visible INTEGER NOT NULL DEFAULT 1,
        public INTEGER NOT NULL DEFAULT 1,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL,
        synced_at INTEGER,
        is_dirty INTEGER NOT NULL DEFAULT 0
      )
    ''');
  }

  // Music tables
  Future<void> _createMusicTables(Database db) async {
    // Note: last_listened_songs table is already created in _createAboutTables
    // No additional music tables needed here
  }

  // Settings tables
  Future<void> _createSettingsTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS settings (
        id TEXT PRIMARY KEY,
        key TEXT NOT NULL UNIQUE,
        value TEXT NOT NULL,
        updated_at INTEGER NOT NULL
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS social_media (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        github TEXT,
        linkedin TEXT,
        twitter TEXT,
        instagram TEXT,
        youtube TEXT,
        discord TEXT,
        telegram TEXT,
        website TEXT,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS contact_info (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        email TEXT NOT NULL,
        phone TEXT,
        address TEXT,
        city TEXT,
        country TEXT,
        timezone TEXT,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS site_settings (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        theme TEXT DEFAULT 'system' NOT NULL,
        language TEXT DEFAULT 'en' NOT NULL,
        projects_per_page INTEGER DEFAULT 6 NOT NULL,
        maintenance_mode INTEGER DEFAULT 0 NOT NULL,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
      )
    ''');
  }

  // Sync tables
  Future<void> _createSyncTables(Database db) async {
    await db.execute('''
      CREATE TABLE IF NOT EXISTS sync_queue (
        id TEXT PRIMARY KEY,
        table_name TEXT NOT NULL,
        record_id TEXT NOT NULL,
        operation TEXT NOT NULL, -- CREATE, UPDATE, DELETE
        data TEXT, -- JSON data for the operation
        retry_count INTEGER NOT NULL DEFAULT 0,
        max_retries INTEGER NOT NULL DEFAULT 3,
        created_at INTEGER NOT NULL,
        last_retry_at INTEGER,
        status TEXT NOT NULL DEFAULT 'PENDING' -- PENDING, SYNCING, COMPLETED, FAILED
      )
    ''');

    await db.execute('''
      CREATE TABLE IF NOT EXISTS sync_status (
        id TEXT PRIMARY KEY,
        table_name TEXT NOT NULL UNIQUE,
        last_sync_at INTEGER,
        sync_token TEXT,
        created_at INTEGER NOT NULL,
        updated_at INTEGER NOT NULL
      )
    ''');
  }

  // Utility methods
  Future<void> clearAllData() async {
    final db = await database;
    await db.delete('users');
    await db.delete('blog_posts');
    await db.delete('blog_tags');
    await db.delete('projects');
    await db.delete('frameworks');
    await db.delete('biography');
    await db.delete('anime_list');
    await db.delete('music_genres');
    await db.delete('education');
    await db.delete('experiences');
    await db.delete('ideas');
    await db.delete('idea_nodes');
    await db.delete('ai_chats');
    await db.delete('testimonials');
    await db.delete('settings');
    await db.delete('social_media');
    await db.delete('contact_info');
    await db.delete('site_settings');
    await db.delete('last_listened_songs');
    await db.delete('sync_queue');
    await db.delete('sync_status');
  }

  /// Completely resets the database by deleting the file and recreating it
  Future<void> resetDatabase() async {
    // Close the current database connection
    await close();
    
    // Get the database path
    String path = join(await getDatabasesPath(), 'woragis.db');
    
    // Delete the database file
    await deleteDatabase(path);
    
    // Reset the database instance
    _database = null;
    
    // Recreate the database
    await database;
  }

  Future<void> close() async {
    final db = _database;
    if (db != null) {
      await db.close();
      _database = null;
    }
  }
}
