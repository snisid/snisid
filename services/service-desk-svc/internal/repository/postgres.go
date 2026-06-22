package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/snisid/service-desk-svc/internal/domain"
)

type Repository interface {
	CreateCase(ctx context.Context, c *domain.SupportCase) error
	GetCase(ctx context.Context, caseID uuid.UUID) (*domain.SupportCase, error)
	ListCases(ctx context.Context, status domain.CaseStatus) ([]domain.SupportCase, error)
	UpdateCaseStatus(ctx context.Context, caseID uuid.UUID, status domain.CaseStatus) error
	CreateChallenge(ctx context.Context, ch *domain.VerificationChallenge) error
	ResolveChallenge(ctx context.Context, challengeID uuid.UUID) error
	CreateRecoveryRequest(ctx context.Context, req *domain.IdentityRecoveryRequest) error
	VerifyRecoveryRequest(ctx context.Context, requestID uuid.UUID) error
	AddCaseNote(ctx context.Context, note *domain.CaseNote) error
	CreateResolution(ctx context.Context, res *domain.Resolution) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateCase(ctx context.Context, c *domain.SupportCase) error {
	query := `INSERT INTO support_cases (case_id, citizen_id, subject, description, status, assigned_to, created_at, updated_at, resolved_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		c.CaseID, c.CitizenID, c.Subject, c.Description, c.Status,
		c.AssignedTo, c.CreatedAt, c.UpdatedAt, c.ResolvedAt,
	)
	if err != nil {
		return fmt.Errorf("insert support_case: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetCase(ctx context.Context, caseID uuid.UUID) (*domain.SupportCase, error) {
	query := `SELECT case_id, citizen_id, subject, description, status, assigned_to, created_at, updated_at, resolved_at
		FROM support_cases WHERE case_id = $1`
	c := &domain.SupportCase{}
	err := r.db.QueryRowContext(ctx, query, caseID).Scan(
		&c.CaseID, &c.CitizenID, &c.Subject, &c.Description, &c.Status,
		&c.AssignedTo, &c.CreatedAt, &c.UpdatedAt, &c.ResolvedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("case not found: %s", caseID)
		}
		return nil, fmt.Errorf("query case: %w", err)
	}
	return c, nil
}

func (r *postgresRepo) ListCases(ctx context.Context, status domain.CaseStatus) ([]domain.SupportCase, error) {
	query := `SELECT case_id, citizen_id, subject, description, status, assigned_to, created_at, updated_at, resolved_at
		FROM support_cases WHERE status = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("query cases: %w", err)
	}
	defer rows.Close()

	var cases []domain.SupportCase
	for rows.Next() {
		var c domain.SupportCase
		if err := rows.Scan(&c.CaseID, &c.CitizenID, &c.Subject, &c.Description, &c.Status,
			&c.AssignedTo, &c.CreatedAt, &c.UpdatedAt, &c.ResolvedAt); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, rows.Err()
}

func (r *postgresRepo) UpdateCaseStatus(ctx context.Context, caseID uuid.UUID, status domain.CaseStatus) error {
	query := `UPDATE support_cases SET status = $1, updated_at = $2 WHERE case_id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now().UTC(), caseID)
	if err != nil {
		return fmt.Errorf("update case status: %w", err)
	}
	return nil
}

func (r *postgresRepo) CreateChallenge(ctx context.Context, ch *domain.VerificationChallenge) error {
	query := `INSERT INTO verification_challenges (challenge_id, case_id, method, challenge, expires_at, is_resolved, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		ch.ChallengeID, ch.CaseID, ch.Method, ch.Challenge, ch.ExpiresAt, ch.IsResolved, ch.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert challenge: %w", err)
	}
	return nil
}

func (r *postgresRepo) ResolveChallenge(ctx context.Context, challengeID uuid.UUID) error {
	query := `UPDATE verification_challenges SET is_resolved = TRUE WHERE challenge_id = $1`
	_, err := r.db.ExecContext(ctx, query, challengeID)
	if err != nil {
		return fmt.Errorf("resolve challenge: %w", err)
	}
	return nil
}

func (r *postgresRepo) CreateRecoveryRequest(ctx context.Context, req *domain.IdentityRecoveryRequest) error {
	query := `INSERT INTO identity_recovery_requests (request_id, case_id, citizen_id, preferred_method, verified_methods, is_verified, created_at, resolved_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		req.RequestID, req.CaseID, req.CitizenID, req.PreferredMethod,
		req.VerifiedMethods, req.IsVerified, req.CreatedAt, req.ResolvedAt,
	)
	if err != nil {
		return fmt.Errorf("insert recovery request: %w", err)
	}
	return nil
}

func (r *postgresRepo) VerifyRecoveryRequest(ctx context.Context, requestID uuid.UUID) error {
	query := `UPDATE identity_recovery_requests SET is_verified = TRUE, resolved_at = $1 WHERE request_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now().UTC(), requestID)
	if err != nil {
		return fmt.Errorf("verify recovery request: %w", err)
	}
	return nil
}

func (r *postgresRepo) AddCaseNote(ctx context.Context, note *domain.CaseNote) error {
	query := `INSERT INTO case_notes (note_id, case_id, author, content, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, note.NoteID, note.CaseID, note.Author, note.Content, note.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert case note: %w", err)
	}
	return nil
}

func (r *postgresRepo) CreateResolution(ctx context.Context, res *domain.Resolution) error {
	query := `INSERT INTO case_resolutions (resolution_id, case_id, action, details, resolved_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		res.ResolutionID, res.CaseID, res.Action, res.Details, res.ResolvedBy, res.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert resolution: %w", err)
	}
	return nil
}
