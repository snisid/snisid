package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/fir-svc/internal/domain"
)

type CriminalRecordRepo struct {
	pool *pgxpool.Pool
}

func NewCriminalRecordRepo(pool *pgxpool.Pool) *CriminalRecordRepo {
	return &CriminalRecordRepo{pool: pool}
}

func (r *CriminalRecordRepo) Create(ctx context.Context, record *domain.CriminalRecord) error {
	query := `
		INSERT INTO fir_criminal_records 
			(record_id, national_fir_id, snisid_person_id, afis_subject_id, 
			 is_haitian_national, aliases, is_active, is_expunged, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		record.RecordID, record.NationalFIRID, record.SNISIDPersonID,
		record.AfisSubjectID, record.IsHaitianNational, record.Aliases,
		record.IsActive, record.IsExpunged, record.CreatedAt, record.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create criminal record: %w", err)
	}
	return nil
}

func (r *CriminalRecordRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.CriminalRecord, error) {
	query := `
		SELECT record_id, national_fir_id, snisid_person_id, afis_subject_id,
			   is_haitian_national, aliases, is_active, is_expunged, created_at, updated_at
		FROM fir_criminal_records
		WHERE record_id = $1
	`
	record := &domain.CriminalRecord{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&record.RecordID, &record.NationalFIRID, &record.SNISIDPersonID,
		&record.AfisSubjectID, &record.IsHaitianNational, &record.Aliases,
		&record.IsActive, &record.IsExpunged, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find criminal record: %w", err)
	}
	return record, nil
}

func (r *CriminalRecordRepo) FindByPersonID(ctx context.Context, personID uuid.UUID) (*domain.CriminalRecord, error) {
	query := `
		SELECT record_id, national_fir_id, snisid_person_id, afis_subject_id,
			   is_haitian_national, aliases, is_active, is_expunged, created_at, updated_at
		FROM fir_criminal_records
		WHERE snisid_person_id = $1
	`
	record := &domain.CriminalRecord{}
	err := r.pool.QueryRow(ctx, query, personID).Scan(
		&record.RecordID, &record.NationalFIRID, &record.SNISIDPersonID,
		&record.AfisSubjectID, &record.IsHaitianNational, &record.Aliases,
		&record.IsActive, &record.IsExpunged, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find criminal record by person: %w", err)
	}
	return record, nil
}

func (r *CriminalRecordRepo) FindByFIRID(ctx context.Context, firID string) (*domain.CriminalRecord, error) {
	query := `
		SELECT record_id, national_fir_id, snisid_person_id, afis_subject_id,
			   is_haitian_national, aliases, is_active, is_expunged, created_at, updated_at
		FROM fir_criminal_records
		WHERE national_fir_id = $1
	`
	record := &domain.CriminalRecord{}
	err := r.pool.QueryRow(ctx, query, firID).Scan(
		&record.RecordID, &record.NationalFIRID, &record.SNISIDPersonID,
		&record.AfisSubjectID, &record.IsHaitianNational, &record.Aliases,
		&record.IsActive, &record.IsExpunged, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find criminal record by FIR ID: %w", err)
	}
	return record, nil
}

func (r *CriminalRecordRepo) Update(ctx context.Context, record *domain.CriminalRecord) error {
	query := `
		UPDATE fir_criminal_records
		SET afis_subject_id = $3, is_haitian_national = $4, aliases = $5,
			is_active = $6, is_expunged = $7, updated_at = $8
		WHERE record_id = $1 AND snisid_person_id = $2
	`
	record.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query,
		record.RecordID, record.SNISIDPersonID, record.AfisSubjectID,
		record.IsHaitianNational, record.Aliases, record.IsActive,
		record.IsExpunged, record.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update criminal record: %w", err)
	}
	return nil
}

func (r *CriminalRecordRepo) NextSequence(ctx context.Context) (int64, error) {
	var nextVal int64
	err := r.pool.QueryRow(ctx, "SELECT nextval('fir_record_seq')").Scan(&nextVal)
	if err != nil {
		return 0, fmt.Errorf("failed to get next sequence: %w", err)
	}
	return nextVal, nil
}

func (r *CriminalRecordRepo) SaveCertificate(ctx context.Context, cert *domain.Certificate) error {
	query := `
		INSERT INTO fir_certificates 
			(cert_id, record_id, snisid_person_id, certificate_number, 
			 issued_for, result, issued_by, issuing_office, issued_at, expires_at, qr_code_ref)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.pool.Exec(ctx, query,
		cert.CertID, cert.RecordID, cert.SNISIDPersonID, cert.CertificateNumber,
		cert.IssuedFor, cert.Result, cert.IssuedBy, cert.IssuingOffice,
		cert.IssuedAt, cert.ExpiresAt, cert.QRCodeRef,
	)
	if err != nil {
		return fmt.Errorf("failed to save certificate: %w", err)
	}
	return nil
}

func (r *CriminalRecordRepo) FindCertificateByNumber(ctx context.Context, num string) (*domain.Certificate, error) {
	query := `
		SELECT cert_id, record_id, snisid_person_id, certificate_number,
			   issued_for, result, issued_by, issuing_office, issued_at, expires_at, qr_code_ref
		FROM fir_certificates
		WHERE certificate_number = $1
	`
	cert := &domain.Certificate{}
	err := r.pool.QueryRow(ctx, query, num).Scan(
		&cert.CertID, &cert.RecordID, &cert.SNISIDPersonID, &cert.CertificateNumber,
		&cert.IssuedFor, &cert.Result, &cert.IssuedBy, &cert.IssuingOffice,
		&cert.IssuedAt, &cert.ExpiresAt, &cert.QRCodeRef,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find certificate: %w", err)
	}
	return cert, nil
}
