package impactmetrics

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// MetricType represents the type of impact metric.
type MetricType string

const (
	MetricTypeProjectsDelivered    MetricType = "projects_delivered"
	MetricTypeUsersImpacted       MetricType = "users_impacted"
	MetricTypePerformanceImprovement MetricType = "performance_improvement"
	MetricTypeCostSavings         MetricType = "cost_savings"
	MetricTypeTimeSaved           MetricType = "time_saved"
)

// MetricUnit represents the unit of measurement for a metric.
type MetricUnit string

const (
	MetricUnitCount      MetricUnit = "count"
	MetricUnitPercentage MetricUnit = "percentage"
	MetricUnitCurrency   MetricUnit = "currency"
	MetricUnitHours      MetricUnit = "hours"
	MetricUnitDays       MetricUnit = "days"
	MetricUnitMonths     MetricUnit = "months"
	MetricUnitYears      MetricUnit = "years"
	MetricUnitMilliseconds MetricUnit = "milliseconds"
	MetricUnitSeconds    MetricUnit = "seconds"
	MetricUnitMinutes    MetricUnit = "minutes"
)

// EntityType represents the type of entity being linked to a metric.
type EntityType string

const (
	EntityTypeProject         EntityType = "project"
	EntityTypeProblemSolution EntityType = "problem_solution"
	EntityTypeCaseStudy       EntityType = "case_study"
	EntityTypeSystemDesign    EntityType = "system_design"
)

// ImpactMetric represents a single impact metric entry.
type ImpactMetric struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Type        MetricType `gorm:"column:type;type:varchar(50);not null;index" json:"type"`
	Value       float64    `gorm:"column:value;not null" json:"value"`
	Unit        MetricUnit `gorm:"column:unit;type:varchar(30);not null" json:"unit"`
	Description string     `gorm:"column:description;type:text" json:"description,omitempty"`
	// Optional entity linking
	EntityType *EntityType `gorm:"column:entity_type;type:varchar(50);index" json:"entityType,omitempty"`
	EntityID   *uuid.UUID  `gorm:"column:entity_id;type:uuid;index" json:"entityId,omitempty"`
	// Time period for the metric
	PeriodStart *time.Time `gorm:"column:period_start;type:date" json:"periodStart,omitempty"`
	PeriodEnd   *time.Time `gorm:"column:period_end;type:date" json:"periodEnd,omitempty"`
	// Metadata
	Featured    bool      `gorm:"column:featured;not null;default:false;index" json:"featured"`
	DisplayOrder int      `gorm:"column:display_order;not null;default:0;index" json:"displayOrder"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for ImpactMetric.
func (ImpactMetric) TableName() string {
	return "impact_metrics"
}

// NewImpactMetric creates a new impact metric entity.
func NewImpactMetric(userID uuid.UUID, metricType MetricType, value float64, unit MetricUnit) (*ImpactMetric, error) {
	metric := &ImpactMetric{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      metricType,
		Value:     value,
		Unit:      unit,
		Featured:  false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return metric, metric.Validate()
}

// Validate ensures impact metric invariants hold.
func (m *ImpactMetric) Validate() error {
	if m == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilMetric)
	}
	if m.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyMetricID)
	}
	if m.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if !isValidMetricType(m.Type) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedMetricType)
	}
	if !isValidMetricUnit(m.Unit) {
		return NewDomainError(ErrCodeInvalidUnit, ErrUnsupportedMetricUnit)
	}
	if m.Value < 0 {
		return NewDomainError(ErrCodeInvalidValue, ErrNegativeValue)
	}
	if m.EntityType != nil && !isValidEntityType(*m.EntityType) {
		return NewDomainError(ErrCodeInvalidEntityType, ErrUnsupportedEntityType)
	}
	if m.PeriodStart != nil && m.PeriodEnd != nil && m.PeriodEnd.Before(*m.PeriodStart) {
		return NewDomainError(ErrCodeInvalidDate, ErrPeriodEndBeforeStart)
	}
	return nil
}

// UpdateDetails updates metric details.
func (m *ImpactMetric) UpdateDetails(description string, periodStart, periodEnd *time.Time) {
	if description != "" {
		m.Description = strings.TrimSpace(description)
	}
	if periodStart != nil {
		m.PeriodStart = periodStart
	}
	if periodEnd != nil {
		m.PeriodEnd = periodEnd
	}
	m.UpdatedAt = time.Now().UTC()
}

// SetValue updates the metric value.
func (m *ImpactMetric) SetValue(value float64) error {
	if value < 0 {
		return NewDomainError(ErrCodeInvalidValue, ErrNegativeValue)
	}
	m.Value = value
	m.UpdatedAt = time.Now().UTC()
	return nil
}

// SetEntityLink updates the entity link.
func (m *ImpactMetric) SetEntityLink(entityType EntityType, entityID uuid.UUID) error {
	if !isValidEntityType(entityType) {
		return NewDomainError(ErrCodeInvalidEntityType, ErrUnsupportedEntityType)
	}
	m.EntityType = &entityType
	m.EntityID = &entityID
	m.UpdatedAt = time.Now().UTC()
	return nil
}

// ClearEntityLink removes the entity link.
func (m *ImpactMetric) ClearEntityLink() {
	m.EntityType = nil
	m.EntityID = nil
	m.UpdatedAt = time.Now().UTC()
}

// SetFeatured updates the featured flag.
func (m *ImpactMetric) SetFeatured(featured bool) {
	m.Featured = featured
	m.UpdatedAt = time.Now().UTC()
}

// SetDisplayOrder updates the display order.
func (m *ImpactMetric) SetDisplayOrder(order int) {
	m.DisplayOrder = order
	m.UpdatedAt = time.Now().UTC()
}

// Validation helpers

func isValidMetricType(mt MetricType) bool {
	switch mt {
	case MetricTypeProjectsDelivered, MetricTypeUsersImpacted, MetricTypePerformanceImprovement,
		MetricTypeCostSavings, MetricTypeTimeSaved:
		return true
	}
	return false
}

func isValidMetricUnit(mu MetricUnit) bool {
	switch mu {
	case MetricUnitCount, MetricUnitPercentage, MetricUnitCurrency,
		MetricUnitHours, MetricUnitDays, MetricUnitMonths, MetricUnitYears,
		MetricUnitMilliseconds, MetricUnitSeconds, MetricUnitMinutes:
		return true
	}
	return false
}

func isValidEntityType(et EntityType) bool {
	switch et {
	case EntityTypeProject, EntityTypeProblemSolution, EntityTypeCaseStudy, EntityTypeSystemDesign:
		return true
	}
	return false
}

