package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/bio-ht/internal/domain"
)

type Repository interface {
	CreateTemplate(ctx context.Context, t *domain.BioTemplate) error
	GetTemplate(ctx context.Context, templateID uuid.UUID) (*domain.BioTemplate, error)
	GetActiveTemplatesByCitizen(ctx context.Context, citizenID uuid.UUID) ([]domain.BioTemplate, error)
	DeactivateTemplate(ctx context.Context, templateID uuid.UUID) error
	LogVerification(ctx context.Context, log *domain.VerificationLog) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateTemplate(ctx context.Context, t *domain.BioTemplate) error {
	query := `INSERT INTO bio_templates (template_id, citizen_id, modality, milvus_vector_id, quality_score, capture_device, capture_location, captured_by, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query,
		t.TemplateID, t.CitizenID, t.Modality, t.MilvusVectorID, t.QualityScore,
		t.CaptureDevice, t.CaptureLocation, t.CapturedBy, t.IsActive, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert template: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetTemplate(ctx context.Context, templateID uuid.UUID) (*domain.BioTemplate, error) {
	query := `SELECT template_id, citizen_id, modality, milvus_vector_id, quality_score, capture_device, capture_location, captured_by, is_active, superseded_by_template_id, created_at
		FROM bio_templates WHERE template_id = $1`

	t := &domain.BioTemplate{}
	err := r.db.QueryRowContext(ctx, query, templateID).Scan(
		&t.TemplateID, &t.CitizenID, &t.Modality, &t.MilvusVectorID, &t.QualityScore,
		&t.CaptureDevice, &t.CaptureLocation, &t.CapturedBy, &t.IsActive,
		&t.SupersededByTemplateID, &t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("query template: %w", err)
	}
	return t, nil
}

func (r *postgresRepo) GetActiveTemplatesByCitizen(ctx context.Context, citizenID uuid.UUID) ([]domain.BioTemplate, error) {
	query := `SELECT template_id, citizen_id, modality, milvus_vector_id, quality_score, capture_device, capture_location, captured_by, is_active, superseded_by_template_id, created_at
		FROM bio_templates WHERE citizen_id = $1 AND is_active = TRUE`

	rows, err := r.db.QueryContext(ctx, query, citizenID)
	if err != nil {
		return nil, fmt.Errorf("query templates: %w", err)
	}
	defer rows.Close()

	var templates []domain.BioTemplate
	for rows.Next() {
		var t domain.BioTemplate
		if err := rows.Scan(&t.TemplateID, &t.CitizenID, &t.Modality, &t.MilvusVectorID, &t.QualityScore,
			&t.CaptureDevice, &t.CaptureLocation, &t.CapturedBy, &t.IsActive,
			&t.SupersededByTemplateID, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (r *postgresRepo) DeactivateTemplate(ctx context.Context, templateID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE bio_templates SET is_active = FALSE WHERE template_id = $1`, templateID)
	return err
}

func (r *postgresRepo) LogVerification(ctx context.Context, logEntry *domain.VerificationLog) error {
	query := `INSERT INTO bio_verification_log (verification_id, citizen_id, modality, requesting_module, match_score, is_match, verified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		logEntry.VerificationID, logEntry.CitizenID, logEntry.Modality,
		logEntry.RequestingModule, logEntry.MatchScore, logEntry.IsMatch, time.Now().UTC(),
	)
	return err
}
