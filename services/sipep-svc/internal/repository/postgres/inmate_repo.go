package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sipep-svc/internal/domain"
)

type InmateRepo struct {
	pool *pgxpool.Pool
}

func NewInmateRepo(pool *pgxpool.Pool) *InmateRepo {
	return &InmateRepo{pool: pool}
}

func (r *InmateRepo) Create(ctx context.Context, inmate *domain.Inmate) error {
	query := `
		INSERT INTO sipep_inmates 
			(inmate_id, national_inmate_id, snisid_person_id, fir_record_id, afis_subject_id,
			 current_facility, current_dept_code, cell_block, is_currently_detained, is_minor,
			 is_female, has_special_needs, special_needs_notes, intake_date, expected_release_date,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.pool.Exec(ctx, query,
		inmate.InmateID, inmate.NationalInmateID, inmate.SNISIDPersonID,
		inmate.FIRRecordID, inmate.AFISSubjectID, inmate.CurrentFacility,
		inmate.CurrentDeptCode, inmate.CellBlock, inmate.IsCurrentlyDetained,
		inmate.IsMinor, inmate.IsFemale, inmate.HasSpecialNeeds,
		inmate.SpecialNeedsNotes, inmate.IntakeDate, inmate.ExpectedReleaseDate,
		inmate.CreatedAt, inmate.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create inmate: %w", err)
	}
	return nil
}

func (r *InmateRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Inmate, error) {
	query := `
		SELECT inmate_id, national_inmate_id, snisid_person_id, fir_record_id, afis_subject_id,
			   current_facility, current_dept_code, cell_block, is_currently_detained, is_minor,
			   is_female, has_special_needs, special_needs_notes, intake_date, expected_release_date,
			   created_at, updated_at
		FROM sipep_inmates
		WHERE inmate_id = $1
	`
	inmate := &domain.Inmate{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&inmate.InmateID, &inmate.NationalInmateID, &inmate.SNISIDPersonID,
		&inmate.FIRRecordID, &inmate.AFISSubjectID, &inmate.CurrentFacility,
		&inmate.CurrentDeptCode, &inmate.CellBlock, &inmate.IsCurrentlyDetained,
		&inmate.IsMinor, &inmate.IsFemale, &inmate.HasSpecialNeeds,
		&inmate.SpecialNeedsNotes, &inmate.IntakeDate, &inmate.ExpectedReleaseDate,
		&inmate.CreatedAt, &inmate.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find inmate: %w", err)
	}
	return inmate, nil
}

func (r *InmateRepo) FindByPersonID(ctx context.Context, personID uuid.UUID) (*domain.Inmate, error) {
	query := `
		SELECT inmate_id, national_inmate_id, snisid_person_id, fir_record_id, afis_subject_id,
			   current_facility, current_dept_code, cell_block, is_currently_detained, is_minor,
			   is_female, has_special_needs, special_needs_notes, intake_date, expected_release_date,
			   created_at, updated_at
		FROM sipep_inmates
		WHERE snisid_person_id = $1 AND is_currently_detained = TRUE
	`
	inmate := &domain.Inmate{}
	err := r.pool.QueryRow(ctx, query, personID).Scan(
		&inmate.InmateID, &inmate.NationalInmateID, &inmate.SNISIDPersonID,
		&inmate.FIRRecordID, &inmate.AFISSubjectID, &inmate.CurrentFacility,
		&inmate.CurrentDeptCode, &inmate.CellBlock, &inmate.IsCurrentlyDetained,
		&inmate.IsMinor, &inmate.IsFemale, &inmate.HasSpecialNeeds,
		&inmate.SpecialNeedsNotes, &inmate.IntakeDate, &inmate.ExpectedReleaseDate,
		&inmate.CreatedAt, &inmate.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find inmate by person: %w", err)
	}
	return inmate, nil
}

func (r *InmateRepo) Update(ctx context.Context, inmate *domain.Inmate) error {
	query := `
		UPDATE sipep_inmates
		SET current_facility = $3, current_dept_code = $4, cell_block = $5,
			is_currently_detained = $6, is_minor = $7, is_female = $8,
			has_special_needs = $9, special_needs_notes = $10, expected_release_date = $11,
			updated_at = $12
		WHERE inmate_id = $1 AND snisid_person_id = $2
	`
	_, err := r.pool.Exec(ctx, query,
		inmate.InmateID, inmate.SNISIDPersonID, inmate.CurrentFacility,
		inmate.CurrentDeptCode, inmate.CellBlock, inmate.IsCurrentlyDetained,
		inmate.IsMinor, inmate.IsFemale, inmate.HasSpecialNeeds,
		inmate.SpecialNeedsNotes, inmate.ExpectedReleaseDate, inmate.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update inmate: %w", err)
	}
	return nil
}

func (r *InmateRepo) Search(ctx context.Context, query string) ([]*domain.Inmate, error) {
	sql := `
		SELECT inmate_id, national_inmate_id, snisid_person_id, fir_record_id, afis_subject_id,
			   current_facility, current_dept_code, cell_block, is_currently_detained, is_minor,
			   is_female, has_special_needs, special_needs_notes, intake_date, expected_release_date,
			   created_at, updated_at
		FROM sipep_inmates
		WHERE national_inmate_id ILIKE $1 OR current_facility ILIKE $1
		LIMIT 50
	`
	rows, err := r.pool.Query(ctx, sql, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search inmates: %w", err)
	}
	defer rows.Close()

	var inmates []*domain.Inmate
	for rows.Next() {
		inmate := &domain.Inmate{}
		err := rows.Scan(
			&inmate.InmateID, &inmate.NationalInmateID, &inmate.SNISIDPersonID,
			&inmate.FIRRecordID, &inmate.AFISSubjectID, &inmate.CurrentFacility,
			&inmate.CurrentDeptCode, &inmate.CellBlock, &inmate.IsCurrentlyDetained,
			&inmate.IsMinor, &inmate.IsFemale, &inmate.HasSpecialNeeds,
			&inmate.SpecialNeedsNotes, &inmate.IntakeDate, &inmate.ExpectedReleaseDate,
			&inmate.CreatedAt, &inmate.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inmate: %w", err)
		}
		inmates = append(inmates, inmate)
	}
	return inmates, nil
}
