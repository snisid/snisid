package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/snisid/fisa-court-svc/internal/domain"
)

type Repository interface {
	CreateWarrant(ctx context.Context, w *domain.SurveillanceWarrant) error
	UpdateWarrant(ctx context.Context, w *domain.SurveillanceWarrant) error
	GetWarrant(ctx context.Context, id uuid.UUID) (*domain.SurveillanceWarrant, error)
	GetActiveWarrants(ctx context.Context) ([]domain.SurveillanceWarrant, error)
	CreateReport(ctx context.Context, r *domain.SurveillanceReport) error
	GetDocketByTerm(ctx context.Context, term string) (*domain.FISADocket, error)
	UpsertDocket(ctx context.Context, d *domain.FISADocket) error
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateWarrant(ctx context.Context, w *domain.SurveillanceWarrant) error {
	query := `INSERT INTO fisa_warrants
		(id, warrant_id, warrant_type, target_identity, target_details, issuing_court, judge_name,
		 applicant_agency, applicant_officer, probable_cause_summary, duration_days,
		 authorized_start, authorized_end, renewals, status, review_required_at,
		 emergency_authorized, emergency_approved_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`
	_, err := r.db.ExecContext(ctx, query,
		w.ID, w.WarrantID, w.WarrantType, w.TargetIdentity, w.TargetDetails,
		w.IssuingCourt, w.JudgeName, w.ApplicantAgency, w.ApplicantOfficer,
		w.ProbableCauseSummary, w.DurationDays, w.AuthorizedStart, w.AuthorizedEnd,
		w.Renewals, w.Status, w.ReviewRequiredAt, w.EmergencyAuthorized,
		w.EmergencyApprovedBy, w.CreatedAt, w.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) UpdateWarrant(ctx context.Context, w *domain.SurveillanceWarrant) error {
	query := `UPDATE fisa_warrants SET
		status = $2, authorized_start = $3, authorized_end = $4, renewals = $5,
		review_required_at = $6, emergency_authorized = $7, emergency_approved_by = $8,
		updated_at = $9 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		w.ID, w.Status, w.AuthorizedStart, w.AuthorizedEnd, w.Renewals,
		w.ReviewRequiredAt, w.EmergencyAuthorized, w.EmergencyApprovedBy, w.UpdatedAt,
	)
	return err
}

func (r *postgresRepo) GetWarrant(ctx context.Context, id uuid.UUID) (*domain.SurveillanceWarrant, error) {
	query := `SELECT id, warrant_id, warrant_type, target_identity, target_details, issuing_court, judge_name,
		applicant_agency, applicant_officer, probable_cause_summary, duration_days,
		authorized_start, authorized_end, renewals, status, review_required_at,
		emergency_authorized, emergency_approved_by, created_at, updated_at
		FROM fisa_warrants WHERE id = $1`
	w := &domain.SurveillanceWarrant{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.WarrantID, &w.WarrantType, &w.TargetIdentity, &w.TargetDetails,
		&w.IssuingCourt, &w.JudgeName, &w.ApplicantAgency, &w.ApplicantOfficer,
		&w.ProbableCauseSummary, &w.DurationDays, &w.AuthorizedStart, &w.AuthorizedEnd,
		&w.Renewals, &w.Status, &w.ReviewRequiredAt, &w.EmergencyAuthorized,
		&w.EmergencyApprovedBy, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("warrant not found")
		}
		return nil, err
	}
	return w, nil
}

func (r *postgresRepo) GetActiveWarrants(ctx context.Context) ([]domain.SurveillanceWarrant, error) {
	query := `SELECT id, warrant_id, warrant_type, target_identity, target_details, issuing_court, judge_name,
		applicant_agency, applicant_officer, probable_cause_summary, duration_days,
		authorized_start, authorized_end, renewals, status, review_required_at,
		emergency_authorized, emergency_approved_by, created_at, updated_at
		FROM fisa_warrants WHERE status IN ('PENDING','APPROVED','ACTIVE') ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warrants []domain.SurveillanceWarrant
	for rows.Next() {
		var w domain.SurveillanceWarrant
		if err := rows.Scan(
			&w.ID, &w.WarrantID, &w.WarrantType, &w.TargetIdentity, &w.TargetDetails,
			&w.IssuingCourt, &w.JudgeName, &w.ApplicantAgency, &w.ApplicantOfficer,
			&w.ProbableCauseSummary, &w.DurationDays, &w.AuthorizedStart, &w.AuthorizedEnd,
			&w.Renewals, &w.Status, &w.ReviewRequiredAt, &w.EmergencyAuthorized,
			&w.EmergencyApprovedBy, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		warrants = append(warrants, w)
	}
	return warrants, rows.Err()
}

func (r *postgresRepo) CreateReport(ctx context.Context, rep *domain.SurveillanceReport) error {
	query := `INSERT INTO fisa_reports
		(id, warrant_id, reporting_period_start, reporting_period_end, communications_intercepted,
		 minimization_applied, incidental_collection, us_person_identities, results_summary,
		 submitted_by, submitted_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.db.ExecContext(ctx, query,
		rep.ID, rep.WarrantID, rep.ReportingPeriodStart, rep.ReportingPeriodEnd,
		rep.CommunicationsIntercepted, rep.MinimizationApplied, rep.IncidentalCollection,
		rep.USPersonIdentities, rep.ResultsSummary, rep.SubmittedBy, rep.SubmittedAt,
	)
	return err
}

func (r *postgresRepo) GetDocketByTerm(ctx context.Context, term string) (*domain.FISADocket, error) {
	query := `SELECT id, docket_number, court_term, judge_presiding, applications_filed,
		applications_approved, applications_modified, applications_denied, total_targets,
		foreign_targets, us_person_targets, sealed_until, created_at
		FROM fisa_dockets WHERE court_term = $1`
	d := &domain.FISADocket{}
	err := r.db.QueryRowContext(ctx, query, term).Scan(
		&d.ID, &d.DocketNumber, &d.CourtTerm, &d.JudgePresiding,
		&d.ApplicationsFiled, &d.ApplicationsApproved, &d.ApplicationsModified,
		&d.ApplicationsDenied, &d.TotalTargets, &d.ForeignTargets,
		&d.USPersonTargets, &d.SealedUntil, &d.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("docket not found")
		}
		return nil, err
	}
	return d, nil
}

func (r *postgresRepo) UpsertDocket(ctx context.Context, d *domain.FISADocket) error {
	query := `INSERT INTO fisa_dockets
		(id, docket_number, court_term, judge_presiding, applications_filed, applications_approved,
		 applications_modified, applications_denied, total_targets, foreign_targets,
		 us_person_targets, sealed_until, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		ON CONFLICT (court_term) DO UPDATE SET
		applications_filed = EXCLUDED.applications_filed,
		applications_approved = EXCLUDED.applications_approved,
		applications_modified = EXCLUDED.applications_modified,
		applications_denied = EXCLUDED.applications_denied,
		total_targets = EXCLUDED.total_targets,
		foreign_targets = EXCLUDED.foreign_targets,
		us_person_targets = EXCLUDED.us_person_targets`
	_, err := r.db.ExecContext(ctx, query,
		d.ID, d.DocketNumber, d.CourtTerm, d.JudgePresiding,
		d.ApplicationsFiled, d.ApplicationsApproved, d.ApplicationsModified,
		d.ApplicationsDenied, d.TotalTargets, d.ForeignTargets,
		d.USPersonTargets, d.SealedUntil, d.CreatedAt,
	)
	return err
}
