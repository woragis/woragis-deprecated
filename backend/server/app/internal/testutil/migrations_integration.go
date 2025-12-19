//go:build integration && !performance_test
// +build integration,!performance_test

package testutil

import (
	"testing"

	"gorm.io/gorm"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	casestudiesdomain "github.com/woragis/backend/server/app/internal/domains/casestudies"
	certificationsdomain "github.com/woragis/backend/server/app/internal/domains/certifications"
	creativeassetsdomain "github.com/woragis/backend/server/app/internal/domains/creativeassets"
	experiencesdomain "github.com/woragis/backend/server/app/internal/domains/experiences"
	interestsdomain "github.com/woragis/backend/server/app/internal/domains/interests"
	postcommentsdomain "github.com/woragis/backend/server/app/internal/domains/posts/comments"
	postsdomain "github.com/woragis/backend/server/app/internal/domains/posts"
	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	skillsdomain "github.com/woragis/backend/server/app/internal/domains/skills"
	testimonialsdomain "github.com/woragis/backend/server/app/internal/domains/testimonials"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	userpreferencesdomain "github.com/woragis/backend/server/app/internal/domains/userpreferences"
	userprofilesdomain "github.com/woragis/backend/server/app/internal/domains/userprofiles"
)

// MigrateTestDB runs database migrations for integration tests
// Excludes socialmediaposts domain to avoid import cycles
// Accepts testing.TB to support both *testing.T and *testing.B
func MigrateTestDB(t testing.TB, db *gorm.DB) error {
	err := db.AutoMigrate(
		&authdomain.User{},
		&userprofilesdomain.UserProfile{},
		&userpreferencesdomain.UserPreferences{},
		&postsdomain.Post{},
		&postcommentsdomain.Comment{},
		&projectsdomain.Project{},
		&interestsdomain.Interest{},
		&experiencesdomain.Experience{},
		&casestudiesdomain.CaseStudy{},
		&certificationsdomain.Certification{},
		&skillsdomain.Skill{},
		&testimonialsdomain.Testimonial{},
		&creativeassetsdomain.CreativeAsset{},
		&apikeysdomain.APIKey{},
		&translationsdomain.Translation{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
		return err
	}
	return nil
}
