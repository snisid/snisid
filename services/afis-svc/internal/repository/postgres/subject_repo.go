package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/snisid/platform/services/afis-svc/internal/domain"
)

type SubjectRepo struct {
	pool *pgxpool.Pool
}

func NewSubjectRepo(pool *pgxpool.Pool) *SubjectRepo {
	return &SubjectRepo{pool: pool}
}

func (r *SubjectRepo) Create(ctx context.Context, s *domain.SubjectProfile) error {
	query := `INSERT INTO afis_subjects (subject_id, snisid_person_id, fir_record_id, subject_type, national_afis_id, enrolling_unit, enrolling_officer)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, s.SubjectID, s.SNISIDPersonID, s.FIRRecordID, s.SubjectType, s.NationalAFISID, s.EnrollingUnit, uuid.Nil)
	if err != nil {
		return fmt.Errorf("insert subject: %w", err)
	}
	return nil
}

func (r *SubjectRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.SubjectProfile, error) {
	query := `SELECT subject_id, snisid_person_id, fir_record_id, subject_type, national_afis_id, enrolling_unit FROM afis_subjects WHERE subject_id = $1`
	s := &domain.SubjectProfile{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.SubjectID, &s.SNISIDPersonID, &s.FIRRecordID, &s.SubjectType, &s.NationalAFISID, &s.EnrollingUnit,
	)
	if err != nil {
		return nil, fmt.Errorf("get subject: %w", err)
	}
	return s, nil
}

func (r *SubjectRepo) GetBySNISIDPersonID(ctx context.Context, personID uuid.UUID) (*domain.SubjectProfile, error) {
	query := `SELECT subject_id, snisid_person_id, fir_record_id, subject_type, national_afis_id, enrolling_unit FROM afis_subjects WHERE snisid_person_id = $1`
	s := &domain.SubjectProfile{}
	err := r.pool.QueryRow(ctx, query, personID).Scan(
		&s.SubjectID, &s.SNISIDPersonID, &s.FIRRecordID, &s.SubjectType, &s.NationalAFISID, &s.EnrollingUnit,
	)
	if err != nil {
		return nil, fmt.Errorf("get subject by snisid: %w", err)
	}
	return s, nil
}
