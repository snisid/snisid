package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/snisid/bug-bounty-svc/internal/domain"
)

type Repository interface {
	CreateProgram(ctx context.Context, scope *domain.ProgramScope) error
	ListPrograms(ctx context.Context) ([]domain.ProgramScope, error)
	CreateReport(ctx context.Context, r *domain.VulnerabilityReport) error
	FindReportByID(ctx context.Context, id uuid.UUID) (*domain.VulnerabilityReport, error)
	SaveTriageResult(ctx context.Context, t *domain.TriageResult) error
	IssueReward(ctx context.Context, r *domain.Reward) error
	CreatePentestEngagement(ctx context.Context, e *domain.PentestEngagement) error
	FindPentestByID(ctx context.Context, id uuid.UUID) (*domain.PentestEngagement, error)
	SaveRetestSchedule(ctx context.Context, s *domain.RetestSchedule) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateProgram(ctx context.Context, scope *domain.ProgramScope) error {
	q := `INSERT INTO bb_programs (scope_id, program_id, target, scope_type, in_scope, reward_min, reward_max)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, q,
		scope.ScopeID, scope.ProgramID, scope.Target, scope.ScopeType,
		scope.InScope, scope.RewardMin, scope.RewardMax,
	)
	if err != nil {
		return fmt.Errorf("insert program: %w", err)
	}
	return nil
}

func (r *postgresRepo) ListPrograms(ctx context.Context) ([]domain.ProgramScope, error) {
	q := `SELECT scope_id, program_id, target, scope_type, in_scope, reward_min, reward_max FROM bb_programs WHERE in_scope = true ORDER BY target`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list programs: %w", err)
	}
	defer rows.Close()

	var programs []domain.ProgramScope
	for rows.Next() {
		var p domain.ProgramScope
		if err := rows.Scan(&p.ScopeID, &p.ProgramID, &p.Target, &p.ScopeType, &p.InScope, &p.RewardMin, &p.RewardMax); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, rows.Err()
}

func (r *postgresRepo) CreateReport(ctx context.Context, rep *domain.VulnerabilityReport) error {
	q := `INSERT INTO bb_reports (report_id, program_id, submitter, title, description, severity, scope_id, status, submitted_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, q,
		rep.ReportID, rep.ProgramID, rep.Submitter, rep.Title, rep.Description,
		rep.Severity, rep.ScopeID, rep.Status, time.Now().UTC(), time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert report: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindReportByID(ctx context.Context, id uuid.UUID) (*domain.VulnerabilityReport, error) {
	q := `SELECT report_id, program_id, submitter, title, description, severity, scope_id, status, submitted_at, updated_at
		FROM bb_reports WHERE report_id = $1`
	rep := &domain.VulnerabilityReport{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&rep.ReportID, &rep.ProgramID, &rep.Submitter, &rep.Title, &rep.Description,
		&rep.Severity, &rep.ScopeID, &rep.Status, &rep.SubmittedAt, &rep.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("report not found: %s", id)
		}
		return nil, fmt.Errorf("query report: %w", err)
	}
	return rep, nil
}

func (r *postgresRepo) SaveTriageResult(ctx context.Context, t *domain.TriageResult) error {
	q := `INSERT INTO bb_triage_results (triage_id, report_id, triager, severity, reproducible, duplicate_of, notes, triaged_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, q,
		t.TriageID, t.ReportID, t.Triager, t.Severity, t.Reproducible, t.DuplicateOf, t.Notes, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert triage: %w", err)
	}

	upd := `UPDATE bb_reports SET status = 'TRIAGED', updated_at = $1 WHERE report_id = $2`
	r.db.ExecContext(ctx, upd, time.Now().UTC(), t.ReportID)
	return nil
}

func (r *postgresRepo) IssueReward(ctx context.Context, rew *domain.Reward) error {
	q := `INSERT INTO bb_rewards (reward_id, report_id, amount, currency, paid_to, approved_by, paid_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, q,
		rew.RewardID, rew.ReportID, rew.Amount, rew.Currency, rew.PaidTo, rew.ApprovedBy, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert reward: %w", err)
	}

	upd := `UPDATE bb_reports SET status = 'REWARDED', updated_at = $1 WHERE report_id = $2`
	r.db.ExecContext(ctx, upd, time.Now().UTC(), rew.ReportID)
	return nil
}

func (r *postgresRepo) CreatePentestEngagement(ctx context.Context, e *domain.PentestEngagement) error {
	q := `INSERT INTO bb_pentest_engagements (engagement_id, program_id, title, scope, start_date, end_date, team_lead, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, q,
		e.EngagementID, e.ProgramID, e.Title, e.Scope, e.StartDate, e.EndDate,
		e.TeamLead, e.Status, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert pentest: %w", err)
	}
	return nil
}

func (r *postgresRepo) FindPentestByID(ctx context.Context, id uuid.UUID) (*domain.PentestEngagement, error) {
	q := `SELECT engagement_id, program_id, title, scope, start_date, end_date, team_lead, status, created_at
		FROM bb_pentest_engagements WHERE engagement_id = $1`
	e := &domain.PentestEngagement{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&e.EngagementID, &e.ProgramID, &e.Title, &e.Scope, &e.StartDate,
		&e.EndDate, &e.TeamLead, &e.Status, &e.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pentest not found: %s", id)
		}
		return nil, fmt.Errorf("query pentest: %w", err)
	}
	return e, nil
}

func (r *postgresRepo) SaveRetestSchedule(ctx context.Context, s *domain.RetestSchedule) error {
	q := `INSERT INTO bb_retest_schedules (schedule_id, report_id, scheduled_for, completed_at, assigned_to, status)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, q,
		s.ScheduleID, s.ReportID, s.ScheduledFor, s.CompletedAt, s.AssignedTo, s.Status,
	)
	if err != nil {
		return fmt.Errorf("insert retest: %w", err)
	}
	return nil
}
