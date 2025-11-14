package languages

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// StudySession tracks focused study periods per language.
type StudySession struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	LanguageCode string    `gorm:"column:language_code;size:8;not null" json:"languageCode"`
	SkillFocus   string    `gorm:"column:skill_focus;size:32" json:"skillFocus"`
	DurationMin  int       `gorm:"column:duration_min;not null" json:"durationMin"`
	Notes        string    `gorm:"column:notes;size:255" json:"notes"`
	CompletedAt  time.Time `gorm:"column:completed_at;not null" json:"completedAt"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// VocabularyEntry keeps vocabulary and review metadata.
type VocabularyEntry struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	LanguageCode string    `gorm:"column:language_code;size:8;index;not null" json:"languageCode"`
	Term         string    `gorm:"column:term;size:120;not null" json:"term"`
	Translation  string    `gorm:"column:translation;size:255;not null" json:"translation"`
	Context      string    `gorm:"column:context;size:255" json:"context"`
	AddedAt      time.Time `gorm:"column:added_at;not null" json:"addedAt"`
	ReviewAt     time.Time `gorm:"column:review_at;not null" json:"reviewAt"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
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
