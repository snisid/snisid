package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/gang-svc/internal/domain"
)

type IncidentRepo struct {
	pool *pgxpool.Pool
}

func NewIncidentRepo(pool *pgxpool.Pool) *IncidentRepo {
	return &IncidentRepo{pool: pool}
}

func (r *IncidentRepo) Create(ctx context.Context, incident *domain.Incident) error {
	query := `
		INSERT INTO gang_incidents 
			(incident_id, gang_id, incident_type, incident_date, location_desc,
			 dept_code, commune, lat, lng, casualties, victim_ids, sivc_alert_id,
			 description, intelligence_source, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.pool.Exec(ctx, query,
		incident.IncidentID, incident.GangID, incident.IncidentType,
		incident.IncidentDate, incident.LocationDesc, incident.DeptCode,
		incident.Commune, incident.Lat, incident.Lng, incident.Casualties,
		incident.VictimIDs, incident.SIVCAlertID, incident.Description,
		incident.IntelSource, incident.CreatedBy, incident.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create incident: %w", err)
	}
	return nil
}

func (r *IncidentRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	query := `
		SELECT incident_id, gang_id, incident_type, incident_date, location_desc,
			   dept_code, commune, lat, lng, casualties, victim_ids, sivc_alert_id,
			   description, intelligence_source, created_by, created_at
		FROM gang_incidents
		WHERE incident_id = $1
	`
	incident := &domain.Incident{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&incident.IncidentID, &incident.GangID, &incident.IncidentType,
		&incident.IncidentDate, &incident.LocationDesc, &incident.DeptCode,
		&incident.Commune, &incident.Lat, &incident.Lng, &incident.Casualties,
		&incident.VictimIDs, &incident.SIVCAlertID, &incident.Description,
		&incident.IntelSource, &incident.CreatedBy, &incident.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find incident: %w", err)
	}
	return incident, nil
}

func (r *IncidentRepo) FindByGangID(ctx context.Context, gangID uuid.UUID) ([]*domain.Incident, error) {
	query := `
		SELECT incident_id, gang_id, incident_type, incident_date, location_desc,
			   dept_code, commune, lat, lng, casualties, victim_ids, sivc_alert_id,
			   description, intelligence_source, created_by, created_at
		FROM gang_incidents
		WHERE gang_id = $1
		ORDER BY incident_date DESC
	`
	return r.queryIncidents(ctx, query, gangID)
}

func (r *IncidentRepo) FindByDeptCode(ctx context.Context, deptCode string) ([]*domain.Incident, error) {
	query := `
		SELECT incident_id, gang_id, incident_type, incident_date, location_desc,
			   dept_code, commune, lat, lng, casualties, victim_ids, sivc_alert_id,
			   description, intelligence_source, created_by, created_at
		FROM gang_incidents
		WHERE dept_code = $1
		ORDER BY incident_date DESC
		LIMIT 100
	`
	return r.queryIncidents(ctx, query, deptCode)
}

func (r *IncidentRepo) queryIncidents(ctx context.Context, query string, args ...interface{}) ([]*domain.Incident, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query incidents: %w", err)
	}
	defer rows.Close()

	var incidents []*domain.Incident
	for rows.Next() {
		incident := &domain.Incident{}
		err := rows.Scan(
			&incident.IncidentID, &incident.GangID, &incident.IncidentType,
			&incident.IncidentDate, &incident.LocationDesc, &incident.DeptCode,
			&incident.Commune, &incident.Lat, &incident.Lng, &incident.Casualties,
			&incident.VictimIDs, &incident.SIVCAlertID, &incident.Description,
			&incident.IntelSource, &incident.CreatedBy, &incident.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}
		incidents = append(incidents, incident)
	}
	return incidents, nil
}
