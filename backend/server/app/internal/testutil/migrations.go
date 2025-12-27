//go:build !performance_test && !integration
// +build !performance_test,!integration

package testutil

import (
	"testing"

	"gorm.io/gorm"

	aimlintegrationsdomain "github.com/woragis/backend/server/app/internal/domains/aimlintegrations"
	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	casestudiesdomain "github.com/woragis/backend/server/app/internal/domains/casestudies"
	certificationsdomain "github.com/woragis/backend/server/app/internal/domains/certifications"
	chatsdomain "github.com/woragis/backend/server/app/internal/domains/chats"
	clientsdomain "github.com/woragis/backend/server/app/internal/domains/clients"
	creativeassetsdomain "github.com/woragis/backend/server/app/internal/domains/creativeassets"
	experiencesdomain "github.com/woragis/backend/server/app/internal/domains/experiences"
	financesdomain "github.com/woragis/backend/server/app/internal/domains/finances"
	ideasdomain "github.com/woragis/backend/server/app/internal/domains/ideas"
	impactmetricsdomain "github.com/woragis/backend/server/app/internal/domains/impactmetrics"
	interestsdomain "github.com/woragis/backend/server/app/internal/domains/interests"
	jobapplicationsdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications"
	responsesdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications/responses"
	jobapplicationstagesdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications/interviewstages"
	jobwebsitesdomain "github.com/woragis/backend/server/app/internal/domains/jobwebsites"
	languagesdomain "github.com/woragis/backend/server/app/internal/domains/languages"
	commentsdomain "github.com/woragis/backend/server/app/internal/domains/posts/comments"
	postsdomain "github.com/woragis/backend/server/app/internal/domains/posts"
	problemsolutionsdomain "github.com/woragis/backend/server/app/internal/domains/problemsolutions"
	projectcasestudiesdomain "github.com/woragis/backend/server/app/internal/domains/projects/projectcasestudies"
	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	reportsdomain "github.com/woragis/backend/server/app/internal/domains/reports"
	resumesdomain "github.com/woragis/backend/server/app/internal/domains/resumes"
	schedulerdomain "github.com/woragis/backend/server/app/internal/domains/scheduler"
	skillsdomain "github.com/woragis/backend/server/app/internal/domains/skills"
	socialmediapostsdomain "github.com/woragis/backend/server/app/internal/domains/socialmediaposts"
	systemdesignsdomain "github.com/woragis/backend/server/app/internal/domains/systemdesigns"
	technicalwritingsdomain "github.com/woragis/backend/server/app/internal/domains/technicalwritings"
	testimonialsdomain "github.com/woragis/backend/server/app/internal/domains/testimonials"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	userprofilesdomain "github.com/woragis/backend/server/app/internal/domains/userprofiles"
	userpreferencesdomain "github.com/woragis/backend/server/app/internal/domains/userpreferences"
)

// MigrateTestDB runs database migrations for all domains including socialmediaposts
// This function is excluded when building with performance_test tag to avoid import cycles
// Accepts testing.TB to support both *testing.T and *testing.B
func MigrateTestDB(t testing.TB, db *gorm.DB) error {
	// Migrate all domains
	domains := []interface{}{
		// Auth and user domains
		&authdomain.User{},
		&userprofilesdomain.UserProfile{},
		&userpreferencesdomain.UserPreferences{},
		
		// Content domains
		&postsdomain.Post{},
		&commentsdomain.Comment{},
		&projectsdomain.Project{},
		&projectcasestudiesdomain.ProjectCaseStudy{},
		&skillsdomain.Skill{},
		&interestsdomain.Interest{},
		&experiencesdomain.Experience{},
		&certificationsdomain.Certification{},
		&casestudiesdomain.CaseStudy{},
		&testimonialsdomain.Testimonial{},
		&creativeassetsdomain.CreativeAsset{},
		
		// Job application domains
		&jobapplicationsdomain.JobApplication{},
		&responsesdomain.Response{},
		&jobapplicationstagesdomain.InterviewStage{},
		&jobwebsitesdomain.JobWebsite{},
		
		// Other domains
		&apikeysdomain.APIKey{},
		&translationsdomain.Translation{},
		&languagesdomain.StudySession{},
		&languagesdomain.VocabularyEntry{},
		&resumesdomain.Resume{},
		&clientsdomain.Client{},
		&financesdomain.Transaction{},
		&financesdomain.RecurringTemplate{},
		&ideasdomain.Idea{},
		&problemsolutionsdomain.ProblemSolution{},
		&systemdesignsdomain.SystemDesign{},
		&technicalwritingsdomain.TechnicalWriting{},
		&reportsdomain.ReportDefinition{},
		&reportsdomain.ReportSchedule{},
		&reportsdomain.ReportDelivery{},
		&reportsdomain.ReportRun{},
		&aimlintegrationsdomain.AIMLIntegration{},
		&chatsdomain.Conversation{},
		&chatsdomain.Message{},
		&chatsdomain.ConversationTranscript{},
		&chatsdomain.ConversationAssignment{},
		&schedulerdomain.Schedule{},
		&schedulerdomain.ExecutionRun{},
		&impactmetricsdomain.ImpactMetric{},
		
		// Social media posts (included in normal builds, excluded in performance_test builds)
		&socialmediapostsdomain.SocialMediaPost{},
	}
	
	if err := db.AutoMigrate(domains...); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
		return err
	}
	
	return nil
}
