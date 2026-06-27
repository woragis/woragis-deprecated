package files

import (
	"gorm.io/gorm"

	"woragis-files-service/internal/domains/files"
)

// MigrateFilesTables runs database migrations for files service
func MigrateFilesTables(db *gorm.DB) error {
	// Enable UUID extension if not already enabled
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Enable gen_random_uuid function if not already available
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		return err
	}

	// Migrate files tables
	if err := db.AutoMigrate(
		&files.File{},
	); err != nil {
		return err
	}

	return nil
}
