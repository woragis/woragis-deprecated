//go:build performance_test
// +build performance_test

package testutil

import (
	"testing"

	"gorm.io/gorm"

	aimlintegrationsdomain "github.com/woragis/backend/server/app/internal/domains/aimlintegrations"
	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	casestudiesdomain "github.com/woragis/backend/server/app/internal/domains/casestudies"
	certificationsdomain "github.com/woragis/backend/server/app/internal/domains/certifications"
	clientsdomain "github.com/woragis/backend/server/app/internal/domains/clients"
	creativeassetsdomain "github.com/woragis/backend/server/app/internal/domains/creativeassets"
	experiencesdomain "github.com/woragis/backend/server/app/internal/domains/experiences"
	ideasdomain "github.com/woragis/backend/server/app/internal/domains/ideas"
	impactmetricsdomain "github.com/woragis/backend/server/app/internal/domains/impactmetrics"
	interestsdomain "github.com/woragis/backend/server/app/internal/domains/interests"
	jobapplicationsdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications"
	jobapplicationstagesdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications/interviewstages"
	jobwebsitesdomain "github.com/woragis/backend/server/app/internal/domains/jobwebsites"
	postcommentsdomain "github.com/woragis/backend/server/app/internal/domains/posts/comments"
	postsdomain "github.com/woragis/backend/server/app/internal/domains/posts"
	problemsolutionsdomain "github.com/woragis/backend/server/app/internal/domains/problemsolutions"
	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	resumesdomain "github.com/woragis/backend/server/app/internal/domains/resumes"
	skillsdomain "github.com/woragis/backend/server/app/internal/domains/skills"
	systemdesignsdomain "github.com/woragis/backend/server/app/internal/domains/systemdesigns"
	technicalwritingsdomain "github.com/woragis/backend/server/app/internal/domains/technicalwritings"
	testimonialsdomain "github.com/woragis/backend/server/app/internal/domains/testimonials"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	userprofilesdomain "github.com/woragis/backend/server/app/internal/domains/userprofiles"
	userpreferencesdomain "github.com/woragis/backend/server/app/internal/domains/userpreferences"
)

// MigrateTestDB runs database migrations for all domains except socialmediaposts
// This version is used for performance tests to avoid import cycles
// Accepts testing.TB to support both *testing.T and *testing.B
func MigrateTestDB(t testing.TB, db *gorm.DB) error {
	// Migrate all domains except socialmediaposts (to avoid import cycle)
	// Only migrate essential domains for performance tests
	// Some domains may not have entities or may cause issues, so we keep it minimal
	domains := []interface{}{
		// Auth and user domains
		&authdomain.User{},
		&userprofilesdomain.UserProfile{},
		&userpreferencesdomain.UserPreferences{}, // Note: plural name
		
		// Content domains
		&postsdomain.Post{},
		&postcommentsdomain.Comment{}, // Note: just "Comment", not "PostComment"
		&projectsdomain.Project{},
		&skillsdomain.Skill{},
		&interestsdomain.Interest{},
		&experiencesdomain.Experience{},
		&certificationsdomain.Certification{},
		&casestudiesdomain.CaseStudy{},
		&testimonialsdomain.Testimonial{},
		&creativeassetsdomain.CreativeAsset{},
		
		// Job application domains
		&jobapplicationsdomain.JobApplication{},
		&jobapplicationstagesdomain.InterviewStage{},
		&jobwebsitesdomain.JobWebsite{},
		
		// Other essential domains
		&apikeysdomain.APIKey{},
		&translationsdomain.Translation{},
		&resumesdomain.Resume{},
		&clientsdomain.Client{},
		&ideasdomain.Idea{},
		&problemsolutionsdomain.ProblemSolution{},
		&systemdesignsdomain.SystemDesign{},
		&technicalwritingsdomain.TechnicalWriting{},
		&aimlintegrationsdomain.AIMLIntegration{},
		&impactmetricsdomain.ImpactMetric{},
		
		// Note: socialmediaposts excluded to avoid import cycle
		// Note: Some domains excluded if their entity types don't exist or cause compilation issues
		// Note: jobapplicationresponses, languages, finances, reports, chats, scheduler excluded for now
	}
	
	if err := db.AutoMigrate(domains...); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
		return err
	}
	
	return nil
}
