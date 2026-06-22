package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/accessibility-svc/internal/domain"
)

type Repository interface {
	CreateAuditRun(ctx context.Context, a *domain.AuditRun) error
	FindAuditRunByID(ctx context.Context, id uuid.UUID) (*domain.AuditRun, error)
	ListAuditRuns(ctx context.Context) ([]domain.AuditRun, error)
	CreateViolation(ctx context.Context, v *domain.Violation) error
	ListViolationsByAudit(ctx context.Context, auditID uuid.UUID) ([]domain.Violation, error)
	MarkViolationRemediated(ctx context.Context, id uuid.UUID) error
	GetComplianceOverview(ctx context.Context) ([]domain.AccessibilityReport, error)
	CreateAuditSchedule(ctx context.Context, s *domain.AuditSchedule) error
	ListAuditSchedules(ctx context.Context) ([]domain.AuditSchedule, error)
	GetDashboard(ctx context.Context) ([]domain.AccessibilityReport, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateAuditRun(ctx context.Context, a *domain.AuditRun) error {
	q := `INSERT INTO acc_audit_runs (audit_run_id, target_url, wcag_level, status, total_violations, passed, failed, started_at, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, q,
		a.AuditRunID, a.TargetURL, a.WCAGLevel, a.Status, a.TotalViolations,
		a.Passed, a.Failed, a.StartedAt, a.CompletedAt, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert audit run: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindAuditRunByID(ctx context.Context, id uuid.UUID) (*domain.AuditRun, error) {
	q := `SELECT audit_run_id, target_url, wcag_level, status, total_violations, passed, failed, started_at, completed_at, created_at
		FROM acc_audit_runs WHERE audit_run_id = $1`
	a := &domain.AuditRun{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&a.AuditRunID, &a.TargetURL, &a.WCAGLevel, &a.Status, &a.TotalViolations,
		&a.Passed, &a.Failed, &a.StartedAt, &a.CompletedAt, &a.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("audit run not found: %s", id)
		}
		return nil, fmt.Errorf("query audit run: %w", err)
	}
	return a, nil
}

func (r *postgresRepo) ListAuditRuns(ctx context.Context) ([]domain.AuditRun, error) {
	q := `SELECT audit_run_id, target_url, wcag_level, status, total_violations, passed, failed, started_at, completed_at, created_at
		FROM acc_audit_runs ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list audit runs: %w", err)
	}
	defer rows.Close()

	var runs []domain.AuditRun
	for rows.Next() {
		var a domain.AuditRun
		if err := rows.Scan(&a.AuditRunID, &a.TargetURL, &a.WCAGLevel, &a.Status, &a.TotalViolations,
			&a.Passed, &a.Failed, &a.StartedAt, &a.CompletedAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		runs = append(runs, a)
	}
	return runs, rows.Err()
}

func (r *postgresRepo) CreateViolation(ctx context.Context, v *domain.Violation) error {
	q := `INSERT INTO acc_violations (violation_id, audit_run_id, wcag_level, guideline, description, element, severity, remediated, remediated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, q,
		v.ViolationID, v.AuditRunID, v.WCAGLevel, v.Guideline, v.Description,
		v.Element, v.Severity, v.Remediated, v.RemediatedAt, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert violation: %w", err)
	}
	return nil
}

func (r *postgresRepo) ListViolationsByAudit(ctx context.Context, auditID uuid.UUID) ([]domain.Violation, error) {
	q := `SELECT violation_id, audit_run_id, wcag_level, guideline, description, element, severity, remediated, remediated_at, created_at
		FROM acc_violations WHERE audit_run_id = $1 ORDER BY severity`
	rows, err := r.db.QueryContext(ctx, q, auditID)
	if err != nil {
		return nil, fmt.Errorf("list violations: %w", err)
	}
	defer rows.Close()

	var violations []domain.Violation
	for rows.Next() {
		var v domain.Violation
		if err := rows.Scan(&v.ViolationID, &v.AuditRunID, &v.WCAGLevel, &v.Guideline, &v.Description,
			&v.Element, &v.Severity, &v.Remediated, &v.RemediatedAt, &v.CreatedAt); err != nil {
			return nil, err
		}
		violations = append(violations, v)
	}
	return violations, rows.Err()
}

func (r *postgresRepo) MarkViolationRemediated(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	q := `UPDATE acc_violations SET remediated = true, remediated_at = $1 WHERE violation_id = $2`
	_, err := r.db.ExecContext(ctx, q, now, id)
	if err != nil {
		return fmt.Errorf("mark remediated: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetComplianceOverview(ctx context.Context) ([]domain.AccessibilityReport, error) {
	q := `SELECT a.audit_run_id, a.target_url, a.wcag_level,
		COUNT(v.violation_id) AS total_violations,
		COUNT(*) FILTER (WHERE v.remediated = true) AS passed,
		COUNT(*) FILTER (WHERE v.remediated = false) AS failed
		FROM acc_audit_runs a LEFT JOIN acc_violations v ON a.audit_run_id = v.audit_run_id
		GROUP BY a.audit_run_id, a.target_url, a.wcag_level ORDER BY a.wcag_level`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("compliance overview: %w", err)
	}
	defer rows.Close()

	var reports []domain.AccessibilityReport
	for rows.Next() {
		var rep domain.AccessibilityReport
		var total, passed, failed int
		if err := rows.Scan(&rep.AuditRunID, &rep.TargetURL, &rep.WCAGLevel, &total, &passed, &failed); err != nil {
			return nil, err
		}
		rep.GeneratedAt = time.Now().UTC()
		if total > 0 {
			rep.PassRate = float64(passed) / float64(total) * 100
		} else {
			rep.PassRate = 100
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

func (r *postgresRepo) CreateAuditSchedule(ctx context.Context, s *domain.AuditSchedule) error {
	q := `INSERT INTO acc_audit_schedules (schedule_id, target_url, wcag_level, cron_expr, enabled, last_run_at, next_run_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, q,
		s.ScheduleID, s.TargetURL, s.WCAGLevel, s.CronExpr, s.Enabled,
		s.LastRunAt, s.NextRunAt, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert schedule: %w", err)
	}
	return nil
}

func (r *postgresRepo) ListAuditSchedules(ctx context.Context) ([]domain.AuditSchedule, error) {
	q := `SELECT schedule_id, target_url, wcag_level, cron_expr, enabled, last_run_at, next_run_at, created_at
		FROM acc_audit_schedules ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list schedules: %w", err)
	}
	defer rows.Close()

	var schedules []domain.AuditSchedule
	for rows.Next() {
		var s domain.AuditSchedule
		if err := rows.Scan(&s.ScheduleID, &s.TargetURL, &s.WCAGLevel, &s.CronExpr, &s.Enabled,
			&s.LastRunAt, &s.NextRunAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, rows.Err()
}

func (r *postgresRepo) GetDashboard(ctx context.Context) ([]domain.AccessibilityReport, error) {
	return r.GetComplianceOverview(ctx)
}
