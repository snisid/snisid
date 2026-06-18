package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/sipep-svc/internal/domain"
)

type DetentionRepo struct {
	pool *pgxpool.Pool
}

func NewDetentionRepo(pool *pgxpool.Pool) *DetentionRepo {
	return &DetentionRepo{pool: pool}
}

func (r *DetentionRepo) Create(ctx context.Context, detention *domain.Detention) error {
	query := `
		INSERT INTO sipep_detentions 
			(detention_id, inmate_id, facility, detention_basis, legal_status,
			 case_reference, court_name, arresting_authority, warrant_number,
			 intake_date, intake_officer, sentence_duration_days, release_date,
			 release_type, releasing_authority, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.pool.Exec(ctx, query,
		detention.DetentionID, detention.InmateID, detention.Facility,
		detention.DetentionBasis, detention.LegalStatus, detention.CaseReference,
		detention.CourtName, detention.ArrestingAuthority, detention.WarrantNumber,
		detention.IntakeDate, detention.IntakeOfficer, detention.SentenceDurationDays,
		detention.ReleaseDate, detention.ReleaseType, detention.ReleasingAuthority,
		detention.Notes, detention.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create detention: %w", err)
	}
	return nil
}

func (r *DetentionRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Detention, error) {
	query := `
		SELECT detention_id, inmate_id, facility, detention_basis, legal_status,
			   case_reference, court_name, arresting_authority, warrant_number,
			   intake_date, intake_officer, sentence_duration_days, release_date,
			   release_type, releasing_authority, notes, created_at
		FROM sipep_detentions
		WHERE detention_id = $1
	`
	detention := &domain.Detention{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&detention.DetentionID, &detention.InmateID, &detention.Facility,
		&detention.DetentionBasis, &detention.LegalStatus, &detention.CaseReference,
		&detention.CourtName, &detention.ArrestingAuthority, &detention.WarrantNumber,
		&detention.IntakeDate, &detention.IntakeOfficer, &detention.SentenceDurationDays,
		&detention.ReleaseDate, &detention.ReleaseType, &detention.ReleasingAuthority,
		&detention.Notes, &detention.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find detention: %w", err)
	}
	return detention, nil
}

func (r *DetentionRepo) FindByInmateID(ctx context.Context, inmateID uuid.UUID) ([]*domain.Detention, error) {
	query := `
		SELECT detention_id, inmate_id, facility, detention_basis, legal_status,
			   case_reference, court_name, arresting_authority, warrant_number,
			   intake_date, intake_officer, sentence_duration_days, release_date,
			   release_type, releasing_authority, notes, created_at
		FROM sipep_detentions
		WHERE inmate_id = $1
		ORDER BY intake_date DESC
	`
	rows, err := r.pool.Query(ctx, query, inmateID)
	if err != nil {
		return nil, fmt.Errorf("failed to query detentions: %w", err)
	}
	defer rows.Close()

	var detentions []*domain.Detention
	for rows.Next() {
		detention := &domain.Detention{}
		err := rows.Scan(
			&detention.DetentionID, &detention.InmateID, &detention.Facility,
			&detention.DetentionBasis, &detention.LegalStatus, &detention.CaseReference,
			&detention.CourtName, &detention.ArrestingAuthority, &detention.WarrantNumber,
			&detention.IntakeDate, &detention.IntakeOfficer, &detention.SentenceDurationDays,
			&detention.ReleaseDate, &detention.ReleaseType, &detention.ReleasingAuthority,
			&detention.Notes, &detention.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan detention: %w", err)
		}
		detentions = append(detentions, detention)
	}
	return detentions, nil
}

func (r *DetentionRepo) GetActiveDetention(ctx context.Context, inmateID uuid.UUID) (*domain.Detention, error) {
	query := `
		SELECT detention_id, inmate_id, facility, detention_basis, legal_status,
			   case_reference, court_name, arresting_authority, warrant_number,
			   intake_date, intake_officer, sentence_duration_days, release_date,
			   release_type, releasing_authority, notes, created_at
		FROM sipep_detentions
		WHERE inmate_id = $1 AND release_date IS NULL
		ORDER BY intake_date DESC
		LIMIT 1
	`
	detention := &domain.Detention{}
	err := r.pool.QueryRow(ctx, query, inmateID).Scan(
		&detention.DetentionID, &detention.InmateID, &detention.Facility,
		&detention.DetentionBasis, &detention.LegalStatus, &detention.CaseReference,
		&detention.CourtName, &detention.ArrestingAuthority, &detention.WarrantNumber,
		&detention.IntakeDate, &detention.IntakeOfficer, &detention.SentenceDurationDays,
		&detention.ReleaseDate, &detention.ReleaseType, &detention.ReleasingAuthority,
		&detention.Notes, &detention.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find active detention: %w", err)
	}
	return detention, nil
}

func (r *DetentionRepo) Update(ctx context.Context, detention *domain.Detention) error {
	query := `
		UPDATE sipep_detentions
		SET facility = $3, detention_basis = $4, legal_status = $5,
			case_reference = $6, court_name = $7, arresting_authority = $8,
			warrant_number = $9, sentence_duration_days = $10, release_date = $11,
			release_type = $12, releasing_authority = $13, notes = $14
		WHERE detention_id = $1 AND inmate_id = $2
	`
	_, err := r.pool.Exec(ctx, query,
		detention.DetentionID, detention.InmateID, detention.Facility,
		detention.DetentionBasis, detention.LegalStatus, detention.CaseReference,
		detention.CourtName, detention.ArrestingAuthority, detention.WarrantNumber,
		detention.SentenceDurationDays, detention.ReleaseDate, detention.ReleaseType,
		detention.ReleasingAuthority, detention.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to update detention: %w", err)
	}
	return nil
}
