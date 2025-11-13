package languages

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// StudySession tracks focused study periods per language.
type StudySession struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;index;not null"`
	LanguageCode string    `gorm:"size:8;not null"`
	SkillFocus   string    `gorm:"size:32"`
	DurationMin  int       `gorm:"not null"`
	Notes        string    `gorm:"size:255"`
	CompletedAt  time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// VocabularyEntry keeps vocabulary and review metadata.
type VocabularyEntry struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;index;not null"`
	LanguageCode string    `gorm:"size:8;index;not null"`
	Term         string    `gorm:"size:120;not null"`
	Translation  string    `gorm:"size:255;not null"`
	Context      string    `gorm:"size:255"`
	AddedAt      time.Time `gorm:"not null"`
	ReviewAt     time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewStudySession builds a new StudySession aggregate.
func NewStudySession(userID uuid.UUID, languageCode, skillFocus string, durationMin int, notes string, completedAt time.Time) (*StudySession, error) {
	session := &StudySession{
		ID:           uuid.New(),
		UserID:       userID,
		LanguageCode: normalizeLang(languageCode),
		SkillFocus:   strings.ToLower(strings.TrimSpace(skillFocus)),
		DurationMin:  durationMin,
		Notes:        strings.TrimSpace(notes),
		CompletedAt:  completedAt.UTC(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	return session, session.Validate()
}

// Validate enforces invariants for StudySession.
func (s *StudySession) Validate() error {
	if s == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilStudySession)
	}

	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySessionID)
	}

	if s.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if s.LanguageCode == "" {
		return NewDomainError(ErrCodeInvalidLanguage, ErrEmptyLanguageCode)
	}

	if len(s.LanguageCode) < 2 || len(s.LanguageCode) > 8 {
		return NewDomainError(ErrCodeInvalidLanguage, ErrInvalidLanguageCode)
	}

	if s.DurationMin <= 0 {
		return NewDomainError(ErrCodeInvalidDuration, ErrDurationMustBePositive)
	}

	if s.CompletedAt.IsZero() {
		return NewDomainError(ErrCodeInvalidCompletedAt, ErrCompletedAtRequired)
	}

	return nil
}

// NewVocabularyEntry creates a vocabulary entry with spaced-review defaults.
func NewVocabularyEntry(userID uuid.UUID, languageCode, term, translation, context string, reviewAt time.Time) (*VocabularyEntry, error) {
	entry := &VocabularyEntry{
		ID:           uuid.New(),
		UserID:       userID,
		LanguageCode: normalizeLang(languageCode),
		Term:         strings.TrimSpace(term),
		Translation:  strings.TrimSpace(translation),
		Context:      strings.TrimSpace(context),
		AddedAt:      time.Now().UTC(),
		ReviewAt:     reviewAt.UTC(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	return entry, entry.Validate()
}

// Validate ensures vocabulary entry fields are consistent.
func (v *VocabularyEntry) Validate() error {
	if v == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilVocabularyEntry)
	}

	if v.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyVocabularyID)
	}

	if v.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if v.LanguageCode == "" {
		return NewDomainError(ErrCodeInvalidLanguage, ErrEmptyLanguageCode)
	}

	if v.Term == "" {
		return NewDomainError(ErrCodeInvalidVocabulary, ErrEmptyTerm)
	}

	if v.Translation == "" {
		return NewDomainError(ErrCodeInvalidVocabulary, ErrEmptyTranslation)
	}

	if v.ReviewAt.IsZero() {
		return NewDomainError(ErrCodeInvalidReviewAt, ErrReviewAtRequired)
	}

	return nil
}

func normalizeLang(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}
