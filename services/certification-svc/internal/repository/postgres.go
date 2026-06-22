package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/certification-svc/internal/domain"
)

type Repository interface {
	CreateProfile(ctx context.Context, profile *domain.AssuranceProfile) error
	FindProfileByIDentityID(ctx context.Context, identityID uuid.UUID) (*domain.AssuranceProfile, error)
	UpdateIAL(ctx context.Context, identityID uuid.UUID, ial domain.IALLevel, updatedBy string) error
	UpdateAAL(ctx context.Context, identityID uuid.UUID, aal domain.AALLevel, updatedBy string) error
	UpdateFAL(ctx context.Context, identityID uuid.UUID, fal domain.FALLevel, updatedBy string) error
	CreateTrustFrameworkClaim(ctx context.Context, claim *domain.TrustFrameworkClaim) error
	FindClaimsByIDentityID(ctx context.Context, identityID uuid.UUID) ([]domain.TrustFrameworkClaim, error)
	CreateAuditEntry(ctx context.Context, audit *domain.CertificationAudit) error
	FindAuditByIDentityID(ctx context.Context, identityID uuid.UUID) ([]domain.CertificationAudit, error)
	FindAllAudit(ctx context.Context) ([]domain.CertificationAudit, error)
	CreateComplianceCheck(ctx context.Context, check *domain.ComplianceCheck) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateProfile(ctx context.Context, profile *domain.AssuranceProfile) error {
	query := `INSERT INTO certification_profiles (profile_id, identity_id, ial, aal, fal, is_active, valid_from, valid_until, last_assessed, assessor_id, assessor_org, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.ExecContext(ctx, query,
		profile.ProfileID, profile.IdentityID, profile.IAL, profile.AAL, profile.FAL,
		profile.IsActive, profile.ValidFrom, profile.ValidUntil, profile.LastAssessed,
		profile.AssessorID, profile.AssessorOrg, profile.CreatedAt, profile.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert certification_profile: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindProfileByIDentityID(ctx context.Context, identityID uuid.UUID) (*domain.AssuranceProfile, error) {
	query := `SELECT profile_id, identity_id, ial, aal, fal, is_active, valid_from, valid_until, last_assessed, assessor_id, assessor_org, created_at, updated_at
		FROM certification_profiles WHERE identity_id = $1`
	p := &domain.AssuranceProfile{}
	err := r.db.QueryRowContext(ctx, query, identityID).Scan(
		&p.ProfileID, &p.IdentityID, &p.IAL, &p.AAL, &p.FAL, &p.IsActive,
		&p.ValidFrom, &p.ValidUntil, &p.LastAssessed, &p.AssessorID,
		&p.AssessorOrg, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("certification profile not found: %s", identityID)
		}
		return nil, fmt.Errorf("query certification_profile: %w", err)
	}
	return p, nil
}

func (r *postgresRepo) UpdateIAL(ctx context.Context, identityID uuid.UUID, ial domain.IALLevel, updatedBy string) error {
	query := `UPDATE certification_profiles SET ial = $1, updated_at = $2 WHERE identity_id = $3`
	_, err := r.db.ExecContext(ctx, query, ial, time.Now().UTC(), identityID)
	if err != nil {
		return fmt.Errorf("update IAL: %w", err)
	}
	return nil
}

func (r *postgresRepo) UpdateAAL(ctx context.Context, identityID uuid.UUID, aal domain.AALLevel, updatedBy string) error {
	query := `UPDATE certification_profiles SET aal = $1, updated_at = $2 WHERE identity_id = $3`
	_, err := r.db.ExecContext(ctx, query, aal, time.Now().UTC(), identityID)
	if err != nil {
		return fmt.Errorf("update AAL: %w", err)
	}
	return nil
}

func (r *postgresRepo) UpdateFAL(ctx context.Context, identityID uuid.UUID, fal domain.FALLevel, updatedBy string) error {
	query := `UPDATE certification_profiles SET fal = $1, updated_at = $2 WHERE identity_id = $3`
	_, err := r.db.ExecContext(ctx, query, fal, time.Now().UTC(), identityID)
	if err != nil {
		return fmt.Errorf("update FAL: %w", err)
	}
	return nil
}

func (r *postgresRepo) CreateTrustFrameworkClaim(ctx context.Context, claim *domain.TrustFrameworkClaim) error {
	query := `INSERT INTO certification_claims (claim_id, identity_id, framework_name, claim_type, claim_value, issuer, issued_at, expires_at, is_verified, verification_ref, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		claim.ClaimID, claim.IdentityID, claim.FrameworkName, claim.ClaimType,
		claim.ClaimValue, claim.Issuer, claim.IssuedAt, claim.ExpiresAt,
		claim.IsVerified, claim.VerificationRef, claim.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert certification_claim: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindClaimsByIDentityID(ctx context.Context, identityID uuid.UUID) ([]domain.TrustFrameworkClaim, error) {
	query := `SELECT claim_id, identity_id, framework_name, claim_type, claim_value, issuer, issued_at, expires_at, is_verified, verification_ref, created_at
		FROM certification_claims WHERE identity_id = $1 ORDER BY issued_at DESC`
	rows, err := r.db.QueryContext(ctx, query, identityID)
	if err != nil {
		return nil, fmt.Errorf("query certification_claims: %w", err)
	}
	defer rows.Close()

	var claims []domain.TrustFrameworkClaim
	for rows.Next() {
		var c domain.TrustFrameworkClaim
		if err := rows.Scan(
			&c.ClaimID, &c.IdentityID, &c.FrameworkName, &c.ClaimType, &c.ClaimValue,
			&c.Issuer, &c.IssuedAt, &c.ExpiresAt, &c.IsVerified, &c.VerificationRef, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		claims = append(claims, c)
	}
	return claims, rows.Err()
}

func (r *postgresRepo) CreateAuditEntry(ctx context.Context, audit *domain.CertificationAudit) error {
	query := `INSERT INTO certification_audit (audit_id, identity_id, action, field, old_value, new_value, performed_by, performed_at, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		audit.AuditID, audit.IdentityID, audit.Action, audit.Field,
		audit.OldValue, audit.NewValue, audit.PerformedBy, audit.PerformedAt, audit.Notes,
	)
	if err != nil {
		return fmt.Errorf("insert certification_audit: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindAuditByIDentityID(ctx context.Context, identityID uuid.UUID) ([]domain.CertificationAudit, error) {
	query := `SELECT audit_id, identity_id, action, field, old_value, new_value, performed_by, performed_at, notes
		FROM certification_audit WHERE identity_id = $1 ORDER BY performed_at DESC`
	rows, err := r.db.QueryContext(ctx, query, identityID)
	if err != nil {
		return nil, fmt.Errorf("query certification_audit: %w", err)
	}
	defer rows.Close()

	var entries []domain.CertificationAudit
	for rows.Next() {
		var e domain.CertificationAudit
		if err := rows.Scan(
			&e.AuditID, &e.IdentityID, &e.Action, &e.Field,
			&e.OldValue, &e.NewValue, &e.PerformedBy, &e.PerformedAt, &e.Notes,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *postgresRepo) FindAllAudit(ctx context.Context) ([]domain.CertificationAudit, error) {
	query := `SELECT audit_id, identity_id, action, field, old_value, new_value, performed_by, performed_at, notes
		FROM certification_audit ORDER BY performed_at DESC LIMIT 100`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all audit: %w", err)
	}
	defer rows.Close()

	var entries []domain.CertificationAudit
	for rows.Next() {
		var e domain.CertificationAudit
		if err := rows.Scan(
			&e.AuditID, &e.IdentityID, &e.Action, &e.Field,
			&e.OldValue, &e.NewValue, &e.PerformedBy, &e.PerformedAt, &e.Notes,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *postgresRepo) CreateComplianceCheck(ctx context.Context, check *domain.ComplianceCheck) error {
	query := `INSERT INTO certification_compliance (check_id, identity_id, check_type, requirement, is_compliant, details, checked_at, checked_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		check.CheckID, check.IdentityID, check.CheckType, check.Requirement,
		check.IsCompliant, check.Details, check.CheckedAt, check.CheckedBy,
	)
	if err != nil {
		return fmt.Errorf("insert certification_compliance: %w", err)
	}
	return nil
}
