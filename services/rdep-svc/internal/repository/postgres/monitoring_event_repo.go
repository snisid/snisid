package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/rdep-svc/internal/domain"
)

type MonitoringEventRepo struct {
	pool *pgxpool.Pool
}

func NewMonitoringEventRepo(pool *pgxpool.Pool) *MonitoringEventRepo {
	return &MonitoringEventRepo{pool: pool}
}

func (r *MonitoringEventRepo) Create(ctx context.Context, event *domain.MonitoringEvent) error {
	query := `
		INSERT INTO rdep_monitoring_events 
			(event_id, deportee_id, event_type, event_date, location_lat,
			 location_lng, notes, reported_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		event.EventID, event.DeporteeID, event.EventType, event.EventDate,
		event.LocationLat, event.LocationLng, event.Notes, event.ReportedBy, event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create monitoring event: %w", err)
	}
	return nil
}

func (r *MonitoringEventRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.MonitoringEvent, error) {
	query := `
		SELECT event_id, deportee_id, event_type, event_date, location_lat,
			   location_lng, notes, reported_by, created_at
		FROM rdep_monitoring_events
		WHERE event_id = $1
	`
	event := &domain.MonitoringEvent{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&event.EventID, &event.DeporteeID, &event.EventType, &event.EventDate,
		&event.LocationLat, &event.LocationLng, &event.Notes, &event.ReportedBy, &event.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find monitoring event: %w", err)
	}
	return event, nil
}

func (r *MonitoringEventRepo) FindByDeporteeID(ctx context.Context, deporteeID uuid.UUID) ([]*domain.MonitoringEvent, error) {
	query := `
		SELECT event_id, deportee_id, event_type, event_date, location_lat,
			   location_lng, notes, reported_by, created_at
		FROM rdep_monitoring_events
		WHERE deportee_id = $1
		ORDER BY event_date DESC
	`
	rows, err := r.pool.Query(ctx, query, deporteeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query monitoring events: %w", err)
	}
	defer rows.Close()

	var events []*domain.MonitoringEvent
	for rows.Next() {
		event := &domain.MonitoringEvent{}
		err := rows.Scan(
			&event.EventID, &event.DeporteeID, &event.EventType, &event.EventDate,
			&event.LocationLat, &event.LocationLng, &event.Notes, &event.ReportedBy, &event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monitoring event: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}
