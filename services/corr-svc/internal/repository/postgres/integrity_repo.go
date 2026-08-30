package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/corr-svc/internal/domain"
)

type integrityRepo struct {
	pool *pgxpool.Pool
}

func NewIntegrityRepo(pool *pgxpool.Pool) *integrityRepo {
	return &integrityRepo{pool: pool}
}

func (r *integrityRepo) CreateCase(c *domain.IntegrityCase) (*domain.IntegrityCase, error) {
	ctx := context.Background()
	c.ID = uuid.New()
	c.NationalCorrID = "CORR-HT-" + time.Now().Format("2006") + "-" + c.ID.String()[:6]
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO corr_integrity_cases
		 (case_id, national_corr_id, officer_snisid_id, officer_badge, officer_unit, officer_rank,
		  allegation_type, severity, status, allegation_summary, gang_id, reporting_date, created_by, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		c.ID, c.NationalCorrID, c.OfficerSNISIDID, c.OfficerBadge, c.OfficerUnit,
		c.OfficerRank, c.AllegationType, c.Severity, c.Status, c.AllegationSummary,
		c.GangID, c.ReportingDate, c.CreatedBy, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *integrityRepo) FindByID(id uuid.UUID) (*domain.IntegrityCase, error) {
	ctx := context.Background()
	c := &domain.IntegrityCase{}
	err := r.pool.QueryRow(ctx,
		`SELECT case_id, national_corr_id, officer_snisid_id, allegation_type, severity, status,
		        allegation_summary, reporting_date, created_at
		 FROM corr_integrity_cases WHERE case_id = $1`, id).Scan(
		&c.ID, &c.NationalCorrID, &c.OfficerSNISIDID, &c.AllegationType, &c.Severity,
		&c.Status, &c.AllegationSummary, &c.ReportingDate, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *integrityRepo) FindActive() ([]domain.IntegrityCase, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT case_id, national_corr_id, officer_snisid_id, allegation_type, severity, status,
		        allegation_summary, reporting_date, created_at
		 FROM corr_integrity_cases
		 WHERE status IN ('REPORTED','UNDER_INVESTIGATION')
		 ORDER BY reporting_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cases []domain.IntegrityCase
	for rows.Next() {
		var c domain.IntegrityCase
		if err := rows.Scan(&c.ID, &c.NationalCorrID, &c.OfficerSNISIDID, &c.AllegationType,
			&c.Severity, &c.Status, &c.AllegationSummary, &c.ReportingDate, &c.CreatedAt); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, nil
}

func (r *integrityRepo) CreateWhistleblowerReport(wr *domain.WhistleblowerReport) (*domain.WhistleblowerReport, error) {
	ctx := context.Background()
	wr.ID = uuid.New()
	wr.SubmissionDate = time.Now()
	wr.CreatedAt = time.Now()
	isProcessed := false
	wr.Processed = &isProcessed

	_, err := r.pool.Exec(ctx,
		`INSERT INTO corr_whistleblower_reports
		 (report_id, report_token, allegation_type, severity_estimate, description,
		  submission_date, processed, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		wr.ID, wr.ReportToken, wr.AllegationType, wr.SeverityEstimate,
		wr.Description, wr.SubmissionDate, wr.Processed, wr.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return wr, nil
}

func (r *integrityRepo) FindByToken(token string) (*domain.WhistleblowerReport, error) {
	ctx := context.Background()
	wr := &domain.WhistleblowerReport{}
	err := r.pool.QueryRow(ctx,
		`SELECT report_id, report_token, allegation_type, description, submission_date, processed, created_at
		 FROM corr_whistleblower_reports WHERE report_token = $1`, token).Scan(
		&wr.ID, &wr.ReportToken, &wr.AllegationType, &wr.Description,
		&wr.SubmissionDate, &wr.Processed, &wr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return wr, nil
}

func (r *integrityRepo) FindBehavioralAlerts() ([]domain.BehavioralAlert, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT alert_id, officer_snisid_id, alert_type, description, module_source, risk_score, created_at
		 FROM corr_behavioral_alerts
		 WHERE reviewed = FALSE
		 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.BehavioralAlert
	for rows.Next() {
		var a domain.BehavioralAlert
		if err := rows.Scan(&a.ID, &a.OfficerSNISIDID, &a.AlertType, &a.Description,
			&a.ModuleSource, &a.RiskScore, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r *integrityRepo) CreateAssetDeclaration(ad *domain.AssetDeclaration) (*domain.AssetDeclaration, error) {
	ctx := context.Background()
	ad.ID = uuid.New()
	ad.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO corr_asset_declarations
		 (declaration_id, officer_snisid_id, declaration_year, real_estate_usd, vehicles_usd,
		  bank_accounts_usd, other_assets_usd, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		ad.ID, ad.OfficerSNISIDID, ad.DeclarationYear, ad.RealEstateUSD, ad.VehiclesUSD,
		ad.BankAccountsUSD, ad.OtherAssetsUSD, ad.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (r *integrityRepo) FindFlaggedDeclarations() ([]domain.AssetDeclaration, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT declaration_id, officer_snisid_id, declaration_year, real_estate_usd, vehicles_usd,
		        bank_accounts_usd, other_assets_usd, is_flagged, created_at
		 FROM corr_asset_declarations WHERE is_flagged = TRUE ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var declarations []domain.AssetDeclaration
	for rows.Next() {
		var d domain.AssetDeclaration
		if err := rows.Scan(&d.ID, &d.OfficerSNISIDID, &d.DeclarationYear, &d.RealEstateUSD,
			&d.VehiclesUSD, &d.BankAccountsUSD, &d.OtherAssetsUSD, &d.IsFlagged, &d.CreatedAt); err != nil {
			return nil, err
		}
		declarations = append(declarations, d)
	}
	return declarations, nil
}
