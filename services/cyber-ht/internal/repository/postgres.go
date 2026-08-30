package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/snisid/cyber-ht/internal/domain"
)

type Repository interface {
	CreateIncident(ctx context.Context, inc *domain.Incident) error
	GetActiveIncidents(ctx context.Context) ([]domain.Incident, error)
	CreatePolicy(ctx context.Context, p *domain.ZeroTrustPolicy) error
	CheckThreatIndicator(ctx context.Context, indicator string) (*domain.ThreatIndicator, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateIncident(ctx context.Context, inc *domain.Incident) error {
	query := `INSERT INTO cyber_incidents (id, title, description, severity, status, source_ip, target_asset, detected_by, assigned_to, created_at, updated_at, closed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		inc.ID, inc.Title, inc.Description, inc.Severity, inc.Status,
		inc.SourceIP, inc.TargetAsset, inc.DetectedBy, inc.AssignedTo,
		inc.CreatedAt, inc.UpdatedAt, inc.ClosedAt,
	)
	return err
}

func (r *postgresRepo) GetActiveIncidents(ctx context.Context) ([]domain.Incident, error) {
	query := `SELECT id, title, description, severity, status, source_ip, target_asset, detected_by, assigned_to, created_at, updated_at, closed_at
		FROM cyber_incidents WHERE status != 'CLOSED' ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []domain.Incident
	for rows.Next() {
		var inc domain.Incident
		if err := rows.Scan(
			&inc.ID, &inc.Title, &inc.Description, &inc.Severity, &inc.Status,
			&inc.SourceIP, &inc.TargetAsset, &inc.DetectedBy, &inc.AssignedTo,
			&inc.CreatedAt, &inc.UpdatedAt, &inc.ClosedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	return incidents, rows.Err()
}

func (r *postgresRepo) CreatePolicy(ctx context.Context, p *domain.ZeroTrustPolicy) error {
	query := `INSERT INTO cyber_zero_trust_policies (id, name, description, policy_type, rules, enabled, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.Name, p.Description, p.PolicyType, pq.StringArray(p.Rules),
		p.Enabled, p.CreatedBy, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) CheckThreatIndicator(ctx context.Context, indicator string) (*domain.ThreatIndicator, error) {
	query := `SELECT id, indicator, type, threat_level, source, description, tags, expires_at, created_at
		FROM cyber_threat_indicators WHERE indicator = $1 AND (expires_at IS NULL OR expires_at > NOW())`
	ti := &domain.ThreatIndicator{}
	var tags pq.StringArray
	err := r.db.QueryRowContext(ctx, query, indicator).Scan(
		&ti.ID, &ti.Indicator, &ti.Type, &ti.ThreatLevel, &ti.Source,
		&ti.Description, &tags, &ti.ExpiresAt, &ti.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("threat indicator not found")
		}
		return nil, err
	}
	ti.Tags = []string(tags)
	return ti, nil
}
