package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/pki-bridge-svc/internal/domain"
)

type Repository interface {
	CreateForeignCA(ctx context.Context, ca *domain.ForeignCA) error
	FindForeignCAByID(ctx context.Context, id uuid.UUID) (*domain.ForeignCA, error)
	ListForeignCAs(ctx context.Context) ([]domain.ForeignCA, error)
	CreateCrossCert(ctx context.Context, cert *domain.CrossCertificate) error
	FindCrossCertBySubject(ctx context.Context, subject string) (*domain.CrossCertificate, error)
	ListTrustAnchors(ctx context.Context) ([]domain.TrustAnchor, error)
	SavePathValidation(ctx context.Context, pv *domain.PathValidation) error
	CreateBridgeAgreement(ctx context.Context, a *domain.BridgeAgreement) error
	ListBridgeAgreements(ctx context.Context) ([]domain.BridgeAgreement, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateForeignCA(ctx context.Context, ca *domain.ForeignCA) error {
	q := `INSERT INTO pki_foreign_cas (ca_id, name, country, public_key_pem, cert_policy, registered_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, q,
		ca.CAID, ca.Name, ca.Country, ca.PublicKeyPEM, ca.CertPolicy, time.Now().UTC(), ca.Status,
	)
	if err != nil {
		return fmt.Errorf("insert foreign ca: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindForeignCAByID(ctx context.Context, id uuid.UUID) (*domain.ForeignCA, error) {
	q := `SELECT ca_id, name, country, public_key_pem, cert_policy, registered_at, status FROM pki_foreign_cas WHERE ca_id = $1`
	ca := &domain.ForeignCA{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&ca.CAID, &ca.Name, &ca.Country, &ca.PublicKeyPEM, &ca.CertPolicy, &ca.RegisteredAt, &ca.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("foreign ca not found: %s", id)
		}
		return nil, fmt.Errorf("query foreign ca: %w", err)
	}
	return ca, nil
}

func (r *postgresRepo) ListForeignCAs(ctx context.Context) ([]domain.ForeignCA, error) {
	q := `SELECT ca_id, name, country, public_key_pem, cert_policy, registered_at, status FROM pki_foreign_cas ORDER BY name`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list foreign cas: %w", err)
	}
	defer rows.Close()

	var cas []domain.ForeignCA
	for rows.Next() {
		var ca domain.ForeignCA
		if err := rows.Scan(&ca.CAID, &ca.Name, &ca.Country, &ca.PublicKeyPEM, &ca.CertPolicy, &ca.RegisteredAt, &ca.Status); err != nil {
			return nil, err
		}
		cas = append(cas, ca)
	}
	return cas, rows.Err()
}

func (r *postgresRepo) CreateCrossCert(ctx context.Context, cert *domain.CrossCertificate) error {
	q := `INSERT INTO pki_cross_certs (cross_cert_id, subject, issuer_ca_id, serial_number, not_before, not_after, certificate_pem, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, q,
		cert.CrossCertID, cert.Subject, cert.IssuerCAID, cert.SerialNumber,
		cert.NotBefore, cert.NotAfter, cert.CertificatePEM, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert cross cert: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindCrossCertBySubject(ctx context.Context, subject string) (*domain.CrossCertificate, error) {
	q := `SELECT cross_cert_id, subject, issuer_ca_id, serial_number, not_before, not_after, certificate_pem, created_at
		FROM pki_cross_certs WHERE subject = $1 ORDER BY created_at DESC LIMIT 1`
	cert := &domain.CrossCertificate{}
	err := r.db.QueryRowContext(ctx, q, subject).Scan(
		&cert.CrossCertID, &cert.Subject, &cert.IssuerCAID, &cert.SerialNumber,
		&cert.NotBefore, &cert.NotAfter, &cert.CertificatePEM, &cert.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cross cert not found: %s", subject)
		}
		return nil, fmt.Errorf("query cross cert: %w", err)
	}
	return cert, nil
}

func (r *postgresRepo) ListTrustAnchors(ctx context.Context) ([]domain.TrustAnchor, error) {
	q := `SELECT anchor_id, subject, public_key_pem, added_at, expires_at FROM pki_trust_anchors ORDER BY subject`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list trust anchors: %w", err)
	}
	defer rows.Close()

	var anchors []domain.TrustAnchor
	for rows.Next() {
		var a domain.TrustAnchor
		if err := rows.Scan(&a.AnchorID, &a.Subject, &a.PublicKeyPEM, &a.AddedAt, &a.ExpiresAt); err != nil {
			return nil, err
		}
		anchors = append(anchors, a)
	}
	return anchors, rows.Err()
}

func (r *postgresRepo) SavePathValidation(ctx context.Context, pv *domain.PathValidation) error {
	q := `INSERT INTO pki_path_validations (validation_id, path_id, result, errors, validated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, q,
		pv.ValidationID, pv.PathID, pv.Result, stringSlice(pv.Errors), time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert path validation: %w", err)
	}
	return nil
}

func (r *postgresRepo) CreateBridgeAgreement(ctx context.Context, a *domain.BridgeAgreement) error {
	q := `INSERT INTO pki_bridge_agreements (agreement_id, name, partner_ca, policy_id, signed_at, expires_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, q,
		a.AgreementID, a.Name, a.PartnerCA, a.PolicyID, time.Now().UTC(), a.ExpiresAt, a.Status,
	)
	if err != nil {
		return fmt.Errorf("insert bridge agreement: %w", err)
	}
	return nil
}

func (r *postgresRepo) ListBridgeAgreements(ctx context.Context) ([]domain.BridgeAgreement, error) {
	q := `SELECT agreement_id, name, partner_ca, policy_id, signed_at, expires_at, status FROM pki_bridge_agreements ORDER BY signed_at DESC`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list agreements: %w", err)
	}
	defer rows.Close()

	var agreements []domain.BridgeAgreement
	for rows.Next() {
		var a domain.BridgeAgreement
		if err := rows.Scan(&a.AgreementID, &a.Name, &a.PartnerCA, &a.PolicyID, &a.SignedAt, &a.ExpiresAt, &a.Status); err != nil {
			return nil, err
		}
		agreements = append(agreements, a)
	}
	return agreements, rows.Err()
}

func stringSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
