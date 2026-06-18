package repository

import (
	"context"

	"github.com/snisid/platform/internal/domain/audit/entity"
	"gorm.io/gorm"
)

type AuditRepository interface {
	// Append strictly appends a new event. Updates and Deletes are not supported.
	Append(ctx context.Context, event *entity.AuditEvent) error
	GetLastEvent(ctx context.Context) (*entity.AuditEvent, error)
	GetEventsByCorrelationID(ctx context.Context, correlationID string) ([]entity.AuditEvent, error)
	GetEventsBySequenceRange(ctx context.Context, start, end int64) ([]entity.AuditEvent, error)
}

type postgresAuditRepo struct {
	db *gorm.DB
}

func NewPostgresAuditRepository(db *gorm.DB) AuditRepository {
	return &postgresAuditRepo{db: db}
}

func (r *postgresAuditRepo) Append(ctx context.Context, event *entity.AuditEvent) error {
	// Using a transaction to ensure no dirty reads on last event could be added here
	// But assuming singleton ingester, a simple Create is fine.
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *postgresAuditRepo) GetLastEvent(ctx context.Context) (*entity.AuditEvent, error) {
	var event entity.AuditEvent
	err := r.db.WithContext(ctx).Order("sequence_id desc").First(&event).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // Valid for the first event
	}
	return &event, err
}

func (r *postgresAuditRepo) GetEventsByCorrelationID(ctx context.Context, correlationID string) ([]entity.AuditEvent, error) {
	var events []entity.AuditEvent
	err := r.db.WithContext(ctx).Where("correlation_id = ?", correlationID).Order("sequence_id asc").Find(&events).Error
	return events, err
}

func (r *postgresAuditRepo) GetEventsBySequenceRange(ctx context.Context, start, end int64) ([]entity.AuditEvent, error) {
	var events []entity.AuditEvent
	err := r.db.WithContext(ctx).Where("sequence_id >= ? AND sequence_id <= ?", start, end).Order("sequence_id asc").Find(&events).Error
	return events, err
}
