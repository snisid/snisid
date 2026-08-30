package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/data-ht/internal/domain"
)

type Repository interface {
	ListPipelines(ctx context.Context) ([]domain.Pipeline, error)
	CreateModel(ctx context.Context, m *domain.MLModel) error
	GetModel(ctx context.Context, id uuid.UUID) (*domain.MLModel, error)
	GetGovernanceAuditsByModel(ctx context.Context, modelID uuid.UUID) ([]domain.GovernanceAudit, error)
	GetNationalDashboard(ctx context.Context) (*domain.NationalDashboard, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) ListPipelines(ctx context.Context) ([]domain.Pipeline, error) {
	query := `SELECT id, name, source_topics, destination, config, is_active, created_at
		FROM data_pipelines ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query pipelines: %w", err)
	}
	defer rows.Close()

	var pipelines []domain.Pipeline
	for rows.Next() {
		var p domain.Pipeline
		if err := rows.Scan(
			&p.ID, &p.Name, &p.SourceTopics, &p.Destination, &p.Config, &p.IsActive, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		pipelines = append(pipelines, p)
	}
	return pipelines, rows.Err()
}

func (r *postgresRepo) CreateModel(ctx context.Context, m *domain.MLModel) error {
	query := `INSERT INTO data_ml_models (id, name, model_type, version, mlflow_run_id, bias_metric, bias_score, training_date, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.Name, m.ModelType, m.Version, m.MlflowRunID, m.BiasMetric, m.BiasScore, m.TrainingDate, m.IsActive, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert model: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetModel(ctx context.Context, id uuid.UUID) (*domain.MLModel, error) {
	query := `SELECT id, name, model_type, version, mlflow_run_id, bias_metric, bias_score, training_date, is_active, created_at
		FROM data_ml_models WHERE id = $1`

	m := &domain.MLModel{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.ModelType, &m.Version, &m.MlflowRunID,
		&m.BiasMetric, &m.BiasScore, &m.TrainingDate, &m.IsActive, &m.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("model not found")
		}
		return nil, fmt.Errorf("query model: %w", err)
	}
	return m, nil
}

func (r *postgresRepo) GetGovernanceAuditsByModel(ctx context.Context, modelID uuid.UUID) ([]domain.GovernanceAudit, error) {
	query := `SELECT id, model_id, audit_type, findings, conducted_by, conducted_at
		FROM data_governance_audit WHERE model_id = $1 ORDER BY conducted_at DESC`

	rows, err := r.db.QueryContext(ctx, query, modelID)
	if err != nil {
		return nil, fmt.Errorf("query audits: %w", err)
	}
	defer rows.Close()

	var audits []domain.GovernanceAudit
	for rows.Next() {
		var a domain.GovernanceAudit
		if err := rows.Scan(&a.ID, &a.ModelID, &a.AuditType, &a.Findings, &a.ConductedBy, &a.ConductedAt); err != nil {
			return nil, err
		}
		audits = append(audits, a)
	}
	return audits, rows.Err()
}

func (r *postgresRepo) GetNationalDashboard(ctx context.Context) (*domain.NationalDashboard, error) {
	dash := &domain.NationalDashboard{
		ModelTypeBreakdown: make(map[string]int),
	}

	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(COUNT(*), 0) FROM data_pipelines WHERE is_active = TRUE`,
	).Scan(&dash.TotalPipelines)
	if err != nil {
		return nil, fmt.Errorf("query pipeline count: %w", err)
	}

	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(COUNT(*), 0) FROM data_ml_models WHERE is_active = TRUE`,
	).Scan(&dash.ActiveModels)
	if err != nil {
		return nil, fmt.Errorf("query model count: %w", err)
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT model_type, COUNT(*) as cnt FROM data_ml_models WHERE is_active = TRUE GROUP BY model_type`,
	)
	if err != nil {
		return nil, fmt.Errorf("query model breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var modelType string
		var count int
		if err := rows.Scan(&modelType, &count); err != nil {
			return nil, err
		}
		dash.ModelTypeBreakdown[modelType] = count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(COUNT(*), 0) FROM data_governance_audit WHERE conducted_at > NOW() - INTERVAL '30 days'`,
	).Scan(&dash.RecentAudits)
	if err != nil {
		return nil, fmt.Errorf("query recent audits: %w", err)
	}

	return dash, nil
}
