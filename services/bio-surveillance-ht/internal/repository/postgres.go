package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/snisid/bio-surveillance-ht/internal/domain"
)

type Repository interface {
	CreateAlert(ctx context.Context, a *domain.DiseaseAlert) error
	GetActiveAlerts(ctx context.Context) ([]domain.DiseaseAlert, error)
	GetAlertsByRegion(ctx context.Context, region string) ([]domain.DiseaseAlert, error)
	CreateCampaign(ctx context.Context, c *domain.VaccinationCampaign) error
	GetCampaignCoverage(ctx context.Context, id uuid.UUID) (*domain.VaccinationCampaign, error)
	UpdateFacilityStock(ctx context.Context, id uuid.UUID, req domain.UpdateFacilityStockRequest) (*domain.HealthFacility, error)
	GetDashboardNational(ctx context.Context) (*domain.DashboardNational, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateAlert(ctx context.Context, a *domain.DiseaseAlert) error {
	query := `INSERT INTO biosurv_disease_alerts (id, disease_name, pathogen_type, icd10_code, alert_level, first_case_detected_at, symptoms_hallmark, transmission_mode, incubation_days, fatality_rate, cases_confirmed, cases_suspected, cases_deaths, affected_regions, source_lab, who_alert_ref, containment_measures, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.DiseaseName, a.PathogenType, a.Icd10Code, a.AlertLevel, a.FirstCaseDetected,
		a.SymptomsHallmark, a.TransmissionMode, a.IncubationDays, a.FatalityRate,
		a.CasesConfirmed, a.CasesSuspected, a.CasesDeaths, pq.Array(a.AffectedRegions),
		a.SourceLab, a.WhoAlertRef, a.ContainmentMeasures, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert alert: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetActiveAlerts(ctx context.Context) ([]domain.DiseaseAlert, error) {
	query := `SELECT id, disease_name, pathogen_type, icd10_code, alert_level, first_case_detected_at, symptoms_hallmark, transmission_mode, incubation_days, fatality_rate, cases_confirmed, cases_suspected, cases_deaths, affected_regions, source_lab, who_alert_ref, containment_measures, created_at
		FROM biosurv_disease_alerts WHERE alert_level IN ('YELLOW', 'ORANGE', 'RED') ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query active alerts: %w", err)
	}
	defer rows.Close()

	var alerts []domain.DiseaseAlert
	for rows.Next() {
		var a domain.DiseaseAlert
		if err := rows.Scan(
			&a.ID, &a.DiseaseName, &a.PathogenType, &a.Icd10Code, &a.AlertLevel, &a.FirstCaseDetected,
			&a.SymptomsHallmark, &a.TransmissionMode, &a.IncubationDays, &a.FatalityRate,
			&a.CasesConfirmed, &a.CasesSuspected, &a.CasesDeaths, pq.Array(&a.AffectedRegions),
			&a.SourceLab, &a.WhoAlertRef, &a.ContainmentMeasures, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (r *postgresRepo) GetAlertsByRegion(ctx context.Context, region string) ([]domain.DiseaseAlert, error) {
	query := `SELECT id, disease_name, pathogen_type, icd10_code, alert_level, first_case_detected_at, symptoms_hallmark, transmission_mode, incubation_days, fatality_rate, cases_confirmed, cases_suspected, cases_deaths, affected_regions, source_lab, who_alert_ref, containment_measures, created_at
		FROM biosurv_disease_alerts WHERE $1 = ANY(affected_regions) ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, region)
	if err != nil {
		return nil, fmt.Errorf("query alerts by region: %w", err)
	}
	defer rows.Close()

	var alerts []domain.DiseaseAlert
	for rows.Next() {
		var a domain.DiseaseAlert
		if err := rows.Scan(
			&a.ID, &a.DiseaseName, &a.PathogenType, &a.Icd10Code, &a.AlertLevel, &a.FirstCaseDetected,
			&a.SymptomsHallmark, &a.TransmissionMode, &a.IncubationDays, &a.FatalityRate,
			&a.CasesConfirmed, &a.CasesSuspected, &a.CasesDeaths, pq.Array(&a.AffectedRegions),
			&a.SourceLab, &a.WhoAlertRef, &a.ContainmentMeasures, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (r *postgresRepo) CreateCampaign(ctx context.Context, c *domain.VaccinationCampaign) error {
	query := `INSERT INTO biosurv_vaccination_campaigns (id, campaign_name, target_disease, vaccine_type, target_population, doses_administered, coverage_pct, regions_active, start_date, end_date, coordinator_agency, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.CampaignName, c.TargetDisease, c.VaccineType, c.TargetPopulation,
		c.DosesAdministered, c.CoveragePct, pq.Array(c.RegionsActive),
		c.StartDate, c.EndDate, c.CoordinatorAgency, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert campaign: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetCampaignCoverage(ctx context.Context, id uuid.UUID) (*domain.VaccinationCampaign, error) {
	query := `SELECT id, campaign_name, target_disease, vaccine_type, target_population, doses_administered, coverage_pct, regions_active, start_date, end_date, coordinator_agency, created_at
		FROM biosurv_vaccination_campaigns WHERE id = $1`

	c := &domain.VaccinationCampaign{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.CampaignName, &c.TargetDisease, &c.VaccineType, &c.TargetPopulation,
		&c.DosesAdministered, &c.CoveragePct, pq.Array(&c.RegionsActive),
		&c.StartDate, &c.EndDate, &c.CoordinatorAgency, &c.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("query campaign coverage: %w", err)
	}
	return c, nil
}

func (r *postgresRepo) UpdateFacilityStock(ctx context.Context, id uuid.UUID, req domain.UpdateFacilityStockRequest) (*domain.HealthFacility, error) {
	query := `UPDATE biosurv_facilities SET stock_status = $1, beds_available = $2, last_report_at = $3 WHERE id = $4
		RETURNING id, facility_name, facility_type, region, commune, dept_code, capacity_beds, beds_available, stock_status, has_ventilators, has_ambulance, last_report_at, created_at`

	f := &domain.HealthFacility{}
	err := r.db.QueryRowContext(ctx, query, req.StockStatus, req.BedsAvailable, time.Now().UTC(), id).Scan(
		&f.ID, &f.FacilityName, &f.FacilityType, &f.Region, &f.Commune, &f.DeptCode,
		&f.CapacityBeds, &f.BedsAvailable, &f.StockStatus, &f.HasVentilators, &f.HasAmbulance,
		&f.LastReportAt, &f.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update facility stock: %w", err)
	}
	return f, nil
}

func (r *postgresRepo) GetDashboardNational(ctx context.Context) (*domain.DashboardNational, error) {
	query := `SELECT
		COALESCE((SELECT COUNT(*) FROM biosurv_disease_alerts), 0),
		COALESCE((SELECT COUNT(*) FROM biosurv_disease_alerts WHERE alert_level IN ('YELLOW','ORANGE','RED')), 0),
		COALESCE((SELECT COUNT(*) FROM biosurv_vaccination_campaigns), 0),
		COALESCE((SELECT COUNT(*) FROM biosurv_facilities), 0),
		COALESCE((SELECT COUNT(*) FROM biosurv_facilities WHERE stock_status = 'CRITICAL' OR stock_status = 'OUT_OF_STOCK'), 0),
		COALESCE((SELECT AVG(coverage_pct) FROM biosurv_vaccination_campaigns), 0.0),
		COALESCE((SELECT SUM(cases_confirmed) FROM biosurv_disease_alerts), 0),
		COALESCE((SELECT SUM(cases_deaths) FROM biosurv_disease_alerts), 0)`

	d := &domain.DashboardNational{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&d.TotalAlerts, &d.ActiveAlerts, &d.TotalCampaigns, &d.TotalFacilities,
		&d.CriticalFacilities, &d.AvgCoveragePct, &d.TotalCasesConfirmed, &d.TotalDeaths,
	)
	if err != nil {
		return nil, fmt.Errorf("query dashboard national: %w", err)
	}
	return d, nil
}
