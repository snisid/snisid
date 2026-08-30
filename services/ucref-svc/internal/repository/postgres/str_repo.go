package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/snisid/platform/services/ucref-svc/internal/domain"
)

type strRepo struct {
	pool *pgxpool.Pool
}

func NewSTRRepository(pool *pgxpool.Pool) *strRepo {
	return &strRepo{pool: pool}
}

func (r *strRepo) CreateSTR(report *domain.STRReport) error {
	ctx := context.Background()
	report.StrID = uuid.New()
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()
	if report.TransactionCurrency == "" {
		report.TransactionCurrency = "HTG"
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO ucref_str_reports (
			str_id, national_str_id, report_type, status, reporting_institution,
			institution_type, report_date, transaction_date, transaction_amount,
			transaction_currency, transaction_amount_usd, subject_snisid_ids,
			subject_names, subject_accounts, suspicious_activity, ml_typology,
			predicate_crime, gang_id, fpr_person_ids, sanc_match_ids, analyst_id,
			analysis_notes, disseminated_to, created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25
		)`,
		report.StrID, report.NationalStrID, report.ReportType, report.Status,
		report.ReportingInstitution, report.InstitutionType, report.ReportDate,
		report.TransactionDate, report.TransactionAmount, report.TransactionCurrency,
		report.TransactionAmountUSD, report.SubjectSnisidIDs, report.SubjectNames,
		report.SubjectAccounts, report.SuspiciousActivity, report.MLTypology,
		report.PredicateCrime, report.GangID, report.FPRPersonIDs, report.SancMatchIDs,
		report.AnalystID, report.AnalysisNotes, report.DisseminatedTo,
		report.CreatedAt, report.UpdatedAt,
	)
	return err
}

func (r *strRepo) FindByID(id uuid.UUID) (*domain.STRReport, error) {
	ctx := context.Background()
	row := r.pool.QueryRow(ctx,
		`SELECT str_id, national_str_id, report_type, status, reporting_institution,
			institution_type, report_date, transaction_date, transaction_amount,
			transaction_currency, transaction_amount_usd,
			COALESCE(array_to_string(subject_snisid_ids, ','), ''),
			COALESCE(array_to_string(subject_names, ','), ''),
			COALESCE(array_to_string(subject_accounts, ','), ''),
			suspicious_activity, COALESCE(ml_typology, ''),
			COALESCE(predicate_crime, ''), gang_id,
			COALESCE(array_to_string(fpr_person_ids, ','), ''),
			COALESCE(array_to_string(sanc_match_ids, ','), ''),
			analyst_id, COALESCE(analysis_notes, ''),
			COALESCE(array_to_string(disseminated_to, ','), ''),
			disseminated_at, created_at, updated_at
		FROM ucref_str_reports WHERE str_id = $1`, id)
	return scanSTRReport(row)
}

func (r *strRepo) GetFinancialProfile(personID uuid.UUID) (*domain.FinancialProfile, error) {
	ctx := context.Background()
	profile := &domain.FinancialProfile{}

	var knownBusinesses string
	err := r.pool.QueryRow(ctx,
		`SELECT profile_id, snisid_person_id, total_str_count, total_ctr_count,
			COALESCE(estimated_illegal_assets_usd, 0), known_accounts, known_properties,
			COALESCE(array_to_string(known_businesses, ','), ''),
			COALESCE(ml_risk_score, 0), is_pep, last_updated
		FROM ucref_financial_profiles WHERE snisid_person_id = $1`, personID).Scan(
		&profile.ProfileID, &profile.SNISIDPersonID, &profile.TotalSTRCount,
		&profile.TotalCTRCount, &profile.EstimatedIllegalAssetsUSD,
		&profile.KnownAccounts, &profile.KnownProperties, &knownBusinesses,
		&profile.MLRiskScore, &profile.IsPEP, &profile.LastUpdated,
	)
	if err != nil {
		return nil, err
	}

	if knownBusinesses != "" {
		profile.KnownBusinesses = strings.Split(knownBusinesses, ",")
	}

	return profile, nil
}

func (r *strRepo) CreateMonCashPattern(pattern *domain.MonCashPattern) error {
	ctx := context.Background()
	pattern.PatternID = uuid.New()
	pattern.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO ucref_moncash_patterns (
			pattern_id, str_id, phone_number, snisid_person_id, pattern_type,
			transaction_count, total_amount_htg, period_start, period_end, notes, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		pattern.PatternID, pattern.STRID, pattern.PhoneNumber, pattern.SNISIDPersonID,
		pattern.PatternType, pattern.TransactionCount, pattern.TotalAmountHTG,
		pattern.PeriodStart, pattern.PeriodEnd, pattern.Notes, pattern.CreatedAt,
	)
	return err
}

func (r *strRepo) GetUnanalyzedSTRs() ([]domain.STRReport, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT str_id, national_str_id, report_type, status, reporting_institution,
			institution_type, report_date, transaction_date, transaction_amount,
			transaction_currency, transaction_amount_usd,
			COALESCE(array_to_string(subject_snisid_ids, ','), ''),
			COALESCE(array_to_string(subject_names, ','), ''),
			COALESCE(array_to_string(subject_accounts, ','), ''),
			suspicious_activity, COALESCE(ml_typology, ''),
			COALESCE(predicate_crime, ''), gang_id,
			COALESCE(array_to_string(fpr_person_ids, ','), ''),
			COALESCE(array_to_string(sanc_match_ids, ','), ''),
			analyst_id, COALESCE(analysis_notes, ''),
			COALESCE(array_to_string(disseminated_to, ','), ''),
			disseminated_at, created_at, updated_at
		FROM ucref_str_reports
		WHERE status = 'RECEIVED'
		ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []domain.STRReport
	for rows.Next() {
		report, err := scanSTRReport(rows)
		if err != nil {
			return nil, err
		}
		reports = append(reports, *report)
	}
	return reports, nil
}

func (r *strRepo) DisseminateSTR(id uuid.UUID, disseminatedTo []string) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		`UPDATE ucref_str_reports
		SET status = 'DISSEMINATED', disseminated_to = $1, disseminated_at = NOW(), updated_at = NOW()
		WHERE str_id = $2`,
		disseminatedTo, id,
	)
	return err
}

func (r *strRepo) GetGangFinances(gangID uuid.UUID) ([]domain.FinancialProfile, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx,
		`SELECT fp.profile_id, fp.snisid_person_id, fp.total_str_count, fp.total_ctr_count,
			COALESCE(fp.estimated_illegal_assets_usd, 0), fp.known_accounts, fp.known_properties,
			COALESCE(array_to_string(fp.known_businesses, ','), ''),
			COALESCE(fp.ml_risk_score, 0), fp.is_pep, fp.last_updated
		FROM ucref_financial_profiles fp
		INNER JOIN ucref_str_reports sr ON sr.subject_snisid_ids @> ARRAY[fp.snisid_person_id]
		WHERE sr.gang_id = $1
		GROUP BY fp.profile_id, fp.snisid_person_id, fp.total_str_count, fp.total_ctr_count,
			fp.estimated_illegal_assets_usd, fp.known_accounts, fp.known_properties,
			fp.known_businesses, fp.ml_risk_score, fp.is_pep, fp.last_updated`, gangID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []domain.FinancialProfile
	for rows.Next() {
		profile := domain.FinancialProfile{}
		var knownBusinesses string

		err := rows.Scan(
			&profile.ProfileID, &profile.SNISIDPersonID, &profile.TotalSTRCount,
			&profile.TotalCTRCount, &profile.EstimatedIllegalAssetsUSD,
			&profile.KnownAccounts, &profile.KnownProperties, &knownBusinesses,
			&profile.MLRiskScore, &profile.IsPEP, &profile.LastUpdated,
		)
		if err != nil {
			return nil, err
		}

		if knownBusinesses != "" {
			profile.KnownBusinesses = strings.Split(knownBusinesses, ",")
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (r *strRepo) GetNextSequence(year string) (int, error) {
	ctx := context.Background()
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ucref_str_reports
		WHERE national_str_id LIKE $1`,
		"STR-HT-"+year+"-%",
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count + 1, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanSTRReport(row scannable) (*domain.STRReport, error) {
	report := &domain.STRReport{}
	var subjectSnisidIDs, subjectNames, subjectAccounts, fprPersonIDs, sancMatchIDs, disseminatedTo string

	err := row.Scan(
		&report.StrID, &report.NationalStrID, &report.ReportType, &report.Status,
		&report.ReportingInstitution, &report.InstitutionType, &report.ReportDate,
		&report.TransactionDate, &report.TransactionAmount, &report.TransactionCurrency,
		&report.TransactionAmountUSD, &subjectSnisidIDs, &subjectNames, &subjectAccounts,
		&report.SuspiciousActivity, &report.MLTypology, &report.PredicateCrime,
		&report.GangID, &fprPersonIDs, &sancMatchIDs, &report.AnalystID,
		&report.AnalysisNotes, &disseminatedTo, &report.DisseminatedAt,
		&report.CreatedAt, &report.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	report.SubjectSnisidIDs = splitCSV(subjectSnisidIDs)
	report.SubjectNames = splitCSV(subjectNames)
	report.SubjectAccounts = splitCSV(subjectAccounts)
	report.FPRPersonIDs = splitCSV(fprPersonIDs)
	report.SancMatchIDs = splitCSV(sancMatchIDs)
	report.DisseminatedTo = splitCSV(disseminatedTo)

	return report, nil
}

func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
