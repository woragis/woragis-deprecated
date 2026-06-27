package skills

import (
	"gorm.io/gorm"

	"woragis-skills-service/internal/domains/skills"
	"woragis-skills-service/internal/domains/interests"
)

// MigrateSkillsTables runs database migrations for skills service
func MigrateSkillsTables(db *gorm.DB) error {
	// Enable UUID extension if not already enabled
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Enable gen_random_uuid function if not already available
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		return err
	}

	// Migrate skills tables
	if err := db.AutoMigrate(
		&skills.Skill{},
		&skills.ProjectSkill{},
	); err != nil {
		return err
	}

	// Migrate interests tables
	if err := db.AutoMigrate(
		&interests.Interest{},
	); err != nil {
		return err
	}

	return nil
}
