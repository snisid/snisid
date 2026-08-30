package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sigeo-svc/internal/domain"
)

type geoRepo struct {
	pool *pgxpool.Pool
}

func NewGeoRepo(pool *pgxpool.Pool) *geoRepo {
	return &geoRepo{pool: pool}
}

func (r *geoRepo) CreateIncident(incident *domain.Incident) (*domain.Incident, error) {
	ctx := context.Background()
	incident.ID = uuid.New()
	incident.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO sigeo_incidents_unified
		 (event_id, source_module, source_record_id, event_type, event_date, lat, lng,
		  dept_code, commune, severity, gang_id, description, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		incident.ID, incident.SourceModule, incident.SourceRecordID, incident.EventType,
		incident.EventDate, incident.Lat, incident.Lng, incident.DeptCode, incident.Commune,
		incident.Severity, incident.GangID, incident.Description, incident.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return incident, nil
}

func (r *geoRepo) FindIncidents(deptCode string, since time.Time) ([]domain.Incident, error) {
	ctx := context.Background()
	query := `SELECT event_id, source_module, source_record_id, event_type, event_date, lat, lng,
	                 dept_code, commune, severity, gang_id, description, created_at
	          FROM sigeo_incidents_unified WHERE event_date >= $1`
	args := []interface{}{since}

	if deptCode != "" {
		query += " AND dept_code = $2"
		args = append(args, deptCode)
	}
	query += " ORDER BY event_date DESC LIMIT 500"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []domain.Incident
	for rows.Next() {
		var inc domain.Incident
		if err := rows.Scan(&inc.ID, &inc.SourceModule, &inc.SourceRecordID, &inc.EventType,
			&inc.EventDate, &inc.Lat, &inc.Lng, &inc.DeptCode, &inc.Commune,
			&inc.Severity, &inc.GangID, &inc.Description, &inc.CreatedAt); err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	return incidents, nil
}

func (r *geoRepo) FindCheckpoints() ([]domain.Checkpoint, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT cp_id, cp_type, dept_code, road_number, description, controlling_gang_id, is_active, created_at
		 FROM sigeo_checkpoints WHERE is_active = TRUE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkpoints []domain.Checkpoint
	for rows.Next() {
		var cp domain.Checkpoint
		if err := rows.Scan(&cp.ID, &cp.CPType, &cp.DeptCode, &cp.RoadNumber,
			&cp.Description, &cp.ControllingGangID, &cp.IsActive, &cp.CreatedAt); err != nil {
			return nil, err
		}
		checkpoints = append(checkpoints, cp)
	}
	return checkpoints, nil
}

func (r *geoRepo) CountIncidentsByZone(deptCode string, since time.Time) (int, error) {
	ctx := context.Background()
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM sigeo_incidents_unified WHERE dept_code = $1 AND event_date >= $2`,
		deptCode, since).Scan(&count)
	return count, err
}
