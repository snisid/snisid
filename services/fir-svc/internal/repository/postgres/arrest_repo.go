package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/fir-svc/internal/domain"
)

type ArrestRepo struct {
	pool *pgxpool.Pool
}

func NewArrestRepo(pool *pgxpool.Pool) *ArrestRepo {
	return &ArrestRepo{pool: pool}
}

func (r *ArrestRepo) Create(ctx context.Context, arrest *domain.Arrest) error {
	query := `
		INSERT INTO fir_arrests 
			(arrest_id, record_id, arrest_date, arresting_unit, arresting_officer,
			 arrest_location, dept_code, charges_text, offense_class, case_reference,
			 release_date, release_reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.pool.Exec(ctx, query,
		arrest.ArrestID, arrest.RecordID, arrest.ArrestDate, arrest.ArrestingUnit,
		arrest.ArrestingOfficer, arrest.ArrestLocation, arrest.DeptCode,
		arrest.ChargesText, arrest.OffenseClass, arrest.CaseReference,
		arrest.ReleaseDate, arrest.ReleaseReason, arrest.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create arrest: %w", err)
	}
	return nil
}

func (r *ArrestRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Arrest, error) {
	query := `
		SELECT arrest_id, record_id, arrest_date, arresting_unit, arresting_officer,
			   arrest_location, dept_code, charges_text, offense_class, case_reference,
			   release_date, release_reason, created_at
		FROM fir_arrests
		WHERE arrest_id = $1
	`
	arrest := &domain.Arrest{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&arrest.ArrestID, &arrest.RecordID, &arrest.ArrestDate, &arrest.ArrestingUnit,
		&arrest.ArrestingOfficer, &arrest.ArrestLocation, &arrest.DeptCode,
		&arrest.ChargesText, &arrest.OffenseClass, &arrest.CaseReference,
		&arrest.ReleaseDate, &arrest.ReleaseReason, &arrest.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find arrest: %w", err)
	}
	return arrest, nil
}

func (r *ArrestRepo) FindByRecordID(ctx context.Context, recordID uuid.UUID) ([]*domain.Arrest, error) {
	query := `
		SELECT arrest_id, record_id, arrest_date, arresting_unit, arresting_officer,
			   arrest_location, dept_code, charges_text, offense_class, case_reference,
			   release_date, release_reason, created_at
		FROM fir_arrests
		WHERE record_id = $1
		ORDER BY arrest_date DESC
	`
	rows, err := r.pool.Query(ctx, query, recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to query arrests: %w", err)
	}
	defer rows.Close()

	var arrests []*domain.Arrest
	for rows.Next() {
		arrest := &domain.Arrest{}
		err := rows.Scan(
			&arrest.ArrestID, &arrest.RecordID, &arrest.ArrestDate, &arrest.ArrestingUnit,
			&arrest.ArrestingOfficer, &arrest.ArrestLocation, &arrest.DeptCode,
			&arrest.ChargesText, &arrest.OffenseClass, &arrest.CaseReference,
			&arrest.ReleaseDate, &arrest.ReleaseReason, &arrest.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan arrest: %w", err)
		}
		arrests = append(arrests, arrest)
	}
	return arrests, nil
}

func (r *ArrestRepo) Update(ctx context.Context, arrest *domain.Arrest) error {
	query := `
		UPDATE fir_arrests
		SET arrest_date = $3, arresting_unit = $4, arresting_officer = $5,
			arrest_location = $6, dept_code = $7, charges_text = $8,
			offense_class = $9, case_reference = $10, release_date = $11,
			release_reason = $12
		WHERE arrest_id = $1 AND record_id = $2
	`
	_, err := r.pool.Exec(ctx, query,
		arrest.ArrestID, arrest.RecordID, arrest.ArrestDate, arrest.ArrestingUnit,
		arrest.ArrestingOfficer, arrest.ArrestLocation, arrest.DeptCode,
		arrest.ChargesText, arrest.OffenseClass, arrest.CaseReference,
		arrest.ReleaseDate, arrest.ReleaseReason,
	)
	if err != nil {
		return fmt.Errorf("failed to update arrest: %w", err)
	}
	return nil
}
