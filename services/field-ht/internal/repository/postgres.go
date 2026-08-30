package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/snisid/field-ht/internal/domain"
)

type Repository interface {
	CreateMission(ctx context.Context, m *domain.Mission) error
	GetActiveMissions(ctx context.Context) ([]domain.Mission, error)
	CreateMissionLog(ctx context.Context, l *domain.MissionLog) error
	GetCoverageStats(ctx context.Context) (*domain.CoverageStats, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateMission(ctx context.Context, m *domain.Mission) error {
	query := `INSERT INTO field_missions (id, title, description, status, assigned_unit_id, dept_code, started_at, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.Title, m.Description, m.Status, m.AssignedUnitID, m.DeptCode,
		m.StartedAt, m.CompletedAt, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert mission: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetActiveMissions(ctx context.Context) ([]domain.Mission, error) {
	query := `SELECT id, title, description, status, assigned_unit_id, dept_code, started_at, completed_at, created_at
		FROM field_missions WHERE status IN ('PLANNED', 'IN_PROGRESS') ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query active missions: %w", err)
	}
	defer rows.Close()

	var missions []domain.Mission
	for rows.Next() {
		var m domain.Mission
		if err := rows.Scan(
			&m.ID, &m.Title, &m.Description, &m.Status, &m.AssignedUnitID, &m.DeptCode,
			&m.StartedAt, &m.CompletedAt, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		missions = append(missions, m)
	}
	return missions, rows.Err()
}

func (r *postgresRepo) CreateMissionLog(ctx context.Context, l *domain.MissionLog) error {
	query := `INSERT INTO field_mission_logs (id, mission_id, logged_by, action, latitude, longitude, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		l.ID, l.MissionID, l.LoggedBy, l.Action, l.Latitude, l.Longitude, l.Notes, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert mission log: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetCoverageStats(ctx context.Context) (*domain.CoverageStats, error) {
	query := `SELECT
		COALESCE((SELECT COUNT(*) FROM field_missions WHERE status IN ('PLANNED', 'IN_PROGRESS')), 0) AS total_missions,
		COALESCE((SELECT COUNT(*) FROM field_mobile_units WHERE is_active = TRUE), 0) AS active_units,
		0.0 AS coverage_lat,
		0.0 AS coverage_lng,
		0.0 AS coverage_radius_km`

	stats := &domain.CoverageStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalMissions, &stats.ActiveUnits, &stats.CoverageLat, &stats.CoverageLng, &stats.CoverageRadiusKm,
	)
	if err != nil {
		return nil, fmt.Errorf("query coverage stats: %w", err)
	}
	return stats, nil
}
