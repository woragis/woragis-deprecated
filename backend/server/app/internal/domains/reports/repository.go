package reports

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DefinitionFilters provides list filters.
type DefinitionFilters struct {
	Search          string
	IncludeArchived bool
	OnlyFavorites   bool
	Channel         string
	Limit           int
	Offset          int
}

// RunFilters provides run filtering.
type RunFilters struct {
	Status string
	Limit  int
	Offset int
}

// Repository defines persistence operations for reports.
type Repository interface {
	CreateDefinition(ctx context.Context, def *ReportDefinition) error
	UpdateDefinition(ctx context.Context, def *ReportDefinition) error
	GetDefinition(ctx context.Context, id, userID uuid.UUID) (*ReportDefinition, error)
	ListDefinitions(ctx context.Context, userID uuid.UUID, filters DefinitionFilters) ([]ReportDefinition, error)
	BulkArchiveDefinitions(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	BulkRestoreDefinitions(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	BulkDeleteDefinitions(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	SetFavorite(ctx context.Context, userID uuid.UUID, id uuid.UUID, favorite bool) error

	CreateSchedule(ctx context.Context, schedule *ReportSchedule) error
	UpdateSchedule(ctx context.Context, schedule *ReportSchedule) error
	GetSchedule(ctx context.Context, id uuid.UUID) (*ReportSchedule, error)
	ListSchedules(ctx context.Context, reportID uuid.UUID) ([]ReportSchedule, error)
	DeleteSchedule(ctx context.Context, id uuid.UUID) error

	CreateDelivery(ctx context.Context, delivery *ReportDelivery) error
	UpdateDelivery(ctx context.Context, delivery *ReportDelivery) error
	GetDelivery(ctx context.Context, id uuid.UUID) (*ReportDelivery, error)
	ListDeliveries(ctx context.Context, reportID uuid.UUID) ([]ReportDelivery, error)
	DeleteDelivery(ctx context.Context, id uuid.UUID) error

	CreateRun(ctx context.Context, run *ReportRun) error
	UpdateRun(ctx context.Context, run *ReportRun) error
	ListRuns(ctx context.Context, reportID uuid.UUID, filters RunFilters) ([]ReportRun, error)
}

// gormRepository implements Repository.
type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository constructs a repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateDefinition(ctx context.Context, def *ReportDefinition) error {
	if err := def.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(def).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateDefinition(ctx context.Context, def *ReportDefinition) error {
	if err := def.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Model(def).
		Where("id = ? AND deleted_at IS NULL", def.ID).
		Updates(map[string]any{
			"name":        def.Name,
			"description": def.Description,
			"sections":    def.Sections,
			"filters":     def.Filters,
			"is_favorite": def.IsFavorite,
			"updated_at":  def.UpdatedAt,
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) GetDefinition(ctx context.Context, id, userID uuid.UUID) (*ReportDefinition, error) {
	var def ReportDefinition
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&def).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrReportDefinitionNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &def, nil
}

func (r *gormRepository) ListDefinitions(ctx context.Context, userID uuid.UUID, filters DefinitionFilters) ([]ReportDefinition, error) {
	var defs []ReportDefinition

	db := r.db.WithContext(ctx).
		Where("user_id = ?", userID)

	if !filters.IncludeArchived {
		db = db.Where("archived_at IS NULL")
	}

	if filters.OnlyFavorites {
		db = db.Where("is_favorite = ?", true)
	}

	if search := strings.TrimSpace(filters.Search); search != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		db = db.Where("(LOWER(name) LIKE ? OR LOWER(description) LIKE ?)", pattern, pattern)
	}

	if filters.Channel != "" {
		channel := strings.ToLower(strings.TrimSpace(filters.Channel))
		db = db.Where("id IN (?)",
			r.db.Table("report_deliveries").
				Select("report_id").
				Where("LOWER(channel) = ? AND deleted_at IS NULL", channel))
	}

	if filters.Limit > 0 {
		db = db.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		db = db.Offset(filters.Offset)
	}

	if err := db.Order("updated_at DESC").
		Find(&defs).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return defs, nil
}

func (r *gormRepository) BulkArchiveDefinitions(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	now := time.Now().UTC()
	if err := r.db.WithContext(ctx).
		Model(&ReportDefinition{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Update("archived_at", now).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) BulkRestoreDefinitions(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&ReportDefinition{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Update("archived_at", nil).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) BulkDeleteDefinitions(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND id IN ?", userID, ids).
		Delete(&ReportDefinition{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) SetFavorite(ctx context.Context, userID uuid.UUID, id uuid.UUID, favorite bool) error {
	if err := r.db.WithContext(ctx).
		Model(&ReportDefinition{}).
		Where("user_id = ? AND id = ?", userID, id).
		Update("is_favorite", favorite).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) CreateSchedule(ctx context.Context, schedule *ReportSchedule) error {
	if err := schedule.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(schedule).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateSchedule(ctx context.Context, schedule *ReportSchedule) error {
	if err := schedule.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Model(&ReportSchedule{}).
		Where("id = ?", schedule.ID).
		Updates(map[string]any{
			"cron":       schedule.Cron,
			"frequency":  schedule.Frequency,
			"timezone":   schedule.Timezone,
			"next_run":   schedule.NextRun,
			"enabled":    schedule.Enabled,
			"meta":       schedule.Meta,
			"updated_at": time.Now().UTC(),
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) GetSchedule(ctx context.Context, id uuid.UUID) (*ReportSchedule, error) {
	var schedule ReportSchedule
	if err := r.db.WithContext(ctx).First(&schedule, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrScheduleNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &schedule, nil
}

func (r *gormRepository) ListSchedules(ctx context.Context, reportID uuid.UUID) ([]ReportSchedule, error) {
	var schedules []ReportSchedule
	if err := r.db.WithContext(ctx).
		Where("report_id = ?", reportID).
		Order("created_at ASC").
		Find(&schedules).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return schedules, nil
}

func (r *gormRepository) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&ReportSchedule{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) CreateDelivery(ctx context.Context, delivery *ReportDelivery) error {
	if err := delivery.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(delivery).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateDelivery(ctx context.Context, delivery *ReportDelivery) error {
	if err := delivery.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).
		Model(&ReportDelivery{}).
		Where("id = ?", delivery.ID).
		Updates(map[string]any{
			"channel":    delivery.Channel,
			"target":     delivery.Target,
			"template":   delivery.Template,
			"enabled":    delivery.Enabled,
			"updated_at": time.Now().UTC(),
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) GetDelivery(ctx context.Context, id uuid.UUID) (*ReportDelivery, error) {
	var delivery ReportDelivery
	if err := r.db.WithContext(ctx).First(&delivery, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrDeliveryNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &delivery, nil
}

func (r *gormRepository) ListDeliveries(ctx context.Context, reportID uuid.UUID) ([]ReportDelivery, error) {
	var deliveries []ReportDelivery
	if err := r.db.WithContext(ctx).
		Where("report_id = ?", reportID).
		Order("created_at ASC").
		Find(&deliveries).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return deliveries, nil
}

func (r *gormRepository) DeleteDelivery(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&ReportDelivery{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) CreateRun(ctx context.Context, run *ReportRun) error {
	if err := r.db.WithContext(ctx).Create(run).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateRun(ctx context.Context, run *ReportRun) error {
	if err := r.db.WithContext(ctx).
		Model(&ReportRun{}).
		Where("id = ?", run.ID).
		Updates(map[string]any{
			"status":          run.Status,
			"started_at":      run.StartedAt,
			"completed_at":    run.CompletedAt,
			"output_location": run.OutputLocation,
			"error_message":   run.ErrorMessage,
			"metadata":        run.Metadata,
			"updated_at":      time.Now().UTC(),
		}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) ListRuns(ctx context.Context, reportID uuid.UUID, filters RunFilters) ([]ReportRun, error) {
	var runs []ReportRun

	db := r.db.WithContext(ctx).
		Where("report_id = ?", reportID)

	if status := strings.TrimSpace(filters.Status); status != "" {
		db = db.Where("LOWER(status) = ?", strings.ToLower(status))
	}

	if filters.Limit > 0 {
		db = db.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		db = db.Offset(filters.Offset)
	}

	if err := db.Order("created_at DESC").Find(&runs).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return runs, nil
}
