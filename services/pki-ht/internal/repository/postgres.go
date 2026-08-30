package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/pki-ht/internal/domain"
)

type Repository interface {
	CreateCertificate(ctx context.Context, cert *domain.IssuedCertificate) error
	FindBySerial(ctx context.Context, serial string) (*domain.IssuedCertificate, error)
	RevokeCertificate(ctx context.Context, serial string, reason string) error
	GetActiveCRL(ctx context.Context, caID uuid.UUID) (*domain.CRL, error)
	UpdateCRL(ctx context.Context, crl *domain.CRL) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateCertificate(ctx context.Context, cert *domain.IssuedCertificate) error {
	q := `INSERT INTO pki_issued_certificates (cert_id, serial_number, issuing_ca_id, subject_type, subject_ref, common_name, status, valid_from, valid_until, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, q, cert.CertID, cert.SerialNumber, cert.IssuingCAID,
		cert.SubjectType, cert.SubjectRef, cert.CommonName, cert.Status,
		cert.ValidFrom, cert.ValidUntil, time.Now().UTC())
	return err
}

func (r *postgresRepo) FindBySerial(ctx context.Context, serial string) (*domain.IssuedCertificate, error) {
	q := `SELECT cert_id, serial_number, issuing_ca_id, subject_type, subject_ref, common_name, status, valid_from, valid_until, revoked_at, revocation_reason, created_at FROM pki_issued_certificates WHERE serial_number = $1`
	c := &domain.IssuedCertificate{}
	err := r.db.QueryRowContext(ctx, q, serial).Scan(
		&c.CertID, &c.SerialNumber, &c.IssuingCAID, &c.SubjectType, &c.SubjectRef,
		&c.CommonName, &c.Status, &c.ValidFrom, &c.ValidUntil, &c.RevokedAt,
		&c.RevocationReason, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("certificate not found")
		}
		return nil, err
	}
	return c, nil
}

func (r *postgresRepo) RevokeCertificate(ctx context.Context, serial string, reason string) error {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `UPDATE pki_issued_certificates SET status = 'REVOKED', revoked_at = $1, revocation_reason = $2 WHERE serial_number = $3`, now, reason, serial)
	return err
}

func (r *postgresRepo) GetActiveCRL(ctx context.Context, caID uuid.UUID) (*domain.CRL, error) {
	q := `SELECT crl_id, ca_id, crl_number, revoked_serials, published_at, next_update FROM pki_crl WHERE ca_id = $1 ORDER BY published_at DESC LIMIT 1`
	crl := &domain.CRL{}
	err := r.db.QueryRowContext(ctx, q, caID).Scan(&crl.CRLID, &crl.CAID, &crl.CRLNumber, &crl.RevokedSerials, &crl.PublishedAt, &crl.NextUpdate)
	if err != nil {
		return nil, fmt.Errorf("no CRL found for CA %s", caID)
	}
	return crl, nil
}

func (r *postgresRepo) UpdateCRL(ctx context.Context, crl *domain.CRL) error {
	q := `INSERT INTO pki_crl (crl_id, ca_id, crl_number, revoked_serials, published_at, next_update) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.ExecContext(ctx, q, crl.CRLID, crl.CAID, crl.CRLNumber, crl.RevokedSerials, crl.PublishedAt, crl.NextUpdate)
	return err
}
