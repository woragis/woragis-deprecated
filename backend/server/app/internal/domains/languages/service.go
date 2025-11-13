package languages

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates commands for language learning workflows.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// StudySessionRequest carries the payload for logging sessions.
type StudySessionRequest struct {
	UserID       uuid.UUID
	LanguageCode string
	SkillFocus   string
	DurationMin  int
	Notes        string
	CompletedAt  time.Time
}

// VocabularyRequest carries the payload for vocabulary registration.
type VocabularyRequest struct {
	UserID       uuid.UUID
	LanguageCode string
	Term         string
	Translation  string
	Context      string
	ReviewAt     time.Time
}

// LogStudySession persists a new study session entry.
func (s *Service) LogStudySession(ctx context.Context, req StudySessionRequest) (*StudySession, error) {
	if req.CompletedAt.IsZero() {
		req.CompletedAt = time.Now()
	}

	session, err := NewStudySession(
		req.UserID,
		req.LanguageCode,
		req.SkillFocus,
		req.DurationMin,
		req.Notes,
		req.CompletedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateStudySession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// AddVocabularyEntry stores a vocabulary entry with review metadata.
func (s *Service) AddVocabularyEntry(ctx context.Context, req VocabularyRequest) (*VocabularyEntry, error) {
	if req.ReviewAt.IsZero() {
		req.ReviewAt = time.Now().Add(24 * time.Hour)
	}

	entry, err := NewVocabularyEntry(
		req.UserID,
		req.LanguageCode,
		req.Term,
		req.Translation,
		req.Context,
		req.ReviewAt,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateVocabularyEntry(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// GetStudySessions returns user sessions with optional filters.
func (s *Service) GetStudySessions(ctx context.Context, userID uuid.UUID, language string, from, to time.Time) ([]StudySession, error) {
	return s.repo.ListStudySessions(ctx, userID, language, from, to)
}

// GetVocabulary returns vocabulary entries optionally filtered by language or due status.
func (s *Service) GetVocabulary(ctx context.Context, userID uuid.UUID, language string, dueOnly bool) ([]VocabularyEntry, error) {
	return s.repo.ListVocabularyEntries(ctx, userID, language, dueOnly)
}

// GetSummary aggregates progress metrics per language.
func (s *Service) GetSummary(ctx context.Context, userID uuid.UUID) ([]LanguageSummary, error) {
	summaries, err := s.repo.AggregateSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	if s.logger != nil {
		for _, summary := range summaries {
			s.logger.Debug("languages summary",
				slog.String("user_id", userID.String()),
				slog.String("language_code", summary.LanguageCode),
				slog.Int64("total_minutes", summary.TotalMinutes),
				slog.Int64("session_count", summary.SessionCount),
				slog.Int64("vocabulary_count", summary.VocabularyCount),
			)
		}
	}

	return summaries, nil
}
