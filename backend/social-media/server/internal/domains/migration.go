package socialmedia

import (
	"gorm.io/gorm"

	"woragis-social-media-service/internal/domains/socialmediaposts"
	"woragis-social-media-service/internal/domains/socialmediaposts/analytics"
	"woragis-social-media-service/internal/domains/socialmediaposts/assets"
	"woragis-social-media-service/internal/domains/socialmediaposts/content"
	"woragis-social-media-service/internal/domains/socialmediaposts/platforms"
	"woragis-social-media-service/internal/domains/socialmediaposts/scheduling"
	"woragis-social-media-service/internal/domains/creativeassets"
)

// MigrateSocialMediaTables runs database migrations for social media service
func MigrateSocialMediaTables(db *gorm.DB) error {
	// Enable UUID extension if not already enabled
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Enable gen_random_uuid function if not already available
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		return err
	}

	// Migrate social media posts tables
	if err := db.AutoMigrate(
		&socialmediaposts.SocialMediaPost{},
	); err != nil {
		return err
	}

	// Migrate subdomain tables for social media posts
	if err := db.AutoMigrate(
		&analytics.PostAnalytics{},
		&assets.ContentAsset{},
		&content.ContentPost{},
		&content.ContentRepurposing{},
		&platforms.PlatformConfig{},
		&scheduling.ScheduledPost{},
	); err != nil {
		return err
	}

	// Migrate creative assets tables
	if err := db.AutoMigrate(
		&creativeassets.CreativeAsset{},
	); err != nil {
		return err
	}

	return nil
}
