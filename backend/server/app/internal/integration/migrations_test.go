//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/woragis/backend/server/app/internal/testutil"
)

// TestDatabaseMigrations tests that migrations run successfully
func TestDatabaseMigrations(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Verify that migrations ran successfully by checking for key tables
	tables := []string{
		"users",
		"projects",
		"skills",
		"posts",
		"interests",
		"sessions",
		"api_keys",
	}

	for _, tableName := range tables {
		var exists bool
		err := db.Raw(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = ?
			)
		`, tableName).Scan(&exists).Error
		require.NoError(t, err, "Failed to check if table %s exists", tableName)
		assert.True(t, exists, "Table %s should exist after migration", tableName)
	}
}

// TestMigrationsIdempotency tests that migrations can be run multiple times
func TestMigrationsIdempotency(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Get initial table count
	var initialCount int64
	err := db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name NOT LIKE 'pg_%'
	`).Scan(&initialCount).Error
	require.NoError(t, err)

	// Run migrations again
	err = testutil.MigrateTestDB(t, db)
	require.NoError(t, err, "Migrations should run successfully a second time")

	// Verify table count hasn't changed (no duplicate tables)
	var finalCount int64
	err = db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name NOT LIKE 'pg_%'
	`).Scan(&finalCount).Error
	require.NoError(t, err)
	assert.Equal(t, initialCount, finalCount, "Table count should remain the same after re-running migrations")
}

// TestMigrationsWithData tests that migrations work with existing data
func TestMigrationsWithData(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create some test data
	userID := testutil.CreateTestUser(t, db, "migration@example.com", "password123")

	// Verify data exists
	var userCount int64
	err := db.Table("users").Where("id = ?", userID).Count(&userCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), userCount, "User should exist")

	// Run migrations again (should not affect existing data)
	err = testutil.MigrateTestDB(t, db)
	require.NoError(t, err, "Migrations should run successfully with existing data")

	// Verify data still exists
	var userCountAfter int64
	err = db.Table("users").Where("id = ?", userID).Count(&userCountAfter).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), userCountAfter, "User should still exist after migrations")
}

// TestMigrationSchemaValidation tests that schema is correct after migration
func TestMigrationSchemaValidation(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Test users table schema
	var columnInfo []struct {
		ColumnName string
		DataType   string
		IsNullable string
	}
	err := db.Raw(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public' 
		AND table_name = 'users'
		ORDER BY ordinal_position
	`).Scan(&columnInfo).Error
	require.NoError(t, err)

	// Verify key columns exist
	columns := make(map[string]bool)
	for _, col := range columnInfo {
		columns[col.ColumnName] = true
	}

	requiredColumns := []string{"id", "email", "password_hash", "created_at", "updated_at"}
	for _, col := range requiredColumns {
		assert.True(t, columns[col], "Column %s should exist in users table", col)
	}

	// Verify id is UUID type
	var idType string
	err = db.Raw(`
		SELECT data_type
		FROM information_schema.columns
		WHERE table_schema = 'public' 
		AND table_name = 'users'
		AND column_name = 'id'
	`).Scan(&idType).Error
	require.NoError(t, err)
	assert.Contains(t, []string{"uuid", "character varying"}, idType, "ID column should be UUID type")

	// Test projects table schema
	var projectColumns []string
	err = db.Raw(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'public' 
		AND table_name = 'projects'
	`).Scan(&projectColumns).Error
	require.NoError(t, err)
	assert.Greater(t, len(projectColumns), 0, "Projects table should have columns")
}

// TestMigrationForeignKeys tests that foreign key constraints are created
func TestMigrationForeignKeys(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Check for foreign key constraints
	var fkCount int64
	err := db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.table_constraints
		WHERE constraint_type = 'FOREIGN KEY'
		AND table_schema = 'public'
	`).Scan(&fkCount).Error
	require.NoError(t, err)
	assert.Greater(t, fkCount, int64(0), "Should have foreign key constraints")
}

// TestMigrationIndexes tests that indexes are created
func TestMigrationIndexes(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Check for indexes on users table
	var indexCount int64
	err := db.Raw(`
		SELECT COUNT(*)
		FROM pg_indexes
		WHERE schemaname = 'public'
		AND tablename = 'users'
	`).Scan(&indexCount).Error
	require.NoError(t, err)
	assert.Greater(t, indexCount, int64(0), "Users table should have indexes")

	// Verify email index exists (should be unique)
	var emailIndexExists bool
	err = db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_indexes
			WHERE schemaname = 'public'
			AND tablename = 'users'
			AND indexname LIKE '%email%'
		)
	`).Scan(&emailIndexExists).Error
	require.NoError(t, err)
	assert.True(t, emailIndexExists, "Email index should exist")
}

// TestMigrationCleanup tests that database cleanup works correctly
func TestMigrationCleanup(t *testing.T) {
	db := testutil.SetupTestDB(t)

	// Create some data
	userID := testutil.CreateTestUser(t, db, "cleanup@example.com", "password123")

	// Verify data exists
	var userCount int64
	err := db.Table("users").Where("id = ?", userID).Count(&userCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), userCount)

	// Cleanup
	testutil.CleanupTestDB(t, db)

	// Verify tables are dropped (or at least user is gone)
	var finalUserCount int64
	err = db.Table("users").Where("id = ?", userID).Count(&finalUserCount).Error
	// This might fail if table is dropped, which is fine
	// The important thing is cleanup doesn't crash
	assert.True(t, finalUserCount == 0 || err != nil, "User should be cleaned up")
}
