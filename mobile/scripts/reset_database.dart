#!/usr/bin/env dart

import 'dart:developer' as developer;
import 'dart:io';
import '../lib/core/database/database_helper.dart';

/// Simple command-line script to reset the database
/// Run with: dart scripts/reset_database.dart
void main() async {
  developer.log('🔄 Starting database reset...');
  
  try {
    final dbHelper = DatabaseHelper();
    
    // Reset the database
    await dbHelper.resetDatabase();
    
    developer.log('✅ Database reset completed successfully!');
    developer.log('📊 Database has been recreated with fresh schema');
    developer.log('🚀 You can now run your Flutter app');
    
  } catch (e) {
    developer.log('❌ Error resetting database: $e');
    exit(1);
  }
}
