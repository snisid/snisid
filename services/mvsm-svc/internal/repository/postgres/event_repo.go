package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/mvsm-svc/internal/domain"
)

type eventRepo struct {
	pool *pgxpool.Pool
}

func NewEventRepo(pool *pgxpool.Pool) *eventRepo {
	return &eventRepo{pool: pool}
}

func (r *eventRepo) Create(event *domain.Event) (*domain.Event, error) {
	ctx := context.Background()
	event.ID = uuid.New()
	event.NationalMvsmID = "MVSM-HT-" + time.Now().Format("2006") + "-" + event.ID.String()[:6]
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	status := "PLANNED"
	event.Status = &status

	_, err := r.pool.Exec(ctx,
		`INSERT INTO mvsm_events
		 (event_id, national_mvsm_id, event_type, event_name, risk_level, status,
		  organizer_name, gang_id, scheduled_date, location_desc, dept_code, commune,
		  lat, lng, estimated_crowd, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
		event.ID, event.NationalMvsmID, event.EventType, event.EventName, event.RiskLevel,
		event.Status, event.OrganizerName, event.GangID, event.ScheduledDate,
		event.LocationDesc, event.DeptCode, event.Commune, event.Lat, event.Lng,
		event.EstimatedCrowd, event.CreatedBy, event.CreatedAt, event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *eventRepo) FindUpcoming() ([]domain.Event, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT event_id, national_mvsm_id, event_type, event_name, risk_level, status,
		        scheduled_date, location_desc, dept_code, commune, estimated_crowd, created_at
		 FROM mvsm_events
		 WHERE scheduled_date >= NOW() AND status IN ('PLANNED','ACTIVE')
		 ORDER BY scheduled_date ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.ID, &e.NationalMvsmID, &e.EventType, &e.EventName,
			&e.RiskLevel, &e.Status, &e.ScheduledDate, &e.LocationDesc,
			&e.DeptCode, &e.Commune, &e.EstimatedCrowd, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *eventRepo) FindActive() ([]domain.Event, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT event_id, national_mvsm_id, event_type, event_name, risk_level, status,
		        scheduled_date, location_desc, dept_code, commune, estimated_crowd, created_at
		 FROM mvsm_events WHERE status = 'ACTIVE' ORDER BY scheduled_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.ID, &e.NationalMvsmID, &e.EventType, &e.EventName,
			&e.RiskLevel, &e.Status, &e.ScheduledDate, &e.LocationDesc,
			&e.DeptCode, &e.Commune, &e.EstimatedCrowd, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *eventRepo) AddUpdate(update *domain.RealTimeUpdate) error {
	ctx := context.Background()
	update.UpdateID = uuid.New()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO mvsm_real_time_updates
		 (update_id, event_id, update_time, current_crowd_est, situation, risk_change,
		  action_taken, reported_by, lat, lng)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		update.UpdateID, update.EventID, update.UpdateTime, update.CurrentCrowdEst,
		update.Situation, update.RiskChange, update.ActionTaken, update.ReportedBy,
		update.Lat, update.Lng,
	)
	return err
}

func (r *eventRepo) UpdateRiskLevel(id uuid.UUID, level domain.RiskLevel) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE mvsm_events SET risk_level = $1, updated_at = NOW() WHERE event_id = $2`,
		level, id)
	return err
}
