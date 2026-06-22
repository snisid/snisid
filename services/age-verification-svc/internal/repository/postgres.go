package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/age-verification-svc/internal/domain"
)

type Repository interface {
	InsertAttestation(ctx context.Context, att *domain.AgeAttestation) error
	FindAttestationByID(ctx context.Context, attestationID uuid.UUID) (*domain.AgeAttestation, error)
	FindAttestationByIdentity(ctx context.Context, identityID uuid.UUID) (*domain.AgeAttestation, error)
	UpdateAttestationRevoked(ctx context.Context, attestationID uuid.UUID) error
	InsertAgeClaim(ctx context.Context, claim *domain.AgeClaim) error
	FindClaimsByAttestation(ctx context.Context, attestationID uuid.UUID) ([]domain.AgeClaim, error)
	InsertVerifierRequest(ctx context.Context, req *domain.VerifierRequest) error
	FindVerifierRequest(ctx context.Context, verifierID string, bracket domain.AgeBracket) (*domain.VerifierRequest, error)
	UpdateVerifierRequestApproved(ctx context.Context, requestID uuid.UUID, approved bool) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) InsertAttestation(ctx context.Context, att *domain.AgeAttestation) error {
	query := `INSERT INTO age_attestations (attestation_id, identity_id, date_of_birth, issued_at, expires_at, is_revoked, revoked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		att.AttestationID, att.IdentityID, att.DateOfBirth, att.IssuedAt, att.ExpiresAt, att.IsRevoked, att.RevokedAt,
	)
	if err != nil {
		return fmt.Errorf("insert attestation: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindAttestationByID(ctx context.Context, attestationID uuid.UUID) (*domain.AgeAttestation, error) {
	query := `SELECT attestation_id, identity_id, date_of_birth, issued_at, expires_at, is_revoked, revoked_at
		FROM age_attestations WHERE attestation_id = $1`
	a := &domain.AgeAttestation{}
	err := r.db.QueryRowContext(ctx, query, attestationID).Scan(
		&a.AttestationID, &a.IdentityID, &a.DateOfBirth, &a.IssuedAt, &a.ExpiresAt, &a.IsRevoked, &a.RevokedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("attestation not found: %s", attestationID)
		}
		return nil, fmt.Errorf("query attestation: %w", err)
	}
	return a, nil
}

func (r *postgresRepo) FindAttestationByIdentity(ctx context.Context, identityID uuid.UUID) (*domain.AgeAttestation, error) {
	query := `SELECT attestation_id, identity_id, date_of_birth, issued_at, expires_at, is_revoked, revoked_at
		FROM age_attestations WHERE identity_id = $1 ORDER BY issued_at DESC LIMIT 1`
	a := &domain.AgeAttestation{}
	err := r.db.QueryRowContext(ctx, query, identityID).Scan(
		&a.AttestationID, &a.IdentityID, &a.DateOfBirth, &a.IssuedAt, &a.ExpiresAt, &a.IsRevoked, &a.RevokedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("attestation not found for identity: %s", identityID)
		}
		return nil, fmt.Errorf("query attestation by identity: %w", err)
	}
	return a, nil
}

func (r *postgresRepo) UpdateAttestationRevoked(ctx context.Context, attestationID uuid.UUID) error {
	query := `UPDATE age_attestations SET is_revoked = TRUE, revoked_at = $1 WHERE attestation_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now().UTC(), attestationID)
	if err != nil {
		return fmt.Errorf("update attestation revoked: %w", err)
	}
	return nil
}

func (r *postgresRepo) InsertAgeClaim(ctx context.Context, claim *domain.AgeClaim) error {
	query := `INSERT INTO age_claims (claim_id, attestation_id, verifier_id, bracket, is_satisfied, claimed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		claim.ClaimID, claim.AttestationID, claim.VerifierID, claim.Bracket, claim.IsSatisfied, claim.ClaimedAt,
	)
	if err != nil {
		return fmt.Errorf("insert age claim: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindClaimsByAttestation(ctx context.Context, attestationID uuid.UUID) ([]domain.AgeClaim, error) {
	query := `SELECT claim_id, attestation_id, verifier_id, bracket, is_satisfied, claimed_at
		FROM age_claims WHERE attestation_id = $1 ORDER BY claimed_at DESC`
	rows, err := r.db.QueryContext(ctx, query, attestationID)
	if err != nil {
		return nil, fmt.Errorf("query claims: %w", err)
	}
	defer rows.Close()
	var claims []domain.AgeClaim
	for rows.Next() {
		var c domain.AgeClaim
		if err := rows.Scan(&c.ClaimID, &c.AttestationID, &c.VerifierID, &c.Bracket, &c.IsSatisfied, &c.ClaimedAt); err != nil {
			return nil, err
		}
		claims = append(claims, c)
	}
	return claims, rows.Err()
}

func (r *postgresRepo) InsertVerifierRequest(ctx context.Context, req *domain.VerifierRequest) error {
	query := `INSERT INTO verifier_requests (request_id, verifier_id, bracket, requested_at, is_approved)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		req.RequestID, req.VerifierID, req.Bracket, req.RequestedAt, req.IsApproved,
	)
	if err != nil {
		return fmt.Errorf("insert verifier request: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindVerifierRequest(ctx context.Context, verifierID string, bracket domain.AgeBracket) (*domain.VerifierRequest, error) {
	query := `SELECT request_id, verifier_id, bracket, requested_at, is_approved
		FROM verifier_requests WHERE verifier_id = $1 AND bracket = $2 ORDER BY requested_at DESC LIMIT 1`
	vr := &domain.VerifierRequest{}
	err := r.db.QueryRowContext(ctx, query, verifierID, bracket).Scan(
		&vr.RequestID, &vr.VerifierID, &vr.Bracket, &vr.RequestedAt, &vr.IsApproved,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("verifier request not found")
		}
		return nil, fmt.Errorf("query verifier request: %w", err)
	}
	return vr, nil
}

func (r *postgresRepo) UpdateVerifierRequestApproved(ctx context.Context, requestID uuid.UUID, approved bool) error {
	query := `UPDATE verifier_requests SET is_approved = $1 WHERE request_id = $2`
	_, err := r.db.ExecContext(ctx, query, approved, requestID)
	if err != nil {
		return fmt.Errorf("update verifier request: %w", err)
	}
	return nil
}
