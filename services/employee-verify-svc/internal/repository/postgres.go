package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/employee-verify-svc/internal/domain"
)

type Repository interface {
	InsertEmployer(ctx context.Context, emp *domain.EmployerRegistration) error
	FindEmployerByEIN(ctx context.Context, ein string) (*domain.EmployerRegistration, error)
	FindEmployerByID(ctx context.Context, employerID uuid.UUID) (*domain.EmployerRegistration, error)
	InsertCase(ctx context.Context, vreq *domain.VerificationRequest) error
	FindCaseByTCN(ctx context.Context, tcn string) (*domain.VerificationRequest, error)
	FindCasesByEmployer(ctx context.Context, employerEIN string) ([]domain.VerificationRequest, error)
	UpdateCaseStatus(ctx context.Context, tcn string, status domain.VerificationStatus) error
	InsertVerificationResult(ctx context.Context, result *domain.VerificationResult) error
	FindVerificationResult(ctx context.Context, tcn string) (*domain.VerificationResult, error)
	InsertCaseHistory(ctx context.Context, history *domain.CaseHistory) error
	FindCaseHistory(ctx context.Context, tcn string) ([]domain.CaseHistory, error)
	GetStats(ctx context.Context) (map[string]int, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) InsertEmployer(ctx context.Context, emp *domain.EmployerRegistration) error {
	query := `INSERT INTO employers (employer_id, company_name, ein, address, contact_email, contact_phone, registered_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		emp.EmployerID, emp.CompanyName, emp.EIN, emp.Address, emp.ContactEmail, emp.ContactPhone, emp.RegisteredAt, emp.IsActive,
	)
	if err != nil {
		return fmt.Errorf("insert employer: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindEmployerByEIN(ctx context.Context, ein string) (*domain.EmployerRegistration, error) {
	query := `SELECT employer_id, company_name, ein, address, contact_email, contact_phone, registered_at, is_active
		FROM employers WHERE ein = $1`
	e := &domain.EmployerRegistration{}
	err := r.db.QueryRowContext(ctx, query, ein).Scan(
		&e.EmployerID, &e.CompanyName, &e.EIN, &e.Address, &e.ContactEmail, &e.ContactPhone, &e.RegisteredAt, &e.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employer not found: %s", ein)
		}
		return nil, fmt.Errorf("query employer: %w", err)
	}
	return e, nil
}

func (r *postgresRepo) FindEmployerByID(ctx context.Context, employerID uuid.UUID) (*domain.EmployerRegistration, error) {
	query := `SELECT employer_id, company_name, ein, address, contact_email, contact_phone, registered_at, is_active
		FROM employers WHERE employer_id = $1`
	e := &domain.EmployerRegistration{}
	err := r.db.QueryRowContext(ctx, query, employerID).Scan(
		&e.EmployerID, &e.CompanyName, &e.EIN, &e.Address, &e.ContactEmail, &e.ContactPhone, &e.RegisteredAt, &e.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employer not found: %s", employerID)
		}
		return nil, fmt.Errorf("query employer by id: %w", err)
	}
	return e, nil
}

func (r *postgresRepo) InsertCase(ctx context.Context, vreq *domain.VerificationRequest) error {
	query := `INSERT INTO verification_cases (tcn, employer_id, employee_name, document_number, document_type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		vreq.TCN, vreq.EmployerID, vreq.EmployeeName, vreq.DocumentNumber, vreq.DocumentType, vreq.Status, vreq.CreatedAt, vreq.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert case: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindCaseByTCN(ctx context.Context, tcn string) (*domain.VerificationRequest, error) {
	query := `SELECT tcn, employer_id, employee_name, document_number, document_type, status, created_at, updated_at
		FROM verification_cases WHERE tcn = $1`
	v := &domain.VerificationRequest{}
	err := r.db.QueryRowContext(ctx, query, tcn).Scan(
		&v.TCN, &v.EmployerID, &v.EmployeeName, &v.DocumentNumber, &v.DocumentType, &v.Status, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("case not found: %s", tcn)
		}
		return nil, fmt.Errorf("query case: %w", err)
	}
	return v, nil
}

func (r *postgresRepo) FindCasesByEmployer(ctx context.Context, employerEIN string) ([]domain.VerificationRequest, error) {
	query := `SELECT v.tcn, v.employer_id, v.employee_name, v.document_number, v.document_type, v.status, v.created_at, v.updated_at
		FROM verification_cases v JOIN employers e ON v.employer_id = e.employer_id
		WHERE e.ein = $1 ORDER BY v.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, employerEIN)
	if err != nil {
		return nil, fmt.Errorf("query cases by employer: %w", err)
	}
	defer rows.Close()
	var cases []domain.VerificationRequest
	for rows.Next() {
		var c domain.VerificationRequest
		if err := rows.Scan(&c.TCN, &c.EmployerID, &c.EmployeeName, &c.DocumentNumber, &c.DocumentType, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, rows.Err()
}

func (r *postgresRepo) UpdateCaseStatus(ctx context.Context, tcn string, status domain.VerificationStatus) error {
	query := `UPDATE verification_cases SET status = $1, updated_at = $2 WHERE tcn = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now().UTC(), tcn)
	if err != nil {
		return fmt.Errorf("update case status: %w", err)
	}
	return nil
}

func (r *postgresRepo) InsertVerificationResult(ctx context.Context, result *domain.VerificationResult) error {
	query := `INSERT INTO verification_results (result_id, tcn, ssa_match, dhs_match, is_eligible, reason, completed_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		result.ResultID, result.TCN, result.SSAMatch, result.DHSMatch, result.IsEligible, result.Reason, result.CompletedAt, result.Status,
	)
	if err != nil {
		return fmt.Errorf("insert verification result: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindVerificationResult(ctx context.Context, tcn string) (*domain.VerificationResult, error) {
	query := `SELECT result_id, tcn, ssa_match, dhs_match, is_eligible, reason, completed_at, status
		FROM verification_results WHERE tcn = $1 ORDER BY completed_at DESC LIMIT 1`
	res := &domain.VerificationResult{}
	err := r.db.QueryRowContext(ctx, query, tcn).Scan(
		&res.ResultID, &res.TCN, &res.SSAMatch, &res.DHSMatch, &res.IsEligible, &res.Reason, &res.CompletedAt, &res.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("result not found: %s", tcn)
		}
		return nil, fmt.Errorf("query result: %w", err)
	}
	return res, nil
}

func (r *postgresRepo) InsertCaseHistory(ctx context.Context, history *domain.CaseHistory) error {
	query := `INSERT INTO case_history (history_id, tcn, action, actioned_by, actioned_at, details)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		history.HistoryID, history.TCN, history.Action, history.ActionedBy, history.ActionedAt, history.Details,
	)
	if err != nil {
		return fmt.Errorf("insert case history: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindCaseHistory(ctx context.Context, tcn string) ([]domain.CaseHistory, error) {
	query := `SELECT history_id, tcn, action, actioned_by, actioned_at, details
		FROM case_history WHERE tcn = $1 ORDER BY actioned_at ASC`
	rows, err := r.db.QueryContext(ctx, query, tcn)
	if err != nil {
		return nil, fmt.Errorf("query case history: %w", err)
	}
	defer rows.Close()
	var hist []domain.CaseHistory
	for rows.Next() {
		var h domain.CaseHistory
		if err := rows.Scan(&h.HistoryID, &h.TCN, &h.Action, &h.ActionedBy, &h.ActionedAt, &h.Details); err != nil {
			return nil, err
		}
		hist = append(hist, h)
	}
	return hist, rows.Err()
}

func (r *postgresRepo) GetStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM verification_cases`).Scan(&stats["total_cases"])
	if err != nil {
		return nil, fmt.Errorf("query total cases: %w", err)
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM verification_cases WHERE status = 'VERIFIED'`).Scan(&stats["verified"])
	if err != nil {
		return nil, fmt.Errorf("query verified: %w", err)
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM verification_cases WHERE status = 'NOT_VERIFIED'`).Scan(&stats["not_verified"])
	if err != nil {
		return nil, fmt.Errorf("query not verified: %w", err)
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM employers`).Scan(&stats["total_employers"])
	if err != nil {
		return nil, fmt.Errorf("query total employers: %w", err)
	}
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM verification_cases WHERE status = 'PENDING'`).Scan(&stats["pending"])
	if err != nil {
		return nil, fmt.Errorf("query pending: %w", err)
	}
	return stats, nil
}
