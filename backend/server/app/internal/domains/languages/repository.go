package languages

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence behavior for the languages domain.
type Repository interface {
	CreateStudySession(ctx context.Context, session *StudySession) error
	ListStudySessions(ctx context.Context, userID uuid.UUID, language string, from, to time.Time) ([]StudySession, error)
	CreateVocabularyEntry(ctx context.Context, entry *VocabularyEntry) error
	ListVocabularyEntries(ctx context.Context, userID uuid.UUID, language string, dueOnly bool) ([]VocabularyEntry, error)
	AggregateSummary(ctx context.Context, userID uuid.UUID) ([]LanguageSummary, error)
}

// LanguageSummary aggregates learning metrics per language.
type LanguageSummary struct {
	LanguageCode    string
	TotalMinutes    int64
	SessionCount    int64
	VocabularyCount int64
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository provides a GORM-backed repository implementation.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateStudySession(ctx context.Context, session *StudySession) error {
	if err := session.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) ListStudySessions(ctx context.Context, userID uuid.UUID, language string, from, to time.Time) ([]StudySession, error) {
	var sessions []StudySession

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if language != "" {
		query = query.Where("language_code = ?", normalizeLang(language))
	}
	if !from.IsZero() {
		query = query.Where("completed_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("completed_at <= ?", to)
	}

	if err := query.Order("completed_at desc").Find(&sessions).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return sessions, nil
}

func (r *gormRepository) CreateVocabularyEntry(ctx context.Context, entry *VocabularyEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) ListVocabularyEntries(ctx context.Context, userID uuid.UUID, language string, dueOnly bool) ([]VocabularyEntry, error) {
	var entries []VocabularyEntry

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if language != "" {
		query = query.Where("language_code = ?", normalizeLang(language))
	}
	if dueOnly {
		query = query.Where("review_at <= ?", time.Now().UTC())
	}

	if err := query.Order("review_at asc").Find(&entries).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return entries, nil
}

func (r *gormRepository) AggregateSummary(ctx context.Context, userID uuid.UUID) ([]LanguageSummary, error) {
	type row struct {
		LanguageCode    string
		TotalMinutes    int64
		SessionCount    int64
		VocabularyCount int64
	}

	var rows []row

	subSessions := r.db.
		Select("language_code, SUM(duration_min) as total_minutes, COUNT(*) as session_count").
		Model(&StudySession{}).
		Where("user_id = ?", userID).
		Group("language_code")

	subVocabulary := r.db.
		Select("language_code, COUNT(*) as vocabulary_count").
		Model(&VocabularyEntry{}).
		Where("user_id = ?", userID).
		Group("language_code")

	err := r.db.WithContext(ctx).
		Table("(?) as sessions", subSessions).
		Select("sessions.language_code, sessions.total_minutes, sessions.session_count, COALESCE(vocab.vocabulary_count, 0) as vocabulary_count").
		Joins("LEFT JOIN (?) as vocab ON vocab.language_code = sessions.language_code", subVocabulary).
		Scan(&rows).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeSummaryFailure, ErrUnableToSummarize)
	}

	summaries := make([]LanguageSummary, 0, len(rows))
	for _, r := range rows {
		summaries = append(summaries, LanguageSummary{
			LanguageCode:    r.LanguageCode,
			TotalMinutes:    r.TotalMinutes,
			SessionCount:    r.SessionCount,
			VocabularyCount: r.VocabularyCount,
		})
	}

	return summaries, nil
}
