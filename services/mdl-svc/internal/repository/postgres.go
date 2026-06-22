package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/mdl-svc/internal/domain"
)

type Repository interface {
	CreateIssuance(ctx context.Context, issuance *domain.MDLIssuance) error
	FindIssuanceByID(ctx context.Context, issuanceID uuid.UUID) (*domain.MDLIssuance, error)
	FindIssuanceByIdentity(ctx context.Context, identityID uuid.UUID) (*domain.MDLIssuance, error)
	UpdateIssuanceRevoked(ctx context.Context, issuanceID uuid.UUID) error
	InsertPresentation(ctx context.Context, pres *domain.MDLPresentation) error
	FindPresentationsByIssuance(ctx context.Context, issuanceID uuid.UUID) ([]domain.MDLPresentation, error)
	FindPresentationByID(ctx context.Context, presID uuid.UUID) (*domain.MDLPresentation, error)
	UpdatePresentationVerification(ctx context.Context, presID uuid.UUID, verified bool, result string) error
	InsertDataElement(ctx context.Context, elem *domain.MDLDataElement) error
	FindDataElementsByIssuance(ctx context.Context, issuanceID uuid.UUID) ([]domain.MDLDataElement, error)
	InsertDeviceEngagement(ctx context.Context, eng *domain.DeviceEngagement) error
	FindEngagementByIssuance(ctx context.Context, issuanceID uuid.UUID) (*domain.DeviceEngagement, error)
	InsertQRBarcode(ctx context.Context, qr *domain.QRBarcode) error
	FindQRByEngagement(ctx context.Context, engagementID uuid.UUID) (*domain.QRBarcode, error)
	InsertTrustRegistry(ctx context.Context, entry *domain.MDLTrustRegistry) error
	FindTrustRegistry(ctx context.Context) ([]domain.MDLTrustRegistry, error)
	FindTrustEntryByReaderID(ctx context.Context, readerID string) (*domain.MDLTrustRegistry, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateIssuance(ctx context.Context, issuance *domain.MDLIssuance) error {
	query := `INSERT INTO mdl_issuances (issuance_id, identity_id, device_id, issued_at, expires_at, is_revoked, revoked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		issuance.IssuanceID, issuance.IdentityID, issuance.DeviceID,
		issuance.IssuedAt, issuance.ExpiresAt, issuance.IsRevoked, issuance.RevokedAt,
	)
	if err != nil {
		return fmt.Errorf("insert issuance: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindIssuanceByID(ctx context.Context, issuanceID uuid.UUID) (*domain.MDLIssuance, error) {
	query := `SELECT issuance_id, identity_id, device_id, issued_at, expires_at, is_revoked, revoked_at
		FROM mdl_issuances WHERE issuance_id = $1`
	i := &domain.MDLIssuance{}
	err := r.db.QueryRowContext(ctx, query, issuanceID).Scan(
		&i.IssuanceID, &i.IdentityID, &i.DeviceID, &i.IssuedAt, &i.ExpiresAt, &i.IsRevoked, &i.RevokedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("issuance not found: %s", issuanceID)
		}
		return nil, fmt.Errorf("query issuance: %w", err)
	}
	return i, nil
}

func (r *postgresRepo) FindIssuanceByIdentity(ctx context.Context, identityID uuid.UUID) (*domain.MDLIssuance, error) {
	query := `SELECT issuance_id, identity_id, device_id, issued_at, expires_at, is_revoked, revoked_at
		FROM mdl_issuances WHERE identity_id = $1 ORDER BY issued_at DESC LIMIT 1`
	i := &domain.MDLIssuance{}
	err := r.db.QueryRowContext(ctx, query, identityID).Scan(
		&i.IssuanceID, &i.IdentityID, &i.DeviceID, &i.IssuedAt, &i.ExpiresAt, &i.IsRevoked, &i.RevokedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("issuance not found for identity: %s", identityID)
		}
		return nil, fmt.Errorf("query issuance by identity: %w", err)
	}
	return i, nil
}

func (r *postgresRepo) UpdateIssuanceRevoked(ctx context.Context, issuanceID uuid.UUID) error {
	query := `UPDATE mdl_issuances SET is_revoked = TRUE, revoked_at = $1 WHERE issuance_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now().UTC(), issuanceID)
	if err != nil {
		return fmt.Errorf("update issuance revoked: %w", err)
	}
	return nil
}

func (r *postgresRepo) InsertPresentation(ctx context.Context, pres *domain.MDLPresentation) error {
	query := `INSERT INTO mdl_presentations (presentation_id, issuance_id, reader_id, presented_at, is_verified, verification_result)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		pres.PresentationID, pres.IssuanceID, pres.ReaderID,
		pres.PresentedAt, pres.IsVerified, pres.VerificationResult,
	)
	if err != nil {
		return fmt.Errorf("insert presentation: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindPresentationsByIssuance(ctx context.Context, issuanceID uuid.UUID) ([]domain.MDLPresentation, error) {
	query := `SELECT presentation_id, issuance_id, reader_id, presented_at, is_verified, verification_result
		FROM mdl_presentations WHERE issuance_id = $1 ORDER BY presented_at DESC`
	rows, err := r.db.QueryContext(ctx, query, issuanceID)
	if err != nil {
		return nil, fmt.Errorf("query presentations: %w", err)
	}
	defer rows.Close()
	var pres []domain.MDLPresentation
	for rows.Next() {
		var p domain.MDLPresentation
		if err := rows.Scan(&p.PresentationID, &p.IssuanceID, &p.ReaderID, &p.PresentedAt, &p.IsVerified, &p.VerificationResult); err != nil {
			return nil, err
		}
		pres = append(pres, p)
	}
	return pres, rows.Err()
}

func (r *postgresRepo) FindPresentationByID(ctx context.Context, presID uuid.UUID) (*domain.MDLPresentation, error) {
	query := `SELECT presentation_id, issuance_id, reader_id, presented_at, is_verified, verification_result
		FROM mdl_presentations WHERE presentation_id = $1`
	p := &domain.MDLPresentation{}
	err := r.db.QueryRowContext(ctx, query, presID).Scan(
		&p.PresentationID, &p.IssuanceID, &p.ReaderID, &p.PresentedAt, &p.IsVerified, &p.VerificationResult,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("presentation not found: %s", presID)
		}
		return nil, fmt.Errorf("query presentation: %w", err)
	}
	return p, nil
}

func (r *postgresRepo) UpdatePresentationVerification(ctx context.Context, presID uuid.UUID, verified bool, result string) error {
	query := `UPDATE mdl_presentations SET is_verified = $1, verification_result = $2 WHERE presentation_id = $3`
	_, err := r.db.ExecContext(ctx, query, verified, result, presID)
	if err != nil {
		return fmt.Errorf("update presentation verification: %w", err)
	}
	return nil
}

func (r *postgresRepo) InsertDataElement(ctx context.Context, elem *domain.MDLDataElement) error {
	query := `INSERT INTO mdl_data_elements (element_id, issuance_id, element_name, element_value, is_mandatory)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		elem.ElementID, elem.IssuanceID, elem.ElementName, elem.ElementValue, elem.IsMandatory,
	)
	if err != nil {
		return fmt.Errorf("insert data element: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindDataElementsByIssuance(ctx context.Context, issuanceID uuid.UUID) ([]domain.MDLDataElement, error) {
	query := `SELECT element_id, issuance_id, element_name, element_value, is_mandatory
		FROM mdl_data_elements WHERE issuance_id = $1`
	rows, err := r.db.QueryContext(ctx, query, issuanceID)
	if err != nil {
		return nil, fmt.Errorf("query data elements: %w", err)
	}
	defer rows.Close()
	var elems []domain.MDLDataElement
	for rows.Next() {
		var e domain.MDLDataElement
		if err := rows.Scan(&e.ElementID, &e.IssuanceID, &e.ElementName, &e.ElementValue, &e.IsMandatory); err != nil {
			return nil, err
		}
		elems = append(elems, e)
	}
	return elems, rows.Err()
}

func (r *postgresRepo) InsertDeviceEngagement(ctx context.Context, eng *domain.DeviceEngagement) error {
	query := `INSERT INTO mdl_device_engagements (engagement_id, issuance_id, qr_payload, engagement_code, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		eng.EngagementID, eng.IssuanceID, eng.QRPayload, eng.EngagementCode, eng.CreatedAt, eng.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert device engagement: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindEngagementByIssuance(ctx context.Context, issuanceID uuid.UUID) (*domain.DeviceEngagement, error) {
	query := `SELECT engagement_id, issuance_id, qr_payload, engagement_code, created_at, expires_at
		FROM mdl_device_engagements WHERE issuance_id = $1 ORDER BY created_at DESC LIMIT 1`
	e := &domain.DeviceEngagement{}
	err := r.db.QueryRowContext(ctx, query, issuanceID).Scan(
		&e.EngagementID, &e.IssuanceID, &e.QRPayload, &e.EngagementCode, &e.CreatedAt, &e.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("engagement not found for issuance: %s", issuanceID)
		}
		return nil, fmt.Errorf("query engagement: %w", err)
	}
	return e, nil
}

func (r *postgresRepo) InsertQRBarcode(ctx context.Context, qr *domain.QRBarcode) error {
	query := `INSERT INTO mdl_qr_barcodes (barcode_id, engagement_id, encoded_data, format, generated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		qr.BarcodeID, qr.EngagementID, qr.EncodedData, qr.Format, qr.GeneratedAt,
	)
	if err != nil {
		return fmt.Errorf("insert qr barcode: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindQRByEngagement(ctx context.Context, engagementID uuid.UUID) (*domain.QRBarcode, error) {
	query := `SELECT barcode_id, engagement_id, encoded_data, format, generated_at
		FROM mdl_qr_barcodes WHERE engagement_id = $1 ORDER BY generated_at DESC LIMIT 1`
	q := &domain.QRBarcode{}
	err := r.db.QueryRowContext(ctx, query, engagementID).Scan(
		&q.BarcodeID, &q.EngagementID, &q.EncodedData, &q.Format, &q.GeneratedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("qr barcode not found for engagement: %s", engagementID)
		}
		return nil, fmt.Errorf("query qr barcode: %w", err)
	}
	return q, nil
}

func (r *postgresRepo) InsertTrustRegistry(ctx context.Context, entry *domain.MDLTrustRegistry) error {
	query := `INSERT INTO mdl_trust_registry (entry_id, reader_id, reader_name, public_key, is_trusted, registered_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		entry.EntryID, entry.ReaderID, entry.ReaderName, entry.PublicKey, entry.IsTrusted, entry.RegisteredAt, entry.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert trust registry: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindTrustRegistry(ctx context.Context) ([]domain.MDLTrustRegistry, error) {
	query := `SELECT entry_id, reader_id, reader_name, public_key, is_trusted, registered_at, expires_at
		FROM mdl_trust_registry WHERE is_trusted = TRUE ORDER BY registered_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query trust registry: %w", err)
	}
	defer rows.Close()
	var entries []domain.MDLTrustRegistry
	for rows.Next() {
		var e domain.MDLTrustRegistry
		if err := rows.Scan(&e.EntryID, &e.ReaderID, &e.ReaderName, &e.PublicKey, &e.IsTrusted, &e.RegisteredAt, &e.ExpiresAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *postgresRepo) FindTrustEntryByReaderID(ctx context.Context, readerID string) (*domain.MDLTrustRegistry, error) {
	query := `SELECT entry_id, reader_id, reader_name, public_key, is_trusted, registered_at, expires_at
		FROM mdl_trust_registry WHERE reader_id = $1`
	e := &domain.MDLTrustRegistry{}
	err := r.db.QueryRowContext(ctx, query, readerID).Scan(
		&e.EntryID, &e.ReaderID, &e.ReaderName, &e.PublicKey, &e.IsTrusted, &e.RegisteredAt, &e.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trust entry not found: %s", readerID)
		}
		return nil, fmt.Errorf("query trust entry: %w", err)
	}
	return e, nil
}
